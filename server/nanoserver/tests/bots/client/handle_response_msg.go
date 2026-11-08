package client

import (
	"nanoserver/protocol/protobuf/pb"
	"nanoserver/tests/bots/internal/message"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func (c *Client) makeResponseHandler(factory func() proto.Message, handlerFunc func(proto.Message)) func([]byte) {
	return func(data []byte) {
		msg := factory()
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Error().Int64("uid", c.uid).Str("name", c.name).Msgf("Failed to unmarshal message: %T", msg)
			return
		}
		handlerFunc(msg)
	}
}

// initResponseHandlers initializes the response handlers for the client.
func (c *Client) initResponseHandlers() {
	c.staticResponseHandlers = map[uint64]func([]byte){
		uint64(pb.CmdId_C2S_Login):           c.makeResponseHandler(func() proto.Message { return &pb.LoginToGameServerResponse{} }, c.LoginToGameServerResponse),
		uint64(pb.CmdId_C2S_Join):            c.makeResponseHandler(func() proto.Message { return &pb.JoinDeskResponse{} }, c.JoinDeskResponse),
		uint64(pb.CmdId_C2S_Exit):            c.makeResponseHandler(func() proto.Message { return &pb.ExitResponse{} }, c.ExitResponse),
		uint64(pb.CmdId_C2S_Move):            c.makeResponseHandler(func() proto.Message { return &pb.MoveResponse{} }, c.MoveResponse),
		uint64(pb.CmdId_C2S_Attack):          c.makeResponseHandler(func() proto.Message { return &pb.AttackResponse{} }, c.AttackResponse),
		uint64(pb.CmdId_C2S_GetBag):          c.makeResponseHandler(func() proto.Message { return &pb.GetBagResponse{} }, c.GetBagResponse),
		uint64(pb.CmdId_C2S_ItemUse):         c.makeResponseHandler(func() proto.Message { return &pb.ItemUseResponse{} }, c.ItemUseResponse),
		uint64(pb.CmdId_C2S_ItemDestroy):     c.makeResponseHandler(func() proto.Message { return &pb.ItemDestroyResponse{} }, c.ItemDestroyResponse),
		uint64(pb.CmdId_C2S_ItemDisassembly): c.makeResponseHandler(func() proto.Message { return &pb.ItemDisassemblyResponse{} }, c.ItemDisassemblyResponse),
		uint64(pb.CmdId_C2S_InstanceTimeAdd): c.makeResponseHandler(func() proto.Message { return &pb.InstanceTimeAddResponse{} }, c.InstanceTimeAddResponse),
		uint64(pb.CmdId_C2S_Revive):          c.makeResponseHandler(func() proto.Message { return &pb.ReviveResponse{} }, c.ReviveResponse),
		uint64(pb.CmdId_C2S_ChangeEquip):     c.makeResponseHandler(func() proto.Message { return &pb.ChangeEquipResponse{} }, c.ChangeEquipResponse),
		uint64(pb.CmdId_C2S_TaskComplate):    c.makeResponseHandler(func() proto.Message { return &pb.TaskComplateResponse{} }, c.TaskCompleteResponse),
		uint64(pb.CmdId_C2S_MatchingRequest): c.makeResponseHandler(func() proto.Message { return &pb.MatchingResponse{} }, c.MatchingResponse),
		// uint64(pb.CmdId_C2S_MatchingReady):         c.makeResponseHandler(func() proto.Message { return &pb.MatchingReadyResponse{} }, c.MatchingReadyResponse),
		uint64(pb.CmdId_C2S_ChangeRoleName):      c.makeResponseHandler(func() proto.Message { return &pb.ChangeRoleNameResponse{} }, c.ChangeRoleNameResponse),
		uint64(pb.CmdId_C2S_Enhance):             c.makeResponseHandler(func() proto.Message { return &pb.EnhanceResponse{} }, c.EnhanceResponse),
		uint64(pb.CmdId_C2S_Stars):               c.makeResponseHandler(func() proto.Message { return &pb.StarsResponse{} }, c.StarsResponse),
		uint64(pb.CmdId_C2S_EquipRefine):         c.makeResponseHandler(func() proto.Message { return &pb.EquipRefineResponse{} }, c.EquipRefineResponse),
		uint64(pb.CmdId_C2S_EquipRefineReplace):  c.makeResponseHandler(func() proto.Message { return &pb.EquipRefineReplaceResponse{} }, c.EquipRefineReplaceResponse),
		uint64(pb.CmdId_C2S_Chat):                c.makeResponseHandler(func() proto.Message { return &pb.ChatResponse{} }, c.ChatResponse),
		uint64(pb.CmdId_C2S_SkillChange):         c.makeResponseHandler(func() proto.Message { return &pb.SkillChangeResponse{} }, c.SkillChangeResponse),
		uint64(pb.CmdId_C2S_SkillUpgrade):        c.makeResponseHandler(func() proto.Message { return &pb.SkillUpgradeResponse{} }, c.SkillUpgradeResponse),
		uint64(pb.CmdId_C2S_SkillAdvanced):       c.makeResponseHandler(func() proto.Message { return &pb.SkillAdvancedResponse{} }, c.SkillAdvancedResponse),
		uint64(pb.CmdId_C2S_EndlesstrialsReward): c.makeResponseHandler(func() proto.Message { return &pb.EndlesstrialsRewardResponse{} }, c.EndlessTrialsRewardResponse),
		uint64(pb.CmdId_C2S_StarsLevel):          c.makeResponseHandler(func() proto.Message { return &pb.StarsLevelResponse{} }, c.StarsLevelResponse),
		uint64(pb.CmdId_C2S_PVPInfo):             c.makeResponseHandler(func() proto.Message { return &pb.PVPInfoResponse{} }, c.PVPInfoResponse),
		uint64(pb.CmdId_C2S_Compose):             c.makeResponseHandler(func() proto.Message { return &pb.ComposeResponse{} }, c.ComposeResponse),
		uint64(pb.CmdId_C2S_CheckBoss):           c.makeResponseHandler(func() proto.Message { return &pb.CheckBossResponse{} }, c.CheckBossResponse),
		uint64(pb.CmdId_C2S_GemstoneChange):      c.makeResponseHandler(func() proto.Message { return &pb.GemstoneChangeResponse{} }, c.GemstoneChangeResponse),
		uint64(pb.CmdId_C2S_DoubleReward):        c.makeResponseHandler(func() proto.Message { return &pb.DoubleRewardResponse{} }, c.DoubleRewardResponse),
		uint64(pb.CmdId_C2S_CopySweep):           c.makeResponseHandler(func() proto.Message { return &pb.CopySweepResponse{} }, c.CopySweepResponse),
		uint64(pb.CmdId_C2S_GemstoneChangeAll):   c.makeResponseHandler(func() proto.Message { return &pb.GemstoneChangeAllResponse{} }, c.GemstoneChangeAllResponse),
		uint64(pb.CmdId_C2S_Mining):              c.makeResponseHandler(func() proto.Message { return &pb.MiningResponse{} }, c.MiningResponse),
		uint64(pb.CmdId_C2S_DailyMapInfo):        c.makeResponseHandler(func() proto.Message { return &pb.DailyMapInfoResponse{} }, c.DailyMapInfoResponse),
		uint64(pb.CmdId_C2S_AutoPlay):            c.makeResponseHandler(func() proto.Message { return &pb.AutoPlayResponse{} }, c.AutoPlayResponse),
		uint64(pb.CmdId_C2S_AutoDisassembly):     c.makeResponseHandler(func() proto.Message { return &pb.AutoDisassemblyResponse{} }, c.AutoDisassemblyResponse),
		uint64(pb.CmdId_C2S_HomeOpenBox):         c.makeResponseHandler(func() proto.Message { return &pb.HomeOpenBoxResponse{} }, c.HomeOpenBoxResponse),
		uint64(pb.CmdId_C2S_UpgradeLevel):        c.makeResponseHandler(func() proto.Message { return &pb.UpgradeLevelResponse{} }, c.UpgradeLevelResponse),
		uint64(pb.CmdId_C2S_BoxResIntoBag):       c.makeResponseHandler(func() proto.Message { return &pb.BoxResIntoBagResponse{} }, c.BoxResIntoBagResponse),
		uint64(pb.CmdId_C2S_PartnerStar):         c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerStarResponse),
		uint64(pb.CmdId_C2S_PartnerCombat):       c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerCombatResponse),
		uint64(pb.CmdId_C2S_PartnerRest):         c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerRestResponse),
		uint64(pb.CmdId_C2S_PartnerList):         c.makeResponseHandler(func() proto.Message { return &pb.PartnerListResponse{} }, c.PartnerListResponse),
		uint64(pb.CmdId_C2S_PartnerSlotInfo):     c.makeResponseHandler(func() proto.Message { return &pb.PartnerSlotInfoResponse{} }, c.PartnerSlotInfoResponse),
		uint64(pb.CmdId_C2S_EquipBest):           c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.EquipBestResponse),
		// uint64(pb.CmdId_C2S_PetCombat):             c.makeResponseHandler(func() proto.Message { return &pb.PetCombatResponse{} }, c.PetCombatResponse),
		// uint64(pb.CmdId_C2S_PetRest):               c.makeResponseHandler(func() proto.Message { return &pb.PetRestResponse{} }, c.PetRestResponse),
		uint64(pb.CmdId_C2S_PetStar):               c.makeResponseHandler(func() proto.Message { return &pb.PetStarResponse{} }, c.PetStarResponse),
		uint64(pb.CmdId_C2S_PetLevel):              c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PetLevelResponse),
		uint64(pb.CmdId_C2S_PetMove):               c.makeResponseHandler(func() proto.Message { return &pb.PetMoveResponse{} }, c.PetMoveResponse),
		uint64(pb.CmdId_C2S_InstanceInfo):          c.makeResponseHandler(func() proto.Message { return &pb.InstanceInfoResponse{} }, c.InstanceInfoResponse),
		uint64(pb.CmdId_C2S_BattleReward):          c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.BattleRewardResponse),
		uint64(pb.CmdId_C2S_PartnerCompose):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerComposeResponse),
		uint64(pb.CmdId_C2S_HomeUpgradeBox):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.HomeUpgradeBoxResponse),
		uint64(pb.CmdId_C2S_HomeUpgradeBoxQuick):   c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.HomeUpgradeBoxQuickResponse),
		uint64(pb.CmdId_C2S_GuildList):             c.makeResponseHandler(func() proto.Message { return &pb.GuildListResponse{} }, c.GuildListResponse),
		uint64(pb.CmdId_C2S_GuildCreate):           c.makeResponseHandler(func() proto.Message { return &pb.GuildCreateResponse{} }, c.GuildCreateResponse),
		uint64(pb.CmdId_C2S_GuildApply):            c.makeResponseHandler(func() proto.Message { return &pb.GuildApplyResponse{} }, c.GuildApplyResponse),
		uint64(pb.CmdId_C2S_GuildCancelApply):      c.makeResponseHandler(func() proto.Message { return &pb.GuildCancelApplyResponse{} }, c.GuildCancelApplyResponse),
		uint64(pb.CmdId_C2S_GuildDismiss):          c.makeResponseHandler(func() proto.Message { return &pb.GuildDismissResponse{} }, c.GuildDismissResponse),
		uint64(pb.CmdId_C2S_GuildLeave):            c.makeResponseHandler(func() proto.Message { return &pb.GuildLeaveResponse{} }, c.GuildLeaveResponse),
		uint64(pb.CmdId_C2S_GuildKick):             c.makeResponseHandler(func() proto.Message { return &pb.GuildKickResponse{} }, c.GuildKickResponse),
		uint64(pb.CmdId_C2S_GuildPromote):          c.makeResponseHandler(func() proto.Message { return &pb.GuildPromoteResponse{} }, c.GuildPromoteResponse),
		uint64(pb.CmdId_C2S_GuildApplyList):        c.makeResponseHandler(func() proto.Message { return &pb.GuildApplyListResponse{} }, c.GuildApplyListResponse),
		uint64(pb.CmdId_C2S_GuildHandleApply):      c.makeResponseHandler(func() proto.Message { return &pb.GuildHandleApplyResponse{} }, c.GuildHandleApplyResponse),
		uint64(pb.CmdId_C2S_GuildInfo):             c.makeResponseHandler(func() proto.Message { return &pb.GuildInfoResponse{} }, c.GuildInfoResponse),
		uint64(pb.CmdId_C2S_GuildMemberList):       c.makeResponseHandler(func() proto.Message { return &pb.GuildMemberListResponse{} }, c.GuildMemberListResponse),
		uint64(pb.CmdId_C2S_GuildDonateInfo):       c.makeResponseHandler(func() proto.Message { return &pb.GuildDonateInfoResponse{} }, c.GuildDonateInfoResponse),
		uint64(pb.CmdId_C2S_GuildDonate):           c.makeResponseHandler(func() proto.Message { return &pb.GuildDonateInfoResponse{} }, c.GuildDonateResponse),
		uint64(pb.CmdId_C2S_GuildShopBuyInfo):      c.makeResponseHandler(func() proto.Message { return &pb.GuildShopBuyInfoResponse{} }, c.GuildShopBuyInfoResponse),
		uint64(pb.CmdId_C2S_GuildShopBuy):          c.makeResponseHandler(func() proto.Message { return &pb.GuildShopBuyInfoResponse{} }, c.GuildShopBuyResponse),
		uint64(pb.CmdId_C2S_GuildUpgradeSkill):     c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.GuildUpgradeSkillResponse),
		uint64(pb.CmdId_C2S_ObservePlayerInfo):     c.makeResponseHandler(func() proto.Message { return &pb.ObservePlayerInfoResponse{} }, c.ObservePlayerInfoResponse),
		uint64(pb.CmdId_C2S_HomeBoxSetting):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.HomeBoxSettingResponse),
		uint64(pb.CmdId_C2S_PartnerReforge):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerReforgeResponse),
		uint64(pb.CmdId_C2S_PartnerDisassemble):    c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerDisassembleResponse),
		uint64(pb.CmdId_C2S_PvpMatchingRequest):    c.makeResponseHandler(func() proto.Message { return &pb.PvpMatchingResponse{} }, c.PvpMatchingResponse),
		uint64(pb.CmdId_C2S_PVPChallengeRequest):   c.makeResponseHandler(func() proto.Message { return &pb.PVPChallengeResponse{} }, c.PVPChallengeResponse),
		uint64(pb.CmdId_C2S_PartnerReplace):        c.makeResponseHandler(func() proto.Message { return &pb.CommonResponse{} }, c.PartnerReplaceResponse),
		uint64(pb.CmdId_C2S_PvpRecoder):            c.makeResponseHandler(func() proto.Message { return &pb.PvpRecorderResponse{} }, c.PvpRecorderResponse),
		// uint64(pb.CmdId_C2S_ActivityRewardReceive): c.makeResponseHandler(func() proto.Message { return &pb.ActivityRewardReceiveResponse{} }, c.ActivityRewardReceiveResponse),
		uint64(pb.CmdId_C2S_GuidanceUpdate):        c.makeResponseHandler(func() proto.Message { return &pb.GuidanceUpdateResponse{} }, c.GuidanceUpdateResponse),
		uint64(pb.CmdId_C2S_TurnBasedSlot):         c.makeResponseHandler(func() proto.Message { return &pb.TurnBasedSlotResponse{} }, c.TurnBasedSlotResponse),
		uint64(pb.CmdId_C2S_TurnBasedChallenge):    c.makeResponseHandler(func() proto.Message { return &pb.TurnBasedChallengeResponse{} }, c.TurnBasedChallengeResponse),
		uint64(pb.CmdId_C2S_ChatHistory):           c.makeResponseHandler(func() proto.Message { return &pb.ChatHistoryResponse{} }, c.ChatHistoryResponse),
		uint64(pb.CmdId_C2S_GuardIslandMultiplier): c.makeResponseHandler(func() proto.Message { return &pb.GuardIslandMultiplierResponse{} }, c.GuardIslandMultiplierResponse),
		uint64(pb.CmdId_C2S_ElementInfo):           c.makeResponseHandler(func() proto.Message { return &pb.ElementInfoResponse{} }, c.ElementInfoResponse),
	}
}

