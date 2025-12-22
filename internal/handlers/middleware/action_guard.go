package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// simple TTL cache for role/action checks
type allowEntry struct {
	allow bool
	exp   time.Time
}

type allowCache struct {
	mu   sync.RWMutex
	data map[string]allowEntry
	ttl  time.Duration
}

func (c *allowCache) get(key string) (val allowEntry, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok = c.data[key]
	if !ok {
		return
	}
	if time.Now().After(val.exp) {
		ok = false
	}
	return
}

func (c *allowCache) set(key string, val allowEntry) {
	c.mu.Lock()
	c.data[key] = val
	c.mu.Unlock()
}

var (
	cpsApproveRepo storage.CPSActionApproveIndexRepository
	cpsGuardCache  = &allowCache{data: make(map[string]allowEntry), ttl: 5 * time.Minute}
	guardLogger    utils.Logger
)

// InitCPSActionGuard configures repository and cache TTL for the CPS action guard
func InitCPSActionGuard(repo storage.CPSActionApproveIndexRepository, ttl time.Duration, logger utils.Logger) {
	cpsApproveRepo = repo
	if ttl > 0 {
		cpsGuardCache.ttl = ttl
	}
	guardLogger = logger
}

// RequireCPSAction returns a middleware that allows the request only if the current
// role_id is authorized for the provided actionName (checked in cps_action_approver_index).
func RequireCPSAction(actionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cpsApproveRepo == nil {
				// repository not configured; deny to be safe
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			roleID, _ := r.Context().Value(constants.ContextKey("role_id")).(string)
			if roleID == "" {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			action := strings.ToUpper(strings.TrimSpace(actionName))
			key := roleID + ":" + action

			if ent, ok := cpsGuardCache.get(key); ok {
				if ent.allow {
					next.ServeHTTP(w, r)
					return
				}
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			allowed, err := cpsApproveRepo.ExistsByRoleAndAction(r.Context(), roleID, action)
			if err != nil {
				if guardLogger != nil {
					guardLogger.Errorf("action guard lookup failed: %v", err)
				}
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			// cache result
			cpsGuardCache.set(key, allowEntry{allow: allowed, exp: time.Now().Add(cpsGuardCache.ttl)})
			if !allowed {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
