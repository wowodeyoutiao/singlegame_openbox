package game

// import (
// 	"encoding/json"
// 	"nanoserver/db"
// 	"nanoserver/db/model"
// 	"nanoserver/global"
// 	"time"
// )

// type (
// 	ItemRecoder struct {
// 		UserID    int64
// 		ItemCfgId int32
// 		ItemNum   int32
// 		Time      time.Time
// 		MapId     int32
// 		Reason    string
// 	}

// 	RefreshTime struct {
// 		Uid        int32
// 		ItemCfgIds []int32   `json:"cfgids"`
// 		NextTime   time.Time `json:"-"`
// 		Tick       int32     `json:"-"`
// 	}

// 	ItemRecoderMgr struct {
// 		refresh []RefreshTime

// 		server map[int32]int32

// 		personnal map[int64]map[int32]int32
// 	}
// )

// var itemRecoderMgr = ItemRecoderManager()

// func ItemRecoderManager() *ItemRecoderMgr {
// 	mgr := &ItemRecoderMgr{}

// 	return mgr
// }

// func (irm *ItemRecoderMgr) load() {

// 	irm.refresh = []RefreshTime{}
// 	irm.server = make(map[int32]int32)
// 	irm.personnal = make(map[int64]map[int32]int32)

// 	for _, v := range global.ItemDropControlCfg {
// 		r := &model.ItemControlRefresh{Uid: v.Uid}
// 		ids, err := json.Marshal(v.ItemIDs)
// 		if err == nil {
// 			r.ItemCfgIds = string(ids)
// 		}
// 		r.RecoderTime = time.Now().Add(time.Duration(v.ControlTime) * time.Second).Unix()

// 		if data, err := db.QueryRecoder(v.Uid); err != nil {
// 			go db.InsertRecoder(r)
// 			rt := RefreshTime{Uid: v.Uid, ItemCfgIds: v.ItemIDs, NextTime: time.Now().Add(time.Duration(v.ControlTime) * time.Second), Tick: v.ControlTime}
// 			irm.refresh = append(irm.refresh, rt)
// 		} else {
// 			rt := RefreshTime{Uid: v.Uid, Tick: v.ControlTime}

// 			if err = json.Unmarshal([]byte(data.ItemCfgIds), &rt.ItemCfgIds); data.ItemCfgIds != "" && err != nil {
// 				logger.Panicf("掉落控制加载失败 err: %s", err.Error())
// 				return
// 			}

// 			rt.NextTime = time.Unix(data.RecoderTime, 0)

// 			irm.refresh = append(irm.refresh, rt)
// 		}
// 	}

// 	logger.Printf("掉落控制加载完成  %v", irm.refresh)

// 	if r, err := db.QueryAllControlLast(); err == nil {
// 		for _, v := range r {
// 			if v.IsServer == 0 {
// 				irm.server[v.ItemCfgId] = v.ItemNum
// 			} else {
// 				if _, ok := irm.personnal[v.Uid]; ok {
// 					irm.personnal[v.Uid][v.ItemCfgId] = v.ItemNum
// 				} else {
// 					irm.personnal[v.Uid] = make(map[int32]int32)
// 					irm.personnal[v.Uid][v.ItemCfgId] = v.ItemNum
// 				}
// 			}
// 		}
// 	}
// 	logger.Printf(" 个人掉落记录加载完成  %v", irm.personnal)
// 	logger.Printf(" 服务器掉落记录加载完成  %v", irm.server)
// 	//irm.chRecoder = make(chan *ItemRecoder, 100)
// }

// func (irm *ItemRecoderMgr) control(cfgId int32) *global.ItemDropControlConfig {
// 	for _, v := range global.ItemDropControlCfg {
// 		for _, v1 := range v.ItemIDs {

// 			if v1 == cfgId {
// 				return v
// 			}
// 		}
// 	}
// 	return nil
// }

