package interface_repo

import "context"

type InterfaceWriterRepo interface {
	WriterMsg(ctx context.Context, event, msg string) error
}
