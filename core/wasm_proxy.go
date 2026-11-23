package core

import (
	"blockexchange/metrics"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// websocket proxy for the luanti wasm client (protocol of the emsocket proxy in luanti-wasm-playground).
//
// The first text message selects the target:
//   - "NEWADDR"                          -> "ADDR ..." (address assignment, unused)
//   - "PROXY IPV4 TCP 10.0.0.1 53"       -> dns lookup, then "<hostname> <...>" -> [count, version, address...]
//   - "PROXY IPV4 TCP 10.0.0.1 8080"     -> http proxy (CONNECT or absolute-form requests)
//
// Unlike the playground proxy this one is not an open relay: dns lookups and http(s) requests are only
// allowed to the host of the configured BASE_URL (the blockexchange mod talks to this server),
// udp and plain tcp connections to other hosts are refused.

var (
	// virtual addresses of the proxy services inside the wasm client
	wasmProxyHosts = []string{"10.0.0.1", "fd00::1"}
)

const (
	wasmProxyDNSPort  = "53"
	wasmProxyHTTPPort = "8080"

	wasmProxyDialTimeout = 10 * time.Second
	// max size of a single websocket message from the client
	wasmProxyReadLimit = 1024 * 1024
)

var ErrWasmProxyNotAllowed = errors.New("target not allowed")

type WasmProxy struct {
	// hostname and port of BASE_URL, empty if not configured
	host string
	port string
}

func (c *Core) NewWasmProxy() *WasmProxy {
	p := &WasmProxy{}
	u, err := url.Parse(c.cfg.BaseURL)
	if err != nil || u.Hostname() == "" {
		return p
	}
	p.host = strings.ToLower(u.Hostname())
	p.port = u.Port()
	if p.port == "" {
		p.port = "80"
		if u.Scheme == "https" {
			p.port = "443"
		}
	}
	return p
}

// handles a single proxy connection, blocks until it is closed
func (p *WasmProxy) Handle(ws *websocket.Conn) error {
	defer ws.Close()
	ws.SetReadLimit(wasmProxyReadLimit)

	metrics.WasmSessionsActive.Inc()
	defer metrics.WasmSessionsActive.Dec()

	_, data, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	parts := strings.Split(string(data), " ")
	if parts[0] == "NEWADDR" {
		metrics.WasmProxyConnections.WithLabelValues("newaddr").Inc()
		return ws.WriteMessage(websocket.TextMessage, []byte("ADDR xy abc ?"))
	}
	if parts[0] != "PROXY" || len(parts) != 5 {
		metrics.WasmProxyConnections.WithLabelValues("invalid").Inc()
		return fmt.Errorf("invalid command: '%s'", data)
	}
	protocol, host, port := parts[2], parts[3], parts[4]

	if protocol != "TCP" || !slices.Contains(wasmProxyHosts, host) {
		metrics.WasmProxyRejected.WithLabelValues("target").Inc()
		return fmt.Errorf("%w: %s %s:%s", ErrWasmProxyNotAllowed, protocol, host, port)
	}

	switch port {
	case wasmProxyDNSPort:
		metrics.WasmProxyConnections.WithLabelValues("dns").Inc()
		return p.handleDNS(ws)
	case wasmProxyHTTPPort:
		metrics.WasmProxyConnections.WithLabelValues("http").Inc()
		err = ws.WriteMessage(websocket.TextMessage, []byte("PROXY OK"))
		if err != nil {
			return err
		}
		return p.handleHTTP(newWSConn(ws))
	default:
		metrics.WasmProxyRejected.WithLabelValues("target").Inc()
		return fmt.Errorf("%w: %s %s:%s", ErrWasmProxyNotAllowed, protocol, host, port)
	}
}

// resolves the allowed host only
func (p *WasmProxy) handleDNS(ws *websocket.Conn) error {
	err := ws.WriteMessage(websocket.TextMessage, []byte("PROXY OK"))
	if err != nil {
		return err
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		return err
	}
	parts := strings.Split(string(data), " ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid dns request: '%x'", data)
	}
	if p.host == "" || !strings.EqualFold(parts[0], p.host) {
		metrics.WasmProxyRejected.WithLabelValues("dns").Inc()
		return fmt.Errorf("%w: dns lookup of '%s'", ErrWasmProxyNotAllowed, parts[0])
	}

	ips, err := net.LookupIP(p.host)
	if err != nil {
		return err
	}
	// prefer ipv4, record-count, ip-version, address
	for _, ip := range ips {
		if ip4 := ip.To4(); ip4 != nil {
			return ws.WriteMessage(websocket.BinaryMessage, append([]byte{0x01, 0x04}, ip4...))
		}
	}
	for _, ip := range ips {
		if ip16 := ip.To16(); ip16 != nil {
			return ws.WriteMessage(websocket.BinaryMessage, append([]byte{0x01, 0x06}, ip16...))
		}
	}
	return fmt.Errorf("host '%s' not found", p.host)
}

