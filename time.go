package kit

import (
	"time"
)

const (
	LAYOUT_DATE            = "2006-01-02"
	LAYOUT_TIME            = "15:04:05"
	LAYOUT_DATETIME        = "2006-01-02 15:04:05"
	LAYOUT_DATETIME_LENGTH = len(LAYOUT_DATETIME)
)

func FormatDate(t time.Time) string {
	return t.Format(LAYOUT_DATE)
}

func FormatTime(t time.Time) string {
	return t.Format(LAYOUT_TIME)
}

func FormatDateTime(t time.Time) string {
	return t.Format(LAYOUT_DATETIME)
}

func ParseDate(v string) (ret time.Time) {
	ret, _ = time.ParseInLocation(LAYOUT_DATE, v, time.Local)
	return
}

func ParseTime(v string) (ret time.Time) {
	ret, _ = time.ParseInLocation(LAYOUT_TIME, v, time.Local)
	return
}

func ParseDateTime(v string) (ret time.Time) {
	ret, _ = time.ParseInLocation(LAYOUT_DATETIME, v, time.Local)
	return
}

func ParseDateTimeExt(v string) (ret time.Time) {
	vln := len(v)
	if vln == LAYOUT_DATETIME_LENGTH {
		ret, _ = time.ParseInLocation(LAYOUT_DATETIME, v, time.Local)
	} else if vln > LAYOUT_DATETIME_LENGTH {
		ret, _ = time.ParseInLocation(LAYOUT_DATETIME, v[:LAYOUT_DATETIME_LENGTH], time.Local)
	} else {
		ret, _ = time.ParseInLocation(LAYOUT_DATETIME[:vln], v, time.Local)
	}
	return
}
