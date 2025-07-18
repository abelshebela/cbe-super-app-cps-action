package inbound

import (
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
)

type Inbound interface {
	FetchServices(w http.ResponseWriter, r *http.Request)
	EnableDisableServicesMaker(w http.ResponseWriter, r *http.Request)
	EnableDisableServicesChecker(w http.ResponseWriter, r *http.Request, req dto.EnableDisableServiceCheckerDtoRequest)
	DisableServicesChecker(w http.ResponseWriter, r *http.Request)
	EnableServicesChecker(w http.ResponseWriter, r *http.Request)
	SearchAccountByAccountNumber(w http.ResponseWriter, r *http.Request)
	RemoveCifMaker(w http.ResponseWriter, r *http.Request)
	RemoveCifChecker(w http.ResponseWriter, r *http.Request, req dto.CifRemoveCheckerRequest)
	RemoveCifCheckerApprove(w http.ResponseWriter, r *http.Request)
	RemoveCifCheckerReject(w http.ResponseWriter, r *http.Request)
}
