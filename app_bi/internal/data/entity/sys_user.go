package entity

import (
	"context"
	"github.com/google/uuid"
	"github.com/leaf-rain/raindata/app_bi/third_party/hash"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"github.com/leaf-rain/raindata/common/ecode"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Login interface {
	GetUsername() string
	GetNickname() string
	GetUUID() uuid.UUID
	GetUserId() uint
	GetAuthorityId() uint
	GetUserInfo() any
}

var _ Login = (*SysUser)(nil)

type SysUser struct {
	gorm.Model
	UUID        uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`                // 用户UUID
	Username    string         `json:"userName" gorm:"index;comment:用户登录名"`             // 用户登录名
	Password    string         `json:"password"  gorm:"comment:用户登录密码"`                 // 用户登录密码
	NickName    string         `json:"nickName" gorm:"comment:用户昵称"`                    // 用户昵称
	SideMode    string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`     // 用户侧边主题
	HeaderImg   string         `json:"headerImg" gorm:"comment:用户头像"`                   // 用户头像
	BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`      // 基础颜色
	AuthorityId uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`   // 用户角色ID
	Phone       string         `json:"phone"  gorm:"comment:用户手机号"`                     // 用户手机号
	Email       string         `json:"email"  gorm:"comment:用户邮箱"`                      // 用户邮箱
	Enable      int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
	Authority   SysAuthority   `json:"authority" gorm:"-"`
	Authorities []SysAuthority `json:"authorities" gorm:"-"`
}

func (SysUser) TableName() string {
	return "sys_users"
}

func (s *SysUser) GetUsername() string {
	return s.Username
}

func (s *SysUser) GetNickname() string {
	return s.NickName
}

func (s *SysUser) GetUUID() uuid.UUID {
	return s.UUID
}

func (s *SysUser) GetUserId() uint {
	return s.ID
}

func (s *SysUser) GetAuthorityId() uint {
	return s.AuthorityId
}

func (s *SysUser) GetUserInfo() any {
	return *s
}

var _ initDb = (*EntitySysUser)(nil)

type EntitySysUser struct {
	data  *Data
	Model *SysUser
}

func NewEntitySysUser(data *Data) *EntitySysUser {
	return &EntitySysUser{
		data: data,
	}
}

func (entity *EntitySysUser) MigrateTable(ctx context.Context) error {
	return entity.data.SqlClient.AutoMigrate(&SysUser{})
}

func (entity *EntitySysUser) TableCreated(context.Context) bool {
	return entity.data.SqlClient.Migrator().HasTable(&SysUser{})
}

func (entity EntitySysUser) InitializerName() string {
	return SysUser{}.TableName()
}

func (entity *EntitySysUser) InitializeData(ctx context.Context) (next context.Context, err error) {
	ap := ctx.Value("adminPassword")
	apStr, ok := ap.(string)
	if !ok {
		apStr = "123456"
	}
	adminPassword := utils.BcryptHash(apStr)
	entities := []SysUser{
		{
			UUID:        uuid.Must(uuid.NewV7()),
			Username:    "admin",
			Password:    adminPassword,
			NickName:    "root",
			HeaderImg:   "https://img1.baidu.com/it/u=1657712229,2620982189&fm=253&app=120&size=w931&n=0&f=JPEG&fmt=auto?sec=1724173200&t=990ed5fcca0d90a914a0afb7c1a3b3b8",
			AuthorityId: 888,
			Phone:       "17611111111",
			Email:       "111111111@qq.com",
		},
	}
	if err = entity.data.SqlClient.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, SysUser{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, entity.InitializerName(), entities)
	authorityEntities, ok := ctx.Value(entity.Model.TableName()).([]SysAuthority)
	if !ok {
		return next, errors.Wrap(ecode.ErrMissingDependentContext, "创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err = entity.data.SqlClient.Model(&entities[0]).Association("Authorities").Replace(authorityEntities); err != nil {
		return next, err
	}
	if err = entity.data.SqlClient.Model(&entities[1]).Association("Authorities").Replace(authorityEntities[:1]); err != nil {
		return next, err
	}
	return next, err
}

func (entity *EntitySysUser) ReloadByDb() error {
	var user *SysUser
	err := entity.data.SqlClient.Where("username = ?", entity.Model.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := hash.BcryptCheck(entity.Model.Password, user.Password); !ok {
			return ecode.ERR_USER_AUTH
		}
	}
	entity.Model = user
	return err
}
