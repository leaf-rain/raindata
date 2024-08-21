package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"net"
)

func (b *UserApi) ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.dto.Host)
	if err != nil {
		host = c.dto.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("x-token", "", -1, "/", host, false, false)
	}
}

func (b *UserApi) SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.dto.Host)
	if err != nil {
		host = c.dto.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}

func (b *UserApi) GetToken(c *gin.Context) string {
	token, _ := c.Cookie("x-token")
	if token == "" {
		token = c.dto.Header.Get("x-token")
	}
	return token
}

func (b *UserApi) GetClaims(c *gin.Context) (*data.CustomClaims, error) {
	token := b.GetToken(c)
	j := data.NewJWT(b.data)
	claims, err := j.ParseToken(token)
	if err != nil {
		b.log.Error("[GetClaims]从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func (b *UserApi) GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := b.GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.ID
		}
	} else {
		waitUse := claims.(*data.CustomClaims)
		return waitUse.BaseClaims.ID
	}
}

// GetUserUuid 从Gin的Context中获取从jwt解析出来的用户UUID
func (b *UserApi) GetUserUuid(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := b.GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*data.CustomClaims)
		return waitUse.UUID
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func (b *UserApi) GetUserAuthorityId(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := b.GetClaims(c); err != nil {
			return 0
		} else {
			return cl.AuthorityId
		}
	} else {
		waitUse := claims.(*data.CustomClaims)
		return waitUse.AuthorityId
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func (b *UserApi) GetUserInfoByCtx(c *gin.Context) *data.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := b.GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*data.CustomClaims)
		return waitUse
	}
}

// GetUserName 从Gin的Context中获取从jwt解析出来的用户名
func (b *UserApi) GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := b.GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*data.CustomClaims)
		return waitUse.Username
	}
}

func (b *UserApi) LoginToken(user data.Login) (token string, claims data.CustomClaims, err error) {
	j := data.NewJWT(b.data) // 唯一签名
	claims = j.CreateClaims(data.BaseClaims{
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
