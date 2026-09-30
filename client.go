// Package fgis provides a small client for the public Rosaccreditation
// declaration and certificate registries.
package fgis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL = "https://pub.fsa.gov.ru"
	loginPath      = "/login"
	declarations   = "/api/v1/rds/common/declarations/get"
	certificates   = "/api/v1/rss/common/certificates/get"
)

// Config contains credentials used by the site's login endpoint. Keep them in
// a secret store or environment variables; do not commit them to source code.
type Config struct {
	BaseURL    string
	StartPath  string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// Client maintains the site's session cookie and Bearer token in memory.
// A single Client should be reused for requests to both registries.
type Client struct {
	base       *url.URL
	startPath  string
	username   string
	password   string
	httpClient *http.Client

	mu    sync.Mutex
	token string
}

// New creates a client. If HTTPClient has no cookie jar, New installs one on
// a shallow copy so the bootstrap cookie is carried through /login and later
// registry requests.
func New(cfg Config) (*Client, error) {
	baseText := cfg.BaseURL
	if baseText == "" {
		baseText = defaultBaseURL
	}
	base, err := url.Parse(baseText)
	if err != nil || base.Host == "" || (base.Scheme != "https" && base.Scheme != "http") {
		return nil, fmt.Errorf("invalid FGIS base URL")
	}
	base.Path = strings.TrimRight(base.Path, "/")

	startPath := cfg.StartPath
	if startPath == "" {
		startPath = "/rds/declaration"
	}
	if !strings.HasPrefix(startPath, "/") {
		return nil, fmt.Errorf("start path must begin with slash")
	}

	hc := &http.Client{Timeout: 30 * time.Second}
	if cfg.HTTPClient != nil {
		copy := *cfg.HTTPClient
		hc = &copy
		if hc.Timeout == 0 {
			hc.Timeout = 30 * time.Second
		}
	}
	if hc.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, fmt.Errorf("create cookie jar: %w", err)
		}
		hc.Jar = jar
	}

	return &Client{
		base:       base,
		startPath:  startPath,
		username:   cfg.Username,
		password:   cfg.Password,
		httpClient: hc,
	}, nil
}

// Authenticate creates a site session and obtains the Bearer token returned
// by POST /login in the Authorization response header.
func (c *Client) Authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.authenticateLocked(ctx)
}

func (c *Client) authenticateLocked(ctx context.Context) error {
	if c.username == "" || c.password == "" {
		return errors.New("FGIS username and password are required")
	}

	pageURL := c.endpoint(c.startPath)
	bootstrapReq, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return fmt.Errorf("create bootstrap request: %w", err)
	}
	setNavigationHeaders(bootstrapReq)
	bootstrapResp, err := c.httpClient.Do(bootstrapReq)
	if err != nil {
		return fmt.Errorf("load registry page: %w", err)
	}
	_, readErr := io.Copy(io.Discard, bootstrapResp.Body)
	closeErr := bootstrapResp.Body.Close()
	if readErr != nil {
		return fmt.Errorf("read registry page response: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close registry page response: %w", closeErr)
	}
	if bootstrapResp.StatusCode < 200 || bootstrapResp.StatusCode >= 300 {
		return &HTTPError{Method: http.MethodGet, Path: c.startPath, StatusCode: bootstrapResp.StatusCode}
	}

	loginBody, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{Username: c.username, Password: c.password})
	if err != nil {
		return fmt.Errorf("encode login request: %w", err)
	}
	loginURL := c.endpoint(loginPath)
	loginReq, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, bytes.NewReader(loginBody))
	if err != nil {
		return fmt.Errorf("create login request: %w", err)
	}
	setBrowserHeaders(loginReq, c.base, pageURL)
	loginReq.Header.Set("Authorization", "Bearer null")
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := c.httpClient.Do(loginReq)
	if err != nil {
		return fmt.Errorf("send login request: %w", err)
	}
	_, readErr = io.Copy(io.Discard, loginResp.Body)
	closeErr = loginResp.Body.Close()
	if readErr != nil {
		return fmt.Errorf("read login response: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close login response: %w", closeErr)
	}
	if loginResp.StatusCode < 200 || loginResp.StatusCode >= 300 {
		return &HTTPError{Method: http.MethodPost, Path: loginPath, StatusCode: loginResp.StatusCode}
	}

	authorization := loginResp.Header.Get("Authorization")
	scheme, token, ok := strings.Cut(authorization, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.EqualFold(token, "null") {
		return errors.New("login succeeded without a Bearer token in the Authorization response header")
	}
	c.token = authorization
	return nil
}

// SearchDeclarations submits a caller-provided JSON search body to the
// declaration registry and returns its page envelope.
func (c *Client) SearchDeclarations(ctx context.Context, body any) (ListResponse[json.RawMessage], error) {
	return postList[json.RawMessage](ctx, c, declarations, body)
}

// SearchCertificates submits a caller-provided JSON search body to the
// certificate registry and returns its page envelope.
func (c *Client) SearchCertificates(ctx context.Context, body any) (ListResponse[json.RawMessage], error) {
	return postList[json.RawMessage](ctx, c, certificates, body)
}

// ListResponse is the common list envelope returned by registry endpoints.
type ListResponse[T any] struct {
	Items []T `json:"items"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// HTTPError describes an upstream HTTP error without including response
// bodies, which can contain credentials or personal data.
type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("FGIS %s %s returned HTTP %d", e.Method, e.Path, e.StatusCode)
}

func postList[T any](ctx context.Context, c *Client, path string, body any) (ListResponse[T], error) {
	var result ListResponse[T]
	if err := c.ensureAuthenticated(ctx); err != nil {
		return result, err
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return result, fmt.Errorf("encode search request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(path), bytes.NewReader(payload))
	if err != nil {
		return result, fmt.Errorf("create search request: %w", err)
	}
	setBrowserHeaders(req, c.base, c.endpoint(c.startPath))
	req.Header.Set("Authorization", c.bearer())
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("send search request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return result, &HTTPError{Method: http.MethodPost, Path: path, StatusCode: resp.StatusCode}
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("decode FGIS list response: %w", err)
	}
	return result, nil
}

func (c *Client) ensureAuthenticated(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" {
		return nil
	}
	return c.authenticateLocked(ctx)
}

func (c *Client) bearer() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.token
}

func (c *Client) endpoint(path string) string {
	u := *c.base
	u.Path = strings.TrimRight(c.base.Path, "/") + path
	return u.String()
}

func setBrowserHeaders(req *http.Request, base *url.URL, referer string) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", base.Scheme+"://"+base.Host)
	req.Header.Set("Referer", referer)
	req.Header.Set("lkId", "")
	req.Header.Set("orgId", "")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
}

func setNavigationHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
}
