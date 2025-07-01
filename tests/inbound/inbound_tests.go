package inbound_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"cbe-super-app-cps-action/internal/application/dto"
	mock_inbound "cbe-super-app-cps-action/mocks/port/inbound/bulk_services"
)

func TestFetchServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInbound := mock_inbound.NewMockInbound(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cbesuperapp/cps_config/fetch_service?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	mockInbound.EXPECT().FetchServices(w, req).Times(1)

	mockInbound.FetchServices(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestEnableDisableServicesMaker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInbound := mock_inbound.NewMockInbound(ctrl)

	payload := dto.EnableDisableServiceMakerDtoRequest{
		ServiceId:     "123",
		ServiceAction: true,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cbesuperapp/cps_config/enable_disable_service_maker", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockInbound.EXPECT().EnableDisableServicesMaker(w, req).Times(1)

	mockInbound.EnableDisableServicesMaker(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestEnableDisableServicesChecker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInbound := mock_inbound.NewMockInbound(ctrl)

	payload := dto.EnableDisableServiceCheckerDtoRequest{
		Action_Id:     "123",
		ServiceAction: false,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cbesuperapp/cps_config/enable_disable_service_checker", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockInbound.EXPECT().EnableDisableServicesChecker(w, req).Times(1)

	mockInbound.EnableDisableServicesChecker(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
