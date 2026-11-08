using JEngine.Core;
using JEngine.Helper;
using JEngine.Singleton;

using MetaFramework.UI;

using pb;

using System;
using System.Collections.Generic;
using System.Linq;

namespace HotUpdateScripts.ProtoM
{
    public class RedDotM : SingletonTemplate<RedDotM>
    {
        //VIP
        public int GetVIPRedNum(int vipLv = -1)
        {
            vipLv = vipLv == -1 ? RoleM.Instance.Vip : vipLv;
            if (vipLv <= 0 || RoleM.Instance.Vip < vipLv)
                return 0;

            var StoreList = CFG_Store.GetListByStoreType(CFG_Store.StoreType.VIPGift);

            for (int i = 0; i < vipLv; i++)
            {
                if (GameUtil.Instance.ComCoinToCheck(StoreList[i].ID, null, false) && StoreM.Instance.CheckInventory(StoreList[i].ID))
                {
                    return 1;
                }
            }
            return 0;
        }


        public int TurnWheelId = -1;
        //点券大转盘
        public int GetTurnWheelRedNum()
        {
            if (!ActivityM.Instance.GetAcitvityState(ActivityType.TurnWheel) || TurnWheelId == -1)
            {
                return 0;
            }
            var CFGAct = CFG_Activity.GetById(TurnWheelId);
            if (CFGAct.PreLimit > RoleM.Instance.Level)
            {
                return 0;
            }
            if (!GameUtil.Instance.ComCoinToCheck(int.Parse(CFGAct.PayItem[0].ElementAt(0).Key), CFGAct.PayItem[0].ElementAt(0).Value, null, false))
            {
                return 0;
            }
            return 1;
        }

