package newscategory_dto

type CreateNewsCategoryRequest struct {
	CategoryName []string `json:"category_name"`
}
type UpdateNewsCategoryRequest struct {
	CategoryName string `json:"category_name"`
}
