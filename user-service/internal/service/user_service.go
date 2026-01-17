package service

import (
	"context"

	"go.uber.org/zap"
	"shared.local/pkg/logger"
	"shared.local/pkg/trace"

	"user-service/internal/repo"
)

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) string {
	logger.L().Info("get user in service",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.String("user_id", id),
	)

	return s.repo.GetUser(ctx, id)
}
