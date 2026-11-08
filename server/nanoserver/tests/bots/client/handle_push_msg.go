package client

import (
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"nanoserver/tests/bots/internal/message"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (c *Client) makePushHandler(factory func() proto.Message, handlerFunc func(proto.Message)) func([]byte) {
	return func(data []byte) {
		msg := factory()
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Error().
				Err(err).
				Str("client", c.name).
				Int64("uid", c.uid).
				Msgf("Failed to unmarshal message: %T", msg)
			return
		}
		handlerFunc(msg)
	}
}

func (c *Client) initPushHandlers() {
	c.pushHandlers = map[string]func([]byte){
		protocol.RoutePlayerEnterDeskPush:      c.makePushHandler(func() proto.Message { return &pb.PlayerEnterDeskPush{} }, c.onPlayerEnterDesk),
		protocol.RouteExitPush:                 c.makePushHandler(func() proto.Message { return &pb.ExitPush{} }, c.onExitRoompush),
		protocol.RouteMovePush:                 c.makePushHandler(func() proto.Message { return &pb.MovePush{} }, c.onMove),
		protocol.RouteAttackPush:               c.makePushHandler(func() proto.Message { return &pb.AttackPush{} }, c.onAttack),
		protocol.RouteMonsterInPush:            c.makePushHandler(func() proto.Message { return &pb.MonsterInPush{} }, c.onMonsterIn),
		protocol.RouteMonsterOutPush:           c.makePushHandler(func() proto.Message { return &pb.MonsterOutPush{} }, c.onMonsterOut),
		protocol.RouteObjectHpChangePush:       c.makePushHandler(func() proto.Message { return &pb.ObjectHpChangePush{} }, c.onMonsterHpChange),
		protocol.RouteObjectDiePush:            c.makePushHandler(func() proto.Message { return &pb.ObjectDiePush{} }, c.onObjectDiePush),
		protocol.RouteAddSourcePush:            c.makePushHandler(func() proto.Message { return &pb.AddSourcePush{} }, c.onAddSource),
		protocol.RouteTakeSourcePush:           c.makePushHandler(func() proto.Message { return &pb.TakeSourcePush{} }, c.onTakeSource),
		protocol.RouteSyncUserBasePush:         c.makePushHandler(func() proto.Message { return &pb.SyncUserBasePush{} }, c.onSyncUserBase),
		protocol.RouteReJoinDeskPush:           c.makePushHandler(func() proto.Message { return &pb.ReJoinDeskPush{} }, c.onReJoinDesk),
		protocol.RouteSyncUserAdvanceInfoPush:  c.makePushHandler(func() proto.Message { return &pb.SyncUserAdvanceInfoPush{} }, c.onSyncUserAdvanceInfo),
		protocol.RouteSkillCooldownPush:        c.makePushHandler(func() proto.Message { return &pb.SkillCooldownPush{} }, c.onSkillCooldown),
		protocol.RouteAddBuffPush:              c.makePushHandler(func() proto.Message { return &pb.AddBuffPush{} }, c.onAddBuff),
		protocol.RouteRemoveBuffPush:           c.makePushHandler(func() proto.Message { return &pb.RemoveBuffPush{} }, c.onRemoveBuff),
		protocol.RouteMatchingResultPush:       c.makePushHandler(func() proto.Message { return &pb.MatchingResultPush{} }, c.onMatchingResult),
		protocol.RouteMatchingStartPush:        c.makePushHandler(func() proto.Message { return &pb.MatchingStartPush{} }, c.onMatchingStart),
		protocol.RouteMatchRoomJoinPush:        c.makePushHandler(func() proto.Message { return &pb.MatchRoomJoinPush{} }, c.onMatchRoomJoin),
		protocol.RouteGameResultMatchPush:      c.makePushHandler(func() proto.Message { return &pb.GameResultMatchPush{} }, c.onMatchFinish),
		protocol.RouteTaskUpdatePush:           c.makePushHandler(func() proto.Message { return &pb.TaskUpdatePush{} }, c.onTaskUpdatePush),
		protocol.RouteChatPush:                 c.makePushHandler(func() proto.Message { return &pb.ChatPush{} }, c.onChatPush),
		protocol.RouteNewSkill:                 c.makePushHandler(func() proto.Message { return &pb.NewSkillPush{} }, c.onNewSkill),
		protocol.RouteItemDropPush:             c.makePushHandler(func() proto.Message { return &pb.ItemDropPush{} }, c.onItemDropPush),
		protocol.RouteGameResultMatchMeleePush: c.makePushHandler(func() proto.Message { return &pb.GameResultMatchMeleePush{} }, c.onGameResultMatchMeleePush),
		protocol.RouteGameObjectAddPush:        c.makePushHandler(func() proto.Message { return &pb.GameObjectAddPush{} }, c.onGameObjectAddPush),
		protocol.RouteGameObjectDelPush:        c.makePushHandler(func() proto.Message { return &pb.GameObjectDelPush{} }, c.onGameObjectDelPush),
		protocol.RouteTeleportPush:             c.makePushHandler(func() proto.Message { return &pb.TeleportPush{} }, c.onTeleportPush),
		protocol.RoutePartnerSlotUpdatePush:    c.makePushHandler(func() proto.Message { return &pb.PartnerSlotUpdatePush{} }, c.onPartnerSlotUpdatePush),
		protocol.SyncUserSkillSlotPush:         c.makePushHandler(func() proto.Message { return &pb.SyncUserSkillSlotPush{} }, c.onSyncUserSkillSlotPush),
		protocol.RouteSyncPartnerListPush:      c.makePushHandler(func() proto.Message { return &pb.SyncPartnerListPush{} }, c.onSyncUserSkillPush),
		protocol.RouteSyncUserStatsPush:        c.makePushHandler(func() proto.Message { return &pb.SyncUserStatsPush{} }, c.onSyncUserStatsPush),
		protocol.RouteSyncUserPetInfoPush:      c.makePushHandler(func() proto.Message { return &pb.SyncUserPetInfoPush{} }, c.onSyncUserPetInfoPush),
		protocol.RouteSyncUserBoxInfoPush:      c.makePushHandler(func() proto.Message { return &pb.SyncUserBoxInfoPush{} }, c.onSyncUserBoxInfoPush),
		protocol.RouteHomeAutoEncounterPush:    c.makePushHandler(func() proto.Message { return &pb.HomeAutoEncounterPush{} }, c.onHomeAutoEncounterPush),
		protocol.RouteBuffListPush:             c.makePushHandler(func() proto.Message { return &pb.BuffListPush{} }, c.onBuffListPush),
		protocol.RouteLevelUpPush:              c.makePushHandler(func() proto.Message { return &pb.LevelUpPush{} }, c.onLevelUpPush),
		protocol.RouteTaskAcceptPush:           c.makePushHandler(func() proto.Message { return &pb.TaskAcceptPush{} }, c.onTaskAcceptPush),
		protocol.RouteSystemNotificationPush:   c.makePushHandler(func() proto.Message { return &pb.SystemNotificationPush{} }, c.onSystemNotificationPush),
	}
}

