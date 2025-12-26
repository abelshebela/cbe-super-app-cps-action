package localization

import "time"

// ResponseCode represents a standardized response code with status and message
type ResponseCode struct {
	Code       string    `json:"code"`
	StatusCode int       `json:"status_code"`
	TimeStamp  time.Time `json:"timestamp"`
	Message    string    `json:"message"`
	Type       string    `json:"type"` // "success", "error", "warning", "info"
}

// Error implements error.
func (r ResponseCode) Error() string {
	// panic("unimplemented")
	return r.Message
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
	SuccessDonationCompanyFetched,
	SuccessDonationCategoryEnableRequestSent,
	SuccessDonationCategoryDisableRequestSent,
	SuccessDonationCompanyEnableRequestSent,
	SuccessDonationCompanyDisableRequestSent,
	SuccessAccountInfoFetched,

	SuccessDonationImageUploaded,
	SuccessDonationImagesUpdated,
	SuccessNotificationConstructed,
	SuccessNotificationUpdated,
	SuccessFeedbackSavedToDatabase,
	SuccessAvatarCreated,
	SuccessDeleteBanksRequest,
	SuccessDeleteRequestCreated,
	SuccessGetOneBank,
	SuccessBankDisableRequestCreated,
	SuccessBankEnableRequestCreated,
	SuccessBankUpdatedRequestSent,
	SuccessWalletEnableRequestSubmitted,
	SuccessWalletDisableRequestSubmitted,
	SuccessWalletCreationRequestSent,
	SuccessWalletUpdateRequestSent,
	SuccessWalletDeleted,
	SuccessWalletsRetrieved,
	SuccessWalletRetrieved,

	SuccessEcommerceMerchantCreatedSuccessfully,
	SuccessEcommerceMerchantUpdatedSuccessfully,
	SuccessEcommerceMerchantDeletedSuccessfully,
	SuccessEcommerceMerchantEnableSuccessfully,
	SuccessEcommerceMerchantDisableSuccessfully,
	SuccessEcommerceMerchantFetchedSuccessfully,
	SuccessEcommerceMerchantsFetchedSuccessfully,
	SuccessEcommerceMerchantLookup,

	SuccessTopupEnableRequestSubmitted,
	SuccessTopupDisableRequestSubmitted,
	SuccessTopupCreationRequestSent,
	SuccessTopupUpdateRequestSent,
	SuccessTopupDeleted,
	SuccessTopupsRetrieved,
	SuccessTopupRetrieved,

	//BPS_USer
	SuccessBpsUserEnableRequestSent,
	SuccessBpsUserDisableRequestSent,

	// CPS Roles
	SuccessCPSRoleCreated,
	SuccessCPSRoleUpdated,
	SuccessCPSRolesFetched,
	SuccessCPSRoleFetched,
	SuccessCPSRoleEnabled,
	SuccessCPSRoleDisabled,

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
	SuccessAvatarEnabled,
	SuccessAvatarDisabled,
	SuccessAvatarDeleted,

	SuccessBranchRetrieved,
	SuccessBranchesRetrieved,
	SuccessBranchesEnabled,
	SuccessBranchesDisabled,
	SuccessEnableBranchesRequestSent,
	SuccessDisableBranchesRequestSent,

	SuccessRegionRetrieved,
	SuccessRegionsRetrieved,
	SuccessRegionsEnabled,
	SuccessRegionsDisabled,
	SuccessEnableRegionsRequestSent,
	SuccessDisableRegionsRequestSent,

	SuccessDistrictRetrieved,
	SuccessDistrictsRetrieved,
	SuccessDistrictsEnabled,
	SuccessDistrictsDisabled,
	SuccessEnableDistrictsRequestSent,
	SuccessDisableDistrictsRequestSent,

	SuccessCityRetrieved,
	SuccessCitiesRetrieved,
	SuccessCitiesEnabled,
	SuccessCitiesDisabled,
	SuccessEnableCitiesRequestSent,
	SuccessDisableCitiesRequestSent,

	// department related success response
	SuccessGetAllDepartments,
	SuccessDepartmentCreateRequestCreated,
	SuccessDepartmentUpdateRequestCreated,
	SuccessDepartmentDisableRequestCreated,
	SuccessDepartmentEnableRequestCreated,
	SuccessDepartmentUpdateRequestCreated,
	SuccessDepartmentCreateRequestCreated,
	ErrorDepartmentNotFound,

	//hq related success response codes
	SuccessHQArchiveTimeFetched,
	SuccessHQBlockTimeFetched,
	SuccessHQPasswordExpiryFetched,
	SuccessHQBlockTimeUpdateRequestSubmitted,
	SuccessHQArchiveTimeUpdateRequestSubmitted,
	SuccessHQPasswordExpiryUpdateRequestSubmitted,
	//mini app related success response codes
	SuccessMiniAppDisableRequestSubmitted,
	SuccessMiniAppEnableRequestSubmitted,
	SuccessMiniAppMerchantCreateRequestCreated,
	SuccessUpdateMiniAppRequestCreated,
	SuccessMiniAppDeletedRequestCreated,
	SuccessMiniAppsRetrieved,
	SuccessMiniAppFetchedByID,
	SuccessMiniAppActionCompleted,
	ErrorMiniAppMerchantNotFound,
	ErrorMiniAppMerchantDisabled,

	// Sitota Related success
	SuccessAllSitotasRetrieved,
	SuccessSitotaRetrieved,

	// Encryption
	SuccessEncryptionGenerated,

	// Action Role success response codes
	SuccessActionRolesFetched,
	SuccessActionRoleFetched,
	SuccessActionRoleCreateRequestCreated,
	SuccessActionRoleUpdateRequestCreated,
	SuccessActionRoleEnableRequestCreated,
	SuccessActionRoleDisableRequestCreated,

	// event merchant success response codes
	SuccessEventMerchantCreated,
	SuccessEventMerchantDisabled,
	SuccessEventMerchantEnabled,
	SuccessEventMerchantUpdated,
	SuccessEventMerchantDeleted,
	SuccessEventMerchantFetched,

	// Error codes
	ErrorDeviceVersionAlreadyExists,
	ErrorDeviceVersionAlreadyEnabled,
	ErrorDeviceVersionAlreadyDisabled,
	ErrorDeviceVersionNotFound,
	ErrorDuplicateProductName,
	ErrorDeviceVersionUpdateFailed,
	ErrorDeviceVersionDeleteFailed,
	ErrorDeviceVersionEnableFailed,
	ErrorDeviceVersionDisableFailed,
	ErrorInvalidKey,
	ErrorInvalidEncData,
	ErrorInvalidPadding,
	ErrorUserNotFound,
	ErrorUserAlreadyExists,
	ErrorPermissionGroupAlreadyExists,
	ErrorPermissionCatagoryNotFound,
	ErrorUserUnauthorized,
	ErrorThirdAPIRequestNotUnauthorized,
	ErrorUserForbidden,
	ErrorUserInvalidCredentials,
	ErrorUserAccountBlocked,
	ErrorUserSessionExpired,
	ErrorProductCodeDatabaseQueryFailed,
	ErrorProductCodeCountFailed,
	ErrorProductCodeDatabaseUpdateFailed,
	ErrorProductCodeNotFound,
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
	ErrorNotificationAlreadyExists,
	ErrorFeedbackSavedToDatabase,
	ErrorHQApproved,
	ErrorUnlinkFaild,
	ErrorCPSActionStatusInvalid,
	ErrorAdvertCreated,
	ErrorAdvertUpdate,
	ErrorAdvertUpdate,
	ErrorAvatarAlreadyEnabled,
	ErrorAvatarAlreadyDisabled,
	ErrorAdvertAlreadyEnabled,
	ErrorTitleLength3To20,
	ErrorDescriptionLength30To100,
	ErrorAdvertAlreadyDisabled,
	ErrorAdvertTitleAlreadyExists,
	ErrorBudgetCategoryNameAlreadyExists,
	ErrorAdvertTitleNotChanged,
	ErrorValidationRuleApproved,
	ErrorAccountNumberRequired,

	ErrorActionNotFound,
	ErrorPendingCpsActionExists,
	ErrorUserAlreadyEnabled,
	ErrorUserAlreadyDisabled,
	ErrorBankImageMissingOrInvalid,
	ErrorBankUpdateFailed,
	ErrorValidationFailed,
	ErrorRequiredFieldMissing,
	ErrorBankDeleteRequestFailed,
	ErrorBankImageMissingOrInvalid,
	ErrorGetAllBanksFailed,
	ErrorGetAllBanksFailed,
	ErrorGetOneBank,
	ErrorBankDisableRequest,
	ErrorBankEnableRequestFailed,
	ErrorUserCodeRequired,

	ErrorEventNameRequired,
	ErrorEventAlreadyExists,
	ErrorCoverImageRequired,
	ErrorInvalidAmounts,
	ErrorAuthTierAlreadyExists,
	ErrorInvalidMethod,
	ErrorMerchantNotFound,
	ErrorEventNotFound,
	ErrorEventAlreadyEnabled,
	ErrorEventAlreadyDisabled,
	ErrorUnhandledServer,
	ErrorInvalidMerchantID,
	ErrorInvalidDateFormat,
	ErrorInvalidFileUpload,
	ErrorInvalidNumberFormat,
	ErrorUpdateEventEmptyPayload,
	ErrorMerchantIDRequired,
	ErrorEventVenueRequired,
	ErrorCustomerIDRequired,
	ErrorStartDateRequired,
	ErrorDueDateRequired,
	ErrorTotalTicketCountRequired,
	ErrorInvalidTicketCount,
	ErrorEventCityRequired,
	ErrorEventDescriptionRequired,
	ErrorTicketsRequired,
	ErrorBankAlreadyEnabled,
	ErrorAmountTierNotFound,
	ErrorBankAlreadyDisabled,
	ErrorHealthCheck,
	ErrorActionAlreadyExists,
	ErrorNoDataProvidedForCreate,
	ErrorNoDataProvidedForUpdate,
	ErrorNoDataProvidedForBankUpdate,
	ErrorBankWithCodeAlreadyExists,
	ErrorBankWithBICAlreadyExists,
	ErrorBankWithNameAlreadyExists,
	ErrorSessionRetrievalFailed,
	ErrorInvalidToken,
	ErrorFileInvalidType,
	ErrorFileUploadFailed,
	ErrorMiniAppNameAlreadyExists,

	ErrorFileParseFailed,
	ErrorResourceNotFound,
	ErrorActionNameAlreadyExists,
	ErrorOnDisablingExistingDeviceControl,
	ErrorInvalidInputParameter,
	ErrorActionNameIsRequired,
	ErrorInvalidInputParameters,
	ErrorMissingOrInvalidImage,
	ErrorPendingCpsActionExists,
	ErrorUnexpectedError,
	ErrorExternalServiceError,
	ErrorFileNotFound,
	ErrorInvalidID,
	ErrorInvalidBooleanFormat,
	ErrorInvalidJSONPayload,
	ErrorInvalidAction,
	ErrorIncompleteUserInfo,
	ErrorPendingActionExists,
	ErrorDuplicateColorExists,
	ErrorInvalidActionFormat,
	ErrorMissingFile,
	ErrorBucketNotFound,
	ErrorFailedToBucket,

	// Add more as needed...
	//wallet related error codes
	ErrorWalletNameRequired,
	ErrorWalletCodeRequired,
	ErrorWalletTypeRequired,
	ErrorWalletAvatarRequired,
	ErrorWalletAvatarInvalid,
	ErrorWalletAvatarTooLarge,
	ErrorWalletAvatarInvalidType,
	ErrorWalletRechangeOption,
	ErrorAvatarAlreadyExist,
	ErrorWalletNameAlreadyExists,
	ErrorWalletCodeAlreadyExists,
	ErrorWalletAlreadyDisabled,
	ErrorWalletAlreadyEnabled,
	ErrorWalletIDRequired,
	ErrorWalletNotFound,
	ErrorWalletUpdateEmptyPayload,
	ErrorInvalidWalletCode,
	ErrorInvalidWalletType,
	ErrorInvalidWalletName,

	//topup related error codes
	ErrorTopupNameRequired,
	ErrorTopupCodeRequired,
	ErrorTopupAvatarRequired,
	ErrorTopupAvatarInvalid,
	ErrorTopupAvatarTooLarge,
	ErrorTopupAvatarInvalidType,
	ErrorTopupServiceOption,
	ErrorAvatarAlreadyExist,
	ErrorTopupNameAlreadyExists,
	ErrorTopupCodeAlreadyExists,
	ErrorTopupAlreadyDisabled,
	ErrorTopupAlreadyEnabled,
	ErrorTopupIDRequired,
	ErrorTopupNotFound,
	ErrorTopupUpdateEmptyPayload,

	ErrorInvalidRequestBody,
	ErrorInvalidRequest,
	ErrorBulkServiceAlreadyDisabled,
	ErrorInvalidPaginationParams,
	ErrorInvalidDistrict,
	ErrorInvalidRegion,
	ErrorInvalidDistrictOrRegionCodeLength,

	ErrorBranchCodeRequired,
	ErrorDistrictCodeRequired,
	ErrorRegionCodeRequired,
	ErrorCityCodeRequired,

	ErrorBranchNotFound,
	ErrorDistrictNotFound,
	ErrorCannotEnableDistrict,
	ErrorCannotEnableBranch,
	ErrorRegionNotFound,
	ErrorCityNotFound,

	ErrorDuplicateAction,
	ErrorOneOrMoreInvalidCodes,

	ErrorAlreadyEnabled,
	ErrorAlreadyDisabled,
	ErrorInvalidBulkServiceKey,
	ErrorInvalidRequiredAction,
	ErrorFailToUpdateParent,
	ErrorFailToUpdateChild,

	// department related error
	ErrorDepartmentCreateRequest,
	ErrorInvalidFormatForDepartmentPermissionGroups,
	ErrorInvalidFormatForDepartmentPortalCards,
	ErrorInvalidDepartmentPortalCard,
	ErrorInvalidDepartmentPermissionGroup,
	ErrorDepartmentAlreadyDisabled,
	ErrorDepartmentAlreadyEnabled,
	ErrorDepartmentInvalidID,
	ErrorDepartmentWithNameAlreadyExists,
	// Fayda
	ErrorFaydaUserAccountEnabled,
	ErrorFaydaUserAccountDisabled,
	ErrorNotFaydaUser,

	//hq related error response codes
	ErrorHQNotFound,
	ErrorInvalidHQRequest,
	ErrorCPSActionFailed,
	ErrorFailedToParseJson,

	//miniapp related errors
	ErrorMiniAppNotFound,
	ErrorMiniAppAlreadyExists,
	ErrorMiniAppAlreadyEnabled,
	ErrorMiniAppAlreadyDisabled,
	ErrorMiniAppAlreadyDeleted,
	ErrorMiniAppIDRequired,
	ErrorMiniAppMerchantIDRequired,
	ErrorMiniAppNameRequired,
	ErrorMiniAppDescriptionRequired,
	ErrorAppIconRequired,
	ErrorInvalidAppNameFormat,
	ErrorAppViewTypeRequired,
	ErrorBannerImageRequired,
	ErrorMiniAppURLRequired,
	ErrorMiniAppMerchantIDRequired,
	ErrorMiniAppNameRequired,
	ErrorInvalidBooleanFormat,
	ErrorInvalidImageFormat,
	ErrorFileTooLarge,
	ErrorInvalidURL,
	ErrorIncompleteBranchProductCodes,
	ErrorNoProductCodesProvided,
	ErrorBothProductCodesRequired,
	ErrorAppViewTypeInvalidOrMissing,
	ErrorInvalidAppViewType,
	ErrorExclusiveAppFlags,
	ErrorUpdateMiniAppEmptyPayload,
	//customer  and bulk relatedcode
	UserNotFoundWithGivenID,
	ErrorFailedToGetCustomerDetail,
	ErrorKeyRequiredForBulkService,
	SuccessCustomerDetailSuccessfullyFetched,
	CustomerEnableRequestSessionCreatedSuccessfully,
	CustomerDisableRequestCreatedSuccessfully,
	CustomerEnableRequestCreatedSuccessfully,
	ErrorCustomerAlreadyDisabled,
	ErrorCustomerAlreadyEnabled,
	ErrorIdNotSetOnQueryParam,
	CustomerDetailSuccessfullyFetched,
	ErrorFailedToGetBlockedCustomer,
	SuccessFullyFetchBlockCustomer,
	UnableToFetchBulkService,
	BulkServiceFetchSuccessfully,
	BulkServiceEnableRequestSuccess,
	BulkServiceDisableRequestSuccess,
	ErrorAvatarNotExist,
	ErrorBulkServiceAlreadyEnabled,
	ErrorDuplicateCBEIFBProductCode,
	ErrorCustomerAccountNumberMustContainOnlyNumbers,
	ErrorCustomerCIFMustContainOnlyNumbers,
	//service details
	ErrorSingleMaxTransferCannotBeLessOrEqualToMinAmount,
	ErrorTotalMaxTransferCannotBeLessExistTransfers,
	ErrorMinAmountCanNotBeGreaterThanCap,
	ErrorPermissionCategoryNotFound,
	ErrorPermissionGroupNotFound,
	ErrorPermissionGroupRequired,
	ErrorCityAlreadyDisabled,
	ErrorCityAlreadyEnabled,
	ErrorRegionAlreadyDisabled,
	ErrorRegionAlreadyEnabled,
	ErrorDistrictAlreadyDisabled,
	ErrorDistrictAlreadyEnabled,
	ErrorBranchAlreadyDisabled,
	ErrorBranchAlreadyEnabled,
	ErrorSingleTransferCanNotBeGreaterThanCap,
	ErrorMinAmountCanNotBeGreaterThanTotal,
	ErrorNoChangesDetected,
	ErrorDonationCategoryLookupFailed,
	ErrorNoChangesToUpdate,
	ErrorBudgetCategoryAlreadyEnabled,
	ErrorBudgetCategoryAlreadyDisabled,
	ErrorDonationParseStartDateFailed,
	ErrorDonationParseEndDateFailed,
	ErrorDonationFetchCompanyFailed,
	ErrorDonationFetchCategoryFailed,
	ErrorDonationUpdateFailed,
	ErrorDonationImageUpdateFailed,
	ErrorDonationImageDeleteFailed,
	ErrorDonationImageAddFailed,
	ErrorBankFileParseFailed,
	ErrorBankRejectionPayloadDecodeFailed,
	ErrorBankLogoUpdateFailed,
	ErrorBankUpdateFailed,
	ErrorDonationCategoryNameDuplicated,
	ErrorDonationCompanyLookupFailed,
	ErrorDonationCategoryIDRequired,
	ErrorDuplicateCBEProductCode,

	// Donation Company Error Codes
	ErrorCompanyNameAlreadyExists,
	ErrorCompanyCodeAlreadyExists,
	ErrorAccountNumberAlreadyExists,
	ErrorEmailAlreadyExist,
	ErrorPhonenumberAlreadyExist,
	ErrorCodeAlreadyExist,
	ErrorLogoIsRequired,
	ErrorAccountNumberValidationFailed,
	ErrorAccountNumberNotActive,
	ErrorAccountNumberNotFound,
	ErrorDonationCompanyIdRequired,
	ErrorDonationAlreadyEnabled,
	ErrorDonationAlreadyDisabled,
	SuccessDonationCompanyUpdated,
	ErrorDonationTitleDuplicated,
	ErrorDonationImageUploaded,
	ErrorDonationImagesUpdated,
	ErrorImageRequired,
	ErrorDonationCompanyNotFound,
	ErrorCompanyIsNotEnabled,
	ErrorCategoryIsNotEnabled,
	ErrorDonationCategoryNotFound,
	ErrorDonationLookupFailed,
	ErrorMiniAppMerchantEnableFailed,
	ErrorMiniAppMerchantDisableFailed,
	ErrorMiniAppMerchantDeleteFailed,
	ErrorMiniAppMerchantUpdateFailed,
	ErrorExistEmail,
	ErrorInvalidPhoneNumber,
	ErrorExistPhoneNumber,
	ErrorInvalidEmail,
	ErrorCPSRoleAlreadyEnabled,
	ErrorCPSRoleAlreadyDisabled,

	// OTP related error codes
	ErrorOTPExpired,
	ErrorUnsupportedAction,
	ErrorOTPInvalid,
	ErrorOTPAlreadyExists,
	ErrorOTPTooManyAttempts,
	ErrorOTPSendFailed,

	ErrorNewsCategoryInvalidID,
	ErrorNewsTagInvalidID,
	ErrorNewsTagWithNameAlreadyExists,
	ErrorNewsCategoryWithNameAlreadyExists,

	// Sitota Related errors
	ErrorSitotaRequired,

	// Encryption
	ErrConfigIsEmpty,
	ErrMarshalingData,
	ErrInvalidKeyOrIv,

	// Vault related
	SuccessVaultGroupCategoryCreationRequestSubmitted,
	SuccessVaultGroupCategoriesRetrieved,
	SuccessVaultGroupCategoryRetrieved,
	SuccessVaultGroupCategoryUpdateRequestSubmitted,
	SuccessVaultGroupCategoryDeleteRequestSubmitted,
	SuccessVaultGroupCategoryEnableRequestSubmitted,
	SuccessVaultGroupCategoryDisableRequestSubmitted,
	SuccessVaultAmountTierCreationRequestSubmitted,
	SuccessVaultAmountTierFetchedSuccessfully,
	SuccessVaultAmountTierUpdateRequestSubmitted,
	SuccessVaultAmountTierDeleteRequestSubmitted,
	SuccessVaultAmountTierDisableRequestSubmitted,
	ErrorFailedToBeingTransaction,
	ErrorDuplicateBankProduct,
	ErrorVaultGroupCategoryNotFound,
	ErrorGroupVaultNotFound,
	ErrorCannotDeleteActiveVaultGroupCategory,
	ErrorNoBankProductFound,
	ErrorCannotDeletedBankProduct,
	ErrorCannotEnableOrDisable,
	ErrorBankVaultProductAlreadyDeleted,
	ErrorDuplicateGroupVaultCategory,
	ErrorBankVaultAlreadyEnabled,
	ErrorBankVaultAlreadyDisabled,
	ErrorVaultGroupAlreadyEnabled,
	ErrorVaultGroupAlreadyDisabled,
	ErrorBankAlreadyEnabled,
	ErrorBankAlreadyDisabled,
	ErrorBankWithNameAlreadyExists,
	ErrorBankWithBICAlreadyExists,
	ErrorBankWithCodeAlreadyExists,

	// BPS Action Role related error codes
	ErrorBpsActionRoleNotFound,

	// transaction related responses
	SuccessTransactionRetrieved,
	ErrorTransactionIDRequired,
	ErrorTransactionIdentifierRequired,

	// bank related errors
	ErrorInvalidAccountNumberFormat,
	ErrorInvalidFormatForBIC,
	ErrorInvalidFormatForType,
	ErrorInvalidFormatForCode,
	ErrorInvalidFormatForName,

	// event mercahnt error
	ErrorEventMerchantInvalidMerchantID,
	ErrorEventMerchantInvalidMerchantType,
	ErrorEventMerchantInvalidSettlementMethod,
	ErrorEventMerchantInvalidMerchantName,
	ErrorEventMerchantInvalidBankAccountNumber,
	ErrorEventMerchantInvalidEmail,
	ErrorEventMerchantInvalidPhoneNumber,

	// Access List Segmentation Success Codes
	SuccessAccessListSegmentationCreated,
	SuccessAccessListSegmentationUpdated,
	SuccessAccessListSegmentationEnabled,
	SuccessAccessListSegmentationDisabled,
	SuccessAccessListSegmentationsRetrieved,
	SuccessAccessListSegmentationRetrieved,

	// Access List Segmentation Error Codes
	ErrorAccessListSegmentationNotFound,
	ErrorAccessListSegmentationAlreadyEnabled,
	ErrorAccessListSegmentationAlreadyDisabled,
	ErrorAccessListSegmentationInvalidID,
	ErrorAccessListSegmentationInvalidID,
	ErrorAccessListSegmentationAlreadyEnabled,
	ErrorAccessListSegmentationNotFound,
	UnableToCreateAccessListSegmentation,
	ErrorAccessListSegmentationInvalidID,
	ErrorAccessListSegmentationIDSRequired,
	ErrorAccessListSegmentationNameAlreadyExists,
	ErrorCustomerSegmentationCodeNotFound,
	ErrorServiceIdRequired,

	// Access List Segmentaion Success Code
	SuccessAccessListSegmentationRetrieved,
	SuccessAccessListSegmentationDisabled,
	SuccessAccessListSegmentationEnabled,
	SuccessAccessListSegmentationUpdated,
	SuccessAccessListSegmentationCreated,
	AccessListSegmentationCreatedSuccessfully,

	// Customer segmentations
	CustomerSegmentationCreationSubmittedSuccessfully,
	CustomerSegmentationUpdateSubmittedSuccessfully,
	CustomerSegmentationFetchedSuccessfully,
	CustomerSegmentationDeleteddSuccessfully,
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

	SuccessAvatarUpdated = ResponseCode{
		Code:       "SUCCESS_AVATAR_UPDATED",
		StatusCode: StatusCreated,
		Message:    MsgAvatarUpdateSuccessfully,
		Type:       "success",
	}

	SuccessAvatarRetrieved = ResponseCode{
		Code:       "SUCCESS_AVATAR_RETRIVED",
		StatusCode: StatusCreated,
		Message:    MsgAvatarRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessAvatarEnabled = ResponseCode{
		Code:       "SUCCESS_AVATAR_ENABLED",
		StatusCode: StatusAccepted,
		Message:    MsgAvatarEnabledSuccessfully,
		Type:       "success",
	}

	SuccessAvatarDisabled = ResponseCode{
		Code:       "SUCCESS_AVATAR_DISABLED",
		StatusCode: StatusAccepted,
		Message:    MsgAvatarDisabledSuccessfully,
		Type:       "success",
	}

	SuccessAvatarDeleted = ResponseCode{
		Code:       "SUCCESS_AVATAR_DELETED",
		StatusCode: StatusNoContent,
		Message:    MsgAvatarDeletedSuccessfully,
		Type:       "success",
	}

	SuccessUserUpdated = ResponseCode{
		Code:       "SUCCESS_USER_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgUserUpdatedSuccessfully,
		Type:       "success",
	}

	SuccessAmountBasedAuthRequestSent = ResponseCode{
		Code:       "SUCCESS_AMOUNT_BASED_AUTH_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgAmountBasedSuccessfullySent,
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

	SuccessRoleCreatedRequestSent = ResponseCode{
		Code:       "SUCCESS_Role_CREATED_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgRoleCreatedRequestSent,
		Type:       "success",
	}

	SuccessRoleUpdatedRequestSent = ResponseCode{
		Code:       "SUCCESS_Role_UPDATED_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgRoleUpdatedRequestSent,
		Type:       "success",
	}

	SuccessJobRoleCreatedRequestSent = ResponseCode{
		Code:       "SUCCESS_JOB_ROLE_CREATED_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgJobRoleCreatedRequestSent,
		Type:       "success",
	}
	SuccessJobRoleUpdatedRequestSent = ResponseCode{
		Code:       "SUCCESS_JOB_ROLE_UPDATED_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgJobRoleUpdateRequestSent,
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
	SuccessDonationCategoryEnableRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORY_ENABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryEnableRequestSent,
		Type:       "success",
	}

	SuccessDonationCategoryDisableRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_CATEGORY_DISABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationCategoryDisableRequestSent,
		Type:       "success",
	}
	SuccessDonationCompanyEnableRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_ENABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyEnableRequestSent,
		Type:       "success",
	}

	SuccessDonationCompanyDisableRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_DISABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyDisableRequestSent,
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
	SuccessAccountInfoFetched = ResponseCode{
		Code:       "SUCCESS_ACCOUNT_INFO_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgAccountInfoFetched,
		Type:       "success",
	}

	SuccessDonationCompanyUpdatedRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_UPDATED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationCompanyUpdatedRequestSent,
		Type:       "success",
	}

	SuccessDonationCompanyUpdated = ResponseCode{
		Code:       "SUCCESS_DONATION_COMPANY_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgSuccessDonationCompanyUpdated,
		Type:       "success",
	}

	SuccessKYCApproved = ResponseCode{
		Code:       "SUCCESS_KYC_APPROVE_REQUESTED",
		StatusCode: StatusOK,
		Message:    MsgKYCApproved,
		Type:       "success",
	}

	SuccessKYCUpdatedRequestSent = ResponseCode{
		Code:       "SUCCESS_KYC_UPDATED_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgKYCUpdatedRequestSent,
		Type:       "success",
	}

	SuccessKYCFetched = ResponseCode{
		Code:       "SUCCESS_KYC_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgKYCFetched,
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

	SuccessDonationEnableRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_ENABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationEnableRequestSent,
		Type:       "success",
	}

	SuccessDonationDisableRequestSent = ResponseCode{
		Code:       "SUCCESS_DONATION_DISABLE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDonationDisableRequestSent,
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

	SuccessCPSActionReversed = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_REVERSED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionReversedSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionRejected = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_REJECTED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionRejectedSuccessfully,
		Type:       "success",
	}
	SuccessCPSActionCanceled = ResponseCode{
		Code:       "SUCCESS_CPS_ACTION_CANCELED",
		StatusCode: StatusOK,
		Message:    MsgCPSActionCanceledSuccessfully,
		Type:       "success",
	}

	SuccessCPSActionsRetrievedSuccessfully = ResponseCode{
		Code:       "SUCCESS_CPS_ACTIONS_RETRIEVED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgCPSActionsRetrieved,
		Type:       "success",
	}
	// cps_user related success response codes
	SuccessCpsUserCreationRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_CPS_USER_CREATION_REQUEST_SUBMITTED",
		StatusCode: StatusCreated,
		Message:    MsgCpsUserCreationRequestSubmitted,
		Type:       "success",
	}
	SuccessCpsUserUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_CPS_USER_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgCpsUserUpdateRequestSubmitted,
		Type:       "success",
	}
	SuccessCpsUserDeleted = ResponseCode{
		Code:       "SUCCESS_CPS_USER_DELETED",
		StatusCode: StatusOK,
		Message:    MsgCpsUserDeletedSuccessfully,
		Type:       "success",
	}
	SuccessCpsUsersRetrieved = ResponseCode{
		Code:       "SUCCESS_CPS_USERS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgCpsUsersRetrievedSuccessfully,
		Type:       "success",
	}
	SuccessCpsUserRetrieved = ResponseCode{
		Code:       "SUCCESS_CPS_USER_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgCpsUserRetrievedSuccessfully,
		Type:       "success",
	}
	SucessCpsUserEnabled = ResponseCode{
		Code:       "SUCCESS_CPS_USER_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgCpsUserEnabledSuccessfully,
		Type:       "success",
	}
	SuccessCpsUserDisabled = ResponseCode{
		Code:       "SUCCESS_CPS_USER_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgCpsUserDisabledSuccessfully,
		Type:       "success",
	}

	SuccessBpsUserEnableRequestSent = ResponseCode{
		Code:       "SUCCESS_BPS_USER_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgBpsUserEnabledRequestedSuccessfully,
		Type:       "success",
	}
	SuccessBpsUserDisableRequestSent = ResponseCode{
		Code:       "SUCCESS_BPS_USER_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgBpsUserDisabledRequestedSuccessfully,
		Type:       "success",
	}

	// CPS Roles
	SuccessCPSRoleCreated = ResponseCode{
		Code:       "SUCCESS_CPS_ROLE",
		StatusCode: StatusOK,
		Message:    MsgCpsRoleCreated,
		Type:       "success",
	}
	SuccessCPSRoleUpdated = ResponseCode{
		Code:       "SUCCESS_CPS_ROLE",
		StatusCode: StatusOK,
		Message:    MsgCpsRoleUpdated,
		Type:       "success",
	}
	SuccessCPSRolesFetched = ResponseCode{
		Code:       "SUCCESS_CPS_ROLES_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgCpsRolesFetched,
		Type:       "success",
	}
	SuccessCPSRoleFetched = ResponseCode{
		Code:       "SUCCESS_CPS_ROLE_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgCpsRoleFetched,
		Type:       "success",
	}
	SuccessCPSRoleEnabled = ResponseCode{
		Code:       "SUCCESS_CPS_ROLE_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgCpsRoleEnabled,
		Type:       "success",
	}
	SuccessCPSRoleDisabled = ResponseCode{
		Code:       "SUCCESS_CPS_ROLE_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgCpsRoleDisabled,
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
		Code:       "SUCCESS_WALLET_DELETE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgWalletDeleteRequestSent,
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

	ErrorWalletIDRequired = ResponseCode{
		Code:       "ERROR_WALLET_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Wallet ID is required",
		Type:       "error",
	}

	SuccessWalletEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_WALLET_ENABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    "Wallet enable request sent successfully",
		Type:       "success",
	}

	SuccessWalletDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_WALLET_DISABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    "Wallet disable request sent successfully",
		Type:       "success",
	}

	//topup realted success response codes
	SuccessTopupCreationRequestSent = ResponseCode{
		Code:       "SUCCESS_TOPUP_CREATION_REQUEST_SENT",
		StatusCode: StatusCreated,
		Message:    MsgTopupCreationRequestSent,
		Type:       "success",
	}

	SuccessTopupUpdateRequestSent = ResponseCode{
		Code:       "SUCCESS_TOPUP_UPDATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgTopupUpdateRequestSent,
		Type:       "success",
	}

	SuccessTopupDeleted = ResponseCode{
		Code:       "SUCCESS_TOPUP_DELETE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgTopupDeleteRequestSent,
		Type:       "success",
	}

	SuccessTopupsRetrieved = ResponseCode{
		Code:       "SUCCESS_TopupS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgTopupsRetrievedSuccessfully,
		Type:       "success",
	}

	SuccessTopupRetrieved = ResponseCode{
		Code:       "SUCCESS_TOPUP_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgTopupRetrievedSuccessfully,
		Type:       "success",
	}

	ErrorTopupIDRequired = ResponseCode{
		Code:       "ERROR_TOPUP_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Topup ID is required",
		Type:       "error",
	}

	SuccessTopupEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_TOPUP_ENABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    "Topup enable request sent successfully",
		Type:       "success",
	}

	SuccessTopupDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_TOPUP_DISABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    "Topup disable request sent successfully",
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
	// bankvault related
	SuccessBankVaultCreationRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_BANKVAULT_CREATION_REQUEST_SUBMITTED",
		StatusCode: StatusCreated,
		Message:    MsgBankVaultCreationRequestSubmitted,
		Type:       "success",
	}
	SuccessBankVaultsRetrieved = ResponseCode{
		Code:       "SUCCESS_BANKVAULTS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgBankVaultsRetrievedSuccessfully,
		Type:       "success",
	}
	SuccessBankVaultRetrieved = ResponseCode{
		Code:       "SUCCESS_BANKVAULT_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgBankVaultRetrievedSuccessfully,
		Type:       "success",
	}
	SuccessBankVaultUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_BANKVAULT_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgBankVaultUpdateRequestSubmitted,
		Type:       "success",
	}
	SuccessBankVaultDeleteRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_BANKVAULT_DELETE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgBankVaultDeleteRequestSubmitted,
		Type:       "success",
	}
	SuccessBankVaultEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_BANKVAULT_ENABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgBankVaultEnableRequestSubmitted,
		Type:       "success",
	}
	SuccessBankVaultDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_BANKVAULT_DISABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgBankVaultDisableRequestSubmitted,
		Type:       "success",
	}
	SuccessAllBankLockedVaultsRetrievedSuccessfully = ResponseCode{
		Code:       "SUCCESS_BANK_LOCKED_VAULTS_RETRIEVED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgBankLockedVaultsRetrievedSuccessfully,
		Type:       "success",
	}
	SuccessBankLockedVaultsRetrievedSuccessfully = ResponseCode{
		Code:       "SUCCESS_BANK_LOCKED_VAULT_RETRIEVED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgBankLockedVaultRetrievedSuccessfully,
		Type:       "success",
	}
	SuccessGroupVaultsRetrievedSuccessfully = ResponseCode{
		Code:       "SUCCESS_GROUP_VAULTS_RETRIEVED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgGroupVaultsRetrievedSuccessfully,
		Type:       "success",
	}
	SuccessGroupVaultRetrievedSuccessfully = ResponseCode{
		Code:       "SUCCESS_GROUP_VAULT_RETRIEVED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgGroupVaultRetrievedSuccessfully,
		Type:       "success",
	}
	// vaultgroup category related
	SuccessVaultGroupCategoryCreationRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULTGROUPCATEGORY_CREATION_REQUEST_SUBMITTED",
		StatusCode: StatusCreated,
		Message:    MsgVaultGroupCategoryCreationRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultGroupCategoriesRetrieved = ResponseCode{
		Code:       "SUCCESS_VAULTGROUPCATEGORIES_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgVaultGroupCategoriesRetrieved,
		Type:       "success",
	}
	SuccessVaultGroupCategoryRetrieved = ResponseCode{
		Code:       "SUCCESS_VAULTGROUPCATEGORY_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgVaultGroupCategoryRetrieved,
		Type:       "success",
	}
	SuccessVaultGroupCategoryUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULTGROUPCATEGORY_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgVaultGroupCategoryUpdateRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultGroupCategoryDeleteRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULTGROUPCATEGORY_DELETE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgVaultGroupCategoryDeleteRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultGroupCategoryEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULTGROUPCATEGORY_ENABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgVaultGroupCategoryEnableRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultGroupCategoryDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULTGROUPCATEGORY_DISABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgVaultGroupCategoryDisableRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultAmountTierCreationRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULT_AMOUNT_TIER_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgVaultAmountTierCreated,
		Type:       "success",
	}
	SuccessVaultAmountTierFetchedSuccessfully = ResponseCode{
		Code:       "SUCCESS_VAULT_RETRIEVED_SUCCEFULLY",
		StatusCode: StatusOK,
		Message:    MsgVaultAmountTierFetchedSuccessfully,
		Type:       "success",
	}
	SuccessVaultAmountTierUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULT_AMOUNT_TIER_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgVaultAmountTierUpdateRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultAmountTierDeleteRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULT_AMOUNT_TIER_DELETED",
		StatusCode: StatusOK,
		Message:    MsgVaultAmountTierDeleteRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultAmountTierDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULT_AMOUNT_TIER_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgVaultAmountTierDisableRequestSubmitted,
		Type:       "success",
	}
	SuccessVaultAmountTierEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_VAULT_AMOUNT_TIER_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgVaultAmountTierEnableRequestSubmitted,
		Type:       "success",
	}

	// Event related error response codes for bankvault
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

	ErrorAuthTierAlreadyExists = ResponseCode{
		Code:       "ERROR_AUTH_TIER_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    "Auth tier already exists with the same values",
		Type:       "error",
	}

	ErrorInvalidMethod = ResponseCode{
		Code:       "ERROR_INVALID_METHOD",
		StatusCode: StatusBadRequest,
		Message:    "Invalid method value",
		Type:       "error",
	}

	ErrorInvalidAmounts = ResponseCode{
		Code:       "ERROR_INVALID_AMOUNTS",
		StatusCode: StatusBadRequest,
		Message:    "Invalid amounts",
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

	ErrorInvalidMerchantID = ResponseCode{
		Code:       "ERROR_INVALID_MERCHANT_ID",
		StatusCode: StatusBadRequest,
		Message:    "Invalid merchant ID",
		Type:       "error",
	}

	ErrorAdvertTitleAlreadyExists = ResponseCode{
		Code:       "ERROR_ADVERT_TITLE_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    "Advert title already exists",
		Type:       "error",
	}
	ErrorBudgetCategoryNameAlreadyExists = ResponseCode{
		Code:       "ERROR_BUDGET_CATEGORY_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    "budget category name already exists",
		Type:       "error",
	}

	ErrorAdvertTitleNotChanged = ResponseCode{
		Code:       "ERROR_ADVERT_TITLE_NOT_CHANGED",
		StatusCode: StatusBadRequest,
		Message:    "Advert title not changed",
		Type:       "error",
	}

	ErrorInvalidDateFormat = ResponseCode{
		Code:       "ERROR_INVALID_DATE_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    "Invalid date format",
		Type:       "error",
	}
	ErrorInvalidFileUpload = ResponseCode{
		Code:       "IMAGE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "App icon image is required",
		Type:       "error",
	}
	ErrorInvalidBooleanFormat = ResponseCode{
		Code:       "ERROR_INVALID_BOOLEAN_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    "Invalid boolean format",
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

	ErrorCustomerIDRequired = ResponseCode{
		Code:       "ERROR_CUSTOMER_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Customer ID is required in param",
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

	// product code related error response codes
	ErrorNoProductCodesProvided = ResponseCode{
		Code:       "ERROR_NO_PRODUCT_CODES_PROVIDED",
		StatusCode: StatusBadRequest,
		Message:    "No product codes provided",
		Type:       "error",
	}
	ErrorProductCodeValidationError = ResponseCode{
		Code:       "ERROR_PRODUCT_CODE_VALIDATION_ERROR",
		StatusCode: StatusBadRequest,
		Message:    "error validating product code ",
		Type:       "error",
	}
	ErrorProductCodesValidationError = ResponseCode{
		Code:       "ERROR_PRODUCT_CODES_VALIDATION_ERROR",
		StatusCode: StatusBadRequest,
		Message:    "error validating product codes ",
		Type:       "error",
	}
	ErrorProductCodeUpdateRequestValidationErrorAtLeastOne = ResponseCode{
		Code:       "ERROR_PRODUCT_CODE_UPDATE_REQUEST_VALIDATION_ERROR_ATLEAST_ONE",
		StatusCode: StatusBadRequest,
		Message:    "at least one of ProductName,CBEProductCodes,CBEIFBProductCodes should be present",
		Type:       "error",
	}
	ErrorBothProductCodesRequired = ResponseCode{
		Code:       "ERROR_BOTH_PRODUCT_CODES_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Product codes for both branches required",
		Type:       "error",
	}

	//
	ErrorTopupNameAlreadyExists = ResponseCode{
		Code:       "ERROR_TOPUP_WITH_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    "Topup with the given name already exists",
		Type:       "error",
	}

	ErrorTopupcoDEAlreadyExists = ResponseCode{
		Code:       "ERROR_TOPUP_WITH_CODE_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    "Topup with the given code already exists",
		Type:       "error",
	}

	ErrorTopupCodeAlreadyExists = ResponseCode{
		Code:       "ERROR_TOPUP_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    "Topup with the given code already exists",
		Type:       "error",
	}
	ErrorTopupNotFound = ResponseCode{
		Code:       "ERROR_TOPUP_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "Topup not found",
		Type:       "error",
	}

	ErrorTopupAlreadyEnabled = ResponseCode{
		Code:       "ERROR_TOPUP_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    "Topup is already enabled",
		Type:       "error",
	}
	ErrorTopupUpdateEmptyPayload = ResponseCode{
		Code:       "ERROR_TOPUP_UPDATE_EMPTY_PAYLOAD",
		StatusCode: StatusBadRequest,
		Message:    "No data provided for Topup update",
		Type:       "error",
	}

	ErrorTopupAlreadyDisabled = ResponseCode{
		Code:       "ERROR_TOPUP_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    "Topup is already disabled",
		Type:       "error",
	}
	//wallet related error codes

	ErrorWalletNameAlreadyExists = ResponseCode{
		Code:       "ERROR_WALLET_WITH_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    "Wallet with the given name or code already exists",
		Type:       "error",
	}
	ErrorWalletCodeAlreadyExists = ResponseCode{
		Code:       "ERROR_WALLET_WITH_CODE_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    "Wallet with the given code already exists",
		Type:       "error",
	}
	ErrorWalletNotFound = ResponseCode{
		Code:       "ERROR_WALLET_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "Wallet not found",
		Type:       "error",
	}

	ErrorWalletAlreadyEnabled = ResponseCode{
		Code:       "ERROR_WALLET_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    "Wallet is already enabled",
		Type:       "error",
	}
	ErrorWalletUpdateEmptyPayload = ResponseCode{
		Code:       "ERROR_WALLET_UPDATE_EMPTY_PAYLOAD",
		StatusCode: StatusBadRequest,
		Message:    "No data provided for wallet update",
		Type:       "error",
	}

	ErrorWalletAlreadyDisabled = ResponseCode{
		Code:       "ERROR_WALLET_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    "Wallet is already disabled",
		Type:       "error",
	}

	ErrorInvalidID = ResponseCode{
		Code:       "ERROR_INVALID_ID",
		StatusCode: 400,
		Message:    "Invalid ID",
		Type:       "error",
	}

	ErrorNoDataProvided = ResponseCode{
		Code:       "ERROR_NO_DATA_PROVIDED",
		StatusCode: 400,
		Message:    "No data provided for update",
		Type:       "error",
	}

	ErrorWalletNameRequired = ResponseCode{
		Code:       "ERROR_WALLET_NAME_REQUIRED",
		StatusCode: 400,
		Message:    "Wallet name is required",
		Type:       "error",
	}

	ErrorInvalidWalletName = ResponseCode{
		Code:       "ERROR_INVALID_WALLET_NAME",
		StatusCode: 400,
		Message:    "Invalid wallet name, special characters are not allowed",
		Type:       "error",
	}

	ErrorInvalidWalletCode = ResponseCode{
		Code:       "ERROR_INVALID_WALLET_CODE",
		StatusCode: 400,
		Message:    "Invalid wallet code, wallet should be only characters and exactly 6 characters long",
		Type:       "error",
	}
	ErrorInvalidWalletType = ResponseCode{
		Code:       "ERROR_INVALID_WALLET_TYPE",
		StatusCode: 400,
		Message:    "Invalid wallet type, special characters are not allowed",
		Type:       "error",
	}

	ErrorWalletCodeRequired = ResponseCode{
		Code:       "ERROR_WALLET_CODE_REQUIRED",
		StatusCode: 400,
		Message:    "Wallet code is required",
		Type:       "error",
	}
	ErrorWalletTypeRequired = ResponseCode{
		Code:       "ERROR_WALLET_TYPE_REQUIRED",
		StatusCode: 400,
		Message:    "Wallet type is required",
		Type:       "error",
	}

	ErrorWalletAvatarRequired = ResponseCode{
		Code:       "ERROR_WALLET_AVATAR_REQUIRED",
		StatusCode: 400,
		Message:    "Wallet avatar is required",
		Type:       "error",
	}
	ErrorWalletAvatarInvalid = ResponseCode{
		Code:       "ERROR_WALLET_AVATAR_INVALID",
		StatusCode: 400,
		Message:    "Invalid wallet avatar",
		Type:       "error",
	}

	ErrorWalletAvatarTooLarge = ResponseCode{
		Code:       "ERROR_WALLET_AVATAR_TOO_LARGE",
		StatusCode: 400,
		Message:    "Wallet avatar file size exceeds the limit",
		Type:       "error",
	}

	ErrorWalletAvatarInvalidType = ResponseCode{
		Code:       "ERROR_WALLET_AVATAR_INVALID_TYPE",
		StatusCode: 400,
		Message:    "Wallet avatar must be of type jpeg, png, gif, or webp",
		Type:       "error",
	}

	ErrorWalletRechangeOption = ResponseCode{
		Code:       "ERROR_WALLET_RECHARGE_OPTION_INVALID_VALUES",
		StatusCode: 400,
		Message:    "At Least one of the three rechange options should be enabled(self,other,agent)",
		Type:       "error",
	}

	ErrorAvatarAlreadyExist = ResponseCode{
		Code:       "ERROR_AVATAR_ALREADY_EXIST",
		StatusCode: StatusConflict,
		Message:    "Avatar label already exist",
		Type:       "error",
	}

	ErrorAvatarNotExist = ResponseCode{
		Code:       "ERROR_AVATAR_NOT_EXIST",
		StatusCode: StatusConflict,
		Message:    "Avatar label not exist",
		Type:       "error",
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
	// cps user related success response codes
	ErrorGroupNameRequired = ResponseCode{
		Code:       "ERROR_GROUP_NAME_REQUIRED",
		StatusCode: StatusNotFound,
		Message:    MsgGroupNameRequired,
		Type:       "error",
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
	SuccessMiniAppCreateRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_CREATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgMiniAppCreateRequestSuccessfully,
		Type:       "success",
	}
	SuccessMiniAppUpdateRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_UPDATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgMiniAppMerchantUpdateRequestSuccessfully,
		Type:       "success",
	}

	// Mini App Merchant related success response codes
	SuccessMiniAppMerchantCreateRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_MERCHANT_CREATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgMiniAppMerchantCreateRequestSuccessfully,
		Type:       "success",
	}
	SuccessMiniAppMerchantUpdateRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_MERCHANT_UPDATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgMiniAppMerchantUpdateRequestSuccessfully,
		Type:       "success",
	}
	SuccessMiniAppMerchantEnableRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_MERCHANT_ENABLE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgEnableRequestSuccessfullyCreated,
		Type:       "success",
	}
	SuccessMiniAppMerchantDisableRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_MERCHANT_DISABLE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDisableRequestSuccessfullyCreated,
		Type:       "success",
	}
	SuccessMiniAppMerchantDeleteRequestCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_MERCHANT_DELETE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDeleteRequestSuccessfullyCreated,
		Type:       "success",
	}
	SuccessMiniAppMerchantCreated = ResponseCode{
		Code:       "SUCCESS_MINI_APP_MERCHANT_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgMiniAppMerchantCreatedSuccessfully,
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
	SuccessMiniAppEnable = ResponseCode{
		Code:       "SUCCESS_ENABLE_SUCCESS",
		StatusCode: StatusOK,
		Message:    MsgMiniAppMerchantEnableSuccessfully,
		Type:       "success",
	}
	SuccessMiniAppDesable = ResponseCode{
		Code:       "SUCCESS_DISABLE_SUCCESS",
		StatusCode: StatusOK,
		Message:    MsgDisableRequestSuccessfullyCreated,
		Type:       "success",
	}

	SuccessBankDisableRequestCreated = ResponseCode{
		Code:       "SUCCESS_BANK_DISABLE_REQUEST_CREATED",
		StatusCode: StatusOK,
		Message:    MsgBankDisableRequestSent,
		Type:       "success",
	}

	SuccessBankEnableRequestCreated = ResponseCode{
		Code:       "SUCCESS_BANK_ENABLE_REQUEST_CREATED",
		StatusCode: StatusOK,
		Message:    MsgBankEnableRequestSent,
		Type:       "success",
	}

	// Ecommerce merchant
	SuccessEcommerceMerchantCreatedSuccessfully = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANT_CREATED",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantCreated,
		Type:       "success",
	}
	SuccessEcommerceMerchantUpdatedSuccessfully = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANT_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantUpdated,
		Type:       "success",
	}
	SuccessEcommerceMerchantDeletedSuccessfully = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANT_DELETED",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantDeleted,
		Type:       "success",
	}
	SuccessEcommerceMerchantEnableSuccessfully = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANT_ENABLE",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantEnable,
		Type:       "success",
	}
	SuccessEcommerceMerchantDisableSuccessfully = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANT_DISABLE",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantDisable,
		Type:       "success",
	}
	SuccessEcommerceMerchantFetchedSuccessfully = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANT_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantFetched,
		Type:       "success",
	}
	SuccessEcommerceMerchantsFetchedSuccessfully = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANTS_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantsFetched,
		Type:       "success",
	}
	SuccessEcommerceMerchantLookup = ResponseCode{
		Code:       "SUCCESS_ECOMMERCE_MERCHANT_LOOKUP",
		StatusCode: StatusOK,
		Message:    MsgEcommerceMerchantLookup,
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

	// department related success response codes
	SuccessGetAllDepartments = ResponseCode{
		Code:       "SUCCESS_GET_ALL_DEPARTMENTS",
		StatusCode: StatusOK,
		Message:    msgGetAllDepartmentsSuccess,
		Type:       "success",
	}

	SuccessGetDepartments = ResponseCode{
		Code:       "SUCCESS_GET_DEPARTMENTS",
		StatusCode: StatusOK,
		Message:    msgGetDepartmentsSuccess,
		Type:       "success",
	}
	SuccessDepartmentCreateRequestCreated = ResponseCode{
		Code:       "SUCCESS_DEPARTMENT_CREATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentCreateRequestedSuccessfully,
		Type:       "success",
	}
	SuccessDepartmentUpdateRequestCreated = ResponseCode{
		Code:       "SUCCESS_DEPARTMENT_UPDATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentUpdateRequestedSuccessfully,
		Type:       "success",
	}
	SuccessDepartmentEnableRequestCreated = ResponseCode{
		Code:       "SUCCESS_DEPARTMENT_UPDATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentEnableRequestedSuccessfully,
		Type:       "success",
	}
	SuccessDepartmentDisableRequestCreated = ResponseCode{
		Code:       "SUCCESS_DEPARTMENT_UPDATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentDisableRequestedSuccessfully,
		Type:       "success",
	}

	SuccessDepartmentUpdateCPSActionCreated = ResponseCode{
		Code:       "SUCCESS_DEPARTMENT_UPDATE_CPS_ACTION_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgDepartmentUpdateCPSActionCreated,
		Type:       "success",
	}

	// Budget Category related success response codes
	SuccessBudgetCategoryRequestSubmittedForApproval = ResponseCode{
		Code:       "SUCCESS_BUDGET_CATEGORY_REQUEST_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusCreated,
		Message:    MsgBudgetCategoryRequestSubmittedForApprovalSuccess,
		Type:       "success",
	}

	SuccessBudgetCategoryUpdateSubmittedForApproval = ResponseCode{
		Code:       "SUCCESS_BUDGET_CATEGORY_UPDATE_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoryUpdateSubmittedForApprovalSuccess,
		Type:       "success",
	}

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
		Message:    MsgValidationRuleSuccessFech,
		Type:       "success",
	}

	// Password rule
	SuccessFetchAllPasswordRules = ResponseCode{
		Code:       "SUCCESS_ALL_PASSWORD_RULE_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgFetchAllPasswordRules,
		Type:       "success",
	}
	SuccessUpdatePasswordRule = ResponseCode{
		Code:       "SUCCESS_UPDATE_PASSWORD_RULE",
		StatusCode: StatusOK,
		Message:    MsgUpdatePasswordRule,
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
		Message:    MsgDeleteRequestSuccessfullyCreated,
		Type:       "success",
	}

	SuccessMiniAppDisableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_MINI_APP_DISABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgDisableRequestSuccessfullyCreated,
		Type:       "success",
	}
	SuccessMiniAppEnableRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_MINI_APP_ENABLE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgEnableRequestSuccessfullyCreated,
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
	// mini app handler related error response codes

	// Encryption
	SuccessEncryptionGenerated = ResponseCode{
		Code:       "SUCCESS_ENCRYPTION",
		StatusCode: StatusOK,
		Message:    MsgEncryptionSuccessfully,
		Type:       "success",
	}

	ErrorMiniAppNotFound = ResponseCode{
		Code:       "ERROR_MINI_APP_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "Mini App not found",
		Type:       "error",
	}

	ErrorMiniAppAlreadyExists = ResponseCode{
		Code:       "ERROR_MINI_APP_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    "A Mini App with this name already exists",
		Type:       "error",
	}
	ErrorMiniAppAlreadyEnabled = ResponseCode{
		Code:       "ERROR_MINI_APP_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App is already enabled",
		Type:       "error",
	}
	ErrorMiniAppAlreadyDisabled = ResponseCode{
		Code:       "ERROR_MINI_APP_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App is already disabled",
		Type:       "error",
	}
	ErrorMiniAppAlreadyDeleted = ResponseCode{
		Code:       "ERROR_MINI_APP_ALREADY_DELETED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App is already deleted",
		Type:       "error",
	}
	ErrorMiniAppIDRequired = ResponseCode{
		Code:       "ERROR_MINI_APP_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App ID is required",
		Type:       "error",
	}
	ErrorMiniAppMerchantIDRequired = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App Merchant ID is required",
		Type:       "error",
	}
	ErrorMiniAppNameRequired = ResponseCode{
		Code:       "ERROR_MINI_APP_NAME_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App Name is required",
		Type:       "error",
	}
	ErrorMiniAppDescriptionRequired = ResponseCode{
		Code:       "ERROR_MINI_APP_DESCRIPTION_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App Description is required",
		Type:       "error",
	}
	ErrorMiniAppURLRequired = ResponseCode{
		Code:       "ERROR_MINI_APP_URL_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Mini App URL is required",
		Type:       "error",
	}

	ErrorAppIconRequired = ResponseCode{
		Code:       "ERROR_APP_ICON_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "App Icon is required",
		Type:       "error",
	}
	ErrorInvalidAppNameFormat = ResponseCode{
		Code:       "ERROR_INVALID_APP_NAME_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    "App Name must not contain special characters",
		Type:       "error",
	}
	ErrorAppViewTypeRequired = ResponseCode{
		Code:       "ERROR_APP_VIEW_TYPE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "App View Type is required",
		Type:       "error",
	}
	ErrorBannerImageRequired = ResponseCode{
		Code:       "ERROR_BANNER_IMAGE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    "Banner Image is required",
		Type:       "error",
	}

	ErrorInvalidImageFormat = ResponseCode{
		Code:       "ERROR_INVALID_IMAGE_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    "Image must be a valid format (jpeg, png, gif)",
		Type:       "error",
	}

	ErrorInvalidURL = ResponseCode{
		Code:       "ERROR_INVALID_URL",
		StatusCode: StatusBadRequest,
		Message:    "Invalid URL format",
		Type:       "error",
	}
	ErrorIncompleteBranchProductCodes = ResponseCode{
		Code:       "ERROR_INCOMPLETE_BRANCH_PRODUCT_CODES",
		StatusCode: StatusBadRequest,
		Message:    "Incomplete product codes for branch",
		Type:       "error",
	}

	ErrorAppViewTypeInvalidOrMissing = ResponseCode{
		Code:       "ERROR_APP_VIEW_TYPE_INVALID_OR_MISSING",
		StatusCode: StatusBadRequest,
		Message:    "App View Type is invalid or missing",
		Type:       "error",
	}
	ErrorInvalidAppViewType = ResponseCode{
		Code:       "ERROR_INVALID_APP_VIEW_TYPE",
		StatusCode: StatusBadRequest,
		Message:    "Invalid App View Type",
		Type:       "error",
	}
	ErrorExclusiveAppFlags = ResponseCode{
		Code:       "ERROR_EXCLUSIVE_APP_FLAGS",
		StatusCode: StatusBadRequest,
		Message:    "IsEventMiniApp and IsThreeClick cannot both be true",
		Type:       "error",
	}
	ErrorUpdateMiniAppEmptyPayload = ResponseCode{
		Code:       "ERROR_UPDATE_MINI_APP_EMPTY_PAYLOAD",
		StatusCode: StatusBadRequest,
		Message:    "No data provided for mini app update",
		Type:       "error",
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
	SuccessHQBlockTimeFetched = ResponseCode{
		Code:       "SUCCESS_HQ_BLOCK_TIME_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgHQBlockTImeFetchedSuccessfully,
		Type:       "success",
	}

	SuccessHQArchiveTimeFetched = ResponseCode{
		Code:       "SUCCESS_HQ_ARCHIVE_TIME_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgHQArchiveTimeFetchedSuccessfully,
		Type:       "success",
	}
	SuccessHQPasswordExpiryFetched = ResponseCode{
		Code:       "SUCCESS_HQ_PASSWORD_EXPIRY_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgHQPasswordExpiryFetchedSuccessfully,
		Type:       "success",
	}

	SuccessHQBlockTimeUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_HQ_BLOCK_TIME_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgHQBlockTimeUpdateRequestSubmitted,
		Type:       "success",
	}

	SuccessHQArchiveTimeUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_HQ_ARCHIVE_TIME_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgHQArchiveTimeUpdateRequestSubmitted,
		Type:       "success",
	}

	SuccessHQPasswordExpiryUpdateRequestSubmitted = ResponseCode{
		Code:       "SUCCESS_HQ_PASSWORD_EXPIRY_UPDATE_REQUEST_SUBMITTED",
		StatusCode: StatusOK,
		Message:    MsgHQPasswordExpiryUpdateRequestSubmitted,
		Type:       "success",
	}

	ErrorFailedToParseJson = ResponseCode{Code: "ERROR_FAILED_TO_PARSE_JSON", StatusCode: 500, Message: "Failed to parse json", Type: "error"}

	// HQ related error response codes

	ErrorHQNotFound       = ResponseCode{Code: "ERROR_HQ_NOT_FOUND", StatusCode: 404, Message: "HQ record not found", Type: "error"}
	ErrorInvalidHQRequest = ResponseCode{Code: "ERROR_INVALID_HQ_REQUEST", StatusCode: 400, Message: "Invalid HQ request", Type: "error"}
	ErrorCPSActionFailed  = ResponseCode{Code: "ERROR_CPS_ACTION_FAILED", StatusCode: 500, Message: "Failed to handle CPS action", Type: "error"}
	ErrorHQIDRequired     = ResponseCode{Code: "ERROR_HQ_ID_REQUIRED", StatusCode: 400, Message: "HQ ID is required", Type: "error"}

	ErrorInvalidBlockTime      = ResponseCode{Code: "ERROR_INVALID_BLOCK_TIME", StatusCode: 400, Message: "BlockTime must be greater than 0", Type: "error"}
	ErrorInvalidArchiveTime    = ResponseCode{Code: "ERROR_INVALID_ARCHIVE_TIME", StatusCode: 400, Message: "ArchiveTime must be greater than 0", Type: "error"}
	ErrorInvalidPasswordExpiry = ResponseCode{Code: "ERROR_INVALID_PASSWORD_EXPIRY", StatusCode: 400, Message: "PasswordExpiry must be greater than 0", Type: "error"}

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

	SuccessBudgetIconRequestSubmittedForApproval = ResponseCode{
		Code:       "BUDGET_ICON_REQUEST_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusOK,
		Message:    MsgBudgetIconRequestSubmittedForApprovalSuccess,
		Type:       "success",
	}

	SuccessBudgetRequestSubmittedForApproval = ResponseCode{
		Code:       "BUDGET_REQUEST_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusOK,
		Message:    MsgBudgetRequestSubmittedForApprovalSuccess,
		Type:       "success",
	}

	SuccessBudgetUpdateSubmittedForApproval = ResponseCode{
		Code:       "BUDGET_UPDATE_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusOK,
		Message:    MsgBudgetUpdateSubmittedForApprovalSuccess,
		Type:       "success",
	}

	SuccessBudgetCategoryFetched = ResponseCode{
		Code:       "BUDGET_CATEGORY_FETCHED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoryFetchedSuccessfully,
		Type:       "success",
	}

	SuccessBudgetCategoriesFetched = ResponseCode{
		Code:       "BUDGET_CATEGORIES_FETCHED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoriesFetchedSuccessfully,
		Type:       "success",
	}

	SuccessBudgetCategoryDeleteSubmittedForApproval = ResponseCode{
		Code:       "BUDGET_CATEGORY_DELETE_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoryDeleteSubmittedForApprovalSuccess,
		Type:       "success",
	}

	SuccessBudgetCategoryEnableSubmittedForApproval = ResponseCode{
		Code:       "BUDGET_CATEGORY_ENABLE_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoryEnableSubmittedForApprovalSuccess,
		Type:       "success",
	}
	SuccessBudgetCategoryDisableSubmittedForApproval = ResponseCode{
		Code:       "BUDGET_CATEGORY_DISABLE_SUBMITTED_FOR_APPROVAL",
		StatusCode: StatusOK,
		Message:    MsgBudgetCategoryDisableSubmittedForApprovalSuccess,
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

	SuccessBranchRetrieved = ResponseCode{
		Code:       "SUCCESS_BRANCH_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgBranchSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessBranchesRetrieved = ResponseCode{
		Code:       "SUCCESS_BRANCHES_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgBranchesSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessBranchesEnabled = ResponseCode{
		Code:       "SUCCESS_BRANCHES_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgBranchesSuccessfullyEnabled,
		Type:       "success",
	}

	SuccessBranchesDisabled = ResponseCode{
		Code:       "SUCCESS_BRANCHES_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgBranchesSuccessfullyDisabled,
		Type:       "success",
	}

	SuccessEnableBranchesRequestSent = ResponseCode{
		Code:       "SUCCESS_ENABLE_BRANCHES_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgEnableBranchesRequestSent,
		Type:       "success",
	}

	SuccessDisableBranchesRequestSent = ResponseCode{
		Code:       "SUCCESS_DISABLE_BRANCHES_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDisableBranchesRequestSent,
		Type:       "success",
	}

	SuccessRegionRetrieved = ResponseCode{
		Code:       "SUCCESS_REGION_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgRegionSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessRegionsRetrieved = ResponseCode{
		Code:       "SUCCESS_REGIONS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgRegionsSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessRegionsEnabled = ResponseCode{
		Code:       "SUCCESS_REGIONSS_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgRegionsSuccessfullyEnabled,
		Type:       "success",
	}

	SuccessRegionsDisabled = ResponseCode{
		Code:       "SUCCESS_REGIONSS_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgRegionsSuccessfullyDisabled,
		Type:       "success",
	}

	SuccessEnableRegionsRequestSent = ResponseCode{
		Code:       "SUCCESS_ENABLE_REGIONS_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgEnableRegionsRequestSent,
		Type:       "success",
	}

	SuccessDisableRegionsRequestSent = ResponseCode{
		Code:       "SUCCESS_DISABLE_REGIONS_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDisableRegionsRequestSent,
		Type:       "success",
	}

	SuccessDistrictRetrieved = ResponseCode{
		Code:       "SUCCESS_DISTRICT_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgDistrictSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessDistrictsRetrieved = ResponseCode{
		Code:       "SUCCESS_DISTRICTS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgDistrictsSuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessDistrictsEnabled = ResponseCode{
		Code:       "SUCCESS_DISTRICTS_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgDistrictsSuccessfullyEnabled,
		Type:       "success",
	}

	SuccessDistrictsDisabled = ResponseCode{
		Code:       "SUCCESS_DISTRICTS_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgDistrictsSuccessfullyDisabled,
		Type:       "success",
	}

	SuccessEnableDistrictsRequestSent = ResponseCode{
		Code:       "SUCCESS_ENABLE_DISTRICTS_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgEnableDistrictsRequestSent,
		Type:       "success",
	}

	SuccessDisableDistrictsRequestSent = ResponseCode{
		Code:       "SUCCESS_DISABLE_DISTRICTS_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDisableDistrictsRequestSent,
		Type:       "success",
	}

	SuccessCityRetrieved = ResponseCode{
		Code:       "SUCCESS_CITY_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgCitySuccessfullyRetrieved,
		Type:       "success",
	}

	SuccessCitiesRetrieved = ResponseCode{
		Code:       "SUCCESS_CITIES_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgCitiesSuccessfullyRetrieved,
		Type:       "success",
	}
	SuccessCitiesEnabled = ResponseCode{
		Code:       "SUCCESS_CITIES_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgCitiesSuccessfullyEnabled,
		Type:       "success",
	}

	SuccessCitiesDisabled = ResponseCode{
		Code:       "SUCCESS_CITIES_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgCitiesSuccessfullyDisabled,
		Type:       "success",
	}
	SuccessEnableCitiesRequestSent = ResponseCode{
		Code:       "SUCCESS_ENABLE_CITIES_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgEnableCitiesRequestSent,
		Type:       "success",
	}

	SuccessDisableCitiesRequestSent = ResponseCode{
		Code:       "SUCCESS_DISABLE_CITIES_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgDisableCitiesRequestSent,
		Type:       "success",
	}

	SuccessGetAllBanks = ResponseCode{
		Code:       "SUCCESS_GET_ALL_BANKS",
		StatusCode: StatusOK,
		Message:    msgGetAllBanksSuccess,
		Type:       "success",
	}

	SuccessGetOneBank = ResponseCode{
		Code:       "SUCCESS_GET_ONE_BANK",
		StatusCode: StatusOK,
		Message:    msgGetOneBankSuccess,
		Type:       "success",
	}

	SuccessDeleteBanksRequest = ResponseCode{
		Code:       "SUCCESS_DELETE_BANK_request",
		StatusCode: StatusNoContent,
		Message:    msgDeleteBankRequestSuccess,
		Type:       "success",
	}

	// Fayda Account
	SuccessFaydaEnableActionCreated = ResponseCode{
		Code:       "SUCCESS_FAYDA_ENABLE_ACTION_CREATED",
		StatusCode: StatusOK,
		Message:    MsgFaydaAccountEnableCreatedSuccessfully,
		Type:       "success",
	}
	SuccessFaydaDisableActionCreated = ResponseCode{
		Code:       "SUCCESS_FAYDA_DISABLE_ACTION_CREATED",
		StatusCode: StatusOK,
		Message:    MsgFaydaAccountDisableCreatedSuccessfully,
		Type:       "success",
	}
)

