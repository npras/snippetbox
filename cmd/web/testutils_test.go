package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/npras/snippetbox/internal/models/mocks"
)

func extractCSRFToken(t *testing.T, body string) string {
	csrfTokenRX := regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(matches[1])
}

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		snippet:        &mocks.SnippetModel{},
		user:           &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

//

type testServer struct {
	*httptest.Server
	t *testing.T
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts, t}
}

func (ts *testServer) get(path string) (int, http.Header, string) {
	resp, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		ts.t.Fatal(err)
	}
	defer resp.Body.Close()
	return ts.readBodyAndExtractStuff(resp)
}

func (ts *testServer) postForm(path string, f url.Values) (int, http.Header, string) {
	resp, err := ts.Client().PostForm(ts.URL+path, f)
	if err != nil {
		ts.t.Fatal(err)
	}
	defer resp.Body.Close()
	return ts.readBodyAndExtractStuff(resp)
}

func (ts *testServer) readBodyAndExtractStuff(r *http.Response) (int, http.Header, string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ts.t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return r.StatusCode, r.Header, string(body)
}
