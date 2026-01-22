package repo

import (
	"context"
	"errors"

	"tenant-service/internal/model"

	"gorm.io/gorm"
)

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) Create(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *TenantRepo) GetTenantByID(ctx context.Context, id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.WithContext(ctx).First(&tenant, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tenant, err
}

func (r *TenantRepo) GetTenantByDomain(ctx context.Context, domain string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.WithContext(ctx).Where("domain = ?", domain).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tenant, err
}

func (r *TenantRepo) Update(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *TenantRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Tenant{}, id).Error
}

func (r *TenantRepo) List(ctx context.Context, offset, limit int) ([]model.Tenant, int64, error) {
	var tenants []model.Tenant
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Tenant{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&tenants).Error
	return tenants, total, err
}
