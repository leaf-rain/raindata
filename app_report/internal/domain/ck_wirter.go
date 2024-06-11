package domain

import (
	"context"
	"github.com/leaf-rain/fastjson"
	"go.uber.org/zap"
)

type CkWriter struct {
	logger *zap.Logger
}

func NewCkWriter(logger *zap.Logger) *CkWriter {
	return &CkWriter{
		logger: logger,
	}
}

func (cw *CkWriter) DataCleaning(ctx context.Context, msg string) error {
	// todo: 数据清洗， 从base中获取事件属性并进行对比，
	panic("implement me")
}

func (cw *CkWriter) WriterMsg(ctx context.Context, msg string) error {
	v, err := fastjson.Parse(msg)
	if err != nil || v == nil {
		cw.logger.Error("[WriterMsg] fastjson.Parse failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	// 校验数据
	v.Range(func(key []byte, v *fastjson.Value) {

	})
	// todo: 校验数据
	panic("implement me")
}
