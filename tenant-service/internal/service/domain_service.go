package service

import "tenant-service/internal/repo"

type DomainService struct {
	repo *repo.DomainRepo
}

func NewDomainService(repo *repo.DomainRepo) *DomainService {
	return &DomainService{repo: repo}
}
