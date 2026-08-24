package server

import (
	"bytes"
	"encoding/json"
)

// jsonUnmarshal 用严格模式解析 JSON：拒绝未知字段与重复键，
// 便于把拼写错误即时暴露给调用方。
func jsonUnmarshal(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
