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

func postList[T any](ctx context.Context, c *Client, path, referer string, query any) (Page[T], error) {
	var result Page[T]
	payload, err := json.Marshal(query)
	if err != nil {
		return result, fmt.Errorf("encode search request: %w", err)
	}
	token, err := c.ensureToken(ctx)
	if err != nil {
		return result, err
	}
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(path), bytes.NewReader(payload))
		if err != nil {
			return result, fmt.Errorf("create search request: %w", err)
		}
		setBrowserHeaders(req, c.base, c.endpoint(referer))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return result, fmt.Errorf("send search request: %w", err)
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
			return result, &HTTPError{Method: http.MethodPost, Path: path, StatusCode: resp.StatusCode}
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		closeErr := resp.Body.Close()
		if err != nil {
			return result, fmt.Errorf("decode FGIS list response: %w", err)
		}
		if closeErr != nil {
			return result, fmt.Errorf("close FGIS list response: %w", closeErr)
		}
		return result, nil
	}
	return result, fmt.Errorf("FGIS list retry exhausted")
}