func (c *Client) HandleResponseMsg(msg *message.Message) {
	responseHandler, ok := <-c.dynamicResponseChan
	if !ok {
		log.Warn().Int64("uid", c.uid).Str("name", c.name).Msgf("No handler for dynamic response message ID: %d", msg.ID)
		return
	}

	if responseHandler.DynamicID != msg.ID {
		log.Warn().Int64("uid", c.uid).Str("name", c.name).Msgf("Dynamic response ID mismatch: expected=%d, got=%d, staticID=%d", responseHandler.DynamicID, msg.ID, responseHandler.StaticID)
		return
	}

	responseHandler.Handler(msg.Data)
}

func (c *Client) LoginToGameServerResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("LoginToGameServerResponse: %+v", resp)
}

func (c *Client) JoinDeskResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("JoinDeskResponse: %+v", resp)
}

func (c *Client) ChatResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ChatResponse: %+v", resp)
}

func (c *Client) ExitResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ExitResponse: %+v", resp)
}

func (c *Client) MoveResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("MoveResponse: %+v", resp)
}

func (c *Client) AttackResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("AttackResponse: %+v", resp)
}

func (c *Client) GetBagResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GetBagResponse: %+v", resp)
}

func (c *Client) ItemUseResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ItemUseResponse: %+v", resp)
}

