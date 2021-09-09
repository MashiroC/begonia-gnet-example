package common

import (
	"errors"
	"github.com/MashiroC/begonia/tool/qconv"
)

// frame_request.go something

const (
	// request的typCode
	requestTypCode = 0
)

// Request Request的frame实现
//
// opcode4 length8 extendLength16
// req:server fun reqId param
//     4      4         8       0 || 16   [              length                  ]
// {opcode}{version}{length}{extendLength}{reqId}0x49{server}0x49{fun}0x49{param}
//
type Request struct {
	ReqID   string // 请求id
	Service string // 要调用的服务
	Fun     string // 要调用的函数
	Params  []byte // 入参

	v      []byte // 序列化后的payload，这里是一个缓存
	opcode int    // 序列化后的opcode，初始化为-1
}

// NewRequest 创建一个新的Request
func NewRequest(reqID, service, fun string, params []byte) *Request {
	return &Request{
		ReqID:   reqID,
		Service: service,
		Fun:     fun,
		Params:  params,
		opcode:  -1,
	}
}

// unMarshalRequest 根据payload去反序列化出一个request
func UnMarshalRequest(data []byte) (req *Request, err error) {

	req = &Request{}

	var pos int
	pos = -1

	req.ReqID, pos, err = findInBytes(data, pos)
	if err != nil || len(req.ReqID)==0 {
		return
	}
	if len(req.ReqID) == 0 {
		err = errors.New("frame unmarshal error: reqID len can not be zero")
		return
	}

	req.Service, pos, err = findInBytes(data, pos)
	if err != nil {
		return
	}
	if len(req.Service) == 0 {
		err = errors.New("frame unmarshal error: service len can not be zero")
		return
	}

	req.Fun, pos, err = findInBytes(data, pos)
	if err!=nil {
		return
	}
	if len(req.Fun) == 0 {
		err = errors.New("frame unmarshal error: fun len can not be zero")
		return
	}

	req.Params = data[pos+1:]

	req.v = data
	req.opcode = -1

	return
}

// Marshal 序列化payload
//
//      4      4         8       0 || 16   [              length                  ]
//	{opcode}{version}{length}{extendLength}{reqId}0x49{server}0x49{fun}0x49{param}
//
func (r *Request) Marshal() []byte {

	if r.v == nil {
		buf := make([]byte, 0, 128)

		buf = append(buf, qconv.Qs2b(r.ReqID)...)
		buf = append(buf, breakByte)

		buf = append(buf, qconv.Qs2b(r.Service)...)
		buf = append(buf, breakByte)

		buf = append(buf, qconv.Qs2b(r.Fun)...)
		buf = append(buf, breakByte)

		buf = append(buf, r.Params...)

		r.v = buf
	}

	return r.v
}

// Opcode 组装出一个opcode
func (r *Request) Opcode() int {
	if r.opcode == -1 {
		r.opcode = makeOpcode(requestTypCode)
	}

	return r.opcode
}
