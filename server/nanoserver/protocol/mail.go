package protocol

type CreateMailRequest struct {
	MailId string //来自哪个渠道的订单
	Uid    int64  //Token
	Offset int    `json:"offset"`
	Count  int    `json:"count"`
}
type MailInfo struct {
	MailId            string       `json:"mail_id"`      //订单号
	SenderUid         int64        `json:"sender_uid"`   //发送者ID
	RecipientUid      int64        `json:"RecipientUid"` //接受者ID
	Title             string       `json:"title"`        //标题
	Desc              string       `json:"desc"`         //详情
	IsRead            int          `json:"is_read"`      //是否阅读
	IsOfficial        int          `json:"is_official"`
	IsExtract         int          `json:"is_extract"`          //是否提取
	IsAttachments     int          `json:"is_attachments"`      //是否有附件
	IsDeleteSender    int          `json:"is_delete_sender"`    //是否对于发送者删除
	IsDeleteRecipient int          `json:"is_delete_recipient"` //是否对于接受者删除
	CreatedAt         int64        `json:"created_at"`          //创建时间
	ExpiryAt          int64        `json:"expiry_at"`           //截止时间
	Attachments       []Attachment `json:"attachments"`
}
type MailAdd struct {
	Uid          int64        `json:"uid"`           //发送者id，官方的话建议-1
	RecipientUid int64        `json:"recipient_uid"` //接受者id
	Title        string       `json:"title"`         //标题
	Desc         string       `json:"desc"`          //详情
	ChannelCode  string       `json:"channel_code"`  //来源
	Attachments  []Attachment `json:"attachments"`   //附件
}

// Attachment 代表邮件附件结构
type Attachment struct {
	Id          int
	MailId      string `json:"mail_id"`
	UserPropsId int64  `json:"user_props_id"`
	Num         int64  `json:"num"`
}

// OrderListRequest 订单列表
type MailListRequest struct {
	Uid    int64  `json:"uid"` //用户id
	MailId string //来自哪个渠道的订单
	Offset int    `json:"offset"`
	Count  int    `json:"count"`
}
type MailListResponse struct {
	Code  int        `json:"code"`  //状态码
	Data  []MailInfo `json:"data"`  //渠道列表
	Total int        `json:"total"` //总数量
}

type MailDetailResponse struct {
	Code int      `json:"code"` //状态码
	Data MailInfo `json:"data"` //返回信息
}
type MailOperateResponse struct {
	Code int    `json:"code"` //状态码
	Data string `json:"data"` //返回信息
}
type MailOperateRequest struct {
	MailId string `json:"mail_id"` //来自哪个渠道的订单
}
type MailAddRequest struct {
	Uid          int64        `json:"uid"`           //发送者id，官方的话建议-1
	RecipientUid int64        `json:"recipient_uid"` //接受者id，
	Title        string       `json:"title"`         //标题
	Desc         string       `json:"desc"`          //详情
	ChannelCode  string       `json:"channel_code"`  //来源
	Attachment   []Attachment `json:"attachment"`    //附件
}
