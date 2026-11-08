package client

import (
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/rs/zerolog/log"
)

// initResponseHandlers initializes the response handlers for the client.
// func (c *Client) initResponseHandlers() {
// 	c.staticResponseHandlers = map[uint64]func([]byte){
// 		uint64(pb.CmdId_C2S_Login):           c.makeResponseHandler(func() proto.Message { return &pb.LoginToGameServerResponse{} }, c.LoginToGameServerResponse),
// 		uint64(pb.CmdId_C2S_Join):            c.makeResponseHandler(func() proto.Message { return &pb.JoinDeskResponse{} }, c.JoinDeskResponse),
// 		uint64(pb.CmdId_C2S_Exit):            c.makeResponseHandler(func() proto.Message { return &pb.ExitResponse{} }, c.ExitResponse),
// 		uint64(pb.CmdId_C2S_Move):            c.makeResponseHandler(func() proto.Message { return &pb.MoveResponse{} }, c.MoveResponse),
// 		uint64(pb.CmdId_C2S_Attack):          c.makeResponseHandler(func() proto.Message { return &pb.AttackResponse{} }, c.AttackResponse),
// 		uint64(pb.CmdId_C2S_GetBag):          c.makeResponseHandler(func() proto.Message { return &pb.GetBagResponse{} }, c.GetBagResponse),
// 		uint64(pb.CmdId_C2S_ItemUse):         c.makeResponseHandler(func() proto.Message { return &pb.ItemUseResponse{} }, c.ItemUseResponse),
// 		uint64(pb.CmdId_C2S_ItemDestroy):     c.makeResponseHandler(func() proto.Message { return &pb.ItemDestroyResponse{} }, c.ItemDestroyResponse),
// 		uint64(pb.CmdId_C2S_ItemDisassembly): c.makeResponseHandler(func() proto.Message { return &pb.ItemDisassemblyResponse{} }, c.ItemDisassemblyResponse),
// 		uint64(pb.CmdId_C2S_InstanceTimeAdd): c.makeResponseHandler(func() proto.Message { return &pb.InstanceTimeAddResponse{} }, c.InstanceTimeAddResponse),
// 		uint64(pb.CmdId_C2S_Revive):          c.makeResponseHandler(func() proto.Message { return &pb.ReviveResponse{} }, c.ReviveResponse),
// 		uint64(pb.CmdId_C2S_ChangeEquip):     c.makeResponseHandler(func() proto.Message { return &pb.ChangeEquipResponse{} }, c.ChangeEquipResponse),
// 		uint64(pb.CmdId_C2S_TaskComplate):    c.makeResponseHandler(func() proto.Message { return &pb.TaskComplateResponse{} }, c.TaskCompleteResponse),
// 		uint64(pb.CmdId_C2S_MatchingRequest): c.makeResponseHandler(func() proto.Message { return &pb.MatchingResponse{} }, c.MatchingResponse),
// 		// uint64(pb.CmdId_C2S_MatchingReady):         c.makeResponseHandler(func() proto.Message { return &pb.MatchingReadyResponse{} }, c.MatchingReadyResponse),
// 		uint64(pb.CmdId_C2S_ChangeRoleName):      c.makeResponseHandler(func() proto.Message { return &pb.ChangeRoleNameResponse{} }, c.ChangeRoleNameResponse),
// 		uint64(pb.CmdId_C2S_Enhance):             c.makeResponseHandler(func() proto.Message { return &pb.EnhanceResponse{} }, c.EnhanceResponse),
// 		uint64(pb.CmdId_C2S_Stars):               c.makeResponseHandler(func() proto.Message { return &pb.StarsResponse{} }, c.StarsResponse),
// 		uint64(pb.CmdId_C2S_EquipRefine):         c.makeResponseHandler(func() proto.Message { return &pb.EquipRefineResponse{} }, c.EquipRefineResponse),
// 		uint64(pb.CmdId_C2S_EquipRefineReplace):  c.makeResponseHandler(func() proto.Message { return &pb.EquipRefineReplaceResponse{} }, c.EquipRefineReplaceResponse),
// 		uint64(pb.CmdId_C2S_Chat):                c.makeResponseHandler(func() proto.Message { return &pb.ChatResponse{} }, c.ChatResponse),
// 		uint64(pb.CmdId_C2S_SkillChange):         c.makeResponseHandler(func() proto.Message { return &pb.SkillChangeResponse{} }, c.SkillChangeResponse),
// 		uint64(pb.CmdId_C2S_SkillUpgrade):        c.makeResponseHandler(func() proto.Message { return &pb.SkillUpgradeResponse{} }, c.SkillUpgradeResponse),
// 		uint64(pb.CmdId_C2S_SkillAdvanced):       c.makeResponseHandler(func() proto.Message { return &pb.SkillAdvancedResponse{} }, c.SkillAdvancedResponse),
// 		uint64(pb.CmdId_C2S_EndlesstrialsReward): c.makeResponseHandler(func() proto.Message { return &pb.EndlesstrialsRewardResponse{} }, c.EndlessTrialsRewardResponse),
// 		uint64(pb.CmdId_C2S_StarsLevel):          c.makeResponseHandler(func() proto.Message { return &pb.StarsLevelResponse{} }, c.StarsLevelResponse),
// 		uint64(pb.CmdId_C2S_PVPInfo):             c.makeResponseHandler(func() proto.Message { return &pb.PVPInfoResponse{} }, c.PVPInfoResponse),
// 		uint64(pb.CmdId_C2S_Compose):             c.makeResponseHandler(func() proto.Message { return &pb.ComposeResponse{} }, c.ComposeResponse),
// 		uint64(pb.CmdId_C2S_CheckBoss):           c.makeResponseHandler(func() proto.Message { return &pb.CheckBossResponse{} }, c.CheckBossResponse),
// 		uint64(pb.CmdId_C2S_GemstoneChange):      c.makeResponseHandler(func() proto.Message { return &pb.GemstoneChangeResponse{} }, c.GemstoneChangeResponse),
// 		uint64(pb.CmdId_C2S_DoubleReward):        c.makeResponseHandler(func() proto.Message { return &pb.DoubleRewardResponse{} }, c.DoubleRewardResponse),
// 		uint64(pb.CmdId_C2S_CopySweep):           c.makeResponseHandler(func() proto.Message { return &pb.CopySweepResponse{} }, c.CopySweepResponse),
// 		uint64(pb.CmdId_C2S_GemstoneChangeAll):   c.makeResponseHandler(func() proto.Message { return &pb.GemstoneChangeAllResponse{} }, c.GemstoneChangeAllResponse),
// 		uint64(pb.CmdId_C2S_Mining):              c.makeResponseHandler(func() proto.Message { return &pb.MiningResponse{} }, c.MiningResponse),
// 		uint64(pb.CmdId_C2S_DailyMapInfo):        c.makeResponseHandler(func() proto.Message { return &pb.DailyMapInfoResponse{} }, c.DailyMapInfoResponse),
// 		uint64(pb.CmdId_C2S_AutoPlay):            c.makeResponseHandler(func() proto.Message { return &pb.AutoPlayResponse{} }, c.AutoPlayResponse),
// 		uint64(pb.CmdId_C2S_AutoDisassembly):     c.makeResponseHandler(func() proto.Message { return &pb.AutoDisassemblyResponse{} }, c.AutoDisassemblyResponse),
// 		uint64(pb.CmdId_C2S_HomeOpenBox):         c.makeResponseHandler(func() proto.Message { return &pb.HomeOpenBoxResponse{} }, c.HomeOpenBoxResponse),
// 		uint64(pb.CmdId_C2S_UpgradeLevel):        c.makeResponseHandler(func() proto.Message { return &pb.UpgradeLevelResponse{} }, c.UpgradeLevelResponse),
// 		uint64(pb.CmdId_C2S_BoxResIntoBag):       c.makeResponseHandler(func() proto.Message { return &pb.BoxResIntoBagResponse{} }, c.BoxResIntoBagResponse),
// 		uint64(pb.CmdId_C2S_PartnerStar):         c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerStarResponse),
// 		uint64(pb.CmdId_C2S_PartnerCombat):       c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerCombatResponse),
// 		uint64(pb.CmdId_C2S_PartnerRest):         c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerRestResponse),
// 		uint64(pb.CmdId_C2S_PartnerList):         c.makeResponseHandler(func() proto.Message { return &pb.PartnerListResponse{} }, c.PartnerListResponse),
// 		uint64(pb.CmdId_C2S_PartnerSlotInfo):     c.makeResponseHandler(func() proto.Message { return &pb.PartnerSlotInfoResponse{} }, c.PartnerSlotInfoResponse),
// 		uint64(pb.CmdId_C2S_EquipBest):           c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.EquipBestResponse),
// 		// uint64(pb.CmdId_C2S_PetCombat):             c.makeResponseHandler(func() proto.Message { return &pb.PetCombatResponse{} }, c.PetCombatResponse),
// 		// uint64(pb.CmdId_C2S_PetRest):               c.makeResponseHandler(func() proto.Message { return &pb.PetRestResponse{} }, c.PetRestResponse),
// 		uint64(pb.CmdId_C2S_PetStar):               c.makeResponseHandler(func() proto.Message { return &pb.PetStarResponse{} }, c.PetStarResponse),
// 		uint64(pb.CmdId_C2S_PetLevel):              c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PetLevelResponse),
// 		uint64(pb.CmdId_C2S_PetMove):               c.makeResponseHandler(func() proto.Message { return &pb.PetMoveResponse{} }, c.PetMoveResponse),
// 		uint64(pb.CmdId_C2S_InstanceInfo):          c.makeResponseHandler(func() proto.Message { return &pb.InstanceInfoResponse{} }, c.InstanceInfoResponse),
// 		uint64(pb.CmdId_C2S_BattleReward):          c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.BattleRewardResponse),
// 		uint64(pb.CmdId_C2S_PartnerCompose):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerComposeResponse),
// 		uint64(pb.CmdId_C2S_HomeUpgradeBox):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.HomeUpgradeBoxResponse),
// 		uint64(pb.CmdId_C2S_HomeUpgradeBoxQuick):   c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.HomeUpgradeBoxQuickResponse),
// 		uint64(pb.CmdId_C2S_GuildList):             c.makeResponseHandler(func() proto.Message { return &pb.GuildListResponse{} }, c.GuildListResponse),
// 		uint64(pb.CmdId_C2S_GuildCreate):           c.makeResponseHandler(func() proto.Message { return &pb.GuildCreateResponse{} }, c.GuildCreateResponse),
// 		uint64(pb.CmdId_C2S_GuildApply):            c.makeResponseHandler(func() proto.Message { return &pb.GuildApplyResponse{} }, c.GuildApplyResponse),
// 		uint64(pb.CmdId_C2S_GuildCancelApply):      c.makeResponseHandler(func() proto.Message { return &pb.GuildCancelApplyResponse{} }, c.GuildCancelApplyResponse),
// 		uint64(pb.CmdId_C2S_GuildDismiss):          c.makeResponseHandler(func() proto.Message { return &pb.GuildDismissResponse{} }, c.GuildDismissResponse),
// 		uint64(pb.CmdId_C2S_GuildLeave):            c.makeResponseHandler(func() proto.Message { return &pb.GuildLeaveResponse{} }, c.GuildLeaveResponse),
// 		uint64(pb.CmdId_C2S_GuildKick):             c.makeResponseHandler(func() proto.Message { return &pb.GuildKickResponse{} }, c.GuildKickResponse),
// 		uint64(pb.CmdId_C2S_GuildPromote):          c.makeResponseHandler(func() proto.Message { return &pb.GuildPromoteResponse{} }, c.GuildPromoteResponse),
// 		uint64(pb.CmdId_C2S_GuildApplyList):        c.makeResponseHandler(func() proto.Message { return &pb.GuildApplyListResponse{} }, c.GuildApplyListResponse),
// 		uint64(pb.CmdId_C2S_GuildHandleApply):      c.makeResponseHandler(func() proto.Message { return &pb.GuildHandleApplyResponse{} }, c.GuildHandleApplyResponse),
// 		uint64(pb.CmdId_C2S_GuildInfo):             c.makeResponseHandler(func() proto.Message { return &pb.GuildInfoResponse{} }, c.GuildInfoResponse),
// 		uint64(pb.CmdId_C2S_GuildMemberList):       c.makeResponseHandler(func() proto.Message { return &pb.GuildMemberListResponse{} }, c.GuildMemberListResponse),
// 		uint64(pb.CmdId_C2S_GuildDonateInfo):       c.makeResponseHandler(func() proto.Message { return &pb.GuildDonateInfoResponse{} }, c.GuildDonateInfoResponse),
// 		uint64(pb.CmdId_C2S_GuildDonate):           c.makeResponseHandler(func() proto.Message { return &pb.GuildDonateInfoResponse{} }, c.GuildDonateResponse),
// 		uint64(pb.CmdId_C2S_GuildShopBuyInfo):      c.makeResponseHandler(func() proto.Message { return &pb.GuildShopBuyInfoResponse{} }, c.GuildShopBuyInfoResponse),
// 		uint64(pb.CmdId_C2S_GuildShopBuy):          c.makeResponseHandler(func() proto.Message { return &pb.GuildShopBuyInfoResponse{} }, c.GuildShopBuyResponse),
// 		uint64(pb.CmdId_C2S_GuildUpgradeSkill):     c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.GuildUpgradeSkillResponse),
// 		uint64(pb.CmdId_C2S_ObservePlayerInfo):     c.makeResponseHandler(func() proto.Message { return &pb.ObservePlayerInfoResponse{} }, c.ObservePlayerInfoResponse),
// 		uint64(pb.CmdId_C2S_HomeBoxSetting):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.HomeBoxSettingResponse),
// 		uint64(pb.CmdId_C2S_PartnerReforge):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerReforgeResponse),
// 		uint64(pb.CmdId_C2S_PartnerDisassemble):    c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerDisassembleResponse),
// 		uint64(pb.CmdId_C2S_PvpMatchingRequest):    c.makeResponseHandler(func() proto.Message { return &pb.PvpMatchingResponse{} }, c.PvpMatchingResponse),
// 		uint64(pb.CmdId_C2S_PVPChallengeRequest):   c.makeResponseHandler(func() proto.Message { return &pb.PVPChallengeResponse{} }, c.PVPChallengeResponse),
// 		uint64(pb.CmdId_C2S_PartnerReplace):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerReplaceResponse),
// 		uint64(pb.CmdId_C2S_PvpRecoder):            c.makeResponseHandler(func() proto.Message { return &pb.PvpRecorderResponse{} }, c.PvpRecorderResponse),
// 		uint64(pb.CmdId_C2S_ActivityRewardReceive): c.makeResponseHandler(func() proto.Message { return &pb.ActivityRewardReceiveResponse{} }, c.ActivityRewardReceiveResponse),
// 		uint64(pb.CmdId_C2S_GuidanceUpdate):        c.makeResponseHandler(func() proto.Message { return &pb.GuidanceUpdateResponse{} }, c.GuidanceUpdateResponse),
// 		uint64(pb.CmdId_C2S_TurnBasedSlot):         c.makeResponseHandler(func() proto.Message { return &pb.TurnBasedSlotResponse{} }, c.TurnBasedSlotResponse),
// 		uint64(pb.CmdId_C2S_TurnBasedChallenge):    c.makeResponseHandler(func() proto.Message { return &pb.TurnBasedChallengeResponse{} }, c.TurnBasedChallengeResponse),
// 		uint64(pb.CmdId_C2S_ChatHistory):           c.makeResponseHandler(func() proto.Message { return &pb.ChatHistoryResponse{} }, c.ChatHistoryResponse),
// 		uint64(pb.CmdId_C2S_GuardIslandMultiplier): c.makeResponseHandler(func() proto.Message { return &pb.GuardIslandMultiplierResponse{} }, c.GuardIslandMultiplierResponse),
// 		uint64(pb.CmdId_C2S_ElementInfo):           c.makeResponseHandler(func() proto.Message { return &pb.ElementInfoResponse{} }, c.ElementInfoResponse),
// 	}
// }

