package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/transport"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/leaf-rain/raindata/common/str"
	"strconv"
	"strings"
	"time"
)

type AuthUser struct {
	Path   string
	Method string
	Domain string
	Token  string

	IssTime     string
	AuthorityId string
}

// ParseFromContext 从上下文中解析出JWT Claims和操作信息，并填充到AuthUser中。
func (su *AuthUser) ParseFromContext(ctx context.Context, secretKey []byte) error {
	// 从上下文中获取操作信息
	header, ok := transport.FromServerContext(ctx)
	if ok {
		su.Path = header.Operation()
		su.Method = string(header.Kind())
		su.Domain = header.Endpoint()
		su.Token = header.RequestHeader().Get("Authorization")
		su.Token = strings.TrimPrefix(su.Token, "Bearer ")
	} else {
		return errors.New("jwt claim missing")
	}
	// 校验jwt
	parseAuth, err := jwtv5.Parse(su.Token, func(*jwtv5.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	// 从上下文中获取JWT Claims
	claims, _ := parseAuth.Claims.(jwtv5.MapClaims)
	if err := su.ParseAccessJwtToken(claims); err != nil {
		return err
	}
	return nil
}

// 以下方法获取AuthUser的属性值

func (su *AuthUser) GetDomain() string {
	return su.Domain
}

func (su *AuthUser) GetSubject() string {
	return su.AuthorityId
}

func (su *AuthUser) GetObject() string {
	return su.Path
}

func (su *AuthUser) GetAction() string {
	return su.Method
}

// CreateAccessJwtToken 使用给定的密钥为AuthUser创建一个JWT访问令牌。
func (su *AuthUser) CreateAccessJwtToken(secretKey []byte) string {
	claims := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256,
		jwtv5.MapClaims{
			"AuthorityId": su.AuthorityId,
			"IssTime":     strconv.FormatInt(time.Now().Unix(), 10),
		})
	signedToken, err := claims.SignedString(secretKey)
	if err != nil {
		return ""
	}
	return signedToken
}

// ParseAccessJwtTokenFromString 解析给定的JWT令牌字符串，并填充到AuthUser中。
func (su *AuthUser) ParseAccessJwtTokenFromString(token string, secretKey []byte) error {
	parseAuth, err := jwtv5.Parse(token, func(*jwtv5.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	claims, ok := parseAuth.Claims.(jwtv5.MapClaims)
	if !ok {
		return errors.New("no jwt token in context")
	}
	return su.ParseAccessJwtToken(claims)
}

// ParseAccessJwtToken 从JWT Claims解析并填充到AuthUser中。
func (su *AuthUser) ParseAccessJwtToken(claims jwtv5.Claims) error {
	if claims == nil {
		return errors.New("claims is nil")
	}
	mc, ok := claims.(jwtv5.MapClaims)
	if !ok {
		return errors.New("claims is not map claims")
	}
	su.AuthorityId = str.GetDataByInterface("AuthorityId", mc)
	su.IssTime = str.GetDataByInterface("IssTime", mc)
	return nil
}
