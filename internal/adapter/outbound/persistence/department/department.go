// Package department provides persistence logic for department and CPS action entities.
package department

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repository "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"time"

	"errors"

	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cpsconstants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
)

type DepartmentPersistence struct {
	departmentdal dal.MongoDal[entities.Department, entities.Department]
	cpsdal        dal.MongoDal[cpsactions.CPSAction, cpsactions.CPSAction]
	modelCpsdal   dal.MongoDal[model.CPSAction, model.CPSAction]
	timeout       time.Duration
	logger        utils.Logger
}

type CreateDepartmentRequest struct {
	Department       string          `json:"department"`
	PortalCards      []string        `json:"portal_cards"`
	PermissionGroups []bson.ObjectID `json:"permission_groups"`
}

var _ repository.DepartmentRepository = (*DepartmentPersistence)(nil)
var _ repository.CPSActionRepository = (*DepartmentPersistence)(nil)

func InitDepartment(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *DepartmentPersistence {
	departmentdal := dal.NewMongoDal[entities.Department, entities.Department](client, dbName, "department")
	cpsdal := dal.NewMongoDal[cpsactions.CPSAction, cpsactions.CPSAction](client, dbName, "cps_actions")
	modelCpsdal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, "cps_actions")
	return &DepartmentPersistence{
		departmentdal: departmentdal,
		cpsdal:        cpsdal,
		modelCpsdal:   modelCpsdal,
		timeout:       timeout,
		logger:        logger,
	}
}
func (r *DepartmentPersistence) CheckRequestExists(ctx context.Context, action cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	res, err := r.cpsdal.FindOne(ctx, bson.M{
		"request_action": action.RequestAction,
		"action_status":  action.ActionStatus,
		"maker_id":       action.MakerID,
	}, bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBInsertFailed)
	}
	return res, nil
}

func (r *DepartmentPersistence) CheckDepartmentExists(ctx context.Context, dept string) (bool, error) {
	_, err := r.departmentdal.FindOne(ctx, bson.M{"department": dept}, bson.M{})
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}

	return true, nil
}

func (r *DepartmentPersistence) CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid department id: %s", id)
	}
	_, err = r.departmentdal.FindOne(ctx, bson.M{"_id": objID}, bson.M{})
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}
	return true, nil
}

func (r *DepartmentPersistence) CreateCPSAction(ctx context.Context, department string, portalCards []string, permissionGroups []string, action cpsactions.CPSAction) error {
	modelAction := model.CPSAction{
		ActionCode:         action.ActionCode,
		UniqueId:           action.UniqueID,
		MakerID:            action.MakerID,
		MakerName:          action.MakerName,
		MakerPhoneNumber:   action.MakerPhoneNumber,
		CheckerID:          action.CheckerID,
		CheckerName:        action.CheckerName,
		CheckerPhoneNumber: action.CheckerPhoneNumber,
		Department:         action.Department,
		RejectionReason:    action.RejectionReason,
		PreviousAction:     action.PreviousAction,
		CurrentAction:      action.CurrentAction,
		ActionStatus:       string(action.ActionStatus),
		ActionType:         string(action.ActionType),
		RequestAction:      string(action.RequestAction),
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
		MakerActionTime:    action.MakerActionTime,
		CheckerActionTime:  &action.CheckerActionTime,
	}
	if len(permissionGroups) > 0 {
		modelAction.CurrentAction = map[string]interface{}{
			"department":        department,
			"portal_cards":      portalCards,
			"permission_groups": permissionGroups,
		}
	}
	_, err := r.modelCpsdal.InsertOne(ctx, modelAction)
	return err
}

func (r *DepartmentPersistence) CreateDepartment(ctx context.Context, dept entities.Department) error {

	_, err := r.departmentdal.InsertOne(ctx, dept)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBInsertFailed)
	}
	return nil
}

func (r *DepartmentPersistence) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*cpsactions.CPSAction, error) {

	filter := bson.M{
		"action_code":   actionCode,
		"action_status": cpsconstants.ActionPending,
	}
	var action *cpsactions.CPSAction
	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf(error_codes.ActionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}

	if action == nil {
		return nil, fmt.Errorf(error_codes.ActionNotFound)
	}

	return action, nil
}

func (r *DepartmentPersistence) UpdateDepartment(ctx context.Context, id string, department string, portalCards []string, permissionGroups []string) (*entities.Department, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid department id: %s", id)
	}
	var groupIDs []bson.ObjectID
	for _, idStr := range permissionGroups {
		obj, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid permission group id: %s", idStr)
		}
		groupIDs = append(groupIDs, obj)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"department":        department,
		"portal_cards":      portalCards,
		"permission_groups": groupIDs,
		"last_modified":     time.Now(),
	}
	data, err := r.departmentdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return &data, nil
}
func (r *DepartmentPersistence) FindByActionCode(ctx context.Context, code string) (*cpsactions.CPSAction, error) {
	filter := bson.M{"action_code": code}
	var action *cpsactions.CPSAction

	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf(error_codes.ActionNotFound)
	}

	return action, nil
}

