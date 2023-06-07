// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package logsegment

import (
	flatbuffers "github.com/google/flatbuffers/go"

	types "github.com/stealthrocket/timecraft/format/types"
)

type RecordBatch struct {
	_tab flatbuffers.Table
}

const RecordBatchIdentifier = "TL.0"

func GetRootAsRecordBatch(buf []byte, offset flatbuffers.UOffsetT) *RecordBatch {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RecordBatch{}
	x.Init(buf, n+offset)
	return x
}

func FinishRecordBatchBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	identifierBytes := []byte(RecordBatchIdentifier)
	builder.FinishWithFileIdentifier(offset, identifierBytes)
}

func RecordBatchBufferHasIdentifier(buf []byte) bool {
	return flatbuffers.BufferHasIdentifier(buf, RecordBatchIdentifier)
}

func GetSizePrefixedRootAsRecordBatch(buf []byte, offset flatbuffers.UOffsetT) *RecordBatch {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RecordBatch{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRecordBatchBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	identifierBytes := []byte(RecordBatchIdentifier)
	builder.FinishSizePrefixedWithFileIdentifier(offset, identifierBytes)
}

func SizePrefixedRecordBatchBufferHasIdentifier(buf []byte) bool {
	return flatbuffers.SizePrefixedBufferHasIdentifier(buf, RecordBatchIdentifier)
}

func (rcv *RecordBatch) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RecordBatch) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RecordBatch) FirstOffset() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateFirstOffset(n int64) bool {
	return rcv._tab.MutateInt64Slot(4, n)
}

func (rcv *RecordBatch) FirstTimestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateFirstTimestamp(n int64) bool {
	return rcv._tab.MutateInt64Slot(6, n)
}

func (rcv *RecordBatch) LastTimestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateLastTimestamp(n int64) bool {
	return rcv._tab.MutateInt64Slot(8, n)
}

func (rcv *RecordBatch) CompressedSize() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateCompressedSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func (rcv *RecordBatch) UncompressedSize() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateUncompressedSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(12, n)
}

func (rcv *RecordBatch) NumRecords() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateNumRecords(n uint32) bool {
	return rcv._tab.MutateUint32Slot(14, n)
}

func (rcv *RecordBatch) Checksum() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateChecksum(n uint32) bool {
	return rcv._tab.MutateUint32Slot(16, n)
}

func (rcv *RecordBatch) Compression() types.Compression {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return types.Compression(rcv._tab.GetUint32(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *RecordBatch) MutateCompression(n types.Compression) bool {
	return rcv._tab.MutateUint32Slot(18, uint32(n))
}

func RecordBatchStart(builder *flatbuffers.Builder) {
	builder.StartObject(8)
}
func RecordBatchAddFirstOffset(builder *flatbuffers.Builder, firstOffset int64) {
	builder.PrependInt64Slot(0, firstOffset, 0)
}
func RecordBatchAddFirstTimestamp(builder *flatbuffers.Builder, firstTimestamp int64) {
	builder.PrependInt64Slot(1, firstTimestamp, 0)
}
func RecordBatchAddLastTimestamp(builder *flatbuffers.Builder, lastTimestamp int64) {
	builder.PrependInt64Slot(2, lastTimestamp, 0)
}
func RecordBatchAddCompressedSize(builder *flatbuffers.Builder, compressedSize uint32) {
	builder.PrependUint32Slot(3, compressedSize, 0)
}
func RecordBatchAddUncompressedSize(builder *flatbuffers.Builder, uncompressedSize uint32) {
	builder.PrependUint32Slot(4, uncompressedSize, 0)
}
func RecordBatchAddNumRecords(builder *flatbuffers.Builder, numRecords uint32) {
	builder.PrependUint32Slot(5, numRecords, 0)
}
func RecordBatchAddChecksum(builder *flatbuffers.Builder, checksum uint32) {
	builder.PrependUint32Slot(6, checksum, 0)
}
func RecordBatchAddCompression(builder *flatbuffers.Builder, compression types.Compression) {
	builder.PrependUint32Slot(7, uint32(compression), 0)
}
func RecordBatchEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type Record struct {
	_tab flatbuffers.Table
}

func GetRootAsRecord(buf []byte, offset flatbuffers.UOffsetT) *Record {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Record{}
	x.Init(buf, n+offset)
	return x
}

func FinishRecordBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRecord(buf []byte, offset flatbuffers.UOffsetT) *Record {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Record{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRecordBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Record) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Record) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Record) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Record) MutateTimestamp(n int64) bool {
	return rcv._tab.MutateInt64Slot(4, n)
}

func (rcv *Record) FunctionId() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Record) MutateFunctionId(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *Record) FunctionCall(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *Record) FunctionCallLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Record) FunctionCallBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Record) MutateFunctionCall(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func RecordStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func RecordAddTimestamp(builder *flatbuffers.Builder, timestamp int64) {
	builder.PrependInt64Slot(0, timestamp, 0)
}
func RecordAddFunctionId(builder *flatbuffers.Builder, functionId uint32) {
	builder.PrependUint32Slot(1, functionId, 0)
}
func RecordAddFunctionCall(builder *flatbuffers.Builder, functionCall flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(functionCall), 0)
}
func RecordStartFunctionCallVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func RecordEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
