package service

import (
	"context"
	apperror "tenant-service/internal/errors"
	"tenant-service/internal/model/dto"
	"tenant-service/internal/model/entity"
	"tenant-service/internal/repo"

	"go.uber.org/zap"
	"shared.local/pkg/logger"
	"shared.local/pkg/pagination"
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

// ListDomains 分页查询域名列表
func (s *DomainService) ListDomains(ctx context.Context, req dto.ListDomainsRequest) (*dto.ListDomainsResponse, error) {
	logger.L().Info("list domains service",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.Any("tenantID", req.TenantID),
		zap.String("domain", req.Domain),
	)

	offset, limit := pagination.CalculatePaginationParams(req.Page, req.PageSize) // 计算分页参数

	domains, total, err := s.domainRepo.List(ctx, offset, limit, req.TenantID, req.Domain) // 查询数据库
	if err != nil {
		return nil, err
	}

	domainResponses := make([]dto.DomainResponse, 0, len(domains)) // 转换为响应 DTO
	for _, d := range domains {
		domainResponses = append(domainResponses, dto.DomainResponse{
			ID:        d.ID,
			Domain:    d.Domain,
			Status:    d.Status,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}

	totalPages := pagination.CalculateTotalPages(total, req.PageSize) // 计算总页数

	return &dto.ListDomainsResponse{
		Items:      domainResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
