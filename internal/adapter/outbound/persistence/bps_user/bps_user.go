package persistence

import (
	"context"
	"fmt"
	"strconv"
	"time"

	bsonv2 "go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bps_user"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BpsPersistence struct {
	bpsDal  dal.MongoDal[model.BPSUser, model.BPSUser]
	cpsDal  dal.MongoDal[model.CPSAction, model.CPSAction]
	logger  utils.Logger
	timeout time.Duration
	client  *mongo.Client
}

func NewBpsPersistence(client *mongo.Client, database string, timeout time.Duration, logger utils.Logger) *BpsPersistence {
	return &BpsPersistence{
		bpsDal:  dal.NewMongoDal[model.BPSUser, model.BPSUser](client, database, "branch_user"),
		cpsDal:  dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, "cps_actions"),
		client:  client,
		timeout: timeout,
		logger:  logger,
	}
}

// GetBPSUserByUserCode retrieves a BPS user by user code
func (o *BpsPersistence) GetBPSUserByUserCode(ctx context.Context, userCode string) (*bps_user.BPSUser, error) {
	filter := bsonv2.M{"user_code": userCode}

	user, err := o.bpsDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("DATABASE_ERROR_FINDING_USER")
	}

	// Convert model.BPSUser to domain.BPSUser
	domainUser := &bps_user.BPSUser{
		ID:          user.ID,
		UserCode:    user.UserCode,
		FullName:    user.FullName,
		Username:    user.UserName,
		PhoneNumber: user.PhoneNumber,
		BranchCode:  user.BranchCode,
		BranchName:  user.BranchName,
		HomeBranch:  user.HomeBranch,
		Role:        user.Role,
		Realm:       string(user.Realm),
		Enabled:     user.Enabled,
		IsDeleted:   user.IsDeleted,
	}

	return domainUser, nil
}

func (o *BpsPersistence) GetAllBPSUsers(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*bps_user.BPSUser], error) {
	filter := bsonv2.M{}

	if filterParams.Search != "" {
		searchRegex := bsonv2.Regex{Pattern: filterParams.Search, Options: "i"}
		filter["$or"] = []bsonv2.M{
			{"user_code": searchRegex},
			{"full_name": searchRegex},
			{"phone_number": searchRegex},
			{"username": searchRegex},
			{"role": searchRegex},
			{"branch_name": searchRegex},
		}
	}

	if filterParams.Filters != nil {
		allowedKeys := []string{
			"user_code", "full_name", "username", "phone_number",
			"branch_code", "branch_name", "home_branch", "role",
			"realm", "enabled", "is_deleted",
		}

		handlers := map[string]func(interface{}) interface{}{
			"enabled": func(value interface{}) interface{} {
				if str, ok := value.(string); ok {
					if parsed, err := strconv.ParseBool(str); err == nil {
						return parsed
					}
				}
				return value
			},
			"is_deleted": func(value interface{}) interface{} {
				if str, ok := value.(string); ok {
					if parsed, err := strconv.ParseBool(str); err == nil {
						return parsed
					}
				}
				return value
			},
		}

		enhancedFilter := common_util.BuildMongoFilterWithHandlers(filterParams.Filters, allowedKeys, handlers)
		for key, value := range enhancedFilter {
			if _, exists := filter[key]; !exists {
				filter[key] = value
			}
		}
	}

	skip := int64((filterParams.Page - 1) * filterParams.PerPage)
	limit := int64(filterParams.PerPage)

	users, err := o.bpsDal.FindAllWithPagination(ctx, filter, bsonv2.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	domainUsers := make([]*bps_user.BPSUser, 0, len(users))
	for i := range users {
		u := users[i]
		if u == nil {
			continue
		}
		domainUsers = append(domainUsers, &bps_user.BPSUser{
			ID:          u.ID,
			UserCode:    u.UserCode,
			FullName:    u.FullName,
			Username:    u.UserName,
			PhoneNumber: u.PhoneNumber,
			BranchCode:  u.BranchCode,
			BranchName:  u.BranchName,
			HomeBranch:  u.HomeBranch,
			Role:        u.Role,
			Realm:       string(u.Realm),
			Enabled:     u.Enabled,
			IsDeleted:   u.IsDeleted,
		})
	}

	totalDocs, err := o.bpsDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := common_util.BuildPaginationMeta(totalDocs, filterParams.Page, filterParams.PerPage)

	return &common_util.PaginatedResponse[[]*bps_user.BPSUser]{
		Data: domainUsers,
		Meta: meta,
	}, nil
}

// UpdateBPSUserStatus updates the enabled status of a BPS user
func (o *BpsPersistence) UpdateBPSUserStatus(ctx context.Context, userCode string, enabled bool, updatedAt time.Time) error {
	filter := bsonv2.M{"user_code": userCode}
	update := bsonv2.M{

		"enabled":          enabled,
		"last_modified_at": updatedAt,
	}

	_, err := o.bpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		o.logger.Errorf("failed to update BPS user status: %v", err)
		return fmt.Errorf("FAILED_TO_UPDATE_USER")
	}

	return nil
}

