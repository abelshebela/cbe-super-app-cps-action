package branch

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BranchMapper maps a Branch model to a bson.M for updates
func BranchMapper(branch model.Branch) bson.M {
	return bson.M{
		"$set": bson.M{
			"branch_code":    branch.BranchCode,
			"branch_name":    branch.BranchName,
			"branch_address": branch.BranchAddress,
			"district_code":  branch.DistrictCode,
			"district_name":  branch.DistrictName,
			"region_name":  branch.RegionName,
			"record_stat":    branch.RecordStat,
			"enabled":        branch.Enabled,
			"updated_at":     time.Now(),
		},
	}
}
