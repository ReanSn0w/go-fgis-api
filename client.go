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

	// Start is serialized separately from cancellation, so Close can cancel
	// a login even while that login holds the authentication mutex.
	startMu     sync.Mutex
	lifecycleMu sync.Mutex
	loop        *backgroundLoop
}

type backgroundLoop struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
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

// Start obtains a token and starts one background refresh loop. A running
// loop is reused; after its context is canceled, Start can start a new loop.
func (c *Client) Start(ctx context.Context) error {
	c.startMu.Lock()
	defer c.startMu.Unlock()

	for {
		c.lifecycleMu.Lock()
		if c.closed.Load() {
			c.lifecycleMu.Unlock()
			return ErrClosed
		}
		if err := ctx.Err(); err != nil {
			c.lifecycleMu.Unlock()
			return err
		}
		previous := c.loop
		if previous != nil {
			select {
			case <-previous.done:
				c.loop = nil
			default:
				c.lifecycleMu.Unlock()
				if previous.ctx.Err() == nil {
					return nil
				}
				select {
				case <-previous.done:
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		loopCtx, cancel := context.WithCancel(ctx)
		loop := &backgroundLoop{ctx: loopCtx, cancel: cancel, done: make(chan struct{})}
		c.loop = loop
		c.lifecycleMu.Unlock()

		_, err := c.ensureToken(loopCtx)
		if c.closed.Load() {
			err = ErrClosed
		} else if err == nil {
			err = loopCtx.Err()
		}
		if err != nil {
			cancel()
			close(loop.done)
			return err
		}
		go c.refreshLoop(loopCtx, loop.done)
		return nil
	}
}

// Close permanently rejects new operations, cancels background work and waits
// for it to finish. In-flight caller requests retain their own contexts.
// Close is safe to call repeatedly, including before Start.
func (c *Client) Close() error {
	c.lifecycleMu.Lock()
	c.closed.Store(true)
	loop := c.loop
	if loop != nil {
		loop.cancel()
	}
	c.lifecycleMu.Unlock()
	if loop != nil {
		<-loop.done
	}
	return nil
}
