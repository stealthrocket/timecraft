// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package logindex

import (
	flatbuffers "github.com/google/flatbuffers/go"

	types "github.com/stealthrocket/timecraft/pkg/format/types"
)

type RecordIndex struct {
	_tab flatbuffers.Table
}

func GetRootAsRecordIndex(buf []byte, offset flatbuffers.UOffsetT) *RecordIndex {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RecordIndex{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsRecordIndex(buf []byte, offset flatbuffers.UOffsetT) *RecordIndex {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RecordIndex{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *RecordIndex) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RecordIndex) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RecordIndex) ProcessId(obj *types.Hash) *types.Hash {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(types.Hash)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *RecordIndex) Segment() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordIndex) MutateSegment(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *RecordIndex) Keys(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *RecordIndex) KeysLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RecordIndex) MutateKeys(j int, n uint64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func (rcv *RecordIndex) Values(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *RecordIndex) ValuesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RecordIndex) MutateValues(j int, n uint64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func RecordIndexStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func RecordIndexAddProcessId(builder *flatbuffers.Builder, processId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(processId), 0)
}
func RecordIndexAddSegment(builder *flatbuffers.Builder, segment uint32) {
	builder.PrependUint32Slot(1, segment, 0)
}
func RecordIndexAddKeys(builder *flatbuffers.Builder, keys flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(keys), 0)
}
func RecordIndexStartKeysVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func RecordIndexAddValues(builder *flatbuffers.Builder, values flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(values), 0)
}
func RecordIndexStartValuesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func RecordIndexEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
