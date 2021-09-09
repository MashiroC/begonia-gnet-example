package common

import (
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet"
)

var EOF = errors.New("eof")

type Codec struct {
}

func (c2 *Codec) Encode(c gnet.Conn, buf []byte) (res []byte, err error) {
	res = buf
	return
}

func (c2 *Codec) Decode(c gnet.Conn) (frame []byte, err error) {
	frame, err = c2.read(c, 1)
	if err != nil {
		c.ResetBuffer()
		return nil, nil
	}

	buf, err := c2.read(c, 1)
	if err != nil {
		panic(err)
	}

	length := buf[0]
	payloadLen := int(length)
	if length == baseLenMaxByte {
		var extendLen []byte
		extendLen, err = c2.read(c, 2)
		if err != nil {
			panic(err)
		}

		payloadLen = int(binary.BigEndian.Uint16(extendLen))
		// 我们不支持超过一定大小的包
		if payloadLen >= extendLengthMax {
			err = errors.New("payload len oversize")
			return
		}
	}

	payload, err := c2.read(c, payloadLen)
	if err != nil {
		panic(err)
	}
	frame = append(frame, payload...)

	return
}

func (c2 *Codec) read(c gnet.Conn, n int) ([]byte, error) {
	size, buf := c.ReadN(n)
	if size == 0 {
		return nil, EOF
	} else {
		c.ShiftN(size)
	}

	length := size

	for length < n {
		var buf2 []byte
		size, buf2 = c.ReadN(n - length)
		buf = append(buf, buf2...)
		length += size
		if size == 0 {
			panic("read error")
		}
		c.ShiftN(size)
	}

	return buf, nil
}
