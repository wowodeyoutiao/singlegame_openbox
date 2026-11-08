using HotUpdateScripts.ProtoM;

using JEngine.Helper;

using System;
using System.Collections.Generic;

namespace MetaFramework.UI
{
    public class RedPointC
    {
        public static Dictionary<RedPointConst, Func<int>> RedPointNumDic = new Dictionary<RedPointConst, Func<int>>
        {
            [RedPointConst.Role_Equip] = () => RedDotM.Instance.GetBeatEquipRed(),
            [RedPointConst.Pet_Award] = () => RedDotM.Instance.GetPetAwardNum(),
            [RedPointConst.Pet_UpStar_1] = () => RedDotM.Instance.GetPetUpStar(),

            [RedPointConst.Partner_Into] = () => RedDotM.Instance.PartnerIntoRedDot(),
            [RedPointConst.Partner_Call_Free] = () => RedDotM.Instance.GetPartnerCallRedDot(),
            [RedPointConst.Partner_UpStar] = () => RedDotM.Instance.GetPartnerUpStarRedDot(),

            [RedPointConst.Equip_Reinforce] = () => RedDotM.Instance.GetAllEquipSlotReinforceRedNum(),
            [RedPointConst.Equip_Insert] = () => RedDotM.Instance.GetAllEquipInsetRedNum(),
            [RedPointConst.Equip_Star] = () => RedDotM.Instance.GetEquipStarRedNum(),
            [RedPointConst.Equip_Refine] = () => RedDotM.Instance.GetEquipRefineNum(),

            [RedPointConst.Skill_UpLv] = () => RedDotM.Instance.GetAllSkillUpLvRedNum(),
            [RedPointConst.Skill_Advance] = () => RedDotM.Instance.GetAllSkillAdvanceRedNum(),
            [RedPointConst.Skill_Setting] = () => RedDotM.Instance.GetSkillSettingRedNum(),

            //[RedPointConst.ExpPool] = () => RedDotM.Instance.GetExpPoolRedNum(),

            [RedPointConst.Challenge_Fight_Click] = () => FightPlayM.Instance.GetChallenge_FightRedNum(),
            [RedPointConst.Challenge_Fight_Award] = () => FightPlayM.Instance.GetChallenge_FightAwardRedNum(),
            [RedPointConst.Challenge_Boss] = () => FightPlayM.Instance.GetChallenge_BossRedNum(),
            [RedPointConst.Challenge_Tower_Click] = () => FightPlayM.Instance.GetChallenge_TowerRedNum(),
            [RedPointConst.Challenge_Relice] = () => FightPlayM.Instance.GetChallenge_ReliceRedNum(),
            [RedPointConst.Challenge_OnHook] = () => FightPlayM.Instance.GetChallenge_OnHook(),

            [RedPointConst.Bag_Syn] = () => RedDotM.Instance.GetSynRedNum(),

            [RedPointConst.Rank_HuiHe] = () => RedDotM.Instance.GetRankHuiHe(),
            [RedPointConst.Rank_Power] = () => RedPointManager.Instance.GetRedNum(RedPointConst.Rank_Power.ToString()),

            [RedPointConst.Rank_ZhanHun] = () =>
            {
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false) || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false))
                    return 0;
                else
                    return RedPointManager.Instance.GetRedNum(RedPointConst.Rank_ZhanHun.ToString());
            },

            [RedPointConst.WarFight] = () => RedDotM.Instance.GetFightWarRedNum(),

            [RedPointConst.Role_Title_Active] = () => RedDotM.Instance.GetTitleActiveRedNum(),
            [RedPointConst.Role_Title_UpStar] = () => RedDotM.Instance.GetTitleUpStarNum(),

            [RedPointConst.ActiveCenter_StarWish] = () => RedDotM.Instance.GetStarWishOrCrazyRefineRedNum(ActivityType.StarWish),
            [RedPointConst.ActiveCenter_CrazyRefine] = () => RedDotM.Instance.GetStarWishOrCrazyRefineRedNum(ActivityType.CrazyRefine),
        };
    }
}
