package fgis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newTestClient(t *testing.T, cfg Config) *Client {
	t.Helper()
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://fgis.test"
	}
	c, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func waitSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(4 * time.Second):
		t.Fatal("timed out waiting for client lifecycle")
	}
}

func currentLoop(c *Client) *backgroundLoop {
	c.lifecycleMu.Lock()
	defer c.lifecycleMu.Unlock()
	return c.loop
}

func triggerRefresh(c *Client) {
	c.mu.Lock()
	c.refreshAt = time.Now().Add(-time.Second)
	c.mu.Unlock()
	select {
	case c.wake <- struct{}{}:
	default:
	}
}

func TestNewRefreshSettingsAndHTTPClientCopy(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	hc := &http.Client{Jar: jar}
	c := newTestClient(t, Config{HTTPClient: hc})
	if c.refreshBefore != 5*time.Minute {
		t.Fatalf("default refresh before = %v", c.refreshBefore)
	}
	if c.httpClient == hc || hc.Timeout != 0 || c.httpClient.Timeout != 30*time.Second || c.httpClient.Jar != jar {
		t.Fatal("New must copy the HTTP client, preserve its jar and default only the copy's timeout")
	}
	c = newTestClient(t, Config{RefreshBefore: time.Minute})
	if c.refreshBefore != time.Minute {
		t.Fatalf("configured refresh before = %v", c.refreshBefore)
	}
	if _, err := New(Config{RefreshBefore: -time.Second}); err == nil {
		t.Fatal("accepted a negative refresh interval")
	}
}

func TestCredentialsProviderPrecedenceAndReload(t *testing.T) {
	var calls, logins atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
		case "/login":
			n := logins.Add(1)
			var got Credentials
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Error(err)
			}
			want := Credentials{Username: "provider", Password: fmt.Sprintf("password-%d", n)}
			if got != want {
				t.Errorf("credentials = %+v, want provider credentials for login %d", got, n)
			}
			w.Header().Set("Authorization", "Bearer token")
		default:
			http.NotFound(w, r)
		}
	})
	c := newTestClient(t, Config{
		Username: "static", Password: "static", HTTPClient: mockHTTPClient(handler),
		CredentialsProvider: func(context.Context) (Credentials, error) {
			n := calls.Add(1)
			return Credentials{Username: "provider", Password: fmt.Sprintf("password-%d", n)}, nil
		},
	})
	for i := 0; i < 2; i++ {
		if err := c.Authenticate(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 2 || logins.Load() != 2 {
		t.Fatalf("provider/login calls = %d/%d", calls.Load(), logins.Load())
	}
}

func TestCredentialsProviderFailureDoesNotUseStaticCredentials(t *testing.T) {
	want := errors.New("credentials unavailable")
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"provider error", want},
		{"empty credentials", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests atomic.Int32
			c := newTestClient(t, Config{
				Username: "static", Password: "static",
				HTTPClient:          mockHTTPClient(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests.Add(1) })),
				CredentialsProvider: func(context.Context) (Credentials, error) { return Credentials{}, tc.err },
			})
			err := c.Authenticate(context.Background())
			if err == nil || (tc.err != nil && !errors.Is(err, want)) {
				t.Fatalf("authentication error = %v", err)
			}
			if requests.Load() != 0 {
				t.Fatal("provider failure must not fall back to static credentials")
			}
		})
	}
}

func TestConcurrentStartAndRestartAfterCancellation(t *testing.T) {
	var logins atomic.Int32
	c := newTestClient(t, Config{
		Username: "user", Password: "password",
		HTTPClient: mockHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				logins.Add(1)
				w.Header().Set("Authorization", "Bearer "+testJWT(time.Now().Add(time.Hour)))
			} else if r.URL.Path == declarations {
				_, _ = w.Write([]byte(`{"items":[],"size":10,"total":0}`))
			}
		})),
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.Start(ctx); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	first := currentLoop(c)
	if first == nil || logins.Load() != 1 {
		t.Fatalf("loop = %v, login calls = %d", first, logins.Load())
	}
	if err := c.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if currentLoop(c) != first {
		t.Fatal("repeated Start replaced the running loop")
	}
	cancel()
	// Restart must wait for the canceled loop rather than run alongside it.
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	waitSignal(t, first.done)
	second := currentLoop(c)
	if second == nil || second == first {
		t.Fatal("Start did not create a new loop after cancellation")
	}
	if _, err := c.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest()); err != nil {
		t.Fatalf("search after loop cancellation: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	waitSignal(t, second.done)
}

func TestStartFailureAllowsRetry(t *testing.T) {
	var logins atomic.Int32
	c := newTestClient(t, Config{
		Username: "user", Password: "password",
		HTTPClient: mockHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				if logins.Add(1) == 1 {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				w.Header().Set("Authorization", "Bearer token")
			}
		})),
	})
	var upstream *HTTPError
	if err := c.Start(context.Background()); !errors.As(err, &upstream) || upstream.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("first Start error = %v", err)
	}
	failed := currentLoop(c)
	waitSignal(t, failed.done)
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if currentLoop(c) == failed || logins.Load() != 2 {
		t.Fatal("failed Start prevented a new login and loop")
	}
}