func (c *Client) ItemDestroyResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ItemDestroyResponse: %+v", resp)
}

func (c *Client) ItemDisassemblyResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ItemDisassemblyResponse: %+v", resp)
}

func (c *Client) InstanceTimeAddResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("InstanceTimeAddResponse: %+v", resp)
}

func (c *Client) ReviveResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ReviveResponse: %+v", resp)
}

func (c *Client) ChangeEquipResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ChangeEquipResponse: %+v", resp)
}

func (c *Client) TaskCompleteResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("TaskCompleteResponse: %+v", resp)
}

func (c *Client) MatchingResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("MatchingResponse: %+v", resp)
}

func (c *Client) ChangeRoleNameResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ChangeRoleNameResponse: %+v", resp)
}

func (c *Client) EnhanceResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("EnhanceResponse: %+v", resp)
}

func (c *Client) StarsResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("StarsResponse: %+v", resp)
}

func (c *Client) EquipRefineResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("EquipRefineResponse: %+v", resp)
}

func (c *Client) EquipRefineReplaceResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("EquipRefineReplaceResponse: %+v", resp)
}

func (c *Client) SkillChangeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("SkillChangeResponse: %+v", resp)
}

func (c *Client) SkillUpgradeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("SkillUpgradeResponse: %+v", resp)
}

