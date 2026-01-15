package service

import "tenant-service/internal/repo"

type TenantService struct {
	repo *repo.TenantRepo
}

func NewCommonService(repo *repo.TenantRepo) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) GetCommon(id string) string {
	return s.repo.GetTenant(id)
}
