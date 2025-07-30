package mappers

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

// DomainToModel maps domain.User to model.User
func DomainToUserModel(u *domain.User) (*model.User, error) {
	var id bson.ObjectID
	if u.ID == "" {
		id = bson.NewObjectID()
	} else {
		ID, err := bson.ObjectIDFromHex(u.ID)

		if err != nil {
			return nil, err
		}

		id = ID
	}

	return &model.User{
		ID:                id,
		UserCode:          u.UserCode,
		FullName:          u.FullName,
		MotherName:        u.MotherName,
		Nationality:       u.Nationality,
		BirthDate:         u.BirthDate,
		ResidentialStatus: u.ResidentialStatus,
		IssuedDate:        u.IssuedDate,
		PhoneNumber:       u.PhoneNumber,
		Gender:            model.Gender(u.Gender),
		MaritalStatus:     model.MaritalStatus(u.MaritalStatus),
		Fayda: struct {
			FaydaID          string `json:"id_number" bson:"id_number"`
			FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
			EmploymentStatus string `json:"employment_status" bson:"employement_status"`
			EmployerName     string `json:"employer_name" bson:"employer_name"`
			IssuedBy         string `json:"issued_by" bson:"issued_by"`
			MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
		}{
			FaydaID:          u.Fayda.FaydaID,
			FaydaAccessToken: u.Fayda.FaydaAccessToken,
			EmploymentStatus: u.Fayda.EmploymentStatus,
			EmployerName:     u.Fayda.EmployerName,
			IssuedBy:         u.Fayda.IssuedBy,
			MonthlyIncome:    u.Fayda.MonthlyIncome,
		},

		Address:          model.Address(u.Address),
		DocumentFront:    u.DocumentFront,
		DocumentBack:     u.DocumentBack,
		Photo:            u.Photo,
		Signature:        u.Signature,
		Avatar:           u.Avatar,
		Email:            u.Email,
		PushToken:        u.PushToken,
		Realm:            model.Realm(u.Realm),
		PermissionGroup:  stringIDsToObjectIDs(u.PermissionGroup),
		Permissions:      stringIDsToObjectIDs(u.Permissions),
		IsAccountBlocked: u.IsAccountBlocked,
		IsAccountLinked:  u.IsAccountLinked,
		MemberType:       model.MemberType(u.MemberType),
		RegistrationType: model.RegistrationType(u.RegistrationType),
		AccountStatus:    model.AccountStatus(u.AccountStatus),
		KYCLevel:         u.KYCLevel,

		KYC: struct {
			KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
			KYCStatus            model.KYCStatus     `json:"kyc_status" bson:"kyc_status"`
			KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
			KYCIsApproved        bool                `json:"kyc_approved" bson:"kyc_approved"`
			KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by"`
			// Level                int                 `json:"level" bson:"level"`
		}{
			KYCRejectReasonField: u.KYC.KYCRejectReasonField,
			KYCStatus:            model.KYCStatus(u.KYC.KYCStatus),
			KYCRejectReason:      u.KYC.KYCRejectReason,
			KYCIsApproved:        u.KYC.KYCIsApproved,
			KYCActivityBy:        u.KYC.KYCActivityBy,
			// Level:                0,
		},

		IsBranchApproved:   u.IsBranchApproved,
		IsVerified:         u.IsVerified,
		IsBlocked:          u.IsBlocked,
		IsAccessRestricted: u.IsAccessRestricted,
		BlockedAt:          u.BlockedAt,
		RegisterBy:         struct{}{},
		LoginAttemptCount:  u.LoginAttemptCount,
		LastLoginAttempt:   u.LastLoginAttempt,
		LastOnlineDate:     u.LastOnlineDate,
		LastLogin:          u.LastLogin,

		BPSStatus:  model.BPSStatus(u.BPSStatus),
		LoginPIN:   model.LoginPIN(u.LoginPIN),
		DeviceUUID: u.DeviceUUID,

		Device: struct {
			DevicePlatform      string    `json:"device_platform" bson:"device_platform"`
			AppVersion          string    `json:"app_version" bson:"app_version"`
			APPInstallationDate time.Time `json:"application_installation_date" bson:"application_installation_date"`
		}{
			DevicePlatform:      u.Device.DevicePlatform,
			AppVersion:          u.Device.AppVersion,
			APPInstallationDate: u.Device.APPInstallationDate,
		},

		CustomerNumber:        u.CustomerNumber,
		InitialLinkedDate:     u.InitialLinkedDate,
		PrimaryAuthentication: model.PrimaryAuthentication(u.PrimaryAuthentication),
		DeviceStatus:          model.DeviceStatus(u.DeviceStatus),
		Enabled:               u.Enabled,
		IsDeleted:             u.IsDeleted,
		PINChangedAt:          u.PINChangedAt,
		OTPLastTriedAt:        u.OTPLastTriedAt,
		OTPLastVerifiedAt:     u.OTPLastVerifiedAt,
		OTPVerifyCount:        u.OTPVerifyCount,
		InitialiLinkedAt:      u.InitialLinkedAt,
		CreatedAt:             u.CreatedAt,
		DeletedAt:             u.DeletedAt,
		LastModifiedAt:        u.LastModifiedAt,
	}, nil
}

