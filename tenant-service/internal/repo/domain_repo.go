package repo

import (
	"context"
	"tenant-service/internal/model/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"shared.local/pkg/logger"
	"shared.local/pkg/trace"
)

type DomainRepo struct {
	db *gorm.DB
}

func NewDomainRepo(db *gorm.DB) *DomainRepo {
	return &DomainRepo{db: db}
}

func (r *DomainRepo) CreateDomain(ctx context.Context, domain *entity.Domain) error {
	logger.L().Info("create domain DB",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("domain", domain),
	)

	return r.db.WithContext(ctx).Create(domain).Error
}

// List 分页查询域名列表
func (r *DomainRepo) List(ctx context.Context, offset, limit int, tenantID *uint, domain string) ([]entity.Domain, int64, error) {
	var domains []entity.Domain
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Domain{})

	// 构建条件查询
	if tenantID != nil {
		db = db.Where("tenant_id = ?", *tenantID)
	}
	if domain != "" {
		db = db.Where("domain LIKE ?", "%"+domain+"%")
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		logger.L().Error("count domains error",
			zap.String("trace_id", trace.FromContext(ctx)),
			zap.Error(err),
		)
		return nil, 0, err
	}

	// 查询分页数据
	err := db.Offset(offset).Limit(limit).Find(&domains).Error
	if err != nil {
		logger.L().Error("list domains error",
			zap.String("trace_id", trace.FromContext(ctx)),
			zap.Error(err),
		)
		return nil, 0, err
	}

	return domains, total, nil
}
