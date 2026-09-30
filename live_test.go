package fgis

import (
	"context"
	"os"
	"testing"
	"time"
)

// Run explicitly with FGIS_LIVE_TEST=1 and credentials in the environment.
// Normal go test and GitHub Actions never contact the real registry.
func TestLiveDeclarationsAPI(t *testing.T) {
	if os.Getenv("FGIS_LIVE_TEST") != "1" {
		t.Skip("set FGIS_LIVE_TEST=1 to call the real FGIS API")
	}
	username, password := os.Getenv("FGIS_USERNAME"), os.Getenv("FGIS_PASSWORD")
	if username == "" || password == "" {
		t.Fatal("FGIS_USERNAME and FGIS_PASSWORD are required for live test")
	}
	client, err := New(Config{Username: username, Password: password})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	page, err := client.SearchDeclarations(ctx, DefaultDeclarationSearchRequest())
	if err != nil {
		t.Fatal(err)
	}
	if page.Total < 1 || len(page.Items) < 1 || page.Items[0].ID < 1 {
		t.Fatalf("real API returned no identifiable declaration (total=%d, items=%d)", page.Total, len(page.Items))
	}
}
