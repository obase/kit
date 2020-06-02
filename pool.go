package kit

import (
	"bytes"
	"github.com/obase/conf"
	"sync"
)

const POOL_CKEY = "kit.pool"

type PoolConfig struct {
	BytesBufferInitSize  int `json:"bytesBufferInitSize" yaml:"bytesBufferInitSize"`
	StringBufferInitSize int `json:"stringBufferInitSize" yaml:"stringBufferInitSize"`
	BlockBufferInitSize  int `json:"blockBufferInitSize" yaml:"blockBufferInitSize"`
}

func init() {
	var config PoolConfig
	conf.Bind(POOL_CKEY, &config)
	if config.BytesBufferInitSize > 0 {
		bytesBufferInitSIze = config.BytesBufferInitSize
	} else {
		bytesBufferInitSIze = 1024 // 默认1K
	}
	if config.StringBufferInitSize > 0 {
		stringBufferInitSize = config.StringBufferInitSize
	} else {
		stringBufferInitSize = 1024 //默认1K
	}
	if config.BlockBufferInitSize > 0 {
		blockBufferInitSize = config.BlockBufferInitSize
	} else {
		blockBufferInitSize = 32 * 1024 //默认32K
	}
}

var (
	bytesBufferInitSIze  int // 默认1K
	stringBufferInitSize int // 默认1K
	blockBufferInitSize  int // 默认32K

	bytesBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, bytesBufferInitSIze))
		},
	}
	stringBufferPool = sync.Pool{
		New: func() interface{} {
			return newStringBuffer(stringBufferInitSize)
		},
	}
	blockBufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, blockBufferInitSize)
		},
	}
)

func GetBytesBuffer() (ret *bytes.Buffer) {
	ret = bytesBufferPool.Get().(*bytes.Buffer)
	ret.Reset() // 需要重置
	return
}

func PutBytesBuffer(buf *bytes.Buffer) {
	bytesBufferPool.Put(buf)
}

func GetStringBuffer() (ret *StringBuffer) {
	ret = stringBufferPool.Get().(*StringBuffer)
	ret.Reset()
	return
}

func PutStringBuffer(buf *StringBuffer) {
	stringBufferPool.Put(buf)
}

func GetBlockBuffer() (ret []byte) {
	ret = blockBufferPool.Get().([]byte)
	return
}

func PutBlockBuffer(buf []byte) {
	blockBufferPool.Put(buf)
}
