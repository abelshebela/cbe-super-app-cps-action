package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/budget_category"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func ParseRequestFromMultipartForm(r *http.Request, isCreate bool) (budget_category.CreateBudgetRequest, error) {
	var req budget_category.CreateBudgetRequest

	if err := r.ParseMultipartForm(315 << 20); err != nil {
		return req, errors.New("failed to parse multipart form")
	}

	req.Name = r.FormValue("name")
	req.Color = r.FormValue("color")

	_, iconHeader, err := utils.ParseMultipartFormFile(r, "icon", int64(constants.MaxMemoryForUpload))
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return req, errors.New("icon file is required")
		}
		return req, err
	}
	req.Icon = iconHeader

	return req, nil
}

func ParseUpdateRequestFromMultipartForm(r *http.Request) (budget_category.UpdateBudgetRequest, error) {
	var req budget_category.UpdateBudgetRequest

	if err := r.ParseMultipartForm(315 << 20); err != nil {
		return req, errors.New("failed to parse multipart form")
	}

	if name := r.FormValue("name"); name != "" {
		req.Name = name
	}
	if color := r.FormValue("color"); color != "" {
		req.Color = color
	}

	_, iconHeader, err := utils.ParseMultipartFormFile(r, "icon", int64(constants.MaxMemoryForUpload))
	if err != nil {
		fmt.Println("///error for icon", err)
		if !errors.Is(err, http.ErrMissingFile) {
			return req, err
		}
	} else {
		req.Icon = iconHeader
	}

	return req, nil
}

func ExtractIDFromURL(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
