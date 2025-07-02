package initiator

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)


func InitMinio(minio_endpoint, minio_access_key, minio_secret_key string,logger utils.Logger) config.MinioClientInterface {
	minioClient, err := config.NewMinioClient(&config.VaultConfig{
		MinioEndPoint:  minio_endpoint,
		MinioAccessKey: minio_access_key,
		MinioSecretKey: minio_secret_key,
	})
	if err != nil {
		logger.Fatalf("failed to initialize minio clinet", err)
	}

	return minioClient
}