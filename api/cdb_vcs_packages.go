package api

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	cdbVcsPackagesURL = "https://content.luanti.org/api/packages/?fmt=vcs"
	cdbVcsPackagesTTL = 24 * time.Hour
)

var cdbVcsPackages = &cdbCache{url: cdbVcsPackagesURL, ttl: cdbVcsPackagesTTL}

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

	data, err := c.fetch()
	if err != nil {
		if c.data != nil {
			// serve stale data rather than failing
			logrus.WithError(err).Warn("cdb vcs-packages refresh failed, serving stale data")
			return c.data, nil
		}
		return nil, err
	}

	c.data = data
	c.fetched = time.Now()
	return data, nil
}

func (c *cdbCache) fetch() ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from cdb: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (api *Api) GetCDBVcsPackages(w http.ResponseWriter, r *http.Request) {
	data, err := cdbVcsPackages.Get()
	if err != nil {
		SendError(w, http.StatusBadGateway, err.Error())
		return
	}
	SendRawJson(w, data)
}
