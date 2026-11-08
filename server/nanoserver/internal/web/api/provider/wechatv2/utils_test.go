package wechatv2

// import (
// 	"testing"
// )

// func TestGenerateHMACSignature(t *testing.T) {
// 	jsonData := `{"kv_list":[{"key":"score","value":"100"},{"key":"gold","value":"3000"}]}`
// 	sessionKey := "9hAb/NEYUlkaMBEsmFgzig=="
// 	signature := GenerateHMACSignature(sessionKey, jsonData)
// 	except := "282c0d87a489e3cf6cb58d93c747716563ed67fb13ba38c2433327097adbe02c"

// 	if signature != except {
// 		t.Errorf("expected signature %s, got %s", except, signature)
// 	}
// }

// func TestJSONToMap(t *testing.T) {
// 	// Test valid JSON
// 	jsonStr := `{"name":"John","age":30,"city":"New York"}`
// 	result, err := JSONToMap(jsonStr)
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}

// 	if result["name"] != "John" {
// 		t.Errorf("expected name to be John, got %v", result["name"])
// 	}
// 	if result["age"] != float64(30) {
// 		t.Errorf("expected age to be 30, got %v", result["age"])
// 	}
// 	if result["city"] != "New York" {
// 		t.Errorf("expected city to be New York, got %v", result["city"])
// 	}

// 	// Test invalid JSON
// 	invalidJSON := `{"name":"John","age":}`
// 	_, err = JSONToMap(invalidJSON)
// 	if err == nil {
// 		t.Error("expected error for invalid JSON, got nil")
// 	}

// 	// Test empty JSON object
// 	emptyJSON := `{}`
// 	result, err = JSONToMap(emptyJSON)
// 	if err != nil {
// 		t.Errorf("expected no error for empty JSON, got %v", err)
// 	}
// 	if len(result) != 0 {
// 		t.Errorf("expected empty map, got %v", result)
// 	}
// }

// func TestMapToJSON(t *testing.T) {
// 	// Test valid map
// 	data := map[string]interface{}{
// 		"name": "John",
// 		"age":  30,
// 		"city": "New York",
// 	}
// 	result, err := MapToJSON(data)
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}

// 	// Parse back to verify
// 	parsed, err := JSONToMap(result)
// 	if err != nil {
// 		t.Errorf("failed to parse result back to map: %v", err)
// 	}
// 	if parsed["name"] != "John" {
// 		t.Errorf("expected name to be John, got %v", parsed["name"])
// 	}
// 	if parsed["age"] != float64(30) {
// 		t.Errorf("expected age to be 30, got %v", parsed["age"])
// 	}
// 	if parsed["city"] != "New York" {
// 		t.Errorf("expected city to be New York, got %v", parsed["city"])
// 	}

// 	// Test empty map
// 	emptyData := map[string]interface{}{}
// 	result, err = MapToJSON(emptyData)
// 	if err != nil {
// 		t.Errorf("expected no error for empty map, got %v", err)
// 	}
// 	if result != "{}" {
// 		t.Errorf("expected empty JSON object, got %s", result)
// 	}

// 	// Test map with nil value
// 	nilData := map[string]interface{}{
// 		"name": "John",
// 		"data": nil,
// 	}
// 	_, err = MapToJSON(nilData)
// 	if err != nil {
// 		t.Errorf("expected no error for map with nil value, got %v", err)
// 	}

// 	// Test map with nested data
// 	nestedData := map[string]interface{}{
// 		"user": map[string]interface{}{
// 			"name": "John",
// 			"age":  30,
// 		},
// 		"items": []string{"apple", "banana"},
// 	}
// 	_, err = MapToJSON(nestedData)
// 	if err != nil {
// 		t.Errorf("expected no error for nested data, got %v", err)
// 	}
// }
