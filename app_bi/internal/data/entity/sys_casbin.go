package entity

import (
	"context"
	"errors"
	adapter "github.com/casbin/gorm-adapter/v3"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ initDb = (*EntitySysCasbin)(nil)

type EntitySysCasbin struct {
	data *Data
}

func NewEntitySysCasbin(data *Data) *EntitySysCasbin {
	return &EntitySysCasbin{
		data: data,
	}
}

func (i *EntitySysCasbin) MigrateTable(ctx context.Context) error {
	return i.data.SqlClient.AutoMigrate(&adapter.CasbinRule{})
}

func (i *EntitySysCasbin) TableCreated(context.Context) bool {
	return i.data.SqlClient.Migrator().HasTable(&adapter.CasbinRule{})
}

func (i *EntitySysCasbin) InitializerName() string {
	var entity adapter.CasbinRule
	return entity.TableName()
}

func (i *EntitySysCasbin) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/user/admin_register", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/createApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/getApiList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/getApiById", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/deleteApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/updateApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/getAllApis", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/deleteApisByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/api/syncApi", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/getApiGroups", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/enterSyncApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/ignoreApi", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/authority/copyAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/updateAuthority", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/authority/createAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/deleteAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/getAuthorityList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/setDataAuthority", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/user/setUserInfo", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/user/setSelfInfo", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/user/getUserList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/deleteUser", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/user/changePassword", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/setUserAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/setUserAuthorities", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/resetPassword", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/findFile", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/breakpointContinueFinish", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/breakpointContinue", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/removeChunk", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/upload", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/deleteFile", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/editFileName", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/getFileList", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/casbin/updateCasbin", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/casbin/getPolicyPathByAuthorityId", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/jwt/jsonInBlacklist", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/system/getSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/system/setSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/system/getServerInfo", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/customer/customerList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/findSysOperationRecord", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/updateSysOperationRecord", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/createSysOperationRecord", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/getSysOperationRecordList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/deleteSysOperationRecord", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/deleteSysOperationRecordByIds", V2: "DELETE"},

		{Ptype: "p", V0: "888", V1: "/email/emailTest", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/email/sendEmail", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/simpleUploader/upload", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/simpleUploader/checkFileMd5", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/simpleUploader/mergeFileMd5", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/authorityBtn/setAuthorityBtn", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authorityBtn/getAuthorityBtn", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authorityBtn/canRemoveAuthorityBtn", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/createSysExportTemplate", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/deleteSysExportTemplate", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/deleteSysExportTemplateByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/updateSysExportTemplate", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/findSysExportTemplate", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/getSysExportTemplateList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/exportExcel", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/exportTemplate", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/importExcel", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/info/createInfo", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/info/deleteInfo", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/info/deleteInfoByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/info/updateInfo", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/info/findInfo", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/info/getInfoList", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/user/admin_register", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/createApi", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/getApiList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/getApiById", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/deleteApi", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/updateApi", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/getAllApis", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/createAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/deleteAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/getAuthorityList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/setDataAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/user/changePassword", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/user/getUserList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/user/setUserAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/upload", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/getFileList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/deleteFile", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/editFileName", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/casbin/updateCasbin", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/casbin/getPolicyPathByAuthorityId", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/jwt/jsonInBlacklist", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/system/getSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/system/setSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "PUT"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "DELETE"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/customer/customerList", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/user/getUserInfo", V2: "GET"},

		{Ptype: "p", V0: "9528", V1: "/user/admin_register", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/createApi", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/getApiList", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/getApiById", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/deleteApi", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/updateApi", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/getAllApis", V2: "POST"},

		{Ptype: "p", V0: "9528", V1: "/authority/createAuthority", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/authority/deleteAuthority", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/authority/getAuthorityList", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/authority/setDataAuthority", V2: "POST"},

		{Ptype: "p", V0: "9528", V1: "/menu/getMenu", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/getMenuList", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/addBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/getBaseMenuTree", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/addMenuAuthority", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/getMenuAuthority", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/deleteBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/updateBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/menu/getBaseMenuById", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/user/changePassword", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/user/getUserList", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/user/setUserAuthority", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/fileUploadAndDownload/upload", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/fileUploadAndDownload/getFileList", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/fileUploadAndDownload/deleteFile", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/fileUploadAndDownload/editFileName", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/casbin/updateCasbin", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/casbin/getPolicyPathByAuthorityId", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/jwt/jsonInBlacklist", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/system/getSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/system/setSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/customer/customer", V2: "PUT"},
		{Ptype: "p", V0: "9528", V1: "/customer/customer", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/customer/customer", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/customer/customer", V2: "DELETE"},
		{Ptype: "p", V0: "9528", V1: "/customer/customerList", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/autoCode/createTemp", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/user/getUserInfo", V2: "GET"},
	}
	if err := i.data.SqlClient.Create(&entities).Error; err != nil {
		return ctx, errors2.Wrap(err, "Casbin 表 ("+i.InitializerName()+") 数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *EntitySysCasbin) DataInserted(ctx context.Context) bool {
	if errors.Is(i.data.SqlClient.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/user/getUserInfo", V2: "GET"}).
		First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
