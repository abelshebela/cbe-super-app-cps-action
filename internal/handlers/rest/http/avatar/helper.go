package avatar

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/avatar"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"net/http"
)

func ReqFileParse(r *http.Request) (avatar.AvatarDTO, error) {
	var req avatar.AvatarDTO

	file, fileHeader, err := local_util.ParseMultipartFormFile(r, "avatar", int64(constants.MaxMemoryForUpload))
	if err != nil && (err.Error() != localization.ErrorFileNotFound.Code) {
		return avatar.AvatarDTO{}, errors.New(localization.ErrorFileNotFound.Code)
	}
	if file != nil {
		defer file.Close()
		req.Avatar = fileHeader
	}

	req.Label = r.FormValue("label")
	return req, nil
}
