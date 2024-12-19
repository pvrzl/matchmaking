package monitoring

import (
	"context"
	"net/http"
)

type (
	noopAPM         struct{}
	noopTransaction struct{}
	noopSegment     struct{}
)

func (nr *noopAPM) HandlerWrapper(path string, handler http.HandlerFunc) (string, http.HandlerFunc) {
	return path, handler
}

func (nr *noopAPM) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
		return http.HandlerFunc(fn)
	}
}

func (nr *noopAPM) StartTX(ctx context.Context, name string) (context.Context, Transaction) {
	return ctx, &noopTransaction{}
}

func (nr *noopAPM) GetHTTPParentTX(ctx context.Context, httpHeader http.Header, name string) (context.Context, Transaction) {
	return ctx, &noopTransaction{}
}

func (nr *noopAPM) InjectHTTPHeader(ctx context.Context, header http.Header) {}

func (nr *noopAPM) GetTraceParentKV(ctx context.Context) (key string, val string) {
	return "", ""
}

func (nr *noopAPM) GetTX(ctx context.Context) Transaction {
	return &noopTransaction{}
}

func (nr *noopAPM) WithHTTPRoundTripper(ctx context.Context, client *http.Client) {

}

func (nr *noopAPM) IsEnabled() bool {
	return false
}

func (t *noopTransaction) Finish(ctx context.Context) {}

func (t *noopTransaction) GetContext(ctx context.Context) context.Context {
	return ctx
}

func (t *noopTransaction) StartSegment(ctx context.Context, name string) (context.Context, Segment) {
	return ctx, &noopSegment{}
}

func (t *noopTransaction) SetWebResponse(w http.ResponseWriter) {}

func (t *noopTransaction) SetWebRequest(req *http.Request) *http.Request {
	return req
}

func (t *noopTransaction) AddAttribute(ctx context.Context, key string, value interface{}) {}

func (t *noopTransaction) NoticeError(err error) {}

func (s *noopSegment) Finish(ctx context.Context) {}

func (s *noopSegment) StartSegment(ctx context.Context, name string) (context.Context, Segment) {
	return ctx, &noopSegment{}
}

func (s *noopSegment) GetContext(ctx context.Context) context.Context {
	return ctx
}

func (s *noopSegment) AddAttribute(ctx context.Context, key string, value interface{}) {}

func (s *noopSegment) NoticeError(err error) {}

func (s *noopSegment) StartMessageProducerSegment(ctx context.Context, library string, topic string) (context.Context, Segment) {
	return ctx, &noopSegment{}
}

func (s *noopSegment) FinishMessageProducer(ctx context.Context) {}

func (s *noopSegment) StartExternalSegment(ctx context.Context, req *http.Request) (context.Context, Segment) {
	return ctx, &noopSegment{}
}

func (s *noopSegment) SetExternalSegmentResponse(ctx context.Context, resp *http.Response) {}

func (s *noopSegment) FinishExternal(ctx context.Context) {}
