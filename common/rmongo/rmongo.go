package rmongo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var ProviderSet = wire.NewSet(NewMongo)

var mc *mongo.Client
var mcOnce = sync.Once{}

func NewMongo(ctx context.Context, mogConf *MongoCfg, logger log.Logger) (*mongo.Database, error) {
	if mc != nil {
		return mc.Database(mogConf.Db), nil
	}
	var err error
	mcOnce.Do(func() {
		clientOption := options.Client().ApplyURI(mogConf.Url).SetConnectTimeout(time.Duration(mogConf.ConnectTimeoutMS) * time.Millisecond)
		if mogConf.MaxPoolSize > 0 {
			clientOption = clientOption.SetMaxPoolSize(mogConf.MaxPoolSize)
		}
		if mogConf.MaxPoolSize > 0 {
			clientOption = clientOption.SetMinPoolSize(mogConf.MaxPoolSize)
		}
		if mogConf.MaxIdleTimeMS > 0 {
			clientOption = clientOption.SetMaxConnIdleTime(time.Duration(mogConf.MaxIdleTimeMS) * time.Millisecond)
		}
		mc, err = mongo.Connect(ctx, clientOption)
	})
	if err != nil {
		log.NewHelper(logger).Errorf("mongo open error: %s, mogConf %#v \n", err.Error(), mogConf)
		return nil, err
	}
	return mc.Database(mogConf.Db), nil
}
