using JEngine.Core;
using JEngine.Helper;
using JEngine.Singleton;

using System.Collections.Generic;
using System.Linq;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 活动
    /// </summary>
    public class ActivityM : SingletonTemplate<ActivityM>
    {
        public enum SignState
        {
            CanGet = 1,
            Geted = 2,
            PayGet = 4,
            NoGet = 3,
        }

        private List<CFGItem_Activity> showActivityList; // 活动表

        #region 七日 30日 签到 

        private Dictionary<int, HTTP_ActivityListDay> SignStateDic = new Dictionary<int, HTTP_ActivityListDay>();
        public void SetSignDic(HTTP_ActivityListResponse msg)
        {
            SignStateDic.Clear();
            if (msg.data != null)
            {
                for (int i = 0; i < msg.data.Length; i++)
                {
                    SignStateDic.Add(msg.data[i].day, msg.data[i]);
                }
            }
        }

        public HTTP_ActivityListDay GetSignStateByDay(int day)
        {
            HTTP_ActivityListDay a = new HTTP_ActivityListDay { day = 0, status = (int)SignState.CanGet, buy = (int)SignState.CanGet };
            if (!SignStateDic.ContainsKey(day))
            {
                return a;
            }
            return SignStateDic[day];
        }

        public int GetSign7DotNum()
        {
            if (!GetAcitvityState(ActivityType.Sign7))
                return 0;

            showActivityList = CFG_Activity.GetActiveListByType((CFG_Activity.Type)1);
            if (showActivityList == null || showActivityList.Count == 0)
                return 0;

            foreach (var item in SignStateDic.Values)
            {
                if ((SignState)item.status == SignState.CanGet)
                    return 1;

                if ((SignState)item.status == SignState.PayGet && int.Parse(showActivityList[item.day - 1].CompensationConsume[0].ElementAt(0).Key) == 51002)
                    return 1;

            }
            return 0;
        }
        #endregion

        #region 等级礼包

        private Dictionary<int, int> LevelGiftDic = new Dictionary<int, int>();

        public void SetLevelGiftDic(HTTP_ActivityListGrade[] data)
        {
            LevelGiftDic.Clear();
            foreach (var item in data)
            {
                LevelGiftDic.Add(item.grade, item.status);
            }
        }

        public SignState GetOneLevelGiftStat(int grade)
        {
            return (SignState)LevelGiftDic[grade];
        }

        public int GetLevelDotNum()
        {
            foreach (var item in LevelGiftDic.Values)
            {
                if ((SignState)item == SignState.CanGet)
                {
                    return 1;
                }

            }
            return 0;
        }
        #endregion


        //活动开启状态
        private Dictionary<ActivityType, HTTP_ItemActivityAll> DicActivityState = new Dictionary<ActivityType, HTTP_ItemActivityAll>();

        public void SetActivityState(HTTP_ItemActivityAll[] AllState)
        {
            DicActivityState.Clear();
            foreach (var item in AllState)
            {
                int addKet = int.Parse($"{item.id}{item.act_type}");
                Log.Print($"活动状态开关： id = {addKet}, time = {item.end_time}");
                DicActivityState.Add((ActivityType)addKet, item);
                //UnityEngine.Debug.LogError($"活动状态开关： id = {addKet}, time = {item.end_time}");
            }
        }

        public bool GetAcitvityState(ActivityType activityType)
        {
            if (DicActivityState.ContainsKey(activityType))
            {
                if (DicActivityState[activityType].end_time == -1)
                    return true;
                else
                {
                    return GameUtil.Instance.GetCurrentTimeSeconds() < DicActivityState[activityType].end_time;
                }

            }
            else
            {
                return false;
            }
        }

        public long GetAcitvityEndTime(ActivityType activityType)
        {
            if (DicActivityState.ContainsKey(activityType))
            {

                return DicActivityState[activityType].end_time - GameUtil.Instance.GetCurrentTimeSeconds();
            }
            else
            {
                return -1;
            }
        }

        public void CloseAcitvity(ActivityType activityType)
        {
            if (DicActivityState.ContainsKey(activityType))
            {
                var data = DicActivityState[activityType];
                data.end_time = GameUtil.Instance.GetCurrentTimeSeconds() - 100;
                DicActivityState[activityType] = data;
            }
        }


        #region 试炼基金可领取奖励数量

        List<CFGItem_Activity> cfg_TrialFundItems;

        public void GetTrialFundAwardItems()
        {
            cfg_TrialFundItems = CFG_Activity.GetActiveListByType(CFG_Activity.Type.TrialFund);
            cfg_TrialFundItems.Sort((a, b) => a.Content.CompareTo(b.Content));
        }

        public int GetTrialFundCanGetAwardCount(HTTP_ActivityTrialFundAwardResponse response)
        {
            GetTrialFundAwardItems();
            int count = 0;

            for (int i = 0; i < cfg_TrialFundItems.Count; i++)
            {
                bool isCanReceivedTheDay = cfg_TrialFundItems[i].PreLimit <= response.day;
                bool isReceived = i < response.data.Length;
                int status = 3;
                int currentFloorIsBuy = isReceived ? response.data[i].buy : response.buy;

                if (isReceived)
                    status = 2;

                else if (isCanReceivedTheDay)
                    status = 1;

                for (int j = 0; j < 4; j++)
                {
                    bool canGet;

                    if (i == 0)
                        canGet = status == 1;
                    else
                        canGet = (status == 1 && currentFloorIsBuy == 2) || (status == 2 && currentFloorIsBuy == 1 && response.buy == 2);

                    if (canGet)
                        count++;

                }
            }

            return count;
        }

        #endregion

        public bool IsBuyFirstPayGift()
        {
            var cfg_activeItems = CFG_Activity.GetActiveListByType(CFG_Activity.Type.FirstPay);
            long buyCount = 0;
            foreach (var item in cfg_activeItems)
            {
                buyCount += StoreM.Instance.GetBuyNumByID(item.PreLimit);
            }
            return buyCount > 0;
        }

    }


    /// <summary>
    /// 设置活动开启类型 
    /// 由 HTTP_ItemActivityAll id 和 act_type 拼合
    /// </summary>
    public enum ActivityType
    {
        NO = 0,
        FirstPay = 15,
        Sign30 = 12,
        Sign7 = 11,
        Level = 13,
        GiftPay = 14,//特惠礼包
        LimitBounty = 17,//限时悬赏
        TurnWheel = 16,//点券大转盘
        SinglePay = 18,//单笔充值
        CumulativePay = 19,//累计充值
        //SLJJ = 20,//试炼基金
        SignDay = 60,//签到
        PartnerGift = 111,//伙伴礼包
        PetGift = 112,//宠物礼包
        DailyGift = 113, // 每日礼包
        Gemstone = 114, // 魂石礼包
        GrowthFund = 115, // 成长基金
        Element = 117, // 元素
        SignHappy7 = 119, // 七日狂欢
        AdGift = 120, // 广告礼包
        PowerGift = 121, // 体力礼包
        BlessingGift = 122, // 福袋礼包
        TimeGift = 123, // 时间礼包
        BattleSupportGift = 124, // 助战礼包
        RefineGift = 125, // 洗练礼包
        StarWish = 126, // 星愿召唤
        CrazyRefine = 127, // 疯狂洗练
        RefineFund = 128, // 洗练基金
        CallFund = 129, // 召唤基金
        CopyFund = 130, // 副本基金
        Month = 10000, // 月卡
    }

    /// <summary>
    /// 等级领取 活动信息列表
    /// </summary>
    public struct HTTP_ActivityGradeResponse
    {
        public int code { get; set; }
        public string data { get; set; }
        public HTTP_ActivityListGrade[] activity_list_grade { get; set; }

    }

    public struct HTTP_ActivityListGrade
    {
        //等级
        public int grade { get; set; }

        //状态 1-待领取，2-已领取
        public int status { get; set; }
    }

    //等级礼包领取
    public struct HTTP_ActivityGradeRequset
    {
        //活动id
        public int activity_id { get; set; }

        public int uuid { get; set; }
    }


    /// <summary>
    /// 连续登陆 
    /// </summary>
    public struct HTTP_ActivityListRequset
    {
        public int activity_id { get; set; }
    }

    /// <summary>
    /// 连续登陆 累计登录 活动信息列表
    /// </summary>
    public struct HTTP_ActivityListResponse
    {
        public int code { get; set; }
        public int day { get; set; }
        public int buy { get; set; }
        public HTTP_ActivityListDay[] data { get; set; }

    }

    public struct HTTP_ActivityListDay
    {
        /// <summary>
        /// 第几天
        /// </summary>
        public int day { get; set; }
        /// <summary>
        /// 状态 1-待领取，2-已领取，3-不可领取，4-需要补签
        /// </summary>
        public int status { get; set; }

        /// <summary>
        /// 1-待领取，2-已领取
        /// </summary>
        public int buy { get; set; }
    }

    /// <summary>
    /// 连续登陆的签到
    /// </summary>
    public struct HTTP_ActivityReceiveDayRequset
    {
        public int ad { get; set; } // > 0消耗物品，其它不消化
        public int uuid { get; set; }
    }

    /// <summary>
    /// 累计 连续登陆的签到 等级礼包领取
    /// </summary>
    public struct HTTP_ActivityReceiveDayResponse
    {
        public int code { get; set; }
        public string data { get; set; }
    }

    /// <summary>
    /// 活动开启详情
    /// </summary>
    public struct HTTP_ActivityAllResponse
    {
        public int code { get; set; }
        public HTTP_ItemActivityAll[] data { get; set; }
    }

    /// <summary>
    /// 活动开启详情
    /// </summary>
    public struct HTTP_ItemActivityAll
    {
        public int id { get; set; }
        public int act_type { get; set; }

        public long start_time { get; set; }

        public long end_time { get; set; }

    }

    /// <summary>
    /// 首充列表
    /// </summary>
    public struct HTTP_ActivityChargeResponse
    {
        public int code { get; set; }
        public int day { get; set; }
        public string data { get; set; }

        //1-待领取，2-已领取，3-不可领取（没充值）
        public int status { get; set; }
    }

    /// <summary>
    /// 百宝箱抽奖列表
    /// </summary>
    public struct HTTP_ActivityListTreasureBoxResponse
    {
        public int code { get; set; }
        public string data { get; set; }

        public int can_play { get; set; }
        public int need_num { get; set; }
    }

    /// <summary>
    /// 百宝箱
    /// </summary>
    public struct HTTP_ActivityTreasureBoxRequest
    {
        // 1-抽   2-10抽
        public int type { get; set; }
    }

    /// <summary>
    /// 百宝箱 返回
    /// </summary>
    public struct HTTP_ActivityTreasureBoxResponse
    {
        public int code { get; set; }
        public string data { get; set; }

        public int[] position { get; set; }
    }

    //点券大转盘列表 点击 send
    public struct HTTP_ActivityTruntableRequest
    {
        public int activity_id { get; set; }
    }

    //点券大转盘列表 返回
    public struct HTTP_ActivityTruntableResponse
    {
        public int code { get; set; }

        //当前准备抽的内容编号 -1就是已经没有可以抽的了
        public int content { get; set; }

        public string error { get; set; }
    }

    //点券大转盘 返回 点击 
    public struct HTTP_ActivityTruntableClickResponse
    {
        public int code { get; set; }
        public string data { get; set; }

        public int puid { get; set; }
    }
}

