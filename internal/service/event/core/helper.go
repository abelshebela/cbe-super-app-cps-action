package core

import (
	"cbe-super-app-cps-action/internal/constants"
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	"cbe-super-app-cps-action/internal/constants/lib"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"context"
	cRand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/mongo"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func NonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func NonZeroTime(t, fallback time.Time) time.Time {
	if !t.IsZero() {
		return t
	}
	return fallback
}

func NonZeroUint64(n, fallback uint64) uint64 {
	if n != 0 {
		return n
	}
	return fallback
}

func NonEmptyTickets(tickets, fallback []types.Ticket) []types.Ticket {
	if len(tickets) > 0 {
		return tickets
	}
	return fallback
}

func GeneratePrefixedName(prefix, value string, logger shared_utils.Logger) (string, error) {
	logger.Infof("Generating prefixed name", "prefix", prefix, "value", value)

	if prefix == "" || value == "" {
		logger.Errorf("Invalid input for GeneratePrefixedName", "prefix", prefix, "value", value)
		return "", fmt.Errorf("prefix and value must not be empty")
	}

	digits := "0123456789"
	max := big.NewInt(int64(len(digits)))
	code := make([]byte, 7)

	for i := range code {
		n, err := cRand.Int(cRand.Reader, max)
		if err != nil {
			logger.Errorf("Failed to generate random digit", "error", err)
			return "", fmt.Errorf("failed to generate random digit: %v", err)
		}
		code[i] = digits[n.Int64()]
	}

	value = strings.ReplaceAll(value, " ", "")
	prefix = strings.ReplaceAll(prefix, " ", "")
	result := strings.Join([]string{prefix, value, string(code)}, "-")
	logger.Infof("Successfully generated prefixed name", "result", result)
	return result, nil
}
func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func GenerateEvent(event model.Event) *model.Event {
	return &model.Event{
		ID:            event.ID,
		EventCode:     event.EventCode,
		EventName:     event.EventName,
		EventCity:     event.EventCity,
		EventVenue:    event.EventVenue,
		Status:        constants.EventUpcomming,
		AccountNumber: event.AccountNumber,
		MerchantInformation: types.MerchantInformation{
			MerchantID:          event.MerchantInformation.MerchantID,
			MercahntName:        event.MerchantInformation.MercahntName,
			MerchantPhoneNumber: event.MerchantInformation.MerchantPhoneNumber,
			MerchantEmail:       event.MerchantInformation.MerchantEmail,
		},
		EventInformation: types.EventInformation{
			StartDate:   event.EventInformation.StartDate,
			DueDate:     event.EventInformation.DueDate,
			Description: event.EventInformation.Description,
			Cover:       event.EventInformation.Cover,
		},
		TicketInformation: types.TicketInformation{
			TotalNumberOfTicket: event.TicketInformation.TotalNumberOfTicket,
		},
		Ticket:         event.Ticket,
		CreatedAt:      event.CreatedAt,
		LastModifiedAt: event.LastModifiedAt,
	}
}

func SetMerchantDetails(ctx context.Context, merchantService service.MiniAppMerchantService, event *eventdto.EventRequest) error {
	if event.MerchantID == "" {
		return nil
	}

	merchant, err := merchantService.FindByID(ctx, event.MerchantID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("Failed to get merchant details", "merchantID", event.MerchantID, "error", err)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Println("  merchant details merchant", merchant)
	event.MercahntName = merchant.MerchantName
	event.MerchantEmail = merchant.Email
	event.MerchantPhoneNumber = merchant.PhoneNumber
	event.AccountNumber = merchant.BankAccountNumber

	log.Println("  merchant details MercahntName: ", event.MercahntName)
	log.Println("  merchant details MerchantPhoneNumber: ", event.MerchantPhoneNumber)
	log.Println("  merchant details MerchantEmail: ", event.MerchantEmail)

	return nil
}
func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(constants.ActionDelete))

	log.Println("Creating CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID, "actionType", actionType)
	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func EventMapperForUpdate(prevEvent *model.Event, update eventdto.EventRequest, coverURL string) model.Event {
	return model.Event{
		ID:         prevEvent.ID,
		EventCode:  prevEvent.EventCode,
		EventName:  NonEmptyString(update.EventName, prevEvent.EventName),
		EventCity:  NonEmptyString(update.EventCity, prevEvent.EventCity),
		EventVenue: NonEmptyString(update.EventVenue, prevEvent.EventVenue),
		Status:     constants.EventUpcomming,
		MerchantInformation: types.MerchantInformation{
			MerchantID:          NonEmptyString(update.MerchantID, prevEvent.MerchantInformation.MerchantID),
			MercahntName:        NonEmptyString(update.MercahntName, prevEvent.MerchantInformation.MercahntName),
			MerchantPhoneNumber: NonEmptyString(update.MerchantPhoneNumber, prevEvent.MerchantInformation.MerchantPhoneNumber),
			MerchantEmail:       NonEmptyString(update.MerchantEmail, prevEvent.MerchantInformation.MerchantEmail),
		},
		EventInformation: types.EventInformation{
			StartDate:   NonZeroTime(update.StartDate, prevEvent.EventInformation.StartDate),
			DueDate:     NonZeroTime(update.DueDate, prevEvent.EventInformation.DueDate),
			Description: NonEmptyString(update.EventDescription, prevEvent.EventInformation.Description),
			Cover:       NonEmptyString(coverURL, prevEvent.EventInformation.Cover),
		},
		TicketInformation: types.TicketInformation{
			TotalNumberOfTicket: NonZeroUint64(uint64(update.TotalTicketCount), prevEvent.TicketInformation.TotalNumberOfTicket),
		},
		Ticket:         NonEmptyTickets(update.Tickets, prevEvent.Ticket),
		CreatedAt:      prevEvent.CreatedAt,
		LastModifiedAt: time.Now(),
		AccountNumber:  NonEmptyString(update.AccountNumber, prevEvent.AccountNumber),
	}
}

func CreateEventMapper(event eventdto.EventRequest, code string, coverURL string) *model.Event {
	return &model.Event{
		EventCode:     code,
		EventName:     event.EventName,
		EventCity:     event.EventCity,
		EventVenue:    event.EventVenue,
		Status:        constants.EventUpcomming,
		AccountNumber: event.AccountNumber,
		MerchantInformation: types.MerchantInformation{
			MerchantID:          event.MerchantID,
			MercahntName:        event.MercahntName,
			MerchantPhoneNumber: event.MerchantPhoneNumber,
			MerchantEmail:       event.MerchantEmail,
		},
		EventInformation: types.EventInformation{
			StartDate:   event.StartDate,
			DueDate:     event.DueDate,
			Description: event.EventDescription,
			Cover:       coverURL,
		},
		TicketInformation: types.TicketInformation{
			TotalNumberOfTicket:          uint64(event.TotalTicketCount),
			TotalNumberOfAvailableTicket: uint64(event.TotalTicketCount),
		},
		Ticket:         event.Tickets,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
		Enabled:        true,
	}
}
