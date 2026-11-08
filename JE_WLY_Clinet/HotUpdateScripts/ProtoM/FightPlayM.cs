using JEngine.Helper;
using JEngine.Singleton;
using MetaFramework.UI;
using pb;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 玩法信息
    /// </summary>
    public class FightPlayM : SingletonTemplate<FightPlayM>
    {
        //玩法信息
        private InstanceInfoResponse _instanceInfo;

        private Dictionary<int, int> ReclicLevelStarDic = new Dictionary<int, int>();
        private Dictionary<int, int> ReclicChapterStarGetDic = new Dictionary<int, int>();



        public void SetInfo(InstanceInfoResponse data)
        {
            _instanceInfo = data;
            ReclicLevelStarDic.Clear();
            for (int i = 0; i < data.relic_map.Evaluation.Count; i++)
            {
                ReclicLevelStarDic.Add(data.relic_map.Evaluation[i].K, data.relic_map.Evaluation[i].V);
            }

            ReclicChapterStarGetDic.Clear();
            for (int i = 0; i < data.relic_map.ChapterReward.Count; i++)
            {
                ReclicLevelStarDic.Add(data.relic_map.ChapterReward[i].K, data.relic_map.ChapterReward[i].V);
            }
        }

        public InstanceInfoResponse GetInfo()
        {
            return _instanceInfo;
        }

        public int GetTrainLayers()
        {
            var battid = GetInfo().trial_tower_map.BattleId;
            if (battid == 0)
                return 0;
            return CFG_Battle.GetById(battid).Layers;
        }

        //获取遗迹关卡评价
        public int GetRelicLevelStar(int level)
        {
            if (!ReclicLevelStarDic.ContainsKey(level))
                return 0;
            return ReclicLevelStarDic[level];
        }

        //获取遗迹章节已经领取的章节
        public int GetRelicChapterStarGet(int chapter)
        {
            if (!ReclicLevelStarDic.ContainsKey(chapter))
                return 0;
            return ReclicLevelStarDic[chapter];
        }

        //战役点击红点
        public int GetChallenge_FightRedNum()
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play, false))
                return 0;
            return RedPointManager.Instance.GetRedNum(RedPointConst.Challenge_Fight_Click.ToString());
        }

        //战役领取红点
        public int GetChallenge_FightAwardRedNum()
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play, false))
                return 0;

            for (int i = 0; i < CFG_Battle.AllFightChapterAwardList.Count; i++)
            {
                if (GetInfo().battle_map.RewardStage < CFG_Battle.AllFightChapterAwardList[i].UID && CFG_Battle.AllFightChapterAwardList[i].UID <= GetInfo().battle_map.BattleId)
                    return 1;
            }
            return 0;
        }

        //首领副本红点
        public int GetChallenge_BossRedNum()
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss, false))
                return 0;
            if (GetInfo().boss_map.ChallengeTime > 0)
                return 1;
            return 0;
        }

        //试炼塔副本红点
        public int GetChallenge_TowerRedNum()
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower, false))
                return 0;
            return RedPointManager.Instance.GetRedNum(RedPointConst.Challenge_Tower_Click.ToString());
        }

        //遗迹副本红点
        public int GetChallenge_ReliceRedNum()
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_Relice, false))
                return 0;
            if (GetInfo().relic_map.ChallengeTime > 0)
                return 1;
            return 0;
        }

        // 挂机升级红点
        public int GetChallenge_OnHook()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.DefendIsland, false) || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.DefendIsland, false))
                return 0;
            if (GetInfo().afk.DailyTime > 0)
                return 1;
            return 0;
        }

        public int GetFightEndUseTime(int battleId)
        {
            List<KeyValue> fightEnd = _instanceInfo.to_end.Evaluation;
            var cfg_battle = CFG_Battle.GetById(battleId);
            if (cfg_battle.Layers > fightEnd.Count)
                return 0;
            return _instanceInfo.to_end.Evaluation[cfg_battle.Layers - 1].K;
        }

        public bool IsFightEndNewRecord(int battleId, long useTime)
        {
            int time = GetFightEndUseTime(battleId);
            return time == 0 || useTime < time;
        }


        /// <summary>
        /// 战场副本列表
        /// </summary>
        private List<pb.WarRoomListResponse.RoomInfo> WarFightRoomList = new List<pb.WarRoomListResponse.RoomInfo>();

        /// <summary>
        /// 设置战场副本列表
        /// </summary>
        /// <param name="roomInfos"></param>
        public void SetWarFightRoomList(List<pb.WarRoomListResponse.RoomInfo> roomInfos)
        {
            WarFightRoomList = roomInfos;
        }

        public pb.WarRoomListResponse.RoomInfo GetOneWarFightRoomInfo(int index)
        {
            return WarFightRoomList.Find(x => x.RoomCfgId == index);
        }


        public bool CheakWarFightOpenTime()
        {
            var xml = CFG_XMLConfig.xmlConfig.WarFight;
            for (int i = 0; i < xml.openDays.Length; i++)
            {
                if (GameUtil.Instance.CheakCurTimeWeek(xml.openDays[i]))
                {
                    return GameUtil.Instance.CheakCurTime(xml.startTimes[0], xml.endTimes[0]) || GameUtil.Instance.CheakCurTime(xml.startTimes[1], xml.endTimes[1]);
                }
            }
            return false;
        }
    }

}

