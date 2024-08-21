package data

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
	AuthorityId uint           `json:"authorityId" gorm:"default:0;comment:用户角色ID"`     // 用户角色ID
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
	*Data
	Model *SysUser
}

func NewEntitySysUser(data *Data) *EntitySysUser {
	return &EntitySysUser{
		Data: data,
	}
}

func (entity *EntitySysUser) MigrateTable(ctx context.Context) error {
	return entity.SqlClient.AutoMigrate(&SysUser{})
}

func (entity *EntitySysUser) TableCreated(context.Context) bool {
	return entity.SqlClient.Migrator().HasTable(&SysUser{})
}

func (entity EntitySysUser) InitializerName() string {
	return SysUser{}.TableName()
}

func (entity *EntitySysUser) InitializeData(ctx context.Context) (err error) {
	apStr := "yeyangfengqi"
	adminPassword := utils.BcryptHash(apStr)
	entities := []SysUser{
		{
			UUID:        uuid.Must(uuid.NewV7()),
			Username:    "admin",
			Password:    adminPassword,
			NickName:    "root",
			HeaderImg:   "https://img1.baidu.com/it/u=1657712229,2620982189&fm=253&app=120&size=w931&n=0&f=JPEG&fmt=auto?sec=1724173200&t=990ed5fcca0d90a914a0afb7c1a3b3b8",
			AuthorityId: 1,
			Phone:       "17611111111",
			Email:       "111111111@qq.com",
		},
	}
	if err = entity.SqlClient.Create(&entities).Error; err != nil {
		return errors.Wrap(err, SysUser{}.TableName()+"表数据初始化失败!")
	}
	return err
}

func (entity *EntitySysUser) ReloadByDb() error {
	var user *SysUser
	err := entity.SqlClient.Where("username = ?", entity.Model.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := hash.BcryptCheck(entity.Model.Password, user.Password); !ok {
			return ecode.ERR_USER_AUTH
		} else {
			entity.Model = user
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ecode.ERR_USER_NOTFOUND
	}
	return err
}

func (entity *EntitySysUser) CreateUser() error {
	return entity.SqlClient.Create(&entity.Model).Error
}
