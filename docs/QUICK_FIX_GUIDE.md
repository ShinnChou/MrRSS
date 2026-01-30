# 前5阶段关键问题修复指南

**快速索引**:
- [问题1: Handler 仍直接依赖](#问题1-handler-仍直接依赖)
- [问题2: article_db.go 过大](#问题2-articledbo-过大)
- [问题3: Response 格式未统一](#问题3-response-格式未统一)
- [修复优先级](#修复优先级表)

---

## 问题1: Handler 仍直接依赖

**状态**: 🔴 需立即修复
**影响**: 20+ handler 文件
**难度**: ⭐⭐ (中等)
**工作量**: 1-2 天

### 现状分析

Handler 结构仍保留原有直接依赖:

```go
// handler.go 第 50-65 行
type Handler struct {
    Services *svc.Registry    // ✅ 新增

    // ❌ 这些仍需移除或便捷化
    DB               *database.DB
    Fetcher          *feed.Fetcher
    Translator       translation.Translator
    AITracker        *aiusage.Tracker
    DiscoveryService *discovery.Service
    // ...
}
```

Handler 方法中的使用模式:

```go
// ❌ 当前: 直接使用 h.DB
func HandleGetArticles(h *core.Handler, w http.ResponseWriter, r *http.Request) {
    articles, err := h.DB.GetArticles(filter, feedID, category, showHidden, limit, offset)
    // ...
}

// ✅ 应该改为: 通过服务层
func HandleGetArticles(h *core.Handler, w http.ResponseWriter, r *http.Request) {
    articles, err := h.Services.Article().GetArticles(ctx, ArticleQueryOptions{...})
    // ...
}
```

### 修复方案

#### 方案 A: 逐步迁移 (推荐)

1. **保留 Handler 直接依赖** (向后兼容)
   ```go
   type Handler struct {
       Services *svc.Registry

       // 为向后兼容，保留但标记为 Deprecated
       DB               *database.DB        `deprecated:"use Services.DB()"`
       Fetcher          *feed.Fetcher       `deprecated:"use Services.Feed()"`
       Translator       translation.Translator `deprecated:"use Services.Translation()"`
       // ...
   }
   ```

2. **添加便捷访问方法**
   ```go
   // handler.go
   func (h *Handler) GetDB() *database.DB {
       if h.DB != nil {
           return h.DB
       }
       // 从 service 中获取 (如果已迁移到 registry)
       return h.Services.GetDB()
   }
   ```

3. **逐文件迁移** (按优先级)
   ```
   优先级高 (核心功能):
   - handlers/article/*.go
   - handlers/feed/*.go
   - handlers/ai/*.go

   优先级中 (配置):
   - handlers/settings/*.go
   - handlers/translation/*.go

   优先级低 (辅助):
   - handlers/window/*.go
   - handlers/browser/*.go
   ```

#### 方案 B: 一次性重构 (快速但风险高)

```go
// 直接移除所有直接依赖
type Handler struct {
    Services *svc.Registry
    App      interface{}  // 保留 Wails app 实例
}
```

**风险**: 需要同时修改所有 handler，容易出错

### 修复检查清单

- [ ] 创建迁移计划文档
- [ ] 识别所有受影响的 handler 文件 (20+)
- [ ] 为每个文件创建对应的服务调用示例
- [ ] 逐文件编辑并测试
- [ ] 添加单元测试
- [ ] 进行集成测试

### 快速命令

```bash
# 找出所有直接使用 h.DB 的地方
grep -r "h\.DB\." internal/handlers/ --include="*.go" | head -20

# 统计受影响的文件数
grep -r "h\.DB\." internal/handlers/ --include="*.go" | cut -d: -f1 | sort -u | wc -l

# 统计受影响的行数
grep -r "h\.DB\." internal/handlers/ --include="*.go" | wc -l
```

---

## 问题2: article_db.go 过大

**状态**: 🔴 需立即修复
**当前行数**: 924 行
**目标**: 拆分为 4 个文件
**难度**: ⭐⭐ (中等)
**工作量**: 1 天

### 拆分方案

```plaintext
article_db.go (924行)
├── article_crud_db.go (~200行)
│   ├── Save(ctx, article)
│   ├── SaveBatch(ctx, articles)
│   ├── GetByID(id)
│   ├── Delete(id)
│   └── // 其他基本 CRUD
│
├── article_query_db.go (~300行)
│   ├── GetArticles(filter, feedID, ...)
│   ├── Query(opts)
│   ├── Search(keyword)
│   ├── Count(opts)
│   └── // 复杂查询逻辑
│
├── article_batch_db.go (~200行)
│   ├── SaveBatch()
│   ├── DeleteBatch()
│   ├── UpdateBatch()
│   └── // 批量操作
│
└── article_status_db.go (~100+行)
    ├── MarkRead(id, read)
    ├── MarkFavorite(id, favorite)
    ├── MarkHidden(id, hidden)
    └── // 状态更新方法
```

### 拆分步骤

1. **分析 article_db.go 中的函数分类**
   ```bash
   grep -n "^func.*DB.*Article" internal/database/article_db.go | head -30
   ```

2. **创建新文件，移动相应函数**
   - 每个文件保留相同的 receiver
   - 保持相同的包名 `database`

3. **验证编译**
   ```bash
   go build ./...
   go test ./internal/database/...
   ```

4. **检查是否有遗漏**
   ```bash
   # 确保所有 ArticleDB 方法都被找到
   grep "article_crud_db.go\|article_query_db.go\|article_batch_db.go\|article_status_db.go" go.mod
   ```

### 快速命令

```bash
# 获取 article_db.go 中所有公开函数
grep -n "^func.*DB.*\(.*)" internal/database/article_db.go

# 统计各类函数
echo "CRUD functions:"
grep "^func.*DB.*\(Save\|Get\|Delete\)" internal/database/article_db.go | wc -l

echo "Query functions:"
grep "^func.*DB.*\(Search\|Query\|Count\|Filter\)" internal/database/article_db.go | wc -l

echo "Status functions:"
grep "^func.*DB.*Mark\|Update" internal/database/article_db.go | wc -l
```

---

## 问题3: Response 格式未统一

**状态**: ⚠️ 部分完成
**受影响**: 20+ handler
**难度**: ⭐ (简单)
**工作量**: 1-2 天

### 现状

当前响应格式混乱:

```go
// 方式1: 直接 http.Error
http.Error(w, "Invalid feed ID", http.StatusBadRequest)

// 方式2: JSON 响应
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{
    "success": true,
    "data": articles,
})

// 方式3: 无格式化响应
w.WriteHeader(http.StatusInternalServerError)
fmt.Fprintf(w, "error: %v", err)
```

### 目标格式

```json
// 成功响应
{
    "success": true,
    "data": { /* ... */ }
}

// 错误响应
{
    "success": false,
    "error": {
        "code": "INVALID_INPUT",
        "message": "Invalid feed ID",
        "detail": "Feed ID must be a positive integer"
    }
}
```

### 实现步骤

1. **创建响应工具** (已有框架，需完善)

```go
// internal/handlers/response/response.go (待完善)
package response

import (
    "encoding/json"
    "net/http"
    apperrors "MrRSS/internal/errors"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// Success 返回成功响应
func Success(w http.ResponseWriter, data interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(APIResponse{
        Success: true,
        Data:    data,
    })
}

// Error 返回错误响应
func Error(w http.ResponseWriter, err error) {
    w.Header().Set("Content-Type", "application/json")

    var appErr *apperrors.AppError
    if e, ok := err.(*apperrors.AppError); ok {
        appErr = e
    } else {
        appErr = apperrors.NewAppError(
            apperrors.ErrCodeInternal,
            "An error occurred",
            err,
        )
    }

    statusCode := errorToHTTPStatus(appErr.Code)
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(APIResponse{
        Success: false,
        Error: &ErrorInfo{
            Code:    string(appErr.Code),
            Message: appErr.Message,
        },
    })
}

func errorToHTTPStatus(code apperrors.ErrorCode) int {
    switch code {
    case apperrors.ErrCodeNotFound:
        return http.StatusNotFound
    case apperrors.ErrCodeInvalidInput:
        return http.StatusBadRequest
    case apperrors.ErrCodeUnauthorized:
        return http.StatusUnauthorized
    case apperrors.ErrCodeForbidden:
        return http.StatusForbidden
    default:
        return http.StatusInternalServerError
    }
}
```

2. **更新 Handler 使用方式**

```go
// ❌ 旧方式
http.Error(w, "Invalid input", http.StatusBadRequest)

// ✅ 新方式
response.Error(w, apperrors.NewAppError(
    apperrors.ErrCodeInvalidInput,
    "Invalid feed ID",
    fmt.Errorf("feed id: %w", err),
))

// 成功响应
response.Success(w, articles, http.StatusOK)
```

3. **按优先级迁移 Handler**
   - 核心 API 优先 (article, feed, ai)
   - 配置 API (settings, translation)
   - 其他 API (window, browser)

### 修复检查清单

- [ ] 完善 `response.go` 工具函数
- [ ] 逐个 handler 文件迁移
- [ ] 测试所有 API 端点返回格式
- [ ] 更新 API 文档/Swagger

### 快速验证

```bash
# 测试 API 响应格式
curl -X GET http://localhost:1234/api/articles 2>/dev/null | jq .

# 应输出:
# {
#   "success": true,
#   "data": [...]
# }
```

---

## 修复优先级表

| 优先级 | 问题 | 工作量 | 依赖 | 建议执行顺序 |
|--------|------|--------|------|-----------|
| 🔴 高 | Handler 迁移 | 1-2d | 无 | 1 |
| 🔴 高 | Response 格式 | 1-2d | 1 | 2 |
| 🔴 高 | article_db 拆分 | 1d | 无 | 3 |
| 🟡 中 | Repository 接口 | 1d | 3 | 4 |
| 🟡 中 | 服务完善 | 1d | 1, 4 | 5 |
| 🟡 中 | Provider 验证 | 0.5d | 无 | 6 |

**总工作量**: 约 5-7 天

---

## 快速修复命令集

```bash
# 1. 识别问题区域
echo "=== Handler 直接依赖使用 ==="
grep -r "h\.DB\." internal/handlers/ --include="*.go" | wc -l

echo "=== Response 格式混乱 ==="
grep -r "http\.Error\|WriteHeader.*Status\|json\.NewEncoder" internal/handlers/ --include="*.go" | wc -l

echo "=== 大文件列表 ==="
find internal/database -name "*.go" -exec wc -l {} + | sort -rn | head -10

# 2. 测试编译
go build ./...

# 3. 运行现有测试
go test ./...
```

---

## 下一步行动

**立即 (今天)**:
1. 确认上述问题分析准确无误
2. 制定详细的修复计划
3. 按优先级创建任务

**短期 (3-5 天)**:
1. 完成高优先级问题修复
2. 进行集成测试
3. 验证 API 功能

**中期 (1 周)**:
1. 完成所有修复
2. 补充单元测试
3. 准备合并回 main 分支
