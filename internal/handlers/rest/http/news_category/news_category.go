package newscategory_handler

import (
	newscategory_dto "cbe-super-app-cps-action/internal/constants/dto/news_category"
	newscategory_adaptor "cbe-super-app-cps-action/internal/constants/interfaces/news_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/news_category/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	_ "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type NewsCategoryHandler struct {
	service service.NewsCategoryService
	logger  utils.Logger
}

func NewNewsCategoryHandler(newsCategoryService service.NewsCategoryService, logger utils.Logger) newscategory_adaptor.NewsCategoryAdaptor {
	return NewsCategoryHandler{
		service: newsCategoryService,
		logger:  logger,
	}
}

// CreateNewsCategory creates a new news category
//
//	@Summary		Create news category
//	@Description	Creates a new news category
//	@Tags			NewsCategory
//	@Accept			json
//	@Produce		json
//	@Param			body	body		newscategory_dto.CreateNewsCategoryRequest	true	"Create News Category Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}		"News category created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/category/create [post]
//
// CreateNewsCategory implements newscategory_adaptor.NewsCategoryAdaptor.
func (n NewsCategoryHandler) CreateNewsCategory(w http.ResponseWriter, r *http.Request) {
	var req newscategory_dto.CreateNewsCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		n.logger.Errorf("failed to decode create news category request: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}

	if core.ValidateString(req.CategoryName) != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	if err := n.service.CreateNewsCategory(r.Context(), req.CategoryName); err != nil {
		n.logger.Errorf("[CreateNewsCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	n.logger.Infof("[CreateNewsCategory] request sent successfully for category_name: %s", req.CategoryName)
	localization.SendSuccessResponse(w, localization.SuccessNewsCategoryCreated, nil)
}

// DeleteNewsCategory deletes a news category by ID
//
//	@Summary		Delete news category
//	@Description	Deletes a news category by its ID
//	@Tags			NewsCategory
//	@Produce		json
//	@Param			id	path		string									true	"News Category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"News category deleted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request - ID required"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/category/{id} [delete]
func (n NewsCategoryHandler) DeleteNewsCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		n.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	if err := n.service.DeleteNewsCategory(r.Context(), id); err != nil {
		n.logger.Errorf("[DeleteNewsCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	n.logger.Infof("[DeleteNewsCategory] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessNewsCategoryDeleted, nil)
}

// FetchNewsCategories returns a paginated list of news categories
//
//	@Summary		List news categories
//	@Description	Retrieves a paginated list of news categories
//	@Tags			NewsCategory
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int															false	"Page number"
//	@Param			per_page	query		int															false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=[]model.NewsCategory}	"List of news categories"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}						"Bad request - Invalid pagination params"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}						"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/category [get]
func (n NewsCategoryHandler) FetchNewsCategories(w http.ResponseWriter, r *http.Request) {
	filterPtr := local_util.ExtractFilterParams(r)
	filter := *filterPtr

	if filter.Page < 0 || filter.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidPaginationParams, nil, nil)
		return
	}

	list, err := n.service.FindAllWithPagination(r.Context(), filter)
	if err != nil {
		n.logger.Errorf("[FetchNewsCategories] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	n.logger.Infof("[FetchNewsCategories] retrieved %d news categories", len(list.Data))
	// return as-is
	localization.SendSuccessResponse(w, localization.SuccessNewsCategoryFetched, list)
}

// GetNewsCategoryByID returns a news category by its ID
//
//	@Summary		Get news category by ID
//	@Description	Retrieves a single news category by its ID
//	@Tags			NewsCategory
//	@Produce		json
//	@Param			id	path		string													true	"News Category ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.NewsCategory}	"News category details"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}					"Bad request - ID required"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/category/{id} [get]
//
// GetNewsCategoryByID implements newscategory_adaptor.NewsCategoryAdaptor.
func (n NewsCategoryHandler) GetNewsCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		n.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	data, err := n.service.GetNewsCategoryByID(r.Context(), id)
	if err != nil {
		n.logger.Errorf("[GetNewsCategoryByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	n.logger.Infof("[GetNewsCategoryByID] news category retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessNewsCategoryFetched, data)
}

// UpdateNewsCategory updates a news category by ID
//
//	@Summary		Update news category
//	@Description	Updates a news category by its ID
//	@Tags			NewsCategory
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string										true	"News Category ID"
//	@Param			body	body		newscategory_dto.UpdateNewsCategoryRequest	true	"Update News Category Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}		"News category updated successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/category/{id} [patch]
//
// UpdateNewsCategory implements newscategory_adaptor.NewsCategoryAdaptor.
func (n NewsCategoryHandler) UpdateNewsCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		n.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var req newscategory_dto.UpdateNewsCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		n.logger.Errorf("failed to decode update news category request: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}

	if req.CategoryName == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	// Note: service.UpdateNewsCategory signature accepts only the category name.
	if err := n.service.UpdateNewsCategory(r.Context(), id, req.CategoryName); err != nil {
		n.logger.Errorf("[UpdateNewsCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	n.logger.Infof("[UpdateNewsCategory] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessNewsCategoryUpdated, nil)
}
