package util

import (
	"encoding/json"
	"os"
)

var (
	debugJSON = false
)

func init() {
	debugJSON = debugJSON || (os.Getenv("ZENLOG_DEBUG_JSON") == "1")
}

// MustMarshal is a must version of json.Marshal.
func MustMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	Check(err, "json.Marshal failed")
	if debugJSON {
		DebugfForce("Marshal -> \"%s\"", data)
	}
	return string(data)
}

// MustUnmarshal is a must version of json.Unmarshal.
func MustUnmarshal(data string, v interface{}) {
	if debugJSON {
		DebugfForce("Unarshal <- \"%s\"", data)
	}
	err := json.Unmarshal([]byte(data), v)
	Check(err, "json.Unmarshal failed, input=\"%s\"", data)
}

// TryUnmarshal is a json.Unmarshal wrapper that returns whether succeeded or not.
func TryUnmarshal(data string, v interface{}) bool {
	if debugJSON {
		DebugfForce("Unarshal <- \"%s\"", data)
	}
	err := json.Unmarshal([]byte(data), v)
	Warn(err, "json.Unmarshal failed, input=\"%s\"", data)
	return err == nil
}
