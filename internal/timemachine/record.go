package timemachine

import (
	"io"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stealthrocket/timecraft/format/logsegment"
)

// Record is a read-only record from the log.
type Record struct {
	Offset       int64
	Time         time.Time
	FunctionID   int
	FunctionCall []byte
}

// RecordBuilder is a builder for records.
type RecordBuilder struct {
	startTime    time.Time
	builder      *flatbuffers.Builder
	timestamp    int64
	functionID   uint32
	functionCall []byte
	finished     bool
}

// Reset resets the builder.
func (b *RecordBuilder) Reset(startTime time.Time) {
	b.startTime = startTime
	if b.builder == nil {
		b.builder = flatbuffers.NewBuilder(defaultBufferSize)
	} else {
		b.builder.Reset()
	}
	b.timestamp = 0
	b.functionID = 0
	b.functionCall = nil
	b.finished = false
}

// SetTimestamp sets the timestamp.
func (b *RecordBuilder) SetTimestamp(t time.Time) {
	if b.finished {
		panic("builder must be reset before timestamp can be set")
	}
	b.timestamp = int64(t.Sub(b.startTime))
}

// SetFunctionID sets the function ID.
func (b *RecordBuilder) SetFunctionID(id int) {
	if b.finished {
		panic("builder must be reset before function ID can be set")
	}
	b.functionID = uint32(id)
}

// SetFunctionCall sets the function call.
//
// The provided slice is retained until Bytes() is called and the record is
// serialized.
func (b *RecordBuilder) SetFunctionCall(functionCall []byte) {
	if b.finished {
		panic("builder must be reset before function call can be set")
	}
	b.functionCall = functionCall
}

// Bytes returns the serialized representation of the record.
func (b *RecordBuilder) Bytes() []byte {
	if !b.finished {
		b.build()
		b.finished = true
	}
	return b.builder.FinishedBytes()
}

// Write writes the serialized representation of the record
// to the specified writer.
func (b *RecordBuilder) Write(w io.Writer) (int, error) {
	return w.Write(b.Bytes())
}

func (b *RecordBuilder) build() {
	if b.builder == nil {
		b.builder = flatbuffers.NewBuilder(defaultBufferSize)
	}
	functionCall := b.builder.CreateByteVector(b.functionCall)
	logsegment.RecordStart(b.builder)
	logsegment.RecordAddTimestamp(b.builder, b.timestamp)
	logsegment.RecordAddFunctionId(b.builder, b.functionID)
	logsegment.RecordAddFunctionCall(b.builder, functionCall)
	b.builder.FinishSizePrefixed(logsegment.RecordEnd(b.builder))
}
