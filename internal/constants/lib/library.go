package lib

import (
	"bytes"
	"cbe-super-app-cps-action/internal/constants"

	// "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CpsModelBuilder(unique string, makerUser types.UserContext, prevAction, currentAction any, requestAction, actionType string) model.CPSAction {
	return model.CPSAction{
		ActionCode:       local_util.GenerateActionCode(),
		UniqueId:         unique,
		MakerID:          makerUser.UserID,
		MakerName:        makerUser.FullName,
		MakerPhoneNumber: makerUser.PhoneNumber,
		Department:       makerUser.Department,
		PreviousAction:   prevAction,
		CurrentAction:    currentAction,
		ActionStatus:     string(constants.Pending),
		ActionType:       actionType,
		RequestAction:    requestAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		MakerActionTime:  time.Now(),
	}
}

func GoRoutinBaker(opts types.BakerOptions, tasks ...func()) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	runtime.GOMAXPROCS(runtime.NumCPU())
	if opts.Sequential {
		for _, task := range tasks {
			if opts.UseMutex {
				mu.Lock()
				task()
				mu.Unlock()
			} else {
				task()
			}
		}
		return
	}

	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			if opts.UseMutex {
				mu.Lock()
				t()
				mu.Unlock()
			} else {
				t()
			}
		}(task)
	}
	wg.Wait()
}

func FilterBuilder(filterParam types.Filter, searchKeys bson.M, allowedKeys []string) (bson.M, int64, int64) {
	var skip, limit int64
	filter := bson.M{}

	if filterParam.Search != "" {
		maps.Copy(filter, searchKeys)
	}

	if filterParam.Filters != nil {

		handler := map[string]func(interface{}) interface{}{}
		includedKeys := []string{
			"enabled",
			"enable",
			"is_enabled",
			"is_deleted",
			"is_blocked",
			"ussd_enabled",
			"is_account_active",
			"is_main",
			"last_linked_status",
			"is_verified",
			"is_blocked",
			"active_account",
			"account_frozen",
			"account_dormant",
			"debit_allowed",
			"credit_allowed",
			"has_restriction",
		}
		for _, key := range includedKeys {
			for _, allowedKey := range allowedKeys {
				if allowedKey == key {
					handler[key] = func(value interface{}) interface{} {
						if str, ok := value.(string); ok {
							if parsed, err := strconv.ParseBool(str); err == nil {
								return parsed
							}
						}
						return value
					}
				}
			}
		}
		enhancedFilter := local_util.BuildMongoFilterWithKeys(filterParam.Filters, allowedKeys, handler)

		maps.Copy(filter, enhancedFilter)
	}

	skip = int64((filterParam.Page - 1) * filterParam.PerPage)
	limit = int64(filterParam.PerPage)

	return filter, skip, limit
}

func UploadFileToMinio(
	ctx context.Context,
	s3Client *s3.Client,
	bucketName string,
	fileHeader *multipart.FileHeader,
	prefix string,
	env config.VaultConfig,
	objectkey string,
	logger interface {
		Errorf(format string, args ...any)
	},
) (string, error) {

	// Open file and buffer its content
	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorf("failed to open file: %v", err)
		return "", errors.New(err.Error())
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, file)
	if err != nil {
		logger.Errorf("failed to read uploaded file error: %v", err)
		return "", errors.New(err.Error())
	}

	// Build object key under a folder (bucketName used as folder/prefix)
	// and generate unique filename based on prefix and timestamp to avoid collisions
	genName := fmt.Sprintf("%s-%d-%s", prefix, time.Now().UnixNano(), fileHeader.Filename)
	key := genName
	// key := path.Join(bucketName, genName)

	// Determine content type
	contentType := fileHeader.Header.Get("Content-Type")
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	// Upload using the buffered bytes to avoid EOF issues
	cl := n
	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName), // secrets.AWS_BUCKET_NAME
		Key:           aws.String(key),
		Body:          bytes.NewReader(buf.Bytes()),
		ContentType:   aws.String(contentType),
		ContentLength: &cl,
	}

	if _, err := s3Client.PutObject(context.TODO(), putInput); err != nil {
		logger.Errorf("upload failed error: %v", err)
		return "", errors.New(err.Error())
	}

	// Build streamed URL served by the uploader service
	url := fmt.Sprintf("%s/%s", env.MinioPublicEndPoint, strings.TrimPrefix(key, "/"))
	return url, nil

}

func RemoveFileFromMino(ctx context.Context, client aws.Config, bucketName string, objectkey string,
	logger interface {
		Errorf(format string, args ...any)
	},
) error {
	// exist, err := client.BucketExist(ctx, bucketName)
	// if err != nil {
	// 	logger.Errorf("failed to check bucket '%s': '%v'", bucketName, err)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// if !exist {
	// 	logger.Errorf("bucket '%s' does not exist", bucketName)
	// 	return errors.New(localization.ErrorBucketNotFound.Code)
	// }

	// isDeleted, err := client.DeleteObject(ctx, config.DeleteObjectBody{
	// 	BucketName: bucketName,
	// 	ObjectName: objectkey,
	// })
	// if err != nil {
	// 	logger.Errorf("failed to delete object '%s' from bucket '%s': %v", objectkey, bucketName, err)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// if !isDeleted {
	// 	logger.Errorf("object '%s' could not be deleted from bucket '%s'", objectkey, bucketName)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	return nil
}
