package fgis

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Authenticate forces a fresh page visit and login. Failed proactive refresh
// leaves a still-valid token available to other requests.
func (c *Client) Authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed.Load() {
		return ErrClosed
	}
	return c.authenticateLocked(ctx)
}

func (c *Client) ensureToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed.Load() {
		return "", ErrClosed
	}
	if err := c.ensureTokenLocked(ctx); err != nil {
		return "", err
	}
	return c.token, nil
}

func (c *Client) ensureTokenLocked(ctx context.Context) error {
	now := time.Now()
	if c.token != "" && now.Before(c.refreshAt) && (c.expiresAt.IsZero() || now.Before(c.expiresAt)) {
		return nil
	}
	if err := c.authenticateLocked(ctx); err != nil {
		if c.token == "" || (!c.expiresAt.IsZero() && !time.Now().Before(c.expiresAt)) {
			return err
		}
		c.refreshAt = time.Now().Add(30 * time.Second)
	}
	return nil
}

func (c *Client) refreshAfterUnauthorized(ctx context.Context, failedToken string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed.Load() {
		return "", ErrClosed
	}
	if c.token != failedToken && c.token != "" {
		return c.token, nil
	}
	c.token, c.expiresAt, c.refreshAt = "", time.Time{}, time.Time{}
	if err := c.authenticateLocked(ctx); err != nil {
		return "", err
	}
	return c.token, nil
}

func (c *Client) authenticateLocked(ctx context.Context) error {
	credentials := c.credentials
	if c.credentialsProvider != nil {
		var err error
		credentials, err = c.credentialsProvider(ctx)
		if err != nil {
			return fmt.Errorf("load FGIS credentials: %w", err)
		}
	}
	if credentials.Username == "" || credentials.Password == "" {
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

	loginBody, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("encode login request: %w", err)
	}
	loginReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(loginPath), bytes.NewReader(loginBody))
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
	issuedAt := time.Now()
	expiresAt, hasExpiry := jwtExpiry(token)
	if hasExpiry && !expiresAt.After(issuedAt) {
		return errors.New("login returned an expired Bearer token")
	}
	refreshAt := issuedAt.Add(time.Hour)
	if hasExpiry {
		refreshAt = expiresAt.Add(-c.refreshBefore)
		if !refreshAt.After(issuedAt) {
			refreshAt = issuedAt.Add(expiresAt.Sub(issuedAt) / 2)
		}
	}
	c.token, c.expiresAt, c.refreshAt = authorization, expiresAt, refreshAt
	select {
	case c.wake <- struct{}{}:
	default:
	}
	return nil
}

// The unverified exp claim is used only for scheduling. The token itself is
// supplied by the HTTPS login response and is never exposed through this API.
func jwtExpiry(token string) (time.Time, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, false
	}
	var claims struct {
		Exp json.Number `json:"exp"`
	}
	if json.Unmarshal(payload, &claims) != nil || claims.Exp == "" {
		return time.Time{}, false
	}
	seconds, err := claims.Exp.Int64()
	if err != nil {
		return time.Time{}, false
	}
	return time.Unix(seconds, 0), true
}

func (c *Client) refreshLoop(ctx context.Context, done chan struct{}) {
	defer close(done)
	for {
		c.mu.Lock()
		refreshAt := c.refreshAt
		c.mu.Unlock()
		delay := time.Until(refreshAt)
		if delay < time.Second {
			delay = time.Second
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-c.wake:
			timer.Stop()
			continue
		case <-timer.C:
		}
		c.mu.Lock()
		if c.closed.Load() || ctx.Err() != nil {
			c.mu.Unlock()
			return
		}
		if time.Now().Before(c.refreshAt) {
			c.mu.Unlock()
			continue
		}
		err := c.authenticateLocked(ctx)
		if err != nil {
			c.refreshAt = time.Now().Add(30 * time.Second)
		}
		c.mu.Unlock()
		if err != nil && ctx.Err() == nil && c.onRefreshError != nil {
			c.onRefreshError(err)
		}
	}
}

func setBrowserHeaders(req *http.Request, base *url.URL, referer string) {
	setCommonHeaders(req)
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
	setCommonHeaders(req)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
}

func setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:156.0) Gecko/20100101 Firefox/156.0")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
}
