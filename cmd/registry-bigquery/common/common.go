package common

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/googleapi"
)

// Used by subcommands, now provides a single common timestap for a command invocation.
var Now = time.Now()

// Get a BigQuery dataset by name and create it if it doesn't exist.
func GetOrCreateDataset(ctx context.Context, client *bigquery.Client, name string) (*bigquery.Dataset, error) {
	dataset := client.Dataset(name)
	if err := dataset.Create(ctx, nil); err != nil {
		switch v := err.(type) {
		case *googleapi.Error:
			if v.Code != http.StatusConflict { // already exists
				return nil, err
			}
		default:
			return nil, err
		}
	}
	return dataset, nil
}

// Get a BigQuery table by name and create it if it doesn't exist.
func GetOrCreateTable(ctx context.Context, dataset *bigquery.Dataset, name string, prototype interface{}) (*bigquery.Table, error) {
	table := dataset.Table(name)
	schema, err := bigquery.InferSchema(prototype)
	if err != nil {
		return nil, err
	}
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		switch v := err.(type) {
		case *googleapi.Error:
			if v.Code != http.StatusConflict { // already exists
				return nil, err
			}
		default:
			return nil, err
		}
	}
	return table, nil
}
