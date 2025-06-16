package department

import (
  "context"
  "go.mongodb.org/mongo-driver/v2/bson"
  "go.mongodb.org/mongo-driver/v2/mongo"
  "errors"
  "log"

  "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/department/entities"
  repository "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/department"
 "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

 "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

  "time"
)

type DepartmentPersistence struct {
    departmentdal dal.MongoDal[entities.Department, entities.Department]
    cpsdal    dal.MongoDal[entities.CPSAction, entities.CPSAction]
    timeout   time.Duration
	logger    utils.Logger
}

var _ repository.DepartmentRepository = (*DepartmentPersistence)(nil)
var _ repository.CPSActionRepository = (*DepartmentPersistence)(nil)

func InitDepartment(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *DepartmentPersistence {
    departmentdal := dal.NewMongoDal[entities.Department, entities.Department](client, dbName, "department")
    cpsdal := dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, "cps_actions")
    return &DepartmentPersistence{
        departmentdal: departmentdal,
        cpsdal:       cpsdal,
        timeout:      timeout,
		logger: logger,
    }
}
func (r *DepartmentPersistence) CheckRequestExists(action entities.CPSAction) (bool, error) {
	ctx := context.Background()
 _, err := r.cpsdal.FindOne(ctx, bson.M{
    "request_action": action.RequestAction,
    "action_status":       action.ActionStatus,
    "maker_id": action.MakerID,
}, bson.M{})
if err == mongo.ErrNoDocuments {
    return false, nil
}
return err == nil, err

}

func (r *DepartmentPersistence) CheckDepartmentExists(dept string) (bool, error) {
  ctx := context.Background()
 _, err := r.departmentdal.FindOne(ctx, bson.M{"department": dept}, bson.M{})
if err == mongo.ErrNoDocuments {
    return false, nil
}
return err == nil, err
}

func (r DepartmentPersistence) CreateCPSAction(department string, portalCards []string, action entities.CPSAction) error {
	ctx := context.Background()
	action.CreatedAt = time.Now()
	    action.LastModifiedAt = time.Now()
  _, err := r.cpsdal.InsertOne(ctx, action)
  return err
}

func (r *DepartmentPersistence) CreateDepartment(dept entities.Department) error {
	ctx := context.Background()
	_, err := r.departmentdal.InsertOne(ctx, dept)
	return err
}

func (r *DepartmentPersistence) ValidateActionRequest(actionCode string, userDept string) (*entities.CPSAction, error) {
	ctx := context.Background()
	filter := bson.M{
		"action_code":    actionCode,
		"action_status":   entities.ActionPending,
	}
	var action *entities.CPSAction
	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if action == nil {
		return nil, errors.New("action not found")
	}
	return action, nil
}

func (r *DepartmentPersistence) UpdateDepartment(code string, department string, portalCards []string) error {
	ctx := context.Background()
	filter := bson.M{"department_code": code}
	update := bson.M{
		"department":   department,
		"portal_cards": portalCards,
	}
	_, err := r.departmentdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	// if result.MatchedCount == 0 {
	// 	return errors.New("department not found")
	// }
	return nil
}





func (r *DepartmentPersistence) FindByActionCode(code string) (*entities.CPSAction, error) {
	ctx := context.Background()
	log.Println("Finding action by code:", code)
	filter := bson.M{"action_code": code}
	var action *entities.CPSAction
	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	log.Println("Found action:", action)
	return action, err
}

func (r *DepartmentPersistence) UpdateActionStatus(actionCode string, status string) error {
	ctx := context.Background()
	filter := bson.M{"action_code": actionCode}
	update := bson.M{"action_status": status, "last_modified_at": time.Now()}
	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	// if result.MatchedCount == 0 {
	// 	return errors.New("action not found")
	// }
	return nil
}



func (r *DepartmentPersistence) ApproveActionRequest(actionCode string, user entities.CPSAction) error {			
	
	log.Println("user", user)
	ctx := context.Background()
	filter := bson.M{"action_code": actionCode}
	update := bson.M{
        "checker_name": user.CheckerName,
		"checker_id": user.CheckerID,
		"checker_phone_number": user.CheckerPhoneNumber,
		"action_status": entities.ActionApproved,
		"checker_action_time": time.Now(),
	}
	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}