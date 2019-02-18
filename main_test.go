package main

import (
	"github.com/gorilla/handlers"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSelfRoot(t *testing.T) {
	var (
		r   *http.Request
		err error
	)

	// Test Default
	r, err = http.NewRequest("GET", "/authorize", nil)
	r.Host = "foo.bar"
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "http://foo.bar", SelfRoot(r))

	// Test Manual forwarded proto
	r, err = http.NewRequest("GET", "/validate", nil)
	r.Host = "foo.bar"
	r.Header.Set("X-Forwarded-Proto", "https")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "https://foo.bar", SelfRoot(r))

	// Test ProxyHeader handler
	rr := httptest.NewRecorder()
	r, err = http.NewRequest("GET", "/validate", nil)
	r.Header.Set("X-Forwarded-Host", "foo.bar")
	r.Header.Set("X-Forwarded-Proto", "https")
	handlers.ProxyHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rr, r)
	assert.Equal(t, "https://foo.bar", SelfRoot(r))
}

func TestAllowedHostsHandler_single_hostname(t *testing.T) {
	f := allowedHostsHandler("foo.bar")

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Host = "foo.bar"

	f(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rr, r)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func TestAllowedHostsHandler_multiple_hostnames(t *testing.T) {
	f := allowedHostsHandler("foo.bar, bar.foo")

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Host = "bar.foo"

	f(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rr, r)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func TestAllowedHostsHandler_mismatch_hostname(t *testing.T) {
	f := allowedHostsHandler("unknown.host")

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Host = "known.host"

	f(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rr, r)
	assert.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
}
