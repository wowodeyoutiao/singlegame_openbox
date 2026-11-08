package wechatv2

// import (
// 	"sync"
// 	"time"
// )

// // MockWeChatStore 模拟存储实现，用于测试
// type MockWeChatStore struct {
// 	rwLock sync.RWMutex
// 	tokens map[string]string
// }

// func NewMockWeChatStore() *MockWeChatStore {
// 	return &MockWeChatStore{
// 		tokens: make(map[string]string),
// 	}
// }

// func (m *MockWeChatStore) SetAccessToken(PID, accessToken string, expired time.Duration) error {
// 	m.rwLock.Lock()
// 	defer m.rwLock.Unlock()
// 	m.tokens[PID] = accessToken
// 	return nil
// }

// func (m *MockWeChatStore) GetAccessToken(PID string) (string, error) {
// 	m.rwLock.RLock()
// 	defer m.rwLock.RUnlock()
// 	token, exists := m.tokens[PID]
// 	if !exists {
// 		return "", nil
// 	}
// 	return token, nil
// }

// func (m *MockWeChatStore) SetSessionKey(PID, sessionKey string, expired time.Duration) error {
// 	m.rwLock.Lock()
// 	defer m.rwLock.Unlock()
// 	m.tokens[PID] = sessionKey
// 	return nil
// }

// func (m *MockWeChatStore) GetSessionKey(PID string) (string, error) {
// 	m.rwLock.RLock()
// 	defer m.rwLock.RUnlock()
// 	sessionKey, exists := m.tokens[PID]
// 	if !exists {
// 		return "", nil
// 	}
// 	return sessionKey, nil
// }

// func TestWeChat_GameMessagePushSimulation(t *testing.T) {
// 	// 创建模拟存储
// 	mockStore := NewMockWeChatStore()

// 	// 创建微信实例
// 	wechat := NewWeChat("test_app_id", "test_offer_id", "test_app_key", "test_app_secret", mockStore)

// 	testCase := struct {
// 		name      string
// 		token     string
// 		timestamp string
// 		nonce     string
// 		signature string
// 		expected  bool
// 		desc      string
// 	}{
// 		name:      "官方示例测试用例",
// 		token:     "AAAAA",
// 		timestamp: "1714036504",
// 		nonce:     "1514711492",
// 		signature: "f464b24fc39322e44b38aa78f5edd27bd1441696",
// 		expected:  true,
// 		desc:      "根据官方文档提供的示例数据，验证签名计算的正确性",
// 	}

// 	t.Run(testCase.name, func(t *testing.T) {
// 		result := wechat.HandleWechatMessagePushVerifyURL(testCase.token, testCase.timestamp, testCase.nonce, testCase.signature)
// 		if result != testCase.expected {
// 			t.Errorf("GameMessagePushSimulation() = %v, expected %v\nDescription: %s",
// 				result, testCase.expected, testCase.desc)
// 		}
// 	})
// }
