package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"

	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func ApplyOpenUpdate(openTier *model.AuthTier, pinTier *model.AuthTier) error {

	if openTier.Method != shared_constant.OPEN || pinTier.Method != shared_constant.PIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	if openTier.MinAmount >= openTier.MaxAmount {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	pinTier.MinAmount = openTier.MaxAmount
	return nil
}

// ApplyPinUpdate validates PIN tier values against OPEN and OTP_PIN constraints
// but does NOT force PIN values - it validates user input is within valid ranges
func ApplyPinUpdate(pinTier *model.AuthTier, openTier *model.AuthTier, otpPinTier *model.AuthTier) error {
	if pinTier.Method != shared_constant.PIN || openTier.Method != shared_constant.OPEN || otpPinTier.Method != shared_constant.OTPANDPIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	// Validate referenced tiers are valid
	if openTier.MinAmount >= openTier.MaxAmount {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	// Validate PIN tier values are within acceptable ranges
	if pinTier.MinAmount >= pinTier.MaxAmount {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	openTier.MaxAmount = pinTier.MinAmount
	otpPinTier.MinAmount = pinTier.MaxAmount

	return nil
}

func ApplyOtpPinUpdate(otpPinTier *model.AuthTier, pinTier *model.AuthTier) error {
	if otpPinTier.Method != shared_constant.OTPANDPIN || pinTier.Method != shared_constant.PIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	pinTier.MaxAmount = otpPinTier.MinAmount
	return nil
}
