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
	"reflect"

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

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jung-kurt/gofpdf"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type FileType string

const (
	FileTypeCSV FileType = "csv"
	FileTypePDF FileType = "pdf"
)

// Generic upload function
type UploadFunc func(ctx context.Context, file *os.File, size int64, objectName string, fileType FileType) (string, error)

// Generic config
type FileProducerConfig struct {
	FilePrefix string
	Header     []string
	FileType   FileType
	ObjectName string
}

// extractFieldsFilter normalizes filterMap.Filters["fields"] into []string.
func extractFieldsFilter(filters map[string]interface{}) []string {
	if filters == nil {
		return nil
	}
	return local_util.StringSliceFromFilterValue(filters["fields"])
}

// PDFLayout bundles font + row sizing for a tabular PDF export. It is derived from the
// column count so wider tables automatically get a smaller body font + row height.
type PDFLayout struct {
	HeaderFontPt float64
	BodyFontPt   float64
	HeaderRowMM  float64
	BodyRowMM    float64
}

// CalcPDFLayout picks responsive font + row sizing based on the number of columns.
// Rationale: at 9pt Arial a character is ~2mm wide; with 14 columns in 277mm landscape
// each column is only ~19mm, so 8-character budget wasn't enough. Scale font instead.
func CalcPDFLayout(cols int) PDFLayout {
	switch {
	case cols <= 4:
		return PDFLayout{HeaderFontPt: 11, BodyFontPt: 10, HeaderRowMM: 9, BodyRowMM: 8}
	case cols <= 7:
		return PDFLayout{HeaderFontPt: 10, BodyFontPt: 9, HeaderRowMM: 8, BodyRowMM: 7}
	case cols <= 10:
		return PDFLayout{HeaderFontPt: 9, BodyFontPt: 8, HeaderRowMM: 7, BodyRowMM: 6}
	case cols <= 13:
		return PDFLayout{HeaderFontPt: 8, BodyFontPt: 7, HeaderRowMM: 6.5, BodyRowMM: 5.5}
	default:
		return PDFLayout{HeaderFontPt: 7, BodyFontPt: 6, HeaderRowMM: 6, BodyRowMM: 5}
	}
}

// pdfCharBudget estimates how many characters of `fontPt` Arial fit into colWidthMM
// after reserving 2mm of horizontal padding inside the cell.
func pdfCharBudget(colWidthMM, fontPt float64) int {
	// Arial regular at 9pt ≈ 2mm per char; width scales linearly with font size.
	charMM := fontPt * (2.0 / 9.0)
	if charMM <= 0 {
		return 4
	}
	budget := int((colWidthMM - 2.0) / charMM)
	if budget < 4 {
		budget = 4
	}
	return budget
}

// truncatePDFCell keeps cell text from overflowing a fixed-width column. ASCII ellipsis
// ("...") is used instead of the Unicode "…" because gofpdf's core Arial font uses
// CP-1252 and renders U+2026 as the garbled "â€¦" sequence.
func truncatePDFCell(s string, colWidthMM, fontPt float64) string {
	if s == "" {
		return s
	}
	budget := pdfCharBudget(colWidthMM, fontPt)
	if len(s) <= budget {
		return s
	}
	if budget <= 3 {
		return s[:budget]
	}
	return s[:budget-3] + "..."
}

// PDFExportOptions configures page size for tabular PDF exports. Orientation and column
// paging are chosen automatically from the column count.
type PDFExportOptions struct {
	PageSize string // A4 (default) or A5
}

const (
	pdfSideMarginMM      = 10.0
	pdfBannerHeightMM    = 22.0
	pdfFooterMarginMM    = 18.0
	pdfMinColWidthMM     = 14.0
	pdfMaxColsPerSection = 10
	pdfCellPadHMM        = 1.5 // horizontal inset inside each cell
	pdfCellPadVMM        = 1.5 // vertical inset (top and bottom)
)

// pdfPageSizeMM returns full page width and height in millimeters.
func pdfPageSizeMM(size, orientation string) (widthMM, heightMM float64) {
	size = strings.ToUpper(strings.TrimSpace(size))
	if size == "" {
		size = "A4"
	}
	switch size {
	case "A5":
		if orientation == "L" {
			return 210, 148
		}
		return 148, 210
	default:
		if orientation == "L" {
			return 297, 210
		}
		return 210, 297
	}
}

// pdfEffectiveLayout picks orientation and page size so every column gets at least
// pdfMinColWidthMM of width; wide tables use landscape A4 and may split across sections.
func pdfEffectiveLayout(requestedSize string, colCount int) (pageSize, orientation string, sections [][2]int) {
	requestedSize = strings.ToUpper(strings.TrimSpace(requestedSize))
	if requestedSize == "" {
		requestedSize = "A4"
	}
	pageSize = requestedSize
	orientation = "P"

	if colCount <= 0 {
		return pageSize, orientation, [][2]int{{0, 0}}
	}

	tryOrient := func(size, orient string) float64 {
		w, _ := pdfPageSizeMM(size, orient)
		return (w - 2*pdfSideMarginMM) / float64(colCount)
	}

	if tryOrient(pageSize, "P") < pdfMinColWidthMM || colCount > 7 {
		orientation = "L"
	}
	if pageSize == "A5" && (orientation == "L" || colCount > 6) {
		pageSize = "A4"
	}
	if tryOrient(pageSize, orientation) < pdfMinColWidthMM && pageSize != "A4" {
		pageSize = "A4"
		orientation = "L"
	}

	perSection := colCount
	if orientation == "P" && colCount > pdfMaxColsPerSection {
		perSection = pdfMaxColsPerSection
	} else if orientation == "L" && colCount > pdfMaxColsPerSection+2 {
		perSection = pdfMaxColsPerSection + 2
	}
	if perSection < 1 {
		perSection = 1
	}

	for start := 0; start < colCount; start += perSection {
		end := start + perSection
		if end > colCount {
			end = colCount
		}
		sections = append(sections, [2]int{start, end})
	}
	return pageSize, orientation, sections
}

func pdfLineHeightMM(fontPt float64) float64 {
	if fontPt <= 0 {
		fontPt = 8
	}
	// Convert pt → mm and add ~15% leading so wrapped lines don't collide.
	return fontPt * 25.4 / 72.0 * 1.15
}

// pdfWrapFriendlyText inserts break hints in long unbroken tokens (emails, IDs, codes)
// so SplitText wraps on punctuation instead of mid-word.
func pdfWrapFriendlyText(s string) string {
	if s == "" || strings.Contains(s, " ") || len(s) <= 14 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	for i, r := range s {
		b.WriteRune(r)
		if i < len(s)-1 {
			switch r {
			case '@', '.', '_', '-', ':', '/', ',':
				b.WriteByte(' ')
			}
		}
	}
	return b.String()
}

