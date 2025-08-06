package types

import "time"

type Address struct {
	Zone        string `json:"zone" bson:"zone"`
	Wereda      string `json:"wereda" bson:"wereda"`
	Kebele      string `json:"kebele" bson:"kebele"`
	Region      string `json:"region" bson:"region"`
	City        string `json:"city" bson:"city"`
	SubCity     string `json:"sub_city" bson:"sub_city"`
	StreetName  string `json:"street_name" bson:"street_name"`
	HouseNumber string `json:"house_number" bson:"house_number"`
}

type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_created_at" bson:"last_pin_created_at"`
}

type RegistrationRecord struct {
	ID          string
	PhoneNumber string
	DeviceUUID  string
	Platform    string
	FullName    string
	OTP         string
	OTPFor      string
	Status      string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	Attempts    int
	MaxAttempts int
}

type UserInfo struct {
	UserID      string
	FullName    string
	PhoneNumber string
	Action      string
}
