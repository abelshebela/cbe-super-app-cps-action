package service_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/mocks"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/service"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.uber.org/mock/gomock"
)

func createMockFileHeader(filename string, content []byte, size int64) *multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Simulate file field
	part, _ := writer.CreateFormFile("file", filename)
	part.Write(content)
	writer.Close()

	// Create fake request
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Parse and extract FileHeader
	req.ParseMultipartForm(10 << 20)
	fileHeader := req.MultipartForm.File["file"][0]

	// Override the size with the parameter
	fileHeader.Size = size

	return fileHeader
}

// MockLogger implements utils.Logger interface for testing
type MockLogger struct{}

func (l *MockLogger) Debug(args ...interface{})                 {}
func (l *MockLogger) Debugf(format string, args ...interface{}) {}
func (l *MockLogger) Info(args ...interface{})                  {}
func (l *MockLogger) Infof(format string, args ...interface{})  {}
func (l *MockLogger) Warn(args ...interface{})                  {}
func (l *MockLogger) Warnf(format string, args ...interface{})  {}
func (l *MockLogger) Error(args ...interface{})                 {}
func (l *MockLogger) Errorf(format string, args ...interface{}) {}
func (l *MockLogger) Fatal(args ...interface{})                 {}
func (l *MockLogger) Fatalf(format string, args ...interface{}) {}
func (l *MockLogger) Sync() error                               { return nil }

// MockMinioClient implements config.MinioClientInterface for testing
type MockMinioClient struct{}

// SaveObjectN implements config.MinioClientInterface.
func (m *MockMinioClient) SaveObjectN(ctx context.Context, obj config.SaveObjectBodyN) (*minio.UploadInfo, error) {
	return &minio.UploadInfo{
		Bucket: "test-bucket",
		Key:    "test-key",
	}, nil
}

func (m *MockMinioClient) BucketExist(ctx context.Context, bucketName string) (bool, error) {
	return true, nil
}

func (m *MockMinioClient) MakeBucket(ctx context.Context, bucketName string) (bool, error) {
	return true, nil
}

func (m *MockMinioClient) SaveObject(ctx context.Context, body config.SaveObjectBody) (*config.SaveObjectResponse, error) {
	return &config.SaveObjectResponse{
		Bucket: "test-bucket",
		Key:    "test-key",
	}, nil
}

func (m *MockMinioClient) GetObject(ctx context.Context, bucketName string, objectName string) (*minio.Object, error) {
	return nil, nil
}

func (m *MockMinioClient) DeleteObject(ctx context.Context, body config.DeleteObjectBody) (bool, error) {
	return true, nil
}

func (m *MockMinioClient) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return []minio.BucketInfo{
		{
			Name:         "test-bucket",
			CreationDate: time.Now(),
		},
	}, nil
}

func TestWalletDomain_CreateWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	// Create test data with full struct values
	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	testCreateWalletRequest := dto.CreateWalletRequest{
		Name:   "TestWallet",
		Avatar: createMockFileHeader("test-avatar.png", []byte("test-avatar"), 1024),
		Code:   "TW001",
	}

	testCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT001",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestCreateWallet,
		ActionType:      model.ActionCreate,
		ActionData:      testCreateWalletRequest,
		MakerActionTime: testTime,
	}

	testWalletEntity := entity.Wallet{
		ID:             "WALLET001",
		Name:           "TestWallet",
		Avatar:         "test-bucket/test-key",
		Code:           "TW001",
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      testTime,
		LastModifiedAt: testTime,
		DeletedAt:      time.Time{},
	}

	tests := []struct {
		name    string
		req     model.CreateCPSAction
		mock    func()
		want    *model.CPSAction
		wantErr bool
	}{
		{
			name: "success",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil) // Assuming no existing action with the same code
				mockRepo.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Return(&model.CPSAction{
						ID:              "CPS001",
						ActionCode:      "ACT001",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestCreateWallet,
						ActionType:      model.ActionCreate,
						ActionData:      testWalletEntity,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CPSAction{
				ID:              "CPS001",
				ActionCode:      "ACT001",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestCreateWallet,
				ActionType:      model.ActionCreate,
				ActionData:      testWalletEntity,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.CreateWallet(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ActionCode, got.ActionCode)
				assert.Equal(t, tt.want.ActionType, got.ActionType)
			}
		})
	}
}

