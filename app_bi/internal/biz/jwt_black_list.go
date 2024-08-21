package biz

import (
	"context"
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

type Jwt struct {
	*Business
}

func NewJwt(biz *Business) *Jwt {
	return &Jwt{
		biz,
	}
}

//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList data.JwtBlacklist
//@return: err error

func (b *Jwt) JsonInBlacklist(jwtList data.JwtBlacklist) (err error) {
	err = b.data.SqlClient.Create(&jwtList).Error
	if err != nil {
		return
	}
	return
}

//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt.proto string
//@return: bool

func (b *Jwt) IsBlacklist(jwt string) bool {
	err := b.data.SqlClient.Where("jwt.proto = ?", jwt).First(&data.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}

//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (b *Jwt) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = b.data.RdClient.Get(context.Background(), userName).Result()
	return redisJWT, err
}

//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt.proto string, userName string
//@return: err error

func (b *Jwt) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(b.data.Config.Jwt.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = b.data.RdClient.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

func (b *Jwt) LoadAll() {
	var records []string
	err := b.data.SqlClient.Model(&data.JwtBlacklist{}).Select("jwt.proto").Find(&records).Error
	if err != nil {
		b.logger.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
}
