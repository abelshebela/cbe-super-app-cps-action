package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

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
	cpsApproveRepo              storage.CPSActionApproveIndexRepository
	bpsApproveRepo              storage.BPSActionApproveIndexRepository
	roleRepo                    storage.RoleRepository
	clientOrchestrationProducer *kafka.ClientOrchestrationProducer
	cpsGuardCache               = &allowCache{data: make(map[string]allowEntry), ttl: 5 * time.Minute}
	roleEnabledCache            = &allowCache{data: make(map[string]allowEntry), ttl: 5 * time.Minute}
	guardLogger                 utils.Logger
)

func InitCPSActionGuard(cpsRepo storage.CPSActionApproveIndexRepository, bpsRepo storage.BPSActionApproveIndexRepository, rRepo storage.RoleRepository, ttl time.Duration, logger utils.Logger) {
	cpsApproveRepo = cpsRepo
	bpsApproveRepo = bpsRepo
	roleRepo = rRepo
	if ttl > 0 {
		cpsGuardCache.ttl = ttl
		roleEnabledCache.ttl = ttl
	}
	guardLogger = logger
}

func InitClientOrchestrationProducer(producer *kafka.ClientOrchestrationProducer) {
	clientOrchestrationProducer = producer
}

func RequireCPSAction(actionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cpsApproveRepo == nil {
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
			cpsGuardCache.set(key, allowEntry{allow: allowed, exp: time.Now().Add(cpsGuardCache.ttl)})
			if !allowed {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
