package fgis

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func assertJSONEqual(t *testing.T, got []byte, want string) {
	t.Helper()
	var gotValue, wantValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("JSON mismatch\ngot: %s\nwant: %s", got, want)
	}
}

func TestObservedSearchBodies(t *testing.T) {
	decl := DefaultDeclarationSearchRequest()
	defaultJSON, err := json.Marshal(decl)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, defaultJSON, `{"size":10,"page":0,"count":0,"filter":{"status":[],"idDeclType":[],"idCertObjectType":[],"idProductType":[],"idGroupRU":[],"idGroupEEU":[],"idTechReg":[],"idApplicantType":[],"regDate":{"minDate":"","maxDate":""},"endDate":{"minDate":"","maxDate":""},"columnsSearch":[]},"columnsSort":[{"column":"declDate","sort":"DESC"}]}`)
	decl.Page = 1
	decl.Count = 1000
	decl.Filter.AddColumnSearch(DeclarationColumnSearch{Name: "productFullName", Search: "проектор", Type: 0})
	filteredJSON, err := json.Marshal(decl)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, filteredJSON, `{"size":10,"page":1,"count":1000,"filter":{"status":[],"idDeclType":[],"idCertObjectType":[],"idProductType":[],"idGroupRU":[],"idGroupEEU":[],"idTechReg":[],"idApplicantType":[],"regDate":{"minDate":"","maxDate":""},"endDate":{"minDate":"","maxDate":""},"columnsSearch":[{"name":"productFullName","search":"проектор","type":0}]},"columnsSort":[{"column":"declDate","sort":"DESC"}]}`)
	cert := DefaultCertificateQuery()
	cert.Page = 1
	cert.Filter.Status = []int{6}
	cert.Filter.IDCertType = []int{1}
	cert.Filter.IDCertObjectType = []int{3}
	cert.Filter.ColumnsSearch = []CertificateColumnSearch{{Column: "manufacterName", Search: "XGIMI"}, {Column: "productFullName", Search: "проектор"}}
	certJSON, err := json.Marshal(cert)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, certJSON, `{"size":10,"page":1,"filter":{"status":[6],"idCertType":[1],"idCertObjectType":[3],"idCertScheme":[],"regDate":{"startDate":null,"endDate":null},"endDate":{"startDate":null,"endDate":null},"columnsSearch":[{"column":"manufacterName","search":"XGIMI"},{"column":"productFullName","search":"проектор"}]},"columnsSort":[{"column":"date","sort":"DESC"}]}`)
}

func TestSearchPageNumberAndCountsArePreserved(t *testing.T) {
	declFixture, err := os.ReadFile("testdata/research/declarations-list.json")
	if err != nil {
		t.Fatal(err)
	}
	certFixture, err := os.ReadFile("testdata/research/certificates-list.json")
	if err != nil {
		t.Fatal(err)
	}
	var pages []int
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rds/declaration":
			http.SetCookie(w, &http.Cookie{Name: "session-cookie", Value: "test", Path: "/"})
		case "/login":
			w.Header().Set("Authorization", "Bearer test-token")
		case declarations, certificates:
			var request struct {
				Page int `json:"page"`
			}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Error(err)
			}
			pages = append(pages, request.Page)
			if r.URL.Path == declarations {
				_, _ = w.Write(declFixture)
			} else {
				_, _ = w.Write(certFixture)
			}
		default:
			http.NotFound(w, r)
		}
	})
	client, err := New(Config{BaseURL: "https://fgis.test", HTTPClient: mockHTTPClient(handler), Username: "user", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	for _, page := range []int{0, 1} {
		decl := DefaultDeclarationSearchRequest()
		decl.Page = page
		got, err := client.SearchDeclarations(context.Background(), decl)
		if err != nil || len(got.Items) != 2 || got.Size != 25 || got.Total != 25 {
			t.Fatalf("declaration page %d: %+v, %v", page, got, err)
		}
		cert := DefaultCertificateQuery()
		cert.Page = page
		found, err := client.SearchCertificateSummaries(context.Background(), cert)
		if err != nil || len(found.Items) != 2 || found.Size != 25 || found.Total != 25 {
			t.Fatalf("certificate page %d: %+v, %v", page, found, err)
		}
	}
	if !reflect.DeepEqual(pages, []int{0, 0, 1, 1}) {
		t.Fatalf("wire pages = %v", pages)
	}
}

func TestObservedSummaryFieldsHaveJSONTags(t *testing.T) {
	for _, item := range []struct {
		fixture string
		model   any
	}{
		{"testdata/research/declarations-list.json", Declaration{}},
		{"testdata/research/certificates-list.json", CertificateSummary{}},
	} {
		data, err := os.ReadFile(item.fixture)
		if err != nil {
			t.Fatal(err)
		}
		var page struct {
			Items []map[string]json.RawMessage `json:"items"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			t.Fatal(err)
		}
		tags := make(map[string]struct{})
		typ := reflect.TypeOf(item.model)
		for i := 0; i < typ.NumField(); i++ {
			tags[typ.Field(i).Tag.Get("json")] = struct{}{}
		}
		for _, row := range page.Items {
			for key := range row {
				if _, ok := tags[key]; !ok {
					t.Errorf("%s has unmapped JSON field %s", item.fixture, key)
				}
			}
		}
	}
}
