package protocol

import "time"

type CreateRankRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type RankInfo struct {
}
type RankOperateRequest struct {
	Uid      int64  //Token
	Name     string //个人姓名
	Type     int64  `json:"type"`      //排行榜类型 1-等级，2-战力榜，3-乱斗榜，4-试炼榜 20-世界boss 21-世界boss跨服
	All      int    `json:"all"`       //1-所有数据（排行榜+个人），2-排行榜，3-个人
	ServerId int32  `json:"server_id"` //服务器ID
	ZoneId   int32  `json:"zone_id"`   //战区ID
	Count    int    `json:"count"`     //个数
}
type RankOperateResponse struct {
	Code int  `json:"code"` //状态码
	Data bool `json:"data"` //渠道列表
}

type RankListRequest struct {
	Uid  int64 //Token
	Type int64 `json:"type"` //排行榜类型 1-等级，2-战力榜，3-乱斗榜，4-试炼榜
	All  int   `json:"all"`  //1-所有数据（排行榜+个人），2-排行榜，3-个人
}

// RankListResponse 订单列表
type RankListResponse struct {
	Code int        `json:"code"` //状态码
	Data RankDetail `json:"data"` //详情
}
type RankDetail struct {
	RankAll []RankUserInfo `json:"rank_all"` //总排行榜
	RankOwn RankUserInfo   `json:"rank_own"` //个人榜单信息
}
type RankUserInfo struct {
	Uid        int64  `json:"uid"`         //
	Num        int64  `json:"num"`         //名次
	Type       int64  `json:"type"`        //排行榜类型 1-等级，2-战力榜，3-乱斗榜，4-试炼榜
	Avatar     string `json:"avatar"`      //头像
	Name       string `json:"name"`        //昵称
	Value      int64  `json:"value"`       //数值
	Sex        int32  `json:"sex"`         //性别
	Job        int32  `json:"job"`         //职业
	RobotCId   int32  `json:"robot_cid"`   //机器人ID
	ServerName string `json:"server_name"` //服务器名称
	AvatarV2   int32  `json:"avatar_v2"`   //头像
	Frame      int32  `json:"frame"`       //头像框
	// ZoneId   int32  `json:"zone_id"`   //战区ID
}
type RankUserInfoInput struct {
	Uid          int64  `json:"uid"`            //
	Type         int32  `json:"type"`           //排行榜类型 1-等级，2-战力榜，3-乱斗榜，4-试炼榜
	Compare      bool   `json:"compare"`        //是否需要比较
	Mysql        bool   `json:"mysql"`          //数据是否需要落地到数据库
	Avatar       string `json:"avatar"`         //头像
	AvatarV2     int32  `json:"avatar_v2"`      //头像string `json:"avatar"`         //头像
	Frame        int32  `json:"frame"`          //头像框
	Name         string `json:"name"`           //昵称
	Sex          int32  `json:"sex"`            //性别
	Job          int32  `json:"job"`            //职业
	Value        int64  `json:"value"`          //数值
	ServerId     int32  `json:"server_id"`      //服务器ID
	RankServerId int32  `json:"rank_server_id"` //排行榜服务器ID
	UserServerId int32  `json:"user_server_id"` //用户服务器ID
	RobotCId     int32  `json:"robot_cid"`      //机器人ID
	ServerName   string `json:"server_name"`    //服务器名称
	Count        int    `json:"count"`          //个数
	ZoneId       int32  `json:"zone_id"`        //战区ID
}

//	type RankListCharge struct {
//		Day    int32 `json:"day"`    //第几天
//		Status int   `json:"status"` //状态 1-待领取，2-已领取，3-不可领取（没充值）
//	}
type RankDetailDay struct {
	Pid     int64 `json:"pid"`      //内容编号
	Buy     int   `json:"buy"`      //1-未购买，2-已购买
	Get     int   `json:"get"`      //状态 1-待领取，2-已领取，
	EndTime int64 `json:"end_time"` //到期时间
}
type RankDetailTCP struct {
	Pid     int64     `json:"pid"`      //内容编号
	EndTime time.Time `json:"end_time"` //到期时间
}

// RankListResponse 订单列表
type RankChargeResponse struct {
	Code   int    `json:"code"`   //状态码
	Day    int64  `json:"day"`    //当前天数
	Data   string `json:"data"`   //返回信息
	Status int    `json:"status"` //状态 1-待领取，2-已领取，3-不可领取（没充值）
}
