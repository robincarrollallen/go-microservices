package repo

import "context"

type TenantRepo struct{}

func NewTenantRepo() *TenantRepo {
	return &TenantRepo{}
}

func (r *TenantRepo) GetTenant(ctx context.Context, id string) string {
	return "tenant-" + id
}
