package service

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"shared.local/pkg/logger"
	"shared.local/pkg/trace"

	"tenant-service/internal/model"
	"tenant-service/internal/repo"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrDomainExists   = errors.New("domain already exists")
)

type TenantService struct {
	repo *repo.TenantRepo
}

func NewTenantService(repo *repo.TenantRepo) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) GetTenantByID(ctx context.Context, id uint) (*model.Tenant, error) {
	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
			return nil, err
	}
	if tenant == nil {
			return nil, ErrTenantNotFound
	}
	return tenant, nil
}

func (s *TenantService) GetTenantByDomain(ctx context.Context, Domain string) (*model.Tenant, error) {
	logger.L().Info("get tenant in service by domain",
		zap.String("domain", Domain),
		zap.String("trace_id", trace.FromContext(ctx)),
	)

	tenant, err := s.repo.GetTenantByDomain(ctx, Domain)

	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, ErrTenantNotFound
	}

	return tenant, nil
}
