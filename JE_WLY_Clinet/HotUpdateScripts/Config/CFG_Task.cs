using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using MetaFramework.UI;
using pb;
using System;
using System.Collections.Generic;

/// <summary>
/// 任务
/// </summary>
public class CFGItem_Task : CFGItemBase
{
    public long UUID { set; get; }
    public override long GetIdx()
    {
        return UUID;
    }

    public int Class;

    public int Type;

    public decimal SevenDay;

    public long Typeson;

    public long limit { set; get; }

    public Dictionary<string, int>[] Reward;

    public decimal NextAchievement;

    public decimal VIPLevel;

    public decimal ResetMark;

    public ActivityReward_Item ActivityReward;

    public class ActivityReward_Item
    {
        public int key;
        public int value;
    }

    public string LinkPanel;

    public string TaskDes;

    public LimitTimeItem LimitTime;

    public struct LimitTimeItem
    {
        public string start;
        public string end;
    }
}

/// <summary>
/// CFG_Task 配置表
/// </summary>
public class CFG_Task : CFGTableBase<CFGItem_Task>
{
    private static CFG_Task _AllCFG;
    private static CFG_Task AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Task, CFGItem_Task>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Task Get()
    {
        return AllCFG;
    }

    public static Dictionary<long, CFGItem_Task> GetDic()
    {
        return AllCFG.DicDatas;
    }


    public static CFGItem_Task GetById(long id)
    {
        return AllCFG.DicDatas[id];
    }

    //public static string StrTaskDes(long id)
    //{
    //    string str = GetById(id).TaskDes;

    //    if (str.Contains("{1}"))
    //    {
    //        if (GetById(id).Type == 123)
    //        {
    //            str = string.Format(str, GameUtil.Instance.FormatNumWithUnit((int)GetById(id).limit), GameUtil.Instance.GetItemName((int)GetById(id).Typeson));
    //        }
    //        else
    //            str = string.Format(str, CFG_Skill.GetById((int)GetById(id).Typeson).Name, GameUtil.Instance.FormatNumWithUnit((int)GetById(id).limit));
    //    }
    //    else if (str.Contains("{0}"))
    //    {
    //        if (GetById(id).Type == 151 || GetById(id).Type == 154)
    //        {
    //            CFGItem_Battle cFGItem = CFG_Battle.GetById((int)GetById(id).Typeson);
    //            str = string.Format(str, $"{cFGItem.Chapter}章{cFGItem.Layers}节");
    //        }
    //        else if (GetById(id).Type == 152 || GetById(id).Type == 153)
    //        {
    //            CFGItem_Battle cFGItem = CFG_Battle.GetById((int)GetById(id).Typeson);
    //            str = string.Format(str, $"{cFGItem.Layers}");
    //        }
    //        else if (GetById(id).Type == 163 || GetById(id).Type == 161)
    //        {
    //            str = string.Format(str, GameUtil.Instance.FormatNumWithUnit((int)GetById(id).Typeson));
    //        }
    //        else
    //            str = string.Format(str, GameUtil.Instance.FormatNumWithUnit((int)GetById(id).limit));
    //    }

    //    return str;
    //}

    public static string StrTaskDes(long id)
    {
        var taskConfig = GetById(id);

        if (taskConfig == null)
        {
            return ""; //没有知道到Task信息
        }

        switch (taskConfig.Type)
        {
            case 123:
                return string.Format(taskConfig.TaskDes, GameUtil.Instance.FormatNumWithUnit((int)GetById(id).limit), GameUtil.Instance.GetItemName((int)GetById(id).Typeson));
            case 151:
            case 154:
                CFGItem_Battle cFGItem = CFG_Battle.GetById((int)GetById(id).Typeson);
                return string.Format(taskConfig.TaskDes, $"{cFGItem.Chapter}章{cFGItem.Layers}节");
            case 152:
            case 153:
            case 197:
                CFGItem_Battle cFGLayerItem = CFG_Battle.GetById((int)GetById(id).Typeson);
                return string.Format(taskConfig.TaskDes, $"{cFGLayerItem.Layers}");
            case 161:
            case 163:
            case 200:
                return string.Format(taskConfig.TaskDes, GameUtil.Instance.FormatNumWithUnit((int)GetById(id).Typeson));
            case 182:
                var mainBattleInfo = CFG_MainBattle.GetById((int)taskConfig.Typeson);
                return string.Format(taskConfig.TaskDes, mainBattleInfo.name); //返回TaskDes和主线管卡Name的情况
            case 191:
                return taskConfig.TaskDes; //直接返回TaskDes的情况
        }

        return string.Format(taskConfig.TaskDes, GameUtil.Instance.FormatNumWithUnit((int)GetById(id).limit)); //处理直接把TaskDes和limit组合的情况
    }

