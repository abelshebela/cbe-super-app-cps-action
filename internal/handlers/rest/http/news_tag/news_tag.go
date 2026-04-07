package newstag_handler

import (
	newstag_dto "cbe-super-app-cps-action/internal/constants/dto/news_tag"
	newstag_adaptor "cbe-super-app-cps-action/internal/constants/interfaces/news_tag"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/news_tag/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	_ "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type NewsTagHandler struct {
	service service.NewsTagService
	logger  utils.Logger
}

func NewNewsTagHandler(newsTagService service.NewsTagService, logger utils.Logger) newstag_adaptor.NewsTagAdaptor {
	return NewsTagHandler{
		service: newsTagService,
		logger:  logger,
	}
}

// CreateNewsTags creates a new news tag
//
//	@Summary		Create news tag
//	@Description	Creates a new news tag
//	@Tags			NewsTags
//	@Accept			json
//	@Produce		json
//	@Param			body	body		newstag_dto.CreateNewsTagRequest		true	"Create News Tag Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"News tag created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/tags/create [post]
func (n NewsTagHandler) CreateNewsTags(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), n.logger)
	var req newstag_dto.CreateNewsTagRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[NewsTagH][Create] decode body err: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}

	if core.ValidateString(req.TagName) != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	if err := n.service.CreateNewsTags(r.Context(), req.TagName); err != nil {
		log.Errorf("[CreateNewsTags] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[CreateNewsTags] request sent successfully for tag_name: %s", req.TagName)
	localization.SendSuccessResponse(w, localization.SuccessNewsTagCreated, nil)
}

// DeleteNewsTag deletes a news tag by ID
//
//	@Summary		Delete news tag
//	@Description	Deletes a news tag by its ID
//	@Tags			NewsTags
//	@Produce		json
//	@Param			id	path		string									true	"News Tag ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"News tag deleted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request - ID required"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/tags/{id} [delete]
//
// DeleteNewsTag implements newstag_adaptor.NewsTagAdaptor.
func (n NewsTagHandler) DeleteNewsTag(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), n.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[NewsTagH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	if err := n.service.DeleteNewsTag(r.Context(), id); err != nil {
		log.Errorf("[DeleteNewsTag] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[DeleteNewsTag] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessNewsTagDeleted, nil)
}

// FetchNewsTags returns a paginated list of news tags
//
//	@Summary		List news tags
//	@Description	Retrieves a paginated list of news tags
//	@Tags			NewsTags
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int													false	"Page number"
//	@Param			per_page	query		int													false	"Items per page"
//	@Param			search		query		string												false	"search field tag_name"
//	@Param			tag_name	query		string												false	"filter key list tag_name"
//	@Success		200			{object}	localization.StandardResponse{data=[]model.NewsTag}	"List of news tags"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}				"Bad request - Invalid pagination params"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/tags [get]
func (n NewsTagHandler) FetchNewsTags(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), n.logger)
	filterPtr := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	ftr := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(ftr); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	filter := *filterPtr

	if filter.Page < 0 || filter.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidPaginationParams, nil, nil)
		return
	}

	list, err := n.service.FindAllWithPagination(r.Context(), filter)
	if err != nil {
		log.Errorf("[FetchNewsTags] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[FetchNewsTags] retrieved %d news tags", len(list.Data))
	localization.SendSuccessResponse(w, localization.SuccessNewsTagFetched, list)
}

// GetNewsTagByID returns a news tag by its ID
//
//	@Summary		Get news tag by ID
//	@Description	Retrieves a single news tag by its ID
//	@Tags			NewsTags
//	@Produce		json
//	@Param			id	path		string												true	"News Tag ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.NewsTag}	"News tag details"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}				"Bad request - ID required"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/tags/{id} [get]
//
// GetNewsTagByID implements newstag_adaptor.NewsTagAdaptor.
func (n NewsTagHandler) GetNewsTagByID(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), n.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[NewsTagH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	data, err := n.service.GetNewsTagByID(r.Context(), id)
	if err != nil {
		log.Errorf("[GetNewsTagByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[GetNewsTagByID] news tag retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessNewsTagFetched, data)
}

// UpdateNewsTag updates a news tag by ID
//
//	@Summary		Update news tag
//	@Description	Updates a news tag by its ID
//	@Tags			NewsTags
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"News Tag ID"
//	@Param			body	body		newstag_dto.UpdateNewsTagRequest		true	"Update News Tag Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"News tag updated successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/news/tags/{id} [patch]
//
// UpdateNewsTag implements newstag_adaptor.NewsTagAdaptor.
func (n NewsTagHandler) UpdateNewsTag(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), n.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[NewsTagH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var req newstag_dto.UpdateNewsTagRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[NewsTagH][Update] decode body err: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}

	if req.TagName == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	if err := n.service.UpdateNewsTag(r.Context(), id, req.TagName); err != nil {
		log.Errorf("[UpdateNewsTag] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[UpdateNewsTag] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessNewsTagUpdated, nil)
}
