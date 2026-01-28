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
