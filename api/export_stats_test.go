package api

import (
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/stretchr/testify/assert"
)

func TestClientIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.0.0.5:1234"
	assert.Equal(t, "10.0.0.5", clientIP(r))

	r.Header.Set("X-Forwarded-For", " 1.2.3.4 , 5.6.7.8")
	assert.Equal(t, "1.2.3.4", clientIP(r))
}

func TestMarkFirstHit(t *testing.T) {
	cache := expirable.NewLRU[string, bool](10, nil, 0)

	a := httptest.NewRequest("GET", "/", nil)
	a.RemoteAddr = "10.0.0.1:1"
	b := httptest.NewRequest("GET", "/", nil)
	b.RemoteAddr = "10.0.0.2:1"

	assert.True(t, markFirstHit(cache, "schema1", a))
	assert.False(t, markFirstHit(cache, "schema1", a))
	// other client and other schema are counted separately
	assert.True(t, markFirstHit(cache, "schema1", b))
	assert.True(t, markFirstHit(cache, "schema2", a))
}
