package kit

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	str := "2020-05-29 23:59:59"
	ts := ParseDateTime(str)
	tz := ts.Unix()
	fmt.Println(ts.Unix())

	ts = time.Unix(tz, 0)
	fmt.Println(time.Now().Zone())

}
