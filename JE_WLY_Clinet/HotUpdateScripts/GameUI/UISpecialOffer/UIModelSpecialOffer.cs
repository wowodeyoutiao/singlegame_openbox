using HotUpdateScripts.ProtoM;
using MetaFramework.UI;
using System.Collections.Generic;
using UnityEngine;

public class UIModelSpecialOffer : IUIModel
{
    public struct ActiveTitle
    {
        public ActivityType activeType;
        public int cfgId;
        public bool ActiveOpen;
        public GameObject BtnOn;
        public string RedKey;
        public string activeName;
        public string iconName;
        public int shopType;
    }

    public Dictionary<ActivityType, ActiveTitle> ActiveTitlesDic = new Dictionary<ActivityType, ActiveTitle>()
    {
        [ActivityType.Month] = new ActiveTitle
        {
            activeType = ActivityType.Month,
            ActiveOpen = true,
            RedKey = RedPointConst.PrivilegeCard.ToString(),
            activeName = "月卡",
            iconName = "",
        },
        [ActivityType.GrowthFund] = new ActiveTitle
        {
            activeType = ActivityType.GrowthFund,
            cfgId = 15,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_GrowthFund.ToString(),
            activeName = "成长基金"
        },
        [ActivityType.PowerGift] = new ActiveTitle
        {
            activeType = ActivityType.PowerGift,
            cfgId = 21,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_Power.ToString(),
            activeName = "体力礼包",
            shopType = 29,
            iconName = "bg_Gift_Package_banner_8",
        },
        [ActivityType.TimeGift] = new ActiveTitle
        {
            activeType = ActivityType.TimeGift,
            cfgId = 23,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_Time.ToString(),
            activeName = "挂机礼包",
            shopType = 32,
            iconName = "bg_Gift_Package_banner_10",
        },
        [ActivityType.GiftPay] = new ActiveTitle
        {
            activeType = ActivityType.GiftPay,
            cfgId = 4,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_GiftPay.ToString(),
            activeName = "每日特价", // 原特惠礼包
            shopType = 7,
            iconName = "bg_Gift_Package_banner_3",
        },[ActivityType.DailyGift] = new ActiveTitle
        {
            activeType = ActivityType.DailyGift,
            cfgId = 13,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_DailyGift.ToString(),
            activeName = "每日礼包",
            shopType = 20,
            iconName = "bg_Gift_Package_banner_1",
        },
        [ActivityType.AdGift] = new ActiveTitle
        {
            activeType = ActivityType.AdGift,
            cfgId = 20,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_AdGift.ToString(),
            activeName = "广告礼包",
            shopType = 27,
            iconName = "bg_Gift_Package_banner_7",
        },
    };
}