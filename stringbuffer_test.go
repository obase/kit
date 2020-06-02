package kit

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestNewStringBuffer(t *testing.T) {
	var ls = GetStringBuffer()
	ls.WriteString("this is a")
	ls.WriteString(" another...")
	ls.Write([]byte(" and testign..."))
	var str1 string = ls.Intern()
	fmt.Println(str1)
	ls.buf[0] = 'A'
	io.Copy(ls, strings.NewReader("this is a test"))
	var str2 string = ls.Intern()

	fmt.Println(str1)
	fmt.Println(str2)
}
