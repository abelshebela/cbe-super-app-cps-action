package core

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/pkgs/utils"
	"strings"
)

func PhoneModifier(filter *types.Filter) *types.Filter {
	if filter.Filters["phone"] == "" {
		return filter
	}

	if strings.HasPrefix(filter.Search, "+") {
		filter.Filters["phone"] = utils.FormatPhoneNumber(filter.Search)
	}

	if strings.HasPrefix(filter.Search, "09") {
		filter.Filters["phone"] = "+251" + filter.Search[1:]
	}

	if strings.HasPrefix(filter.Search, "9") {
		filter.Filters["phone"] = "+251" + filter.Search
	}

	if strings.HasPrefix(filter.Search, "07") {
		filter.Filters["phone"] = "+251" + filter.Search[1:]
	}

	if strings.HasPrefix(filter.Search, "7") {
		filter.Filters["phone"] = "+251" + filter.Search
	}

	return filter
}
