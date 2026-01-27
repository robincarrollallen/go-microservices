package service

import (
	"context"
	apperror "tenant-service/internal/errors"

	"go.uber.org/zap"

	"shared.local/pkg/logger"
	"shared.local/pkg/trace"

	"tenant-service/internal/model/dto"
	"tenant-service/internal/model/entity"
	"tenant-service/internal/repo"
)

type TenantService struct {
	repo *repo.TenantRepo
}

func NewTenantService(repo *repo.TenantRepo) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) CreateTenant(ctx context.Context, req dto.CreateTenantRequest) (*entity.Tenant, error) {
	logger.L().Info("create tenant service",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("body", req),
	)

	// 检查 Name 是否已存在
	existingTenant, err := s.repo.GetTenantByName(ctx, req.Name)
	// 如果 Name 已存在，返回错误
	if existingTenant != nil {
		return nil, apperror.ErrNameExists
	}
	// 如果获取 Name 失败，返回错误
	if err != nil {
		return nil, err
	}

	// 处理 Status：如果未传值，使用数据库默认值 1
	status := uint8(1)
	if req.Status != nil {
		status = *req.Status
	}

	tenant := &entity.Tenant{
		Name:   req.Name,
		Status: status,
	}

	// 如果有 Domains，使用事务处理
	if len(req.Domains) > 0 {
		return s.repo.CreateTenantWithDomains(ctx, tenant, req.Domains)
	}

	// 否则普通创建
	if err := s.repo.CreateTenant(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *TenantService) GetTenantByID(ctx context.Context, id uint) (*dto.TenantResponse, error) {
	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, apperror.ErrTenantNotFound
	}
	return tenant, nil
}

func (s *TenantService) GetTenantByDomain(ctx context.Context, Domain string) (*dto.TenantResponse, error) {
	logger.L().Info("get tenant in service by domain",
		zap.String("domain", Domain),
		zap.String("trace_id", trace.FromContext(ctx)),
	)

	tenant, err := s.repo.GetTenantByDomain(ctx, Domain)

	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, apperror.ErrTenantNotFound
	}

	return tenant, nil
}
