package service

import (
	"context"

	"shared.local/pkg/logger"
	"shared.local/pkg/trace"
	"go.uber.org/zap"

	"tenant-service/internal/repo"
)

type TenantService struct {
	repo *repo.TenantRepo
}

func NewTenantService(repo *repo.TenantRepo) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) GetTenant(ctx context.Context, id string) string {
	logger.L().Info("get tenant in service",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.String("tenant_id", id),
	)

	return s.repo.GetTenant(ctx, id)
}
