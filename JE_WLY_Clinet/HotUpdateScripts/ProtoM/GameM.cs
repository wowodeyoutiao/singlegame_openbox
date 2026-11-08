using JEngine.Singleton;
using System.Collections.Generic;
using static SDKCtr;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 游戏总信息
    /// </summary>
    public class GameM : SingletonTemplate<GameM>
    {
        public string[] EquipName = { "", "武器", "衣服", "头盔", "项链", "左戒", "右戒", "左手镯", "右手镯", "腰带", "靴子", "护腿", "饰品" };

        public int SeverIdx
        {
            get
            {
                return PlayerPrefsExtend.GetInt("SeverIdx", 0);
            }
            set { PlayerPrefsExtend.SetInt("SeverIdx", value); }
        }

        public int SeverDicIdx
        {
            get
            {
                return PlayerPrefsExtend.GetInt("SeverDicIdx", 0);
            }
            set { PlayerPrefsExtend.SetInt("SeverDicIdx", value); }
        }

        private static string BaseUrl = "http://127.0.0.1:18080";
        public string SeverInfoUrl = BaseUrl + "/goadminapi/api/v1/servername";
        public string WXInfoUrl = "https://payment.shbinqigame.com/v1/wx/user-login/code2session";
        public readonly string SessionsCheakUrl = "https://payment.shbinqigame.com/v1/wx/user-login/check-session-key"; //put

        //版本检测
        public static string VersionUrl = BaseUrl + "/goadminapi/api/v1/servername/version";

        public string WXOpneID = "";
        public string WXCode = "";

        //服务器主动关闭
        public bool isSeverClose = false;

        public float BaseScale = 0;

        //总服务器列表
        public Dictionary<int, List<SeverItem>> SeverInfoDic = new Dictionary<int, List<SeverItem>>();

        public int SeverCount = 0;

        private bool _isAudit = false;

        public int VersionStatus = 0; // 1:线上版本 2:审核版本3 体验服

        public bool SkipAD = false;

        /// <summary>
        /// 是否是审核版本
        /// </summary>
        public bool IsAudit
        {
            get
            {
                if (_isAudit)
                    return _isAudit;

                if (VersionStatus == 2 && (PlayerPrefsExtend.GetString(SaveKey.WX_Platform.ToString()) == "ios" || PlayerPrefsExtend.GetString(SaveKey.WX_Platform.ToString()) == "mac"))
                {
                    _isAudit = true;
                    return true;
                }
                else
                    return false;
            }

        }

        public class SeverItem
        {
            public SeverItem(int _id, string _des, string _url, int _severType, int _server_id)
            {
                name = _des;
                ipaddr = _url;
                status = _severType;
                id = _id;
                server_id = _server_id;
            }
            public string name;
            public string ipaddr;
            public int status;
            public int id;
            public int server_id;
            public long open_by;
            public string websocket;
        }

        // ===== Local single-server support (no go-admin) =====
        public const int LocalServerId = 9;
        public const string LocalSeverUrl = "http://127.0.0.1:12307";
        public const string LocalWebSocket = "127.0.0.1:7000";

        public void InitLocalSeverInfo()
        {
            SeverInfoDic.Clear();
            SeverInfoDic[0] = new List<SeverItem> { CreateLocalSeverItem() };
            SeverInfoDic[1] = new List<SeverItem> { CreateLocalSeverItem() };
            SeverIdx = 0;
            SeverDicIdx = 0;
            SeverCount = 0;
        }

        private static SeverItem CreateLocalSeverItem()
        {
            var item = new SeverItem(LocalServerId, "本地服", LocalSeverUrl, 1, LocalServerId);
            item.server_id = LocalServerId;
            item.open_by = 0;
            item.websocket = LocalWebSocket;
            return item;
        }

        public string HttpUrl
        {
            get
            {

                //return NewSeverItems[SeverIdx].ipaddr;
                return SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx].ipaddr;
            }
        }

        public readonly string LoginUrl = "/v1/user/login/guest";
        public readonly string Login3rdUrl = "/v1/user/login/3rd";
        public readonly string RealNameUrl = "/v1/user/login/realname";
        public readonly string Role_CreatUrl = "/v1/user/login/create_char";
        public readonly string Role_SelectUrl = "/v1/user/login/select_char";
        public readonly string Role_DeleteUrl = "/v1/user/login/del_char";

        public readonly string Store_GetListUrl = "/v1/store/list";
        public readonly string Store_BuyUrl = "/v1/store/buy";

        public readonly string CDK_Url = "/v1/activation/check";

        public readonly string Mail_GetListUrl = "/v1/mail/list";
        public readonly string Mail_DetailUrl = "/v1/mail/detail";
        public readonly string Mail_ExtractUrl = "/v1/mail/extract";
        public readonly string Mail_ExtractAllUrl = "/v1/mail/extractall";
        public readonly string Mail_DelUrl = "/v1/mail/del";
        public readonly string Mail_DelAllUrl = "/v1/mail/delall";

        public readonly string Activity_GetSignListUrl = "/v1/activity/listcontinuity"; //
        public readonly string Activity_SignUrl = "/v1/activity/receiveday"; //PUT
        public readonly string Activity_GetSign30ListUrl = "/v1/activity/listcumulative"; //
        public readonly string Activity_Sign30GetUrl = "/v1/activity/receivecumulative"; //PUT
        public readonly string Activity_ActiveStateUrl = "/v1/activity/listall"; //GET
        public readonly string Activity_GetLevelGiftListUrl = "/v1/activity/listgrade";
        public readonly string Activity_GetLevelGiftUrl = "/v1/activity/receivegrade"; //PUT

        public readonly string Activity_FirstPayListUrl = "/v1/activity/listcharge"; //
        public readonly string Activity_FirstPayGetUrl = "/v1/activity/receivecharge"; //PUT

        public readonly string Activity_TrialFundListUrl = "/v1/activity/listtrial"; //
        public readonly string Activity_TrialFundGetUrl = "/v1/activity/receivetrial"; //PUT

        // 成长基金
        public readonly string Activity_GrowthFundListUrl = "/v1/activity/groupup"; //
        public readonly string Activity_GrowthFundGetUrl = "/v1/activity/receivegroupup"; //PUT

        public readonly string Sign_ListUrl = "/v1/sign/list"; //GET
        public readonly string Sign_SignDayUrl = "/v1/sign/add";//PUT
        public readonly string Sign_SignExtractUrl = "/v1/sign/extract";//post

        public readonly string Privilege_ListUrl = "/v1/privilege/listcontinuity";//GET
        public readonly string Privilege_SignUrl = "/v1/privilege/receive";//PUT

        public readonly string Rank_ListUrl = "/v1/rank/list";//GET

        public readonly string Store_ExtractAccount = "/v1/store/extractaccount";//GET
        public readonly string Store_ExtractAccountUpdate = "/v1/store/extractaccountupdate";//Put

        public readonly string Activity_ReceivetreasureBoxUrl = "/v1/activity/receivetreasurebox";//PUT
        public readonly string Activity_GetData_ReceivetreasureBoxUrl = "/v1/activity/listtreasurebox";//GET

        public readonly string Activity_ListTurntableUrl = "/v1/activity/listturntable";//Get
        public readonly string Activity_ReceiveTurntableUrl = "/v1/activity/receiveturntable";//PUT

        public readonly string Fate_UpLvUrl = "/v1/fate/up";//POST
        public readonly string Fate_UpAdvanceUrl = "/v1/fate/break";//PUT

        public readonly string ShuShu_Common = "/v1/window/shushu";//POST
        public readonly string ShuShu_AppHide = "/v1/window/ta_mp_hide";//POST
        public readonly string ShuShu_GuideCompleted = "/v1/window/completed";//POST
        public readonly string ShuShu_OpenBox = "/v1/window/openBox";//POST
        public readonly string ShuShu_AdShow = "/v1/window/ad_show";//POST
        public readonly string ShuShu_adClose = "/v1/window/ad_close";//POST
        public readonly string ShuShu_clickEvent = "/v1/window/click_event";//POST

        public readonly string ChangeRoleName = "/v1/sensitive/role/name";//PUT 改名字
        public readonly string Chat = "/v1/sensitive/messages";//POST 聊天
        public readonly string CreateGuild = "/v1/sensitive/guilds";//POST 创建工会

        /// <summary>
        /// 服务器时间偏差
        /// </summary>
        public long SeverTimeOff = 0;

        public double SeverTime { get; set; }
    }

}

