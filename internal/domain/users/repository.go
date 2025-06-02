package users

import (
    "context"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities"
)

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*entities.User, error)
}