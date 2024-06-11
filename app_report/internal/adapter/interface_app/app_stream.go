package interface_app

import "context"

type InterfaceStream interface {
	Stream(ctx context.Context, msg string) error
}
