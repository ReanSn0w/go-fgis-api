package fgis

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestObservedListModels(t *testing.T) {
	var declarations Page[Declaration]
	declFixture, err := os.ReadFile("testdata/research/declarations-list.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(declFixture, &declarations); err != nil {
		t.Fatal(err)
	}
	if len(declarations.Items) != 2 || declarations.Items[0].DeclarationReplacedNumber == "" {
		t.Fatalf("rare declaration replacement was lost: %+v", declarations.Items)
	}

	var certificates Page[CertificateSummary]
	certFixture, err := os.ReadFile("testdata/research/certificates-list.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(certFixture, &certificates); err != nil {
		t.Fatal(err)
	}
	if len(certificates.Items) != 2 || certificates.Items[0].CertificateReplacedNumber == nil || certificates.Items[1].CertificateReplacedNumber != nil {
		t.Fatalf("nullable certificate replacement was lost: %+v", certificates.Items)
	}
	if certificates.Size != 25 || certificates.Total != 25 {
		t.Fatalf("page counts changed: %+v", certificates)
	}
}

func TestDefaultCertificateQueryMatchesObservedShape(t *testing.T) {
	got, err := json.Marshal(DefaultCertificateQuery())
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"size":10,"page":0,"filter":{"idCertScheme":[],"regDate":{"startDate":null,"endDate":null},"endDate":{"startDate":null,"endDate":null},"columnsSearch":[]},"columnsSort":[{"column":"date","sort":"DESC"}]}`)
	var gotMap, wantMap map[string]any
	if err := json.Unmarshal(got, &gotMap); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &wantMap); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotMap, wantMap) {
		t.Fatalf("query = %s, want %s", got, want)
	}
}

func TestCertificateFilterExtraCannotReplaceKnownField(t *testing.T) {
	q := DefaultCertificateQuery()
	q.Filter.Extra = map[string]json.RawMessage{"status": json.RawMessage(`[6]`)}
	if _, err := json.Marshal(q); err == nil {
		t.Fatal("known certificate filter field was replaced")
	}
	q.Filter.Extra = map[string]json.RawMessage{"idGroupRU": json.RawMessage(`[7]`)}
	if _, err := json.Marshal(q); err != nil {
		t.Fatal(err)
	}
}

func TestRegistryColumnSearchShapesStayDistinct(t *testing.T) {
	declarations := DefaultDeclarationSearchRequest()
	declarations.Filter.AddColumnSearch(DeclarationColumnSearch{Name: "productFullName", Search: "проектор", Type: 0})
	declJSON, err := json.Marshal(declarations)
	if err != nil {
		t.Fatal(err)
	}
	certificates := DefaultCertificateQuery()
	certificates.Filter.ColumnsSearch = append(certificates.Filter.ColumnsSearch, CertificateColumnSearch{Column: "productFullName", Search: "проектор"})
	certJSON, err := json.Marshal(certificates)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(declJSON) || !json.Valid(certJSON) {
		t.Fatal("invalid search JSON")
	}
	var decl, cert map[string]any
	_ = json.Unmarshal(declJSON, &decl)
	_ = json.Unmarshal(certJSON, &cert)
	d := decl["filter"].(map[string]any)["columnsSearch"].([]any)[0].(map[string]any)
	c := cert["filter"].(map[string]any)["columnsSearch"].([]any)[0].(map[string]any)
	if d["name"] != "productFullName" || d["type"] != float64(0) || c["column"] != "productFullName" {
		t.Fatalf("unexpected search forms: declaration=%v, certificate=%v", d, c)
	}
	if _, exists := c["type"]; exists {
		t.Fatalf("certificate search has declaration-only type: %v", c)
	}
}