func TestWalletDomain_GetAllWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testWallets := []*entity.Wallet{
		{
			ID:             "WALLET001",
			Name:           "TestWallet1",
			Avatar:         "test-bucket/avatar1.png",
			Code:           "TW001",
			Enabled:        true,
			IsDeleted:      false,
			CreatedAt:      testTime,
			LastModifiedAt: testTime,
			DeletedAt:      time.Time{},
		},
		{
			ID:             "WALLET002",
			Name:           "TestWallet2",
			Avatar:         "test-bucket/avatar2.png",
			Code:           "TW002",
			Enabled:        true,
			IsDeleted:      false,
			CreatedAt:      testTime,
			LastModifiedAt: testTime,
			DeletedAt:      time.Time{},
		},
	}

	// Note: The mock returns BankResponse, but we'll convert it to WalletResponse for testing
	testWalletResponse := &entity.WalletResponse{
		Page:    1,
		Wallets: testWallets,
		Limit:   10,
		Total:   2,
	}

	tests := []struct {
		name         string
		filterParams *constant.Filter
		mock         func()
		want         *entity.WalletResponse
		wantErr      bool
	}{
		{
			name: "success",
			filterParams: &constant.Filter{
				Page: 1,
			},
			mock: func() {
				mockRepo.EXPECT().
					GetAllWallet(gomock.Any(), &constant.Filter{Page: 1}).
					Return(testWalletResponse, nil)
			},
			want: &entity.WalletResponse{
				Page:    1,
				Wallets: testWallets,
				Limit:   10,
				Total:   2,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			filterParams: &constant.Filter{
				Page: 1,
			},
			mock: func() {
				mockRepo.EXPECT().
					GetAllWallet(gomock.Any(), &constant.Filter{Page: 1}).
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.GetAllWallet(context.Background(), tt.filterParams)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWalletDomain_GetWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testWallet := &entity.Wallet{
		ID:             "WALLET001",
		Name:           "TestWallet",
		Avatar:         "test-bucket/avatar.png",
		Code:           "TW001",
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      testTime,
		LastModifiedAt: testTime,
		DeletedAt:      time.Time{},
	}

	tests := []struct {
		name    string
		id      string
		mock    func()
		want    *entity.Wallet
		wantErr bool
	}{
		{
			name: "success",
			id:   "WALLET001",
			mock: func() {
				mockRepo.EXPECT().
					GetWallet(gomock.Any(), "WALLET001").
					Return(testWallet, nil)
			},
			want:    testWallet,
			wantErr: false,
		},
		{
			name: "error from database",
			id:   "WALLET001",
			mock: func() {
				mockRepo.EXPECT().
					GetWallet(gomock.Any(), "WALLET001").
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.GetWallet(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWalletDomain_UpdateWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	testCPSAction := model.CreateCPSAction{
		ActionCode:    "ACT002",
		MakerUser:     testUser,
		Department:    "IT",
		Status:        model.ActionPending,
		RequestAction: model.RequestUpdateWallet,
		ActionType:    model.ActionUpdate,
		ActionData: entity.Wallet{
			ID:             "WALLET001",
			Name:           "UpdatedWallet",
			Avatar:         "test-bucket/updated-avatar.png",
			Code:           "UW001",
			Enabled:        true,
			IsDeleted:      false,
			CreatedAt:      testTime,
			LastModifiedAt: testTime,
			DeletedAt:      time.Time{},
		},
		MakerActionTime: testTime,
	}

	tests := []struct {
		name    string
		id      string
		req     model.CreateCPSAction
		mock    func()
		want    *model.CPSAction
		wantErr bool
	}{
		{
			name: "success",
			id:   "WALLET001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil) // Assuming no existing action with the same code
				mockRepo.EXPECT().
					UpdateWallet(gomock.Any(), "WALLET001", testCPSAction).
					Return(&model.CPSAction{
						ID:              "CPS002",
						ActionCode:      "ACT002",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestUpdateWallet,
						ActionType:      model.ActionUpdate,
						ActionData:      testCPSAction.ActionData,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CPSAction{
				ID:              "CPS002",
				ActionCode:      "ACT002",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestUpdateWallet,
				ActionType:      model.ActionUpdate,
				ActionData:      testCPSAction.ActionData,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			id:   "WALLET001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(errors.New("action already exists"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.UpdateWallet(context.Background(), tt.id, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWalletDomain_DeleteWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	testCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT003",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestDeleteWallet,
		ActionType:      model.ActionDelete,
		ActionData:      nil,
		MakerActionTime: testTime,
	}

	tests := []struct {
		name    string
		id      string
		req     model.CreateCPSAction
		mock    func()
		want    *model.CPSAction
		wantErr bool
	}{
		{
			name: "success",
			id:   "WALLET001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					DeleteWallet(gomock.Any(), "WALLET001", testCPSAction).
					Return(&model.CPSAction{
						ID:              "CPS003",
						ActionCode:      "ACT003",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestDeleteWallet,
						ActionType:      model.ActionDelete,
						ActionData:      nil,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CPSAction{
				ID:              "CPS003",
				ActionCode:      "ACT003",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestDeleteWallet,
				ActionType:      model.ActionDelete,
				ActionData:      nil,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			id:   "WALLET001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					DeleteWallet(gomock.Any(), "WALLET001", testCPSAction).
					Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.DeleteWallet(context.Background(), tt.id, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWalletDomain_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testMakerUser := model.User{
		UserCode:    "MAKER001",
		FullName:    "Maker User",
		PhoneNumber: "+251911234567",
	}

	testCheckerUser := model.User{
		UserCode:    "CHECKER001",
		FullName:    "Checker User",
		PhoneNumber: "+251911234568",
	}

	testAuthorizeReq := model.AuthorizeCPSAction{
		ActionCode:        "ACT001",
		Department:        "IT",
		CheckerUser:       testCheckerUser,
		CheckerActionTime: testTime,
	}

	tests := []struct {
		name    string
		req     model.AuthorizeCPSAction
		mock    func()
		want    *model.CPSAction
		wantErr bool
	}{
		{
			name: "success",
			req:  testAuthorizeReq,
			mock: func() {
				mockRepo.EXPECT().
					Authorize(gomock.Any(), testAuthorizeReq).
					Return(&model.CPSAction{
						ID:                "CPS001",
						ActionCode:        "ACT001",
						MakerUser:         testMakerUser,
						CheckerUser:       testCheckerUser,
						Department:        "IT",
						Status:            model.ActionApproved,
						RequestAction:     model.RequestCreateBank,
						ActionType:        model.ActionCreate,
						ActionData:        nil,
						MakerActionTime:   testTime,
						CheckerActionTime: testTime,
					}, nil)
			},
			want: &model.CPSAction{
				ID:                "CPS001",
				ActionCode:        "ACT001",
				MakerUser:         testMakerUser,
				CheckerUser:       testCheckerUser,
				Department:        "IT",
				Status:            model.ActionApproved,
				RequestAction:     model.RequestCreateBank,
				ActionType:        model.ActionCreate,
				ActionData:        nil,
				MakerActionTime:   testTime,
				CheckerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			req:  testAuthorizeReq,
			mock: func() {
				mockRepo.EXPECT().
					Authorize(gomock.Any(), testAuthorizeReq).
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.Authorize(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWalletDomain_Reject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testMakerUser := model.User{
		UserCode:    "MAKER001",
		FullName:    "Maker User",
		PhoneNumber: "+251911234567",
	}

	testCheckerUser := model.User{
		UserCode:    "CHECKER001",
		FullName:    "Checker User",
		PhoneNumber: "+251911234568",
	}

	testRejectReq := model.RejectCPSAction{
		CreateCPSAction: model.CreateCPSAction{
			ActionCode:      "ACT001",
			MakerUser:       testMakerUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestCreateBank,
			ActionType:      model.ActionCreate,
			ActionData:      nil,
			MakerActionTime: testTime,
		},
		CheckerUser:       testCheckerUser,
		RejectedReason:    "Invalid wallet information Invalid wallet information",
		CheckerActionTime: testTime,
	}

	tests := []struct {
		name    string
		req     model.RejectCPSAction
		mock    func()
		want    *model.CPSAction
		wantErr bool
	}{
		{
			name: "success",
			req:  testRejectReq,
			mock: func() {
				mockRepo.EXPECT().
					Reject(gomock.Any(), testRejectReq).
					Return(&model.CPSAction{
						ID:                "CPS001",
						ActionCode:        "ACT001",
						MakerUser:         testMakerUser,
						CheckerUser:       testCheckerUser,
						Department:        "IT",
						Status:            model.ActionRejected,
						RequestAction:     model.RequestCreateBank,
						ActionType:        model.ActionCreate,
						ActionData:        nil,
						RejectedReason:    "Invalid wallet information",
						MakerActionTime:   testTime,
						CheckerActionTime: testTime,
					}, nil)
			},
			want: &model.CPSAction{
				ID:                "CPS001",
				ActionCode:        "ACT001",
				MakerUser:         testMakerUser,
				CheckerUser:       testCheckerUser,
				Department:        "IT",
				Status:            model.ActionRejected,
				RequestAction:     model.RequestCreateBank,
				ActionType:        model.ActionCreate,
				ActionData:        nil,
				RejectedReason:    "Invalid wallet information",
				MakerActionTime:   testTime,
				CheckerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			req:  testRejectReq,
			mock: func() {
				mockRepo.EXPECT().
					Reject(gomock.Any(), testRejectReq).
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.Reject(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWalletDomain_EnableOrDisableWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	walletService := service.InitWalletDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	testEnableCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT004",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestEnableWallet,
		ActionType:      model.ActionUpdate,
		ActionData:      nil,
		MakerActionTime: testTime,
	}

	testDisableCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT004",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestDisableWallet,
		ActionType:      model.ActionUpdate,
		ActionData:      nil,
		MakerActionTime: testTime,
	}

	tests := []struct {
		name          string
		id            string
		requestAction model.RequestAction
		cpsReq        model.CreateCPSAction
		mock          func()
		want          *model.CPSAction
		wantErr       bool
	}{
		{
			name:          "enable wallet success",
			id:            "WALLET001",
			requestAction: model.RequestEnableWallet,
			cpsReq:        testEnableCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testEnableCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					EnableOrDisableWallet(gomock.Any(), "WALLET001", model.RequestEnableWallet, testEnableCPSAction).
					Return(&model.CPSAction{
						ID:              "CPS004",
						ActionCode:      "ACT004",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestEnableWallet,
						ActionType:      model.ActionUpdate,
						ActionData:      nil,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CPSAction{
				ID:              "CPS004",
				ActionCode:      "ACT004",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestEnableWallet,
				ActionType:      model.ActionUpdate,
				ActionData:      nil,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name:          "disable wallet success",
			id:            "WALLET001",
			requestAction: model.RequestDisableWallet,
			cpsReq:        testDisableCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testDisableCPSAction).
					Return(nil) // Assuming no existing action with the same code
				mockRepo.EXPECT().
					EnableOrDisableWallet(gomock.Any(), "WALLET001", model.RequestDisableWallet, testDisableCPSAction).
					Return(&model.CPSAction{
						ID:              "CPS004",
						ActionCode:      "ACT004",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestDisableWallet,
						ActionType:      model.ActionUpdate,
						ActionData:      nil,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CPSAction{
				ID:              "CPS004",
				ActionCode:      "ACT004",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestDisableWallet,
				ActionType:      model.ActionUpdate,
				ActionData:      nil,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name:          "error from database",
			id:            "WALLET001",
			requestAction: model.RequestEnableWallet,
			cpsReq:        testEnableCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testEnableCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					EnableOrDisableWallet(gomock.Any(), "WALLET001", model.RequestEnableWallet, testEnableCPSAction).
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := walletService.EnableOrDisableWallet(context.Background(), tt.id, tt.requestAction, tt.cpsReq)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
