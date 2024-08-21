package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"go.uber.org/zap"
	"net"
	"time"
)

type JWT struct {
	SigningKey  []byte
	BufferTime  string
	ExpiresTime string
	Issuer      string
	logger      *zap.Logger
}

func NewJWT(signing []byte, logger *zap.Logger) *JWT {
	return &JWT{
		SigningKey: signing,
		logger:     logger,
	}
}

func (j *JWT) CreateClaims(baseClaims dto.BaseClaims) dto.CustomClaims {
	bf, _ := ParseDuration(j.BufferTime)
	ep, _ := ParseDuration(j.ExpiresTime)
	claims := dto.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                   // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    j.Issuer,                                  // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims dto.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims dto.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*dto.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ERR_TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ERR_TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ERR_TokenNotValidYet
			} else {
				return nil, ERR_TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*dto.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ERR_TokenInvalid

	} else {
		return nil, ERR_TokenInvalid
	}
}

func ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("x-token", "", -1, "/", host, false, false)
	}
}

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}

func GetToken(c *gin.Context) string {
	token, _ := c.Cookie("x-token")
	if token == "" {
		token = c.Request.Header.Get("x-token")
	}
	return token
}

func GetClaims(c *gin.Context) (*dto.CustomClaims, error) {
	token := GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.ID
		}
	} else {
		waitUse := claims.(*dto.CustomClaims)
		return waitUse.BaseClaims.ID
	}
}

// GetUserUuid 从Gin的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*dto.CustomClaims)
		return waitUse.UUID
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.AuthorityId
		}
	} else {
		waitUse := claims.(*dto.CustomClaims)
		return waitUse.AuthorityId
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *dto.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*dto.CustomClaims)
		return waitUse
	}
}

// GetUserName 从Gin的Context中获取从jwt解析出来的用户名
func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*dto.CustomClaims)
		return waitUse.Username
	}
}

func LoginToken(user data.Login) (token string, claims dto.CustomClaims, err error) {
	j := &JWT{SigningKey: []byte(svc.conf.JWT.SigningKey)} // 唯一签名
	claims = j.CreateClaims(dto.BaseClaims{
		UUID:        user.GetUUID(),
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
