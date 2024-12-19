package monitoring

import (
	"context"
	"net/http"
)

type Config struct {
	Name    string
	License string
	Enabled bool
}

type Transaction interface {
	Finish(context.Context)
	GetContext(context.Context) context.Context
	StartSegment(ctx context.Context, name string) (context.Context, Segment)
	AddAttribute(ctx context.Context, key string, value interface{})
	SetWebResponse(w http.ResponseWriter)
	SetWebRequest(req *http.Request) *http.Request
	NoticeError(error)
}

type Segment interface {
	Finish(context.Context)
	StartSegment(ctx context.Context, name string) (context.Context, Segment)
	GetContext(context.Context) context.Context
	AddAttribute(ctx context.Context, key string, value interface{})
	StartMessageProducerSegment(ctx context.Context, library string, topic string) (context.Context, Segment)
	FinishMessageProducer(ctx context.Context)
	StartExternalSegment(ctx context.Context, req *http.Request) (context.Context, Segment)
	SetExternalSegmentResponse(ctx context.Context, resp *http.Response)
	FinishExternal(ctx context.Context)
	NoticeError(error)
}

type Monitoring interface {
	HandlerWrapper(path string, handler http.HandlerFunc) (string, http.HandlerFunc)
	Middleware() func(next http.Handler) http.Handler
	StartTX(ctx context.Context, name string) (context.Context, Transaction)
	GetHTTPParentTX(ctx context.Context, httpHeader http.Header, name string) (context.Context, Transaction)
	InjectHTTPHeader(ctx context.Context, header http.Header)
	GetTraceParentKV(ctx context.Context) (key string, val string)
	GetTX(ctx context.Context) Transaction
	WithHTTPRoundTripper(ctx context.Context, client *http.Client)
	IsEnabled() bool
}

var defaultAPM Monitoring = &noopAPM{}

func SetAPM(apm Monitoring) {
	if apm.IsEnabled() {
		defaultAPM = apm
	}
}

func APM() Monitoring {
	return defaultAPM
}
