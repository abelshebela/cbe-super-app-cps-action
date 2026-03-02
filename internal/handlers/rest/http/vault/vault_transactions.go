package vault

import (
	localization "cbe-super-app-cps-action/internal/constants/localization"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
)

func (h *handler) GetVaultTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getVaultTransactions", "handler", "getVaultTransactions")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	params := common_utils.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := common_utils.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := common_utils.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, err := h.service.FindAllVaultTransactions(ctx, params)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[getVaultTransactions] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessVaultTransactionsRetrievedS, result)
}

func (h *handler) GetVaultTransaction(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "GetVaultTransaction", "handler", "GetVaultTransaction")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "transaction_id")

	span.SetAttributes(attribute.String("vault_transaction.id", id))
	result, err := h.service.FindVaultTransaction(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetVaultTransaction] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[VaultTxnH][GetByID] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultTransactionRetrievedS, result)

}
