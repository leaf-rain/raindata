package biz

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"strconv"
	"sync"

	adapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"go.uber.org/zap"
)

//@function: UpdateCasbin
//@description: 更新casbin权限
//@param: authorityId string, casbinInfos []dto.CasbinInfo
//@return: error

type Casbin struct {
	*Business
}

func NewCasbin(biz *Business) *Casbin {
	return &Casbin{
		biz,
	}
}

func (b *Casbin) UpdateCasbin(AuthorityID uint, casbinInfos []dto.CasbinInfo) error {
	authorityId := strconv.Itoa(int(AuthorityID))
	b.ClearCasbin(0, authorityId)
	rules := [][]string{}
	//做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfos {
		key := authorityId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
	}
	e := b.Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

//@function: UpdateCasbinApi
//@description: API更新随动
//@param: oldPath string, newPath string, oldMethod string, newMethod string
//@return: error

func (b *Casbin) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := b.data.SqlClient.Model(&adapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	e := b.Casbin()
	err = e.LoadPolicy()
	if err != nil {
		return err
	}
	return err
}

//@function: GetPolicyPathByAuthorityId
//@description: 获取权限列表
//@param: authorityId string
//@return: pathMaps []dto.CasbinInfo

func (b *Casbin) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []dto.CasbinInfo) {
	e := b.Casbin()
	authorityId := strconv.Itoa(int(AuthorityID))
	list := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, dto.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

//@function: ClearCasbin
//@description: 清除匹配的权限
//@param: v int, p ...string
//@return: bool

func (b *Casbin) ClearCasbin(v int, p ...string) bool {
	e := b.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

//@function: RemoveFilteredPolicy
//@description: 使用数据库方法清理筛选的politicy 此方法需要调用FreshCasbin方法才可以在系统中即刻生效
//@param: db *gorm.DB, authorityId string
//@return: error

func (b *Casbin) RemoveFilteredPolicy(authorityId string) error {
	return b.data.SqlClient.Delete(&adapter.CasbinRule{}, "v0 = ?", authorityId).Error
}

//@function: SyncPolicy
//@description: 同步目前数据库的policy 此方法需要调用FreshCasbin方法才可以在系统中即刻生效
//@param: db *gorm.DB, authorityId string, rules [][]string
//@return: error

func (b *Casbin) SyncPolicy(authorityId string, rules [][]string) error {
	err := b.RemoveFilteredPolicy(authorityId)
	if err != nil {
		return err
	}
	return b.AddPolicies(rules)
}

//@function: AddPolicies
//@description: 添加匹配的权限
//@param: v int, p ...string
//@return: bool

func (b *Casbin) AddPolicies(rules [][]string) error {
	var casbinRules []adapter.CasbinRule
	for i := range rules {
		casbinRules = append(casbinRules, adapter.CasbinRule{
			Ptype: "p",
			V0:    rules[i][0],
			V1:    rules[i][1],
			V2:    rules[i][2],
		})
	}
	return b.data.SqlClient.Create(&casbinRules).Error
}

func (b *Casbin) FreshCasbin() (err error) {
	e := b.Casbin()
	err = e.LoadPolicy()
	return err
}

//@function: Casbin
//@description: 持久化到数据库  引入自定义规则
//@return: *casbin.Enforcer

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

func (b *Casbin) Casbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		a, err := adapter.NewAdapterByDB(b.data.SqlClient)
		if err != nil {
			zap.L().Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", zap.Error(err))
			return
		}
		text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		m, err := model.NewModelFromString(text)
		if err != nil {
			zap.L().Error("字符串加载模型失败!", zap.Error(err))
			return
		}
		syncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		_ = syncedCachedEnforcer.LoadPolicy()
	})
	return syncedCachedEnforcer
}
