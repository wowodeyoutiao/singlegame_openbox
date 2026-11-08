using Cysharp.Threading.Tasks;

using HotUpdateScripts.Config;
using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;

using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using JEngine.Singleton;

using MetaFramework.UI;

using pb;
using System.Collections.Generic;

using Debug = UnityEngine.Debug;

namespace HotUpdateScripts.Proto
{
    //[Subscriber(ThreadMode.Main)]
    public class ProtoC : SingletonTemplate<ProtoC>
    {

        public void OnAddEvent()
        {
            UniGameEvent.AddEventListener<LoginToGameServerResponse>(Sever_LoginToGameServerResponse);
            UniGameEvent.AddEventListener<CompulsoryWithdrawalPush>(Sever_CompulsoryWithdrawalPush);
            UniGameEvent.AddEventListener<SystemNotificationPush>(Sever_SystemNotificationPush);
            UniGameEvent.AddEventListener<HomeOpenBoxResponse>(HomeOpenBoxResponse);

            UniGameEvent.AddEventListener<PvpRecorderResponse>(Sever_PvpRecorderResponse);
            UniGameEvent.AddEventListener<PVPInfoResponse>(Sever_PVPInfoResponse);
            UniGameEvent.AddEventListener<PVPChallengeResponse>(Sever_PVPChallengeResponse);
            UniGameEvent.AddEventListener<PvpMatchingResponse>(Sever_PvpMatchingResponse);
            UniGameEvent.AddEventListener<MatchingResponse>("MatchingResponse", Sever_MatchingResponse);

            UniGameEvent.AddEventListener<PVPMeleeKillNumPush>("PVPMeleeKillNumPush", Sever_PVPMeleeKillNumPush);
            UniGameEvent.AddEventListener<MatchingResultPush>("MatchingResultPush", Sever_MatchingResultPush);
            UniGameEvent.AddEventListener<MatchRoomJoinPush>("MatchRoomJoinPush", Sever_MatchRoomJoinPush);

            UniGameEvent.AddEventListener<MatchingStartPush>(Sever_MatchingStartPush);
            UniGameEvent.AddEventListener<GameResultMatchPush>(Sever_GameResultMatchPush);
            UniGameEvent.AddEventListener<GameResultMatchMeleePush>(Sever_GameResultMatchMeleePush);
            UniGameEvent.AddEventListener<DailyMapInfoResponse>(Sever_DailyMapInfoResponse);
            UniGameEvent.AddEventListener<MiningResponse>(Sever_MiningResponse);
            UniGameEvent.AddEventListener<DoubleRewardResponse>(Sever_DoubleRewardResponse);
            UniGameEvent.AddEventListener<GameResultDailyMapPush>(Sever_GameResultDailyMapPush);
            UniGameEvent.AddEventListener<DailyMapScorePush>(Sever_DailyMapScorePush);

            UniGameEvent.AddEventListener<WarRoomListResponse>(Sever_WarRoomListResponse);
            UniGameEvent.AddEventListener<JoinDeskResponse>(Sever_C2S_Join);
            UniGameEvent.AddEventListener<ReJoinDeskPush>(Sever_S2C_ReJoin);

            UniGameEvent.AddEventListener<CheckBossResponse>(Sever_CheckBossResponse);
            UniGameEvent.AddEventListener<GameObjectAddPush>(Sever_GameObjectAddPush);
            UniGameEvent.AddEventListener<GameObjectDelPush>(Sever_GameObjectDelPush);
            UniGameEvent.AddEventListener<ExitResponse>(Server_C2S_EXIT);

            UniGameEvent.AddEventListener<AutoPlayResponse>(AutoPlayResponse);
            UniGameEvent.AddEventListener<PlayerEnterDeskPush>(Sever_S2C_PlayerEnter);
            UniGameEvent.AddEventListener<ExitPush>(Sever_ExitPush);
            UniGameEvent.AddEventListener<LevelUpPush>(Sever_LevelUpPush);
            UniGameEvent.AddEventListener<TeleportPush>(Sever_Teleport);

            UniGameEvent.AddEventListener<MovePush>(Sever_MovePush);
            UniGameEvent.AddEventListener<AttackResponse>(Sever_AttackResponse);
            UniGameEvent.AddEventListener<AttackPush>(Sever_AttackPush);

            UniGameEvent.AddEventListener<ObjectHpChangePush>(Sever_ObjectHpChangePush);
            UniGameEvent.AddEventListener<MonsterInPush>(Sever_MonsterInPush);
            UniGameEvent.AddEventListener<ObjectDiePush>(Sever_ObjectDiePush);
            UniGameEvent.AddEventListener<RevivePush>(Sever_RevivePush);
            UniGameEvent.AddEventListener<ReviveResponse>(Sever_ReviveResponse);
            UniGameEvent.AddEventListener<InstanceTimeAddResponse>(Sever_InstanceTimeAddResponse);
            UniGameEvent.AddEventListener<InstanceTimeAddPush>(Sever_InstanceTimeAddPush);
            UniGameEvent.AddEventListener<MonsterOutPush>(Sever_MonsterOutPush);
            UniGameEvent.AddEventListener<MonsterUpdatePush>(Server_MonsterUpdatePush);
            UniGameEvent.AddEventListener<MonsterRefreshPush>(Server_MonsterRefreshPush);
            UniGameEvent.AddEventListener<MiningPush>(Server_MiningPush);

            UniGameEvent.AddEventListener<SkillCooldownPush>(Sever_SkillCooldownPush);
            UniGameEvent.AddEventListener<AddBuffPush>(Server_AddBuffPush);
            UniGameEvent.AddEventListener<BuffListPush>(Server_BuffListPush);
            UniGameEvent.AddEventListener<RemoveBuffPush>(Server_RemoveBuffPush);
            UniGameEvent.AddEventListener<ItemDropPush>(Server_ItemDropPush);
            UniGameEvent.AddEventListener<ConsecutiveKillsPush>(Server_ConsecutiveKillsPush);

            //包裹资源
            UniGameEvent.AddEventListener<GetBagResponse>(Sever_GetBagResponse);
            UniGameEvent.AddEventListener<ItemUseResponse>(Sever_ItemUseResponse);
            UniGameEvent.AddEventListener<ItemDestroyResponse>(Sever_ItemDestroyResponse);
            UniGameEvent.AddEventListener<TakeSourcePush>(Sever_TakeSourcePush);
            UniGameEvent.AddEventListener<AddSourcePush>(Sever_AddSourcePush);
            UniGameEvent.AddEventListener<ComposeResponse>(Sever_ComposeResponse);

            //充值信息刷新推送
            UniGameEvent.AddEventListener<PayRefreshPush>(Sever_PayRefreshPush);

            //任务
            UniGameEvent.AddEventListener<TaskAcceptPush>(Sever_TaskAcceptPush);
            UniGameEvent.AddEventListener<TaskRemovePush>(Sever_TaskRemovePush);
            UniGameEvent.AddEventListener<TaskUpdatePush>(Sever_TaskUpdatePush);
            UniGameEvent.AddEventListener<TaskComplateResponse>(Sever_TaskComplateResponse);

            //个人信息
            UniGameEvent.AddEventListener<SyncUserBasePush>(Sever_SyncUserBasePush);
            UniGameEvent.AddEventListener<SyncUserStatsPush>(Sever_SyncUserStatsPush);
            UniGameEvent.AddEventListener<SyncUserAdvanceInfoPush>(Sever_SyncUserAdvanceInfoPush);
            UniGameEvent.AddEventListener<ChangeRoleNameResponse>(Sever_ChangeRoleNameResponse);
            UniGameEvent.AddEventListener<CommonResponse>(Sever_CommonResponse);

            //开箱
            UniGameEvent.AddEventListener<HomeAutoEncounterPush>(Sever_HomeAutoEncounterPush);
            UniGameEvent.AddEventListener<BoxResIntoBagResponse>(Sever_BoxResIntoBagResponse);

            //伙伴
            UniGameEvent.AddEventListener<SyncPartnerListPush>(Sever_SyncPartnerListPush);
            UniGameEvent.AddEventListener<PartnerListResponse>(Sever_PartnerListResponse);
            UniGameEvent.AddEventListener<PartnerUpdatePush>(Sever_PartnerUpdatePush);
            UniGameEvent.AddEventListener<PartnerSlotInfoResponse>(Sever_PartnerSlotInfoResponse);
            UniGameEvent.AddEventListener<PartnerSlotUpdatePush>(Sever_PartnerSlotUpdatePush);

            //装备
            UniGameEvent.AddEventListener<ChangeEquipResponse>(Sever_ChangeEquipResponse);
            UniGameEvent.AddEventListener<ItemDisassemblyResponse>(Sever_ItemDisassemblyResponse);
            UniGameEvent.AddEventListener<AutoDisassemblyResponse>(Sever_AutoDisassemblyResponse);
            UniGameEvent.AddEventListener<EquipRefineResponse>(Sever_EquipRefineResponse);
            UniGameEvent.AddEventListener<EquipRefineReplaceResponse>(Sever_EquipRefineReplaceResponse);
            UniGameEvent.AddEventListener<EnhanceResponse>(Sever_EnhanceResponse);
            UniGameEvent.AddEventListener<StarsResponse>(Sever_StarsResponse);
            UniGameEvent.AddEventListener<SyncUserEquipsInfoPush>(Sever_SyncUserEquipsInfoPush);
            UniGameEvent.AddEventListener<GemstoneChangeResponse>(Sever_GemstoneChangeResponse);
            UniGameEvent.AddEventListener<GemstoneChangeAllResponse>(Sever_GemstoneChangeAllResponse);

            //技能
            UniGameEvent.AddEventListener<SyncUserSkillSlotPush>(SyncUserSkillSlotPush);
            UniGameEvent.AddEventListener<NewSkillPush>(Sever_NewSkillPush);
            UniGameEvent.AddEventListener<SkillUpgradeResponse>(Sever_SkillUpgradeResponse);
            UniGameEvent.AddEventListener<SkillChangeResponse>(Sever_SkillChangeResponse);
            UniGameEvent.AddEventListener<SkillAdvancedResponse>(Sever_SkillAdvancedResponse);

            //聊天
            UniGameEvent.AddEventListener<ChatResponse>(Sever_ChatResponse);
            UniGameEvent.AddEventListener<ChatPush>(Sever_ChatPush);
            UniGameEvent.AddEventListener<LampPush>(Sever_LampPush);

            //战宠 宠物
            UniGameEvent.AddEventListener<SyncUserPetInfoPush>(Sever_SyncUserPetInfoPush);
            UniGameEvent.AddEventListener<PetMoveResponse>(Sever_PetMoveResponse);
            UniGameEvent.AddEventListener<PetStarResponse>(Sever_PetStarResponse);

            //玩法信息
            UniGameEvent.AddEventListener<InstanceInfoResponse>(Sever_C2S_InstanceInfo);
            UniGameEvent.AddEventListener<GameResultBattlePush>(Sever_S2C_GameResultBattleMap);
            UniGameEvent.AddEventListener<CopySweepResponse>(Sever_CopySweepResponse);
            UniGameEvent.AddEventListener<GuardIslandMultiplierResponse>(Sever_GuardIslandMultiplierResponse);

            //宝箱
            UniGameEvent.AddEventListener<SyncUserBoxInfoPush>(Sever_SyncUserBoxInfoPush);

            //工会
            UniGameEvent.AddEventListener<GuildListResponse>(Sever_GuildListResponse);
            UniGameEvent.AddEventListener<GuildCreateResponse>(Sever_GuildCreateResponse);
            UniGameEvent.AddEventListener<GuildApplyResponse>(Sever_GuildApplyResponse);
            UniGameEvent.AddEventListener<GuildCancelApplyResponse>(Sever_GuildCancelApplyResponse);
            UniGameEvent.AddEventListener<GuildDismissResponse>(Sever_GuildDismissResponse);
            UniGameEvent.AddEventListener<GuildLeaveResponse>(Sever_GuildLeaveResponse);
            UniGameEvent.AddEventListener<GuildKickResponse>(Sever_GuildKickResponse);
            UniGameEvent.AddEventListener<GuildPromoteResponse>(Sever_GuildPromoteResponse);
            UniGameEvent.AddEventListener<GuildApplyListResponse>(Sever_GuildApplyListResponse);
            UniGameEvent.AddEventListener<GuildHandleApplyResponse>(Sever_GuildHandleApplyResponse);
            UniGameEvent.AddEventListener<GuildInfoResponse>(Sever_GuildInfoResponse);
            UniGameEvent.AddEventListener<GuildMemberListResponse>(Sever_GuildMemberListResponse);
            UniGameEvent.AddEventListener<GuildDonateInfoResponse>(Sever_GuildDonateInfoResponse);
            UniGameEvent.AddEventListener<GuildShopBuyInfoResponse>(Sever_GuildShopBuyInfoResponse);
            UniGameEvent.AddEventListener<GuildCopyRewardPush>(Sever_GuildCopyRewardPush);
            UniGameEvent.AddEventListener<ObservePlayerInfoResponse>(Sever_ObservePlayerInfoResponse);
            UniGameEvent.AddEventListener<GuidanceUpdateResponse>(Sever_GuidanceUpdateResponse);
            UniGameEvent.AddEventListener<ChatHistoryResponse>(Sever_ChatHistoryResponse);

            //回合制玩法
            UniGameEvent.AddEventListener<TurnBasedSlotResponse>(Sever_TurnBasedSlotResponse);
            UniGameEvent.AddEventListener<TurnBasedChallengeResponse>(Sever_TurnBasedChallengeResponse);
            UniGameEvent.AddEventListener<ElementInfoResponse>(Sever_ElementInfoResponse);
            UniGameEvent.AddEventListener<ElementUpgradeResponse>(Sever_ElementUpgradeResponse);
            UniGameEvent.AddEventListener<ElementPartnerSlotResponse>(Sever_ElementPartnerSlotResponse);
            UniGameEvent.AddEventListener<ElementBreakthroughResponse>(Sever_ElementBreakthroughResponse);
            UniGameEvent.AddEventListener<ElementChallengeResponse>(Sever_ElementChallengeResponse);
            UniGameEvent.AddEventListener<ElementSweepForWardResponse>(Sever_ElementSweepForWardResponse);

            //爬塔
            UniGameEvent.AddEventListener<EndlesstrialsRewardResponse>(Sever_EndlesstrialsRewardResponse);
            UniGameEvent.AddEventListener<GameResultEndlesstrialsPush>(Sever_GameResultEndlesstrialsPush);

            //称号
            UniGameEvent.AddEventListener<TitleListResponse>(Sever_TitleListResponse);
            UniGameEvent.AddEventListener<TitleComposeResponse>(Sever_TitleComposeResponse);
            UniGameEvent.AddEventListener<TitleStartResponse>(Sever_TitleStartResponse);
            UniGameEvent.AddEventListener<TitleEquipResponse>(Sever_TitleEquipResponse);

            // 召唤
            UniGameEvent.AddEventListener<LotteryInfoResponse>(Sever_LotteryInfoResponse);

            //新功能开启
            UniGameEvent.AddEventListener<ModuleOpenRecordResponse>(Server_ModuleOpenRecordResponse);

            // 洗练
            UniGameEvent.AddEventListener<RefineResponse>(Server_RefineResponse);
            UniGameEvent.AddEventListener<RefineInfoResponse>(Server_RefineInfoResponse);
            UniGameEvent.AddEventListener<RefineSelectResponse>(Server_RefineSelectResponse);
            // 头像与头像框
            UniGameEvent.AddEventListener<SyncUserAvatarPush>(Server_SyncUserAvatarPush);
            UniGameEvent.AddEventListener<UseAvatarResponse>(Server_UseAvatarResponse);

            // 幸运金币
            UniGameEvent.AddEventListener<FateCoinPush>(Server_FateCoinPush);
            UniGameEvent.AddEventListener<FateCoinResponse>(Server_FateCoinResponse);
            UniGameEvent.AddEventListener<CommonNotificationPush>(Server_CommonNotificationPush);

            UniGameEvent.AddEventListener<ModuleNewPush>(Server_ModuleNewPush);
            UniGameEvent.AddEventListener<ModuleNewResponse>(Server_ModuleNewResponse);
        }

