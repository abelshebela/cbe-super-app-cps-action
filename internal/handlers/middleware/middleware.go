package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/golang-jwt/jwt/v5"

	cps_auth "cbe-super-app-cps-action/grpc/auth/proto"
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func CORS(cfg *config.VaultConfig) func(http.Handler) http.Handler {
	var allowedOrigins []string
	allowCredentials := true

	switch cfg.GoEnv {
	case "production":
		allowedOrigins = []string{"https://superappcps.cbe.com.et"}
	case "uat":
		allowedOrigins = []string{"localhost:3000", "http://localhost:3000", "https://uat-cbe-super-app-central-portal.vercel.app", "https://dev-cbe-super-app-central-portal.vercel.app", "https://cpsportal-uat.cbe.com.et"}
	case "staging":
		allowedOrigins = []string{"0.0.0.0:3000", "https://staging-cbe-super-app-central-portal.vercel.app", "https://cpsportal-stg.cbe.com.et"}
	case "qa":
		allowedOrigins = []string{"localhost:3000", "http://localhost:3000", "0.0.0.0:3000", "https://qa-cbe-super-app-central-portal.vercel.app", "https://dev-cbe-super-app-central-portal.vercel.app"}
	case "dev":
		allowedOrigins = []string{"localhost:3000", "http://localhost:3000", "0.0.0.0:3000", "https://dev-cbe-super-app-central-portal.vercel.app", "https://dev-cbe-super-app-central-portal.vercel.app/"}
	default:
		allowedOrigins = []string{"*"}
		allowCredentials = false
	}

	return cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "X-CSRF-Token", "Accept", "Origin", "x-api-applicationid", "x-source-secret", "X-Api-Key", "x-api-key", "enable_encryption", "X-Session-ID"},
		ExposedHeaders:   []string{"X-Refreshed-Token", "x-action-code"},
		AllowCredentials: allowCredentials,
		MaxAge:           300,
	})
}

