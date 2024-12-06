package ent

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	_ "github.com/glebarez/go-sqlite"
)

type BuildConfig struct {
	DBEngine string
	DBURL    string
}

func Build(ctx context.Context, conf *BuildConfig) (*ent.Client, error) {
	drv, err := sql.Open(conf.DBEngine, conf.DBURL)
	if err != nil {
		return nil, err
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, drv)))
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
