package repository

import (
	"context"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	"go.uber.org/zap"
)

var _ interface_repo.InterfaceWriterRepo = (*Writer)(nil)

func NewCkWriter(logger *zap.Logger, ckCluster *clickhouse_sqlx.ClickhouseCluster) *Writer {
	return &Writer{
		logger:    logger,
		ckCluster: ckCluster,
	}
}

type Writer struct {
	logger    *zap.Logger
	ckCluster *clickhouse_sqlx.ClickhouseCluster
}

func (w Writer) WriterMsg(ctx context.Context, event, msg string) error {
	//TODO implement me
	panic("implement me")
}