func (r *DepartmentPersistence) UpdateActionStatus(ctx context.Context, actionCode string, status string) (*cpsactions.CPSAction, error) {

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"action_status": status, "last_modified_at": time.Now()}

	data, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return &data, nil
}

func (r *DepartmentPersistence) ApproveActionRequest(ctx context.Context, actionCode string, action cpsactions.CPSAction) error {

	// filter := bson.M{"action_code": actionCode}
	// update := bson.M{
	// 	"checker_name":         action.CheckerName,
	// 	"checker_id":           action.CheckerID,
	// 	"checker_phone_number": action.CheckerPhoneNumber,
	// 	"action_status":        entities.ActionApproved,
	// 	"checker_action_time":  time.Now(),
	// }

	// _, err := r.cpsdal.UpdateOne(ctx, filter, update)

	// if err != nil {
	// 	return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	// }

	// Extract department data from CurrentAction map
	// currentActionMap, ok := actionData.CurrentAction.(map[string]interface{})
	// if !ok {
	// 	return fmt.Errorf("invalid current action data format")
	// }

	// departmentName, ok := currentActionMap["department"].(string)
	// if !ok {
	// 	return fmt.Errorf("invalid department name in current action")
	// }

	// portalCardsInterface, ok := currentActionMap["portal_cards"].([]interface{})
	// if !ok {
	// 	return fmt.Errorf("invalid portal cards in current action")
	// }

	// // Convert portal cards to string slice
	// var portalCards []string
	// for _, card := range portalCardsInterface {
	// 	if cardStr, ok := card.(string); ok {
	// 		portalCards = append(portalCards, cardStr)
	// 	}
	// }

	// data := entities.Department{
	// 	Department:     departmentName,
	// 	DepartmentCode: utils.RandomGenerator(20),
	// 	PortalCards:    portalCards,
	// 	Enabled:        false,
	// 	CreatedAt:      time.Now().UTC(),
	// }

	// _, err = r.departmentdal.InsertOne(ctx, data)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (r *DepartmentPersistence) RejectActionRequest(ctx context.Context, actionCode string, action cpsactions.CPSAction) error {

	filter := bson.M{"action_code": actionCode}
	update := bson.M{
		"checker_name":         action.CheckerName,
		"checker_id":           action.CheckerID,
		"checker_phone_number": action.CheckerPhoneNumber,
		"action_status":        cpsconstants.ActionRejected,
		"checker_action_time":  time.Now(),
	}

	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return nil
}

// CreateDepartmentUpdateCPSAction creates a CPS action for department updates
func (r *DepartmentPersistence) CreateDepartmentUpdateCPSAction(ctx context.Context, req cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	// Check for pending update action
	filter := bson.M{
		"maker_phone_number": req.MakerPhoneNumber,
		"action_status":      cpsconstants.ActionPending,
		"department":         req.Department,
		"request_action":     string(req.RequestAction),
	}

	projection := bson.M{}
	existing, err := r.cpsdal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		r.logger.Errorf("failed to get cps action: %v", err)
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("pending department update action present")
	}

	// Set up the CPS action for department update using the model.CPSAction struct (like CreateCPSAction)
	modelAction := model.CPSAction{
		ActionCode:         req.ActionCode,
		UniqueId:           req.UniqueID,
		MakerID:            req.MakerID,
		MakerName:          req.MakerName,
		MakerPhoneNumber:   req.MakerPhoneNumber,
		CheckerID:          req.CheckerID,
		CheckerName:        req.CheckerName,
		CheckerPhoneNumber: req.CheckerPhoneNumber,
		Department:         req.Department,
		RejectionReason:    req.RejectionReason,
		PreviousAction:     req.PreviousAction,
		CurrentAction:      req.CurrentAction,
		ActionStatus:       string(req.ActionStatus),
		ActionType:         string(req.ActionType),
		RequestAction:      string(req.RequestAction),
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
		MakerActionTime:    req.MakerActionTime,
		CheckerActionTime:  &req.CheckerActionTime,
	}
	_, err = r.modelCpsdal.InsertOne(ctx, modelAction)
	if err != nil {
		r.logger.Errorf("failed to create cps action: %v", err)
		return nil, err
	}
	return &req, nil
}

