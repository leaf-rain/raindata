package domain

import (
	"context"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"go.uber.org/zap"
)

var _ interface_domain.InterfaceWriter = (*Writer)(nil)

type Writer struct {
	domain *Domain
}

func NewCkWriter(domain *Domain) *Writer {
	return &Writer{
		domain: domain,
	}
}

func (logic *Writer) WriterMsg(ctx context.Context, appid int64, msg string) error {
	if len(msg) == 0 {
		logic.domain.logger.Error("[WriterMsg] msg is empty", zap.String("msg", msg))
		return nil
	}
	// 获取唯一id
	id, err := logic.domain.repoId.InitId(ctx, appid)
	if err != nil {
		logic.domain.logger.Error("[WriterMsg] logic.repoId.InitId failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	var entityMsg *EntityMsg
	entityMsg, err = NewEntityMsg(id, appid, msg)
	if err != nil {
		logic.domain.logger.Error("[WriterMsg] NewEntityMsg failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	// 写入数据
	err = logic.domain.repoWriterMsg.WriterMsg(ctx, appid, entityMsg.event, entityMsg.getString())
	if err != nil {
		logic.domain.logger.Error("[WriterMsg] logic.repoWriterMsg.WriterMsg failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	return nil
}
