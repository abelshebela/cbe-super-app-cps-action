package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	mock_repository "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/mocks/domain/action"
)

func TestGetServicePaginated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)
	serviceStore := action.NewService(mockRepo)

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedServices := []action.Service{
		{
			ID:          stringToPointer("1"),
			Key:         "service1",
			ServiceName: "Service 1",
			Flag:        true,
		},
		{
			ID:          stringToPointer("2"),
			Key:         "service2",
			ServiceName: "Service 2",
			Flag:        false,
		},
	}

	mockRepo.EXPECT().GetAllHqServicesPaginated(ctx, offset, limit).Return(expectedServices, nil)

	services, err := serviceStore.GetServicePaginated(ctx, limit, offset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(services) != len(expectedServices) {
		t.Errorf("expected %d services, got %d", len(expectedServices), len(services))
	}
}

func TestUpdateServiceFlag(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mock_repository.NewMockRepository(ctrl)
    serviceStore := action.NewService(mockRepo)

    ctx := context.Background()
    serviceID := "1"
    action_taken := true

    service := action.Service{
        ID:          stringToPointer(serviceID),
        Key:         "service1",
        ServiceName: "Service 1",
        Flag:        false,
    }

    // Add this expectation if your code calls FetchCpsActionById
    mockRepo.EXPECT().FetchCpsActionById(ctx, serviceID).Return(action.CPSAction{}, nil)

    mockRepo.EXPECT().GetHqServiceById(ctx, serviceID).Return(service, nil)
    mockRepo.EXPECT().UpdateHqService(ctx, gomock.Any()).Return(nil)

    err := serviceStore.UpdateServiceFlag(ctx, serviceID, action_taken, "")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestUpdateServiceFlag_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)
	serviceStore := action.NewService(mockRepo)

	ctx := context.Background()
	serviceID := "1"
	action_taken := true

	mockRepo.EXPECT().GetHqServiceById(ctx, serviceID).Return(action.Service{}, errors.New("service not found"))

	err := serviceStore.UpdateServiceFlag(ctx, serviceID, action_taken, "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func stringToPointer(s string) *string {
	return &s
}
