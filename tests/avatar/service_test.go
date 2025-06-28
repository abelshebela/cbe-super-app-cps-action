package avatar

import (
	"context"
	"fmt"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/minio/minio-go/v7"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/avatar"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/tests/avatar/mocks"
	sharedconfig "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

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
			Avatar: &multipart.FileHeader{Filename: "avatar.png", Size: 1024},
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
			Avatar: "test-bucket/test-key.png",
			Label:  "Test Avatar",
			Enable: true,
		}
		expectedCpsAction := model.CpsAction{
			ID:              "CPS001",
			ActionCode:      "ACT001",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestCreateAvatar,
			ActionType:      model.ActionCreate,
			ActionData:      expectedAvatar,
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
			Avatar: &multipart.FileHeader{Filename: "avatar.png", Size: 1024},
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
		assert.Equal(t, model.CpsAction{}, result)
	})

	t.Run("validation error file size should be less than 2MB", func(t *testing.T) {
		testCreateAvatar := avatar.CreateAvatar{
			Label:  "Test Avatar",
			Avatar: &multipart.FileHeader{Filename: "avatar.png", Size: 3 << 20}, // 3MB, triggers validation error
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
		assert.Equal(t, model.CpsAction{}, result)
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
		expectedCpsAction := model.CpsAction{
			ID:              "CPS003",
			ActionCode:      "ACT003",
			MakerUser:       testUser,
			Department:      "IT",
			Status:          "PENDING",
			RequestAction:   "DELETE_AVATAR",
			ActionType:      "DELETE",
			ActionData:      nil,
			MakerActionTime: testTime,
		}

		mockRepo.EXPECT().CPSActionExists(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().DeleteAvatar(gomock.Any(), avatarID, gomock.Any()).Return(expectedCpsAction, nil)

		result, err := avatarService.DeleteAvatar(context.Background(), avatarID, testCPSAction)
		assert.NoError(t, err)
		assert.Equal(t, expectedCpsAction, result)
	})

	t.Run("CPS action exists error", func(t *testing.T) {
		avatarID := "AVATAR002"
		testCPSAction := model.CreateCPSAction{
			ActionCode:      "ACT004",
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
		assert.Equal(t, model.CpsAction{}, result)
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
		mockRepo.EXPECT().DeleteAvatar(gomock.Any(), avatarID, gomock.Any()).Return(model.CpsAction{}, fmt.Errorf("avatar not found"))

		result, err := avatarService.DeleteAvatar(context.Background(), avatarID, testCPSAction)
		assert.Error(t, err)
		assert.Equal(t, model.CpsAction{}, result)
		assert.Contains(t, err.Error(), "avatar not found")
	})
}
