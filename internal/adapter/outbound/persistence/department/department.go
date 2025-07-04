// Package department provides persistence logic for department and CPS action entities.
package department

import (
	"context"
	"fmt"
	"log"

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
func (r *DepartmentPersistence) CheckRequestExists(ctx context.Context, action entities.CPSAction) (bool, error) {
	_, err := r.cpsdal.FindOne(ctx, bson.M{
		"request_action": action.RequestAction,
		"action_status":  action.ActionStatus,
		"maker_id":       action.MakerID,
	}, bson.M{})
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}
	return true, nil

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
		return nil, fmt.Errorf(error_codes.GeneralDBInsertFailed)
	}

	return action, nil
}

func (r *DepartmentPersistence) UpdateDepartment(ctx context.Context, code string, department string, portalCards []string) error {

	filter := bson.M{"department_code": code}
	update := bson.M{
		"department":   department,
		"portal_cards": portalCards,
	}
	_, err := r.departmentdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return nil
}

func (r *DepartmentPersistence) FindByActionCode(ctx context.Context, code string) (*entities.CPSAction, error) {

	log.Println("Finding action by code:", code)
	filter := bson.M{"action_code": code}
	var action *entities.CPSAction

	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf(error_codes.ActionNotFound)
	}

	log.Println("Found action:", action)
	return action, nil
}

func (r *DepartmentPersistence) UpdateActionStatus(ctx context.Context, actionCode string, status string) error {

	filter := bson.M{"action_code": actionCode}
	update := bson.M{"action_status": status, "last_modified_at": time.Now()}

	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}

	return nil
}

func (r *DepartmentPersistence) ApproveActionRequest(ctx context.Context, actionCode string, user entities.CPSAction) error {

	filter := bson.M{"action_code": actionCode}
	update := bson.M{
		"checker_name":         user.CheckerName,
		"checker_id":           user.CheckerID,
		"checker_phone_number": user.CheckerPhoneNumber,
		"action_status":        entities.ActionApproved,
		"checker_action_time":  time.Now(),
	}

	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return nil
}
