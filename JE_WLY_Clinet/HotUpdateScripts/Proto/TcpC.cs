using JEngine.MyNet;
using MetaFramework.UI;
using pb;

using System.Collections.Generic;

using UnityEngine;

namespace HotUpdateScripts.ProtoM
{
    //TCP
    public class TcpC
    {
        //角色升级
        public static void C2S_UpgradeLevel(int idx, long Uid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_UpgradeLevel, MType.Request,
                TCPClientMgr.ProtoSerialize(new UpgradeLevelRequest
                {
                    Type = idx == 0 ? ObjectType.PLAYER : ObjectType.PARTNER,
                    Uid = Uid,
                }));
        }

        #region PVP
        //pvp信息
        public static void OnPVPInfo(int _type)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PVPInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new PVPInfoRequest
                {
                    MatchType = (MatchType)_type
                }));
        }

        /// <summary>
        /// PVP 匹配对手列表
        /// </summary>
        public static void C2S_PvpMatching(int matchType)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PvpMatchingRequest, MType.Request,
                TCPClientMgr.ProtoSerialize(new PvpMatchingRequest
                {
                    MatchType = matchType
                }));
        }

        /// <summary>
        /// PVP挑战对手
        /// </summary>
        public static void C2S_PVPChallenge(long Uid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PVPChallengeRequest, MType.Request,
                TCPClientMgr.ProtoSerialize(new PVPChallengeRequest { Uid = Uid }));
        }

        public static void C2S_MatchingRequest(int Roomtype, bool Cancel)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_MatchingRequest, MType.Request,
                TCPClientMgr.ProtoSerialize(new MatchingRequest
                {
                    MatchType = (MatchType)Roomtype,
                    Cancel = Cancel
                }));
        }

        public static void C2S_PvpRecorderRequest(int battletype)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PvpRecoder, MType.Request,
                TCPClientMgr.ProtoSerialize(new PvpRecorderRequest
                {
                    Type = battletype,
                }));
        }
        #endregion

        //爬塔信息 1预览 ,2领奖
        public static void GetTrainInfo(int _type)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_EndlesstrialsReward, MType.Request,
                TCPClientMgr.ProtoSerialize(new EndlesstrialsRewardRequest
                {
                    Type = _type
                }));
        }

        /// <summary>
        /// 群雄战场
        /// </summary>
        public static void Get_C2S_WarRoomList()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_WarRoomList, MType.Request,
                TCPClientMgr.ProtoSerialize(new WarRoomListRequest()));
        }

        public static void GetGameJoin(int roomId, int _battld)
        {
            GameRoomM.Instance._GameState = GameRoomM.GameState.No;
            byte[] bytes = TCPClientMgr.ProtoSerialize(new JoinDeskRequest { RoomId = roomId, BattleId = _battld });
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Join, MType.Request, bytes);
        }


        public static void C2S_ItemUse(long Uid, int UsedNum, int SelectItemID = 0)
        {
            if (UsedNum > CFG_XMLConfig.xmlConfig.BaseValue.Maxuselimit)
            {
                UsedNum = CFG_XMLConfig.xmlConfig.BaseValue.Maxuselimit;
                UITool.ShowTips("单次使用上限为：" + UsedNum);
            }
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ItemUse, MType.Request,
                    TCPClientMgr.ProtoSerialize(new ItemUseRequest()
                    {
                        ItemUid = Uid,
                        Num = UsedNum,
                        SelectItemID = SelectItemID,
                    })
                    );
        }
        //获取背包信息
        public static void GetBagInfo()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GetBag, MType.Request, TCPClientMgr.ProtoSerialize(new GetBagRequest()));
        }

        //合成
        public static void SetCompose(List<ComposeInfo> list)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Compose, MType.Request, TCPClientMgr.ProtoSerialize(new ComposeRequest { Info = list }));
        }

        //魂石镶嵌
        public static void SetGemstoneChange(GemstoneChangeRequest info)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GemstoneChange, MType.Request, TCPClientMgr.ProtoSerialize(info));
        }

        //一键魂石镶嵌
        public static void SetAllGemstoneChange(EquipSlot slot, int TeamIdx)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GemstoneChangeAll, MType.Request, TCPClientMgr.ProtoSerialize(new GemstoneChangeAllRequest { OpType = 1, Slot = slot, PartnerSlotId = TeamIdx }));
        }

        // 查看boss信息
        public static void C2S_CheckBoss(int cid)
        {
            //if (CFG_XMLConfig.CheakLock(CFG_XMLConfig.XML_Config.SwitchOpne.World_Boss))
            //    return;

            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_CheckBoss, MType.Request, TCPClientMgr.ProtoSerialize(new CheckBossRequest { BossCfgId = cid }));
        }

        // 每日副本挑战次数
        public static void C2S_DailyMapInfo()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_DailyMapInfo, MType.Request, TCPClientMgr.ProtoSerialize(new DailyMapInfoRequest()));
        }


        // 游戏退出
        public static void C2S_Exit()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Exit, MType.Request,
               TCPClientMgr.ProtoSerialize(new ExitRequest()));
        }

        //移动
        public static void C2S_Move(MoveRequest moveRequest)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Move, MType.Request,
                   TCPClientMgr.ProtoSerialize(moveRequest));
        }

        //攻击
        public static void C2S_Attack(AttackRequest attackRequest)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Attack, MType.Request,
                   TCPClientMgr.ProtoSerialize(attackRequest));
        }

        //复活
        public static void C2S_Revive(int ReviveType)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Revive, MType.Request,
                TCPClientMgr.ProtoSerialize(new ReviveRequest { ReviveType = ReviveType }));
        }



        // 双倍领取
        public static void C2S_DoubleReward(bool isDouble)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_DoubleReward, MType.Request,
               TCPClientMgr.ProtoSerialize(new DoubleRewardRequest { Double = isDouble }));
        }

        //采集
        public static void C2S_Mining(long uid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Mining, MType.Request,
               TCPClientMgr.ProtoSerialize(new MiningRequest { Uid = uid }));
        }

        public static void C2S_AutoPlay()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_AutoPlay, MType.Request,
               TCPClientMgr.ProtoSerialize(new AutoPlayRequest()));
        }

        // 自动分解 
        public static void C2S_AutoDisassembly(AutoDisassemblyInfo autoDisassemblyInfo)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_AutoDisassembly, MType.Request,
               TCPClientMgr.ProtoSerialize(new AutoDisassemblyRequest { Info = autoDisassemblyInfo }));
        }

        //首页开宝箱
        public static void C2S_HomeBoxSetting(int level, int job)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_HomeBoxSetting, MType.Request,
               TCPClientMgr.ProtoSerialize(new HomeBoxSettingRequest { level = level, Job = job }));
        }

        //首页开宝箱
        public static void C2S_HomeOpenBox(int _Num, int level, int job)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_HomeOpenBox, MType.Request,
               TCPClientMgr.ProtoSerialize(new HomeOpenBoxRequest { Num = _Num, level = level, Job = job }));
        }

        public static void C2S_HomeUpgradeBox()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_HomeUpgradeBox, MType.Request,
                TCPClientMgr.ProtoSerialize(new HomeUpgradeBoxRequest()));
        }

        public static void C2S_HomeUpgradeBoxQuick(int num)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_HomeUpgradeBoxQuick, MType.Request,
                TCPClientMgr.ProtoSerialize(new HomeUpgradeBoxQuickRequest() { Num = num }));
        }

        //主界面挂机
        public static void C2S_MainBattle(int id)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_MainBattle, MType.Request,
                TCPClientMgr.ProtoSerialize(new MainBattleRequest() { MainBattleId = id }));
        }

        public static void C2S_MainBattlePrepare(int job)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_MainBattlePrepare, MType.Request,
                TCPClientMgr.ProtoSerialize(new MainBattlePrepareRequest() { job = job }));
        }

        #region 任务

        //任务领取
        public static void C2S_TaskComplate(int num)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_TaskComplate, MType.Request, TCPClientMgr.ProtoSerialize(
           new TaskComplateRequest { TaskId = num }));
        }
        #endregion

        #region 装备
        //装备回收
        public static void C2S_ItemDisassembly(long[] uids)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ItemDisassembly, MType.Request, TCPClientMgr.ProtoSerialize(new ItemDisassemblyRequest { ItemUids = uids }));
        }

        //装备放置背包
        public static void C2S_BoxResIntoBag(long uid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_BoxResIntoBag, MType.Request, TCPClientMgr.ProtoSerialize(new BoxResIntoBagRequest { Uid = uid }));
        }

        /// <summary>
        /// 装备穿脱
        /// </summary>
        public static void C2S_ChangeEquip(long ItemUid, int OpType, int Slot, int PartnerSlotId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ChangeEquip, MType.Request, TCPClientMgr.ProtoSerialize(new ChangeEquipRequest
            {
                ItemUid = ItemUid,
                OpType = OpType,
                Slot = (EquipSlot)Slot,
                PartnerSlotId = PartnerSlotId,
            }));
        }
        /// <summary>
        /// 一键装备
        /// </summary>
        public static void C2S_EquipBest(ObjectType objectType, int soltID)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_EquipBest, MType.Request, TCPClientMgr.ProtoSerialize(new EquipBestRequest { Type = objectType, PartnerSlotId = soltID }));
        }
        #endregion

        #region 伙伴
        /// <summary>
        /// 获取伙伴信息
        /// </summary>
        public static void C2S_PartnerList()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerList, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerListRequest()));
        }
        /// <summary>
        /// 取消上阵
        /// </summary>
        public static void C2S_PartnerRest(int SlotId, long soltID)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerRest, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerRestRequest { SlotId = SlotId, PartnerUid = soltID }));
        }

        /// <summary>
        /// 上阵
        /// </summary>
        public static void C2S_PartnerCombat(int SlotId, long soltID)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerCombat, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerCombatRequest { SlotId = SlotId, PartnerUid = soltID }));
        }

        /// <summary>
        /// 伙伴替换
        /// </summary>
        public static void C2S_PartnerReplace(int SlotId, long soltID)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerReplace, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerReplaceRequest { SlotId = SlotId, PartnerUid = soltID }));
        }

        /// <summary>
        /// 伙伴升星
        /// </summary>
        public static void C2S_PartnerStar(int SlotId, long soltID)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerStar, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerStarRequest { SlotId = SlotId, PartnerUid = soltID }));
        }

        /// <summary>
        /// 伙伴合成
        /// </summary>
        public static void C2S_PartnerCompose(long PartnerCfgId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerCompose, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerComposeRequest { PartnerCfgId = PartnerCfgId }));
        }

        /// <summary>
        /// 伙伴重铸
        /// </summary>
        public static void C2S_PartnerReforge(long PartnerUid, bool isFree)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerReforge, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerReforgeRequest { PartnerUid = PartnerUid, IsFree = isFree }));
        }

        /// <summary>
        /// 伙伴分解
        /// </summary>
        public static void C2S_PartnerDisassemble(long PartnerUid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PartnerDisassemble, MType.Request, TCPClientMgr.ProtoSerialize(new PartnerDisassembleRequest { PartnerUid = PartnerUid }));
        }
        #endregion

        #region 宠物

        //宠物移动
        public static void C2S_PetMove(int SlotId, int TargetSlotId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PetMove, MType.Request,
               TCPClientMgr.ProtoSerialize(new PetMoveRequest { SlotId = SlotId, TargetSlotId = TargetSlotId }));
        }

        //宠物升级
        public static void C2S_PetLevel()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PetLevel, MType.Request,
               TCPClientMgr.ProtoSerialize(new PetLevelRequest { SlotId = 0 }));
        }

        //宠物升星 
        public static void C2S_PetStar(int cfgId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_PetStar, MType.Request,
               TCPClientMgr.ProtoSerialize(new PetStarRequest { SlotId = 0, CfgId = cfgId, Num = 1 }));
        }

        #endregion

        #region 玩法选择信息
        //领取章节奖励
        public static void C2S_BattleReward(CFG_Room.RoomType roomType, int chapterId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_BattleReward, MType.Request,
               TCPClientMgr.ProtoSerialize(new BattleRewardRequest { MType = (MapType)roomType, ChapterId = chapterId }
               ));
        }

        /// <summary>
        /// 副本信息
        /// </summary>
        public static void C2S_InstanceInfo()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_InstanceInfo, MType.Request,
               TCPClientMgr.ProtoSerialize(new InstanceInfoRequest()));
        }

        /// <summary>
        /// 扫荡
        /// </summary>
        public static void C2S_CopySweep(int battId, int SweepTimes)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_CopySweep, MType.Request,
               TCPClientMgr.ProtoSerialize(new CopySweepRequest
               {
                   RoomId = CFG_Room.GetUIDByRoomType((CFG_Room.RoomType)CFG_Battle.GetById(battId).Copytype),
                   BattleId = battId,
                   SweepTimes = SweepTimes
               }));
        }

        // 扫荡
        public static void C2S_CopySweep(int roomId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_CopySweep, MType.Request,
               TCPClientMgr.ProtoSerialize(new CopySweepRequest { RoomId = roomId }));
        }

        /// <summary>
        /// 海岛挂机 倍率
        /// </summary>
        public static void C2S_GuardIslandMultiplier()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuardIslandMultiplier, MType.Request,
               TCPClientMgr.ProtoSerialize(new GuardIslandMultiplierRequest { ID = 0 }));
        }

        #endregion

        #region 工会

        // 工会列表
        public static void C2S_GuildListRequest(int page, int pageSize)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildList, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildListRequest() { Page = page, PageSize = pageSize }));
        }

        // 创建工会
        public static void C2S_GuildCreateRequest(string name, string description)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildCreate, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildCreateRequest() { Name = name, Description = description }));
        }

        // 申请加入工会
        public static void C2S_GuildApplyRequest(long[] guildIds)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildApply, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildApplyRequest() { GuildIds = guildIds }));
        }

        // 取消申请
        public static void C2S_GuildCancelApplyRequest(long guildId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildCancelApply, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildCancelApplyRequest() { GuildId = guildId }));
        }

        // 解散工会
        public static void C2S_GuildDismissRequest(long guildId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildDismiss, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildDismissRequest() { GuildId = guildId }));
        }

        // 离开工会
        public static void C2S_GuildLeaveRequest(long guildId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildLeave, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildLeaveRequest() { GuildId = guildId }));
        }

        // 踢出工会
        public static void C2S_GuildKickRequest(long guildId, long memberId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildKick, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildKickRequest() { GuildId = guildId, MemberId = memberId }));
        }

        // 提升工会职位
        public static void C2S_GuildPromoteRequest(long guildId, long memberId, int role)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildPromote, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildPromoteRequest() { GuildId = guildId, MemberId = memberId, Role = role }));
        }

        // 申请列表
        public static void C2S_GuildApplyListRequest()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildApplyList, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildApplyListRequest()));
        }

        // 处理申请
        public static void C2S_GuildHandleApplyRequest(long[] applicantIds, int accept)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildHandleApply, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildHandleApplyRequest() { ApplicantIds = applicantIds, Accept = accept }));
        }

        // 工会信息
        public static void C2S_GuildInfoRequest(long guildId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildInfoRequest() { GuildId = guildId }));
        }

        // 工会成员列表
        public static void C2S_GuildMemberListRequest(long guildId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildMemberList, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildMemberListRequest() { GuildId = guildId }));
        }

        // 工会捐献信息
        public static void C2S_GuildDonateInfoRequest()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildDonateInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildDonateInfoRequest()));
        }

        // 给工会捐献/签到
        public static void C2S_GuildDonateRequest(int type)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildDonate, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildDonateRequest() { Type = type }));
        }

        // 工会商城购买情况
        public static void C2S_GuildShopBuyInfoRequest()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildShopBuyInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildShopBuyInfoRequest()));
        }

        // 工会商城购买
        public static void C2S_GuildShopBuyRequest(int id, int buyNum)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildShopBuy, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildShopBuyRequest() { Id = id, BuyNum = buyNum }));
        }

        // 工会升级技能
        public static void C2S_GuildUpgradeSkillRequest(int id, int upgradeCount)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildUpgradeSkill, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuildUpgradeSkillRequest() { UUId = id, UpgradeCount = upgradeCount }));
        }

        // 观察玩家详情
        public static void C2S_ObservePlayerInfoRequest(long uid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ObservePlayerInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new ObservePlayerInfoRequest() { Uid = uid }));
        }

        // 保存玩家引导
        public static void C2S_GuidanceUpdateRequest(int[] guideTypes)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuidanceUpdate, MType.Request,
                TCPClientMgr.ProtoSerialize(new GuidanceUpdateRequest() { Guidance = guideTypes }));
        }

        // 拉取玩家聊天记录
        public static void C2S_ChatHistoryRequest(int uid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ChatHistory, MType.Request,
                TCPClientMgr.ProtoSerialize(new ChatHistoryRequest() { Channel = uid }));
        }

        // 拉取玩家聊天记录
        public static void C2S_ModuleNew(int confID, ModuleState val)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ModuleNew, MType.Request,
                TCPClientMgr.ProtoSerialize(new ModuleNewRequest() { Info = new ModuleInfo() { CfgId = confID, State = val } }));
        }
        #endregion

        #region 回合制玩法
        // 回合制阵容设置
        public static void C2S_TurnBasedSlotRequest(long[] uids)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_TurnBasedSlot, MType.Request,
                TCPClientMgr.ProtoSerialize(new TurnBasedSlotRequest() { Uids = uids }));
        }

        // 回合制挑战
        public static void C2S_TurnBasedChallengeRequest(long uid)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_TurnBasedChallenge, MType.Request,
                TCPClientMgr.ProtoSerialize(new TurnBasedChallengeRequest() { Uid = uid }));
        }

        // 五行试炼信息
        public static void C2S_ElementInfoRequest()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ElementInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new ElementInfoRequest()));
        }

        // 元素升级
        public static void C2S_ElementUpgradeRequest(int slotId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ElementUpgrade, MType.Request,
                TCPClientMgr.ProtoSerialize(new ElementUpgradeRequest() { SlotId = slotId }));
        }

        // 元素突破
        public static void C2S_ElementBreakthroughRequest(int slotId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ElementBreakthrough, MType.Request,
                TCPClientMgr.ProtoSerialize(new ElementBreakthroughRequest() { SlotId = slotId }));
        }

        // 五行元素出战设置
        public static void C2S_ElementPartnerSlotRequest(long[] uids)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ElementPartnerSlot, MType.Request,
                TCPClientMgr.ProtoSerialize(new ElementPartnerSlotRequest() { Uids = uids }));
        }

        // 五行试炼战斗
        public static void C2S_ElementChallengeRequest()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ElementChallenge, MType.Request,
                TCPClientMgr.ProtoSerialize(new ElementChallengeRequest()));
        }

        // 五行试炼扫荡
        public static void C2S_ElementSweepRequest()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ElementSweep, MType.Request,
                TCPClientMgr.ProtoSerialize(new ElementSweepRequest()));
        }

        // 五行试炼单层扫荡
        public static void C2S_ElementSweepForWardRequest(int num)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ElementSweepForward, MType.Request,
                TCPClientMgr.ProtoSerialize(new ElementSweepForWardRequest { Num = num }));
        }
        #endregion

        #region 称号

        /// <summary>
        /// 获取称号列表
        /// </summary>
        public static void C2S_TitleList()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_TitleList, MType.Request,
                TCPClientMgr.ProtoSerialize(new TitleListRequest()));
        }

        /// <summary>
        /// 称号合成
        /// </summary>
        public static void C2S_TitleCompose(int cfgID)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_TitleCompose, MType.Request,
                TCPClientMgr.ProtoSerialize(new TitleComposeRequest { CfgId = cfgID }));
        }

        /// <summary>
        /// 称号升级
        /// </summary>
        /// <param name="cfgID"></param>
        public static void C2S_TitleStart(int cfgID, int luckItemCfgId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_TitleStart, MType.Request,
                TCPClientMgr.ProtoSerialize(new TitleStartRequest { CfgId = cfgID, LuckItemCfgId = luckItemCfgId }));
        }

        /// <summary>
        /// 称号装备
        /// </summary>
        /// <param name="cfgID">0 为卸下</param>
        /// <param name="soltID"></param>
        public static void C2S_TitleEquip(int cfgID, int soltID)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_TitleEquip, MType.Request,
                TCPClientMgr.ProtoSerialize(new TitleEquipRequest { CfgId = cfgID, Slot = soltID }));
        }

        #endregion

        #region 召唤
        /*
          C2S_LotteryInfo = 1101,
        C2S_LotteryRefresh = 1102,
        C2S_Lottery = 1103,
        C2S_LotterySelect = 1104,
         */
        // 查询
        public static void C2S_LotteryInfo()
        {
            Debug.Log("C2S_LotteryInfo");
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_LotteryInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new LotteryInfoRequest()));
        }

        // 主动刷新
        public static void C2S_LotteryRefresh()
        {
            Debug.Log("C2S_LotteryRefresh");
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_LotteryRefresh, MType.Request,
                TCPClientMgr.ProtoSerialize(new LotteryRefreshRequest()));
        }

        // 抽奖
        public static void C2S_Lottery()
        {
            Debug.Log("C2S_Lottery");
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Lottery, MType.Request,
                TCPClientMgr.ProtoSerialize(new LotteryRequest()));
        }

        // 选择指定奖励
        public static void C2S_LotterySelect(int itemId)
        {
            Debug.Log("C2S_LotterySelect");
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_LotterySelect, MType.Request,
                TCPClientMgr.ProtoSerialize(new LotterySelectRequest() { ItemId = itemId }));
        }
        #endregion

        #region 装备槽位洗炼
        // 槽位洗炼
        public static void C2S_RefineRequest(RefineType refineType, CharacterSlot characterSlot, EquipSlot equipSlot, RefinedSlot refinedSlot)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Refine, MType.Request,
                TCPClientMgr.ProtoSerialize(new RefineRequest() { type = refineType, character_slot = characterSlot, equip_slot = equipSlot, refined_slot = refinedSlot }));
        }

        // 洗炼信息
        public static void C2S_RefineInfoRequest()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_RefineInfo, MType.Request,
                TCPClientMgr.ProtoSerialize(new RefineInfoRequest()));
        }

        // 放弃或保存洗练
        public static void C2S_RefineSelectRequest(CharacterSlot characterSlot, EquipSlot equipSlot, RefinedSlot refinedSlot, Stats stats)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_RefineSelect, MType.Request,
                TCPClientMgr.ProtoSerialize(new RefineSelectRequest() { character_slot = characterSlot, equip_slot = equipSlot, refined_slot = refinedSlot, selected_attributes = stats }));
        }

        #endregion

        #region 新功能开启
        public static void C2S_RecordOpenID(int moduleId)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ModuleOpenRecord, MType.Request,
                TCPClientMgr.ProtoSerialize(new ModuleOpenRecordRequest() { ModuleId = moduleId }));
        }
        #endregion

        #region 头像与头像框
        public static void C2S_UseAvatarRequest(AvatarType avatarType, int id)
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_UseAvatar, MType.Request,
                TCPClientMgr.ProtoSerialize(new UseAvatarRequest() { type = avatarType, id = id }));
        }
        #endregion

        #region 幸运金币
        public static void C2S_FateCoin()
        {
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_FateCoin, MType.Request,
                TCPClientMgr.ProtoSerialize(new FateCoinRequest()));
        }

        #endregion
    }
}