package account

import "context"


type AccountAPIPort interface {
	LookupAccountByPhone(ctx context.Context, phoneNumber string,PhoneLookupUrl string ) (bool, error)
}