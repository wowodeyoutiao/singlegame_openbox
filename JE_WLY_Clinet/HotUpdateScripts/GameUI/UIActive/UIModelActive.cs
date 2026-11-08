using HotUpdateScripts.ProtoM;

using MetaFramework.UI;

using pb;

using System.Collections.Generic;

using UnityEngine;

public class UIModelActive : IUIModel
{
    public struct ActiveTitle
    {
        public ActivityType activeType;
        public int cfgId;
        public bool ActiveOpen;
        public GameObject BtnOn;
        public string RedKey;
        public string activeName;
        public int shopType;
    }

    public Dictionary<ActivityType, ActiveTitle> ActiveTitlesDic = new Dictionary<ActivityType, ActiveTitle>()
    {
        [ActivityType.Sign7] = new ActiveTitle
        {
            activeType = ActivityType.Sign7,
            cfgId = 1,
            ActiveOpen = ActivityM.Instance.GetAcitvityState(ActivityType.Sign7),
            RedKey = RedPointConst.ActiveCenter_Sign7.ToString(),
            activeName = "七日签到",
        },

        [ActivityType.SignHappy7] = new ActiveTitle
        {
            activeType = ActivityType.SignHappy7,
            cfgId = 19,
            ActiveOpen = ActivityM.Instance.GetAcitvityState(ActivityType.SignHappy7),
            RedKey = TaskClass.SEVENDAY.ToString(),
            activeName = "七日狂欢",
            shopType = 25,
        },

        [ActivityType.StarWish] = new ActiveTitle
        {
            activeType = ActivityType.StarWish,
            cfgId = 26,
            ActiveOpen = ActivityM.Instance.GetAcitvityState(ActivityType.StarWish),
            RedKey = RedPointConst.ActiveCenter_StarWish.ToString(),
            activeName = "星愿召唤",
            shopType = 36,
        },

        [ActivityType.CrazyRefine] = new ActiveTitle
        {
            activeType = ActivityType.CrazyRefine,
            cfgId = 27,
            ActiveOpen = ActivityM.Instance.GetAcitvityState(ActivityType.CrazyRefine),
            RedKey = RedPointConst.ActiveCenter_CrazyRefine.ToString(),
            activeName = "疯狂洗练",
            shopType = 37,
        },

        [ActivityType.RefineFund] = new ActiveTitle
        {
            activeType = ActivityType.RefineFund,
            cfgId = 28,
            ActiveOpen = ActivityM.Instance.GetAcitvityState(ActivityType.RefineFund),
            RedKey = RedPointConst.ActiveCenter_RefineFund.ToString(),
            activeName = "洗练基金",
            shopType = 38,
        },

        [ActivityType.CallFund] = new ActiveTitle
        {
            activeType = ActivityType.CallFund,
            cfgId = 29,
            ActiveOpen = ActivityM.Instance.GetAcitvityState(ActivityType.CallFund),
            RedKey = RedPointConst.ActiveCenter_CallFund.ToString(),
            activeName = "召唤基金",
            shopType = 39,
        },

        [ActivityType.CopyFund] = new ActiveTitle
        {
            activeType = ActivityType.CopyFund,
            cfgId = 30,
            ActiveOpen = ActivityM.Instance.GetAcitvityState(ActivityType.CopyFund),
            RedKey = RedPointConst.ActiveCenter_CopyFund.ToString(),
            activeName = "副本基金",
            shopType = 40,
        },
    };
}