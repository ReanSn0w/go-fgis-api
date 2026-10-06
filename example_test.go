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
