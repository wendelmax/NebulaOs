package api

import (
	"encoding/json"
	"net/http"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type BillingHandler struct {
	billingMgr domain.BillingManager
}

func NewBillingHandler(mgr domain.BillingManager) *BillingHandler {
	return &BillingHandler{billingMgr: mgr}
}

func (h *BillingHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = "default-tenant"
	}

	report, err := h.billingMgr.GenerateReport(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
