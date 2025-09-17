package utils

type DecisonEnum string

const (
	DecisionApproved DecisonEnum = "APPROVED"
	DecisionDenied   DecisonEnum = "DENIED"
)

func Decison(res bool) DecisonEnum {
	switch res {
	case true:
		return DecisionApproved
	case false:
		return DecisionDenied
	}

	return ""
}