    public static (int, int) GetTaskPro(long id, int Amount, int TargetNum)
    {
        string str = GetById(id).TaskDes;

        if (GetById(id).Type >= 151 && GetById(id).Type <= 154)
        {
            if (Amount == TargetNum)
            {
                return (1, 1);
            }
            else
                return (0, 1);
        }
        else
            return (Amount, TargetNum);
    }

    public static void JumpUI(int cfgId, Action jumpHandle = null, string linkPanel = "")
    {
        JumpUIByType((TargetType)GetById(cfgId).Type, jumpHandle, linkPanel);
    }

    public static void JumpUIByType(int taskType, Action jumpHandle = null, string linkPanel = "")
    {
        JumpUIByType((TargetType)taskType, jumpHandle, linkPanel);
    }

    public static void JumpUIByType(TargetType targetType, Action jumpHandle = null, string linkPanel = "")
    {
        switch (targetType)
        {
            #region  old Type
            case TargetType.TASK_TARGET_TYPE_INVALID:
                UITool.ShowTips("任务类型错误");
                break;
            case TargetType.ROLE_POWER_REQUIREMENT:

                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
                    return;

                jumpHandle?.Invoke();
                UITool.Open<UIBuild>();
                break;
            case TargetType.KILL_MONSTER_REQUIREMENT:
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }
                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>();
                break;
            case TargetType.ROLE_LEVEL_REQUIREMENT:
            case TargetType.UPGRADE_ANY_PARTNER_LEVEL: // 升级任意一个伙伴等级至X级
                //if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.ExpPool))
                //    return;

                jumpHandle?.Invoke();
                //UITool.Open<UIExpPool>();
                break;
            //case TargetType.PVE_INSTANCE_ENTRY_REQUIREMENT:
            //    if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Relice))
            //        return;

            //    if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
            //    {
            //        UITool.ShowTips("已在战斗中");
            //        return;
            //    }

            //    jumpHandle?.Invoke();
            //    UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.Relic);
            //    break;
            case TargetType.TOTAL_PHYSICAL_STRENGTH_CONSUMED_REQUIREMENT:
                break;
            //case TargetType.EQUIPMENT_STRENGTHENING_LEVEL_REQUIREMENT:
            //    if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
            //        return;
            //    if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
            //        return;

            //    jumpHandle?.Invoke();
            //    UITool.Open<UIBuild>();
            //    break;
            case TargetType.EQUIPMENT_UPGRADING_STAR_LEVEL_REQUIREMENT:
            case TargetType.ANY_EQUIP_STAR_UP_TIMES:// 任意装备升星X次
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar))
                    return;

                jumpHandle?.Invoke();
                UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Star);
                break;
            case TargetType.SKILL_UPGRADING_REQUIREMENT:
            case TargetType.SKILL_UPGRADING_TIMES:// 升级X次主角技能
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
                    return;

                jumpHandle?.Invoke();
                UITool.Open<UIBuild>(null, UIBuild.SelectOne.Skill, UIBuild.SelectTwo.Skill_UpLv);
                break;
            case TargetType.TOTAL_GOLD_COINS_CONSUMED_REQUIREMENT:
                jumpHandle?.Invoke();
                UITool.Open<UIShop>(null, CFG_Store.StoreType.Coin);
                break;
            case TargetType.TIME_RANGE:
                break;
            case TargetType.DAILY_ACTIVITY_LEVEL:
                break;
            case TargetType.WEEKLY_ACTIVITY_LEVEL:
                break;
            case TargetType.TOTAL_BOOST_PURCHASES:
                break;
            case TargetType.TOTAL_CHALLENGES_IN_ARENA:
            case TargetType.TOTAL_ACQUIRED_SOULS:
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                    return;

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.PvP);
                break;
            case TargetType.TOTAL_ATTEMPTS_IN_ENDLESS_TRIALS:

                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                UITool.Close<UIShop>();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.SLTower);
                break;
            case TargetType.TOTAL_KILLS_OF_NORMAL_MONSTERS:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>();
                break;
            case TargetType.TOTAL_KILLS_OF_BOSS_MONSTERS:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss))
                    return;

                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.Boss);
                break;
            case TargetType.TOTAL_SHOP_PURCHASES_WITH_GOLD:
                jumpHandle?.Invoke();
                UITool.Open<UIShop>(null, CFG_Store.StoreType.Coin);
                break;
            case TargetType.TOTAL_SHOP_PURCHASES_WITH_COUPONS:
                jumpHandle?.Invoke();
                UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                break;
            case TargetType.TOTAL_SHOP_PURCHASES_WITH_SOULS:
                jumpHandle?.Invoke();
                UITool.Open<UIShop>(null, CFG_Store.StoreType.Diamanta);
                break;
            case TargetType.TOTAL_SHOP_PURCHASES_WITH_TRIAL_TOKENS:
                jumpHandle?.Invoke();
                UITool.Open<UIShop>(null, CFG_Store.StoreType.SiLian);
                break;
            case TargetType.TOTAL_ACQUIRED_TRIAL_TOKENS:
                break;
            case TargetType.TOTAL_ITEM_CONSUMPTION:
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Bag))
                    return;

                TcpC.GetBagInfo();
                UITool.Open<UIBag>();
                break;
            case TargetType.TOTAL_BOX_OPENED:
                jumpHandle?.Invoke();
                break;
            case TargetType.BATTLE_CHAPTERS_CLEARED:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>();
                break;
            case TargetType.BOSS_LEVELS_CLEARED:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.Boss);
                break;
            case TargetType.TRIAL_TOWER_FLOORS_CLEARED:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower))
                    return;

                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.SLTower);
                break;
            case TargetType.PVE_INSTANCE_ENTRY_REQUIREMENT:
            case TargetType.RELIC_CHAPTERS_CLEARED:
            case TargetType.RELIC_TOTAL_STAR://遗迹
            case TargetType.TOTAL_RELIC_CHALLENGE:

                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Relice))
                    return;

                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Relice))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.Relic);
                break;
            case TargetType.VIP_LEVEL:
                jumpHandle?.Invoke();
                UITool.Open<UIVIP>();
                break;
            case TargetType.TOTAL_MATERIAL_INSTANCE_CHALLENGE:
            case TargetType.TOTAL_EXPERIENCE_FB_CLEARED: // 通关X次经验副本
            case TargetType.TOTAL_GOLD_FB_CLEARED: // 通关X次经验副本
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_QHSAndGold))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_QHSAndGold))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.FuBen, UIFightChoice.ChanllengeType.DayDo);
                break;
            case TargetType.TOTAL_MULTI_VS_MULTI_ARENA_CHALLENGE:

                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                    return;

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.PvP);
                break;
            case TargetType.TOTAL_BOSS_CHALLENGE:

                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss))
                    return;

                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }

                jumpHandle?.Invoke();
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.Boss);
                break;
            case TargetType.TASK_TARGET_TYPE_MAX:
                break;

            case TargetType.UNLOCK_PARTNER_SLOT: // 解锁第X个伙伴栏
                jumpHandle?.Invoke();
                UITool.Open<UIRole>(null, -2);
                break;

            case TargetType.TOTAL_PET_EGGS_OPENED: // 累计打开X个宠物蛋
            case TargetType.ANY_PET_UPGRADED_TO_LEVEL: // 任意宠物升至X级
                {
                    if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet))
                        return;
                    if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet))
                        return;
                    jumpHandle?.Invoke();
                    UITool.Open<UIPet>();
                }
                break;
            case TargetType.EQUIPMENT_STRENGTHENING_LEVEL_REQUIREMENT://代表全身装备强化等级需求
            case TargetType.ANY_EQUIP_STRENGTHENED_TO_LEVEL: // 任意装备格强化至X
            //case TargetType.ANY_EQUIP_STRENGTHENED_TO_LEVEL_Count://任意装备格强化N次
            case TargetType.TOTAL_STRENGTHEN_TIMES: // 装备强化
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build))
                    return;

                jumpHandle?.Invoke();
                UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Reinforce);
                break;

            case TargetType.CHEST_LEVEL_UPGRADED: // 宝箱等级升至X级
                jumpHandle?.Invoke();
                UITool.Open<UIBoxUpgrade>();
                break;
            #endregion
            /* Kanxf add new type*/
            case TargetType.SEVENDAY_TASK_DAY_ALL_COMPLETED:
            case TargetType.SEVENDAY_TASK_ALL_COMPLETED:
            case TargetType.SEVENDAY_TASK_GIFT_PURCHASED:
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.SignDay))
                    return;

                UITool.Open<UIActive>();
                break;

            case TargetType.MAIN_BATTLE:
                var uihall = UITool.Get<UIHall>();
                if (uihall != null)
                {
                    uihall.OnClickAttackButton();
                    jumpHandle?.Invoke();
                }
                break;
            case TargetType.PARTNER_RECEIVE://抽取伙伴
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.CallRole))
                    return;

                UITool.Open<UIPartnerCall>();
                break;
            case TargetType.WORLD_BOSS_PARTICIPATION://世界BOSS
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.WorldBoss);
                break;
            case TargetType.GUARD_ISLAND_INSTANCE://挂机副本
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.DefendIsland))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.HaiDao);
                break;
            case TargetType.PET_COMPOSE_TIMES: //宠物合成
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet))
                    return;

                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIPet>();
                });
                break;
            case TargetType.TURN_BASED_PK_TIMES: //回合制比赛
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.PvP);
                break;
            case TargetType.SOUL_STONE_INSTANCE: //魂石副本
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_HS))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.FuBen, UIFightChoice.ChanllengeType.DayDo);
                break;
            case TargetType.PARTNER_ON_STAGE_TIMES: //伙伴上阵
                UITool.Open<UIRole>();
                break;
            case TargetType.SKILL_ADVANCED_TIMES://进阶技能
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillStar))
                    return;

                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Skill, UIBuild.SelectTwo.Skill_Advance);
                });
                break;
            case TargetType.JOIN_GUILD: //行会
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Guild))
                    return;

                UITool.Open<UIGuild>();
                break;
            case TargetType.NORMAL_REFINE_TIMES: //普通洗炼
            case TargetType.SPECIAL_REFINE_TIMES: //特殊洗炼
                //if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Guild))
                //    return;
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Refined))
                    return;
                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Refine);
                });
                break;
            case TargetType.STAR_UP_STONE_INSTANCE://升星石副本
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_ExpAndSXS))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.FuBen, UIFightChoice.ChanllengeType.DayDo);
                break;
            case TargetType.BOOK_PAGE_INSTANCE: //书页副本
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SY))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.FuBen, UIFightChoice.ChanllengeType.DayDo);
                break;
            case TargetType.ENHANCE_STONE_INSTANCE://强化石副本
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_QHSAndGold))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.FuBen, UIFightChoice.ChanllengeType.DayDo);
                break;
            case TargetType.ENDLESS_TRIALS_STAGE: //一战到底
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.FightEnd);
                break;
            case TargetType.FIVE_ELEMENTS_INSTANCE: //五行副本
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WuXing))
                    return;

                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.WuXing);
                break;
            case TargetType.GEMSTONE_INLAY_TIMES: //镶嵌宝石
                //if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillStar))
                //    return;

                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Inset);
                });
                break;
            case TargetType.EQUIP_STAR_TO_X: //装备升星
                if (systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar))
                    return;

                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Star);
                });
                break;
            default:
                break;
        }
    }

    private static bool systemIsLock(CFG_XMLConfig.XML_Config.SwitchOpne type)
    {
        if (GameUtil.Instance.CheckOpenServerIsLock(type))
            return true;
        if (GameUtil.Instance.CheakIsLock(type))
            return true;

        return false;
    }
}