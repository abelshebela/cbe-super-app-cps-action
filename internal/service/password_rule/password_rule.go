package passwordrule

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/password_rule/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"unicode"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type passwordService struct {
	repo       storage.PasswordRuleRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewPasswordRuleService(repo storage.PasswordRuleRepository, cpsService service.CPSActionService, logger utils.Logger) service.PasswordRuleService {
	return &passwordService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (p *passwordService) GetAllPasswordRules(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.PasswordRule], error) {
	result, err := p.repo.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		p.logger.Errorf("[GetAllPasswordRules] failed to fetch password rules: %v", err)
		return nil, err
	}
	p.logger.Infof("[GetAllPasswordRules] retrieved %d password rules", len(result.Data))
	return result, nil
}

func (p *passwordService) RequestPasswordRuleUpdate(ctx context.Context, id string, body model.PasswordRule) error {
	existingRule, err := p.repo.FindCurrentRule(ctx)
	if err != nil {
		p.logger.Errorf("[RequestPasswordRuleUpdate] failed to fetch current password rule: %v", err)
		return errors.New(localization.ErrorNoPasswordRule.Code)
	}

	err = core.HandleCPSAction(ctx, p.cpsService, existingRule.ID.Hex(), constants.RequestUpdatePasswordRule, body, existingRule, constants.ActionUpdate)
	if err != nil {
		p.logger.Errorf("[RequestPasswordRuleUpdate] failed to create CPS action: %v", err)
		return err
	}

	p.logger.Infof("[RequestPasswordRuleUpdate] password rule update request created successfully")
	return nil
}

func (p *passwordService) CheckPasswordRule(ctx context.Context, password string) (bool, string) {
	rule, err := p.repo.FindCurrentRule(ctx)
	if err != nil || rule == nil {
		p.logger.Errorf("[CheckPasswordRule] failed to retrieve password rule: %v", err)
		return false, "could not retrieve password rule"
	}

	if len(password) < rule.MinLength {
		return false, fmt.Sprintf("password must be at least %d characters", rule.MinLength)
	}

	if rule.MaxLength > 0 && len(password) > rule.MaxLength {
		return false, fmt.Sprintf("password must be at most %d characters", rule.MaxLength)
	}

	if rule.Numbers {
		hasNumber := false
		for _, c := range password {
			if c >= '0' && c <= '9' {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return false, "password must contain at least one number"
		}
	}

	if rule.CapitalLetters {
		hasUpper := false
		for _, c := range password {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return false, "password must contain at least one uppercase letter"
		}
	}

	if rule.SmallLetters {
		hasLower := false
		for _, c := range password {
			if unicode.IsLower(c) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return false, "password must contain at least one lowercase letter"
		}
	}

	if rule.Characters {
		hasSpecial := false
		for _, c := range password {
			if (c >= 33 && c <= 47) || (c >= 58 && c <= 64) || (c >= 91 && c <= 96) || (c >= 123 && c <= 126) {
				hasSpecial = true
				break
			}
		}
		if !hasSpecial {
			return false, "password must contain at least one special character"
		}
	}
	return true, "Password is valid"
}

func (p *passwordService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	p.logger.Infof("[Authorize] authorizing password rule action: %s", cpsAction.RequestAction)

	passwordRule, err := local_util.JsonUnmarshal[model.PasswordRule](cpsAction.CurrentAction)
	if err != nil {
		p.logger.Errorf("[Authorize] failed to unmarshal password rule from action: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	err = p.repo.Update(ctx, cpsAction.UniqueId, passwordRule)
	if err != nil {
		p.logger.Errorf("[Authorize] failed to update password rule: %v", err)
		return nil, err
	}

	cpsAction.CurrentAction = passwordRule
	p.logger.Infof("[Authorize] password rule authorized successfully for id: %s", cpsAction.UniqueId)
	return cpsAction, nil
}
