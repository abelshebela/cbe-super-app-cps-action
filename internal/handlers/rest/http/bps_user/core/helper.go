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
	if strings.HasPrefix(filter.Filters["phone"].(string), "+") {
		filter.Filters["phone"] = utils.FormatPhoneNumber(filter.Filters["phone"].(string))
	}
	if strings.HasPrefix(filter.Filters["phone"].(string), "09") {
		filter.Filters["phone"] = "+251" + filter.Filters["phone"].(string)[1:]
	}
	if strings.HasPrefix(filter.Filters["phone"].(string), "07") {
		filter.Filters["phone"] = "+251" + filter.Filters["phone"].(string)[1:]
	}

	return filter
}
