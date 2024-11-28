package ent

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"github.com/gelleson/changescout/changescout/internal/platform/logger"
	"go.uber.org/zap"

	_ "github.com/glebarez/go-sqlite"
)

type BuildConfig struct {
	DBEngine string
	DBURL    string
}

func Build(ctx context.Context, conf *BuildConfig) *ent.Client {
	log := logger.L("ent")
	drv, err := sql.Open(conf.DBEngine, conf.DBURL)
	if err != nil {
		log.Fatal("failed to open database connection", zap.Error(err))
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, drv)))
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal("failed to create schema resources", zap.Error(err))
	}

	return client
}
