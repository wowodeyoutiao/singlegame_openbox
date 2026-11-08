package bootstrap

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"nanoserver/global"
	"nanoserver/pkg/astar"
	"nanoserver/pkg/excel"
	"nanoserver/pkg/sensitive"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("component", "global")

func Initialize() {
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			global.BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			global.BasePath = curPath
		}
	} else {
		log.Fatal("初始化项目根目录失败")
	}

	src := rand.NewSource(time.Now().UnixNano())
	global.Rng = rand.New(src)
	global.LoginOpen = true

	sensitive.Load()
	// 加载配置
	LoadGlobalCfgXML()

	loadRoomCfg()
	loadMonsterCfg()
	loadMapCfg()
	LoadItemCfg()
	LoadRoleCfg()
	LoadSkillCfg()
	LoadSignCfg()
	LoadStoreCfg()
	LoadItemUseCfg()
	LoadActivityCfg()
	LoadBuffCfg()
	LoadTaskCfg()
	LoadPrivilegeCfg()
	LoadRankingAwardCfg()
	LoadRandomCreatCfg()
	LoadRandomDistribCfg()
	LoadVipCfg()
	LoadRobotCfg()
	LoadEnhanceCfg()
	LoadStarsCfg()
	LoadEndlesstrialsCfg()
	LoadEquipMasterCfg()
	LoadPrizeWheelCfg()
	LoadTreasureBoxCfg()
	LoadRandNameCfg()
	LoadFateCfg()
	LoadComposeCfg()
	LoadItemGroupCfg()
	LoadLevelGiftCfg()
	LoadPartnerCfg()
	LoadPartnerTalentCfg()
	LoadMaintriggerCfg()
	LoadOpenBoxCfg()
	LoadBoxLevelCfg()
	LoadSummonCfg()
	LoadPetCfg()
	LoadBattleCfg()
	LoadGuildCfg()
	LoadGuildStoreCfg()
	LoadGuildSkillCfg()
	LoadLampCfg()
	ActionsetCfg()
	LoadBonusCfg()
	LoadBondCfg()
	LoadMainBattleCfg()
	// LoadActivityControllerCfg()
	// LoadActivityQuestCfg()
	LoadElementCfg()
	LoadElementTrialsCfg()
	LoadTitleCfg()
	LoadTitleActivationCfg()
	LoadActiveCodeCfg()
	LoadKillMonControlCfg()
	LoadLotteryCfg()
	LoadAvatarCfg()
	LoadRandomDistribV2()
	LoadRandomCreatV2()
	LoadFateCoinCfg()
	LoadArtifactStrengthCfg()
	LoadArtifactAdvanceCfg()
	LoadPushingWallMonsterCfg()
	LoadOfficialCfg()
}

func LoadGlobalCfgXML() {
	// 加载通用配置
	data, err := os.ReadFile(global.BasePath + "/configs/gameconfig/config.xml")
	if err != nil {
		log.Fatalf("读取文件%v失败: %v\n", "config.xml", err)
	}

	if err = xml.Unmarshal(data, &global.GlobalCfg); err != nil {
		log.Fatal("XML数据解析失败", err)
	}

	global.SwitchCfg = make(map[int32]*global.SwitchConfig)
	for _, v := range global.GlobalCfg.Switch.Module {
		global.SwitchCfg[v.ID] = &global.SwitchConfig{
			OpenLevel: v.Openlevel,
			OpenDay:   v.OpenDay,
			MinBattel: v.MinBattel,
			ServerDay: v.Severday,
		}
	}

	global.GuardIslandCfg = make(map[int64]*global.GuardIslandRecord, 5)
	for _, v := range global.GlobalCfg.GuardIsland.Hangup {
		global.GuardIslandCfg[int64(v.ID)] = &global.GuardIslandRecord{
			ID:         v.ID,
			Type:       v.Type,
			Multiplier: v.Multiplier,
		}
	}

	sort.Slice(global.GlobalCfg.GVG.GVGScore.Items, func(i, j int) bool {
		return global.GlobalCfg.GVG.GVGScore.Items[i].Id > global.GlobalCfg.GVG.GVGScore.Items[j].Id // 降序
	})

	for k, v := range global.GlobalCfg.War.OpenTime {
		var err error
		var starttime, endtime time.Time
		if starttime, err = time.Parse(time.TimeOnly, v.StartTime); err != nil {
			logger.Panic("战场时间配置错误")
		}

		if endtime, err = time.Parse(time.TimeOnly, v.EndTime); err != nil {
			logger.Panic("战场时间配置错误")
		}

		global.GlobalCfg.War.OpenTime[k].StartHour = starttime.Hour()
		global.GlobalCfg.War.OpenTime[k].StartMinute = starttime.Minute()
		global.GlobalCfg.War.OpenTime[k].EndHour = endtime.Hour()
		global.GlobalCfg.War.OpenTime[k].EndMinute = endtime.Minute()
	}
	logger.Debug("全局配置加载完毕")
}

