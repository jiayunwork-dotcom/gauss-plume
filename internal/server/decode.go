package server

import (
	"fmt"
	"io"
	"net/http"
)

// decodeJSON 读取并解析请求 JSON 体。未知字段一律拒绝，
// 解析失败返回带说明的 error。
func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("读取请求体失败：%v", err)
	}
	if len(body) == 0 {
		return fmt.Errorf("请求体为空")
	}
	if err := jsonUnmarshal(body, dst); err != nil {
		return fmt.Errorf("请求 JSON 解析失败：%v", err)
	}
	return nil
}
