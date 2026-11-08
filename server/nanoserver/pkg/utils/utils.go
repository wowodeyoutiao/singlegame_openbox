package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
)

// ConvertJSONToMap 将 JSON 字符串安全转换为 map[string]interface{}
func ConvertJSONToMap(jsonStr string) (map[string]interface{}, error) {
	// 输入验证
	if jsonStr == "" {
		return nil, fmt.Errorf("empty JSON string")
	}
	// 创建目标 map
	data := make(map[string]interface{})

	// 安全解析 JSON
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	return data, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateNonceStr(length int) (string, error) {
	if length < 1 || length > 32 {
		return "", fmt.Errorf("length must be between 1 and 32")
	}

	b := make([]byte, length)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", err
		}
		b[i] = letterBytes[idx.Int64()]
	}

	return string(b), nil
}

func RemoveElementInPlace[T any](slice []T, index int) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return slice, fmt.Errorf("index %d out of range [0, %d)", index, len(slice))
	}

	copy(slice[index:], slice[index+1:])
	return slice[:len(slice)-1], nil
}
