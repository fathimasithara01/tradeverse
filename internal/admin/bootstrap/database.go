package bootstrap

import (
	"context"
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/database"
	"gorm.io/gorm"
)

func InitDatabase(ctx context.Context, cfg *config.Config) *gorm.DB {

	db, err := database.ConnectDB(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	return db
}
