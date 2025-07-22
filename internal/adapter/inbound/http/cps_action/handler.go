package cpsaction

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_action"
	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	inboundCPS "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/cps_actions"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type cpsActionAdapter struct {
	cpsActionApplication app.CPSActionApplication
	logger               utils.Logger
}

func InitCPSActionAdapter(cpsActionApplication app.CPSActionApplication, logger utils.Logger) inboundCPS.CPSActionAdapter {
	return &cpsActionAdapter{
		logger:               logger,
		cpsActionApplication: cpsActionApplication,
	}
}

func (a *cpsActionAdapter) parseUserContext(r *http.Request) (ctx_util.UserContext, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return ctx_util.UserContext{}, fmt.Errorf(common_util.IncompleteUserInfo)
	}

	return userContext, nil
}

func (a *cpsActionAdapter) createCPSUserForAuthorize(r *http.Request) (*model.AuthorizeCPSAction, error) {
	userContext, err := a.parseUserContext(r)
	if err != nil {
		return nil, err
	}
	return &model.AuthorizeCPSAction{
		CheckerUser: model.User{
			UserCode:    userContext.UserCode,
			FullName:    userContext.FullName,
			PhoneNumber: userContext.PhoneNumber,
		},
		Department: userContext.Department,
	}, nil
}

func (a *cpsActionAdapter) getCPSAction(w http.ResponseWriter, r *http.Request, paramKey string, fetchFunc func(ctx context.Context, key string) (*model.CPSAction, error)) {
	id, ok := common_util.GetParam(r, paramKey)
	if !ok {
		a.logger.Errorf("missing or invalid parameter '%s'", paramKey)
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsAction, err := fetchFunc(r.Context(), id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "CPS Actions retrieved successfully")
}

func (a *cpsActionAdapter) parseCPSRequest(w http.ResponseWriter, r *http.Request) (*model.AuthorizeCPSAction, bool) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return nil, false
	}

	cpsReq, err := a.createCPSUserForAuthorize(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return nil, false
	}
	cpsReq.ActionCode = actionCode
	return cpsReq, true
}

func (a *cpsActionAdapter) ApproveCPSAction(w http.ResponseWriter, r *http.Request) {
	cpsReq, ok := a.parseCPSRequest(w, r)
	if !ok {
		return
	}

	cpsAction, err := a.cpsActionApplication.ApproveCPSAction(r.Context(), cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "CPS Action authorized successfully")
}

func (a *cpsActionAdapter) RejectCPSAction(w http.ResponseWriter, r *http.Request) {
	cpsReq, ok := a.parseCPSRequest(w, r)
	if !ok {
		return
	}

	var rejectPayload model.AuthorizeCPSAction
	if err := json.NewDecoder(r.Body).Decode(&rejectPayload); err != nil {
		a.logger.Errorf("failed to decode rejection payload: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidReq, 0, nil)
		return
	}
	cpsReq.RejectionReason = rejectPayload.RejectionReason

	cpsAction, err := a.cpsActionApplication.RejectCPSAction(r.Context(), cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "CPS Action rejected successfully")
}

func (a *cpsActionAdapter) GetCPSActionsByDepartment(w http.ResponseWriter, r *http.Request) {
	user, err := a.parseUserContext(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	filterParams, status := ExtractStatusFilterParams(r)
	cpsActions, err := a.cpsActionApplication.GetCPSActionsByDepartment(r.Context(), user.Department, status, filterParams)

	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsActions, "CPS actions retrieved successfully")
}

func (a *cpsActionAdapter) GetCPSActionByID(w http.ResponseWriter, r *http.Request) {
	a.getCPSAction(w, r, "id", a.cpsActionApplication.GetCPSActionByID)
}

func (a *cpsActionAdapter) GetCPSActionByActionCode(w http.ResponseWriter, r *http.Request) {
	a.getCPSAction(w, r, "action_code", a.cpsActionApplication.GetCPSActionByActionCode)
}

func ExtractStatusFilterParams(r *http.Request) (*constant.Filter, string) {
	query := r.URL.Query()

	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}

	perPage := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt > 0 {
		perPage = perPageInt
	}

	status := query.Get("status")

	return &constant.Filter{
		Page:    page,
		PerPage: perPage,
		Search:  query.Get("search"),
		Filters: query.Get("filter"),
	}, status
}