        //C2S_Login
        private void Sever_LoginToGameServerResponse(LoginToGameServerResponse data)
        {
            Log.Print($"--------LoginToGameServerResponse------ {data.Code}");

            if (data.Code != 0)
            {
                UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack
                {
                    showStr = "账号登录信息异常，请重新登录！",
                    timeCount = 3,
                    cType = 1,
                    endFun = () =>
                    {
                        GameM.Instance.isSeverClose = false;
                        GameHelper.Instance.BackToLogin(true);
                    },
                });
                return;
            }
            RoleM.Instance.Nickname = data.Name;
            RoleM.Instance.Uid = data.Uid;
            RoleM.Instance.HeadUrl = data.HeadUrl;
            GameRoomM.Instance.AutoPlay = false;
            GameRoomM.Instance.AutoPlayLock = false;

            GameRoomM.Instance._GameState = GameRoomM.GameState.No;


            GameM.Instance.SeverTimeOff = GameUtil.Instance.GetCurrentTimeSeconds() - data.Timestamp;
            GameM.Instance.SeverTime = data.Timestamp;
            //Log.Print($"--------GameM.Instance.SeverTimeOff------ {GameM.Instance.SeverTimeOff} data.GameStatus {data.GameStatus}");

            TDAnalyticsUtil.Login(data.Uid.ToString());
            string account = PlayerPrefsExtend.GetString("TextAcount");
            TDAnalyticsUtil.SetDistinctId(account);
            //Debug.Log("TDAnalyticsUtil.Login() ++++");

            if (data.GameStatus == 0)
            {
                JEvent.defaultEvent.Post(new SenceChange
                {
                    scencePath = GameConst.HallPath,
                    CallBack = () =>
                    {
                        UITool.Open<UIHall>((ui) =>
                        {
                            if (RoleM.Instance.IsReceivedGuideData)
                            {
                                var _ui = ui as UIHall;
                                _ui.LoadChatHistyory();
                                //if (!_ui.CheckStartGuideOpenBox())
                                //    _ui.CheckNewFunctionTip();
                            }
                        });
                    },
                });
            }

            TcpC.C2S_PartnerList();

            GameHelper.Instance.SetZeroCheakTime();
        }

        //充值信息刷新推送
        private void Sever_PayRefreshPush(PayRefreshPush data)
        {
            Log.Print($"--------充值信息刷新推送------ {data.request_id}");
            StoreM.Instance.OnDoneBuyEvent(data.request_id);
        }