// CheckPendingAction checks if there's a pending action for the given request action and department
func (o *BpsPersistence) CheckPendingAction(ctx context.Context, requestAction string, department string) (bool, error) {
	filter := bsonv2.M{
		"action_status":  "PENDING",
		"request_action": requestAction,
		"department":     department,
	}

	existing, err := o.cpsDal.FindOne(ctx, filter, nil)
	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return existing != nil, nil
}

// EnableBPSUser enables a BPS user
func (o *BpsPersistence) EnableBPSUser(ctx context.Context, userCode string) error {
	return o.UpdateBPSUserStatus(ctx, userCode, true, time.Now())
}

// DisableBPSUser disables a BPS user
func (o *BpsPersistence) DisableBPSUser(ctx context.Context, userCode string) error {
	return o.UpdateBPSUserStatus(ctx, userCode, false, time.Now())
}

// CreateCPSAction creates a new CPS action
func (o *BpsPersistence) CreateCPSAction(ctx context.Context, action model.CPSAction) (*model.CPSAction, error) {
	action.ID = bsonv2.NewObjectID()
	action.CreatedAt = time.Now()
	action.LastModifiedAt = time.Now()

	_, err := o.cpsDal.InsertOne(ctx, action)
	if err != nil {
		o.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return &action, nil
}

// GetCPSActionByID retrieves a CPS action by ID
func (o *BpsPersistence) GetCPSActionByID(ctx context.Context, actionID string) (*model.CPSAction, error) {
	filter := bsonv2.M{"action_code": actionID}
	action, err := o.cpsDal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return action, nil
}

// AuthorizeBPSUserEnable authorizes enabling a BPS user
func (o *BpsPersistence) AuthorizeBPSUserEnable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	// Extract user code from the action
	var currentAction struct {
		BPSUser struct {
			UserCode string `json:"user_code"`
		} `json:"bps_user"`
	}

	currentActionBytes, err := bsonv2.Marshal(action.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_MARSHAL_CURRENT_ACTION")
	}

	if err := bsonv2.Unmarshal(currentActionBytes, &currentAction); err != nil {
		return nil, fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
	}

	// Update the user status
	if err := o.EnableBPSUser(ctx, currentAction.BPSUser.UserCode); err != nil {
		return nil, err
	}
	return action, nil
}

// AuthorizeBPSUserDisable authorizes disabling a BPS user
func (o *BpsPersistence) AuthorizeBPSUserDisable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	var currentAction struct {
		BPSUser struct {
			UserCode string `json:"user_code"`
		} `json:"bps_user"`
	}

	currentActionBytes, err := bsonv2.Marshal(action.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_MARSHAL_CURRENT_ACTION")
	}

	if err := bsonv2.Unmarshal(currentActionBytes, &currentAction); err != nil {
		return nil, fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
	}

	// Update the user status
	if err := o.DisableBPSUser(ctx, currentAction.BPSUser.UserCode); err != nil {
		return nil, err
	}
	return action, nil
}
