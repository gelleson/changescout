package testdb

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

func NewEntClient() *ent.Client {
	drv, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal(err)
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, drv)))
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}
