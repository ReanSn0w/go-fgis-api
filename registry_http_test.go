package fgis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestRequestJSONGetRetriesAfterUnauthorizedWithoutBody(t *testing.T) {
	logins, reads := 0, 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			logins++
			w.Header().Set("Authorization", fmt.Sprintf("Bearer token-%d", logins))
		case "/api/v1/rds/common/declarations/123":
			reads++
			if r.Method != http.MethodGet || r.ContentLength != 0 || r.Header.Get("Content-Type") != "" {
				t.Errorf("detail request has wrong method/body: %s %d %q", r.Method, r.ContentLength, r.Header.Get("Content-Type"))
			}
			if reads == 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_, _ = w.Write([]byte(`{"idDeclaration":123}`))
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	got, err := requestJSON[struct {
		ID int `json:"idDeclaration"`
	}](context.Background(), client, http.MethodGet, "/api/v1/rds/common/declarations/123", "/rds/declaration/view/123/common", nil)
	if err != nil || got.ID != 123 || logins != 2 || reads != 2 {
		t.Fatalf("GET retry: value=%+v err=%v logins=%d reads=%d", got, err, logins, reads)
	}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := requestJSON[map[string]any](context.Background(), client, http.MethodGet, "/api/v1/rds/common/identifiers", "/rds/declaration", nil); !errors.Is(err, ErrClosed) {
		t.Fatalf("closed client error = %v", err)
	}
}

func TestDetailPathRejectsNonPositiveID(t *testing.T) {
	for _, id := range []int64{-1, 0} {
		if _, err := detailPath("/api/v1/rds/common/declarations", id); err == nil {
			t.Fatalf("accepted ID %d", id)
		}
	}
	path, err := detailPath("/api/v1/rds/common/declarations", 123)
	if err != nil || path != "/api/v1/rds/common/declarations/123" {
		t.Fatalf("detail path = %q, %v", path, err)
	}
}
