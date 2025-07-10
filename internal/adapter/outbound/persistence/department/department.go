// Package department provides persistence logic for department and CPS action entities.
package department

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repository "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"time"
)

type DepartmentPersistence struct {
	departmentdal dal.MongoDal[entities.Department, entities.Department]
	cpsdal        dal.MongoDal[entities.CPSAction, entities.CPSAction]
	timeout       time.Duration
	logger        utils.Logger
}

type CreateDepartmentRequest struct {
	Department  string   `json:"department"`
	PortalCards []string `json:"portal_cards"`
}

var _ repository.DepartmentRepository = (*DepartmentPersistence)(nil)
var _ repository.CPSActionRepository = (*DepartmentPersistence)(nil)

func InitDepartment(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *DepartmentPersistence {
	departmentdal := dal.NewMongoDal[entities.Department, entities.Department](client, dbName, "department")
	cpsdal := dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, "cps_actions")
	return &DepartmentPersistence{
		departmentdal: departmentdal,
		cpsdal:        cpsdal,
		timeout:       timeout,
		logger:        logger,
	}
}
func (r *DepartmentPersistence) CheckRequestExists(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error) {
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

func (r DepartmentPersistence) CreateCPSAction(ctx context.Context, department string, portalCards []string, action entities.CPSAction) error {

	action.CreatedAt = time.Now()
	action.LastModifiedAt = time.Now()
	_, err := r.cpsdal.InsertOne(ctx, action)
	return err
}

func (r *DepartmentPersistence) CreateDepartment(ctx context.Context, dept entities.Department) error {

	_, err := r.departmentdal.InsertOne(ctx, dept)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBInsertFailed)
	}
	return nil
}

func (r *DepartmentPersistence) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error) {

	filter := bson.M{
		"action_code":   actionCode,
		"action_status": entities.ActionPending,
	}
	var action *entities.CPSAction
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

func (r *DepartmentPersistence) UpdateDepartment(ctx context.Context, code string, department string, portalCards []string) (*entities.Department, error) {

	filter := bson.M{"department_code": code}
	update := bson.M{
		"department":   department,
		"portal_cards": portalCards,
	}
	data, err := r.departmentdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return &data, nil
}
func (r *DepartmentPersistence) FindByActionCode(ctx context.Context, code string) (*entities.CPSAction, error) {
	filter := bson.M{"action_code": code}
	var action *entities.CPSAction

	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf(error_codes.ActionNotFound)
	}

	return action, nil
}

func (r *DepartmentPersistence) UpdateActionStatus(ctx context.Context, actionCode string, status string) (*entities.CPSAction, error) {

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"action_status": status, "last_modified_at": time.Now()}

	data, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return &data, nil
}

func (r *DepartmentPersistence) ApproveActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error {

	filter := bson.M{"action_code": actionCode}
	update := bson.M{
		"checker_name":         action.CheckerName,
		"checker_id":           action.CheckerID,
		"checker_phone_number": action.CheckerPhoneNumber,
		"action_status":        entities.ActionApproved,
		"checker_action_time":  time.Now(),
	}

	_, err := r.cpsdal.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}

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

func (r *DepartmentPersistence) RejectActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error {

	filter := bson.M{"action_code": actionCode}
	update := bson.M{
		"checker_name":         action.CheckerName,
		"checker_id":           action.CheckerID,
		"checker_phone_number": action.CheckerPhoneNumber,
		"action_status":        entities.ActionRejected,
		"checker_action_time":  time.Now(),
	}

	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return nil
}

// CreateDepartmentUpdateCPSAction creates a CPS action for department updates
func (r *DepartmentPersistence) CreateDepartmentUpdateCPSAction(ctx context.Context, req entities.CPSAction) (*entities.CPSAction, error) {
	// Check for pending update action
	filter := bson.M{
		"maker_phone_number": req.MakerPhoneNumber,
		"action_status":      entities.ActionPending,
		"department":         req.Department,
		"request_action":     entities.RequestDepartment,
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

	// Set up the CPS action for department update
	req.ActionStatus = entities.ActionPending
	req.RequestAction = entities.RequestDepartment
	req.ActionType = entities.ActionUpdate
	req.MakerActionTime = time.Now()
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()

	cpsAction, err := r.cpsdal.InsertOne(ctx, req)
	if err != nil {
		r.logger.Errorf("failed to create cps action: %v", err)
		return nil, err
	}
	return &cpsAction, nil
}

// ApproveDepartmentUpdate approves a department update CPS action
func (r *DepartmentPersistence) ApproveDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"action_status": entities.ActionPending,
	}
	update := bson.M{
		"action_status":        entities.ActionApproved,
		"checker_id":           cpsAction.CheckerID,
		"checker_name":         cpsAction.CheckerName,
		"checker_phone_number": cpsAction.CheckerPhoneNumber,
		"checker_action_time":  time.Now(),
		"last_modified_at":     time.Now(),
	}
	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to approve cps action: %v", err)
		return err
	}

	// Extract department data from current action and update the department
	if deptData, ok := cpsAction.CurrentAction.(map[string]interface{}); ok {
		departmentCode, ok := deptData["department_code"].(string)
		if !ok {
			return fmt.Errorf("invalid department_code in current action")
		}

		department, _ := deptData["department"].(string)
		portalCardsInterface, _ := deptData["portal_cards"].([]interface{})

		// Convert portal cards to string slice
		var portalCards []string
		for _, card := range portalCardsInterface {
			if cardStr, ok := card.(string); ok {
				portalCards = append(portalCards, cardStr)
			}
		}

		// Handle different action types
		switch cpsAction.ActionType {
		case entities.ActionCreate:
			// For create actions, insert a new department
			data := entities.Department{
				Department:     department,
				DepartmentCode: utils.RandomGenerator(20),
				PortalCards:    portalCards,
				Enabled:        true,
				CreatedAt:      time.Now().UTC(),
				LastModified:   time.Now().UTC(),
			}
			_, err = r.departmentdal.InsertOne(ctx, data)
			if err != nil {
				r.logger.Errorf("failed to create department: %v", err)
				return err
			}
		case entities.ActionUpdate:
			// For update actions, update the existing department
			_, err = r.UpdateDepartment(ctx, departmentCode, department, portalCards)
			if err != nil {
				r.logger.Errorf("failed to update department: %v", err)
				return err
			}
		case entities.ActionDelete:
			// For delete actions, mark the department as deleted
			deleteFilter := bson.M{"department_code": departmentCode}
			deleteUpdate := bson.M{
				"enabled":       false,
				"deleted_at":    time.Now().UTC(),
				"last_modified": time.Now().UTC(),
			}
			_, err = r.departmentdal.UpdateOne(ctx, deleteFilter, deleteUpdate)
			if err != nil {
				r.logger.Errorf("failed to delete department: %v", err)
				return err
			}
		default:
			return fmt.Errorf("unsupported action type: %s", cpsAction.ActionType)
		}
	}

	return nil
}

// RejectDepartmentUpdate rejects a department update CPS action
func (r *DepartmentPersistence) RejectDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"action_status": entities.ActionPending,
	}
	update := bson.M{
		"action_status":        entities.ActionRejected,
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

func (r *DepartmentPersistence) GetAllDepartments(ctx context.Context) ([]entities.Department, error) {
	filter := bson.M{}
	departments, err := r.departmentdal.Find(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return departments, nil
}
