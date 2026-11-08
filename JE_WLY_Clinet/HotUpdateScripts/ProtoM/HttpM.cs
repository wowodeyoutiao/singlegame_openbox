using JEngine.Singleton;
using pb;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// Http信息
    /// </summary>
    public class HttpM : SingletonTemplate<HttpM>
    {
        public readonly string PrivilegeDetailUrl = "/v1/prize_draw/detail";
        public readonly string PrivilegeAwardGetUrl = "/v1/prize_draw/partner_receive";
        public readonly string PetReceiveGetUrl = "/v1/prize_draw/pet_receive";
    }

    /// <summary>
    /// 性别枚举
    /// </summary>
    public enum SexType
    {
        SexNo = 0,
        SexBoy = 1,
        SexGril = 2,
    }

    /// <summary>
    /// 职业枚举
    /// </summary>
    public enum JobType
    {
        /// <summary>
        /// 战士
        /// </summary>
        Warrior = 1,
        /// <summary>
        /// 法师
        /// </summary>
        Mage = 2,

        /// <summary>
        /// 祭祀
        /// </summary>
        Sacrifice = 3,
    }

    public struct HTTP_VersionRequest
    {
        public int code;
        public HTTP_VersionItem data;
        public string msg;

        public struct HTTP_VersionItem
        {
            public long created_at;
            public string version;
        }
    }

    public struct HTTP_WXSeeionRequest
    {
        public string openid;
    }

    public struct HTTP_WXSeeionResponse
    {
        public int code;
        public string data;
    }

    public struct HTTP_WXLoginResponse
    {
        public int code;
        public string msg;
        public HTTP_WXLoginItem data;

        public struct HTTP_WXLoginItem
        {
            public string openid;
            public int errcode;
            public string errmsg;
        }

    }

    public struct HTTP_WXLoginRequest
    {
        public string code;
        public string clickid;
    }

    public struct HTTP_SeverRequest
    {
        public int pageIndex;
        public int pageSize;
        public string code;
        public int show;
        public string version;
    }

    public struct HTTP_SeverResponse
    {
        public int code;
        public string msg;
        public HTTP_SeverResponseItem data;
        public struct HTTP_SeverResponseItem
        {
            public int count;
            public int pageIndex;
            public int pageSize;
            public List<GameM.SeverItem> list;
            public List<GameM.SeverItem> mylist;
            /// <summary>
            /// 版本状态 2 为审核
            /// </summary>
            public int version_status;
        }

    }


    /// <summary>
    /// HTTP 登录消息发送
    /// </summary>
    public struct HTTP_LoginRequest
    {
        public string appid;
        public string channelId;
        public string imei;
        public string platform;
        public string openid; //微信平台openid
        public string access_token; //微信AccessToken
        public int serverId;
        public HTTP_Device device;
        public string clickid;

        /// <summary>
        /// 设备信息
        /// </summary>
        public struct HTTP_Device
        {
            public string imei;     //设备的imei号
            public string os;       //os版本号
            public string model;    //硬件型号
            public string ip;       //内网IP
            public string remote;   //外网IP
        }
    }

    /// <summary>
    /// HTTP 实名认证
    /// </summary>
    public struct HTTP_RealNameRequest
    {
        public string username { get; set; }

        public string idcard { get; set; }
    }

    /// <summary>
    /// HTTP 实名认证
    /// </summary>
    public struct HTTP_RealNameResponse
    {
        public int code { get; set; }

        public string message { get; set; }
    }

    /// <summary>
    /// HTTP 登录消息返回
    /// </summary>
    public struct HTTP_LoginResponse
    {
        public int code { get; set; }
        public List<HTTP_CharacterInfo> characterList { get; set; }
        public HTTP_ClientConfig config { get; set; }

        public List<string> messages { get; set; }
        public string token { get; set; }

        public bool realname { get; set; }
    }

    /// <summary>
    /// 角色信息
    /// </summary>
    public class HTTP_CharacterInfo
    {
        public string name;
        public long uid;
        public string headUrl;
        public long coin;   //金币
        public long diamanta; //钻石
        public long soul; //战魂
        public long coupon; //点券
        public long energy; //体力
        public int job;   //职业
        public int sex; // 0: Unknown, 1: Male, 2: Female
        public long exp; //经验
        public int level; //等级
        public int vip; //vip等级
        public int recharge; //vip经验

        public long registerAt; // 角色注册时间

        public string ip;
        public int port;
        public string playerIp;

        public int debug;
    }

    public struct HTTP_ClientConfig
    {
        public string version { get; set; }
        public string android { get; set; }
        public string ios { get; set; }
        public int heartbeat { get; set; }
        public bool forceUpdate { get; set; }
        // 分享标题
        public string title { get; set; }
        // 分享描述
        public string desc { get; set; }
        public string daili1 { get; set; }
        public string daili2 { get; set; }
        public string kefu1 { get; set; }
        public string appId { get; set; }
        public string appKey { get; set; }
    }

    /// <summary>
    /// 创建角色
    /// </summary>
    public struct HTTP_CreateCharacterRequest
    {
        public int sex { get; set; }         // [0]未知 [1]男 [2]女
        public string name { get; set; }     // 名字 N1--N2 个字符
        public int job { get; set; }       // 职业
        public string clickid { get; set; }
    }

    /// <summary>
    /// 创建角色 返回
    /// </summary>
    public struct HTTP_CreateCharacterResponse
    {
        public int code { get; set; }
        public string date { get; set; }

        public HTTP_CharacterInfo characterList { get; set; }
    }

    /// <summary>
    /// 选择角色
    /// </summary>
    public struct HTTP_SelectCharacterRequest
    {
        public long uid { get; set; }
    }

    /// <summary>
    /// 选择角色 返回
    /// </summary>
    public struct HTTP_SelectCharacterResponse
    {
        public int code { get; set; }
        public string token { get; set; }
    }


    /// <summary>
    /// 角色删除
    /// </summary>
    public struct HTTP_LoginrateRequest
    {
        public long uid { get; set; }
    }

    /// <summary>
    /// 角色删除 返回
    /// </summary>
    public struct HTTP_LoginrateResponse
    {
        public int code { get; set; }
    }


    /// <summary>
    /// 兑换码
    /// </summary>
    public struct HTTP_ActivationOperateRequest
    {
        public string activation_code { get; set; }
    }

    /// <summary>
    /// 兑换码返回
    /// </summary>
    public struct HTTP_ActivationOperateResponse
    {
        public int code { get; set; }
        public string data { get; set; }
    }


    //命格

    /// <summary>
    /// 升级 命格
    /// </summary>
    public struct HTTP_FateUpLvRequest
    {
        public long uid { get; set; }
        public int config_id { get; set; }
        public int num { get; set; }
    }
    /// <summary>
    /// 升级 命格
    /// </summary>
    public struct HTTP_FateUpLvResponse
    {
        public int code { get; set; }
        public int[] per { get; set; }
        public int num { get; set; }
    }

    /// <summary>
    /// 升阶 命格
    /// </summary>
    public struct HTTP_FateUpAdvanceResponse
    {
        public int code { get; set; }
        public string data { get; set; }
    }

    #region 伙伴抽取

    /// <summary>
    /// 伙伴召唤列表获取
    /// </summary>
    public struct HTTP_PrizedrawListResponse
    {
        public int code { get; set; }
        public string data { get; set; }
        public PrizedrawDetail[] canplay { get; set; }
        public struct PrizedrawDetail
        {
            public int type;// 类型 1-青铜白银黄金
            public int can_play; // 剩余可抽奖次数
            public int need_num; //还剩多少抽
        }

    }


    public struct HTTP_PrizedrawListRequest
    {
        public int module { get; set; }  //1-伙伴，2-宠物
    }

    public struct HTTP_PrizedrawRequest
    {
        public int ptype { get; set; }
        public int type { get; set; }
    }
    //伙伴抽奖领取
    public struct HTTP_PrizedrawResponse
    {
        public int code { get; set; }
        public string data { get; set; }
        public PartnerAward[] award { get; set; }

        public struct PartnerAward
        {
            public int kind;//  类型 1-碎片 2-整卡
            public int id; // 物品id
            public int num; //数量
        }
    }

    #endregion

    #region 试炼基金

    //试炼基金 点击
    public struct HTTP_ActivityTrialFundRequset
    {
        public int activity_id { get; set; }
    }

    //试炼基金 奖励列表
    public struct HTTP_ActivityTrialFundAwardResponse
    {
        public int code { get; set; }
        public int day { get; set; }
        public int buy { get; set; }
        public HTTP_ActivityTrialFundDays[] data { get; set; }
    }

    //试炼基金 详情
    public struct HTTP_ActivityTrialFundDays
    {
        public int day { get; set; }
        public int buy { get; set; }
        public int status { get; set; }
    }

    //试炼基金 一键领取
    public struct HTTP_ActivityTrialFundGetAwardResponse
    {
        public int code { get; set; }
        public string data { get; set; }
    }
    #endregion

    #region ShuShu
    public struct HTTP_WindowShushuResponse
    {
        public int Code { get; set; }
    }

    public struct HTTP_WindowShushuRequest
    {
        public string Name { get; set; }
    }

    public struct HTTP_WindowShushuCompletedRequest
    {
        public string name { get; set; }
        public int guide_step { get; set; }
        public string guide_name { get; set; }
    }

    public struct HTTP_WindowShusOpenBoxRequest
    {
        public string name { get; set; }
        public HTTP_BoxItemInfo[] box_item_info { get; set; }
        public bool is_auto { get; set; }
        public int power_up { get; set; }
        public int open_num { get; set; }
        public bool role_power_stop { get; set; }
        public bool partner_power_stop { get; set; }
        public int box_num_before { get; set; }
        public int box_num_after { get; set; }
        public string box_job_limited { get; set; }
        public int box_level_limited { get; set; }
        public string retain_property { get; set; }
        public string retain_property_num { get; set; }
    }

    public struct HTTP_BoxItemInfo
    {
        public string item_id { get; set; }
        public int item_level { get; set; }
        public int item_quality { get; set; }
        public string item_name { get; set; }
        public int count { get; set; }
        public string item_property { get; set; }
    }

    public struct HTTP_WindowShushuAdShowRequest
    {
        public string name;
        public string ad_id;
        public string ad_type;
        public string ad_scene;
        public string ad_position;
    }

    public struct HTTP_WindowShushuAdCloseRequest
    {
        public string name;
        public string ad_id;
        public string ad_type;
        public string ad_scene;
        public int ad_time;
        public int ad_duration;
        public bool view_complete;
        public string reward_type;
        public string ad_position;
        public HTTP_GetItemInfo[] get_item_info;
    }

    public struct HTTP_GetItemInfo
    {
        public string item_id;
        public string item_name;
        public int item_num;
    }

    public struct HTTP_WindowShushuClickEventRequest
    {
        public string name;
        public string page;
        public string component;
        public string position;
        public string event_click;
    }
    #endregion


    #region 成长基金
    // 获取已经领取的情况
    public struct HTTP_GrowthFundRequset
    {
        public int activity_id { get; set; }
        public int atype { get; set; }
    }

    public struct HTTP_GrowthFundResponse
    {
        public int code { get; set; }
        public int day { get; set; }
        public int buy { get; set; } // 2 已购买
        public int num { get; set; } // 2 完成数量
        public HTTP_GrowthFundInfo[] data { get; set; }
    }

    // 每个等级的详情
    public struct HTTP_GrowthFundInfo
    {
        public int day { get; set; }
        public int buy { get; set; }
        public int status { get; set; } // 状态 1-待领取，2-已领取，3-不可领取，4-需要补签
    }

    //一键领取
    public struct HTTP_GrowthFundGetAwardResponse
    {
        public int code { get; set; }
        public string data { get; set; }
    }
    #endregion

    #region 修改名字
    public struct HTTP_ChangeRoleNameRequest
    {
        public string RoleName { get; set; }
    }

    public struct HTTP_ChangeRoleNameResponse
    {
        public int Code { get; set; }
        public string Error { get; set; }
        public string RoleName { get; set; }
    }
    #endregion

    #region 修改名字
    public struct HTTP_ChatRequest
    {
        public ChatType Type { get; set; }
        public string Context { get; set; }
    }

    public struct HTTP_ChatResponse
    {
        public int Code { get; set; }
        public string Error { get; set; }
    }
    #endregion

    #region 行会创建
    public struct HTTP_GuildCreateRequest
    {
        public string Name { get; set; }
        public string Description { get; set; }
    }

    public struct HTTP_GuildCreateResponse
    {
        public int Code { get; set; }
        public string Error { get; set; }
        public GuildInfo Guild { get; set; }
    }
    #endregion
}

