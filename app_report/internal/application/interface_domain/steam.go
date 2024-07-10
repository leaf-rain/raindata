package interface_domain

import "context"

type InterfaceWriter interface {
	WriterMsg(ctx context.Context, appid int64, msg string) error
}
