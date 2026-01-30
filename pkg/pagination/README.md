# Pagination Package

通用分页工具包，用于微服务中的分页查询和响应处理。

## 功能

### 数据结构

#### `PaginationRequest`
标准的分页请求参数：
- `Page`: 页码（从 1 开始，必填，最小值 1）
- `PageSize`: 每页数量（必填，范围 1-100）

#### `PaginationResponse[T]`
通用的分页响应模板（泛型）：
- `Items`: 数据项列表
- `Total`: 总数量
- `Page`: 当前页码
- `PageSize`: 每页数量
- `TotalPages`: 总页数

### 工具函数

#### `CalculatePaginationParams(page, pageSize int) (offset, limit int)`
计算数据库查询参数
- 参数验证和默认值处理
- 返回 offset 和 limit

#### `CalculateTotalPages(total int64, pageSize int) int`
计算总页数
- 使用数学方法向上取整

## 使用示例

### 1. 定义业务特定的请求/响应类型

```go
package dto

import "shared.local/pkg/pagination"

// 继承通用分页请求
type ListUsersRequest struct {
    pagination.PaginationRequest
    Name   string `form:"name"`
    Status int    `form:"status"`
}

// 定义响应类型
type ListUsersResponse struct {
    Items      []UserResponse `json:"items"`
    Total      int64          `json:"total"`
    Page       int            `json:"page"`
    PageSize   int            `json:"pageSize"`
    TotalPages int            `json:"totalPages"`
}
```

### 2. 在 Service 层使用

```go
package service

import "shared.local/pkg/pagination"

func (s *UserService) ListUsers(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
    // 计算分页参数
    offset, limit := pagination.CalculatePaginationParams(req.Page, req.PageSize)
    
    // 查询数据库
    users, total, err := s.repo.List(ctx, offset, limit, req.Name, req.Status)
    if err != nil {
        return nil, err
    }
    
    // 构建响应
    totalPages := pagination.CalculateTotalPages(total, req.PageSize)
    return &dto.ListUsersResponse{
        Items:      users,
        Total:      total,
        Page:       req.Page,
        PageSize:   req.PageSize,
        TotalPages: totalPages,
    }, nil
}
```

## 跨服务使用

由于此包位于 `pkg` 目录，所有微服务都可以使用：

```go
import "shared.local/pkg/pagination"
```

## 参数验证规则

分页参数会在以下位置验证：

1. **HTTP 请求绑定**：使用 Gin 的 binding 标签
   - `form:"page" binding:"required,min=1"`
   - `form:"pageSize" binding:"required,min=1,max=100"`

2. **函数内部**：`CalculatePaginationParams` 会处理异常情况
   - PageSize <= 0：默认为 10
   - PageSize > 100：限制为 100
   - Page <= 0：默认为 1

## 注意事项

- `PaginationRequest` 应该只在 API 请求中使用
- 数据库查询总数应使用 COUNT 语句获取准确值
- offset 和 limit 仅用于 SQL LIMIT 和 OFFSET 子句
