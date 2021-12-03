// Copyright 2021 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/apigee/registry-experimental/server/search/internal/storage/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Client represents a connection to a storage provider.
type Client struct {
	db *gorm.DB
}

var mutex sync.Mutex
var disableMutex bool

func lock() {
	if !disableMutex {
		mutex.Lock()
	}
}

func unlock() {
	if !disableMutex {
		mutex.Unlock()
	}
}

func defaultConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // https://gorm.io/docs/logger.html
	}
}

// NewClient creates a new database session using the provided driver and data source name.
// Driver must be one of [ postgres, cloudsqlpostgres ]. DSN format varies per database driver.
//
// PostgreSQL DSN Reference: See "Connection Strings" at https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
// SQLite DSN Reference: See "URI filename examples" at https://www.sqlite.org/c3ref/open.html
func NewClient(ctx context.Context, driver, dsn string) (*Client, error) {
	lock()
	switch driver {
	case "postgres", "cloudsqlpostgres":
		db, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: driver,
			DSN:        dsn,
		}), defaultConfig())
		if err != nil {
			c := &Client{db: db}
			c.close()
			unlock()
			return nil, err
		}
		unlock()
		// postgres runs in a separate process and seems to have no problems
		// with concurrent access and modifications.
		disableMutex = true
		c := &Client{db: db}
		c.ensure()
		return c, nil
	default:
		unlock()
		return nil, fmt.Errorf("search: unsupported database %s", driver)
	}
}

// Close closes a database session.
func (c *Client) Close() {
	lock()
	defer unlock()
	c.close()
}

func (c *Client) close() {
	sqlDB, _ := c.db.DB()
	sqlDB.Close()
}

func (c *Client) ensureTable(v interface{}) {
	lock()
	defer unlock()
	if !c.db.Migrator().HasTable(v) {
		_ = c.db.Migrator().CreateTable(v)
	}
}

func (c *Client) ensure() {
	c.ensureTable(&models.Document{})
}

// IsNotFound returns true if an error is due to an entity not being found.
func (c *Client) IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}

// Get gets an entity using the storage client.
func (c *Client) Get(ctx context.Context, k *Key, v interface{}) error {
	lock()
	defer unlock()
	return c.db.Where("key = ?", k.Name).First(v).Error
}

// Put puts an entity using the storage client.
func (c *Client) PutDocument(ctx context.Context, r *models.Document) error {
	lock()
	defer unlock()
	_ = c.db.Transaction(func(tx *gorm.DB) error {
		// Update all fields from model: https://gorm.io/docs/update.html#Update-Selected-Fields
		rowsAffected := tx.Model(r).Select("*").Where("key = ?", r.Key).Updates(r).RowsAffected
		if rowsAffected == 0 {
			tx.Create(r)
		}
		return nil
	})
	return nil
}

// Delete deletes all entities matching a key.
func (c *Client) Delete(ctx context.Context, q *Query) error {
	op := c.db
	for _, r := range q.Requirements {
		op = op.Where(r.Name+" = ?", r.Value)
	}
	switch q.Kind {
	case "Document":
		return op.Delete(models.Document{}).Error
	}
	return nil
}

// Run runs a query using the storage client, returning an iterator.
func (c *Client) Run(ctx context.Context, q *Query) *Iterator {
	lock()
	defer unlock()

	// Filtering is currently implemented by skipping iterator elements that
	// don't match the filter criteria, and expects to only reach the end of
	// the iterator if there are no more resources to consider. Previously,
	// the entire table would be read into memory. This limit should maintain
	// that behavior until we improve our iterator implementation.
	op := c.db.Offset(q.Offset).Limit(100000)
	for _, r := range q.Requirements {
		op = op.Where(r.Name+" = ?", r.Value)
	}

	if order := q.Order; order != "" {
		op = op.Order(order)
	} else {
		op = op.Order("key")
	}

	switch q.Kind {
	case "Document":
		var v []models.Document
		_ = op.Find(&v).Error
		return &Iterator{Client: c, Values: v, Index: 0}
	default:
		return nil
	}
}

type RawRows interface {
	Append(rows *sql.Rows) error
}

func (c *Client) Raw(ctx context.Context, target RawRows, sql string, values ...interface{}) error {
	lock()
	defer unlock()

	rows, err := c.db.Raw(sql, values).Rows()
	if err != nil {
		return err
	}

	defer rows.Close()
	return target.Append(rows)
}
