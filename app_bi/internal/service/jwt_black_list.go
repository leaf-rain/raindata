package service

import (
	"context"
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

type JwtService struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var JwtServiceApp = new(JwtService)

//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList data.JwtBlacklist
//@return: err error

func (svc *JwtService) JsonInBlacklist(jwtList entity.JwtBlacklist) (err error) {
	err = svc.data.SqlClient.Create(&jwtList).Error
	if err != nil {
		return
	}
	return
}

//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (svc *JwtService) IsBlacklist(jwt string) bool {
	err := svc.data.SqlClient.Where("jwt = ?", jwt).First(&entity.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}

//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (svc *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = svc.data.RdClient.Get(context.Background(), userName).Result()
	return redisJWT, err
}

//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func (svc *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(svc.conf.Jwt.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = svc.data.RdClient.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

func (svc *JwtService) LoadAll() {
	var records []string
	err := svc.data.SqlClient.Model(&entity.JwtBlacklist{}).Select("jwt").Find(&records).Error
	if err != nil {
		svc.log.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
}