func loadRoomCfg() {
	global.RoomCfg = []*global.RoomConfig{}
	if data, err := excel.LoadExcel("room.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}
			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}
			gametype, err := strconv.Atoi(row[kv["GameType"]])
			if row[kv["GameType"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}
			minEntryLimit, err := strconv.Atoi(row[kv["MinEntryLimit"]])
			if row[kv["MinEntryLimit"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}
			maxEntryLimit, err := strconv.Atoi(row[kv["MaxEntryLimit"]])
			if row[kv["MaxEntryLimit"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}
			totalTables, err := strconv.Atoi(row[kv["TotalTables"]])
			if row[kv["TotalTables"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}
			playerNum, err := strconv.Atoi(row[kv["PlayerNum"]])
			if row[kv["PlayerNum"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}
			mapId, err := strconv.Atoi(row[kv["MapId"]])
			if row[kv["MapId"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}

			roomtype, err := strconv.Atoi(row[kv["Roomtype"]])
			if row[kv["Roomtype"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}

			needenerge, err := strconv.Atoi(row[kv["Needenerge"]])
			if row[kv["Needenerge"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}

			addtime, err := strconv.Atoi(row[kv["Addtime"]])
			if row[kv["Addtime"]] != "" && err != nil {
				logger.Panicf("room 加载失败 ID :%d", id)
			}

			global.RoomCfg = append(global.RoomCfg, &global.RoomConfig{
				ID: int32(id), GameType: protocol.GameType(gametype),
				MinEntryLimit: minEntryLimit, MaxEntryLimit: maxEntryLimit, Name: row[kv["Name"]],
				TotalTables: totalTables, PlayerNum: int32(playerNum), MapId: int32(mapId), Roomtype: int32(roomtype),
				Needenerge: int32(needenerge), Addtime: int32(addtime),
			})
		}
		logger.Infof("roomconfig 加载完毕,共 [ %d ]  ", len(global.RoomCfg))
	} else {
		logger.Panic("room.xlsx 加载失败")
	}
}

func loadMonsterCfg() {
	global.MonsterCfg = make(map[int32]*global.MonsterConfig)

	if data, err := excel.LoadExcel("monster.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			level, err := strconv.Atoi(row[kv["Level"]])
			if err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}

			var bassProperty []map[int]int
			if err = json.Unmarshal([]byte(row[kv["BassProperty"]]), &bassProperty); row[kv["BassProperty"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 UUID :%d", level)
			}

			var drop []global.DropConfig
			if err := json.Unmarshal([]byte(row[kv["Dorp"]]), &drop); row[kv["Dorp"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}

			basePro := []global.KeyValue{}
			for _, m := range bassProperty {
				for key, value := range m {
					// 添加到结果slice中
					basePro = append(basePro, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			relive, err := strconv.Atoi(row[kv["Relive"]])
			if row[kv["Relive"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			delay, err := strconv.Atoi(row[kv["Delay"]])
			if row[kv["Delay"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			boss, err := strconv.Atoi(row[kv["Isboss"]])
			if row[kv["Isboss"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			isdisplaybar, err := strconv.Atoi(row[kv["Isdisplaybar"]])
			if row[kv["Isdisplaybar"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			trigger, err := strconv.Atoi(row[kv["Trigger"]])
			if row[kv["Trigger"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			experience, err := strconv.Atoi(row[kv["Experience"]])
			if row[kv["Experience"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			viewRange, err := strconv.Atoi(row[kv["ViewRange"]])
			if row[kv["ViewRange"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			attackRange, err := strconv.Atoi(row[kv["AttackRange"]])
			if row[kv["AttackRange"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			ai, err := strconv.Atoi(row[kv["AI"]])
			if row[kv["AI"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			skills := []int{}
			if err = json.Unmarshal([]byte(row[kv["Skills"]]), &skills); row[kv["Skills"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			AttackDistance, err := strconv.Atoi(row[kv["AttackDistance"]])
			if row[kv["AttackDistance"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			chase, err := strconv.Atoi(row[kv["Chase"]])
			if row[kv["Chase"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			distributeRange, err := strconv.Atoi(row[kv["Range"]])
			if row[kv["Range"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}
			quality, err := strconv.Atoi(row[kv["Quality"]])
			if row[kv["Quality"]] != "" && err != nil {
				logger.Panicf("monster 加载失败 ID :%d", id)
			}

			global.MonsterCfg[int32(id)] = &global.MonsterConfig{
				ID: int32(id), Name: row[1], Level: int32(level), BassProperty: basePro,
				Drop: drop, Relive: int32(relive), Delay: int32(delay), MType: int32(boss),
				Isdisplaybar: isdisplaybar == 1, Trigger: int32(trigger),
				Experience: int32(experience), ViewRange: viewRange, AttackRange: attackRange,
				Ai: ai, Skills: skills, AttackDistance: AttackDistance, Chase: int32(chase),
				DistributeRange: distributeRange, Quality: pb.ItemQuality(quality),
			}
		}
		logger.Infof("monstercfg 加载完毕,共 [ %d ] ", len(global.RoomCfg))
	} else {
		logger.Panic("monster.xlsx 加载失败")
	}
}

func loadMapCfg() {
	global.MapCfg = []*global.MapConfig{}
	if data, err := excel.LoadExcel("map.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("map 加载失败 ID :%d", id)
			}
			high, err := strconv.Atoi(row[kv["High"]])
			if row[kv["High"]] != "" && err != nil {
				logger.Panicf("map 加载失败 ID :%d", id)
			}
			wide, err := strconv.Atoi(row[kv["Wide"]])
			if row[kv["Wide"]] != "" && err != nil {
				logger.Panicf("map 加载失败 ID :%d", id)
			}

			var monsters []global.IdNumP
			if err = json.Unmarshal([]byte(row[kv["MonsterIds"]]), &monsters); row[kv["MonsterIds"]] != "" && err != nil {
				logger.Panicf("map 加载失败 ID :%d", id)
			}

			maxMonster, err := strconv.Atoi(row[kv["MaxMonster"]])
			if row[kv["MaxMonster"]] != "" && err != nil {
				logger.Panicf("map 加载失败 ID :%d", id)
			}

			var bornPoints []global.PointI
			if err = json.Unmarshal([]byte(row[kv["Born"]]), &bornPoints); row[kv["Born"]] != "" && err != nil {
				logger.Panicf("map 加载失败 ID :%d", id)
			}

			mapId, err := strconv.Atoi(row[kv["MapID"]])
			if row[kv["MapID"]] != "" && err != nil {
				logger.Panicf("map 加载失败 ID :%d", id)
			}

			mapCfg_ := &global.MapConfig{
				ID: int32(id), High: high, Wide: wide,
				MonsterIds: monsters, MaxMonster: int32(maxMonster), AstarMgr: astar.Map{}, BornPoints: bornPoints, MapID: int32(mapId),
			}

			mapCfg_.Cells = make(map[global.PointI]bool)
			for x := 0; x < int(mapCfg_.Wide); x++ {
				for y := 0; y < int(mapCfg_.High); y++ {
					mapCfg_.Cells[global.PointI{X: x, Y: y}] = true
				}
			}
			// astarMgr := astar.NewMapExx(mapCfg.High, mapCfg.Wide, mapCfg.StrCells)
			astarMgr := astar.NewMapEx(mapCfg_.High, mapCfg_.Wide)
			mapCfg_.AstarMgr = astarMgr

			global.MapCfg = append(global.MapCfg, mapCfg_)
		}
		logger.Infof("map 加载完毕,共 [ %d ]  ", len(global.MapCfg))
	} else {
		logger.Panic("map.xlsx 加载失败")
	}

	//真实加载地图，由于地图无阻挡点，所以加载后直接全部为可走
	/*mapPath := global.BasePath + "/configs/gameconfig/map"
	err := filepath.Walk(mapPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".xml" {
			xmlData, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err)
				return err
			}

			var mapinfo global.MapInfo
			err = xml.Unmarshal(xmlData, &mapinfo)
			if err != nil {
				log.Fatalf("Error parsing XML: %s\n", err)
			}
			valuesStr := strings.Split(mapinfo.Values, ",")

			maps := global.GetMapCfgByMapID(int32(mapinfo.MapId))

			for _, mapid := range maps {
				mapCfg := global.GetMapCfg(mapid)
				if mapCfg == nil {
					logger.Panicf("map.xlsx 加载失败 mapid : %d", mapinfo.MapId)
				}
				mapCfg.High = mapinfo.Height
				mapCfg.Wide = mapinfo.Wide
				mapCfg.StrCells = valuesStr
				mapCfg.Cells = make(map[global.PointI]bool)

				for _, v := range valuesStr {
					_, err := strconv.Atoi(v)
					if err != nil {
						return err
					}
				}
				for x := 0; x < int(mapCfg.Wide); x++ {
					for y := 0; y < int(mapCfg.High); y++ {
						//mapCfg.Cells = append(mapCfg.Cells, global.Cell{L: global.PointI{X: x, Y: y}, V: mapCfg.StrCells[x*mapCfg.High+y]})
						mapCfg.Cells[global.PointI{X: x, Y: y}] = mapCfg.StrCells[x*mapCfg.High+y] == "0"
					}
				}
				//astarMgr := astar.NewMapExx(mapCfg.High, mapCfg.Wide, mapCfg.StrCells)
				astarMgr := astar.NewMapEx(mapCfg.High, mapCfg.Wide)
				mapCfg.AstarMgr = astarMgr
			}
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path: %s\n", err)
	}*/
}

func LoadItemCfg() {
	global.ItemCfg = make(map[int32]*global.ItemConfig)
	// global.ItemStoneCfg = make(map[int32]map[int32]int32)
	if data, err := excel.LoadExcel("item.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}

			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}

			mainType, err := strconv.Atoi(row[kv["Type"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}

			sonType, err := strconv.Atoi(row[kv["SonType"]])
			if row[kv["SonType"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}

			maxNum, err := strconv.Atoi(row[kv["Maxnum"]])
			if row[kv["Maxnum"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}
			if maxNum == 0 || err != nil {
				maxNum = 1
			}
			quality, err := strconv.Atoi(row[kv["Quality"]])
			if row[kv["Quality"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}
			priority, err := strconv.Atoi(row[kv["Priority"]])
			if row[kv["Priority"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}

			du := global.Duration{}
			if err = json.Unmarshal([]byte(row[kv["Timeliness"]]), &du); row[kv["Timeliness"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}
			equiptype, err := strconv.Atoi(row[kv["Equiptype"]])
			if row[kv["Equiptype"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}
			levellimit, err := strconv.Atoi(row[kv["Levellimit"]])
			if row[kv["Levellimit"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}
			occupylimit, err := strconv.Atoi(row[kv["Occupylimit"]])
			if row[kv["Occupylimit"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}
			genderlimit, err := strconv.Atoi(row[kv["Genderlimit"]])
			if row[kv["Genderlimit"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}

			tmpStats := []global.Stats{}
			phyatk, err := strconv.Atoi(row[kv["Phyatk"]])
			if row[kv["Phyatk"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if phyatk > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_PHYSICAL_ATTACK, Value: int32(phyatk)})
			}

			magatk, err := strconv.Atoi(row[kv["Magatk"]])
			if row[kv["Magatk"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if magatk > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_MAGIC_ATTACK, Value: int32(magatk)})
			}

			Phydef, err := strconv.Atoi(row[kv["Phydef"]])
			if row[kv["Phydef"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if Phydef > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_PHYSICAL_DEFENSE, Value: int32(Phydef)})
			}
			Magdef, err := strconv.Atoi(row[kv["Magdef"]])
			if row[kv["Magdef"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if Magdef > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_MAGIC_DEFENSE, Value: int32(Magdef)})
			}
			hp, err := strconv.Atoi(row[kv["HP"]])
			if row[kv["HP"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if hp > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_MAX_HEALTH, Value: int32(hp)})
			}
			hit, err := strconv.Atoi(row[kv["Hit"]])
			if row[kv["Hit"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if hit > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_HIT, Value: int32(hit)})
			}
			miss, err := strconv.Atoi(row[kv["Miss"]])
			if row[kv["Miss"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if miss > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_DODGE, Value: int32(miss)})
			}
			crit, err := strconv.Atoi(row[kv["Crit"]])
			if row[kv["Crit"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if crit > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_CRITICAL_STRIKE, Value: int32(crit)})
			}
			decrit, err := strconv.Atoi(row[kv["Decrit"]])
			if row[kv["Decrit"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if decrit > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_TOUGHNESS, Value: int32(decrit)})
			}
			pvpdamage, err := strconv.Atoi(row[kv["PVPdamage"]])
			if row[kv["PVPdamage"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if pvpdamage > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_PVP_DAMAGE_PERCENTAGE, Value: int32(pvpdamage)})
			}
			pvedamage, err := strconv.Atoi(row[kv["PVEdamage"]])
			if row[kv["PVEdamage"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if pvedamage > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_PVE_DAMAGE_PERCENTAGE, Value: int32(pvedamage)})
			}

			hppercent, err := strconv.Atoi(row[kv["HPPercent"]])
			if row[kv["HPPercent"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			} else if hppercent > 0 {
				tmpStats = append(tmpStats, global.Stats{Key: pb.Attribute_HEALTH_PERCENTAGE, Value: int32(hppercent)})
			}
			var decompose []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Decompose"]]), &decompose); row[kv["Decompose"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}

			RandomVariable, err := strconv.Atoi(row[kv["RandomVariable"]])
			if row[kv["RandomVariable"]] != "" && err != nil {
				logger.Panicf("item 加载失败 ID :%d", id)
			}

			NatureMark, err := strconv.Atoi(row[kv["NatureMark"]])
			if row[kv["NatureMark"]] != "" && err != nil {
				logger.Panicf("item 加载失败 NatureMark :%d", NatureMark)
			}

			Itemlevel, err := strconv.Atoi(row[kv["itemlevel"]])
			if row[kv["itemlevel"]] != "" && err != nil {
				logger.Panicf("item 加载失败 Itemlevel :%d", Itemlevel)
			}

			var DiamondPosition []int
			if err = json.Unmarshal([]byte(row[kv["DiamondPosition"]]), &DiamondPosition); row[kv["DiamondPosition"]] != "" && err != nil {
				logger.Panicf("item 加载失败 DiamondPosition :%d", DiamondPosition)
			}
			userNow, err := strconv.Atoi(row[kv["Usenow"]])
			if row[kv["Usenow"]] != "" && err != nil {
				logger.Panicf("item 加载失败 UserNow :%d", userNow)
			}

			global.ItemCfg[int32(id)] = &global.ItemConfig{
				ID: int32(id), Name: row[1], MainType: int32(mainType), SonType: int32(sonType),
				Maxnum: int64(maxNum), Quality: int32(quality), Priority: int32(priority), Desc: row[4], Icon: row[5],
				Timeliness: du, Equiptype: equiptype, Levellimit: levellimit, Occupylimit: occupylimit, Genderlimit: genderlimit,
				Stats: tmpStats, Decompose: []global.KeyValue{}, RandomVariable: int32(RandomVariable), NatureMark: int32(NatureMark), Itemlevel: int32(Itemlevel), UserNow: int32(userNow),
			}
			for _, m := range decompose {
				for key, value := range m {
					global.ItemCfg[int32(id)].Decompose = append(global.ItemCfg[int32(id)].Decompose, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			if len(DiamondPosition) > 0 {
				for _, n := range DiamondPosition {
					// if len(global.ItemStoneCfg[int32(n)]) == 0 {
					// 	global.ItemStoneCfg[int32(n)] = make(map[int32]int32, 0)
					// }
					// global.ItemStoneCfg[int32(n)][int32(id)] = int32(Itemlevel)
					global.ItemCfg[int32(id)].DiamondPosition = append(global.ItemCfg[int32(id)].DiamondPosition, n)
				}
			}

		}
		logger.Infof("itemcfg 加载完毕,共 [ %d ] ", len(global.ItemCfg))
	} else {
		logger.Panic("item.xlsx 加载失败")
	}
}

func LoadRoleCfg() {
	global.RoleCfg = make(map[global.RoleKey]*global.RoleConfig)

	if data, err := excel.LoadExcel("role.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}
			level, err := strconv.Atoi(row[kv["Level"]])
			if row[kv["Level"]] != "" && err != nil {
				logger.Panicf("role 加载失败 ID :%d", level)
			}
			job, err := strconv.Atoi(row[kv["job"]])
			if row[kv["job"]] != "" && err != nil {
				logger.Panicf("role 加载失败 ID :%d", level)
			}
			exp, err := strconv.Atoi(row[kv["NeedEXP"]])
			if row[kv["NeedEXP"]] != "" && err != nil {
				logger.Panicf("role 加载失败 ID :%d", level)
			}

			var warriorProperty []map[int]int
			if err = json.Unmarshal([]byte(row[kv["WarriorBassProperty"]]), &warriorProperty); row[kv["WarriorBassProperty"]] != "" && err != nil {
				logger.Panicf("role 加载失败 UUID :%d", level)
			}

			// var wizardProperty []map[int]int
			// if err = json.Unmarshal([]byte(row[kv["WizardBassProperty"]]), &wizardProperty); row[kv["WizardBassProperty"]] != "" && err != nil {
			// 	logger.Panicf("role 加载失败 UUID :%d", level)
			// }

			cfg := &global.RoleConfig{Level: int32(level), Job: int32(job), NeedExp: int64(exp)}
			for _, m := range warriorProperty {
				for key, value := range m {
					cfg.BassProperty = append(cfg.BassProperty, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}
			// for _, m := range wizardProperty {
			// 	for key, value := range m {
			// 		cfg.WizardBassProperty = append(cfg.WizardBassProperty, global.KeyValue{
			// 			Key:   int32(key),
			// 			Value: int32(value),
			// 		})
			// 	}
			// }
			key := global.RoleKey{Job: int32(job), Level: int32(level)}
			global.RoleCfg[key] = cfg
		}
		logger.Infof("rolecfg 加载完毕,共 [ %d ] ", len(global.RoleCfg))
	} else {
		logger.Panic("role.xlsx 加载失败")
	}
}

func LoadSkillCfg() {
	global.SkillCfg = make(map[int32]*global.SkillConfig)

	if data, err := excel.LoadExcel("skill.xlsx"); err == nil {
		kv := make(map[string]int)
		lastUpperlevelSetting := global.SkillUpperlevel{}
		var lastSkillId int
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			if uid == 0 {
				continue
			}
			id, err := strconv.Atoi(row[kv["SkillID"]])
			if row[kv["SkillID"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			stype, err := strconv.Atoi(row[kv["Skilltype"]])
			if row[kv["Skilltype"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 Skilltype :%d", uid)
			}
			level, err := strconv.Atoi(row[kv["SkillLevel"]])
			if row[kv["SkillLevel"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			roleLevellimit, err := strconv.Atoi(row[kv["RoleLevellimit"]])
			if row[kv["RoleLevellimit"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			activeorPassive, err := strconv.Atoi(row[kv["ActiveorPassive"]])
			if row[kv["ActiveorPassive"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			pc := []global.IdNum{}
			if err = json.Unmarshal([]byte(row[kv["PromoteConsume"]]), &pc); row[kv["PromoteConsume"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			var effect []map[int]int
			if err = json.Unmarshal([]byte(row[kv["SkillEffect"]]), &effect); row[kv["SkillEffect"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			colDown, err := strconv.Atoi(row[kv["CoolDown"]])
			if row[kv["CoolDown"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			delayTime, err := strconv.Atoi(row[kv["Delay"]])
			if row[kv["Delay"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			maxDistance, err := strconv.Atoi(row[kv["MaxDistance"]])
			if row[kv["MaxDistance"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			radius, err := strconv.Atoi(row[kv["Radius"]])
			if row[kv["Radius"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			var triggerBuff []int32
			err = json.Unmarshal([]byte(row[kv["TriggerBuff"]]), &triggerBuff)
			if row[kv["TriggerBuff"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			skillProfession, err := strconv.Atoi(row[kv["SkillProfession"]])
			if row[kv["SkillProfession"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			canupperLevel, err := strconv.Atoi(row[kv["CanupperLevel"]])
			if row[kv["CanupperLevel"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			var upperlevelConsume []global.IdNum
			err = json.Unmarshal([]byte(row[kv["UpperlevelConsume"]]), &upperlevelConsume)
			if row[kv["UpperlevelConsume"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 upperlevelConsume :%d", uid)
			}

			var upperlevelSetting global.SkillUpperlevel
			err = json.Unmarshal([]byte(row[kv["UpperlevelSetting"]]), &upperlevelSetting)
			if row[kv["UpperlevelSetting"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			if level != 0 {
				if row[kv["UpperlevelSetting"]] == "" {
					if lastSkillId == id {
						upperlevelSetting = lastUpperlevelSetting
					} else {
						lastUpperlevelSetting.DamagePer = 0
						lastUpperlevelSetting.DurationCon = 0
						lastUpperlevelSetting.Distance = 0
						lastUpperlevelSetting.Control = 0
						lastUpperlevelSetting.Trigger = 0
					}
				} else {
					lastUpperlevelSetting = upperlevelSetting
				}
			} else {
				lastUpperlevelSetting.DamagePer = 0
				lastUpperlevelSetting.DurationCon = 0
				lastUpperlevelSetting.Distance = 0
				lastUpperlevelSetting.Control = 0
				lastUpperlevelSetting.Trigger = 0
			}
			lastSkillId = id

			Powerfactor, err := strconv.Atoi(row[kv["Powerfactor"]])
			if row[kv["Powerfactor"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			var Skillbuff []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Skillbuff"]]), &Skillbuff); row[kv["Skillbuff"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			Object, err := strconv.Atoi(row[kv["Object"]])
			if row[kv["Object"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			Needaim, err := strconv.Atoi(row[kv["Needaim"]])
			if row[kv["Needaim"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			Effecttrigger, err := strconv.Atoi(row[kv["Effecttrigger"]])
			if row[kv["Effecttrigger"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			Bufftime, err := strconv.Atoi(row[kv["Bufftime"]])
			if row[kv["Bufftime"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			Skilltrigger, err := strconv.Atoi(row[kv["Skilltrigger"]])
			if row[kv["Skilltrigger"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			effectnum, err := strconv.Atoi(row[kv["Effectnum"]])
			if row[kv["Effectnum"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			buffround, err := strconv.Atoi(row[kv["Buffround"]])
			if row[kv["Buffround"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			roundMark, err := strconv.Atoi(row[kv["RoundMark"]])
			if row[kv["RoundMark"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}
			roundCD, err := strconv.Atoi(row[kv["RoundCD"]])
			if row[kv["RoundCD"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			var specialEffect global.SpecialEffect
			err = json.Unmarshal([]byte(row[kv["SpecialEffect"]]), &specialEffect)
			if row[kv["SpecialEffect"]] != "" && err != nil {
				logger.Panicf("skill 加载失败 UUID :%d", uid)
			}

			cfg := &global.SkillConfig{
				UID:               int32(uid),
				SkillID:           int32(id),
				Level:             int32(level),
				Name:              row[kv["Name"]],
				SType:             int32(stype),
				ActiveorPassive:   int32(activeorPassive),
				RoleLevelLimit:    int32(roleLevellimit),
				PromoteConsume:    pc,
				Description:       row[7],
				CoolDown:          int32(colDown),
				DelayTime:         int32(delayTime),
				MaxDistance:       int32(maxDistance),
				Radius:            int32(radius),
				TriggerBuff:       triggerBuff,
				SkillProfession:   int32(skillProfession),
				CanupperLevel:     int32(canupperLevel),
				UpperlevelSetting: lastUpperlevelSetting,
				UpperlevelConsume: upperlevelConsume,
				Powerfactor:       int32(Powerfactor),
				Object:            int32(Object),
				Needaim:           int32(Needaim),
				Effecttrigger:     int32(Effecttrigger),
				Skilltrigger:      int32(Skilltrigger),
				SpecialEffect:     specialEffect,
				Bufftime:          int32(Bufftime),
				Effectnum:         int32(effectnum),
				Buffround:         int32(buffround),
				RoundMark:         int32(roundMark),
				RoundCD:           int32(roundCD),
			}

			for _, m := range effect {
				for key, value := range m {
					cfg.Effect = append(cfg.Effect, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			for _, m := range Skillbuff {
				for key, value := range m {
					cfg.Skillbuff = append(cfg.Skillbuff, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			global.SkillCfg[int32(uid)] = cfg
		}
		logger.Infof("skill 加载完毕,共 [ %d ] ", len(global.RoleCfg))
	} else {
		logger.Panic("skill.xlsx 加载失败")
	}
}

func LoadSignCfg() {
	global.SignCfg = make(map[string]map[int32]*global.SignConfig)

	if data, err := excel.LoadExcel("sign.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panic("sign.xlsx 加载失败")
			}
			round, err := strconv.Atoi(row[kv["Round"]])
			if row[kv["Round"]] != "" && err != nil {
				logger.Panic("sign.xlsx 加载失败")
			}
			days, err := strconv.Atoi(row[kv["Days"]])
			if row[kv["Days"]] != "" && err != nil {
				logger.Panic("sign.xlsx 加载失败")
			}
			var award []global.Award
			err = json.Unmarshal([]byte(row[kv["Award"]]), &award)
			if row[kv["Award"]] != "" && err != nil {
				logger.Panic("sign.xlsx 加载失败")
			}
			var accumulateaward []global.Award
			err = json.Unmarshal([]byte(row[kv["AccumulateAward"]]), &accumulateaward)
			if row[kv["AccumulateAward"]] != "" && err != nil {
				logger.Panic("sign.xlsx 加载失败")
			}

			if len(accumulateaward) > 0 {
				some2 := "accumulate"
				if _, ok := global.SignCfg[some2]; !ok {
					global.SignCfg[some2] = make(map[int32]*global.SignConfig)
				}
				global.SignCfg[some2][int32(id)] = &global.SignConfig{
					Id:       int32(id),
					Round:    int32(round),
					Days:     int32(days),
					VIPLimit: row[4],
					Content:  accumulateaward,
				}
			}
			some := "day"
			if _, ok := global.SignCfg[some]; !ok {
				global.SignCfg[some] = make(map[int32]*global.SignConfig)
			}
			global.SignCfg[some][int32(id)] = &global.SignConfig{
				Id:       int32(id),
				Round:    int32(round),
				Days:     int32(days),
				VIPLimit: row[kv["VIPLimit"]],
				Content:  award,
			}
		}
		logger.Infof("signcfg 加载完毕,共 [ %d ] ", len(global.SignCfg))
	} else {
		panic("item.xlsx 加载失败")
	}
}

func LoadItemUseCfg() {
	global.ItemUseCfg = make(map[int32]*global.ItemUseConfig)

	if data, err := excel.LoadExcel("itemuse.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("itemuse 加载失败 ID :%d", uid)
			}
			if uid == 0 {
				continue
			}

			itemUseType, err := strconv.Atoi(row[kv["Type"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("itemuse 加载失败 Type :%d", uid)
			}

			ingw := []global.IdNumJobWeight{}
			err = json.Unmarshal([]byte(row[kv["Givegoods"]]), &ingw)
			if row[kv["Givegoods"]] != "" && err != nil {
				logger.Panicf("itemuse 加载失败 Givegoods :%d", uid)
			}

			ikn := []global.IdKeyNum{}
			err = json.Unmarshal([]byte(row[kv["Buffeffect"]]), &ikn)
			if row[kv["Buffeffect"]] != "" && err != nil {
				logger.Panicf("itemuse 加载失败 Buffeffect :%d", uid)
			}

			useLimit, err := strconv.Atoi(row[kv["UseLimit"]])
			if err != nil {
				logger.Panicf("itemuse 加载失败 UseLimit :%d", uid)
			}

			cfg := &global.ItemUseConfig{ID: int32(uid), Name: row[kv["Name"]], Type: int32(itemUseType), GiveGoods: ingw, Buffeffect: ikn, UseLimit: int64(useLimit)}

			global.ItemUseCfg[int32(uid)] = cfg
		}
		logger.Infof("itemuse 加载完毕,共 [ %d ] ", len(global.ItemUseCfg))
	} else {
		logger.Panic("itemuse.xlsx 加载失败")
	}
}

func LoadStoreCfg() {
	global.StoreCfg = make(map[int32]*global.StoreConfig)

	if data, err2 := excel.LoadExcel("store.xlsx"); err2 == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panic("store.xlsx 加载失败")
			}
			storeType, err := strconv.Atoi(row[kv["StoreType"]])
			if row[kv["StoreType"]] != "" && err != nil {
				logger.Panic("store.xlsx 加载失败")
			}
			goodsId, err := strconv.Atoi(row[kv["GoodsID"]])
			if row[kv["GoodsID"]] != "" && err != nil {
				logger.Panic("store.xlsx 加载失败")
			}
			var limit global.Limit
			err = json.Unmarshal([]byte(row[kv["Limite"]]), &limit)
			if row[kv["Limite"]] != "" && err != nil {
				logger.Panic("store.xlsx 加载失败")
			}
			var props []global.Props
			err = json.Unmarshal([]byte(row[kv["Solditems"]]), &props)
			if row[kv["Solditems"]] != "" && err != nil {
				logger.Panic("store.xlsx 加载失败")
			}

			preLimit, err := strconv.Atoi(row[kv["Prelimit"]])
			if row[kv["Prelimit"]] != "" && err != nil {
				logger.Panic("store.xlsx 加载失败")
			}
			entrance := row[kv["entrance"]]
			method := row[kv["method"]]

			var Consumechange []global.KeyValue
			err = json.Unmarshal([]byte(row[kv["Consumechange"]]), &Consumechange)
			if row[kv["Consumechange"]] != "" && err != nil {
				logger.Panic("store.xlsx 加载失败")
			}

			//global.SignCfg = append(global.SignCfg, &global.SignConfig{
			//	ID:              int32(id),
			//	Round:           int32(round),
			//	Days:            int32(days),
			//	Award:           award,
			//	VIPLimit:        row[4],
			//	AccumulateAward: accumulateaward,
			//})

			global.StoreCfg[int32(id)] = &global.StoreConfig{
				Id:            int32(id),
				StoreType:     int32(storeType),
				GoodsId:       int64(goodsId),
				PreLimit:      int32(preLimit),
				Limit:         limit,
				Props:         props,
				Entrance:      entrance,
				Method:        method,
				Consumechange: Consumechange,
			}
		}
		logger.Infof("storecfg 加载完毕,共 [ %d ] ", len(global.SignCfg))
	} else {
		panic("store.xlsx 加载失败")
	}
}

func LoadActivityCfg() {
	global.ActivityCfg = make(map[int32]*global.ActivityConfig)
	global.ActivityCfgByType = make(map[int32][]*global.ActivityConfig)

	if data, err := excel.LoadExcel("activity.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}
			UUID, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 ID :%d", UUID)
			}
			Type, err := strconv.Atoi(row[kv["Type"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 Type :%d", Type)
			}
			Content, err := strconv.Atoi(row[kv["Content"]])
			if row[kv["Content"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 ID :%d", Content)
			}
			Acttag, err := strconv.Atoi(row[kv["Acttag"]])
			if row[kv["Acttag"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 ID :%d", Acttag)
			}
			PreLimit, err := strconv.Atoi(row[kv["PreLimit"]])
			if row[kv["PreLimit"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 PreLimit :%d", PreLimit)
			}
			ActivityTime := row[kv["ActivityTime"]]
			if row[kv["ActivityTime"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 ActivityTime :%s", ActivityTime)
			}
			RuleDes := row[kv["RuleDes"]]
			if row[kv["RuleDes"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 RuleDes :%s", RuleDes)
			}
			Name := row[kv["Name"]]
			if row[kv["Name"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 Name :%s", Name)
			}
			ResetTime, err := strconv.Atoi(row[kv["ResetTime"]])
			if row[kv["ResetTime"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 ResetTime :%d", ResetTime)
			}
			StartTime, err := strconv.Atoi(row[kv["StartTime"]])
			if row[kv["StartTime"]] != "" && err != nil {
				logger.Panicf("Activity 加载失败 StartTime :%d", StartTime)
			}
			var CommonItem []map[int]int
			if err = json.Unmarshal([]byte(row[kv["CommonItem"]]), &CommonItem); row[kv["CommonItem"]] != "" && err != nil {
				logger.Panicln("activity 加载失败 CommonItem ", key, CommonItem, row[kv["CommonItem"]], err.Error())
			}

			var PayItem []map[int]int
			if err = json.Unmarshal([]byte(row[kv["PayItem"]]), &PayItem); row[kv["PayItem"]] != "" && err != nil {
				logger.Panicf("activity 加载失败 PayItem :%d", PayItem)
			}

			var CompensationConsume []map[int]int
			if err = json.Unmarshal([]byte(row[kv["CompensationConsume"]]), &CompensationConsume); row[kv["CompensationConsume"]] != "" && err != nil {
				logger.Panicf("activity 加载失败 CompensationConsume :%d", CompensationConsume)
			}

			cfg := &global.ActivityConfig{
				UUID:         int64(UUID),
				Type:         int32(Type),
				Content:      int32(Content),
				Acttag:       int32(Acttag),
				PreLimit:     int32(PreLimit),
				ActivityTime: ActivityTime,
				RuleDes:      RuleDes,
				Name:         Name,
				ResetTime:    int32(ResetTime),
				StartTime:    int32(StartTime),
			}
			for _, m := range CommonItem {
				for key2, value := range m {
					cfg.CommonItem = append(cfg.CommonItem, global.KeyValue{
						Key:   int32(key2),
						Value: int32(value),
					})
				}
			}

			for _, m := range PayItem {
				for key2, value := range m {
					cfg.PayItem = append(cfg.PayItem, global.KeyValue{
						Key:   int32(key2),
						Value: int32(value),
					})
				}
			}

			for _, m := range CompensationConsume {
				for key2, value := range m {
					cfg.CompensationConsume = append(cfg.CompensationConsume, global.KeyValue{
						Key:   int32(key2),
						Value: int32(value),
					})
				}
			}
			global.ActivityCfg[int32(UUID)] = cfg

			if _, ok := global.ActivityCfgByType[int32(Type)]; !ok {
				global.ActivityCfgByType[int32(Type)] = make([]*global.ActivityConfig, 0, 10)
			}
			global.ActivityCfgByType[int32(Type)] = append(global.ActivityCfgByType[int32(Type)], cfg)
		}
		logger.Infof("activitycfg 加载完毕,共 [ %d ] ", len(global.ActivityCfg))
	} else {
		logger.Panic("activity.xlsx 加载失败")
	}
}

func LoadBuffCfg() {
	global.BuffCfg = make(map[int32]*global.BuffConfig)
	if data, err := excel.LoadExcel("buffeffect.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}

			buffType, err := strconv.Atoi(row[kv["Type"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}

			mark, err := strconv.Atoi(row[kv["Mark"]])
			if row[kv["Mark"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}

			describe := row[kv["Describe"]]

			var values []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Value"]]), &values); row[kv["Value"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 UUID :%d", uid)
			}

			rate, err := strconv.Atoi(row[kv["Rate"]])
			if row[kv["Rate"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}
			lastTime, err := strconv.Atoi(row[kv["LastTime"]])
			if row[kv["LastTime"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}
			cDTime, err := strconv.Atoi(row[kv["CDTime"]])
			if row[kv["CDTime"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}
			icon := row[kv["Icon"]]

			buffEffect := row[kv["BuffEffect"]]

			canPickUp, err := strconv.Atoi(row[kv["CanPick"]])
			if row[kv["CanPick"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}

			effectTime, err := strconv.Atoi(row[kv["Effecttime"]])
			if row[kv["Effecttime"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}
			effectRange, err := strconv.Atoi(row[kv["EffectRange"]])
			if row[kv["EffectRange"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}
			damIncrease, err := strconv.Atoi(row[kv["DamIncrease"]])
			if row[kv["DamIncrease"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}
			pretype, err := strconv.Atoi(row[kv["Pretype"]])
			if row[kv["Pretype"]] != "" && err != nil {
				logger.Panicf("buffeffect 加载失败 ID :%d", uid)
			}

			cfg := &global.BuffConfig{
				ID:          int32(uid),
				Type:        int32(buffType),
				Mark:        int32(mark),
				Describe:    describe,
				Value:       []global.KeyValue{},
				Rate:        int32(rate),
				LastTime:    int32(lastTime),
				CDTime:      int32(cDTime),
				Icon:        icon,
				BuffEffect:  buffEffect,
				CanPickUp:   int32(canPickUp),
				EffectTime:  int32(effectTime),
				EffectRange: int32(effectRange),
				DamIncrease: int32(damIncrease),
				Pretype:     int32(pretype),
			}
			for _, v := range values {
				for key2, value := range v {
					cfg.Value = append(cfg.Value, global.KeyValue{
						Key:   int32(key2),
						Value: int32(value),
					})
				}
			}

			global.BuffCfg[int32(uid)] = cfg
		}
		logger.Infof("buffeffect 加载完毕,共 [ %d ]  ", len(global.RoomCfg))
	} else {
		logger.Panic("buffeffect.xlsx 加载失败")
	}
}

func LoadTaskCfg() {
	global.TaskCfg = make(map[int32]*global.TaskConfig)

	if data, err := excel.LoadExcel("task.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}

			Class, err := strconv.Atoi(row[kv["Class"]])
			if row[kv["Class"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}

			Type, err := strconv.Atoi(row[kv["Type"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}

			Typeson, err := strconv.Atoi(row[kv["Typeson"]])
			if row[kv["Typeson"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}

			Limit, err := strconv.Atoi(row[kv["limit"]])
			if row[kv["limit"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}

			var values []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Reward"]]), &values); row[kv["Reward"]] != "" && err != nil {
				logger.Panicf("task 加载失败 UUID :%d", uid)
			}

			NextAchievement, err := strconv.Atoi(row[kv["NextAchievement"]])
			if row[kv["NextAchievement"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}

			var Level []int32
			if err = json.Unmarshal([]byte(row[kv["Level"]]), &Level); row[kv["Level"]] != "" && err != nil {
				logger.Panicf("task 加载失败 UUID :%d", uid)
			}

			Job, err := strconv.Atoi(row[kv["Job"]])
			if row[kv["Job"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}
			VipLevel, err := strconv.Atoi(row[kv["VIPLevel"]])
			if row[kv["VIPLevel"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}
			SevenDay, err := strconv.Atoi(row[kv["SevenDay"]])
			if row[kv["SevenDay"]] != "" && err != nil {
				logger.Panicf("task 加载失败 SevenDay :%d", uid)
			}
			ActivID, err := strconv.Atoi(row[kv["ActivID"]])
			if row[kv["ActivID"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}
			ResetMark, err := strconv.Atoi(row[kv["ResetMark"]])
			if row[kv["ResetMark"]] != "" && err != nil {
				logger.Panicf("task 加载失败 ID :%d", uid)
			}
			TaskDes := row[kv["TaskDes"]]
			var ActivityReward global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["ActivityReward"]]), &ActivityReward); row[kv["ActivityReward"]] != "" && err != nil {
				logger.Panicf("task 加载失败 UUID :%d", uid)
			}

			var start_end global.StartEnd
			if err = json.Unmarshal([]byte(row[kv["LimitTime"]]), &start_end); row[kv["LimitTime"]] != "" && err != nil {
				logger.Panicf("task 加载失败 UUID :%d", uid)
			}

			var firstRound int
			if err = json.Unmarshal([]byte(row[kv["TaskMark"]]), &firstRound); row[kv["TaskMark"]] != "" && err != nil {
				logger.Panicf("task 加载失败 UUID :%d", uid)
			}

			cfg := &global.TaskConfig{
				UUID:  int32(uid),
				Class: int32(Class),
				Type:  int32(Type), Typeson: int32(Typeson),
				Limit: int32(Limit), Reward: []global.KeyValue{},
				NextAchievement: int32(NextAchievement), Level: Level,
				Job: int32(Job), VipLevel: int32(VipLevel), FirstRound: firstRound > 0,
				ResetMark: ResetMark == 1, LimitTime: start_end,
				ActivityReward: ActivityReward,
				TaskDes:        TaskDes,
				ActivID:        int32(ActivID),
				SevenDay:       int32(SevenDay),
			}

			for _, v := range values {
				for key2, value := range v {
					cfg.Reward = append(cfg.Reward, global.KeyValue{
						Key:   int32(key2),
						Value: int32(value),
					})
				}
			}
			global.TaskCfg[int32(uid)] = cfg
		}
		logger.Infof("task 加载完毕,共 [ %d ]  ", len(global.TaskCfg))
	} else {
		logger.Panic("task.xlsx 加载失败")
	}
}

func LoadPrivilegeCfg() {
	global.PrivilegeCfg = make(map[int32]*global.PrivilegeConfig)

	if data, err1 := excel.LoadExcel("privilege.xlsx"); err1 == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			pid, err := strconv.Atoi(row[kv["PID"]])
			if row[kv["PID"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载ID失败")
			}
			sid, err := strconv.Atoi(row[kv["SID"]])
			if row[kv["SID"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载\"SID\"失败")
			}
			pname := row[kv["Pname"]]
			if row[kv["Pname"]] == "" {
				logger.Panic("privilege.xlsx 加载Pname失败")
			}
			var dailyaward []global.Award
			err = json.Unmarshal([]byte(row[kv["Dailyaward"]]), &dailyaward)
			if row[kv["Dailyaward"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载Dailyaward失败")
			}

			// pdisplay := row[kv["Pdisplay"]]
			// if row[kv["Pdisplay"]] == "" {
			// 	logger.Panic("privilege.xlsx 加载Pdisplay失败")
			// }
			var pcontent []global.KeyValue
			err = json.Unmarshal([]byte(row[kv["Pcontent"]]), &pcontent)
			if row[kv["Pcontent"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载Pcontent失败")
			}
			duration, err := strconv.Atoi(row[kv["Duration"]])
			if row[kv["Duration"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载Duration失败")
			}
			energystate, err := strconv.Atoi(row[kv["Energystate"]])
			if row[kv["Energystate"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载Energystate失败")
			}

			SkipAd, err := strconv.Atoi(row[kv["SkipAd"]])
			if row[kv["SkipAd"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载ID失败")
			}

			StaminaUp, err := strconv.Atoi(row[kv["StaminaUp"]])
			if row[kv["StaminaUp"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载ID失败")
			}
			QuicklyIsiand, err := strconv.Atoi(row[kv["QuicklyIsiand"]])
			if row[kv["QuicklyIsiand"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载ID失败")
			}
			QuicklyDaily, err := strconv.Atoi(row[kv["QuicklyDaily"]])
			if row[kv["QuicklyDaily"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载ID失败")
			}
			AutoBattle, err := strconv.Atoi(row[kv["AutoBattle"]])
			if row[kv["AutoBattle"]] != "" && err != nil {
				logger.Panic("privilege.xlsx 加载ID失败")
			}

			global.PrivilegeCfg[int32(pid)] = &global.PrivilegeConfig{
				Pid:           int32(pid),
				Sid:           int32(sid),
				Pname:         pname,
				Dailyaward:    dailyaward,
				Pcontent:      pcontent,
				Duration:      int32(duration),
				Energystate:   int32(energystate),
				SkipAd:        int32(SkipAd),
				StaminaUp:     int32(StaminaUp),
				QuicklyIsiand: int32(QuicklyIsiand),
				QuicklyDaily:  int32(QuicklyDaily),
				AutoBattle:    int32(AutoBattle),
			}
		}
		logger.Infof("privilegecfg 加载完毕,共 [ %d ] ", len(global.PrivilegeCfg))
	} else {
		panic("privilege.xlsx 加载失败")
	}
}

func LoadRankingAwardCfg() {
	global.RankingAwardCfg = make(map[int32]*global.RankingAwardConfig)

	if data, err1 := excel.LoadExcel("rankingaward.xlsx"); err1 == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uuid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panic("rankingaward.xlsx 加载UUID失败")
			}
			rtype, err := strconv.Atoi(row[kv["type"]])
			if row[kv["type"]] != "" && err != nil {
				logger.Panic("rankingaward.xlsx 加载\"type\"失败")
			}
			rank, err := strconv.Atoi(row[kv["rank"]])
			if row[kv["rank"]] != "" && err != nil {
				logger.Panic("rankingaward.xlsx 加载\"rank\"失败")
			}
			var award []map[int]int
			err = json.Unmarshal([]byte(row[kv["award"]]), &award)
			if row[kv["award"]] != "" && err != nil {
				logger.Panic("rankingaward.xlsx award")
			}
			cfg := &global.RankingAwardConfig{
				UUID:  int32(uuid),
				Type:  int32(rtype),
				Rank:  int32(rank),
				Award: []global.KeyValue{},
			}
			for _, v := range award {
				for key2, value := range v {
					cfg.Award = append(cfg.Award, global.KeyValue{
						Key:   int32(key2),
						Value: int32(value),
					})
				}
			}
			global.RankingAwardCfg[int32(uuid)] = cfg

		}
		logger.Infof("RankingAwardcfg 加载完毕,共 [ %d ] ", len(global.SignCfg))
	} else {
		panic("rankingaward.xlsx 加载失败")
	}
}

func LoadRandomCreatCfg() {
	global.RandomCreatCfg = []*global.RandomCreat{}

	if data, err := excel.LoadExcel("randomcreat.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			RandomID, err := strconv.Atoi(row[kv["RandomID"]])
			if row[kv["RandomID"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", RandomID)
			}

			var SetID []int32
			if err = json.Unmarshal([]byte(row[kv["SetID"]]), &SetID); row[kv["SetID"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", RandomID)
			}

			var CatalogueNum []map[int]int
			if err = json.Unmarshal([]byte(row[kv["CatalogueNum"]]), &CatalogueNum); row[kv["CatalogueNum"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", RandomID)
			}

			var KindNum []map[int]int
			if err = json.Unmarshal([]byte(row[kv["KindNum"]]), &KindNum); row[kv["KindNum"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", RandomID)
			}
			var DropSet []map[int]int
			if err = json.Unmarshal([]byte(row[kv["DropSet"]]), &DropSet); row[kv["DropSet"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", RandomID)
			}

			var RefingSet []map[int]int
			if err = json.Unmarshal([]byte(row[kv["RefingSet"]]), &RefingSet); row[kv["RefingSet"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", RandomID)
			}
			quality, err := strconv.Atoi(row[kv["quality"]])
			if row[kv["quality"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", RandomID)
			}

			var weight int
			cfg := &global.RandomCreat{
				RandomID:     int32(RandomID),
				CatalogueNum: []global.KeyValue{},
				KindNum:      []global.KeyValue{},
				DropSet:      []global.KeyValue{},
				RefiningSet:  []global.KeyValue{},
				SetID:        SetID,
				Quality:      int32(quality),
			}

			for _, v := range CatalogueNum {
				for key, value := range v {
					weight += value
					cfg.CatalogueNum = append(cfg.CatalogueNum, global.KeyValue{
						Key:   int32(key),
						Value: int32(weight),
					})
				}
			}
			weight = 0
			for _, v := range KindNum {
				for key, value := range v {
					weight += value
					cfg.KindNum = append(cfg.KindNum, global.KeyValue{
						Key:   int32(key),
						Value: int32(weight),
					})
				}
			}
			weight = 0
			for _, v := range DropSet {
				for key, value := range v {
					weight += value
					cfg.DropSet = append(cfg.DropSet, global.KeyValue{
						Key:   int32(key),
						Value: int32(weight),
					})
				}
			}
			weight = 0
			for _, v := range RefingSet {
				for key, value := range v {
					weight += value
					cfg.RefiningSet = append(cfg.RefiningSet, global.KeyValue{
						Key:   int32(key),
						Value: int32(weight),
					})
				}
			}

			global.RandomCreatCfg = append(global.RandomCreatCfg, cfg)
		}
		logger.Infof("randomCreat 加载完毕,共 [ %d ]  ", len(global.RandomCreatCfg))
	} else {
		logger.Panic("randomCreat.xlsx 加载失败")
	}
}

func LoadRandomDistribCfg() {
	global.RandomDistribCfg = []*global.RandomDistrib{}

	if data, err := excel.LoadExcel("randomdistrib.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			UUID, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			Kind, err := strconv.Atoi(row[kv["Kind"]])
			if row[kv["Kind"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			GroupID, err := strconv.Atoi(row[kv["GroupID"]])
			if row[kv["GroupID"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			ValueID, err := strconv.Atoi(row[kv["ValueID"]])
			if row[kv["ValueID"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			// PieceMin_Blue, err := strconv.Atoi(row[kv["PieceMin_Blue"]])
			// if row[kv["PieceMin_Blue"]] != "" && err != nil {
			// 	logger.Panicf("task 加载失败 RandomID :%d", UUID)
			// }
			// PieceMax_Blue, err := strconv.Atoi(row[kv["PieceMax_Blue"]])
			// if row[kv["PieceMax_Blue"]] != "" && err != nil {
			// 	logger.Panicf("task 加载失败 RandomID :%d", UUID)
			// }
			PieceMin_Purple, err := strconv.Atoi(row[kv["PieceMin_Purple"]])
			if row[kv["PieceMin_Purple"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			PieceMax_Purple, err := strconv.Atoi(row[kv["PieceMax_Purple"]])
			if row[kv["PieceMax_Purple"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			PieceMin_Pink, err := strconv.Atoi(row[kv["PieceMin_Pink"]])
			if row[kv["PieceMin_Pink"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			PieceMax_Pink, err := strconv.Atoi(row[kv["PieceMax_Pink"]])
			if row[kv["PieceMax_Pink"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			PieceMin_Gold, err := strconv.Atoi(row[kv["PieceMin_Gold"]])
			if row[kv["PieceMin_Gold"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			PieceMax_Gold, err := strconv.Atoi(row[kv["PieceMax_Gold"]])
			if row[kv["PieceMax_Gold"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}

			PieceMin_Red, err := strconv.Atoi(row[kv["PieceMin_Red"]])
			if row[kv["PieceMin_Red"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			PieceMax_Red, err := strconv.Atoi(row[kv["PieceMax_Red"]])
			if row[kv["PieceMax_Red"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}

			PieceMin_Yellow, err := strconv.Atoi(row[kv["PieceMin_Yellow"]])
			if row[kv["PieceMin_Yellow"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}
			PieceMax_Yellow, err := strconv.Atoi(row[kv["PieceMax_Yellow"]])
			if row[kv["PieceMax_Yellow"]] != "" && err != nil {
				logger.Panicf("task 加载失败 RandomID :%d", UUID)
			}

			Interval := map[int32]global.Interval{}
			// Interval[int32(pb.ItemQuality_BLUE)] = global.Interval{Min: int32(PieceMin_Blue), Max: int32(PieceMax_Blue)}
			Interval[int32(pb.ItemQuality_PURPLE)] = global.Interval{Min: int32(PieceMin_Purple), Max: int32(PieceMax_Purple)}
			Interval[int32(pb.ItemQuality_PINK)] = global.Interval{Min: int32(PieceMin_Pink), Max: int32(PieceMax_Pink)}
			Interval[int32(pb.ItemQuality_ORANGE)] = global.Interval{Min: int32(PieceMin_Gold), Max: int32(PieceMax_Gold)}
			Interval[int32(pb.ItemQuality_RED)] = global.Interval{Min: int32(PieceMin_Red), Max: int32(PieceMax_Red)}
			Interval[int32(pb.ItemQuality_Rainbow)] = global.Interval{Min: int32(PieceMin_Yellow), Max: int32(PieceMax_Yellow)}

			global.RandomDistribCfg = append(global.RandomDistribCfg, &global.RandomDistrib{
				UUID:     int32(UUID),
				Kind:     int32(Kind),
				GroupID:  int32(GroupID),
				ValueID:  int32(ValueID),
				Interval: Interval,
			})
		}
		logger.Infof("randomDistrib 加载完毕,共 [ %d ]  ", len(global.RandomDistribCfg))
	} else {
		logger.Panic("randomDistrib.xlsx 加载失败")
	}
}

func LoadVipCfg() {
	global.VIPCfg = []*global.VIPInfo{}
	if data, err := excel.LoadExcel("vip.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			VIPLevel, err := strconv.Atoi(row[kv["VIPLevel"]])
			if row[kv["VIPLevel"]] != "" && err != nil {
				logger.Panicf("vip 加载失败 VIPLevel :%d", VIPLevel)
			}

			NeedExp, err := strconv.Atoi(row[kv["NeedEXP"]])
			if row[kv["NeedExp"]] != "" && err != nil {
				logger.Panicf("vip 加载失败 NeedExp :%d", NeedExp)
			}

			var Pcontent []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Pcontent"]]), &Pcontent); row[kv["Pcontent"]] != "" && err != nil {
				logger.Panicf("vip 加载失败 Pcontent :%d", Pcontent)
			}
			Partnerslot, err := strconv.Atoi(row[kv["Partnerslot"]])
			if row[kv["Partnerslot"]] != "" && err != nil {
				logger.Panicf("vip 加载失败 Partnerslot :%d", Partnerslot)
			}
			Energyadd, err := strconv.Atoi(row[kv["energyadd"]])
			if row[kv["energyadd"]] != "" && err != nil {
				logger.Panicf("vip 加载失败 Energyadd :%d", Energyadd)
			}

			var dailytime []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["dailytime"]]), &dailytime); row[kv["dailytime"]] != "" && err != nil {
				logger.Panicf("vip 加载失败 dailytime :%d", dailytime)
			}

			cfg := &global.VIPInfo{VIPLevel: int32(VIPLevel), NeedExp: int32(NeedExp), Pcontent: Pcontent, DailyTime: dailytime, Partnerslot: int32(Partnerslot), Energyadd: int32(Energyadd)}
			global.VIPCfg = append(global.VIPCfg, cfg)
		}
		logger.Infof("vip 加载完毕,共 [ %d ]  ", len(global.VIPCfg))
	} else {
		logger.Panic("vip.xlsx 加载失败")
	}
}

func LoadRobotCfg() {
	global.RobotCfg = []*global.RobotConfig{}
	if data, err := excel.LoadExcel("robot.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			UUID, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("robot 加载失败 UUID :%d", UUID)
			}

			Job, err := strconv.Atoi(row[kv["Job"]])
			if row[kv["Job"]] != "" && err != nil {
				logger.Panicf("robot 加载失败 UUID :%d", UUID)
			}

			Power, err := strconv.Atoi(row[kv["Power"]])
			if row[kv["Power"]] != "" && err != nil {
				logger.Panicf("robot 加载失败 UUID :%d", UUID)
			}
			Level, err := strconv.Atoi(row[kv["Level"]])
			if row[kv["Level"]] != "" && err != nil {
				logger.Panicf("robot 加载失败 UUID :%d", UUID)
			}

			RankMark, err := strconv.Atoi(row[kv["RankMark"]])
			if row[kv["RankMark"]] != "" && err != nil {
				logger.Panicf("robot 加载失败 UUID :%d", UUID)
			}

			officialEXP, err := strconv.Atoi(row[kv["OfficialEXP"]])
			if row[kv["OfficialEXP"]] != "" && err != nil {
				logger.Panicf("robot 加载失败 UUID :%d", UUID)
			}

			var Feature []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Feature"]]), &Feature); row[kv["Feature"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", UUID)
			}

			var Skill []int
			if err = json.Unmarshal([]byte(row[kv["Skill"]]), &Skill); row[kv["Skill"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", UUID)
			}

			var partners []global.IdKeyNum
			if err = json.Unmarshal([]byte(row[kv["Partners"]]), &partners); row[kv["Partners"]] != "" && err != nil {
				logger.Panicf("randomCreat 加载失败 UUID :%d", UUID)
			}

			Name := row[kv["Name"]]

			cfg := &global.RobotConfig{
				UUID:        int32(UUID),
				Job:         int32(Job),
				Power:       int32(Power),
				Level:       int32(Level),
				Feature:     []global.KeyValue{},
				Skill:       Skill,
				Partners:    partners,
				Name:        Name,
				RankMark:    int32(RankMark),
				OfficialEXP: int64(officialEXP),
			}
			for _, v := range Feature {
				for key, value := range v {
					cfg.Feature = append(cfg.Feature, global.KeyValue{Key: int32(key), Value: int32(value)})
				}
			}

			global.RobotCfg = append(global.RobotCfg, cfg)
		}
		logger.Infof("robot 加载完毕,共 [ %d ]  ", len(global.RobotCfg))
	} else {
		logger.Panic("robot.xlsx 加载失败")
	}
}

func LoadEnhanceCfg() {
	global.EnhanceCfg = make(map[int32]*global.EnhanceConfig)

	if data, err1 := excel.LoadExcel("equipstrengthen.xlsx"); err1 == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uuid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panic("equipstrengthen.xlsx 加载UUID失败")
			}
			rtype, err := strconv.Atoi(row[kv["Job"]])
			if row[kv["Job"]] != "" && err != nil {
				logger.Panic("equipstrengthen.xlsx 加载\"type\"失败")
			}
			rank, err := strconv.Atoi(row[kv["Strengthenlevel"]])
			if row[kv["Strengthenlevel"]] != "" && err != nil {
				logger.Panic("equipstrengthen.xlsx 加载\"rank\"失败")
			}
			var consume []global.IdNum
			err = json.Unmarshal([]byte(row[kv["Consume"]]), &consume)
			if row[kv["Consume"]] != "" && err != nil {
				logger.Panic("equipstrengthen.xlsx Consume")
			}
			var attribute []global.Attribute
			err = json.Unmarshal([]byte(row[kv["Attribute"]]), &attribute)
			if row[kv["Attribute"]] != "" && err != nil {
				logger.Panic("equipstrengthen.xlsx Attribute", err)
			}
			cfg := &global.EnhanceConfig{
				UUID:            int32(uuid),
				Job:             int32(rtype),
				Strengthenlevel: int32(rank),
				Consume:         []global.IdNum{},
				Attribute:       []global.Attribute{},
			}
			for _, v := range consume {
				cfg.Consume = append(cfg.Consume, global.IdNum{
					ID:  v.ID,
					Num: v.Num,
				})
			}

			for _, v := range attribute {
				cfg.Attribute = append(cfg.Attribute, global.Attribute{
					Slot:      v.Slot,
					Attribute: v.Attribute,
					Value:     v.Value,
				})
			}

			global.EnhanceCfg[int32(uuid)] = cfg

		}
		logger.Infof("equipstrengthen 加载完毕,共 [ %d ] ", len(global.SignCfg))
	} else {
		panic("equipstrengthen.xlsx 加载失败")
	}
}

func LoadStarsCfg() {
	global.StarsCfg = make(map[int32]*global.StarsConfig)

	if data, err1 := excel.LoadExcel("equipstar.xlsx"); err1 == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uuid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx 加载UUID失败")
			}
			job, err := strconv.Atoi(row[kv["Job"]])
			if row[kv["Job"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx 加载\"Job\"失败")
			}
			starLevel, err := strconv.Atoi(row[kv["StarLevel"]])
			if row[kv["StarLevel"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx 加载\"StarLevel\"失败")
			}
			decLevel, err := strconv.Atoi(row[kv["DecLevel"]])
			if row[kv["DecLevel"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx 加载\"DecLevel\"失败")
			}
			failLevel, err := strconv.Atoi(row[kv["FailLevel"]])
			if row[kv["FailLevel"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx 加载\"FailLevel\"失败")
			}
			successRate, err := strconv.Atoi(row[kv["SuccessRate"]])
			if row[kv["SuccessRate"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx 加载\"SuccessRate\"失败")
			}
			var materialID []global.IdNum
			err = json.Unmarshal([]byte(row[kv["MaterialID"]]), &materialID)
			if row[kv["MaterialID"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx MaterialID")
			}
			var attribute []global.Attribute
			err = json.Unmarshal([]byte(row[kv["Attribute"]]), &attribute)
			if row[kv["Attribute"]] != "" && err != nil {
				logger.Panic("equipstar.xlsx Attribute")
			}
			cfg := &global.StarsConfig{
				UUID:        int32(uuid),
				Job:         int32(job),
				StarLevel:   int32(starLevel),
				DecLevel:    int32(decLevel),
				FailLevel:   int32(failLevel),
				SuccessRate: int32(successRate),
				MaterialID:  []global.IdNum{},
				Attribute:   []global.Attribute{},
			}
			for _, v := range materialID {
				cfg.MaterialID = append(cfg.MaterialID, global.IdNum{
					ID:  v.ID,
					Num: v.Num,
				})
			}

			for _, v := range attribute {
				cfg.Attribute = append(cfg.Attribute, global.Attribute{
					Slot:      v.Slot,
					Attribute: v.Attribute,
					Value:     v.Value,
				})
			}
			global.StarsCfg[int32(uuid)] = cfg
		}
		logger.Infof("equipstar 加载完毕,共 [ %d ] ", len(global.SignCfg))
	} else {
		panic("equipstar.xlsx 加载失败")
	}
}

func LoadEndlesstrialsCfg() {
	global.EndlesstrialsCfg = make(map[int32]*global.EndlesstrialsConfig)

	if data, err := excel.LoadExcel("endlesstrials.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			Layers, err := strconv.Atoi(row[kv["Layers"]])
			if row[kv["Layers"]] != "" && err != nil {
				logger.Panicf("endlesstrials 加载失败 ID :%d", Layers)
			}
			Levellimit, err := strconv.Atoi(row[kv["Levellimit"]])
			if row[kv["Levellimit"]] != "" && err != nil {
				logger.Panicf("endlesstrials 加载失败 ID :%d", Layers)
			}

			Timelimit, err := strconv.Atoi(row[kv["Timelimit"]])
			if row[kv["Timelimit"]] != "" && err != nil {
				logger.Panicf("endlesstrials 加载失败 ID :%d", Layers)
			}

			var Monster []global.IdNumP
			if err = json.Unmarshal([]byte(row[kv["Monster"]]), &Monster); row[kv["Monster"]] != "" && err != nil {
				logger.Panicf("endlesstrials 加载失败 ID :%d", Layers)
			}

			var Settlement []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Settlement"]]), &Settlement); row[kv["Settlement"]] != "" && err != nil {
				logger.Panicf("endlesstrials 加载失败 ID :%d", Layers)
			}

			var Offlineraward []global.Offlineraward
			if err = json.Unmarshal([]byte(row[kv["Offlineraward"]]), &Offlineraward); row[kv["Offlineraward"]] != "" && err != nil {
				logger.Panicf("endlesstrials 加载失败 ID :%d", Layers)
			}

			global.EndlesstrialsCfg[int32(Layers)] = &global.EndlesstrialsConfig{
				Layers:        int32(Layers),
				LevelLimit:    int32(Levellimit),
				TimeLimit:     int32(Timelimit),
				Monster:       Monster,
				Settlement:    Settlement,
				OfflineReward: Offlineraward,
			}
		}
		logger.Infof("endlesstrials 加载完毕,共 [ %d ]  ", len(global.EndlesstrialsCfg))
	} else {
		logger.Panic("endlesstrials.xlsx 加载失败")
	}
}

func LoadEquipMasterCfg() {
	global.EquipMasterCfgMap = make(map[int32]*global.EquipMasterConfig)
	global.EquipMasterCfg = []*global.EquipMasterConfig{}

	if data, err := excel.LoadExcel("equipmaster.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			UUID, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("equipmaster 加载失败 ID :%d", UUID)
			}
			Category, err := strconv.Atoi(row[kv["Category"]])
			if row[kv["Category"]] != "" && err != nil {
				logger.Panicf("equipmaster 加载失败 ID :%d", UUID)
			}

			Limit, err := strconv.Atoi(row[kv["Limit"]])
			if row[kv["Limit"]] != "" && err != nil {
				logger.Panicf("equipmaster 加载失败 ID :%d", UUID)
			}

			var Attribute []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Attribute"]]), &Attribute); row[kv["Attribute"]] != "" && err != nil {
				logger.Panicf("equipmaster 加载失败 ID :%d", UUID)
			}

			global.EquipMasterCfg = append(global.EquipMasterCfg,
				&global.EquipMasterConfig{
					UUID:      int32(UUID),
					Category:  int32(Category),
					Limit:     int32(Limit),
					Attribute: Attribute,
				})

			global.EquipMasterCfgMap[int32(UUID)] = &global.EquipMasterConfig{
				UUID:      int32(UUID),
				Category:  int32(Category),
				Limit:     int32(Limit),
				Attribute: Attribute,
			}

		}
		logger.Infof("equipmaster 加载完毕,共 [ %d ]  ", len(global.EquipMasterCfg))
	} else {
		logger.Panic("equipmaster.xlsx 加载失败")
	}
}

func LoadPrizeWheelCfg() {
	global.PrizeWheelCfg = make(map[int32]*global.PrizeWheelConfig)

	if data, err := excel.LoadExcel("prizewheel.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			UID, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("PrizeWheel 加载失败 UID :%d", UID)
			}
			Round, err := strconv.Atoi(row[kv["Round"]])
			if row[kv["Round"]] != "" && err != nil {
				logger.Panicf("PrizeWheel 加载失败 Round :", Round)
			}

			Consume, err := strconv.Atoi(row[kv["Consume"]])
			if row[kv["Consume"]] != "" && err != nil {
				logger.Panicf("PrizeWheel 加载失败 Consume :%d", Consume)
			}
			Rate, err := strconv.ParseFloat(row[kv["rate"]], 64)
			if row[kv["rate"]] != "" && err != nil {
				logger.Panicf("PrizeWheel 加载失败 rate :", Rate)
			}
			Priority, err := strconv.Atoi(row[kv["Priority"]])
			if row[kv["Priority"]] != "" && err != nil {
				logger.Panicf("PrizeWheel 加载失败 Priority :%d", Priority)
			}
			LevelLimit, err := strconv.Atoi(row[kv["LevelLimit"]])
			if row[kv["LevelLimit"]] != "" && err != nil {
				logger.Panicf("PrizeWheel 加载失败 LevelLimit :%d", LevelLimit)
			}

			global.PrizeWheelCfg[int32(UID)] = &global.PrizeWheelConfig{
				Uid:        int32(UID),
				Round:      int32(Round),
				Consume:    int32(Consume),
				Priority:   int32(Priority),
				LevelLimit: int32(LevelLimit),
				Rate:       Rate,
			}
		}
		logger.Infof("PrizeWheel 加载完毕,共 [ %d ]  ", len(global.PrizeWheelCfg))
	} else {
		logger.Panic("PrizeWheel.xlsx 加载失败")
	}
}

func LoadTreasureBoxCfg() {
	global.TreasureBoxCfg = make(map[int32]*global.TreasureBoxConfig)

	if data, err := excel.LoadExcel("treasurebox.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			UID, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("TreasureBox 加载失败 UID :%d", UID)
			}
			LevelRange, err := global.ConvertRangeToSlice(row[kv["LevelRange"]])
			if row[kv["LevelRange"]] != "" && err != nil {
				logger.Panicf("TreasureBox 加载失败 LevelRange :%d", LevelRange)
			}

			var Consume []global.TreasureBoxConsume
			if err = json.Unmarshal([]byte(row[kv["Consume"]]), &Consume); row[kv["Consume"]] != "" && err != nil {
				logger.Panicf("TreasureBox 加载失败 Consume :%d", Consume)
			}
			var CommonPool []global.TreasureBoxCommonPool
			if err = json.Unmarshal([]byte(row[kv["CommonPool"]]), &CommonPool); row[kv["CommonPool"]] != "" && err != nil {
				logger.Panicf("TreasureBox 加载失败 CommonPool :%d", CommonPool)
			}
			var SpecialPool []global.TreasureBoxCommonPool
			if err = json.Unmarshal([]byte(row[kv["SpecialPool"]]), &SpecialPool); row[kv["SpecialPool"]] != "" && err != nil {
				logger.Panicf("TreasureBox 加载失败 SpecialPool :%d", SpecialPool)
			}
			Node, err := global.ConvertRangeToSlice(row[kv["Node"]])
			if row[kv["Node"]] != "" && err != nil {
				logger.Panicf("TreasureBox 加载失败 Node :%d", Node)
			}
			MaxTime, err := strconv.Atoi(row[kv["MaxTime"]])
			if row[kv["MaxTime"]] != "" && err != nil {
				logger.Panicf("TreasureBox 加载失败 MaxTime :%d", MaxTime)
			}

			global.TreasureBoxCfg[int32(UID)] = &global.TreasureBoxConfig{
				Uid:         int32(UID),
				LevelRange:  LevelRange,
				Consume:     Consume,
				CommonPool:  CommonPool,
				SpecialPool: SpecialPool,
				Node:        Node,
				MaxTime:     int32(MaxTime),
			}
		}
		logger.Infof("TreasureBox 加载完毕,共 [ %d ]  ", len(global.TreasureBoxCfg))
	} else {
		logger.Panic("TreasureBox.xlsx 加载失败")
	}
}

func LoadRandNameCfg() {
	global.RandNameCfg = []*global.RandNameConfig{}
	if data, err := excel.LoadExcel("name.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			ID, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("vip 加载失败 ID :%d", ID)
			}

			FirstName := row[kv["Firstname"]]
			Familyname := row[kv["Familyname"]]
			Lastname := row[kv["Lastname"]]

			cfg := &global.RandNameConfig{ID: int32(ID), Firstname: FirstName, Familyname: Familyname, Lastname: Lastname}
			global.RandNameCfg = append(global.RandNameCfg, cfg)
		}
		logger.Infof("name 加载完毕,共 [ %d ]  ", len(global.RandNameCfg))
	} else {
		logger.Panic("name.xlsx 加载失败")
	}
}

func LoadFateCfg() {
	global.FateCfg = make(map[int32]*global.FateConfig)

	if data, err := excel.LoadExcel("fate.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			FateLevel, err := strconv.Atoi(row[kv["FateLevel"]])
			if row[kv["FateLevel"]] != "" && err != nil {
				logger.Panicf("Fate 加载失败 FateLevel :%d", FateLevel)
			}
			Multiplicity := row[kv["Multiplicity"]]
			if row[kv["Multiplicity"]] == "" {
				logger.Panic("Fate 加载失败 Multiplicity: ", Multiplicity)
			}
			Level := row[kv["Level"]]
			if row[kv["Level"]] == "" {
				logger.Panic("Fate 加载失败 Level: ", Level)
			}
			FateExp, err := strconv.Atoi(row[kv["FateExp"]])
			if row[kv["FateExp"]] != "" && err != nil {
				logger.Panicf("Fate 加载失败 FateExp :%d", FateExp)
			}
			var FateUpValue []global.FateUpValue
			if err = json.Unmarshal([]byte(row[kv["UpValue"]]), &FateUpValue); row[kv["UpValue"]] != "" && err != nil {
				logger.Panicf("Fate 加载失败 UpValue :%d", FateUpValue)
			}
			var AdvancedNeed []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["AdvancedNeed"]]), &AdvancedNeed); row[kv["AdvancedNeed"]] != "" && err != nil {
				logger.Panicf("Fate 加载失败 AdvancedNeed :%d", AdvancedNeed)
			}
			var AdvancedValue []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["AdvancedValue"]]), &AdvancedValue); row[kv["AdvancedValue"]] != "" && err != nil {
				logger.Panicf("Fate 加载失败 AdvancedValue :%d", AdvancedValue)
			}

			global.FateCfg[int32(FateLevel)] = &global.FateConfig{
				FateLevel:     int32(FateLevel),
				Multiplicity:  Multiplicity,
				Level:         Level,
				FateExp:       int32(FateExp),
				UpValue:       FateUpValue,
				AdvancedNeed:  AdvancedNeed,
				AdvancedValue: AdvancedValue,
			}
		}
		logger.Infof("Fate 加载完毕,共 [ %d ]  ", len(global.FateCfg))
	} else {
		logger.Panic("Fate.xlsx 加载失败")
	}
}

func LoadComposeCfg() {
	global.ComposeCfg = make(map[int32]*global.ComposeConfig)

	if data, err := excel.LoadExcel("compose.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			ComposeID, err := strconv.Atoi(row[kv["ComposeID"]])
			if row[kv["ComposeID"]] != "" && err != nil {
				logger.Panicf("compose 加载失败 ComposeID :%d", ComposeID)
			}
			var Raw []int32
			if err = json.Unmarshal([]byte(row[kv["Raw"]]), &Raw); row[kv["Raw"]] != "" && err != nil {
				logger.Panicf("compose 加载失败 Raw :%d", Raw)
			}
			var Produced []global.ProducedConfig
			if err = json.Unmarshal([]byte(row[kv["Produced"]]), &Produced); row[kv["Produced"]] != "" && err != nil {
				logger.Panicf("compose 加载失败 Produced :%d", Produced)
			}
			var NeedGold global.IdNum
			if err = json.Unmarshal([]byte(row[kv["NeedGold"]]), &NeedGold); row[kv["NeedGold"]] != "" && err != nil {
				logger.Panicf("compose 加载失败 NeedGold :%d", NeedGold)
			}
			NeedNum, err := strconv.Atoi(row[kv["NeedNum"]])
			if row[kv["NeedNum"]] != "" && err != nil {
				logger.Panicf("compose 加载失败 NeedNum :%d", NeedNum)
			}
			SuccessRate, err := strconv.Atoi(row[kv["SuccessRate"]])
			if row[kv["SuccessRate"]] != "" && err != nil {
				logger.Panicf("compose 加载失败 SuccessRate :%d", SuccessRate)
			}
			// 失败返还
			var Return []global.KeyValueEx2
			if err = json.Unmarshal([]byte(row[kv["Return"]]), &Return); row[kv["Return"]] != "" && err != nil {
				logger.Panicf("compose 加载失败 Return :%d", Return)
			}

			global.ComposeCfg[int32(ComposeID)] = &global.ComposeConfig{
				ComposeID:   int32(ComposeID),
				Raw:         Raw,
				Produced:    Produced,
				NeedGold:    NeedGold,
				NeedNum:     int32(NeedNum),
				SuccessRate: int32(SuccessRate),
				Return:      Return,
			}
		}
		logger.Infof("compose 加载完毕,共 [ %d ]  ", len(global.ComposeCfg))
	} else {
		logger.Panic("compose.xlsx 加载失败")
	}
}

func LoadItemGroupCfg() {
	global.ItemDropControlCfg = []*global.ItemDropControlConfig{}
	global.ItemDropControlCfgMap = make(map[string]*global.ItemDropControlConfig, 128)
	if data, err := excel.LoadExcel("dropcontrol.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			Uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("dropcontrol 加载失败 %d", Uid)
			}

			var IDs []int32
			if err = json.Unmarshal([]byte(row[kv["ItemGroup"]]), &IDs); row[kv["ItemGroup"]] != "" && err != nil {
				logger.Panicf("dropcontrol 加载失败  %d", Uid)
			}

			ControlTime, err := strconv.Atoi(row[kv["ControlTime"]])
			if row[kv["ControlTime"]] != "" && err != nil {
				logger.Panicf("dropcontrol 加载失败  %d", Uid)
			}

			var ControlLevel global.Interval
			if err = json.Unmarshal([]byte(row[kv["ControlLevel"]]), &ControlLevel); row[kv["ControlLevel"]] != "" && err != nil {
				logger.Panicf("dropcontrol 加载失败  %d", Uid)
			}

			Maxnum, err := strconv.Atoi(row[kv["Maxnum"]])
			if row[kv["Maxnum"]] != "" && err != nil {
				logger.Panicf("dropcontrol 加载失败  %d", Uid)
			}

			var Candropmaps []int32
			if err = json.Unmarshal([]byte(row[kv["Candropmap"]]), &Candropmaps); row[kv["Candropmap"]] != "" && err != nil {
				logger.Panicf("dropcontrol 加载失败  %d", Uid)
			}

			Scope, err := strconv.Atoi(row[kv["Scope"]])
			if row[kv["Scope"]] != "" && err != nil {
				logger.Panicf("dropcontrol 加载失败  %d", Uid)
			}

			cfg := &global.ItemDropControlConfig{
				Uid:          int32(Uid),
				ItemIDs:      IDs,
				ControlTime:  int32(ControlTime),
				ControlLevel: ControlLevel,
				Maxnum:       int32(Maxnum),
				CandropMaps:  Candropmaps,
				Scope:        int32(Scope),
			}

			global.ItemDropControlCfg = append(global.ItemDropControlCfg, cfg)

			for _, itemID := range IDs {
				for _, mapID := range Candropmaps {
					key := fmt.Sprintf("%d:%d", itemID, mapID)
					global.ItemDropControlCfgMap[key] = cfg
				}
			}
		}
		logger.Infof("dropcontrol 加载完毕,共 [ %d ]  ", len(global.ItemDropControlCfg))
	} else {
		logger.Panic("dropcontrol.xlsx 加载失败")
	}
}

func LoadLevelGiftCfg() {
	global.LevelGiftCfg = make(map[int32]*global.LevelGiftConfig)
	if data, err := excel.LoadExcel("levelgift.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			Uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("levelgift 加载失败 %d", Uid)
			}

			Level, err := strconv.Atoi(row[kv["Level"]])
			if row[kv["Level"]] != "" && err != nil {
				logger.Panicf("levelgift 加载失败  %d", Uid)
			}

			Gift := []global.IdNumJobGender{}
			if err = json.Unmarshal([]byte(row[kv["Gift"]]), &Gift); row[kv["Gift"]] != "" && err != nil {
				logger.Panicf("levelgift 加载失败  %d", Uid)
			}

			cfg := &global.LevelGiftConfig{
				Uid:   int32(Uid),
				Level: int32(Level),
				Gift:  []global.IdNumJobGender{},
			}
			cfg.Gift = Gift

			global.LevelGiftCfg[cfg.Level] = cfg
		}
		logger.Infof("levelgift 加载完毕,共 [ %d ]  ", len(global.LevelGiftCfg))
	} else {
		logger.Panic("levelgift.xlsx 加载失败")
	}
}

func LoadPartnerCfg() {
	global.PartnerCfg = make(map[int32]*global.PartnerConfig)
	if data, err := excel.LoadExcel("partner.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			cid, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}
			name := row[kv["Name"]]

			star, err := strconv.Atoi(row[kv["Star"]])
			if row[kv["Star"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			quality, err := strconv.Atoi(row[kv["Quality"]])
			if row[kv["Quality"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			job, err := strconv.Atoi(row[kv["job"]])
			if row[kv["job"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			attackRange, err := strconv.Atoi(row[kv["AttackRange"]])
			if row[kv["AttackRange"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			attackDistance, err := strconv.Atoi(row[kv["AttackDistance"]])
			if row[kv["AttackDistance"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			ai, err := strconv.Atoi(row[kv["AI"]])
			if row[kv["AI"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			itemID, err := strconv.Atoi(row[kv["ItemID"]])
			if row[kv["ItemID"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			var skills []map[int32]int32
			if err = json.Unmarshal([]byte(row[kv["Skills"]]), &skills); row[kv["Skills"]] != "" && err != nil {
				logger.Panicf("partner 加载失败 :%d", cid)
			}

			var bornValue []map[int32]int32
			if err = json.Unmarshal([]byte(row[kv["BornValue"]]), &bornValue); row[kv["BornValue"]] != "" && err != nil {
				logger.Panicf("partner 加载失败 :%d", cid)
			}

			var starConsume []global.KeyValueEx
			if err = json.Unmarshal([]byte(row[kv["StarConsume"]]), &starConsume); row[kv["StarConsume"]] != "" && err != nil {
				logger.Panicf("partner 加载失败 :%d", cid)
			}

			var chipNum []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["ChipNum"]]), &chipNum); row[kv["ChipNum"]] != "" && err != nil {
				logger.Panicf("partner 加载失败 :%d", cid)
			}

			var starAttribute []global.KeyValueEx
			if err = json.Unmarshal([]byte(row[kv["StarAttribute"]]), &starAttribute); row[kv["StarAttribute"]] != "" && err != nil {
				logger.Panicf("partner 加载失败 :%d", cid)
			}
			extropower, err := strconv.Atoi(row[kv["extropower"]])
			if row[kv["extropower"]] != "" && err != nil {
				logger.Panicf("partner 加载失败  %d", cid)
			}

			var talentPoints []map[int32]int32
			if err = json.Unmarshal([]byte(row[kv["talent"]]), &talentPoints); row[kv["talent"]] != "" && err != nil {
				logger.Panicf("partner 加载失败 :%d", cid)
			}

			var talentlist []int32
			if err = json.Unmarshal([]byte(row[kv["talentlist"]]), &talentlist); row[kv["talentlist"]] != "" && err != nil {
				logger.Panicf("partner 加载失败 :%d", cid)
			}

			pskill := []global.KeyValue{}
			for _, m := range skills {
				for key, value := range m {
					// 添加到结果slice中
					pskill = append(pskill, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			pbornValue := []global.KeyValue{}
			for _, m := range bornValue {
				for key, value := range m {
					// 添加到结果slice中
					pbornValue = append(pbornValue, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			pTalentPoints := []global.KeyValue{}
			for _, m := range talentPoints {
				for key, value := range m {
					// 添加到结果slice中
					pTalentPoints = append(pTalentPoints, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			cfg := &global.PartnerConfig{
				Cid:            int32(cid),
				Name:           name,
				Star:           int32(star),
				Quality:        int32(quality),
				Job:            int32(job),
				AttackRange:    int32(attackRange),
				AttackDistance: attackDistance,
				AI:             int32(ai),
				Skills:         pskill,
				StarConsume:    starConsume,
				ChipNum:        chipNum,
				StarAttribute:  starAttribute,
				ItemID:         int32(itemID),
				BornValue:      pbornValue,
				Extropower:     int32(extropower),
				TalentPoints:   pTalentPoints,
				TalentList:     talentlist,
			}

			global.PartnerCfg[cfg.Cid] = cfg
		}
		logger.Infof("PartnerCfg 加载完毕,共 [ %d ]  ", len(global.PartnerCfg))
	} else {
		logger.Panic("PartnerCfg.xlsx 加载失败")
	}
}

// 伙伴天赋
func LoadPartnerTalentCfg() {
	global.PartnerTalentCfg = make(map[int32]*global.PartnerTalentConfig)
	if data, err := excel.LoadExcel("talent.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("partnertalent 加载失败 UID :%d", uid)
			}

			name := row[kv["Name"]]
			if row[kv["Name"]] == "" {
				logger.Panicf("partnertalent 加载失败 Name 为空")
			}

			slot, err := strconv.Atoi(row[kv["Slot"]])
			if row[kv["Slot"]] != "" && err != nil {
				logger.Panicf("partnertalent 加载失败 Slot :%d", slot)
			}

			level, err := strconv.Atoi(row[kv["Level"]])
			if row[kv["Level"]] != "" && err != nil {
				logger.Panicf("partnertalent 加载失败 Level :%d", level)
			}

			var values []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Value"]]), &values); row[kv["Value"]] != "" && err != nil {
				logger.Panicf("partnertalent 加载失败 Value :%d", uid)
			}

			var consume global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Consume"]]), &consume); row[kv["Consume"]] != "" && err != nil {
				logger.Panicf("partnertalent 加载失败 Consume :%d", uid)
			}

			var limit global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Limit"]]), &limit); row[kv["Limit"]] != "" && err != nil {
				logger.Panicf("partnertalent 加载失败 Limit :%d", uid)
			}

			path := row[kv["Path"]]
			icon := row[kv["Icon"]]

			cfg := &global.PartnerTalentConfig{
				UID:     int32(uid),
				Name:    name,
				Slot:    int32(slot),
				Level:   int32(level),
				Value:   values,
				Consume: consume,
				Limit:   limit,
				Path:    path,
				Icon:    icon,
			}

			global.PartnerTalentCfg[cfg.UID] = cfg
		}
		logger.Infof("partnertalent 加载完毕,共 [ %d ]  ", len(global.PartnerTalentCfg))
	} else {
		logger.Panic("partnertalent.xlsx 加载失败")
	}
}

func LoadMaintriggerCfg() {
	global.MaintriggerCfg = make(map[int32]*global.MaintriggerConfig)
	if data, err := excel.LoadExcel("maintrigger.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("maintrigger 加载失败  %d", id)
			}
			var level []global.LevelInterval
			if err = json.Unmarshal([]byte(row[kv["Level"]]), &level); row[kv["Level"]] != "" && err != nil {
				logger.Panicf("maintrigger 加载失败 UUID :%d", level)
			}
			var monsterIds []global.MonsterIds
			if err = json.Unmarshal([]byte(row[kv["MonsterIds"]]), &monsterIds); row[kv["MonsterIds"]] != "" && err != nil {
				logger.Panicf("maintrigger 加载失败 UUID :%d", monsterIds)
			}
			backgroundmapID, err := strconv.Atoi(row[kv["BackgroundmapID"]])
			if row[kv["BackgroundmapID"]] != "" && err != nil {
				logger.Panicf("maintrigger 加载BackgroundmapID失败  %d", backgroundmapID)
			}

			cfg := &global.MaintriggerConfig{
				ID:              int32(id),
				Level:           level,
				MonsterIds:      monsterIds,
				BackgroundmapID: int32(backgroundmapID),
			}

			global.MaintriggerCfg[cfg.ID] = cfg
		}
		logger.Infof("maintrigger 加载完毕,共 [ %d ]  ", len(global.PartnerCfg))
	} else {
		logger.Panic("maintrigger.xlsx 加载失败")
	}
}

func LoadOpenBoxCfg() {
	global.OpenBoxCfg = make(map[int32]*global.OpenBoxConfig)
	if data, err := excel.LoadExcel("openbox.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("openbox 加载失败  %d", id)
			}
			var level []global.LevelInterval
			if err = json.Unmarshal([]byte(row[kv["Level"]]), &level); row[kv["Level"]] != "" && err != nil {
				logger.Panicf("openbox 加载失败 level :%d", level)
			}
			var equiprandom []global.IdRateJob
			if err = json.Unmarshal([]byte(row[kv["Equiprandom"]]), &equiprandom); row[kv["Equiprandom"]] != "" && err != nil {
				logger.Panicf("openbox 加载失败 equiprandom :%d", equiprandom)
			}

			cfg := &global.OpenBoxConfig{
				ID:          int32(id),
				Level:       level,
				Equiprandom: equiprandom,
			}

			global.OpenBoxCfg[cfg.ID] = cfg
		}
		logger.Infof("openbox 加载完毕,共 [ %d ]  ", len(global.OpenBoxCfg))
	} else {
		logger.Panic("openbox.xlsx 加载失败")
	}
}

func LoadBoxLevelCfg() {
	global.BoxLevelCfg = make(map[int32]*global.BoxLevelConfig)
	if data, err := excel.LoadExcel("boxlevel.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败  %d", id)
			}

			level, err := strconv.Atoi(row[kv["Level"]])
			if row[kv["Level"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败 level %d", level)
			}
			rolelevel, err := strconv.Atoi(row[kv["Rolelevel"]])
			if row[kv["Rolelevel"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败  Rolelevel%d", rolelevel)
			}
			battleID, err := strconv.Atoi(row[kv["BattleID"]])
			if row[kv["BattleID"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败  BattleID%d", battleID)
			}
			needtime, err := strconv.Atoi(row[kv["Needtime"]])
			if row[kv["Needtime"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败 Needtime %d", needtime)
			}
			var consume []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Consume"]]), &consume); row[kv["Consume"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败 Consume :%d", consume)
			}
			maxtime, err := strconv.Atoi(row[kv["Maxtime"]])
			if row[kv["Maxtime"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败  Maxtime%d", maxtime)
			}
			var qualityrandom []global.IdRate
			if err = json.Unmarshal([]byte(row[kv["qualityrandom"]]), &qualityrandom); row[kv["qualityrandom"]] != "" && err != nil {
				logger.Panicf("boxlevel 加载失败 qualityrandom :%d", qualityrandom)
			}

			cfg := &global.BoxLevelConfig{
				ID:            int32(id),
				Level:         int32(level),
				Rolelevel:     int32(rolelevel),
				BattleID:      int32(battleID),
				Needtime:      int32(needtime),
				Consume:       consume,
				Maxtime:       int32(maxtime),
				Qualityrandom: qualityrandom,
			}

			global.BoxLevelCfg[cfg.ID] = cfg
		}
		logger.Infof("boxlevel 加载完毕,共 [ %d ]  ", len(global.BoxLevelCfg))
	} else {
		logger.Panic("boxlevel.xlsx 加载失败")
	}
}

func LoadSummonCfg() {
	global.SummonCfg = make(map[int32]*global.SummonConfig)

	if data, err := excel.LoadExcel("summon.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			UID, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 UID :%d", UID)
			}
			DrawModule, err := strconv.Atoi(row[kv["DrawModule"]])
			if row[kv["DrawModule"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 DrawModule :%d", DrawModule)
			}
			Drawtype, err := strconv.Atoi(row[kv["Drawtype"]])
			if row[kv["Drawtype"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 drawtype :%d", Drawtype)
			}
			LevelRange, err := global.ConvertRangeToSlice(row[kv["LevelRange"]])
			if row[kv["LevelRange"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 LevelRange :%d", LevelRange)
			}

			var Consume []global.SummonConsume
			if err = json.Unmarshal([]byte(row[kv["Consume"]]), &Consume); row[kv["Consume"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 Consume :%d", Consume)
			}
			var CommonPool []global.SummonCommonPool
			if err = json.Unmarshal([]byte(row[kv["CommonPool"]]), &CommonPool); row[kv["CommonPool"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 CommonPool :%d", CommonPool)
			}
			var SpecialPool []global.SummonCommonPool
			if err = json.Unmarshal([]byte(row[kv["SpecialPool"]]), &SpecialPool); row[kv["SpecialPool"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 SpecialPool :%d", SpecialPool)
			}
			Node, err := global.ConvertRangeToSlice(row[kv["Node"]])
			if row[kv["Node"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 Node :%d", Node)
			}
			MaxTime, err := strconv.Atoi(row[kv["MaxTime"]])
			if row[kv["MaxTime"]] != "" && err != nil {
				logger.Panicf("Summon 加载失败 MaxTime :%d", MaxTime)
			}

			global.SummonCfg[int32(UID)] = &global.SummonConfig{
				Uid:         int32(UID),
				DrawModule:  int32(DrawModule),
				LevelRange:  LevelRange,
				Drawtype:    int32(Drawtype),
				Consume:     Consume,
				CommonPool:  CommonPool,
				SpecialPool: SpecialPool,
				Node:        Node,
				MaxTime:     int32(MaxTime),
			}
		}
		logger.Infof("Summon 加载完毕,共 [ %d ]  ", len(global.SummonCfg))
	} else {
		logger.Panic("Summon.xlsx 加载失败")
	}
}

func LoadPetCfg() {
	global.PetCfg = make(map[int32]*global.PetConfig)

	if data, err := excel.LoadExcel("pet.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			ID, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 ID :%d", ID)
			}
			Name := row[kv["Name"]]
			PetKind := global.KeyValue{}
			if err = json.Unmarshal([]byte(row[kv["PetKind"]]), &PetKind); row[kv["PetKind"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 PetKind :%d", PetKind)
			}
			Quality, err := strconv.Atoi(row[kv["Quality"]])
			if row[kv["Quality"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 Quality :%d", Quality)
			}
			Mark, err := strconv.Atoi(row[kv["Mark"]])
			if row[kv["Mark"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 Mark :%d", Mark)
			}
			var StarLimit []global.IdValues
			if err = json.Unmarshal([]byte(row[kv["StarLimit"]]), &StarLimit); row[kv["StarLimit"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 StarLimit :%d", StarLimit)
			}
			var StarConsume []global.IdValues
			if err = json.Unmarshal([]byte(row[kv["StarConsume"]]), &StarConsume); row[kv["StarConsume"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 StarConsume :%d", StarConsume)
			}
			var StarAttribute []global.KeyValueEx
			if err = json.Unmarshal([]byte(row[kv["StarAttribute"]]), &StarAttribute); row[kv["StarAttribute"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 StarAttribute :%d", StarAttribute)
			}
			var OfflineAward []global.KeyValuesOffline
			if err = json.Unmarshal([]byte(row[kv["OfflineAward"]]), &OfflineAward); row[kv["OfflineAward"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 OfflineAward :%d", OfflineAward)
			}
			Icon := row[kv["Icon"]]
			Iconname := row[kv["Iconname"]]
			ItemID, err := strconv.Atoi(row[kv["ItemID"]])
			if row[kv["ItemID"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 ItemID :%d", ItemID)
			}

			var handleTalent []global.Talent
			if err = json.Unmarshal([]byte(row[kv["Talent"]]), &handleTalent); row[kv["Talent"]] != "" && err != nil {
				logger.Panicf("Pet 加载失败 Talent :%v", row[kv["Talent"]])
			}

			cfg := &global.PetConfig{
				ID:            int32(ID),
				StarLimit:     StarLimit,
				Name:          Name,
				Mark:          int32(Mark),
				PetKind:       PetKind,
				Quality:       int32(Quality),
				StarConsume:   StarConsume,
				StarAttribute: StarAttribute,
				OfflineAward:  OfflineAward,
				Icon:          Icon,
				Iconname:      Iconname,
				ItemID:        int32(ItemID),
				Talent:        handleTalent,
			}
			// for _, m := range handleTalent {
			// 	for key, value := range m.Value {
			// 		cfg.Talent = append(cfg.Talent,
			// 			global.Talent{
			// 				ID:    m.ID,
			// 				Key:   key,
			// 				Value: value,
			// 			})
			// 	}
			// }
			global.PetCfg[cfg.ID] = cfg
		}
		logger.Infof("Pet 加载完毕,共 [ %d ]  ", len(global.PetCfg))
	} else {
		logger.Panic("pet.xlsx 加载失败")
	}
}

func LoadBattleCfg() {
	global.BattleCfg = make(map[int32]*global.BattleConfig)

	if data, err := excel.LoadExcel("battle.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			UID, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 UID :%d", UID)
			}
			Copytype, err := strconv.Atoi(row[kv["Copytype"]])
			if row[kv["Copytype"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Copytype :%d", Copytype)
			}
			Powerlimit, err := strconv.Atoi(row[kv["Powerlimit"]])
			if row[kv["Powerlimit"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Powerlimit :%d", Powerlimit)
			}
			Levellimit, err := strconv.Atoi(row[kv["Levellimit"]])
			if row[kv["Levellimit"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Levellimit :%d", Levellimit)
			}
			Chapter, err := strconv.Atoi(row[kv["Chapter"]])
			if row[kv["Chapter"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Chapter :%d", Chapter)
			}
			Layers, err := strconv.Atoi(row[kv["Layers"]])
			if row[kv["Layers"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Layers :%d", Layers)
			}
			MapMark, err := strconv.Atoi(row[kv["MapMark"]])
			if row[kv["MapMark"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Layers :%d", Layers)
			}

			var Monster []global.IdNumP
			if err = json.Unmarshal([]byte(row[kv["Monster"]]), &Monster); row[kv["Monster"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Monster :%d", Monster)
			}

			Timelimit, err := strconv.Atoi(row[kv["Timelimit"]])
			if row[kv["Timelimit"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Timelimit :%d", Timelimit)
			}

			var Settlement []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Settlement"]]), &Settlement); row[kv["Settlement"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Settlement :%d", Settlement)
			}

			var Randomaward []global.IdNumRate
			if err = json.Unmarshal([]byte(row[kv["Randomaward"]]), &Randomaward); row[kv["Randomaward"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Randomaward :%d", Randomaward)
			}

			var ChapterAward []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["ChapterAward"]]), &ChapterAward); row[kv["ChapterAward"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 ChapterAward :%d", ChapterAward)
			}

			var StarAward []global.KeyValueEx2
			if err = json.Unmarshal([]byte(row[kv["StarAward"]]), &StarAward); row[kv["StarAward"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 StarAward :%d", StarAward)
			}

			var BattleEffect []global.IdNumRate
			if err = json.Unmarshal([]byte(row[kv["Battleeffect"]]), &BattleEffect); row[kv["Battleeffect"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 BattleEffect :%d", BattleEffect)
			}

			recommendPOW, err := strconv.Atoi(row[kv["recommendPOW"]])
			if row[kv["recommendPOW"]] != "" && err != nil {
				logger.Panicf("battle 加载失败 Layers :%d", Layers)
			}

			global.BattleCfg[int32(UID)] = &global.BattleConfig{
				UID:          int32(UID),
				CopyType:     int32(Copytype),
				PowerLimit:   int32(Powerlimit),
				LevelLimit:   int32(Levellimit),
				Chapter:      int32(Chapter),
				Layers:       int32(Layers),
				Monster:      Monster,
				TimeLimit:    int32(Timelimit),
				Settlement:   Settlement,
				Randomaward:  Randomaward,
				ChapterAward: ChapterAward,
				StarAward:    StarAward,
				MapMark:      int32(MapMark),
				BattleEffect: BattleEffect,
				RecommendPOW: int32(recommendPOW),
			}
		}
		logger.Infof("battle 加载完毕,共 [ %d ]  ", len(global.BattleCfg))
	} else {
		logger.Panic("battle.xlsx 加载失败")
	}
}

func LoadGuildCfg() {
	global.GuildCfg = make(map[int32]*global.GuildConfig)

	if data, err := excel.LoadExcel("guildEXP.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			guildLevel, err := strconv.Atoi(row[kv["GuildLevel"]])
			if row[kv["GuildLevel"]] != "" && err != nil {
				logger.Panicf("guild 加载失败 GuildLevel :%d", guildLevel)
			}

			needEXP, err := strconv.Atoi(row[kv["NeedEXP"]])
			if row[kv["NeedEXP"]] != "" && err != nil {
				logger.Panicf("guild 加载失败 NeedEXP :%d", int32(needEXP))
			}

			var attributeAdd []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["AttributeAdd"]]), &attributeAdd); row[kv["AttributeAdd"]] != "" && err != nil {
				logger.Panicf("guild 加载失败 AttributeAdd :%d", attributeAdd)
			}

			privilegeDes := row[kv["PrivilegeDes"]]

			warehouseAdd, err := strconv.Atoi(row[kv["WarehouseAdd"]])
			if row[kv["WarehouseAdd"]] != "" && err != nil {
				logger.Panicf("guild 加载失败 WarehouseAdd :%d", warehouseAdd)
			}

			memberNum, err := strconv.Atoi(row[kv["MemberNum"]])
			if row[kv["MemberNum"]] != "" && err != nil {
				logger.Panicf("guild 加载失败 MemberNum :%d", memberNum)
			}

			global.GuildCfg[int32(guildLevel)] = &global.GuildConfig{
				GuildLevel:   int32(guildLevel),
				NeedEXP:      int32(needEXP),
				AttributeAdd: attributeAdd,
				PrivilegeDes: privilegeDes,
				WarehouseAdd: int32(warehouseAdd),
				MemberNum:    int32(memberNum),
			}
		}
		logger.Infof("guild 加载完毕,共 [ %d ]  ", len(global.GuildCfg))
	} else {
		logger.Panic("guild.xlsx 加载失败")
	}
}

func LoadGuildStoreCfg() {
	global.GuildStoreCfg = make(map[int32]*global.GuildStoreConfig)
	if data, err := excel.LoadExcel("guildstore.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("guildstore 加载失败 ID :%d", id)
			}

			var goods []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Goods"]]), &goods); row[kv["Goods"]] != "" && err != nil {
				logger.Panicf("guildstore 加载失败 Goods :%d", goods)
			}

			var price []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Price"]]), &price); row[kv["Price"]] != "" && err != nil {
				logger.Panicf("guildstore 加载失败 Price :%d", price)
			}

			guildLimit, err := strconv.Atoi(row[kv["GuildLimit"]])
			if row[kv["GuildLimit"]] != "" && err != nil {
				logger.Panicf("guildstore 加载失败 GuildLimit :%d", guildLimit)
			}

			personalLimit, err := strconv.Atoi(row[kv["PersonalLimit"]])
			if row[kv["PersonalLimit"]] != "" && err != nil {
				logger.Panicf("guildstore 加载失败 PersonalLimit :%d", personalLimit)
			}

			unlockLevel, err := strconv.Atoi(row[kv["UnlockLevel"]])
			if row[kv["UnlockLevel"]] != "" && err != nil {
				logger.Panicf("guildstore 加载失败 UnlockLevel :%d", unlockLevel)
			}

			global.GuildStoreCfg[int32(id)] = &global.GuildStoreConfig{
				ID:            int32(id),
				Goods:         goods,
				Price:         price,
				GuildLimit:    int32(guildLimit),
				PersonalLimit: int32(personalLimit),
				UnlockLevel:   int32(unlockLevel),
			}
		}
		logger.Infof("guildstore 加载完毕,共 [ %d ]  ", len(global.GuildStoreCfg))
	} else {
		logger.Panic("guildstore.xlsx 加载失败")
	}
}

func LoadGuildSkillCfg() {
	global.GuildSkillCfg = make(map[int32]*global.GuildSkillConfig)

	if data, err := excel.LoadExcel("guildskill.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uuid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("guildskill 加载失败 UUID :%d", uuid)
			}

			skillID, err := strconv.Atoi(row[kv["SkillID"]])
			if row[kv["SkillID"]] != "" && err != nil {
				logger.Panicf("guildskill 加载失败 SkillID :%d", skillID)
			}

			skillName := row[kv["SkillName"]]

			skillLevel, err := strconv.Atoi(row[kv["SkillLevel"]])
			if row[kv["SkillLevel"]] != "" && err != nil {
				logger.Panicf("guildskill 加载失败 SkillLevel :%d", skillLevel)
			}

			exp, err := strconv.Atoi(row[kv["EXP"]])
			if row[kv["EXP"]] != "" && err != nil {
				logger.Panicf("guildskill 加载失败 EXP :%d", exp)
			}

			var valueAdd []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["ValueAdd"]]), &valueAdd); row[kv["ValueAdd"]] != "" && err != nil {
				logger.Panicf("guildskill 加载失败 ValueAdd :%v", valueAdd)
			}

			global.GuildSkillCfg[int32(uuid)] = &global.GuildSkillConfig{
				UUID:       int32(uuid),
				SkillID:    int32(skillID),
				SkillName:  skillName,
				SkillLevel: int32(skillLevel),
				EXP:        int32(exp),
				ValueAdd:   valueAdd,
			}
		}
		logger.Infof("guildskill 加载完毕,共 [ %d ]  ", len(global.GuildSkillCfg))
	} else {
		logger.Panic("guildskill.xlsx 加载失败")
	}
}

func LoadLampCfg() {
	global.LampCfg = make(map[int32]*global.LampConfig)
	if data, err := excel.LoadExcel("lamp.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("lamp 加载失败 ID :%d", id)
			}

			receptor, err := strconv.Atoi(row[kv["Receptor"]])
			if row[kv["Receptor"]] != "" && err != nil {
				logger.Panicf("lamp 加载失败 ID :%d", id)
			}

			global.LampCfg[int32(id)] = &global.LampConfig{
				ID:       int32(id),
				Receptor: int32(receptor),
			}
		}
	} else {
		logger.Panic("guildstore.xlsx 加载失败")
	}
}

// 技能释放时间延迟
func ActionsetCfg() {
	global.ActionsetCfg = make(map[int32]*global.ActionsetConfig)

	if data, err := excel.LoadExcel("actionset.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			skillID, err := strconv.Atoi(row[kv["SkillID"]])
			if err != nil {
				logger.Panicf("actionset 加载失败 SkillID :%d", skillID)
			}

			lasttime, err := strconv.Atoi(row[kv["Lasttime"]])
			if err != nil {
				logger.Panicf("actionset 加载失败 Lasttime :%d", skillID)
			}

			job, err := strconv.Atoi(row[kv["Job"]])
			if err != nil {
				logger.Panicf("actionset 加载失败 Job :%d", skillID)
			}

			cfg := &global.ActionsetConfig{
				SkillID:  int32(skillID),
				Lasttime: int32(lasttime),
				Job:      int32(job),
			}

			global.ActionsetCfg[int32(skillID)] = cfg
		}
		logger.Infof("actionset 加载完毕,共 [ %d ]  ", len(global.ActionsetCfg))
	} else {
		logger.Panic("actionset.xlsx 加载失败")
	}
}

// 活动控制表
func LoadActivityControllerCfg() {
	global.ActivityControllerCfg = make(map[int32]*global.ActivityControllerConfig, 10)

	if data, err := excel.LoadExcel("activity_controller.xlsx"); err == nil {

		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uuid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("activity_controller 加载失败 UUID :%d", uuid)
			}

			payID, err := strconv.Atoi(row[kv["PayID"]])
			if row[kv["PayID"]] != "" && err != nil {
				logger.Panicf("activity_controller 加载失败 PayID :%d", payID)
			}

			startTimeType, err := strconv.Atoi(row[kv["StartTimeType"]])
			if row[kv["StartTimeType"]] != "" && err != nil {
				logger.Panicf("activity_controller 加载失败 StartTimeType :%d", startTimeType)
			}

			startTime, err := strconv.Atoi(row[kv["StartTime"]])
			if row[kv["StartTime"]] != "" && err != nil {
				logger.Panicf("activity_controller 加载失败 StartTime :%d", startTime)
			}

			endtimeType, err := strconv.Atoi(row[kv["EndTimeType"]])
			if row[kv["EndTimeType"]] != "" && err != nil {
				logger.Panicf("activity_controller 加载失败 EndTimeType :%d", endtimeType)
			}

			endTime, err := strconv.Atoi(row[kv["EndTime"]])
			if row[kv["EndTime"]] != "" && err != nil {
				logger.Panicf("activity_controller 加载失败 EndTime :%d", endTime)
			}

			cfg := &global.ActivityControllerConfig{
				UUID:          int64(uuid),
				PayId:         int64(payID),
				StartTimeType: int32(startTimeType),
				StartTime:     int64(startTime),
				EndTimeType:   int32(endtimeType),
				EndTime:       int64(endTime),
			}
			global.ActivityControllerCfg[int32(uuid)] = cfg
		}

		logger.Infof("activity_controller 加载完毕,共 [ %d ]  ", len(global.ActivityControllerCfg))
	} else {
		logger.Panic("actionset.xlsx 加载失败")
	}
}

// 活动任务表
func LoadActivityQuestCfg() {
	global.ActivityQuestCfg = make(map[int32]*global.ActivityQuestConfig, 64)

	if data, err := excel.LoadExcel("activity_quest.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uuid, err := strconv.Atoi(row[kv["UUID"]])
			if row[kv["UUID"]] != "" && err != nil {
				logger.Panicf("activity_quest 加载失败 UUID :%d", uuid)
			}

			activityControllerID, err := strconv.Atoi(row[kv["ActivityControllerID"]])
			if row[kv["ActivityControllerID"]] != "" && err != nil {
				logger.Panicf("activity_quest 加载失败 ActivityControllerID :%d", activityControllerID)
			}

			subgroupID, err := strconv.Atoi(row[kv["SubgroupID"]])
			if row[kv["SubgroupID"]] != "" && err != nil {
				logger.Panicf("activity_quest 加载失败 SubgroupID :%d", subgroupID)
			}

			day, err := strconv.Atoi(row[kv["Day"]])
			if row[kv["Day"]] != "" && err != nil {
				logger.Panicf("activity_quest 加载失败 Day :%d", day)
			}

			questType, err := strconv.Atoi(row[kv["Type"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("activity_quest 加载失败 Type :%d", questType)
			}

			limit, err := strconv.Atoi(row[kv["Limit"]])
			if row[kv["Limit"]] != "" && err != nil {
				logger.Panicf("activity_quest 加载失败 Limit :%d", limit)
			}

			var Reward []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Reward"]]), &Reward); row[kv["Reward"]] != "" && err != nil {
				logger.Panicln("activity 加载失败 Reward ", key, Reward, row[kv["Reward"]], err.Error())
			}

			activity, err := strconv.Atoi(row[kv["Activity"]])
			if row[kv["Activity"]] != "" && err != nil {
				logger.Panicf("activity_quest 加载失败 Activity :%d", activity)
			}

			cfg := &global.ActivityQuestConfig{
				UUID:                 int32(uuid),
				ActivityControllerID: int32(activityControllerID),
				SubGroupID:           int32(subgroupID),
				Day:                  int32(day),
				Type:                 int32(questType),
				Limit:                int32(limit),
				// Reward:               int32(reward),
				Activity: int32(activity),
			}

			for _, m := range Reward {
				for key2, value := range m {
					cfg.Reward = append(cfg.Reward, global.KeyValue{
						Key:   int32(key2),
						Value: int32(value),
					})
				}
			}

			global.ActivityQuestCfg[int32(uuid)] = cfg
		}

		logger.Infof("activity_quest 加载完毕,共 [ %d ]  ", len(global.ActivityQuestCfg))
	} else {
		logger.Panic("activity_quest.xlsx 加载失败")
	}
}
func LoadBonusCfg() {
	global.BonusCfg = make(map[int32]*global.BonusConfig)
	if data, err := excel.LoadExcel("bonus.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("bonus 加载失败 ID :%d", id)
			}

			level, err := strconv.Atoi(row[kv["Level"]])
			if row[kv["Level"]] != "" && err != nil {
				logger.Panicf("bonus 加载失败 Level :%d", level)
			}

			var gift []global.IdNumJobGender
			if err = json.Unmarshal([]byte(row[kv["Gift"]]), &gift); row[kv["Gift"]] != "" && err != nil {
				logger.Panicf("bonus 加载失败 Gift :%d", id)
			}

			cfg := &global.BonusConfig{
				ID:    int32(id),
				Level: int32(level),
				Gift:  gift,
			}

			global.BonusCfg[cfg.Level] = cfg
		}
		logger.Infof("bonus 加载完毕,共 [ %d ]  ", len(global.BonusCfg))
	} else {
		logger.Panic("bonus.xlsx 加载失败")
	}
}

// 羁绊表
func LoadBondCfg() {
	global.PartnerBondCfg = []*global.PartnerBondConfig{}

	if data, err := excel.LoadExcel("partnerlink.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("bond 加载失败 ID :%d", id)
			}

			bondID, err := strconv.Atoi(row[kv["BondID"]])
			if row[kv["BondID"]] != "" && err != nil {
				logger.Panicf("bond 加载失败 BondID :%d", bondID)
			}

			bondName := row[kv["BondName"]]

			bondLevel, err := strconv.Atoi(row[kv["BondLevel"]])
			if row[kv["BondLevel"]] != "" && err != nil {
				logger.Panicf("bond 加载失败 BondLevel :%d", bondLevel)
			}

			var partnerIDs []int32
			if err = json.Unmarshal([]byte(row[kv["PartnerIDs"]]), &partnerIDs); row[kv["PartnerIDs"]] != "" && err != nil {
				logger.Panicf("bond 加载失败 PartnerIDs :%d", id)
			}

			star, err := strconv.Atoi(row[kv["Star"]])
			if row[kv["Star"]] != "" && err != nil {
				logger.Panicf("bond 加载失败 Star :%d", star)
			}

			var attribute []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Attribute"]]), &attribute); row[kv["Attribute"]] != "" && err != nil {
				logger.Panicln("activity 加载失败 Attribute ", key, attribute, row[kv["Attribute"]], err.Error())
			}

			cfg := &global.PartnerBondConfig{
				ID:         int32(id),
				BondID:     int32(bondID),
				BondName:   bondName,
				BondLevel:  int32(bondLevel),
				PartnerIDs: partnerIDs,
				Star:       int32(star),
			}
			cfg.Attribute = []global.KeyValue{}

			for _, m := range attribute {
				for k, v := range m {
					cfg.Attribute = append(cfg.Attribute, global.KeyValue{
						Key:   int32(k),
						Value: int32(v),
					})
				}
			}

			global.PartnerBondCfg = append(global.PartnerBondCfg, cfg)
		}
		logger.Infof("bond 加载完毕,共 [ %d ]  ", len(global.PartnerBondCfg))
	} else {
		logger.Panic("bond.xlsx 加载失败")
	}
}

// // 主界面战斗
// MainBattleConfig struct {
// 	ID           int32   // 唯一ID
// 	Name         string  // 名称
// 	Flag         int32   // 标识
// 	Cost         int32   //消耗
// 	FirstReward  []IdNum // 首次奖励
// 	NormalReward []IdNum // 普通奖励
// }

//主界面战斗配置表

func LoadMainBattleCfg() {
	global.MainBattleCfg = make(map[int32]*global.MainBattleConfig)

	if data, err := excel.LoadExcel("mainbattle.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}
			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("mainbattle 加载失败 ID :%d", id)
			}
			name := row[kv["name"]]
			flag, err := strconv.Atoi(row[kv["flag"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("mainbattle 加载失败 ID :%d", id)
			}

			cost, err := strconv.Atoi(row[kv["cost"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("mainbattle 加载失败 ID :%d", id)
			}
			var firstReward []map[int]int
			if err = json.Unmarshal([]byte(row[kv["firstReward"]]), &firstReward); row[kv["firstReward"]] != "" && err != nil {
				logger.Panicln("mainbattle 加载失败 ID ", key, row[kv["firstReward"]], err.Error())
			}

			var drop []map[int]int
			if err = json.Unmarshal([]byte(row[kv["drop"]]), &drop); row[kv["drop"]] != "" && err != nil {
				logger.Panicln("mainbattle 加载失败 ID ", key, row[kv["drop"]], err.Error())
			}
			var monster [][]int32
			if err = json.Unmarshal([]byte(row[kv["monster"]]), &monster); row[kv["monster"]] != "" && err != nil {
				logger.Panicln("mainbattle 加载失败 ID ", key, row[kv["monster"]], err.Error())
			}
			cfg := &global.MainBattleConfig{
				ID:          int32(id),
				Name:        name,
				Flag:        int32(flag),
				Cost:        int32(cost),
				FirstReward: []global.IdNum{},
				Drop:        []global.IdNum{},
				Monster:     [][]int32{},
			}

			for _, m := range firstReward {
				for k, v := range m {
					cfg.FirstReward = append(cfg.FirstReward, global.IdNum{
						ID:  int32(k),
						Num: int32(v),
					})
				}
			}

			for _, m := range drop {
				for k, v := range m {
					cfg.Drop = append(cfg.Drop, global.IdNum{
						ID:  int32(k),
						Num: int32(v),
					})
				}
			}

			for _, v := range monster {
				cfg.Monster = append(cfg.Monster, v)
			}

			global.MainBattleCfg[int32(id)] = cfg
		}
	} else {
		logger.Panic("mainbattle.xlsx 加载失败")
	}
}
func LoadElementCfg() {
	global.ElementCfg = make(map[int32]*global.ElementConfig)
	if data, err := excel.LoadExcel("element.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("element 加载失败 UID :%d", uid)
			}
			part, err := strconv.Atoi(row[kv["Part"]])
			if row[kv["Part"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Part :%d", part)
			}
			partname := row[kv["Partname"]]
			level, err := strconv.Atoi(row[kv["Level"]])
			if row[kv["Level"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Level :%d", level)
			}
			rate, err := strconv.Atoi(row[kv["Rate"]])
			if row[kv["Rate"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Rate :%d", rate)
			}

			var levelvalue []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Levelvalue"]]), &levelvalue); row[kv["Levelvalue"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Levelvalue :%v", row[kv["Levelvalue"]])
			}

			var levelconsume []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Levelconsume"]]), &levelconsume); row[kv["Levelconsume"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Levelconsume :%v", row[kv["Levelconsume"]])
			}

			breaklevel, err := strconv.Atoi(row[kv["Breaklevel"]])
			if row[kv["Breaklevel"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Breaklevel :%d", breaklevel)
			}

			var breakconsume []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Breakconsume"]]), &breakconsume); row[kv["Breakconsume"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Breakconsume :%v", row[kv["Breakconsume"]])
			}

			var breakvalue []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Breakvalue"]]), &breakvalue); row[kv["Breakvalue"]] != "" && err != nil {
				logger.Panicf("element 加载失败 Breakvalue :%v", row[kv["Breakvalue"]])
			}

			cfg := &global.ElementConfig{
				UID:          int32(uid),
				Part:         int32(part),
				Partname:     partname,
				Level:        int32(level),
				Rate:         int32(rate),
				Levelvalue:   levelvalue,
				Levelconsume: levelconsume,
				Breaklevel:   int32(breaklevel),
				Breakconsume: breakconsume,
				Breakvalue:   breakvalue,
			}
			global.ElementCfg[int32(uid)] = cfg
		}
		logger.Infof("element 加载完毕,共 [ %d ]  ", len(global.ElementCfg))
	} else {
		logger.Panic("element.xlsx 加载失败")
	}
}
func LoadElementTrialsCfg() {
	global.ElementTrialsCfg = make(map[int32]*global.ElementTrialsConfig)
	if data, err := excel.LoadExcel("elementtrials.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 UID :%d", uid)
			}
			levelLimit, err := strconv.Atoi(row[kv["Levellimit"]])
			if row[kv["Levellimit"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 Levellimit :%d", levelLimit)
			}
			layers, err := strconv.Atoi(row[kv["Layers"]])
			if row[kv["Layers"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 Layers :%d", layers)
			}

			var monster []global.IdPosition
			if err = json.Unmarshal([]byte(row[kv["Monster"]]), &monster); row[kv["Monster"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 Monster :%v", row[kv["Monster"]])
			}

			timeLimit, err := strconv.Atoi(row[kv["Timelimit"]])
			if row[kv["Timelimit"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 Timelimit :%d", timeLimit)
			}

			var settlement []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Settlement"]]), &settlement); row[kv["Settlement"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 Settlement :%v", row[kv["Settlement"]])
			}

			var firstPassAward []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Firstpassraward"]]), &firstPassAward); row[kv["Firstpassraward"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 Firstpassraward :%v", row[kv["Firstpassraward"]])
			}

			var preview []int32
			if err = json.Unmarshal([]byte(row[kv["Preview"]]), &preview); row[kv["Preview"]] != "" && err != nil {
				logger.Panicf("elementtrials 加载失败 Preview :%v", row[kv["Preview"]])
			}

			cfg := &global.ElementTrialsConfig{
				UID:            int32(uid),
				LevelLimit:     int32(levelLimit),
				Layers:         int32(layers),
				Monster:        monster,
				TimeLimit:      int32(timeLimit),
				Settlement:     settlement,
				FirstPassAward: firstPassAward,
				Preview:        preview,
			}
			global.ElementTrialsCfg[int32(uid)] = cfg
		}
		logger.Infof("elementtrials 加载完毕,共 [ %d ]  ", len(global.ElementTrialsCfg))
	} else {
		logger.Panic("elementtrials.xlsx 加载失败")
	}
}

func LoadTitleCfg() {
	global.TitleCfg = make(map[int32]*global.TitleConfig)

	if data, err := excel.LoadExcel("title.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("title 加载失败 UID :%d", uid)
			}

			name := row[kv["Name"]]
			slot, err := strconv.Atoi(row[kv["Slot"]])
			if row[kv["Slot"]] != "" && err != nil {
				logger.Panicf("title 加载失败 Slot :%d", uid)
			}

			var goods []global.KeyValueEx
			if err = json.Unmarshal([]byte(row[kv["Goods"]]), &goods); row[kv["Goods"]] != "" && err != nil {
				logger.Panicf("title 加载失败 Goods :%d", uid)
			}

			var limit []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Limit"]]), &limit); row[kv["Limit"]] != "" && err != nil {
				logger.Panicf("title 加载失败 Limit :%d", uid)
			}

			var activity []global.KeyValueEx
			if err = json.Unmarshal([]byte(row[kv["Activity"]]), &activity); row[kv["Activity"]] != "" && err != nil {
				logger.Panicf("title 加载失败 Activity :%d", uid)
			}

			var starRate []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["StarRate"]]), &starRate); row[kv["StarRate"]] != "" && err != nil {
				logger.Panicf("title 加载失败 StarRate :%d", uid)
			}

			var starConsume []global.KeyValueEx
			if err = json.Unmarshal([]byte(row[kv["StarConsume"]]), &starConsume); row[kv["StarConsume"]] != "" && err != nil {
				logger.Panicf("title 加载失败 StarConsume :%d", uid)
			}

			period, err := strconv.Atoi(row[kv["Period"]])
			if row[kv["Period"]] != "" && err != nil {
				logger.Panicf("title 加载失败 Period :%d", uid)
			}

			path := row[kv["Path"]]
			icon := row[kv["Icon"]]

			cfg := &global.TitleConfig{
				UID:         int32(uid),
				Name:        name,
				Slot:        int32(slot),
				Goods:       goods,
				Limit:       limit,
				Activity:    activity,
				StarRate:    starRate,
				StarConsume: starConsume,
				Period:      int32(period),
				Path:        path,
				Icon:        icon,
			}

			global.TitleCfg[int32(uid)] = cfg
		}
		logger.Infof("title 加载完毕,共 [ %d ]  ", len(global.TitleCfg))
	} else {
		logger.Panic("title.xlsx 加载失败")
	}
}

func LoadTitleActivationCfg() {
	global.TitleActivationCfg = make(map[int32]*global.TitleActivationConfig)

	if data, err := excel.LoadExcel("titleactivity.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("titleactivity 加载失败 ID :%d", id)
			}

			slot, err := strconv.Atoi(row[kv["Slot"]])
			if row[kv["Slot"]] != "" && err != nil {
				logger.Panicf("titleactivity 加载失败 Slot :%d", id)
			}

			limit := row[kv["Limit"]]

			starLimit, err := strconv.Atoi(row[kv["StarLimit"]])
			if row[kv["StarLimit"]] != "" && err != nil {
				logger.Panicf("titleactivity 加载失败 StarLimit :%d", id)
			}

			var attribute []global.KeyValue
			if err = json.Unmarshal([]byte(row[kv["Attribute"]]), &attribute); row[kv["Attribute"]] != "" && err != nil {
				logger.Panicf("titleactivity 加载失败 Attribute :%d", id)
			}

			cfg := &global.TitleActivationConfig{
				ID:        int32(id),
				Slot:      int32(slot),
				Limit:     limit,
				StarLimit: int32(starLimit),
				Attribute: attribute,
			}

			global.TitleActivationCfg[int32(id)] = cfg
		}
		logger.Infof("titleactivity 加载完毕,共 [ %d ]  ", len(global.TitleActivationCfg))
	} else {
		logger.Panic("titleactivity.xlsx 加载失败")
	}
}

func LoadActiveCodeCfg() {
	global.ActivationCodeCfg = make(map[string]*global.ActivationCodeConfig)

	if data, err := excel.LoadExcel("activecode.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			channelTypeStr := row[kv["ChannelType"]]
			channelType, err := strconv.Atoi(channelTypeStr)
			if channelTypeStr != "" && err != nil {
				logger.Panicf("activecode 加载失败 ChannelType :%s [行:%d 列:%s]", channelTypeStr, key, "ChannelType")
			}

			id := row[kv["ID"]]

			codeTypeStr := row[kv["CodeType"]]
			codeType, err := strconv.Atoi(codeTypeStr)
			if codeTypeStr != "" && err != nil {
				logger.Panicf("activecode 加载失败 CodeType :%s [行:%d 列:%s]", codeTypeStr, key, "CodeType")
			}
			startStr := row[kv["Start"]]
			start, err := time.Parse("2006-1-2", startStr)
			if startStr != "" && err != nil {
				logger.Panicf("activecode 加载失败 Start :%s [行:%d 列:%s]", startStr, key, "Start")
			}
			endStr := row[kv["End"]]
			end, err := time.Parse("2006-1-2", endStr)
			if endStr != "" && err != nil {
				logger.Panicf("activecode 加载失败 End :%s [行:%d 列:%s]", endStr, key, "End")
			}
			rewardStr := row[kv["Reward"]]
			var reward []global.IdNum
			if err = json.Unmarshal([]byte(rewardStr), &reward); rewardStr != "" && err != nil {
				logger.Panicf("activecode 加载失败 Reward :%s [行:%d 列:%s]", rewardStr, key, "Reward")
			}

			cfg := &global.ActivationCodeConfig{
				ChannelType: int32(channelType),
				CodeType:    int32(codeType),
				Start:       start,
				End:         end,
				Reward:      reward,
			}
			global.ActivationCodeCfg[id] = cfg
		}
		logger.Infof("activecode 加载完毕,共 [ %d ]  ", len(global.ActivationCodeCfg))
	} else {
		logger.Panic("activecode.xlsx 加载失败")
	}
}

func LoadKillMonControlCfg() {
	global.KillMonControlCfg = []*global.KillMonControlConfig{}
	if data, err := excel.LoadExcel("killmoncontrol.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			Uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("killmoncontrol 加载失败 %d", Uid)
			}

			var IDs []int32
			if err = json.Unmarshal([]byte(row[kv["MonID"]]), &IDs); row[kv["MonID"]] != "" && err != nil {
				logger.Panicf("killmoncontrol 加载失败  %d", Uid)
			}

			Maxnum, err := strconv.Atoi(row[kv["Maxnum"]])
			if row[kv["Maxnum"]] != "" && err != nil {
				logger.Panicf("killmoncontrol 加载失败  %d", Uid)
			}

			cfg := &global.KillMonControlConfig{
				Uid:        int32(Uid),
				MonsterIDs: IDs,
				Maxnum:     int32(Maxnum),
			}

			global.KillMonControlCfg = append(global.KillMonControlCfg, cfg)
		}
		logger.Infof("killmoncontrol 加载完毕,共 [ %d ]  ", len(global.KillMonControlCfg))
	} else {
		logger.Panic("killmoncontrol.xlsx 加载失败")
	}
}
func LoadLotteryCfg() {
	global.LotteryCfg = make(map[int32]*global.LotteryConfig)
	global.LotteryCfgByType = make(map[int32][]*global.LotteryConfig)
	if data, err := excel.LoadExcel("dropbag.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 ID :%d", id)
			}
			bagId, err := strconv.Atoi(row[kv["bagId"]])
			if row[kv["bagId"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 bagId :%d", bagId)
			}
			name := row[kv["name"]]
			studioname, err := strconv.Atoi(row[kv["studioname"]])
			if row[kv["studioname"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 studioname :%d", studioname)
			}
			ifBag, err := strconv.Atoi(row[kv["ifBag"]])
			if row[kv["ifBag"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 ifBag :%d", ifBag)
			}
			itemId, err := strconv.Atoi(row[kv["itemId"]])
			if row[kv["itemId"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 itemId :%d", itemId)
			}
			itemName := row[kv["itemName"]]
			itemNum, err := strconv.Atoi(row[kv["itemNum"]])
			if row[kv["itemNum"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 itemNum :%d", itemNum)
			}
			probability, err := strconv.Atoi(row[kv["probability"]])
			if row[kv["probability"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 probability :%d id:%d", probability, id)
			}
			typ, err := strconv.Atoi(row[kv["type"]])
			if row[kv["type"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 type :%d", typ)
			}
			timesOut, err := strconv.Atoi(row[kv["timesOut"]])
			if row[kv["timesOut"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 timesOut :%d", timesOut)
			}
			countNum, err := strconv.Atoi(row[kv["countNum"]])
			if row[kv["countNum"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 countNum :%d", countNum)
			}
			show, err := strconv.Atoi(row[kv["show"]])
			if row[kv["show"]] != "" && err != nil {
				logger.Panicf("lottery 加载失败 show :%d", show)
			}

			global.LotteryCfgByType[int32(bagId)] = append(global.LotteryCfgByType[int32(bagId)], &global.LotteryConfig{
				ID:          int32(id),
				BagId:       int32(bagId),
				Name:        name,
				Studioname:  int32(studioname),
				IfBag:       int32(ifBag),
				ItemId:      int32(itemId),
				ItemName:    itemName,
				ItemNum:     int32(itemNum),
				Probability: int32(probability),
				Type:        int32(typ),
				TimesOut:    int32(timesOut),
				CountNum:    int32(countNum),
				Show:        int32(show),
			})
			global.LotteryCfg[int32(id)] = &global.LotteryConfig{
				ID:          int32(id),
				BagId:       int32(bagId),
				Name:        name,
				Studioname:  int32(studioname),
				IfBag:       int32(ifBag),
				ItemId:      int32(itemId),
				ItemName:    itemName,
				ItemNum:     int32(itemNum),
				Probability: int32(probability),
				Type:        int32(typ),
				TimesOut:    int32(timesOut),
				CountNum:    int32(countNum),
				Show:        int32(show),
			}
		}
		logger.Infof("lottery 加载完毕,共 [ %d ]  ", len(global.LotteryCfg))
	} else {
		logger.Panic("lottery.xlsx 加载失败")
	}
}

func LoadAvatarCfg() {
	global.AvatarCfg = make(map[int32]*global.AvatarConfig)
	if data, err := excel.LoadExcel("roleicon.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("roleicon 加载失败 ID :%d", id)
			}
			name := row[kv["Name"]]
			typ, err := strconv.Atoi(row[kv["ItemType"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("roleicon 加载失败 Type :%d", typ)
			}
			var value []global.AttrWithJob
			if err = json.Unmarshal([]byte(row[kv["Value"]]), &value); row[kv["Value"]] != "" && err != nil {
				logger.Panicf("roleicon 加载失败 Value :%d", id)
			}
			period, err := strconv.Atoi(row[kv["Period"]])
			if row[kv["Period"]] != "" && err != nil {
				logger.Panicf("roleicon 加载失败 Period :%d", period)
			}
			cfg := &global.AvatarConfig{
				UID:      int32(id),
				Name:     name,
				ItemType: int32(typ),
				Value:    value,
				Period:   int32(period),
			}
			global.AvatarCfg[int32(id)] = cfg
		}
		logger.Infof("roleicon 加载完毕,共 [ %d ]  ", len(global.AvatarCfg))
	}
}

// LoadRandomCreatV2 加载 randomcreatV2.xlsx 配置文件
func LoadRandomCreatV2() {
	global.RandomCreatCfgV2 = make(map[int32][]*global.RandomCreatV2Config, 16)

	data, err := excel.LoadExcel("randomcreatV2.xlsx")
	if err != nil {
		logger.Panic("randomcreatv2.xlsx 加载失败: ", err)
		return
	}

	if len(data) <= 2 {
		logger.Warn("randomcreatv2.xlsx 数据为空或格式不正确")
		return
	}

	// 解析表头
	header := data[2]
	columnIndex := make(map[string]int)
	for i, colName := range header {
		columnIndex[colName] = i
	}

	// 校验必要列是否存在
	requiredColumns := []string{"EquipSlot", "MaxLevel", "AttributeTypeSet", "RefiningSegmentRandomSet", "RefiningSegmentDirectionalSet"}
	for _, col := range requiredColumns {
		if _, exists := columnIndex[col]; !exists {
			logger.Panicf("randomcreatv2.xlsx 缺少必要列: %s", col)
			return
		}
	}

	// 解析数据行
	for key, row := range data {
		if key <= 2 {
			continue // 跳过标题行和表头行
		}

		equipSlot, err := strconv.Atoi(row[columnIndex["EquipSlot"]])
		if err != nil || row[columnIndex["EquipSlot"]] == "" {
			logger.Panicf("randomcreatv2 加载失败: EquipSlot 无效 (行: %d)", key)
		}

		maxLevel, err := strconv.Atoi(row[columnIndex["MaxLevel"]])
		if err != nil || row[columnIndex["MaxLevel"]] == "" {
			logger.Panicf("randomcreatv2 加载失败: MaxLevel 无效 (行: %d)", key)
		}

		var attributeTypeSet []map[int]int
		if err := json.Unmarshal([]byte(row[columnIndex["AttributeTypeSet"]]), &attributeTypeSet); err != nil && row[columnIndex["KindSet"]] != "" {
			logger.Panicf("randomcreatv2 加载失败: AttributeTypeSet 无效 (行: %d)", key)
		}

		var refiningSegmentRandomSet []map[int]int
		if err := json.Unmarshal([]byte(row[columnIndex["RefiningSegmentRandomSet"]]), &refiningSegmentRandomSet); err != nil && row[columnIndex["RefiningSegmentRandomSet"]] != "" {
			logger.Panicf("randomcreatv2 加载失败: RefiningSegmentRandomSet 无效 (行: %d)", key)
		}

		var refiningSegmentDirectionalSet []map[int]int
		if err := json.Unmarshal([]byte(row[columnIndex["RefiningSegmentDirectionalSet"]]), &refiningSegmentDirectionalSet); err != nil && row[columnIndex["RefiningSegmentDirectionalSet"]] != "" {
			logger.Panicf("randomcreatv2 加载失败: RefiningSegmentDirectionalSet 无效 (行: %d)", key)
		}

		cfg := &global.RandomCreatV2Config{
			EquipSlot: int32(equipSlot),
			MaxLevel:  int32(maxLevel),
		}

		cfg.AttributeTypeSet = []global.KeyValue{}
		for _, m := range attributeTypeSet {
			for k, v := range m {
				cfg.AttributeTypeSet = append(cfg.AttributeTypeSet, global.KeyValue{
					Key:   int32(k),
					Value: int32(v),
				})
			}
		}

		cfg.RefiningSegmentRandomSet = []global.KeyValue{}
		for _, m := range refiningSegmentRandomSet {
			for k, v := range m {
				cfg.RefiningSegmentRandomSet = append(cfg.RefiningSegmentRandomSet, global.KeyValue{
					Key:   int32(k),
					Value: int32(v),
				})
			}
		}

		cfg.RefiningSegmentDirectionalSet = []global.KeyValue{}
		for _, m := range refiningSegmentDirectionalSet {
			for k, v := range m {
				cfg.RefiningSegmentDirectionalSet = append(cfg.RefiningSegmentDirectionalSet, global.KeyValue{
					Key:   int32(k),
					Value: int32(v),
				})
			}
		}

		global.RandomCreatCfgV2[int32(equipSlot)] = append(global.RandomCreatCfgV2[int32(equipSlot)], cfg)
	}

	logger.Infof("randomcreatv2 加载完毕，共 [ %d ] 条记录", len(global.RandomCreatCfgV2))
}

// LoadRandomDistribV2 加载 randomdistribV2.xlsx 配置文件
func LoadRandomDistribV2() {
	global.AttributeSegmentCfgV2 = make(map[string]*global.AttributeSegmentV2Config)
	global.AttributeSegmentCfgV2MaxLevel = []int32{}
	uniqueMaxLevel := make(map[int32]struct{})

	data, err := excel.LoadExcel("randomdistribV2.xlsx")
	if err != nil {
		logger.Panic("randomdistribv2.xlsx 加载失败: ", err)
		return
	}

	if len(data) <= 2 {
		logger.Warn("randomdistribv2.xlsx 数据为空或格式不正确")
		return
	}

	// 解析表头
	header := data[2]
	columnIndex := make(map[string]int)
	for i, colName := range header {
		columnIndex[colName] = i
	}

	// 校验必要列是否存在
	requiredColumns := []string{"AttributeType", "MaxLevel", "RefingSegment", "ValueMin", "ValueMax", "QualityColor"}
	for _, col := range requiredColumns {
		if _, exists := columnIndex[col]; !exists {
			logger.Panicf("randomdistribv2.xlsx 缺少必要列: %s", col)
			return
		}
	}

	// 解析数据行
	for key, row := range data {
		if key <= 2 {
			continue // 跳过标题行和表头行
		}
		attributeType, err := strconv.Atoi(row[columnIndex["AttributeType"]])
		if err != nil || row[columnIndex["AttributeType"]] == "" {
			logger.Panicf("randomdistribv2 加载失败: AttributeType 无效 (行: %d)", key)
		}
		maxLevel, err := strconv.Atoi(row[columnIndex["MaxLevel"]])
		if err != nil || row[columnIndex["MaxLevel"]] == "" {
			logger.Panicf("randomdistribv2 加载失败: MaxLevel 无效 (行: %d)", key)
		}
		refingSegment, err := strconv.Atoi(row[columnIndex["RefingSegment"]])
		if err != nil || row[columnIndex["RefingSegment"]] == "" {
			logger.Panicf("randomdistribv2 加载失败: RefingSegment 无效 (行: %d)", key)
		}
		valueMin, err := strconv.ParseInt(row[columnIndex["ValueMin"]], 10, 64)
		if err != nil || row[columnIndex["ValueMin"]] == "" {
			logger.Panicf("randomdistribv2 加载失败: ValueMin 无效 (行: %d)", key)
		}
		valueMax, err := strconv.ParseInt(row[columnIndex["ValueMax"]], 10, 64)
		if err != nil || row[columnIndex["ValueMax"]] == "" {
			logger.Panicf("randomdistribv2 加载失败: ValueMax 无效 (行: %d)", key)
		}
		qualityColor, err := strconv.ParseUint(row[columnIndex["QualityColor"]], 10, 32)
		if err != nil || row[columnIndex["QualityColor"]] == "" {
			logger.Panicf("randomdistribv2 加载失败: QualityColor 无效 (行: %d)", key)
		}

		cfg := &global.AttributeSegmentV2Config{
			AttributeType: int32(attributeType),
			MaxLevel:      int32(maxLevel),
			RefingSegment: int32(refingSegment),
			ValueMin:      valueMin,
			ValueMax:      valueMax,
			QualityColor:  uint32(qualityColor),
		}

		// 生成唯一键
		keyStr := fmt.Sprintf("%d_%d_%d", cfg.AttributeType, cfg.MaxLevel, cfg.RefingSegment)
		if _, exists := global.AttributeSegmentCfgV2[keyStr]; exists {
			logger.Panicf("randomdistribv2 加载失败: 重复的配置 (AttributeType: %d, MaxLevel: %d, RefingSegment: %d)", cfg.AttributeType, cfg.MaxLevel, cfg.RefingSegment)
		}
		global.AttributeSegmentCfgV2[keyStr] = cfg

		// 收集唯一的 MaxLevel
		if _, exists := uniqueMaxLevel[cfg.MaxLevel]; !exists {
			uniqueMaxLevel[cfg.MaxLevel] = struct{}{}
			global.AttributeSegmentCfgV2MaxLevel = append(global.AttributeSegmentCfgV2MaxLevel, cfg.MaxLevel)
		}
	}

	// 对 MaxLevel 列表进行排序，方便后续使用
	sort.Slice(global.AttributeSegmentCfgV2MaxLevel, func(i, j int) bool {
		return global.AttributeSegmentCfgV2MaxLevel[i] < global.AttributeSegmentCfgV2MaxLevel[j]
	})

	logger.Infof("randomdistribv2 加载完毕，共 [ %d ] 条记录", len(global.AttributeSegmentCfgV2))
}

func LoadFateCoinCfg() {
	global.FateCoinCfg = make(map[int32]*global.FateCoinConfig)
	if data, err := excel.LoadExcel("fatecoin.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["ID"]])
			if row[kv["ID"]] != "" && err != nil {
				logger.Panicf("fatecoin 加载失败 ID :%d", id)
			}
			typ, err := strconv.Atoi(row[kv["Type"]])
			if row[kv["Type"]] != "" && err != nil {
				logger.Panicf("fatecoin 加载失败 Type :%d", typ)
			}
			buffID, err := strconv.Atoi(row[kv["BuffID"]])
			if row[kv["BuffID"]] != "" && err != nil {
				logger.Panicf("roleicon 加载失败 buffID :%d", buffID)
			}
			effectNum, err := strconv.Atoi(row[kv["Effectnum"]])
			if row[kv["Effectnum"]] != "" && err != nil {
				logger.Panicf("effectNum 加载失败 effectNum :%d", effectNum)
			}
			buffTime, err := strconv.Atoi(row[kv["Bufftime"]])
			if row[kv["Bufftime"]] != "" && err != nil {
				logger.Panicf("roleicon 加载失败 bufftime :%d", buffTime)
			}
			var skillBuff []map[int]int
			if err = json.Unmarshal([]byte(row[kv["Skillbuff"]]), &skillBuff); row[kv["Skillbuff"]] != "" && err != nil {
				logger.Panicf("roleicon 加载失败 skillBuff :%v", skillBuff)
			}
			cfg := &global.FateCoinConfig{
				ID:         int32(id),
				EffectType: pb.FateCoinEffect(typ),
				BuffID:     int32(buffID),
				EffectNum:  int32(effectNum),
				BuffTime:   int32(buffTime),
			}
			for _, m := range skillBuff {
				for key, value := range m {
					cfg.SkillBuff = append(cfg.SkillBuff, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}
			global.FateCoinCfg[int32(id)] = cfg
		}
		logger.Infof("BuffID 加载完毕,共 [ %d ]  ", len(global.AvatarCfg))
	}
}

// 神兵进阶表

func LoadArtifactAdvanceCfg() {
	global.ArtifactAdvanceCfg = make(map[int32]*global.ArtifactAdvanceConfig)
	if data, err := excel.LoadExcel("epicweaponadvance.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 UID :%d", uid)
			}

			weaponID, err := strconv.Atoi(row[kv["WeaponID"]])
			if row[kv["WeaponID"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 WeaponID :%d", weaponID)
			}

			rank, err := strconv.Atoi(row[kv["Rank"]])
			if row[kv["Rank"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 Rank :%d", rank)
			}

			weaponName := row[kv["Wwaponname"]]
			if row[kv["Wwaponname"]] == "" {
				logger.Panic("epicweaponadvance 加载失败 Wwaponname 为空")
			}

			skillName := row[kv["Skillname"]]
			if row[kv["Skillname"]] == "" {
				logger.Panic("epicweaponadvance 加载失败 Skillname 为空")
			}

			openLimit, err := strconv.Atoi(row[kv["OpenLimit"]])
			if row[kv["OpenLimit"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 OpenLimit :%d", openLimit)
			}

			var advanceLimit []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["AdvanceLimit"]]), &advanceLimit); row[kv["AdvanceLimit"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 AdvanceLimit :%d", uid)
			}

			var consume []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["Consume"]]), &consume); row[kv["Consume"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 Consume :%d", uid)
			}

			var advanceValue []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["AdvanceValue"]]), &advanceValue); row[kv["AdvanceValue"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 AdvanceValue :%d", uid)
			}

			skillID, err := strconv.Atoi(row[kv["SkillID"]])
			if row[kv["SkillID"]] != "" && err != nil {
				logger.Panicf("epicweaponadvance 加载失败 SkillID :%d", skillID)
			}

			cfg := &global.ArtifactAdvanceConfig{
				UID:          int32(uid),
				WeaponID:     int32(weaponID),
				Rank:         int32(rank),
				WeaponName:   weaponName,
				SkillName:    skillName,
				OpenLimit:    int32(openLimit),
				AdvanceLimit: advanceLimit,
				Consume:      consume,
				AdvanceValue: advanceValue,
				SkillID:      int32(skillID),
			}

			global.ArtifactAdvanceCfg[int32(uid)] = cfg
		}
		logger.Infof("epicweaponadvance 加载完毕,共 [ %d ]  ", len(global.ArtifactAdvanceCfg))
	} else {
		logger.Panic("epicweaponadvance.xlsx 加载失败")
	}
}

// 神兵强化表
func LoadArtifactStrengthCfg() {
	global.ArtifactStrengthCfg = make(map[int32]*global.ArtifactStrengthConfig)
	if data, err := excel.LoadExcel("epicweaponstrength.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			uid, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("epicweaponstrength 加载失败 UID :%d", uid)
			}

			weaponID, err := strconv.Atoi(row[kv["WeaponID"]])
			if row[kv["WeaponID"]] != "" && err != nil {
				logger.Panicf("epicweaponstrength 加载失败 WeaponID :%d", weaponID)
			}

			weaponName := row[kv["WeaponName"]]
			if row[kv["WeaponName"]] == "" {
				logger.Panic("epicweaponstrength 加载失败 WeaponName 为空")
			}

			openLimit, err := strconv.Atoi(row[kv["OpenLimit"]])
			if row[kv["OpenLimit"]] != "" && err != nil {
				logger.Panicf("epicweaponstrength 加载失败 OpenLimit :%d", openLimit)
			}

			levelLimit, err := strconv.Atoi(row[kv["Levellimit"]])
			if row[kv["Levellimit"]] != "" && err != nil {
				logger.Panicf("epicweaponstrength 加载失败 Levellimit :%d", levelLimit)
			}

			levelExp, err := strconv.Atoi(row[kv["LevelEXP"]])
			if row[kv["LevelEXP"]] != "" && err != nil {
				logger.Panicf("epicweaponstrength 加载失败 LevelEXP :%d", levelExp)
			}

			var levelValue []global.IdNum
			if err = json.Unmarshal([]byte(row[kv["LevelValue"]]), &levelValue); row[kv["LevelValue"]] != "" && err != nil {
				logger.Panicf("epicweaponstrength 加载失败 LevelValue :%d", uid)
			}

			cfg := &global.ArtifactStrengthConfig{
				UID:        int32(uid),
				WeaponID:   int32(weaponID),
				WeaponName: weaponName,
				OpenLimit:  int32(openLimit),
				LevelLimit: int32(levelLimit),
				LevelExp:   int64(levelExp),
				LevelValue: levelValue,
			}

			global.ArtifactStrengthCfg[int32(uid)] = cfg
		}
		logger.Infof("epicweaponstrength 加载完毕,共 [ %d ]  ", len(global.ArtifactStrengthCfg))
	} else {
		logger.Panic("epicweaponstrength.xlsx 加载失败")
	}
}

// LoadPushingWallMonsterCfg 加载 推塔怪物 表
func LoadPushingWallMonsterCfg() {
	global.PushingWallMonsterCfg = make(map[int32]*global.PushingWallMonsterConfig)
	if data, err := excel.LoadExcel("pushingwallmonster.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			// 去掉标题行
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["id"]])
			if row[kv["id"]] != "" && err != nil {
				logger.Panicf("pushingwallmonster 加载失败 ID :%d", id)
			}

			level, err := strconv.Atoi(row[kv["level"]])
			if row[kv["level"]] != "" && err != nil {
				logger.Panicf("pushingwallmonster 加载失败 Level :%d", level)
			}

			var item []map[int]int
			if err = json.Unmarshal([]byte(row[kv["item"]]), &item); row[kv["item"]] != "" && err != nil {
				logger.Panicf("pushingwallmonster 加载失败 Item :%d", id)
			}

			cfg := &global.PushingWallMonsterConfig{
				ID:    int32(id),
				Level: int32(level),
				Item:  []global.KeyValue{},
			}

			for _, m := range item {
				for key, value := range m {
					cfg.Item = append(cfg.Item, global.KeyValue{
						Key:   int32(key),
						Value: int32(value),
					})
				}
			}

			global.PushingWallMonsterCfg[int32(id)] = cfg
		}
		logger.Infof("pushingwallmonster 加载完毕,共 [ %d ] ", len(global.PushingWallMonsterCfg))
	} else {
		logger.Panic("pushingwallmonster.xlsx 加载失败")
	}
}

func LoadOfficialCfg() {
	global.OfficialCfg = make(map[int32]*global.OfficialConfig)
	global.OfficalListCfg = make([]*global.OfficialConfig, 0, 32)
	if data, err := excel.LoadExcel("official.xlsx"); err == nil {
		kv := make(map[string]int)
		for key, row := range data {
			if key <= 1 {
				continue
			}
			if key == 2 {
				for i, v := range row {
					kv[v] = i
				}
				continue
			}

			id, err := strconv.Atoi(row[kv["UID"]])
			if row[kv["UID"]] != "" && err != nil {
				logger.Panicf("official 加载失败 UID :%d", id)
			}
			name := row[kv["Name"]]
			exp, err := strconv.ParseInt(row[kv["EXP"]], 10, 64)
			if row[kv["EXP"]] != "" && err != nil {
				logger.Panicf("official 加载失败 EXP :%d", id)
			}
			quality, err := strconv.Atoi(row[kv["Quality"]])
			if row[kv["Quality"]] != "" && err != nil {
				logger.Panicf("official 加载失败 Quality :%d", id)
			}

			cfg := &global.OfficialConfig{
				UID:     int32(id),
				Name:    name,
				EXP:     exp,
				Quality: int32(quality),
			}
			global.OfficialCfg[int32(id)] = cfg
			global.OfficalListCfg = append(global.OfficalListCfg, cfg)
		}
		logger.Infof("official 加载完毕,共 [ %d ]  ", len(global.OfficialCfg))
	} else {
		logger.Panic("official.xlsx 加载失败")
	}
}
