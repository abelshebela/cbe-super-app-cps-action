package inbound

import "net/http"

type CPSUserMakerHandler interface {
    CreateUserRequest(w http.ResponseWriter, r *http.Request)
    UpdateUserRequest(w http.ResponseWriter, r *http.Request)
    ApproveUserAction(w http.ResponseWriter, r *http.Request)
    GetPendingUserActions(w http.ResponseWriter, r *http.Request)
}