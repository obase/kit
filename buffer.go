package kit

import (
	"bytes"
	"github.com/obase/conf"
	"sync"
)

const BUFFER_CKEY = "kitset.buffer"

var (
	defaultBufferBool sync.Pool
	defaultBytesPool  sync.Pool
)

func init() {
	var config *BufferConfig
	conf.Bind(BUFFER_CKEY, &config)
	SetupBuffer(config)
}

type BufferConfig struct {
	BufferInitSize int `json:"bufferInitSize" bson:"bufferInitSize"`
	BytesInitSize  int `json:"bytesInitSize" bson:"bytesInitSize"`
}

func mergeBufferConfig(c *BufferConfig) *BufferConfig {
	if c == nil {
		c = new(BufferConfig)
	}
	return c
}

func SetupBuffer(c *BufferConfig) {
	c = mergeBufferConfig(c)
	defaultBufferBool = sync.Pool{
		New: func(size int) func() interface{} {
			if size > 0 {
				return func() interface{} {
					return bytes.NewBuffer(make([]byte, 0, size))
				}
			} else {
				return func() interface{} {
					return new(bytes.Buffer)
				}
			}
		}(c.BufferInitSize),
	}
	defaultBytesPool = sync.Pool{
		New: func(size int) func() interface{} {
			if size > 0 {
				return func() interface{} {
					return make([]byte, size)
				}
			} else {
				return func() interface{} {
					return make([]byte, 32*1024)
				}
			}
		}(c.BytesInitSize),
	}
}

func BorrowBuffer() (ret *bytes.Buffer) {
	ret = defaultBufferBool.Get().(*bytes.Buffer)
	ret.Reset() // 需要重置
	return
}

func ReturnBuffer(buf *bytes.Buffer) {
	defaultBufferBool.Put(buf)
}

func BorrowBytes() (ret []byte) {
	ret = defaultBytesPool.Get().([]byte) // 不会重置
	return
}

func ReturnBytes(v []byte) {
	defaultBytesPool.Put(v)
}

type httpBufferPool struct {
	Proxy *sync.Pool
}

func (s *httpBufferPool) Get() []byte {
	return s.Proxy.Get().([]byte)
}
func (s *httpBufferPool) Put(v []byte) {
	s.Proxy.Put(v)
}
