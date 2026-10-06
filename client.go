// Package fgis provides a small client for the public Rosaccreditation
// declaration and certificate registries.
package fgis

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBaseURL = "https://pub.fsa.gov.ru"
	loginPath      = "/login"
	declarations   = "/api/v1/rds/common/declarations/get"
	certificates   = "/api/v1/rss/common/certificates/get"
)

// Credentials contains the login payload. Keep these values out of logs.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Config contains credentials used by the site's login endpoint. Keep them in
// a secret store or environment variables; do not commit them to source code.
type Config struct {
	BaseURL    string
	StartPath  string
	Username   string
	Password   string
	HTTPClient *http.Client

	// CredentialsProvider takes precedence over Username and Password and is
	// called for every login. It must respect cancellation of ctx.
	CredentialsProvider func(context.Context) (Credentials, error)
	// RefreshBefore defaults to five minutes. Negative values are invalid.
	RefreshBefore time.Duration
	// OnRefreshError runs in the background loop. It must not call Close
	// synchronously, because Close waits for that loop to exit.
	OnRefreshError func(error)
}

// Client maintains the site's session cookie and Bearer token in memory.
// A single Client should be reused for requests to both registries.
type Client struct {
	base       *url.URL
	startPath  string
	httpClient *http.Client

	credentials         Credentials
	credentialsProvider func(context.Context) (Credentials, error)
	refreshBefore       time.Duration
	onRefreshError      func(error)
	wake                chan struct{}
	closed              atomic.Bool

	mu        sync.Mutex
	token     string
	expiresAt time.Time
	refreshAt time.Time
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

	refreshBefore := cfg.RefreshBefore
	if refreshBefore < 0 {
		return nil, fmt.Errorf("refresh before must not be negative")
	}
	if refreshBefore == 0 {
		refreshBefore = 5 * time.Minute
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
		base:                base,
		startPath:           startPath,
		httpClient:          hc,
		credentials:         Credentials{Username: cfg.Username, Password: cfg.Password},
		credentialsProvider: cfg.CredentialsProvider,
		refreshBefore:       refreshBefore,
		onRefreshError:      cfg.OnRefreshError,
		wake:                make(chan struct{}, 1),
	}, nil
}

func (c *Client) endpoint(path string) string {
	u := *c.base
	u.Path = strings.TrimRight(c.base.Path, "/") + path
	return u.String()
}
