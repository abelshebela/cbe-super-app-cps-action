package localization

// Success Messages
const (
	// User related success messages
	MsgUserCreatedSuccessfully     = "User create request sent  successfully"
	MsgAvatarCreatedSuccessfully   = "Avatar create request sent successfully"
	MsgAvatarUpdateSuccessfully    = "Avatar update request sent successfully"
	MsgAvatarEnabledSuccessfully   = "Avatar enable request sent successfully"
	MsgAvatarDisabledSuccessfully  = "Avatar diabled request sent successfully"
	MsgAvatarDeletedSuccessfully   = "Avatar delete request sent successfully"
	MsgAvatarRetrievedSuccessfully = "Avatar retrieved successfully"
	MsgUserUpdatedSuccessfully     = "User update request sent successfully"
	MsgAmountBasedSuccessfullySent = "Amount-based authentication update request sent successfully"
	MsgUserDeletedSuccessfully     = "User delete request sent successfully"
	MsgUserRetrievedSuccessfully   = "User retrieved successfully"
	MsgUserLoginSuccessfully       = "User logged in successfully"
	MsgUserLogoutSuccessfully      = "User logged out successfully"
	MsgUserProfileUpdated          = "User profile updated successfully"
	MsgUserPasswordChanged         = "Password changed successfully"
	MsgUserPINChanged              = "PIN changed successfully"
	MsgUserDeviceLinked            = "Device linked successfully"
	MsgUserDeviceUnlinked          = "Device unlinked successfully"
	MSGAccountNumberRequired       = "Account Number is required"
	MsgFeedbackIDRequired          = "Feedback ID is required"
	MsgInvalidIDFormat             = "Invalid ID format"
	MSGIncompleteUserInfo          = "Incomplte user info"
	MsgInvalidJSONPayload          = "Invalid json payload"
	MSGUserCodeIsRequired          = "User code is required"
	MsgActionAlreadyExists         = "Action already requested wait for checker approval"
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
	MsgMissingContentTypeHeader  = "Missing Content-Type header"

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

	MsgJobRoleCreatedRequestSent = "Job role created request sent successfully"
	MsgJobRoleUpdateRequestSent  = "Job role updated request sent successfully"
	MsgRoleCreatedRequestSent    = "Rule created    request sent successfully"
	MsgRoleUpdatedRequestSent    = "Rule updated request sent successfully"

	// Donation related success messages
	MsgDonationCategoryCreateRequestSent     = "Donation category create request submitted successfully"
	MsgDonationCategoriesFetched             = "Donation categories fetched successfully"
	MsgDonationCategoryFetched               = "Donation category fetched successfully"
	MsgDonationCategoryUpdated               = "Donation category update request submitted successfully"
	MsgDonationCategoryDisableRequestSent    = "Donation category Disable request submitted successfully"
	MsgDonationCompanyEnableRequestSent      = "Donation Company Enable request submitted successfully"
	MsgDonationCompanyDisableRequestSent     = "Donation company Disable request submitted successfully"
	MsgDonationCategoryEnableRequestSent     = "Donation category Enable request submitted successfully"
	MsgDonationCompanyCreateRequestSent      = "Donation company create request submitted successfully"
	MsgDonationCompaniesFetched              = "Donation companies fetched successfully"
	MsgDonationCompanyFetched                = "Donation company fetched successfully"
	MsgAccountInfoFetched                    = "Account info fetched successfully"
	MsgDonationCompanyUpdatedRequestSent     = "Donation company updated request sent successfully"
	MsgSuccessDonationCompanyUpdated         = "Donation company updated request submitted successfully"
	MsgKYCApproved                           = "KYC approve request sent successfully"
	MsgKYCUpdatedRequestSent                 = "KYC updated request sent successfully"
	MsgKYCFetched                            = "KYC fetched successfully"
	MsgDonationCreateRequestSent             = "Donation create request submitted successfully"
	MsgDonationsFetched                      = "Donations fetched successfully"
	MsgDonationFetched                       = "Donation fetched successfully"
	MsgDonationUpdateRequestSent             = "Donation update request submitted successfully"
	MsgDonationImageUpdateRequestSent        = "Donation image update request submitted successfully"
	MsgDonationImageDeleteRequestSent        = "Donation image delete request submitted successfully"
	MsgDonationImageAddRequestSent           = "Donation image add request submitted successfully"
	MsgDonationEnableRequestSent             = "Donation enable request submitted successfully"
	MsgDonationDisableRequestSent            = "Donation disable request submitted successfully"
	MsgSuccessDonationCreateRequestSent      = "Donation create request sent successfully"
	MsgSuccessDonationUpdateRequestSent      = "Donation update request sent successfully"
	MsgSuccessDonationFetchRequestSent       = "Donation fetch request sent successfully"
	MsgSuccessDonationImageUpdateRequestSent = "Donation image update request sent successfully"
	MsgSuccessDonationImageDeleteRequestSent = "Donation image delete request sent successfully"
	MsgSuccessDonationImageAddRequestSent    = "Donation image add request sent successfully"
	MsgSuccessDonationEnableRequestSent      = "Donation enable request sent successfully"
	MsgSuccessDonationDisableRequestSent     = "Donation disable request sent successfully"

	// CPS Action related success messages
	MsgCPSActionsRetrievedSuccessfully = "CPS Actions retrieved successfully"
	MsgCPSActionAuthorizedSuccessfully = "CPS Action authorized successfully"
	MsgCPSActionReversedSuccessfully   = "CPS Action reversed successfully"
	MsgCPSActionRejectedSuccessfully   = "CPS Action rejected successfully"
	MsgCPSActionCanceledSuccessfully   = "CPS Action Canceled Successfully"
	MsgCPSActionsRetrieved             = "CPS actions retrieved successfully"

	// Wallet related success messages
	MsgWalletCreationRequestSent    = "Wallet creation request sent successfully"
	MsgWalletUpdateRequestSent      = "Wallet update request sent successfully"
	MsgWalletDeleteRequestSent      = "Wallet delete request sent successfully"
	MsgWalletDeletedSuccessfully    = "Wallet deleted successfully"
	MsgWalletsRetrievedSuccessfully = "Wallets retrieved successfully"
	MsgWalletRetrievedSuccessfully  = "Wallet retrieved successfully"
	//topup related success messages
	MsgTopupCreationRequestSent    = "Topup creation request sent successfully"
	MsgTopupUpdateRequestSent      = "Topup update request sent successfully"
	MsgTopupDeleteRequestSent      = "Topup delete request sent successfully"
	MsgTopupDeletedSuccessfully    = "Topup deleted successfully"
	MsgTopupsRetrievedSuccessfully = "Topups retrieved successfully"
	MsgTopupRetrievedSuccessfully  = "Topup retrieved successfully"
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
	MsgPermissionGroupRequestDeleted = "Permission group Request deleted successfully"
	MsgPermissionCategoriesFetched   = "Permission categories fetched successfully"

	// Mini App Merchant related success messages
	MsgMiniAppMerchantCreateRequestSuccessfully = "Create MiniAppMerchant request submitted"
	MsgMiniAppMerchantUpdateRequestSuccessfully = "Update request submitted for approval"
	MsgUpdateRequestSuccessfullyCreated         = "Update request successfully created"
	MsgDeleteRequestSuccessfullyCreated         = "Delete request successfully created"
	MsgEnableRequestSuccessfullyCreated         = "enable request successfully created"
	MsgDisableRequestSuccessfullyCreated        = "Disable request successfully created"
	MsgMiniAppMerchantDisable                   = "Create MiniAppMerchant request successfully created"
	MsgMiniAppMerchantEnableSuccessfully        = "Create MiniAppMerchant request successfully created"
	MsgMiniAppMerchantCreatedSuccessfully       = "Mini app merchant created successfully"

	// mini app
	MsgMiniAppCreateRequestSuccessfully = "mini App request successfully created"

	// Notification related success messages
	MsgNotificationCreationRequestSubmitted = "Notification creation request submitted successfully"
	MsgNotificationUpdateRequestSubmitted   = "Notification update request submitted successfully"
	MsgNotificationDeleteRequestSubmitted   = "Notification delete request submitted successfully"

	// Department related success messages
	MsgDepartmentCreateRequestedSuccessfully  = "Department create requested successfully"
	MsgDepartmentUpdateRequestedSuccessfully  = "Department update requested successfully"
	MsgDepartmentEnableRequestedSuccessfully  = "Department enable requested successfully"
	MsgDepartmentDisableRequestedSuccessfully = "Department disable requested successfully"

	// Bank Vault related success messages
	MsgBankVaultCreationRequestSubmitted     = "Bank vault creation request submitted successfully"
	MsgBankVaultUpdateRequestSubmitted       = "Bank vault update request submitted successfully"
	MsgBankVaultDeleteRequestSubmitted       = "Bank vault deleted  request submitted successfully"
	MsgBankVaultDisableRequestSubmitted      = "Bank vault disabled request submitted successfully"
	MsgBankVaultEnableRequestSubmitted       = "Bank vault enabled request submitted successfully"
	MsgBankVaultsRetrievedSuccessfully       = "Bank vaults retrieved successfully"
	MsgBankVaultRetrievedSuccessfully        = "Bank vault retrieved successfully"
	MsgBankLockedVaultsRetrievedSuccessfully = "Bank locked vaults retrieved successfully"
	MsgBankLockedVaultRetrievedSuccessfully  = "Bank locked vault retrieved successfully"
	MsgGroupVaultsRetrievedSuccessfully      = "Group vaults retrieved successfully"
	MsgGroupVaultRetrievedSuccessfully       = "Group vault retrieved successfully"

	// vault group category sucess messages
	MsgVaultGroupCategoryCreationRequestSubmitted = "Vault group category creation request submitted successfully"
	MsgVaultGroupCategoriesRetrieved              = "Vault group categories retrieved successfully"
	MsgVaultGroupCategoryRetrieved                = "Vault group category retrieved successfully"
	MsgVaultGroupCategoryUpdateRequestSubmitted   = "Vault group category update request submitted successfully"
	MsgVaultGroupCategoryDeleteRequestSubmitted   = "Vault group category deleted  request submitted successfully"
	MsgVaultGroupCategoryEnableRequestSubmitted   = "Vault group category enabled request submitted successfully"
	MsgVaultGroupCategoryDisableRequestSubmitted  = "Vault group category disabled request submitted successfully"
	MsgBudgetCategoryUpdatedSuccessfully          = "Budget category updated successfully"

	MsgVaultAmountTierCreated                 = "Vault amount tier creation request submitted successfully"
	MsgVaultAmountTierFetchedSuccessfully     = "Vault amount tier retrieved successfully"
	MsgVaultAmountTierUpdateRequestSubmitted  = "Vault amount tier update request submitted successfully"
	MsgVaultAmountTierDeleteRequestSubmitted  = "Vault amount tier delete request submitted successfully"
	MsgVaultAmountTierDisableRequestSubmitted = "Vault amount tier disable request submitted successfully"
	MsgVaultAmountTierEnableRequestSubmitted  = "Vault amount tier enable request submitted successfully"
	// BPS User related success messages
	MsgBPSUserApprovedSuccessfully = "BPS user approved successfully"

	// Account Validation related success messages
	MsgValidationRuleApprovedSuccessfully = "Update request submitted for approval"

	// Feedback related success messages
	MsgFeedbackSavedSuccessfully = "Feedback saved successfully"

	// HQ related success messages
	MsgHQApprovedSuccessfully                 = "HQ approved successfully"
	MsgHQArchiveTimeUpdateRequestSubmitted    = "Update archive time request submitted for approval"
	MsgHQPasswordExpiryUpdateRequestSubmitted = "Update password expiry request submitted for approval"
	MsgHQBlockTimeUpdateRequestSubmitted      = "Update block time request submitted for approval"

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

	MsgUpdateMiniAppRequestCreatedSuccessfully    = "Update Mini App request created successfully"
	MsgMiniAppDeletedRequestCreatedSuccessfully   = "Mini App deleted request created successfully"
	MsgMiniAppDisableRequestSubmittedSuccessfully = "Mini App disable request submitted successfully"
	MsgMiniAppEnableRequestSubmittedSuccessfully  = "Mini App enable request submitted successfully"
	MsgCreateMiniAppRequestCreatedSuccessfully    = "Mini App create request sent successfully"
	MsgMiniAppsSuccessfullyRetrieved              = "Mini Apps successfully retrieved"
	MsgMiniAppFetchedByIDSuccessfully             = "Mini App fetched by ID successfully"
	MsgMiniAppActionCompletedSuccessfully         = "Mini App action completed successfully"

	// cps_user related success messages
	MsgCpsUserCreationRequestSubmitted = "CPS user creation request submitted successfully"
	MsgCpsUserUpdateRequestSubmitted   = "CPS user update request submitted successfully"
	MsgCpsUserDeletedSuccessfully      = "CPS user deleted request sent successfully"
	MsgCpsUsersRetrievedSuccessfully   = "CPS users fetched successfully"
	MsgCpsUserRetrievedSuccessfully    = "CPS user retrieved successfully"
	MsgCpsUserEnabledSuccessfully      = "CPS user enabled request sent successfully"
	MsgCpsUserDisabledSuccessfully     = "CPS user disable request sent successfully"

	MsgBpsUserEnabledRequestedSuccessfully  = "BPS user enabled request sent successfully"
	MsgBpsUserDisabledRequestedSuccessfully = "BPS user disabled request sent successfully"

	// CPS Roles
	MsgCpsRoleCreated         = "CPS role create request submitted successfully"
	MsgCpsRoleUpdated         = "CPS role update request submitted successfully"
	MsgCpsRolesFetched        = "CPS roles fetched successfully"
	MsgCpsRoleFetched         = "CPS role fetched successfully"
	MsgCpsRoleEnabled         = "CPS role enable request submitted successfully"
	MsgCpsRoleDisabled        = "CPS role disable request submitted successfully"
	MsgCpsRoleAlreadyEnabled  = "CPS role already enabled"
	MsgCpsRoleAlreadyDisabled = "CPS role already disabled"

	// Feedback Handler related success messages
	MsgFeedbackCreatedSuccessfully = "Feedback created successfully"
	MsgFeedbackFetchedSuccessfully = "Feedback fetched successfully"

	// HQ related success messages
	MsgHQFetchedSuccessfully               = "HQ fetched successfully"
	MsgHQsFetchedSuccessfully              = "HQs fetched successfully"
	MsgHQBlockTImeFetchedSuccessfully      = "Block time fetched successfully"
	MsgHQArchiveTimeFetchedSuccessfully    = "Archive time fetched successfully"
	MsgHQPasswordExpiryFetchedSuccessfully = "Password expiry fetched successfully"

	// Wallet related success messages
	MsgWalletActionRequestSentSuccessfully = "Wallet action request sent successfully"

	// permission related error messages
	MsgGroupNameRequired = "Group name is required"
)

