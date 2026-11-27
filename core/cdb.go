package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	CDBBaseURL         = "https://content.luanti.org"
	cdbVcsPackagesTTL  = 24 * time.Hour
	cdbRequestTimeout  = 30 * time.Second
	cdbMaxResponseSize = 64 * 1024 * 1024
)

// list of all packages with their git repository ("author", "name", "title", "type", "repo"), cached for a day
var CDBVcsPackages = &cdbCache{url: CDBBaseURL + "/api/packages/?fmt=vcs", ttl: cdbVcsPackagesTTL}

var cdbClient = &http.Client{Timeout: cdbRequestTimeout}

// in-memory cache of a single remote json document
type cdbCache struct {
	url     string
	ttl     time.Duration
	mu      sync.Mutex
	data    []byte
	fetched time.Time
}

func (c *cdbCache) Get() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.data != nil && time.Since(c.fetched) < c.ttl {
		return c.data, nil
	}

	data, err := cdbGet(context.Background(), c.url)
	if err != nil {
		if c.data != nil {
			// serve stale data rather than failing
			logrus.WithError(err).Warn("cdb refresh failed, serving stale data")
			return c.data, nil
		}
		return nil, err
	}

	c.data = data
	c.fetched = time.Now()
	return data, nil
}

func cdbGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cdbClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from cdb: %d", resp.StatusCode)
	}

	return io.ReadAll(io.LimitReader(resp.Body, cdbMaxResponseSize))
}

type CDBVcsPackage struct {
	Author string `json:"author"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Repo   string `json:"repo"`
}

// comparable form of a git url: no scheme, credentials, ".git" suffix or trailing slash, lower case
func normalizeGitURL(u string) string {
	u = strings.ToLower(strings.TrimSpace(u))
	if parsed, err := url.Parse(u); err == nil && parsed.Host != "" {
		u = parsed.Host + parsed.Path
	}
	u = strings.TrimSuffix(u, "/")
	u = strings.TrimSuffix(u, ".git")
	return strings.TrimSuffix(u, "/")
}

// contentdb package ("author/name") with the given git repository, empty if none matches
func findCDBPackageByRepo(packages []CDBVcsPackage, git_url string) string {
	n := normalizeGitURL(git_url)
	for _, p := range packages {
		if p.Repo != "" && normalizeGitURL(p.Repo) == n {
			return p.Author + "/" + p.Name
		}
	}
	return ""
}

func lookupCDBPackageByRepo(git_url string) (string, error) {
	data, err := CDBVcsPackages.Get()
	if err != nil {
		return "", err
	}
	packages := []CDBVcsPackage{}
	err = json.Unmarshal(data, &packages)
	if err != nil {
		return "", fmt.Errorf("parse vcs packages: %w", err)
	}
	return findCDBPackageByRepo(packages, git_url), nil
}

type CDBPackageLicenses struct {
	License      string `json:"license"`
	MediaLicense string `json:"media_license"`
}

// license and media license from the package detail ("author/name")
func fetchCDBPackageLicenses(ctx context.Context, pkg string) (*CDBPackageLicenses, error) {
	author, name, found := strings.Cut(pkg, "/")
	if !found || author == "" || name == "" || strings.Contains(name, "/") {
		return nil, fmt.Errorf("invalid contentdb package '%s', expected 'author/name'", pkg)
	}
	data, err := cdbGet(ctx, fmt.Sprintf("%s/api/packages/%s/%s/", CDBBaseURL, url.PathEscape(author), url.PathEscape(name)))
	if err != nil {
		return nil, err
	}
	l := &CDBPackageLicenses{}
	err = json.Unmarshal(data, l)
	if err != nil {
		return nil, fmt.Errorf("parse package detail: %w", err)
	}
	return l, nil
}
