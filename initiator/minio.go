package initiator

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	// "github.com/aws/aws-sdk-go-v2/aws"
	configAw "github.com/aws/aws-sdk-go-v2/config"
)

func InitMinio(cfgMain config.VaultConfig, logger utils.Logger) *s3.Client {

	cfg, err := configAw.LoadDefaultConfig(context.TODO(),
		configAw.WithRegion("us-east-1"),
		configAw.WithBaseEndpoint(cfgMain.S3BucketURL),
		configAw.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfgMain.S3AccessKeyID,
				cfgMain.S3SecretAccessKey,
				"",
			),
		),
	)

	if err != nil {
		logger.Errorf("failed to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.DisableLogOutputChecksumValidationSkipped = true
	})
	return s3Client
}
