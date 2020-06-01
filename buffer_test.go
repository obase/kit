package kit

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

const (
	sss = "xfoasneobfasieongasbg"
	cnt = 10
)

var (
	bbb      = []byte(sss)
	expected = strings.Repeat(sss, cnt)
)

func BenchmarkCopyPreAllocate(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		bs := make([]byte, cnt*len(sss))
		bl := 0
		for i := 0; i < cnt; i++ {
			bl += copy(bs[bl:], sss)
		}
		result = string(bs)
	}
	b.StopTimer()
	if result != expected {
		b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkAppendPreAllocate(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		data := make([]byte, 0, cnt*len(sss))
		for i := 0; i < cnt; i++ {
			data = append(data, sss...)
		}
		result = string(data)
	}
	b.StopTimer()
	if result != expected {
		b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkBufferPreAllocate(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(make([]byte, 0, cnt*len(sss)))
		for i := 0; i < cnt; i++ {
			buf.WriteString(sss)
		}
		result = buf.String()
	}
	b.StopTimer()
	if result != expected {
		b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkCopy(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		data := make([]byte, 0, 64) // same size as bootstrap array of bytes.Buffer
		for i := 0; i < cnt; i++ {
			off := len(data)
			if off+len(sss) > cap(data) {
				temp := make([]byte, 2*cap(data)+len(sss))
				copy(temp, data)
				data = temp
			}
			data = data[0 : off+len(sss)]
			copy(data[off:], sss)
		}
		result = string(data)
	}
	b.StopTimer()
	if result != expected {
		b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkAppend(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		data := make([]byte, 0, 64)
		for i := 0; i < cnt; i++ {
			data = append(data, sss...)
		}
		result = string(data)
	}
	b.StopTimer()
	if result != expected {
		b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkBufferWrite(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		for i := 0; i < cnt; i++ {
			buf.Write(bbb)
		}
		result = buf.String()
	}
	b.StopTimer()
	if result != expected {
		b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkBufferWriteString(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		for i := 0; i < cnt; i++ {
			buf.WriteString(sss)
			buf.WriteString(sss)
			buf.WriteString(strconv.Itoa(i))
		}
		result = buf.String()
	}
	b.StopTimer()
	if result != expected {
		//b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkConcat(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			str = sss + sss + strconv.Itoa(i)
		}
		result = str
	}
	b.StopTimer()
	if result != expected {
		//b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkConcatJoin(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			str = strings.Join([]string{sss, sss, strconv.Itoa(i)}, "")
		}
		result = str
	}
	b.StopTimer()
	if result != expected {
		//b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkConcat2(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			str = sss
			str = str + sss
			str = str + strconv.Itoa(i)
		}
		result = str
	}
	b.StopTimer()
	if result != expected {
		//b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkConcat3(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			str = sss
			str += sss
			str += strconv.Itoa(i)
		}
		result = str
	}
	b.StopTimer()
	if result != expected {
		//b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkConcatBuffer(b *testing.B) {
	var result string
	buf := new(bytes.Buffer)
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			buf.Reset()
			buf.WriteString(sss)
			buf.WriteString(sss)
			buf.WriteString(strconv.Itoa(i))
			str = buf.String()
		}
		result = str
	}
	b.StopTimer()
	if result != expected {
		//b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	var result string
	buf := new(strings.Builder)
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			buf.Reset()
			buf.WriteString(sss)
			buf.WriteString(sss)
			buf.WriteString(strconv.Itoa(i))
			str = buf.String()
		}
		result = str
	}
	b.StopTimer()
	if result != expected {
		//b.Errorf("unexpected result; got=%s, want=%s", string(result), expected)
	}
}

func BenchmarkConcatOperator(b *testing.B) {
	var result string
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			str = sss + sss + strconv.Itoa(i)
		}
		result = str
	}
	b.StopTimer()
	_ = result
	//fmt.Println(result)
}

func BenchmarkConcatCopy(b *testing.B) {
	var result string
	var bs []byte //= make([]byte, 1024)
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < cnt; i++ {
			bs = bs[:0]
			bs = append(bs, sss...)
			bs = append(bs, sss...)
			bs = append(bs, strconv.Itoa(i)...)
			str = *(*string)(unsafe.Pointer(&bs))
		}
		result = str
	}
	b.StopTimer()
	_ = result
	//fmt.Println(result)
}

func TestLocalBytes_String(t *testing.T) {
	var ls StringBuilder
	ls = append(ls, "this is a"...)
	ls = append(ls, " another..."...)
	fmt.Println(ls.String())
	fmt.Println(ls.String() == "this is a another...")
}
