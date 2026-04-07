package passwordrule

import "net/http"

type PasswordRule interface {
	GetPasswordRule(w http.ResponseWriter, r *http.Request)
	RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request)
	CheckPasswordRule(w http.ResponseWriter, r *http.Request)
}
