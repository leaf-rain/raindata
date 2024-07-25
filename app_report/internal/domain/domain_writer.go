package domain

import (
	"context"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/consts"
	"github.com/tidwall/gjson"
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

func (logic *Writer) WriterMsg(ctx context.Context, msg string) error {
	if len(msg) == 0 {
		logic.domain.logger.Error("[WriterMsg] msg is empty", zap.String("msg", msg))
		return nil
	}
	var appid = gjson.Get(msg, consts.KeyAppidForMsg).Int()
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
	// 获取元数据
	keys := entityMsg.getKeys()
	var newKeys map[string]string
	newKeys, err = logic.domain.repoMetadata.MetadataPut(keys)
	if err != nil || len(newKeys) != len(keys) {
		logic.domain.logger.Error("[WriterMsg] logic.repoMetadata.MetadataPut failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	if !entityMsg.putFields(newKeys) {
		logic.domain.logger.Error("[WriterMsg] entityMsg.putFields failed", zap.String("msg", msg))
		return err
	}
	// 写入数据
	err = logic.domain.repoWriterMsg.WriterMsg(ctx, appid, entityMsg.event, entityMsg.getString())
	if err != nil {
		logic.domain.logger.Error("[WriterMsg] logic.repoWriterMsg.WriterMsg failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	entityMsg.put()
	return nil
}
