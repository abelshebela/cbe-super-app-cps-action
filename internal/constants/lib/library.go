package lib

import (
	"bytes"
	"cbe-super-app-cps-action/internal/constants"
	"encoding/csv"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

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

// func FilterBuilder(filterParam types.Filter, searchKeys bson.M, allowedKeys []string) (bson.M, int64, int64) {
// 	var skip, limit int64
// 	filter := bson.M{}

// 	if filterParam.Search != "" {
// 		maps.Copy(filter, searchKeys)
// 	}

// 	if filterParam.Filters != nil {

// 		// --- date range filters ---
// 		// For each allowed key, check if _from / _to variants exist in the filters.
// 		allowedSet := make(map[string]bool, len(allowedKeys))
// 		for _, k := range allowedKeys {
// 			allowedSet[k] = true
// 		}

// 		for _, ak := range allowedKeys {
// 			fromKey := ak + "_from"
// 			toKey := ak + "_to"
// 			dateFilter := bson.M{}

// 			// exact date match: ?created_at=2026-01-05 → full day range
// 			if raw, ok := filterParam.Filters[ak]; ok {
// 				if str, ok := raw.(string); ok && str != "" {
// 					if t, err := parseDateInput(str); err == nil {
// 						if !strings.Contains(str, "T") {
// 							// date-only → match the whole day
// 							dateFilter["$gte"] = t
// 							dateFilter["$lte"] = t.Add(24*time.Hour - time.Millisecond)
// 						} else {
// 							// exact datetime
// 							dateFilter["$eq"] = t
// 						}
// 						delete(filterParam.Filters, ak)
// 					}
// 				}
// 			}

// 			// range: ?created_at_from=...&created_at_to=...
// 			if raw, ok := filterParam.Filters[fromKey]; ok {
// 				if str, ok := raw.(string); ok && str != "" {
// 					if t, err := parseDateInput(str); err == nil {
// 						dateFilter["$gte"] = t
// 					}
// 				}
// 				delete(filterParam.Filters, fromKey)
// 			}

// 			if raw, ok := filterParam.Filters[toKey]; ok {
// 				if str, ok := raw.(string); ok && str != "" {
// 					if t, err := parseDateInput(str); err == nil {
// 						// if date-only (no time component), set to end of day
// 						if !strings.Contains(str, "T") {
// 							t = t.Add(24*time.Hour - time.Millisecond)
// 						}
// 						dateFilter["$lte"] = t
// 					}
// 				}
// 				delete(filterParam.Filters, toKey)
// 			}

// 			if len(dateFilter) > 0 {
// 				filter[ak] = dateFilter
// 			}
// 		}

// 		handler := map[string]func(interface{}) interface{}{}
// 		includedKeys := []string{
// 			"enabled",
// 			"enable",
// 			"is_enabled",
// 			"is_deleted",
// 			"is_blocked",
// 			"ussd_enabled",
// 			"is_account_active",
// 			"is_main",
// 			"last_linked_status",
// 			"is_verified",
// 			"is_blocked",
// 			"active_account",
// 			"account_frozen",
// 			"account_dormant",
// 			"debit_allowed",
// 			"credit_allowed",
// 			"has_restriction",
// 			"advert_for",
// 		}
// 		for _, key := range includedKeys {
// 			for _, allowedKey := range allowedKeys {
// 				if allowedKey == key {
// 					handler[key] = func(value interface{}) interface{} {
// 						if str, ok := value.(string); ok {
// 							if parsed, err := strconv.ParseBool(str); err == nil {
// 								return parsed
// 							}
// 						}
// 						return value
// 					}
// 				}
// 			}
// 		}
// 		enhancedFilter := local_util.BuildMongoFilterWithKeys(filterParam.Filters, allowedKeys, handler)

// 		maps.Copy(filter, enhancedFilter)
// 	}

// 	skip = int64((filterParam.Page - 1) * filterParam.PerPage)
// 	limit = int64(filterParam.PerPage)

//		return filter, skip, limit
//	}
func FilterBuilder(filterParam types.Filter, searchKeys bson.M, allowedKeys []string) (bson.M, int64, int64) {
	var skip, limit int64
	filter := bson.M{}

	if filterParam.Search != "" {
		maps.Copy(filter, searchKeys)
	}

	if filterParam.Filters != nil {

		// --- NEW: CPS Action Date Range Filter ---
		var startDate, endDate time.Time
		var hasStart, hasEnd bool

		if raw, ok := filterParam.Filters["created_at_from"]; ok {
			if str, ok := raw.(string); ok && str != "" {
				if t, err := parseDateInput(str); err == nil {
					startDate = t
					hasStart = true
				}
			}
			// delete(filterParam.Filters, "created_at_from")
		}

		if raw, ok := filterParam.Filters["created_at_to"]; ok {
			if str, ok := raw.(string); ok && str != "" {
				if t, err := parseDateInput(str); err == nil {
					// include full day if date-only
					if !strings.Contains(str, "T") {
						t = t.Add(24*time.Hour - time.Millisecond)
					}
					endDate = t
					hasEnd = true
				}
			}
			// delete(filterParam.Filters, "created_at_to")
		}

		if hasStart && hasEnd {
			dateRangeFilter := BuildCPSActionDateRangeFilter(&filterParam, startDate, endDate)

			// merge with existing filter using $and
			if len(filter) > 0 {
				filter = bson.M{
					"$and": []bson.M{
						filter,
						dateRangeFilter,
					},
				}
			} else {
				filter = dateRangeFilter
			}
		}

		// --- EXISTING DATE FILTER LOGIC (per-field) ---
		allowedSet := make(map[string]bool, len(allowedKeys))
		for _, k := range allowedKeys {
			allowedSet[k] = true
		}

		for _, ak := range allowedKeys {
			fromKey := ak + "_from"
			toKey := ak + "_to"
			dateFilter := bson.M{}

			if raw, ok := filterParam.Filters[ak]; ok {
				if str, ok := raw.(string); ok && str != "" {
					if t, err := parseDateInput(str); err == nil {
						if !strings.Contains(str, "T") {
							dateFilter["$gte"] = t
							dateFilter["$lte"] = t.Add(24*time.Hour - time.Millisecond)
						} else {
							dateFilter["$eq"] = t
						}
						delete(filterParam.Filters, ak)
					}
				}
			}

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

		// --- BOOLEAN HANDLER ---
		handler := map[string]func(interface{}) interface{}{}
		includedKeys := []string{
			"enabled", "enable", "is_enabled", "is_deleted", "is_blocked",
			"ussd_enabled", "is_account_active", "is_main", "last_linked_status",
			"is_verified", "active_account", "account_frozen", "account_dormant",
			"debit_allowed", "credit_allowed", "has_restriction", "advert_for",
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
func BuildOracleFilter(
	filterParam types.Filter,
	searchKeys map[string]string, // keep as map
	allowedKeys []string, // keep as slice
) (string, []interface{}, int64, int64) {

	var filters []string
	var args []interface{}
	idx := 1

	filters = append(filters, "1=1")

	// --- SEARCH ---
	if filterParam.Search != "" {
		search := "%" + strings.ToUpper(filterParam.Search) + "%"
		searchParts := []string{}

		for _, column := range searchKeys { // keep map structure
			searchParts = append(searchParts,
				fmt.Sprintf("UPPER(%s) LIKE :%d", column, idx))
			args = append(args, search)
			idx++
		}

		if len(searchParts) > 0 {
			filters = append(filters, "("+strings.Join(searchParts, " OR ")+")")
		}
	}

	// --- FILTERS ---
	if filterParam.Filters != nil {
		allowedSet := make(map[string]bool)
		for _, k := range allowedKeys { // keep slice as input
			allowedSet[k] = true
		}

		boolKeys := map[string]bool{
			"enabled": true, "enable": true, "is_enabled": true,
			"is_deleted": true, "is_blocked": true,
			"is_verified": true, "active_account": true,
		}

		for key, val := range filterParam.Filters {
			if !allowedSet[key] {
				continue
			}

			// --- DATE RANGE ---
			if strings.HasSuffix(key, "_from") {
				column := strings.TrimSuffix(key, "_from")
				if str, ok := val.(string); ok && str != "" {
					if t, err := parseDateInput(str); err == nil {
						filters = append(filters,
							fmt.Sprintf("%s >= :%d", column, idx))
						args = append(args, t)
						idx++
					}
				}
				continue
			}

			if strings.HasSuffix(key, "_to") {
				column := strings.TrimSuffix(key, "_to")
				if str, ok := val.(string); ok && str != "" {
					if t, err := parseDateInput(str); err == nil {
						if !strings.Contains(str, "T") {
							t = t.Add(24*time.Hour - time.Millisecond)
						}
						filters = append(filters,
							fmt.Sprintf("%s <= :%d", column, idx))
						args = append(args, t)
						idx++
					}
				}
				continue
			}

			// --- BOOLEAN ---
			if boolKeys[key] {
				if str, ok := val.(string); ok {
					if parsed, err := strconv.ParseBool(str); err == nil {
						filters = append(filters,
							fmt.Sprintf("%s = :%d", key, idx))
						args = append(args, parsed)
						idx++
						continue
					}
				}
			}

			// --- DEFAULT ---
			filters = append(filters,
				fmt.Sprintf("%s = :%d", key, idx))
			args = append(args, val)
			idx++
		}
	}

	// --- PAGINATION ---
	offset := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	return strings.Join(filters, " AND "), args, offset, limit
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

func BuildCPSActionDateRangeFilter(filterMap *types.Filter, startDate, endDate time.Time) bson.M {
	rangeFilter := bson.M{
		"$gte": startDate,
		"$lte": endDate,
	}

	// Some records use different timestamp fields; include all known variants.
	return bson.M{
		"$or": []bson.M{
			{"created_at": rangeFilter},
			// {"action_created_at": rangeFilter},
			// {"maker_action_time": rangeFilter},
			// {"last_modified_at": rangeFilter},
		},
	}
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

// UploadCSVToMinio uploads a CSV file (from an io.Reader) to MinIO and returns
func UploadCSVToMinio(
	ctx context.Context,
	s3Client *s3.Client,
	bucketName string,
	body io.Reader,
	contentLength int64,
	env config.VaultConfig,
	objectKey string,
	logger interface {
		Errorf(format string, args ...any)
	},
) (string, error) {

	contentType := "text/csv; charset=utf-8"

	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: &contentLength,
	}
	if _, err := s3Client.PutObject(ctx, putInput); err != nil {
		logger.Errorf("upload CSV failed error: %v", err)
		return "", err
	}

	baseURL := strings.TrimSuffix(env.MinioPublicEndPoint, "/")
	url := fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(objectKey, "/"))
	return url, nil
}

// ExportCSVAndUpload is a shared utility that handles the full CSV-export-to-MinIO pipeline:
//  1. Creates a temporary CSV file
//  2. Writes the provided CSV header
//  3. Invokes the writeRows callback to stream data rows into the CSV writer
//  4. Flushes the CSV writer
//  5. Uploads the resulting file to MinIO and returns the public URL
//
// The writeRows callback receives a *csv.Writer and is responsible for writing
// all data rows (e.g., by streaming from a repository).
func ExportCSVAndUpload(
	ctx context.Context,
	s3Client *s3.Client,
	bucketName string,
	env config.VaultConfig,
	objectKey string,
	headers []string,
	writeRows func(writer *csv.Writer) error,
	logger interface {
		Errorf(format string, args ...any)
	},
) (string, error) {

	// 1. Create temp CSV file
	tmpFile, err := os.CreateTemp("", "export_*.csv")
	if err != nil {
		logger.Errorf("[ExportCSVAndUpload] create temp file: %v", err)
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)

	// 2. Write CSV header
	if err := writer.Write(headers); err != nil {
		logger.Errorf("[ExportCSVAndUpload] write header: %v", err)
		return "", fmt.Errorf("write header: %w", err)
	}

	// 3. Stream rows via callback
	if err := writeRows(writer); err != nil {
		logger.Errorf("[ExportCSVAndUpload] write rows: %v", err)
		return "", fmt.Errorf("stream data: %w", err)
	}

	// 4. Flush CSV writer
	writer.Flush()
	if err := writer.Error(); err != nil {
		logger.Errorf("[ExportCSVAndUpload] flush csv: %v", err)
		return "", fmt.Errorf("flush csv: %w", err)
	}

	// 5. Seek to beginning and get file size
	if _, err := tmpFile.Seek(0, 0); err != nil {
		logger.Errorf("[ExportCSVAndUpload] seek temp file: %v", err)
		return "", fmt.Errorf("seek temp file: %w", err)
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		logger.Errorf("[ExportCSVAndUpload] stat temp file: %v", err)
		return "", fmt.Errorf("stat temp file: %w", err)
	}

	// 6. Upload to MinIO
	url, err := UploadCSVToMinio(ctx, s3Client, bucketName, tmpFile, stat.Size(), env, objectKey, logger)
	if err != nil {
		return "", fmt.Errorf("upload to minio: %w", err)
	}

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

func PublishMerchantChangeToERP(ctx context.Context, cfg *config.VaultConfig, body erp_merchant_update_dto.ERPUpdateRequest, merchantID string, isEventMerchant bool, logger utils.Logger) error {
	logger.Infof("Publishing merchant change to ERP for merchant %s with body %+v", merchantID, body)
	ctx, span := local_util.TraceLogger(ctx, "core", "UpdateERP", "LogisticsMerchant", "UpdateERP")
	defer span.End()

	base := ""
	if cfg == nil {
		logger.Debugf("env config is nil, using hardcoded base url and cannot proceed without api key")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if cfg.OddoEcommerceBaseUrl != "" {
		base = cfg.OddoEcommerceBaseUrl
	} else {
		logger.Debugf("env url for ecommerce merchant publish not found using hardcoded")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if isEventMerchant {
		base += "/cps/event/merchant/update/" + merchantID
	} else {
		base += "/cps/merchant/update/" + merchantID
	}

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
