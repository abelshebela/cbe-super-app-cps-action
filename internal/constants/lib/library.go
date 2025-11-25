package lib

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"maps"
	"mime/multipart"
	"runtime"
	"strconv"
	"sync"
	"time"

	config "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

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
	uploader config.MinioClientInterface,
	bucketName string,
	fileHeader *multipart.FileHeader,
	prefix string,
	minioEndpoint string,
	objectkey string,
	logger interface {
		Errorf(format string, args ...any)
	},
) (string, error) {
	// Ensure bucket exists
	exist, err := uploader.BucketExist(ctx, bucketName)
	if err != nil {
		fmt.Println("=====fileName=====", bucketName, err)
		logger.Errorf("failed to check bucket '%s': %v", bucketName, err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	if !exist {
		created, err := uploader.MakeBucket(ctx, bucketName)
		if err != nil || !created {
			logger.Errorf("failed to create bucket '%s': %v", bucketName, err)
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorf("failed to open file: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer file.Close()

	var fileName string
	if objectkey != "" {
		fileName = objectkey
	} else {
		extension := fileHeader.Filename[len(fileHeader.Filename)-4:]
		fileName = fmt.Sprintf("%s-%d.%s", prefix, time.Now().UnixNano(), extension)
	}

	// Upload file
	saveObj, err := uploader.SaveObjectN(ctx, config.SaveObjectBodyN{
		BucketName:  bucketName,
		ObjectName:  fileName,
		Reader:      file,
		Size:        fileHeader.Size,
		ContentType: config.ContentType(fileHeader.Header.Get("Content-Type")),
	})

	if err != nil {
		logger.Errorf("failed to upload file to MinIO: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Return full URL
	url := fmt.Sprintf("%s/%s/%s", minioEndpoint, saveObj.Bucket, saveObj.Key)
	return url, nil
}

func RemoveFileFromMino(ctx context.Context, client config.MinioClientInterface, bucketName string, objectkey string,
	logger interface {
		Errorf(format string, args ...any)
	},
) error {
	exist, err := client.BucketExist(ctx, bucketName)
	if err != nil {
		logger.Errorf("failed to check bucket '%s': '%v'", bucketName, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if !exist {
		logger.Errorf("bucket '%s' does not exist", bucketName)
		return errors.New(localization.ErrorBucketNotFound.Code)
	}

	isDeleted, err := client.DeleteObject(ctx, config.DeleteObjectBody{
		BucketName: bucketName,
		ObjectName: objectkey,
	})
	if err != nil {
		logger.Errorf("failed to delete object '%s' from bucket '%s': %v", objectkey, bucketName, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if !isDeleted {
		logger.Errorf("object '%s' could not be deleted from bucket '%s'", objectkey, bucketName)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
