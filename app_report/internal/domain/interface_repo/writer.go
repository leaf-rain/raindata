package interface_repo

import "context"

type InterfaceWriterRepo interface {
	WriterMsg(ctx context.Context, appid int64, event, msg string) error
}
