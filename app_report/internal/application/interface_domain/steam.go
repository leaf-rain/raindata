package interface_domain

import "context"

type InterfaceWriter interface {
	WriterMsg(ctx context.Context, msg string) error
}
