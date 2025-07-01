package service_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"cbe-super-app-cps-action/internal/adapter/outbound/model"
	"cbe-super-app-cps-action/internal/domain/bank/dto"
	"cbe-super-app-cps-action/internal/domain/bank/entity"
	"cbe-super-app-cps-action/internal/domain/bank/mocks"
	"cbe-super-app-cps-action/internal/domain/bank/service"
	constant "cbe-super-app-cps-action/utils"

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

func TestBankDomain_CreateOneBank(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	// Create test data with full struct values
	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	testCreateBankRequest := dto.CreateBankRequest{
		Name: "TestBank",
		Logo: createMockFileHeader("test-logo.png", []byte("test-logo"), 1024),
		Code: "TB001",
		BIC:  "TESTBIC123",
	}

	testCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT001",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestCreateBank,
		ActionType:      model.ActionCreate,
		ActionData:      testCreateBankRequest,
		MakerActionTime: testTime,
	}

	testBankEntity := entity.Bank{
		ID:             "BANK001",
		Name:           "TestBank",
		Logo:           "test-bucket/test-key",
		Code:           "TB001",
		BIC:            "TESTBIC123",
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      testTime,
		LastModifiedAt: testTime,
	}

	tests := []struct {
		name    string
		req     model.CreateCPSAction
		mock    func()
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					CreateBank(gomock.Any(), gomock.Any()).
					Return(&model.CpsAction{
						ID:              "CPS001",
						ActionCode:      "ACT001",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestCreateBank,
						ActionType:      model.ActionCreate,
						ActionData:      testBankEntity,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CpsAction{
				ID:              "CPS001",
				ActionCode:      "ACT001",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestCreateBank,
				ActionType:      model.ActionCreate,
				ActionData:      testBankEntity,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from repository",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					CreateBank(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := bankService.CreateOneBank(context.Background(), tt.req)
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

func TestBankDomain_GetAllBank(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testBanks := []*entity.Bank{
		{
			ID:             "BANK001",
			Name:           "TestBank1",
			Logo:           "test-bucket/logo1.png",
			Code:           "TB001",
			BIC:            "TESTBIC123",
			Enabled:        true,
			IsDeleted:      false,
			CreatedAt:      testTime,
			LastModifiedAt: testTime,
		},
		{
			ID:             "BANK002",
			Name:           "TestBank2",
			Logo:           "test-bucket/logo2.png",
			Code:           "TB002",
			BIC:            "TESTBIC456",
			Enabled:        true,
			IsDeleted:      false,
			CreatedAt:      testTime,
			LastModifiedAt: testTime,
		},
	}

	tests := []struct {
		name         string
		filterParams *constant.Filter
		mock         func()
		want         *entity.BankResponse
		wantErr      bool
	}{
		{
			name: "success",
			filterParams: &constant.Filter{
				Page: 1,
			},
			mock: func() {
				mockRepo.EXPECT().
					GetAllBanks(gomock.Any(), &constant.Filter{Page: 1}).
					Return(&entity.BankResponse{
						Page:  1,
						Banks: testBanks,
						Limit: 10,
						Total: 2,
					}, nil)
			},
			want: &entity.BankResponse{
				Page:  1,
				Banks: testBanks,
				Limit: 10,
				Total: 2,
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
					GetAllBanks(gomock.Any(), &constant.Filter{Page: 1}).
					Return(nil, errors.New("internal server error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := bankService.GetAllBank(context.Background(), tt.filterParams)
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

func TestBankDomain_GetOneBank(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testBank := &entity.Bank{
		ID:             "BANK001",
		Name:           "TestBank",
		Logo:           "test-bucket/logo.png",
		Code:           "TB001",
		BIC:            "TESTBIC123",
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      testTime,
		LastModifiedAt: testTime,
	}

	tests := []struct {
		name    string
		id      string
		mock    func()
		want    *entity.Bank
		wantErr bool
	}{
		{
			name: "success",
			id:   "BANK001",
			mock: func() {
				mockRepo.EXPECT().
					GetBank(gomock.Any(), "BANK001").
					Return(testBank, nil)
			},
			want:    testBank,
			wantErr: false,
		},
		{
			name: "error from database",
			id:   "BANK001",
			mock: func() {
				mockRepo.EXPECT().
					GetBank(gomock.Any(), "BANK001").
					Return(nil, errors.New("internal server erro"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := bankService.GetOneBank(context.Background(), tt.id)
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

func TestBankDomain_UpdateOneBank(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

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
		RequestAction: model.RequestUpdateBank,
		ActionType:    model.ActionUpdate,
		ActionData: entity.Bank{
			ID:             "BANK001",
			Name:           "UpdatedBank",
			Logo:           "test-bucket/updated-logo.png",
			Code:           "UB001",
			BIC:            "UPDATEDBIC123",
			Enabled:        true,
			IsDeleted:      false,
			CreatedAt:      testTime,
			LastModifiedAt: testTime,
		},
		MakerActionTime: testTime,
	}

	tests := []struct {
		name    string
		id      string
		req     model.CreateCPSAction
		mock    func()
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			id:   "BANK001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					UpdateBank(gomock.Any(), "BANK001", testCPSAction).
					Return(&model.CpsAction{
						ID:              "CPS002",
						ActionCode:      "ACT002",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestUpdateBank,
						ActionType:      model.ActionUpdate,
						ActionData:      testCPSAction.ActionData,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CpsAction{
				ID:              "CPS002",
				ActionCode:      "ACT002",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestUpdateBank,
				ActionType:      model.ActionUpdate,
				ActionData:      testCPSAction.ActionData,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			id:   "BANK001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					UpdateBank(gomock.Any(), "BANK001", testCPSAction).
					Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := bankService.UpdateOneBank(context.Background(), tt.id, tt.req)
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

func TestBankDomain_DeleteOneBank(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

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
		RequestAction:   model.RequestDeleteBank,
		ActionType:      model.ActionDelete,
		ActionData:      nil,
		MakerActionTime: testTime,
	}

	tests := []struct {
		name    string
		id      string
		req     model.CreateCPSAction
		mock    func()
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			id:   "BANK001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					DeleteBank(gomock.Any(), "BANK001", testCPSAction).
					Return(&model.CpsAction{
						ID:              "CPS003",
						ActionCode:      "ACT003",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestDeleteBank,
						ActionType:      model.ActionDelete,
						ActionData:      nil,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CpsAction{
				ID:              "CPS003",
				ActionCode:      "ACT003",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestDeleteBank,
				ActionType:      model.ActionDelete,
				ActionData:      nil,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name: "error from database",
			id:   "BANK001",
			req:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					DeleteBank(gomock.Any(), "BANK001", testCPSAction).
					Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := bankService.DeleteOneBank(context.Background(), tt.id, tt.req)
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

func TestBankDomain_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

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
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			req:  testAuthorizeReq,
			mock: func() {
				mockRepo.EXPECT().
					Authorize(gomock.Any(), testAuthorizeReq).
					Return(&model.CpsAction{
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
			want: &model.CpsAction{
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
			got, err := bankService.Authorize(context.Background(), tt.req)
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

func TestBankDomain_Reject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

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
		RejectedReason:    "Invalid bank information Invalid bank information Invalid bank information Invalid bank information",
		CheckerActionTime: testTime,
	}

	tests := []struct {
		name    string
		req     model.RejectCPSAction
		mock    func()
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			req:  testRejectReq,
			mock: func() {
				mockRepo.EXPECT().
					Reject(gomock.Any(), testRejectReq).
					Return(&model.CpsAction{
						ID:                "CPS001",
						ActionCode:        "ACT001",
						MakerUser:         testMakerUser,
						CheckerUser:       testCheckerUser,
						Department:        "IT",
						Status:            model.ActionRejected,
						RequestAction:     model.RequestCreateBank,
						ActionType:        model.ActionCreate,
						ActionData:        nil,
						RejectedReason:    "Invalid bank information Invalid bank information Invalid bank information Invalid bank information",
						MakerActionTime:   testTime,
						CheckerActionTime: testTime,
					}, nil)
			},
			want: &model.CpsAction{
				ID:                "CPS001",
				ActionCode:        "ACT001",
				MakerUser:         testMakerUser,
				CheckerUser:       testCheckerUser,
				Department:        "IT",
				Status:            model.ActionRejected,
				RequestAction:     model.RequestCreateBank,
				ActionType:        model.ActionCreate,
				ActionData:        nil,
				RejectedReason:    "Invalid bank information Invalid bank information Invalid bank information Invalid bank information",
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
					Return(nil, errors.New("internal server eerror"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := bankService.Reject(context.Background(), tt.req)
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

func TestBankDomain_EnableOrDisableBank(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bankService := service.InitBankDomain(mockRepo, mockMinioClient, "test-bucket", mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	testCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT004",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestEnableBank,
		ActionType:      model.ActionUpdate,
		ActionData:      nil,
		MakerActionTime: testTime,
	}

	testDisableCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT004",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestDisableBank,
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
		want          *model.CpsAction
		wantErr       bool
	}{
		{
			name:          "enable wallet success",
			id:            "WALLET001",
			requestAction: model.RequestEnableBank,
			cpsReq:        testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					EnableOrDisableBank(gomock.Any(), "WALLET001", model.RequestEnableBank, testCPSAction).
					Return(&model.CpsAction{
						ID:              "CPS004",
						ActionCode:      "ACT004",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestEnableBank,
						ActionType:      model.ActionUpdate,
						ActionData:      nil,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CpsAction{
				ID:              "CPS004",
				ActionCode:      "ACT004",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestEnableBank,
				ActionType:      model.ActionUpdate,
				ActionData:      nil,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name:          "disable wallet success",
			id:            "WALLET001",
			requestAction: model.RequestDisableBank,
			cpsReq:        testDisableCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testDisableCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					EnableOrDisableBank(gomock.Any(), "WALLET001", model.RequestDisableBank, testDisableCPSAction).
					Return(&model.CpsAction{
						ID:              "CPS004",
						ActionCode:      "ACT004",
						MakerUser:       testUser,
						Department:      "IT",
						Status:          model.ActionPending,
						RequestAction:   model.RequestDisableBank,
						ActionType:      model.ActionUpdate,
						ActionData:      nil,
						MakerActionTime: testTime,
					}, nil)
			},
			want: &model.CpsAction{
				ID:              "CPS004",
				ActionCode:      "ACT004",
				MakerUser:       testUser,
				Department:      "IT",
				Status:          model.ActionPending,
				RequestAction:   model.RequestDisableBank,
				ActionType:      model.ActionUpdate,
				ActionData:      nil,
				MakerActionTime: testTime,
			},
			wantErr: false,
		},
		{
			name:          "error from databas",
			id:            "WALLET001",
			requestAction: model.RequestEnableBank,
			cpsReq:        testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					CPSActionExists(gomock.Any(), testCPSAction).
					Return(nil)
				mockRepo.EXPECT().
					EnableOrDisableBank(gomock.Any(), "WALLET001", model.RequestEnableBank, testCPSAction).
					Return(nil, errors.New("internal server eror"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := bankService.EnableOrDisableBank(context.Background(), tt.id, tt.requestAction, tt.cpsReq)
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
