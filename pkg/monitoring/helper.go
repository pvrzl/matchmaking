package monitoring

import "context"

type Helper struct{}

func (b Helper) Monitor(ctx context.Context, name string) (context.Context, Segment) {
	tx := defaultAPM.GetTX(ctx)
	if tx == nil {
		ctx, tx = defaultAPM.StartTX(ctx, name)
	}
	return tx.StartSegment(ctx, name)
}
