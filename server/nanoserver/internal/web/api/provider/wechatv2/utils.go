package wechatv2

// import (
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// )

// // GenerateHMACSignature 生成 HMAC-SHA256 签名 用户登录态签名
// // https://developers.weixin.qq.com/minigame/dev/guide/open-ability/signature.html
// func GenerateHMACSignature(key, encodeJsonData string) string {
// 	mac := hmac.New(sha256.New, []byte(key))
// 	mac.Write([]byte(encodeJsonData))
// 	return hex.EncodeToString(mac.Sum(nil))
// }

// // JSONToMap 将JSON字符串转换为map[string]interface{}
// func JSONToMap(jsonStr string) (map[string]interface{}, error) {
// 	var data map[string]interface{}
// 	err := json.Unmarshal([]byte(jsonStr), &data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

// // MapToJSON 将map[string]interface{}转换为JSON字符串
// func MapToJSON(data map[string]interface{}) (string, error) {
// 	jsonBytes, err := json.Marshal(data)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(jsonBytes), nil
// }