// const (
// 	RouteManagerCheckOrder            = "Manager.CheckOrder"
// 	RouteManagerLogin                 = "Manager.Login"
// 	RouteManagerGetBag                = "Manager.GetBag"
// 	RouteManagerItemUse               = "Manager.ItemUse"
// 	RouteManagerItemDestroy           = "Manager.ItemDestroy"
// 	RouteManagerItemDisassembly       = "Manager.ItemDisassembly"
// 	RouteManagerChangeEquip           = "Manager.ChangeEquip"
// 	RouteManagerTaskComplate          = "Manager.TaskComplate"
// 	RouteManagerChangeRoleName        = "Manager.ChangeRoleName"
// 	RouteManagerEquipRefine           = "Manager.EquipRefine"
// 	RouteManagerEquipRefineReplace    = "Manager.EquipRefineReplace"
// 	RouteManagerChat                  = "Manager.Chat"
// 	RouteManagerEnhance               = "Manager.Enhance"
// 	RouteManagerStars                 = "Manager.Stars"
// 	RouteManagerSkillsChange          = "Manager.SkillChange"
// 	RouteManagerSkillsUpgrade         = "Manager.SkillUpgrade"
// 	RouteManagerSkillsAdvanced        = "Manager.SkillAdvanced"
// 	RouteManagerStarsLevel            = "Manager.StarsLevel"
// 	RouteManagerCompose               = "Manager.Compose"
// 	RouteManagerGemstone              = "Manager.GemstoneChange"
// 	RouteManagerGemstoneAll           = "Manager.GemstoneChangeAll"
// 	RouteManagerAutoDisassembly       = "Manager.AutoDisassembly"
// 	RouteHomeOpenBox                  = "Manager.HomeOpenBox"
// 	RouteManagerBoxResIntoBag         = "Manager.BoxResIntoBag"
// 	RouteManagerUpgradeLevel          = "Manager.UpgradeLevel"
// 	RouteManagerPartnerStar           = "Manager.PartnerStar"
// 	RouteManagerPartnerCombat         = "Manager.PartnerCombat"
// 	RouteManagerPartnerRest           = "Manager.PartnerRest"
// 	RouteManagerPartnerReplace        = "Manager.PartnerReplace"
// 	RouteManagerPartnerList           = "Manager.PartnerList"
// 	RouteManagerPartnerSlotInfo       = "Manager.PartnerSlotInfo"
// 	RouteManagerPartnerEquipBest      = "Manager.EquipBest"
// 	RouteManagerInstanceInfo          = "Manager.InstanceInfo"
// 	RouteManagerPartnerCompose        = "Manager.PartnerCompose"
// 	RouteManagerHomeUpgradeBox        = "Manager.HomeUpgradeBox"
// 	RouteManagerHomeUpgradeBoxQuick   = "Manager.HomeUpgradeBoxQuick"
// 	RouteManagerObservePlayerInfo     = "Manager.ObservePlayerInfo"
// 	RouteManagerHomeBoxSetting        = "Manager.HomeBoxSetting"
// 	RouteManagerPartnerReforge        = "Manager.PartnerReforge"
// 	RouteManagerPartnerDisassemble    = "Manager.PartnerDisassemble"
// 	RouteManagerActivityRewardReceive = "Manager.ActivityRewardReceive"
// 	RouteManagerGuidanceUpdate        = "Manager.GuidanceUpdate"
// 	RouteChatHistory                  = "Manager.ChatHistory"
// 	RouteGuardIslandMultiplier        = "Manager.GuardIslandMultiplier"

