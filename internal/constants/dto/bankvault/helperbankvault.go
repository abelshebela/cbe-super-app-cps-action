package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	"strings"
)

func ToDomainMethod(s string) constants.AccrualMethod {
	switch strings.ToUpper(s) {
	case "SIMPLE":
		return constants.AccrualMethodSimple
	case "COMPOUND":
		return constants.AccrualMethodCompound
	default:
		return constants.AccrualMethod(s)
	}
}

// func ToDomainFrequency(s string) constants.AccrualFrequency {
// 	switch strings.ToUpper(s) {
// 	case "DAILY":
// 		return constants.AccrualFreqDaily
// 	case "MONTHLY":
// 		return constants.AccrualFreqMonthly
// 	case "QUARTERLY":
// 		return constants.AccrualFreqQuarterly
// 	case "ANNUALLY":
// 		return constants.AccrualFreqAnnually
// 	default:
// 		return constants.AccrualFrequency(s)
// 	}
// }
