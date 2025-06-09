package users

import "net/http"

type UserPortHandler interface {
    FetchLinkedAccounts(w http.ResponseWriter, r *http.Request)
    GenerateEmailOTP(w http.ResponseWriter, r *http.Request)
    VerifyEmailOTP(w http.ResponseWriter, r *http.Request)
}