        //单个装备位的强化红点
        public int GetOneEquipSlotReinforceRedNum(EquipSlot equipSlot, int TeamIdx)
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build, false))
                return 0;

            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build, false))
                return 0;

            //if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build])
            //    return 0;

            if (TeamIdx != 0 && (PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx) == null || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid == 0))
                return 0;

            if (TeamIdx == 0 && !RoleM.Instance.isEquip((int)equipSlot))
                return 0;

            var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(equipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)equipSlot);

            int lv = TeamIdx == 0 ? RoleM.Instance.Level : PartnerM.Instance.GetOnePartnerInfoByUid(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid).Level;

            if (_SoltInfo.Enhance >= lv * CFG_XMLConfig.xmlConfig.BaseValue.equipStrenghthen)
                return 0;
            var _cfgInfo = CFG_EquipStrengthen.GetByJobAndLv(RoleM.Instance.Job, _SoltInfo.Enhance);

            for (int i = 0; i < _cfgInfo.Consume.Length; i++)
            {
                if (!GameUtil.Instance.ComCoinToCheck(_cfgInfo.Consume[i].id, _cfgInfo.Consume[i].num, null, false))
                {
                    return 0;
                }
            }
            return 1;
        }

        //单个伙伴装备位的强化红点
        public int GetOneTeamAllEquipSlotReinforceRedNum(int TeamIdx)
        {
            for (int i = (int)EquipSlot.Weapon; i < (int)EquipSlot.MaxEquipSlot; i++)
            {
                if (GetOneEquipSlotReinforceRedNum((EquipSlot)i, TeamIdx) > 0)
                    return 1;
            }
            return 0;
        }

        //所有伙伴装备位的强化红点
        public int GetAllEquipSlotReinforceRedNum()
        {
            for (int i = 0; i < 5; i++)
            {
                if (GetOneTeamAllEquipSlotReinforceRedNum(i) > 0)
                {
                    return 1;
                }
            }
            return 0;
        }

        //装备 升星
        public int GetEquipStarRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar, false))
                return 0;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar, false))
                return 0;
            //if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar])
            //    return 0;
            return RedPointManager.Instance.GetRedNum(RedPointConst.Equip_Star.ToString());
        }

        // 装备 洗炼
        public int GetEquipRefineNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Refined, false))
                return 0;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Refined, false))
                return 0;
            //if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipRefine])
            //    return 0;
            return RedPointManager.Instance.GetRedNum(RedPointConst.Equip_Refine.ToString());
        }

        //装备 单个红点镶嵌
        public int GetOneEquipInsetRedNum(int Solt, int TeamIdx)
        {
            int roleLv = RoleM.Instance.Level;
            var xmlConfig = CFG_XMLConfig.xmlConfig;
            var switchOpenLvDic = xmlConfig.SwitchOpenLvDic;

            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster, false))
                return 0;

            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster, false))
                return 0;
            //if (roleLv < switchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster])
            //    return 0;

            PartnerSlotInfo partnerSlot = null;
            if (TeamIdx != 0)
            {
                partnerSlot = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
                if (partnerSlot == null || partnerSlot.PartnerUid == 0)
                    return 0;
            }

            int[] gemstone = TeamIdx == 0
                ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)Solt).Gemstone
                : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, Solt).Gemstone;

            var diamondOpenLvDic = xmlConfig.Diamond_OpenlvDic;
            List<int> cfgids = CFG_Item.GetCfgIdsByEquipSolt(Solt);

            for (int i = 0; i < 4; i++)
            {
                //bool isOpen = diamondOpenLvDic[i + 1] <= roleLv;
                //if (isOpen && BsedInset(Solt, Gemstone[i]))
                //{
                //    return 1;
                //}

                if (diamondOpenLvDic[i + 1] > roleLv)
                    continue;

                int id = gemstone[i];
                int baseLv = id == 0 ? 0 : (int)CFG_Item.GetById(id).itemlevel;
                for (int j = 0; j < cfgids.Count; j++)
                {
                    int cfgid = cfgids[j];
                    if ((int)CFG_Item.GetById(cfgid).itemlevel > baseLv && BagM.Instance.GetNumByCfgID(cfgid) > 0)
                        return 1;
                }
            }
            return 0;
        }

        //private static bool BsedInset(int Solt, int id)
        //{
        //    int baseLv = id == 0 ? 0 : (int)CFG_Item.GetById(id).itemlevel;
        //    List<int> cfgids = CFG_Item.GetCfgIdsByEquipSolt(Solt);
        //    for (int i = 0; i < cfgids.Count; i++)
        //    {
        //        int cfgid = cfgids[i];
        //        if ((int)CFG_Item.GetById(cfgid).itemlevel > baseLv && BagM.Instance.GetNumByCfgID(cfgid) > 0)
        //            return true;
        //    }
        //    return false;
        //}

        //装备 单个成员 单个红点镶嵌
        public int GetOneAllEquipInsetRedNum(int TeamIdx)
        {
            for (int i = (int)EquipSlot.Weapon; i < (int)EquipSlot.MaxEquipSlot; i++)
            {
                if (GetOneEquipInsetRedNum(i, TeamIdx) > 0)
                    return 1;
            }
            return 0;
        }

        //所有装备位的 红点镶嵌
        public int GetAllEquipInsetRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster, false))
                return 0;


            //if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster])
            //    return 0;

            for (int i = 0; i < AllTeamCount; i++)
            {
                if (GetOneAllEquipInsetRedNum(i) > 0)
                {
                    return 1;
                }
            }
            return 0;
        }

        //技能 检测单个技能升级红点
        public int GetOneSkillUpLvRedNum(int skillId, int TeamIdx, bool needBase = true)
        {
            if (needBase)
            {
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv, false))
                    return 0;

                if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv])
                    return 0;

                if (TeamIdx != 0 && (PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx) == null || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid == 0))
                    return 0;
            }

            var skillInfo = TeamIdx == 0 ? SkillM.Instance.GetSkillInfoBySkillID(skillId) : PartnerM.Instance.GetOneSkillInfo(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid, skillId);
            var cfgInfo = CFG_Skill.GetById(skillInfo.CfgId);
            var nextCfgInfo = CFG_Skill.GetNextLvByUid(skillInfo.CfgId);

            int lv = TeamIdx == 0 ? RoleM.Instance.Level : PartnerM.Instance.GetOnePartnerInfoByUid(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid).Level;

            if (nextCfgInfo.SkillLevel == cfgInfo.SkillLevel || nextCfgInfo.RoleLevellimit > lv)
                return 0;

            for (int i = 0; i < cfgInfo.PromoteConsume.Length; i++)
            {
                if (!GameUtil.Instance.ComCoinToCheck(cfgInfo.PromoteConsume[i].id, cfgInfo.PromoteConsume[i].num, null, false))
                {
                    return 0;
                }
            }
            return 1;
        }

        //技能 单个伙伴检测所有技能升级红点
        public int GetOneTeamAllSkillUpLvRedNum(int TeamIdx, bool needBase = true)
        {

            if (needBase && RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv])
                return 0;

            if (TeamIdx != 0 && (PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx) == null || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid == 0))
                return 0;

            var skillList = TeamIdx == 0 ? SkillM.Instance.GetSkillInfoList() : PartnerM.Instance.GetSkillInfoByPartnerId(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid);
            for (int i = 0; i < skillList.Count; i++)
            {
                var cfgInfo = CFG_Skill.GetById(skillList[i].CfgId);
                if (cfgInfo.Skilltype != (int)CFG_Skill.SkillType.Normall && cfgInfo.ActiveorPassive == 1 && GetOneSkillUpLvRedNum(cfgInfo.SkillID, TeamIdx, false) > 0)
                    return 1;
            }
            return 0;
        }

        /// <summary>
        /// 技能 检测技能设置
        /// </summary>
        /// <returns></returns>
        public int GetSkillSettingRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv, false))
                return 0;
            if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv])
                return 0;
            return RedPointManager.Instance.GetRedNum(RedPointConst.Skill_Setting.ToString());
        }

        int AllTeamCount = 5;
        //技能 检测所有技能升级红点
        public int GetAllSkillUpLvRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv, false))
                return 0;

            if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv])
                return 0;

            for (int i = 0; i < AllTeamCount; i++)
            {
                if (GetOneTeamAllSkillUpLvRedNum(i, false) > 0)
                {
                    return 1;
                }
            }
            return 0;
        }

        //技能 检测单个技能升阶红点
        public int GetOneSkillAdvanceRedNum(int skillId, int TeamIdx, bool needBase = true)
        {
            if (needBase)
            {
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillStar, false))
                    return 0;

                if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillStar])
                    return 0;

                if (TeamIdx != 0 && (PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx) == null || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid == 0))
                    return 0;
            }

            var skillInfo = TeamIdx == 0 ? SkillM.Instance.GetSkillInfoBySkillID(skillId) : PartnerM.Instance.GetOneSkillInfo(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid, skillId);
            var cfgInfo = CFG_Skill.GetById(skillInfo.CfgId);
            var nextCfgInfo = CFG_Skill.GetNextLvByUid(skillInfo.CfgId);

            var nowAdvCfgInfo = CFG_Skill.GetAdvanceByUidAndAdvan(skillInfo.CfgId, skillInfo.Advanced);
            var nextAdvCfgInfo = CFG_Skill.GetAdvanceByUidAndAdvan(skillInfo.CfgId, skillInfo.Advanced + 1);
            nextAdvCfgInfo = nextAdvCfgInfo == null ? nowAdvCfgInfo : nextAdvCfgInfo;

            if (skillInfo.Advanced == nextAdvCfgInfo.CanupperLevel || cfgInfo.SkillLevel < nextAdvCfgInfo.SkillLevel)
                return 0;

            for (int i = 0; i < nextAdvCfgInfo.UpperlevelConsume.Length; i++)
            {
                if (!GameUtil.Instance.ComCoinToCheck(nextAdvCfgInfo.UpperlevelConsume[i].id, nextAdvCfgInfo.UpperlevelConsume[i].num, null, false))
                {
                    return 0;
                }
            }
            return 1;
        }

        //技能 单个伙伴检测所有技能升阶红点
        public int GetOneTeamAllSkillAdvanceRedNum(int TeamIdx, bool needBase = true)
        {
            if (needBase && RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillStar])
                return 0;

            if (TeamIdx != 0 && (PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx) == null || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid == 0))
                return 0;

            var skillList = TeamIdx == 0 ? SkillM.Instance.GetSkillInfoList() : PartnerM.Instance.GetSkillInfoByPartnerId(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid);
            for (int i = 0; i < skillList.Count; i++)
            {
                var cfgInfo = CFG_Skill.GetById(skillList[i].CfgId);
                if (cfgInfo.Skilltype != (int)CFG_Skill.SkillType.Normall && cfgInfo.ActiveorPassive == 1 && GetOneSkillAdvanceRedNum(cfgInfo.SkillID, TeamIdx, false) > 0)
                    return 1;
            }
            return 0;
        }

        //技能 检测所有技能升阶红点
        public int GetAllSkillAdvanceRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillStar, false))
                return 0;

            if (RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillStar])
                return 0;

            for (int i = 0; i < AllTeamCount; i++)
            {
                if (GetOneTeamAllSkillAdvanceRedNum(i, false) > 0)
                {
                    return 1;
                }
            }
            return 0;
        }

        //排行红点
        public int GetRandRedNumByType(string key)
        {
            string timeStr = GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear.ToString();
            return PlayerPrefsExtend.GetString(key + RoleM.Instance.Uid) != timeStr ? 1 : 0;
        }

        //排行红点关闭
        public void SetRandRedNumByType(string key, Action callBack)
        {
            string timeStr = GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear.ToString();
            PlayerPrefsExtend.SetString(key + RoleM.Instance.Uid, timeStr);
            callBack?.Invoke();
        }

        //世界boss
        //public int OnWorldBossRedDot(int num)
        //{
        //    //if (CFG_XMLConfig.CheakLock(CFG_XMLConfig.XML_Config.SwitchOpne.World_Boss))
        //    //    return 0;
        //    if (!GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.WorldBoss.startTime, CFG_XMLConfig.xmlConfig.WorldBoss.endTime))
        //        return 0;

        //    return CFG_XMLConfig.xmlConfig.WorldBoss.enterTime - num;
        //}

        // 每日必做
        public int OnDayDoRedNum()
        {
            foreach (var item in UIDayDo.RoomToLockDic)
            {
                if (!GameUtil.Instance.CheckOpenServerIsLock(item.Value, false) && !GameUtil.Instance.CheakIsLock(item.Value, false) && RoleM.Instance.GetCanTimeByRoomType((int)item.Key) > 0)
                {
                    return 1;
                }
            }
            return 0;
        }

        public int GetWuXingChallengeNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WuXing, false) || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WuXing, false))
                return 0;

            return TurnBasedM.Instance.elementTimes;
        }

        #region 角色红点

        public int GetOneTeamBeatEquipRed(int TeamIdx)
        {
            if (TeamIdx != 0)
            {
                if (CFG_XMLConfig.GetOpenSoltLv(TeamIdx) > RoleM.Instance.Level)
                    return 0;
                var soltInfo = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
                if (soltInfo.PartnerUid == 0)
                    return 0;
            }
            var Equips = BagM.Instance.GetEquipItemList();
            for (int i = 0; i < Equips.Count; i++)
            {
                (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(Equips[i]);
                if (powerNum > 0 && soltIdx == TeamIdx)
                    return 1;
            }
            return 0;
        }

        public int GetBeatEquipRed()
        {
            var Equips = BagM.Instance.GetEquipItemList();
            for (int i = 0; i < Equips.Count; i++)
            {
                (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(Equips[i]);
                if (powerNum > 0)
                    return 1;
            }
            return 0;
        }

        //伙伴上阵检测
        public int PartnerIntoRedDot()
        {
            int SoltCount = CFG_XMLConfig.GetOpenSolt(RoleM.Instance.Level);

            for (int i = 1; i <= SoltCount; i++)
            {
                if (OnePartnerIntoRedDot(i) > 0 && CFG_XMLConfig.GetOpenSolMainStaget(i) <= RoleM.Instance.MainBattleId)
                {
                    return 1;
                }
            }
            return 0;
        }

        public int OnePartnerIntoRedDot(int TeamIdx)
        {
            if (CFG_XMLConfig.CheckPartnerSoltIsOpenServer(TeamIdx, false) || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx) == null || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid != 0)
                return 0;
            var AllPartners = PartnerM.Instance.GetPartners();
            for (int i = 0; i < AllPartners.Count; i++)
            {
                if (!PartnerM.Instance.IsInSlot(AllPartners[i].Uid))
                {
                    return 1;
                }
            }
            return 0;
        }

        //伙伴召唤红点
        public int GetPartnerCallRedDot()
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.CallRole, false))
                return 0;

            if (CFG_XMLConfig.GetOpenSolt(RoleM.Instance.Level) <= 0)
                return 0;

            //for (int i = 0; i < UIPartnerCall.UsedItemIDs.Length; i++)
            //{
            //    if (BagM.Instance.GetNumByCfgID(UIPartnerCall.UsedItemIDs[i]) > 0 && PartnerM.Instance.CanGetNums[i] > 0)
            //        return 1;
            //}
            return 0;
        }

        /// <summary>
        /// 伙伴升星红点
        /// </summary>
        /// <param name="TeamIdx"></param>
        /// <returns></returns>
        public int GetOnePartnerUpStarRedDot(int TeamIdx)
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Partner_Star, false))
                return 0;

            var partners = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
            if (partners.PartnerUid != 0)
            {
                var partner = PartnerM.Instance.GetOnePartnerInfoByUid(partners.PartnerUid);
                var cfgInfo = CFG_Partner.GetById(partner.CfgId);
                if (partner.Star != cfgInfo.StarConsume.Length)
                {
                    var userdData = cfgInfo.GetStarConsumeByLv(partner.Star + 1);

                    bool isCanUp = true;
                    for (int j = 0; j < userdData.Length; j++)
                    {
                        if (!GameUtil.Instance.ComCoinToCheck(userdData[j].key, userdData[j].value, null, false))
                        {
                            isCanUp = false;
                            break;
                        }
                    }
                    if (isCanUp)
                        return 1;
                }
            }
            return 0;
        }

        /// <summary>
        /// 伙伴升星红点
        /// </summary>
        /// <returns></returns>
        public int GetPartnerUpStarRedDot()
        {
            var partners = PartnerM.Instance.GetPartnerSlotInfo();
            for (int i = 1; i <= 4; i++)
            {
                if (GetOnePartnerUpStarRedDot(i) > 0)
                    return 1;
            }
            return 0;
        }

        #endregion

        #region 宠物

        private int[] PetUsedIds = { 80, 81, 82 };
        public int GetPetAwardNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet, false))
                return 0;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet, false))
                return 0;
            for (int i = 0; i < PetUsedIds.Length; i++)
            {
                if (BagM.Instance.GetNumByCfgID(PetUsedIds[i]) > 0 && PetM.Instance.CanGetNums[i] > 0)
                    return 1;
            }
            return 0;
        }

        public int GetPetUpStar()
        {
            var Pet = PetM.Instance.GetOnePetByIdx(0);
            if (Pet.CfgId == 0)
                return 0;

            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet, false))
                return 0;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet, false))
                return 0;

            var cfgInfo = CFG_Pet.GetById(Pet.CfgId);
            int nextBreak = Pet.Star == cfgInfo.StarConsume.Length ? Pet.Star : (Pet.Star + 1);
            if (Pet.Star == nextBreak)
                return 0;
            if (Pet.ExpStar != cfgInfo.GetStarNeedAllPro(nextBreak))
            {
                for (int i = 0; i < CFG_XMLConfig.xmlConfig.PetAdvanceUseds.Length; i++)
                {
                    if (BagM.Instance.GetNumByCfgID(CFG_XMLConfig.xmlConfig.PetAdvanceUseds[i].id) > 0)
                        return 1;
                }
            }
            return 0;
        }
        #endregion


        /// <summary>
        /// 经验池红点
        /// </summary>
        /// <returns></returns>
        public int GetExpPoolRedNum()
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.ExpPool, false))
                return 0;

            bool isPartner = false;
            for (int i = 1; i <= 4; i++)
            {
                var soltInfo = PartnerM.Instance.GetOnePartnerSlotInfo(i);
                if (soltInfo != null && soltInfo.PartnerUid != 0)
                {
                    isPartner = true;

                    var partner = PartnerM.Instance.GetOnePartnerInfoByUid(soltInfo.PartnerUid);
                    if (partner == null)
                        continue;

                    var cfg_partner = CFG_Partner.GetById(partner.CfgId);
                    if (CFG_Role.GetByJobAndLevel(cfg_partner.job, partner.Level).NeedEXP <= RoleM.Instance.Exp)
                    {
                        return 1;
                    }
                }
            }
            if (isPartner)
            {
                if (CFG_Role.GetByJobAndLevel(RoleM.Instance.Job, RoleM.Instance.Level).NeedEXP <= RoleM.Instance.Exp)
                {
                    return 1;
                }
            }
            return 0;
        }

        /// <summary>
        /// 获取合成红点数量
        /// </summary>
        /// <returns></returns>
        public int GetSynRedNum()
        {
            if (RedPointManager.Instance.GetRedNum(RedPointConst.Bag_Syn.ToString()) == 0)
                return 0;

            int[] synTypes = new int[] { 1, 3 };

            for (int i = 0; i < synTypes.Length; i++)
            {
                var synList = BagM.Instance.GetSynListByType(synTypes[i], null);
                if (synList.Count > 0)
                {
                    for (int j = 0; j < synList.Count; j++)
                    {
                        if (synList[j].Num >= 3)
                            return 1;
                    }
                }
            }
            return 0;
        }

        #region 排行

        public int GetRankHuiHe()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                return 0;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                return 0;

            return RedPointManager.Instance.GetRedNum(RedPointConst.Rank_HuiHe.ToString());
        }

        #endregion

        /// <summary>
        /// 战场红点
        /// </summary>
        /// <returns></returns>
        public int GetFightWarRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight, false)
               || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight, false))
                return 0;
            return (FightPlayM.Instance.CheakWarFightOpenTime() && RedPointManager.Instance.GetRedNum(RedPointConst.WarFight.ToString()) == 1) ? 1 : 0;
        }

        #region 称号
        public int GetOneTitleAtiveRedNum(int titleId, CFGItem_Title.JosnItem2[] josnItem2s = null, bool checkUpStar = false)
        {
            if (checkUpStar)
            {
                if (!TitleM.Instance.CheckTitleIsGet(titleId))
                    return 0;

                if (TitleM.Instance.GetTitleStar(titleId) >= CFG_XMLConfig.xmlConfig.WarFight.maxStar)
                    return 0;
            }
            else
            {
                if (TitleM.Instance.CheckTitleIsGet(titleId))
                    return 0;
            }

            var CFG_TitileActive = josnItem2s != null ? CFG_Title.GetById(titleId).Activity : josnItem2s;
            for (int i = 0; i < CFG_TitileActive.Length; i++)
            {
                var Iteminfo = CFG_TitileActive[i].values[0];
                long itemNum = BagM.Instance.GetNumByCfgID(Iteminfo.key);
                if (itemNum < Iteminfo.value)
                {
                    return 0;
                }
            }
            return 1;
        }

        public int GetOneTeamTitleActiveRedNum(int _TeamIdx, bool checkUpStar)
        {
            if (_TeamIdx > 0 && (CFG_XMLConfig.CheckPartnerSoltIsOpenServer(_TeamIdx, false) || PartnerM.Instance.GetOnePartnerSlotInfo(_TeamIdx) == null || PartnerM.Instance.GetOnePartnerSlotInfo(_TeamIdx).PartnerUid == 0))
                return 0;

            var cFGItem_Titles = CFG_Title.GetActiveListByType(_TeamIdx + 1);
            for (int i = 0; i < cFGItem_Titles.Count; i++)
            {
                if (GetOneTitleAtiveRedNum(cFGItem_Titles[i].UID, cFGItem_Titles[i].Activity, checkUpStar) > 0)
                    return 1;
            }
            return 0;
        }

        public int GetTitleActiveRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Title, false)
                || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Title, false))
                return 0;

            for (int i = 0; i < 5; i++)
            {
                if (GetOneTeamTitleActiveRedNum(i, false) > 0)
                    return 1;
            }
            return 0;
        }

        public int GetTitleUpStarNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Title, false)
                || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Title, false))
                return 0;

            for (int i = 0; i < 5; i++)
            {
                if (GetOneTeamTitleActiveRedNum(i, true) > 0)
                    return 1;
            }
            return 0;
        }


        #endregion

        #region 活动

        /// <summary>
        /// 活动红点
        /// </summary>
        /// <param name="activityType"></param>
        /// <returns></returns>
        public int GetActivityRedNum(ActivityType activityType)
        {
            int redNum = 0;

            switch (activityType)
            {
                case ActivityType.StarWish:
                case ActivityType.CrazyRefine:
                    redNum = GetStarWishOrCrazyRefineRedNum(activityType);
                    break;
                default:
                    break;
            }
            return redNum;
        }

        /// <summary>
        /// 星愿洗练、疯狂洗练 红点
        /// </summary>
        /// <param name="activityType"></param>
        /// <returns></returns>
        public int GetStarWishOrCrazyRefineRedNum(ActivityType activityType)
        {
            if (ActivityM.Instance.GetAcitvityEndTime(activityType) <= 0)
                return 0;

            var storeType = activityType switch
            {
                ActivityType.StarWish => CFG_Store.StoreType.StarWish,
                ActivityType.CrazyRefine => CFG_Store.StoreType.CrazyRefine,
                _ => CFG_Store.StoreType.No
            };

            var cfgStoreItems = CFG_Store.GetListByStoreType(storeType);
            int num = 0;
            for (int i = 0; i < cfgStoreItems.Count; i++)
            {
                int canbuy = StoreM.Instance.GetCanBuyNum(cfgStoreItems[i].ID);
                if (canbuy > 0)
                {
                    if (cfgStoreItems[i].Limite.Num == 0)
                        num = 1;

                    break;
                }
            }

            Log.Print($"red: {activityType}, {num}");
            return num;
        }
        #endregion
    }


}

