package localization

// Success Messages
const (
	// User related success messages
	MsgUserCreatedSuccessfully   = "User created successfully"
	MsgAvatarCreatedSuccessfully = "Avatar created successfully"
	MsgUserUpdatedSuccessfully   = "User updated successfully"
	MsgUserDeletedSuccessfully   = "User deleted successfully"
	MsgUserRetrievedSuccessfully = "User retrieved successfully"
	MsgUserLoginSuccessfully     = "User logged in successfully"
	MsgUserLogoutSuccessfully    = "User logged out successfully"
	MsgUserProfileUpdated        = "User profile updated successfully"
	MsgUserPasswordChanged       = "Password changed successfully"
	MsgUserPINChanged            = "PIN changed successfully"
	MsgUserDeviceLinked          = "Device linked successfully"
	MsgUserDeviceUnlinked        = "Device unlinked successfully"
	MSGAccountNumberRequired     = "Account Number is required"
	MsgFeedbackIDRequired        = "Feedback ID is required"
	MsgInvalidIDFormat           = "Invalid ID format"
	MSGIncompleteUserInfo        = "Incomplte user info"
	MSGUserCodeIsRequired        = "User code is required"
	MsgActionAlreadyExists       = "Action already requested wait for checker approval"
	// OTP related success messages
	MsgOTPSentSuccessfully     = "OTP sent successfully"
	MsgOTPVerifiedSuccessfully = "OTP verified successfully"
	MsgOTPResentSuccessfully   = "OTP resent successfully"

	// Session related success messages
	MsgSessionCreatedSuccessfully    = "Session created successfully"
	MsgSessionRefreshedSuccessfully  = "Session refreshed successfully"
	MsgSessionTerminatedSuccessfully = "Session terminated successfully"

	// File related success messages
	MsgFileUploadedSuccessfully  = "File uploaded successfully"
	MsgFileDeletedSuccessfully   = "File deleted successfully"
	MsgFileRetrievedSuccessfully = "File retrieved successfully"

	// General success messages
	MsgOperationCompleted        = "Operation completed successfully"
	MsgDataRetrievedSuccessfully = "Data retrieved successfully"
	MsgDataSavedSuccessfully     = "Data saved successfully"
	MsgDataUpdatedSuccessfully   = "Data updated successfully"
	MsgDataDeletedSuccessfully   = "Data deleted successfully"
	MsgValidationPassed          = "Validation passed successfully"
	MsgHealthCheckPassed         = "Health check passed"
	MsgHealthCheckFailed         = "Health check failed"
	MsgBadRequest                = "Bad request"

	// Bank related success messages
	MsgBankCreatedRequestSent     = "Bank created request sent successfully"
	MsgBankUpdatedRequestSent     = "Bank updated request sent successfully"
	MsgBankDeletedRequestSent     = "Bank deleted request sent successfully"
	MsgBanksRetrievedRequestSent  = "Banks retrieved request sent successfully"
	MsgBankRetrievedSuccessfully  = "Bank retrieved successfully"
	MsgBankRejectedRequestSent    = "Bank rejected request sent successfully"
	MsgBankDisableRequestSent     = "Bank disable request sent successfully"
	MsgBankEnableRequestSent      = "Bank enable request sent successfully"
	MsgBankLogoUpdatedRequestSent = "Bank logo updated request sent successfully"

	// Donation related success messages
	MsgDonationCategoryCreateRequestSent = "Donation category create request sent successfully"
	MsgDonationCategoriesFetched         = "Donation categories fetched successfully"
	MsgDonationCategoryFetched           = "Donation category fetched successfully"
	MsgDonationCategoryUpdated           = "Donation category updated successfully"
	MsgDonationCompanyCreateRequestSent  = "Donation company create request sent successfully"
	MsgDonationCompaniesFetched          = "Donation companies fetched successfully"
	MsgDonationCompanyFetched            = "Donation company fetched successfully"
	MsgDonationCompanyUpdatedRequestSent = "Donation company updated request sent successfully"
	MsgDonationCreateRequestSent         = "Donation create request sent successfully"
	MsgDonationsFetched                  = "Donations fetched successfully"
	MsgDonationFetched                   = "Donation fetched successfully"
	MsgDonationUpdateRequestSent         = "Donation update request sent successfully"
	MsgDonationImageUpdateRequestSent    = "Donation image update request sent successfully"
	MsgDonationImageDeleteRequestSent    = "Donation image delete request sent successfully"
	MsgDonationImageAddRequestSent       = "Donation image add request sent successfully"

	// CPS Action related success messages
	MsgCPSActionsRetrievedSuccessfully = "CPS Actions retrieved successfully"
	MsgCPSActionAuthorizedSuccessfully = "CPS Action authorized successfully"
	MsgCPSActionRejectedSuccessfully   = "CPS Action rejected successfully"
	MsgCPSActionsRetrieved             = "CPS actions retrieved successfully"

	// Wallet related success messages
	MsgWalletCreationRequestSent    = "Wallet creation request sent successfully"
	MsgWalletUpdateRequestSent      = "Wallet update request sent successfully"
	MsgWalletDeletedSuccessfully    = "Wallet deleted successfully"
	MsgWalletsRetrievedSuccessfully = "Wallets retrieved successfully"
	MsgWalletRetrievedSuccessfully  = "Wallet retrieved successfully"

	// Event related success messages
	MsgEventCreationRequestSubmitted = "Event creation request submitted successfully"
	MsgEventUpdateRequestSubmitted   = "Event update request submitted successfully"
	MsgEventDeleteRequestSubmitted   = "Event delete request submitted successfully"
	MsgEventEnableRequestSubmitted   = "Event enable request submitted successfully"
	MsgEventDisableRequestSubmitted  = "Event disable request submitted successfully"
	MsgEventSuccessfullyRetrieved    = "Event successfully retrieved"
	MsgEventsSuccessfullyRetrieved   = "Events successfully retrieved"

	// Product Code related success messages
	MsgProductCodeUpdateRequestSubmitted = "Product code update request submitted successfully"
	MsgProductCodeSuccessfullyRetrieved  = "Product code successfully retrieved"
	MsgProductCodesSuccessfullyRetrieved = "Product codes successfully retrieved"

	// Portal Card related success messages
	MsgPortalCardsFetchedSuccessfully = "Portal cards fetched successfully"

	// Password Rule related success messages
	MsgPasswordRulesSuccessfullyFetched = "Password Rules Successfully Fetched"
	MsgPasswordRuleUpdateActionFetched  = "Fetched password rule update action successfully"
	MsgUpdateActionFetchedSuccessfully  = "Fetched update action successfully"

	// Permission related success messages
	MsgPermissionGroupRequestCreated = "Permission group Request created successfully"
	MsgPermissionGroupsFetched       = "Permission groups fetched successfully"
	MsgPermissionGroupFetched        = "Permission group fetched successfully"
	MsgPermissionGroupRequestUpdated = "Permission group Request updated successfully"
	MsgPermissionCategoriesFetched   = "Permission categories fetched successfully"

	// Mini App Merchant related success messages
	MsgMiniAppMerchantCreateRequestSuccessfully = "Create MiniAppMerchant request successfully created"
	MsgUpdateRequestSuccessfullyCreated         = "Update request successfully created"
	MsgDeleteRequestSuccessfullyCreated         = "Delete request successfully created"
	MsgEnableRequestSuccessfullyCreated         = "Enable request successfully created"
	MsgDisableRequestSuccessfullyCreated        = "Disable request successfully created"

	// Notification related success messages
	MsgNotificationCreationRequestSubmitted = "Notification creation request submitted successfully"
	MsgNotificationUpdateRequestSubmitted   = "Notification update request submitted successfully"
	MsgNotificationDeleteRequestSubmitted   = "Notification delete request submitted successfully"

	// Department related success messages
	MsgDepartmentCreatedSuccessfully = "Department created successfully"

	// Budget Category related success messages
	MsgBudgetCategoryUpdatedSuccessfully = "Budget category updated successfully"

	// BPS User related success messages
	MsgBPSUserApprovedSuccessfully = "BPS user approved successfully"

	// Account Validation related success messages
	MsgValidationRuleApprovedSuccessfully = "Validation rule approved successfully"

	// Feedback related success messages
	MsgFeedbackSavedSuccessfully = "Feedback saved successfully"

	// HQ related success messages
	MsgHQApprovedSuccessfully = "HQ approved successfully"

	// Mini App related success messages
	MsgMiniAppDetailsFetchedSuccessfully = "MiniApp details fetched successfully"
	MsgCredentialInformationGenerated    = "Credential information generated successfully"
	MsgFileUploadedToMinIOSuccessfully   = "File uploaded to MinIO successfully"

	// Key Generation related success messages
	MsgEd25519KeyPairGeneratedSuccessfully = "Ed25519 key pair generated successfully"
	MsgDataSignedSuccessfully              = "Data signed successfully"
	MsgSignatureVerificationSuccessful     = "Signature verification successful"
	MsgFabricIDGeneratedSuccessfully       = "Fabric ID generated successfully"
	MsgNumericCodeGeneratedSuccessfully    = "Numeric code generated successfully"
	MsgAppSecretGeneratedSuccessfully      = "App secret generated successfully"
	MsgAppSecretHashedSuccessfully         = "App secret hashed successfully"
	MsgAppSecretVerificationSuccessful     = "App secret verification successful"
	MsgPrefixedNameGeneratedSuccessfully   = "Prefixed name generated successfully"
	MsgSMSsentSuccessfully                 = "SMS sent successfully"

	// Service related success messages
	MsgServiceRetrievedSuccessfully                 = "Service retrieved successfully"
	MsgServiceTireUpdatedSuccessfully               = "Service tire updated successfully"
	MsgServiceSingleTransferUpdatedSuccessfully     = "Service single transfer updated successfully"
	MsgServiceMinAmountTransferUpdatedSuccessfully  = "Service min amount transfer updated successfully"
	MsgServiceTierDeleteRequestSuccessfully         = "Service tier delete request successfully"
	MsgServiceFeeDetailFetchedSuccessfully          = "Service fee detail fetched successfully"
	MsgCPSActionCreatedForUpdateServiceFee          = "CPS action created for update service fee successfully"
	MsgCPSActionCreatedForUpdateSingleMaxTransfer   = "CPS action created for update single max transfer successfully"
	MsgCPSActionCreatedForUpdateTotalMaxTransferCap = "CPS action created for update total max transfer cap successfully"
	MsgCPSActionCreatedForUpdateMinimumTransferCap  = "CPS action created for update minimum transfer cap successfully"
	MsgCPSActionCreatedForDeleteServiceFeeTire      = "CPS action created for delete service fee tire successfully"
	MsgPreviousServiceDataFetchedSuccessfully       = "Previous service data fetched successfully"
	MsgMaxTotalCapValidatedSuccessfully             = "Max total cap validated successfully"

	// Product Code related success messages
	MsgProductCodeFetchedSuccessfully  = "Product code fetched successfully"
	MsgProductCodesFetchedSuccessfully = "Product codes fetched successfully"
	MsgProductCodeUpdatedSuccessfully  = "Product code updated successfully"

	// Mini App Merchant related success messages
	MsgMiniAppAddedSuccessfully               = "Mini app added successfully"
	MsgMiniAppEnabledStateUpdatedSuccessfully = "Mini app enabled state updated successfully"
	MsgMiniAppSoftDeletedSuccessfully         = "Mini app soft deleted successfully"
	MsgTransactionAbortedSuccessfully         = "Transaction aborted successfully"
	MsgTransactionCommittedSuccessfully       = "Transaction committed successfully"

	// CPS Action related success messages
	MsgCPSActionCreatedSuccessfully       = "CPS action created successfully"
	MsgCPSActionUpdatedSuccessfully       = "CPS action updated successfully"
	MsgCPSActionStatusUpdatedSuccessfully = "CPS action status updated successfully"
	MsgCPSActionFetchedSuccessfully       = "CPS action fetched successfully"

	// Ad related success messages
	MsgAdvertCreatedSuccessfully     = "Advert created successfully"
	MsgAdvertCreatedRequestSent      = "Advert create request sent successfully"
	MsgAdvertUpdatedSuccessfully     = "Advert updated successfully"
	MsgAdvertUpdateRequestSent       = "Advert update request sent successfully"
	MsgAdvertDeletedSuccessfully     = "Advert deleted successfully"
	MsgAdvertDeleteRequestSent       = "Advert delete request sent successfully"
	MsgAdvertEnableDisableSuccessful = "Advert enable/disable successful"
	MsgAdvertFetchedSuccessfully     = "Advert fetched successfully"
	MsgAdvertsFetchedSuccessfully    = "Adverts fetched successfully"
	MsgAdvertEnableRequestSent       = "Advert enable request sent successfully"
	MsgAdvertDisableRequestSent      = "Advert disable request sent successfully"

	// BPS Calls related success messages
	MsgLinkedAccountFetchedSuccessfully = "Linked account fetched successfully"

	// Unlink related success messages
	MsgUserSuccessfullyRetrieved             = "User successfully retrieved"
	MsgUnlinkCifRequestSentSuccessfully      = "Unlink CIF request sent successfully"
	MsgCPSActionCreatedForUnlinkSuccessfully = "CPS action created for unlink successfully"
	MsgDocumentHardDeletedSuccessfully       = "Document hard deleted successfully"

	// Permission related success messages
	MsgPermissionGroupRequestCreatedSuccessfully = "Permission group request created successfully"
	MsgPermissionGroupRequestUpdatedSuccessfully = "Permission group request updated successfully"
	MsgPermissionGroupUpdatedSuccessfully        = "Permission group updated successfully"
	MsgPermissionCategoriesFetchedSuccessfully   = "Permission categories fetched successfully"
	MsgPermissionsFetchedForCategorySuccessfully = "Permissions fetched for category successfully"

	// Notification related success messages
	MsgNotificationEnableRequestSubmitted  = "Notification enable request submitted successfully"
	MsgNotificationDisableRequestSubmitted = "Notification disable request submitted successfully"
	MsgNotificationSuccessfullyRetrieved   = "Notification successfully retrieved"
	MsgNotificationsSuccessfullyRetrieved  = "Notifications successfully retrieved"

	// Mini App Handler related success messages
	MsgUpdateMiniAppRequestCreatedSuccessfully  = "Update Mini App request created successfully"
	MsgMiniAppDeletedRequestCreatedSuccessfully = "Mini App deleted request created successfully"
	MsgMiniAppsSuccessfullyRetrieved            = "Mini Apps successfully retrieved"
	MsgMiniAppFetchedByIDSuccessfully           = "Mini App fetched by ID successfully"
	MsgMiniAppActionCompletedSuccessfully       = "Mini App action completed successfully"

	// Feedback Handler related success messages
	MsgFeedbackCreatedSuccessfully = "Feedback created successfully"
	MsgFeedbackFetchedSuccessfully = "Feedback fetched successfully"

	// HQ related success messages
	MsgHQFetchedSuccessfully  = "HQ fetched successfully"
	MsgHQsFetchedSuccessfully = "HQs fetched successfully"

	// Wallet related success messages
	MsgWalletActionRequestSentSuccessfully = "Wallet action request sent successfully"
)

