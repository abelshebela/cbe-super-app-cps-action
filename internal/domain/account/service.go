package account

import (
	accountRepoPort "cbe-super-app-member-users/internal/port/outbound/account"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type AccountService struct {
	repository accountRepoPort.AccountRepositoryPort
	apiClient  accountRepoPort.AccountAPIPort
	logger     interface {
		Errorf(string, ...interface{})
		Warnf(string, ...interface{})
		Infof(string, ...interface{})
	}
	cfg *config.VaultConfig
}

func NewAccountService(repository accountRepoPort.AccountRepositoryPort, apiClient accountRepoPort.AccountAPIPort, logger interface {
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
}, cfg *config.VaultConfig) *AccountService {
	return &AccountService{
		repository: repository,
		apiClient:  apiClient,
		logger:     logger,
		cfg:        cfg,
	}
}

var ErrNotFound = errors.New("not found")

func (s *AccountService) CreateAccount(ctx context.Context, userID string) (*AccountCreationResult, error) {
	user, err := s.repository.FindAccountUserByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to fetch user with ID %s: %s", userID, err.Error())
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	if user.KYCLevel == 0 {
		s.logger.Warnf("User %s has KYC level 0", userID)
		return nil, fmt.Errorf("USER_KYC_LEVEL_ZERO")
	}

	exists, err := s.apiClient.LookupAccountByPhone(ctx, user.PhoneNumber, "" /*s.cfg.PhoneLookupUrl*/)
	if err != nil {
		s.logger.Errorf("Phone lookup failed for user %s: %s", userID, err.Error())
		return nil, fmt.Errorf("API_REQUEST_FAILED")
	}
	if exists {
		s.logger.Warnf("Phone number %s already exists for user %s", user.PhoneNumber, userID)
		return nil, fmt.Errorf("USER_PHONE_EXISTS")
	}

	sif, accountNumber, err := s.GenerateSIFAndAccountNumber()
	if err != nil {
		s.logger.Errorf("Failed to generate SIF and account number: %s", err.Error())
		return nil, fmt.Errorf("GENERAL_SIF_GENERATION_FAILED")
	}

	if err := s.repository.UpdateUserCustomerNumber(ctx, userID, sif); err != nil {
		s.logger.Errorf("Failed to update user %s: %s", userID, err.Error())
		return nil, fmt.Errorf("GENERAL_DB_UPDATE_FAILED")
	}

	var accountType string
	if user.RegistrationType != "" {
		accountType = user.RegistrationType
	} else {
		accountType = "Savings"
	}

	linkedAccount := &LinkedAccount{
		UserID:            userID,
		CustomerNumber:    sif,
		AccountNumber:     accountNumber,
		AccountHolderName: user.FullName,
		AccountType:       accountType,
		BranchCode:        user.BranchCode,
		RegistrationType:  user.RegistrationType,
		AndOrStatus:       user.AndOrStatus,
		CurrencyCode:      "USD",
		IsMain:            true,
	}

	linkedAccountRepo := &accountRepoPort.LinkedAccount{
		UserID:            linkedAccount.UserID,
		CustomerNumber:    linkedAccount.CustomerNumber,
		AccountNumber:     linkedAccount.AccountNumber,
		AccountHolderName: linkedAccount.AccountHolderName,
		AccountType:       linkedAccount.AccountType,
		BranchCode:        linkedAccount.BranchCode,
		RegistrationType:  linkedAccount.RegistrationType,
		AndOrStatus:       linkedAccount.AndOrStatus,
		CurrencyCode:      linkedAccount.CurrencyCode,
		IsMain:            linkedAccount.IsMain,
	}

	if err := s.repository.CreateLinkedAccount(ctx, linkedAccountRepo); err != nil {
		s.logger.Errorf("Failed to create linked account for user %s: %s", userID, err.Error())
		return nil, fmt.Errorf("GENERAL_DB_INSERT_FAILED")
	}

	s.logger.Infof("Account created successfully for user ID %s: SIF=%s, AccountNumber=%s", userID, sif, accountNumber)
	return &AccountCreationResult{
		CustomerNumber: sif,
		AccountNumber:  accountNumber,
	}, nil
}

func (s *AccountService) GenerateSIFAndAccountNumber() (string, string, error) {
	sif := "SIF-" + uuid.New().String()[:8]
	accountNumber := "ACC-" + uuid.New().String()[:10]
	return sif, accountNumber, nil
}
