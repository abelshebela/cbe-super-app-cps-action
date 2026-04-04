package initiator

import (
	"log"

	logge "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils/logger"
	"go.uber.org/zap"
)

func InitLogger() logge.Logger {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initiate zap logger %v", err)
	}

	return logge.InitLogger(zapLogger)
}
