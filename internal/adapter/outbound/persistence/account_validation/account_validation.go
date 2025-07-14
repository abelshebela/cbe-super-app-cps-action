package account_validation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"

	// outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_validation"
	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"

	actions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccountValidationRepo struct {
	client        *mongo.Client
	validationDal dal.MongoDal[model.ValidationRule, model.ValidationRule]
	cpsActionDal  dal.MongoDal[model.CPSAction, model.CPSAction]
	timeout       time.Duration
	logger        utils.Logger
}

//	type CPSActionRepo struct {
//		client   *mongo.Client
//		mongoDal dal.MongoDal[model.ValidationRule, model.ValidationRule,model.CPSAction]
//		timeout  time.Duration
//		logger   utils.Logger
//	}
// var _ outbound.OutboundInfra = (*AccountValidationRepo)(nil)
// var _ account_validation.Repository = (*AccountValidationRepo)(nil)

func InitAccountValidationPersistence(client *mongo.Client, database string, timeout time.Duration, logger utils.Logger) *AccountValidationRepo {
	validationDal := dal.NewMongoDal[model.ValidationRule, model.ValidationRule](client, database, "validation_rule")
	cpsActionDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, "cps_actions")
	return &AccountValidationRepo{
		client:        client,
		validationDal: validationDal,
		cpsActionDal:  cpsActionDal,
		timeout:       timeout,
		logger:        logger,
	}
}

func (r *AccountValidationRepo) GetAccountValidationByID(ctx context.Context, id string) (account_validation.ValidationRule, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if id == "" {
		r.logger.Errorf("invalid validation rule ID: empty")
		return account_validation.ValidationRule{}, errors.New("validation rule ID cannot be empty")
	}

	var objID bson.ObjectID
	var err error
	if objID, err = bson.ObjectIDFromHex(id); err != nil {
		filter := bson.M{"_id": id}
		projection := bson.M{}
		rule, err := r.validationDal.FindOne(ctx, filter, projection)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				r.logger.Errorf("validation rule not found: id=%s", id)
				return account_validation.ValidationRule{}, errors.New("validation rule not found")
			}
			r.logger.Errorf("failed to get validation rule: %v", err)
			return account_validation.ValidationRule{}, errors.New("failed to get validation rule")
		}

		if rule == nil {
			r.logger.Errorf("validation rule not found: id=%s", id)
			return account_validation.ValidationRule{}, errors.New("validation rule not found")
		}

		return account_validation.ValidationRule{
			ID:             rule.ID.Hex(),
			EntityType:     rule.EntityType,
			ValidationFor:  rule.ValidationFor,
			Identifier:     rule.Identifier,
			MinLength:      uint8(rule.MinLength),
			MaxLength:      uint8(rule.MaxLength),
			Enabled:        rule.Enabled,
			IsDeleted:      rule.IsDeleted,
			CreatedAt:      rule.CreatedAt,
			LastModifiedAt: rule.LastModifiedAt,
			ServiceID:      rule.ServiceID,
		}, nil
	}

	// If it is a valid ObjectID, use that
	filter := bson.M{"_id": objID}
	projection := bson.M{}
	rule, err := r.validationDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Errorf("validation rule not found: id=%s", id)
			return account_validation.ValidationRule{}, errors.New("validation rule not found")
		}
		r.logger.Errorf("failed to get validation rule: %v", err)
		return account_validation.ValidationRule{}, errors.New("failed to get validation rule")
	}

	if rule == nil {
		r.logger.Errorf("validation rule not found: id=%s", id)
		return account_validation.ValidationRule{}, errors.New("validation rule not found")
	}

	if rule.MinLength < 0 || rule.MinLength > 255 || rule.MaxLength < 0 || rule.MaxLength > 255 {
		r.logger.Errorf("validation rule length out of uint8 range: id=%s, min_length=%d, max_length=%d", id, rule.MinLength, rule.MaxLength)
		return account_validation.ValidationRule{}, errors.New("validation failed: min length cannot exceed max length")
	}

	return account_validation.ValidationRule{
		ID:             rule.ID.Hex(),
		EntityType:     rule.EntityType,
		ValidationFor:  rule.ValidationFor,
		Identifier:     rule.Identifier,
		MinLength:      uint8(rule.MinLength),
		MaxLength:      uint8(rule.MaxLength),
		Enabled:        rule.Enabled,
		IsDeleted:      rule.IsDeleted,
		CreatedAt:      rule.CreatedAt,
		LastModifiedAt: rule.LastModifiedAt,
		ServiceID:      rule.ServiceID,
	}, nil
}

