package repo

type TenantRepo struct{}

func NewTenantRepo() *TenantRepo {
	return &TenantRepo{}
}

func (r *TenantRepo) GetTenant(id string) string {
	return "user-" + id
}