func (c *Client) HandlePushMsg(msg *message.Message) {
	if handler, exists := c.pushHandlers[msg.Route]; exists {
		handler(msg.Data)
	} else {
		log.Warn().Str("client", c.name).Int64("uid", c.uid).Str("route", msg.Route).Msg("Unhandled push route")
	}
}

func (c *Client) onPlayerEnterDesk(msg proto.Message) {
	push := msg.(*pb.PlayerEnterDeskPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RoutePlayerEnterDeskPush).Msgf("Player entered desk | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onExitRoompush(msg proto.Message) {
	push := msg.(*pb.ExitPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteExitPush).Msgf("Exit room | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMove(msg proto.Message) {
	push := msg.(*pb.MovePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteMovePush).Msgf("Player moved | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onAttack(msg proto.Message) {
	push := msg.(*pb.AttackPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteAttackPush).Msgf("Player attacked | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMonsterIn(msg proto.Message) {
	push := msg.(*pb.MonsterInPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteMonsterInPush).Msgf("Monster entered | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMonsterOut(msg proto.Message) {
	push := msg.(*pb.MonsterOutPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteMonsterOutPush).Msgf("Monster exited | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMonsterHpChange(msg proto.Message) {
	push := msg.(*pb.ObjectHpChangePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteObjectHpChangePush).Msgf("Monster HP changed | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onObjectDiePush(msg proto.Message) {
	push := msg.(*pb.ObjectDiePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteObjectDiePush).Msgf("Object died | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onAddSource(msg proto.Message) {
	push := msg.(*pb.AddSourcePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteAddSourcePush).Msgf("Source added | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onTakeSource(msg proto.Message) {
	push := msg.(*pb.TakeSourcePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteTakeSourcePush).Msgf("Source taken | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSyncUserBase(msg proto.Message) {
	push := msg.(*pb.SyncUserBasePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteSyncUserBasePush).Msgf("User base synced | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onReJoinDesk(msg proto.Message) {
	push := msg.(*pb.ReJoinDeskPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteReJoinDeskPush).Msgf("Rejoined desk | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSyncUserAdvanceInfo(msg proto.Message) {
	push := msg.(*pb.SyncUserAdvanceInfoPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteSyncUserAdvanceInfoPush).Msgf("User advance info synced | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSkillCooldown(msg proto.Message) {
	push := msg.(*pb.SkillCooldownPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteSkillCooldownPush).Msgf("Skill cooldown | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onAddBuff(msg proto.Message) {
	push := msg.(*pb.AddBuffPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteAddBuffPush).Msgf("Buff added | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onRemoveBuff(msg proto.Message) {
	push := msg.(*pb.RemoveBuffPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteRemoveBuffPush).Msgf("Buff removed | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMatchingResult(msg proto.Message) {
	push := msg.(*pb.MatchingResultPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteMatchingResultPush).Msgf("Matching result | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMatchingStart(msg proto.Message) {
	push := msg.(*pb.MatchingStartPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteMatchingStartPush).Msgf("Matching started | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMatchRoomJoin(msg proto.Message) {
	push := msg.(*pb.MatchRoomJoinPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteMatchRoomJoinPush).Msgf("Joined match room | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onMatchFinish(msg proto.Message) {
	push := msg.(*pb.GameResultMatchPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteGameResultMatchPush).Msgf("Match finished | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onTaskUpdatePush(msg proto.Message) {
	push := msg.(*pb.TaskUpdatePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteTaskUpdatePush).Msgf("Task updated | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onChatPush(msg proto.Message) {
	push := msg.(*pb.ChatPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteChatPush).Msgf("Chat message | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onNewSkill(msg proto.Message) {
	push := msg.(*pb.NewSkillPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteNewSkill).Msgf("New skill acquired | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onItemDropPush(msg proto.Message) {
	push := msg.(*pb.ItemDropPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteItemDropPush).Msgf("Item dropped | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onGameResultMatchMeleePush(msg proto.Message) {
	push := msg.(*pb.GameResultMatchMeleePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteGameResultMatchMeleePush).Msgf("Melee match result | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onGameObjectAddPush(msg proto.Message) {
	push := msg.(*pb.GameObjectAddPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteGameObjectAddPush).Msgf("Game object added | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onGameObjectDelPush(msg proto.Message) {
	push := msg.(*pb.GameObjectDelPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteGameObjectDelPush).Msgf("Game object deleted | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onTeleportPush(msg proto.Message) {
	push := msg.(*pb.TeleportPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteTeleportPush).Msgf("Player teleported | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onPartnerSlotUpdatePush(msg proto.Message) {
	push := msg.(*pb.PartnerSlotUpdatePush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RoutePartnerSlotUpdatePush).Msgf("Partner slot updated | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSyncUserSkillSlotPush(msg proto.Message) {
	push := msg.(*pb.SyncUserSkillSlotPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.SyncUserSkillSlotPush).Msgf("User skill slot synced | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSyncUserSkillPush(msg proto.Message) {
	push := msg.(*pb.SyncPartnerListPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.SyncUserSkillSlotPush).Msgf("User skill synced | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSyncUserStatsPush(msg proto.Message) {
	push := msg.(*pb.SyncUserStatsPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteSyncUserStatsPush).Msgf("User stats synced | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSyncUserPetInfoPush(msg proto.Message) {
	push := msg.(*pb.SyncUserPetInfoPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteSyncUserPetInfoPush).Msgf("User pet info synced | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSyncUserBoxInfoPush(msg proto.Message) {
	push := msg.(*pb.SyncUserBoxInfoPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteSyncUserBoxInfoPush).Msgf("User box info synced | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onHomeAutoEncounterPush(msg proto.Message) {
	push := msg.(*pb.HomeAutoEncounterPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteHomeAutoEncounterPush).Msgf("Home auto encounter | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onBuffListPush(msg proto.Message) {
	push := msg.(*pb.BuffListPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteBuffListPush).Msgf("Buff list | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onLevelUpPush(msg proto.Message) {
	push := msg.(*pb.LevelUpPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteLevelUpPush).Msgf("Level up | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onTaskAcceptPush(msg proto.Message) {
	push := msg.(*pb.TaskAcceptPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteTaskAcceptPush).Msgf("Task accepted | Data: %s", marshalOpts.Format(push))
}

func (c *Client) onSystemNotificationPush(msg proto.Message) {
	push := msg.(*pb.SystemNotificationPush)
	marshalOpts := protojson.MarshalOptions{Multiline: false}
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Str("route", protocol.RouteSystemNotificationPush).Msgf("System notification | Data: %s", marshalOpts.Format(push))
}
