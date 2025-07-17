package inbound

import (
	"net/http"
)

type PasswordRuleInbound interface {
	GetPasswordRule(w http.ResponseWriter, r *http.Request)
	RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request)
	ApproveOrRejectPasswordRuleAction(w http.ResponseWriter, r *http.Request)
	GetPasswordRuleUpdateActionByID(w http.ResponseWriter, r *http.Request)
	GetUpdateAction(w http.ResponseWriter, r *http.Request)
	CheckPasswordRule(w http.ResponseWriter, r *http.Request)
}
