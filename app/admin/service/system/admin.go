package system

import (
	"errors"

	"github.com/zhoudm1743/go-flow/app/admin/schemas/req"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/resp"
	"github.com/zhoudm1743/go-flow/app/models"
	"github.com/zhoudm1743/go-flow/core/database"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"gorm.io/gorm"
)

type AdminService interface {
	All(queryReq *req.SystemAdminListReq) ([]resp.SystemAdminResp, error)
	List(pageReq *req.PageReq, queryReq *req.SystemAdminListReq) (response.PageResp, error)
	Find(id uint) (resp.SystemAdminResp, error)
	FindByUsername(username string) (resp.SystemAdminResp, error)
	Add(req *req.SystemAdminAddReq) error
	Edit(req *req.SystemAdminEditReq) error
	Delete(id uint) error
	ChangeStatus(id uint) error
}

type adminService struct {
	db database.Database
}

// Add implements AdminService.
func (a *adminService) Add(req *req.SystemAdminAddReq) error {
	panic("unimplemented")
}

// All implements AdminService.
func (a *adminService) All(queryReq *req.SystemAdminListReq) ([]resp.SystemAdminResp, error) {
	chain := a.db.GetDB().Model(&models.SystemAdmin{})
	if queryReq.Username != "" {
		chain = chain.Where("username LIKE ?", "%"+queryReq.Username+"%")
	}
	if queryReq.Nickname != "" {
		chain = chain.Where("nickname LIKE ?", "%"+queryReq.Nickname+"%")
	}
	if queryReq.Status != -1 {
		chain = chain.Where("status = ?", queryReq.Status)
	}
	if queryReq.Role != "" {
		chain = chain.Where("role = ?", queryReq.Role)
	}
	chain.Order("id DESC")
	var admins []models.SystemAdmin
	if err := chain.Find(&admins).Error; err != nil {
		return nil, err
	}
	var respAdmins []resp.SystemAdminResp
	response.Copy(&respAdmins, admins)
	return respAdmins, nil
}

// ChangeStatus implements AdminService.
func (a *adminService) ChangeStatus(id uint) error {
	var admin models.SystemAdmin
	if err := a.db.GetDB().Model(&models.SystemAdmin{}).Where("id = ?", id).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("管理员不存在")
		}
		return err
	}
	admin.Status = 1 - admin.Status
	if err := a.db.GetDB().Save(&admin).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements AdminService.
func (a *adminService) Delete(id uint) error {
	var admin models.SystemAdmin
	if err := a.db.GetDB().Model(&models.SystemAdmin{}).Where("id = ?", id).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("管理员不存在")
		}
		return err
	}
	if err := a.db.GetDB().Delete(&admin).Error; err != nil {
		return err
	}
	return nil
}

// Edit implements AdminService.
func (a *adminService) Edit(req *req.SystemAdminEditReq) error {
	var admin models.SystemAdmin
	if err := a.db.GetDB().Model(&models.SystemAdmin{}).Where("id = ?", req.ID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("管理员不存在")
		}
		return err
	}
	response.Copy(&admin, req)
	if err := a.db.GetDB().Save(&admin).Error; err != nil {
		return err
	}
	return nil
}

// Find implements AdminService.
func (a *adminService) Find(id uint) (resp.SystemAdminResp, error) {
	var admin models.SystemAdmin
	if err := a.db.GetDB().Model(&models.SystemAdmin{}).Where("id = ?", id).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.SystemAdminResp{}, errors.New("管理员不存在")
		}
		return resp.SystemAdminResp{}, err
	}
	var respAdmin resp.SystemAdminResp
	response.Copy(&respAdmin, admin)
	return respAdmin, nil
}

// FindByUsername implements AdminService.
func (a *adminService) FindByUsername(username string) (resp.SystemAdminResp, error) {
	var admin models.SystemAdmin
	if err := a.db.GetDB().Model(&models.SystemAdmin{}).Where("username = ?", username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.SystemAdminResp{}, errors.New("管理员不存在")
		}
		return resp.SystemAdminResp{}, err
	}
	var respAdmin resp.SystemAdminResp
	response.Copy(&respAdmin, admin)
	return respAdmin, nil
}

// List implements AdminService.
func (a *adminService) List(pageReq *req.PageReq, queryReq *req.SystemAdminListReq) (response.PageResp, error) {
	chain := a.db.GetDB().Model(&models.SystemAdmin{})
	if queryReq.Username != "" {
		chain = chain.Where("username LIKE ?", "%"+queryReq.Username+"%")
	}
	if queryReq.Nickname != "" {
		chain = chain.Where("nickname LIKE ?", "%"+queryReq.Nickname+"%")
	}
	if queryReq.Status != -1 {
		chain = chain.Where("status = ?", queryReq.Status)
	}
	if queryReq.Role != "" {
		chain = chain.Where("role = ?", queryReq.Role)
	}
	chain.Order("id DESC")
	var count int64
	chain.Count(&count)
	var admins []models.SystemAdmin
	if err := chain.Find(&admins).Error; err != nil {
		return response.PageResp{}, err
	}
	var respAdmins []resp.SystemAdminResp
	response.Copy(&respAdmins, admins)
	return response.PageResp{
		PageNo:   pageReq.PageNo,
		PageSize: pageReq.PageSize,
		Count:    count,
		Lists:    respAdmins,
	}, nil
}

// NewAdminService 创建管理员服务
func NewAdminService(db database.Database) AdminService {
	return &adminService{
		db: db,
	}
}
