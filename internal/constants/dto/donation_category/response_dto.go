package donation_category

type DonationCategoryListResponse struct {
	ID             string `json:"id" bson:"id"`
	CategoryName   string `json:"category_name" bson:"category_name"`
	Icon           string `json:"donation_icon" bson:"donation_icon"`
	IsDeleted      bool   `json:"is_deleted" bson:"is_deleted"`
	Enabled        bool   `json:"enabled" bson:"enabled"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
}
type DonationCategoryResponse struct {
	CategoryName string `json:"category_name" bson:"category_name"`
	Icon         string `json:"donation_icon" bson:"donation_icon"`
	Enabled      bool   `json:"enabled" bson:"enabled"`
}
