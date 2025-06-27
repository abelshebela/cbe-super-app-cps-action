package entity

import "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"

type CustomerRespose struct {
	Page      int            `json:"page"`
	Customers []*member.User `json:"customer"`
	Limit     int            `json:"limit"`
	Total     int64          `json:"total"`
}
