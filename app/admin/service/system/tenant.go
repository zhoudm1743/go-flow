package system

import (
	"errors"

	"github.com/zhoudm1743/go-flow/app/admin/schemas/req"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/resp"
	"github.com/zhoudm1743/go-flow/app/models"
	"github.com/zhoudm1743/go-flow/core/database"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/util"
	"gorm.io/gorm"
)

type TenantService interface {
	All(queryReq *req.SystemTenantListReq) ([]resp.SystemTenantResp, error)
	List(pageReq *req.PageReq, queryReq *req.SystemTenantListReq) (response.PageResp, error)
	Find(id uint) (resp.SystemTenantResp, error)
	Add(req *req.SystemTenantAddReq) error
	Edit(req *req.SystemTenantEditReq) error
	Delete(id uint) error
	ChangeStatus(id uint) error
	TestConnection(id uint) error
}

type tenantService struct {
	db              database.Database
	tenantDBService TenantDatabaseService
}

// All 获取所有租户
func (t *tenantService) All(queryReq *req.SystemTenantListReq) ([]resp.SystemTenantResp, error) {
	chain := t.db.GetDB().Model(&models.SystemTenant{})
	if queryReq.Name != "" {
		chain = chain.Where("name LIKE ?", "%"+queryReq.Name+"%")
	}
	if queryReq.Code != "" {
		chain = chain.Where("code LIKE ?", "%"+queryReq.Code+"%")
	}
	if queryReq.Type != "" {
		chain = chain.Where("type = ?", queryReq.Type)
	}
	if queryReq.Status != -1 {
		chain = chain.Where("status = ?", queryReq.Status)
	}
	chain.Order("id DESC")

	var tenants []models.SystemTenant
	if err := chain.Find(&tenants).Error; err != nil {
		return nil, err
	}

	var respTenants []resp.SystemTenantResp
	response.Copy(&respTenants, tenants)
	return respTenants, nil
}

// List 获取租户列表（分页）
func (t *tenantService) List(pageReq *req.PageReq, queryReq *req.SystemTenantListReq) (response.PageResp, error) {
	chain := t.db.GetDB().Model(&models.SystemTenant{})
	if queryReq.Name != "" {
		chain = chain.Where("name LIKE ?", "%"+queryReq.Name+"%")
	}
	if queryReq.Code != "" {
		chain = chain.Where("code LIKE ?", "%"+queryReq.Code+"%")
	}
	if queryReq.Type != "" {
		chain = chain.Where("type = ?", queryReq.Type)
	}
	if queryReq.Status != -1 {
		chain = chain.Where("status = ?", queryReq.Status)
	}
	chain.Order("id DESC")

	var count int64
	chain.Count(&count)

	offset := (pageReq.PageNo - 1) * pageReq.PageSize
	chain = chain.Offset(offset).Limit(pageReq.PageSize)

	var tenants []models.SystemTenant
	if err := chain.Find(&tenants).Error; err != nil {
		return response.PageResp{}, err
	}

	var respTenants []resp.SystemTenantResp
	response.Copy(&respTenants, tenants)

	return response.PageResp{
		PageNo:   pageReq.PageNo,
		PageSize: pageReq.PageSize,
		Count:    count,
		Lists:    respTenants,
	}, nil
}

// Find 根据ID获取租户详情
func (t *tenantService) Find(id uint) (resp.SystemTenantResp, error) {
	var tenant models.SystemTenant
	if err := t.db.GetDB().Where("id = ?", id).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.SystemTenantResp{}, errors.New("租户不存在")
		}
		return resp.SystemTenantResp{}, err
	}

	var respTenant resp.SystemTenantResp
	response.Copy(&respTenant, tenant)
	return respTenant, nil
}

