package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/mocks"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/stretchr/testify/assert"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"
)

// MockLogger is a mock implementation of the Logger interface
type MockLogger struct{}

func (m *MockLogger) Infof(msg string, args ...interface{})  {}
func (m *MockLogger) Errorf(msg string, args ...interface{}) {}
func (m *MockLogger) Debugf(msg string, args ...interface{}) {}
func (m *MockLogger) Fatalf(msg string, args ...interface{}) {}
func (m *MockLogger) Warnf(msg string, args ...interface{})  {}
func (m *MockLogger) Sync() error {
	return nil
}

func TestGetCustomersDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testID := "665c17a7e629d4b3b437d2f9"
	objectID, _ := bson.ObjectIDFromHex(testID)

	tests := []struct {
		name           string
		filterParams   *utils.Filter
		mockResponse   *entity.CustomerRespose
		mockError      error
		expectedResult *entity.CustomerRespose
		expectedError  error
	}{
		{
			name: "Success",
			filterParams: &utils.Filter{
				Page:    1,
				PerPage: 5,
			},
			mockResponse: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Active",
					},
					{
						ID:                objectID,
						UserCode:          "USR456",
						FullName:          "Sara Johnson",
						MotherName:        "Mary Johnson",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1995, time.July, 24, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-3, 0, 0),
						PhoneNumber:       "+251911000002",
						Gender:            "Female",
						MaritalStatus:     "Married",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA654321",
							FaydaAccessToken: "token456",
							EmploymentStatus: "Self-Employed",
							EmployerName:     "Sara Beauty Salon",
							IssuedBy:         "Gov",
							MonthlyIncome:    8000,
						},
						Email:            "sara.johnson@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Active",
					},
				},
				Limit: 10,
				Total: 2,
			},
			mockError: nil,
			expectedResult: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Active",
					},
					{
						ID:                objectID,
						UserCode:          "USR456",
						FullName:          "Sara Johnson",
						MotherName:        "Mary Johnson",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1995, time.July, 24, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-3, 0, 0),
						PhoneNumber:       "+251911000002",
						Gender:            "Female",
						MaritalStatus:     "Married",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA654321",
							FaydaAccessToken: "token456",
							EmploymentStatus: "Self-Employed",
							EmployerName:     "Sara Beauty Salon",
							IssuedBy:         "Gov",
							MonthlyIncome:    8000,
						},
						Email:            "sara.johnson@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Active",
					},
				},
				Limit: 10,
				Total: 2,
			},
			expectedError: nil,
		},
		{
			name: "No Customers Found",
			filterParams: &utils.Filter{
				Page:    1,
				PerPage: 10,
			},
			mockResponse: &entity.CustomerRespose{
				Page:      1,
				Limit:     10,
				Customers: nil,
				Total:     0,
			},
			mockError: nil,
			expectedResult: &entity.CustomerRespose{
				Page:      1,
				Limit:     10,
				Customers: nil,
				Total:     0,
			},
			expectedError: nil,
		},
		{
			name: "Filter Pending Account Status",
			filterParams: &utils.Filter{
				Page:    1,
				PerPage: 5,
				Filters: "Pending",
			},
			mockResponse: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Pending",
					},
				},
				Limit: 10,
				Total: 1,
			},
			mockError: nil,
			expectedResult: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Pending",
					},
				},
				Limit: 10,
				Total: 1,
			},
			expectedError: nil,
		},
		{
			name: "Search by user code",
			filterParams: &utils.Filter{
				Page:    1,
				PerPage: 5,
				Search:  "USR123",
			},
			mockResponse: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Pending",
					},
				},
				Limit: 10,
				Total: 1,
			},
			mockError: nil,
			expectedResult: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Pending",
					},
				},
				Limit: 10,
				Total: 1,
			},
			expectedError: nil,
		},
		{
			name: "Search by name",
			filterParams: &utils.Filter{
				Page:    1,
				PerPage: 5,
				Search:  "John Smith Doe",
			},
			mockResponse: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Pending",
					},
				},
				Limit: 10,
				Total: 1,
			},
			mockError: nil,
			expectedResult: &entity.CustomerRespose{
				Page: 1,
				Customers: []*member.User{
					{
						ID:                objectID,
						UserCode:          "USR123",
						FullName:          "John Smith Doe",
						MotherName:        "Jane Doe",
						Nationality:       "Ethiopian",
						BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
						ResidentialStatus: "Resident",
						IssuedDate:        time.Now().AddDate(-5, 0, 0),
						PhoneNumber:       "+251911000001",
						Gender:            "Male",
						MaritalStatus:     "Single",
						Fayda: struct {
							FaydaID          string `json:"id_number" bson:"id_number"`
							FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
							EmploymentStatus string `json:"employment_status" bson:"employement_status"`
							EmployerName     string `json:"employer_name" bson:"employer_name"`
							IssuedBy         string `json:"issued_by" bson:"issued_by"`
							MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
						}{
							FaydaID:          "FA123456",
							FaydaAccessToken: "token123",
							EmploymentStatus: "Employed",
							EmployerName:     "Tech Co",
							IssuedBy:         "Gov",
							MonthlyIncome:    12000,
						},
						Email:            "john.doe@example.com",
						IsAccountBlocked: false,
						CreatedAt:        time.Now(),
						AccountStatus:    "Pending",
					},
				},
				Limit: 10,
				Total: 1,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockCustomerRepository(ctrl)
			mockLogger := &MockLogger{}
			service := IntiCustomerDomain(mockRepo, mockLogger)

			mockRepo.EXPECT().
				GetCustomersDetail(gomock.Any(), tt.filterParams).
				Return(tt.mockResponse, tt.mockError)

			result, err := service.GetCustomersDetail(context.Background(), tt.filterParams)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Page, result.Page)
				assert.Equal(t, tt.expectedResult.Limit, result.Limit)
				assert.Equal(t, tt.expectedResult.Total, result.Total)
				for i := range tt.expectedResult.Customers {
					assert.Equal(t, tt.expectedResult.Customers[i].ID, result.Customers[i].ID)
					assert.Equal(t, tt.expectedResult.Customers[i].FullName, result.Customers[i].FullName)
					assert.Equal(t, tt.expectedResult.Customers[i].AccountStatus, result.Customers[i].AccountStatus)
					assert.Equal(t, tt.expectedResult.Customers[i].UserCode, tt.expectedResult.Customers[i].UserCode)
				}
			}
		})
	}
}

func TestGetCustomerByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testID := "665c17a7e629d4b3b437d2f9"
	objectID, _ := bson.ObjectIDFromHex(testID)

	tests := []struct {
		name           string
		customerID     string
		mockResponse   *member.User
		mockError      error
		expectedResult *member.User
		expectedError  error
	}{
		{
			name:       "Success",
			customerID: objectID.String(),
			mockResponse: &member.User{
				ID:                objectID,
				UserCode:          "USR123",
				FullName:          "John Smith Doe",
				MotherName:        "Jane Doe",
				Nationality:       "Ethiopian",
				BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
				ResidentialStatus: "Resident",
				IssuedDate:        time.Now().AddDate(-5, 0, 0),
				PhoneNumber:       "+251911000001",
				Gender:            "Male",
				MaritalStatus:     "Single",
				Fayda: struct {
					FaydaID          string `json:"id_number" bson:"id_number"`
					FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
					EmploymentStatus string `json:"employment_status" bson:"employement_status"`
					EmployerName     string `json:"employer_name" bson:"employer_name"`
					IssuedBy         string `json:"issued_by" bson:"issued_by"`
					MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
				}{
					FaydaID:          "FA123456",
					FaydaAccessToken: "token123",
					EmploymentStatus: "Employed",
					EmployerName:     "Tech Co",
					IssuedBy:         "Gov",
					MonthlyIncome:    12000,
				},
				Email:            "john.doe@example.com",
				IsAccountBlocked: false,
				CreatedAt:        time.Now(),
				AccountStatus:    "Pending",
			},
			mockError: nil,
			expectedResult: &member.User{
				ID:                objectID,
				UserCode:          "USR123",
				FullName:          "John Smith Doe",
				MotherName:        "Jane Doe",
				Nationality:       "Ethiopian",
				BirthDate:         time.Date(1990, time.March, 10, 0, 0, 0, 0, time.UTC),
				ResidentialStatus: "Resident",
				IssuedDate:        time.Now().AddDate(-5, 0, 0),
				PhoneNumber:       "+251911000001",
				Gender:            "Male",
				MaritalStatus:     "Single",
				Fayda: struct {
					FaydaID          string `json:"id_number" bson:"id_number"`
					FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
					EmploymentStatus string `json:"employment_status" bson:"employement_status"`
					EmployerName     string `json:"employer_name" bson:"employer_name"`
					IssuedBy         string `json:"issued_by" bson:"issued_by"`
					MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
				}{
					FaydaID:          "FA123456",
					FaydaAccessToken: "token123",
					EmploymentStatus: "Employed",
					EmployerName:     "Tech Co",
					IssuedBy:         "Gov",
					MonthlyIncome:    12000,
				},
				Email:            "john.doe@example.com",
				IsAccountBlocked: false,
				CreatedAt:        time.Now(),
				AccountStatus:    "Pending",
			},
			expectedError: nil,
		},
		{
			name:           "No Customer Found",
			customerID:     objectID.String(),
			mockResponse:   nil,
			mockError:      errors.New("sustomer not found"),
			expectedResult: nil,
			expectedError:  errors.New("sustomer not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockCustomerRepository(ctrl)
			mockLogger := &MockLogger{}
			service := IntiCustomerDomain(mockRepo, mockLogger)

			mockRepo.EXPECT().
				GetCustomerByID(gomock.Any(), tt.customerID).
				Return(tt.mockResponse, tt.mockError)

			result, err := service.GetCustomerByID(context.Background(), tt.customerID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.FullName, result.FullName)
				assert.Equal(t, tt.expectedResult.UserCode, result.UserCode)
				assert.Equal(t, tt.expectedResult.AccountStatus, result.AccountStatus)
			}
		})
	}
}
