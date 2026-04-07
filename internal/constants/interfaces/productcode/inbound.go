package productcode

import "net/http"

type ProductCodeAdapter interface {
	UpdateProductCode(w http.ResponseWriter, r *http.Request)
	FetchProductCodeByID(w http.ResponseWriter, r *http.Request)
	FetchProductCodes(w http.ResponseWriter, r *http.Request)
}