func (c *Client) SkillAdvancedResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("SkillAdvancedResponse: %+v", resp)
}

func (c *Client) EndlessTrialsRewardResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("EndlessTrialsRewardResponse: %+v", resp)
}

func (c *Client) EndlesstrialsRewardResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("EndlesstrialsRewardResponse: %+v", resp)
}

func (c *Client) StarsLevelResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("StarsLevelResponse: %+v", resp)
}

func (c *Client) PVPInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PVPInfoResponse: %+v", resp)
}

func (c *Client) ComposeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ComposeResponse: %+v", resp)
}

func (c *Client) CheckBossResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("CheckBossResponse: %+v", resp)
}

func (c *Client) GemstoneChangeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GemstoneChangeResponse: %+v", resp)
}

func (c *Client) DoubleRewardResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("DoubleRewardResponse: %+v", resp)
}

func (c *Client) CopySweepResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("CopySweepResponse: %+v", resp)
}

func (c *Client) GemstoneChangeAllResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GemstoneChangeAllResponse: %+v", resp)
}

func (c *Client) MiningResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("MiningResponse: %+v", resp)
}

func (c *Client) DailyMapInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("DailyMapInfoResponse: %+v", resp)
}

func (c *Client) AutoPlayResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("AutoPlayResponse: %+v", resp)
}

