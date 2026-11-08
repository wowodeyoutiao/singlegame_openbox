package protocol

type ConversionEventRequest struct {
	EventType string  `json:"event_type"`
	Context   Context `json:"context"`
	Timestamp int64   `json:"timestamp"`
}

type Context struct {
	Ad Ad `json:"ad"`
}

type Ad struct {
	Callback string `json:"callback"`
}

type ConversionEventResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type LaunchOptionsSyncQuery struct {
	ClueToken    string `protobuf:"bytes,1,opt,name=clue_token,json=clueToken,proto3" json:"clue_token,omitempty"`           // 线索 token
	ProjectId    string `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`           // 项目 ID
	PromotionId  string `protobuf:"bytes,3,opt,name=promotion_id,json=promotionId,proto3" json:"promotion_id,omitempty"`     // 推广 ID
	Mid1         string `protobuf:"bytes,4,opt,name=mid1,proto3" json:"mid1,omitempty"`                                      // 扩展字段 mid1
	Mid2         string `protobuf:"bytes,5,opt,name=mid2,proto3" json:"mid2,omitempty"`                                      // 扩展字段 mid2
	Mid3         string `protobuf:"bytes,6,opt,name=mid3,proto3" json:"mid3,omitempty"`                                      // 扩展字段 mid3
	Mid4         string `protobuf:"bytes,7,opt,name=mid4,proto3" json:"mid4,omitempty"`                                      // 扩展字段 mid4
	Mid5         string `protobuf:"bytes,8,opt,name=mid5,proto3" json:"mid5,omitempty"`                                      // 扩展字段 mid5
	ReqId        string `protobuf:"bytes,9,opt,name=req_id,json=reqId,proto3" json:"req_id,omitempty"`                       // 请求 ID
	AdvertiserId string `protobuf:"bytes,10,opt,name=advertiser_id,json=advertiserId,proto3" json:"advertiser_id,omitempty"` // 广告主 ID
	Clickid      string `protobuf:"bytes,11,opt,name=clickid,proto3" json:"clickid,omitempty"`                               // 点击 ID
}
