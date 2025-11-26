package initiator

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// "github.com/aws/aws-sdk-go-v2/aws"
	configAw "github.com/aws/aws-sdk-go-v2/config"
)

func InitMinio(cfgMain config.VaultConfig, logger utils.Logger) aws.Config{

	cfg, err:= configAw.LoadDefaultConfig(context.TODO(),
		configAw.WithRegion(cfgMain.S3BucketName),      // secrets.AWS_REGION
		configAw.WithBaseEndpoint(cfgMain.S3BucketURL), // secrets.AWS_BUCKET_URL
		configAw.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfgMain.S3AccessKeyID,     // secrets.AWS_ACCESS_KEY_ID
				cfgMain.S3SecretAccessKey, // secrets.AWS_SECRET_ACCESS_KEY
				"",
			),
		),
	)
if err != nil {
	logger.Errorf("Failed to load AWS config: %v", err)
}
	return cfg
}