// func (irm *ItemRecoderMgr) update() {
// 	for k, v := range irm.refresh {
// 		if time.Now().After(v.NextTime) {
// 			irm.refresh[k].NextTime = time.Now().Add(time.Duration(v.Tick) * time.Second)

// 			for _, id := range v.ItemCfgIds {
// 				// 清空servercontrol
// 				if _, ok := irm.server[int32(id)]; ok {
// 					if irm.server[int32(id)] > 0 {
// 						irm.server[int32(id)] = 0
// 						go db.UpdateControlLast(&model.ItemControlLast{ItemCfgId: int32(id), ItemNum: 0})

// 						logger.Printf(" 重置服务器掉落记录 ItemId: %d ", id)
// 					}

// 				}

// 				//清空personctrol
// 				for k, r1 := range irm.personnal {
// 					if _, ok := r1[int32(id)]; ok {
// 						if r1[int32(id)] > 0 {
// 							r1[int32(id)] = 0
// 							go db.UpdateControlLast(&model.ItemControlLast{Uid: k, ItemCfgId: int32(id), ItemNum: 0})

// 							logger.Printf(" 重置个人掉落记录 ItemId: %d ", id)
// 						}
// 					}
// 				}
// 			}

// 			r := &model.ItemControlRefresh{Uid: v.Uid, RecoderTime: irm.refresh[k].NextTime.Unix()}
// 			go db.UpdateRecoder(r)
// 		}
// 	}
// }

// func (irm *ItemRecoderMgr) add(uid int64, itemid int32, itemnum int32, mapid int32) {

// 	if cfg := irm.control(itemid); cfg != nil {

// 		if cfg.Scope == 0 {
// 			irm.server[itemid] += itemnum

// 			rl := model.ItemControlLast{ItemCfgId: itemid, IsServer: cfg.Scope}
// 			if data, err := db.QueryControlLast(&rl); err == nil {
// 				data.ItemNum = irm.server[itemid]
// 				go db.UpdateControlLast(data)
// 			} else {
// 				rl.ItemNum = irm.server[itemid]
// 				go db.InsertControlLast(&rl)
// 			}
// 		} else if cfg.Scope == 1 {
// 			if _, ok := irm.personnal[uid]; ok {
// 				irm.personnal[uid][itemid] += itemnum
// 			} else {
// 				irm.personnal[uid] = make(map[int32]int32)
// 				irm.personnal[uid][itemid] = itemnum
// 			}

// 			rl := model.ItemControlLast{ItemCfgId: itemid, IsServer: cfg.Scope, Uid: uid}
// 			if data, err := db.QueryControlLast(&rl); err == nil {
// 				data.ItemNum = irm.personnal[uid][itemid]
// 				go db.UpdateControlLast(data)
// 			} else {
// 				rl.ItemNum = irm.personnal[uid][itemid]
// 				go db.InsertControlLast(&rl)
// 			}
// 		}
// 	}
// }

// func (irm *ItemRecoderMgr) check(uid int64, level int32, itemid int32, itemnum int32, mapid int32) bool {
// 	if cfg := irm.control(itemid); cfg != nil {

// 		var validMap bool
// 		for _, v := range cfg.CandropMaps {
// 			if v == mapid {
// 				validMap = true
// 			}
// 		}
// 		if !validMap {
// 			return false
// 		}

// 		if level > cfg.ControlLevel.Max || level < cfg.ControlLevel.Min {
// 			return false
// 		}

// 		if r, ok := irm.server[itemid]; ok {
// 			if r >= cfg.Maxnum {
// 				return false
// 			}
// 		}

// 		if pr, ok := irm.personnal[uid]; ok {
// 			if r, ok2 := pr[itemid]; ok2 {
// 				if r >= cfg.Maxnum {
// 					return false
// 				}
// 			}
// 		}

// 		irm.add(uid, itemid, itemnum, mapid)
// 	}
// 	return true
// }
