package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// genCmd 代码生成命令
var genCmd = &cobra.Command{
	Use:   "gen [module] [name] [comment]",
	Short: "生成代码",
	Long:  `生成模块代码，包括控制器、服务、仓库等`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		module := args[0]
		name := args[1]
		comment := args[2]

		err := GenerateCode(name, module, comment)
		if err != nil {
			fmt.Printf("生成代码失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("代码生成成功! 模块: %s, 名称: %s\n", module, name)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}

// 控制器模板
var controllerTmpl = `package controller

import (
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/schemas/req"
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/service"
	"github.com/zhoudm1743/go-frame/pkg/http/unified"
	"github.com/zhoudm1743/go-frame/pkg/response"
	"github.com/zhoudm1743/go-frame/util"
)

// {{.Name}}Controller {{.Comment}}控制器
type {{.Name}}Controller struct {
	service *service.{{.Name}}Service
}

// New{{.Name}}Controller 创建{{.Comment}}控制器
func New{{.Name}}Controller(service *service.{{.Name}}Service) *{{.Name}}Controller {
	return &{{.Name}}Controller{
		service: service,
	}
}

// List 获取列表
func (c *{{.Name}}Controller) List(ctx unified.Context) error {
	// 这里不需要验证参数，直接获取所有数据
	items, err := c.service.GetAll()
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, items)
}

// Get 获取单个记录
func (c *{{.Name}}Controller) Get(ctx unified.Context) error {
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	item, err := c.service.GetByID(idReq.ID)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, item)
}

// Create 创建记录
func (c *{{.Name}}Controller) Create(ctx unified.Context) error {
	var createReq req.{{.Name}}CreateReq
	if err := util.VerifyUtil.VerifyJSON(ctx, &createReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	created, err := c.service.Create(&createReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, created)
}

// Update 更新记录
func (c *{{.Name}}Controller) Update(ctx unified.Context) error {
	// 验证路径参数
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	// 验证请求体
	var updateReq req.{{.Name}}UpdateReq
	if err := util.VerifyUtil.VerifyJSON(ctx, &updateReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	updated, err := c.service.Update(idReq.ID, &updateReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, updated)
}

// Delete 删除记录
func (c *{{.Name}}Controller) Delete(ctx unified.Context) error {
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	err := c.service.Delete(idReq.ID)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOk(ctx)
}

// Page 分页查询
func (c *{{.Name}}Controller) Page(ctx unified.Context) error {
	var pageReq req.PageReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &pageReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	pageResult, err := c.service.GetPage(&pageReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, pageResult)
}
`

// 路由模板
var routerTmpl = `package controller

import (
	"github.com/zhoudm1743/go-frame/pkg/http/unified"
)

// {{.Name}}Router {{.Comment}}路由
type {{.Name}}Router struct {
	controller *{{.Name}}Controller
}

// New{{.Name}}Router 创建{{.Comment}}路由
func New{{.Name}}Router(controller *{{.Name}}Controller) *{{.Name}}Router {
	return &{{.Name}}Router{
		controller: controller,
	}
}

// RegisterRoutes 注册路由
func (r *{{.Name}}Router) RegisterRoutes(router unified.Router) {
	group := router.Group("/api/{{.LowerName}}ies")

	group.GET("", r.controller.List)
	group.GET("/:id", r.controller.Get)
	group.POST("", r.controller.Create)
	group.PUT("/:id", r.controller.Update)
	group.DELETE("/:id", r.controller.Delete)
	group.GET("/page", r.controller.Page)
}
`

// 服务模板
var serviceTmpl = `package service

import (
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/model"
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/repository"
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/schemas/req"
	"github.com/zhoudm1743/go-frame/pkg/response"
)

// {{.Name}}Service {{.Comment}}服务
type {{.Name}}Service struct {
	repo *repository.{{.Name}}Repository
}

// New{{.Name}}Service 创建{{.Comment}}服务
func New{{.Name}}Service(repo *repository.{{.Name}}Repository) *{{.Name}}Service {
	return &{{.Name}}Service{
		repo: repo,
	}
}

// GetAll 获取所有{{.Comment}}
func (s *{{.Name}}Service) GetAll() ([]*model.{{.Name}}, error) {
	return s.repo.FindAll()
}

// GetByID 根据ID获取{{.Comment}}
func (s *{{.Name}}Service) GetByID(id uint) (*model.{{.Name}}, error) {
	return s.repo.FindByID(id)
}

// Create 创建{{.Comment}}
func (s *{{.Name}}Service) Create(req *req.{{.Name}}CreateReq) (*model.{{.Name}}, error) {
	// 将请求转换为模型
	{{.LowerName}} := &model.{{.Name}}{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		// 根据需要添加更多字段
	}

	// 创建记录
	return s.repo.Create({{.LowerName}})
}

// Update 更新{{.Comment}}
func (s *{{.Name}}Service) Update(id uint, req *req.{{.Name}}UpdateReq) (*model.{{.Name}}, error) {
	// 查找记录
	{{.LowerName}}, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	{{.LowerName}}.Name = req.Name
	{{.LowerName}}.Description = req.Description
	{{.LowerName}}.Status = req.Status
	// 根据需要添加更多字段

	// 保存更新
	return s.repo.Update({{.LowerName}})
}

// Delete 删除{{.Comment}}
func (s *{{.Name}}Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetPage 分页查询{{.Comment}}
func (s *{{.Name}}Service) GetPage(req *req.PageReq) (*response.PageResult[*model.{{.Name}}], error) {
	return s.repo.FindPage(req)
}
`

// 仓库模板
var repositoryTmpl = `package repository

import (
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/model"
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/schemas/req"
	"github.com/zhoudm1743/go-frame/pkg/response"
	"gorm.io/gorm"
)

// {{.Name}}Repository {{.Comment}}仓库
type {{.Name}}Repository struct {
	db *gorm.DB
}

// New{{.Name}}Repository 创建{{.Comment}}仓库
func New{{.Name}}Repository(db *gorm.DB) *{{.Name}}Repository {
	// 自动迁移数据库模型
	_ = db.AutoMigrate(&model.{{.Name}}{})
	
	return &{{.Name}}Repository{
		db: db,
	}
}

// FindAll 查询所有{{.Comment}}
func (r *{{.Name}}Repository) FindAll() ([]*model.{{.Name}}, error) {
	var items []*model.{{.Name}}
	err := r.db.Find(&items).Error
	return items, err
}

// FindByID 根据ID查询{{.Comment}}
func (r *{{.Name}}Repository) FindByID(id uint) (*model.{{.Name}}, error) {
	var item model.{{.Name}}
	err := r.db.First(&item, id).Error
	return &item, err
}

// Create 创建{{.Comment}}
func (r *{{.Name}}Repository) Create(item *model.{{.Name}}) (*model.{{.Name}}, error) {
	err := r.db.Create(item).Error
	return item, err
}

// Update 更新{{.Comment}}
func (r *{{.Name}}Repository) Update(item *model.{{.Name}}) (*model.{{.Name}}, error) {
	err := r.db.Save(item).Error
	return item, err
}

// Delete 删除{{.Comment}}
func (r *{{.Name}}Repository) Delete(id uint) error {
	return r.db.Delete(&model.{{.Name}}{}, id).Error
}

// FindPage 分页查询{{.Comment}}
func (r *{{.Name}}Repository) FindPage(req *req.PageReq) (*response.PageResult[*model.{{.Name}}], error) {
	var items []*model.{{.Name}}
	var total int64

	// 查询总数
	query := r.db.Model(&model.{{.Name}}{})
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 查询数据
	err = query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return &response.PageResult[*model.{{.Name}}]{
		Total: total,
		Items: items,
	}, nil
}
`

// 模型模板
var modelTmpl = `package model

import (
	"github.com/zhoudm1743/go-frame/pkg/types"
)

// {{.Name}} {{.Comment}}模型
type {{.Name}} struct {
	types.GormModel
	Name        string ` + "`" + `json:"name" gorm:"size:100;not null;comment:{{.Comment}}名称"` + "`" + `
	Description string ` + "`" + `json:"description" gorm:"size:500;comment:{{.Comment}}描述"` + "`" + `
	Status      int8    ` + "`" + `json:"status" gorm:"default:1;comment:状态 1-启用 0-禁用"` + "`" + `
	// 可以根据需要添加更多字段
}
`

// 请求模板
var requestTmpl = `package req

// {{.Name}}CreateReq {{.Comment}}创建请求
type {{.Name}}CreateReq struct {
	Name        string ` + "`" + `json:"name" binding:"required" msg:"请输入{{.Comment}}名称"` + "`" + `
	Description string ` + "`" + `json:"description"` + "`" + `
	Status      int8    ` + "`" + `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"` + "`" + `
	// 可以根据需要添加更多字段
}

// {{.Name}}UpdateReq {{.Comment}}更新请求
type {{.Name}}UpdateReq struct {
	Name        string ` + "`" + `json:"name" binding:"required" msg:"请输入{{.Comment}}名称"` + "`" + `
	Description string ` + "`" + `json:"description"` + "`" + `
	Status      int8    ` + "`" + `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"` + "`" + `
	// 可以根据需要添加更多字段
}
`

// 基础请求模板
var baseRequestTmpl = `package req

// IdReq ID请求
type IdReq struct {
	ID uint ` + "`" + `uri:"id" form:"id" binding:"required,min=1"` + "`" + `
}

// PageReq 分页请求
type PageReq struct {
	Page     int ` + "`" + `form:"page" binding:"required,min=1"` + "`" + `
	PageSize int ` + "`" + `form:"pageSize" binding:"required,min=1,max=100"` + "`" + `
}

// GetOffset 获取偏移量
func (r *PageReq) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// GetLimit 获取限制
func (r *PageReq) GetLimit() int {
	return r.PageSize
}
`

// 模块模板
var moduleTmpl = `package {{.Module}}

import (
	"fmt"

	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/controller"
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/repository"
	"github.com/zhoudm1743/go-frame/internal/{{.Module}}/service"
	"github.com/zhoudm1743/go-frame/pkg/http"
	"github.com/zhoudm1743/go-frame/pkg/http/unified"
	"github.com/zhoudm1743/go-frame/pkg/log"
	"go.uber.org/fx"
)

// Router 路由接口
type Router interface {
	RegisterRoutes(router unified.Router)
}

// Module 模块定义
type Module struct {
	moduleID string
	name     string
}

// NewModule 创建模块
func NewModule() *Module {
	return &Module{
		moduleID: "{{.Module}}",
		name:     "{{.ModuleName}}",
	}
}

// Name 模块名称
func (m *Module) Name() string {
	return m.name
}

// RoutePrefix 获取模块路由前缀
func (m *Module) RoutePrefix() string {
	return "/{{.Module}}"
}

// Options 模块配置选项
func (m *Module) Options() fx.Option {
	// 将模块ID作为fx命名的一部分，避免多个模块实例冲突
	name := fmt.Sprintf("module_%s", m.moduleID)

	return fx.Module(
		name,
		// 提供所有依赖
		fx.Provide(
			repository.New{{.Name}}Repository,
			service.New{{.Name}}Service,
			controller.New{{.Name}}Controller,
			controller.New{{.Name}}Router,
			// 注册路由组
			fx.Annotate(
				func(router *controller.{{.Name}}Router) Router {
					return router
				},
				fx.ResultTags(fmt.Sprintf(` + "`" + `group:"%s_routers"` + "`" + `, m.moduleID)),
			),
		),
		// 注册路由
		fx.Invoke(
			fx.Annotate(
				func(server http.Server, logger log.Logger, routers []Router) {
					logger.Infof("注册模块: %s", m.Name())
					for _, router := range routers {
						router.RegisterRoutes(server.Router())
					}
					logger.Infof("模块 %s 路由注册完成", m.moduleID)
				},
				fx.ParamTags(` + "`" + `` + "`" + `, ` + "`" + `` + "`" + `, fmt.Sprintf(` + "`" + `group:"%s_routers"` + "`" + `, m.moduleID)),
			),
		),
	)
}
`

// GenerateCode 生成代码
func GenerateCode(name, module, comment string) error {
	// 转换名称为首字母大写
	name = strings.Title(name)

	// 生成输出目录
	outputDir := filepath.Join("internal", module)

	// 生成各个文件
	if err := generateController(name, module, comment, outputDir); err != nil {
		return fmt.Errorf("生成控制器失败: %w", err)
	}

	if err := generateRouter(name, module, comment, outputDir); err != nil {
		return fmt.Errorf("生成路由失败: %w", err)
	}

	if err := generateService(name, module, comment, outputDir); err != nil {
		return fmt.Errorf("生成服务失败: %w", err)
	}

	if err := generateRepository(name, module, comment, outputDir); err != nil {
		return fmt.Errorf("生成仓库失败: %w", err)
	}

	if err := generateModel(name, module, comment, outputDir); err != nil {
		return fmt.Errorf("生成模型失败: %w", err)
	}

	if err := generateRequest(name, module, comment, outputDir); err != nil {
		return fmt.Errorf("生成请求失败: %w", err)
	}

	// 检查模块文件是否存在
	moduleFilePath := filepath.Join(outputDir, "module.go")
	if _, err := os.Stat(moduleFilePath); os.IsNotExist(err) {
		// 模块文件不存在，生成新的模块文件
		if err := generateModuleFile(name, module, comment, outputDir); err != nil {
			return fmt.Errorf("生成模块失败: %w", err)
		}
	}

	// 检查基础请求文件是否存在
	baseReqFilePath := filepath.Join(outputDir, "schemas", "req", "base.go")
	if _, err := os.Stat(baseReqFilePath); os.IsNotExist(err) {
		// 基础请求文件不存在，生成新的基础请求文件
		if err := generateBaseRequest(module, outputDir); err != nil {
			return fmt.Errorf("生成基础请求失败: %w", err)
		}
	}

	// 确保响应目录存在
	respDir := filepath.Join(outputDir, "schemas", "resp")
	if err := os.MkdirAll(respDir, 0755); err != nil {
		return fmt.Errorf("生成响应目录失败: %w", err)
	}

	return nil
}

// 生成控制器文件
func generateController(name, module, comment, outputDir string) error {
	data := struct {
		Name    string
		Module  string
		Comment string
		Package string
	}{
		Name:    name,
		Module:  module,
		Comment: comment,
		Package: "controller",
	}

	// 解析模板
	tmpl, err := template.New("controller").Parse(controllerTmpl)
	if err != nil {
		return err
	}

	// 生成文件路径
	filePath := filepath.Join(outputDir, "controller", strings.ToLower(name)+"_controller.go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 生成文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// 生成路由文件
func generateRouter(name, module, comment, outputDir string) error {
	data := struct {
		Name      string
		Module    string
		Comment   string
		Package   string
		LowerName string
	}{
		Name:      name,
		Module:    module,
		Comment:   comment,
		Package:   "controller",
		LowerName: strings.ToLower(name),
	}

	// 解析模板
	tmpl, err := template.New("router").Parse(routerTmpl)
	if err != nil {
		return err
	}

	// 生成文件路径
	filePath := filepath.Join(outputDir, "controller", strings.ToLower(name)+"_router.go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 生成文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// 生成服务文件
func generateService(name, module, comment, outputDir string) error {
	data := struct {
		Name      string
		Module    string
		Comment   string
		LowerName string
	}{
		Name:      name,
		Module:    module,
		Comment:   comment,
		LowerName: strings.ToLower(name),
	}

	// 解析模板
	tmpl, err := template.New("service").Parse(serviceTmpl)
	if err != nil {
		return err
	}

	// 生成文件路径
	filePath := filepath.Join(outputDir, "service", strings.ToLower(name)+"_service.go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 生成文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// 生成仓库文件
func generateRepository(name, module, comment, outputDir string) error {
	data := struct {
		Name    string
		Module  string
		Comment string
	}{
		Name:    name,
		Module:  module,
		Comment: comment,
	}

	// 解析模板
	tmpl, err := template.New("repository").Parse(repositoryTmpl)
	if err != nil {
		return err
	}

	// 生成文件路径
	filePath := filepath.Join(outputDir, "repository", strings.ToLower(name)+"_repository.go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 生成文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// 生成模型文件
func generateModel(name, module, comment, outputDir string) error {
	data := struct {
		Name    string
		Module  string
		Comment string
	}{
		Name:    name,
		Module:  module,
		Comment: comment,
	}

	// 解析模板
	tmpl, err := template.New("model").Parse(modelTmpl)
	if err != nil {
		return err
	}

	// 生成文件路径
	filePath := filepath.Join(outputDir, "model", strings.ToLower(name)+".go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 生成文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// 生成请求文件
func generateRequest(name, module, comment, outputDir string) error {
	data := struct {
		Name    string
		Module  string
		Comment string
	}{
		Name:    name,
		Module:  module,
		Comment: comment,
	}

	// 解析模板
	tmpl, err := template.New("request").Parse(requestTmpl)
	if err != nil {
		return err
	}

	// 生成文件路径
	filePath := filepath.Join(outputDir, "schemas", "req", strings.ToLower(name)+".go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 生成文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// 生成基础请求文件
func generateBaseRequest(module, outputDir string) error {
	// 生成文件路径
	filePath := filepath.Join(outputDir, "schemas", "req", "base.go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, []byte(baseRequestTmpl), 0644)
}

// 生成模块文件
func generateModuleFile(name, module, comment, outputDir string) error {
	data := struct {
		Name       string
		Module     string
		Comment    string
		ModuleName string
	}{
		Name:       name,
		Module:     module,
		Comment:    comment,
		ModuleName: strings.Title(module),
	}

	// 解析模板
	tmpl, err := template.New("module").Parse(moduleTmpl)
	if err != nil {
		return err
	}

	// 生成文件路径
	filePath := filepath.Join(outputDir, "module.go")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// 生成文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
