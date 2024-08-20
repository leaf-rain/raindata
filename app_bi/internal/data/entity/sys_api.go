package entity

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SysApi struct {
	gorm.Model
	Path        string `json:"path" gorm:"comment:api路径"`             // api路径
	Description string `json:"description" gorm:"comment:api中文描述"`    // api中文描述
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`          // api组
	Method      string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (SysApi) TableName() string {
	return "sys_apis"
}

var _ initDb = (*EntitySysApi)(nil)

type EntitySysApi struct {
	data  *Data
	Model *SysApi
}

func NewEntitySysApi(data *Data) *EntitySysApi {
	return &EntitySysApi{
		data: data,
	}
}

func (i *EntitySysApi) MigrateTable(ctx context.Context) error {
	return i.data.SqlClient.AutoMigrate(&SysApi{})
}

func (i *EntitySysApi) TableCreated(context.Context) bool {
	return i.data.SqlClient.Migrator().HasTable(&SysApi{})
}

func (i *EntitySysApi) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "jwt加入黑名单(退出，必选)"},

		{ApiGroup: "系统用户", Method: "DELETE", Path: "/user/deleteUser", Description: "删除用户"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/admin_register", Description: "用户注册"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/getUserList", Description: "获取用户列表"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setUserInfo", Description: "设置用户信息"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfInfo", Description: "设置自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "GET", Path: "/user/getUserInfo", Description: "获取自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthorities", Description: "设置权限组"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/changePassword", Description: "修改密码（建议选择)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthority", Description: "修改用户角色(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/resetPassword", Description: "重置用户密码"},

		{ApiGroup: "api", Method: "POST", Path: "/api/createApi", Description: "创建api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/deleteApi", Description: "删除Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/updateApi", Description: "更新Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiList", Description: "获取api列表"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getAllApis", Description: "获取所有api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiById", Description: "获取api详细信息"},
		{ApiGroup: "api", Method: "DELETE", Path: "/api/deleteApisByIds", Description: "批量删除api"},
		{ApiGroup: "api", Method: "GET", Path: "/api/syncApi", Description: "获取待同步API"},
		{ApiGroup: "api", Method: "GET", Path: "/api/getApiGroups", Description: "获取路由组"},
		{ApiGroup: "api", Method: "POST", Path: "/api/enterSyncApi", Description: "确认同步API"},
		{ApiGroup: "api", Method: "POST", Path: "/api/ignoreApi", Description: "忽略API"},

		{ApiGroup: "角色", Method: "POST", Path: "/authority/copyAuthority", Description: "拷贝角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/createAuthority", Description: "创建角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/deleteAuthority", Description: "删除角色"},
		{ApiGroup: "角色", Method: "PUT", Path: "/authority/updateAuthority", Description: "更新角色信息"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/getAuthorityList", Description: "获取角色列表"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/setDataAuthority", Description: "设置角色资源权限"},

		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "更改角色api权限"},
		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "获取权限列表"},

		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/upload", Description: "文件上传示例"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/deleteFile", Description: "删除文件"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/editFileName", Description: "文件名或者备注编辑"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/getFileList", Description: "获取上传文件列表"},

		{ApiGroup: "系统服务", Method: "POST", Path: "/system/getServerInfo", Description: "获取服务器信息"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/system/getSystemConfig", Description: "获取配置文件内容"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/system/setSystemConfig", Description: "设置配置文件内容"},

		{ApiGroup: "客户", Method: "PUT", Path: "/customer/customer", Description: "更新客户"},
		{ApiGroup: "客户", Method: "POST", Path: "/customer/customer", Description: "创建客户"},
		{ApiGroup: "客户", Method: "DELETE", Path: "/customer/customer", Description: "删除客户"},
		{ApiGroup: "客户", Method: "GET", Path: "/customer/customer", Description: "获取单一客户"},
		{ApiGroup: "客户", Method: "GET", Path: "/customer/customerList", Description: "获取客户列表"},

		{ApiGroup: "操作记录", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "新增操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "根据ID获取操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "获取操作记录列表"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "删除操作记录"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "批量删除操作历史"},

		{ApiGroup: "email", Method: "POST", Path: "/email/emailTest", Description: "发送测试邮件"},
		{ApiGroup: "email", Method: "POST", Path: "/email/sendEmail", Description: "发送邮件"},

		{ApiGroup: "公告", Method: "POST", Path: "/info/createInfo", Description: "新建公告"},
		{ApiGroup: "公告", Method: "DELETE", Path: "/info/deleteInfo", Description: "删除公告"},
		{ApiGroup: "公告", Method: "DELETE", Path: "/info/deleteInfoByIds", Description: "批量删除公告"},
		{ApiGroup: "公告", Method: "PUT", Path: "/info/updateInfo", Description: "更新公告"},
		{ApiGroup: "公告", Method: "GET", Path: "/info/findInfo", Description: "根据ID获取公告"},
		{ApiGroup: "公告", Method: "GET", Path: "/info/getInfoList", Description: "获取公告列表"},
	}
	if err := i.data.SqlClient.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, SysApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.Model.TableName(), entities)
	return next, nil
}

func (i *EntitySysApi) DataInserted(ctx context.Context) bool {
	if errors.Is(i.data.SqlClient.Where("path = ? AND method = ?", "/authorityBtn/canRemoveAuthorityBtn", "POST").
		First(&SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
