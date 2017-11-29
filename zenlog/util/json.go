package util

import "encoding/json"

func MustMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	Check(err, "json.Marshal failed")
	return string(data)
}

func MustUnmarshal(data string, v interface{}) {
	err := json.Unmarshal([]byte(data), v)
	Check(err, "json.Unmarshal failed")
}

func TryUnmarshal(data string, v interface{}) bool {
	err := json.Unmarshal([]byte(data), v)
	Warn(err, "json.Unmarshal failed")
	return err != nil
}
