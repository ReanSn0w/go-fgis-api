package fgis

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func TestObservedDetailShapes(t *testing.T) {
	for _, name := range []string{"declaration-serial.json", "declaration-projector.json"} {
		data, err := os.ReadFile("testdata/research/" + name)
		if err != nil {
			t.Fatal(err)
		}
		var detail DeclarationDetails
		if err := json.Unmarshal(data, &detail); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if detail.IDDeclaration == 0 || len(detail.Product.Identifications) == 0 || len(detail.TestingLabs) == 0 || !json.Valid(detail.Raw) {
			t.Fatalf("%s lost detail fields", name)
		}
	}
	for _, name := range []string{"certificate-projector.json", "certificate-replaced.json", "certificate-replacement.json"} {
		data, err := os.ReadFile("testdata/research/" + name)
		if err != nil {
			t.Fatal(err)
		}
		var detail CertificateDetails
		if err := json.Unmarshal(data, &detail); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if detail.IDCertificate == 0 || len(detail.Product.Identifications) == 0 || len(detail.TestingLabs) == 0 || !json.Valid(detail.Raw) {
			t.Fatalf("%s lost detail fields", name)
		}
		if name == "certificate-projector.json" && (len(detail.TestingLabs) != 2 || len(detail.Product.Identifications[0].Standards) != 8) {
			t.Fatalf("projector nested data incomplete: labs=%d standards=%d", len(detail.TestingLabs), len(detail.Product.Identifications[0].Standards))
		}
		if name == "certificate-replaced.json" && (detail.ReplacedBy == nil || detail.ReplacedBy.CertSource != 100 || detail.ReplacedBy.CertTarget != 101) {
			t.Fatalf("replacement link reversed or missing: %+v", detail.ReplacedBy)
		}
		if name == "certificate-replacement.json" && (detail.Replacement == nil || detail.Replacement.CertSource != 100 || detail.Replacement.CertTarget != 101) {
			t.Fatalf("reverse replacement link missing: %+v", detail.Replacement)
		}
	}
}

func TestGetDetailsUsesObservedRoutes(t *testing.T) {
	declaration, err := os.ReadFile("testdata/research/declaration-projector.json")
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := os.ReadFile("testdata/research/certificate-projector.json")
	if err != nil {
		t.Fatal(err)
	}
	reads := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			w.Header().Set("Authorization", "Bearer test-token")
		case "/api/v1/rds/common/declarations/1":
			reads++
			if r.Method != http.MethodGet {
				t.Errorf("declaration method = %s", r.Method)
			}
			_, _ = w.Write(declaration)
		case "/api/v1/rss/common/certificates/102":
			reads++
			if r.Method != http.MethodGet {
				t.Errorf("certificate method = %s", r.Method)
			}
			_, _ = w.Write(certificate)
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	if _, err := client.GetDeclaration(context.Background(), 0); err == nil || reads != 0 {
		t.Fatalf("invalid ID triggered HTTP: %v, reads=%d", err, reads)
	}
	decl, err := client.GetDeclaration(context.Background(), 1)
	if err != nil || decl.IDDeclaration != 1 || len(decl.Raw) == 0 {
		t.Fatalf("declaration = %+v, %v", decl, err)
	}
	cert, err := client.GetCertificate(context.Background(), 102)
	if err != nil || cert.IDCertificate != 102 || len(cert.Raw) == 0 || reads != 2 {
		t.Fatalf("certificate ID=%d raw=%d reads=%d err=%v", cert.IDCertificate, len(cert.Raw), reads, err)
	}
}
