/*
必须注意: StringBuffer使用了unsafe避免string内存复制,但是其String()必须在局部范围使用即在GetStringBuffer()与PutStringBuffer()之间. 否则请使用bytes.Buffer
*/
package kit

import (
	"unsafe"
)

type StringBuffer struct {
	len int
	cap int
	buf []byte
}

func newStringBuffer(cap int) *StringBuffer {
	return &StringBuffer{
		len: 0,
		cap: cap,
		buf: make([]byte, cap),
	}
}

func (sb *StringBuffer) String() string {
	return string(sb.buf[:sb.len])
}

// 减少内存复制,有效提高性能,但必须确保返回结果只能在PutStringBuffer()前使用,而且结果是可变的.
func (sb *StringBuffer) UnsafeString() string {
	buf := sb.buf[:sb.len]
	return *(*string)(unsafe.Pointer(&buf))
}
func (sb *StringBuffer) Len() int { return sb.len }
func (sb *StringBuffer) Cap() int { return sb.cap }
func (sb *StringBuffer) Reset() {
	sb.len = 0
}
func (sb *StringBuffer) Grow(n int) {
	sb.cap += n
	buf := make([]byte, sb.cap)
	if sb.len > 0 {
		copy(buf, sb.buf)
	}
	sb.buf = buf
}
func (sb *StringBuffer) Write(p []byte) (int, error) {
	pln := len(p)
	sln := sb.len
	sb.len += pln
	if sb.len > sb.cap {
		buf := make([]byte, sb.len)
		if sln > 0 {
			copy(buf, sb.buf)
		}
		sb.buf = buf
		sb.cap = sb.len
	}
	copy(sb.buf[sln:sb.len], p)
	return pln, nil
}
func (sb *StringBuffer) WriteByte(c byte) error {
	sln := sb.len
	sb.len++
	if sb.len > sb.cap {
		buf := make([]byte, sb.len)
		if sln > 0 {
			copy(buf, sb.buf)
		}
		sb.buf = buf
		sb.cap = sb.len
	}
	sb.buf[sln] = c
	return nil
}

func (sb *StringBuffer) WriteString(p string) (int, error) {
	pln := len(p)
	sln := sb.len
	sb.len += pln
	if sb.len > sb.cap {
		buf := make([]byte, sb.len)
		if sln > 0 {
			copy(buf, sb.buf)
		}
		sb.buf = buf
		sb.cap = sb.len
	}
	copy(sb.buf[sln:sb.len], p)
	return pln, nil
}
