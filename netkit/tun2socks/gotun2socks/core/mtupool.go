package core

import (
	"sync"
)

var mtuBufferPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, MTU)
	},
}

func GetMTU() []byte {
	return mtuBufferPool.Get().([]byte)
}

func PutMTU(buf []byte) {
	mtuBufferPool.Put(buf)
}
