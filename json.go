package kit

import "encoding/json"

func Json(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}
