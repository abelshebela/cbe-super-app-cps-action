package productcode

import (
	"context"
	"testing"
	"time"

	"cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/errors"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	pcode "cbe-super-app-cps-action/internal/service/productcode"
	"cbe-super-app-cps-action/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// MockProductCodeRepository is a mock implementation of the ProductCodeRepository interface
type MockProductCodeRepository struct {
	storage.ProductCodeRepository
	FetchByIDFunc func(ctx context.Context, id string) (*model.ProductCode, error)
	UpdateFunc    func(ctx context.Context, p *model.ProductCode) error
}

func (m *MockProductCodeRepository) FetchByID(ctx context.Context, id string) (*model.ProductCode, error) {
	return m.FetchByIDFunc(ctx, id)
}

func (m *MockProductCodeRepository) Update(ctx context.Context, p *model.ProductCode) error {
	return m.UpdateFunc(ctx, p)
}

// MockCPSActionService is a mock implementation of the CPSActionService interface
type MockCPSActionService struct {
	service.CPSActionService
	CreateCPSActionFunc func(ctx context.Context, action *model.CPSAction) error
}

func (m *MockCPSActionService) CreateCPSAction(ctx context.Context, action *model.CPSAction) error {
	return m.CreateCPSActionFunc(ctx, action)
}

func TestUpdateProductCode(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	logger := utils.NewLogger()

	existingProductCode := &model.ProductCode{
		ID:          "test-id",
		ProductName: "Test Product",
		CBEProductCodes: model.ProductCodes{
			PRD:    "123",
			VATPRD: "456",
			SFPRD:  "789",
			TRXN:   "101",
		},
		CBEIFBProductCodes: model.ProductCodes{
			PRD:    "ifb-123",
			VATPRD: "ifb-456",
			SFPRD:  "ifb-789",
			TRXN:   "ifb-101",
		},
		CreatedAt:     now,
		LastUpdatedAt: now,
	}

	updateRequest := productcode.UpdateProductCodeRequest{
		ID:          "test-id",
		ProductName: "Updated Test Product",
	}

	t.Run("successful update", func(t *testing.T) {
		mockRepo := &MockProductCodeRepository{
			FetchByIDFunc: func(ctx context.Context, id string) (*model.ProductCode, error) {
				return existingProductCode, nil
			},
		}
		mockCPSActionService := &MockCPSActionService{
			CreateCPSActionFunc: func(ctx context.Context, action *model.CPSAction) error {
				return nil
			},
		}

		productCodeService := pcode.NewProductCodeService(mockRepo, mockCPSActionService, logger)

		existing, updated, err := productCodeService.UpdateProductCode(ctx, updateRequest)

		assert.NoError(t, err)
		assert.NotNil(t, existing)
		assert.NotNil(t, updated)
		assert.Equal(t, "Updated Test Product", updated.ProductName)
	})

	t.Run("product code not found", func(t *testing.T) {
		mockRepo := &MockProductCodeRepository{
			FetchByIDFunc: func(ctx context.Context, id string) (*model.ProductCode, error) {
				return nil, mongo.ErrNoDocuments
			},
		}

		productCodeService := pcode.NewProductCodeService(mockRepo, nil, logger)

		existing, updated, err := productCodeService.UpdateProductCode(ctx, updateRequest)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrProductCodeNotFound, err)
		assert.Nil(t, existing)
		assert.Nil(t, updated)
	})

	t.Run("no update", func(t *testing.T) {
		noUpdateRequest := productcode.UpdateProductCodeRequest{
			ID:          "test-id",
			ProductName: "Test Product",
		}
		mockRepo := &MockProductCodeRepository{
			FetchByIDFunc: func(ctx context.Context, id string) (*model.ProductCode, error) {
				return existingProductCode, nil
			},
		}

		productCodeService := pcode.NewProductCodeService(mockRepo, nil, logger)

		existing, updated, err := productCodeService.UpdateProductCode(ctx, noUpdateRequest)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrNoUpdate, err)
		assert.Nil(t, existing)
		assert.Nil(t, updated)
	})
}
