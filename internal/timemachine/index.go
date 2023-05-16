package timemachine

import (
	"io"
	"sort"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stealthrocket/timecraft/format/logindex"
)

type RecordIndex struct {
	keys   []uint64
	values []uint64
}

func NewRecordIndex(data []byte) *RecordIndex {
	table := flatbuffers.Table{Bytes: data}
	return &RecordIndex{
		keys:   uint64Vector(table, 8),
		values: uint64Vector(table, 10),
	}
}

func ReadRecordIndex(input io.ReaderAt, size int64) (*RecordIndex, error) {
	b, err := readAll(input, size)
	if err != nil {
		return nil, err
	}
	return NewRecordIndex(b), nil
}

type RecordIndexBuilder struct {
	index RecordIndex
}

func (b *RecordIndexBuilder) Push(key, value uint64) {
	b.index.keys = append(b.index.keys, key)
	b.index.values = append(b.index.values, value)
}

func (b *RecordIndexBuilder) Reset() {
	b.index.keys = b.index.keys[:0]
	b.index.values = b.index.values[:0]
}

func (b *RecordIndexBuilder) RecordIndex() *RecordIndex {
	sort.Sort(recordIndexOrder{&b.index})
	return &b.index
}

type recordIndexOrder struct {
	*RecordIndex
}

func (index recordIndexOrder) Len() int {
	return len(index.keys)
}

func (index recordIndexOrder) Less(i, j int) bool {
	return index.keys[i] < index.keys[j]
}

func (index recordIndexOrder) Swap(i, j int) {
	swap(index.keys, i, j)
	swap(index.values, i, j)
}

func swap(s []uint64, i, j int) {
	s[i], s[j] = s[j], s[i]
}

type RecordIndexWriter struct {
	output  io.Writer
	builder *flatbuffers.Builder
}

func NewRecordIndexWriter(output io.Writer) *RecordIndexWriter {
	return NewRecordIndexWriterSize(output, defaultBufferSize)
}

func NewRecordIndexWriterSize(output io.Writer, bufferSize int) *RecordIndexWriter {
	return &RecordIndexWriter{
		output:  output,
		builder: flatbuffers.NewBuilder(bufferSize),
	}
}

func (w *RecordIndexWriter) Reset(output io.Writer) {
	w.output = output
	w.builder.Reset()
}

func (w *RecordIndexWriter) WriteRecordIndex(header *LogHeader, index *RecordIndex) error {
	w.builder.Reset()

	processID := header.Process.ID.prepend(w.builder)
	keys := prependUint64Vector(w.builder, index.keys)
	values := prependUint64Vector(w.builder, index.values)

	logindex.RecordIndexStart(w.builder)
	logindex.RecordIndexAddProcessId(w.builder, processID)
	logindex.RecordIndexAddSegment(w.builder, header.Segment)
	logindex.RecordIndexAddKeys(w.builder, keys)
	logindex.RecordIndexAddValues(w.builder, values)
	indexOffset := logindex.RecordIndexEnd(w.builder)

	w.builder.FinishWithFileIdentifier(indexOffset, tl2)

	_, err := w.output.Write(w.builder.FinishedBytes())
	return err
}