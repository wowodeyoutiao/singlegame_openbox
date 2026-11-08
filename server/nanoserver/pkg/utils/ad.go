package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type AdIdentifier string

const (
	AdGDT       AdIdentifier = "gdt"
	AdByteDance AdIdentifier = "bytedance"
)

func (a AdIdentifier) String() string {
	return string(a)
}

// GDTParams {\"wx_aid\":\"66110715840\",\"uxinfo\":\"66110715840%7Cwx05hjojixwm4z2m01%7C66110715840%7C0%7C1761722093%7C0%7C2%7C4003158499171286%7C%7CAgI7RnBBTyAy3UavV0jZsU7Aq1c0Cq35Ja4mNlpt890bCEx7%2Fhowm4rtaa3O7u0gS1E6ZzcgLqHIwckt0Q6EC6VA7nykMEgKMpVzF%2Fw%3D%7C66114124930%7C21%7C0%7CoDdoCt1MsIfYHxJQYw8jpEMWc3xI%7C3%7C5668452122\",\"gdt_vid\":\"wx05hjojixwm4z2m01\",\"weixinadinfo\":\"66110715840.wx05hjojixwm4z2m01.0.0\",\"wx_traceid\":\"wx05hjojixwm4z2m01\"}
type GDTParams struct {
	WxAid        string `json:"wx_aid"`
	GdtVid       string `json:"gdt_vid"`
	Weixinadinfo string `json:"weixinadinfo"`
	WxTraceid    string `json:"wx_traceid"`
	Uxinfo       string `json:"uxinfo,omitempty"` // 可选
}

// BytedanceParams {\"clue_token\":\"1840409168876777,1846829806532212\",\"clickid\":\"EIiCrO7g9aMDGJTIgbubjMAFILDN4IjGrOQGMAw4wbgCQiIyMDI1MTAyNDExMjYyMjdGQUVCODAwQzc3MUZFREI1REM0SMG4ApABAA\",\"promotion_id\":\"7564322620218933257\",\"advertiser_id\":\"1840409166240896\",\"req_id\":\"202510241126227FAEB800C771FEDB5DC4\",\"clue_id\":\"1846829806532212\",\"project_id\":\"7564319835098939430\",\"mid3\":\"7563548944574087209\",\"adv_id\":\"1840409166240896\",\"mid2\":\"7550548580387586094\",\"mid1\":\"0\",\"mid5\":\"7564605681439391750\",\"mid4\":\"__MID4__\"}
type BytedanceParams struct {
	Clickid      string `json:"clickid"`
	ClueToken    string `json:"clue_token"`
	PromotionId  string `json:"promotion_id"`
	AdvertiserId string `json:"advertiser_id"`
	ReqId        string `json:"req_id"`
	ClueId       string `json:"clue_id"`
	ProjectId    string `json:"project_id"`
	Mid1         string `json:"mid1"`
	Mid2         string `json:"mid2"`
	Mid3         string `json:"mid3"`
	Mid4         string `json:"mid4"`
	Mid5         string `json:"mid5"`
}

type ParseResult struct {
	Source    AdIdentifier
	GDT       *GDTParams
	ByteDance *BytedanceParams
}

// ParseAdQuery 解析广告参数 JSON 字符串，自动识别来源（GDT / ByteDance）
func ParseAdQuery(queryJSON string) (*ParseResult, error) {
	trimmed := strings.TrimSpace(queryJSON)
	if trimmed == "" || trimmed == "{}" {
		return nil, fmt.Errorf("empty input")
	}

	payload := []byte(trimmed)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w, queryJSON: %s", err, queryJSON)
	}

	isGDT := hasNonEmptyField(raw, "gdt_vid") || (hasNonEmptyField(raw, "wx_aid") || hasNonEmptyField(raw, "weixinadinfo"))
	isByteDance := hasNonEmptyField(raw, "clickid") || (hasNonEmptyField(raw, "clue_token") || hasNonEmptyField(raw, "promotion_id"))

	switch {
	case isGDT:
		var g GDTParams
		if err := json.Unmarshal(payload, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal GDT params: %w, queryJSON: %s", err, queryJSON)
		}
		return &ParseResult{Source: AdGDT, GDT: &g}, nil
	case isByteDance:
		var b BytedanceParams
		if err := json.Unmarshal(payload, &b); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ByteDance params: %w, queryJSON: %s", err, queryJSON)
		}
		return &ParseResult{Source: AdByteDance, ByteDance: &b}, nil
	default:
		return nil, fmt.Errorf("unrecognized ad parameters, queryJSON: %s", queryJSON)
	}
}

func hasNonEmptyField(raw map[string]json.RawMessage, key string) bool {
	value, ok := raw[key]
	if !ok || len(value) == 0 {
		return false
	}

	trimmed := bytes.TrimSpace(value)
	if bytes.Equal(trimmed, []byte("null")) {
		return false
	}

	if len(trimmed) >= 2 && trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"' {
		var s string
		if err := json.Unmarshal(value, &s); err != nil {
			return false
		}
		return s != ""
	}

	return true
}
