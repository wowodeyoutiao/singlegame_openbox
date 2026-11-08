package game

type LotteryInfo struct {
	FreeNum       int32   //免费次数
	RefreshTime   int64   //刷新时间
	PrizePool     []int32 //奖励池水
	HavePrizePool []int32 //已有槽位
	Progress      int32   //进度
	NeedItemId    int32   //需要道具id
}

// func (p *Player) lotteryInfolotteryToPb() *pb.lotteryEmpowerList {
// 	empower := &pb.lotteryEmpowerList{}
// 	for _, v := range p.LotteryInfo.lottery {
// 		empower.lotterys = append(empower.lotterys, &pb.lotteryEmpower{
// 			SlotId:     v.SlotId,
// 			Level:      v.Level,
// 			BreakLevel: v.Break,
// 		})
// 	}
// 	return empower
// }
// func (p *Player) lotteryInfolotteryInworkToPb() *pb.lotteryUids {
// 	uid := &pb.lotteryUids{}
// 	for _, v := range p.lotteryInfo.lotteryInwork {
// 		uid.Uids = append(uid.Uids, v)
// 	}
// 	return uid
// }