func (r *AccountValidationRepo) UpdateAccountValidation(ctx context.Context, id string, rule account_validation.ValidationRule) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if id == "" {
		r.logger.Errorf("invalid validation rule ID: empty")
		return errors.New("validation rule ID cannot be empty")
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil || objID.IsZero() {
		r.logger.Errorf("failed to convert id to object id: %v", err)
		return errors.New("validation rule ID cannot be empty")
	}
	update := bson.M{

		"entity_type":      rule.EntityType,
		"validation_for":   rule.ValidationFor,
		"identifier":       rule.Identifier,
		"min_length":       int(rule.MinLength),
		"max_length":       int(rule.MaxLength),
		"enabled":          rule.Enabled,
		"is_deleted":       rule.IsDeleted,
		"service_id":       rule.ServiceID,
		"last_modified_at": rule.LastModifiedAt,
	}

	filter := bson.M{"_id": objID}
	_, err = r.validationDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Errorf("validation rule not found for update: id=%s", id)
			return errors.New("validation rule not found")
		}
		r.logger.Errorf("failed to update validation rule: %v", err)
		return errors.New("failed to update validation rule")
	}

	r.logger.Infof("successfully updated validation rule: id=%s", id)
	return nil
}

func (r *AccountValidationRepo) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]actions.ActionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if uniqueID == "" {
		r.logger.Errorf("invalid unique ID: empty")
		return nil, fmt.Errorf("invalid unique ID provided %w", local_utils.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid unique ID",
		})
	}

	filter := bson.M{
		"unique_id":     uniqueID,
		"action_status": actions.ActionPending,
	}
	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	r.logger.Infof("fetching first pending action with filter: %+v", filter)
	doc, err := r.cpsActionDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return []actions.ActionResponse{}, nil
		}

		r.logger.Errorf("failed to fetch pending action for unique_id=%s: %v", uniqueID, err)
		return nil, fmt.Errorf("FAILED_TO_FETCH: %w", err)
	}

	if doc == nil {
		return []actions.ActionResponse{}, nil

	}

	response := actions.ActionResponse{
		ID:       doc.ID.Hex(),
		ActionId: doc.ActionCode,
	}
	return []actions.ActionResponse{response}, nil
}

func (r *AccountValidationRepo) FetchAllActionsByUniqueID(ctx context.Context, uniqueID string) ([]actions.ActionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if uniqueID == "" {
		r.logger.Errorf("invalid unique ID: empty")
		return nil, fmt.Errorf("invalid unique ID provided %w", local_utils.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid unique ID",
		})
	}

	filter := bson.M{
		"unique_id": uniqueID,
	}
	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	r.logger.Infof("fetching all actions with filter: %+v", filter)
	docs, err := r.cpsActionDal.FindAll(ctx, filter, projection)
	if err != nil {
		r.logger.Errorf("failed to fetch actions for unique_id=%s: %v", uniqueID, err)
		return nil, fmt.Errorf("FAILED_TO_FETCH: %w", err)
	}

	var responses []actions.ActionResponse
	for _, doc := range docs {
		responses = append(responses, actions.ActionResponse{
			ID:       doc.ID.Hex(),
			ActionId: doc.ActionCode,
		})
	}
	return responses, nil
}
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
