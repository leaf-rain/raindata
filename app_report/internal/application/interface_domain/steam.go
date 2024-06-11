package interface_domain

import "context"

type InterfaceEventManager interface {
	StorageEvent(ctx context.Context, msg string) error
}

type InterfaceWriter interface {
	WriterMsg(ctx context.Context, msg string) error
}
