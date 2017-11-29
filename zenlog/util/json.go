package util

import (
	"encoding/json"
	"os"
)

var (
	debugJson = false
)

func init() {
	debugJson = debugJson || (os.Getenv("ZENLOG_DEBUG_JSON") == "1")
}

func MustMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	Check(err, "json.Marshal failed")
	if debugJson {
		DebugfForce("Marshal -> \"%s\"", data)
	}
	return string(data)
}

func MustUnmarshal(data string, v interface{}) {
	if debugJson {
		DebugfForce("Unarshal <- \"%s\"", data)
	}
	err := json.Unmarshal([]byte(data), v)
	Check(err, "json.Unmarshal failed, input=\"%s\"", data)
}

func TryUnmarshal(data string, v interface{}) bool {
	if debugJson {
		DebugfForce("Unarshal <- \"%s\"", data)
	}
	err := json.Unmarshal([]byte(data), v)
	Warn(err, "json.Unmarshal failed, input=\"%s\"", data)
	return err == nil
}
