package utils

import "math"

type PaginationParams struct {
	Page     int
	PageSize int
}

// CalculatePaginationParams 计算分页参数
func CalculatePaginationParams(page, pageSize int) (offset int, limit int) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}

	offset = (page - 1) * pageSize
	return offset, pageSize
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}