// 	RouteManagerPetStar                 = "Manager.PetStar"
// 	RouteManagerPetLevel                = "Manager.PetLevel"
// 	RouteManagerPetMove                 = "Manager.PetMove"
// 	RouteManagerGuildList               = "GuildManager.GuildList"
// 	RouteManagerGuildCreate             = "GuildManager.GuildCreate"
// 	RouteManagerGuildApply              = "GuildManager.GuildApply"
// 	RouteManagerGuildCancelApply        = "GuildManager.GuildCancelApply"
// 	RouteManagerGuildDismiss            = "GuildManager.GuildDismiss"
// 	RouteManagerGuildLeave              = "GuildManager.GuildLeave"
// 	RouteManagerGuildKick               = "GuildManager.GuildKick"
// 	RouteManagerGuildPromote            = "GuildManager.GuildPromote"
// 	RouteManagerGuildApplyList          = "GuildManager.GuildApplyList"
// 	RouteManagerGuildHandleApply        = "GuildManager.GuildHandleApply"
// 	RouteManagerGuildInfo               = "GuildManager.GuildInfo"
// 	RouteManagerGuildMemberList         = "GuildManager.GuildMemberList"
// 	RouteManagerGuildDonateInfo         = "GuildManager.GuildDonateInfo"
// 	RouteManagerGuildDonate             = "GuildManager.GuildDonate"
// 	RouteManagerGuildShopBuyInfo        = "GuildManager.GuildShopBuyInfo"
// 	RouteManagerGuildShopBuy            = "GuildManager.GuildShopBuy"
// 	RouteManagerGuildUpgradeSkill       = "GuildManager.GuildUpgradeSkill"
// 	RouteRoomManagerJoin                = "RoomManager.Join"
// 	RouteRoomManagerExit                = "RoomManager.Exit"
// 	RouteRoomManagerMove                = "RoomManager.Move"
// 	RouteRoomManagerAttack              = "RoomManager.Attack"
// 	RouteRoomManagerInstanceTimeAdd     = "RoomManager.InstanceTimeAdd"
// 	RouteRoomManagerRevive              = "RoomManager.Revive"
// 	RouteRoomManagerEndlesstrialsReward = "RoomManager.EndlesstrialsReward"
// 	RouteRoomManagerDoubleReward        = "RoomManager.DoubleReward"
// 	RouteRoomManagerCopySweep           = "RoomManager.CopySweep"
// 	RouteRoomManagerMining              = "RoomManager.Mining"
// 	RouteRoomManagerDailyMapInfo        = "RoomManager.DailyMapInfo"
// 	RouteRoomManagerAutoPlay            = "RoomManager.AutoPlay"
// 	RouteRoomManagerBattleReward        = "RoomManager.BattleReward"

