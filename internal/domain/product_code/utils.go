package productcode

func nonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

// nonEmptyProductCodes updates ProductCodes fields if provided, otherwise preserves existing values
func nonEmptyProductCodes(request, existing ProductCodes) ProductCodes {
	return ProductCodes{
		PRD:    nonEmptyString(request.PRD, existing.PRD),
		VATPRD: nonEmptyString(request.VATPRD, existing.VATPRD),
		SFPRD:  nonEmptyString(request.SFPRD, existing.SFPRD),
		TRXN:   nonEmptyString(request.TRXN, existing.TRXN),
	}
}