// Error Messages
const (
	MsgInvalidKey     = "Invalid key"
	MsgInvalidEncData = "Invalid encryption data"
	MsgInvalidPadding = "Invalid padding"

	MsgDepartmentIDRequired = "Department ID is required"

	// User related error messages
	MsgUserNotFound                 = "User not found"
	MsgUserAlreadyExists            = "User already exists"
	MsgUserAccountBlocked           = "User account is blocked"
	MsgUserUnauthorized             = "User is not authorized"
	MsgThirdApiRequestNotAuthorized = "Third API request is not authorized"
	MsgUserForbidden                = "User access is forbidden"
	MsgUserInvalidCredentials       = "Invalid credentials"
	MsgUserSessionExpired           = "User session has expired"
	MsgUserTooManyLoginAttempts     = "Too many login attempts"
	MsgUserDeviceMismatch           = "Device mismatch detected"
	MsgUserDeviceNotFound           = "User device not found"
	MsgUserDeviceAlreadyExists      = "User device already exists"
	MsgSamePIN                      = "You used the same pin as the old one"

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
	MsgMiniAppNameAlreadyExists    = "Mini App name already exists"
	MsgBankDeleteRequestFailed     = "Bank delete request failed"
	MsgBankImageRequiredOrMissing  = "Image invalid or missing"
	MsgTopupImageRequiredOrMissing = "Required topup image invalid or missing"

	MsgWalletImageRequiredOrMissing          = "Required wallet image invalid or missing"
	MsgValidationFailed                      = "Validation failed"
	MsgRequiredFieldMissing                  = "Required field is missing"
	MsgServiceExists                         = "Service already exists"
	MsgInvalidFormat                         = "Invalid format"
	MsgInvalidInputParameter                 = "Invalid input parameter"
	MsgActionNameIsRequired                  = "Action name is required"
	MsgInvalidEmail                          = "Invalid email format"
	MsgInvalidPhoneNumber                    = "Invalid phone number format"
	MsgExistEmail                            = "cps user with this email already exist"
	MsgExistPhoneNumber                      = "cps user with this  phone Number already exist"
	MsgInvalidDate                           = "Invalid date format"
	MsgInvalidAction                         = "Invalid Action"
	MsgFieldTooLong                          = "Field value is too long"
	MsgFieldTooShort                         = "Field value is too short"
	MsgBankDisableRequestSuccessfullyCreated = "Bank Disable request successfully created"
	MsgBankDisableRequestFailed              = "Bank Disable request  failed"

	MsgBankEnableRequestSuccessfullyCreated = "Bank Enable request successfully created"
	MsgBankEnableRequestFailed              = "Bank Enable request  failed"
	MsgInvalidActionData                    = "Invalid action"
	MsgUnsupportedAction                    = "Unsupported action"
	MsgMissingFile                          = "Missing file"
	MsgFileParseFailed                      = "Failed to parse file"
	MsgInvalidID                            = "Invalid ID format"
	MsgInvalidActionFormat                  = "Invalid action format"

	// System error messages
	MsgInternalServerError         = "Internal server error occurred"
	MsgCPSActionWaitPrevious       = "Wait for previous checker approval to be completed"
	MsgCPSActionApprovedByThisRole = "Action Approval on your role is already completed"
	MsgServiceUnavailable          = "Service is temporarily unavailable"
	MsgDatabaseError               = "Database operation failed"
	MsgNetworkError                = "Network error occurred"
	MsgTimeoutError                = "Request timeout occurred"
	MsgUnexpectedError             = "An unexpected error occurred"
	MsgConfigurationError          = "Configuration error"
	MsgExternalServiceError        = "External service error"

	// Business logic error messages
	MsgInsufficientBalance = "Insufficient balance"
	MsgTransactionFailed   = "Transaction failed"
	MsgLimitExceeded       = "Limit exceeded"
	MsgOperationNotAllowed = "Operation not allowed"
	MsgMaintenanceMode     = "System is in maintenance mode"

	// Bank related error messages
	MsgAmountTierNotFound        = "Vault's amount tier not found"
	MsgBankNotFound              = "Bank not found"
	MsgBankAlreadyExists         = "Bank already exists"
	MsgBankAlreadyEnabled        = "Bank already enabled"
	MsgBankAlreadyDisabled       = "Bank already disabled"
	MsgBankNameAlreadyExists     = "Bank name already exists"
	MsgBankBICAlreadyExists      = "Bank BIC already exists"
	MsgBankCodeAlreadyExists     = "Bank Code already exists"
	MsgBankFetchFailed           = "Bank fetch failed"
	MsgBankIDRequired            = "Bank ID is required"
	MsgBankNameRequired          = "Bank name is required"
	MsgBankCodeRequired          = "Bank code is required"
	MsgBankBICRequired           = "Bank BIC is required"
	MsgBankAccountLengthRequired = "Bank account length is required"
	MsgBankAccountLengthInvalid  = "Bank account length is invalid"
	MsgInvalidBankName           = "Invalid bank name"
	MsgBankNameTooShort          = "Bank name too short"
	MsgBankNameTooLong           = "Bank name too long"
	MsgBankNameInvalidCharacters = "Bank name contains invalid characters"
	msgGetAllBanksFailed         = "Get all banks failed"
	msgGetAllBanksSuccess        = "Successfully banks retrived"
	msgGetOneBankSuccess         = "Successfully bank retrived"
	msgDeleteBankRequestSuccess  = "Successfully bank delete request created"
	msgGetOneBankFailed          = "Get one bank failed"

	// Donation related error messages
	MsgDonationCategoryNameAlreadyExists = "Category name already exists"
	MsgDonationCompanyNameAlreadyExists  = "Company name already exists"
	MsgCompanyNameAlreadyExists          = "Company name already exists"
	MsgCompanyCodeAlreadyExists          = "Company code already exists"
	MsgDonationTitleAlreadyExists        = "Donation title already exists"
	MsgAccountNumberAlreadyExists        = "company with this account number already exists"
	MsgEmailAlreadyExists                = "Email already exists"
	MsgPhonenumberAlreadyExists          = "Phone number already exists"
	MsgCodeAlreadyExists                 = "Ecommerce merchant with this code/id already exists"
	MsgLogoIsRequired                    = "Company logo is required"
	MsgAccountNumberValidationFailed     = "Account number validation failed"
	MsgAccountNotFound                   = "Account number not found"
	MsgAccountNotActive                  = "Account number is not active"
	MsgDonationIconRequired              = "Icon is required"
	MsgDonationLogoRequired              = "Logo is required"
	MsgDonationImagesRequired            = "Donation images are required"
	MsgDonationInvalidAmount             = "Invalid donation amount"
	MsgDonationInvalidStartDate          = "Invalid start date format"
	MsgDonationInvalidEndDate            = "Invalid end date format"

	// Ecommerce Merchant
	MsgEcommerceMerchantCreated  = "Ecommerce merchant created request submitted successfully"
	MsgEcommerceMerchantUpdated  = "Ecommerce merchant updated request submitted successfully"
	MsgEcommerceMerchantDeleted  = "Ecommerce merchant delete request submitted successfully"
	MsgEcommerceMerchantEnable   = "Ecommerce merchant enable request submitted successfully"
	MsgEcommerceMerchantDisable  = "Ecommerce merchant disable request submitted successfully"
	MsgEcommerceMerchantFetched  = "Ecommerce merchant fetched successfully"
	MsgEcommerceMerchantsFetched = "Ecommerce merchants fetched successfully"
	MsgEcommerceMerchantLookup   = "Ecommerce merchant lookup successfully"

	MsgDonationAccountNotFound    = "Account not found"
	MsgDonationUploadFailed       = "Failed to upload donation files"
	MsgDonationTitleDuplicated    = "Donation title already exists"
	MsgDonationCategoryNotFound   = "Donation category not found"
	MsgDonationCompanyNotFound    = "Donation company not found"
	MsgImageRequired              = "Image is required"
	MsgDonationCompanyNotEnabled  = "The selected donating company is not enabled"
	MsgDonationCategoryNotEnabled = "The selected donation category is not enabled"
	MsgDonationMinIOError         = "MinIO operation failed"
	MsgDonationConfigurationError = "Configuration error"
	MsgDonationTimeoutError       = "Operation timed out"

	// CPS Action related error messages
	MsgPendingCPSActionExists        = "Pending CPS action already exists"
	MsgCPSActionNotFound             = "CPS action not found"
	MsgBpsUserAlreadyEnabled         = "User already enabled"
	MsgBpsUserAlreadyDisabled        = "User already disabled"
	MsgCPSActionNotPending           = "Action is not in pending status"
	MsgCPSActionAlreadyApproved      = "Action already approved"
	MsgCPSActionAlreadyRejected      = "Action already rejected"
	MsgCPSActionAlreadyCanceled      = "Action already canceled"
	MsgCPSActionCreationFailed       = "Failed to create CPS action"
	MsgCPSActionUpdateFailed         = "Failed to update CPS action"
	MsgCPSActionApprovalFailed       = "Failed to approve action"
	MsgCPSActionRejectionFailed      = "Failed to reject action"
	MsgCPSActionInvalidStatus        = "Invalid action status"
	MsgCPSActionInvalidRequestAction = "Invalid request action"
	MsgBankInvalidRequestAction      = "Invalid request action"
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
	MsgProductCodeNotFound        = "Product code not found"
	MsgProductCodeUpdateFailed    = "Failed to update product code"
	MsgProductCodeFetchFailed     = "Failed to fetch product code"
	MsgProductCodeInvalidFormat   = "Invalid product code format"
	MsgProductCodeRequired        = "Product code is required"
	MsgDuplicateProductName       = "Product name already exists"
	MsgDuplicateCBEProductCode    = "CBE product code (PRD) already exists"
	MsgDuplicateCBEIFBProductCode = "CBE IFB product code (PRD) already exists"

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
	MsgPermissionGroupNotFound         = "Permission group not found"
	MsgPermissionGroupAlreadyExists    = "Permission group already exists"
	MsgPermissionGroupCreationFailed   = "Failed to create permission group"
	MsgPermissionGroupUpdateFailed     = "Failed to update permission group"
	MsgPermissionGroupFetchFailed      = "Failed to fetch permission groups"
	MsgPermissionCategoryNotFound      = "Permission category not found"
	MsgPermissionCategoryFetchFailed   = "Failed to fetch permission categories"
	MsgPermissionGroupValidationFailed = "Permission group validation failed"
	MsgPermissionGroupRequired         = "Permission group is required"
	MsgPermissionCatagoryNotFound      = "Permission category not found"

	// Mini App Merchant related error messages
	MsgMiniAppMerchantNotFound       = "Mini app merchant not found"
	MsgMiniAppMerchantCreationFailed = "Failed to create mini app merchant"
	MsgMiniAppMerchantDeletionFailed = "Failed to delete mini app merchant"
	MsgMiniAppMerchantDisableFailed  = "ecommerce merchant is already disabled"
	MsgMiniAppMerchantEnableFailed   = "ecommerce merchant is already enabled"
	MsgMiniAppMerchantInvalidFormat  = "Invalid mini app merchant format"
	MsgMiniAppMerchantRequired       = "Mini app merchant is required"

	// Notification related error messages
	MsgNotificationNotFound       = "Notification not found"
	MsgNotificationCreationFailed = "Failed to create notification"
	MsgNotificationUpdateFailed   = "Failed to update notification"
	MsgNotificationDeletionFailed = "Failed to delete notification"
	MsgNotificationInvalidFormat  = "Invalid notification format"
	MsgNotificationRequired       = "Notification is required"

	// department related error messages
	MsgDepartmentNotFound              = "Department not found"
	MsgDepartmentAlreadyExists         = "Department already exists"
	MsgDepartmentCreationFailed        = "Failed to create department"
	MsgDepartmentUpdateFailed          = "Failed to update department"
	MsgDepartmentInvalidFormat         = "Invalid department format"
	MsgDepartmentRequired              = "Department is required"
	msgGetAllDepartmentsSuccess        = "Successfully got all Departments"
	msgGetDepartmentsSuccess           = "Successfully got a Department"
	MsgDepartmentAlreadyEnabled        = "Department already enabled"
	MsgDepartmentAlreadyDisabled       = "Department already disabled"
	MsgDepartmentInvalidID             = "Department Invalid or Missing id"
	MsgDepartmentWithNameAlreadyExists = "Department with name already exists"
	MsgDepartmentInvalidRequestAction  = "Invalid request action"

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

	// bank vault related error messages
	MsgBankVaultAlreadyEnabled  = "Bank vault already enabled"
	MsgBankVaultAlreadyDisabled = "Bank vault already disabled"
	MsgFailedToBeingTransaction = "Failed to begin transaction"
	Msg

	// vault group category related error messages
	MsgVaultGroupAlreadyEnabled  = "Vault group category already enabled"
	MsgVaultGroupAlreadyDisabled = "Vault group category already disabled"

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
	MsgDuplicateKey                     = "Duplicate key error"
	MsgDuplicateAction                  = "Duplicate action error"
	MsgOneOrMoreInvalidCodes            = "One or more invalid ids"
	MsgInvalidInput                     = "Invalid input provided"
	MsgMissingRequiredFields            = "Missing required fields"
	MsgResourceNotFound                 = "Resource not found"
	MsgActionNameAlreadyExists          = "Action name already exists"
	MsgOnDisablingExistingDeviceControl = "Failed to disable existing device control"
	MsgResourceAlreadyExists            = "Resource already exists"
	MsgResourceUpdateFailed             = "Failed to update resource"
	MsgResourceCreationFailed           = "Failed to create resource"
	MsgResourceDeletionFailed           = "Failed to delete resource"
	MsgResourceFetchFailed              = "Failed to fetch resource"
	MsgResourceInvalidFormat            = "Invalid resource format"
	MsgResourceRequired                 = "Resource is required"
	MsgResourceBusy                     = "Resource is busy"
	MsgResourceLocked                   = "Resource is locked"
	MsgResourceExpired                  = "Resource has expired"
	MsgResourceInvalid                  = "Resource is invalid"
	MsgResourceUnauthorized             = "Resource access unauthorized"
	MsgResourceForbidden                = "Resource access forbidden"
	MsgResourceConflict                 = "Resource conflict"
	MsgResourceTooManyRequests          = "Too many requests for resource"
	MsgResourceUnprocessable            = "Resource cannot be processed"
	MsgResourceNotImplemented           = "Resource operation not implemented"
	MsgResourceBadGateway               = "Resource gateway error"
	MsgResourceServiceUnavailable       = "Resource service unavailable"
	MsgResourceGatewayTimeout           = "Resource gateway timeout"
	MsgResourceNetworkError             = "Resource network error"
	MsgResourceConfigurationError       = "Resource configuration error"
	MsgResourceMaintenanceMode          = "Resource in maintenance mode"
	MsgInvalidRequest                   = "Invalid request data"
	MsgInvalidRequestOnParam            = "Invalid request on page and per page"
	MsgInvalidBankRequest               = "Invalid bank request"
	MsgInvalidRequestBankName           = "invalid format for Name: only letters, numbers, and spaces are allowed"
	MsgInvalidRequestBankCode           = "invalid format for Code: only letters, numbers, and spaces are allowed"
	MsgInvalidRequestBankBIC            = "invalid format for BIC: only letters, numbers, and spaces are allowed"
	MsgInvalidRequestBankType           = "invalid format for Bank Type: only WALLET,BANK OR MFI allowed"
	MsgBankTypeRequired                 = "Bank type is required"
	MsgNoDataProvidedForBankUpdate      = "No data provided for Bank update"
	MsgInvalidToken                     = "Invalid token"
	MsgCustomerPermanentlyDisabled      = "Customer account is permanently disabled"

	// DEPARTMENT RELATED MESSAGES
	MsgInvalidRequestDepartmentName             = "invalid format for department: only letters, numbers, and spaces are allowed"
	MsgInvalidRequestDepartmentPortalCards      = "invalid format for Portal cards: only letters, numbers, and spaces are allowed"
	MsgInvalidRequestDepartmentPermissionGroups = "invalid format for permission group: only letters, numbers, and spaces are allowed"
	MsgInvalidDepartmentPermissionGroup         = "invalid permission group id or permission group with id not found"
	MsgInvalidDepartmentPortalCard              = "invalid portal card id or portal card with id not found"

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
	MsgServiceIdRequered                  = "Service ID is required"

	// Product Code related error messages
	MsgProductCodeParseFailed          = "Failed to parse product code ID"
	MsgProductCodeDatabaseQueryFailed  = "Database query failed for product code"
	MsgProductCodeCountFailed          = "Failed to count product codes"
	MsgProductCodeDatabaseUpdateFailed = "Database update failed for product code"

	// Mini App Merchant related error messages
	MsgMiniAppMerchantCheckPendingFailed     = "Failed to check pending request"
	MsgMiniAppMerchantExistsCheckFailed      = "Failed to check miniapp merchant existence"
	MsgMiniAppMerchantFetchCategoriesFailed  = "Failed to fetch permission categories"
	MsgMiniAppMerchantUnexpectedDBError      = "Unexpected database error occurred"
	MsgMiniAppMerchantMarshalFailed          = "Failed to marshal current action"
	MsgMiniAppMerchantUnmarshalFailed        = "Failed to unmarshal current action"
	MsgMiniAppMerchantCreateFromActionFailed = "Failed to create permission group from action"
	MsgMiniAppMerchantUnmarshalActionFailed  = "Failed to unmarshal action data"
	MsgMiniAppMerchantUpdateFailed           = "Failed to update permission group"
	MsgMiniAppMerchantFetchPermissionsFailed = "Failed to fetch permissions for category"
	MsgMiniAppMerchantDeleteFailed           = "Failed to fetch delete category"

	// Device Version related error messages
	MsgDeviceVersionAlreadyExists   = "Device version already exists"
	MsgDeviceVersionNotFound        = "Device version not found"
	MsgDeviceVersionUpdateFailed    = "Device version update failed"
	MsgDeviceVersionDeleteFailed    = "Device version delete failed"
	MsgDeviceVersionEnableFailed    = "Device version enable failed"
	MsgDeviceVersionDisableFailed   = "Device version disable failed"
	MsgDeviceVersionAlreadyEnabled  = "Device version already enabled"
	MsgDeviceVersionAlreadyDisabled = "Device version already disabled"
	// Notification related error messages
	MsgNotificationMapFailed       = "Failed to map notification to document"
	MsgNotificationAlreadyExists   = "Notification already exists"
	MsgNotificationInsertFailed    = "Failed to insert notification"
	MsgNotificationInvalidIDFormat = "Invalid notification ID format"
	MsgNotificationFetchFailed     = "Failed to fetch notification"
	MsgNotificationDeleteFailed    = "Failed to delete notification"

	// Customer related error messages
	MsgCustomerCountFailed                         = "Failed to count customers"
	MsgCustomerConvertIDFailed                     = "Failed to convert customer ID"
	MsgCustomerNotFound                            = "Customer not found"
	MsgCustomerFetchFailed                         = "Failed to get customer"
	MsgCustomerCIFMustContainOnlyNumbers           = "CIF must contain only numbers"
	MsgCustomerAccountNumberMustContainOnlyNumbers = "Account number must contain only numbers"

	// Account Validation related error messages
	MsgValidationRuleGetFailed                = "Failed to get validation rule"
	MsgValidationRuleMinMaxLengthMismatch     = "Validation failed: min length cannot exceed max length"
	MsgValidationRuleConvertIDFailed          = "Failed to convert validation rule ID"
	MsgValidationRuleUpdateFailed             = "Failed to update validation rule"
	MsgValidationRulePendingActionFetchFailed = "Failed to fetch pending action"
	MsgValidationRuleActionsFetchFailed       = "Failed to fetch actions"

	// Donation related error messages
	MsgDonationCategoryLookupFailed   = "Category lookup failed"
	MsgDonationCategoryIdRequired     = "category id is required"
	MsgDonationCompanyIdRequired      = "company id is required"
	MsgAccountNumberRequired          = "Account Number is required"
	MsgDonationCompanyLookupFailed    = "Company lookup failed"
	MsgDonationLookupFailed           = "Donation lookup failed"
	MsgDonationAlreadyEnabled         = "Donation already enabled"
	MsgDonationAlreadyDisabled        = "Donation already disabled"
	MsgDonationCategoryNameDuplicated = "donation category name already exist"
	MsgNoChangesToUpdate              = "no changes to update"
	MsgBudgetCategoryAlreadyEnabled   = "Budget category already enabled"
	MsgBudgetCategoryAlreadyDisable   = "Budget category already disabled"
	MsgDonationParseStartDateFailed   = "Failed to parse start date"
	MsgDonationParseEndDateFailed     = "Failed to parse end date"
	MsgDonationFetchCompanyFailed     = "Failed to fetch company"
	MsgDonationFetchCategoryFailed    = "Failed to fetch category"
	MsgDonationUpdateFailed           = "Failed to update donation"
	MsgDonationImageUpdateFailed      = "Failed to update donation image"
	MsgDonationImageDeleteFailed      = "Failed to delete donation image"
	MsgDonationImageAddFailed         = "Failed to add donation image"

	MsgCannotGetRole         = "CPS user's role cannot be found"
	MsgJobTitleAlreadyExists = "Job title already exists"
	MsgRoleAlreadyExists     = "Role already exists"

	// Bank related error messages
	MsgBankFileParseFailed              = "Failed to parse bank file"
	MsgBankRejectionPayloadDecodeFailed = "Failed to decode rejection payload"
	MsgBankLogoUpdateFailed             = "Failed to update bank logo"
	MsgBankUpdateFailed                 = "Failed to update bank"

	// CPS Action related error messages
	MsgCPSActionRejectionPayloadDecodeFailed = "Failed to decode rejection payload"
	MsgCPSActionRejectionReason              = "rejection_reason is required"
	MsgCPSActionRejectionReasonLength        = "rejection_reason must be between 10 and 300 characters"

	// Permission related error messages
	MsgPermissionGroupRequestCreationFailed = "Failed to create permission group request"
	MsgPermissionGroupRequestUpdateFailed   = "Failed to update permission group request"

	// Department related error messages
	MsgDepartmentCreateRequestSuccess = "Department created Request successfully"
	MsgDepartmentCreateRequestFail    = "Department created Request successfully"

	MsgDepartmentUpdateCPSActionCreated = "Department update CPS action created successfully"

	// Budget Category related error messages
	MsgBudgetCategoryUpdateSuccess = "Budget category updated successfully"

	// Donation Service related error messages
	MsgDonationCategoryIconUploadSuccess             = "Successfully uploaded donation category icon"
	MsgDonationCategoryIconUpdateSuccess             = "Successfully uploaded updated donation category icon"
	MsgDonationCPSActionApprovedSuccess              = "Successfully approved and authorized CPS action"
	MsgDonationCompanyLogoUploadSuccess              = "Successfully uploaded donation company logo"
	MsgBudgetIconRequestSubmittedForApprovalSuccess  = "Successfully budget icon request submitted for approval"
	MsgBudgetCheckerActionApprovedSuccess            = "Successfully budget icon checker submitted for approval"
	MsgBudgetColorUpdateSubmittedForApprovalSuccess  = "Color creation request submitted for approval"
	MsgBudgetIconUpdateSubmittedForApprovalSuccess   = "Successfully budget icon update submitted for approval"
	MsgBudgetColorRequestSubmittedForApprovalSuccess = "Color creation request submitted for approval"
	MsgBudgetIconsFetchedSuccessfully                = "Successfully fetched budget icons"
	MsgBudgetColorsFetchedSuccessfully               = "Successfully fetched budget colors"

	// Unified Budget Operations Messages
	MsgBudgetRequestSubmittedForApprovalSuccess         = "Budget creation request submitted for approval"
	MsgBudgetUpdateSubmittedForApprovalSuccess          = "Budget update request submitted for approval"
	MsgBudgetCategoryFetchedSuccessfully                = "Budget category fetched successfully"
	MsgBudgetCategoriesFetchedSuccessfully              = "Budget categories fetched successfully"
	MsgBudgetCategoryRequestSubmittedForApprovalSuccess = "Budget category creation request submitted for approval"
	MsgBudgetCategoryUpdateSubmittedForApprovalSuccess  = "Budget category update request submitted for approval"
	MsgBudgetCategoryDeleteSubmittedForApprovalSuccess  = "Budget category deletion request submitted for approval"
	MsgBudgetCategoryEnableSubmittedForApprovalSuccess  = "Budget category enable request submitted for approval"
	MsgBudgetCategoryDisableSubmittedForApprovalSuccess = "Budget category disable request submitted for approval"
	MsgDonationCompanyLogoUpdateSuccess                 = "Successfully uploaded updated donation company logo"
	MsgDonationImageUploadSuccess                       = "Successfully uploaded donation image"
	MsgDonationImagesUpdateSuccess                      = "Successfully uploaded updated donation images"

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
	MsgBPSUserApprovedSuccess   = "BPS user approved successfully"
	MsgUserNotHaveLinkedAccount = "User Doesn't have any linked account"

	MsgUserUnlinkFailed = "User Doesn't have any linked account"
	// Ad Service related success messages
	MsgAdvertCreatedSuccess = "Advert created successfully"
	MsgAdvertUpdateSuccess  = "Advert updated successfully"

	// Ad Service related error messages
	MsgAdvertCreateError     = "Advert creation failed"
	MsgAdvertUpdateError     = "Advert update failed"
	MsgAdvertAlreadyEnabled  = "Advert already enabled"
	MsgAvatarAlreadyEnabled  = "Avatar already enabled"
	MsgAdvertAlreadyDisabled = "Advert already disabled"
	MsgAvatarAlreadyDisabled = "Avatar already disabled"
	MsgAdvertNotFound        = "Advert not found"

	// Account Validation Service related error messages
	MsgValidationRuleApprovedSuccess = "Validation rule approved successfully"
	MsgValidationRuleSuccessFech     = "Account fetched successfully"

	MsgPendingActionExists  = "Pending action exists"
	MsgDuplicateColorExists = "Duplicate color exists"

	MsgInvalidRequestBody      = "Invalid request body"
	MsgErrorBulkDisabled       = "Error Disabled password rule"
	MsgInvalidPaginationParams = "Invalid pagination parameters"

	MsgBranchSuccessfullyRetrieved   = "Branch retrieved successfully"
	MsgBranchesSuccessfullyRetrieved = "Branches retrieved successfully"
	MsgBranchesSuccessfullyEnabled   = "Branches enabled successfully"
	MsgBranchesSuccessfullyDisabled  = "Branches disabled successfully"
	MsgEnableBranchesRequestSent     = "Request to enable branches sent successfully"
	MsgDisableBranchesRequestSent    = "Request to disable branches sent successfully"

	MsgRegionSuccessfullyRetrieved  = "Region retrieved successfully"
	MsgRegionsSuccessfullyRetrieved = "Regions retrieved successfully"
	MsgRegionsSuccessfullyEnabled   = "Regions enabled successfully"
	MsgRegionsSuccessfullyDisabled  = "Regions disabled successfully"
	MsgEnableRegionsRequestSent     = "Request to enable regions sent successfully"
	MsgDisableRegionsRequestSent    = "Request to disable regions sent successfully"

	MsgDistrictSuccessfullyRetrieved  = "District retrieved successfully"
	MsgDistrictsSuccessfullyRetrieved = "Districts retrieved successfully"
	MsgDistrictsSuccessfullyEnabled   = "Districts enabled successfully"
	MsgDistrictsSuccessfullyDisabled  = "Districts disabled successfully"
	MsgEnableDistrictsRequestSent     = "Request to enable districts sent successfully"
	MsgDisableDistrictsRequestSent    = "Request to disable districts sent successfully"

	MsgCitySuccessfullyRetrieved   = "City retrieved successfully"
	MsgCitiesSuccessfullyRetrieved = "Cities retrieved successfully"
	MsgCitiesSuccessfullyEnabled   = "Cities enabled successfully"
	MsgCitiesSuccessfullyDisabled  = "Cities disabled successfully"
	MsgEnableCitiesRequestSent     = "Request to enable cities sent successfully"
	MsgDisableCitiesRequestSent    = "Request to disable cities sent successfully"

	MsgInvalidDistrict                   = "Invalid district"
	MsgInvalidRegion                     = "Invalid region"
	MsgInvalidDistrictOrRegionCodeLength = "Region and District id must be at least 3 characters long"
	MsgBranchCodeRequired                = "Branch id is required"
	MsgBranchAlreadyEnabled              = "Branch already enabled"
	MsgBranchAlreadyDisabled             = "Branch already disabled"
	MsgDistrictCodeRequired              = "District id is required"
	MsgDistrictAlreadyEnabled            = "District already enabled"
	MsgDistrictAlreadyDisabled           = "District already disabled"
	MsgRegionCodeRequired                = "Region id is required"
	MsgRegionAlreadyEnabled              = "Region already enabled"
	MsgRegionAlreadyDisabled             = "Region already disabled"
	MsgCityCodeRequired                  = "City id is required"
	MsgCityAlreadyEnabled                = "City already enabled"
	MsgCityAlreadyDisabled               = "City already disabled"

	MsgBranchNotFound         = "Branch not found"
	MsgDistrictNotFound       = "District not found"
	MsgRegionNotFound         = "Region not found"
	MsgCityNotFound           = "City not found"
	MsgInvalidInputParameters = "Invalid input parameters provided"
	MsgMissingOrInvalidImage  = "Missing or invalid image"
	MsgCannotEnableDistrict   = "Cannot enable district, its region is blocked"
	MsgCannotEnableCity       = "Cannot enable city, its district is blocked"
	MsgCannotEnableBranch     = "Cannot enable branch, its city is already blocked"

	// Password Rule
	MsgFetchAllPasswordRules = "Password Rules Successfully Fetched"
	MsgUpdatePasswordRule    = "Update request submitted for approval"

	// bank vault
	MsgCannotDeleteActiveBankVault             = "Cannot delete active bank vault"
	MsgBankVaultProductAlreadyDeleted          = "Bank vault product already deleted"
	MsgDuplicatebankProduct                    = "There is an existing bank product"
	MsgNoBankProductFound                      = "No bank Product found"
	MsgCannotDeleteBankProduct                 = "Cannot delete bank product"
	MsgcannotEnableOrDisableDeletedBankProduct = "Cannot enable or disable deleted bank product"
	MsgDuplicateGroupVaultCategory             = "There is an existing group vault category"

	// Fayda Account
	MsgFaydaAccountEnableCreatedSuccessfully  = "Fayda Account enable action submitted successfully"
	MsgFaydaAccountDisableCreatedSuccessfully = "Fayda Account disable action submitted successfully"
	MsgUserFaydaAccountAlreadyEnabled         = "Fayda user account already enabled"
	MsgUserFaydaAccountAlreadyDisabled        = "Fayda user account already disabled"
	MsgNotFaydaUser                           = "This user is not fayda user"

	MsgAlreadyEnabled  = "Already enabled"
	MsgAlreadyDisabled = "Already disabled"

	MsgBucketNotFound = "bucket not found"

	MsgNoUpdateDetected = "No update detected"

	MsgNewsCategoryInvalidID             = "Invalid or missing news category id"
	MsgNewsTagInvalidID                  = "Invalid or missing news tag id"
	MsgNewsCategoryWithNameAlreadyExists = "News category with this name already exists"
	MsgNewsTagWithNameAlreadyExists      = "News tag with this name already exists"
	MsgNewsCategoryDeletedSuccess        = "News category deleted request sent successfully"
	MsgNewsCategoryUpdatedSuccess        = "News category updated request sent successfully"
	MsgNewsCategoryFetchedSuccess        = "News category fetched successfully"
	MsgNewsCategoryCreatedSuccess        = "News category created request sent successfully"
	MsgNewsTagDeletedSuccess             = "News tag deleted request sent successfully"
	MsgNewsTagUpdatedSuccess             = "News tag updated request sent successfully"
	MsgNewsTagFetchedSuccess             = "News tag fetched successfully"
	MsgNewsTagCreatedSuccess             = "News tag created request sent successfully"

	// Action Role success messages
	MsgActionRolesFetchedSuccess       = "Action roles fetched successfully"
	MsgActionRoleFetchedSuccess        = "Action role fetched successfully"
	MsgActionRoleCreateRequestCreated  = "Action role create request sent successfully"
	MsgActionRoleUpdateRequestCreated  = "Action role update request sent successfully"
	MsgActionRoleEnableRequestCreated  = "Action role enable request sent successfully"
	MsgActionRoleDisableRequestCreated = "Action role disable request sent successfully"

	// Sitota related
	MsgAllSitotasRetrievedSuccess = "All sistotas retrieved successfully"
	MsgSitotaRetrievedSuccess     = "Sitota retrieved successfully"
	MsgSitotaRequired             = "Sitota ID is required"

	// Encryption
	MsgEncryptionSuccessfully = "Encryption generated successfully"

	// bps action role related
	MsgBpsActionRoleNotFound = "BPS action role not found"
	MsgCpsActionRoleNotFound = "CPS action role not found"
	MsgRoleNotFound          = "Role not found"
	MsgExistingUsedRole      = "You used Existing Role data in the request"
	// transaction related messages
	MsgTransactionsRetrievedSuccess  = "Transactions retrieved successfully"
	MsgTransactionRetrievedSuccess   = "Transaction retrieved successfully"
	MsgTransactionIDRequired         = "Transaction ID is required"
	MsgTransactionIdentifierRequired = "Transaction Identifier is required"

	// event merchant message
	MsgEventMerchantCreatedSuccessfully  = "Event merchant created successfully"
	MsgEventMerchantInvalidID            = "merchant id invalid or empty"
	MsgEventMerchantInvalidType          = "merchant type invalid or empty"
	MsgEventMerchantInvalidMethod        = "merchant method invalid or empty"
	MsgEventMerchantInvalidName          = "merchant name invalid or empty"
	MsgEventMerchantInvalidAccountNumber = "merchant account number invalid or empty"
	MsgEventMerchantInvalidEmail         = "merchant email invalid or empty"
	MsgEventMerchantInvalidPhoneNumber   = "merchant phone number invalid or empty"
	MsgEventMerchantEnableFailed         = "event merchant already enabled"
	MsgEventMerchantDisableFailed        = "event merchant already disabled"
	MsgEventMerchantNotFound             = "event merchant not found"
	MsgEventMerchantDisabledSuccessfully = "event merchant disabled successfully"
	MsgEventMerchantEnabledSuccessfully  = "event merchant enabled successfully"
	MsgEventMerchantUpdatedSuccessfully  = "event merchant updated successfully"
	MsgEventMerchantDeletedSuccessfully  = "event merchant deleted successfully"
	MsgEventMerchantFetchedSuccessfully  = "event merchant fetched successfully"

	// access list segmentation related messages
	MsgAccessListSegmentationCreatedSuccessfully    = "Access list segmentation created successfully"
	MsgAccessListSegmentationUpdatedSuccessfully    = "Access list segmentation updated successfully"
	MsgAccessListSegmentationEnabledSuccessfully    = "Access list segmentation enabled successfully"
	MsgAccessListSegmentationDisabledSuccessfully   = "Access list segmentation disabled successfully"
	MsgAccessListSegmentationsRetrievedSuccessfully = "Access list segmentations retrieved successfully"

	MsgAccessListSegmentationRetrievedSuccessfully = "Access list segmentation retrieved successfully"
	MsgAccessListSegmentationNotFound              = "Access list segmentation not found"
	MsgAccessListSegmentationAlreadyEnabled        = "Access list segmentation is already enabled"
	MsgAccessListSegmentationAlreadyDisabled       = "Access list segmentation is already disabled"
	MsgAccessListSegmentationInvalidID             = "Invalid access list segmentation ID"
	MsgAccessListSegmentationIDsRequired           = "Access list segmentation IDs are required"
	MsgAccessListSegmentationNameAlreadyExists     = "Access list segmentation name already exists"
	MsgCustomerSegmentationCodeNotFound            = "Customer segmentation code not found"
	MsgServiceIdRequired                           = "Service ID is required"
)
