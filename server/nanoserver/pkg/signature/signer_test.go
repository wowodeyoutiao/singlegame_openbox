package signature

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestSigner_Sign(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	method := "POST"
	path := "/api/v1/users"
	queryParams := url.Values{}
	queryParams.Set("name", "alice")
	queryParams.Add("tags", "admin")
	queryParams.Add("tags", "user")

	body := []byte(`{"email":"alice@example.com"}`)

	result := signer.Sign(method, path, queryParams, body)
	if result.Signature == "" {
		t.Error("expected signature to be non-empty")
	}
	if result.Timestamp == "" {
		t.Error("expected timestamp to be non-empty")
	}
	if result.Nonce == "" {
		t.Error("expected nonce to be non-empty")
	}
}

func TestSigner_Verify_Success(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	method := "GET"
	path := "/api/v1/users"
	queryParams := url.Values{}
	queryParams.Set("name", "alice")

	body := []byte{}

	result := signer.Sign(method, path, queryParams, body)

	valid, err := signer.Verify(method, path, result.Signature, result.Timestamp, result.Nonce, queryParams, body)
	if !valid {
		t.Errorf("expected verification to succeed, got error: %v", err)
	}
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestSigner_Verify_InvalidSignature(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	method := "GET"
	path := "/api/v1/users"
	queryParams := url.Values{}
	queryParams.Set("name", "alice")

	body := []byte{}

	result := signer.Sign(method, path, queryParams, body)

	// 修改签名
	invalidSignature := "invalid-signature"

	valid, err := signer.Verify(method, path, invalidSignature, result.Timestamp, result.Nonce, queryParams, body)
	if valid {
		t.Error("expected verification to fail due to invalid signature")
	}
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSigner_Verify_MissingParameters(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	valid, err := signer.Verify("GET", "/path", "", "1234567890", "nonce123", nil, nil)
	if valid {
		t.Error("expected verification to fail due to missing signature")
	}
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "missing signature parameters" {
		t.Errorf("expected 'missing signature parameters', got: %v", err)
	}
}

func TestSigner_Verify_InvalidTimestamp(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	valid, err := signer.Verify("GET", "/path", "sig", "not-a-number", "nonce123", nil, nil)
	if valid {
		t.Error("expected verification to fail due to invalid timestamp")
	}
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "invalid timestamp" {
		t.Errorf("expected 'invalid timestamp', got: %v", err)
	}
}

func TestSigner_Verify_RequestExpired(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret, WithTimestampWindow(1*time.Second)) // 小窗口

	method := "GET"
	path := "/api/v1/users"
	queryParams := url.Values{}
	queryParams.Set("name", "alice")
	body := []byte{}

	result := signer.Sign(method, path, queryParams, body)

	time.Sleep(2 * time.Second) // 超时

	valid, err := signer.Verify(method, path, result.Signature, result.Timestamp, result.Nonce, queryParams, body)
	if valid {
		t.Error("expected verification to fail due to expired request")
	}
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "request expired" {
		t.Errorf("expected 'request expired', got: %v", err)
	}
}

func TestSigner_getSortedQueryString(t *testing.T) {
	signer := &Signer{}

	params := url.Values{}
	params.Add("z", "3")
	params.Add("a", "1")
	params.Add("b", "2")
	params.Add("a", "0") // 多个值

	expected := "a=0&a=1&b=2&z=3"
	actual := signer.getSortedQueryString(params)
	if actual != expected {
		t.Errorf("expected sorted query string to be %s, got %s", expected, actual)
	}
}

func TestSigner_generateNonce(t *testing.T) {
	signer := &Signer{}
	nonce := signer.generateNonce()
	if len(nonce) != 32 { // hex encoded 16 bytes = 32 chars
		t.Errorf("expected nonce length to be 32, got %d", len(nonce))
	}
}

func TestSigner_computeHMAC(t *testing.T) {
	secret := "mykey"
	signer := &Signer{secret: secret}
	data := "hello world"
	signature := signer.computeHMAC(data)

	// 使用标准库计算期望值进行验证
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expected := hex.EncodeToString(h.Sum(nil))

	if signature != expected {
		t.Errorf("expected HMAC to be %s, got %s", expected, signature)
	}

	// 验证签名长度
	if len(signature) != 64 {
		t.Errorf("expected signature length to be 64, got %d", len(signature))
	}
}

func TestSigner_SignAndVerify_WithDifferentMethods(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	testCases := []struct {
		method string
		body   []byte
	}{
		{"GET", nil},
		{"POST", []byte(`{"data":"value"}`)},
		{"PUT", []byte(`{"id":1}`)},
		{"DELETE", nil},
		{"PATCH", []byte(`{"update":"field"}`)},
		{"HEAD", nil}, // default case
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Method_%s", tc.method), func(t *testing.T) {
			path := "/test"
			queryParams := url.Values{}
			queryParams.Set("q", "search")

			result := signer.Sign(tc.method, path, queryParams, tc.body)
			valid, err := signer.Verify(tc.method, path, result.Signature, result.Timestamp, result.Nonce, queryParams, tc.body)
			if !valid {
				t.Errorf("verification failed for method %s: %v", tc.method, err)
			}
			if err != nil {
				t.Errorf("unexpected error for method %s: %v", tc.method, err)
			}
		})
	}
}

func TestSigner_Verify_WithEmptyParams(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	method := "POST"
	path := "/api/test"
	body := []byte(`{"test": "data"}`)

	result := signer.Sign(method, path, nil, body)

	// 测试 nil queryParams
	valid, err := signer.Verify(method, path, result.Signature, result.Timestamp, result.Nonce, nil, body)
	if !valid {
		t.Errorf("verification failed with nil queryParams: %v", err)
	}

	// 测试 empty queryParams
	valid, err = signer.Verify(method, path, result.Signature, result.Timestamp, result.Nonce, url.Values{}, body)
	if !valid {
		t.Errorf("verification failed with empty queryParams: %v", err)
	}
}

func TestSigner_Verify_WithSpecialCharacters(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	method := "GET"
	path := "/api/test"
	queryParams := url.Values{}
	queryParams.Set("special", "a&b=c d")
	queryParams.Set("chinese", "中文")
	queryParams.Set("emoji", "😀")

	body := []byte{}

	result := signer.Sign(method, path, queryParams, body)
	valid, err := signer.Verify(method, path, result.Signature, result.Timestamp, result.Nonce, queryParams, body)
	if !valid {
		t.Errorf("verification failed with special characters: %v", err)
	}
}

func BenchmarkSigner_Sign(b *testing.B) {
	secret := "benchmark-secret"
	signer := NewSigner(secret)
	method := "POST"
	path := "/api/v1/benchmark"
	params := url.Values{}
	params.Set("key", "value")
	body := make([]byte, 1024)
	rand.Read(body)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signer.Sign(method, path, params, body)
	}
}

func BenchmarkSigner_Verify(b *testing.B) {
	secret := "benchmark-secret"
	signer := NewSigner(secret)
	method := "POST"
	path := "/api/v1/benchmark"
	params := url.Values{}
	params.Set("key", "value")
	body := make([]byte, 1024)
	rand.Read(body)

	result := signer.Sign(method, path, params, body)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = signer.Verify(method, path, result.Signature, result.Timestamp, result.Nonce, params, body)
	}
}
