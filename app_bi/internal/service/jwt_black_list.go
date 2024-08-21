package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
)

type JwtService struct {
	*Service
}

func NewJwtService(service *Service) *JwtService {
	return &JwtService{
		service,
	}
}

//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList data.JwtBlacklist
//@return: err error

func (svc *JwtService) JsonInBlacklist(jwtList data.JwtBlacklist) (err error) {
	b := biz.NewJwt(svc.biz)
	return b.JsonInBlacklist(jwtList)
}

//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt.proto string
//@return: bool

func (svc *JwtService) IsBlacklist(jwt string) bool {
	b := biz.NewJwt(svc.biz)
	return b.IsBlacklist(jwt)
}

//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (svc *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	b := biz.NewJwt(svc.biz)
	return b.GetRedisJWT(userName)
}

//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt.proto string, userName string
//@return: err error

func (svc *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	b := biz.NewJwt(svc.biz)
	return b.SetRedisJWT(jwt, userName)
}

func (svc *JwtService) LoadAll() {
	b := biz.NewJwt(svc.biz)
	b.LoadAll()
}
