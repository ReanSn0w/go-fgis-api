package fgis

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestDetailMissingNullEmptyAndUnknownFields(t *testing.T) {
	var detail CertificateDetails
	input := []byte(`{"idCertificate":5,"number":"","replacement":null,"annexes":[],"futureField":{"enabled":true}}`)
	if err := json.Unmarshal(input, &detail); err != nil {
		t.Fatal(err)
	}
	if detail.IDCertificate != 5 || detail.Number != "" || detail.Replacement != nil || detail.ReplacedBy != nil || detail.Annexes == nil || len(detail.Annexes) != 0 {
		t.Fatalf("missing, null and empty fields: %+v", detail)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(detail.Raw, &raw); err != nil || string(raw["futureField"]) != `{"enabled":true}` {
		t.Fatalf("unknown raw field: %s, %v", detail.Raw, err)
	}
	var declaration DeclarationDetails
	if err := json.Unmarshal([]byte(`{"idDeclaration":7,"annexes":[{"idAnnex":1,"ord":"2"}],"future":null}`), &declaration); err != nil {
		t.Fatal(err)
	}
	if declaration.Annexes[0].Ord != "2" || !json.Valid(declaration.Raw) {
		t.Fatalf("declaration annex/raw: %+v", declaration)
	}
	if err := json.Unmarshal([]byte(`{"idCertificate":5,"annexes":[{"idAnnex":1,"ord":2}]}`), &detail); err != nil || detail.Annexes[0].Ord != 2 {
		t.Fatalf("certificate annex: %+v, %v", detail, err)
	}
}

func TestIdentifierRoutesAndRawNSIGroup(t *testing.T) {
	paths := []string{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			w.Header().Set("Authorization", "Bearer test-token")
		case "/api/v1/rds/common/identifiers":
			paths = append(paths, r.Method+" "+r.URL.Path)
			_, _ = w.Write([]byte(`{"status":{"active":{"id":6,"name":"Active"}},"future":{"id":8}}`))
		case "/api/v1/rss/common/identifiers":
			paths = append(paths, r.Method+" "+r.URL.Path)
			_, _ = w.Write([]byte(`{"status":{"active":{"id":6,"name":"Active"}},"future":{"id":9}}`))
		case "/nsi/api/multi":
			paths = append(paths, r.Method+" "+r.URL.Path)
			var body NSIQuery
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Items["future"]) != 1 {
				t.Errorf("NSI body: %+v, %v", body, err)
			}
			_, _ = w.Write([]byte(`{"future":[{"id":8,"value":"kept"}]}`))
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	declaration, err := client.GetDeclarationIdentifiers(context.Background())
	if err != nil || declaration.Status["active"].ID != 6 || len(declaration.Raw) == 0 {
		t.Fatalf("declaration identifiers: %+v, %v", declaration, err)
	}
	certificate, err := client.GetCertificateIdentifiers(context.Background())
	if err != nil || certificate.Status["active"].ID != 6 || len(certificate.Raw) == 0 {
		t.Fatalf("certificate identifiers: %+v, %v", certificate, err)
	}
	response, err := client.GetNSIMulti(context.Background(), NSIQuery{Items: map[string][]NSISelection{"future": {{ID: []any{8}, Fields: []string{"id", "value"}}}}})
	if err != nil || string(response.Groups["future"]) != `[{"id":8,"value":"kept"}]` {
		t.Fatalf("unknown NSI group: %+v, %v", response, err)
	}
	if want := []string{"GET /api/v1/rds/common/identifiers", "GET /api/v1/rss/common/identifiers", "POST /nsi/api/multi"}; !reflect.DeepEqual(paths, want) {
		t.Fatalf("paths = %v, want %v", paths, want)
	}
}

func TestLegacySearchFormsRemainUsable(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			w.Header().Set("Authorization", "Bearer test-token")
		case declarations:
			var body DeclarationSearchRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || string(body.Filter.ColumnsSearch[0]) != `{"name":"productFullName","search":"legacy","type":0}` {
				t.Errorf("legacy declaration body: %+v, %v", body, err)
			}
			_, _ = w.Write([]byte(`{"items":[],"size":0,"total":0}`))
		case certificates:
			var body CertificateSearchRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Filter["custom"] != "legacy" {
				t.Errorf("legacy certificate body: %+v, %v", body, err)
			}
			_, _ = w.Write([]byte(`{"items":[{"future":"still raw"}],"size":1,"total":1}`))
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	decl := DefaultDeclarationSearchRequest()
	decl.Filter.ColumnsSearch = []json.RawMessage{json.RawMessage(`{"name":"productFullName","search":"legacy","type":0}`)}
	if _, err := client.SearchDeclarations(context.Background(), decl); err != nil {
		t.Fatal(err)
	}
	cert, err := client.SearchCertificates(context.Background(), CertificateSearchRequest{Size: 10, Filter: map[string]any{"custom": "legacy"}})
	if err != nil || string(cert.Items[0]["future"]) != `"still raw"` {
		t.Fatalf("raw certificate row: %+v, %v", cert, err)
	}
}

func TestNewRequestErrors(t *testing.T) {
	for _, status := range []int{404, 500, 401} {
		for _, method := range []string{"GET", "POST"} {
			t.Run(method+"/"+http.StatusText(status), func(t *testing.T) {
				reads := 0
				client := newErrorTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					reads++
					w.WriteHeader(status)
				})
				defer client.Close()
				err := callNewRead(context.Background(), client, method)
				var httpErr *HTTPError
				if !errors.As(err, &httpErr) || httpErr.StatusCode != status || httpErr.Method != method {
					t.Fatalf("HTTP error = %v", err)
				}
				wantReads := 1
				if status == 401 {
					wantReads = 2
				}
				if reads != wantReads {
					t.Fatalf("reads = %d, want %d", reads, wantReads)
				}
			})
		}
	}
	for _, method := range []string{"GET", "POST"} {
		t.Run(method+"/invalid-JSON", func(t *testing.T) {
			client := newErrorTestClient(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`{invalid`)) })
			defer client.Close()
			if err := callNewRead(context.Background(), client, method); err == nil {
				t.Fatal("invalid JSON accepted")
			}
		})
		t.Run(method+"/cancelled", func(t *testing.T) {
			client := newErrorTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Error("cancelled request reached endpoint") })
			defer client.Close()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			if err := callNewRead(ctx, client, method); !errors.Is(err, context.Canceled) {
				t.Fatalf("cancelled request = %v", err)
			}
		})
		t.Run(method+"/closed", func(t *testing.T) {
			client := newErrorTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Error("closed client reached endpoint") })
			if err := client.Close(); err != nil {
				t.Fatal(err)
			}
			if err := callNewRead(context.Background(), client, method); !errors.Is(err, ErrClosed) {
				t.Fatalf("closed client = %v", err)
			}
		})
	}
}

func newErrorTestClient(t *testing.T, endpoint http.HandlerFunc) *Client {
	t.Helper()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			w.Header().Set("Authorization", "Bearer test-token")
		case "/api/v1/rds/common/declarations/1", "/nsi/api/multi":
			endpoint(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func callNewRead(ctx context.Context, client *Client, method string) error {
	if method == "GET" {
		_, err := client.GetDeclaration(ctx, 1)
		return err
	}
	_, err := client.GetNSIMulti(ctx, NSIQuery{Items: map[string][]NSISelection{"tnved": {{ID: []any{1}, Fields: []string{"name"}}}}})
	return err
}