        private void Sever_CompulsoryWithdrawalPush(CompulsoryWithdrawalPush data)
        {
            //Log.Print($"--------CompulsoryWithdrawalPush------ {data.Status}");
            GameM.Instance.isSeverClose = true;

            UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack
            {
                showStr = data.Status == CompulsoryWithdrawalPush.WithdrawalStatus.CloseServer ? "服务器暂时维护，即将退出" : "检测到您的账号已在其他地方登录！",
                timeCount = 3,
                cType = 1,
                endFun = () =>
                {
                    GameM.Instance.isSeverClose = false;
                    GameHelper.Instance.BackToLogin(true);
                },
            });
        }

        private void Sever_SystemNotificationPush(SystemNotificationPush data)
        {
            // Log.Print($"--------LoginToGameServerResponse------ {data.Code}");
            if (data.Type == Notification.NEW_MAIL)
            {
                MailM.Instance.SetMailNew(true);
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Set_Email.ToString(), MailM.Instance.GetDotNum());
            }
            else if (data.Type == Notification.MAIL_FULL)
            {
                UITool.ShowTips("邮件已满，请及时清理！");
            }
            else if (data.Type == Notification.BAG_FULL)
            {
                UITool.ShowTips("包裹已满，请及时清理！");
            }
            else if (data.Type == Notification.INVALID_Notification)
            {
                UITool.ShowTips(data.Msg);
            }
        }

        private void HomeOpenBoxResponse(HomeOpenBoxResponse data)
        {
            //Log.Print($"--------HomeOpenBoxResponse------ {data.Code}");
            if (data.Code == 0)
            {
                RoleM.Instance.OpenBoxItems = data.ItemInfo;
                RoleM.Instance.OpenBoxItems.Sort((x, y) =>
                {
                    if (x.Type == ItemType.EQUIPMENT && y.Type == ItemType.EQUIPMENT)
                    {
                        int qualityX = (int)x.Equip.Quality;
                        int qualityY = (int)y.Equip.Quality;
                        return qualityY.CompareTo(qualityX);
                    }
                    else if (x.Type != ItemType.EQUIPMENT && y.Type != ItemType.EQUIPMENT)
                    {
                        int xQuality = CFG_Item.GetById(x.CfgId).Quality;
                        int yQuality = CFG_Item.GetById(y.CfgId).Quality;
                        return yQuality.CompareTo(xQuality);
                    }
                    else
                    {
                        int xType = (int)x.Type;
                        int yType = (int)y.Type;
                        return yType.CompareTo(xType);
                    }
                });
            }

            JEvent.defaultEvent.Post(new GameEvent.Sync_HomeOpenBox { data = data });
        }

        #region 爬塔

        // 爬塔离线奖励 预览/领取
        private void Sever_EndlesstrialsRewardResponse(EndlesstrialsRewardResponse data)
        {
            // Log.Print($"--------爬塔离线奖励------ ,data.Ready." + data.Ready);
            RoleM.Instance.SetTrainInfo(data.Ready, data.Layers, data.Reward);
            JEvent.defaultEvent.Post(new GameEvent.Sync_Endlesstrials());

            // RedPointManager.Instance.RPS.Dispatch(RedPointConst.Tower_Get.ToString(), RoleM.Instance.GetTrainRedNum());
        }

        // 游戏结算
        private void Sever_GameResultEndlesstrialsPush(GameResultEndlesstrialsPush data)
        {
            // Log.Print($"--------爬塔 游戏结算------ ");
            GameRoomM.Instance._GameState = GameRoomM.GameState.Train_Over;
        }

        // 


        #endregion

        #region PVP & GVG

        private void Sever_PvpRecorderResponse(PvpRecorderResponse data)
        {
            if (data.Code == 0)
            {
                // 1计时赛 2回合制
                UITool.Open<UIFightRecord>(null, data.Recorders, data.Type == 1 ? (int)CFG_Room.RoomType.PVP_One : (int)CFG_Room.RoomType.PVP_More);
            }
        }

        private void Sever_PVPChallengeResponse(PVPChallengeResponse data)
        {
            UITool.ShowTips(data.Error);
        }

        // pvp信息
        private void Sever_PVPInfoResponse(PVPInfoResponse data)
        {
            //  Log.Print($"--------pvp信息------  data.DailyMatchTimes{ data.DailyMatchTimes}  data.MatchType{ data.MatchType  }  data.DailyBuyNum{data.DailyBuyNum}");
            RoleM.Instance.SetAddPVPInfo(data.MatchType, data.DailyMatchTimes, data.DailyBuyNum);
            JEvent.defaultEvent.Post(new GameEvent.Sync_PVPInfo());

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.PVP_One.ToString(), RoleM.Instance.GetPVPRedNum());
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.PVP_More.ToString(), RoleM.Instance.GetGVGRedNum());
        }

        private void Sever_PvpMatchingResponse(PvpMatchingResponse data)
        {
            //Log.Print($"--------PvpMatchingResponse------  data.Code{data.Code},  data.MatchType{data.MatchType}, data.Player Count{data.Players.Count} ");

            if (data.Code == 0)
            {
                if (data.MatchType == 1)
                    UITool.Open<UIFightRank>(null, data.Players, CFG_Room.RoomType.PVP_One);
                else if (data.MatchType == 2)
                    UITool.Open<UIFightRank>(null, data.Players, CFG_Room.RoomType.PVP_More);
                //UITool.Open<UIFightRank>(null, data.Players, CFG_Room.RoomType.PVP_More);
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // pvp匹配 
        private void Sever_MatchingResponse(MatchingResponse data)
        {
            Log.Print($"--------pvp匹配------ data.Code{data.Code} data.Cancel{data.Cancel}data.MatchType{data.MatchType} ");
            if (data.Code == 0)
            {
                if (data.Cancel)
                    UITool.Close<UIPVPOneReady>();
                else
                    UITool.Open<UIPVPOneReady>(null, data.MatchType);
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        //混战击杀排行
        private void Sever_PVPMeleeKillNumPush(PVPMeleeKillNumPush data)
        {
            //Log.Print($"--------混战击杀排行------  ");
            JEvent.defaultEvent.Post(new GameEvent.PVP_MeleeKillNum { killInfos = data.List });
        }

        //匹配结果 
        private void Sever_MatchingResultPush(MatchingResultPush data)
        {
            //Log.Print($"--------pvp 匹配结果------ data.Status{data.Status} data.Players.Count {data.Players.Count} " );
            if (data.Status != 2)
            {
                GameRoomM.Instance.Add_PlayerInfo(data.Players);
                List<PlayerInfo> Emenys = new List<PlayerInfo>();
                for (int i = 0; i < data.Players.Count; i++)
                {
                    if (data.Players[i].Uid != RoleM.Instance.Uid)
                    {
                        Emenys.Add(data.Players[i]);
                    }
                }
                UITool.Get<UIPVPOneReady>().OnMatchResult(new GameEvent.PVP_MatchResult { Enemys = Emenys, IsBegin = data.Status == 0 });
                // JEvent.defaultEvent.Post(new GameEvent.PVP_MatchResult { Enemys = Emenys });
            }
            else
            {
                UITool.ShowTips("未匹配到玩家");
                UITool.Close<UIPVPOneReady>();
            }
        }

        //加入比赛房间 
        private void Sever_MatchRoomJoinPush(MatchRoomJoinPush data)
        {
            Log.Print($"--------pvp 加入比赛房间------  data.RoomCfgId{data.RoomCfgId}");

            GameRoomM.Instance.ClearObjInfo();
            SkillM.Instance.OnInitGameSkill();
            MonsterManager.Instance.clearMonsterList();

            GameRoomM.Instance.RoomTableInfo = new TableInfo
            {
                RoomCfgId = data.RoomCfgId,
                RoomUid = data.RoomUid,
                BattleId = 0,
                MonsterList = new List<MonsterInfo>(),
            };
            GameRoomM.Instance.CFGRoom = CFG_Room.GetById(data.RoomCfgId);

            int duration = GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_One ? CFG_XMLConfig.xmlConfig.PVPInfo.duration : CFG_XMLConfig.xmlConfig.GVGInfo.duration;
            GameRoomM.Instance.GameFightTime = GameUtil.Instance.GetCurrentTimeSeconds() + duration;

            if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_One)
            {
                JEvent.defaultEvent.Post(new SenceChange
                {
                    scencePath = GameConst.NorFightPath,
                    CallBack = () =>
                    {
                        HotUpdateScripts.Game.FightController.Instance.OnInit();
                    }
                });
            }
        }

        //比赛开始 
        private void Sever_MatchingStartPush(MatchingStartPush data)
        {
            // Log.Print($"--------pvp 比赛开始------  ");
            UITool.Open<UIPVPBeginTime>(null, UIPVPBeginTime.ShowType.Begin);
        }

        //游戏结算 
        private void Sever_GameResultMatchPush(GameResultMatchPush data)
        {
            // Log.Print($"--------pvp 游戏结算------  ");
            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                GameRoomM.Instance._GameState = GameRoomM.GameState.PVP_Over;
            UITool.Open<UIPVPGameOver>(null, data.Win, data.Score);
        }

        //游戏结算 混战
        private void Sever_GameResultMatchMeleePush(GameResultMatchMeleePush data)
        {
            // Log.Print($"--------混战 游戏结算------  ");
            GameRoomM.Instance._GameState = GameRoomM.GameState.PVP_Over;
            UITool.Open<UIGVGGameOver>(null, data.Result);
        }

        #endregion

        #region 每日必做类副本

        private void Sever_DailyMapInfoResponse(DailyMapInfoResponse data)
        {
            // Log.Print($"--------DailyMapInfoResponse------ {data.Code} ");

            RoleM.Instance.AddChallengeInfo(data.Info);
            JEvent.defaultEvent.Post(new GameEvent.Daily_MapInfo());

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.DayDo.ToString(), RedDotM.Instance.OnDayDoRedNum());
        }


        // 评分
        private void Sever_DailyMapScorePush(DailyMapScorePush data)
        {
            //Log.Print($"--------评分------ data.TotalDamage{data.TotalDamage}data.MonsterNum{data.MonsterNum}");
            GameRoomM.Instance.MonsterAllCount = data.MonsterNum;
            JEvent.defaultEvent.Post(new GameEvent.DailyMap_Score
            {
                TotalDamage = data.TotalDamage,
                Scoring = data.Scoring,
                Round = data.Round,
                MonsterNum = data.MonsterNum,
            });
        }

        //挖矿
        private void Sever_MiningResponse(MiningResponse data)
        {
            //Log.Print($"--------挖矿 ------ {data.Code} data.Error {data.Error}");

            UITool.ShowTips(data.Error);
        }


        //双倍领取
        private void Sever_DoubleRewardResponse(DoubleRewardResponse data)
        {
            UITool.ShowTips(data.Error);
            if (data.Code == 0)
            {
                UITool.Close<UIGameFightOver>();
                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                    FightController.Instance.OnExitGame();
            }
        }

        //游戏结算 每日必做
        private void Sever_GameResultDailyMapPush(GameResultDailyMapPush data)
        {
            Log.Print($"-------- 游戏结算 每日必做------  ");
            //if (GameRoomM.Instance._GameState == GameRoomM.GameState.No)
            //    UITool.Open<UIDayDoAward>(null, data.Items, data.RoomId);
            //else
            //    GameRoomM.Instance.gameResultDailyMap = data;
        }


        #endregion

        #region 即时战斗

        //战场房间信息
        private void Sever_WarRoomListResponse(WarRoomListResponse data)
        {
            //Log.Print($"-------战场房间信息------");
            if (data.Code == 0)
            {
                FightPlayM.Instance.SetWarFightRoomList(data.Rooms);
                JEvent.defaultEvent.Post(new GameEvent.Refresh_WarRoomInfo());
            }
            else
                UITool.ShowTips(data.Error);
        }

        /// <summary>
        /// 加入游戏
        /// </summary>
        /// <param name="data"></param>
        private void Sever_C2S_Join(JoinDeskResponse data)
        {
            //Log.Print($"--------JoinDeskResponse 加入游戏------ data.LeftTime{data.LeftTime } data.Code{data.Code} {data.Error}");

            if (data.Code == 0)
            {
                SkillM.Instance.OnInitGameSkill();

                MonsterManager.Instance.clearMonsterList();
                GameRoomM.Instance.RoomTableInfo = data.TableInfo;
                GameRoomM.Instance.GameFightTime = GameUtil.Instance.GetCurrentTimeSeconds() + data.LeftTime;
                GameRoomM.Instance.ClearObjInfo();

                UITool.Open<UILoading>(null, new UILoading.Loading { _time = 5f });

                JEvent.defaultEvent.Post(new SenceChange
                {
                    scencePath = GameConst.NorFightPath,
                    CallBack = () =>
                    {
                        FightController.Instance.OnInit();
                    },
                });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }


        private void Sever_S2C_ReJoin(ReJoinDeskPush data)
        {
            //Log.Print($"jerryc ReJoinDeskPush 断线重连------ data.LeftTime{data.LeftTime}");

            MonsterManager.Instance.clearMonsterList();
            GameRoomM.Instance.RoomTableInfo = data.TableInfo;
            GameRoomM.Instance.GameFightTime = GameUtil.Instance.GetCurrentTimeSeconds() + data.LeftTime;
            GameRoomM.Instance.ClearObjInfo();
            SkillM.Instance.OnInitGameSkill();

            JEvent.defaultEvent.Post(new SenceChange
            {
                scencePath = GameConst.NorFightPath,
                CallBack = () =>
                {
                    FightController.Instance.OnInit();
                },
            });


        }

        // 查看boss信息
        private void Sever_CheckBossResponse(CheckBossResponse data)
        {
            //Debug.Log($"查看boss信息 {data.CfgId} {data.Hp} ");

            UniGameEvent.Send(new GameEvent.Check_WorldBoss { BossID = data.CfgId, HP = data.Hp, UsedCount = data.Challengetime });
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.WorldBoss.ToString(), RedDotM.Instance.OnWorldBossRedDot(data.Challengetime));
        }

        //添加实例
        private void Sever_GameObjectAddPush(GameObjectAddPush data)
        {
            //Log.Print($"GameObjectAddPush 添加实例------ ");
            GameRoomM.Instance.Add_OneObjInfo(data.Info);
            FightController.Instance.SetGameObjectEff(data.Info, true);
        }

        //删除实例
        private void Sever_GameObjectDelPush(GameObjectDelPush data)
        {
            //Log.Print($"GameObjectDelPush 删除实例------ ");
            GameRoomM.Instance.Remove_OneObjInfo(data.Info);
            FightController.Instance.SetGameObjectEff(data.Info, false);

            UniGameEvent.Send(new GameEvent.Fight_Die { id = data.Info.Uid });

        }

        //退出房间
        private void Server_C2S_EXIT(ExitResponse data)
        {
            //Log.Print($"jerryc ExitResponse 退出游戏------ {data.Code} {data.Error}");
            if (data.Code == 0)
                UITool.Close<UIGameFightRevive>();
            if (GameRoomM.Instance._GameState != GameRoomM.GameState.Over
                && GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.PVP_One
                && GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.QHS
                && GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.SY
                && GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.Gold
                && GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.EXP)
            {
                GameRoomM.Instance._GameState = GameRoomM.GameState.Over;
                FightController.Instance.OnExitGame();
            }
            else
                GameRoomM.Instance._GameState = GameRoomM.GameState.Over;

        }

        //自动战斗
        private void AutoPlayResponse(AutoPlayResponse data)
        {
            GameRoomM.Instance.AutoPlayLock = data.Code == 0;

            // Log.Print($"AutoPlayResponse 自动攻击------ {data.Code} {data.Error}");
            //UITool.ShowTips(data.Error);
        }

        //玩家进入
        private void Sever_S2C_PlayerEnter(PlayerEnterDeskPush data)
        {
            Log.Print("--------S2C_PlayerEnter  玩家进入------");

            GameRoomM.Instance.Add_PlayerInfo(data.Data);
            FightController.Instance.SetInitPlayerAsync().Forget();
        }

        //玩家退出
        private void Sever_ExitPush(ExitPush data)
        {
            Log.Print("--------ExitPush------ " + data.Code + " ID " + data.Uid);

            if (data.Uid == MainRole._MainRole.targetID)
                MainRole._MainRole.targetID = 0;

            if (data.Uid == RoleM.Instance.Uid)
            {
                if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.SXS || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.HS ||
                    GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_Peace || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Guild ||
                    GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.WorldBoss || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace)
                {
                    GameRoomM.Instance._GameState = GameRoomM.GameState.Over;
                    if (data.Code != 2)
                    {
                        //RoleM.Instance.Guidance = false;
                        string[] Strs =  {
                                "",
                                "<size=28>副本时间已到！</size>\n<color=#F9FF00>{0}</color>秒后自动退出",
                                "",
                                "<size=28>怪物已被击杀！</size>\n<color=#F9FF00>{0}</color>秒后自动退出",
                                "<size=28>角色已死亡！</size>\n<color=#F9FF00>{0}</color>秒后自动退出",
                            };
                        UIMessageBoxOne.CallBack _call = new UIMessageBoxOne.CallBack
                        {
                            cType = 1,
                            timeCount = 3,
                            showStr = Strs[data.Code],
                            endFun = () =>
                            {
                                FightController.Instance.OnExitGame();
                            }
                        };
                        UITool.Open<UIMessageBoxOne>(null, _call);
                    }
                }
            }
            else
            {
                GameRoomM.Instance.Remove_PlayerInfo(data.Uid);
                FightController.Instance.OnPlayerExit(data.Uid);
                UniGameEvent.Send(new GameEvent.Fight_Die { id = data.Uid });
            }
        }

        //升级
        private void Sever_LevelUpPush(LevelUpPush data)
        {
            //Log.Print($"--------LevelUpPush------{data.Uid}  data.NewLevel {data.NewLevel}");

            if (data.Uid == RoleM.Instance.Uid)
            {
                RoleM.Instance.Level = data.NewLevel;
            }
            else
            {
                PartnerM.Instance.SetOnePartnerLv(data.Uid, data.NewLevel);
            }

            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
            {
                FightController.Instance.OnPlayerUpLv(data.Uid, data.NewLevel);
            }

            UniGameEvent.Send(new GameEvent.Refresh_LvUp { uid = data.Uid });
        }

        //传送 
        private void Sever_Teleport(TeleportPush data)
        {
            //Log.Print($"--------传送------ " + data.Uid);
            if (!MonsterManager.Instance.SetTeleport(data.Uid, data.Location, data.Direction, data.Type))
            {
                FightController.Instance.SetTeleport(data.Uid, data.Location, data.Direction, data.Type);
            }
        }

        //移动
        private void Sever_MovePush(MovePush data)
        {
            //Log.Print("--------MovePush------" + data.Uid);

            if (MonsterManager.Instance.monsterMove(data.Uid, data.Location, data.Direction))
            {
                return;
            }
            else
            {
                FightController.Instance.SetRoleMove(data);
            }
        }

        private void Sever_AttackResponse(AttackResponse data)
        {
            // Log.Print("jerryc--------AttackResponse------"+ data.Code);
            if (data.Error != "" && data.Error != "技能冷却时间未到")
            {
                //Log.Print("AttackResponse" + data.Code  + " data.Error "+ data.Error);
                UITool.ShowTips(data.Error);

                if (data.Code == -12)
                {
                    var strs = data.Error.Split(':');
                    long uid = long.Parse(strs[(strs.Length - 1)]);
                    UITool.ShowTips((GameRoomM.Instance.DeathHashSet.Contains(uid) ? "已处理怪物死亡" : "未收到服务器信息"));
                }
            }
        }

        //攻击
        private void Sever_AttackPush(AttackPush data)
        {
            //Log.PrintYellow($"--------AttackPush------{data.AttackerId}");
            FightController.Instance.SetRoleAttack(data);
        }

        //血量变化
        private void Sever_ObjectHpChangePush(ObjectHpChangePush data)
        {
            //Log.Print($"--------血量变化------TargetId{ data.TargetId }ChangeHp{ data.ChangeHp}");
            MonsterManager.Instance.SetMonsterHPChange(data.TargetId, data.ChangeHp, data.MaxHp, data.AttackerId, data.SkillId);
            FightController.Instance.SetRoleHpChange(data.TargetId, data.ChangeHp, data.MaxHp, data.AttackerId, data.SkillId);

            UniGameEvent.Send(new GameEvent.Fight_HPChange { tid = data.TargetId, changeNum = data.ChangeHp, Critical = data.Critical, Miss = data.Miss, AttackerId = data.AttackerId });
        }

        //怪物加入
        private void Sever_MonsterInPush(MonsterInPush data)
        {
            // Log.Print("--------MonsterInPush------");

            for (int i = 0; i < data.MonsterList.Count; i++)
            {
                if (GameRoomM.Instance._GameState == GameRoomM.GameState.No)
                    MonsterManager.Instance.addCreateMonsterList(data.MonsterList[i]);
                else
                    MonsterManager.Instance.createMonster(data.MonsterList[i]).Forget();
            }

            UniGameEvent.Send(new GameEvent.Fight_MonsterAdd());
        }

        // 怪物/角色死亡
        private void Sever_ObjectDiePush(ObjectDiePush data)
        {
            //Log.Print("怪物死亡 " + data.ObjectId);
            if (GameRoomM.Instance.DeathHashSet.Contains(data.ObjectId))
                Log.Print("----重复死亡 " + data.ObjectId);
            else
                GameRoomM.Instance.DeathHashSet.Add(data.ObjectId);
            if (GameRoomM.Instance._GameState == GameRoomM.GameState.Play)
            {
                if (data.ObjectId == MainRole._MainRole.targetID)
                    MainRole._MainRole.targetID = 0;
                if (data.ObjectId == MainRole._MainRole.targetCollectID)
                    MainRole._MainRole.targetCollectID = 0;
            }

            if (!MonsterManager.Instance.removeMonster(data.ObjectId))
            {
                FightController.Instance.SetPlayerDeath(data.ObjectId);

                if (data.ObjectId == RoleM.Instance.Uid)
                {
                    if ((GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace || GameRoomM.Instance.IsDayDo()) && data.ReviveTip)
                    {
                        UITool.Open<UIGameFightRevive>(null, data);
                    }
                    if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Guide)
                    {
                        UITool.Open<UIPVPBeginTime>(null, UIPVPBeginTime.ShowType.DieToLife);
                    }
                }
            }
            UniGameEvent.Send(new GameEvent.Fight_Die { id = data.ObjectId });
        }

        // 复活
        private void Sever_RevivePush(RevivePush data)
        {
            //Log.PrintYellow("复活 push------------"+ data.Uid);
            if (data.Uid == RoleM.Instance.Uid)
            {
                if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More && UITool.Check<UIPVPBeginTime>())
                {
                    UITool.Close<UIPVPBeginTime>();
                }
            }

            if (GameRoomM.Instance.DeathHashSet.Contains(data.Uid))
                GameRoomM.Instance.DeathHashSet.Remove(data.Uid);

            GameRoomM.Instance.SetPlayerHP(data.Uid, data.CurHp, 0, true);
            FightController.Instance.SetPlayerRevive(data);
            UniGameEvent.Send(new GameEvent.Fight_Revive { id = data.Uid });
        }

        // 复活
        private void Sever_ReviveResponse(ReviveResponse data)
        {
            // Log.Print($"复活 Back ----data.Code  {data.Code } data.Error {data.Error}" );
            if (data.Code != 0)
            {
                UITool.ShowTips(data.Error);
            }
        }

        //增加副本时间
        private void Sever_InstanceTimeAddResponse(InstanceTimeAddResponse data)
        {
            // Log.Print($"增加副本时间 ----data.Code  {data.Code } data.Error {data.Error}");

            if (data.Code == 0)
            {
                GameRoomM.Instance.TimeAutoBuy = data.AutoBuy;
                GameRoomM.Instance.GameFightTime = GameUtil.Instance.GetCurrentTimeSeconds() + data.InstanceTimeLeft;
            }
            else
                UITool.ShowTips(data.Error);
        }


        //增加副本时间
        private void Sever_InstanceTimeAddPush(InstanceTimeAddPush data)
        {
            // Log.Print($"增加副本时间 Push ----");
            GameRoomM.Instance.GameFightTime = GameUtil.Instance.GetCurrentTimeSeconds() + data.InstanceTimeLeft;
        }

        //怪物(角色)离开
        private void Sever_MonsterOutPush(MonsterOutPush data)
        {
            Log.Print("--------MonsterOutPush------");
            //if (data.MonsterId == MainRole._MainRole.targetID)
            //    MainRole._MainRole.targetID = 0;

            //MonsterManager.Instance.removeMonster(data.MonsterId);
            //JEvent.defaultEvent.Post(new GameEvent.Fight_Die { id = data.MonsterId });
        }

        //采集广播
        private void Server_MiningPush(MiningPush data)
        {
            //Log.Print("--------采集广播------" + data.Action + " " + data.TargetUid + " " + data.PlayerUid);

            if (data.PlayerUid == RoleM.Instance.Uid)
            {
                if (data.Action == MiningPush.MiningAction.LIMIT)
                    MainRole._MainRole.targetCollectID = 0;

                if (data.Action == MiningPush.MiningAction.BREAK || data.Action == MiningPush.MiningAction.LIMIT)
                    MonsterManager.Instance.OnRefreshCollectlMonser(data.TargetUid);
            }
            if (data.Action == MiningPush.MiningAction.START)
            {
                var palyer = GameRoomM.Instance.getPlayerInfoByID(data.PlayerUid);
                GameRoomM.Instance.CollectMonsterDic.Add(data.TargetUid, palyer == null ? string.Empty : palyer.Nickname);
            }
            else if (data.Action == MiningPush.MiningAction.BREAK || data.Action == MiningPush.MiningAction.LIMIT)
            {
                GameRoomM.Instance.CollectMonsterDic.Remove(data.TargetUid);
            }

        }

        //怪物更新
        private void Server_MonsterUpdatePush(MonsterUpdatePush data)
        {
            //Log.Print("--------怪物更新------");
            MonsterManager.Instance.OnRefreshMonser(data.Monster);
            UniGameEvent.Send(new GameEvent.Update_OneMonster { uid = data.Monster.Uid });
        }

        //怪物刷新
        private void Server_MonsterRefreshPush(MonsterRefreshPush data)
        {
            //Log.Print("--------怪物刷新------");
            for (int i = 0; i < data.List.Count; i++)
            {
                if (CFG_Monster.GetById(data.List[i].CfgId).Isboss == 4)
                {
                    GameRoomM.Instance.MosterBossRefreshTime = data.List[i].RefreshTime;
                    break;
                }
            }
            //UniGameEvent.Send(new GameEvent.Refresh_Monsters());
        }


        //技能刷新
        private void Sever_SkillCooldownPush(SkillCooldownPush data)
        {
            // Log.PrintYellow($"--------SkillCooldownPush------ data.CfgId {data.CfgId} data.Delay {data.Delay}");
            SkillM.Instance.SetGameSkillDealy(data.CfgId, GameUtil.Instance.GetCurrentTimeSeconds() + data.Delay / 1000);
        }

        //add buffer
        private void Server_AddBuffPush(AddBuffPush data)
        {
            //Log.PrintYellow("--------AddBuffPush------"+ data.Buff.BuffId);
            if (data.TargetId == RoleM.Instance.Uid)
                RoleM.Instance.AddBuffInfo(data.Buff);

            GameRoomM.Instance.Add_PlayerBuffer(data.TargetId, data.Buff);

            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
            {
                FightController.Instance.AddOrMoveOneBuffer(data.TargetId, data.Buff.BuffId, true);
                UniGameEvent.Send(new GameEvent.Refresh_Buffer { Uid = data.TargetId });
            }

        }

        //add buffer
        private void Server_BuffListPush(BuffListPush data)
        {
            //Log.PrintYellow("--------BuffListPush------");
            RoleM.Instance.RefreshBuffInfo(data.BuffList);
        }

        //删除buff S2C_RemoveBuff
        private void Server_RemoveBuffPush(RemoveBuffPush data)
        {
            //Log.PrintYellow("-------- RemoveBuffPush ------"+ data.BuffId);
            if (data.TargetId == RoleM.Instance.Uid)
            {
                RoleM.Instance.RemoveBuffInfo(data.BuffId);
            }
            GameRoomM.Instance.Remove_PlayerBuffer(data.TargetId, data.BuffId);

            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
            {
                FightController.Instance.AddOrMoveOneBuffer(data.TargetId, data.BuffId, false);
                UniGameEvent.Send(new GameEvent.Refresh_Buffer { Uid = data.TargetId });
            }
        }

        //死亡 掉落
        private void Server_ItemDropPush(ItemDropPush data)
        {
            //Log.PrintYellow("--------ItemDropPush------"+ data.Location.X + "  data.Location.Y "+ data.Location.Y);
            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
            {
                MainRole._MainRole.SetDropEff(data.Location, data.Items);
            }
            GameRoomM.Instance.Add_DropItemInfo(data.Items);

            JEvent.defaultEvent.Post(new GameEvent.Refresh_ItemDrop());
        }

        //连战
        private void Server_ConsecutiveKillsPush(ConsecutiveKillsPush data)
        {
            if (data.Uid == RoleM.Instance.Uid)
                JEvent.defaultEvent.Post(new GameEvent.Refresh_ConsecutiveKills { count = data.Count });
        }

        #endregion

        #region 包裹 资源
        private void Sever_GetBagResponse(GetBagResponse data)
        {
            // Log.Print("--------GetBagResponse------" + data.ItemList);
            BagM.Instance.SetBagInfo(data.ItemList.ItemList);
            BagM.Instance.SetDailyUsedItems(data.DailyItemUsedList);
            //JEvent.defaultEvent.Post(new Event_Source_Refresh { isChange = false, _ItemInfo = data.ItemList.ItemList }) ;
            UniGameEvent.Send(new Event_Source_Refresh { isChange = false, _ItemInfo = data.ItemList.ItemList });

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Equip);
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_Award);
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Call);
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Reinforce);
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Insert);
        }

        //使用道具获得
        private void Sever_ItemUseResponse(ItemUseResponse data)
        {
            //Log.Print("--------ItemUseResponse------" + data.Code);
            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }

            if (data.Code == 0 && UITool.Check<UIItemUesTips>())
            {
                UITool.Get<UIItemUesTips>().RemoveTime();
            }

            if (data.Code == 0 && UITool.Check<UIRole>())
            {
                UITool.ShowTips($"使用成功!");
            }

            foreach (var item in data.ItemList)
            {
                // 副本挂机、首领、遗迹 挑战次数刷新
                if (item.CfgId == 56 || item.CfgId == 25073 || item.CfgId == 25074)
                {
                    TcpC.C2S_InstanceInfo();
                    if (item.CfgId == 56)
                        TcpC.GetBagInfo();
                }
                // 副本挑战次数刷新
                else if (item.CfgId >= 25075 && item.CfgId <= 25080)
                    TcpC.C2S_DailyMapInfo();
                // 五行挑战次数刷新
                else if (item.CfgId == 25081)
                    TcpC.C2S_ElementInfoRequest();
                // 即时挑战次数刷新
                else if (item.CfgId == 25082)
                    TcpC.OnPVPInfo((int)MatchType.SINGLE_COMBAT);
                // 回合挑战次数刷新
                else if (item.CfgId == 25083)
                    TcpC.OnPVPInfo((int)MatchType.FREE_FOR_ALL);
            }
        }

        private void Sever_ItemDestroyResponse(ItemDestroyResponse data)
        {
            // Log.Print("--------ItemDestroyResponse------" + data.Code);

            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }
        }

        //增加资源
        private void Sever_AddSourcePush(AddSourcePush data)
        {
            //Log.Print($"--------增加资源---- Type{data.Type}  --{data.Item.Count}");
            if (data.Type == 4)
                return;


            bool isGoldSound = false;
            for (int i = 0; i < data.Item.Count; i++)
            {
                if (data.Item[i].CfgId == 1)
                    isGoldSound = true;

                if (data.Type != 2)
                {
                    UITool.ShowTips($"获得{GameUtil.Instance.GetItemName(data.Item[i].CfgId, ((int)CFG_Item.ItemType.Equip == CFG_Item.GetById(data.Item[i].CfgId).Type ? (int)data.Item[i].Equip.Quality : -1))}{GameUtil.Instance.FormatNumWithUnit(data.Item[i].Num)}");

                    if (data.Type == 1 && CFG_Item.CheckIsItemTipUsed(data.Item[i].CfgId) && CFG_Item.GetById(data.Item[i].CfgId).Type != (int)ItemType.EQUIPMENT)
                    {
                        UITool.ItemUseTips(data.Item[i]);
                    }
                }
            }

            if (isGoldSound)
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_SellAndGotCoin); //播放获取金币音效

            if (data.Type != 3)
            {
                BagM.Instance.AddItemsChange(data.Item);
                UniGameEvent.Send(new Event_Source_Refresh { isChange = true, _ItemInfo = data.Item });

                if (BagM.Instance.GetEquipChange())
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Equip);

            }
        }


        //扣除资源
        private void Sever_TakeSourcePush(TakeSourcePush data)
        {
            //Log.Print("--------TakeSourcePush------" + data.GetType().Name);
            BagM.Instance.DelItemsChange(data.Item);
            UniGameEvent.Send(new GameEvent.Sync_TakeSource());
            //JEvent.defaultEvent.Post(new GameEvent.Sync_TakeSource());

            if (BagM.Instance.GetEquipChange())
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Equip);

            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Call);

        }

        //合成
        private void Sever_ComposeResponse(ComposeResponse data)
        {
            Log.Print($"合成--------{data.GetType().Name}------ Code{data.Code}  {data.Error}");

            JEvent.defaultEvent.Post(new GameEvent.Sync_ComposeResponse { data = data });
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Insert);
        }

        #endregion

        #region 任务

        //接任务
        private void Sever_TaskAcceptPush(TaskAcceptPush data)
        {
            //Log.Print($"--------{"接任务"}{data.GetType().Name}------");
            TaskM.Instance.AddTask(data.task);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_Task { ConfigId = data.task.CfgId, TaskClass = data.task.Class });

            TaskM.Instance.SendTaskRed(data.task.Class);
        }

        //删任务
        private void Sever_TaskRemovePush(TaskRemovePush data)
        {
            //Log.Print($"--------{"删任务"}{data.GetType().Name}------");
            TaskM.Instance.RemoveTask(data.task);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_Task { ConfigId = data.task.CfgId, TaskClass = data.task.Class });
        }

        // 更新任务
        private void Sever_TaskUpdatePush(TaskUpdatePush data)
        {
            // Log.Print($"--------{"更新任务"}{data.GetType().Name}{data.task.CfgId}------");

            TaskM.Instance.UpdateTask(data.task);

            //周常、日常,播放开箱子声音
            var taskInfo = CFG_Task.GetById(data.task.CfgId);
            if (taskInfo != null && data.task.Status == TaskStatus.REWARDED && (taskInfo.Class == 5 || taskInfo.Class == 6))
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_OpenBox);
            }

            JEvent.defaultEvent.Post(new GameEvent.Refresh_Task { ConfigId = data.task.CfgId, TaskClass = data.task.Class });
            TaskM.Instance.SendTaskRed(data.task.Class);
        }

        //完成任务
        private void Sever_TaskComplateResponse(TaskComplateResponse data)
        {
            //Log.Print($"--------{"完成任务"}{data.GetType().Name}------ data.Code {data.Code}");
            if (data.Code == 0)
            {
                TaskM.Instance.SendTaskRed(TaskClass.INVALID_TaskClass);
            }
            JEvent.defaultEvent.Post(new GameEvent.Refresh_Task { IsErr = data.Code == 0 });
            UITool.ShowTips(data.Error);
        }

        #endregion


        #region 个人信息
        //同步角色基础信息
        private void Sever_SyncUserBasePush(SyncUserBasePush data)
        {
            //Log.Print($"------ SyncUserBasePush --{data.GetType().Name}------{data.PetExp}data.PetExp");
            RoleM.Instance.Coin = data.Coin;
            RoleM.Instance.Soul = data.Soul;
            RoleM.Instance.Coupon = data.Coupon;
            RoleM.Instance.Energy = data.Energy;
            RoleM.Instance.Vip = data.Vip;
            RoleM.Instance.Exp = data.Exp;
            RoleM.Instance.Level = data.Level;
            RoleM.Instance.EnergyRecoveryTime = data.EnergyRecoveryTime;
            RoleM.Instance.Combat = data.Combat;
            RoleM.Instance.VipExp = data.VipExp;

            RoleM.Instance.FateLevel = data.FateLevel;
            RoleM.Instance.FateExp = (int)data.FateExp;

            RoleM.Instance.TrialVoucher = data.TrialVoucher;
            RoleM.Instance.TreasuerBox = data.TreasuerBox;
            RoleM.Instance.BoxTimeAllSecound = (int)data.HomeAFK;
            RoleM.Instance.PetUsedExp = data.PetExp;
            RoleM.Instance.BoxQuick = data.BoxQuick;
            RoleM.Instance.GuildId = data.GuildId;
            RoleM.Instance.GuildContribution = data.GuildContribution;
            RoleM.Instance.GuildWish = data.GuildWish;
            RoleM.Instance.WishRecoveryTime = data.WishRecoveryTime;

            RoleM.Instance.MainBattleId = data.MainBattleId;

            RoleM.Instance.TitleId = data.TitleId;
            RoleM.Instance.StarWishCoin = data.StarWishCoin;

            if (RoleM.Instance.IsFirstReceiveBaseInfo)
            {
                RoleM.Instance.guidePartnerRecordLevel = data.Level;
                RoleM.Instance.IsFirstReceiveBaseInfo = false;
            }

            UniGameEvent.Send(new GameEvent.Sync_UserBase());

            FuncOpenShowWidget.Instance.InitRecord(data.ModuleIds);
        }

        //同步角色属性
        private void Sever_SyncUserStatsPush(SyncUserStatsPush data)
        {
            //Log.Print($"同步角色属性--------{data.TotalCombat}------");
            RoleM.Instance.SetRoleStats(data.Attributes);
            PartnerM.Instance.SetPartnerStatsDic(data.Partners);
            if (RoleM.Instance.TotalCombat > 0 && RoleM.Instance.TotalCombat < data.TotalCombat)
            {
                if (!UIPowerTip.IsOpne)
                {
                    UIPowerTip.IsOpne = true;
                    UITool.Open<UIPowerTip>(null, data.TotalCombat, RoleM.Instance.TotalCombat);
                }
                else if (UITool.Check<UIPowerTip>())
                    UITool.Get<UIPowerTip>().OnShow(data.TotalCombat, RoleM.Instance.TotalCombat);
            }
            RoleM.Instance.Combat = data.Combat;
            RoleM.Instance.TotalCombat = data.TotalCombat;
            JEvent.defaultEvent.Post(new GameEvent.Sync_UserStatsPush());

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Equip);

            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
            {
                FightController.Instance.RefrshBuffer();
            }
        }

        //同步角色高级属性
        private void Sever_SyncUserAdvanceInfoPush(SyncUserAdvanceInfoPush data)
        {
            //Log.Print($"同步角色高级属性-------- {data.GetType().Name} ------- guides = {data.Guidance}");
            SkillM.Instance.SetSkillInfo(data.Skills);
            RoleM.Instance.SetRoleEquipDic(data.Equips);
            TaskM.Instance.SetAllTasks(data.Tasks);
            //SkillM.Instance.SetSkillEquipDic(data.SkillsSlot);
            TaskM.Instance.SendTaskRed(TaskClass.INVALID_TaskClass);
            //GuidanceM.Instance.SetGuideTypes(data.Guidance);

            NewPlayerGuideWidget.Instance.CacheGuideStep(data.Guidance);

            RoleM.Instance.IsReceivedGuideData = true;

            BagM.Instance.AutoDelSetInfo = data.AutoDisassembly;

            TurnBasedM.Instance.pvpFormationUids = data.TurnBasedUids;
            //Log.Print($"tcpc: TurnBasedUids-------- {string.Join(",", data.TurnBasedUids)} -------");
            JEvent.defaultEvent.Post(new GameEvent.Sync_UserAdvanceInfoPush());
        }

        //改名
        private void Sever_ChangeRoleNameResponse(ChangeRoleNameResponse data)
        {
            // Log.Print($"改名--------{data.GetType().Name}------");

            UITool.ShowTips(data.Error);
            JEvent.defaultEvent.Post(new GameEvent.Sync_ChangeName { isWin = data.Code == 0 });

            if (data.Code == 0)
            {
                UITool.ShowTips("改名成功");
                RoleM.Instance.Nickname = data.RoleName;
                UniGameEvent.Send(new GameEvent.Sync_UserBase());
            }
        }

        #endregion

        //通用返回
        private void Sever_CommonResponse(CommonResponse data)
        {
            //Log.Print($"通用返回--------{data.RequestCmdId}------ {data.Code}  {data.Error} ");
            UITool.ShowTips(data.Error);

            if (data.Code == 0)
            {
                if ((CmdId)data.RequestCmdId == CmdId.C2S_BattleReward)
                {
                    TcpC.C2S_InstanceInfo();
                }
                else if ((CmdId)data.RequestCmdId == CmdId.C2S_PartnerCombat || (CmdId)data.RequestCmdId == CmdId.C2S_PartnerRest)
                {
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Into);
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Equip);

                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Reinforce);
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Insert);

                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_UpLv);
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Advance);
                }
            }
            //JEvent.defaultEvent.Post(new GameEvent.Refrsh_Common { cid = (CmdId)data.RequestCmdId, code = data.Code });
            UniGameEvent.Send(new GameEvent.Refrsh_Common { cid = (CmdId)data.RequestCmdId, code = data.Code });
        }


        #region 大厅开箱

        //大厅开箱事件
        private void Sever_HomeAutoEncounterPush(HomeAutoEncounterPush data)
        {
            //  Log.Print($"大厅事件--------{data.GetType().Name}------ Type{data.Type} ---{data.Id}");

            JEvent.defaultEvent.Post(new GameEvent.Sync_HomeAutoEncounter { data = data });

        }

        //开箱装备保存
        private void Sever_BoxResIntoBagResponse(BoxResIntoBagResponse data)
        {
            //Log.Print($"装备保存--------{data.GetType().Name}------ ");

            if (RoleM.Instance.OpenBoxItems != null && RoleM.Instance.OpenBoxItems.Count > 0)
            {
                bool isRemove = false;
                for (int i = 0; i < RoleM.Instance.OpenBoxItems.Count; i++)
                {
                    if (RoleM.Instance.OpenBoxItems[i].Uid == data.Uid)
                    {
                        isRemove = true;
                        BagM.Instance.AddOneItemChange(RoleM.Instance.OpenBoxItems[i]);
                        RoleM.Instance.OpenBoxItems.RemoveAt(i);
                        break;
                    }
                }
                if (isRemove)
                    UniGameEvent.Send(new GameEvent.Refresh_HomeOpen());
            }
        }

        #endregion

        #region 伙伴

        //伙伴列表
        private void Sever_SyncPartnerListPush(SyncPartnerListPush data)
        {
            //Log.Print($"伙伴列表SyncPartnerListPush--------{data.GetType().Name}------ ");
            PartnerM.Instance.SetPartners(data.List);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_PartnerList { isOne = false });
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Into);
        }

        //伙伴列表
        private void Sever_PartnerListResponse(PartnerListResponse data)
        {
            //Log.Print($"伙伴列表--------{data.GetType().Name}------ ");
            PartnerM.Instance.SetPartners(data.List);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_PartnerList { isOne = false });
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Into);
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.ExpPool);
        }

        //伙伴更新
        private void Sever_PartnerUpdatePush(PartnerUpdatePush data)
        {
            //Log.Print($"伙伴更新--------{data.GetType().Name}------ ");
            PartnerM.Instance.SetOnePartners(data.Partner);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_PartnerList { isOne = true });
        }

        //伙伴槽位信息 主动拉取
        private void Sever_PartnerSlotInfoResponse(PartnerSlotInfoResponse data)
        {
            //Log.Print($"主动拉取 伙伴槽位信息--------{data.GetType().Name}------ ");
            PartnerM.Instance.SetPartnerSlotInfo(data.List);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_PartnerSlot());
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.ExpPool);
        }

        //伙伴槽位信息
        private void Sever_PartnerSlotUpdatePush(PartnerSlotUpdatePush data)
        {
            //Log.Print($"伙伴槽位信息--------{data.GetType().Name}------ ");
            PartnerM.Instance.SetPartnerSlotInfo(data.List);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_PartnerSlot());
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.ExpPool);
        }

        #endregion

        #region 装备

        //穿脱装备
        private void Sever_ChangeEquipResponse(ChangeEquipResponse data)
        {
            //Log.Print($"穿脱装备--------{data.Code} {data.Slot}------");
            if (data.Code == 0)
            {
                if (data.PartnerSlotId == 0)
                {
                    RoleM.Instance.SetRoleOneEquip(data.Item, data.Slot, data.OpType == 1);
                }

                UniGameEvent.Send(new GameEvent.ChangeEquip { item = data.Item, soldIdx = (int)data.Slot });
                //JEvent.defaultEvent.Post(new GameEvent.ChangeEquip { item = data.Item, soldIdx = (int)data.Slot });

                if (RoleM.Instance.OpenBoxItems != null && RoleM.Instance.OpenBoxItems.Count > 0)
                {
                    bool isChange = false;
                    for (int i = 0; i < RoleM.Instance.OpenBoxItems.Count; i++)
                    {
                        if (RoleM.Instance.OpenBoxItems[i].Uid == data.Item.Uid)
                        {
                            isChange = true;

                            //当前槽位是空的时候开箱数据直接删除，不用替换
                            //UnityEngine.Debug.Log($"开箱数量：{RoleM.Instance.OpenBoxItems.Count}");
                            //UnityEngine.Debug.Log($"data：{data== null}");
                            //UnityEngine.Debug.Log($"data.OldItem：{data.OldItem == null }");
                            //UnityEngine.Debug.Log($"OldItem.Uid：{data.OldItem.Uid}");
                            if (RoleM.Instance.OpenBoxItems[i].Uid == data.OldItem.Uid)
                                RoleM.Instance.OpenBoxItems.RemoveAt(i);
                            else
                                RoleM.Instance.OpenBoxItems[i] = data.OldItem;

                            break;
                        }
                    }
                    if (isChange)
                    {
                        UniGameEvent.Send(new GameEvent.Refresh_HomeOpen());
                    }
                }
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        //装备 批量分解
        private void Sever_ItemDisassemblyResponse(ItemDisassemblyResponse data)
        {
            // Log.Print($"-装备 批量分解-------{data.GetType().Name}------" + data.GetType().Name);
            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }
            else
            {
                if (RoleM.Instance.OpenBoxItems != null && RoleM.Instance.OpenBoxItems.Count > 0)
                {
                    bool isRemoveBox = false;
                    for (int i = 0; i < data.ItemUids.Length; i++)
                    {
                        for (int j = 0; j < RoleM.Instance.OpenBoxItems.Count; j++)
                        {
                            if (data.ItemUids[i] == RoleM.Instance.OpenBoxItems[j].Uid)
                            {
                                isRemoveBox = true;
                                RoleM.Instance.OpenBoxItems.RemoveAt(j);
                                break;
                            }
                        }
                    }
                    if (isRemoveBox)
                        UniGameEvent.Send(new GameEvent.Refresh_HomeOpen());
                }

                UniGameEvent.Send(new GameEvent.OpenBoxDisassemblyEquip { count = data.ItemUids.Length });
            }
        }

        //装备 自动分解
        private void Sever_AutoDisassemblyResponse(AutoDisassemblyResponse data)
        {
            //Log.Print($"--------{data.GetType().Name}------" + data.GetType().Name);

            if (data.Code == 0)
            {
                BagM.Instance.AutoDelSetInfo.HighCombat = data.HighCombat;
                BagM.Instance.AutoDelSetInfo.LevelSwitch = data.LevelSwitch;
                BagM.Instance.AutoDelSetInfo.Level = data.Level;
                BagM.Instance.AutoDelSetInfo.Open = data.Open;
                BagM.Instance.AutoDelSetInfo.Quality = data.Quality;
                BagM.Instance.AutoDelSetInfo.QualitySwitch = data.QualitySwitch;
                UITool.Close<UIAutoDelSetting>();

                JEvent.defaultEvent.Post(new GameEvent.AutoDelSet());
            }
            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }
        }

        //洗炼 装备
        private void Sever_EquipRefineResponse(EquipRefineResponse data)
        {
            // Log.Print($"--------{"洗炼 装备"}{data.GetType().Name}------");
            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Refresh_EquipRefine { PbArrt = data.NewStats });
            }
        }

        //保留洗炼属性
        private void Sever_EquipRefineReplaceResponse(EquipRefineReplaceResponse data)
        {
            // Log.Print($"--------{"保留洗炼属性"}{data.GetType().Name}------");
            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }
            if (data.Code == 0)
            {
                RoleM.Instance.SetRoleOneEquip(data.Item, (pb.EquipSlot)CFG_Item.GetById(data.Item.CfgId).Equiptype, true);
                JEvent.defaultEvent.Post(new GameEvent.Replace_EquipRefine());
            }
        }

        //强化 装备
        private void Sever_EnhanceResponse(EnhanceResponse data)
        {
            //Log.Print($"--------{"强化"}{data.GetType().Name} data.PartnerSlotId{data.PartnerSlotId} data.Slot{data.Slot} data.NewLevel {data.NewLevel} ------");
            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }
            if (data.Code == 0)
            {
                if (data.PartnerSlotId == 0)
                {
                    RoleM.Instance.SetEquipeReinforeceLv(data.Slot, data.NewLevel);
                }
                else
                    PartnerM.Instance.SetEquipeReinforeceLv(data.PartnerSlotId, (int)data.Slot, data.NewLevel);

                JEvent.defaultEvent.Post(new GameEvent.Replace_EquipEnhance());
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Reinforce);
            }
        }

        //升星 装备
        private void Sever_StarsResponse(StarsResponse data)
        {
            //Log.Print($"--------{"升星"}{data.GetType().Name}  data.PartnerSlotId{data.PartnerSlotId} data.Slot{data.Slot} data.NewLevel {data.NewLevel} ------");
            if (data.Error != "")
            {
                UITool.ShowTips(data.Error);
            }
            if (data.Code >= 0)
            {
                if (data.Code == 0)//升星成功
                    AudioManager.Instance.PlaySound(AudioKey.Sound.UI_StrengthStarSuc);
                else
                    AudioManager.Instance.PlaySound(AudioKey.Sound.UI_StrengthStarFaild);

                if (data.PartnerSlotId == 0)
                    RoleM.Instance.SetEquipeStarLv(data.Slot, data.NewLevel);
                else
                    PartnerM.Instance.SetEquipeStarLv(data.PartnerSlotId, (int)data.Slot, data.NewLevel);

                JEvent.defaultEvent.Post(new GameEvent.Replace_EquipStar { isUp = data.Code == 0 });
            }
        }

        // 同步角色装备信息
        private void Sever_SyncUserEquipsInfoPush(SyncUserEquipsInfoPush data)
        {
            //Log.Print($"--------{"同步角色装备信息"}{data.GetType().Name} {data.Equips.Count}------");
            //UITool.ShowTips("升星成功");
            RoleM.Instance.SetRoleEquipDic(data.Equips);
        }

        //宝石 镶嵌
        private void Sever_GemstoneChangeResponse(GemstoneChangeResponse data)
        {
            //Log.Print($"--------{"宝石 镶嵌"} err{data.Error} {data.GetType().Name}------ data.OpType {data.OpType} data.Code{data.Code} data.ItemUid {data.ItemUid}data.Slot {data.Slot} (int)data.GemstoneSlot {(int)data.GemstoneSlot}");
            if (data.Code == 0)
            {
                if (data.OpType == 1)//镶嵌
                    AudioManager.Instance.PlaySound(AudioKey.Sound.UI_GemEquipToSlot);

                if (data.PartnerSlotId == 0)
                    RoleM.Instance.SetRoleOneEquipInset(data.ItemUid, data.Slot, (int)data.GemstoneSlot, data.OpType == 1);
                else
                    PartnerM.Instance.SetRoleOneEquipInset(data.PartnerSlotId, data.ItemUid, (int)data.Slot, (int)data.GemstoneSlot, data.OpType == 1);
                JEvent.defaultEvent.Post(new GameEvent.EquipInsetRefresh());
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Insert);
            }
            else
            {
                UITool.ShowTips(data.Error);
            }

        }

        //一键镶穿
        private void Sever_GemstoneChangeAllResponse(GemstoneChangeAllResponse data)
        {
            //Log.Print($"--------{"宝石 一键镶嵌"}  {data.GetType().Name}------ data.Code{data.Code} data.Slot {data.Slot} ");
            if (data.Code == 0 && data.GemStoneSlot != null)
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_GemEquipToSlot);
                for (int i = 0; i < data.GemStoneSlot.Length; i++)
                {
                    if (data.PartnerSlotId == 0)
                        RoleM.Instance.SetRoleOneEquipInset(data.GemStoneSlot[i], data.Slot, i, true);
                    else
                        PartnerM.Instance.SetRoleOneEquipInset(data.PartnerSlotId, data.GemStoneSlot[i], (int)data.Slot, i, true);
                }
                JEvent.defaultEvent.Post(new GameEvent.EquipInsetRefresh());
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Insert);
            }
            else
            {
                UITool.ShowTips(data.Error);
            }

        }


        #endregion

        #region 技能

        //同步角色技能槽位
        private void SyncUserSkillSlotPush(SyncUserSkillSlotPush data)
        {
            //Log.Print("同步角色技能槽位----------");
            SkillM.Instance.SetSkillEquipDic(data.SkillsSlot);
            JEvent.defaultEvent.Post(new GameEvent.Sync_UserSkillSlot());
        }

        //新技能
        private void Sever_NewSkillPush(NewSkillPush data)
        {
            SkillM.Instance.AddNewSkillInfo(data.skill);
            SkillM.Instance.SetSkillEquipDic(data.SkillsSlot);
            JEvent.defaultEvent.Post(new GameEvent.Sync_GetNewSkill());

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_UpLv);
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Advance);
        }

        //技能升级
        private void Sever_SkillUpgradeResponse(SkillUpgradeResponse data)
        {
            //Log.Print($"技能升级--------{data.GetType().Name}------{data.Code}  data.PartnerUid {data.PartnerUid}");
            if (data.Code == 0)
            {
                if (data.PartnerUid == 0)
                    SkillM.Instance.SetUpOneSkill(data.OldSkillId, data.NewSkillId);
                else
                    PartnerM.Instance.SetUpOneSkill(data.PartnerUid, data.OldSkillId, data.NewSkillId);

                SkillM.Instance.SetUpOneSkill(data.OldSkillId, data.NewSkillId);

                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_UpLv);
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Advance);
            }
            JEvent.defaultEvent.Post(new GameEvent.Sync_SkillUpgrade { isWin = data.Code == 0 });
            UITool.ShowTips(data.Error);
        }

        //技能修改
        private void Sever_SkillChangeResponse(SkillChangeResponse data)
        {
            //Log.Print($"技能修改--------{data.GetType().Name}------");
            UITool.ShowTips(data.Error);
            SkillM.Instance.SetOneSkillEquip(data.SkillId, data.OpType, (int)data.Slot);
            if (data != null && data.OpType == 1)
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_SkillEquipToSlot);
            }
            JEvent.defaultEvent.Post(new GameEvent.Sync_SkillChange { });
        }

        //技能 升阶
        private void Sever_SkillAdvancedResponse(SkillAdvancedResponse data)
        {
            //Log.Print($"--------{data.GetType().Name}------  data.PartnerUid {data.PartnerUid} data.NewLevel{data.NewLevel}");
            UITool.ShowTips(data.Error);
            if (data.PartnerUid == 0)
                SkillM.Instance.SetUpAdvanceSkill(data.OldSkillId, data.NewLevel);
            else
                PartnerM.Instance.SetUpAdvanceSkill(data.PartnerUid, data.OldSkillId, data.NewLevel);

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Advance);

            JEvent.defaultEvent.Post(new GameEvent.Sync_SkillAdvanced { isWin = data.Code == 0 });
        }

        #endregion

        #region 聊天

        private void Sever_ChatResponse(ChatResponse data)
        {
            // Log.Print($"聊天--------{data.GetType().Name}------{data.Error}{data.Code}");
            UITool.ShowTips(data.Error);
            JEvent.defaultEvent.Post(new GameEvent.Sync_ChatSendBack { isWin = data.Code == 0 });
        }

        private void Sever_ChatPush(ChatPush data)
        {
            // Log.Print($"聊天推送--------{data.GetType().Name}----{data.Type}--");

            if (data.Type != ChatType.PROMPT)
                ChatM.Instance.AddOneChat(data);

            //if (data.Type == ChatType.SCROLLING)
            //    UITool.RunNoticeTips(data.Context);
            //else
            if (data.Type == ChatType.PROMPT)
                UITool.ShowTips(data.Context);

            JEvent.defaultEvent.Post(new GameEvent.Sync_ChatPush() { isOnlyHall = false });
        }

        private void Sever_LampPush(LampPush data)
        {
            //Log.Print($"LampPush--------{data.GetType().Name}----{data.Uid}--");
            UITool.RunNoticeTips(data.Uid, data.Context);
        }



        #endregion

        #region 战宠 宠物

        //宠物数据同步
        private void Sever_SyncUserPetInfoPush(SyncUserPetInfoPush data)
        {
            //Log.Print($"宠物数据同步--------"+ data.PetRecoveryTime);
            PetM.Instance.SetPets(data.Pet);
            PetM.Instance.SetGetedPets(data.PetEverHave);
            PetM.Instance.PetRecoveryTime = data.PetRecoveryTime;
            JEvent.defaultEvent.Post(new GameEvent.Refresh_PetsInfo());

            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_Up.ToString(), PetM.Instance.PetUPRedDot());
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_Award);
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_UpStar_1);
        }

        //宠物移动
        private void Sever_PetMoveResponse(PetMoveResponse data)
        {
            //Log.Print($"宠物移动-------*{data.NewPetCfgId} slot"+ data.SlotId);
            if (data.Code == 0)
                JEvent.defaultEvent.Post(new GameEvent.Refresh_PetMove { petCfgId = data.NewPetCfgId, idx = data.SlotId });
            UITool.ShowTips(data.Error);
        }


        //宠物升星
        private void Sever_PetStarResponse(PetStarResponse data)
        {
            //Log.Print($"宠物升星--------");
            UITool.ShowTips(data.Error);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_PetStar { Code = data.Code, ResPer = data.ResPer });
        }

        #endregion

        #region 玩法信息

        //副本信息
        private void Sever_C2S_InstanceInfo(InstanceInfoResponse data)
        {
            Log.Print($"副本信息--------");
            FightPlayM.Instance.SetInfo(data);
            JEvent.defaultEvent.Post(new GameEvent.Refresh_InstanceInfo());

            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Fight_Click);
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Fight_Award);
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Boss);
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Tower_Click);
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Relice);
            //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_OnHook);
        }

        //游戏结算
        private void Sever_S2C_GameResultBattleMap(GameResultBattlePush data)
        {
            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                GameRoomM.Instance._GameState = GameRoomM.GameState.Over;

            if (data.OverType == GameOverType.GameOverTypeCopySweep && data.Items != null && data.Items.Count > 0)
            {
                List<ItemInfo> items = new List<ItemInfo>();
                foreach (ItemInfo item in data.Items)
                {
                    var find = items.Find(x => x.CfgId == item.CfgId);
                    if (find != null)
                        find.Num += item.Num;
                    else
                        items.Add(item);
                }
                data.Items = items;
            }

            UITool.Open<UIGameFightOver>(null, data);
        }

        //扫荡
        private void Sever_CopySweepResponse(CopySweepResponse data)
        {
            //Log.Print($"--------扫荡 ------ ");
            if (data.Code == 0)
            {
                TcpC.C2S_InstanceInfo();
                TcpC.C2S_DailyMapInfo();
            }
            UITool.ShowTips(data.Error);
        }

        //海岛倍率
        private void Sever_GuardIslandMultiplierResponse(GuardIslandMultiplierResponse data)
        {
            //Log.Print($"--------海岛倍率 ------ ");
            if (data.Code == 0)
            {
                TcpC.C2S_InstanceInfo();
                UITool.Close<UIGuardIslandBox>();
                TcpC.GetGameJoin(CFG_Room.GetByRoomType(CFG_Room.RoomType.PVE_Peace, RoleM.Instance.Level).ID, 0);
            }
        }

        #endregion

        #region 宝箱

        private void Sever_SyncUserBoxInfoPush(SyncUserBoxInfoPush data)
        {
            //Log.Print($"-------- 宝箱信息 ------ ");
            RoleM.Instance.BoxNowlevel = data.Box.Nowlevel;
            RoleM.Instance.OpenBoxItems = data.Box.NowOpenRes;
            RoleM.Instance.BoxNextRecoveryTime = data.Box.NextBoxRecoveryTime;
            //RoleM.Instance.NowOpenLevel = data.Box.NowOpenLevel;
            //RoleM.Instance.NowOpenJob = data.Box.Job;

            RoleM.Instance.OpenBoxItems.Sort((x, y) =>
            {
                if (x.Type == ItemType.EQUIPMENT && y.Type == ItemType.EQUIPMENT)
                {
                    int qualityX = (int)x.Equip.Quality;
                    int qualityY = (int)y.Equip.Quality;
                    return qualityY.CompareTo(qualityX);
                }
                else if (x.Type != ItemType.EQUIPMENT && y.Type != ItemType.EQUIPMENT)
                {
                    int xQuality = (int)CFG_Item.GetById(x.CfgId).Quality;
                    int yQuality = (int)CFG_Item.GetById(y.CfgId).Quality;
                    return yQuality.CompareTo(xQuality);
                }
                else
                {
                    int xType = (int)x.Type;
                    int yType = (int)y.Type;
                    return yType.CompareTo(xType);
                }
            });

            JEvent.defaultEvent.Post(new GameEvent.Sync_UserBoxInfo());
        }

        #endregion

        #region 工会

        // 工会列表
        private void Sever_GuildListResponse(GuildListResponse data)
        {
            if (data.Code == 0)
            {
                if (data.Guilds.Count == 0)
                    return;

                //GuildM.Instance.AddGuildInfos(data.Guilds);
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildList { guildsInfo = data.Guilds });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 创建工会
        private void Sever_GuildCreateResponse(GuildCreateResponse data)
        {
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildCreate());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 申请加入工会
        private void Sever_GuildApplyResponse(GuildApplyResponse data)
        {
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildApply());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 取消申请
        private void Sever_GuildCancelApplyResponse(GuildCancelApplyResponse data)
        {
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildCancelApply());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 解散工会
        private void Sever_GuildDismissResponse(GuildDismissResponse data)
        {
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildDismiss());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 离开工会
        private void Sever_GuildLeaveResponse(GuildLeaveResponse data)
        {
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildLeave());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 踢出工会
        private void Sever_GuildKickResponse(GuildKickResponse data)
        {
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildKick());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 提升工会职位
        private void Sever_GuildPromoteResponse(GuildPromoteResponse data)
        {
            if (data.Code == 0)
            {
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildPromote());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 申请列表
        private void Sever_GuildApplyListResponse(GuildApplyListResponse data)
        {
            GuildM.Instance.applyList = data.Applications;
            JEvent.defaultEvent.Post(new GameEvent.Sync_GuildApplyList { applyList = data.Applications });
            GuildM.Instance.CheckGuildApplyRed();
            if (data.Code == 0)
            {
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 处理申请
        private void Sever_GuildHandleApplyResponse(GuildHandleApplyResponse data)
        {
            JEvent.defaultEvent.Post(new GameEvent.Sync_GuildHandleApply());

            if (data.Code == 0)
            {
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 工会信息
        private void Sever_GuildInfoResponse(GuildInfoResponse data)
        {
            if (data.Code == 0)
            {
                GuildM.Instance.guildInfo = data.Guild;
                GuildM.Instance.roleJob = data.Role;
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildInfo { guildInfo = data.Guild, Role = data.Role });
                // 处理红点
                if (!GuildM.Instance.isCheckedStoreRed)
                {
                    if (data.Role == 1 || data.Role == 2)
                        TcpC.C2S_GuildApplyListRequest();

                    GuildM.Instance.CheckGuildStoreRed();
                    GuildM.Instance.isCheckedStoreRed = true;
                }
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 工会成员列表
        private void Sever_GuildMemberListResponse(GuildMemberListResponse data)
        {
            if (data.Code == 0)
            {
                GuildM.Instance.membersInfo = data.Members;
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildMemberList { membersInfo = data.Members });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 工会捐献信息
        private void Sever_GuildDonateInfoResponse(GuildDonateInfoResponse data)
        {
            if (data.Code == 0)
            {
                GuildM.Instance.donateNum = data.Num;
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildDonateInfo { num = data.Num });
                GuildM.Instance.CheckGuildInfoRed();
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 工会商城购买情况
        private void Sever_GuildShopBuyInfoResponse(GuildShopBuyInfoResponse data)
        {
            if (data.Code == 0)
            {
                GuildM.Instance.buyInfo = data.BuyInfo;
                JEvent.defaultEvent.Post(new GameEvent.Sync_GuildShopBuyInfo { buyInfo = data.BuyInfo });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 工会副本奖励
        private void Sever_GuildCopyRewardPush(GuildCopyRewardPush data)
        {
            //Log.Print("工会副本奖励");
            GameRoomM.Instance.guildCopyReward = data;
            JEvent.defaultEvent.Post(new GameEvent.Refrsh_GuildCopyReward());
        }

        // 观察玩家详情
        private void Sever_ObservePlayerInfoResponse(ObservePlayerInfoResponse data)
        {
            if (data.Code == 0)
            {
                if (data.ObservePlayerInfo == null)
                {
                    Debug.LogError("观察的玩家信息为 null");
                    return;
                }
                UITool.Open<UIOtherRole>(null, data.ObservePlayerInfo);
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        private void Sever_GuidanceUpdateResponse(GuidanceUpdateResponse data)
        {
            UniGameEvent.Send(new GameEvent.Receive_Guide_ServerData { Code = data.Code, Error = data.Error, Guidance = data.Guidance });
        }

        private void Sever_ChatHistoryResponse(ChatHistoryResponse data)
        {
            if (data.Code == 0)
            {
                ChatM.Instance.AddChatHistory(data.ChannelHistories);
                JEvent.defaultEvent.Post(new GameEvent.Sync_ChatPush { isOnlyHall = true });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        #endregion

        #region 回合制玩法

        // 回合制站位信息
        private void Sever_TurnBasedSlotResponse(TurnBasedSlotResponse data)
        {
            if (data.Code == 0)
            {
                TurnBasedM.Instance.pvpFormationUids = data.Uids;
                UniGameEvent.Send(new GameEvent.PVP_TurnBasedSlot { uids = data.Uids });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 回合制挑战
        private void Sever_TurnBasedChallengeResponse(TurnBasedChallengeResponse data)
        {
            if (data.Code == 0)
            {
                TurnBasedM.Instance.turnBasedBattleData = data.Info;

                UniGameEvent.Send(new GameEvent.TurnBasedChallenge());
                JEvent.defaultEvent.Post(new SenceChange
                {
                    scencePath = GameConst.TurnBasedFightPath,
                    CallBack = () =>
                    {
                        TurnBasedM.Instance.roomId = TurnBasedM.turnBasedRoomId;
                        TurnBasedFightCtrl.Instance.OnInit().Forget();
                    }
                });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 五行试炼信息
        private void Sever_ElementInfoResponse(ElementInfoResponse data)
        {
            if (data.Code == 0)
            {
                TurnBasedM.Instance.elementScore = data.ElementScore;
                TurnBasedM.Instance.elementTimes = data.ElementTimes;
                TurnBasedM.Instance.elementMain = data.ElementMain;
                TurnBasedM.Instance.empowers = data.Empower.Elements;
                TurnBasedM.Instance.needSweep = data.NeedSweep;
                TurnBasedM.Instance.wuXingFormationUids = data.Uids.Uids;

                UniGameEvent.Send(new GameEvent.WuXingTrialsInfo());
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.WuXing.ToString(), RedDotM.Instance.GetWuXingChallengeNum());
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 五行元素升级
        public void Sever_ElementUpgradeResponse(ElementUpgradeResponse data)
        {
            if (data.Code == 0)
            {
                if (data.Success == 1)
                {
                    //  Debug.Log($"五行元素升级 主元素 --------{data.ElementMain} / id = {data.SlotId} ");
                    TurnBasedM.Instance.elementMain = data.ElementMain;
                    TurnBasedM.Instance.UpdateElementLevel(data.SlotId, data.NewLevel);

                    AudioManager.Instance.PlaySound(AudioKey.Sound.UI_EmpowermentSuc);
                }
            }

            else
            {
                UITool.ShowTips(data.Error);
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_EmpowermentFaild);
            }
            UniGameEvent.Send(new GameEvent.WuXingFuNengUpdate { isBreak = false });
        }

        /// <summary>
        /// 五行元素突破
        /// </summary>
        /// <param name="data"></param>
        public void Sever_ElementBreakthroughResponse(ElementBreakthroughResponse data)
        {
            if (data.Code == 0)
            {
                //Debug.Log($"五行元素突破--------{data.SlotId} {data.NewLevel}");
                TurnBasedM.Instance.UpdateElementBreakLevel(data.SlotId, data.NewLevel);
                UniGameEvent.Send(new GameEvent.WuXingFuNengUpdate { isBreak = true });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 五行元素出战设置
        public void Sever_ElementPartnerSlotResponse(ElementPartnerSlotResponse data)
        {
            //Debug.Log($"uids {string.Join(", ", data.Uids)}");
            if (data.Code == 0)
            {
                TurnBasedM.Instance.wuXingFormationUids = data.Uids;
                UniGameEvent.Send(new GameEvent.WuXing_TurnBasedSlot { uids = data.Uids });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 五行试炼挑战
        public void Sever_ElementChallengeResponse(ElementChallengeResponse data)
        {
            if (data.Code == 0)
            {
                TurnBasedM.Instance.turnBasedBattleData = data.Info;
                UniGameEvent.Send(new GameEvent.TurnBasedChallenge());
                JEvent.defaultEvent.Post(new SenceChange
                {
                    scencePath = GameConst.TurnBasedFightPath,
                    CallBack = () =>
                    {
                        TurnBasedM.Instance.roomId = TurnBasedM.wuXingRoomId;
                        TurnBasedFightCtrl.Instance.OnInit().Forget();
                    }
                });
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        // 五行试炼，单层扫荡
        public void Sever_ElementSweepForWardResponse(ElementSweepForWardResponse data)
        {
            if (data.Code == 0)
            {
                TcpC.C2S_ElementInfoRequest();

                List<ItemInfo> awards = new List<ItemInfo>();
                var items = CFG_ElementTrials.GetById(TurnBasedM.Instance.elementScore).Settlement;
                foreach (var item in items)
                {
                    ItemInfo info = new ItemInfo
                    {
                        CfgId = item.id,
                        Num = item.num * data.Num,
                    };
                    awards.Add(info);
                }

                var gameResult = new GameResultBattlePush
                {
                    RoomId = TurnBasedM.Instance.roomId,
                    BattleId = 0,
                    OverType = GameOverType.GameOverTypeCopySweep,
                    Items = awards
                };
                UITool.Open<UIGameFightOver>(null, gameResult);
            }
            else
            {
                UITool.ShowTips(data.Error);
            }
        }

        #endregion

        #region 称号

        //称号列表
        private void Sever_TitleListResponse(TitleListResponse data)
        {
            //Log.Print($"  称号列表 {data.Titles.Count}");
            if (data.Code == 0)
                TitleM.Instance.SetTitles(data.Titles);
            else
                UITool.ShowTips(data.Error);

            JEvent.defaultEvent.Post(new GameEvent.Refresh_TitleList());
        }

        //称号合成
        private void Sever_TitleComposeResponse(TitleComposeResponse data)
        {
            if (data.Code == 0)
            {
                TcpC.C2S_TitleList();
            }
            else
                UITool.ShowTips(data.Error);
        }

        //称号升星
        private void Sever_TitleStartResponse(TitleStartResponse data)
        {
            if (data.Code == 0)
            {
                UITool.ShowTips(data.Star != TitleM.Instance.GetTitleStar(data.CfgId) ? "升星成功" : "升星失败");
                TitleM.Instance.SetTitleStar(data.CfgId, data.Star);
                JEvent.defaultEvent.Post(new GameEvent.Title_UpStar());
            }
            else
                UITool.ShowTips(data.Error);

        }

        //装备称号
        private void Sever_TitleEquipResponse(TitleEquipResponse data)
        {
            if (data.Code != 0)
                UITool.ShowTips(data.Error);
            else
            {
                TcpC.C2S_TitleList();
                JEvent.defaultEvent.Post(new GameEvent.Title_Equip());
            }
        }

        #endregion

        #region 召唤
        private void Sever_LotteryInfoResponse(LotteryInfoResponse resq)
        {
            if (resq == null || resq.LotteryInfos == null) return;

            if (resq.LotteryInfos.FreeNum == 1)
            {
                //bool isLockServer = GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Role_SkillAndCall, false);
                //bool islockLv = GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Role_SkillAndCall, false);

                bool isLockServer = GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.CallRole, false);
                bool islockLv = GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.CallRole, false);
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Call_Free.ToString(), !isLockServer && !islockLv ? 1 : 0);
            }
        }
        #endregion

        #region 新功能开启
        private void Server_ModuleOpenRecordResponse(ModuleOpenRecordResponse response)
        {
            //不需要这里监听信息，数据放在了 Sever_SyncUserBasePush

            FuncOpenShowWidget.Instance.InitRecord(response.ModuleIds);
        }
        #endregion

        #region 洗练
        private void Server_RefineResponse(RefineResponse data)
        {
            if (data.code == 0)
                RoleM.Instance.refineResuleCache = data;
            else
                UITool.ShowTips(data.error);

            UniGameEvent.Send(new GameEvent.EquipSlotRefine { success = data.code == 0 });

        }

        private void Server_RefineInfoResponse(RefineInfoResponse data)
        {
            if (data.code == 0)
            {
                RoleM.Instance.characterRefineDatas = data.character_refines;
                RoleM.Instance.purchasedRefineSlots = data.purchased_refine_slots;
                UniGameEvent.Send(new GameEvent.EquipSlotRefineInfo());
            }
            else
                UITool.ShowTips(data.error);
        }

        private void Server_RefineSelectResponse(RefineSelectResponse data)
        {
            if (data.code == 0)
            {
                RoleM.Instance.UpdateEquipRefineData(data.character_slot, data.equip_slot, data.refined_slot, data.selected_attributes);
                UniGameEvent.Send(new GameEvent.EquipSlotRefineSelectResult { refinedSlot = data.refined_slot });
            }
            else
                UITool.ShowTips(data.error);
        }
        #endregion

        #region 头像和头像框

        private void Server_SyncUserAvatarPush(SyncUserAvatarPush data)
        {
            RoleM.Instance.avatars = data.avatars;
            RoleM.Instance.frames = data.frames;
            RoleM.Instance.CurrentAvatarId = data.current_avatar_id;
            RoleM.Instance.CurrentFrameId = data.current_frame_id;
            UniGameEvent.Send(new GameEvent.UserAvatarPush());
        }

        private void Server_UseAvatarResponse(UseAvatarResponse data)
        {
            switch (data.type)
            {
                case AvatarType.AVATAR:
                    RoleM.Instance.CurrentAvatarId = data.id;
                    break;
                case AvatarType.FRAME:
                    RoleM.Instance.CurrentFrameId = data.id;
                    break;
                default:
                    break;
            }

            UniGameEvent.Send(new GameEvent.Refresh_UseAvatar { avatarType = data.type });
        }

        #endregion

        #region 幸运金币
        private void Server_FateCoinResponse(FateCoinResponse data)
        {
            if (UITool.Check<UIGameLuckCoin>())
            {

                if (data.Code != 0)
                    UITool.Close<UIGameLuckCoin>();
                else if (data.CfgId != 0)
                    UITool.Get<UIGameLuckCoin>().OnCoinOpne(data.CfgId);
                //UITool.ShowTips("角色死亡状态下无法获取命运效果");
            }
        }

        private void Server_FateCoinPush(FateCoinPush data)
        {
            if (UITool.Check<UIGameLuckCoin>())
                UITool.Close<UIGameLuckCoin>();
            UITool.Open<UIGameLuckCoin>();
        }

        #endregion

        #region 通用购买返回通知
        private void Server_CommonNotificationPush(CommonNotificationPush data)
        {
            UniGameEvent.Send(new GameEvent.CommonNotification { moduleId = data.ModuleID });
        }
        #endregion

        #region 功能开启红点
        private void Server_ModuleNewPush(ModuleNewPush push)
        {
            //string str = "Server_ModuleNewPush Call Back";
            //if (push.Infos != null)
            //    Debug.LogWarning(str + push.Infos.Count);

            FuncOpenShowWidget.Instance.InitFunOpenFlag(push);
        }

        private void Server_ModuleNewResponse(ModuleNewResponse response)
        {
            if (response.ErrorCode != 0) return;

            if (response == null) return;

            FuncOpenShowWidget.Instance.OperateFunOpenFlag(response.Info);
        }
        #endregion
    }
}
