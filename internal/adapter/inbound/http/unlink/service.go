package unlink

import (
	"context"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)


type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction any) error
}

type unlinkAdapter struct {
	logger utils.Logger
	unlinkApp UnlinkAccount
}

func InitAdapterUnlinkService(w http.ResponseWriter,r *http.Request) {
	return 
}

func (ua *unlinkAdapter) GetUserByAccount(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error){
	var req 
	user,err := ua.unlinkApp.GetUserByAccount(ctx,)
}
func (ua *unlinkAdapter) UnlinkUserCif(w http.ResponseWriter,r *http.Request) error{

}
