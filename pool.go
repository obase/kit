package kit

import (
	"bytes"
	"github.com/obase/conf"
	"sync"
)

const POOL_CKEY = "pool"

type PoolConfig struct {
	BytesBufferInitSize  int `json:"bytesBufferInitSize" yaml:"bytesBufferInitSize"`
	StringBufferInitSize int `json:"stringBufferInitSize" yaml:"stringBufferInitSize"`
	BlockBufferInitSize  int `json:"blockBufferInitSize" yaml:"blockBufferInitSize"`
}

var (
	bytesBufferPool  sync.Pool
	stringBufferPool sync.Pool
	blockBufferPool  sync.Pool
)

func SetupPool(c *PoolConfig) {
	if c == nil {
		c = new(PoolConfig)
	}

	if c.BytesBufferInitSize <= 0 {
		c.BytesBufferInitSize = 1024 // 默认1K
	}
	if c.StringBufferInitSize > 0 {
		c.StringBufferInitSize = 1024 //默认1K
	}
	if c.BlockBufferInitSize > 0 {
		c.BlockBufferInitSize = 32 * 1024 //默认32K
	}

	bytesBufferPool = sync.Pool{
		New: func(size int) func() interface{} {
			return func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, size))
			}
		}(c.BytesBufferInitSize),
	}
	stringBufferPool = sync.Pool{
		New: func(size int) func() interface{} {
			return func() interface{} {
				return newStringBuffer(size)
			}
		}(c.StringBufferInitSize),
	}
	blockBufferPool = sync.Pool{
		New: func(size int) func() interface{} {
			return func() interface{} {
				return make([]byte, size)
			}
		}(c.BlockBufferInitSize),
	}
}

func GetBytesBufferN(n int) (ret *bytes.Buffer) {
	ret = bytesBufferPool.Get().(*bytes.Buffer)
	ret.Reset() // 需要重置
	if n -= ret.Cap(); n > 0 {
		ret.Grow(n)
	}
	return
}

func GetBytesBuffer() (ret *bytes.Buffer) {
	ret = bytesBufferPool.Get().(*bytes.Buffer)
	ret.Reset() // 需要重置
	return
}

func PutBytesBuffer(buf *bytes.Buffer) {
	bytesBufferPool.Put(buf)
}

func GetStringBufferN(n int) (ret *StringBuffer) {
	ret = stringBufferPool.Get().(*StringBuffer)
	ret.Reset()
	if n -= ret.Cap(); n > 0 {
		ret.Grow(n)
	}
	return
}

func GetStringBuffer() (ret *StringBuffer) {
	ret = stringBufferPool.Get().(*StringBuffer)
	ret.Reset()
	return
}

func PutStringBuffer(buf *StringBuffer) {
	stringBufferPool.Put(buf)
}

func GetBlockBufferN(n int) (ret []byte) {
	ret = blockBufferPool.Get().([]byte)
	if n-len(ret) > 0 {
		ret = make([]byte, n)
	}
	return
}

func GetBlockBuffer() (ret []byte) {
	ret = blockBufferPool.Get().([]byte)
	return
}

func PutBlockBuffer(buf []byte) {
	blockBufferPool.Put(buf)
}

func init() {
	var config *PoolConfig
	conf.Bind(POOL_CKEY, &config)
	SetupPool(config)
}
