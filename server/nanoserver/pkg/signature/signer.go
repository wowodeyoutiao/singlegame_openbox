package signature

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Signer 签名器
type Signer struct {
	secret          string
	timestampWindow time.Duration
}

// Option 选项函数
type Option func(*Signer)

// WithTimestampWindow 设置时间戳窗口期，默认5分钟
func WithTimestampWindow(window time.Duration) Option {
	return func(s *Signer) {
		s.timestampWindow = window
	}
}

// NewSigner 创建签名器，secret必传，其他为选项
func NewSigner(secret string, opts ...Option) *Signer {
	s := &Signer{
		secret:          secret,
		timestampWindow: 5 * time.Minute, // 默认5分钟
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SignResult 签名结果
type SignResult struct {
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
}

// Sign 生成签名
// method: HTTP方法 (GET/POST/PUT/DELETE等)
// path: 请求路径
// queryParams: 查询参数 (GET/DELETE使用)
// body: 请求体 (POST/PUT/PATCH使用)
func (s *Signer) Sign(method, path string, queryParams url.Values, body []byte) *SignResult {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := s.generateNonce()

	// 构造签名原文
	var raw string
	switch method {
	case "GET", "DELETE":
		raw = timestamp + nonce + path + s.getSortedQueryString(queryParams)
	case "POST", "PUT", "PATCH":
		raw = timestamp + nonce + path + string(body)
	default:
		raw = timestamp + nonce + path
	}

	// 计算签名
	signature := s.computeHMAC(raw)

	return &SignResult{
		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
	}
}

// Verify 验证签名
// method: HTTP方法
// path: 请求路径
// signature: 签名值
// timestamp: 时间戳
// nonce: 随机数
// queryParams: 查询参数 (GET/DELETE使用)
// body: 请求体 (POST/PUT/PATCH使用)
// 返回 (valid bool, err error)
func (s *Signer) Verify(method, path, signature, timestamp, nonce string, queryParams url.Values, body []byte) (bool, error) {
	// 1. 检查必要参数
	if signature == "" || timestamp == "" || nonce == "" {
		return false, errors.New("missing signature parameters")
	}

	// 2. 验证时间戳
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, errors.New("invalid timestamp")
	}

	windowSeconds := int64(s.timestampWindow.Seconds())
	if diff := time.Now().Unix() - ts; diff < -windowSeconds || diff > windowSeconds {
		return false, errors.New("request expired")
	}

	// 3. 构造签名原文
	var raw string
	switch method {
	case "GET", "DELETE":
		raw = timestamp + nonce + path + s.getSortedQueryString(queryParams)
	case "POST", "PUT", "PATCH":
		raw = timestamp + nonce + path + string(body)
	default:
		raw = timestamp + nonce + path
	}

	// 4. 计算期望签名
	expectedSign := s.computeHMAC(raw)

	// 5. 安全比较签名
	if !hmac.Equal([]byte(signature), []byte(expectedSign)) {
		return false, fmt.Errorf("signature mismatch: expected %s, got %s", expectedSign, signature)
	}

	return true, nil
}

// generateNonce 生成随机nonce
func (s *Signer) generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// computeHMAC 计算HMAC-SHA256签名
func (s *Signer) computeHMAC(data string) string {
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// getSortedQueryString 对查询参数进行排序并返回标准化字符串
func (s *Signer) getSortedQueryString(params url.Values) string {
	if len(params) == 0 {
		return ""
	}

	// 收集所有键并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构造键值对
	pairs := make([]string, 0)
	for _, k := range keys {
		// 获取并排序该键的所有值
		values := params[k]
		sort.Strings(values)
		// 为每个值创建键值对
		for _, v := range values {
			encodedKey := url.QueryEscape(k)
			encodedValue := url.QueryEscape(v)
			pairs = append(pairs, encodedKey+"="+encodedValue)
		}
	}

	return strings.Join(pairs, "&")
}
