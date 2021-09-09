package main

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"gnet-example/common"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"time"
)

var count int32 = 0

func PprofWeb() {
	err := http.ListenAndServe(":9909", nil)
	if err != nil {
		panic(err)
	}
}

func main() {

	go PprofWeb()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println(count)
		}

	}()
	s := &testGNet{pool: goroutine.Default()}
	c := &common.Codec{}
	err := gnet.Serve(s, ":12306", gnet.WithMulticore(true), gnet.WithCodec(c), gnet.WithReusePort(true))
	if err != nil {
		panic(err)
	}
}

type testGNet struct {
	*gnet.EventServer
	pool *goroutine.Pool
}

func (t *testGNet) React(data []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	atomic.AddInt32(&count, 1)
	payload := data[1:]
	req, err := common.UnMarshalRequest(payload)
	if err != nil {
		panic(err)
	}
	if req.Service != "TEST" || req.Fun != "FUN" {
		panic("request frame error")
	}
	resp := common.NewResponse(req.ReqID, append([]byte("hello,"), req.Params...), nil)
	msg, err := common.BuildMessage(byte(resp.Opcode()), resp.Marshal())

	if err != nil {
		panic(err)
	}

	go func() {
		err = c.AsyncWrite(msg)
		defer common.BytesPool.Put(msg[:0])
		if err != nil {
			panic(err)
		}
		return
	}()

	//t.pool.Submit(func() {

	//})
	return
}

