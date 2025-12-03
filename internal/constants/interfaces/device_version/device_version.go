package deviceversion

import "net/http"

type DeviceVersionHandler interface {
	CreateDeviceVersion(w http.ResponseWriter, r *http.Request)
	UpdateDeviceVersion(w http.ResponseWriter, r *http.Request)
	GetAllDeviceVersions(w http.ResponseWriter, r *http.Request)
	GetDeviceVersionByID(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