// checks the "host:port" target and dials the configured host
func (p *WasmProxy) dial(ctx context.Context, target string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		return nil, err
	}
	if p.host == "" || port != p.port || !p.isAllowedHost(host) {
		metrics.WasmProxyRejected.WithLabelValues("http").Inc()
		return nil, fmt.Errorf("%w: %s", ErrWasmProxyNotAllowed, target)
	}
	d := net.Dialer{Timeout: wasmProxyDialTimeout}
	return d.DialContext(ctx, "tcp", net.JoinHostPort(p.host, p.port))
}

// the configured hostname or one of its addresses (if the client resolved it before)
func (p *WasmProxy) isAllowedHost(host string) bool {
	host = strings.Trim(strings.ToLower(host), "[]")
	if host == p.host {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	ips, err := net.LookupIP(p.host)
	if err != nil {
		return false
	}
	return slices.ContainsFunc(ips, ip.Equal)
}

// minimal http proxy on the tunneled tcp stream
func (p *WasmProxy) handleHTTP(conn net.Conn) error {
	br := bufio.NewReader(conn)
	transport := &http.Transport{
		DialContext:       func(ctx context.Context, _, addr string) (net.Conn, error) { return p.dial(ctx, addr) },
		Proxy:             nil,
		DisableKeepAlives: false,
	}
	defer transport.CloseIdleConnections()

	for {
		req, err := http.ReadRequest(br)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if req.Method == http.MethodConnect {
			return p.handleConnect(conn, br, req)
		}

		if req.URL.Host == "" {
			return writeProxyError(conn, http.StatusBadRequest, "absolute url required")
		}
		req.RequestURI = ""
		req.Header.Del("Proxy-Connection")
		req.Header.Del("Proxy-Authorization")
		resp, err := transport.RoundTrip(req)
		if errors.Is(err, ErrWasmProxyNotAllowed) {
			return writeProxyError(conn, http.StatusForbidden, err.Error())
		}
		if err != nil {
			return writeProxyError(conn, http.StatusBadGateway, err.Error())
		}

		err = resp.Write(conn)
		resp.Body.Close()
		if err != nil {
			return err
		}
		if req.Close || resp.Close || (resp.ContentLength < 0 && len(resp.TransferEncoding) == 0) {
			// no way to delimit the next response
			return nil
		}
	}
}

func (p *WasmProxy) handleConnect(conn net.Conn, br *bufio.Reader, req *http.Request) error {
	target, err := p.dial(context.Background(), req.Host)
	if errors.Is(err, ErrWasmProxyNotAllowed) {
		return writeProxyError(conn, http.StatusForbidden, err.Error())
	}
	if err != nil {
		return writeProxyError(conn, http.StatusBadGateway, err.Error())
	}
	defer target.Close()

	_, err = io.WriteString(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
	if err != nil {
		return err
	}

	errs := make(chan error, 2)
	go func() {
		// includes the data already buffered by the request reader
		_, err := io.Copy(target, br)
		errs <- err
	}()
	go func() {
		_, err := io.Copy(conn, target)
		errs <- err
	}()
	err = <-errs
	// unblock the other copy
	target.Close()
	conn.Close()
	return err
}

func writeProxyError(conn net.Conn, status int, msg string) error {
	_, err := fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s",
		status, http.StatusText(status), len(msg), msg)
	if err != nil {
		return err
	}
	return fmt.Errorf("proxy error %d: %s", status, msg)
}

// net.Conn on top of binary websocket messages
type wsConn struct {
	ws   *websocket.Conn
	buf  []byte
	wmux sync.Mutex
}

func newWSConn(ws *websocket.Conn) *wsConn {
	return &wsConn{ws: ws}
}

func (c *wsConn) Read(p []byte) (int, error) {
	for len(c.buf) == 0 {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			var ce *websocket.CloseError
			if errors.As(err, &ce) {
				return 0, io.EOF
			}
			return 0, err
		}
		c.buf = data
	}
	n := copy(p, c.buf)
	c.buf = c.buf[n:]
	metrics.WasmProxyBytes.WithLabelValues("in").Add(float64(n))
	return n, nil
}

func (c *wsConn) Write(p []byte) (int, error) {
	c.wmux.Lock()
	defer c.wmux.Unlock()
	err := c.ws.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	metrics.WasmProxyBytes.WithLabelValues("out").Add(float64(len(p)))
	return len(p), nil
}

func (c *wsConn) Close() error                       { return c.ws.Close() }
func (c *wsConn) LocalAddr() net.Addr                { return c.ws.LocalAddr() }
func (c *wsConn) RemoteAddr() net.Addr               { return c.ws.RemoteAddr() }
func (c *wsConn) SetReadDeadline(t time.Time) error  { return c.ws.SetReadDeadline(t) }
func (c *wsConn) SetWriteDeadline(t time.Time) error { return c.ws.SetWriteDeadline(t) }
func (c *wsConn) SetDeadline(t time.Time) error {
	err := c.ws.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return c.ws.SetWriteDeadline(t)
}
