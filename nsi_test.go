package fgis

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestObservedIdentifierAndNSIShapes(t *testing.T) {
	declBytes, err := os.ReadFile("testdata/research/declaration-identifiers.json")
	if err != nil {
		t.Fatal(err)
	}
	var declaration DeclarationIdentifiers
	if err := json.Unmarshal(declBytes, &declaration); err != nil {
		t.Fatal(err)
	}
	if declaration.Status["active"].ID != 6 || declaration.OKSM["russia"].ID == "" || len(declaration.DeclScheme["vpDeclScheme"].Raw) == 0 {
		t.Fatalf("declaration identifiers lost fields: %+v", declaration)
	}
	certBytes, err := os.ReadFile("testdata/research/certificate-identifiers.json")
	if err != nil {
		t.Fatal(err)
	}
	var certificate CertificateIdentifiers
	if err := json.Unmarshal(certBytes, &certificate); err != nil {
		t.Fatal(err)
	}
	if certificate.Status["active"].ID != 6 || certificate.OKSM["russia"].ID == "" || len(certificate.CertType) == 0 {
		t.Fatalf("certificate identifiers lost fields: %+v", certificate)
	}
	for _, fixture := range []string{"nsi-multi-declaration.json", "nsi-multi-certificate.json"} {
		data, err := os.ReadFile("testdata/research/" + fixture)
		if err != nil {
			t.Fatal(err)
		}
		var response NSIResponse
		if err := json.Unmarshal(data, &response); err != nil {
			t.Fatalf("%s: %v", fixture, err)
		}
		if len(response.TNVED) == 0 || len(response.Groups["tnved"]) == 0 || len(response.Status) == 0 {
			t.Fatalf("%s lost NSI groups", fixture)
		}
	}
}

func TestResolveProductCodesPreservesOrderAndMissingIDs(t *testing.T) {
	var gotQuery NSIQuery
	reads := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			w.Header().Set("Authorization", "Bearer test-token")
		case "/nsi/api/multi":
			reads++
			if r.Method != http.MethodPost {
				t.Errorf("NSI method = %s", r.Method)
			}
			if err := json.NewDecoder(r.Body).Decode(&gotQuery); err != nil {
				t.Error(err)
			}
			_, _ = w.Write([]byte(`{"tnved":[{"id":101,"name":"Projectors","code":"8528629009","masterId":"101","hidden":false}],"okpd2":[{"id":201,"name":"Display devices","code":"26.20","hidden":false}],"futureGroup":[{"id":9}]}`))
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	got, err := client.ResolveProductCodes(context.Background(), []int64{101, 999, 101}, []int64{201, 202})
	if err != nil {
		t.Fatal(err)
	}
	if reads != 1 || len(got.TNVED) != 3 || len(got.OKPD2) != 2 || got.TNVED[0].Code == nil || got.TNVED[1].Code != nil || got.TNVED[2].ID != 101 || got.OKPD2[1].Name != nil {
		t.Fatalf("resolved codes = %+v, reads=%d", got, reads)
	}
	if !reflect.DeepEqual(gotQuery.Items["tnved"][0].ID, []any{float64(101), float64(999)}) || !reflect.DeepEqual(gotQuery.Items["okpd2"][0].ID, []any{float64(201), float64(202)}) {
		t.Fatalf("request IDs were not deduplicated: %+v", gotQuery)
	}
	if _, err := client.ResolveProductCodes(context.Background(), []int64{0}, nil); err == nil || reads != 1 {
		t.Fatalf("invalid ID reached NSI: err=%v reads=%d", err, reads)
	}
}
