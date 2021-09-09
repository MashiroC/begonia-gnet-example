package common

import (
	"errors"
	"fmt"
)

const (
	// frame中payload部分，string字段的分隔符
	breakByte = 0x00

	// CtrlDefaultCode 默认的ctrl code
	CtrlDefaultCode = 0

	// PingPongCtrlCode ping-pong的ctrl code
	PingPongCtrlCode = 7 // 0b0111

	// CtrlConnCode 将连接升级为连接池的code
	CtrlConnCode = 1

	// ProtocolVersion 默认的版本
	ProtocolVersion = 0
)

// makeOpcode 使用默认字段构建opcode
func makeOpcode(typCode int) int {
	dispatchCode := CtrlDefaultCode // 0 ~ 7

	version := ProtocolVersion // 0 ~ 15

	return ((typCode<<3)|dispatchCode)<<4 | version
}

func findPosInBytes(data []byte, start int) (pos int) {
	for i := start; i < len(data); i++ {
		if data[i] == breakByte {
			return i
		}
	}

	return -1
}

func findInBytes(data []byte, pos int) (res string, endPos int, err error) {
	tmpPos := findPosInBytes(data, pos+1)

	if tmpPos == -1 {
		err = errors.New(fmt.Sprint("frame unmarshal error: ", data))
		return
	}

	res = string(data[pos+1 : tmpPos])
	endPos = tmpPos
	return
}