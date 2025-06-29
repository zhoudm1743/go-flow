# 多租户权限管理系统

基于Casbin的多租户权限管理系统，支持不同租户使用不同数据库，提供完整的RBAC权限控制。

## 系统架构

### 1. 核心模型

#### 租户模型 (SystemTenant)
```go
type SystemTenant struct {
    BaseModel
    Name             string // 租户名称
    Code             string // 租户编码
    Type             string // 租户类型
    Status           int8   // 状态:1正常,0禁用
    ExpireAt         int64  // 过期时间
    DatabaseHost     string // 数据库主机
    DatabasePort     string // 数据库端口
    DatabaseName     string // 数据库名称
    DatabaseUser     string // 数据库用户
    DatabasePassword string // 数据库密码(加密)
}
```

#### 权限相关模型
- **SystemTenantRole**: 租户角色表
- **SystemTenantPermission**: 租户权限表
- **SystemTenantRolePermission**: 角色权限关联表
- **SystemTenantUserRole**: 用户角色关联表
- **SystemCasbinRule**: Casbin策略规则表

### 2. 服务层架构

#### TenantDatabaseService
负责管理不同租户的数据库连接：
- 连接池管理
- 数据库连接测试
- 租户表结构初始化
- 动态数据库切换

#### TenantService
租户管理服务：
- 租户CRUD操作
- 数据库连接管理
- 状态控制

#### TenantPermissionService
权限管理服务：
- 用户权限验证
- 角色权限分配
- 权限查询

## API接口

### 租户管理接口

#### 1. 获取租户列表
```http
GET /system/tenant/list?page_no=1&page_size=10&name=租户名
Authorization: Bearer {token}
```

#### 2. 添加租户
```http
POST /system/tenant/add
Authorization: Bearer {token}
Content-Type: application/json

{
    "name": "测试租户",
    "code": "test_tenant",
    "type": "enterprise",
    "status": 1,
    "expire_at": 1735689600,
    "database_host": "localhost",
    "database_port": "3306",
    "database_name": "tenant_test",
    "database_user": "root",
    "database_password": "password123"
}
```

#### 3. 编辑租户
```http
POST /system/tenant/edit
Authorization: Bearer {token}
Content-Type: application/json

{
    "id": 1,
    "name": "更新的租户名",
    "code": "updated_tenant",
    // 其他字段...
}
```

#### 4. 删除租户
```http
POST /system/tenant/delete
Authorization: Bearer {token}
Content-Type: application/json

{
    "id": 1
}
```

#### 5. 修改租户状态
```http
POST /system/tenant/change-status
Authorization: Bearer {token}
Content-Type: application/json

{
    "id": 1
}
```

#### 6. 测试数据库连接
```http
POST /system/tenant/test-connection
Authorization: Bearer {token}
Content-Type: application/json

{
    "id": 1
}
```

### 权限管理接口

#### 1. 检查用户权限
```http
POST /system/permission/check
Authorization: Bearer {token}
Content-Type: application/json

{
    "tenant_id": 1,
    "user_id": 1,
    "resource": "/api/users",
    "action": "GET"
}
```

#### 2. 为用户分配角色
```http
POST /system/user/assign-roles
Authorization: Bearer {token}
Content-Type: application/json

{
    "tenant_id": 1,
    "user_id": 1,
    "role_ids": [1, 2, 3]
}
```

#### 3. 获取用户角色
```http
GET /system/user/roles?tenant_id=1&user_id=1
Authorization: Bearer {token}
```

#### 4. 获取用户权限
```http
GET /system/user/permissions?tenant_id=1&user_id=1
Authorization: Bearer {token}
```

#### 5. 为角色分配权限
```http
POST /system/role/assign-permissions
Authorization: Bearer {token}
Content-Type: application/json

{
    "tenant_id": 1,
    "role_id": 1,
    "permission_ids": [1, 2, 3, 4]
}
```

#### 6. 获取角色权限
```http
GET /system/role/permissions?tenant_id=1&role_id=1
Authorization: Bearer {token}
```

## 权限验证中间件

### 1. 全局权限中间件
```go
// 自动验证用户对当前请求路径的权限
app.Use(middleware.TenantPermissionMiddleware(jwtService, permissionService))
```

### 2. 特定权限中间件
```go
// 验证特定资源和操作的权限
app.POST("/api/sensitive", 
    middleware.RequireTenantPermission("/api/sensitive", "POST", jwtService, permissionService),
    handler,
)
```

