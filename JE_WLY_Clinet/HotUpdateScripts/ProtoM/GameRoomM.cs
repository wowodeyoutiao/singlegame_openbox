using HotUpdateScripts.Game;

using JEngine.Singleton;

using pb;

using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 游戏房间内信息
    /// </summary>
    public class GameRoomM : SingletonTemplate<GameRoomM>
    {
        public enum GameState
        {
            No,
            Ready,
            Play,
            Die,
            Over,
            PVP_Over,
            Train_Over,
        }

        public GameState _GameState { get; set; }

        public TableInfo RoomTableInfo { get; set; }

        //副本时间
        public long GameFightTime { get; set; }

        //副本时间增加
        public bool TimeAutoBuy { get; set; }

        public CFGItem_Room CFGRoom { get; set; }

        private bool _attact = false;

        //首先攻击玩家判定
        public bool AttactToPlayer { get { return _attact; } set { _attact = value; } }

        public bool ShortAttactToPlayer;

        //怪物数量
        public int MonsterAllCount { get; set; }

        //自动攻击
        public bool AutoPlay { get; set; }

        public bool AutoPlayLock { get; set; }

        //游戏结算 每日必做
        public GameResultDailyMapPush gameResultDailyMap { get; set; }

        //Boos刷新时间
        public long MosterBossRefreshTime { get; set; }


        /// <summary>
        /// 工会副本奖励
        /// </summary>
        public GuildCopyRewardPush guildCopyReward { get; set; }

        //采集副本 不用攻击
        public bool GetCanNoAttackGame()
        {
            return (CFGRoom.Roomtype == CFG_Room.RoomType.SXS || CFGRoom.Roomtype == CFG_Room.RoomType.HS);
        }

        /// <summary>
        /// 玩家简明信息
        /// </summary>
        private Dictionary<long, PlayerInfo> PInfoDic = new Dictionary<long, PlayerInfo>();

        /// <summary>
        /// 伙伴 分身 实例信息汇总
        /// </summary>
        private Dictionary<long, ObjectInfo> ObjDic = new Dictionary<long, ObjectInfo>();

        //正在被开启宝箱集合
        public Dictionary<long, string> CollectMonsterDic = new Dictionary<long, string>(); //采集怪物集合


        /// <summary>
        /// 战斗掉落
        /// <summary>  

        private List<ItemInfo> DropItemInfos = new List<ItemInfo>();

        public void Add_DropItemInfo(List<ItemInfo> _itemInfo)
        {
            for (int i = 0; i < _itemInfo.Count; i++)
            {
                bool isHave = false;

                for (int j = 0; j < DropItemInfos.Count; j++)
                {
                    if (CFG_Item.GetById(_itemInfo[i].CfgId).Type == (int)CFG_Item.ItemType.Equip)
                        break;

                    if (_itemInfo[i].CfgId == DropItemInfos[j].CfgId)
                    {
                        isHave = true;
                        DropItemInfos[j].Num += _itemInfo[i].Num;
                        break;
                    }
                }
                if (!isHave)
                    DropItemInfos.Add(_itemInfo[i]);
            }
        }

        public List<ItemInfo> GetDropItemInfo()
        {
            return DropItemInfos;
        }

        /// <summary>
        /// 添加玩家信息
        /// </summary>
        /// <param name="pInfo"></param>
        public void Add_PlayerInfo(List<PlayerInfo> pInfo)
        {
            PInfoDic.Clear();
            for (int i = 0; i < pInfo.Count; i++)
            {
                PInfoDic[pInfo[i].Uid] = pInfo[i];
                if (pInfo[i].PartnerList != null)
                {
                    for (int j = 0; j < pInfo[i].PartnerList.Count; j++)
                    {
                        Add_OneObjInfo(pInfo[i].PartnerList[j]);
                    }
                }

            }

        }
        public HashSet<long> DeathHashSet = new HashSet<long>();

        public void ClearObjInfo()
        {
            PInfoDic.Clear();
            MosterBossRefreshTime = 0;
            DeathHashSet.Clear();
            ObjDic.Clear();
            DropItemInfos.Clear();
            CollectMonsterDic.Clear();
        }

        public void Add_OneObjInfo(ObjectInfo objectInfo)
        {
            if (!ObjDic.ContainsKey(objectInfo.Uid))
                ObjDic.Add(objectInfo.Uid, objectInfo);
            else
                ObjDic[objectInfo.Uid] = objectInfo;
        }

        public void Remove_OneObjInfo(ObjectInfo objectInfo)
        {
            if (ObjDic.ContainsKey(objectInfo.Uid))
                ObjDic.Remove(objectInfo.Uid);

        }

        /// <summary>
        /// 获取玩家的GUIDID
        /// </summary>
        /// <param name="PId"></param>
        /// <param name="isPartnher"></param>
        /// <returns></returns>
        public long GetGuiidIdByPid(long PId, bool isPartnher)
        {
            if (isPartnher)
            {
                if (PInfoDic.ContainsKey(PId))
                {
                    return PInfoDic[PId].GuildId;
                }
            }
            else
            {
                if (ObjDic.ContainsKey(PId))
                {
                    if (PInfoDic.ContainsKey(ObjDic[PId].OwnerUid))
                    {
                        return PInfoDic[ObjDic[PId].OwnerUid].GuildId;
                    }
                }
            }
            return 0;
        }

        /// <summary>
        /// 获取玩家的GUID名字
        /// </summary>
        /// <param name="PId"></param>
        /// <param name="isPartnher"></param>
        /// <returns></returns>
        public string GetGuiidNameByPid(long PId, bool isPartnher)
        {
            if (isPartnher)
            {
                if (PInfoDic.ContainsKey(PId))
                {
                    return PInfoDic[PId].GuildName;
                }
            }
            else
            {
                if (ObjDic.ContainsKey(PId))
                {
                    if (PInfoDic.ContainsKey(ObjDic[PId].OwnerUid))
                    {
                        return PInfoDic[ObjDic[PId].OwnerUid].GuildName;
                    }
                }
            }
            return string.Empty;
        }

        public ObjectInfo GetOneObjInfoByID(long uid)
        {
            if (ObjDic.ContainsKey(uid))
                return ObjDic[uid];
            return null;
        }

        public Dictionary<long, ObjectInfo> getObjInfo()
        {
            return ObjDic;
        }

        //添加一个BUffer
        public void Add_PlayerBuffer(long uid, BuffInfo buffInfo)
        {
            var player = getPlayerInfoByID(uid);
            if (player != null)
            {
                if (player.BuffList == null)
                {
                    player.BuffList = new List<BuffInfo>();
                }
                player.BuffList.Add(buffInfo);
            }
            else
            {
                var partner = GetOneObjInfoByID(uid);
                if (partner != null)
                {
                    if (partner.BuffList == null)
                        partner.BuffList = new List<BuffInfo>();

                    partner.BuffList.Add(buffInfo);
                }
            }
        }

        //检测是否有这个Buffer
        public bool CheakHaveBuffer(long pid, int buffid)
        {
            var partner = GetOneObjInfoByID(pid);
            if (partner != null && partner.BuffList != null)
            {
                partner.BuffList.Find(
                delegate (BuffInfo buffer)
                {
                    return buffer.BuffId.Equals(buffid);
                });
            }
            return false;
        }

        //移除一个BUffer
        public void Remove_PlayerBuffer(long uid, int buffId)
        {
            var player = getPlayerInfoByID(uid);
            if (player != null)
            {
                if (player.BuffList != null)
                {
                    for (int i = 0; i < player.BuffList.Count; i++)
                    {
                        if (player.BuffList[i].BuffId == buffId)
                        {
                            player.BuffList.Remove(player.BuffList[i]);
                            return;
                        }
                    }
                }
            }
            else
            {
                var partner = GetOneObjInfoByID(uid);
                if (partner != null)
                {
                    if (partner.BuffList != null)
                    {
                        for (int i = 0; i < partner.BuffList.Count; i++)
                        {
                            if (partner.BuffList[i].BuffId == buffId)
                            {
                                partner.BuffList.Remove(partner.BuffList[i]);
                                return;
                            }
                        }
                    }
                }
            }
        }

        public List<BuffInfo> GetBuffInfosByUid(long uid)
        {
            var player = getPlayerInfoByID(uid);
            if (player != null)
            {
                if (player.BuffList != null)
                {
                    return player.BuffList;
                }
            }
            else
            {
                var partner = GetOneObjInfoByID(uid);
                if (partner != null)
                {
                    if (partner.BuffList != null)
                    {
                        return partner.BuffList;
                    }
                }
            }
            return null;
        }
        public void Remove_PlayerInfo(long uID)
        {
            if (PInfoDic.ContainsKey(uID))
            {
                PInfoDic.Remove(uID);
            }
        }


        public PlayerInfo getPlayerInfoByID(long uID)
        {
            if (PInfoDic.ContainsKey(uID))
            {
                return PInfoDic[uID];
            }

            return null;
        }

        public bool CheakIsMainRole(long uid)
        {
            if (GetOneObjInfoByID(uid) != null)
            {
                return GetOneObjInfoByID(uid).OwnerUid == RoleM.Instance.Uid;
            }
            return uid == RoleM.Instance.Uid;
        }

        public void SetPlayerHP(long uid, int changNum, int maxHp, bool isRefresh = false)
        {
            if (getPlayerInfoByID(uid) != null)
            {
                if (isRefresh)
                    getPlayerInfoByID(uid).CurHp = changNum;
                else
                    getPlayerInfoByID(uid).CurHp += changNum;

                if (maxHp > 0)
                    getPlayerInfoByID(uid).MaxHp = maxHp;
            }
            else if (GetOneObjInfoByID(uid) != null)
            {
                if (isRefresh)
                    GetOneObjInfoByID(uid).CurHp = changNum;
                else
                    GetOneObjInfoByID(uid).CurHp += changNum;

                if (maxHp > 0)
                    GetOneObjInfoByID(uid).MaxHp = maxHp;
            }
        }

        public void SetPlayerLv(long uid, int lv)
        {
            if (getPlayerInfoByID(uid) != null)
            {
                getPlayerInfoByID(uid).Level = lv;
            }
        }

        public int GetPlayerHP(long uid)
        {
            if (getPlayerInfoByID(uid) != null)
            {
                return getPlayerInfoByID(uid).CurHp;
            }
            if (GetOneObjInfoByID(uid) != null)
            {
                return GetOneObjInfoByID(uid).CurHp;
            }
            return 0;
        }

        public Dictionary<long, PlayerInfo> getPlayerInfo()
        {
            return PInfoDic;
        }

        #region 遗迹三星条件

        //血量大于百分之五十
        public bool CheakMyTeamHpHalf()
        {
            float CurMaxHp = getPlayerInfoByID(RoleM.Instance.Uid).CurHp;
            float AllHP = getPlayerInfoByID(RoleM.Instance.Uid).MaxHp;

            foreach (var info in ObjDic.Values)
            {
                if (info.OwnerUid == RoleM.Instance.Uid)
                {
                    CurMaxHp += info.CurHp;
                    AllHP += info.MaxHp;
                }
            }
            if (AllHP <= 0)
                return false;

            return (CurMaxHp / AllHP) >= 0.5f;
        }

        //怪物全部死亡
        public bool CheakMonstersDead()
        {
            return MonsterManager.Instance.MonesterCount <= 0;
        }

        //玩家全部死亡
        public bool CheakMyTeamAllLife()
        {
            if (getPlayerInfoByID(RoleM.Instance.Uid).CurHp <= 0)
                return false;

            foreach (var info in ObjDic.Values)
            {
                if (info.OwnerUid == RoleM.Instance.Uid && info.CurHp <= 0)
                {
                    return false;
                }
            }
            return true;
        }

        #endregion

        public bool IsDayDo(int roomID = -1)
        {
            if (roomID == -1)
            {
                roomID = (int)GameRoomM.Instance.CFGRoom.Roomtype;
            }

            if (roomID == (int)CFG_Room.RoomType.Gold || roomID == (int)CFG_Room.RoomType.EXP
            || roomID == (int)CFG_Room.RoomType.SXS || roomID == (int)CFG_Room.RoomType.HS
            || roomID == (int)CFG_Room.RoomType.QHS || roomID == (int)CFG_Room.RoomType.SY)
                return true;

            return false;
        }

        public bool IsTurnBasedRoom()
        {
            return CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More || CFGRoom.Roomtype == CFG_Room.RoomType.WuXing;
        }

        public bool IsPVPMoreRoom()
        {
            return CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More;
        }

        public bool IsWuXingRoom()
        {
            return CFGRoom.Roomtype == CFG_Room.RoomType.WuXing;
        }

        /// <summary>
        /// 攻击对象类型
        /// </summary>
        public enum AttackClass
        {
            Monster,
            Player,
            All,
            Team,
        }

        public AttackClass GameAttack
        {
            get
            {
                if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_One || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More
                    || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.YZDD)
                {
                    return AttackClass.Player;
                }
                else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace)
                {
                    return AttackClass.Team;
                }
                else
                {
                    return AttackClass.Monster;
                }
            }
        }
    }
}
