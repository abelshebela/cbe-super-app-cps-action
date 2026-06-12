package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleDelegation struct {
	ID bson.ObjectID `json:"id" bson:"_id,omitempty"`

	DelegatedUserID                     string `json:"delegated_user_id" bson:"delegated_user_id"`
	DelegatedUserUserCode               string `json:"delegated_user_user_code" bson:"delegated_user_user_code"`
	DelegatedUserFullName               string `json:"delegated_user_full_name" bson:"delegated_user_full_name"`
	DelegatedUserUserType               string `json:"delegated_user_user_type" bson:"delegated_user_user_type"`
	DelegatedUserDepartmentOrBranch     string `json:"delegated_user_department_or_branch" bson:"delegated_user_department_or_branch"`
	DelegatedUserDepartmentOrBranchName string `json:"delegated_user_department_or_branch_name" bson:"delegated_user_department_or_branch_name"`
	DelegatedUserJobTitle               string `json:"delegated_user_job_title" bson:"delegated_user_job_title"`
	DelegatedUserExistingRole           string `json:"delegated_user_existing_role" bson:"delegated_user_existing_role"`
	DelegatedUserExistingRoleName       string `json:"delegated_user_existing_role_name" bson:"delegated_user_existing_role_name"`
	DelegationType                      string `json:"delegation_type" bson:"delegation_type"`
	DelegatedUserPhoneNumber            string `json:"delegated_user_phone_number" bson:"delegated_user_phone_number"`
	DelegatedUserEmail                  string `json:"delegated_user_email" bson:"delegated_user_email"`

	// taken as value with out modification
	DelegatorUserID                 string `json:"delegator_user_id" bson:"delegator_user_id"`
	DelegatorUserFullName           string `json:"delegator_user_full_name" bson:"delegator_user_full_name"`
	DelegatorUserUserType           string `json:"delegator_user_user_type" bson:"delegator_user_user_type"`
	DelegatorUserDepartmentOrBranch string `json:"delegator_user_department_or_branch" bson:"delegator_user_department_or_branch"`
	DelegatorUserJobTitle           string `json:"delegator_user_job_title" bson:"delegator_user_job_title"`
	DelegatorUserRole               string `json:"delegator_user_role" bson:"delegator_user_role"`
	DelegatorUserRoleName           string `json:"delegator_user_role_name" bson:"delegator_user_role_name"`

	NewRoleID                 string `json:"new_role_id" bson:"new_role_id"`
	NewRoleIDName             string `json:"new_role_id_name" bson:"new_role_id_name"`
	NewDepartmentOrBranch     string `json:"new_department_or_branch" bson:"new_department_or_branch"`
	NewDepartmentOrBranchName string `json:"new_department_or_branch_name" bson:"new_department_or_branch_name"`

	Enable                   bool      `json:"enable" bson:"enable"`
	StartAt                  time.Time `json:"start_at" bson:"start_at"`
	EndAt                    time.Time `json:"end_at" bson:"end_at"`
	Reason                   string    `json:"reason" bson:"reason"`
	CreatedAt                time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" bson:"updated_at"`
	RevokeExistingDelegation bool      `json:"revoke_existing_delegation" bson:"revoke_existing_delegation"`
}
