package interface_repo

import "context"

type DistributedLock interface {
	TryLock(ctx context.Context) (bool, error)
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}