func TestCloseBeforeStartIsPermanentAndIdempotent(t *testing.T) {
	c := newTestClient(t, Config{})
	for i := 0; i < 2; i++ {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
	for name, operation := range map[string]func() error{
		"Start":        func() error { return c.Start(context.Background()) },
		"Authenticate": func() error { return c.Authenticate(context.Background()) },
		"Declarations": func() error {
			_, err := c.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest())
			return err
		},
		"Certificates": func() error {
			_, err := c.SearchCertificates(context.Background(), CertificateSearchRequest{Filter: map[string]any{"invalid": make(chan int)}})
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			if err := operation(); err != ErrClosed {
				t.Fatalf("error = %v, want ErrClosed", err)
			}
		})
	}
}

func TestCloseCancelsInitialLogin(t *testing.T) {
	entered := make(chan struct{})
	c := newTestClient(t, Config{CredentialsProvider: func(ctx context.Context) (Credentials, error) {
		close(entered)
		<-ctx.Done()
		return Credentials{}, ctx.Err()
	}})
	result := make(chan error, 1)
	go func() { result <- c.Start(context.Background()) }()
	waitSignal(t, entered)
	closed := make(chan struct{})
	go func() { _ = c.Close(); close(closed) }()
	waitSignal(t, closed)
	if err := <-result; err != ErrClosed {
		t.Fatalf("Start error = %v, want ErrClosed", err)
	}
}

func TestCloseCancelsBackgroundRequest(t *testing.T) {
	entered, canceled := make(chan struct{}), make(chan struct{})
	var logins atomic.Int32
	c := newTestClient(t, Config{
		Username: "user", Password: "password",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path == "/login" && logins.Add(1) > 1 {
				close(entered)
				<-r.Context().Done()
				close(canceled)
				return nil, r.Context().Err()
			}
			return mockHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/login" {
					w.Header().Set("Authorization", "Bearer token")
				}
			})).Transport.RoundTrip(r)
		})},
	})
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	triggerRefresh(c)
	waitSignal(t, entered)
	closed := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); _ = c.Close() }()
		}
		wg.Wait()
		close(closed)
	}()
	waitSignal(t, closed)
	waitSignal(t, canceled)
	waitSignal(t, currentLoop(c).done)
}

func TestFailedRefreshKeepsOnlyValidToken(t *testing.T) {
	for _, expired := range []bool{false, true} {
		t.Run(fmt.Sprintf("expired=%v", expired), func(t *testing.T) {
			var logins, searches atomic.Int32
			c := newTestClient(t, Config{
				Username: "user", Password: "password",
				HTTPClient: mockHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/login":
						if logins.Add(1) > 1 {
							w.WriteHeader(http.StatusServiceUnavailable)
							return
						}
						w.Header().Set("Authorization", "Bearer "+testJWT(time.Now().Add(time.Hour)))
					case declarations:
						searches.Add(1)
						_, _ = w.Write([]byte(`{"items":[],"total":0,"size":10}`))
					}
				})),
			})
			if err := c.Authenticate(context.Background()); err != nil {
				t.Fatal(err)
			}
			c.mu.Lock()
			token := c.token
			c.refreshAt = time.Now().Add(-time.Second)
			if expired {
				c.expiresAt = time.Now().Add(-time.Second)
			}
			c.mu.Unlock()
			before := time.Now()
			_, err := c.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest())
			if expired {
				var upstream *HTTPError
				if !errors.As(err, &upstream) || upstream.StatusCode != http.StatusServiceUnavailable || searches.Load() != 0 {
					t.Fatalf("expired token: error=%v, searches=%d", err, searches.Load())
				}
				return
			}
			if err != nil || searches.Load() != 1 {
				t.Fatalf("valid token: error=%v, searches=%d", err, searches.Load())
			}
			c.mu.Lock()
			if c.token != token || c.refreshAt.Before(before.Add(30*time.Second)) || c.refreshAt.After(time.Now().Add(30*time.Second)) {
				t.Error("failed refresh did not preserve the token and schedule a retry after 30 seconds")
			}
			c.mu.Unlock()
			if _, err := c.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest()); err != nil {
				t.Fatal(err)
			}
			if logins.Load() != 2 {
				t.Fatal("attempted another login before the retry time")
			}
		})
	}
}

