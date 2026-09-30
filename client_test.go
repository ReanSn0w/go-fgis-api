package fgis

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testJWT(exp time.Time) string {
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, exp.Unix())))
	return "header." + payload + ".signature"
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func mockHTTPClient(handler http.Handler) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		return recorder.Result(), nil
	})}
}

func TestSearchDeclarationsCarriesCookieAndBearer(t *testing.T) {
	token := "Bearer " + testJWT(time.Now().Add(time.Hour))
	var logins atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			if r.Method != http.MethodGet {
				t.Errorf("bootstrap method = %s", r.Method)
			}
			if !strings.Contains(r.UserAgent(), "Firefox/") || r.Header.Get("Accept-Language") == "" {
				t.Error("bootstrap request lacks observed browser headers")
			}
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test-session", Path: "/"})
		case "/login":
			logins.Add(1)
			if r.Method != http.MethodPost {
				t.Errorf("login method = %s", r.Method)
			}
			if cookie, err := r.Cookie("session-cookie"); err != nil || cookie.Value != "test-session" {
				t.Error("login did not carry bootstrap cookie")
			}
			if r.Header.Get("Authorization") != "Bearer null" {
				t.Error("missing initial Bearer header")
			}
			var credentials Credentials
			if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
				t.Error(err)
			}
			if credentials.Username != "user" || credentials.Password != "password" {
				t.Error("wrong login body")
			}
			w.Header().Set("Authorization", token)
		case declarations:
			if cookie, err := r.Cookie("session-cookie"); err != nil || cookie.Value != "test-session" {
				t.Error("list request did not carry session cookie")
			}
			if r.Header.Get("Authorization") != token {
				t.Error("list request did not carry login token")
			}
			var body struct {
				Size   int `json:"size"`
				Filter struct {
					Status []int `json:"status"`
				} `json:"filter"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body.Size != 10 || body.Filter.Status == nil {
				t.Error("default search body is incomplete")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[{"id":42,"declNumber":"TEST-42"}],"size":10,"total":1}`))
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	page, err := client.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest())
	if err != nil {
		t.Fatal(err)
	}
	if logins.Load() != 1 {
		t.Fatalf("login calls = %d, want 1", logins.Load())
	}
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].DeclarationNumber != "TEST-42" {
		t.Fatalf("unexpected declaration page: %+v", page)
	}
}

func TestSearchRetriesOnceAfterUnauthorized(t *testing.T) {
	var logins, searches atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			n := logins.Add(1)
			w.Header().Set("Authorization", fmt.Sprintf("Bearer token-%d", n))
		case declarations:
			searches.Add(1)
			switch r.Header.Get("Authorization") {
			case "Bearer token-1":
				w.WriteHeader(http.StatusUnauthorized)
			case "Bearer token-2":
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"items":[],"size":10,"total":0}`))
			default:
				w.WriteHeader(http.StatusUnauthorized)
			}
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	if _, err := client.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest()); err != nil {
		t.Fatal(err)
	}
	if logins.Load() != 2 || searches.Load() != 2 {
		t.Fatalf("login/list calls = %d/%d, want 2/2", logins.Load(), searches.Load())
	}
}

func TestStartRefreshesBeforeExpiry(t *testing.T) {
	refreshed := make(chan struct{}, 1)
	var logins atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			n := logins.Add(1)
			exp := time.Now().Add(2 * time.Second)
			if n > 1 {
				exp = time.Now().Add(time.Hour)
				select {
				case refreshed <- struct{}{}:
				default:
				}
			}
			w.Header().Set("Authorization", "Bearer "+testJWT(exp))
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{
		BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler),
		Username: "user", Password: "password",
		RefreshBefore: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer client.Close()
	if err := client.Start(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case <-refreshed:
		if logins.Load() != 2 {
			t.Fatalf("login calls = %d, want 2", logins.Load())
		}
	case <-time.After(4 * time.Second):
		t.Fatal("background token refresh did not run")
	}
}

func TestCloseRejectsNewRequests(t *testing.T) {
	client, err := New(Config{Username: "user", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = client.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest())
	if err != ErrClosed {
		t.Fatalf("error = %v, want ErrClosed", err)
	}
}

func TestJWTExpiry(t *testing.T) {
	exp := time.Now().Add(time.Hour).Truncate(time.Second)
	got, ok := jwtExpiry(testJWT(exp))
	if !ok || !got.Equal(exp) {
		t.Fatalf("expiry = %v, valid = %v", got, ok)
	}
	if _, ok := jwtExpiry(strings.Repeat("x", 10)); ok {
		t.Fatal("accepted invalid JWT")
	}
}
