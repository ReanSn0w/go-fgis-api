package fgis_test

import (
	"context"
	"fmt"
	"os"
	"time"

	fgis "github.com/ReanSn0w/go-fgis-api"
)

// This example is compiled but not run: it deliberately has no Output comment.
func ExampleNew() {
	client, err := fgis.New(fgis.Config{
		Username: os.Getenv("FGIS_USERNAME"),
		Password: os.Getenv("FGIS_PASSWORD"),
	})
	if err != nil {
		return
	}
	defer client.Close()

	query := fgis.DefaultDeclarationSearchRequest()
	query.Page = 0
	page, err := client.SearchDeclarations(context.Background(), query)
	if err != nil {
		return
	}
	for _, declaration := range page.Items {
		fmt.Println(declaration.DeclarationNumber)
	}
}

// This example is compiled but not run: it deliberately has no Output comment.
func ExampleClient_Start() {
	client, err := fgis.New(fgis.Config{
		CredentialsProvider: func(ctx context.Context) (fgis.Credentials, error) {
			return fgis.Credentials{
				Username: os.Getenv("FGIS_USERNAME"),
				Password: os.Getenv("FGIS_PASSWORD"),
			}, nil
		},
		RefreshBefore: 5 * time.Minute,
	})
	if err != nil {
		return
	}
	defer client.Close()
	processCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := client.Start(processCtx); err != nil {
		return
	}
}

// The example compiles but needs credentials to run against the public site.
func ExampleClient_GetDeclaration() {
	client, err := fgis.New(fgis.Config{
		Username: os.Getenv("FGIS_USERNAME"),
		Password: os.Getenv("FGIS_PASSWORD"),
	})
	if err != nil {
		return
	}
	defer client.Close()
	ctx := context.Background()
	query := fgis.DefaultDeclarationSearchRequest()
	query.Filter.AddColumnSearch(fgis.DeclarationColumnSearch{Name: "productFullName", Search: "проектор", Type: 0})
	page, err := client.SearchDeclarations(ctx, query)
	if err != nil || len(page.Items) == 0 {
		return
	}
	card, err := client.GetDeclaration(ctx, int64(page.Items[0].ID))
	if err != nil {
		return
	}
	for _, lab := range card.TestingLabs {
		for _, protocol := range lab.Protocols {
			fmt.Println(protocol.Number)
		}
	}
	for _, product := range card.Product.Identifications {
		codes, err := client.ResolveProductCodes(ctx, product.IDTnveds, product.IDOkpds)
		if err != nil {
			return
		}
		for _, code := range codes.TNVED {
			if code.Code != nil {
				fmt.Println(*code.Code)
			}
		}
		for _, code := range codes.OKPD2 {
			if code.Code != nil {
				fmt.Println(*code.Code)
			}
		}
	}
}

// Typed and original raw certificate searches can be used side by side.
func ExampleClient_SearchCertificateSummaries() {
	client, err := fgis.New(fgis.Config{
		Username:  os.Getenv("FGIS_USERNAME"),
		Password:  os.Getenv("FGIS_PASSWORD"),
		StartPath: "/rss/certificate",
	})
	if err != nil {
		return
	}
	defer client.Close()
	ctx := context.Background()
	query := fgis.DefaultCertificateQuery()
	query.Filter.ColumnsSearch = []fgis.CertificateColumnSearch{{Column: "productFullName", Search: "проектор"}}
	typed, err := client.SearchCertificateSummaries(ctx, query)
	if err != nil {
		return
	}
	for _, row := range typed.Items {
		fmt.Println(row.Number)
	}
	raw, err := client.SearchCertificates(ctx, fgis.CertificateSearchRequest{Size: 10, Filter: map[string]any{}})
	if err != nil {
		return
	}
	for _, row := range raw.Items {
		fmt.Println(string(row["number"]))
	}
}
