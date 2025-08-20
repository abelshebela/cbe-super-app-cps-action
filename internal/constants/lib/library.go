package lib

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func CpsModelBuilder(unique string, makerUser types.UserContext, prevAction, currentAction any, requestAction, actionType string) model.CPSAction {

	return model.CPSAction{
		ActionCode:       local_util.GenerateActionCode(),
		UniqueId:         unique,
		MakerID:          makerUser.UserID,
		MakerName:        makerUser.FullName,
		MakerPhoneNumber: makerUser.PhoneNumber,
		Department:       makerUser.Department,
		PreviousAction:   prevAction,
		CurrentAction:    currentAction,
		ActionStatus:     string(constants.Pending),
		ActionType:       actionType,
		RequestAction:    requestAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		MakerActionTime:  time.Now(),
	}
}

func GoRoutinBaker(opts types.BakerOptions, tasks ...func()) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	if opts.Sequential {
		for _, task := range tasks {
			if opts.UseMutex {
				mu.Lock()
				task()
				mu.Unlock()
			} else {
				task()
			}
		}
		return
	}

	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			if opts.UseMutex {
				mu.Lock()
				t()
				mu.Unlock()
			} else {
				t()
			}
		}(task)
	}
	wg.Wait()
}

func FilterBuilder(filterParam types.Filter, searchKeys bson.M, allowedKeys []string) (bson.M, int64, int64) {
	var skip, limit int64
	filter := bson.M{}

	if filterParam.Search != "" {
		for key, value := range searchKeys {
			filter[key] = value
		}
	}

	if filterParam.Filters != nil {

		enhancedFilter := local_util.BuildMongoFilterWithKeys(filterParam.Filters, allowedKeys)
		for key, value := range enhancedFilter {
			filter[key] = value
		}
	}

	skip = int64((filterParam.Page - 1) * filterParam.PerPage)
	limit = int64(filterParam.PerPage)

	return filter, skip, limit
}