// 	RouteMatchManagerMatching    = "MatchManager.Matching"
// 	RouteMatchManagerReady       = "MatchManager.Ready"
// 	RouteMatchPVPInfo            = "MatchManager.PVPInfo"
// 	RouteMatchPvpMatching        = "MatchManager.PvpMatching"
// 	RouteMatchPVPChallenge       = "MatchManager.PVPChallenge"
// 	RouteMatchPvpRecoder         = "MatchManager.PvpRecoder"
// 	RouteMatchTurnBasedSlot      = "MatchManager.TurnBasedSlot"
// 	RouteMatchTurnBasedChallenge = "MatchManager.TurnBasedChallenge"

// 	RouteCheckBoss = "WorldBossManager.CheckBoss"

// 	RoutePlayerEnterDeskPush            = "PlayerEnterDeskPush"
// 	RouteExitPush                       = "ExitPush"
// 	RouteSyncUserBasePush               = "SyncUserBasePush"
// 	SyncUserSkillSlotPush               = "SyncUserSkillSlotPush"
// 	RouteMovePush                       = "MovePush"
// 	RouteAttackPush                     = "AttackPush"
// 	RouteMonsterInPush                  = "MonsterInPush"
// 	RouteMonsterOutPush                 = "MonsterOutPush"
// 	RouteObjectHpChangePush             = "ObjectHpChangePush"
// 	RouteObjectDiePush                  = "ObjectDiePush"
// 	RouteAddSourcePush                  = "AddSourcePush"
// 	RouteTakeSourcePush                 = "TakeSourcePush"
// 	RouteLevelUpPush                    = "LevelUpPush"
// 	RouteSkillCooldownPush              = "SkillCooldownPush"
// 	RouteBuffListPush                   = "BuffListPush"
// 	RouteAddBuffPush                    = "AddBuffPush"
// 	RouteRemoveBuffPush                 = "RemoveBuffPush"
// 	RouteSyncUserStatsPush              = "SyncUserStatsPush"
// 	RouteSyncUserAdvanceInfoPush        = "SyncUserAdvanceInfoPush"
// 	RouteReJoinDeskPush                 = "ReJoinDeskPush"
// 	RouteInstanceTimeAddPush            = "InstanceTimeAddPush"
// 	RouteTaskAcceptPush                 = "TaskAcceptPush"
// 	RouteTaskRemovePush                 = "TaskRemovePush"
// 	RouteTaskUpdatePush                 = "TaskUpdatePush"
// 	RouteMatchingResultPush             = "MatchingResultPush"
// 	RouteMatchingStartPush              = "MatchingStartPush"
// 	RouteMatchRoomJoinPush              = "MatchRoomJoinPush"
// 	RouteGameResultMatchPush            = "GameResultMatchPush"
// 	RouteNewSkill                       = "NewSkillPush"
// 	RouteChatPush                       = "ChatPush"
// 	RouteGameResultEndlesstrialsPush    = "GameResultEndlesstrialsPush"
// 	RouteSyncUserEquipsInfoPush         = "SyncUserEquipsInfoPush"
// 	RouteRevivePush                     = "RevivePush"
// 	RouteItemDropPush                   = "ItemDropPush"
// 	RouteGameResultMatchMeleePush       = "GameResultMatchMeleePush"
// 	RouteGameObjectAddPush              = "GameObjectAddPush"
// 	RouteGameObjectDelPush              = "GameObjectDelPush"
// 	RoutePVPMeleeKillPush               = "PVPMeleeKillPush"
// 	RouteTeleportPush                   = "TeleportPush"
// 	RouteSystemNotificationPush         = "SystemNotificationPush"
// 	RouteGameResultDailyMapPush         = "GameResultDailyMapPush"
// 	RouteDailyMapScorePush              = "DailyMapScorePush"
// 	RouteLampPush                       = "LampPush"
// 	RouteHomeAutoEncounterPush          = "HomeAutoEncounterPush"
// 	RoutePartnerSlotUpdatePush          = "PartnerSlotUpdatePush"
// 	RoutePartnerUpdatePush              = "PartnerUpdatePush"
// 	RouteSyncUserPetInfoPush            = "SyncUserPetInfoPush"
// 	GameResultBattleMapPush             = "GameResultBattleMapPush"
// 	RouteSyncUserBoxInfoPush            = "SyncUserBoxInfoPush"
// 	RouteSyncPartnerListPush            = "PartnerListPush"
// 	RouteGuildCopyRewardPush            = "GuildCopyRewardPush"
// 	RouteActivityCumulativeRechargePush = "ActivityCumulativeRechargePush"
// 	RouteCompulsoryWithdrawalPush       = "CompulsoryWithdrawalPush"
// )

