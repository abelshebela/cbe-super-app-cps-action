package initiator

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitMinio(minioEndpoint, minioAccessKey, minioSecretKey string, logger utils.Logger) config.MinioClientInterface {
	minioClient, err := config.NewMinioClient(&config.VaultConfig{
		MinioEndPoint:  minioEndpoint,
		MinioAccessKey: minioAccessKey,
		MinioSecretKey: minioSecretKey,
	})
	if err != nil {
		logger.Fatalf("failed to initialize minio clinet", err)
		return nil
	}

	return minioClient
}
