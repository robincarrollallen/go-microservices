package repo

import (
	"context"
	"errors"
	"tenant-service/internal/model/dto"
	"time"

	"tenant-service/internal/model/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"shared.local/pkg/logger"
	"shared.local/pkg/trace"
)

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) CreateTenant(ctx context.Context, tenant *entity.Tenant) error {
	logger.L().Info("create tenant from repo",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("tenant", tenant),
	)

	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *TenantRepo) GetTenantByID(ctx context.Context, id uint) (*dto.TenantResponse, error) {
	logger.L().Info("get tenant by id from repo",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Uint("id", id),
	)

	var tenant entity.Tenant
	err := r.db.WithContext(ctx).First(&tenant, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	response := &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}

	return response, nil
}

func (r *TenantRepo) GetTenantByName(ctx context.Context, name string) (*dto.TenantResponse, error) {
	var tenant entity.Tenant
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	response := &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}

	return response, nil
}

func (r *TenantRepo) GetTenantByDomain(ctx context.Context, domain string) (*dto.TenantResponse, error) {
	var domainRecord entity.Domain
	err := r.db.WithContext(ctx).Where("domain = ? AND status = 1", domain).
		First(&domainRecord).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var tenant entity.Tenant
	err = r.db.WithContext(ctx).First(&tenant, domainRecord.TenantID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	response := &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}

	return response, nil
}

func (r *TenantRepo) Update(ctx context.Context, tenant *entity.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *TenantRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Tenant{}, id).Error
}

func (r *TenantRepo) CreateTenantWithDomains(ctx context.Context, tenant *entity.Tenant, domains []string) (*entity.Tenant, error) {
	logger.L().Info("create tenant with domains DB",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("tenant", tenant),
		zap.Strings("domains", domains),
	)

	tx := r.db.WithContext(ctx).Begin()

	// 创建 Tenant
	if err := tx.Create(tenant).Error; err != nil {
		logger.L().Error("failed to create tenant in transaction",
			zap.String("trace_id", trace.FromContext(ctx)),
			zap.Error(err),
		)
		tx.Rollback()
		return nil, err
	}

	// 创建 Domains（如果有的话）
	if len(domains) > 0 {
		for _, domainName := range domains {
			domain := &entity.Domain{
				TenantID: tenant.ID,
				Domain:   domainName,
				Status:   1,
			}
			if err := tx.Create(domain).Error; err != nil {
				logger.L().Error("failed to create domain in transaction",
					zap.String("trace_id", trace.FromContext(ctx)),
					zap.String("domain", domainName),
					zap.Error(err),
				)
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.L().Error("failed to commit transaction",
			zap.String("trace_id", trace.FromContext(ctx)),
			zap.Error(err),
		)
		return nil, err
	}

	logger.L().Info("tenant with domains created successfully",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Uint("tenant_id", tenant.ID),
		zap.Int("domains_count", len(domains)),
	)

	return tenant, nil
}

func (r *TenantRepo) List(ctx context.Context, offset, limit int) ([]entity.Tenant, int64, error) {
	var tenants []entity.Tenant
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Tenant{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&tenants).Error
	return tenants, total, err
}

func (r *TenantRepo) ListTenants(ctx context.Context, id *uint, domain string, startTime, endTime *time.Time, offset, limit int) ([]entity.Tenant, int64, error) {
	logger.L().Info("list tenants from repo",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("id", id),
		zap.String("domain", domain),
		zap.Any("offset", offset),
		zap.Any("limit", limit),
		zap.Any("startTime", startTime),
		zap.Any("endTime", endTime),
	)

	var tenants []entity.Tenant
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Tenant{})

	// 按ID查询
	if id != nil {
		db = db.Where("id = ?", *id)
	}

	// 按Domain查询 - 需要 LEFT JOIN domains表
	if domain != "" {
		db = db.Joins("LEFT JOIN domains ON domains.tenant_id = tenants.id").
			Where("domains.domain = ? AND domains.status = 1", domain).
			Distinct("tenants.id")
	}

	// 按时间范围查询
	if startTime != nil {
		db = db.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		db = db.Where("created_at <= ?", endTime)
	}

	logger.L().Warn("list tenants from repo after count",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.String("domain", domain),
		zap.Bool("has_id_filter", id != nil),

		zap.Bool("has_time_range", startTime != nil || endTime != nil),
		zap.Int64("count", total),
	)

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := db.Offset(offset).Limit(limit).Find(&tenants).Error
	return tenants, total, err
}
