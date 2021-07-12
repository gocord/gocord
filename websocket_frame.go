package gocord

import (
	"encoding/binary"
	"io"
)

type Frame struct {
	op      []byte
	payload []byte
	data    []byte
}

var bits = struct {
	fin  byte
	rsv1 byte
	rsv2 byte
	rsv3 byte
	mask byte
}{
	fin:  byte(1 << 7),
	rsv1: byte(1 << 6),
	rsv2: byte(1 << 5),
	rsv3: byte(1 << 4),
	mask: byte(1 << 7),
}

func (f *Frame) Payload() []byte {
	if len(f.payload) != 0 {
		return f.payload[2:]
	}
	return f.payload
}

// returns amount of bytes for payload length
func (f *Frame) payloadLen() int {
	// second byte = mask&payload_length
	switch f.op[1] & 127 {
	// 127 = 0111 1111
	// should read last 7 bits
	case 127:
		return 8
	case 126:
		return 2
	default:
		return 0
	}
}

func (f *Frame) Len() uint64 {
	length := uint64(f.op[1] & 127)
	switch length {
	case 126:
		return uint64(binary.BigEndian.Uint16(f.op[2:]))
	case 127:
		return binary.BigEndian.Uint64(f.op[2:])
	}

	return length
}

func (f *Frame) Read(r io.Reader) (int64, error) {
	var err error
	var n, m int
	// read first two bytes, (fin, rsvs, opcode, mask, payload length)
	n, err = io.ReadFull(r, f.op[:2])

	if err == nil {
		// get amount of bytes in payload length
		m = f.payloadLen() + 2
		if m > 2 {
			n, err = io.ReadFull(r, f.op[2:m])
		}

		if err == nil {
			frameSize := f.Len()
			if frameSize > 0 {
				_n := int64(frameSize)
				if _n > 0 {

				}
			}
		}
	}
	return int64(n), nil
}
