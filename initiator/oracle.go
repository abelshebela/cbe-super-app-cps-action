package initiator

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	_ "github.com/godror/godror"
	"go.uber.org/zap"
)

func InitOracle(url string, logger utils.Logger) *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	values := strings.Split(url, "@")
	if len(values) != 2 {
		logger.Fatalf("invalid oracle connection URL format", zap.String("url", url))
	}

	userPass := strings.Split(values[0], "/")
	if len(userPass) != 2 {
		logger.Fatalf("invalid oracle connection URL format - missing user/password", zap.String("url", url))
	}

	db, err := sql.Open("godror", url)
	if err != nil {
		logger.Fatalf("failed to open oracle connection", zap.Error(err))
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		logger.Fatalf("failed to ping oracle", zap.Error(err))
	}

	logger.Infof("oracle connection is established using godror driver")
	return db
}

func CloseOracle(db *sql.DB, log utils.Logger) {
	if err := db.Close(); err != nil {
		log.Fatalf("failed to close oracle connection", zap.Error(err))
	}
}
