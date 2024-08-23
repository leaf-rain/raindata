package data

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/leaf-rain/raindata/common/ecode"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type JwtBlacklist struct {
	gorm.Model
	Jwt string `gorm:"type:text;comment:jwt.proto"`
}

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
}

var _ initDb = (*JWT)(nil)

type JWT struct {
	*Data
	Model *BaseClaims
}

func (j *JWT) InitializeData(ctx context.Context) error {
	return nil
}

func NewJWT(data *Data) *JWT {
	return &JWT{
		Data: data,
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	bf, _ := ParseDuration(j.Config.GetJwt().BufferTime)
	ep, _ := ParseDuration(j.Config.GetJwt().ExpiresTime)
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                   // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    j.Config.GetJwt().Issuer,                  // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Config.GetJwt().SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
	v, err, _ := j.SingleflightGroup.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.Config.GetJwt().SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ecode.ERR_TOKEN_MALFORMED
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ecode.ERR_TOKEN_EXPIRED
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ecode.ERR_TOKEN_NOT_VALID_YET
			} else {
				return nil, ecode.ERR_TOKEN_INVALID
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ecode.ERR_TOKEN_INVALID

	} else {
		return nil, ecode.ERR_TOKEN_INVALID
	}
}

func (j *JWT) JsonInBlacklist(jwtList JwtBlacklist) (err error) {
	err = j.SqlClient.Create(&jwtList).Error
	if err != nil {
		return
	}
	return
}

func (j *JWT) IsBlacklist(jwt string) bool {
	err := j.SqlClient.Where("jwt.proto = ?", jwt).First(&JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}

func (j *JWT) GetRedisJWT(ctx context.Context, userName string) (redisJWT string, err error) {
	redisJWT, err = j.RdClient.Get(ctx, userName).Result()
	return redisJWT, err
}

func (j *JWT) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := ParseDuration(j.Config.GetJwt().ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = j.RdClient.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

func (j *JWT) LoadAll() error {
	var data []string
	err := j.SqlClient.Model(&JwtBlacklist{}).Select("jwt.proto").Find(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")

		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			return dr, nil
		}
		return dr + ndr, nil
	}

	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}

func LoginToken(user Login, data *Data) (j *JWT, token string, claims CustomClaims, err error) {
	j = NewJWT(data) // 唯一签名
	claims = j.CreateClaims(BaseClaims{
		ID:          user.GetUserId(),
		NickName:    user.GetNickname(),
		Username:    user.GetUsername(),
		AuthorityId: user.GetAuthorityId(),
	})
	token, err = j.CreateToken(claims)
	if err != nil {
		return
	}
	return
}

func (j *JWT) MigrateTable(ctx context.Context) error {
	return j.SqlClient.AutoMigrate(&JwtBlacklist{})
}

func (j *JWT) TableCreated(context.Context) bool {
	return j.SqlClient.Migrator().HasTable(&JwtBlacklist{})
}
