package kit

import (
	"bytes"
	"github.com/obase/conf"
	"sync"
	"unsafe"
)

const BUFFER_CKEY = "kitset.buffer"

var (
	defaultBufferBool     sync.Pool
	defaultBytesPool      sync.Pool
	defaultBufferInitSize int
	defaultBytesInitSize  int
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
	defaultBufferInitSize = c.BufferInitSize
	defaultBytesInitSize = c.BytesInitSize

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
		}(defaultBufferInitSize),
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
		}(defaultBytesInitSize),
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

type StringBuilder []byte

func (sb *StringBuilder) String() string {
	return *(*string)(unsafe.Pointer(&*sb))
}

func (sb *StringBuilder) Len() int { return len(*sb) }
func (sb *StringBuilder) Cap() int { return cap(*sb) }
func (sb *StringBuilder) Reset()   { *sb = (*sb)[:0] }
func (sb *StringBuilder) grow(n int) {
	buf := make([]byte, len(*sb), 2*cap(*sb)+n)
	copy(buf, *sb)
	*sb = buf
}
func (sb *StringBuilder) copyCheck(n int) {
	if n < 0 {
		panic("StringBuilder.Grow: negative count")
	}
	if cp, ln := cap(*sb), len(*sb); cp-ln < n {
		buf := make([]byte, ln, 2*cp+n)
		copy(buf, *sb)
		*sb = buf
	}
}
func (sb *StringBuilder) Write(p []byte) (int, error) {
	sb.copyCheck()
	b.buf = append(b.buf, p...)
	return len(p), nil
}
