package newstag_dto

type CreateNewsTagRequest struct {
	TagName []string `json:"tag_name"`
}
type UpdateNewsTagRequest struct {
	TagName string `json:"tag_name"`
}
