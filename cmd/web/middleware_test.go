package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/npras/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	respWriter := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(respWriter, req)

	resp := respWriter.Result()

	expectedValue := "default-src 'self';font-src fonts.gstatic.com; style-src 'self' fonts.googleapis.com"
	assert.Equal(t, resp.Header.Get("Content-Security-Policy"), expectedValue)
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, resp.Header.Get("Referrer-Policy"), expectedValue)
	expectedValue = "nosniff"
	assert.Equal(t, resp.Header.Get("X-Content-Type-Options"), expectedValue)
	expectedValue = "deny"
	assert.Equal(t, resp.Header.Get("X-Frame-Options"), expectedValue)
	expectedValue = "0"
	assert.Equal(t, resp.Header.Get("X-XSS-Protection"), expectedValue)
	expectedValue = "Go"
	assert.Equal(t, resp.Header.Get("Server"), expectedValue)

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
