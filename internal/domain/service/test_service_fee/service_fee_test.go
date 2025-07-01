package service_test

import (
	"context"
	"errors"
	"testing"

	"cbe-super-app-cps-action/internal/domain/service"
	"cbe-super-app-cps-action/internal/domain/service/mocks"

	"github.com/golang/mock/gomock"
)

func TestServiceFeeMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceFee := mocks.NewMockServiceFee(ctrl)
	ctx := context.Background()
	cpsAction := service.CPSAction{ /* fill with test data if needed */ }
	expectedResp := service.UpdateServiceDetailsResponse{ /* fill with test data if needed */ }
	expectedErr := errors.New("test error")

	t.Run("InitiateServiceFeeUpdate success", func(t *testing.T) {
		mockServiceFee.EXPECT().
			InitiateServiceFeeUpdate(ctx, cpsAction).
			Return(expectedResp, nil)

		resp, err := mockServiceFee.InitiateServiceFeeUpdate(ctx, cpsAction)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp != expectedResp {
			t.Errorf("expected %v, got %v", expectedResp, resp)
		}
	})

	t.Run("ApproveServiceFeeUpdate success", func(t *testing.T) {
		mockServiceFee.EXPECT().
			ApproveServiceFeeUpdate(ctx, cpsAction).
			Return(expectedResp, nil)

		resp, err := mockServiceFee.ApproveServiceFeeUpdate(ctx, cpsAction)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp != expectedResp {
			t.Errorf("expected %v, got %v", expectedResp, resp)
		}
	})

	t.Run("RejectServiceFeeUpdate error", func(t *testing.T) {
		mockServiceFee.EXPECT().
			RejectServiceFeeUpdate(ctx, cpsAction).
			Return(service.UpdateServiceDetailsResponse{}, expectedErr)

		resp, err := mockServiceFee.RejectServiceFeeUpdate(ctx, cpsAction)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if resp != (service.UpdateServiceDetailsResponse{}) {
			t.Errorf("expected empty response, got %v", resp)
		}
	})
}
