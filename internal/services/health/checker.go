package health

import "context"

type Checker interface {
	Check(ctx context.Context) error
}

type CheckFunc func(ctx context.Context) error

func (fn CheckFunc) Check(ctx context.Context) error {
	return fn(ctx)
}
