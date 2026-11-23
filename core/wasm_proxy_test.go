package core_test

import (
	"blockexchange/core"
	"blockexchange/types"
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// starts the target server (allowed via BASE_URL) and the websocket proxy, returns the proxy ws url
func setupWasmProxy(t *testing.T) (string, string) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from " + r.URL.Path))
	}))
	t.Cleanup(target.Close)

	cfg := types.CreateConfig()
	cfg.BaseURL = target.URL
	c := core.New(cfg, nil)

	upgrader := websocket.Upgrader{}
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c.NewWasmProxy().Handle(conn)
	}))
	t.Cleanup(proxy.Close)

	return "ws" + strings.TrimPrefix(proxy.URL, "http"), strings.TrimPrefix(target.URL, "http://")
}

func dialProxy(t *testing.T, url, cmd string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	t.Cleanup(func() { ws.Close() })
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	assert.NoError(t, ws.WriteMessage(websocket.TextMessage, []byte(cmd)))
	return ws
}

// reads binary messages until the http response is complete
func readResponse(t *testing.T, ws *websocket.Conn) *http.Response {
	pr, pw := io.Pipe()
	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			pw.Write(data)
		}
	}()
	resp, err := http.ReadResponse(bufio.NewReader(pr), nil)
	assert.NoError(t, err)
	return resp
}

func expectOK(t *testing.T, ws *websocket.Conn) {
	_, data, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "PROXY OK", string(data))
}

func TestWasmProxyHTTP(t *testing.T) {
	url, target := setupWasmProxy(t)

	ws := dialProxy(t, url, "PROXY IPV4 TCP 10.0.0.1 8080")
	expectOK(t, ws)
	assert.NoError(t, ws.WriteMessage(websocket.BinaryMessage, []byte("GET http://"+target+"/api/info HTTP/1.1\r\nHost: "+target+"\r\nConnection: close\r\n\r\n")))
	resp := readResponse(t, ws)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "hello from /api/info", string(body))
}

func TestWasmProxyConnect(t *testing.T) {
	url, target := setupWasmProxy(t)

	ws := dialProxy(t, url, "PROXY IPV6 TCP fd00::1 8080")
	expectOK(t, ws)
	// tunnel request and the first request inside the tunnel in one message
	assert.NoError(t, ws.WriteMessage(websocket.BinaryMessage, []byte(
		"CONNECT "+target+" HTTP/1.1\r\nHost: "+target+"\r\n\r\n"+
			"GET /tunnel HTTP/1.1\r\nHost: "+target+"\r\nConnection: close\r\n\r\n")))

	pr, pw := io.Pipe()
	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			pw.Write(data)
		}
	}()
	br := bufio.NewReader(pr)
	resp, err := http.ReadResponse(br, &http.Request{Method: http.MethodConnect})
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = http.ReadResponse(br, nil)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "hello from /tunnel", string(body))
}

func TestWasmProxyRefused(t *testing.T) {
	url, target := setupWasmProxy(t)

	// other hosts via the http proxy
	ws := dialProxy(t, url, "PROXY IPV4 TCP 10.0.0.1 8080")
	expectOK(t, ws)
	assert.NoError(t, ws.WriteMessage(websocket.BinaryMessage, []byte("CONNECT example.com:443 HTTP/1.1\r\nHost: example.com:443\r\n\r\n")))
	assert.Equal(t, 403, readResponse(t, ws).StatusCode)

	ws = dialProxy(t, url, "PROXY IPV4 TCP 10.0.0.1 8080")
	expectOK(t, ws)
	assert.NoError(t, ws.WriteMessage(websocket.BinaryMessage, []byte("GET http://example.com/ HTTP/1.1\r\nHost: example.com\r\n\r\n")))
	assert.Equal(t, 403, readResponse(t, ws).StatusCode)

	// right host, wrong port
	host := strings.Split(target, ":")[0]
	ws = dialProxy(t, url, "PROXY IPV4 TCP 10.0.0.1 8080")
	expectOK(t, ws)
	assert.NoError(t, ws.WriteMessage(websocket.BinaryMessage, []byte("CONNECT "+host+":22 HTTP/1.1\r\nHost: "+host+":22\r\n\r\n")))
	assert.Equal(t, 403, readResponse(t, ws).StatusCode)

	// direct udp/tcp connections are closed without "PROXY OK"
	for _, cmd := range []string{"PROXY IPV4 UDP 1.2.3.4 30000", "PROXY IPV4 TCP " + host + " 80", "PROXY IPV4 TCP 10.0.0.1 22", "FOO"} {
		ws = dialProxy(t, url, cmd)
		_, _, err := ws.ReadMessage()
		assert.Error(t, err, cmd)
	}
}

func TestWasmProxyDNS(t *testing.T) {
	url, target := setupWasmProxy(t)
	host := strings.Split(target, ":")[0]

	ws := dialProxy(t, url, "PROXY IPV4 TCP 10.0.0.1 53")
	expectOK(t, ws)
	assert.NoError(t, ws.WriteMessage(websocket.TextMessage, []byte(host+" A")))
	_, data, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x04, 127, 0, 0, 1}, data)

	// other hosts are not resolved
	ws = dialProxy(t, url, "PROXY IPV4 TCP 10.0.0.1 53")
	expectOK(t, ws)
	assert.NoError(t, ws.WriteMessage(websocket.TextMessage, []byte("example.com A")))
	_, _, err = ws.ReadMessage()
	assert.Error(t, err)
}
