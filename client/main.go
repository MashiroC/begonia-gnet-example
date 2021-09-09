package main

import (
	"gnet-example/common"
	"math/rand"
	"sync"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	c, err := common.Dial(":12306")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			_, payload, err := c.Recv()
			if err != nil {
				panic(err)
			}
			resp, err := common.UnMarshalResponse(payload)
			if resp.Err != "" {
				panic(*resp)
			}
		}
	}()

	Run(50, 100*10000, func() {
		req := common.NewRequest(time.Now().String(), "TEST", "FUN", []byte(RandStringRunes(rand.Intn(49))))
		err := c.Write(byte(req.Opcode()), req.Marshal())
		if err != nil {
			panic(err)
		}
	})

}

func Run(worker, nums int, fun func()) {
	wg := sync.WaitGroup{}
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < nums; i++ {
				fun()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