// Action Role success response codes
var (
	SuccessActionRolesFetched = ResponseCode{
		Code:       "SUCCESS_ACTION_ROLES_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgActionRolesFetchedSuccess,
		Type:       "success",
	}
	SuccessActionRoleFetched = ResponseCode{
		Code:       "SUCCESS_ACTION_ROLE_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgActionRoleFetchedSuccess,
		Type:       "success",
	}
	SuccessActionRoleCreateRequestCreated = ResponseCode{
		Code:       "SUCCESS_ACTION_ROLE_CREATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgActionRoleCreateRequestCreated,
		Type:       "success",
	}
	SuccessActionRoleUpdateRequestCreated = ResponseCode{
		Code:       "SUCCESS_ACTION_ROLE_UPDATE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgActionRoleUpdateRequestCreated,
		Type:       "success",
	}
	SuccessActionRoleEnableRequestCreated = ResponseCode{
		Code:       "SUCCESS_ACTION_ROLE_ENABLE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgActionRoleEnableRequestCreated,
		Type:       "success",
	}
	SuccessActionRoleDisableRequestCreated = ResponseCode{
		Code:       "SUCCESS_ACTION_ROLE_DISABLE_REQUEST_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgActionRoleDisableRequestCreated,
		Type:       "success",
	}

	SuccessEventMerchantCreated = ResponseCode{
		Code:       "SUCCESS_EVENT_MERCHANT_CREATED",
		StatusCode: StatusCreated,
		Message:    MsgEventMerchantCreatedSuccessfully,
		Type:       "success",
	}
	SuccessEventMerchantDeleted = ResponseCode{
		Code:       "SUCCESS_EVENT_MERCHANT_DELETED",
		StatusCode: StatusOK,
		Message:    MsgEventMerchantDeletedSuccessfully,
		Type:       "success",
	}
	SuccessEventMerchantUpdated = ResponseCode{
		Code:       "SUCCESS_EVENT_MERCHANT_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgEventMerchantUpdatedSuccessfully,
		Type:       "success",
	}
	SuccessEventMerchantEnabled = ResponseCode{
		Code:       "SUCCESS_EVENT_MERCHANT_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgEventMerchantEnabledSuccessfully,
		Type:       "success",
	}
	SuccessEventMerchantDisabled = ResponseCode{
		Code:       "SUCCESS_EVENT_MERCHANT_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgEventMerchantDisabledSuccessfully,
		Type:       "success",
	}
	SuccessEventMerchantFetched = ResponseCode{
		Code:       "SUCCESS_EVENT_MERCHANT_FETCHED",
		StatusCode: StatusOK,
		Message:    MsgEventMerchantFetchedSuccessfully,
		Type:       "success",
	}
)