func (c *Client) AutoDisassemblyResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("AutoDisassemblyResponse: %+v", resp)
}

func (c *Client) HomeOpenBoxResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("HomeOpenBoxResponse: %+v", resp)
}

func (c *Client) UpgradeLevelResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("UpgradeLevelResponse: %+v", resp)
}

func (c *Client) BoxResIntoBagResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("BoxResIntoBagResponse: %+v", resp)
}

func (c *Client) PartnerListResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerListResponse: %+v", resp)
}

func (c *Client) PartnerSlotInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerSlotInfoResponse: %+v", resp)
}

func (c *Client) PetStarResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PetStarResponse: %+v", resp)
}

func (c *Client) PetMoveResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PetMoveResponse: %+v", resp)
}

func (c *Client) InstanceInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("InstanceInfoResponse: %+v", resp)
}

func (c *Client) BattleRewardResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("BattleRewardResponse: %+v", resp)
}

func (c *Client) GuildListResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildListResponse: %+v", resp)
}

func (c *Client) GuildCreateResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildCreateResponse: %+v", resp)
}

func (c *Client) GuildApplyResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildApplyResponse: %+v", resp)
}

func (c *Client) GuildCancelApplyResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildCancelApplyResponse: %+v", resp)
}