func pdfSplitCellLines(pdf *gofpdf.Fpdf, text string, innerWidthMM float64) []string {
	text = pdfWrapFriendlyText(text)
	if text == "" {
		return []string{""}
	}
	lines := pdf.SplitText(text, innerWidthMM)
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func pdfRowLineCount(pdf *gofpdf.Fpdf, cells []string, colWidthMM float64) int {
	inner := colWidthMM - 2*pdfCellPadHMM
	if inner < 1 {
		inner = 1
	}
	maxLines := 1
	for _, cell := range cells {
		if n := len(pdfSplitCellLines(pdf, cell, inner)); n > maxLines {
			maxLines = n
		}
	}
	return maxLines
}

func pdfCalcRowHeight(pdf *gofpdf.Fpdf, cells []string, colWidthMM, lineHt, minRowHt float64) float64 {
	lineCount := pdfRowLineCount(pdf, cells, colWidthMM)
	rowHt := 2*pdfCellPadVMM + float64(lineCount)*lineHt
	if rowHt < minRowHt {
		rowHt = minRowHt
	}
	return rowHt
}

// pdfDrawTableRow renders a bordered row with wrapped text; every cell in the row shares
// the same height so columns stay aligned.
func pdfDrawTableRow(pdf *gofpdf.Fpdf, x, y float64, cells []string, colWidthMM, lineHt, minRowHt float64, header bool) float64 {
	innerW := colWidthMM - 2*pdfCellPadHMM
	if innerW < 1 {
		innerW = 1
	}
	rowHt := pdfCalcRowHeight(pdf, cells, colWidthMM, lineHt, minRowHt)

	for i, cell := range cells {
		xi := x + float64(i)*colWidthMM
		if header {
			pdf.SetFillColor(220, 220, 220)
			pdf.Rect(xi, y, colWidthMM, rowHt, "FD")
		} else {
			pdf.Rect(xi, y, colWidthMM, rowHt, "D")
		}
		lines := pdfSplitCellLines(pdf, cell, innerW)
		for li, line := range lines {
			pdf.SetXY(xi+pdfCellPadHMM, y+pdfCellPadVMM+float64(li)*lineHt)
			pdf.Cell(innerW, lineHt, line)
		}
	}
	pdf.SetXY(x, y+rowHt)
	return rowHt
}

// FileExporterForCPSAction streams CPS actions into CSV or PDF, honoring ?fields=... projection.
//
// The CpsActionCSVHeader callback is kept only for backwards compatibility and is NOT used
// when a field is registered in CPSActionFieldRegistry — in that case headers + rows are
// resolved from the registry so column counts always match.
func FileExporterForCPSAction(ctx context.Context, cfg config.VaultConfig, minioClient *s3.Client, buckerName string, filterMap *types.Filter, data []*model.CPSAction, _ func(fields []string) []string, logger utils.Logger) (string, error) {
	if filterMap == nil {
		return "", errors.New(localization.ErrorRequiredFieldMissing.Code)
	}
	if filterMap.Filters == nil {
		filterMap.Filters = map[string]interface{}{}
	}

	// 1. Date range validation.
	createdAtFrom, fromOk := filterMap.Filters["created_at_from"].(string)
	createdAtTo, toOk := filterMap.Filters["created_at_to"].(string)
	if !fromOk || !toOk || createdAtFrom == "" || createdAtTo == "" {
		return "", errors.New(localization.ErrorRequiredFieldMissing.Code)
	}

	startDate, endDate, err := local_util.FormatDateRangeToUTCStrings(createdAtFrom, createdAtTo)
	if err != nil {
		return "", errors.New(localization.ErrorInvalidDateFormat.Code)
	}
	filterMap.Filters["created_at_from"] = startDate
	filterMap.Filters["created_at_to"] = endDate

	if len(data) == 0 {
		return "", errors.New(localization.CpsActionDataNotFoundInDateRange.Code)
	}

	// 2. Resolve output format + field projection.
	fileType, _ := filterMap.Filters["file_type"].(string)
	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType == "" {
		fileType = string(FileTypeCSV)
	}

	requestedFields := extractFieldsFilter(filterMap.Filters)
	resolvedFields := ResolveCPSActionFields(requestedFields)
	if len(requestedFields) > 0 && len(resolvedFields) == 0 {
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}
	header := CPSActionHeadersFromFields(resolvedFields)

	// 3. Build object key.
	ext := "csv"
	if fileType == "pdf" {
		ext = "pdf"
	}
	objectName := fmt.Sprintf(
		"cps_actions_%s_to_%s_%d.%s",
		startDate.Format("20060102"),
		endDate.Format("20060102"),
		time.Now().Unix(),
		ext,
	)

	var url string
	if fileType == "pdf" {
		rows := make([][]string, 0, len(data))
		for _, action := range data {
			row, rerr := BuildCPSActionRowFromFields(action, resolvedFields)
			if rerr != nil {
				return "", rerr
			}
			rows = append(rows, row)
		}

		colFontPt := make([]float64, len(resolvedFields))
		for i, key := range resolvedFields {
			if override, ok := CPSActionFieldFontOverride[key]; ok && override > 0 {
				colFontPt[i] = override
			}
		}

		url, err = ExportPDFAndUpload(ctx, minioClient, buckerName, cfg, objectName, header, rows, PDFExportOptions{PageSize: "A4"}, colFontPt, logger)
		if err != nil {
			logger.Errorf("[CPSExport] PDF export failed: %v", err)
			return "", errors.New(localization.CpsActionDataExportedError.Code)
		}
	} else {
		tmpFile, terr := os.CreateTemp("", "cps_actions_*.csv")
		if terr != nil {
			return "", fmt.Errorf("create temp file: %w", terr)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// UTF-8 BOM for Excel compatibility.
		if _, werr := tmpFile.Write([]byte{0xEF, 0xBB, 0xBF}); werr != nil {
			return "", fmt.Errorf("write BOM: %w", werr)
		}

		writer := csv.NewWriter(tmpFile)
		if werr := writer.Write(header); werr != nil {
			return "", fmt.Errorf("write header: %w", werr)
		}
		for _, action := range data {
			row, rerr := BuildCPSActionRowFromFields(action, resolvedFields)
			if rerr != nil {
				return "", rerr
			}
			if werr := writer.Write(row); werr != nil {
				return "", werr
			}
		}
		writer.Flush()
		if werr := writer.Error(); werr != nil {
			return "", fmt.Errorf("flush writer: %w", werr)
		}

		if _, serr := tmpFile.Seek(0, 0); serr != nil {
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
		stat, serr := tmpFile.Stat()
		if serr != nil {
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
		url, err = UploadCSVToMinio(ctx, minioClient, buckerName, tmpFile, stat.Size(), cfg, objectName, logger)
		if err != nil {
			return "", errors.New(localization.CpsActionDataExportedError.Code)
		}
	}

	baseURL := strings.TrimSuffix(cfg.MinioPublicEndPoint, "/")
	if baseURL != "" {
		url = fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(objectName, "/"))
	}
	return url, nil
}

func FileExporterForBPSAction(ctx context.Context, cfg config.VaultConfig, minioClient *s3.Client, buckerName string, filterMap *types.Filter, data []*imodel.BPSAction, _ func(fields []string) []string, logger utils.Logger) (string, error) {
	if filterMap == nil {
		return "", errors.New(localization.ErrorRequiredFieldMissing.Code)
	}
	if filterMap.Filters == nil {
		filterMap.Filters = map[string]interface{}{}
	}

	// 1. Date range validation.
	createdAtFrom, fromOk := filterMap.Filters["created_at_from"].(string)
	createdAtTo, toOk := filterMap.Filters["created_at_to"].(string)
	if !fromOk || !toOk || createdAtFrom == "" || createdAtTo == "" {
		return "", errors.New(localization.ErrorRequiredFieldMissing.Code)
	}

	startDate, endDate, err := local_util.FormatDateRangeToUTCStrings(createdAtFrom, createdAtTo)
	if err != nil {
		return "", errors.New(localization.ErrorInvalidDateFormat.Code)
	}
	filterMap.Filters["created_at_from"] = startDate
	filterMap.Filters["created_at_to"] = endDate

	if len(data) == 0 {
		return "", errors.New(localization.CpsActionDataNotFoundInDateRange.Code)
	}

	// 2. Resolve output format + field projection.
	fileType, _ := filterMap.Filters["file_type"].(string)
	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType == "" {
		fileType = string(FileTypeCSV)
	}

	requestedFields := extractFieldsFilter(filterMap.Filters)
	resolvedFields := ResolveBPSActionFields(requestedFields)
	if len(requestedFields) > 0 && len(resolvedFields) == 0 {
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}
	header := BPSActionHeadersFromFields(resolvedFields)

	// 3. Build object key.
	ext := "csv"
	if fileType == "pdf" {
		ext = "pdf"
	}
	objectName := fmt.Sprintf(
		"bps_actions_%s_to_%s_%d.%s",
		startDate.Format("20060102"),
		endDate.Format("20060102"),
		time.Now().Unix(),
		ext,
	)

	var url string
	if fileType == "pdf" {
		rows := make([][]string, 0, len(data))
		for _, action := range data {
			row, rerr := BuildBPSActionRowFromFields(action, resolvedFields)
			if rerr != nil {
				return "", rerr
			}
			rows = append(rows, row)
		}

		colFontPt := make([]float64, len(resolvedFields))
		for i, key := range resolvedFields {
			if override, ok := BPSActionFieldFontOverride[key]; ok && override > 0 {
				colFontPt[i] = override
			}
		}

		url, err = ExportPDFAndUpload(ctx, minioClient, buckerName, cfg, objectName, header, rows, PDFExportOptions{PageSize: "A4"}, colFontPt, logger)
		if err != nil {
			logger.Errorf("[BPSExport] PDF export failed: %v", err)
			return "", errors.New(localization.CpsActionDataExportedError.Code)
		}
	} else {
		tmpFile, terr := os.CreateTemp("", "bps_actions_*.csv")
		if terr != nil {
			return "", fmt.Errorf("create temp file: %w", terr)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// UTF-8 BOM for Excel compatibility.
		if _, werr := tmpFile.Write([]byte{0xEF, 0xBB, 0xBF}); werr != nil {
			return "", fmt.Errorf("write BOM: %w", werr)
		}

		writer := csv.NewWriter(tmpFile)
		if werr := writer.Write(header); werr != nil {
			return "", fmt.Errorf("write header: %w", werr)
		}
		for _, action := range data {
			row, rerr := BuildBPSActionRowFromFields(action, resolvedFields)
			if rerr != nil {
				return "", rerr
			}
			if werr := writer.Write(row); werr != nil {
				return "", werr
			}
		}
		writer.Flush()
		if werr := writer.Error(); werr != nil {
			return "", fmt.Errorf("flush writer: %w", werr)
		}

		if _, serr := tmpFile.Seek(0, 0); serr != nil {
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
		stat, serr := tmpFile.Stat()
		if serr != nil {
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
		url, err = UploadCSVToMinio(ctx, minioClient, buckerName, tmpFile, stat.Size(), cfg, objectName, logger)
		if err != nil {
			return "", errors.New(localization.CpsActionDataExportedError.Code)
		}
	}

	baseURL := strings.TrimSuffix(cfg.MinioPublicEndPoint, "/")
	if baseURL != "" {
		url = fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(objectName, "/"))
	}
	return url, nil
}

// CPSActionFieldSpec describes a single exportable column: the user-facing header label
// and the extractor that turns a CPSAction into the string value for that column.
type CPSActionFieldSpec struct {
	Header  string
	Extract func(a *model.CPSAction) string
}

type BPSActionFieldSpec struct {
	Header  string
	Extract func(a *imodel.BPSAction) string
}

// CPSActionExportSchema only exists to define the exported column order and the
// json/bson tags used by the ?fields= filter.
type CPSActionExportSchema struct {
	ID                string `json:"id" bson:"_id"`
	ActionCode        string `json:"action_code" bson:"action_code"`
	MakerID           string `json:"maker_id" bson:"maker_id"`
	MakerName         string `json:"maker_name" bson:"maker_name"`
	MakerPhoneNumber  string `json:"maker_phone_number" bson:"maker_phone_number"`
	CheckerUsers      string `json:"checker_users" bson:"checker_users"`
	AuditorUsers      string `json:"auditor_users" bson:"auditor_users"`
	AuditorNames      string `json:"auditor_names" bson:"auditor_names"`
	AuditorID         string `json:"auditor_id" bson:"auditor_id"`
	AuditorIDs        string `json:"auditor_ids" bson:"auditor_ids"`
	AuditorMark       string `json:"auditor_mark" bson:"auditor_mark"`
	AuditorMarks      string `json:"auditor_marks" bson:"auditor_marks"`
	AuditorStatus     string `json:"auditor_status" bson:"auditor_status"`
	CheckerName       string `json:"checker_name" bson:"checker_name"`
	CheckerID         string `json:"checker_id" bson:"checker_id"`
	CheckerIDs        string `json:"checker_ids" bson:"checker_ids"`
	ActionStatus      string `json:"action_status" bson:"action_status"`
	ActionType        string `json:"action_type" bson:"action_type"`
	RequestAction     string `json:"request_action" bson:"request_action"`
	CreatedAt         string `json:"created_at" bson:"created_at"`
	LastModifiedAt    string `json:"last_modified_at" bson:"last_modified_at"`
	MakerActionTime   string `json:"maker_action_time" bson:"maker_action_time"`
	CheckerActionTime string `json:"checker_action_time" bson:"checker_action_time"`
	AuditorActionTime string `json:"auditor_action_time" bson:"auditor_action_time"`
	RejectionReason   string `json:"rejection_reason" bson:"rejection_reason"`
	CanceledReason    string `json:"canceled_reason" bson:"canceled_reason"`
}

type BPSActionExportSchema struct {
	ID                 string `json:"id" bson:"_id"`
	ActionCode         string `json:"action_code" bson:"action_code"`
	IsAuditorApproved  string `json:"is_auditor_approved" bson:"is_auditor_approved"`
	UserID             string `json:"user_id" bson:"user_information.user_id"`
	UserCode           string `json:"user_code" bson:"user_information.user_code"`
	FullName           string `json:"full_name" bson:"user_information.full_name"`
	AccountNumbers     string `json:"account_numbers" bson:"user_information.account_numbers"`
	PhoneNumbers       string `json:"phone_numbers" bson:"user_information.phone_numbers"`
	UserBranchCode     string `json:"user_branch_code" bson:"user_information.branch_code"`
	BusinessID         string `json:"business_id" bson:"business.business_id"`
	TillNumber         string `json:"till_number" bson:"business.till_number"`
	BusinessName       string `json:"business_name" bson:"business.business_name"`
	CheckersNeeded     string `json:"checkers_needed" bson:"checkers_needed"`
	CheckersApproved   string `json:"checkers_approved" bson:"checkers_approved"`
	CheckerID          string `json:"checker_id" bson:"checker_id"`
	MakerUser          string `json:"maker_user" bson:"maker_user"`
	MakerName          string `json:"maker_name" bson:"maker_name"`
	MakerReason        string `json:"maker_reason" bson:"maker_reason"`
	MakerPhoneNumber   string `json:"maker_phone_number" bson:"maker_phone_number"`
	CheckerName        string `json:"checker_name" bson:"checker_name"`
	CheckerPhoneNumber string `json:"checker_phone_number" bson:"checker_phone_number"`
	ActionType         string `json:"action_type" bson:"action_reason.action_type"`
	ActionNote         string `json:"action_note" bson:"action_reason.action_note"`
	Identifier         string `json:"identifier" bson:"action_reason.identifier"`
	MakerMID           string `json:"maker_mid" bson:"maker_mid"`
	CheckerMID         string `json:"checker_mid" bson:"checker_mid"`
	AuditorMID         string `json:"auditor_mid" bson:"auditor_mid"`
	CheckerNameList    string `json:"checker_name_list" bson:"checker_name_list"`
	AuditorNameList    string `json:"auditor_name_list" bson:"auditor_name_list"`
	AuditorName        string `json:"auditor_name" bson:"auditors.auditor_name"`
	AuditorPhoneNumber string `json:"auditor_phone_number" bson:"auditors.auditor_phone_number"`
	AuditorsRequired   string `json:"auditors_required" bson:"auditors.auditors_required"`
	AuditorID          string `json:"auditor_id" bson:"auditors.auditor_id"`
	Audited            string `json:"audited" bson:"auditors.audited"`
	AuditorApproval    string `json:"auditor_approval" bson:"auditors.auditor_approval"`
	AuditorReason      string `json:"auditor_reason" bson:"auditors.reason"`
	CheckerTime        string `json:"checker_time" bson:"checker_time"`
	AuditorTime        string `json:"auditor_time" bson:"auditor_time"`
	RequestAction      string `json:"request_action" bson:"request_action"`
	Value              string `json:"value" bson:"value"`
	HomeBranch         string `json:"home_branch" bson:"home_branch"`
	AccountBranchCode  string `json:"account_branch_code" bson:"account_branch_code"`
	DistrictCode       string `json:"district_code" bson:"district_code"`
	BranchCode         string `json:"branch_code" bson:"branch_code"`
	LinkedDistrictCode string `json:"linked_district_code" bson:"linked_district_code"`
	AccountNumber      string `json:"account_number" bson:"account_number"`
	AccountHolderName  string `json:"account_holder_name" bson:"account_holder_name"`
	ServiceName        string `json:"service_name" bson:"service_name"`
	CurrentAction      string `json:"current_action" bson:"current_action"`
	PreviousAction     string `json:"previous_action" bson:"previous_action"`
	Time               string `json:"time" bson:"time"`
	Status             string `json:"status" bson:"status"`
	CustomerBarred     string `json:"customer_barred" bson:"customer_barred"`
	CreatedAt          string `json:"created_at" bson:"created_at"`
	LastModifiedAt     string `json:"last_modified_at" bson:"last_modified_at"`
}

var cpsActionFieldTagAliases = func() map[string]string {
	aliases := map[string]string{}
	specType := reflect.TypeOf(CPSActionExportSchema{})
	for i := 0; i < specType.NumField(); i++ {
		field := specType.Field(i)
		canonical := exportFieldKeyFromTag(field)
		if canonical == "" {
			continue
		}
		aliases[canonical] = canonical

		for _, tagName := range []string{"json", "bson"} {
			tag := strings.TrimSpace(field.Tag.Get(tagName))
			if tag == "" || tag == "-" {
				continue
			}
			key := strings.ToLower(strings.TrimSpace(strings.Split(tag, ",")[0]))
			if key != "" {
				aliases[key] = canonical
			}
		}
	}
	return aliases
}()

var bpsActionFieldTagAliases = func() map[string]string {
	aliases := map[string]string{}
	specType := reflect.TypeOf(BPSActionExportSchema{})
	for i := 0; i < specType.NumField(); i++ {
		field := specType.Field(i)
		canonical := exportFieldKeyFromTag(field)
		if canonical == "" {
			continue
		}
		aliases[canonical] = canonical

		for _, tagName := range []string{"json", "bson"} {
			tag := strings.TrimSpace(field.Tag.Get(tagName))
			if tag == "" || tag == "-" {
				continue
			}
			key := strings.ToLower(strings.TrimSpace(strings.Split(tag, ",")[0]))
			if key != "" {
				aliases[key] = canonical
			}
		}
	}

	aliases["maker_id"] = "maker_user"
	aliases["action_status"] = "status"
	aliases["auditor_status"] = "is_auditor_approved"
	aliases["checker_ids"] = "checker_id"
	aliases["auditor_ids"] = "auditor_id"
	aliases["auditor_names"] = "auditor_name_list"
	aliases["checker_action_time"] = "checker_time"
	aliases["auditor_action_time"] = "auditor_time"

	return aliases
}()

// CPSActionFieldRegistry maps field keys (as used in the ?fields=a,b,c query param) to
// their header label and row extractor. Add new exportable fields here.
var CPSActionFieldRegistry = map[string]CPSActionFieldSpec{
	"id":                  {"ID", cpsActionID},
	"action_code":         {"Action Code", func(a *model.CPSAction) string { return a.ActionCode }},
	"maker_id":            {"Maker ID", func(a *model.CPSAction) string { return a.MakerID }},
	"maker_name":          {"Maker Name", func(a *model.CPSAction) string { return a.MakerName }},
	"maker_phone_number":  {"Maker Phone Number", func(a *model.CPSAction) string { return a.MakerPhoneNumber }},
	"checker_users":       {"Checker Users", cpsCheckerUsers},
	"checker_user_info":   {"Checker Info", cpsCheckerUserInfo},
	"auditor_users":       {"Auditor Users", cpsAuditorUsers},
	"auditor_user_info":   {"Auditor Info", cpsAuditorUserInfo},
	"auditor_names":       {"Auditor Names", cpsAuditorNames},
	"auditor_id":          {"Auditor ID", cpsAuditorIDs},
	"auditor_ids":         {"Auditor ID", cpsAuditorIDs},
	"auditor_mark":        {"Auditor Mark", cpsAuditorMarks},
	"auditor_marks":       {"Auditor Mark", cpsAuditorMarks},
	"auditor_status":      {"Auditor Status", func(a *model.CPSAction) string { return string(a.AuditorStatus) }},
	"checker_name":        {"Checker Name", cpsCheckerNames},
	"checker_id":          {"Checker ID", cpsCheckerIDs},
	"checker_ids":         {"Checker ID", cpsCheckerIDs},
	"action_status":       {"Action Status", func(a *model.CPSAction) string { return a.ActionStatus }},
	"action_type":         {"Action Type", func(a *model.CPSAction) string { return a.ActionType }},
	"request_action":      {"Request Action", func(a *model.CPSAction) string { return a.RequestAction }},
	"created_at":          {"Created At", func(a *model.CPSAction) string { return local_util.FormatTime(a.CreatedAt) }},
	"last_modified_at":    {"Last Modified At", func(a *model.CPSAction) string { return local_util.FormatTime(a.LastModifiedAt) }},
	"maker_action_time":   {"Maker Action Time", func(a *model.CPSAction) string { return local_util.FormatTime(a.MakerActionTime) }},
	"checker_action_time": {"Checker Action Time", cpsCheckerActionTimes},
	"auditor_action_time": {"Auditor Action Time", cpsAuditorActionTimes},
	"rejection_reason":    {"Rejection Reason", func(a *model.CPSAction) string { return a.RejectionReason }},
	"canceled_reason":     {"Canceled Reason", func(a *model.CPSAction) string { return a.CanceledReason }},
}

var BPSActionFieldRegistry = map[string]BPSActionFieldSpec{
	"id":                   {"ID", bpsActionID},
	"action_code":          {"Action Code", func(a *imodel.BPSAction) string { return a.ActionCode }},
	"is_auditor_approved":  {"Is Auditor Approved", func(a *imodel.BPSAction) string { return strconv.FormatBool(a.IsAuditorApproved) }},
	"user_id":              {"User ID", func(a *imodel.BPSAction) string { return a.UserInformation.UserID }},
	"user_code":            {"User Code", func(a *imodel.BPSAction) string { return a.UserInformation.UserCode }},
	"full_name":            {"Full Name", func(a *imodel.BPSAction) string { return a.UserInformation.FullName }},
	"account_numbers":      {"Account Numbers", func(a *imodel.BPSAction) string { return strings.Join(a.UserInformation.AccountNumbers, ", ") }},
	"phone_numbers":        {"Phone Numbers", func(a *imodel.BPSAction) string { return a.UserInformation.PhoneNumbers }},
	"user_branch_code":     {"User Branch Code", func(a *imodel.BPSAction) string { return a.UserInformation.BranchCode }},
	"business_id":          {"Business ID", bpsBusinessID},
	"till_number":          {"Till Number", func(a *imodel.BPSAction) string { return a.BusinessInformation.TILLNumber }},
	"business_name":        {"Business Name", func(a *imodel.BPSAction) string { return a.BusinessInformation.BusinessName }},
	"checkers_needed":      {"Checkers Needed", func(a *imodel.BPSAction) string { return strconv.Itoa(a.CheckersNeeded) }},
	"checkers_approved":    {"Checkers Approved", func(a *imodel.BPSAction) string { return strconv.Itoa(a.CheckersApproved) }},
	"checker_id":           {"Checker ID", func(a *imodel.BPSAction) string { return strings.Join(a.CheckerID, ", ") }},
	"maker_user":           {"Maker User", func(a *imodel.BPSAction) string { return a.MakerID }},
	"maker_name":           {"Maker Name", func(a *imodel.BPSAction) string { return a.MakerName }},
	"maker_reason":         {"Maker Reason", func(a *imodel.BPSAction) string { return a.MakerReason }},
	"maker_phone_number":   {"Maker Phone Number", func(a *imodel.BPSAction) string { return a.MakerPhoneNumber }},
	"checker_name":         {"Checker Name", func(a *imodel.BPSAction) string { return a.CheckerName }},
	"checker_phone_number": {"Checker Phone Number", func(a *imodel.BPSAction) string { return a.CheckerPhoneNumber }},
	"action_type":          {"Action Type", func(a *imodel.BPSAction) string { return string(a.ActionReason.ActionType) }},
	"action_note":          {"Action Note", func(a *imodel.BPSAction) string { return a.ActionReason.ActionNote }},
	"identifier":           {"Identifier", func(a *imodel.BPSAction) string { return a.ActionReason.Identifier }},
	"maker_mid":            {"Maker MID", func(a *imodel.BPSAction) string { return a.MakerMID }},
	"checker_mid":          {"Checker MID", func(a *imodel.BPSAction) string { return strings.Join(a.CheckerMID, ", ") }},
	"auditor_mid":          {"Auditor MID", func(a *imodel.BPSAction) string { return strings.Join(a.AuditorMID, ", ") }},
	"checker_name_list":    {"Checker Name List", func(a *imodel.BPSAction) string { return strings.Join(a.CheckerNameList, ", ") }},
	"auditor_name_list":    {"Auditor Name List", func(a *imodel.BPSAction) string { return strings.Join(a.AuditorNameList, ", ") }},
	"auditor_name":         {"Auditor Name", func(a *imodel.BPSAction) string { return a.Auditors.AuditorName }},
	"auditor_phone_number": {"Auditor Phone Number", func(a *imodel.BPSAction) string { return a.Auditors.AuditorPhoneNumber }},
	"auditors_required":    {"Auditors Required", func(a *imodel.BPSAction) string { return strconv.Itoa(a.Auditors.AuditorsRequired) }},
	"auditor_id":           {"Auditor ID", func(a *imodel.BPSAction) string { return strings.Join(a.Auditors.AuditorID, ", ") }},
	"audited":              {"Audited", func(a *imodel.BPSAction) string { return strconv.FormatBool(a.Auditors.Audited) }},
	"auditor_approval":     {"Auditor Approval", func(a *imodel.BPSAction) string { return strconv.FormatBool(a.Auditors.AuditorApproval) }},
	"auditor_reason":       {"Auditor Reason", func(a *imodel.BPSAction) string { return a.Auditors.Reason }},
	"checker_time":         {"Checker Time", func(a *imodel.BPSAction) string { return bpsTimeSlice(a.CheckerTime) }},
	"auditor_time":         {"Auditor Time", func(a *imodel.BPSAction) string { return bpsTimeSlice(a.AuditorTime) }},
	"request_action":       {"Request Action", func(a *imodel.BPSAction) string { return string(a.RequestAction) }},
	"value":                {"Value", func(a *imodel.BPSAction) string { return a.EntityIdentifyer }},
	"home_branch":          {"Home Branch", func(a *imodel.BPSAction) string { return a.HomeBranch }},
	"account_branch_code":  {"Account Branch Code", func(a *imodel.BPSAction) string { return a.AccountBranchCode }},
	"district_code":        {"District Code", func(a *imodel.BPSAction) string { return a.DistrictCode }},
	"branch_code":          {"Branch Code", func(a *imodel.BPSAction) string { return a.BranchCode }},
	"linked_district_code": {"Linked District Code", func(a *imodel.BPSAction) string { return a.LinkedDistrictCode }},
	"account_number":       {"Account Number", func(a *imodel.BPSAction) string { return a.AccountNumber }},
	"account_holder_name":  {"Account Holder Name", func(a *imodel.BPSAction) string { return a.AccountHolderName }},
	"service_name":         {"Service Name", func(a *imodel.BPSAction) string { return a.ServiceName }},
	"current_action":       {"Current Action", func(a *imodel.BPSAction) string { return stringifyExportValue(a.CurrentAction) }},
	"previous_action":      {"Previous Action", func(a *imodel.BPSAction) string { return stringifyExportValue(a.PreviousAction) }},
	"time":                 {"Verified At", func(a *imodel.BPSAction) string { return local_util.FormatTime(a.VerifiedAt) }},
	"status":               {"Status", func(a *imodel.BPSAction) string { return a.Status }},
	"customer_barred":      {"Customer Barred", func(a *imodel.BPSAction) string { return strconv.FormatBool(a.CustomerBarred) }},
	"created_at":           {"Created At", func(a *imodel.BPSAction) string { return local_util.FormatTime(a.CreatedAt) }},
	"last_modified_at":     {"Last Modified At", func(a *imodel.BPSAction) string { return local_util.FormatTime(a.LastModifiedAt) }},
}

// CPSActionFieldFontOverride lets specific columns render in a smaller font than the
// table's responsive default. Useful for columns whose values are long fixed-format
// identifiers (e.g. "SRM26105_145513.440861" for action_code, 22 chars) that we don't
// want to truncate. Keys are CPSActionFieldRegistry keys; values are font points.
var CPSActionFieldFontOverride = map[string]float64{
	"action_code": 5,
}

var BPSActionFieldFontOverride = map[string]float64{
	"action_code":     5,
	"current_action":  5,
	"previous_action": 5,
}

// CPSActionDefaultFieldOrder is the field order used when no ?fields= is provided.
var CPSActionDefaultFieldOrder = []string{
	"id",
	"action_code",
	"maker_id",
	"maker_name",
	"maker_phone_number",
	"checker_user_info",
	"auditor_user_info",
	"auditor_status",
	"action_status",
	"action_type",
	"request_action",
	"created_at",
	"last_modified_at",
	"maker_action_time",
	"checker_action_time",
}

var BPSActionDefaultFieldOrder = []string{
	"id",
	"action_code",
	"user_code",
	"full_name",
	"phone_numbers",
	"maker_user",
	"maker_name",
	"maker_reason",
	"checker_name_list",
	"auditor_name_list",
	"status",
	"request_action",
	"account_number",
	"account_holder_name",
	"service_name",
	"created_at",
	"last_modified_at",
	"checker_time",
	"auditor_time",
}

func cpsAuditorNames(a *model.CPSAction) string {
	names := make([]string, 0, len(a.AuditorUsers))
	for _, au := range a.AuditorUsers {
		if au.AuditorName != "" {
			names = append(names, au.AuditorName)
		}
	}
	return strings.Join(names, ", ")
}

func cpsCheckerUsers(a *model.CPSAction) string {
	parts := make([]string, 0, len(a.CheckerUsers))
	for _, checker := range a.CheckerUsers {
		// entry, ok := formatIDTimestamp(checker.CheckerID, checker.ApprovedAt)
		// if !ok {
		// 	continue
		// }
		parts = append(parts, checker.CheckerID)
	}
	return strings.Join(parts, ", ")
}

func cpsAuditorUsers(a *model.CPSAction) string {
	parts := make([]string, 0, len(a.AuditorUsers))
	for _, auditor := range a.AuditorUsers {
		entry, ok := formatIDTimestamp(auditor.AuditorID, auditor.ApprovedAt)
		if !ok {
			continue
		}
		parts = append(parts, entry)
	}
	return strings.Join(parts, ", ")
}

func cpsAuditorIDs(a *model.CPSAction) string {
	ids := make([]string, 0, len(a.AuditorUsers))
	for _, au := range a.AuditorUsers {
		if au.AuditorID != "" {
			ids = append(ids, au.AuditorID)
		}
	}
	return strings.Join(ids, ", ")
}

func cpsAuditorMarks(a *model.CPSAction) string {
	marks := make([]string, 0, len(a.AuditorUsers))
	for _, au := range a.AuditorUsers {
		if au.AuditorMark != "" {
			marks = append(marks, string(au.AuditorMark))
		}
	}
	return strings.Join(marks, ", ")
}

func cpsCheckerNames(a *model.CPSAction) string {
	names := make([]string, 0, len(a.CheckerUsers))
	for _, c := range a.CheckerUsers {
		if c.CheckerName != "" {
			names = append(names, c.CheckerName)
		}
	}
	return strings.Join(names, ", ")
}

func cpsCheckerIDs(a *model.CPSAction) string {
	ids := make([]string, 0, len(a.CheckerUsers))
	for _, c := range a.CheckerUsers {
		if c.CheckerID != "" {
			ids = append(ids, c.CheckerID)
		}
	}
	return strings.Join(ids, ", ")
}

func cpsCheckerUserInfo(a *model.CPSAction) string {
	parts := make([]string, 0, len(a.CheckerUsers))
	for _, c := range a.CheckerUsers {
		if c.CheckerID == "" && c.CheckerName == "" {
			continue
		}
		entry := fmt.Sprintf("%s (%s / %s)", c.CheckerName, c.CheckerID, c.CheckerPhoneNumber)
		parts = append(parts, entry)
	}
	return strings.Join(parts, " | ")
}

func cpsAuditorUserInfo(a *model.CPSAction) string {
	parts := make([]string, 0, len(a.AuditorUsers))
	for _, au := range a.AuditorUsers {
		if au.AuditorID == "" && au.AuditorName == "" {
			continue
		}
		mark := ""
		if au.AuditorMark != "" {
			mark = " [" + string(au.AuditorMark) + "]"
		}
		entry := fmt.Sprintf("%s (%s / %s)%s", au.AuditorName, au.AuditorID, au.AuditorPhoneNumber, mark)
		parts = append(parts, entry)
	}
	return strings.Join(parts, " | ")
}

func cpsCheckerActionTimes(a *model.CPSAction) string {
	parts := make([]string, 0, len(a.CheckerUsers))
	for _, checker := range a.CheckerUsers {
		// entry, ok := formatIDTimestamp(checker.CheckerID, checker.ApprovedAt)
		// if !ok {
		// 	continue
		// }
		parts = append(parts, checker.ApprovedAt.Format("2006-01-02 15:04:05"))
	}
	return strings.Join(parts, ", ")
}

func cpsAuditorActionTimes(a *model.CPSAction) string {
	parts := make([]string, 0, len(a.AuditorUsers))
	for _, auditor := range a.AuditorUsers {
		entry, ok := formatIDTimestamp(auditor.AuditorID, auditor.ApprovedAt)
		if !ok {
			continue
		}
		parts = append(parts, entry)
	}
	return strings.Join(parts, ", ")
}

func formatIDTimestamp(id string, approvedAt time.Time) (string, bool) {
	id = strings.TrimSpace(id)
	if id == "" || approvedAt.IsZero() {
		return "", false
	}
	return fmt.Sprintf("%s", id), true
	// return fmt.Sprintf("%s:%s", id, local_util.FormatTime(approvedAt)), true
}

// ResolveCPSActionFields returns the effective export column keys. When ?fields= is omitted
// the default column set is used; when provided, only registered keys from that list are
// included, in the same order as the request (unknown keys are skipped).
func ResolveCPSActionFields(fields []string) []string {
	if len(fields) == 0 {
		return append([]string(nil), CPSActionDefaultFieldOrder...)
	}
	out := make([]string, 0, len(fields))
	seen := make(map[string]bool, len(fields))
	for _, f := range fields {
		key := strings.TrimSpace(strings.ToLower(f))
		if key == "" || seen[key] {
			continue
		}
		canonical, ok := cpsActionFieldTagAliases[key]
		if !ok {
			continue
		}
		if seen[canonical] {
			continue
		}
		if _, ok := CPSActionFieldRegistry[canonical]; ok {
			out = append(out, canonical)
			seen[key] = true
			seen[canonical] = true
		}
	}
	return out
}

func ResolveBPSActionFields(fields []string) []string {
	if len(fields) == 0 {
		return append([]string(nil), BPSActionDefaultFieldOrder...)
	}
	out := make([]string, 0, len(fields))
	seen := make(map[string]bool, len(fields))
	for _, f := range fields {
		key := strings.TrimSpace(strings.ToLower(f))
		if key == "" || seen[key] {
			continue
		}
		canonical, ok := bpsActionFieldTagAliases[key]
		if !ok {
			continue
		}
		if seen[canonical] {
			continue
		}
		if _, ok := BPSActionFieldRegistry[canonical]; ok {
			out = append(out, canonical)
			seen[key] = true
			seen[canonical] = true
		}
	}
	return out
}

// CPSActionHeadersFromFields maps the resolved field keys to their display labels.
func CPSActionHeadersFromFields(fields []string) []string {
	resolved := ResolveCPSActionFields(fields)
	out := make([]string, len(resolved))
	for i, key := range resolved {
		out[i] = CPSActionFieldRegistry[key].Header
	}
	return out
}

func BPSActionHeadersFromFields(fields []string) []string {
	resolved := ResolveBPSActionFields(fields)
	out := make([]string, len(resolved))
	for i, key := range resolved {
		out[i] = BPSActionFieldRegistry[key].Header
	}
	return out
}

// BuildCPSActionRowFromFields builds a row containing only the requested fields (in order).
func BuildCPSActionRowFromFields(a *model.CPSAction, fields []string) ([]string, error) {
	resolved := ResolveCPSActionFields(fields)
	row := make([]string, len(resolved))
	for i, key := range resolved {
		row[i] = CPSActionFieldRegistry[key].Extract(a)
	}
	return row, nil
}

// BuildBPSActionRowFromFields builds a row containing only the requested fields (in order).
func BuildBPSActionRowFromFields(a *imodel.BPSAction, fields []string) ([]string, error) {
	resolved := ResolveBPSActionFields(fields)
	row := make([]string, len(resolved))
	for i, key := range resolved {
		row[i] = BPSActionFieldRegistry[key].Extract(a)
	}
	return row, nil
}

func exportFieldKeyFromTag(field reflect.StructField) string {
	for _, tagName := range []string{"json", "bson"} {
		tag := strings.TrimSpace(field.Tag.Get(tagName))
		if tag == "" || tag == "-" {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(strings.Split(tag, ",")[0]))
		if key != "" {
			return key
		}
	}
	return strings.ToLower(field.Name)
}

// BPSActionFieldSpec describes a single exportable column for BPS actions.
// type BPSActionFieldSpec struct {
// 	Header  string
// 	Extract func(a *bps_model.BPSAction) string
// }

// // BPSActionFieldRegistry maps field keys to header + extractor for BPS action export.
// var BPSActionFieldRegistry = map[string]BPSActionFieldSpec{
// 	// --- Action metadata ---
// 	"action_code":    {"Action Code", func(a *bps_model.BPSAction) string { return a.ActionCode }},
// 	"request_action": {"Request Action", func(a *bps_model.BPSAction) string { return string(a.RequestAction) }},
// 	"status":         {"Status", func(a *bps_model.BPSAction) string { return a.Status }},
// 	"service_name":   {"Service Name", func(a *bps_model.BPSAction) string { return a.ServiceName }},
// 	// --- Maker ---
// 	"maker_id":           {"Maker ID", func(a *bps_model.BPSAction) string { return a.MakerID }},
// 	"maker_name":         {"Maker Name", func(a *bps_model.BPSAction) string { return a.MakerName }},
// 	"maker_phone_number": {"Maker Phone", func(a *bps_model.BPSAction) string { return a.MakerPhoneNumber }},
// 	"maker_reason":       {"Maker Reason", func(a *bps_model.BPSAction) string { return a.MakerReason }},
// 	// --- Checker (top-level fields) ---
// 	"checker_id":           {"Checker IDs", func(a *bps_model.BPSAction) string { return strings.Join(a.CheckerID, ", ") }},
// 	"checker_name":         {"Checker Name", func(a *bps_model.BPSAction) string { return a.CheckerName }},
// 	"checker_phone_number": {"Checker Phone", func(a *bps_model.BPSAction) string { return a.CheckerPhoneNumber }},
// 	"checker_name_list":    {"Checker Names", func(a *bps_model.BPSAction) string { return strings.Join(a.CheckerNameList, ", ") }},
// 	// --- Auditor (nested under Auditors) ---
// 	"auditor_id":           {"Auditor IDs", func(a *bps_model.BPSAction) string { return strings.Join(a.Auditors.AuditorID, ", ") }},
// 	"auditor_name":         {"Auditor Name", func(a *bps_model.BPSAction) string { return a.Auditors.AuditorName }},
// 	"auditor_phone_number": {"Auditor Phone", func(a *bps_model.BPSAction) string { return a.Auditors.AuditorPhoneNumber }},
// 	"auditor_approval": {"Auditor Approval", func(a *bps_model.BPSAction) string {
// 		if a.Auditors.AuditorApproval {
// 			return "APPROVED"
// 		}
// 		return "REJECTED"
// 	}},
// 	"auditor_reason": {"Auditor Reason", func(a *bps_model.BPSAction) string { return a.Auditors.Reason }},
// 	// --- Customer / Subject (nested under UserInformation) ---
// 	"customer_name":      {"Customer Name", func(a *bps_model.BPSAction) string { return a.UserInformation.FullName }},
// 	"customer_phone":     {"Customer Phone", func(a *bps_model.BPSAction) string { return a.UserInformation.PhoneNumbers }},
// 	"customer_branch":    {"Customer Branch Code", func(a *bps_model.BPSAction) string { return a.UserInformation.BranchCode }},
// 	"customer_user_code": {"Customer User Code", func(a *bps_model.BPSAction) string { return a.UserInformation.UserCode }},
// 	// --- Account ---
// 	"account_number": {"Account Number", func(a *bps_model.BPSAction) string { return a.AccountNumber }},
// 	"account_holder": {"Account Holder", func(a *bps_model.BPSAction) string { return a.AccountHolderName }},
// 	"customer_barred": {"Customer Barred", func(a *bps_model.BPSAction) string {
// 		if a.CustomerBarred {
// 			return "YES"
// 		}
// 		return "NO"
// 	}},
// 	// --- Timestamps ---
// 	"created_at":      {"Created At", func(a *bps_model.BPSAction) string { return local_util.FormatTime(a.CreatedAt) }},
// 	"last_modified_at": {"Last Modified At", func(a *bps_model.BPSAction) string { return local_util.FormatTime(a.LastModifiedAt) }},
// }

// // BPSActionDefaultFieldOrder is the column order used when no ?fields= is provided.
// var BPSActionDefaultFieldOrder = []string{
// 	"action_code",
// 	"request_action",
// 	"status",
// 	"service_name",
// 	"maker_id",
// 	"maker_name",
// 	"maker_phone_number",
// 	"checker_name",
// 	"checker_phone_number",
// 	"checker_name_list",
// 	"auditor_name",
// 	"auditor_phone_number",
// 	"auditor_approval",
// 	"auditor_reason",
// 	"customer_name",
// 	"customer_phone",
// 	"customer_branch",
// 	"customer_user_code",
// 	"account_number",
// 	"account_holder",
// 	"customer_barred",
// 	"created_at",
// 	"last_modified_at",
// }

// // ResolveBPSActionFields returns the effective export column keys for BPS actions.
// func ResolveBPSActionFields(fields []string) []string {
// 	if len(fields) == 0 {
// 		return append([]string(nil), BPSActionDefaultFieldOrder...)
// 	}
// 	var resolved []string
// 	for _, f := range fields {
// 		key := strings.ToLower(strings.TrimSpace(f))
// 		if _, ok := BPSActionFieldRegistry[key]; ok {
// 			resolved = append(resolved, key)
// 		}
// 	}
// 	return resolved
// }

// FileExporterForBPSAction streams BPS actions into CSV or PDF and uploads to MinIO.
// Triggered when filterMap.Filters["action"] == "export".
// Requires created_at_from + created_at_to. Supports ?fields= and ?file_type=csv|pdf.
// func FileExporterForBPSAction(ctx context.Context, cfg config.VaultConfig, minioClient *s3.Client, buckerName string, filterMap *types.Filter, data []*bps_model.BPSAction, logger utils.Logger) (string, error) {
// 	if filterMap == nil {
// 		return "", errors.New(localization.ErrorRequiredFieldMissing.Code)
// 	}
// 	if filterMap.Filters == nil {
// 		filterMap.Filters = map[string]interface{}{}
// 	}

// 	createdAtFrom, fromOk := filterMap.Filters["created_at_from"].(string)
// 	createdAtTo, toOk := filterMap.Filters["created_at_to"].(string)
// 	if !fromOk || !toOk || createdAtFrom == "" || createdAtTo == "" {
// 		return "", errors.New(localization.ErrorRequiredFieldMissing.Code)
// 	}

// 	startDate, endDate, err := local_util.FormatDateRangeToUTCStrings(createdAtFrom, createdAtTo)
// 	if err != nil {
// 		return "", errors.New(localization.ErrorInvalidDateFormat.Code)
// 	}

// 	if len(data) == 0 {
// 		return "", errors.New(localization.CpsActionDataNotFoundInDateRange.Code)
// 	}

// 	fileType, _ := filterMap.Filters["file_type"].(string)
// 	fileType = strings.ToLower(strings.TrimSpace(fileType))
// 	if fileType == "" {
// 		fileType = string(FileTypeCSV)
// 	}

// 	requestedFields := extractFieldsFilter(filterMap.Filters)
// 	resolvedFields := ResolveBPSActionFields(requestedFields)

// 	headers := make([]string, len(resolvedFields))
// 	for i, key := range resolvedFields {
// 		headers[i] = BPSActionFieldRegistry[key].Header
// 	}

// 	ext := "csv"
// 	if fileType == "pdf" {
// 		ext = "pdf"
// 	}
// 	objectName := fmt.Sprintf(
// 		"bps_actions_%s_to_%s_%d.%s",
// 		startDate.Format("20060102"),
// 		endDate.Format("20060102"),
// 		time.Now().Unix(),
// 		ext,
// 	)

// 	if fileType == "pdf" {
// 		rows := make([][]string, 0, len(data))
// 		for _, action := range data {
// 			row := make([]string, len(resolvedFields))
// 			for i, key := range resolvedFields {
// 				row[i] = BPSActionFieldRegistry[key].Extract(action)
// 			}
// 			rows = append(rows, row)
// 		}
// 		url, perr := ExportPDFAndUpload(ctx, minioClient, buckerName, cfg, objectName, headers, rows, PDFExportOptions{PageSize: "A4"}, nil, logger)
// 		if perr != nil {
// 			logger.Errorf("[BPSExport] PDF export failed: %v", perr)
// 			return "", errors.New(localization.CpsActionDataExportedError.Code)
// 		}
// 		return url, nil
// 	}

// 	tmpFile, terr := os.CreateTemp("", "bps_actions_*.csv")
// 	if terr != nil {
// 		return "", fmt.Errorf("create temp file: %w", terr)
// 	}
// 	defer os.Remove(tmpFile.Name())
// 	defer tmpFile.Close()

// 	if _, werr := tmpFile.Write([]byte{0xEF, 0xBB, 0xBF}); werr != nil {
// 		return "", fmt.Errorf("write BOM: %w", werr)
// 	}

// 	writer := csv.NewWriter(tmpFile)
// 	if werr := writer.Write(headers); werr != nil {
// 		return "", fmt.Errorf("write header: %w", werr)
// 	}
// 	for _, action := range data {
// 		row := make([]string, len(resolvedFields))
// 		for i, key := range resolvedFields {
// 			row[i] = BPSActionFieldRegistry[key].Extract(action)
// 		}
// 		if werr := writer.Write(row); werr != nil {
// 			return "", werr
// 		}
// 	}
// 	writer.Flush()
// 	if werr := writer.Error(); werr != nil {
// 		return "", fmt.Errorf("flush writer: %w", werr)
// 	}

// 	if _, serr := tmpFile.Seek(0, 0); serr != nil {
// 		return "", errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	stat, serr := tmpFile.Stat()
// 	if serr != nil {
// 		return "", errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	url, uerr := UploadCSVToMinio(ctx, minioClient, buckerName, tmpFile, stat.Size(), cfg, objectName, logger)
// 	if uerr != nil {
// 		return "", errors.New(localization.CpsActionDataExportedError.Code)
// 	}
// 	return url, nil
// }

func cpsActionID(a *model.CPSAction) string {
	if a == nil {
		return ""
	}
	return a.ID.Hex()
}

func bpsActionID(a *imodel.BPSAction) string {
	if a == nil {
		return ""
	}
	return a.ID.Hex()
}

func bpsBusinessID(a *imodel.BPSAction) string {
	if a == nil || a.BusinessInformation.BusinessID.IsZero() {
		return ""
	}
	return a.BusinessInformation.BusinessID.Hex()
}

func bpsTimeSlice(items []time.Time) string {
	if len(items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(items))
	for _, item := range items {
		if item.IsZero() {
			continue
		}
		parts = append(parts, local_util.FormatTime(item))
	}
	return strings.Join(parts, ", ")
}

func stringifyExportValue(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(raw)
}

// BuildCPSActionRow keeps the legacy signature (full default row). New callers should use
// BuildCPSActionRowFromFields so that per-field projections work correctly.
func BuildCPSActionRow(a *model.CPSAction) ([]string, error) {
	return BuildCPSActionRowFromFields(a, nil)
}

// 🔥 Generic producer with data
func ProduceFileFromData[T any](
	ctx context.Context,
	cfg FileProducerConfig,
	data []T,
	rowMapper func(T) ([]string, error),
	upload UploadFunc,
) (string, error) {

	ft := FileType(strings.ToLower(string(cfg.FileType)))
	if ft == "" {
		ft = FileTypeCSV
	}

	ext := string(ft)

	// 1️⃣ temp file
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s_*.%s", cfg.FilePrefix, ext))
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)

	// 2️⃣ header
	if ft == FileTypeCSV && len(cfg.Header) > 0 {
		if err := writer.Write(cfg.Header); err != nil {
			return "", fmt.Errorf("write header: %w", err)
		}
	}

	// 3️⃣ loop داخلي 🔥
	rowCount := 0
	for _, item := range data {
		row, err := rowMapper(item)
		if err != nil {
			return "", err
		}

		if ft == FileTypeCSV {
			if err := writer.Write(row); err != nil {
				return "", err
			}
		}

		rowCount++
	}

	writer.Flush()

	if rowCount == 0 {
		return "", fmt.Errorf("no data found")
	}

	// 4️⃣ prepare file
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return "", fmt.Errorf("seek file: %w", err)
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		return "", fmt.Errorf("stat file: %w", err)
	}

	// 5️⃣ object name
	objectName := cfg.ObjectName
	if objectName == "" {
		objectName = fmt.Sprintf("%s_%d.%s", cfg.FilePrefix, stat.ModTime().Unix(), ext)
	}

	if !strings.HasSuffix(objectName, "."+ext) {
		objectName += "." + ext
	}

	// 6️⃣ upload
	url, err := upload(ctx, tmpFile, stat.Size(), objectName, ft)
	if err != nil {
		return "", err
	}

	return url, nil
}

// ==============================================
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

		// --- DATE FILTER LOGIC (per-field) ---
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
			"debit_allowed", "credit_allowed", "has_restriction", "advert_for", "is_expired",
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

		// --- IS_EXPIRED: if is_expired=true, add created_at $lte time.Now() ---
		if isExpired, ok := filter["is_expired"]; ok {
			if expired, ok := isExpired.(bool); ok && expired {
				// if existing, ok := filter["created_at"].(bson.M); ok {
				// 	existing["$lt"] = time.Now()
				// } else {
				// }
				filter["end_date"] = bson.M{"$lt": time.Now()}
			}
		}
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
			"is_verified": true, "active_account": true, "is_expired": true,
		}

		for key, val := range filterParam.Filters {
			if !allowedSet[key] {
				continue
			}

			// --- DERIVED FILTER: is_expired -> created_at ---
			if key == "is_expired" {
				expired, ok := BoolFromInterface(val)
				if ok {
					operator := "<"
					if !expired {
						operator = ">="
					}
					filters = append(filters, fmt.Sprintf("created_at %s :%d", operator, idx))
					args = append(args, time.Now())
					idx++
					continue
				}
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
				if parsed, ok := BoolFromInterface(val); ok {
					filters = append(filters,
						fmt.Sprintf("%s = :%d", key, idx))
					args = append(args, parsed)
					idx++
					continue
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

func BoolFromInterface(v interface{}) (bool, bool) {
	switch b := v.(type) {
	case bool:
		return b, true
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(b))
		if err != nil {
			return false, false
		}
		return parsed, true
	case float64:
		return b != 0, true
	default:
		return false, false
	}
}

// parseDateInput parses a date string that can be either date-only ("2026-01-05")
// or full ISO datetime ("2026-01-05T07:10:33.695+00:00", "2026-01-05T07:10:33").
func parseDateInput(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,                // e.g. export handler FormatDateRangeToUTCStrings
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

func BuildCPSActionDateRangeFilter(filterMap *types.Filter, startDate, endDate time.Time) bson.M {
	rangeFilter := bson.M{
		"$gte": startDate,
		"$lte": endDate,
	}

	// CPS actions may populate different timestamp fields depending on version / pipeline.
	// Match if any known field falls in range (same idea as cps_action.repository buildCPSActionDateRangeFilter).
	return bson.M{
		"$or": []bson.M{
			{"created_at": rangeFilter},
			{"action_created_at": rangeFilter},
			{"maker_action_time": rangeFilter},
			{"last_modified_at": rangeFilter},
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

func UploadPDFToMinio(
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

	// ✅ Correct content type for PDF
	contentType := "application/pdf"

	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: &contentLength,
	}

	if _, err := s3Client.PutObject(ctx, putInput); err != nil {
		logger.Errorf("upload PDF failed error: %v", err)
		return "", err
	}

	baseURL := strings.TrimSuffix(env.MinioPublicEndPoint, "/")
	url := fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(objectKey, "/"))

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

// drawCBEBanner paints the Commercial Bank of Ethiopia branded header band across the
// top of the current PDF page: solid purple background with the bank name in white bold.
//
// pageWidthMM is the full page width (e.g. 210 for portrait A4); heightMM is the banner
// height. Text/fill colors are reset to defaults before this returns so subsequent
// content renders normally.
func drawCBEBanner(pdf *gofpdf.Fpdf, pageWidthMM, heightMM float64) {
	// Purple background band.
	pdf.SetFillColor(123, 45, 142) // #7B2D8E
	pdf.Rect(0, 0, pageWidthMM, heightMM, "F")

	// Bank name centered vertically, slightly inset from the left edge.
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetXY(12, heightMM/2-4)
	pdf.CellFormat(pageWidthMM-24, 8, "Commercial Bank of Ethiopia", "", 0, "L", false, 0, "")

	// Reset so subsequent content uses default black/white.
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)
}

// ExportPDFAndUpload builds a tabular PDF from headers + rows and uploads it to MinIO.
// Orientation, page-size upgrades (A5→A4), and column paging are chosen automatically
// from the column count so headers and body cells always share the same widths.
// colFontPt may be nil or len(headers); non-zero entries override the body font for that column.
func ExportPDFAndUpload(
	ctx context.Context,
	s3Client *s3.Client,
	bucketName string,
	env config.VaultConfig,
	objectKey string,
	headers []string,
	rows [][]string,
	opts PDFExportOptions,
	colFontPt []float64,
	logger interface {
		Errorf(format string, args ...any)
	},
) (string, error) {

	tmpFile, err := os.CreateTemp("", "export_*.pdf")
	if err != nil {
		logger.Errorf("[ExportPDFAndUpload] create temp file: %v", err)
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	colCount := len(headers)
	pageSize, orientation, sections := pdfEffectiveLayout(opts.PageSize, colCount)
	pageWidth, pageHeight := pdfPageSizeMM(pageSize, orientation)

	pdf := gofpdf.New(orientation, "mm", pageSize, "")
	pdf.SetMargins(pdfSideMarginMM, pdfBannerHeightMM+4, pdfSideMarginMM)
	// Row placement is manual so a wrapped row is never split across pages.
	pdf.SetAutoPageBreak(false, 0)
	pdf.AliasNbPages("")

	contentBottomY := pageHeight - pdfFooterMarginMM

	type sectionState struct {
		headers []string
		rows    [][]string
		colFont []float64
		label   string
		colW    float64
		layout  PDFLayout
		lineHt  float64
	}

	buildSection := func(start, end int, label string) sectionState {
		secHeaders := headers[start:end]
		secRows := make([][]string, len(rows))
		for i, row := range rows {
			if end <= len(row) {
				secRows[i] = row[start:end]
			} else if start < len(row) {
				secRows[i] = row[start:]
			}
		}
		var secFont []float64
		if len(colFontPt) > start {
			secFont = colFontPt[start:end]
		}
		usable := pageWidth - 2*pdfSideMarginMM
		colW := usable
		if len(secHeaders) > 0 {
			colW = usable / float64(len(secHeaders))
		}
		layout := CalcPDFLayout(len(secHeaders))
		return sectionState{
			headers: secHeaders,
			rows:    secRows,
			colFont: secFont,
			label:   label,
			colW:    colW,
			layout:  layout,
			lineHt:  pdfLineHeightMM(layout.BodyFontPt),
		}
	}

	var active sectionState
	pdf.SetHeaderFunc(func() {
		drawCBEBanner(pdf, pageWidth, pdfBannerHeightMM)
		y := pdfBannerHeightMM + 2
		if active.label != "" {
			pdf.SetXY(pdfSideMarginMM, y)
			pdf.SetFont("Arial", "I", 8)
			pdf.CellFormat(pageWidth-2*pdfSideMarginMM, 4, active.label, "", 1, "L", false, 0, "")
			y += 5
		}
		pdf.SetXY(pdfSideMarginMM, y)
		pdf.SetFont("Arial", "B", active.layout.HeaderFontPt)
		pdf.SetTextColor(0, 0, 0)
		active.lineHt = pdfLineHeightMM(active.layout.HeaderFontPt)
		pdfDrawTableRow(pdf, pdfSideMarginMM, y, active.headers, active.colW, active.lineHt, active.layout.HeaderRowMM, true)
		pdf.SetFont("Arial", "", active.layout.BodyFontPt)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(-12)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 8, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	for _, bounds := range sections {
		label := ""
		if len(sections) > 1 {
			label = fmt.Sprintf("Columns %d-%d of %d", bounds[0]+1, bounds[1], colCount)
		}
		active = buildSection(bounds[0], bounds[1], label)
		pdf.AddPage()

		pdf.SetFont("Arial", "", active.layout.BodyFontPt)
		active.lineHt = pdfLineHeightMM(active.layout.BodyFontPt)

		for _, row := range active.rows {
			if len(row) < len(active.headers) {
				padded := make([]string, len(active.headers))
				copy(padded, row)
				row = padded
			}
			// Per-column font overrides use the smallest font in the row so wrapped lines align.
			rowFont := active.layout.BodyFontPt
			for i := range row {
				if i < len(active.colFont) && active.colFont[i] > 0 && active.colFont[i] < rowFont {
					rowFont = active.colFont[i]
				}
			}
			if rowFont != active.layout.BodyFontPt {
				pdf.SetFont("Arial", "", rowFont)
			}
			lineHt := pdfLineHeightMM(rowFont)
			rowHt := pdfCalcRowHeight(pdf, row, active.colW, lineHt, active.layout.BodyRowMM)

			// Keep the entire row on one page; header callback redraws column titles.
			if pdf.GetY()+rowHt > contentBottomY {
				pdf.AddPage()
			}

			x := pdfSideMarginMM
			y := pdf.GetY()
			pdfDrawTableRow(pdf, x, y, row, active.colW, lineHt, active.layout.BodyRowMM, false)
			if rowFont != active.layout.BodyFontPt {
				pdf.SetFont("Arial", "", active.layout.BodyFontPt)
			}
		}
	}

	if err := pdf.Error(); err != nil {
		logger.Errorf("[ExportPDFAndUpload] pdf error: %v", err)
		return "", fmt.Errorf("pdf error: %w", err)
	}
	if err := pdf.Output(tmpFile); err != nil {
		logger.Errorf("[ExportPDFAndUpload] output pdf: %v", err)
		return "", fmt.Errorf("output pdf: %w", err)
	}
	if err := tmpFile.Sync(); err != nil {
		return "", fmt.Errorf("sync temp file: %w", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return "", fmt.Errorf("seek temp file: %w", err)
	}
	stat, err := tmpFile.Stat()
	if err != nil {
		return "", fmt.Errorf("stat temp file: %w", err)
	}

	url, err := UploadPDFToMinio(ctx, s3Client, bucketName, tmpFile, stat.Size(), env, objectKey, logger)
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