### 3. 可选权限中间件
```go
// 如果提供token则验证权限，否则跳过
app.Use(middleware.OptionalTenantPermissionMiddleware(jwtService, permissionService))
```

## 使用示例

### 1. 创建租户
```go
// 添加租户请求
addReq := &req.SystemTenantAddReq{
    Name:             "新租户",
    Code:             "new_tenant",
    Type:             "enterprise",
    Status:           1,
    DatabaseHost:     "localhost",
    DatabasePort:     "3306",
    DatabaseName:     "tenant_new",
    DatabaseUser:     "root",
    DatabasePassword: "password123",
}

err := tenantService.Add(addReq)
if err != nil {
    // 处理错误
}
```

### 2. 权限验证
```go
// 检查用户是否有权限访问某资源
hasPermission, err := permissionService.CheckPermission(
    tenantID, 
    userID, 
    "/api/users", 
    "GET",
)
if err != nil {
    // 处理错误
}
if !hasPermission {
    // 没有权限
}
```

### 3. 分配角色权限
```go
// 为角色分配权限
assignReq := &req.SystemTenantRolePermissionReq{
    TenantID:      1,
    RoleID:        1,
    PermissionIDs: []uint{1, 2, 3},
}

err := permissionService.AssignPermissionsToRole(assignReq)
if err != nil {
    // 处理错误
}
```

### 4. 在中间件中获取租户信息
```go
func handler(c *gin.Context) {
    // 从上下文获取租户ID
    tenantID, exists := middleware.GetTenantID(c)
    if !exists {
        // 处理错误
        return
    }
    
    // 从上下文获取用户ID
    userID, exists := middleware.GetUserID(c)
    if !exists {
        // 处理错误
        return
    }
    
    // 业务逻辑...
}
```

## 数据库设计

### 主数据库（存储租户信息）
- `system_tenants`: 租户表
- `system_admins`: 管理员表（包含tenant_id字段）

### 租户数据库（每个租户独立）
- `system_tenant_roles`: 角色表
- `system_tenant_permissions`: 权限表
- `system_tenant_role_permissions`: 角色权限关联表
- `system_tenant_user_roles`: 用户角色关联表
- `casbin_rule`: Casbin规则表

## 权限控制策略

### RBAC模型
```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
```

### 权限层级
1. **租户级别**: 不同租户数据完全隔离
2. **角色级别**: 角色拥有一组权限
3. **用户级别**: 用户通过角色获得权限
4. **资源级别**: 对具体API路径和操作的控制

## 安全特性

1. **数据库隔离**: 每个租户使用独立数据库
2. **连接池管理**: 动态管理数据库连接，防止连接泄露
3. **权限验证**: 多层次权限验证，确保数据安全
4. **密码加密**: 数据库密码使用MD5加密存储
5. **Token验证**: JWT token验证用户身份

## 部署注意事项

1. **数据库配置**: 确保每个租户有独立的数据库
2. **连接池**: 根据租户数量调整连接池配置
3. **缓存策略**: 考虑对权限查询结果进行缓存
4. **监控**: 监控数据库连接数和权限验证性能
5. **备份**: 定期备份租户数据库

## 扩展功能

### 1. 添加Casbin支持
如需要更强大的权限控制，可以安装Casbin依赖：
```bash
go get github.com/casbin/casbin/v2
go get github.com/casbin/gorm-adapter/v3
```

### 2. 权限缓存
可以添加Redis缓存来提高权限查询性能：
```go
// 在权限检查前先查缓存
func (t *tenantPermissionService) CheckPermissionWithCache(tenantID, userID uint, resource, action string) (bool, error) {
    // 先查缓存
    cacheKey := fmt.Sprintf("permission:%d:%d:%s:%s", tenantID, userID, resource, action)
    // ... Redis查询逻辑
    
    // 缓存未命中时查数据库
    return t.CheckPermission(tenantID, userID, resource, action)
}
```

### 3. 权限日志
添加权限访问日志：
```go
func logPermissionCheck(tenantID, userID uint, resource, action string, result bool) {
    log.Printf("Permission check: tenant=%d, user=%d, resource=%s, action=%s, result=%t", 
        tenantID, userID, resource, action, result)
}
```

这个多租户权限管理系统提供了完整的租户隔离和细粒度的权限控制，可以满足企业级应用的安全需求。 