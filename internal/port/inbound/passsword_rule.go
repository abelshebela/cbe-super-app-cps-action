package inbound

import (
    "net/http"
)

type PasswordRuleInbound interface {
    RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request)
    ApproveOrRejectPasswordRuleAction(w http.ResponseWriter, r *http.Request)
    GetPasswordRuleUpdateActionByID(w http.ResponseWriter, r *http.Request)
}