// Error Messages
const (
	// User related error messages
	MsgUserNotFound             = "User not found"
	MsgUserAlreadyExists        = "User already exists"
	MsgUserAccountBlocked       = "User account is blocked"
	MsgUserUnauthorized         = "User is not authorized"
	MsgUserForbidden            = "User access is forbidden"
	MsgUserInvalidCredentials   = "Invalid credentials"
	MsgUserSessionExpired       = "User session has expired"
	MsgUserTooManyLoginAttempts = "Too many login attempts"
	MsgUserDeviceMismatch       = "Device mismatch detected"
	MsgUserDeviceNotFound       = "User device not found"
	MsgUserDeviceAlreadyExists  = "User device already exists"
	MsgSamePIN                  = "You used the same pin as the old one"

	// OTP related error messages
	MsgOTPNotFound        = "OTP not found"
	MsgOTPExpired         = "OTP has expired"
	MsgOTPInvalid         = "Invalid OTP"
	MsgOTPAlreadyExists   = "OTP already exists, wait until it expires"
	MsgOTPTooManyAttempts = "Too many OTP attempts"
	MsgOTPSendFailed      = "Failed to send OTP"

	// PIN related error messages
	MsgPINInvalid       = "Invalid PIN"
	MsgPINMismatch      = "PIN mismatch"
	MsgPINTooWeak       = "PIN is too weak"
	MsgPINInHistory     = "PIN was recently used"
	MsgPINLengthInvalid = "PIN must be exactly 6 digits"
	MsgPINOnlyDigits    = "PIN must contain only digits"
	MsgPINSequential    = "PIN cannot be sequential"
	MsgPINRedundant     = "PIN cannot be redundant"

	// Session related error messages
	MsgSessionNotFound        = "Session not found"
	MsgSessionExpired         = "Session has expired"
	MsgSessionInvalid         = "Invalid session"
	MsgSessionCreationFailed  = "Failed to create session"
	MsgSessionUpdateFailed    = "Failed to update session"
	MsgSessionRetrievalFailed = "Failed to retrieval session"

	// File related error messages
	MsgFileTooLarge     = "file size is too large"
	MsgFileInvalidType  = "invalid file type"
	MsgFileUploadFailed = "file upload failed"
	MsgFileNotFound     = "file not found"
	MsgFileDeleteFailed = "file deletion failed"

	// Validation error messages
	MsgValidationFailed     = "Validation failed"
	MsgRequiredFieldMissing = "Required field is missing"
	MsgInvalidFormat        = "Invalid format"
	MsgInvalidEmail         = "Invalid email format"
	MsgInvalidPhoneNumber   = "Invalid phone number format"
	MsgInvalidDate          = "Invalid date format"
	MsgInvalidActionData    = "Invalid action"
	MsgUnsupportedAction    = "Unsupported action"
	MsgInvalidID            = "Invalid ID"
	MsgMissingFile          = "Missing file"
	MsgFieldTooLong         = "Field value is too long"
	MsgFieldTooShort        = "Field value is too short"

	// System error messages
	MsgInternalServerError  = "Internal server error occurred"
	MsgServiceUnavailable   = "Service is temporarily unavailable"
	MsgDatabaseError        = "Database operation failed"
	MsgNetworkError         = "Network error occurred"
	MsgTimeoutError         = "Request timeout occurred"
	MsgUnexpectedError      = "An unexpected error occurred"
	MsgConfigurationError   = "Configuration error"
	MsgExternalServiceError = "External service error"

	// Business logic error messages
	MsgInsufficientBalance = "Insufficient balance"
	MsgTransactionFailed   = "Transaction failed"
	MsgLimitExceeded       = "Limit exceeded"
	MsgOperationNotAllowed = "Operation not allowed"
	MsgMaintenanceMode     = "System is in maintenance mode"

	// Bank related error messages
	MsgBankNotFound              = "Bank not found"
	MsgBankAlreadyExists         = "Bank already exists"
	MsgBankAlreadyEnabled        = "Bank already enabled"
	MsgBankAlreadyDisabled       = "Bank already disabled"
	MsgBankNameAlreadyExists     = "Bank name already exists"
	MsgBankBICCodeAlreadyExists  = "Bank BIC code already exists"
	MsgBankFetchFailed           = "Bank fetch failed"
	MsgBankIDRequired            = "Bank ID is required"
	MsgBankNameRequired          = "Bank name is required"
	MsgBankCodeRequired          = "Bank code is required"
	MsgBankBICRequired           = "Bank BIC is required"
	MsgInvalidBankName           = "Invalid bank name"
	MsgBankNameTooShort          = "Bank name too short"
	MsgBankNameTooLong           = "Bank name too long"
	MsgBankNameInvalidCharacters = "Bank name contains invalid characters"

	// Donation related error messages
	MsgDonationCategoryNameAlreadyExists  = "Category name already exists"
	MsgDonationCompanyNameAlreadyExists   = "Company name already exists"
	MsgDonationTitleAlreadyExists         = "Donation title already exists"
	MsgDonationAccountNumberAlreadyExists = "Account number already exists"
	MsgDonationIconRequired               = "Icon is required"
	MsgDonationLogoRequired               = "Logo is required"
	MsgDonationImagesRequired             = "Donation images are required"
	MsgDonationInvalidAmount              = "Invalid donation amount"
	MsgDonationInvalidStartDate           = "Invalid start date format"
	MsgDonationInvalidEndDate             = "Invalid end date format"
	MsgDonationCompanyNotFound            = "Company not found"
	MsgDonationCategoryNotFound           = "Category not found"
	MsgDonationAccountNotFound            = "Account not found"
	MsgDonationUploadFailed               = "Failed to upload donation files"
	MsgDonationMinIOError                 = "MinIO operation failed"
	MsgDonationConfigurationError         = "Configuration error"
	MsgDonationTimeoutError               = "Operation timed out"

	// CPS Action related error messages
	MsgPendingCPSActionExists        = "Pending CPS action already exists"
	MsgCPSActionNotFound             = "CPS action not found"
	MsgBpsUserAlreadyEnabled         = "BPS user already enabled"
	MsgBpsUserAlreadyDisabled        = "BPS user already disabled"
	MsgCPSActionNotPending           = "Action is not in pending status"
	MsgCPSActionAlreadyApproved      = "Action already approved"
	MsgCPSActionAlreadyRejected      = "Action already rejected"
	MsgCPSActionCreationFailed       = "Failed to create CPS action"
	MsgCPSActionUpdateFailed         = "Failed to update CPS action"
	MsgCPSActionApprovalFailed       = "Failed to approve action"
	MsgCPSActionRejectionFailed      = "Failed to reject action"
	MsgCPSActionInvalidStatus        = "Invalid action status"
	MsgCPSActionInvalidRequestAction = "Invalid request action"
	MsgCPSActionUnsupportedAction    = "Unsupported request action"
	MsgCPSActionInvalidFormat        = "Invalid action data format"
	MsgCPSActionMissingData          = "Missing action data"
	MsgCPSActionInvalidType          = "Invalid action type"

	// Wallet related error messages
	MsgWalletNotFound                 = "Wallet not found"
	MsgWalletCreationFailed           = "Failed to create wallet"
	MsgWalletUpdateFailed             = "Failed to update wallet"
	MsgWalletDeletionFailed           = "Failed to delete wallet"
	MsgWalletBalanceInsufficient      = "Insufficient wallet balance"
	MsgWalletTransactionFailed        = "Failed to process wallet transaction"
	MsgWalletNotActive                = "Wallet is not active"
	MsgWalletAlreadyExists            = "Wallet already exists"
	MsgWalletTypeNotSupported         = "Wallet type is not supported"
	MsgWalletLimitExceeded            = "Wallet limit exceeded"
	MsgWalletTransactionNotFound      = "Wallet transaction not found"
	MsgWalletTransactionAlreadyExists = "Wallet transaction already exists"

	// Event related error messages
	MsgEventNameAlreadyExists = "Event name already exists"
	MsgEventNotFound          = "Event not found"
	MsgEventCreationFailed    = "Failed to create event"
	MsgEventUpdateFailed      = "Failed to update event"
	MsgEventDeletionFailed    = "Failed to delete event"
	MsgEventEnableFailed      = "Failed to enable event"
	MsgEventDisableFailed     = "Failed to disable event"
	MsgEventFetchFailed       = "Failed to fetch event"

	// Product Code related error messages
	MsgProductCodeNotFound      = "Product code not found"
	MsgProductCodeUpdateFailed  = "Failed to update product code"
	MsgProductCodeFetchFailed   = "Failed to fetch product code"
	MsgProductCodeInvalidFormat = "Invalid product code format"
	MsgProductCodeRequired      = "Product code is required"

	// Portal Card related error messages
	MsgPortalCardNotFound      = "Portal card not found"
	MsgPortalCardFetchFailed   = "Failed to fetch portal cards"
	MsgPortalCardInvalidFormat = "Invalid portal card format"

	// Password Rule related error messages
	MsgPasswordRuleNotFound      = "Password rule not found"
	MsgPasswordRuleFetchFailed   = "Failed to fetch password rules"
	MsgPasswordRuleUpdateFailed  = "Failed to update password rule"
	MsgPasswordRuleInvalidFormat = "Invalid password rule format"

	// Permission related error messages
	MsgPermissionGroupNotFound       = "Permission group not found"
	MsgPermissionGroupAlreadyExists  = "Permission group already exists"
	MsgPermissionGroupCreationFailed = "Failed to create permission group"
	MsgPermissionGroupUpdateFailed   = "Failed to update permission group"
	MsgPermissionGroupFetchFailed    = "Failed to fetch permission groups"
	MsgPermissionCategoryNotFound    = "Permission category not found"
	MsgPermissionCategoryFetchFailed = "Failed to fetch permission categories"

	// Mini App Merchant related error messages
	MsgMiniAppMerchantNotFound       = "Mini app merchant not found"
	MsgMiniAppMerchantCreationFailed = "Failed to create mini app merchant"
	MsgMiniAppMerchantDeletionFailed = "Failed to delete mini app merchant"
	MsgMiniAppMerchantEnableFailed   = "Failed to enable mini app merchant"
	MsgMiniAppMerchantDisableFailed  = "Failed to disable mini app merchant"
	MsgMiniAppMerchantInvalidFormat  = "Invalid mini app merchant format"
	MsgMiniAppMerchantRequired       = "Mini app merchant is required"

	// Notification related error messages
	MsgNotificationNotFound       = "Notification not found"
	MsgNotificationCreationFailed = "Failed to create notification"
	MsgNotificationUpdateFailed   = "Failed to update notification"
	MsgNotificationDeletionFailed = "Failed to delete notification"
	MsgNotificationInvalidFormat  = "Invalid notification format"
	MsgNotificationRequired       = "Notification is required"

	// Department related error messages
	MsgDepartmentNotFound       = "Department not found"
	MsgDepartmentAlreadyExists  = "Department already exists"
	MsgDepartmentCreationFailed = "Failed to create department"
	MsgDepartmentUpdateFailed   = "Failed to update department"
	MsgDepartmentInvalidFormat  = "Invalid department format"
	MsgDepartmentRequired       = "Department is required"

	// Budget Category related error messages
	MsgBudgetCategoryNotFound       = "Budget category not found"
	MsgBudgetCategoryCreationFailed = "Failed to create budget category"
	MsgBudgetCategoryUpdateFailed   = "Failed to update budget category"
	MsgBudgetCategoryInvalidFormat  = "Invalid budget category format"
	MsgBudgetCategoryRequired       = "Budget category is required"

	// BPS User related error messages
	MsgBPSUserNotFound       = "BPS user not found"
	MsgBPSUserCreationFailed = "Failed to create BPS user"
	MsgBPSUserUpdateFailed   = "Failed to update BPS user"
	MsgBPSUserInvalidFormat  = "Invalid BPS user format"
	MsgBPSUserRequired       = "BPS user is required"

	// Ad related error messages
	MsgAdNotFound       = "Ad not found"
	MsgAdCreationFailed = "Failed to create ad"
	MsgAdUpdateFailed   = "Failed to update ad"
	MsgAdInvalidFormat  = "Invalid ad format"
	MsgAdRequired       = "Ad is required"

	// Account Validation related error messages
	MsgValidationRuleNotFound       = "Validation rule not found"
	MsgValidationRuleCreationFailed = "Failed to create validation rule"
	MsgValidationRuleInvalidFormat  = "Invalid validation rule format"
	MsgValidationRuleRequired       = "Validation rule is required"

	// Feedback related error messages
	MsgFeedbackNotFound       = "Feedback not found"
	MsgFeedbackCreationFailed = "Failed to create feedback"
	MsgFeedbackUpdateFailed   = "Failed to update feedback"
	MsgFeedbackInvalidFormat  = "Invalid feedback format"
	MsgFeedbackRequired       = "Feedback is required"

	// HQ related error messages
	MsgHQNotFound       = "HQ not found"
	MsgHQCreationFailed = "Failed to create HQ"
	MsgHQUpdateFailed   = "Failed to update HQ"
	MsgHQInvalidFormat  = "Invalid HQ format"
	MsgHQRequired       = "HQ is required"

	// Mini App related error messages
	MsgMiniAppNotFound       = "Mini app not found"
	MsgMiniAppCreationFailed = "Failed to create mini app"
	MsgMiniAppUpdateFailed   = "Failed to update mini app"
	MsgMiniAppInvalidFormat  = "Invalid mini app format"
	MsgMiniAppRequired       = "Mini app is required"

	// Key Generation related error messages
	MsgKeyGenerationFailed   = "Failed to generate key"
	MsgKeySigningFailed      = "Failed to sign data"
	MsgKeyVerificationFailed = "Failed to verify signature"
	MsgKeyInvalidFormat      = "Invalid key format"
	MsgKeyRequired           = "Key is required"

	// SMS related error messages
	MsgSMSSendingFailed = "Failed to send SMS"
	MsgSMSInvalidFormat = "Invalid SMS format"
	MsgSMSRequired      = "SMS is required"

	// General business error messages
	MsgDuplicateKey               = "Duplicate key error"
	MsgInvalidInput               = "Invalid input provided"
	MsgMissingRequiredFields      = "Missing required fields"
	MsgResourceNotFound           = "Resource not found"
	MsgResourceAlreadyExists      = "Resource already exists"
	MsgResourceUpdateFailed       = "Failed to update resource"
	MsgResourceCreationFailed     = "Failed to create resource"
	MsgResourceDeletionFailed     = "Failed to delete resource"
	MsgResourceFetchFailed        = "Failed to fetch resource"
	MsgResourceInvalidFormat      = "Invalid resource format"
	MsgResourceRequired           = "Resource is required"
	MsgResourceBusy               = "Resource is busy"
	MsgResourceLocked             = "Resource is locked"
	MsgResourceExpired            = "Resource has expired"
	MsgResourceInvalid            = "Resource is invalid"
	MsgResourceUnauthorized       = "Resource access unauthorized"
	MsgResourceForbidden          = "Resource access forbidden"
	MsgResourceConflict           = "Resource conflict"
	MsgResourceTooManyRequests    = "Too many requests for resource"
	MsgResourceUnprocessable      = "Resource cannot be processed"
	MsgResourceNotImplemented     = "Resource operation not implemented"
	MsgResourceBadGateway         = "Resource gateway error"
	MsgResourceServiceUnavailable = "Resource service unavailable"
	MsgResourceGatewayTimeout     = "Resource gateway timeout"
	MsgResourceNetworkError       = "Resource network error"
	MsgResourceConfigurationError = "Resource configuration error"
	MsgResourceMaintenanceMode    = "Resource in maintenance mode"
	MsgInvalidRequestOnParam      = "Invalid request on page and per page"
	MsgInvalidRequest             = "Invalid request"
	MsgInvalidToken               = "Invalid token"

	// Service related error messages
	MsgServiceFetchFailed                 = "Failed to fetch service"
	MsgServiceCountFailed                 = "Failed to count total services"
	MsgServicePendingActionCheckFailed    = "Failed to check pending action"
	MsgServicePreviousDataFetchFailed     = "Failed to fetch previous service data"
	MsgServiceCPSActionCreationFailed     = "Failed to create CPS action for service update"
	MsgServiceValidationFailed            = "Service validation failed"
	MsgServiceHQDataFetchFailed           = "Failed to fetch HQ data for service update"
	MsgServiceUnmarshalFailed             = "Failed to unmarshal service request"
	MsgServiceUpdateFailed                = "Failed to update service"
	MsgServiceTierNotFound                = "Service tier not found"
	MsgServiceApproveFailed               = "Failed to approve service update"
	MsgServiceAuthorizeFailed             = "Failed to authorize service action"
	MsgServiceCPSActionStatusUpdateFailed = "Failed to update CPS action status"
	MsgServiceUnhandledServerError        = "Unhandled server error occurred"
	MsgServiceUpdateFailedError           = "Service update failed"
	MsgNoDataProvidedForUpdate            = "No data provided for service update"
	MsgServiceAuthorizeDeleteFailed       = "Failed to authorize service delete"
	MsgServiceUnknownRequestAction        = "Unknown request action for service"

	// Product Code related error messages
	MsgProductCodeParseFailed          = "Failed to parse product code ID"
	MsgProductCodeDatabaseQueryFailed  = "Database query failed for product code"
	MsgProductCodeCountFailed          = "Failed to count product codes"
	MsgProductCodeDatabaseUpdateFailed = "Database update failed for product code"

	// Mini App Merchant related error messages
	MsgMiniAppMerchantCheckPendingFailed     = "Failed to check pending request"
	MsgMiniAppMerchantExistsCheckFailed      = "Failed to check permission group existence"
	MsgMiniAppMerchantFetchCategoriesFailed  = "Failed to fetch permission categories"
	MsgMiniAppMerchantUnexpectedDBError      = "Unexpected database error occurred"
	MsgMiniAppMerchantMarshalFailed          = "Failed to marshal current action"
	MsgMiniAppMerchantUnmarshalFailed        = "Failed to unmarshal current action"
	MsgMiniAppMerchantCreateFromActionFailed = "Failed to create permission group from action"
	MsgMiniAppMerchantUnmarshalActionFailed  = "Failed to unmarshal action data"
	MsgMiniAppMerchantUpdateFailed           = "Failed to update permission group"
	MsgMiniAppMerchantFetchPermissionsFailed = "Failed to fetch permissions for category"

	// Notification related error messages
	MsgNotificationMapFailed       = "Failed to map notification to document"
	MsgNotificationInsertFailed    = "Failed to insert notification"
	MsgNotificationInvalidIDFormat = "Invalid notification ID format"
	MsgNotificationFetchFailed     = "Failed to fetch notification"

	// Customer related error messages
	MsgCustomerCountFailed     = "Failed to count customers"
	MsgCustomerConvertIDFailed = "Failed to convert customer ID"
	MsgCustomerNotFound        = "Customer not found"
	MsgCustomerFetchFailed     = "Failed to get customer"

	// Account Validation related error messages
	MsgValidationRuleGetFailed                = "Failed to get validation rule"
	MsgValidationRuleMinMaxLengthMismatch     = "Validation failed: min length cannot exceed max length"
	MsgValidationRuleConvertIDFailed          = "Failed to convert validation rule ID"
	MsgValidationRuleUpdateFailed             = "Failed to update validation rule"
	MsgValidationRulePendingActionFetchFailed = "Failed to fetch pending action"
	MsgValidationRuleActionsFetchFailed       = "Failed to fetch actions"

	// Donation related error messages
	MsgDonationCategoryLookupFailed = "Category lookup failed"
	MsgDonationCompanyLookupFailed  = "Company lookup failed"
	MsgDonationLookupFailed         = "Donation lookup failed"
	MsgDonationParseStartDateFailed = "Failed to parse start date"
	MsgDonationParseEndDateFailed   = "Failed to parse end date"
	MsgDonationFetchCompanyFailed   = "Failed to fetch company"
	MsgDonationFetchCategoryFailed  = "Failed to fetch category"
	MsgDonationUpdateFailed         = "Failed to update donation"
	MsgDonationImageUpdateFailed    = "Failed to update donation image"
	MsgDonationImageDeleteFailed    = "Failed to delete donation image"
	MsgDonationImageAddFailed       = "Failed to add donation image"

	// Bank related error messages
	MsgBankFileParseFailed              = "Failed to parse bank file"
	MsgBankRejectionPayloadDecodeFailed = "Failed to decode rejection payload"
	MsgBankLogoUpdateFailed             = "Failed to update bank logo"

	// CPS Action related error messages
	MsgCPSActionRejectionPayloadDecodeFailed = "Failed to decode rejection payload"

	// Permission related error messages
	MsgPermissionGroupRequestCreationFailed = "Failed to create permission group request"
	MsgPermissionGroupRequestUpdateFailed   = "Failed to update permission group request"

	// Department related error messages
	MsgDepartmentCreationSuccess        = "Department created successfully"
	MsgDepartmentUpdateCPSActionCreated = "Department update CPS action created successfully"

	// Budget Category related error messages
	MsgBudgetCategoryUpdateSuccess = "Budget category updated successfully"

	// Donation Service related error messages
	MsgDonationCategoryIconUploadSuccess = "Successfully uploaded donation category icon"
	MsgDonationCategoryIconUpdateSuccess = "Successfully uploaded updated donation category icon"
	MsgDonationCPSActionApprovedSuccess  = "Successfully approved and authorized CPS action"
	MsgDonationCompanyLogoUploadSuccess  = "Successfully uploaded donation company logo"
	MsgDonationCompanyLogoUpdateSuccess  = "Successfully uploaded updated donation company logo"
	MsgDonationImageUploadSuccess        = "Successfully uploaded donation image"
	MsgDonationImagesUpdateSuccess       = "Successfully uploaded updated donation images"

	// Notification Service related error messages
	MsgNotificationConstructedSuccess = "Notification constructed successfully"
	MsgNotificationUpdateSuccess      = "Notification updated successfully"

	// Permission Service related error messages
	MsgPermissionGroupRequestCreatedSuccess = "Successfully created permission group request"

	// Mini App Service related error messages
	MsgMiniAppDetailsFetchedSuccess      = "Successfully fetched MiniApp details"
	MsgMiniAppCredentialGeneratedSuccess = "Successfully generated credential information"
	MsgMiniAppFileUploadedSuccess        = "Successfully uploaded file to MinIO"

	// Feedback Kafka related error messages
	MsgFeedbackSavedToDatabaseSuccess = "Successfully saved feedback to database"

	// HQ Service related error messages
	MsgHQApprovedSuccess = "HQ approved successfully"

	// CPS Action related error messages
	MsgCPSActionStatusInvalid = "Invalid CPS action status"

	// BPS User Service related error messages
	MsgBPSUserApprovedSuccess = "BPS user approved successfully"

	MsgUserUnlinkFailed = "User Unlink Failed"
	// Ad Service related success messages
	MsgAdvertCreatedSuccess = "Advert created successfully"
	MsgAdvertUpdateSuccess  = "Advert updated successfully"

	// Ad Service related error messages
	MsgAdvertCreateError     = "Advert creation failed"
	MsgAdvertUpdateError     = "Advert update failed"
	MsgAdvertAlreadyEnabled  = "Advert already enabled"
	MsgAdvertAlreadyDisabled = "Advert already disabled"
	MsgAdvertNotFound        = "Advert not found"

	// Account Validation Service related error messages
	MsgValidationRuleApprovedSuccess = "Validation rule approved successfully"

	MsgInvalidInputParameters = "Invalid input parameters provided"

	MsgMissingOrInvalidImage = "Missing or invalid image"


	MsgBranchSuccessfullyRetrieved = "Branch retrieved successfully"
	MsgBranchesSuccessfullyRetrieved = "Branches retrieved successfully"
)
