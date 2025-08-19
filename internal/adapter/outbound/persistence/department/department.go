// Package department provides persistence logic for department and CPS action entities.
package department

import (
	"context"
	"fmt"
	"time"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DepartmentPersistence struct {
	departmentdal dal.MongoDal[entities.Department, entities.Department]
	logger        utils.Logger
}

func NewDepartmentPersistence(departmentdal dal.MongoDal[entities.Department, entities.Department], logger utils.Logger) *DepartmentPersistence {
	return &DepartmentPersistence{
		departmentdal: departmentdal,
		logger:        logger,
	}
}

func (r *DepartmentPersistence) CheckDepartmentExists(ctx context.Context, department string) (bool, error) {
	filter := bson.M{"department": department}
	projection := bson.M{}
	_, err := r.departmentdal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		r.logger.Errorf("failed to check department exists: %v", err)
		return false, err
	}
	return true, nil
}

func (r *DepartmentPersistence) CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("INVALID_DEPARTMENT_ID_FORMAT")
	}
	filter := bson.M{"_id": objectID}
	projection := bson.M{}
	_, err = r.departmentdal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		r.logger.Errorf("failed to check department exists by ID: %v", err)
		return false, err
	}
	return true, nil
}

func (r *DepartmentPersistence) CreateDepartment(ctx context.Context, dept entities.Department) error {
	fmt.Println("Creating department:", dept.Department)
	_, err := r.departmentdal.InsertOne(ctx, dept)
	if err != nil {
		r.logger.Errorf("failed to create department: %v", err)
		return err
	}
	return nil
}

func (r *DepartmentPersistence) UpdateDepartment(ctx context.Context, id string, updateData map[string]interface{}) (*entities.Department, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("INVALID_DEPARTMENT_ID_FORMAT")
	}

	updateDoc := bson.M{}
	updateDoc["last_modified"] = time.Now()

	if department, exists := updateData["department"]; exists {
		if strVal, ok := department.(string); ok && strVal != "" {
			updateDoc["department"] = strVal
		}
	}

	if portalCards, exists := updateData["portal_cards"]; exists {
		r.logger.Infof("Processing portal_cards in persistence: %+v (type: %T)", portalCards, portalCards)
		var cardIDs []string

		// Handle different types for portal cards
		if strArr, ok := portalCards.([]string); ok && len(strArr) > 0 {
			r.logger.Infof("Portal cards as []string: %+v", strArr)
			cardIDs = strArr
		} else if interfaceArr, ok := portalCards.([]interface{}); ok && len(interfaceArr) > 0 {
			// Convert []interface{} to []string
			r.logger.Infof("Portal cards as []interface{}: %+v", interfaceArr)
			for _, item := range interfaceArr {
				if str, ok := item.(string); ok && str != "" {
					cardIDs = append(cardIDs, str)
				}
			}
			r.logger.Infof("Converted portal cards: %+v", cardIDs)
		}

		// Process the card IDs if we have any
		if len(cardIDs) > 0 {
			updateDoc["portal_cards"] = cardIDs
			r.logger.Infof("Final portal cards: %+v", cardIDs)
		}
	}

	if permissionGroups, exists := updateData["permission_groups"]; exists {
		r.logger.Infof("Processing permission_groups in persistence: %+v (type: %T)", permissionGroups, permissionGroups)
		var groupIDs []string

		// Handle different types for permission groups
		if strArr, ok := permissionGroups.([]string); ok && len(strArr) > 0 {
			r.logger.Infof("Permission groups as []string: %+v", strArr)
			groupIDs = strArr
		} else if interfaceArr, ok := permissionGroups.([]interface{}); ok && len(interfaceArr) > 0 {
			// Convert []interface{} to []string
			r.logger.Infof("Permission groups as []interface{}: %+v", interfaceArr)
			for _, item := range interfaceArr {
				if str, ok := item.(string); ok && str != "" {
					groupIDs = append(groupIDs, str)
				}
			}
			r.logger.Infof("Converted permission groups: %+v", groupIDs)
		}

		// Process the group IDs if we have any
		if len(groupIDs) > 0 {
			var objectIDs []bson.ObjectID
			for _, groupID := range groupIDs {
				oid, err := bson.ObjectIDFromHex(groupID)
				if err != nil {
					return nil, fmt.Errorf("INVALID_PERMISSION_GROUP_ID")
				}
				objectIDs = append(objectIDs, oid)
			}
			updateDoc["permission_groups"] = objectIDs
			r.logger.Infof("Final permission groups ObjectIDs: %+v", objectIDs)
		}
	}
	fmt.Println("data in persitence", updateData)
	if v, ok := updateData["enabled"].(bool); ok {
		updateDoc["enabled"] = v
	}

	filter := bson.M{"_id": objectID}
	fmt.Println("finall   ", updateData)
	updatedDept, err := r.departmentdal.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		r.logger.Errorf("failed to update department: %v", err)
		return nil, err
	}

	return &updatedDept, nil
}

func (r *DepartmentPersistence) GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error) {
	filter := bson.M{}
	projection := bson.M{}

	// Apply filters if provided
	if filterParams != nil {
		if filterParams.Search != "" {
			filter["department"] = bson.M{"$regex": filterParams.Search, "$options": "i"}
		}
	}

	// Get total count
	totalCount, err := r.departmentdal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to count departments: %v", err)
		return nil, err
	}

	// Apply pagination
	skip := int64((filterParams.Page - 1) * filterParams.PerPage)
	limit := int64(filterParams.PerPage)

	// Find departments with pagination
	departments, err := r.departmentdal.FindAllWithPagination(ctx, filter, projection, skip, limit)
	if err != nil {
		r.logger.Errorf("failed to find departments: %v", err)
		return nil, err
	}

	// Convert to domain entities
	var domainDepartments []*entities.Department
	for _, dept := range departments {
		if dept != nil {
			domainDepartments = append(domainDepartments, dept)
		}
	}

	// Build pagination meta
	meta := common_util.BuildPaginationMeta(totalCount, filterParams.Page, filterParams.PerPage)

	return &common_util.PaginatedResponse[[]*entities.Department]{
		Data: domainDepartments,
		Meta: meta,
	}, nil
}

func (r *DepartmentPersistence) GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("INVALID_DEPARTMENT_ID_FORMAT")
	}

	filter := bson.M{"_id": objectID}
	projection := bson.M{}

	dept, err := r.departmentdal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("DEPARTMENT_NOT_FOUND")
		}
		r.logger.Errorf("failed to get department by ID: %v", err)
		return nil, err
	}

	return dept, nil
}