// Error Response Codes
var (
	ErrorDepartmentIDRequired = ResponseCode{
		Code:       "ERROR_DEPARTMENT_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgDepartmentIDRequired,
		Type:       "error",
	}

	ErrorInvalidKey = ResponseCode{
		Code:       "ERROR_INVALID_KEY",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidKey,
		Type:       "error",
	}
	ErrorInvalidEncData = ResponseCode{
		Code:       "ERROR_INVALID_ENC_DATA",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidEncData,
		Type:       "error",
	}

	ErrorInvalidPadding = ResponseCode{
		Code:       "ERROR_INVALID_PADDING",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidPadding,
		Type:       "error",
	}

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
	ErrorDuplicateColorExists = ResponseCode{
		Code:       "ERROR_DUPLICATE_COLOR",
		StatusCode: StatusBadRequest,
		Message:    MsgDuplicateColorExists,
		Type:       "error",
	}

	ErrorFeedbackIDRequired = ResponseCode{
		Code:       "ERROR_FEEDBACK_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgFeedbackIDRequired,
		Type:       "error",
	}
	ErrorServiceDetailIDRequired = ResponseCode{
		Code:       "ERROR_SERVICE_DETAIL_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgServiceIdRequered,
		Type:       "error",
	}

	ErrorInvalidIDFormat = ResponseCode{
		Code:       "ERROR_INVALID_ID_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidIDFormat,
		Type:       "error",
	}

	ErrorInvalidJSONPayload = ResponseCode{
		Code:       "ERROR_INVALID_JSON_PAYLOAD",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidJSONPayload,
		Type:       "error",
	}

	ErrorIncompleteUserInfo = ResponseCode{
		Code:       "ERROR_INCOMPLTE_USER_INFO",
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

	ErrorThirdAPIRequestNotUnauthorized = ResponseCode{
		Code:       "ERROR_THIRD_API_REQUEST_UNAUTHORIZED",
		StatusCode: StatusUnauthorized,
		Message:    MsgThirdApiRequestNotAuthorized,
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

	ErrorPendingCpsActionExists = ResponseCode{
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

	ErrorBankVaultAlreadyEnabled = ResponseCode{
		Code:       "ERROR_BANK_VAULT_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgBankVaultAlreadyEnabled,
		Type:       "error",
	}
	ErrorBankVaultAlreadyDisabled = ResponseCode{
		Code:       "ERROR_BANK_VAULT_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgBankVaultAlreadyDisabled,
		Type:       "error",
	}

	ErrorVaultGroupAlreadyEnabled = ResponseCode{
		Code:       "ERROR_VAULT_GROUP_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgVaultGroupAlreadyEnabled,
		Type:       "error",
	}
	ErrorVaultGroupAlreadyDisabled = ResponseCode{
		Code:       "ERROR_VAULT_GROUP_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgVaultGroupAlreadyDisabled,
		Type:       "error",
	}

	ErrorBankAlreadyEnabled = ResponseCode{
		Code:       "ERROR_BANK_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgBankAlreadyEnabled,
		Type:       "error",
	}
	ErrorAmountTierNotFound = ResponseCode{
		Code:       "ERROR_AMOUNT_TIER_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgAmountTierNotFound,
		Type:       "error",
	}

	ErrorBankAlreadyDisabled = ResponseCode{
		Code:       "ERROR_Bank_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgBankAlreadyDisabled,
		Type:       "error",
	}

	ErrorBankWithNameAlreadyExists = ResponseCode{
		Code:       "ERROR_Bank_WITH_NAME_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgBankNameAlreadyExists,
		Type:       "error",
	}
	ErrorBankWithBICAlreadyExists = ResponseCode{
		Code:       "ERROR_Bank_WITH_BIC_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgBankBICAlreadyExists,
		Type:       "error",
	}
	ErrorBankWithCodeAlreadyExists = ResponseCode{
		Code:       "ERROR_Bank_WITH_CODE_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgBankCodeAlreadyExists,
		Type:       "error",
	}

	ErrorOTPNotFound = ResponseCode{
		Code:       "ERROR_OTP_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgOTPNotFound,
		Type:       "error",
	}
	ErrorFailedToDecodeRequest = ResponseCode{
		Code:       "ERROR_FAILED_TO_DECODE_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    "Failed to decode request",
		Type:       "error",
	}

	ErrorInvalidRequest = ResponseCode{
		Code:       "ERROR_INVALID_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequest,
		Type:       "error",
	}

	ErrorInvalidBankRequest = ResponseCode{
		Code:       "ERROR_INVALID_BANK_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidBankRequest,
		Type:       "error",
	}

	ErrorNoDataProvidedForUpdate = ResponseCode{
		Code:       "ERROR_NO_DATA_PROVIDED_FOR_UPDATE",
		StatusCode: StatusBadRequest,
		Message:    MsgNoDataProvidedForUpdate,
		Type:       "error",
	}

	ErrorNoDataProvidedForBankUpdate = ResponseCode{
		Code:       "ERROR_NO_DATA_PROVIDED_FOR_UPDATE",
		StatusCode: StatusBadRequest,
		Message:    MsgNoDataProvidedForBankUpdate,
		Type:       "error",
	}

	ErrorInvalidFormatForName = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT_FOR_NAME",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestBankName,
		Type:       "error",
	}

	ErrorInvalidFormatForCode = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT_FOR_CODE",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestBankCode,
		Type:       "error",
	}

	ErrorInvalidFormatForBIC = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT_FOR_BIC",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestBankBIC,
		Type:       "error",
	}
	ErrorInvalidFormatForType = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT_FOR_TYPE",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestBankType,
		Type:       "error",
	}
	ErrorInvalidAccountNumberFormat = ResponseCode{
		Code:       "ERROR_INVALID_ACCOUNT_NUMBER_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    MsgBankAccountLengthInvalid,
		Type:       "error",
	}
	ErrorNoDataProvidedForCreate = ResponseCode{
		Code:       "ERROR_NO_DATA_PROVIDED_FOR_CREATE",
		StatusCode: StatusBadRequest,
		Message:    "No data provided for Create",
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

	ErrorMiniAppNameAlreadyExists = ResponseCode{
		Code:       "ERROR_MINI_APP_NAME_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgMiniAppNameAlreadyExists,
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

	ErrorServiceExists = ResponseCode{
		Code:       "ERROR_SERVICE_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgServiceExists,
		Type:       "error",
	}

	ErrorGetAllBanksFailed = ResponseCode{
		Code:       "ERROR_GET_ALL_BANKS_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    msgGetAllBanksFailed,
		Type:       "error",
	}

	ErrorBankDeleteRequestFailed = ResponseCode{
		Code:       "ERROR_BANK_DELETE_REQUEST_FIELD",
		StatusCode: StatusInternalServerError,
		Message:    MsgBankDeleteRequestFailed,
		Type:       "error",
	}

	ErrorBankImageMissingOrInvalid = ResponseCode{
		Code:       "ERROR_BANK_IMAGE_MISSING_OR_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgBankImageRequiredOrMissing,
		Type:       "error",
	}

	ErrorWalletImageMissingOrInvalid = ResponseCode{
		Code:       "ERROR_WALLET_IMAGE_MISSING_OR_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgWalletImageRequiredOrMissing,
		Type:       "error",
	}

	ErrorInvalidFormat = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidFormat,
		Type:       "error",
	}

	ErrorInvalidInputParameter = ResponseCode{
		Code:       "ERROR_INVALID_INPUT_PARAMETER",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidInputParameter,
		Type:       "error",
	}
	ErrorActionNameIsRequired = ResponseCode{
		Code:       "ERROR_INVALID_ACTION_NAME",
		StatusCode: StatusBadRequest,
		Message:    MsgActionNameIsRequired,
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
	ErrorExistEmail = ResponseCode{
		Code:       "ERROR_EXIST_EMAIL",
		StatusCode: StatusBadRequest,
		Message:    MsgExistEmail,
		Type:       "error",
	}

	ErrorExistPhoneNumber = ResponseCode{
		Code:       "ERROR_EXIST_PHONE_NUMBER",
		StatusCode: StatusBadRequest,
		Message:    MsgExistPhoneNumber,
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

	ErrorMissingFile = ResponseCode{
		Code:       "ERROR_MISSING_FILE",
		StatusCode: StatusBadRequest,
		Message:    MsgMissingFile,
		Type:       "error",
	}

	ErrorFileParseFailed = ResponseCode{
		Code:       "ERROR_FILE_PARSE_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgFileParseFailed,
		Type:       "error",
	}

	ErrorInvalidAction = ResponseCode{
		Code:       "ERROR_INVALID_ACTION",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidAction,
		Type:       "error",
	}
	ErrorInvalidActionFormat = ResponseCode{
		Code:       "ERROR_INVALID_ACTION_FORMAT",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidActionFormat,
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

	ErrorTransactionIdentifierRequired = ResponseCode{
		Code:       "ERROR_TRANSACTION_IDENTIFIER_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgTransactionIdentifierRequired,
		Type:       "error",
	}

	ErrorFailedToBeingTransaction = ResponseCode{
		Code:       "ERROR_FAILED_TO_BE_TRANSACTION",
		StatusCode: StatusInternalServerError,
		Message:    MsgFailedToBeingTransaction,
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

	ErrorResourceNotFound = ResponseCode{
		Code:       "ERROR_RESOURCE_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgResourceNotFound,
		Type:       "error",
	}
	ErrorActionNameAlreadyExists = ResponseCode{
		Code:       "ERROR_ACTION_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgActionNameAlreadyExists,
		Type:       "error",
	}
	ErrorOnDisablingExistingDeviceControl = ResponseCode{
		Code:       "ERROR_ON_DISABLING_EXISTING_DEVICE_CONTROL",
		StatusCode: StatusInternalServerError,
		Message:    MsgOnDisablingExistingDeviceControl,
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

	// Password Rule
	ErrorNoPasswordRule = ResponseCode{
		Code:       "ERROR_PASSWORD_RULE_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgPasswordRuleNotFound,
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

	ErrorProductCodeNotFound = ResponseCode{
		Code:       "ERROR_PRODCUT_CODE_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgProductCodeNotFound,
		Type:       "error",
	}

	ErrorDuplicateProductName = ResponseCode{
		Code:       "ERROR_DUPLICATE_PRODUCT_NAME",
		StatusCode: StatusConflict,
		Message:    MsgDuplicateProductName,
		Type:       "error",
	}

	ErrorDuplicateCBEProductCode = ResponseCode{
		Code:       "ERROR_DUPLICATE_CBE_PRODUCT_CODE",
		StatusCode: StatusConflict,
		Message:    MsgDuplicateCBEProductCode,
		Type:       "error",
	}

	ErrorDuplicateCBEIFBProductCode = ResponseCode{
		Code:       "ERROR_DUPLICATE_CBE_IFB_PRODUCT_CODE",
		StatusCode: StatusConflict,
		Message:    MsgDuplicateCBEIFBProductCode,
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
	ErrorMiniAppMerchantNotFound = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "MiniApp merchant not found",
		Type:       "error",
	}
	ErrorMiniAppMerchantDisabled = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    "MiniApp merchant is Disabled",
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

	ErrorDeviceVersionAlreadyExists = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgDeviceVersionAlreadyExists,
		Type:       "error",
	}
	ErrorDeviceVersionAlreadyEnabled = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgDeviceVersionAlreadyEnabled,
		Type:       "error",
	}
	ErrorDeviceVersionAlreadyDisabled = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgDeviceVersionAlreadyDisabled,
		Type:       "error",
	}
	ErrorDeviceVersionNotFound = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgDeviceVersionNotFound,
		Type:       "error",
	}
	ErrorDeviceVersionUpdateFailed = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDeviceVersionUpdateFailed,
		Type:       "error",
	}
	ErrorDeviceVersionDeleteFailed = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_DELETE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDeviceVersionDeleteFailed,
		Type:       "error",
	}
	ErrorDeviceVersionEnableFailed = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_ENABLE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDeviceVersionEnableFailed,
		Type:       "error",
	}
	ErrorDeviceVersionDisableFailed = ResponseCode{
		Code:       "ERROR_DEVICE_VERSION_DISABLE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDeviceVersionDisableFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantFetchPermissionsFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_FETCH_PERMISSIONS_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantFetchPermissionsFailed,
		Type:       "error",
	}
	ErrorMiniAppMerchantDeleteFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_DELETE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantDeleteFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantEnableFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_ENABLE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantEnableFailed,
		Type:       "error",
	}

	ErrorMiniAppMerchantDisableFailed = ResponseCode{
		Code:       "ERROR_MINI_APP_MERCHANT_DISABLE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgMiniAppMerchantDisableFailed,
		Type:       "error",
	}

	ErrorCPSRoleAlreadyEnabled = ResponseCode{
		Code:       "ERROR_CPS_ROLE_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgCpsRoleAlreadyEnabled,
		Type:       "error",
	}

	ErrorCPSRoleAlreadyDisabled = ResponseCode{
		Code:       "ERROR_CPS_ROLE_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgCpsRoleAlreadyDisabled,
		Type:       "error",
	}

	// Notification related error response codes
	ErrorNotificationMapFailed = ResponseCode{
		Code:       "ERROR_NOTIFICATION_MAP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgNotificationMapFailed,
		Type:       "error",
	}
	ErrorNotificationAlreadyExists = ResponseCode{
		Code:       "ERROR_NOTIFICATION_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgNotificationAlreadyExists,
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
	ErrorNotificationDeleteFailed = ResponseCode{
		Code:       "ERROR_NOTIFICATION_DELETE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgNotificationDeleteFailed,
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

	CustomerSegmentationCreationSubmittedSuccessfully = ResponseCode{
		Code:       "ERROR_CUSTOMER_SEGMENTATION_CREATION_SUBMITTED",
		StatusCode: StatusOK,
		Message:    "Customer segmentation creation request submitted successfully",
		Type:       "error",
	}
	CustomerSegmentationUpdateSubmittedSuccessfully = ResponseCode{
		Code:       "ERROR_CUSTOMER_SEGMENTATION_UPDATE_SUBMITTED",
		StatusCode: StatusOK,
		Message:    "Customer segmentation update request submitted successfully",
		Type:       "error",
	}
	CustomerSegmentationFetchedSuccessfully = ResponseCode{
		Code:       "CUSTOMER_SEGMENTATIONS_FETCHED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    "Customer segmentations fetched successfully",
		Type:       "error",
	}
	CustomerSegmentationDeleteddSuccessfully = ResponseCode{
		Code:       "CUSTOMER_SEGMENTATIONS_DELETED_SUCCESSFULLY",
		StatusCode: StatusOK,
		Message:    "Customer segmentations deleted request submitted successfully",
		Type:       "error",
	}

	// Donation related error response codes
	ErrorDonationCategoryLookupFailed = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_LOOKUP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationCategoryLookupFailed,
		Type:       "error",
	}

	ErrorDonationCategoryIDRequired = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationCategoryIdRequired,
		Type:       "error",
	}

	ErrorDonationCompanyIdRequired = ResponseCode{
		Code:       "ERROR_DONATION_COMPANY_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationCompanyIdRequired,
		Type:       "error",
	}
	ErrorAccountNumberRequired = ResponseCode{
		Code:       "ERROR_ACCOUNT_NUMBER_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgAccountNumberRequired,
		Type:       "error",
	}
	ErrorCannotGetRole = ResponseCode{
		Code:       "ERROR_CANNOT_APPROVE_ACTION",
		StatusCode: StatusBadRequest,
		Message:    MsgCannotGetRole,
		Type:       "error",
	}

	ErrorDonationCompanyLookupFailed = ResponseCode{
		Code:       "ERROR_DONATION_COMPANY_LOOKUP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationCompanyLookupFailed,
		Type:       "error",
	}
	ErrorDonationAlreadyEnabled = ResponseCode{
		Code:       "ERROR_DONATION_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationAlreadyEnabled,
		Type:       "error",
	}
	ErrorDonationAlreadyDisabled = ResponseCode{
		Code:       "ERROR_DONATION_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationAlreadyDisabled,
		Type:       "error",
	}

	ErrorCompanyNameAlreadyExists = ResponseCode{
		Code:       "ERROR_COMPANY_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgCompanyNameAlreadyExists,
		Type:       "error",
	}
	ErrorCompanyCodeAlreadyExists = ResponseCode{
		Code:       "ERROR_COMPANY_CODE_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgCompanyCodeAlreadyExists,
		Type:       "error",
	}

	ErrorAccountNumberAlreadyExists = ResponseCode{
		Code:       "ERROR_ACCOUNT_NUMBER_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgAccountNumberAlreadyExists,
		Type:       "error",
	}
	ErrorEmailAlreadyExist = ResponseCode{
		Code:       "ERROR_EMAIL_ALREADY_EXIST",
		StatusCode: StatusBadRequest,
		Message:    MsgEmailAlreadyExists,
		Type:       "error",
	}
	ErrorPhonenumberAlreadyExist = ResponseCode{
		Code:       "ERROR_PHONENUMBER_ALREADY_EXIST",
		StatusCode: StatusBadRequest,
		Message:    MsgPhonenumberAlreadyExists,
		Type:       "error",
	}
	ErrorCodeAlreadyExist = ResponseCode{
		Code:       "ERROR_CODE_ALREADY_EXIST",
		StatusCode: StatusBadRequest,
		Message:    MsgCodeAlreadyExists,
		Type:       "error",
	}

	ErrorCustomerBlockedPermanently = ResponseCode{
		Code:       "ERROR_CUSTOMER_PERMANENTLY_BLOCKED",
		StatusCode: StatusBadRequest,
		Message:    MsgCustomerPermanentlyDisabled,
		Type:       "error",
	}

	ErrorLogoIsRequired = ResponseCode{
		Code:       "ERROR_LOGO_IS_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgLogoIsRequired,
		Type:       "error",
	}

	ErrorAccountNumberNotFound = ResponseCode{
		Code:       "ERROR_ACCOUNT_NUMBER_NOT_FOUND",
		StatusCode: StatusBadRequest,
		Message:    MsgAccountNotFound,
		Type:       "error",
	}
	ErrorAccountNumberNotActive = ResponseCode{
		Code:       "ERROR_ACCOUNT_NUMBER_NOT_ACTIVE",
		StatusCode: StatusBadRequest,
		Message:    MsgAccountNotActive,
		Type:       "error",
	}
	ErrorAccountNumberValidationFailed = ResponseCode{
		Code:       "ERROR_ACCOUNT_NUMBER_VALIDATION_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgAccountNumberValidationFailed,
		Type:       "error",
	}

	ErrorDonationCategoryNameDuplicated = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_NAME_DUPLICATED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationCategoryNameDuplicated,
		Type:       "error",
	}

	ErrorNoChangesToUpdate = ResponseCode{
		Code:       "ERROR_NO_CHANGES_TO_UPDATE",
		StatusCode: StatusBadRequest,
		Message:    MsgNoChangesToUpdate,
		Type:       "error",
	}

	ErrorBudgetCategoryAlreadyEnabled = ResponseCode{
		Code:       "ERROR_BUDGET_CATEGORY_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgBudgetCategoryAlreadyEnabled,
		Type:       "error",
	}
	ErrorBudgetCategoryAlreadyDisabled = ResponseCode{
		Code:       "ERROR_BUDGET_CATEGORY_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgBudgetCategoryAlreadyDisable,
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
	ErrorBankUpdateFailed = ResponseCode{
		Code:       "ERROR_BANK_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgBankUpdateFailed,
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
	ErrorPermissionGroupNotFound = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgPermissionGroupNotFound,
		Type:       "error",
	}
	ErrorPermissionGroupValidationFailed = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_VALIDATION_FAILED",
		StatusCode: StatusBadRequest,
		Message:    MsgPermissionGroupValidationFailed,
		Type:       "error",
	}
	ErrorPermissionCategoryNotFound = ResponseCode{
		Code:       "ERROR_PERMISSION_CATEGORY_NOT_FOUND",
		StatusCode: StatusBadRequest,
		Message:    MsgPermissionCategoryNotFound,
		Type:       "error",
	}
	ErrorPermissionGroupRequired = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgPermissionGroupRequired,
		Type:       "error",
	}

	ErrorPermissionGroupAlreadyExists = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgPermissionGroupAlreadyExists,
		Type:       "error",
	}
	ErrorPermissionCatagoryNotFound = ResponseCode{
		Code:       "ERROR_PERMISSION_CATAGORY_NOT_FOUND",
		StatusCode: StatusConflict,
		Message:    MsgPermissionCatagoryNotFound,
		Type:       "error",
	}

	ErrorPermissionGroupRequestUpdateFailed = ResponseCode{
		Code:       "ERROR_PERMISSION_GROUP_REQUEST_UPDATE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgPermissionGroupRequestUpdateFailed,
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

	ErrorDonationTitleDuplicated = ResponseCode{
		Code:       "ERROR_DONATION_TITLE_DUPLICATED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationTitleDuplicated,
		Type:       "error",
	}

	ErrorDonationCategoryNotFound = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgDonationCategoryNotFound,
		Type:       "error",
	}

	ErrorDonationCompanyNotFound = ResponseCode{
		Code:       "ERROR_DONATION_COMPANY_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgDonationCompanyNotFound,
		Type:       "error",
	}
	ErrorCompanyIsNotEnabled = ResponseCode{
		Code:       "ERROR_DONATION_COMPANY_NOT_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationCompanyNotEnabled,
		Type:       "error",
	}
	ErrorCategoryIsNotEnabled = ResponseCode{
		Code:       "ERROR_DONATION_CATEGORY_NOT_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgDonationCategoryNotEnabled,
		Type:       "error",
	}

	ErrorImageRequired = ResponseCode{
		Code:       "ERROR_IMAGE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgImageRequired,
		Type:       "error",
	}

	ErrorDonationLookupFailed = ResponseCode{
		Code:       "ERROR_DONATION_LOOKUP_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgDonationLookupFailed,
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
	ErrorNotificationCreated = ResponseCode{
		Code:       "ERROR_NOTIFICATION_CREATED",
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

	// User Unlink error response codes
	ErrorCustomerDoesNotHaveLinkedAccount = ResponseCode{
		Code:       "ERROR_CUSTOMER_DOES_NOT_HAVE_ACCOUNT",
		StatusCode: StatusInternalServerError,
		Message:    MsgUserNotHaveLinkedAccount,
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

	ErrorDescriptionLength30To100 = ResponseCode{
		Code:       "DESCRIPTION_LENGTH_30_TO_100",
		StatusCode: StatusBadRequest,
		Message:    "Description length must be between 30 and 100 characters",
		Type:       "error",
	}

	ErrorTitleLength3To20 = ResponseCode{
		Code:       "TITLE_LENGTH_3_TO_20",
		StatusCode: StatusBadRequest,
		Message:    "Title length must be between 3 and 20 characters",
		Type:       "error",
	}

	ErrorAdvertAlreadyDisabled = ResponseCode{
		Code:       "ERROR_ADVERT_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgAdvertAlreadyDisabled,
		Type:       "error",
	}

	ErrorAvatarAlreadyDisabled = ResponseCode{
		Code:       "ERROR_AVATAR_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgAvatarAlreadyDisabled,
		Type:       "error",
	}

	ErrorAvatarAlreadyEnabled = ResponseCode{
		Code:       "ERROR_AVATAR_ALREADY_Enabled",
		StatusCode: StatusBadRequest,
		Message:    MsgAvatarAlreadyEnabled,
		Type:       "error",
	}

	// Account Validation Service related error response codes
	ErrorValidationRuleApproved = ResponseCode{
		Code:       "ERROR_VALIDATION_RULE_APPROVED",
		StatusCode: StatusOK,
		Message:    MsgValidationRuleApprovedSuccess,
		Type:       "error",
	}

	ErrorGetOneBank = ResponseCode{
		Code:       "ERROR_GET_ONE_BANK",
		StatusCode: StatusInternalServerError,
		Message:    msgGetOneBankFailed,
		Type:       "error",
	}

	ErrorBankDisableRequest = ResponseCode{
		Code:       "ERROR_BANK_DISABLE_REQUEST",
		StatusCode: StatusInternalServerError,
		Message:    MsgBankDisableRequestFailed,
		Type:       "error",
	}

	ErrorBankEnableRequestFailed = ResponseCode{
		Code:       "ERROR_BANK_ENABLE_REQUEST",
		StatusCode: StatusInternalServerError,
		Message:    MsgBankEnableRequestFailed,
		Type:       "error",
	}

	ErrorInvalidInputParameters = ResponseCode{
		Code:       "ERROR_INVAErrorMissingOrInvalidImage,LID_INPUT_PARAMETERS",
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

	ErrorPendingActionExists = ResponseCode{
		Code:       "ERROR_PENDING_ACTION_EXISTS",
		StatusCode: StatusInternalServerError,
		Message:    MsgPendingActionExists,
		Type:       "error",
	}
	ErrorInvalidRequestBody = ResponseCode{
		Code:       "ERROR_INVALID_REQUEST_BODY",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestBody,
		Type:       "error",
	}

	ErrorInvalidPaginationParams = ResponseCode{
		Code:       "ERROR_INVALID_PAGINATION_PARAMS",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidPaginationParams,
		Type:       "error",
	}

	ErrorInvalidDistrict = ResponseCode{
		Code:       "ERROR_INVALID_DISTRICT",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidDistrict,
		Type:       "error",
	}

	ErrorInvalidRegion = ResponseCode{
		Code:       "ERROR_INVALID_REGION",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRegion,
		Type:       "error",
	}

	ErrorInvalidDistrictOrRegionCodeLength = ResponseCode{
		Code:       "ERROR_INVALID_DISTRICT_OR_REGION_CODE_LENGTH",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidDistrictOrRegionCodeLength,
		Type:       "error",
	}

	ErrorBranchCodeRequired = ResponseCode{
		Code:       "ERROR_BRANCH_CODE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgBranchCodeRequired,
		Type:       "error",
	}

	ErrorBranchAlreadyEnabled = ResponseCode{
		Code:       "ERROR_BRANCH_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgBranchAlreadyEnabled,
		Type:       "error",
	}

	ErrorBranchAlreadyDisabled = ResponseCode{
		Code:       "ERROR_BRANCH_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgBranchAlreadyDisabled,
		Type:       "error",
	}

	ErrorDistrictCodeRequired = ResponseCode{
		Code:       "ERROR_DISTRICT_CODE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgDistrictCodeRequired,
		Type:       "error",
	}

	ErrorDistrictAlreadyEnabled = ResponseCode{
		Code:       "ERROR_DISTRICT_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgDistrictAlreadyEnabled,
		Type:       "error",
	}

	ErrorDistrictAlreadyDisabled = ResponseCode{
		Code:       "ERROR_DISTRICT_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgDistrictAlreadyDisabled,
		Type:       "error",
	}

	ErrorRegionCodeRequired = ResponseCode{
		Code:       "ERROR_REGION_CODE_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgRegionCodeRequired,
		Type:       "error",
	}

	ErrorRegionAlreadyEnabled = ResponseCode{
		Code:       "ERROR_REGION_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgRegionAlreadyEnabled,
		Type:       "error",
	}
	ErrorRegionAlreadyDisabled = ResponseCode{
		Code:       "ERROR_REGION_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgRegionAlreadyDisabled,
		Type:       "error",
	}

	ErrorCityCodeRequired = ResponseCode{
		Code:       "ERROR_CITY_CODE_REQUIRED",
		StatusCode: StatusConflict,
		Message:    MsgCityCodeRequired,
		Type:       "error",
	}

	ErrorCityAlreadyEnabled = ResponseCode{
		Code:       "ERROR_CITY_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgCityAlreadyEnabled,
		Type:       "error",
	}

	ErrorCityAlreadyDisabled = ResponseCode{
		Code:       "ERROR_CITY_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgCityAlreadyDisabled,
		Type:       "error",
	}

	ErrorBranchNotFound = ResponseCode{
		Code:       "ERROR_BRANCH_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgBranchNotFound,
		Type:       "error",
	}

	ErrorDistrictNotFound = ResponseCode{
		Code:       "ERROR_DISTRICT_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgDistrictNotFound,
		Type:       "error",
	}
	ErrorCannotEnableDistrict = ResponseCode{
		Code:       "ERROR_CANNOT_ENABLE_DISTRICT",
		StatusCode: StatusBadRequest,
		Message:    MsgCannotEnableDistrict,
		Type:       "error",
	}
	ErrorCannotEnablCity = ResponseCode{
		Code:       "ERROR_CANNOT_ENABLE_CITY",
		StatusCode: StatusBadRequest,
		Message:    MsgCannotEnableCity,
		Type:       "error",
	}

	ErrorCannotEnableBranch = ResponseCode{
		Code:       "ERROR_CANNOT_ENABLE_BRANCH",
		StatusCode: StatusBadRequest,
		Message:    MsgCannotEnableBranch,
		Type:       "error",
	}

	ErrorRegionNotFound = ResponseCode{
		Code:       "ERROR_REGION_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgRegionNotFound,
		Type:       "error",
	}

	ErrorCityNotFound = ResponseCode{
		Code:       "ERROR_CITY_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgCityNotFound,
		Type:       "error",
	}

	ErrorDuplicateAction = ResponseCode{
		Code:       "ERROR_DUPLICATE_ACTION",
		StatusCode: StatusBadRequest,
		Message:    MsgDuplicateAction,
		Type:       "error",
	}

	ErrorAlreadyEnabled = ResponseCode{
		Code:       "ERROR_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgAlreadyEnabled,
		Type:       "error",
	}
	ErrorFailedToUpdateDonation = ResponseCode{
		Code:       "ERROR_FAILED_TO_UPDATE_DONATION_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    "Failed to update donations",
		Type:       "error",
	}

	ErrorAlreadyDisabled = ResponseCode{
		Code:       "ERROR_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    MsgAlreadyDisabled,
		Type:       "error",
	}

	ErrorOneOrMoreInvalidCodes = ResponseCode{
		Code:       "ERROR_ONE_OR_MORE_INVALID_CODES",
		StatusCode: StatusBadRequest,
		Message:    MsgOneOrMoreInvalidCodes,
		Type:       "error",
	}

	// department related error
	ErrorInvalidFormatForDepartmentName = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT_FOR_DEPARTMENT_NAME",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestDepartmentName,
		Type:       "error",
	}
	ErrorInvalidFormatForDepartmentPortalCards = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT_FOR_DEPARTMENT_PORTAL_CARDS",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestDepartmentPortalCards,
		Type:       "error",
	}
	ErrorInvalidFormatForDepartmentPermissionGroups = ResponseCode{
		Code:       "ERROR_INVALID_FORMAT_FOR_DEPARTMENT_PERMISSION_GROUPS",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidRequestDepartmentPermissionGroups,
		Type:       "error",
	}
	ErrorInvalidDepartmentPermissionGroup = ResponseCode{
		Code:       "ERROR_INVALID_DEPARTMENT_PERMISSION_GROUP",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidDepartmentPermissionGroup,
		Type:       "error",
	}
	ErrorInvalidDepartmentPortalCard = ResponseCode{
		Code:       "ERROR_INVALID_DEPARTMENT_PORTAL_CARD",
		StatusCode: StatusBadRequest,
		Message:    MsgInvalidDepartmentPortalCard,
		Type:       "error",
	}
	ErrorDepartmentCreateRequest = ResponseCode{
		Code:       "ERROR_DEPARTMENT_CREATION_REQUEST",
		StatusCode: StatusExpectationFailed,
		Message:    MsgDepartmentCreateRequestFail,
		Type:       "error",
	}
	ErrorDepartmentNotFound = ResponseCode{
		Code:       "ERROR_DEPARTMENT_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgDepartmentNotFound,
		Type:       "error",
	}

	ErrorDepartmentAlreadyEnabled = ResponseCode{
		Code:       "ERROR_DEPARTMENT_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgDepartmentAlreadyEnabled,
		Type:       "error",
	}
	ErrorDepartmentAlreadyDisabled = ResponseCode{
		Code:       "ERROR_DEPARTMENT_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgDepartmentAlreadyDisabled,
		Type:       "error",
	}
	ErrorDepartmentInvalidID = ResponseCode{
		Code:       "ERROR_DEPARTMENT_INVALID_ID",
		StatusCode: StatusBadRequest,
		Message:    MsgDepartmentInvalidID,
		Type:       "error",
	}
	ErrorDepartmentWithNameAlreadyExists = ResponseCode{
		Code:       "ERROR_DEPARTMENT_WITH_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgDepartmentWithNameAlreadyExists,
		Type:       "error",
	}
	// fayda
	ErrorFaydaUserAccountEnabled = ResponseCode{
		Code:       "FAYDA_USER_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgUserFaydaAccountAlreadyEnabled,
		Type:       "error",
	}
	ErrorFaydaUserAccountDisabled = ResponseCode{
		Code:       "FAYDA_USER_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgUserFaydaAccountAlreadyDisabled,
		Type:       "error",
	}
	ErrorNotFaydaUser = ResponseCode{
		Code:       "USER_IS_NOT_FAYDA_USER",
		StatusCode: StatusNotFound,
		Message:    MsgNotFaydaUser,
		Type:       "error",
	}
	ErrorBucketNotFound = ResponseCode{
		Code:       "BUCKET NOT FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgBucketNotFound,
		Type:       "error",
	}
	ErrorFailedToBucket = ResponseCode{
		Code:       "ERROR_FAILED_TO_STORE_IN_BUCKET",
		StatusCode: StatusInternalServerError,
		Message:    "Failed to get store in bucket",
		Type:       "error",
	}

	ErrorKeyRequiredForBulkService = ResponseCode{
		Code:       "ERROR_KEY_REQUIRED_FOR_BULK_SERVICE",
		StatusCode: StatusBadRequest,
		Message:    "Key is required for bulk service operation",
		Type:       "error",
	}

	ErrorFailedToGetCustomerDetail = ResponseCode{
		Code:       "ERROR_FAILED_TO_GET_CUSTOMER_DETAIL",
		StatusCode: StatusInternalServerError,
		Message:    "Failed to get customer detail",
		Type:       "error",
	}
	UserNotFoundWithGivenID = ResponseCode{
		Code:       "USER_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "user not found for given id",
		Type:       "error",
	}

	ErrConfigIsEmpty = ResponseCode{
		Code:       "ERROR_CONFIG_IS_EMPTY",
		StatusCode: StatusInternalServerError,
		Message:    "Config is empty",
		Type:       "error",
	}
	ErrMarshalingData = ResponseCode{
		Code:       "ERROR_MARSHALING_DATA",
		StatusCode: StatusInternalServerError,
		Message:    "Error while marshaling data",
		Type:       "error",
	}
	ErrInvalidKeyOrIv = ResponseCode{
		Code:       "ERROR_INVALID_KEY_OR_IV",
		StatusCode: StatusInternalServerError,
		Message:    "Invalid key or iv",
		Type:       "error",
	}

	SuccessCustomerDetailSuccessfullyFetched = ResponseCode{
		Code:       "SUCCESS_CUSTOMER_DETAIL_FETCHED",
		StatusCode: StatusOK,
		Message:    "Customer detail(s) fetched successfully",
		Type:       "success",
	}

	CustomerEnableRequestCreatedSuccessfully = ResponseCode{
		Code:       "SUCCESS_CUSTOMER_ENABLE_REQUEST_CREATED",
		StatusCode: StatusOK,
		Message:    "Customer Enable request created Successfully",
		Type:       "success",
	}

	CustomerEnableRequestSessionCreatedSuccessfully = ResponseCode{
		Code:       "SUCCESS_CUSTOMER_ENABLE_REQUEST_SESSION_CREATED",
		StatusCode: StatusOK,
		Message:    "Customer Enable request session created Successfully, Verify Otp to continue",
		Type:       "success",
	}

	CustomerDisableRequestCreatedSuccessfully = ResponseCode{
		Code:       "SUCCESS_CUSTOMER_DISABLE_REQUEST_CREATED",
		StatusCode: StatusOK,
		Message:    "Customer Disable request created Successfully",
		Type:       "success",
	}

	ErrorCustomerAlreadyEnabled = ResponseCode{
		Code:       "ERROR_CUSTOMER_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    "Customer is already enabled",
		Type:       "error",
	}

	ErrorCustomerAlreadyDisabled = ResponseCode{
		Code:       "ERROR_CUSTOMER_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    "Customer is already disabled",
		Type:       "error",
	}

	ErrorIdNotSetOnQueryParam = ResponseCode{
		Code:       "ERROR_ID_NOT_SET_ON_QUERY_PARAM",
		StatusCode: StatusBadRequest,
		Message:    "ID not set on query parameter",
		Type:       "error",
	}

	CustomerDetailSuccessfullyFetched = ResponseCode{
		Code:       "SUCCESS_CUSTOMER_DETAIL_FETCHED",
		StatusCode: StatusOK,
		Message:    "Customer detail fetched successfully",
		Type:       "success",
	}

	FaydaCustomerApprovalRequestSent = ResponseCode{
		Code:       "SUCCESS_FAYDA_APPROVAL_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    "Fayda approval request sent successfuly",
		Type:       "success",
	}

	ErrorFailedToGetBlockedCustomer = ResponseCode{
		Code:       "ERROR_FAILED_TO_GET_BLOCKED_CUSTOMER",
		StatusCode: StatusInternalServerError,
		Message:    "Failed to get blocked customer",
		Type:       "error",
	}

	SuccessFullyFetchBlockCustomer = ResponseCode{
		Code:       "SUCCESSFULLY_FETCH_BLOCKED_CUSTOMER",
		StatusCode: StatusOK,
		Message:    "Successfully fetched blocked customer(s)",
		Type:       "success",
	}

	UnableToFetchBulkService = ResponseCode{
		Code:       "ERROR_UNABLE_TO_FETCH_BULK_SERVICE",
		StatusCode: StatusInternalServerError,
		Message:    "Unable to fetch bulk service",
		Type:       "error",
	}

	BulkServiceFetchSuccessfully = ResponseCode{
		Code:       "SUCCESS_BULK_SERVICE_FETCHED",
		StatusCode: StatusOK,
		Message:    "Bulk service(s) fetched successfully",
		Type:       "success",
	}

	BulkServiceEnableRequestSuccess = ResponseCode{
		Code:       "SUCCESS_BULK_SERVICE_ENABLE_REQUEST",
		StatusCode: StatusOK,
		Message:    "Bulk service enable request processed successfully",
		Type:       "success",
	}

	BulkServiceDisableRequestSuccess = ResponseCode{
		Code:       "SUCCESS_BULK_SERVICE_DISABLE_REQUEST",
		StatusCode: StatusOK,
		Message:    "Bulk service disable request processed successfully",
		Type:       "success",
	}
	ErrorSingleMaxTransferCannotBeLessOrEqualToMinAmount = ResponseCode{
		Code:       "ERROR_MAX_TRANSFER_UPDATE_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    "single max transfer can not be less or equal to min amount",
		Type:       "error",
	}

	ErrorTotalMaxTransferCannotBeLessExistTransfers = ResponseCode{
		Code:       "ERROR_TOTAL_CAP_UPDATE_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    "total max transfer can not be less from existing service transfers",
		Type:       "error",
	}
	ErrorMinAmountCanNotBeGreaterThanCap = ResponseCode{
		Code:       "ERROR_MINIMUM_TRANSFER_UPDATE_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    "minimum transfer can not be greater from existing transfer caps",
		Type:       "error",
	}
	ErrorUserNotFaydaRegistered = ResponseCode{
		Code:       "ERROR_NOT_FAYDA_USER",
		StatusCode: StatusBadRequest,
		Message:    "this user is not fayda user",
		Type:       "error",
	}
	ErrorSingleTransferCanNotBeGreaterThanCap = ResponseCode{
		Code:       "ERROR_SINGLE_TRANSFER_UPDATE_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    "single transfer can not be greater from  total cap",
		Type:       "error",
	}
	ErrorMinAmountCanNotBeGreaterThanTotal = ResponseCode{
		Code:       "ERROR_MINIMUM_TRANSFER_UPDATE_REQUEST",
		StatusCode: StatusBadRequest,
		Message:    "minimum transfer can not be greater from total cap",
		Type:       "error",
	}

	ErrorBulkServiceAlreadyEnabled = ResponseCode{
		Code:       "ERROR_BULK_SERVICE_ALREADY_ENABLED",
		StatusCode: StatusBadRequest,
		Message:    "One or more bulk services are already enabled",
		Type:       "error",
	}

	ErrorBulkServiceAlreadyDisabled = ResponseCode{
		Code:       "ERROR_BULK_SERVICE_ALREADY_DISABLED",
		StatusCode: StatusBadRequest,
		Message:    "One or more bulk services are already enabled",
		Type:       "error",
	}

	ErrorInvalidBulkServiceKey = ResponseCode{
		Code:       "ERROR_INVALID_BULK_SERVICE_KEY",
		StatusCode: StatusBadRequest,
		Message:    "One or more provided bulk service keys is/are invalid.",
		Type:       "error",
	}
	ErrorInvalidRequiredAction = ResponseCode{
		Code:       "ERROR_INVALID_REQUIRED_ACTION",
		StatusCode: StatusBadRequest,
		Message:    "The required action is invalid or missing.",
		Type:       "error",
	}

	ErrorFailToUpdateChild = ResponseCode{
		Code:       "ERROR_FAIL_TO_UPDATE_CHILD",
		StatusCode: StatusInternalServerError,
		Message:    "Failed to update child record.",
		Type:       "error",
	}

	ErrorFailToUpdateParent = ResponseCode{
		Code:       "ERROR_FAIL_TO_UPDATE_PARENT",
		StatusCode: StatusInternalServerError,
		Message:    "Failed to update parent record.",
		Type:       "error",
	}
	ErrorNoChangesDetected = ResponseCode{
		Code:       "ERROR_NO_CHANGES_DETECTED",
		StatusCode: StatusBadRequest,
		Message:    "no changes detected to update",
		Type:       "error",
	}
	ErrorVaultGroupCategoryNotFound = ResponseCode{
		Code:       "ERROR_VAULT_GROUP_CATEGORY_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "Vault group category not found.",
		Type:       "error",
	}
	ErrorGroupVaultNotFound = ResponseCode{
		Code:       "ERROR_GROUP_VAULT_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    "Group vault not found",
		Type:       "error",
	}
	ErrorCannotDeleteActiveVaultGroupCategory = ResponseCode{
		Code:       "ERROR_CANNOT_DELETE_ACTIVE_VAULT_GROUP_CATEGORY",
		StatusCode: StatusBadRequest,
		Message:    "Cannot delete an active vault group category.",
		Type:       "error",
	}
	ErrorVaultGroupCategooryAlreadyDeleted = ResponseCode{
		Code:       "ERROR_VAULT_GROUP_CATEGORY_ALREADY_DELETED",
		StatusCode: StatusBadRequest,
		Message:    "Vault group category is already deleted.",
		Type:       "error",
	}
	ErrorVaultCoverImageMissedOrInvalid = ResponseCode{
		Code:       "ERROR_VAULT_COVER_IMAGE_MISSED_OR_INVALID",
		StatusCode: StatusBadRequest,
		Message:    "Vault category cover image missed or invalid",
		Type:       "error",
	}

	ErrorDuplicateBankProduct = ResponseCode{
		Code:       "ERROR_DUPLICATE_BANK_PRODUCT",
		StatusCode: StatusBadRequest,
		Message:    MsgDuplicatebankProduct,
		Type:       "error",
	}

	ErrorNoBankProductFound = ResponseCode{
		Code:       "ERROR_NO_BANK_PRODUCT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgNoBankProductFound,
		Type:       "error",
	}

	ErrorCannotDeletedBankProduct = ResponseCode{
		Code:       "ERROR_CANNOT_DELETED_BANK_PRODUCT",
		StatusCode: StatusBadRequest,
		Message:    MsgCannotDeleteBankProduct,
		Type:       "error",
	}

	ErrorCannotEnableOrDisable = ResponseCode{
		Code:       "ERROR_CANNOT_ENABLE_OR_DISABLE",
		StatusCode: StatusBadRequest,
		Message:    MsgcannotEnableOrDisableDeletedBankProduct,
		Type:       "error",
	}

	ErrorBankVaultProductAlreadyDeleted = ResponseCode{
		Code:       "ERROR_BANK_VAULT_PRODUCT_ALREADY_DELETED",
		StatusCode: StatusBadRequest,
		Message:    MsgBankVaultProductAlreadyDeleted,
		Type:       "error",
	}

	ErrorDuplicateGroupVaultCategory = ResponseCode{
		Code:       "ERROR_DUPLICATE_GROUP_VAULT_CATEGORY",
		StatusCode: StatusBadRequest,
		Message:    MsgDuplicateGroupVaultCategory,
		Type:       "error",
	}

	ErrorTopupNameRequired = ResponseCode{
		Code:       "ERROR_TOPUP_NAME_REQUIRED",
		StatusCode: 400,
		Message:    "Topup name is required",
		Type:       "error",
	}

	ErrorTopupCodeRequired = ResponseCode{
		Code:       "ERROR_TOPUP_CODE_REQUIRED",
		StatusCode: 400,
		Message:    "Topup code is required",
		Type:       "error",
	}

	ErrorTopupAvatarRequired = ResponseCode{
		Code:       "ERROR_TOPUP_AVATAR_REQUIRED",
		StatusCode: 400,
		Message:    "Topup avatar is required",
		Type:       "error",
	}
	ErrorTopupAvatarInvalid = ResponseCode{
		Code:       "ERROR_TOPUP_AVATAR_INVALID",
		StatusCode: 400,
		Message:    "Invalid topup avatar",
		Type:       "error",
	}

	ErrorTopupAvatarTooLarge = ResponseCode{
		Code:       "ERROR_TOPUP_AVATAR_TOO_LARGE",
		StatusCode: 400,
		Message:    "Topup avatar file size exceeds the limit",
		Type:       "error",
	}

	ErrorTopupAvatarInvalidType = ResponseCode{
		Code:       "ERROR_TOPUP_AVATAR_INVALID_TYPE",
		StatusCode: 400,
		Message:    "Topup avatar must be of type jpeg, png, gif, or webp",
		Type:       "error",
	}

	ErrorTopupServiceOption = ResponseCode{
		Code:       "ERROR_TOPUP_SERVICE_OPTION_INVALID_VALUES",
		StatusCode: 400,
		Message:    "At Least one of the three service options should be enabled(self,other,agent)",
		Type:       "error",
	}

	ErrorTopupImageMissingOrInvalid = ResponseCode{
		Code:       "ERROR_Topup_IMAGE_MISSING_OR_INVALID",
		StatusCode: StatusBadRequest,
		Message:    MsgTopupImageRequiredOrMissing,
		Type:       "error",
	}

	ErrorNewsCategoryInvalidID = ResponseCode{
		Code:       "ERROR_NEWS_CATEGORY_INVALID_ID",
		StatusCode: StatusBadRequest,
		Message:    MsgNewsCategoryInvalidID,
		Type:       "error",
	}

	ErrorNewsTagInvalidID = ResponseCode{
		Code:       "ERROR_NEWS_TAG_INVALID_ID",
		StatusCode: StatusBadRequest,
		Message:    MsgNewsTagInvalidID,
		Type:       "error",
	}

	ErrorNewsCategoryWithNameAlreadyExists = ResponseCode{
		Code:       "ERROR_NEWS_CATEGORY_WITH_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgNewsCategoryWithNameAlreadyExists,
		Type:       "error",
	}
	ErrorNewsTagWithNameAlreadyExists = ResponseCode{
		Code:       "ERROR_NEWS_TAG_WITH_NAME_ALREADY_EXISTS",
		StatusCode: StatusBadRequest,
		Message:    MsgNewsTagWithNameAlreadyExists,
		Type:       "error",
	}

	SuccessNewsCategoryCreated = ResponseCode{
		Code:       "SUCCESS_NEWS_CATEGORY_CREATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgNewsCategoryCreatedSuccess,
		Type:       "success",
	}
	SuccessNewsCategoryFetched = ResponseCode{
		Code:       "SUCCESS_NEWS_CATEGORY_FETCHE",
		StatusCode: StatusOK,
		Message:    MsgNewsCategoryFetchedSuccess,
		Type:       "success",
	}
	SuccessNewsCategoryUpdated = ResponseCode{
		Code:       "SUCCESS_NEWS_CATEGORY_UPDATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgNewsCategoryUpdatedSuccess,
		Type:       "success",
	}
	SuccessNewsCategoryDeleted = ResponseCode{
		Code:       "SUCCESS_NEWS_CATEGORY_DELETE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgNewsCategoryDeletedSuccess,
		Type:       "success",
	}

	SuccessNewsTagCreated = ResponseCode{
		Code:       "SUCCESS_NEWS_TAG_CREATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgNewsTagCreatedSuccess,
		Type:       "success",
	}
	SuccessNewsTagFetched = ResponseCode{
		Code:       "SUCCESS_NEWS_TAG_FETCHE",
		StatusCode: StatusOK,
		Message:    MsgNewsTagFetchedSuccess,
		Type:       "success",
	}
	SuccessNewsTagUpdated = ResponseCode{
		Code:       "SUCCESS_NEWS_TAG_UPDATE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgNewsTagUpdatedSuccess,
		Type:       "success",
	}
	SuccessNewsTagDeleted = ResponseCode{
		Code:       "SUCCESS_NEWS_TAG_DELETE_REQUEST_SENT",
		StatusCode: StatusOK,
		Message:    MsgNewsTagDeletedSuccess,
		Type:       "success",
	}
	// Sitota Related Responses
	SuccessAllSitotasRetrieved = ResponseCode{
		Code:       "SUCCESS_ALL_SITOTAS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgAllSitotasRetrievedSuccess,
		Type:       "success",
	}
	SuccessSitotaRetrieved = ResponseCode{
		Code:       "SUCCESS_SITOTA_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgSitotaRetrievedSuccess,
		Type:       "success",
	}

	ErrorSitotaRequired = ResponseCode{
		Code:       "ERROR_SITOTA_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgSitotaRequired,
		Type:       "error",
	}
	ErrorBpsActionRoleNotFound = ResponseCode{
		Code:       "ERROR_BPS_ACTION_ROLE_NOT_FOUND",
		StatusCode: StatusBadRequest,
		Message:    MsgBpsActionRoleNotFound,
		Type:       "error",
	}

	ErrorRoleNotFound = ResponseCode{
		Code:       "ERROR_ROLE_NOT_FOUND",
		StatusCode: StatusBadRequest,
		Message:    MsgRoleNotFound,
		Type:       "error",
	}

	// Transaction Service Related Responses
	ErrorTransactionIDRequired = ResponseCode{
		Code:       "ERROR_TRANSACTION_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgTransactionIDRequired,
		Type:       "error",
	}
	SuccessTransactionRetrieved = ResponseCode{
		Code:       "SUCCESS_TRANSACTION_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgTransactionRetrievedSuccess,
		Type:       "success",
	}

	ErrorEventMerchantInvalidMerchantID = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_INVALID_MERCHANT_ID",
		StatusCode: StatusBadRequest,
		Message:    MsgEventMerchantInvalidID,
		Type:       "error",
	}
	ErrorEventMerchantInvalidMerchantType = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_INVALID_MERCHANT_TYPE",
		StatusCode: StatusBadRequest,
		Message:    MsgEventMerchantInvalidType,
		Type:       "error",
	}
	ErrorEventMerchantInvalidSettlementMethod = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_INVALID_SETTLEMENT_METHOD",
		StatusCode: StatusBadRequest,
		Message:    MsgEventMerchantInvalidMethod,
		Type:       "error",
	}
	ErrorEventMerchantInvalidMerchantName = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_INVALID_MERCHANT_NAME",
		StatusCode: StatusBadRequest,
		Message:    MsgEventMerchantInvalidName,
		Type:       "error",
	}
	ErrorEventMerchantInvalidBankAccountNumber = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_INVALID_BANK_ACCOUNT_NUMBER",
		StatusCode: StatusBadRequest,
		Message:    MsgEventMerchantInvalidAccountNumber,
		Type:       "error",
	}
	ErrorEventMerchantInvalidEmail = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_INVALID_EMAIL",
		StatusCode: StatusBadRequest,
		Message:    MsgEventMerchantInvalidEmail,
		Type:       "error",
	}
	ErrorEventMerchantInvalidPhoneNumber = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_INVALID_PHONE_NUMBER",
		StatusCode: StatusBadRequest,
		Message:    MsgEventMerchantInvalidPhoneNumber,
		Type:       "error",
	}

	ErrorEventMerchantNotFound = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgEventMerchantNotFound,
		Type:       "error",
	}
	ErrorEventMerchantDisableFailed = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_DISABLE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgEventMerchantDisableFailed,
		Type:       "error",
	}
	ErrorEventMerchantEnableFailed = ResponseCode{
		Code:       "ERROR_EVENT_MERCHANT_ENABLE_FAILED",
		StatusCode: StatusInternalServerError,
		Message:    MsgEventMerchantEnableFailed,
		Type:       "error",
	}

	ErrorCustomerAccountNumberMustContainOnlyNumbers = ResponseCode{
		Code:       "ERROR_CUSTOMER_ACCOUNT_NUMBER_MUST_CONTAIN_ONLY_NUMBERS",
		StatusCode: StatusBadRequest,
		Message:    MsgCustomerAccountNumberMustContainOnlyNumbers,
		Type:       "error",
	}
	ErrorCustomerCIFMustContainOnlyNumbers = ResponseCode{
		Code:       "ERROR_CUSTOMER_CIF_MUST_CONTAIN_ONLY_NUMBERS",
		StatusCode: StatusBadRequest,
		Message:    MsgCustomerCIFMustContainOnlyNumbers,
		Type:       "error",
	}

	UnableToCreateAccessListSegmentation = ResponseCode{
		Code:       "ERROR_UNABLE_TO_CREATE_ACCESS_LIST_SEGMENTATION",
		StatusCode: StatusInternalServerError,
		Message:    "Unable to create access list segmentation",
		Type:       "error",
	}
	AccessListSegmentationCreatedSuccessfully = ResponseCode{
		Code:       "SUCCESS_ACCESS_LIST_SEGMENTATION_CREATED",
		StatusCode: StatusOK,
		Message:    "Access list segmentation created successfully",
		Type:       "success",
	}

	// Access List Segmentation Success Codes
	SuccessAccessListSegmentationCreated = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_CREATED",
		StatusCode: StatusOK,
		Message:    MsgAccessListSegmentationCreatedSuccessfully,
		Type:       "success",
	}
	SuccessAccessListSegmentationUpdated = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_UPDATED",
		StatusCode: StatusOK,
		Message:    MsgAccessListSegmentationUpdatedSuccessfully,
		Type:       "success",
	}
	SuccessAccessListSegmentationEnabled = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_ENABLED",
		StatusCode: StatusOK,
		Message:    MsgAccessListSegmentationEnabledSuccessfully,
		Type:       "success",
	}
	SuccessAccessListSegmentationDisabled = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_DISABLED",
		StatusCode: StatusOK,
		Message:    MsgAccessListSegmentationDisabledSuccessfully,
		Type:       "success",
	}
	SuccessAccessListSegmentationsRetrieved = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATIONS_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgAccessListSegmentationsRetrievedSuccessfully,
		Type:       "success",
	}
	SuccessAccessListSegmentationRetrieved = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_RETRIEVED",
		StatusCode: StatusOK,
		Message:    MsgAccessListSegmentationRetrievedSuccessfully,
		Type:       "success",
	}

	// Access List Segmentation Error Codes
	ErrorAccessListSegmentationNotFound = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_NOT_FOUND",
		StatusCode: StatusBadRequest,
		Message:    MsgAccessListSegmentationNotFound,
		Type:       "error",
	}
	ErrorAccessListSegmentationAlreadyEnabled = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_ALREADY_ENABLED",
		StatusCode: StatusConflict,
		Message:    MsgAccessListSegmentationAlreadyEnabled,
		Type:       "error",
	}
	ErrorAccessListSegmentationAlreadyDisabled = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_ALREADY_DISABLED",
		StatusCode: StatusConflict,
		Message:    MsgAccessListSegmentationAlreadyDisabled,
		Type:       "error",
	}
	ErrorAccessListSegmentationInvalidID = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_INVALID_ID",
		StatusCode: StatusBadRequest,
		Message:    MsgAccessListSegmentationInvalidID,
		Type:       "error",
	}
	ErrorAccessListSegmentationNameAlreadyExists = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_NAME_ALREADY_EXISTS",
		StatusCode: StatusConflict,
		Message:    MsgAccessListSegmentationNameAlreadyExists,
		Type:       "error",
	}
	ErrorCustomerSegmentationCodeNotFound = ResponseCode{
		Code:       "CUSTOMER_SEGMENTATION_CODE_NOT_FOUND",
		StatusCode: StatusNotFound,
		Message:    MsgCustomerSegmentationCodeNotFound,
		Type:       "error",
	}
	ErrorAccessListSegmentationIDSRequired = ResponseCode{
		Code:       "ACCESS_LIST_SEGMENTATION_IDS_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgAccessListSegmentationIDsRequired,
		Type:       "error",
	}
	ErrorServiceIdRequired = ResponseCode{
		Code:       "ERROR_SERVICE_ID_REQUIRED",
		StatusCode: StatusBadRequest,
		Message:    MsgServiceIdRequired,
		Type:       "error",
	}
)
