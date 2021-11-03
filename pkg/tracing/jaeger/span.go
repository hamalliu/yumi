package jaeger

import (
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
)

// Span implements opentracing.Span
type Span struct {
	// referenceCounter used to increase the lifetime of
	// the object before return it into the pool.
	referenceCounter int32

	serviceName string

	sync.RWMutex

	// TODO: (breaking change) change to use a pointer
	context SpanContext

	// The name of the "operation" this span is an instance of.
	// Known as a "span name" in some implementations.
	operationName string

	// firstInProcess, if true, indicates that this span is the root of the (sub)tree
	// of spans in the current process. In other words it's true for the root spans,
	// and the ingress spans when the process joins another trace.
	firstInProcess bool

	// startTime is the timestamp indicating when the span began, with microseconds precision.
	startTime time.Time

	// duration returns duration of the span with microseconds precision.
	// Zero value means duration is unknown.
	duration time.Duration

	// tags attached to this span
	tags []Tag

	// The span's "micro-log"
	logs []opentracing.LogRecord

	// The number of logs dropped because of MaxLogsPerSpan.
	numDroppedLogs int

	// references for this span
	references []Reference
}

type Reference struct {
	// TODO:
}

// Tag is a simple key value wrapper.
// TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.
type Tag struct {
	key   string
	value interface{}
}

// NewTag creates a new Tag.
// TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.
func NewTag(key string, value interface{}) Tag {
	return Tag{key: key, value: value}
}
