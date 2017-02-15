package packetutil

import (
	"encoding/binary"
	"math"
)

type PacketWriter struct {
	data       []byte
	packetID   int32
	packetSize int32
}

func CreatePacketWriter(packetID int32) *PacketWriter {
	pw := new(PacketWriter)
	pw.packetID = packetID
	pw.data = make([]byte, 0)
	pw.WriteVarInt(packetID)
	return pw
}

func (pw *PacketWriter) GetPacket() []byte {
	return append(pw.getVarLong(int64(pw.packetSize)), pw.data...)
}

func (pw *PacketWriter) appendByteSlice(data []byte) {
	pw.data = append(pw.data, data...)

	pw.packetSize += int32(len(data))
}

func (pw *PacketWriter) WriteBoolean(val bool) {
	if val {
		pw.WriteUnsignedByte(0x01)
	} else {
		pw.WriteUnsignedByte(0x00)
	}
}

func (pw *PacketWriter) WriteByte(val int8) {
	pw.WriteUnsignedByte(byte(val))
}

func (pw *PacketWriter) WriteUnsignedByte(val byte) {
	pw.appendByteSlice([]byte{val})
}

func (pw *PacketWriter) WriteShort(val int16) {
	pw.WriteUnsignedShort(uint16(val))
}

func (pw *PacketWriter) WriteUnsignedShort(val uint16) {
	buff := make([]byte, 2)
	binary.BigEndian.PutUint16(buff, val)

	pw.appendByteSlice(buff)
}

func (pw *PacketWriter) WriteInt(val int32) {
	pw.writeUnsignedInt(uint32(val))
}

func (pw *PacketWriter) writeUnsignedInt(val uint32) {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, val)

	pw.appendByteSlice(buff)
}

func (pw *PacketWriter) WriteLong(val int64) {
	pw.writeUnsignedLong(uint64(val))
}

func (pw *PacketWriter) writeUnsignedLong(val uint64) {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, val)

	pw.appendByteSlice(buff)
}

func (pw *PacketWriter) WriteFloat(val float32) {
	pw.writeUnsignedInt(math.Float32bits(val))
}

func (pw *PacketWriter) WriteDouble(val float64) {
	pw.writeUnsignedLong(math.Float64bits(val))
}

func (pw *PacketWriter) WriteString(val string) {
	pw.WriteVarInt(int32(len(val)))
	pw.appendByteSlice([]byte(val))
}

func (pw *PacketWriter) WriteVarInt(val int32) {
	pw.WriteVarLong(int64(val))
}

func (pw *PacketWriter) WriteVarLong(val int64) {
	pw.appendByteSlice(pw.getVarLong(val))
}

func (pw *PacketWriter) getVarLong(val int64) []byte {
	var buff []byte
	for {
		temp := byte(val & 0x7F)
		val = int64(uint64(val) >> 7)
		if val != 0 {
			temp |= 0x80
		}
		buff = append(buff, temp)

		if val == 0 {
			break
		}
	}
	return buff
}