// Add 添加租户
func (t *tenantService) Add(addReq *req.SystemTenantAddReq) error {
	// 检查租户编码是否已存在
	var existingTenant models.SystemTenant
	if err := t.db.GetDB().Where("code = ?", addReq.Code).First(&existingTenant).Error; err == nil {
		return errors.New("租户编码已存在")
	}

	// 创建租户对象
	tenant := models.SystemTenant{
		Name:             addReq.Name,
		Code:             addReq.Code,
		Type:             addReq.Type,
		Status:           addReq.Status,
		ExpireAt:         addReq.ExpireAt,
		DatabaseHost:     addReq.DatabaseHost,
		DatabasePort:     addReq.DatabasePort,
		DatabaseName:     addReq.DatabaseName,
		DatabaseUser:     addReq.DatabaseUser,
		DatabasePassword: util.ToolsUtil.MakeMd5(addReq.DatabasePassword), // 加密存储
	}

	// 测试数据库连接
	if err := t.tenantDBService.TestTenantDBConnection(&tenant); err != nil {
		return err
	}

	// 开始事务
	tx := t.db.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存租户
	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建租户数据库连接
	if err := t.tenantDBService.CreateTenantDB(&tenant); err != nil {
		tx.Rollback()
		return err
	}

	// 初始化租户表结构
	if err := t.tenantDBService.InitTenantTables(tenant.ID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Edit 编辑租户
func (t *tenantService) Edit(editReq *req.SystemTenantEditReq) error {
	var tenant models.SystemTenant
	if err := t.db.GetDB().Where("id = ?", editReq.ID).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("租户不存在")
		}
		return err
	}

	// 检查租户编码是否被其他租户使用
	if editReq.Code != tenant.Code {
		var existingTenant models.SystemTenant
		if err := t.db.GetDB().Where("code = ? AND id != ?", editReq.Code, editReq.ID).First(&existingTenant).Error; err == nil {
			return errors.New("租户编码已存在")
		}
	}

	// 更新字段
	tenant.Name = editReq.Name
	tenant.Code = editReq.Code
	tenant.Type = editReq.Type
	tenant.Status = editReq.Status
	tenant.ExpireAt = editReq.ExpireAt
	tenant.DatabaseHost = editReq.DatabaseHost
	tenant.DatabasePort = editReq.DatabasePort
	tenant.DatabaseName = editReq.DatabaseName
	tenant.DatabaseUser = editReq.DatabaseUser

	// 如果密码有变更
	if editReq.DatabasePassword != "" {
		tenant.DatabasePassword = util.ToolsUtil.MakeMd5(editReq.DatabasePassword)
	}

	// 如果数据库连接信息有变更，重新测试连接
	if editReq.DatabaseHost != "" && editReq.DatabasePort != "" &&
		editReq.DatabaseName != "" && editReq.DatabaseUser != "" {
		if err := t.tenantDBService.TestTenantDBConnection(&tenant); err != nil {
			return err
		}
	}

	// 更新数据库
	if err := t.db.GetDB().Save(&tenant).Error; err != nil {
		return err
	}

	// 如果状态变为禁用，移除数据库连接
	if editReq.Status == 0 {
		t.tenantDBService.RemoveTenantDB(tenant.ID)
	}

	return nil
}

// Delete 删除租户
func (t *tenantService) Delete(id uint) error {
	var tenant models.SystemTenant
	if err := t.db.GetDB().Where("id = ?", id).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("租户不存在")
		}
		return err
	}

	// 检查是否还有关联的用户
	var adminCount int64
	if err := t.db.GetDB().Model(&models.SystemAdmin{}).Where("tenant_id = ?", id).Count(&adminCount).Error; err != nil {
		return err
	}

	if adminCount > 0 {
		return errors.New("该租户下还有用户，无法删除")
	}

	// 移除数据库连接
	t.tenantDBService.RemoveTenantDB(id)

	// 删除租户
	if err := t.db.GetDB().Delete(&tenant).Error; err != nil {
		return err
	}

	return nil
}

// ChangeStatus 修改租户状态
func (t *tenantService) ChangeStatus(id uint) error {
	var tenant models.SystemTenant
	if err := t.db.GetDB().Where("id = ?", id).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("租户不存在")
		}
		return err
	}

	// 切换状态
	tenant.Status = 1 - tenant.Status

	if err := t.db.GetDB().Save(&tenant).Error; err != nil {
		return err
	}

	// 如果状态变为禁用，移除数据库连接
	if tenant.Status == 0 {
		t.tenantDBService.RemoveTenantDB(tenant.ID)
	}

	return nil
}

// TestConnection 测试租户数据库连接
func (t *tenantService) TestConnection(id uint) error {
	var tenant models.SystemTenant
	if err := t.db.GetDB().Where("id = ?", id).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("租户不存在")
		}
		return err
	}

	return t.tenantDBService.TestTenantDBConnection(&tenant)
}

// NewTenantService 创建租户服务
func NewTenantService(db database.Database, tenantDBService TenantDatabaseService) TenantService {
	return &tenantService{
		db:              db,
		tenantDBService: tenantDBService,
	}
}