// ModelToDomain maps model.User to domain.User
func ModelToUserDomain(m *model.User) *domain.User {
	return &domain.User{
		ID:                m.ID.Hex(),
		UserCode:          m.UserCode,
		FullName:          m.FullName,
		MotherName:        m.MotherName,
		Nationality:       m.Nationality,
		BirthDate:         m.BirthDate,
		ResidentialStatus: m.ResidentialStatus,
		IssuedDate:        m.IssuedDate,
		PhoneNumber:       m.PhoneNumber,
		Gender:            member.Gender(m.Gender),
		MaritalStatus:     member.MaritalStatus(m.MaritalStatus),

		Fayda: domain.FaydaInfo{
			FaydaID:          m.Fayda.FaydaID,
			FaydaAccessToken: m.Fayda.FaydaAccessToken,
			EmploymentStatus: m.Fayda.EmploymentStatus,
			EmployerName:     m.Fayda.EmployerName,
			IssuedBy:         m.Fayda.IssuedBy,
			MonthlyIncome:    m.Fayda.MonthlyIncome,
		},

		Address:          member.Address(m.Address),
		DocumentFront:    m.DocumentFront,
		DocumentBack:     m.DocumentBack,
		Photo:            m.Photo,
		Signature:        m.Signature,
		Avatar:           m.Avatar,
		Email:            m.Email,
		PushToken:        m.PushToken,
		Realm:            member.Realm(m.Realm),
		PermissionGroup:  objectIDsToStrings(m.PermissionGroup),
		Permissions:      objectIDsToStrings(m.Permissions),
		IsAccountBlocked: m.IsAccountBlocked,
		IsAccountLinked:  m.IsAccountLinked,
		MemberType:       member.MemberType(m.MemberType),
		RegistrationType: member.RegistrationType(m.RegistrationType),
		AccountStatus:    member.AccountStatus(m.AccountStatus),
		KYCLevel:         m.KYCLevel,

		KYC: domain.KYCInfo{
			KYCRejectReasonField: m.KYC.KYCRejectReasonField,
			KYCStatus:            member.KYCStatus(m.KYC.KYCStatus),
			KYCRejectReason:      m.KYC.KYCRejectReason,
			KYCIsApproved:        m.KYC.KYCIsApproved,
			KYCActivityBy:        m.KYC.KYCActivityBy,
		},

		IsBranchApproved:   m.IsBranchApproved,
		IsVerified:         m.IsVerified,
		IsBlocked:          m.IsBlocked,
		IsAccessRestricted: m.IsAccessRestricted,
		BlockedAt:          m.BlockedAt,
		RegisterBy:         struct{}{}, // or actual mapping if needed
		LoginAttemptCount:  m.LoginAttemptCount,
		LastLoginAttempt:   m.LastLoginAttempt,
		LastOnlineDate:     m.LastOnlineDate,
		LastLogin:          m.LastLogin,
		BPSStatus:          member.BPSStatus(m.BPSStatus),
		LoginPIN:           member.LoginPIN(m.LoginPIN),
		DeviceUUID:         m.DeviceUUID,
		Device: domain.Device{
			DevicePlatform:      m.Device.DevicePlatform,
			AppVersion:          m.Device.AppVersion,
			APPInstallationDate: m.Device.APPInstallationDate,
		},
		CustomerNumber:        m.CustomerNumber,
		InitialLinkedDate:     m.InitialLinkedDate,
		PrimaryAuthentication: member.PrimaryAuthentication(m.PrimaryAuthentication),
		DeviceStatus:          member.DeviceStatus(m.DeviceStatus),
		Enabled:               m.Enabled,
		IsDeleted:             m.IsDeleted,
		PINChangedAt:          m.PINChangedAt,
		OTPLastTriedAt:        m.OTPLastTriedAt,
		OTPLastVerifiedAt:     m.OTPLastVerifiedAt,
		OTPVerifyCount:        m.OTPVerifyCount,
		InitialLinkedAt:       m.InitialiLinkedAt,
		CreatedAt:             m.CreatedAt,
		DeletedAt:             m.DeletedAt,
		LastModifiedAt:        m.LastModifiedAt,
	}
}

func stringIDsToObjectIDs(ids []string) []bson.ObjectID {
	var result []bson.ObjectID
	for _, id := range ids {
		objId, err := bson.ObjectIDFromHex(id)
		if err == nil {
			result = append(result, objId)
		}
	}
	return result
}

func objectIDsToStrings(ids []bson.ObjectID) []string {
	var result []string
	for _, id := range ids {
		result = append(result, id.Hex())
	}
	return result
}
