package accountblock

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// AccountBlockResponse represents any region, district, city, or branch
type AccountBlockResponse struct {
	ID         string              `json:"id,omitempty"`
	Name       string              `json:"name"`
	Code       string              `json:"code"`
	Address    string              `json:"address"`
	Slug       string              `json:"slug"`
	ParentID   string              `json:"parent_id,omitempty"`
	Parent     *model.AccountBlock `json:"parent,omitempty"`
	Type       string              `json:"type"` // R=Region, D=District, C=City, B=Branch
	CityID     string              `bson:"city_id,omitempty" json:"city_id,omitempty"`
	RegionID   string              `bson:"region_id,omitempty" json:"region_id,omitempty"`
	DistrictID string              `bson:"district_id,omitempty" json:"district_id,omitempty"`
	IsEnabled  bool                `json:"is_enabled"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// // BranchResponse represents a branch entity
// type BranchResponse struct {
// 	ID            string    `json:"id,omitempty"`
// 	BranchCode    string    `json:"branch_code"`
// 	BranchName    string    `json:"branch_name"`
// 	BranchAddress string    `json:"branch_address"`
// 	DistrictCode  string    `json:"district_code"`
// 	DistrictName  string    `json:"district_name"`
// 	RegionName    string    `json:"region_name"`
// 	RecordStat    string    `json:"record_stat"`
// 	CreatedAt     time.Time `json:"created_at"`
// 	UpdatedAt     time.Time `json:"updated_at"`
// 	Enabled       bool      `json:"enabled"`
// }

// // RegionResponse represents a region entity
// type RegionResponse struct {
// 	ID            string    `json:"id,omitempty"`
// 	RegionCode    string    `json:"region_code"`
// 	RegionName    string    `json:"region_name"`
// 	RegionAddress string    `json:"region_address"`
// 	CreatedAt     time.Time `json:"created_at"`
// 	UpdatedAt     time.Time `json:"updated_at"`
// 	Enabled       bool      `json:"enabled"`
// }

// // DistrictResponse represents a district entity
// type DistrictResponse struct {
// 	ID              string    `json:"id,omitempty"`
// 	DistrictCode    string    `json:"district_code"`
// 	DistrictName    string    `json:"district_name"`
// 	DistrictAddress string    `json:"district_address"`
// 	RegionID        string    `json:"region_id"`
// 	RegionName      string    `json:"region_name"`
// 	CreatedAt       time.Time `json:"created_at"`
// 	UpdatedAt       time.Time `json:"updated_at"`
// 	Enabled         bool      `json:"enabled"`
// }

// // CityResponse represents a city entity
// type CityResponse struct {
// 	ID           string    `json:"id,omitempty"`
// 	CityCode     string    `json:"city_code"`
// 	CityName     string    `json:"city_name"`
// 	CityAddress  string    `json:"city_address"`
// 	DistrictID   string    `json:"district_id"`
// 	DistrictName string    `json:"district_name"`
// 	RegionID     string    `json:"region_id"`
// 	RegionName   string    `json:"region_name"`
// 	CreatedAt    time.Time `json:"created_at"`
// 	UpdatedAt    time.Time `json:"updated_at"`
// 	Enabled      bool      `json:"enabled"`
// }
