using HotUpdateScripts.ProtoM;
using MetaFramework.UI;
using System.Collections.Generic;
using UnityEngine;

public class UIModelGiftBag : IUIModel
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
        [ActivityType.PartnerGift] = new ActiveTitle
        {
            activeType = ActivityType.PartnerGift,
            cfgId = 11,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_PartnerGift.ToString(),
            activeName = "伙伴礼包",
            shopType = 17,
            iconName = "bg_Gift_Package_banner_4",
        },
        [ActivityType.PetGift] = new ActiveTitle
        {
            activeType = ActivityType.PetGift,
            cfgId = 12,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_PetGift.ToString(),
            activeName = "宠物礼包",
            shopType = 18,
            iconName = "bg_Gift_Package_banner_2",
        },
        [ActivityType.Gemstone] = new ActiveTitle
        {
            activeType = ActivityType.Gemstone,
            cfgId = 14,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_Gemstone.ToString(),
            activeName = "魂石礼包",
            shopType = 21,
            iconName = "bg_Gift_Package_banner_5",
        },
        [ActivityType.Element] = new ActiveTitle
        {
            activeType = ActivityType.Element,
            cfgId = 17,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_Element.ToString(),
            activeName = "元素礼包",
            shopType = 24,
            iconName = "bg_Gift_Package_banner_6",
        },
        [ActivityType.BlessingGift] = new ActiveTitle
        {
            activeType = ActivityType.BlessingGift,
            cfgId = 22,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_Blessing.ToString(),
            activeName = "福袋礼包",
            shopType = 31,
            iconName = "bg_Gift_Package_banner_9",
        },
        [ActivityType.RefineGift] = new ActiveTitle
        {
            activeType = ActivityType.RefineGift,
            cfgId = 25,
            ActiveOpen = true,
            RedKey = RedPointConst.GiftBag_Refine.ToString(),
            activeName = "洗练礼包",
            shopType = 35,
            iconName = "bg_Gift_Package_banner_11",
        },
    };
}