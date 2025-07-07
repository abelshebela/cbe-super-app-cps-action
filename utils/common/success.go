package common

type SuccessDefinition struct {
	Code    string
	Message string
}

type SuccessGroup map[string]SuccessDefinition

type SuccessDefinitions struct {
	General     SuccessGroup
	Auth        SuccessGroup
	User        SuccessGroup
	Transaction SuccessGroup
}

var DefineSuccess = SuccessDefinitions{
	General: SuccessGroup{
		"SUCCESS": {
			Code:    "GEN_SUCCESS_000",
			Message: "Successful request.",
		},
	},
}