func (c *Client) GuildDismissResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildDismissResponse: %+v", resp)
}

func (c *Client) GuildLeaveResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildLeaveResponse: %+v", resp)
}

func (c *Client) GuildKickResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildKickResponse: %+v", resp)
}

func (c *Client) GuildPromoteResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildPromoteResponse: %+v", resp)
}

func (c *Client) GuildApplyListResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildApplyListResponse: %+v", resp)
}

func (c *Client) GuildHandleApplyResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildHandleApplyResponse: %+v", resp)
}

func (c *Client) GuildInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildInfoResponse: %+v", resp)
}

func (c *Client) GuildMemberListResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildMemberListResponse: %+v", resp)
}

func (c *Client) GuildDonateInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildDonateInfoResponse: %+v", resp)
}

func (c *Client) GuildDonateResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildDonateResponse: %+v", resp)
}

func (c *Client) GuildShopBuyInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildShopBuyInfoResponse: %+v", resp)
}

func (c *Client) GuildUpgradeSkillResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildUpgradeSkillResponse: %+v", resp)
}

func (c *Client) GuildShopBuyResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuildShopBuyResponse: %+v", resp)
}

func (c *Client) ObservePlayerInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ObservePlayerInfoResponse: %+v", resp)
}

func (c *Client) PvpMatchingResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PvpMatchingResponse: %+v", resp)
}

func (c *Client) PVPChallengeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PVPChallengeResponse: %+v", resp)
}

func (c *Client) PvpRecorderResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PvpRecorderResponse: %+v", resp)
}

func (c *Client) ActivityRewardReceiveResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ActivityRewardReceiveResponse: %+v", resp)
}

func (c *Client) GuidanceUpdateResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuidanceUpdateResponse: %+v", resp)
}

func (c *Client) TurnBasedSlotResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("TurnBasedSlotResponse: %+v", resp)
}

func (c *Client) TurnBasedChallengeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("TurnBasedChallengeResponse: %+v", resp)
}

func (c *Client) ChatHistoryResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ChatHistoryResponse: %+v", resp)
}

func (c *Client) GuardIslandMultiplierResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("GuardIslandMultiplierResponse: %+v", resp)
}

func (c *Client) ElementInfoResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("ElementInfoResponse: %+v", resp)
}

func (c *Client) HomeBoxSettingResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("HomeBoxSettingResponse: %+v", resp)
}

func (c *Client) PartnerReforgeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerReforgeResponse: %+v", resp)
}

func (c *Client) PartnerDisassembleResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerDisassembleResponse: %+v", resp)
}

func (c *Client) PartnerReplaceResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerReplaceResponse: %+v", resp)
}

func (c *Client) EquipBestResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("EquipBestResponse: %+v", resp)
}

func (c *Client) PetCombatResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PetCombatResponse: %+v", resp)
}

func (c *Client) PetRestResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PetRestResponse: %+v", resp)
}

func (c *Client) PetLevelResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PetLevelResponse: %+v", resp)
}

func (c *Client) PartnerStarResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerStarResponse: %+v", resp)
}

func (c *Client) PartnerCombatResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerCombatResponse: %+v", resp)
}

func (c *Client) PartnerRestResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerRestResponse: %+v", resp)
}

func (c *Client) PartnerComposeResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("PartnerComposeResponse: %+v", resp)
}

func (c *Client) HomeUpgradeBoxResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("HomeUpgradeBoxResponse: %+v", resp)
}

func (c *Client) HomeUpgradeBoxQuickResponse(resp proto.Message) {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("HomeUpgradeBoxQuickResponse: %+v", resp)
}
