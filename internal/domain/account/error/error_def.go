package error
type ErrorDefinition struct {
    Code    string
    Message string
}

func (e ErrorDefinition) Error() string {
    return e.Message
}

type ErrorGroup map[string]ErrorDefinition

type ErrorDefinitions struct {
    Account ErrorGroup
}

var DefineError = ErrorDefinitions{
    Account: ErrorGroup{
        "PHONE_LOOKUP_FAILED": {
            Code:    "ACC_001",
            Message: "Failed to perform phone number lookup.",
        },
        "MOCK_DATA_FETCH_FAILED": {
            Code:    "ACC_002",
            Message: "Failed to fetch mock user data.",
        },
        "API_REQUEST_FAILED": {
            Code:    "ACC_003",
            Message: "Failed to communicate with external API.",
        },
    },
}
