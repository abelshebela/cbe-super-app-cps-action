package core

import (
	"cbe-super-app-cps-action/internal/constants"
	amountauthdto "cbe-super-app-cps-action/internal/constants/dto/amount_based_auth"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ValidateTierCascade validates that tiers form a contiguous cascade:
// TierN.max == TierN+1.min for all adjacent tiers.
// Also validates: tier1.min is lowest, no equal min/max within a tier,
// and tierN.max != previous tier's min or max.
func ValidateTierCascade(tiers []amountauthdto.TierInput) error {
	if len(tiers) == 0 {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	for i := 0; i < len(tiers)-1; i++ {
		curr := tiers[i]
		next := tiers[i+1]

		// min < max for non-last tiers
		if curr.MinAmount >= curr.MaxAmount {
			return errors.New(localization.ErrorInvalidAmounts.Code)
		}

		// cascading continuity: curr.max == next.min
		if curr.MaxAmount != next.MinAmount {
			return errors.New(localization.ErrorInvalidTierCascade.Code)
		}

		// tierN.max != previous tier min and max
		if i > 0 {
			prev := tiers[i-1]
			if curr.MaxAmount == prev.MinAmount || curr.MaxAmount == prev.MaxAmount {
				return errors.New(localization.ErrorInvalidAmounts.Code)
			}
		}
	}

	return nil
}

// BuildTiersFromRequest creates AuthTier models from the AddCurrencyRequest input
func BuildTiersFromRequest(currency constants.CurrencyType, tiers []amountauthdto.TierInput) []local_model.AuthTier {
	now := time.Now()
	result := make([]local_model.AuthTier, len(tiers))
	for i, t := range tiers {
		result[i] = local_model.AuthTier{
			ID:           bson.NewObjectID(),
			Currency:     currency,
			MinAmount:    t.MinAmount,
			MaxAmount:    t.MaxAmount,
			Method:       t.Method,
			Enabled:      true,
			IsDeleted:    false,
			CreatedAt:    now,
			LastModified: now,
		}
	}
	return result
}

// ApplyOpenUpdate cascades OPEN tier max to PIN tier min
func ApplyOpenUpdate(openTier *local_model.AuthTier, pinTier *local_model.AuthTier) error {
	if openTier.Method != constants.OPEN || pinTier.Method != constants.PIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	if openTier.MinAmount >= openTier.MaxAmount {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	pinTier.MinAmount = openTier.MaxAmount
	return nil
}

// ApplyPinUpdate validates PIN tier values against OPEN and OTP_PIN constraints
func ApplyPinUpdate(pinTier *local_model.AuthTier, openTier *local_model.AuthTier, otpPinTier *local_model.AuthTier) error {
	if pinTier.Method != constants.PIN || openTier.Method != constants.OPEN || otpPinTier.Method != constants.OTPANDPIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	if openTier.MinAmount >= openTier.MaxAmount {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	if pinTier.MinAmount >= pinTier.MaxAmount {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}

	openTier.MaxAmount = pinTier.MinAmount
	otpPinTier.MinAmount = pinTier.MaxAmount

	return nil
}

// ApplyOtpPinUpdate cascades OTP_PIN tier min to PIN tier max
func ApplyOtpPinUpdate(otpPinTier *local_model.AuthTier, pinTier *local_model.AuthTier) error {
	if otpPinTier.Method != constants.OTPANDPIN || pinTier.Method != constants.PIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}

	pinTier.MaxAmount = otpPinTier.MinAmount
	return nil
}
