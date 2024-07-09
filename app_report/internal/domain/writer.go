package domain

import (
	"context"
	"github.com/leaf-rain/fastjson"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"go.uber.org/zap"
	"unicode"
)

var _ interface_domain.InterfaceWriter = (*Writer)(nil)

type Writer struct {
	logger        *zap.Logger
	repoMetadata  interface_repo.InterfaceMetadataRepo
	repoWriterMsg interface_repo.InterfaceWriterRepo
}

func NewCkWriter(logger *zap.Logger, repoMetadata interface_repo.InterfaceMetadataRepo, repoWriterMsg interface_repo.InterfaceWriterRepo) *Writer {
	return &Writer{
		logger:        logger,
		repoMetadata:  repoMetadata,
		repoWriterMsg: repoWriterMsg,
	}
}

func (logic *Writer) WriterMsg(ctx context.Context, msg string) error {
	if len(msg) == 0 {
		logic.logger.Error("[WriterMsg] msg is empty", zap.String("msg", msg))
		return nil
	}
	// 数据解析重组
	v, err := fastjson.Parse(msg)
	if err != nil || v == nil {
		logic.logger.Error("[WriterMsg] fastjson.Parse failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	var keys []string
	v.Range(func(key []byte, v *fastjson.Value) {
		keys = append(keys, string(key))
	})
	var newKeys map[string]string
	newKeys, err = logic.repoMetadata.MetadataPut(keys)
	if err != nil || len(newKeys) != len(keys) {
		logic.logger.Error("[WriterMsg] logic.repoMetadata.MetadataPut failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	firstChar := rune(0)
	for _, r := range msg {
		if !unicode.IsSpace(r) {
			firstChar = r
		}
	}
	if firstChar == 0 {
		logic.logger.Error("[WriterMsg] msg is empty", zap.String("msg", msg))
		return err
	}
	var newParse *fastjson.Value
	if firstChar == '{' {
		newParse = fastjson.MustParse("{}")
	}
	if firstChar == '[' {
		newParse = fastjson.MustParse("[]")
	}
	v.Range(func(key []byte, v *fastjson.Value) {
		newParse.Set(newKeys[string(key)], v)
	})
	// 添加功能默认字段
	event := string(v.GetStringBytes("event"))
	if event == "" {
		newParse.Set("event", fastjson.MustParse("\"log\""))
		event = "log"
	}
	// 写入数据
	err = logic.repoWriterMsg.WriterMsg(ctx, event, newParse.String())
	if err != nil {
		logic.logger.Error("[WriterMsg] logic.repoWriterMsg.WriterMsg failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	return nil
}
