package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
)

//@function: UpdateCasbin
//@description: 更新casbin权限
//@param: authorityId string, casbinInfos []dto.CasbinInfo
//@return: error

type CasbinService struct {
	*Service
}

func NewCasbinService(service *Service) *CasbinService {
	return &CasbinService{
		service,
	}
}

func (svc *CasbinService) UpdateCasbin(AuthorityID uint, casbinInfos []dto.CasbinInfo) error {
	b := biz.NewCasbin(svc.biz)
	return b.UpdateCasbin(AuthorityID, casbinInfos)
}

//@function: UpdateCasbinApi
//@description: API更新随动
//@param: oldPath string, newPath string, oldMethod string, newMethod string
//@return: error

func (svc *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	b := biz.NewCasbin(svc.biz)
	return b.UpdateCasbinApi(oldPath, newPath, oldMethod, newMethod)
}

//@function: GetPolicyPathByAuthorityId
//@description: 获取权限列表
//@param: authorityId string
//@return: pathMaps []dto.CasbinInfo

func (svc *CasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []dto.CasbinInfo) {
	b := biz.NewCasbin(svc.biz)
	return b.GetPolicyPathByAuthorityId(AuthorityID)
}

//@function: ClearCasbin
//@description: 清除匹配的权限
//@param: v int, p ...string
//@return: bool

func (svc *CasbinService) ClearCasbin(v int, p ...string) bool {
	b := biz.NewCasbin(svc.biz)
	return b.ClearCasbin(v, p...)
}

//@function: RemoveFilteredPolicy
//@description: 使用数据库方法清理筛选的politicy 此方法需要调用FreshCasbin方法才可以在系统中即刻生效
//@param: db *gorm.DB, authorityId string
//@return: error

func (svc *CasbinService) RemoveFilteredPolicy(authorityId string) error {
	b := biz.NewCasbin(svc.biz)
	return b.RemoveFilteredPolicy(authorityId)
}

//@function: SyncPolicy
//@description: 同步目前数据库的policy 此方法需要调用FreshCasbin方法才可以在系统中即刻生效
//@param: db *gorm.DB, authorityId string, rules [][]string
//@return: error

func (svc *CasbinService) SyncPolicy(authorityId string, rules [][]string) error {
	b := biz.NewCasbin(svc.biz)
	return b.SyncPolicy(authorityId, rules)
}

//@function: AddPolicies
//@description: 添加匹配的权限
//@param: v int, p ...string
//@return: bool

func (svc *CasbinService) AddPolicies(rules [][]string) error {
	b := biz.NewCasbin(svc.biz)
	return b.AddPolicies(rules)
}

func (svc *CasbinService) FreshCasbin() (err error) {
	b := biz.NewCasbin(svc.biz)
	return b.FreshCasbin()
}

//@function: Casbin
//@description: 持久化到数据库  引入自定义规则
//@return: *casbin.Enforcer

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

func (svc *CasbinService) Casbin() *casbin.SyncedCachedEnforcer {
	b := biz.NewCasbin(svc.biz)
	return b.Casbin()
}
