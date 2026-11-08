package provider

// func TestPaySig(t *testing.T) {
// 	// 使用示例数据
// 	appKey := "appKey"
// 	sessionKey := "test_session_key"

// 	wechatPay := NewWeChatPay(appKey, sessionKey)

// 	params := map[string]interface{}{
// 		"mode":         "goods",
// 		"offerId":      "123",
// 		"buyQuantity":  1,
// 		"env":          0,
// 		"currencyType": "CNY",
// 		"platform":     "android",
// 		"zoneId":       "1",
// 		"productId":    "testproductId",
// 		"goodsPrice":   10,
// 		"outTradeNo":   "xxxxxx",
// 		"attach":       "testdata",
// 	}

// 	_, signature, err := wechatPay.PaySig(params)
// 	if err != nil {
// 		t.Fatalf("PaySig failed: %v", err)
// 	}

// 	// 期望的签名值（来自Python示例）
// 	expected := "6f9d7b0dbec4b470f1f90930a556a1ff05d2ed3e1b34b2b2bc67d5d818b70c3c"

// 	if signature != expected {
// 		t.Errorf("Expected signature %s, got %s", expected, signature)
// 	}

// 	t.Logf("Generated signature: %s", signature)
// }

// func TestWeChatPay_Signature(t *testing.T) {
// 	weChatPay := NewWeChatPay("appKey", "9hAb/NEYUlkaMBEsmFgzig==")
// 	testData := map[string]interface{}{
// 		"kv_list": []map[string]interface{}{
// 			{
// 				"key":   "score",
// 				"value": "100",
// 			},
// 			{
// 				"key":   "gold",
// 				"value": "3000",
// 			},
// 		},
// 	}

// 	signature, err := weChatPay.Signature(testData)
// 	if err != nil {
// 		t.Errorf("WeChatPay.Signature() error = %v", err)
// 		return
// 	}
// 	expected := "282c0d87a489e3cf6cb58d93c747716563ed67fb13ba38c2433327097adbe02c"
// 	if signature != expected {
// 		t.Errorf("WeChatPay.Signature() = %v, want %v", signature, expected)
// 	}
// }
