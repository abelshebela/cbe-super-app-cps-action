package helper

import (
	"database/sql"
	"fmt"
	"strings"

	account_block_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func AttachParentChain(b *imodel.AccountBlock) {
	if b == nil || b.Type != imodel.TypeBranch {
		return
	}

	federalRegionName := strings.TrimSpace(b.FederalRegionName)
	districtName := strings.TrimSpace(b.DistrictName)

	var federalRegion *imodel.AccountBlock
	if federalRegionName != "" {
		federalRegion = &imodel.AccountBlock{
			ID:                federalRegionName,
			Name:              federalRegionName,
			FederalRegionName: federalRegionName,
			Type:              imodel.TypeRegion,
		}
		frn := federalRegionName
		b.RegionID = &frn
	}

	if districtName != "" {
		district := &imodel.AccountBlock{
			ID:           districtName,
			Name:         districtName,
			DistrictName: districtName,
			Type:         imodel.TypeDistrict,
		}
		if b.FederalRegionName != "" {
			district.FederalRegionName = b.FederalRegionName
			frn := b.FederalRegionName
			district.RegionID = &frn
		}
		if federalRegion != nil {
			district.Parent = federalRegion
		}
		dn := districtName
		b.DistrictID = &dn
		b.Parent = district
		return
	}

	if federalRegion != nil {
		b.Parent = federalRegion
	}
}

func ScanCompanyBranch(scanner interface{ Scan(dest ...any) error }) (*imodel.AccountBlock, error) {
	var ab imodel.AccountBlock
	var districtName, regionName, federalRegionName sql.NullString
	var isEnabledInt int
	var createdAt, updatedAt sql.NullTime

	err := scanner.Scan(
		&ab.ID, &ab.Code, &ab.Name,
		&districtName, &regionName, &federalRegionName,
		&ab.DaoCode, &ab.AccountType, &isEnabledInt,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	ab.Type = imodel.TypeBranch
	ab.IsEnabled = isEnabledInt == 1
	if districtName.Valid {
		ab.DistrictName = districtName.String
	}
	if regionName.Valid {
		ab.RegionName = regionName.String
	}
	if federalRegionName.Valid {
		ab.FederalRegionName = federalRegionName.String
	}
	if createdAt.Valid {
		ab.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		ab.UpdatedAt = updatedAt.Time
	}

	AttachParentChain(&ab)
	return &ab, nil
}

func ScanRegionAggregate(scanner interface{ Scan(dest ...any) error }) (*imodel.AccountBlock, error) {
	var ab imodel.AccountBlock
	var isEnabledInt int
	var createdAt, updatedAt sql.NullTime

	err := scanner.Scan(
		&ab.FederalRegionName,
		&isEnabledInt, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	ab.ID = ab.FederalRegionName
	ab.Name = ab.FederalRegionName
	ab.Type = imodel.TypeRegion
	ab.IsEnabled = isEnabledInt == 1
	frn := ab.FederalRegionName
	ab.RegionID = &frn
	if createdAt.Valid {
		ab.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		ab.UpdatedAt = updatedAt.Time
	}
	return &ab, nil
}

func ScanDistrictAggregate(scanner interface{ Scan(dest ...any) error }) (*imodel.AccountBlock, error) {
	var ab imodel.AccountBlock
	var federalRegionName sql.NullString
	var isEnabledInt int
	var createdAt, updatedAt sql.NullTime

	err := scanner.Scan(
		&federalRegionName, &ab.DistrictName,
		&isEnabledInt, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	ab.ID = ab.DistrictName
	ab.Name = ab.DistrictName
	ab.Type = imodel.TypeDistrict
	ab.IsEnabled = isEnabledInt == 1
	if federalRegionName.Valid {
		ab.FederalRegionName = federalRegionName.String
		frn := federalRegionName.String
		ab.RegionID = &frn
	}
	dn := ab.DistrictName
	ab.DistrictID = &dn
	if createdAt.Valid {
		ab.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		ab.UpdatedAt = updatedAt.Time
	}

	if ab.FederalRegionName != "" {
		ab.Parent = &imodel.AccountBlock{
			ID:                ab.FederalRegionName,
			Name:              ab.FederalRegionName,
			FederalRegionName: ab.FederalRegionName,
			Type:              imodel.TypeRegion,
		}
	}
	return &ab, nil
}

func BuildInClause(column string, values []string, paramPrefix string) (string, []interface{}) {
	if len(values) == 0 {
		return "", nil
	}
	conditions := make([]string, len(values))
	args := make([]interface{}, 0, len(values))
	for i, val := range values {
		name := fmt.Sprintf("%s_%d", paramPrefix, i)
		conditions[i] = fmt.Sprintf("%s = :%s", column, name)
		args = append(args, sql.Named(name, val))
	}
	return "(" + strings.Join(conditions, " OR ") + ")", args
}

func NamesFromFilterValue(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case string:
		var codes []string
		for _, part := range strings.Split(t, ",") {
			if s := strings.TrimSpace(part); s != "" {
				codes = append(codes, s)
			}
		}
		return codes
	case []string:
		codes := make([]string, 0, len(t))
		for _, s := range t {
			if s = strings.TrimSpace(s); s != "" {
				codes = append(codes, s)
			}
		}
		return codes
	case []interface{}:
		var codes []string
		for _, x := range t {
			codes = append(codes, NamesFromFilterValue(x)...)
		}
		return codes
	default:
		if s := strings.TrimSpace(fmt.Sprint(t)); s != "" {
			return []string{s}
		}
		return nil
	}
}

func ConvertCheckers(checkers []model.Checker) []account_block_dto.Checker {
	result := make([]account_block_dto.Checker, 0, len(checkers))
	for _, c := range checkers {
		result = append(result, account_block_dto.Checker{
			CheckerID:          c.CheckerID,
			RoleID:             c.RoleID,
			CheckerIndex:       c.CheckerIndex,
			CheckerName:        c.CheckerName,
			CheckerPhoneNumber: c.CheckerPhoneNumber,
			ApprovedAt:         c.ApprovedAt,
		})
	}
	return result
}

func ConvertAuditors(auditors []model.Auditor) []account_block_dto.Auditor {
	result := make([]account_block_dto.Auditor, 0, len(auditors))
	for _, a := range auditors {
		result = append(result, account_block_dto.Auditor{
			AuditorID:          a.AuditorID,
			RoleID:             a.RoleID,
			AuditorIndex:       a.AuditorIndex,
			AuditorName:        a.AuditorName,
			AuditorPhoneNumber: a.AuditorPhoneNumber,
			AuditorReason:      a.AuditorReason,
			AuditorMark:        account_block_dto.AuditorMark(a.AuditorMark),
			ApprovedAt:         a.ApprovedAt,
		})
	}
	return result
}
