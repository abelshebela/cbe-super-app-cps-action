package service

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/mocks"
	"github.com/minio/minio-go/v7"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	ad "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

// MockMinioClient implements config.MinioClientInterface for testing
type MockMinioClient struct{}

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

func TestADDomain_CreateOneAdvert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdRepository(ctrl)
	mockLogger := &MockLogger{}
	mockMinioClient := &MockMinioClient{}
	service := InitADDomian("ad", mockMinioClient, mockRepo, mockLogger)

	maker_action_time := time.Now()
	started_at := time.Now().Add(20 * time.Minute)
	expired_at := time.Now().Add(48 * time.Hour)
	tests := []struct {
		name      string
		cpsAction model.CreateCPSAction
		mock      func()
		want      *model.CPSAction
		wantErr   bool
	}{
		{
			name: "success",
			cpsAction: model.CreateCPSAction{
				ActionCode: "1234",
				MakerUser: model.User{
					UserCode:    "1234",
					FullName:    "Abebe Kebede",
					PhoneNumber: "251911111111",
				},
				Department: "IT",
				Status:     "PENDING",
				ActionType: "CREATE",
				ActionData: ad.CreateAdvertRequest{
					Title:       "Advert",
					Description: "AdvertDescriptiondescriptionlengthisbetweenand",
					BannerImage: createMockFileHeader("google.png", []byte("test"), 1024),
					AdvertFor:   ad.CB,
					Date: ad.AdvertDate{
						StartedAt: started_at,
						ExpiredAt: expired_at,
					},
				},
				MakerActionTime: maker_action_time,
			},
			mock: func() {
				mockRepo.EXPECT().
					CreateOneAdvert(gomock.Any(), model.CreateCPSAction{
						ActionCode: "1234",
						MakerUser: model.User{
							UserCode:    "1234",
							FullName:    "Abebe Kebede",
							PhoneNumber: "251911111111",
						},
						Department: "IT",
						Status:     model.ActionPending,
						ActionType: "CREATE",
						ActionData: entity.Advert{
							Title:       "Advert",
							Description: "Advert Description",
							BannerImage: "google.png",
							AdvertFor:   entity.CB,
							Date: entity.AdvertDate{
								StartedAt: started_at,
								ExpiredAt: expired_at,
							},
						},
						MakerActionTime: maker_action_time,
					}).
					Return(&model.CPSAction{
						MakerID:          "1234",
						MakerName:        "Abebe Kebede",
						MakerPhoneNumber: "251911111111",
						Department:       "IT",
						ActionStatus:     "PENDING",
						ActionType:       "CREATE",
						CurrentAction: entity.Advert{
							Title:       "Advert",
							Description: "Advert Description",
							BannerImage: "google.png",
							AdvertFor:   entity.CB,
							Date: entity.AdvertDate{
								StartedAt: time.Now(),
								ExpiredAt: time.Now().Add(4 * time.Hour),
							},
						},
						MakerActionTime: maker_action_time}, nil)
			},
			want: &model.CPSAction{
				MakerID:          "1234",
				MakerName:        "Abebe Kebede",
				MakerPhoneNumber: "251911111111",
				Department:       "IT",
				ActionStatus:     "PENDING",
				ActionType:       "CREATE",
				CurrentAction: entity.Advert{
					Title:       "Advert",
					Description: "Advert Description",
					BannerImage: "google.png",
					AdvertFor:   entity.CB,
					Date: entity.AdvertDate{
						StartedAt: time.Now(),
						ExpiredAt: time.Now().Add(4 * time.Hour),
					},
				},
				MakerActionTime: maker_action_time},
			wantErr: false,
		},
		// {
		// 	name: "error from repository",
		// 	cpsAction: model.CreateCPSAction{
		// 		ActionCode: "1234",
		// 		MakerUser: model.User{
		// 			UserCode:    "1234",
		// 			FullName:    "Abebe Kebede",
		// 			PhoneNumber: "251911111111",
		// 		},
		// 		Department: "IT",
		// 		Status:     "PENDING",
		// 		ActionType: "CREATE",
		// 		ActionData: entity.Advert{
		// 			Title:       "Advert",
		// 			Description: "Advert Description",
		// 			BannerImage: "google.png",
		// 			AdvertFor:   entity.CB,
		// 			Date: entity.AdvertDate{
		// 				StartedAt: started_at,
		// 				ExpiredAt: expired_at,
		// 			},
		// 		},
		// 		MakerActionTime: maker_action_time,
		// 	},
		// 	mock: func() {
		// 		mockRepo.EXPECT().
		// 			CreateOneAdvert(gomock.Any(), model.CPSAction{
		// 				ActionCode:       "1234",
		// 				MakerID:          "1234",
		// 				MakerName:        "Abebe Kebede",
		// 				MakerPhoneNumber: "251911111111",
		// 				Department:       "IT",
		// 				ActionStatus:     "PENDING",
		// 				ActionType:       "CREATE",
		// 				CurrentAction: entity.Advert{
		// 					Title:       "Advert",
		// 					Description: "Advert Description",
		// 					BannerImage: "google.png",
		// 					AdvertFor:   entity.CB,
		// 					Date: entity.AdvertDate{
		// 						StartedAt: started_at,
		// 						ExpiredAt: expired_at,
		// 					},
		// 				},
		// 				MakerActionTime: maker_action_time,
		// 			}).
		// 			Return(nil, errors.New("internal server error"))
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.CreateOneAdvert(context.Background(), tt.cpsAction)
			fmt.Println("===========================================", err)
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

// func TestADDomain_DeleteOneAdvert(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockRepository(ctrl)
// 	mockLogger := &MockLogger{}
// 	mockMinioClient := &MockMinioClient{}
// 	service := InitADDomian("ad", mockMinioClient, mockRepo, mockLogger)

// 	tests := []struct {
// 		name      string
// 		cpsAction entity.CPSAction
// 		mock      func()
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			cpsAction: entity.CPSAction{
// 				MakerUser: entity.User{
// 					UserCode:    "1234",
// 					FullName:    "Abebe Kebede",
// 					PhoneNumber: "+251911111111",
// 				},
// 				ActionCode: "12345",
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					DeleteOneAdvert(gomock.Any(), entity.CPSAction{
// 						MakerUser: entity.User{
// 							UserCode:    "1234",
// 							FullName:    "Abebe Kebede",
// 							PhoneNumber: "+251911111111",
// 						},
// 						ActionCode: "12345",
// 					}).
// 					Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "error from repository",
// 			cpsAction: entity.CPSAction{

// 				MakerUser: entity.User{
// 					UserCode:    "1234",
// 					FullName:    "Abebe Kebede",
// 					PhoneNumber: "+251911111111",
// 				},
// 				ActionCode: "12345",
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					DeleteOneAdvert(gomock.Any(), gomock.Any(), entity.CPSAction{
// 						MakerUser: entity.User{
// 							UserCode:    "1234",
// 							FullName:    "Abebe Kebede",
// 							PhoneNumber: "+251911111111",
// 						},
// 						ActionCode: "12345",
// 					}).
// 					Return(errors.New("internal server error"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			_, err := service.DeleteOneAdvert(context.Background(), gomock.Any(), tt.cpsAction)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestADDomain_GetAllAdvert(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockRepository(ctrl)
// 	mockLogger := &MockLogger{}
// 	service := InitADDomian(mockRepo, mockLogger)

// 	started_at := time.Now()
// 	expired_at := time.Now().Add(4 * time.Hour)

// 	tests := []struct {
// 		name         string
// 		filterParams *constant.Filter
// 		mock         func()
// 		want         *entity.AdvertResponse
// 		wantErr      bool
// 	}{
// 		{
// 			name: "success",
// 			filterParams: &constant.Filter{
// 				Page:    1,
// 				PerPage: 3,
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					GetAllAdvert(gomock.Any(), &constant.Filter{Page: 1, PerPage: 3}).
// 					Return(&entity.AdvertResponse{
// 						Page: 1,
// 						Advert: []*entity.Advert{
// 							{

// 								Title:       "example1",
// 								Description: "description",
// 								BannerImage: "google.png",
// 								AdvertFor:   entity.Both,
// 								Date: entity.AdvertDate{
// 									StartedAt: started_at,
// 									ExpiredAt: expired_at,
// 								},
// 							},
// 							{

// 								Title:       "example2",
// 								Description: "descritpion",
// 								BannerImage: "youtube.png",
// 								AdvertFor:   entity.CB,
// 								Date: entity.AdvertDate{
// 									StartedAt: started_at,
// 									ExpiredAt: expired_at,
// 								},
// 							},
// 						},
// 						Limit: 10,
// 						Total: 2,
// 					}, nil)
// 			},
// 			want: &entity.AdvertResponse{
// 				Page: 1,
// 				Advert: []*entity.Advert{
// 					{

// 						Title:       "example1",
// 						Description: "description",
// 						BannerImage: "google.png",
// 						AdvertFor:   entity.Both,
// 						Date: entity.AdvertDate{
// 							StartedAt: started_at,
// 							ExpiredAt: expired_at,
// 						},
// 					},
// 					{

// 						Title:       "example2",
// 						Description: "descritpion",
// 						BannerImage: "youtube.png",
// 						AdvertFor:   entity.CB,
// 						Date: entity.AdvertDate{
// 							StartedAt: started_at,
// 							ExpiredAt: expired_at,
// 						},
// 					},
// 				},
// 				Total: 2,
// 				Limit: 10,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "error from repository",
// 			filterParams: &constant.Filter{
// 				Page:    1,
// 				PerPage: 4,
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					GetAllAdvert(gomock.Any(), &constant.Filter{Page: 1, PerPage: 4}).
// 					Return(nil, errors.New("internal server error"))
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.GetAllAdvert(context.Background(), tt.filterParams)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				assert.Nil(t, got)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want.Page, got.Page)
// 				assert.Equal(t, tt.want.Limit, got.Limit)
// 				for k := range len(got.Advert) {
// 					assert.Equal(t, tt.want.Advert[k].AdvertFor, got.Advert[k].AdvertFor)
// 					assert.Equal(t, tt.want.Advert[k].Title, got.Advert[k].Title)
// 					assert.Equal(t, tt.want.Advert[k].Description, got.Advert[k].Description)
// 				}
// 			}
// 		})
// 	}
// }

// func TestADDomain_GetOneAdvert(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	started_at := time.Now()
// 	expired_at := time.Now().Add(4 * time.Hour)

// 	mockRepo := mocks.NewMockRepository(ctrl)
// 	mockLogger := &MockLogger{}
// 	service := InitADDomian(mockRepo, mockLogger)

// 	tests := []struct {
// 		name    string
// 		id      string
// 		mock    func()
// 		want    *entity.Advert
// 		wantErr bool
// 	}{
// 		{
// 			name: "success",
// 			id:   "123",
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					GetOneAdvert(gomock.Any(), "123").
// 					Return(&entity.Advert{

// 						Title:       "example2",
// 						Description: "descritpion",
// 						BannerImage: "youtube.png",
// 						AdvertFor:   entity.CB,
// 						Date: entity.AdvertDate{
// 							StartedAt: started_at,
// 							ExpiredAt: expired_at,
// 						},
// 					}, nil)
// 			},
// 			want: &entity.Advert{

// 				Title:       "example2",
// 				Description: "descritpion",
// 				BannerImage: "youtube.png",
// 				AdvertFor:   entity.CB,
// 				Date: entity.AdvertDate{
// 					StartedAt: started_at,
// 					ExpiredAt: expired_at,
// 				}},
// 			wantErr: false,
// 		},
// 		{
// 			name: "error from repository",
// 			id:   "123",
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					GetOneAdvert(gomock.Any(), "123").
// 					Return(nil, errors.New("not found"))
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.GetOneAdvert(context.Background(), tt.id)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				assert.Nil(t, got)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want, got)
// 			}
// 		})
// 	}
// }

// func TestADDomain_UpdateOneAdvert(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockRepository(ctrl)
// 	mockLogger := &MockLogger{}
// 	service := InitADDomian(mockRepo, mockLogger)

// 	tests := []struct {
// 		name      string
// 		cpsAction entity.CPSAction
// 		mock      func()
// 		want      *entity.CPSAction
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			cpsAction: entity.CPSAction{
// 				ActionCode: "1234",
// 				MakerUser: entity.User{
// 					UserCode:    "345",
// 					FullName:    "Abebe Alemu",
// 					PhoneNumber: "+251922222321",
// 				},
// 				Department: "IT",
// 				Status:     entity.ActionPending,
// 				ActionData: entity.Advert{

// 					Title: "New Title",
// 				},
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					UpdateOneAdvert(gomock.Any(), entity.CPSAction{
// 						ActionCode: "1234",
// 						MakerUser: entity.User{
// 							UserCode:    "345",
// 							FullName:    "Abebe Alemu",
// 							PhoneNumber: "+251922222321",
// 						},
// 						ActionData: entity.Advert{

// 							Title: "New Title",
// 						},
// 						Department: "IT",
// 						Status:     entity.ActionPending,
// 					}).
// 					Return(&entity.CPSAction{
// 						ActionCode: "1234",
// 						MakerUser: entity.User{
// 							UserCode:    "345",
// 							FullName:    "Abebe Alemu",
// 							PhoneNumber: "+251922222321",
// 						},
// 						Department: "IT",
// 						Status:     entity.ActionPending,
// 						ActionData: entity.Advert{

// 							Title: "New Title",
// 						}}, nil)
// 			},
// 			want: &entity.CPSAction{
// 				ActionCode: "1234",
// 				MakerUser: entity.User{
// 					UserCode:    "345",
// 					FullName:    "Abebe Alemu",
// 					PhoneNumber: "+251922222321",
// 				},
// 				Department: "IT",
// 				Status:     entity.ActionPending,
// 				ActionData: entity.Advert{
// 					Title: "New Title",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "error from repository",
// 			cpsAction: entity.CPSAction{
// 				ID: "123",
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					UpdateOneAdvert(gomock.Any(), entity.CPSAction{ID: "123"}).
// 					Return(nil, errors.New("not found"))
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.UpdateOneAdvert(context.Background(), tt.cpsAction)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				assert.Nil(t, got)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want.ActionCode, got.ActionCode)
// 				assert.Equal(t, tt.want.ActionData.Title, got.ActionData.Title)
// 			}
// 		})
// 	}
// }

// func TestADDomain_Authorize(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockRepository(ctrl)
// 	mockLogger := &MockLogger{}
// 	service := InitADDomian(mockRepo, mockLogger)

// 	tests := []struct {
// 		name      string
// 		cpsAction entity.CPSAction
// 		mock      func()
// 		want      *entity.CPSAction
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			cpsAction: entity.CPSAction{
// 				ActionCode: "123",
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					Authorize(gomock.Any(), entity.CPSAction{
// 						ActionCode: "123",
// 					}).
// 					Return(&entity.CPSAction{
// 						ActionCode: "123",
// 						Status:     entity.ActionApproved,
// 					}, nil)
// 			},
// 			want: &entity.CPSAction{
// 				ActionCode: "123",
// 				Status:     entity.ActionApproved,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "error from repository",
// 			cpsAction: entity.CPSAction{
// 				ActionCode: "1234567",
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					Authorize(gomock.Any(), entity.CPSAction{
// 						ActionCode: "1234567",
// 					}).
// 					Return(nil, errors.New("not found"))
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.Authorize(context.Background(), tt.cpsAction)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				assert.Nil(t, got)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want, got)
// 			}
// 		})
// 	}
// }

// func TestADDomain_Reject(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockRepository(ctrl)
// 	mockLogger := &MockLogger{}
// 	service := InitADDomian(mockRepo, mockLogger)

// 	tests := []struct {
// 		name      string
// 		cpsAction entity.CPSAction
// 		mock      func()
// 		want      *entity.CPSAction
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			cpsAction: entity.CPSAction{
// 				ActionCode: "123",
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					Reject(gomock.Any(), entity.CPSAction{ActionCode: "123"}).
// 					Return(&entity.CPSAction{
// 						ActionCode: "123",
// 						Status:     entity.ActionRejected,
// 					}, nil)
// 			},
// 			want:    &entity.CPSAction{ActionCode: "123", Status: entity.ActionRejected},
// 			wantErr: false,
// 		},
// 		{
// 			name: "error from repository",
// 			cpsAction: entity.CPSAction{
// 				ActionCode: "123",
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().
// 					Reject(gomock.Any(), entity.CPSAction{
// 						ActionCode: "123",
// 					}).
// 					Return(nil, errors.New("not found"))
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.Reject(context.Background(), tt.cpsAction)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				assert.Nil(t, got)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want, got)
// 			}
// 		})
// 	}
// }
