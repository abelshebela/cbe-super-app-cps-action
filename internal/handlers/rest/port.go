package rest

import "net/http"

type Users interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	ChangePin(w http.ResponseWriter, r *http.Request)
	VerifyOtp(w http.ResponseWriter, r *http.Request)
	SetPin(w http.ResponseWriter, r *http.Request)
	UpdateProfilePicture(w http.ResponseWriter, r *http.Request)
	UpdateProfileTheme(w http.ResponseWriter, r *http.Request)
	ForgetPinSendOtp(w http.ResponseWriter, r *http.Request)
	ResetPin(w http.ResponseWriter, r *http.Request)
	PreLogin(w http.ResponseWriter, r *http.Request)
	DeviceLookup(w http.ResponseWriter, r *http.Request)
	Healthcheck(w http.ResponseWriter, r *http.Request)
	VerifyForgetPinOtp(w http.ResponseWriter, r *http.Request)
}

type CPSActionAdapter interface {
	ApproveCPSAction(w http.ResponseWriter, r *http.Request)
	RejectCPSAction(w http.ResponseWriter, r *http.Request)
	GetCPSActionsByDepartment(w http.ResponseWriter, r *http.Request)
	GetCPSActionByID(w http.ResponseWriter, r *http.Request)
	GetCPSActionByActionCode(w http.ResponseWriter, r *http.Request)
}
