package base

import (
	"encoding/json"
	"strconv"
)

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

func JsonMsg(code int, data ...interface{}) []byte {
	msg := GetMsg(code)
	m := map[string]string{
		"code":    strconv.Itoa(code),
		"message": msg,
	}

	if jData, err := json.Marshal(m); err == nil {
		return jData
	}

	return nil
}
