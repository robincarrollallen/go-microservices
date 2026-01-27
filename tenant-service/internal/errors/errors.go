package apperror

import "errors"

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrDomainExists   = errors.New("domain already exists")
	ErrNameExists     = errors.New("tenant name already exists")
)
