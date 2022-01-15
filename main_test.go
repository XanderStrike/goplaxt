package main

import (
	"github.com/gorilla/handlers"
	"github.com/stretchr/testify/assert"

	"context"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xanderstrike/goplaxt/lib/store"
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

func TestAllowedHostsHandler_alwaysAllowHealthcheck(t *testing.T) {
	storage = &MockSuccessStore{}
	f := allowedHostsHandler("unknown.host")

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Host = "known.host"

	f(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rr, r)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

type MockSuccessStore struct{}

func (s MockSuccessStore) Ping(ctx context.Context) error { return nil }
func (s MockSuccessStore) WriteUser(user store.User)      {}
func (s MockSuccessStore) GetUser(id string) *store.User  { return nil }
func (s MockSuccessStore) DeleteUser(id string) bool      { return true }

type MockFailStore struct{}

func (s MockFailStore) Ping(ctx context.Context) error { return errors.New("OH NO") }
func (s MockFailStore) WriteUser(user store.User)      { panic(errors.New("OH NO")) }
func (s MockFailStore) GetUser(id string) *store.User  { panic(errors.New("OH NO")) }
func (s MockFailStore) DeleteUser(id string) bool      { return false }

func TestHealthcheck(t *testing.T) {
	var rr *httptest.ResponseRecorder

	r, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	storage = &MockSuccessStore{}
	rr = httptest.NewRecorder()
	http.Handler(healthcheckHandler()).ServeHTTP(rr, r)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "{\"status\":\"OK\"}\n", rr.Body.String())

	storage = &MockFailStore{}
	rr = httptest.NewRecorder()
	http.Handler(healthcheckHandler()).ServeHTTP(rr, r)
	assert.Equal(t, http.StatusServiceUnavailable, rr.Result().StatusCode)
	assert.Equal(t, "{\"status\":\"Service Unavailable\",\"errors\":{\"storage\":\"OH NO\"}}\n", rr.Body.String())
}
