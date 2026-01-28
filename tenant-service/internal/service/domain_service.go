package service

import (
	"context"
	"tenant-service/internal/errors"
	"tenant-service/internal/model/dto"
	"tenant-service/internal/model/entity"
	"tenant-service/internal/repo"

	"go.uber.org/zap"
	"shared.local/pkg/logger"
	"shared.local/pkg/trace"
)

type DomainService struct {
	domainRepo *repo.DomainRepo
	tenantRepo *repo.TenantRepo
}

func NewDomainService(domainRepo *repo.DomainRepo, tenantRepo *repo.TenantRepo) *DomainService {
	return &DomainService{domainRepo: domainRepo, tenantRepo: tenantRepo}
}

func (s *DomainService) CreateDomain(ctx context.Context, req dto.CreateDomainRequest) (*entity.Domain, error) {
	logger.L().Info("create Domain service",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("body", req),
	)

	tenant, err := s.tenantRepo.GetTenantByID(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, apperror.ErrTenantNotFound
	}

	// 定义查询实例
	status := uint8(1) // 处理 Status：如果未传值，使用数据库默认值 1
	if req.Status != nil {
		status = *req.Status
	}
	domain := &entity.Domain{
		TenantID: req.TenantID,
		Domain:   req.Domain,
		Status:   status,
	}

	if err := s.domainRepo.CreateDomain(ctx, domain); err != nil {
		return nil, err
	}

	return domain, nil
}