func (c *Client) TestLogic() {
	// God
	chatRequest := &pb.ChatRequest{Type: pb.ChatType_WORLD, Context: "@god"}
	if err := c.SendMsg(protocol.RouteManagerChat, chatRequest); err != nil {
		log.Error().Err(err).Msg("failed to send ChatRequest")
		return
	}

	// 挂机副本
	// joinDeskRequest := &pb.JoinDeskRequest{RoomId: 10, BattleId: 0}
	// if err := c.SendMsg(protocol.RouteRoomManagerJoin, joinDeskRequest); err != nil {
	// 	log.Error().Err(err).Msg("failed to send JoinDeskRequest")
	// 	return
	// }

	// 世界BOSS room.xlsx
	// joinBossDeskRequest := &pb.JoinDeskRequest{RoomId: 100, BattleId: 0}
	// if err := c.SendMsg(protocol.RouteRoomManagerJoin, joinBossDeskRequest); err != nil {
	// 	log.Error().Err(err).Msg("failed to send JoinDeskRequest for World Boss")
	// 	return
	// }

	// CheckBoss
	// checkBossRequest := &pb.CheckBossRequest{
	// 	BossCfgId: 500, // 世界BOSS配置ID
	// }
	// if err := c.SendMsg(protocol.RouteCheckBoss, checkBossRequest); err != nil {
	// 	log.Error().Err(err).Msg("failed to send CheckBossRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGetBag, &pb.ChatRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GetBagRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerItemUse, &pb.ItemUseRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ItemUseRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerItemDestroy, &pb.ItemDestroyRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ItemDestroyRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerItemDisassembly, &pb.ItemDisassemblyRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ItemDisassemblyRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerChangeEquip, &pb.ChangeEquipRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ChangeEquipRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerTaskComplate, &pb.TaskComplateRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send TaskComplateRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerChangeRoleName, &pb.ChangeRoleNameRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ChangeRoleNameRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerEquipRefine, &pb.EquipRefineRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send EquipRefineRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerEquipRefineReplace, &pb.EquipRefineReplaceRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send EquipRefineReplaceRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerChat, &pb.ChatRequest{Type: pb.ChatType_WORLD, Context: "Hello, world!"}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ChatRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerEnhance, &pb.EnhanceRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send EnhanceRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerStars, &pb.StarsRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send StarsRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerSkillsChange, &pb.SkillChangeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send SkillChangeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerSkillsUpgrade, &pb.SkillUpgradeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send SkillUpgradeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerSkillsAdvanced, &pb.SkillAdvancedRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send SkillAdvancedRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerStarsLevel, &pb.StarsLevelRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send StarsLevelRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerCompose, &pb.ComposeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ComposeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGemstone, &pb.GemstoneChangeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GemstoneChangeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGemstoneAll, &pb.GemstoneChangeAllRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GemstoneChangeAllRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerAutoDisassembly, &pb.AutoDisassemblyRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send AutoDisassemblyRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteHomeOpenBox, &pb.HomeOpenBoxRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send HomeOpenBoxRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerBoxResIntoBag, &pb.BoxResIntoBagRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send BoxResIntoBagRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerUpgradeLevel, &pb.UpgradeLevelRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send UpgradeLevelRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerStar, &pb.PartnerStarRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerStarRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerCombat, &pb.PartnerCombatRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerCombatRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerRest, &pb.PartnerRestRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerRestRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerReplace, &pb.PartnerReplaceRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerReplaceRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerList, &pb.PartnerListRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerListRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerSlotInfo, &pb.PartnerSlotInfoRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerSlotInfoRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerEquipBest, &pb.EquipBestRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerEquipBestRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerInstanceInfo, &pb.InstanceInfoRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send InstanceInfoRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerCompose, &pb.PartnerComposeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerComposeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerHomeUpgradeBox, &pb.HomeUpgradeBoxRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send HomeUpgradeBoxRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerHomeUpgradeBoxQuick, &pb.HomeUpgradeBoxQuickRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send HomeUpgradeBoxQuickRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerObservePlayerInfo, &pb.ObservePlayerInfoRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ObservePlayerInfoRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerHomeBoxSetting, &pb.HomeBoxSettingRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send HomeBoxSettingRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerReforge, &pb.PartnerReforgeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerReforgeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPartnerDisassemble, &pb.PartnerDisassembleRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PartnerDisassembleRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerActivityRewardReceive, &pb.ActivityRewardReceiveRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ActivityRewardReceiveRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuidanceUpdate, &pb.GuidanceUpdateRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuidanceUpdateRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteChatHistory, &pb.ChatHistoryRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send ChatHistoryRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteGuardIslandMultiplier, &pb.GuardIslandMultiplierRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuardIslandMultiplierRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPetStar, &pb.PetStarRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PetStarRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPetLevel, &pb.PetLevelRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PetLevelRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerPetMove, &pb.PetMoveRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PetMoveRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildList, &pb.GuildListRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildListRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildCreate, &pb.GuildCreateRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildCreateRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildApply, &pb.GuildApplyRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildApplyRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildCancelApply, &pb.GuildCancelApplyRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildCancelApplyRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildDismiss, &pb.GuildDismissRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildDismissRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildLeave, &pb.GuildLeaveRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildLeaveRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildKick, &pb.GuildKickRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildKickRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildPromote, &pb.GuildPromoteRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildPromoteRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildApplyList, &pb.GuildApplyListRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildApplyListRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildHandleApply, &pb.GuildHandleApplyRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildHandleApplyRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildInfo, &pb.GuildInfoRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildInfoRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildMemberList, &pb.GuildMemberListRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildMemberListRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildDonateInfo, &pb.GuildDonateInfoRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildDonateInfoRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildDonate, &pb.GuildDonateRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildDonateRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildShopBuyInfo, &pb.GuildShopBuyInfoRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildShopBuyInfoRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildShopBuy, &pb.GuildShopBuyRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildShopBuyRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteManagerGuildUpgradeSkill, &pb.GuildUpgradeSkillRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send GuildUpgradeSkillRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteMatchPvpMatching, &pb.PvpMatchingRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PvpMatchingRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteMatchPVPChallenge, &pb.PVPChallengeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PVPChallengeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteMatchPvpRecoder, &pb.PvpRecorderRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PvpRecorderRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteMatchTurnBasedSlot, &pb.TurnBasedSlotRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send TurnBasedSlotRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteMatchTurnBasedChallenge, &pb.TurnBasedChallengeRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send TurnBasedChallengeRequest")
	// 	return
	// }

	// if err := c.SendMsg(protocol.RouteMatchManagerMatching, &pb.MatchingRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send MatchingRequest")
	// 	return
	// }

	// // if err := c.SendMsg(protocol.RouteMatchManagerReady, &pb.MatchingReadyNotify{}); err != nil {
	// // 	log.Error().Err(err).Msg("failed to send ReadyRequest")
	// // 	return
	// // }

	// if err := c.SendMsg(protocol.RouteMatchPVPInfo, &pb.PVPInfoRequest{}); err != nil {
	// 	log.Error().Err(err).Msg("failed to send PVPInfoRequest")
	// 	return
	// }

	if err := c.SendMsg(protocol.RouteRoomManagerFateCoin, &pb.FateCoinRequest{}); err != nil {
		log.Error().Err(err).Msg("failed to send FateCoinRequest")
		return
	}

}