func (r *DepartmentPersistence) ApproveDepartmentUpdate(ctx context.Context, cpsAction cpsactions.CPSAction) error {
	if deptData, ok := cpsAction.CurrentAction.(map[string]interface{}); ok {
		department, _ := deptData["department"].(string)
		portalCardsInterface, _ := deptData["portal_cards"].([]interface{})
		permissionGroupsInterface, _ := deptData["permission_groups"].([]interface{})

		var portalCards []string
		for _, card := range portalCardsInterface {
			if cardStr, ok := card.(string); ok {
				portalCards = append(portalCards, cardStr)
			}
		}
		var permissionGroups []string
		for _, group := range permissionGroupsInterface {
			if groupStr, ok := group.(string); ok {
				permissionGroups = append(permissionGroups, groupStr)
			}
		}

		switch cpsAction.ActionType {
		case cpsconstants.ActionCreate:
			// Add 'DEP' prefix to the generated department code
			var permissionGroupIDs []bson.ObjectID
			for _, groupStr := range permissionGroups {
				if objID, err := bson.ObjectIDFromHex(groupStr); err == nil {
					permissionGroupIDs = append(permissionGroupIDs, objID)
				}
			}
			data := entities.Department{
				Department:       department,
				DepartmentCode:   "DEP" + utils.RandomGenerator(20),
				PortalCards:      portalCards,
				PermissionGroups: permissionGroupIDs,
				Enabled:          true,
				CreatedAt:        time.Now().UTC(),
				LastModified:     time.Now().UTC(),
			}
			_, err := r.departmentdal.InsertOne(ctx, data)
			if err != nil {
				r.logger.Errorf("failed to create department: %v", err)
				return err
			}
		case cpsconstants.ActionUpdate, cpsconstants.ActionDelete:
			departmentCode, ok := deptData["department_id"].(string)
			if !ok {
				return fmt.Errorf("invalid department_code in current action")
			}
			if cpsAction.ActionType == cpsconstants.ActionUpdate {
				// For update actions, update the existing department
				var permissionGroups []string
				if groups, ok := deptData["permission_groups"]; ok {
					switch v := groups.(type) {
					case []interface{}:
						for _, g := range v {
							if str, ok := g.(string); ok {
								permissionGroups = append(permissionGroups, str)
							}
						}
					case []string:
						permissionGroups = v
					}
				}
				_, err := r.UpdateDepartment(ctx, departmentCode, department, portalCards, permissionGroups)
				if err != nil {
					r.logger.Errorf("failed to update department: %v", err)
					return err
				}
			} else {
				// For delete actions, mark the department as deleted
				deleteFilter := bson.M{"department_id": departmentCode}
				deleteUpdate := bson.M{
					"enabled":       false,
					"deleted_at":    time.Now().UTC(),
					"last_modified": time.Now().UTC(),
				}
				_, err := r.departmentdal.UpdateOne(ctx, deleteFilter, deleteUpdate)
				if err != nil {
					r.logger.Errorf("failed to delete department: %v", err)
					return err
				}
			}
		default:
			return fmt.Errorf("unsupported action type: %s", cpsAction.ActionType)
		}
	}
	return nil
}

// RejectDepartmentUpdate rejects a department update CPS action
func (r *DepartmentPersistence) RejectDepartmentUpdate(ctx context.Context, cpsAction cpsactions.CPSAction) error {
	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"action_status": cpsconstants.ActionPending,
	}
	update := bson.M{
		"action_status":        cpsconstants.ActionRejected,
		"rejection_reason":     cpsAction.RejectionReason,
		"checker_id":           cpsAction.CheckerID,
		"checker_name":         cpsAction.CheckerName,
		"checker_phone_number": cpsAction.CheckerPhoneNumber,
		"checker_action_time":  time.Now(),
		"last_modified_at":     time.Now(),
	}
	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to reject cps action: %v", err)
		return err
	}
	return nil
}

func (r *DepartmentPersistence) GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	// Apply search if provided
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"department": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"department_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	// Calculate pagination
	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	// Get total count
	totalDocs, err := r.departmentdal.TotalCount(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// Get paginated data
	departments, err := r.departmentdal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departments: %w", err)
	}

	// Build pagination metadata using utility function
	meta := common_util.BuildPaginationMeta(totalDocs, filterParams.Page, filterParams.PerPage)

	if departments == nil {
		departments = []*entities.Department{}
	}

	return &common_util.PaginatedResponse[[]*entities.Department]{
		Data: departments,
		Meta: meta,
	}, nil
}

func (r *DepartmentPersistence) GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("invalid department id provided: %s, error: %v", id, err)
		return nil, fmt.Errorf("INVALID_ID")
	}
	filter := bson.M{"_id": objectID, "is_deleted": false}
	department, err := r.departmentdal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Errorf("department not found, id: %s", id)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		r.logger.Errorf("failed to get department, id: %s, error: %v", id, err)
		return nil, fmt.Errorf("FAILED_TO_GET_DEPARTMENT")
	}
	return department, nil
}
