package department_handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

type DepartmentHandler struct {
	departmentService department.DepartmentService
	logger            utils.Logger
}

func NewDepartmentHTTPHandler(service department.DepartmentService, logger utils.Logger) inbound.DepartmentPortHandler {
	return &DepartmentHandler{
		departmentService: service,
		logger:            logger,
	}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var request department.CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[CreateDepartment] failed to decode request: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           map[string]string{"message": "Invalid JSON payload"},
		}
		resp.SendJSON()
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreateDepartment] validation failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           map[string]string{"message": "Invalid input provided"},
		}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if userID == "" || fullName == "" || phoneNumber == "" || department == "" {
		h.logger.Errorf("[CreateDepartment] incomplete user information")
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusUnauthorized,
			Data:           map[string]string{"message": "Incomplete user information"},
		}
		resp.SendJSON()
		return
	}

	cpsAction := entities.CPSAction{
		MakerID:          userID,
		MakerName:        fullName,
		MakerPhoneNumber: phoneNumber,
		Department:       department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		RequestAction:    entities.RequestDepartment,
	}

	if err := h.departmentService.CreateDepartment(request.Department, request.PortalCards, cpsAction); err != nil {
		h.logger.Errorf("[CreateDepartment] service error: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{
				"message": err.Error(),
			},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[CreateDepartment] request sent successfully by user: %s", userID)
	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data: map[string]string{
			"message": "request sent successfully",
		},
	}
	resp.SendJSON()
}

func (h *DepartmentHandler) ApproveDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if userID == "" || fullName == "" || phoneNumber == "" || department == "" {
		h.logger.Errorf("[ApproveRequest] incomplete user information")
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusUnauthorized,
			Data:           map[string]string{"message": "Incomplete user information"},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[ApproveRequest] user %s is approving action %s", department, actionCode)

	cpsAction, err := h.departmentService.ValidateActionRequest(actionCode, department)
	if err != nil {
		h.logger.Errorf("[ApproveRequest] validation failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           map[string]string{"message": err.Error()},
		}
		resp.SendJSON()
		return
	}
	log.Printf("[ApproveRequest] action validated: %+v", cpsAction)

	var serviceErr error
	switch cpsAction.ActionType {
	case entities.ActionCreate:
		serviceErr = h.approveCreateAction(cpsAction)
	case entities.ActionUpdate:
		serviceErr = h.approveUpdateAction(cpsAction)
	default:
		h.logger.Warnf("[ApproveRequest] invalid action type: %s", cpsAction.ActionType)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           map[string]string{"message": "Invalid action type"},
		}
		resp.SendJSON()
		return
	}

	if serviceErr != nil {
		h.logger.Errorf("[ApproveRequest] service error: %v", serviceErr)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusInternalServerError,
			Data:           map[string]string{"message": serviceErr.Error()},
		}
		resp.SendJSON()
		return
	}

	checker := entities.CPSAction{
		CheckerID:          userID,
		CheckerName:        fullName,
		CheckerPhoneNumber: phoneNumber,
	}

	if err := h.departmentService.ApproveActionRequest(actionCode, checker); err != nil {
		h.logger.Errorf("[ApproveRequest] failed to approve action request: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusInternalServerError,
			Data:           map[string]string{"message": err.Error()},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[ApproveRequest] action %s approved by user %s", actionCode, userID)
	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           map[string]string{"message": "Action approved"},
	}
	resp.SendJSON()
}

func (h *DepartmentHandler) approveCreateAction(action *entities.CPSAction) error {
	var dataMap map[string]interface{}

	bytes, err := json.Marshal(action.CurrentAction)
	if err != nil {
		log.Println("Failed to marshal CurrentAction:", err)
		return errors.New("failed to process action data")
	}

	if err := json.Unmarshal(bytes, &dataMap); err != nil {
		log.Println("Failed to unmarshal CurrentAction into map:", err)
		return errors.New("invalid action data format")
	}

	department, ok := dataMap["department"].(string)
	if !ok {
		log.Println("Missing or invalid 'department'")
		return errors.New("missing or invalid department")
	}

	rawCards, ok := dataMap["portal_cards"].([]interface{})
	if !ok {
		log.Println("Missing or invalid 'portal_cards'")
		return errors.New("missing or invalid portal_cards")
	}

	var portalCards []string
	for _, card := range rawCards {
		if s, ok := card.(string); ok {
			portalCards = append(portalCards, s)
		} else {
			log.Println("Non-string portal card found")
			return errors.New("portal_cards must be strings")
		}
	}

	return h.departmentService.CreateDepartment(department, portalCards, *action)
}

func (h *DepartmentHandler) approveUpdateAction(action *entities.CPSAction) error {
	data, ok := action.CurrentAction.(map[string]interface{})
	if !ok {
		return errors.New("invalid action data")
	}

	deptCode, ok := data["department_code"].(string)
	if !ok {
		return errors.New("missing department_code")
	}

	deptName, ok := data["department"].(string)
	if !ok {
		return errors.New("invalid department name")
	}

	portalCardsRaw, ok := data["portal_cards"].([]interface{})
	if !ok {
		return errors.New("invalid portal_cards data")
	}

	var portalCards []string
	for _, v := range portalCardsRaw {
		if card, ok := v.(string); ok {
			portalCards = append(portalCards, card)
		}
	}

	dept := department.CreateDepartmentRequest{
		Department:  deptName,
		PortalCards: portalCards,
	}

	return h.departmentService.UpdateDepartment(deptCode, dept)
}

/*
func (h *DepartmentHandler) ApproveDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	var request department.ApproveDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[ApproveDepartmentRequest] failed to decode request: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message": "Invalid JSON payload"},
		}
		resp.SendJSON()
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[ApproveDepartmentRequest] validation failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message": "Invalid input provided"},
		}
		resp.SendJSON()
		return
	}

	userCtx := r.Context().Value(middleware.ContextKey("user_payload"))
	user, ok := userCtx.(middleware.UserPayload)
	if !ok {
		h.logger.Errorf("[ApproveDepartmentRequest] failed to extract user from context")
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message": "Invalid user context"},
		}
		resp.SendJSON()
		return
	}

	if err := h.departmentService.ApproveDepartmentRequest(request, user.UserID); err != nil {
		h.logger.Errorf("[ApproveDepartmentRequest] service error: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{
			"message": err.Error(),
		},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[ApproveDepartmentRequest] decision: %s for userCode: %s", request.Decision, request.UserCode)
	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data: map[string]string{
			"message": "Department request processed successfully",
		},
	}
	resp.SendJSON()
}
*/
