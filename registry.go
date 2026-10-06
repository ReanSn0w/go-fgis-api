package fgis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SearchDeclarations returns typed rows from the declaration list.
func (c *Client) SearchDeclarations(ctx context.Context, query DeclarationSearchRequest) (Page[Declaration], error) {
	return postList[Declaration](ctx, c, declarations, "/rds/declaration", query)
}

// SearchCertificates uses an extensible request and raw row fields because a
// certificate list response was not present in the supplied capture.
func (c *Client) SearchCertificates(ctx context.Context, query CertificateSearchRequest) (Page[Certificate], error) {
	return postList[Certificate](ctx, c, certificates, "/rss/certificate", query)
}

// SearchCertificateSummaries returns typed rows using the confirmed certificate
// query shape. SearchCertificates remains available for raw fields.
func (c *Client) SearchCertificateSummaries(ctx context.Context, query CertificateQuery) (Page[CertificateSummary], error) {
	return postList[CertificateSummary](ctx, c, certificates, "/rss/certificate", query)
}

// GetDeclaration retrieves one complete public declaration card by registry ID.
func (c *Client) GetDeclaration(ctx context.Context, id int64) (DeclarationDetails, error) {
	path, err := detailPath("/api/v1/rds/common/declarations", id)
	if err != nil {
		return DeclarationDetails{}, err
	}
	return requestJSON[DeclarationDetails](ctx, c, http.MethodGet, path, fmt.Sprintf("/rds/declaration/view/%d/common", id), nil)
}

// GetCertificate retrieves one complete public certificate card by registry ID.
func (c *Client) GetCertificate(ctx context.Context, id int64) (CertificateDetails, error) {
	path, err := detailPath("/api/v1/rss/common/certificates", id)
	if err != nil {
		return CertificateDetails{}, err
	}
	return requestJSON[CertificateDetails](ctx, c, http.MethodGet, path, fmt.Sprintf("/rss/certificate/view/%d/baseInfo", id), nil)
}

func postList[T any](ctx context.Context, c *Client, path, referer string, query any) (Page[T], error) {
	return requestJSON[Page[T]](ctx, c, http.MethodPost, path, referer, query)
}

// requestJSON is shared by list, detail and dictionary reads. A request body
// is present only for POST; GET requests never send the JSON literal null.
func requestJSON[T any](ctx context.Context, c *Client, method, path, referer string, body any) (T, error) {
	var result T
	if c.closed.Load() {
		return result, ErrClosed
	}
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return result, fmt.Errorf("encode FGIS request: %w", err)
		}
	}
	token, err := c.ensureToken(ctx)
	if err != nil {
		return result, err
	}
	for attempt := 0; attempt < 2; attempt++ {
		var reqBody io.Reader
		if payload != nil {
			reqBody = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.endpoint(path), reqBody)
		if err != nil {
			return result, fmt.Errorf("create FGIS request: %w", err)
		}
		setBrowserHeaders(req, c.base, c.endpoint(referer))
		req.Header.Set("Authorization", token)
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return result, fmt.Errorf("send FGIS request: %w", err)
		}
		if resp.StatusCode == http.StatusUnauthorized && attempt == 0 {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			token, err = c.refreshAfterUnauthorized(ctx, token)
			if err != nil {
				return result, err
			}
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return result, &HTTPError{Method: method, Path: path, StatusCode: resp.StatusCode}
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		closeErr := resp.Body.Close()
		if err != nil {
			return result, fmt.Errorf("decode FGIS response: %w", err)
		}
		if closeErr != nil {
			return result, fmt.Errorf("close FGIS list response: %w", closeErr)
		}
		return result, nil
	}
	return result, fmt.Errorf("FGIS request retry exhausted")
}

func detailPath(prefix string, id int64) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("FGIS detail ID must be positive: %d", id)
	}
	return fmt.Sprintf("%s/%d", prefix, id), nil
}
