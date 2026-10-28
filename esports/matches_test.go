package esports

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func testService(t *testing.T, status int, body string) *Service {
	t.Helper()
	return &Service{token: "test-secret", baseURL: "https://api.pandascore.co", client: &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer test-secret" {
			t.Error("missing authentication")
		}
		if r.URL.Query().Get("per_page") != "5" {
			t.Error("missing pagination limit")
		}
		if !strings.HasPrefix(r.URL.Path, "/lol/matches/") {
			t.Error("wrong endpoint")
		}
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}}
}

func TestMatchesAndCache(t *testing.T) {
	s := testService(t, 200, `[{"id":42,"status":"finished","begin_at":"2026-10-27T12:00:00Z","league":{"name":"Worlds"},"opponents":[{"opponent":{"id":1,"name":"Team A","acronym":"A"}},{"opponent":{"id":2,"name":"Team B","acronym":"B"}}],"results":[{"team_id":1,"score":3},{"team_id":2,"score":1}]}]`)
	res := httptest.NewRecorder()
	s.Handler(res, httptest.NewRequest("GET", "/api/esports/matches", nil))
	if res.Code != 200 {
		t.Fatalf("status %d: %s", res.Code, res.Body.String())
	}
	var feed Feed
	if err := json.Unmarshal(res.Body.Bytes(), &feed); err != nil {
		t.Fatal(err)
	}
	if len(feed.Past) != 1 || len(feed.Running) != 1 || len(feed.Upcoming) != 1 || feed.Past[0].Results[0].Score != 3 || feed.UpdatedAt.IsZero() {
		t.Fatal("incorrect decoded feed")
	}
	if strings.Contains(res.Body.String(), s.token) {
		t.Fatal("token leaked")
	}
	s.client.Transport = transportFunc(func(*http.Request) (*http.Response, error) {
		t.Error("cache should prevent network requests")
		return nil, nil
	})
	cached := httptest.NewRecorder()
	s.Handler(cached, httptest.NewRequest("GET", "/api/esports/matches", nil))
	if cached.Body.String() != res.Body.String() {
		t.Fatal("cache changed")
	}
}

func TestUnavailableAndStale(t *testing.T) {
	for _, status := range []int{401, 403, 429, 500} {
		s := testService(t, status, `{"error":"private upstream details"}`)
		res := httptest.NewRecorder()
		s.Handler(res, httptest.NewRequest("GET", "/", nil))
		if res.Code != 503 || strings.Contains(res.Body.String(), "private") {
			t.Fatalf("unexpected failure response %d %s", res.Code, res.Body.String())
		}
		s.retryAt = time.Time{}
		s.cached = Feed{UpdatedAt: time.Now().Add(-2 * time.Minute), Past: []Match{{ID: 42}}}
		res = httptest.NewRecorder()
		s.Handler(res, httptest.NewRequest("GET", "/", nil))
		var feed Feed
		if err := json.Unmarshal(res.Body.Bytes(), &feed); err != nil {
			t.Fatal(err)
		}
		if res.Code != 200 || !feed.Stale || len(feed.Past) != 1 {
			t.Fatal("expected labeled cached results")
		}
		s.cached.UpdatedAt = time.Now().Add(-16 * time.Minute)
		res = httptest.NewRecorder()
		s.Handler(res, httptest.NewRequest("GET", "/", nil))
		if res.Code != 503 {
			t.Fatal("expired cache served")
		}
	}
}

func TestEmptyMalformedAndMissingToken(t *testing.T) {
	for _, tc := range []struct {
		body   string
		status int
	}{{"[]", 200}, {"null", 200}, {"not json", 503}} {
		s := testService(t, 200, tc.body)
		res := httptest.NewRecorder()
		s.Handler(res, httptest.NewRequest("GET", "/", nil))
		if res.Code != tc.status {
			t.Fatalf("%q: %d", tc.body, res.Code)
		}
	}
	s := &Service{}
	res := httptest.NewRecorder()
	s.Handler(res, httptest.NewRequest("GET", "/", nil))
	if res.Code != 503 {
		t.Fatal("missing token must not fetch")
	}
	res = httptest.NewRecorder()
	s.Handler(res, httptest.NewRequest("POST", "/", nil))
	if res.Code != 405 || res.Header().Get("Allow") != "GET" {
		t.Fatal("method not rejected")
	}
}
