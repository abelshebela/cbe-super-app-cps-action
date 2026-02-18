package lib

import (
	"bytes"
	"cbe-super-app-cps-action/internal/constants"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	// "cbe-super-app-cps-action/internal/constants/localization"
	erp_merchant_update_dto "cbe-super-app-cps-action/internal/constants/dto/erp_merchant_update"
	"cbe-super-app-cps-action/internal/constants/localization"
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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UploadVideoToMinio(
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

	// 1. Open the file
	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorf("failed to open file: %v", err)
		return "", err
	}
	defer file.Close()

	// 2. Create uploader
	uploader := manager.NewUploader(s3Client)

	// 3. Metadata Preparation
	genName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)
	key := genName

	contentType := fileHeader.Header.Get("Content-Type")
	if strings.TrimSpace(contentType) == "" {
		contentType = "video/mp4"
	}

	// 4. Upload using the manager for multipart
	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	if _, err := uploader.Upload(ctx, putInput); err != nil {
		logger.Errorf("failed to upload video %s: %v", key, err)
		return "", err
	}

	baseURL := env.MinioPublicEndPoint
	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}

	// 5. Build public URL
	url := fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(key, "/"))
	return url, nil
}

func CpsModelBuilder(unique string, makerUser types.UserContext, prevAction, currentAction any, requestAction, actionType string) model.CPSAction {
	return model.CPSAction{
		ActionCode:       local_util.GenerateActionCode(),
		UniqueId:         unique,
		MakerID:          makerUser.UserName,
		MakerName:        makerUser.FullName,
		MakerPhoneNumber: makerUser.PhoneNumber,
		PreviousAction:   prevAction,
		AuditorStatus:    model.AuditorStatus(constants.AUDITORNOTCHECKED),
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

		// --- date range filters ---
		// For each allowed key, check if _from / _to variants exist in the filters.
		allowedSet := make(map[string]bool, len(allowedKeys))
		for _, k := range allowedKeys {
			allowedSet[k] = true
		}

		for _, ak := range allowedKeys {
			fromKey := ak + "_from"
			toKey := ak + "_to"
			dateFilter := bson.M{}

			// exact date match: ?created_at=2026-01-05 → full day range
			if raw, ok := filterParam.Filters[ak]; ok {
				if str, ok := raw.(string); ok && str != "" {
					if t, err := parseDateInput(str); err == nil {
						if !strings.Contains(str, "T") {
							// date-only → match the whole day
							dateFilter["$gte"] = t
							dateFilter["$lte"] = t.Add(24*time.Hour - time.Millisecond)
						} else {
							// exact datetime
							dateFilter["$eq"] = t
						}
						delete(filterParam.Filters, ak)
					}
				}
			}

			// range: ?created_at_from=...&created_at_to=...
			if raw, ok := filterParam.Filters[fromKey]; ok {
				if str, ok := raw.(string); ok && str != "" {
					if t, err := parseDateInput(str); err == nil {
						dateFilter["$gte"] = t
					}
				}
				delete(filterParam.Filters, fromKey)
			}

			if raw, ok := filterParam.Filters[toKey]; ok {
				if str, ok := raw.(string); ok && str != "" {
					if t, err := parseDateInput(str); err == nil {
						// if date-only (no time component), set to end of day
						if !strings.Contains(str, "T") {
							t = t.Add(24*time.Hour - time.Millisecond)
						}
						dateFilter["$lte"] = t
					}
				}
				delete(filterParam.Filters, toKey)
			}

			if len(dateFilter) > 0 {
				filter[ak] = dateFilter
			}
		}

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
			"advert_for",
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

// parseDateInput parses a date string that can be either date-only ("2026-01-05")
// or full ISO datetime ("2026-01-05T07:10:33.695+00:00", "2026-01-05T07:10:33").
func parseDateInput(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,                    // 2026-01-05T07:10:33+00:00
		"2006-01-02T15:04:05.000Z07:00", // 2026-01-05T07:10:33.695+00:00
		"2006-01-02T15:04:05.999Z07:00", // milliseconds variant
		"2006-01-02T15:04:05Z07:00",     // without millis
		"2006-01-02T15:04:05",           // no timezone
		"2006-01-02",                    // date only
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// func UploadFileToMinio(
// 	ctx context.Context,
// 	s3Client *s3.Client,
// 	bucketName string,
// 	fileHeader *multipart.FileHeader,
// 	prefix string,
// 	env config.VaultConfig,
// 	objectkey string,
// 	logger interface {
// 		Errorf(format string, args ...any)
// 	},
// ) (string, error) {

// 	// Open file and buffer its content
// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		logger.Errorf("failed to open file: %v", err)
// 		return "", errors.New(err.Error())
// 	}
// 	defer file.Close()

// 	buf := new(bytes.Buffer)
// 	n, err := io.Copy(buf, file)
// 	if err != nil {
// 		logger.Errorf("failed to read uploaded file error: %v", err)
// 		return "", errors.New(err.Error())
// 	}

// 	// Build object key under a folder (bucketName used as folder/prefix)
// 	// and generate unique filename based on prefix and timestamp to avoid collisions
// 	genName := fmt.Sprintf("%s-%d-%s", prefix, time.Now().UnixNano(), fileHeader.Filename)
// 	key := genName
// 	// key := path.Join(bucketName, genName)

// 	// Determine content type
// 	contentType := fileHeader.Header.Get("Content-Type")
// 	if strings.TrimSpace(contentType) == "" {
// 		contentType = "application/octet-stream"
// 	}

// 	// Upload using the buffered bytes to avoid EOF issues
// 	cl := n
// 	putInput := &s3.PutObjectInput{
// 		Bucket:        aws.String(bucketName), // secrets.AWS_BUCKET_NAME
// 		Key:           aws.String(key),
// 		Body:          bytes.NewReader(buf.Bytes()),
// 		ContentType:   aws.String(contentType),
// 		ContentLength: &cl,
// 	}

// 	if _, err := s3Client.PutObject(context.TODO(), putInput); err != nil {
// 		logger.Errorf("upload failed error: %v", err)
// 		return "", errors.New(err.Error())
// 	}

// 	// Build streamed URL served by the uploader service
// 	url := fmt.Sprintf("%s/%s", env.MinioPublicEndPoint, strings.TrimPrefix(key, "/"))
// 	return url, nil
// }

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

	// 1. Open the file
	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorf("failed to open file: %v", err)
		return "", err
	}
	defer file.Close()

	// 2. Read the original bytes into memory first.
	// This acts as our safety "fallback" buffer.
	originalBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Errorf("failed to read file: %v", err)
		return "", err
	}

	// Default values for the upload
	finalBytes := originalBytes
	contentType := fileHeader.Header.Get("Content-Type")

	// 3. Attempt Compression
	// We use bytes.NewReader so we don't exhaust the original stream
	img, format, decodeErr := image.Decode(bytes.NewReader(originalBytes))

	if decodeErr == nil {
		// If decoding succeeded, we try to encode with compression
		buf := new(bytes.Buffer)
		var encodeErr error

		switch format {
		case "jpeg":
			encodeErr = jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})
		case "png":
			enc := png.Encoder{CompressionLevel: png.BestCompression}
			encodeErr = enc.Encode(buf, img)
		case "gif":

		default:
			// No specific compressor for this format?
			// We do nothing and keep originalBytes
			encodeErr = errors.New("no specific encoder")
		}

		// Only swap to compressed bytes if the process actually worked
		if encodeErr == nil {
			finalBytes = buf.Bytes()
		}
	}

	// 4. Metadata Preparation
	// genName := fmt.Sprintf("%s/%d-%s", prefix, time.Now().UnixNano(), fileHeader.Filename)
	genName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)
	key := genName

	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	// 5. Upload the resulting bytes (compressed or original)
	cl := int64(len(finalBytes))
	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          bytes.NewReader(finalBytes),
		ContentType:   aws.String(contentType),
		ContentLength: &cl,
	}
	if _, err := s3Client.PutObject(ctx, putInput); err != nil {
		logger.Errorf("upload failed error: %v", err)
		return "", err
	}
	baseURL := env.MinioPublicEndPoint
	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}
	// Build public URL
	url := fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(key, "/"))
	return url, nil
}
func RemoveFileFromMinio(
	ctx context.Context,
	client *s3.Client,
	bucketName string,
	objectKey string,
	logger interface {
		Errorf(format string, args ...any)
	}) error {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		logger.Errorf("bucket '%s' does not exist", bucketName)
		return errors.New(localization.ErrorBucketNotFound.Code)
	}

	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		logger.Errorf("failed to delete object '%s' from bucket '%s': %v", objectKey, bucketName, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

// fallbackModuleForRA tries to infer a module name from the action string when
// it isn't mapped in RequestActionGroups. Best-effort, case-insensitive.
func FallbackModuleForRA(action constants.RequestAction) (string, bool) {
	s := strings.ToUpper(string(action))
	switch {
	case strings.Contains(s, "WALLET"):
		return "Wallet", true
	case strings.Contains(s, "TOPUP"):
		return "Topup", true
	case strings.Contains(s, "BANK_VAULT"):
		return "BankVault", true
	case strings.Contains(s, "BANK"):
		return "Bank", true
	case strings.Contains(s, "KYC"):
		return "KYCVerifier", true
	case strings.Contains(s, "FAYDA"):
		return "Fayda", true
	case strings.Contains(s, "MINI_APP_MERCHANT"):
		return "MiniAppMerchant", true
	case strings.Contains(s, "MINI_APP_CATEGORY"):
		return "MiniAppCategory", true
	case strings.Contains(s, "MINI_APP"):
		return "MiniApp", true
	case strings.Contains(s, "DEVICE_VERSION"):
		return "DeviceVersion", true
	case strings.Contains(s, "PERMISSION"):
		return "Permission", true
	case strings.Contains(s, "PASSWORD"):
		return "Password", true
	case strings.Contains(s, "AMOUNT_BASED_AUTH") || strings.Contains(s, "AUTHTIER"):
		return "AmountBasedAuth", true
	case strings.Contains(s, "NOTIFICATION"):
		return "Notification", true
	case strings.Contains(s, "AVATAR"):
		return "Avatar", true
	case strings.Contains(s, "DONATION_CATEGORY"):
		return "DonationCategory", true
	case strings.Contains(s, "DONATION_COMPANY"):
		return "DonationCompany", true
	case strings.Contains(s, "DONATION"):
		return "Donation", true
	case strings.Contains(s, "DEPARTMENT"):
		return "Department", true
	case strings.Contains(s, "CPS_USER"):
		return "CPSUser", true
	case strings.Contains(s, "VAULT_GROUP_CATEGORY"):
		return "VaultGroupCategory", true
	case strings.Contains(s, "ARTICLE_CATEGORY"):
		return "ArticleCategory", true
	case strings.Contains(s, "ARTICLE"):
		return "Article", true
	case strings.Contains(s, "SHORT_VIDEO"):
		return "ShortVideo", true
	case strings.Contains(s, "CUSTOMER"):
		return "Customer", true
	case strings.Contains(s, "NEWS_TAG"):
		return "NewsTag", true
	case strings.Contains(s, "NEWS_CATEGORY"):
		return "NewsCategory", true
	case strings.Contains(s, "BUDGET_CATEGORY"):
		return "BudgetCategory", true
	case strings.Contains(s, "SERVICE_FEE") || strings.Contains(s, "DAILY_LIMIT") || strings.Contains(s, "MINIMUM") || strings.Contains(s, "TOTAL") || strings.Contains(s, "ACCESS_CONFIG"):
		return "Service", true
	case strings.Contains(s, "SERVICE"):
		return "ServicesCatalog", true
	case strings.Contains(s, "PRODUCT_CODE"):
		return "ProductCode", true
	case strings.Contains(s, "EVENT"):
		return "Event", true
	case strings.Contains(s, "BULK_SERVICE"):
		return "BulkService", true
	case strings.Contains(s, "UNLINK"):
		return "UnlinkDevice", true
	case strings.Contains(s, "ACTION_ROLE"):
		return "ActionRole", true
	case strings.Contains(s, "CPS_ACTION_ROLE"):
		return "CpsActionRole", true
	}
	return "", false
}

func PublishMerchantChangeToERP(ctx context.Context, cfg *config.VaultConfig, body erp_merchant_update_dto.ERPUpdateRequest, merchantID string, logger utils.Logger) error {
	logger.Infof("Publishing merchant change to ERP for merchant %s with body %+v", merchantID, body)
	ctx, span := local_util.TraceLogger(ctx, "core", "UpdateERP", "LogisticsMerchant", "UpdateERP")
	defer span.End()

	base := "https://qaapisuperapp.cbe.com.et/api/v1/cbesuperapp/ecommerce"
	if cfg == nil {
		logger.Debugf("env config is nil, using hardcoded base url and cannot proceed without api key")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if cfg.OddoEcommerceBaseUrl != "" {
		base = cfg.OddoEcommerceBaseUrl
	} else {
		logger.Debugf("env url for publish not found using hardcoded")
	}
	base += "/cps/merchant/update/" + merchantID
	if cfg.ApiKey == "" {
		logger.Debugf("env api key for publish not found using hardcoded")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	apiKey := cfg.ApiKey
	jsonBody, err := json.Marshal(body)
	if err != nil {
		logger.Errorf("Failed to marshal ERP update body: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, base, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Errorf("Failed to build ERP update request: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-api-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("ERP update request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Errorf("ERP update failed body: %s", string(bodyBytes))
		logger.Errorf("ERP update failed request header: %s", req.Header)
		logger.Errorf("ERP update failed response header: %s", resp.Header)
		logger.Errorf("ERP update failed json body: %s", jsonBody)
		logger.Errorf("ERP UPDATE Used URL %s", base)
		logger.Errorf("ERP update failed api key: %s", cfg.ApiKey)
		return errors.New("ERP update failed")
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	logger.Infof("ERP update successful for merchant %s with response status %d, response body: %s", merchantID, resp.StatusCode, string(bodyBytes))
	return nil
}