func TestBackgroundRefreshErrorCallback(t *testing.T) {
	reported := make(chan error, 1)
	var logins atomic.Int32
	c := newTestClient(t, Config{
		Username: "user", Password: "password", OnRefreshError: func(err error) { reported <- err },
		HTTPClient: mockHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				if logins.Add(1) > 1 {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				w.Header().Set("Authorization", "Bearer token")
			}
		})),
	})
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	before := time.Now()
	triggerRefresh(c)
	select {
	case err := <-reported:
		var upstream *HTTPError
		if !errors.As(err, &upstream) || upstream.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("callback error = %v", err)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("background error was not reported")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "Bearer token" || c.refreshAt.Before(before.Add(30*time.Second)) || c.refreshAt.After(time.Now().Add(30*time.Second)) {
		t.Fatal("background failure did not preserve the token and schedule retry")
	}
}

func TestSecondUnauthorizedReturnsHTTPError(t *testing.T) {
	var logins, searches atomic.Int32
	c := newTestClient(t, Config{
		Username: "user", Password: "password",
		HTTPClient: mockHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/login":
				w.Header().Set("Authorization", fmt.Sprintf("Bearer token-%d", logins.Add(1)))
			case declarations:
				searches.Add(1)
				w.WriteHeader(http.StatusUnauthorized)
			}
		})),
	})
	_, err := c.SearchDeclarations(context.Background(), DefaultDeclarationSearchRequest())
	var upstream *HTTPError
	if !errors.As(err, &upstream) || upstream.StatusCode != http.StatusUnauthorized || upstream.Path != declarations {
		t.Fatalf("error = %v, want declaration HTTP 401", err)
	}
	if logins.Load() != 2 || searches.Load() != 2 {
		t.Fatalf("login/search calls = %d/%d, want 2/2", logins.Load(), searches.Load())
	}
}

func TestSearchCertificatesCarriesFilterAndDecodesRawFields(t *testing.T) {
	var searches atomic.Int32
	c := newTestClient(t, Config{
		Username: "user", Password: "password",
		HTTPClient: mockHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/rds/declaration":
			case "/login":
				w.Header().Set("Authorization", "Bearer certificate-token")
			case certificates:
				searches.Add(1)
				if r.Method != http.MethodPost || r.Referer() != "https://fgis.test/rss/certificate" || r.Header.Get("Authorization") != "Bearer certificate-token" {
					t.Error("wrong certificate method, Referer or Bearer")
				}
				var query CertificateSearchRequest
				if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
					t.Error(err)
				}
				if query.Size != 3 || query.Page != 2 || query.Filter["fixture"] != "example" {
					t.Errorf("query = %+v", query)
				}
				// Synthetic fixture, not a claim about the upstream row schema.
				_, _ = w.Write([]byte(`{"items":[{"arbitrary":"value","nested":{"n":7}}],"size":3,"total":8}`))
			default:
				http.NotFound(w, r)
			}
		})),
	})
	page, err := c.SearchCertificates(context.Background(), CertificateSearchRequest{Size: 3, Page: 2, Filter: map[string]any{"fixture": "example"}})
	if err != nil {
		t.Fatal(err)
	}
	if searches.Load() != 1 || page.Size != 3 || page.Total != 8 || len(page.Items) != 1 || string(page.Items[0]["arbitrary"]) != `"value"` || string(page.Items[0]["nested"]) != `{"n":7}` {
		t.Fatalf("unexpected certificate page: %+v", page)
	}
}
