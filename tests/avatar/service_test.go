package avatar

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/tests/avatar/mocks"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/minio/minio-go/v7"
	sharedconfig "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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

type MockMinioClient struct{}

func (m *MockMinioClient) SaveObjectN(ctx context.Context, obj sharedconfig.SaveObjectBodyN) (*minio.UploadInfo, error) {
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
func (m *MockMinioClient) SaveObject(ctx context.Context, body sharedconfig.SaveObjectBody) (*sharedconfig.SaveObjectResponse, error) {
	return &sharedconfig.SaveObjectResponse{
		Bucket: "test-bucket",
		Key:    "test-key",
	}, nil
}
func (m *MockMinioClient) GetObject(ctx context.Context, bucketName string, objectName string) (interface{}, error) {
	return nil, nil
}
func (m *MockMinioClient) DeleteObject(ctx context.Context, body sharedconfig.DeleteObjectBody) (bool, error) {
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

func TestAvatarDomain_CreateAvatar(t *testing.T) {
	// Mock IsValidImage to always return true for tests
	originalIsValidImage := avatar.IsValidImage
	avatar.IsValidImage = func(fileHeader *multipart.FileHeader) bool { return true }
	defer func() { avatar.IsValidImage = originalIsValidImage }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	t.Run("success", func(t *testing.T) {
		testCreateAvatar := avatar.CreateAvatar{
			Label:  "Test Avatar",
			Avatar: createMockFileHeader("avatar.png", []byte("mock avatar image"), 1024),
		}
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT002",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestCreateAvatar,
			ActionType:      model.ActionCreate,
			ActionData:      testCreateAvatar,
			MakerActionTime: testTime,
		}

		expectedAvatar := avatar.Avatar{
			ID:     "AVATAR001",
			Avatar: "test-bucket/test-key",
			Label:  "Test Avatar",
			Enable: true,
		}
		expectedCpsAction := model.CPSAction{
			ID:              bson.NewObjectID(),
			ActionCode:      "ACT001",
			MakerName:       testUser.FullName,
			Department:      "IT",
			ActionStatus:    string(model.ActionPending),
			RequestAction:   string(model.RequestCreateAvatar),
			ActionType:      string(model.ActionCreate),
			CurrentAction:   expectedAvatar,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().CreateAvatar(gomock.Any(), gomock.Any()).Return(expectedCpsAction, nil)

		result, err := avatarService.CreateAvatar(context.Background(), testCPSAction)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("validation error", func(t *testing.T) {
		testCreateAvatar := avatar.CreateAvatar{
			Label:  "",
			Avatar: createMockFileHeader("avatar.png", []byte("mock avatar image"), 1024),
		}
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT002",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          "PENDING",
			RequestAction:   "CREATE_AVATAR",
			ActionType:      "CREATE",
			ActionData:      testCreateAvatar,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)

		result, err := avatarService.CreateAvatar(context.Background(), testCPSAction)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
	})

	t.Run("validation error file size should be less than 2MB", func(t *testing.T) {
		// Create content larger than 2MB to trigger validation error
		largeContent := make([]byte, 3<<20) // 3MB
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}

		testCreateAvatar := avatar.CreateAvatar{
			Label:  "Test Avatar",
			Avatar: createMockFileHeader("avatar.png", largeContent, 3<<20), // 3MB, triggers validation error
		}
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT002",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          "PENDING",
			RequestAction:   "CREATE_AVATAR",
			ActionType:      "CREATE",
			ActionData:      testCreateAvatar,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)

		result, err := avatarService.CreateAvatar(context.Background(), testCPSAction)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file size should be less than 2MB")
		assert.Equal(t, model.CPSAction{}, result)
	})
}

func TestAvatarDomain_DeleteAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	t.Run("success", func(t *testing.T) {
		avatarID := "AVATAR001"
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT003",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          "PENDING",
			RequestAction:   "DELETE_AVATAR",
			ActionType:      "DELETE",
			ActionData:      nil,
			MakerActionTime: testTime,
		}
		expectedCpsAction := model.CPSAction{

			ID:               bson.NewObjectID(),
			ActionCode:       "ACT003",
			MakerID:          testUser.UserCode,
			MakerName:        testUser.FullName,
			MakerPhoneNumber: testUser.PhoneNumber,
			Department:       "IT",
			ActionStatus:     "PENDING",
			RequestAction:    "DELETE_AVATAR",
			ActionType:       "DELETE",
			CurrentAction:    nil,
			MakerActionTime:  testTime,
		}

		result, err := avatarService.DeleteAvatar(context.Background(), avatarID, testCPSAction)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("CPS action exists error", func(t *testing.T) {
		avatarID := "AVATAR002"
		testCPSAction := model.CreateCPSAction{
			MakerUser:       testUser,
			Department:      "IT",
			Status:          "PENDING",
			RequestAction:   "DELETE_AVATAR",
			ActionType:      "DELETE",
			ActionData:      nil,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(fmt.Errorf("CPS action already exists"))

		result, err := avatarService.DeleteAvatar(context.Background(), avatarID, testCPSAction)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
		assert.Contains(t, err.Error(), "CPS action already exists")
	})

	t.Run("delete avatar repository error", func(t *testing.T) {
		avatarID := "AVATAR003"
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT005",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          "PENDING",
			RequestAction:   "DELETE_AVATAR",
			ActionType:      "DELETE",
			ActionData:      nil,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().DeleteAvatar(gomock.Any(), avatarID, gomock.Any()).Return(model.CPSAction{}, fmt.Errorf("avatar not found"))

		result, err := avatarService.DeleteAvatar(context.Background(), avatarID, testCPSAction)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
		assert.Contains(t, err.Error(), "avatar not found")
	})
}

func TestAvatarDomain_UpdateAvatar(t *testing.T) {
	// Mock IsValidImage to always return true for tests
	originalIsValidImage := avatar.IsValidImage
	avatar.IsValidImage = func(fileHeader *multipart.FileHeader) bool { return true }
	defer func() { avatar.IsValidImage = originalIsValidImage }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	t.Run("success", func(t *testing.T) {
		avatarID := "AVATAR001"
		testUpdateAvatar := avatar.UpdateAvatar{
			Avatar: createMockFileHeader("updated-avatar.png", []byte("updated avatar content"), 1024),
		}
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT006",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestUpdateAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      testUpdateAvatar,
			MakerActionTime: testTime,
		}

		existingAvatar := &avatar.Avatar{
			ID:     avatarID,
			Avatar: "test-bucket/old-avatar.png",
			Label:  "Old Avatar",
			Enable: true,
		}

		expectedAvatar := avatar.Avatar{
			ID:     avatarID,
			Avatar: "test-bucket/test-key",
			Label:  "Old Avatar",
			Enable: true,
		}
		expectedCpsAction := model.CPSAction{
			ID:               bson.NewObjectID(),
			ActionCode:       "ACT006",
			MakerID:          testUser.UserCode,
			MakerName:        testUser.FullName,
			MakerPhoneNumber: testUser.PhoneNumber,
			Department:       "IT",
			ActionStatus:     string(model.ActionPending),
			RequestAction:    string(model.RequestUpdateAvatar),
			ActionType:       string(model.ActionUpdate),
			CurrentAction:    expectedAvatar,
			MakerActionTime:  testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().GetAvatar(gomock.Any(), avatarID).Return(existingAvatar, nil)
		mockRepo.EXPECT().UpdateAvatar(gomock.Any(), avatarID, gomock.Any()).Return(expectedCpsAction, nil)

		result, err := avatarService.UpdateAvatar(context.Background(), avatarID, testCPSAction)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("CPS action exists error", func(t *testing.T) {
		avatarID := "AVATAR002"
		testUpdateAvatar := avatar.UpdateAvatar{
			Avatar: createMockFileHeader("updated-avatar.png", []byte("updated avatar content"), 1024),
		}
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT007",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestUpdateAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      testUpdateAvatar,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(fmt.Errorf("CPS action already exists"))

		result, err := avatarService.UpdateAvatar(context.Background(), avatarID, testCPSAction)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
		assert.Contains(t, err.Error(), "CPS action already exists")
	})

	t.Run("get avatar not found", func(t *testing.T) {
		avatarID := "AVATAR003"
		testUpdateAvatar := avatar.UpdateAvatar{
			Avatar: createMockFileHeader("updated-avatar.png", []byte("updated avatar content"), 1024),
		}
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT008",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestUpdateAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      testUpdateAvatar,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().GetAvatar(gomock.Any(), avatarID).Return(nil, fmt.Errorf("avatar not found"))

		result, err := avatarService.UpdateAvatar(context.Background(), avatarID, testCPSAction)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
		assert.Contains(t, err.Error(), "avatar not found")
	})

	t.Run("validation error file size too large", func(t *testing.T) {
		avatarID := "AVATAR004"
		// Create content larger than 2MB to trigger validation error
		largeContent := make([]byte, 3<<20) // 3MB
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}

		testUpdateAvatar := avatar.UpdateAvatar{
			Avatar: createMockFileHeader("updated-avatar.png", largeContent, 3<<20),
		}
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT009",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestUpdateAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      testUpdateAvatar,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)

		result, err := avatarService.UpdateAvatar(context.Background(), avatarID, testCPSAction)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file size should be less than 2MB")
		assert.Equal(t, model.CPSAction{}, result)
	})
}

func TestAvatarDomain_GetAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	t.Run("success", func(t *testing.T) {
		avatarID := "AVATAR001"
		expectedAvatar := &avatar.Avatar{
			ID:             avatarID,
			Avatar:         "test-bucket/avatar.png",
			Label:          "Test Avatar",
			Enable:         true,
			IsDeleted:      false,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
			DeletedAt:      time.Time{},
		}

		mockRepo.EXPECT().GetAvatar(gomock.Any(), avatarID).Return(expectedAvatar, nil)

		result, err := avatarService.GetAvatar(context.Background(), avatarID)
		assert.NoError(t, err)
		assert.Equal(t, expectedAvatar, result)
	})

	t.Run("avatar not found", func(t *testing.T) {
		avatarID := "AVATAR002"

		mockRepo.EXPECT().GetAvatar(gomock.Any(), avatarID).Return(nil, fmt.Errorf("avatar not found"))

		result, err := avatarService.GetAvatar(context.Background(), avatarID)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "avatar not found")
	})
}

func TestAvatarDomain_GetAllAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	t.Run("success", func(t *testing.T) {
		filterParams := constant.Filter{
			Page:    1,
			PerPage: 10,
			Search:  "",
			Filters: "",
		}

		avatars := []*avatar.Avatar{
			{
				ID:             "AVATAR001",
				Avatar:         "test-bucket/avatar1.png",
				Label:          "Avatar 1",
				Enable:         true,
				IsDeleted:      false,
				CreatedAt:      time.Now(),
				LastModifiedAt: time.Now(),
				DeletedAt:      time.Time{},
			},
			{
				ID:             "AVATAR002",
				Avatar:         "test-bucket/avatar2.png",
				Label:          "Avatar 2",
				Enable:         false,
				IsDeleted:      false,
				CreatedAt:      time.Now(),
				LastModifiedAt: time.Now(),
				DeletedAt:      time.Time{},
			},
		}

		expectedResponse := avatar.AvatarResponse{
			Page:    1,
			Avatars: avatars,
			Limit:   10,
			Total:   2,
		}

		mockRepo.EXPECT().GetAllAvatar(gomock.Any(), filterParams).Return(expectedResponse, nil)

		result, err := avatarService.GetAllAvatar(context.Background(), filterParams)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
	})

	t.Run("no avatars found", func(t *testing.T) {
		filterParams := constant.Filter{
			Page:    1,
			PerPage: 10,
			Search:  "nonexistent",
			Filters: "",
		}

		expectedResponse := avatar.AvatarResponse{
			Page:    1,
			Avatars: []*avatar.Avatar{},
			Limit:   10,
			Total:   0,
		}

		mockRepo.EXPECT().GetAllAvatar(gomock.Any(), filterParams).Return(expectedResponse, nil)

		result, err := avatarService.GetAllAvatar(context.Background(), filterParams)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		assert.Empty(t, result.Avatars)
	})

	t.Run("repository error", func(t *testing.T) {
		filterParams := constant.Filter{
			Page:    1,
			PerPage: 10,
			Search:  "",
			Filters: "",
		}

		mockRepo.EXPECT().GetAllAvatar(gomock.Any(), filterParams).Return(avatar.AvatarResponse{}, fmt.Errorf("database error"))

		result, err := avatarService.GetAllAvatar(context.Background(), filterParams)
		assert.Error(t, err)
		assert.Equal(t, avatar.AvatarResponse{}, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestAvatarDomain_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	testTime := time.Now()
	testCheckerUser := model.User{
		UserCode:    "CHECKER001",
		FullName:    "Jane Checker",
		PhoneNumber: "+251911234568",
	}

	t.Run("success", func(t *testing.T) {
		actionCode := "ACT001"
		testAuthorizeReq := model.AuthorizeCPSAction{
			ActionCode:  actionCode,
			CheckerUser: testCheckerUser,
			Department:  "IT",
		}

		expectedCpsAction := model.CPSAction{
			ID:                 bson.NewObjectID(),
			ActionCode:         actionCode,
			MakerID:            "USER001",
			MakerName:          "John Doe",
			MakerPhoneNumber:   "+251911234567",
			CheckerID:          testCheckerUser.UserCode,
			CheckerName:        testCheckerUser.FullName,
			CheckerPhoneNumber: testCheckerUser.PhoneNumber,
			Department:         "IT",
			ActionStatus:       string(model.ActionApproved),
			RequestAction:      string(model.RequestCreateAvatar),
			ActionType:         string(model.ActionCreate),
			CurrentAction:      avatar.Avatar{ID: "AVATAR001", Avatar: "test-bucket/avatar.png", Label: "Test Avatar", Enable: true},
			MakerActionTime:    testTime,
			CheckerActionTime:  testTime,
		}

		mockRepo.EXPECT().Authorize(gomock.Any(), testAuthorizeReq).Return(expectedCpsAction, nil)

		result, err := avatarService.Authorize(context.Background(), testAuthorizeReq)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("authorization failed", func(t *testing.T) {
		actionCode := "ACT002"
		testAuthorizeReq := model.AuthorizeCPSAction{
			ActionCode:  actionCode,
			CheckerUser: testCheckerUser,
			Department:  "IT",
		}

		mockRepo.EXPECT().Authorize(gomock.Any(), testAuthorizeReq).Return(model.CPSAction{}, fmt.Errorf("action not found"))

		result, err := avatarService.Authorize(context.Background(), testAuthorizeReq)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
		assert.Contains(t, err.Error(), "action not found")
	})
}

func TestAvatarDomain_Reject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	testTime := time.Now()
	testCheckerUser := model.User{
		UserCode:    "CHECKER001",
		FullName:    "Jane Checker",
		PhoneNumber: "+251911234568",
	}

	t.Run("success", func(t *testing.T) {
		actionCode := "ACT001"
		testRejectReq := model.RejectCPSAction{
			CreateCPSAction: model.CreateCPSAction{
				ActionCode: actionCode,
				Department: "IT",
			},
			CheckerUser:    testCheckerUser,
			RejectedReason: "The uploaded image does not meet the required specifications and quality standards.",
		}

		expectedCpsAction := model.CPSAction{
			ID:                 bson.NewObjectID(),
			ActionCode:         actionCode,
			MakerID:            "USER001",
			MakerName:          "John Doe",
			MakerPhoneNumber:   "+251911234567",
			CheckerID:          testCheckerUser.UserCode,
			CheckerName:        testCheckerUser.FullName,
			CheckerPhoneNumber: testCheckerUser.PhoneNumber,
			Department:         "IT",
			ActionStatus:       string(model.ActionRejected),
			RequestAction:      string(model.RequestCreateAvatar),
			ActionType:         string(model.ActionCreate),
			CurrentAction:      avatar.Avatar{ID: "AVATAR001", Avatar: "test-bucket/avatar.png", Label: "Test Avatar", Enable: true},
			RejectionReason:    "The uploaded image does not meet the required specifications and quality standards.",
			MakerActionTime:    testTime,
			CheckerActionTime:  testTime,
		}

		mockRepo.EXPECT().Reject(gomock.Any(), testRejectReq).Return(expectedCpsAction, nil)

		result, err := avatarService.Reject(context.Background(), testRejectReq)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("validation error - missing rejected reason", func(t *testing.T) {
		actionCode := "ACT002"
		testRejectReq := model.RejectCPSAction{
			CreateCPSAction: model.CreateCPSAction{
				ActionCode: actionCode,
				Department: "IT",
			},
			CheckerUser:    testCheckerUser,
			RejectedReason: "", // Empty reason should fail validation
		}

		result, err := avatarService.Reject(context.Background(), testRejectReq)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
	})

	t.Run("rejection failed", func(t *testing.T) {
		actionCode := "ACT003"
		testRejectReq := model.RejectCPSAction{
			CreateCPSAction: model.CreateCPSAction{
				ActionCode: actionCode,
				Department: "IT",
			},
			CheckerUser:    testCheckerUser,
			RejectedReason: "Action not found. The requested action could not be located in the system.",
		}

		mockRepo.EXPECT().Reject(gomock.Any(), testRejectReq).Return(model.CPSAction{}, fmt.Errorf("action not found"))

		result, err := avatarService.Reject(context.Background(), testRejectReq)
		assert.Error(t, err)
		assert.Equal(t, model.CPSAction{}, result)
		assert.Contains(t, err.Error(), "action not found")
	})
}

func TestAvatarDomain_EnableOrDisableAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAvatarRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	bucketName := "test-bucket"
	avatarService := avatar.InitAvatarDomain(mockRepo, mockMinioClient, bucketName, mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}

	t.Run("enable avatar success", func(t *testing.T) {
		avatarID := "AVATAR001"
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT010",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestEnableAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      avatar.Avatar{ID: avatarID},
			MakerActionTime: testTime,
		}

		expectedCpsAction := &model.CPSAction{
			ID:               bson.NewObjectID(),
			ActionCode:       "ACT010",
			MakerID:          "USER001",
			MakerName:        "John Doe",
			MakerPhoneNumber: "+251911234567",
			Department:       "IT",
			ActionStatus:     string(model.ActionPending),
			RequestAction:    string(model.RequestEnableAvatar),
			ActionType:       string(model.ActionUpdate),
			CurrentAction:    avatar.Avatar{ID: avatarID},
			MakerActionTime:  testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().EnableOrDisableAvatar(gomock.Any(), avatarID, model.RequestEnableAvatar, gomock.Any()).Return(expectedCpsAction, nil)

		result, err := avatarService.EnableOrDisableAvatar(context.Background(), avatarID, model.RequestEnableAvatar, testCPSAction)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("disable avatar success", func(t *testing.T) {
		avatarID := "AVATAR002"
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT011",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestDisableAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      avatar.Avatar{ID: avatarID},
			MakerActionTime: testTime,
		}

		expectedCpsAction := &model.CPSAction{
			ID:               bson.NewObjectID(),
			ActionCode:       "ACT011",
			MakerID:          "USER001",
			MakerName:        "John Doe",
			MakerPhoneNumber: "+251911234567",
			Department:       "IT",
			ActionStatus:     string(model.ActionPending),
			RequestAction:    string(model.RequestDisableAvatar),
			ActionType:       string(model.ActionUpdate),
			CurrentAction:    avatar.Avatar{ID: avatarID},
			MakerActionTime:  testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().EnableOrDisableAvatar(gomock.Any(), avatarID, model.RequestDisableAvatar, gomock.Any()).Return(expectedCpsAction, nil)

		result, err := avatarService.EnableOrDisableAvatar(context.Background(), avatarID, model.RequestDisableAvatar, testCPSAction)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("CPS action exists error", func(t *testing.T) {
		avatarID := "AVATAR003"
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT012",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestEnableAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      avatar.Avatar{ID: avatarID},
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(fmt.Errorf("CPS action already exists"))

		result, err := avatarService.EnableOrDisableAvatar(context.Background(), avatarID, model.RequestEnableAvatar, testCPSAction)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "CPS action already exists")
	})

	t.Run("enable/disable operation failed", func(t *testing.T) {
		avatarID := "AVATAR004"
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT013",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestEnableAvatar,
			ActionType:      model.ActionUpdate,
			ActionData:      avatar.Avatar{ID: avatarID},
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().EnableOrDisableAvatar(gomock.Any(), avatarID, model.RequestEnableAvatar, gomock.Any()).Return(nil, fmt.Errorf("avatar not found"))

		result, err := avatarService.EnableOrDisableAvatar(context.Background(), avatarID, model.RequestEnableAvatar, testCPSAction)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "avatar not found")
	})
}
