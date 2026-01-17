package repo

import "context"

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) GetUser(ctx context.Context, id string) string {
	return "user-" + id
}
