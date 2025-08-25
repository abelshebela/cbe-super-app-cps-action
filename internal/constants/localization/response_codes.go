package localization

// ResponseCode represents a standardized response code with status and message
type ResponseCode struct {
	Code       string `json:"code"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Type       string `json:"type"` // "success", "error", "warning", "info"
}

var ResponseCodesList = []ResponseCode{
	// Success codes
	SuccessUserCreated,
	SuccessUserUpdated,
	SuccessUserDeleted,
	SuccessUserRetrieved,
	SuccessUserLogin,
	SuccessUserLogout,
	SuccessOTPSent,
	SuccessDonationCompanyLogoUploaded,
	SuccessDonationCompanyLogoUpdated,
	SuccessDonationImageUploaded,
	SuccessDonationImagesUpdated,
	SuccessNotificationConstructed,
	SuccessNotificationUpdated,
	SuccessFeedbackSavedToDatabase,
	SuccessAvatarCreated,

	// Ad related success response codes
	SuccessAdvertCreated,
	SuccessAdvertCreateRequestSent,
	SuccessAdvertUpdated,
	SuccessAdvertUpdateRequestSent,
	SuccessAdvertDeleted,
	SuccessAdvertDeleteRequestSent,
	SuccessAdvertFetched,
	SuccessAdvertsFetched,
	SuccessAdvertEnableRequestSent,
	SuccessAdvertDisableRequestSent,

	// Error codes
	ErrorUserNotFound,
	ErrorUserAlreadyExists,
	ErrorUserUnauthorized,
	ErrorUserForbidden,
	ErrorUserInvalidCredentials,
	ErrorUserAccountBlocked,
	ErrorUserSessionExpired,
	ErrorProductCodeDatabaseQueryFailed,
	ErrorProductCodeCountFailed,
	ErrorProductCodeDatabaseUpdateFailed,
	ErrorMiniAppMerchantCheckPendingFailed,
	ErrorMiniAppMerchantExistsCheckFailed,
	ErrorMiniAppMerchantFetchCategoriesFailed,
	ErrorMiniAppMerchantUnexpectedDBError,
	ErrorMiniAppMerchantMarshalFailed,
	ErrorMiniAppMerchantUnmarshalFailed,
	ErrorMiniAppMerchantCreateFromActionFailed,
	ErrorMiniAppMerchantUnmarshalActionFailed,
	ErrorMiniAppMerchantUpdateFailed,
	ErrorMiniAppMerchantFetchPermissionsFailed,
	ErrorNotificationMapFailed,
	ErrorFeedbackSavedToDatabase,
	ErrorHQApproved,
	ErrorCPSActionStatusInvalid,
	ErrorAdvertCreated,
	ErrorAdvertUpdate,
	ErrorAdvertUpdate,
	ErrorAdvertAlreadyEnabled,
	ErrorAdvertAlreadyDisabled,
	ErrorValidationRuleApproved,
	ErrorAccountNumberRequired,
	ErrorAccountNumberRequired,
	ErrorActionNotFound,
	ErrorUserAlreadyEnabled,
	ErrorUserAlreadyDisabled,
	ErrorUserCodeRequired,
  
	ErrorEventNameRequired,
	ErrorEventAlreadyExists,
	ErrorCoverImageRequired,
	ErrorMerchantNotFound,
	ErrorEventNotFound,
	ErrorEventAlreadyEnabled,
	ErrorEventAlreadyDisabled,
	ErrorUnhandledServer,
	ErrorInvalidDateFormat,
	ErrorInvalidFileUpload,
	ErrorInvalidNumberFormat,
	ErrorUpdateEventEmptyPayload,
	ErrorMerchantIDRequired,
	ErrorEventVenueRequired,
	ErrorStartDateRequired,
	ErrorDueDateRequired,
	ErrorTotalTicketCountRequired,
	ErrorInvalidTicketCount,
	ErrorEventCityRequired,
	ErrorEventDescriptionRequired,
	ErrorTicketsRequired,
  
	ErrorPendingCpsActionExists,
	ErrorUnexpectedError,
	ErrorFileNotFound,
	ErrorSessionRetrievalFailed,
	ErrorHealthCheck,
	ErrorInvalidToken,

	ErrorInvalidID,
	ErrorMissingFile,

	// Add more as needed...
}

// Success Response Codes
var (
	SuccessUserCreated = ResponseCode{
		Code:       "SUCCESS_USER_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgUserCreatedSuccessfully,
		Type:       "success",
	}

	SuccessAvatarCreated = ResponseCode{
		Code:       "SUCCESS_AVATAR_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgAvatarCreatedSuccessfully,
		Type:       "success",
	}

	SuccessUserUpdated = ResponseCode{
		Code:       "SUCCESS_USER_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgUserUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessUserDeleted = ResponseCode{
		Code:       "SUCCESS_USER_DELETED",
		StatusCode: StatusOK,
		Message:    MsgUserDeletedSuccessfully,
		Type:       "success",
	}

	SuccessUserRetrieved = ResponseCode{
		Code:       "SUCCESS_USER_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgUserRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessUserLogin = ResponseCode{
		Code:       "SUCCESS_USER_LOGIN",
		StatusCode: StatusOK,
		Message:    MsgUserLoginSuccessfully,
		Type:       "success",
	}

	SuccessUserLogout = ResponseCode{
		Code:       "SUCCESS_USER_LOGOUT",
		StatusCode: StatusOK,
		Message:    MsgUserLogoutSuccessfully,
		Type:       "success",
	}

	SuccessOTPSent = ResponseCode{
		Code:       "SUCCESS_OTP_SENT",
		StatusCode: StatusOK,
		Message:    MsgOTPSentSuccessfully,
		Type:       "success",
	}

	SuccessOTPVerified = ResponseCode{
		Code:       "SUCCESS_OTP_VERIFIED",
		StatusCode: StatusOK,
		Message:    MsgOTPVerifiedSuccessfully,
		Type:       "success",
	}

	SuccessOperationCompleted = ResponseCode{
		Code:       "SUCCESS_OPERATION_COMPLETED",
		StatusCode: StatusOK,
		Message:    MsgOperationCompleted,
		Type:       "success",
	}

	SuccessDataRetrieved = ResponseCode{
		Code:       "SUCCESS_DATA_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgDataRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessDataSaved = ResponseCode{
		Code:       "SUCCESS_DATA_SAVED",
		StatusCode: StatusCreated,
		Message:    MsgDataSavedSuccessfully,
		Type:       "success",
	}

	SuccessDataUpdated = ResponseCode{
		Code:       "SUCCESS_DATA_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDataUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessDataDeleted = ResponseCode{
		Code:       "SUCCESS_DATA_DELETED",
		StatusCode: StatusOK,
		Message:    MsgDataDeletedSuccessfully,
		Type:       "success",
	}

	SuccessHealthCheck = ResponseCode{
		Code:       "SUCCESS_HEALTH_CHECK",
		StatusCode: StatusOK,
		Message:    MsgHealthCheckPassed,
		Type:       "success",
	}

	SuccessUserPINChanged = ResponseCode{
		Code:       "SUCCESS_USER_PIN_CHANGED",
		StatusCode: StatusOK,
		Message:    MsgUserPINChanged,
		Type:       "success",
	}

	// Bank related success response codes
	SuccessBankCreatedRequestSent = ResponseCode{
		Code:       "SUCCESS_BANK_CREATED_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgBankCreatedRequestSent,
		Type:       "success",
	}

	SuccessBankUpdatedRequestSent = ResponseCode{
		Code:       "SUCCESS_BANK_UPDATED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgBankUpdatedRequestSent,
		Type:       "success",
	}

	SuccessBankDeletedRequestSent = ResponseCode{
		Code:       "SUCCESS_BANK_DELETED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgBankDeletedRequestSent,
		Type:       "success",
	}

	SuccessBanksRetrievedRequestSent = ResponseCode{
		Code:       "SUCCESS_BANKS_RETRIEVED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgBanksRetrievedRequestSent,
		Type:       "success",
	}

	SuccessBankRetrieved = ResponseCode{
		Code:       "SUCCESS_BANK_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgBankRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessBankRejectedRequestSent = ResponseCode{
		Code:       "SUCCESS_BANK_REJECTED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgBankRejectedRequestSent,
		Type:       "success",
	}

	SuccessBankDisableRequestSent = ResponseCode{
		Code:       "SUCCESS_BANK_DISABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgBankDisableRequestSent,
		Type:       "success",
	}

	SuccessBankEnableRequestSent = ResponseCode{
		Code:       "SUCCESS_BANK_ENABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgBankEnableRequestSent,
		Type:       "success",
	}

	SuccessBankLogoUpdatedRequestSent = ResponseCode{
		Code:       "SUCCESS_BANK_LOGO_UPDATED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgBankLogoUpdatedRequestSent,
		Type:       "success",
	}

	// Donation related success response codes
	SuccessDonationCategoryCreateRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORY_CREATE_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgDonationCategoryCreateRequestSent,
		Type:       "success",
	}

	SuccessDonationCategoriesFetched = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORIES_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoriesFetched,
		Type:       "success",
	}

	SuccessDonationCategoryFetched = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORY_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryFetched,
		Type:       "success",
	}

	SuccessDonationCategoryUpdated = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORY_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryUpdated,
		Type:       "success",
	}

	SuccessDonationCompanyCreateRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_CREATE_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgDonationCompanyCreateRequestSent,
		Type:       "success",
	}

	SuccessDonationCompaniesFetched = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANIES_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgDonationCompaniesFetched,
		Type:       "success",
	}

	SuccessDonationCompanyFetched = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyFetched,
		Type:       "success",
	}

	SuccessDonationCompanyUpdatedRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_UPDATED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyUpdatedRequestSent,
		Type:       "success",
	}

	SuccessDonationCreateRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_CREATE_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgDonationCreateRequestSent,
		Type:       "success",
	}

	SuccessDonationsFetched = ResponseCode{
		Code:       "SUCCESS_DONATIONS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgDonationsFetched,
		Type:       "success",
	}

	SuccessDonationFetched = ResponseCode{
		Code:       "SUCCESS_DONATION_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgDonationFetched,
		Type:       "success",
	}

	SuccessDonationUpdateRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_UPDATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationUpdateRequestSent,
		Type:       "success",
	}

	SuccessDonationImageUpdateRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_IMAGE_UPDATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationImageUpdateRequestSent,
		Type:       "success",
	}

	SuccessDonationImageDeleteRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_IMAGE_DELETE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationImageDeleteRequestSent,
		Type:       "success",
	}

	SuccessDonationImageAddRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_IMAGE_ADD_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationImageAddRequestSent,
		Type:       "success",
	}

	// CPS Action related success response codes
	SuccessCPSActionsRetrieved = ResponseCode{
		Code:       "SUCCESS_CPS_ACTIONS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionsRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionAuthorized = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_AUTHORIZED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionAuthorizedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionRejected = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_REJECTED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionRejectedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionsRetrievedSuccessfully = ResponseCode{
		Code:       "SUCCESS_CPS_ACTIONS_RETRIEVED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgCPSActionsRetrieved,
		Type:       "success",
	}

	// Wallet related success response codes
	SuccessWalletCreationRequestSent = ResponseCode{
		Code:       "SUCCESS_WALLET_CREATION_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgWalletCreationRequestSent,
		Type:       "success",
	}

	SuccessWalletUpdateRequestSent = ResponseCode{
		Code:       "SUCCESS_WALLET_UPDATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgWalletUpdateRequestSent,
		Type:       "success",
	}

	SuccessWalletDeleted = ResponseCode{
		Code:       "SUCCESS_WALLET_DELETED",
		StatusCode: StatusOK,
		Message:    MsgWalletDeletedSuccessfully,
		Type:       "success",
	}

	SuccessWalletsRetrieved = ResponseCode{
		Code:       "SUCCESS_WALLETS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgWalletsRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessWalletRetrieved = ResponseCode{
		Code:       "SUCCESS_WALLET_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgWalletRetrievedSuccessfully,
		Type:       "success",
	}

	// Event related success response codes
	SuccessEventCreationRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_EVENT_CREATION_REQUEST_SUBMITTED",
		StatusCode: StatusCreated,
		Message:    MsgEventCreationRequestSubmitted,
		Type:       "success",
	}

	SuccessEventUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_EVENT_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgEventUpdateRequestSubmitted,
		Type:       "success",
	}

	SuccessEventDeleteRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_EVENT_DELETE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgEventDeleteRequestSubmitted,
		Type:       "success",
	}

	SuccessEventEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_EVENT_ENABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgEventEnableRequestSubmitted,
		Type:       "success",
	}

	SuccessEventDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_EVENT_DISABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgEventDisableRequestSubmitted,
		Type:       "success",
	}

	SuccessEventRetrieved = ResponseCode{
		Code:       "SUCCESS_EVENT_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgEventSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessEventsRetrieved = ResponseCode{
		Code:       "SUCCESS_EVENTS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgEventsSuccessfullyRetrieved,
		Type:       "success",
	}
	ErrorEventNameRequired = ResponseCode{
		Code:       "ERROR_EVENT_NAME_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Event name is required",
		Type:       "error",
	}

	ErrorEventAlreadyExists = ResponseCode{
		Code:       "ERROR_EVENT_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    "An event with this name already exists",
		Type:       "error",
	}

	ErrorCoverImageRequired = ResponseCode{
		Code:       "ERROR_COVER_IMAGE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Cover image is required",
		Type:       "error",
	}

	ErrorMerchantNotFound = ResponseCode{
		Code:       "ERROR_MERCHANT_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "Merchant not found",
		Type:       "error",
	}

	// Update/Delete/Enable/Disable errors
	ErrorEventNotFound = ResponseCode{
		Code:       "ERROR_EVENT_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "Event not found",
		Type:       "error",
	}

	ErrorEventAlreadyEnabled = ResponseCode{
		Code:       "ERROR_EVENT_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    "Event is already enabled",
		Type:       "error",
	}

	ErrorEventAlreadyDisabled = ResponseCode{
		Code:       "ERROR_EVENT_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    "Event is already disabled",
		Type:       "error",
	}

	// General / internal
	ErrorUnhandledServer = ResponseCode{
		Code:       "ERROR_UNHANDLED_SERVER",
		StatusCode: StatusInternalServerError,
		Message:    "Unhandled server error",
		Type:       "error",
	}
	ErrorInvalidDateFormat = ResponseCode{
		Code:       "ERROR_INVALID_DATE_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    "Invalid date format",
		Type:       "error",
	}
	ErrorInvalidFileUpload = ResponseCode{
		Code:       "ERROR_INVALID_FILE_UPLOAD",
		StatusCode: StatusBadRequest,
		Message:    "Invalid or corrupted file upload",
		Type:       "error",
	}

	ErrorInvalidNumberFormat = ResponseCode{
		Code:       "ERROR_INVALID_NUMBER_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    "Invalid number format",
		Type:       "error",
	}
	ErrorMerchantIDRequired = ResponseCode{
		Code:       "ERROR_MERCHANT_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Merchant ID is required",
		Type:       "error",
	}

	ErrorEventVenueRequired = ResponseCode{
		Code:       "ERROR_EVENT_VENUE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Event venue is required",
		Type:       "error",
	}
	ErrorStartDateRequired = ResponseCode{
		Code:       "ERROR_START_DATE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Start date is required",
		Type:       "error",
	}
	ErrorDueDateRequired = ResponseCode{
		Code:       "ERROR_DUE_DATE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Due date is required",
		Type:       "error",
	}
	ErrorTotalTicketCountRequired = ResponseCode{
		Code:       "ERROR_TOTAL_TICKET_COUNT_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Total ticket count is required",
		Type:       "error",
	}
	ErrorInvalidTicketCount = ResponseCode{
		Code:       "ERROR_INVALID_TICKET_COUNT",
		StatusCode: StatusBadRequest,
		Message:    "Total ticket count must be at least 1",
		Type:       "error",
	}
	ErrorEventCityRequired = ResponseCode{
		Code:       "ERROR_EVENT_CITY_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Event city is required",
		Type:       "error",
	}

	ErrorEventDescriptionRequired = ResponseCode{
		Code:       "ERROR_EVENT_DESCRIPTION_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Event description is required",
		Type:       "error",
	}

	ErrorTicketsRequired = ResponseCode{
		Code:       "ERROR_TICKETS_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Tickets are required",
		Type:       "error",
	}

	// Product Code related success response codes
	SuccessProductCodeUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_PRODUCT_CODE_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgProductCodeUpdateRequestSubmitted,
		Type:       "success",
	}

	SuccessProductCodeRetrieved = ResponseCode{
		Code:       "SUCCESS_PRODUCT_CODE_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgProductCodeSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessProductCodesRetrieved = ResponseCode{
		Code:       "SUCCESS_PRODUCT_CODES_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgProductCodesSuccessfullyRetrieved,
		Type:       "success",
	}

	// Portal Card related success response codes
	SuccessPortalCardsFetched = ResponseCode{
		Code:       "SUCCESS_PORTAL_CARDS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgPortalCardsFetchedSuccessfully,
		Type:       "success",
	}

	// Password Rule related success response codes
	SuccessPasswordRulesFetched = ResponseCode{
		Code:       "SUCCESS_PASSWORD_RULES_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgPasswordRulesSuccessfullyFetched,
		Type:       "success",
	}

	SuccessPasswordRuleUpdateActionFetched = ResponseCode{
		Code:       "SUCCESS_PASSWORD_RULE_UPDATE_ACTION_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgPasswordRuleUpdateActionFetched,
		Type:       "success",
	}

	SuccessUpdateActionFetched = ResponseCode{
		Code:       "SUCCESS_UPDATE_ACTION_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgUpdateActionFetchedSuccessfully,
		Type:       "success",
	}

	// Permission related success response codes
	SuccessPermissionGroupRequestCreated = ResponseCode{
		Code:       "SUCCESS_PERMISSION_GROUP_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgPermissionGroupRequestCreated,
		Type:       "success",
	}

	SuccessPermissionGroupsFetched = ResponseCode{
		Code:       "SUCCESS_PERMISSION_GROUPS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgPermissionGroupsFetched,
		Type:       "success",
	}

	SuccessPermissionGroupFetched = ResponseCode{
		Code:       "SUCCESS_PERMISSION_GROUP_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgPermissionGroupFetched,
		Type:       "success",
	}

	SuccessPermissionGroupRequestUpdated = ResponseCode{
		Code:       "SUCCESS_PERMISSION_GROUP_REQUEST_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgPermissionGroupRequestUpdated,
		Type:       "success",
	}

	SuccessPermissionCategoriesFetched = ResponseCode{
		Code:       "SUCCESS_PERMISSION_CATEGORIES_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgPermissionCategoriesFetched,
		Type:       "success",
	}

	// Mini App Merchant related success response codes
	SuccessMiniAppMerchantCreateRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_MERCHANT_CREATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgMiniAppMerchantCreateRequestSuccessfully,
		Type:       "success",
	}

	SuccessUpdateRequestCreated = ResponseCode{
		Code:       "SUCCESS_UPDATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgUpdateRequestSuccessfullyCreated,
		Type:       "success",
	}

	SuccessDeleteRequestCreated = ResponseCode{
		Code:       "SUCCESS_DELETE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDeleteRequestSuccessfullyCreated,
		Type:       "success",
	}

	SuccessEnableRequestCreated = ResponseCode{
		Code:       "SUCCESS_ENABLE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgEnableRequestSuccessfullyCreated,
		Type:       "success",
	}

	SuccessDisableRequestCreated = ResponseCode{
		Code:       "SUCCESS_DISABLE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDisableRequestSuccessfullyCreated,
		Type:       "success",
	}

	// Notification related success response codes
	SuccessNotificationCreationRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_CREATION_REQUEST_SUBMITTED",
		StatusCode: StatusCreated,
		Message:    MsgNotificationCreationRequestSubmitted,
		Type:       "success",
	}

	SuccessNotificationUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgNotificationUpdateRequestSubmitted,
		Type:       "success",
	}

	SuccessNotificationDeleteRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_DELETE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgNotificationDeleteRequestSubmitted,
		Type:       "success",
	}

	// Department related success response codes
	SuccessDepartmentCreated = ResponseCode{
		Code:       "SUCCESS_DEPARTMENT_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentCreatedSuccessfully,
		Type:       "success",
	}

	SuccessDepartmentUpdateCPSActionCreated = ResponseCode{
		Code:       "SUCCESS_DEPARTMENT_UPDATE_CPS_ACTION_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentUpdateCPSActionCreated,
		Type:       "success",
	}

	// Budget Category related success response codes
	SuccessBudgetCategoryUpdated = ResponseCode{
		Code:       "SUCCESS_BUDGET_CATEGORY_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoryUpdatedSuccessfully,
		Type:       "success",
	}

	// BPS User related success response codes
	SuccessBPSUserApproved = ResponseCode{
		Code:       "SUCCESS_BPS_USER_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgBPSUserApprovedSuccessfully,
		Type:       "success",
	}

	// Ad related success response codes
	SuccessAdvertCreated = ResponseCode{
		Code:       "SUCCESS_ADVERT_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgAdvertCreatedSuccessfully,
		Type:       "success",
	}

	SuccessAdvertCreateRequestSent = ResponseCode{
		Code:       "SUCCESS_ADVERT_CREATED_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgAdvertCreatedRequestSent,
		Type:       "success",
	}

	SuccessAdvertUpdated = ResponseCode{
		Code:       "SUCCESS_ADVERT_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgAdvertUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessAdvertUpdateRequestSent = ResponseCode{
		Code:       "SUCCESS_ADVERT_UPDATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgAdvertUpdateRequestSent,
		Type:       "success",
	}

	SuccessAdvertDeleted = ResponseCode{
		Code:       "SUCCESS_ADVERT_DELETED",
		StatusCode: StatusOK,
		Message:    MsgAdvertDeletedSuccessfully,
		Type:       "success",
	}
	SuccessAdvertDeleteRequestSent = ResponseCode{
		Code:       "SUCCESS_ADVERT_DELETE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgAdvertDeleteRequestSent,
		Type:       "success",
	}

	SuccessAdvertFetched = ResponseCode{
		Code:       "SUCCESS_ADVERT_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgAdvertFetchedSuccessfully,
		Type:       "success",
	}

	SuccessAdvertsFetched = ResponseCode{
		Code:       "SUCCESS_ADVERTS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgAdvertsFetchedSuccessfully,
		Type:       "success",
	}

	SuccessAdvertEnableRequestSent = ResponseCode{
		Code:       "SUCCESS_ADVERT_ENABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgAdvertEnableRequestSent,
		Type:       "success",
	}

	SuccessAdvertDisableRequestSent = ResponseCode{
		Code:       "SUCCESS_ADVERT_DISABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgAdvertDisableRequestSent,
		Type:       "success",
	}
	// Account Validation related success response codes
	SuccessValidationRuleApproved = ResponseCode{
		Code:       "SUCCESS_VALIDATION_RULE_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgValidationRuleApprovedSuccessfully,
		Type:       "success",
	}
	SuccessValidationRuleFetched = ResponseCode{
		Code:       "SUCCESS_VALIDATION_RULE_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgValidationRuleApprovedSuccess,
		Type:       "success",
	}


	// Feedback related success response codes
	SuccessFeedbackSaved = ResponseCode{
		Code:       "SUCCESS_FEEDBACK_SAVED",
		StatusCode: StatusCreated,
		Message:    MsgFeedbackSavedSuccessfully,
		Type:       "success",
	}

	// HQ related success response codes
	SuccessHQApproved = ResponseCode{
		Code:       "SUCCESS_HQ_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgHQApprovedSuccessfully,
		Type:       "success",
	}

	// Mini App related success response codes
	SuccessMiniAppDetailsFetched = ResponseCode{
		Code:       "SUCCESS_MINI_APP_DETAILS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppDetailsFetchedSuccessfully,
		Type:       "success",
	}

	SuccessCredentialInformationGenerated = ResponseCode{
		Code:       "SUCCESS_CREDENTIAL_INFORMATION_GENERATED",
		StatusCode: StatusOK,
		Message:    MsgCredentialInformationGenerated,
		Type:       "success",
	}

	SuccessFileUploadedToMinIO = ResponseCode{
		Code:       "SUCCESS_FILE_UPLOADED_TO_MINIO",
		StatusCode: StatusOK,
		Message:    MsgFileUploadedToMinIOSuccessfully,
		Type:       "success",
	}

	// Key Generation related success response codes
	SuccessEd25519KeyPairGenerated = ResponseCode{
		Code:       "SUCCESS_ED25519_KEY_PAIR_GENERATED",
		StatusCode: StatusOK,
		Message:    MsgEd25519KeyPairGeneratedSuccessfully,
		Type:       "success",
	}

	SuccessDataSigned = ResponseCode{
		Code:       "SUCCESS_DATA_SIGNED",
		StatusCode: StatusOK,
		Message:    MsgDataSignedSuccessfully,
		Type:       "success",
	}

	SuccessSignatureVerification = ResponseCode{
		Code:       "SUCCESS_SIGNATURE_VERIFICATION",
		StatusCode: StatusOK,
		Message:    MsgSignatureVerificationSuccessful,
		Type:       "success",
	}

	SuccessFabricIDGenerated = ResponseCode{
		Code:       "SUCCESS_FABRIC_ID_GENERATED",
		StatusCode: StatusOK,
		Message:    MsgFabricIDGeneratedSuccessfully,
		Type:       "success",
	}

	SuccessNumericCodeGenerated = ResponseCode{
		Code:       "SUCCESS_NUMERIC_CODE_GENERATED",
		StatusCode: StatusOK,
		Message:    MsgNumericCodeGeneratedSuccessfully,
		Type:       "success",
	}

	SuccessAppSecretGenerated = ResponseCode{
		Code:       "SUCCESS_APP_SECRET_GENERATED",
		StatusCode: StatusOK,
		Message:    MsgAppSecretGeneratedSuccessfully,
		Type:       "success",
	}

	SuccessAppSecretHashed = ResponseCode{
		Code:       "SUCCESS_APP_SECRET_HASHED",
		StatusCode: StatusOK,
		Message:    MsgAppSecretHashedSuccessfully,
		Type:       "success",
	}

	SuccessAppSecretVerification = ResponseCode{
		Code:       "SUCCESS_APP_SECRET_VERIFICATION",
		StatusCode: StatusOK,
		Message:    MsgAppSecretVerificationSuccessful,
		Type:       "success",
	}

	SuccessPrefixedNameGenerated = ResponseCode{
		Code:       "SUCCESS_PREFIXED_NAME_GENERATED",
		StatusCode: StatusOK,
		Message:    MsgPrefixedNameGeneratedSuccessfully,
		Type:       "success",
	}

	SuccessSMSSent = ResponseCode{
		Code:       "SUCCESS_SMS_SENT",
		StatusCode: StatusOK,
		Message:    MsgSMSsentSuccessfully,
		Type:       "success",
	}

	// Service related success response codes
	SuccessServiceRetrieved = ResponseCode{
		Code:       "SUCCESS_SERVICE_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgServiceRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessServiceTireUpdated = ResponseCode{
		Code:       "SUCCESS_SERVICE_TIRE_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgServiceTireUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessServiceSingleTransferUpdated = ResponseCode{
		Code:       "SUCCESS_SERVICE_SINGLE_TRANSFER_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgServiceSingleTransferUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessServiceMinAmountTransferUpdated = ResponseCode{
		Code:       "SUCCESS_SERVICE_MIN_AMOUNT_TRANSFER_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgServiceMinAmountTransferUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessServiceTierDeleteRequest = ResponseCode{
		Code:       "SUCCESS_SERVICE_TIER_DELETE_REQUEST",
		StatusCode: StatusOK,
		Message:    MsgServiceTierDeleteRequestSuccessfully,
		Type:       "success",
	}

	SuccessServiceFeeDetailFetched = ResponseCode{
		Code:       "SUCCESS_SERVICE_FEE_DETAIL_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgServiceFeeDetailFetchedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionCreatedForUpdateServiceFee = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CREATED_FOR_UPDATE_SERVICE_FEE",
		StatusCode: StatusCreated,
		Message:    MsgCPSActionCreatedForUpdateServiceFee,
		Type:       "success",
	}

	SuccessCPSActionCreatedForUpdateSingleMaxTransfer = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CREATED_FOR_UPDATE_SINGLE_MAX_TRANSFER",
		StatusCode: StatusCreated,
		Message:    MsgCPSActionCreatedForUpdateSingleMaxTransfer,
		Type:       "success",
	}

	SuccessCPSActionCreatedForUpdateTotalMaxTransferCap = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CREATED_FOR_UPDATE_TOTAL_MAX_TRANSFER_CAP",
		StatusCode: StatusCreated,
		Message:    MsgCPSActionCreatedForUpdateTotalMaxTransferCap,
		Type:       "success",
	}

	SuccessCPSActionCreatedForUpdateMinimumTransferCap = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CREATED_FOR_UPDATE_MINIMUM_TRANSFER_CAP",
		StatusCode: StatusCreated,
		Message:    MsgCPSActionCreatedForUpdateMinimumTransferCap,
		Type:       "success",
	}

	SuccessCPSActionCreatedForDeleteServiceFeeTire = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CREATED_FOR_DELETE_SERVICE_FEE_TIRE",
		StatusCode: StatusCreated,
		Message:    MsgCPSActionCreatedForDeleteServiceFeeTire,
		Type:       "success",
	}

	SuccessCPSActionCreated = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgCPSActionCreatedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionUpdated = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionStatusUpdated = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_STATUS_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionStatusUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionFetched = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionFetchedSuccessfully,
		Type:       "success",
	}

	SuccessPreviousServiceDataFetched = ResponseCode{
		Code:       "SUCCESS_PREVIOUS_SERVICE_DATA_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgPreviousServiceDataFetchedSuccessfully,
		Type:       "success",
	}

	SuccessMaxTotalCapValidated = ResponseCode{
		Code:       "SUCCESS_MAX_TOTAL_CAP_VALIDATED",
		StatusCode: StatusOK,
		Message:    MsgMaxTotalCapValidatedSuccessfully,
		Type:       "success",
	}

	// Product Code related success response codes
	SuccessProductCodeFetched = ResponseCode{
		Code:       "SUCCESS_PRODUCT_CODE_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgProductCodeFetchedSuccessfully,
		Type:       "success",
	}

	SuccessProductCodesFetched = ResponseCode{
		Code:       "SUCCESS_PRODUCT_CODES_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgProductCodesFetchedSuccessfully,
		Type:       "success",
	}

	SuccessProductCodeUpdated = ResponseCode{
		Code:       "SUCCESS_PRODUCT_CODE_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgProductCodeUpdatedSuccessfully,
		Type:       "success",
	}

	// Mini App Merchant related success response codes
	SuccessMiniAppAdded = ResponseCode{
		Code:       "SUCCESS_MINI_APP_ADDED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppAddedSuccessfully,
		Type:       "success",
	}

	SuccessMiniAppEnabledStateUpdated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_ENABLED_STATE_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppEnabledStateUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessMiniAppSoftDeleted = ResponseCode{
		Code:       "SUCCESS_MINI_APP_SOFT_DELETED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppSoftDeletedSuccessfully,
		Type:       "success",
	}

	SuccessTransactionAborted = ResponseCode{
		Code:       "SUCCESS_TRANSACTION_ABORTED",
		StatusCode: StatusOK,
		Message:    MsgTransactionAbortedSuccessfully,
		Type:       "success",
	}

	SuccessTransactionCommitted = ResponseCode{
		Code:       "SUCCESS_TRANSACTION_COMMITTED",
		StatusCode: StatusOK,
		Message:    MsgTransactionCommittedSuccessfully,
		Type:       "success",
	}

	// BPS Calls related success response codes
	SuccessLinkedAccountFetched = ResponseCode{
		Code:       "SUCCESS_LINKED_ACCOUNT_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgLinkedAccountFetchedSuccessfully,
		Type:       "success",
	}

	// Unlink related success response codes
	SuccessUnlinkUserRetrieved = ResponseCode{
		Code:       "SUCCESS_UNLINK_USER_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgUserSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessUnlinkCifRequestSent = ResponseCode{
		Code:       "SUCCESS_UNLINK_CIF_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgUnlinkCifRequestSentSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionCreatedForUnlink = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CREATED_FOR_UNLINK",
		StatusCode: StatusCreated,
		Message:    MsgCPSActionCreatedForUnlinkSuccessfully,
		Type:       "success",
	}

	SuccessDocumentHardDeleted = ResponseCode{
		Code:       "SUCCESS_DOCUMENT_HARD_DELETED",
		StatusCode: StatusOK,
		Message:    MsgDocumentHardDeletedSuccessfully,
		Type:       "success",
	}

	// Notification related success response codes
	SuccessNotificationEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_ENABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgNotificationEnableRequestSubmitted,
		Type:       "success",
	}

	SuccessNotificationDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_DISABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgNotificationDisableRequestSubmitted,
		Type:       "success",
	}

	SuccessNotificationRetrieved = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgNotificationSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessNotificationsRetrieved = ResponseCode{
		Code:       "SUCCESS_NOTIFICATIONS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgNotificationsSuccessfullyRetrieved,
		Type:       "success",
	}

	// Mini App Handler related success response codes
	SuccessUpdateMiniAppRequestCreated = ResponseCode{
		Code:       "SUCCESS_UPDATE_MINI_APP_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgUpdateMiniAppRequestCreatedSuccessfully,
		Type:       "success",
	}

	SuccessMiniAppDeletedRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_DELETED_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgMiniAppDeletedRequestCreatedSuccessfully,
		Type:       "success",
	}

	SuccessMiniAppsRetrieved = ResponseCode{
		Code:       "SUCCESS_MINI_APPS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppsSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessMiniAppFetchedByID = ResponseCode{
		Code:       "SUCCESS_MINI_APP_FETCHED_BY_ID",
		StatusCode: StatusOK,
		Message:    MsgMiniAppFetchedByIDSuccessfully,
		Type:       "success",
	}

	SuccessMiniAppActionCompleted = ResponseCode{
		Code:       "SUCCESS_MINI_APP_ACTION_COMPLETED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppActionCompletedSuccessfully,
		Type:       "success",
	}

	// Feedback Handler related success response codes
	SuccessFeedbackCreated = ResponseCode{
		Code:       "SUCCESS_FEEDBACK_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgFeedbackCreatedSuccessfully,
		Type:       "success",
	}

	SuccessFeedbackFetched = ResponseCode{
		Code:       "SUCCESS_FEEDBACK_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgFeedbackFetchedSuccessfully,
		Type:       "success",
	}

	// HQ related success response codes
	SuccessHQFetched = ResponseCode{
		Code:       "SUCCESS_HQ_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgHQFetchedSuccessfully,
		Type:       "success",
	}

	SuccessHQsFetched = ResponseCode{
		Code:       "SUCCESS_HQS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgHQsFetchedSuccessfully,
		Type:       "success",
	}

	// Wallet related success response codes
	SuccessWalletActionRequestSent = ResponseCode{
		Code:       "SUCCESS_WALLET_ACTION_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgWalletActionRequestSentSuccessfully,
		Type:       "success",
	}

	// Donation Service related success response codes
	SuccessDonationCategoryIconUploaded = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORY_ICON_UPLOADED",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryIconUploadSuccess,
		Type:       "success",
	}

	SuccessDonationCategoryIconUpdated = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORY_ICON_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryIconUpdateSuccess,
		Type:       "success",
	}

	SuccessDonationCPSActionApproved = ResponseCode{
		Code:       "SUCCESS_DONATION_CPS_ACTION_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgDonationCPSActionApprovedSuccess,
		Type:       "success",
	}

	SuccessDonationCompanyLogoUploaded = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_LOGO_UPLOADED",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyLogoUploadSuccess,
		Type:       "success",
	}

	SuccessDonationCompanyLogoUpdated = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_LOGO_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyLogoUpdateSuccess,
		Type:       "success",
	}

	SuccessDonationImageUploaded = ResponseCode{
		Code:       "SUCCESS_DONATION_IMAGE_UPLOADED",
		StatusCode: StatusOK,
		Message:    MsgDonationImageUploadSuccess,
		Type:       "success",
	}

	SuccessDonationImagesUpdated = ResponseCode{
		Code:       "SUCCESS_DONATION_IMAGES_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDonationImagesUpdateSuccess,
		Type:       "success",
	}

	// Notification Service related success response codes
	SuccessNotificationConstructed = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_CONSTRUCTED",
		StatusCode: StatusOK,
		Message:    MsgNotificationConstructedSuccess,
		Type:       "success",
	}

	SuccessNotificationUpdated = ResponseCode{
		Code:       "SUCCESS_NOTIFICATION_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgNotificationUpdateSuccess,
		Type:       "success",
	}

	// Feedback Kafka related success response codes
	SuccessFeedbackSavedToDatabase = ResponseCode{
		Code:       "SUCCESS_FEEDBACK_SAVED_TO_DATABASE",
		StatusCode: StatusOK,
		Message:    MsgFeedbackSavedToDatabaseSuccess,
		Type:       "success",
	}
)

// Error Response Codes
var (
	ErrorUserNotFound = ResponseCode{
		Code:       "ERROR_USER_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgUserNotFound,
		Type:       "error",
	}

	ErrorHealthCheck = ResponseCode{
		Code:       "FAILED_HEALTH_CHECK",
		StatusCode: StatusOK,
		Message:    MsgHealthCheckFailed,
		Type:       "error",
	}

	ErrorActionAlreadyExists = ResponseCode{
		Code:       "ERROR_ACTION_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgActionAlreadyExists,
		Type:       "error",
	}

	ErrorAccountNumberRequired = ResponseCode{
		Code:       "ERROR_ACCOUNT_NUMBER_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MSGAccountNumberRequired,
		Type:       "error",
	}
	ErrorFeedbackIDRequired = ResponseCode{
		Code:       "ERROR_FEEDBACK_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgFeedbackIDRequired,
		Type:       "error",
	}

	ErrorInvalidIDFormat = ResponseCode{
		Code:       "ERROR_INVALID_ID_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidIDFormat,
		Type:       "error",
	}

	ErrorIncompleteUserInfo = ResponseCode{
		Code:       "ERROR_INCOMPLET_USER_INFO",
		StatusCode: StatusBadRequest,
		Message:    MSGIncompleteUserInfo,
		Type:       "error",
	}

	ErrorUserCodeRequired = ResponseCode{
		Code:       "ERROR_USER_CODE_IS_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MSGUserCodeIsRequired,
		Type:       "error",
	}

	ErrorUserAlreadyExists = ResponseCode{
		Code:       "ERROR_USER_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgUserAlreadyExists,
		Type:       "error",
	}

	ErrorUserUnauthorized = ResponseCode{
		Code:       "ERROR_USER_UNAUTHORIZED",
		StatusCode: StatusUnauthorized,
		Message:    MsgUserUnauthorized,
		Type:       "error",
	}

	ErrorUserForbidden = ResponseCode{
		Code:       "ERROR_USER_FORBIDDEN",
		StatusCode: StatusForbidden,
		Message:    MsgUserForbidden,
		Type:       "error",
	}

	ErrorUserInvalidCredentials = ResponseCode{
		Code:       "ERROR_USER_INVALID_CREDENTIALS",
		StatusCode: StatusUnauthorized,
		Message:    MsgUserInvalidCredentials,
		Type:       "error",
	}

	ErrorUserAccountBlocked = ResponseCode{
		Code:       "ERROR_USER_ACCOUNT_BLOCKED",
		StatusCode: StatusForbidden,
		Message:    MsgUserAccountBlocked,
		Type:       "error",
	}

	ErrorUserSessionExpired = ResponseCode{
		Code:       "ERROR_USER_SESSION_EXPIRED",
		StatusCode: StatusUnauthorized,
		Message:    MsgUserSessionExpired,
		Type:       "error",
	}

	ErrorUserTooManyLoginAttempts = ResponseCode{
		Code:       "ERROR_USER_TOO_MANY_LOGIN_ATTEMPTS",
		StatusCode: StatusTooManyRequests,
		Message:    MsgUserTooManyLoginAttempts,
		Type:       "error",
	}

	ErrorUserDeviceMismatch = ResponseCode{
		Code:       "ERROR_USER_DEVICE_MISMATCH",
		StatusCode: StatusConflict,
		Message:    MsgUserDeviceMismatch,
		Type:       "error",
	}

	ErrorUserDeviceNotFound = ResponseCode{
		Code:       "ERROR_USER_DEVICE_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgUserDeviceNotFound,
		Type:       "error",
	}

	ErrorUserDeviceAlreadyExists = ResponseCode{
		Code:       "ERROR_USER_DEVICE_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgUserDeviceAlreadyExists,
		Type:       "error",
	}

	ErrorActionNotFound = ResponseCode{
		Code:       "ERROR_ACTION_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgCPSActionNotFound,
		Type:       "error",
	}

	ErrorPendingCpsActionExists =  ResponseCode{
		Code:       "ERROR_PENDING_CPS_ACTION_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgPendingCPSActionExists,
		Type:       "error",
	}

	ErrorUserAlreadyEnabled = ResponseCode{
		Code:       "ERROR_USER_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgBpsUserAlreadyEnabled,
		Type:       "error",
	}

	ErrorUserAlreadyDisabled = ResponseCode{
		Code:       "ERROR_USER_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgBpsUserAlreadyDisabled,
		Type:       "error",
	}

	ErrorOTPNotFound = ResponseCode{
		Code:       "ERROR_OTP_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgOTPNotFound,
		Type:       "error",
	}

	ErrorInvalidRequest = ResponseCode{
		Code:       "ERROR_INVALID_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestOnParam,
		Type:       "error",
	}

	ErrorNoDataProvidedForUpdate = ResponseCode{
		Code:       "ERROR_NO_DATA_PROVIDED_FOR_UPDATE",
		StatusCode: StatusBadRequest,
		Message:    "No data provided for update",
		Type:       "error",
	}

	ErrorEventIDRequired = ResponseCode{
		Code:       "ERROR_EVENT_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Event ID is required",
		Type:       "error",
	}
	ErrorUpdateEventEmptyPayload = ResponseCode{
		Code:       "ERROR_UPDATE_EVENT_EMPTY_PAYLOAD",
		StatusCode: StatusBadRequest,
		Message:    "No data provided to update the event",
		Type:       "error",
	}

	ErrorOTPExpired = ResponseCode{
		Code:       "ERROR_OTP_EXPIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgOTPExpired,
		Type:       "error",
	}

	ErrorOTPInvalid = ResponseCode{
		Code:       "ERROR_OTP_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgOTPInvalid,
		Type:       "error",
	}

	ErrorOTPAlreadyExists = ResponseCode{
		Code:       "ERROR_OTP_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgOTPAlreadyExists,
		Type:       "error",
	}

	ErrorOTPTooManyAttempts = ResponseCode{
		Code:       "ERROR_OTP_TOO_MANY_ATTEMPTS",
		StatusCode: StatusTooManyRequests,
		Message:    MsgOTPTooManyAttempts,
		Type:       "error",
	}

	ErrorOTPSendFailed = ResponseCode{
		Code:       "ERROR_OTP_SEND_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgOTPSendFailed,
		Type:       "error",
	}

	ErrorPINInvalid = ResponseCode{
		Code:       "ERROR_PIN_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgPINInvalid,
		Type:       "error",
	}

	ErrorPINMismatch = ResponseCode{
		Code:       "ERROR_PIN_MISMATCH",
		StatusCode: StatusBadRequest,
		Message:    MsgPINMismatch,
		Type:       "error",
	}

	ErrorPINTooWeak = ResponseCode{
		Code:       "ERROR_PIN_TOO_WEAK",
		StatusCode: StatusBadRequest,
		Message:    MsgPINTooWeak,
		Type:       "error",
	}

	ErrorPINInHistory = ResponseCode{
		Code:       "ERROR_PIN_IN_HISTORY",
		StatusCode: StatusBadRequest,
		Message:    MsgPINInHistory,
		Type:       "error",
	}

	ErrorPINLengthInvalid = ResponseCode{
		Code:       "ERROR_PIN_LENGTH_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgPINLengthInvalid,
		Type:       "error",
	}

	ErrorPINOnlyDigits = ResponseCode{
		Code:       "ERROR_PIN_ONLY_DIGITS",
		StatusCode: StatusBadRequest,
		Message:    MsgPINOnlyDigits,
		Type:       "error",
	}

	ErrorPINSequential = ResponseCode{
		Code:       "ERROR_PIN_SEQUENTIAL",
		StatusCode: StatusBadRequest,
		Message:    MsgPINSequential,
		Type:       "error",
	}

	ErrorPINRedundant = ResponseCode{
		Code:       "ERROR_PIN_REDUNDANT",
		StatusCode: StatusBadRequest,
		Message:    MsgPINRedundant,
		Type:       "error",
	}

	ErrorSamePIN = ResponseCode{
		Code:       "ERROR_SAME_PIN",
		StatusCode: StatusBadRequest,
		Message:    MsgSamePIN,
		Type:       "error",
	}

	ErrorSessionNotFound = ResponseCode{
		Code:       "ERROR_SESSION_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgSessionNotFound,
		Type:       "error",
	}

	ErrorSessionExpired = ResponseCode{
		Code:       "ERROR_SESSION_EXPIRED",
		StatusCode: StatusUnauthorized,
		Message:    MsgSessionExpired,
		Type:       "error",
	}

	ErrorSessionInvalid = ResponseCode{
		Code:       "ERROR_SESSION_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgSessionInvalid,
		Type:       "error",
	}

	ErrorSessionRetrievalFailed = ResponseCode{
		Code:       "ERROR_SESSION_RETRIEVAL_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgSessionRetrievalFailed,
		Type:       "error",
	}

	ErrorSessionCreationFailed = ResponseCode{
		Code:       "ERROR_SESSION_CREATION_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgSessionCreationFailed,
		Type:       "error",
	}

	ErrorSessionUpdateFailed = ResponseCode{
		Code:       "ERROR_SESSION_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgSessionUpdateFailed,
		Type:       "error",
	}

	ErrorFileTooLarge = ResponseCode{
		Code:       "ERROR_FILE_TOO_LARGE",
		StatusCode: StatusRequestEntityTooLarge,
		Message:    MsgFileTooLarge,
		Type:       "error",
	}

	ErrorInvalidToken = ResponseCode{
		Code:       "ERROR_INVALID_TOKEN",
		StatusCode: StatusForbidden,
		Message:    MsgInvalidToken,
		Type:       "error",
	}

	ErrorFileInvalidType = ResponseCode{
		Code:       "ERROR_FILE_INVALID_TYPE",
		StatusCode: StatusBadRequest,
		Message:    MsgFileInvalidType,
		Type:       "error",
	}

	ErrorFileUploadFailed = ResponseCode{
		Code:       "ERROR_FILE_UPLOAD_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgFileUploadFailed,
		Type:       "error",
	}

	ErrorFileNotFound = ResponseCode{
		Code:       "ERROR_FILE_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgFileNotFound,
		Type:       "error",
	}

	ErrorFileDeleteFailed = ResponseCode{
		Code:       "ERROR_FILE_DELETE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgFileDeleteFailed,
		Type:       "error",
	}

	ErrorValidationFailed = ResponseCode{
		Code:       "ERROR_VALIDATION_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgValidationFailed,
		Type:       "error",
	}

	ErrorRequiredFieldMissing = ResponseCode{
		Code:       "ERROR_REQUIRED_FIELD_MISSING",
		StatusCode: StatusBadRequest,
		Message:    MsgRequiredFieldMissing,
		Type:       "error",
	}

	ErrorInvalidFormat = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidFormat,
		Type:       "error",
	}

	ErrorInvalidEmail = ResponseCode{
		Code:       "ERROR_INVALID_EMAIL",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidEmail,
		Type:       "error",
	}

	ErrorInvalidPhoneNumber = ResponseCode{
		Code:       "ERROR_INVALID_PHONE_NUMBER",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidPhoneNumber,
		Type:       "error",
	}

	ErrorInvalidDate = ResponseCode{
		Code:       "ERROR_INVALID_DATE",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidDate,
		Type:       "error",
	}

	ErrorInvalidActionData = ResponseCode{
		Code:       "ERROR_INVALID_ACTION_DATA",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidActionData,
		Type:       "error",
	}
	ErrorUnsupportedAction = ResponseCode{
		Code:       "ERROR_UNSUPPORTED_ACTION",
		StatusCode: StatusBadRequest,
		Message:    MsgUnsupportedAction,
		Type:       "error",
	}

	ErrorInvalidID = ResponseCode{
		Code:       "ERROR_INVALID_ID",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidID,
		Type:       "error",
	}

	ErrorMissingFile = ResponseCode{
		Code:       "ERROR_MISSING_FILE",
		StatusCode: StatusBadRequest,
		Message:    MsgMissingFile,
		Type:       "error",
	}

	ErrorFieldTooLong = ResponseCode{
		Code:       "ERROR_FIELD_TOO_LONG",
		StatusCode: StatusBadRequest,
		Message:    MsgFieldTooLong,
		Type:       "error",
	}

	ErrorFieldTooShort = ResponseCode{
		Code:       "ERROR_FIELD_TOO_SHORT",
		StatusCode: StatusBadRequest,
		Message:    MsgFieldTooShort,
		Type:       "error",
	}

	ErrorInternalServerError = ResponseCode{
		Code:       "ERROR_INTERNAL_SERVER_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgInternalServerError,
		Type:       "error",
	}
	ErrorPendingCPSAction = ResponseCode{
		Code:       "ERROR_PENDING_CPS_ACTION_PRESENT",
		StatusCode: StatusConflict,
		Message:    MsgPendingCPSActionExists,
		Type:       "error",
	}

	ErrorServiceUnavailable = ResponseCode{
		Code:       "ERROR_SERVICE_UNAVAILABLE",
		StatusCode: StatusServiceUnavailable,
		Message:    MsgServiceUnavailable,
		Type:       "error",
	}

	ErrorDatabaseError = ResponseCode{
		Code:       "ERROR_DATABASE_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgDatabaseError,
		Type:       "error",
	}

	ErrorNetworkError = ResponseCode{
		Code:       "ERROR_NETWORK_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgNetworkError,
		Type:       "error",
	}

	ErrorTimeoutError = ResponseCode{
		Code:       "ERROR_TIMEOUT_ERROR",
		StatusCode: StatusRequestTimeout,
		Message:    MsgTimeoutError,
		Type:       "error",
	}

	ErrorUnexpectedError = ResponseCode{
		Code:       "ERROR_UNEXPECTED_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgUnexpectedError,
		Type:       "error",
	}

	ErrorConfigurationError = ResponseCode{
		Code:       "ERROR_CONFIGURATION_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgConfigurationError,
		Type:       "error",
	}

	ErrorExternalServiceError = ResponseCode{
		Code:       "ERROR_EXTERNAL_SERVICE_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgExternalServiceError,
		Type:       "error",
	}

	ErrorInsufficientBalance = ResponseCode{
		Code:       "ERROR_INSUFFICIENT_BALANCE",
		StatusCode: StatusBadRequest,
		Message:    MsgInsufficientBalance,
		Type:       "error",
	}

	ErrorTransactionFailed = ResponseCode{
		Code:       "ERROR_TRANSACTION_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgTransactionFailed,
		Type:       "error",
	}

	ErrorLimitExceeded = ResponseCode{
		Code:       "ERROR_LIMIT_EXCEEDED",
		StatusCode: StatusBadRequest,
		Message:    MsgLimitExceeded,
		Type:       "error",
	}

	ErrorOperationNotAllowed = ResponseCode{
		Code:       "ERROR_OPERATION_NOT_ALLOWED",
		StatusCode: StatusForbidden,
		Message:    MsgOperationNotAllowed,
		Type:       "error",
	}

	ErrorResourceBusy = ResponseCode{
		Code:       "ERROR_RESOURCE_BUSY",
		StatusCode: StatusConflict,
		Message:    MsgResourceBusy,
		Type:       "error",
	}

	ErrorMaintenanceMode = ResponseCode{
		Code:       "ERROR_MAINTENANCE_MODE",
		StatusCode: StatusServiceUnavailable,
		Message:    MsgMaintenanceMode,
		Type:       "error",
	}

	// Service related error response codes
	ErrorServiceFetchFailed = ResponseCode{
		Code:       "ERROR_SERVICE_FETCH_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceFetchFailed,
		Type:       "error",
	}

	ErrorServiceCountFailed = ResponseCode{
		Code:       "ERROR_SERVICE_COUNT_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceCountFailed,
		Type:       "error",
	}

	ErrorServicePendingActionCheckFailed = ResponseCode{
		Code:       "ERROR_SERVICE_PENDING_ACTION_CHECK_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServicePendingActionCheckFailed,
		Type:       "error",
	}

	ErrorServicePreviousDataFetchFailed = ResponseCode{
		Code:       "ERROR_SERVICE_PREVIOUS_DATA_FETCH_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServicePreviousDataFetchFailed,
		Type:       "error",
	}

	ErrorServiceCPSActionCreationFailed = ResponseCode{
		Code:       "ERROR_SERVICE_CPS_ACTION_CREATION_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceCPSActionCreationFailed,
		Type:       "error",
	}

	ErrorServiceValidationFailed = ResponseCode{
		Code:       "ERROR_SERVICE_VALIDATION_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgServiceValidationFailed,
		Type:       "error",
	}

	ErrorServiceHQDataFetchFailed = ResponseCode{
		Code:       "ERROR_SERVICE_HQ_DATA_FETCH_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceHQDataFetchFailed,
		Type:       "error",
	}

	ErrorServiceUnmarshalFailed = ResponseCode{
		Code:       "ERROR_SERVICE_UNMARSHAL_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgServiceUnmarshalFailed,
		Type:       "error",
	}

	ErrorServiceUpdateFailed = ResponseCode{
		Code:       "ERROR_SERVICE_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceUpdateFailed,
		Type:       "error",
	}

	ErrorServiceTierNotFound = ResponseCode{
		Code:       "ERROR_SERVICE_TIER_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgServiceTierNotFound,
		Type:       "error",
	}

	ErrorServiceApproveFailed = ResponseCode{
		Code:       "ERROR_SERVICE_APPROVE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceApproveFailed,
		Type:       "error",
	}

	ErrorServiceAuthorizeFailed = ResponseCode{
		Code:       "ERROR_SERVICE_AUTHORIZE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceAuthorizeFailed,
		Type:       "error",
	}

	ErrorServiceCPSActionStatusUpdateFailed = ResponseCode{
		Code:       "ERROR_SERVICE_CPS_ACTION_STATUS_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceCPSActionStatusUpdateFailed,
		Type:       "error",
	}

	ErrorServiceUnhandledServerError = ResponseCode{
		Code:       "ERROR_SERVICE_UNHANDLED_SERVER_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceUnhandledServerError,
		Type:       "error",
	}

	ErrorServiceUpdateFailedError = ResponseCode{
		Code:       "ERROR_SERVICE_UPDATE_FAILED_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceUpdateFailedError,
		Type:       "error",
	}
	
	ErrorServiceAuthorizeDeleteFailed = ResponseCode{
		Code:       "ERROR_SERVICE_AUTHORIZE_DELETE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgServiceAuthorizeDeleteFailed,
		Type:       "error",
	}

	ErrorServiceUnknownRequestAction = ResponseCode{
		Code:       "ERROR_SERVICE_UNKNOWN_REQUEST_ACTION",
		StatusCode: StatusBadRequest,
		Message:    MsgServiceUnknownRequestAction,
		Type:       "error",
	}

	// Product Code related error response codes
	ErrorProductCodeParseFailed = ResponseCode{
		Code:       "ERROR_PRODUCT_CODE_PARSE_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgProductCodeParseFailed,
		Type:       "error",
	}

	ErrorProductCodeDatabaseQueryFailed = ResponseCode{
		Code:       "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgProductCodeDatabaseQueryFailed,
		Type:       "error",
	}

	ErrorProductCodeCountFailed = ResponseCode{
		Code:       "ERROR_PRODUCT_CODE_COUNT_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgProductCodeCountFailed,
		Type:       "error",
	}

	ErrorProductCodeDatabaseUpdateFailed = ResponseCode{
		Code:       "ERROR_PRODUCT_CODE_DATABASE_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgProductCodeDatabaseUpdateFailed,
		Type:       "error",
	}

	// Mini App Merchant related error response codes
	ErrorMiniAppMerchantCheckPendingFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_CHECK_PENDING_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantCheckPendingFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantExistsCheckFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_EXISTS_CHECK_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantExistsCheckFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantFetchCategoriesFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_FETCH_CATEGORIES_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantFetchCategoriesFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantUnexpectedDBError = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_UNEXPECTED_DB_ERROR",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantUnexpectedDBError,
		Type:       "error",
	}

	ErrorMiniAppMerchantMarshalFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_MARSHAL_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantMarshalFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantUnmarshalFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_UNMARSHAL_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantUnmarshalFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantCreateFromActionFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_CREATE_FROM_ACTION_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantCreateFromActionFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantUnmarshalActionFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_UNMARSHAL_ACTION_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantUnmarshalActionFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantUpdateFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantUpdateFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantFetchPermissionsFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_FETCH_PERMISSIONS_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantFetchPermissionsFailed,
		Type:       "error",
	}

	// Notification related error response codes
	ErrorNotificationMapFailed = ResponseCode{
		Code:       "ERROR_NOTIFICATION_MAP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgNotificationMapFailed,
		Type:       "error",
	}

	ErrorNotificationInsertFailed = ResponseCode{
		Code:       "ERROR_NOTIFICATION_INSERT_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgNotificationInsertFailed,
		Type:       "error",
	}

	ErrorNotificationInvalidIDFormat = ResponseCode{
		Code:       "ERROR_NOTIFICATION_INVALID_ID_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    MsgNotificationInvalidIDFormat,
		Type:       "error",
	}

	ErrorNotificationFetchFailed = ResponseCode{
		Code:       "ERROR_NOTIFICATION_FETCH_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgNotificationFetchFailed,
		Type:       "error",
	}

	// Customer related error response codes
	ErrorCustomerCountFailed = ResponseCode{
		Code:       "ERROR_CUSTOMER_COUNT_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgCustomerCountFailed,
		Type:       "error",
	}

	ErrorCustomerConvertIDFailed = ResponseCode{
		Code:       "ERROR_CUSTOMER_CONVERT_ID_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgCustomerConvertIDFailed,
		Type:       "error",
	}

	ErrorCustomerNotFound = ResponseCode{
		Code:       "ERROR_CUSTOMER_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgCustomerNotFound,
		Type:       "error",
	}

	ErrorCustomerFetchFailed = ResponseCode{
		Code:       "ERROR_CUSTOMER_FETCH_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgCustomerFetchFailed,
		Type:       "error",
	}

	// Account Validation related error response codes
	ErrorValidationRuleGetFailed = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_GET_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgValidationRuleGetFailed,
		Type:       "error",
	}

	ErrorValidationRuleMinMaxLengthMismatch = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_MIN_MAX_LENGTH_MISMATCH",
		StatusCode: StatusBadRequest,
		Message:    MsgValidationRuleMinMaxLengthMismatch,
		Type:       "error",
	}

	ErrorValidationRuleConvertIDFailed = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_CONVERT_ID_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgValidationRuleConvertIDFailed,
		Type:       "error",
	}

	ErrorValidationRuleUpdateFailed = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgValidationRuleUpdateFailed,
		Type:       "error",
	}

	ErrorValidationRulePendingActionFetchFailed = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_PENDING_ACTION_FETCH_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgValidationRulePendingActionFetchFailed,
		Type:       "error",
	}

	ErrorValidationRuleActionsFetchFailed = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_ACTIONS_FETCH_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgValidationRuleActionsFetchFailed,
		Type:       "error",
	}

	// Donation related error response codes
	ErrorDonationCategoryLookupFailed = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_LOOKUP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationCategoryLookupFailed,
		Type:       "error",
	}

	ErrorDonationCompanyLookupFailed = ResponseCode{
		Code:       "ERROR_DONATION_COMPANY_LOOKUP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationCompanyLookupFailed,
		Type:       "error",
	}

	ErrorDonationLookupFailed = ResponseCode{
		Code:       "ERROR_DONATION_LOOKUP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationLookupFailed,
		Type:       "error",
	}

	ErrorDonationParseStartDateFailed = ResponseCode{
		Code:       "ERROR_DONATION_PARSE_START_DATE_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationParseStartDateFailed,
		Type:       "error",
	}

	ErrorDonationParseEndDateFailed = ResponseCode{
		Code:       "ERROR_DONATION_PARSE_END_DATE_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationParseEndDateFailed,
		Type:       "error",
	}

	ErrorDonationFetchCompanyFailed = ResponseCode{
		Code:       "ERROR_DONATION_FETCH_COMPANY_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationFetchCompanyFailed,
		Type:       "error",
	}

	ErrorDonationFetchCategoryFailed = ResponseCode{
		Code:       "ERROR_DONATION_FETCH_CATEGORY_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationFetchCategoryFailed,
		Type:       "error",
	}

	ErrorDonationUpdateFailed = ResponseCode{
		Code:       "ERROR_DONATION_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationUpdateFailed,
		Type:       "error",
	}

	ErrorDonationImageUpdateFailed = ResponseCode{
		Code:       "ERROR_DONATION_IMAGE_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationImageUpdateFailed,
		Type:       "error",
	}

	ErrorDonationImageDeleteFailed = ResponseCode{
		Code:       "ERROR_DONATION_IMAGE_DELETE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationImageDeleteFailed,
		Type:       "error",
	}

	ErrorDonationImageAddFailed = ResponseCode{
		Code:       "ERROR_DONATION_IMAGE_ADD_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationImageAddFailed,
		Type:       "error",
	}

	// Bank related error response codes
	ErrorBankFileParseFailed = ResponseCode{
		Code:       "ERROR_BANK_FILE_PARSE_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgBankFileParseFailed,
		Type:       "error",
	}

	ErrorBankRejectionPayloadDecodeFailed = ResponseCode{
		Code:       "ERROR_BANK_REJECTION_PAYLOAD_DECODE_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgBankRejectionPayloadDecodeFailed,
		Type:       "error",
	}

	ErrorBankLogoUpdateFailed = ResponseCode{
		Code:       "ERROR_BANK_LOGO_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgBankLogoUpdateFailed,
		Type:       "error",
	}

	// CPS Action related error response codes
	ErrorCPSActionRejectionPayloadDecodeFailed = ResponseCode{
		Code:       "ERROR_CPS_ACTION_REJECTION_PAYLOAD_DECODE_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgCPSActionRejectionPayloadDecodeFailed,
		Type:       "error",
	}

	// Permission related error response codes
	ErrorPermissionGroupRequestCreationFailed = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_REQUEST_CREATION_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgPermissionGroupRequestCreationFailed,
		Type:       "error",
	}

	ErrorPermissionGroupRequestUpdateFailed = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_REQUEST_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgPermissionGroupRequestUpdateFailed,
		Type:       "error",
	}

	// Department related error response codes
	ErrorDepartmentCreationSuccess = ResponseCode{
		Code:       "ERROR_DEPARTMENT_CREATION_SUCCESS",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentCreationSuccess,
		Type:       "error",
	}

	ErrorDepartmentUpdateCPSActionCreated = ResponseCode{
		Code:       "ERROR_DEPARTMENT_UPDATE_CPS_ACTION_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentUpdateCPSActionCreated,
		Type:       "error",
	}

	// Budget Category related error response codes
	ErrorBudgetCategoryUpdateSuccess = ResponseCode{
		Code:       "ERROR_BUDGET_CATEGORY_UPDATE_SUCCESS",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoryUpdateSuccess,
		Type:       "error",
	}

	// Donation Service related error response codes
	ErrorDonationCategoryIconUploaded = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_ICON_UPLOADED",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryIconUploadSuccess,
		Type:       "error",
	}

	ErrorDonationCategoryIconUpdated = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_ICON_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryIconUpdateSuccess,
		Type:       "error",
	}

	ErrorDonationCPSActionApproved = ResponseCode{
		Code:       "ERROR_DONATION_CPS_ACTION_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgDonationCPSActionApprovedSuccess,
		Type:       "error",
	}

	ErrorDonationCompanyLogoUploaded = ResponseCode{
		Code:       "ERROR_DONATION_COMPANY_LOGO_UPLOADED",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyLogoUploadSuccess,
		Type:       "error",
	}

	ErrorDonationCompanyLogoUpdated = ResponseCode{
		Code:       "ERROR_DONATION_COMPANY_LOGO_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyLogoUpdateSuccess,
		Type:       "error",
	}

	ErrorDonationImageUploaded = ResponseCode{
		Code:       "ERROR_DONATION_IMAGE_UPLOADED",
		StatusCode: StatusOK,
		Message:    MsgDonationImageUploadSuccess,
		Type:       "error",
	}

	ErrorDonationImagesUpdated = ResponseCode{
		Code:       "ERROR_DONATION_IMAGES_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgDonationImagesUpdateSuccess,
		Type:       "error",
	}

	// Notification Service related error response codes
	ErrorNotificationConstructed = ResponseCode{
		Code:       "ERROR_NOTIFICATION_CONSTRUCTED",
		StatusCode: StatusOK,
		Message:    MsgNotificationConstructedSuccess,
		Type:       "error",
	}

	ErrorNotificationUpdated = ResponseCode{
		Code:       "ERROR_NOTIFICATION_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgNotificationUpdateSuccess,
		Type:       "error",
	}

	// Permission Service related error response codes
	ErrorPermissionGroupRequestCreated = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgPermissionGroupRequestCreatedSuccess,
		Type:       "error",
	}

	// Mini App Service related error response codes
	ErrorMiniAppDetailsFetched = ResponseCode{
		Code:       "ERROR_MINI_APP_DETAILS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppDetailsFetchedSuccess,
		Type:       "error",
	}

	ErrorMiniAppCredentialGenerated = ResponseCode{
		Code:       "ERROR_MINI_APP_CREDENTIAL_GENERATED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppCredentialGeneratedSuccess,
		Type:       "error",
	}

	ErrorMiniAppFileUploaded = ResponseCode{
		Code:       "ERROR_MINI_APP_FILE_UPLOADED",
		StatusCode: StatusOK,
		Message:    MsgMiniAppFileUploadedSuccess,
		Type:       "error",
	}

	// Feedback Kafka related error response codes
	ErrorFeedbackSavedToDatabase = ResponseCode{
		Code:       "ERROR_FEEDBACK_SAVED_TO_DATABASE",
		StatusCode: StatusOK,
		Message:    MsgFeedbackSavedToDatabaseSuccess,
		Type:       "error",
	}

	// HQ Service related error response codes
	ErrorHQApproved = ResponseCode{
		Code:       "ERROR_HQ_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgHQApprovedSuccess,
		Type:       "error",
	}

	// CPS Action related error response codes
	ErrorCPSActionStatusInvalid = ResponseCode{
		Code:       "ERROR_CPS_ACTION_STATUS_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgCPSActionStatusInvalid,
		Type:       "error",
	}

	// User Unlink error response codes
	ErrorUnlinkFaild = ResponseCode{
		Code:       "ERROR_USER_UNLINK_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgUserUnlinkFailed,
		Type:       "error",
	}
	// Ad Service related error response codes
	ErrorAdvertCreated = ResponseCode{
		Code:       "ERROR_ADVERT_CREATED",
		StatusCode: StatusOK,
		Message:    MsgAdvertCreateError,
		Type:       "error",
	}

	ErrorAdvertUpdate = ResponseCode{
		Code:       "ERROR_ADVERT_UPDATE",
		StatusCode: StatusOK,
		Message:    MsgAdvertUpdateError,
		Type:       "error",
	}

	ErrorAdvertAlreadyEnabled = ResponseCode{
		Code:       "ERROR_ADVERT_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgAdvertAlreadyEnabled,
		Type:       "error",
	}
	ErrorAdvertAlreadyDisabled = ResponseCode{
		Code:       "ERROR_ADVERT_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgAdvertAlreadyDisabled,
		Type:       "error",
	}

	// Account Validation Service related error response codes
	ErrorValidationRuleApproved = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgValidationRuleApprovedSuccess,
		Type:       "error",
	}

	ErrorInvalidInputParameters = ResponseCode{
		Code:       "ERROR_INVALID_INPUT_PARAMETERS",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidInputParameters,
		Type:       "error",
	}

	ErrorMissingOrInvalidImage = ResponseCode{
		Code:       "ERROR_MISSING_OR_INVALID_IMAGE",
		StatusCode: StatusBadRequest,
		Message:    MsgMissingOrInvalidImage,
		Type:       "error",
	}
)