type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := JSONResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type UserPayload struct {
	IsERP       bool     `json:"is_erp,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	UserName    string   `json:"username,omitempty"`
	UserRole    string   `json:"user_role,omitempty"`
	RoleId      string   `json:"role_code,omitempty"`
	UserID      string   `json:"user_id,omitempty"`
	UserCode    string   `json:"user_code,omitempty"`
	FullName    string   `json:"full_name,omitempty"`
	Department  string   `json:"department,omitempty"`
	NextStep    string   `json:"next_step,omitempty"`
	Action      string   `json:"action"`
	DeviceID    string   `json:"device_id,omitempty"`
	SessionExp  int64    `json:"session_expiry,omitempty"`
	Environment string   `json:"environment"`
	Permission  []string `json:"permission_group"`
}

type authMiddleware struct {
	client           cps_auth.CpsAuthServiceClient
	logger           utils.Logger
	JWTSecretKey     string
	Key              string
	IV               string
	cfg              config.VaultConfig
	redisRepository  storage.RedisRepository
	approveIndexRepo CPSActionApproveIndexRepository
}

type CPSActionApproveIndexRepository interface {
	FindByRoleAndAction(ctx context.Context, roleID string, actionName string, version int64) (*model.CPSActionApproveIndex, error)
}

type AuthMiddleware interface {
	AccessControl(allowedRoles []string) func(http.Handler) http.Handler
	AuthenticateToken(next http.Handler) http.Handler
	AuthenticateTokenOrMerchantIntegrationAPIKey(next http.Handler) http.Handler
	AuthenticateTempToken(next http.Handler) http.Handler
	RequireFormContentType() func(http.Handler) http.Handler
	ValidateRequiredRoles(next http.Handler) http.Handler
}

func InitAuthMiddleware(client cps_auth.CpsAuthServiceClient, redisRepository storage.RedisRepository, approveIndexRepo CPSActionApproveIndexRepository, secretKey, key, iv string, cfg config.VaultConfig, logger utils.Logger) AuthMiddleware {
	return &authMiddleware{
		client:           client,
		redisRepository:  redisRepository,
		approveIndexRepo: approveIndexRepo,
		JWTSecretKey:     secretKey,
		Key:              key,
		IV:               iv,
		cfg:              cfg,
		logger:           logger,
	}
}

func (a *authMiddleware) RequireFormContentType() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				localization.SendBadRequestResponse(w, localization.ErrorFileInvalidType.Code)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (a *authMiddleware) AuthenticateTempToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := local_util.LoggerFromCtx(r.Context(), a.logger)

		bearer := "Bearer "
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, bearer) {
			log.Warnf("[AuthMW][AuthTempToken] bearer missing")
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearer)
		if tokenString == "" {
			log.Warnf("[AuthMW][AuthTempToken] empty token")
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		data, err := a.validateToken(r.Context(), tokenString)
		if err != nil {
			localization.SendErrorResponse(w, localization.ErrorInvalidToken, nil, nil)
			return
		}

		userPayload, err := a.extractUserPayload(r.Context(), data)
		if err != nil {
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		if userPayload.SessionExp <= time.Now().Unix() {
			log.Warnf("[AuthMW][AuthTempToken] session expired")
			localization.SendUnauthorizedResponse(w, localization.ErrorSessionExpired.Message)
			return
		}

		ctx := a.setUserPayload(r.Context(), userPayload)
		localization.UpdateWriterContext(w, ctx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a *authMiddleware) AccessControl(allowedRoles []string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[strings.ToUpper(r)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		})
	}
}

func (a *authMiddleware) AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := local_util.LoggerFromCtx(r.Context(), a.logger)
		if isErp, ok := r.Context().Value(constants.ContextKey("is_erp")).(bool); ok && isErp {
			log.Infof("[AuthMW][AuthToken] Skipping bearer token validation (isERP=true)")
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Expose-Headers", "X-Refreshed-Token")
		authHeader := r.Header.Get("Authorization")
		bearer := "Bearer "

		if !strings.HasPrefix(authHeader, bearer) {
			log.Warnf("[AuthMW][AuthToken] bearer missing")
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		tokenString := authHeader[len(bearer):]

		if tokenString == "" {
			log.Warnf("[AuthMW][AuthToken] empty token")
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		data, err := a.validateToken(r.Context(), tokenString)
		if err != nil {
			localization.SendErrorResponse(w, localization.ErrorInvalidToken, nil, nil)
			return
		}

		userPayload, err := a.extractUserPayload(r.Context(), data)
		if err != nil {
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		ctx := a.setUserPayload(r.Context(), userPayload)
		log = local_util.LoggerFromCtx(ctx, a.logger)
		localization.UpdateWriterContext(w, ctx)
		now := time.Now().Unix()

		remainTime, err := strconv.Atoi(a.cfg.JWTAccessExpirationMinutesRemain)
		if err != nil || remainTime == 0 {
			remainTime = 120
		} else {
			remainTime *= 60
		}
		redisDeviceIDExpireTime, err := strconv.Atoi(a.cfg.IdelUserTimeoutInMinutes)
		if err != nil || redisDeviceIDExpireTime == 0 {
			redisDeviceIDExpireTime = 15 * 60
		} else {
			redisDeviceIDExpireTime *= 60
		}

		deviceID, err := a.redisRepository.Get(r.Context(), fmt.Sprintf("%s:%s", constants.RedisCPSUserDeviceIDPrefix, userPayload.UserID))
		if err != nil {
			if errors.Is(err, redis.Nil) {
				log.Warnf("[AuthMW][AuthToken] redis device not found: %v", err)
				localization.SendUnauthorizedResponse(w, localization.ErrorSessionExpired.Message)
				return
			}
			log.Warnf("[AuthMW][AuthToken] redis err: %v", err)
		}

		deviceID = strings.Trim(deviceID, "\"")
		ctx = context.WithValue(ctx, constants.ContextKey("session_id"), deviceID)
		ctx = local_util.CtxWithLogger(ctx, local_util.NewSessionLogger(a.logger, deviceID))

		log = local_util.LoggerFromCtx(ctx, a.logger)
		log.Infof("[AuthMW][AuthToken] redis_device: %s payload_device: %s user: %s expire_length: %d", deviceID, userPayload.DeviceID, userPayload.UserID, redisDeviceIDExpireTime)

		if userPayload.SessionExp != 0 {
			log.Infof("[AuthMW][AuthToken] session_exp: %d now: %d remain: %d", userPayload.SessionExp, now, userPayload.SessionExp-now)

			if userPayload.SessionExp <= now {
				log.Warnf("[AuthMW][AuthToken] session expired")
				localization.SendUnauthorizedResponse(w, localization.ErrorSessionExpired.Message)
				return
			} else if deviceID != "" && deviceID == userPayload.DeviceID {
				if err := a.redisRepository.Set(r.Context(), fmt.Sprintf("%s:%s", constants.RedisCPSUserDeviceIDPrefix, userPayload.UserID), deviceID, time.Second*time.Duration(redisDeviceIDExpireTime)); err != nil {
					log.Warnf("[AuthMW][AuthToken] redis set device err user: %s: %v", userPayload.UserID, err)
				} else {
					log.Infof("[AuthMW][AuthToken] device updated user: %s", userPayload.UserID)
				}
			} else {
				log.Warnf("[AuthMW][AuthToken] device mismatch or empty device id, redis_device: %s payload_device: %s", deviceID, userPayload.DeviceID)
				localization.SendUnauthorizedResponse(w, localization.ErrorSessionExpired.Message)
				return
			}
		} else {
			log.Errorf("[AuthMW][AuthToken] unauth redis_device: %s payload_device: %s session: %v remain: %d", deviceID, userPayload.DeviceID, userPayload.SessionExp, remainTime)
			localization.SendUnauthorizedResponse(w, localization.ErrorSessionExpired.Message)
			return
		}

		r = r.WithContext(ctx)
		whitelist := []string{"departments", "access_list_segmentation", "account_block", "cps-action-list", "bps-action-list", "cps_action", "cps_actions", "actions", "bps_actions", "account_lookup", "cps_users", "roles"}
		if err := CPSActionRouteGuard(r, whitelist); err != nil {
			log.Warnf("[AuthMW][AuthToken] route guard blocked access to path: %s error: %v", r.URL.Path, err)
			localization.SendUnauthorizedResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *authMiddleware) validateToken(ctx context.Context, tokenString string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	jwtSecret := []byte(a.JWTSecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			log.Errorf("[AuthMW][ValidateToken] invalid signing algorithm: %v", token.Header["alg"])
			return nil, fmt.Errorf("invalid signing algorithm")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		log.Errorf("[AuthMW][ValidateToken] invalid/expired token: %v", err)
		return "", errors.New(localization.ErrorUserUnauthorized.Code)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Errorf("[AuthMW][ValidateToken] cast claims err")
		return "", errors.New(localization.ErrorUserUnauthorized.Code)
	}

	data, ok := claims["data"].(string)
	if !ok || data == "" {
		log.Errorf("[AuthMW][ValidateToken] no data in token")
		return "", errors.New(localization.ErrorUserUnauthorized.Code)
	}

	return data, nil
}

func (a *authMiddleware) extractUserPayload(ctx context.Context, data string) (UserPayload, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	decryptedUser, err := a.decryptUserData(ctx, data)
	if err != nil {
		return UserPayload{}, errors.New(localization.ErrorUserUnauthorized.Code)
	}

	var userPayload UserPayload
	err = json.Unmarshal([]byte(decryptedUser), &userPayload)
	if err != nil {
		log.Errorf("[AuthMW][ExtractPayload] unmarshal err: %v", err)
		return UserPayload{}, errors.New(localization.ErrorUserUnauthorized.Code)
	}

	return userPayload, nil
}

func (a *authMiddleware) setUserPayload(ctx context.Context, userPayload UserPayload) context.Context {
	if userPayload.IsERP {
		ctx = context.WithValue(ctx, constants.ContextKey("is_erp"), userPayload.IsERP)
	}
	ctx = context.WithValue(ctx, constants.ContextKey("user_role"), userPayload.UserRole)
	if userPayload.RoleId != "" {
		ctx = context.WithValue(ctx, constants.ContextKey("role_code"), userPayload.RoleId)
	}
	ctx = context.WithValue(ctx, constants.ContextKey("user_id"), userPayload.UserID)
	ctx = context.WithValue(ctx, constants.ContextKey("phone_number"), userPayload.PhoneNumber)
	ctx = context.WithValue(ctx, constants.ContextKey("user_code"), userPayload.UserCode)
	ctx = context.WithValue(ctx, constants.ContextKey("full_name"), userPayload.FullName)
	ctx = context.WithValue(ctx, constants.ContextKey("username"), userPayload.UserName)
	ctx = context.WithValue(ctx, constants.ContextKey("department"), userPayload.Department)
	ctx = context.WithValue(ctx, constants.ContextKey("next_step"), userPayload.NextStep)
	ctx = context.WithValue(ctx, constants.ContextKey("action"), userPayload.Action)
	ctx = context.WithValue(ctx, constants.ContextKey("permission"), userPayload.Permission)
	ctx = context.WithValue(ctx, constants.ContextKey("environment"), userPayload.Environment)
	return ctx
}

func (a *authMiddleware) decryptUserData(ctx context.Context, data string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	keyByte := []byte(a.Key)
	ivByte := []byte(a.IV)

	if len(keyByte) != 32 {
		log.Warnf("[AuthMW][Decrypt] invalid key bytes")
		return "", errors.New("key must be 32 bytes for AES-256")
	}
	if len(ivByte) != aes.BlockSize {
		log.Warnf("[AuthMW][Decrypt] invalid IV bytes")
		return "", errors.New("IV must be 16 bytes for AES-256-CBC")
	}

	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		log.Errorf("[AuthMW][Decrypt] hex decode err: %v", err)
		return "", fmt.Errorf("hex decode failed: %w", err)
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		log.Errorf("[AuthMW][Decrypt] cipher err: %v", err)
		return "", fmt.Errorf("NewCipher failed: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, ivByte)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted, err = a.pkcs7Unpad(ctx, decrypted, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpad failed: %w", err)
	}

	result := string(decrypted)
	return result, nil
}

func (a *authMiddleware) pkcs7Unpad(ctx context.Context, data []byte, blockSize int) ([]byte, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	if len(data) == 0 {
		log.Warnf("[AuthMW][Unpad] empty input")
		return nil, errors.New("input data is empty")
	}
	if len(data)%blockSize != 0 {
		log.Warnf("[AuthMW][Unpad] invalid length")
		return nil, errors.New("input length is not a multiple of block size")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		log.Warnf("[AuthMW][Unpad] invalid padding")
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			log.Warnf("[AuthMW][Unpad] invalid padding bytes")
			return nil, errors.New("invalid padding bytes")
		}
	}
	return data[:len(data)-padding], nil
}

func (a *authMiddleware) AuthenticateServiceAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := local_util.LoggerFromCtx(r.Context(), a.logger)
		got := strings.TrimSpace(r.Header.Get("x-api-key"))
		want := strings.TrimSpace(a.cfg.CPSApiTokenForCBE)
		if got == "" || want == "" || len(got) != len(want) ||
			subtle.ConstantTimeCompare([]byte(got), []byte(want)) != 1 {
			log.Warnf("[AuthMW][ServiceAPIKey] unauthorized or missing x-api-key")
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *authMiddleware) ValidateRequiredRoles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := local_util.LoggerFromCtx(r.Context(), a.logger)
		roleCode, ok := r.Context().Value(constants.ContextKey("role_code")).(string)
		if !ok || strings.TrimSpace(roleCode) == "" {
			log.Warnf("[AuthMW][ValidateRequiredRoles] role_code not found in context")
			WriteJSONResponse(w, http.StatusForbidden, "Unauthorized access: User role not found", nil)
			return
		}

		ctx := r.Context()
		hasValidRole, err := a.validateRoleInApproverIndexWithRequest(ctx, roleCode, r)
		if err != nil {
			log.Errorf("[AuthMW][ValidateRequiredRoles] error validating role in approver index: %v", err)
			WriteJSONResponse(w, http.StatusInternalServerError, "Error validating user permissions", nil)
			return
		}

		if !hasValidRole {
			log.Warnf("[AuthMW][ValidateRequiredRoles] unauthorized access attempt for role_code: %s on %s %s", roleCode, r.Method, r.URL.Path)
			WriteJSONResponse(w, http.StatusForbidden, "Unauthorized access: Insufficient permissions", nil)
			return
		}

		log.Infof("[AuthMW][ValidateRequiredRoles] access granted for role_code: %s on %s %s", roleCode, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (a *authMiddleware) validateRoleInApproverIndexWithRequest(ctx context.Context, roleCode string, r *http.Request) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	repo := GetCPSActionApproveRepo()
	if repo == nil {
		log.Warnf("[AuthMW][ValidateRequiredRoles] CPSActionApproveRepo is nil, using fallback validation")
		return a.fallbackRoleValidation(roleCode), nil
	}

	actionName := GetActionNameFromPath(r.Method, r.URL.Path)
	if actionName == "" {
		log.Warnf("[AuthMW][ValidateRequiredRoles] no action name found for path %s %s, trying common actions", r.Method, r.URL.Path)

		commonActions := []string{"CPS_ACTION", "BPS_ACTION", "ACCOUNT_BLOCK", "USER_MANAGEMENT", "CUSTOMER_MANAGEMENT"}

		for _, actionName := range commonActions {
			result, err := repo.FindByRoleAndAction(ctx, roleCode, actionName, 1)
			if err != nil {
				log.Warnf("[AuthMW][ValidateRequiredRoles] error checking role %s for action %s: %v", roleCode, actionName, err)
				continue
			}

			if result != nil {
				if result.ViewerIndex != nil || result.MakerIndex != nil ||
					result.CheckerIndex != nil || result.AuditorIndex != nil {
					log.Infof("[AuthMW][ValidateRequiredRoles] role %s found with approval index for action %s", roleCode, actionName)
					return true, nil
				}
			}
		}

		log.Warnf("[AuthMW][ValidateRequiredRoles] no approval indices found for role %s, trying fallback validation", roleCode)
		return a.fallbackRoleValidation(roleCode), nil
	}

	result, err := repo.FindByRoleAndAction(ctx, roleCode, actionName, 1)
	if err != nil {
		log.Warnf("[AuthMW][ValidateRequiredRoles] error checking role %s for action %s: %v", roleCode, actionName, err)
		return a.fallbackRoleValidation(roleCode), nil
	}

	if result != nil {
		if result.ViewerIndex != nil || result.MakerIndex != nil ||
			result.CheckerIndex != nil || result.AuditorIndex != nil {
			log.Infof("[AuthMW][ValidateRequiredRoles] role %s found with approval index for action %s", roleCode, actionName)
			return true, nil
		}
	}

	log.Warnf("[AuthMW][ValidateRequiredRoles] no approval indices found for role %s for action %s, trying fallback validation", roleCode, actionName)
	return a.fallbackRoleValidation(roleCode), nil
}

func (a *authMiddleware) validateRoleInApproverIndex(ctx context.Context, roleCode string) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	repo := GetCPSActionApproveRepo()
	if repo == nil {
		log.Warnf("[AuthMW][ValidateRequiredRoles] CPSActionApproveRepo is nil, using fallback validation")
		return a.fallbackRoleValidation(roleCode), nil
	}

	commonActions := []string{"CPS_ACTION", "BPS_ACTION", "ACCOUNT_BLOCK", "USER_MANAGEMENT", "CUSTOMER_MANAGEMENT"}

	for _, actionName := range commonActions {
		result, err := repo.FindByRoleAndAction(ctx, roleCode, actionName, 1)
		if err != nil {
			log.Warnf("[AuthMW][ValidateRequiredRoles] error checking role %s for action %s: %v", roleCode, actionName, err)
			continue
		}

		if result != nil {
			if result.ViewerIndex != nil || result.MakerIndex != nil ||
				result.CheckerIndex != nil || result.AuditorIndex != nil {
				log.Infof("[AuthMW][ValidateRequiredRoles] role %s found with approval index for action %s", roleCode, actionName)
				return true, nil
			}
		}
	}

	log.Warnf("[AuthMW][ValidateRequiredRoles] no approval indices found for role %s, trying fallback validation", roleCode)
	return a.fallbackRoleValidation(roleCode), nil
}

func (a *authMiddleware) fallbackRoleValidation(roleCode string) bool {
	validRolePatterns := []string{
		"MAKER_", "CHECKER_", "AUDITOR_", "VIEWER_",
		"IFB_MAKER_", "IFB_CHECKER_",
	}

	roleUpper := strings.ToUpper(roleCode)
	for _, pattern := range validRolePatterns {
		if strings.Contains(roleUpper, pattern) {
			return true
		}
	}

	exactRoles := []string{
		constants.Maker, constants.Checker, constants.Auditor, constants.Viewer,
		constants.IFBMaker, constants.IFBChecker,
	}

	for _, role := range exactRoles {
		if strings.EqualFold(roleCode, role) {
			return true
		}
	}

	return false
}
