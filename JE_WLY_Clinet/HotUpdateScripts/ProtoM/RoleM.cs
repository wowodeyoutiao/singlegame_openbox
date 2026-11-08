using JEngine.Helper;
using JEngine.Singleton;
using MetaFramework.UI;

using pb;
using Spine;

using System;
using System.Collections.Generic;
using UnityEngine;


namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 玩家总信息
    /// </summary>
    public class RoleM : SingletonTemplate<RoleM>
    {
        public string Nickname { get; set; }
        public long Uid { get; set; }
        public string HeadUrl { get; set; }
        public int Sex { get; set; }
        public string IP { get; set; }
        public string Post { get; set; }

        public string Token { get; set; }

        public string Login_Token { get; set; }

        /// <summary>
        /// 战魂 游戏内的兑换货币，对应钻石
        /// </summary>
        public long Diamanta { get; set; }

        /// <summary>
        /// 战魂 游戏内的兑换货币，对应钻石
        /// </summary>
        public long Soul { get; set; }

        /// <summary>
        /// 点券 游戏内的充值货币，对应点券
        /// </summary>
        public long Coupon { get; set; }
        /// <summary>
        /// 体力 进入副本消耗
        /// </summary>
        public long Energy { get; set; } 
        /// <summary>
        /// 职业
        /// </summary>
        public int Job { get; set; }

        /// <summary>
        /// 经验
        /// </summary>
        public long Exp { get; set; }

        /// <summary>
        /// 等级
        /// </summary>
        public int Level { get; set; } 

        /// <summary>
        /// vip等级
        /// </summary>
        public int Vip { get; set; }

        /// <summary>
        /// vip经验
        /// </summary>
        public long VipExp { get; set; }

        /// <summary>
        /// 下次体力恢复剩余时间
        /// </summary>
        public int EnergyRecoveryTime { get; set; }

        /// <summary>
        /// 金币 游戏内的通用货币
        /// </summary>
        public long Coin { get; set; }

        /// <summary>
        /// 个人战力
        /// </summary>
        public int Combat { get; set; }

        /// <summary>
        /// 总战力
        /// </summary>
        public int TotalCombat { get; set; }

        /// <summary>
        /// 命格等级
        /// </summary>
        public int FateLevel { get; set; }

        /// <summary>
        /// 命格经验
        /// </summary>
        public int FateExp { get; set; }

        /// <summary>
        /// 试炼凭证
        /// </summary>
        public long TrialVoucher { get; set; }

        /// <summary>
        /// 宝箱数量
        /// </summary>
        public long TreasuerBox { get; set; }


        /// <summary>
        /// 新手引导
        /// </summary>
        //public bool isGuiding { get; set; }

        /// <summary>
        /// 新手引导步骤
        /// </summary>
        //public int GuideId { get; set; }

        /// <summary>
        /// 引导已完成的类型列表
        /// </summary>
        //public List<int> GuideTypes { get; set; }

        /// <summary>
        /// 开箱数据缓存
        /// </summary>
        public List<ItemInfo> OpenBoxItems { get; set; } = new List<ItemInfo>();

        /// <summary>
        /// 开箱当前等级
        /// </summary>
        public int BoxNowlevel { get; set; }

        /// <summary>
        /// 宝箱升级结束时间
        /// </summary>
        public long BoxNextRecoveryTime {  get; set; }

        /// <summary>
        /// 开箱剩余时间
        /// </summary>
        public int BoxTimeAllSecound { get; set; }

        /// <summary>
        /// 宠物可用经验数值
        /// </summary>
        public long PetUsedExp { get; set; }

        /// <summary>
        /// 加速卷轴道具
        /// </summary>
        public long BoxQuick { get; set; }

        /// <summary>
        /// 所在的行会ID
        /// </summary>
        public long GuildId { get; set; }

        /// <summary>
        /// 工会经验
        /// </summary>
        public long GuildContribution { get; set; }

        /// <summary>
        /// 工会许愿点
        /// </summary>
        public long GuildWish { get; set; }

        /// <summary>
        /// 工会许愿点恢复时间
        /// </summary>
        public int WishRecoveryTime { get; set; }

        /// <summary>
        /// 默认开箱装备等级
        /// </summary>
        public int NowOpenLevel { get; set; }

        /// <summary>
        /// 默认开箱装备职业
        /// </summary>
        public int NowOpenJob { get; set; }

        /// <summary>
        /// 是否收到了引导数据
        /// </summary>
        public bool IsReceivedGuideData { get; set; }


        /// <summary>
        /// 角色登录时的初始等级，用于判断伙伴引导
        /// </summary>
        public int guidePartnerRecordLevel { get; set; }

        /// <summary>
        /// 是否为登录后第一次收到角色基础信息
        /// </summary>
        public bool IsFirstReceiveBaseInfo { get; set; } = true;

        public int MainBattleId { get; set; }

        /// <summary>
        /// 穿戴的称号ID
        /// </summary>
        public int TitleId { get; set; }

        /// <summary>
        /// 角色注册时间
        /// </summary>
        public long RegisterAt { get; set; }

        /// <summary>
        /// 星愿币
        /// </summary>
        public long StarWishCoin { get; set; }

        /// <summary>
        /// 当前头像ID
        /// </summary>
        public int CurrentAvatarId { get; set; }

        /// <summary>
        /// 当前头像框ID
        /// </summary>
        public int CurrentFrameId { get; set; }

        /// <summary>
        /// 玩家头像信息
        /// </summary>
        public List<PlayerAvatarInfo> avatars = new List<PlayerAvatarInfo>();

        /// <summary>
        /// 玩家头像框信息
        /// </summary>
        public List<PlayerAvatarInfo> frames = new List<PlayerAvatarInfo>();

        /// <summary>
        /// 洗练信息
        /// </summary>
        public List<CharacterRefineData> characterRefineDatas = new List<CharacterRefineData>();

        /// <summary>
        /// 洗练信息,购买的槽位
        /// </summary>
        public List<RefineSlot> purchasedRefineSlots = new List<RefineSlot>();

        /// <summary>
        /// 洗练结果缓存
        /// 选择或者放弃洗练时清空
        /// </summary>
        public pb.RefineResponse refineResuleCache;

        private Dictionary<int,BuffInfo> BuffInfos = new Dictionary<int, BuffInfo>();

        public void RefreshBuffInfo(List<BuffInfo> buffs)
        {
            if (buffs != null)
            {
                BuffInfos.Clear();
                for (int i = 0; i < buffs.Count; i++)
                {
                    buffs[i].RemainingTime +=(int)(Time.time * 1000);
                    BuffInfos.Add(buffs[i].BuffId, buffs[i]);
                }
            }

        }

        public void AddBuffInfo(BuffInfo buff)
        {
            buff.RemainingTime += (int)(Time.time * 1000);
            if (!BuffInfos.ContainsKey(buff.BuffId))
            {
                BuffInfos.Add(buff.BuffId, buff);
            }
            else
            {
                BuffInfos[buff.BuffId] = buff;
            }
        }

        public void RemoveBuffInfo(int buffID)
        {
            if (BuffInfos.ContainsKey(buffID))
            {
                BuffInfos.Remove(buffID);
            }
        }

        public bool CheakBuffId(int buffID)
        {
            return BuffInfos.ContainsKey(buffID);
        }

        public bool CheakBuffId(params int[] buffIDs)
        {
            for (int i = 0; i < buffIDs.Length; i++)
            {
                if (BuffInfos.ContainsKey(buffIDs[i]))
                {
                    return true;
                }
            }
            return false;
        }

        public long GetBuffOwnerUid(int buffID)
        {
            if (BuffInfos.ContainsKey(buffID))
            {
                return BuffInfos[buffID].OwnerUid;
            }
            return 0;
        }

        public int GetBuffRemainingTime(int buffID)
        {
            if (BuffInfos.ContainsKey(buffID))
            {
                return BuffInfos[buffID].RemainingTime;
            }
            return 0;
        }

        /// <summary>
        /// 更新某个装备槽位的洗练数据
        /// </summary>
        /// <param name="characterSlot">角色槽位</param>
        /// <param name="equipSlot">装备槽位</param>
        /// <param name="refinedSlot">属性槽位</param>
        /// <param name="stats">属性信息</param>
        public void UpdateEquipRefineData(CharacterSlot characterSlot, EquipSlot equipSlot, RefinedSlot refinedSlot, Stats stats)
        {
            if (stats == null) return;

            CharacterRefineData refineData = characterRefineDatas.Find(x => x.character_slot == characterSlot);
            if (refineData != null)
            {
                var equipRefineData = refineData.equip_refines.Find(x => x.equip_slot == equipSlot);
                if (equipRefineData != null)
                {
                    var data = equipRefineData.refined_stats.Find(x => x.refined_slot == refinedSlot);
                    if (data != null)
                    {
                        data.stats = stats;
                    }
                    else
                    {
                        equipRefineData.refined_stats.Add(new RefineData
                        {
                            refined_slot = refinedSlot,
                            stats = stats
                        });
                    }
                }
                else
                {
                    refineData.equip_refines.Add(new EquipRefineData
                    {
                        equip_slot = equipSlot,
                        refined_stats = new List<RefineData>
                        {
                            new RefineData
                            {
                                refined_slot = refinedSlot,
                                stats = stats
                            }
                        }
                    });
                }
            }
            else
            {
                characterRefineDatas.Add(new CharacterRefineData
                {
                    character_slot = characterSlot,
                    equip_refines = new List<EquipRefineData>
                    {
                        new EquipRefineData
                        {
                            equip_slot = equipSlot,
                            refined_stats = new List<RefineData>
                            {
                                new RefineData
                                {
                                    refined_slot = refinedSlot,
                                    stats = stats
                                }
                            }
                        }
                    }
                });
            }
        }

        public bool HasAvatar(int id)
        {
            return avatars.Exists(x => x.id == id);
        }

        public bool HasFrame(int id)
        {
            return frames.Exists(x => x.id == id);
        }

        public bool TryGetAvatar(int id, out PlayerAvatarInfo info)
        {
            info = avatars.Find(x => x.id == id);
            return info != null;
        }

        public bool TryGetFrame(int id, out PlayerAvatarInfo info)
        {
            info = frames.Find(x => x.id == id);
            return info != null;
        }

        /// <summary>
        /// 装备槽位信息
        /// </summary>
        /// <param name="characterSlot">角色槽位</param>
        /// <param name="equipSlot">装备槽位</param>
        /// <returns></returns>
        public EquipRefineData GetOneEquipSlotData(CharacterSlot characterSlot, EquipSlot equipSlot)
        {
            CharacterRefineData refineData = characterRefineDatas.Find(x => x.character_slot == characterSlot);
            if (refineData != null)
            {
                return refineData.equip_refines.Find(x => x.equip_slot == equipSlot);
            }
            return null;
        }

        /// <summary>
        /// 装备槽位上的一条属性
        /// </summary>
        /// <param name="characterSlot">角色槽位</param>
        /// <param name="equipSlot">装备槽位</param>
        /// <param name="refinedSlot">属性槽位</param>
        /// <returns></returns>
        public RefineData GetEquipSlotOneAttributeData(CharacterSlot characterSlot, EquipSlot equipSlot, RefinedSlot refinedSlot)
        {
            CharacterRefineData refineData = characterRefineDatas.Find(x => x.character_slot == characterSlot);
            if (refineData != null)
            {
                EquipRefineData equipRefineData = refineData.equip_refines.Find(x => x.equip_slot == equipSlot);
                if (equipRefineData != null)
                    return equipRefineData.refined_stats.Find(x => x.refined_slot == refinedSlot);
            }
            return null;
        }

        /// <summary>
        /// 洗练槽位是否已购买
        /// </summary>
        /// <param name="characterSlot">角色槽位</param>
        /// <param name="equipSlot">装备槽位</param>
        /// <param name="refinedSlot">洗练槽位</param>
        /// <returns></returns>
        public bool IsRefineSlotBuy(CharacterSlot characterSlot, EquipSlot equipSlot, RefinedSlot refinedSlot)
        {
            return purchasedRefineSlots.Exists(x => x.character_slot == characterSlot && x.equip_slot == equipSlot && x.refined_slot == refinedSlot);
        }

        public void SetSelectRole(HTTP_CharacterInfo selctInfo)
        {
            Uid = selctInfo.uid;
            HeadUrl = selctInfo.headUrl;
            Sex = selctInfo.sex;
            Nickname = selctInfo.name;
            Coin = selctInfo.coin;
            Soul = selctInfo.soul;
            Coupon = selctInfo.coupon;
            Energy = selctInfo.energy;
            Job = selctInfo.job;
            Exp = selctInfo.exp;
            Level = selctInfo.level;
            Vip = selctInfo.vip;
            VipExp = selctInfo.recharge;
            RegisterAt = selctInfo.registerAt;
        }

        /// <summary>
        /// 登录返回信息
        /// </summary>
        public List<HTTP_CharacterInfo> LoginCharacterInfo { get; set; }
        /// <summary>
        /// Http 获取个人信息
        /// </summary>
        /// <param name="_loginRes"></param>
        public void SetHttpGetInfo(HTTP_LoginResponse _loginRes )
        {
            LoginCharacterInfo = _loginRes.characterList;
            Login_Token = _loginRes.token.ToString();
        }

        public Dictionary<pb.Attribute, long> RoleAttributeDic = new Dictionary<pb.Attribute, long>();

        public List<Stats> Attributes = new List<Stats>();
        /// <summary>
        /// 设置角色属性
        /// </summary>
        /// <param name="_Attributes"></param>
        public void SetRoleStats(List<Stats> _Attributes)
        {
            RoleAttributeDic.Clear();
            Attributes = _Attributes;
            for (int i = 0; i < _Attributes.Count; i++)
            {
                RoleAttributeDic.Add(_Attributes[i].attr, _Attributes[i].Num);
            }
        }

        public List<Stats> GetRoleAllStats()
        {
            return Attributes;
        }

        //获取角色属性数值
        public long GetRoleAttrValueByAttr(pb.Attribute _Attributes)
        {
            long _value = 0;
            if (!RoleAttributeDic.TryGetValue(_Attributes, out _value) )
            {
                Debug.LogWarning($"属性{_Attributes}为空");
            }
            return _value;
        }

        //角色装备信息
        public Dictionary<EquipSlot, SlotInfo> RoleEquipDic = new Dictionary<EquipSlot, SlotInfo>();

        public void SetRoleEquipDic(List<SlotInfo> equipInfo)
        {
            RoleEquipDic.Clear();
            for (int i = 1; i < (int)EquipSlot.MaxEquipSlot; i++)
            {
                RoleEquipDic.Add((EquipSlot)i, equipInfo[i]);
            }
        }

        public bool isEquip(int equipSlot)
        {
            return RoleEquipDic.ContainsKey((EquipSlot)equipSlot);
        }

        //刷新某个部位的装备
        public void SetRoleOneEquip(ItemInfo ItemInfo, EquipSlot equipSlot,bool isAdd)
        {
            if (isAdd)
            {
                if (RoleEquipDic.ContainsKey(equipSlot))
                {
                    RoleEquipDic[equipSlot].Item = ItemInfo;
                }
            }
            else
            {
                if (RoleEquipDic.ContainsKey(equipSlot))
                {
                    RoleEquipDic[equipSlot].Item = default;
                }
            }
        }

        //刷新某个部位的宝石镶嵌
        public void SetRoleOneEquipInset(int cid, EquipSlot equipSlot,int insetIdx, bool isAdd)
        {
            if (isAdd)
            {
                if (RoleEquipDic.ContainsKey(equipSlot))
                {
                    RoleEquipDic[equipSlot].Gemstone[insetIdx] = cid;
                }
            }
            else
            {
                if (RoleEquipDic.ContainsKey(equipSlot))
                {
                    RoleEquipDic[equipSlot].Gemstone[insetIdx] = 0;
                }
            }
        }

        //检测该部位是为空装备
        public bool CheckIsNullEquip(EquipSlot eSlot)
        {
            if (!RoleEquipDic.ContainsKey(eSlot))
            {
                return true;
            }
            return RoleEquipDic[eSlot].Item == null ||  RoleEquipDic[eSlot].Item.CfgId == 0;
        }

        public bool CheckIsBestEquip(EquipSlot eSlot)
        {
            if (CheckIsNullEquip(eSlot))
            {
                return false;
            }
            return RoleEquipDic[eSlot].Item.Equip != null && RoleEquipDic[eSlot].Item.Equip.Attributes.Count != 0;
        }

        //获取某个部位的信息
        public SlotInfo GetSlotInfoBySolt(EquipSlot eSlot)
        {
            return RoleEquipDic[eSlot];
        }

        //设置装备强化
        public void SetEquipeReinforeceLv(EquipSlot equipSlot, int Lv)
        {
            RoleEquipDic[equipSlot].Enhance = Lv;
        }

        //设置装备升星
        public void SetEquipeStarLv(EquipSlot equipSlot, int Lv)
        {
            RoleEquipDic[equipSlot].Stars = Lv;
        }

        //获取全身装备的最小等级（用来比较强化大师）
        public int GetEquipReinforce()
        {
            int reinforceLv = 9999;
            for (int i = 1; i < (int)EquipSlot.MaxEquipSlot; i++)
            {
                if (RoleEquipDic[(EquipSlot)i].Enhance < reinforceLv)
                {
                    reinforceLv = RoleEquipDic[(EquipSlot)i].Enhance;
                }
            }
            return reinforceLv;
        }

        //获取全身装备的最小星级（用来比较升星大师）
        public int GetEquipStar()
        {
            int Lv = 9999;
            for (int i = 1; i < (int)EquipSlot.MaxEquipSlot; i++)
            {
                if (RoleEquipDic[(EquipSlot)i].Stars < Lv)
                {
                    Lv = RoleEquipDic[(EquipSlot)i].Stars;
                }
            }
            return Lv;
        }


        #region PVP

        public class PVPClass
        {
            public MatchType matchType;
            /// <summary>
            /// 今日一挑战次数
            /// </summary>
            public int DailyMatchTimes;

            /// <summary>
            /// 今日已购买次数
            /// </summary>
            public int DailyBuyNum;

            public int CanMatchCount;
        }

        private Dictionary<MatchType, PVPClass> PVPDic = new Dictionary<MatchType, PVPClass>();

        public void SetAddPVPInfo(MatchType _matchType,int dailyMatchTimes,int dailyBuyNum)
        {
            if (!PVPDic.ContainsKey(_matchType))
            {
                PVPDic.Add(_matchType, new PVPClass { matchType = _matchType });
            }
            PVPDic[_matchType].DailyBuyNum = dailyBuyNum;
            PVPDic[_matchType].DailyMatchTimes = dailyMatchTimes;

            int dailyTime = _matchType == MatchType.SINGLE_COMBAT ? CFG_XMLConfig.xmlConfig.PVPInfo.dailyTime : CFG_XMLConfig.xmlConfig.GVGInfo.dailyTime;
            PVPDic[_matchType].CanMatchCount = dailyTime + dailyBuyNum - dailyMatchTimes;
        }

        //获取可挑战次数
        public int GetPVPMatchTimes(int _matchType)
        {
            if (!PVPDic.ContainsKey((MatchType)_matchType))
                return 0;
            else
                return PVPDic[(MatchType)_matchType].CanMatchCount;
        }

        public int GetPVPRedNum()
        {
            if(GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false))
                return 0;
            if(GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false))
                return 0;

            //if (GetPVPOneSeason() == null)
            //    return 0;

            if (!GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.PVPInfo.dailyBegin, CFG_XMLConfig.xmlConfig.PVPInfo.dailyEnd))
                return 0;

            return GetPVPMatchTimes((int)MatchType.SINGLE_COMBAT);
        }

        public int GetGVGRedNum()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                return 0;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                return 0;

            if (!GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.GVGInfo.dailyBegin, CFG_XMLConfig.xmlConfig.PVPInfo.dailyEnd))
                return 0;

            return GetPVPMatchTimes((int)MatchType.FREE_FOR_ALL);
        }

        //检测是否可以购买
        public int GetCanNumBuyPVP(int _matchType)
        {
            int storeId = (MatchType)_matchType == MatchType.SINGLE_COMBAT ? CFG_XMLConfig.xmlConfig.PVPInfo.storeId : CFG_XMLConfig.xmlConfig.GVGInfo.storeId;
            return CFG_Store.GetById(storeId).Limite.limitnum - (int)StoreM.Instance.GetBuyNumByID(storeId);
            //return CFG_Store.GetById(storeId).Limite.limitnum - PVPDic[(MatchType)_matchType].DailyBuyNum;
        }
        
        public CFG_XMLConfig.XML_Config.XML_Item_PVP.PVPseason GetPVPOneSeason()
        {
            foreach (var item in CFG_XMLConfig.xmlConfig.PVPInfo.PVPSeasonDic.Values)
            {
                if (JEngine.Helper.GameUtil.Instance.CheakCurTimeDay(item.starttime,item.endtime))
                    return item;
            }
            return null;
        }



        #endregion


        #region Train

        public class TrainClass 
        {
            /// <summary>
            /// 是否可以领取
            /// </summary>
            public bool Ready;
            public List<ItemInfo> Reward;

            /// <summary>
            /// 已完成
            /// </summary>
            public int Layers;
        }

        private TrainClass _trainClass = new TrainClass();

        public void SetTrainInfo(bool _ready,int _layer,List<ItemInfo> _reward)
        {
            _trainClass.Ready = _ready;
            _trainClass.Layers = _layer;
            _trainClass.Reward = _reward;
        }

        public int GetTrainLayers()
        {
            return _trainClass.Layers;
        }

        public TrainClass GetTrainInfo()
        {
            return _trainClass;
        }


        #endregion


        #region 每日必做
        private Dictionary<int, int> ChallengeDic = new Dictionary<int, int>();

        public void AddChallengeInfo(List<ChallengeInfo> list)
        {
            ChallengeDic.Clear();
            foreach (var item in list)
            {
                ChallengeDic.Add((int)item.MType, item.Times);
            }
        }

        public int GetCanTimeByRoomType(int roomType)
        {
            //int allCount = CFG_XMLConfig.xmlConfig.DayRoomDic[(CFG_Room.RoomType)roomType].challengeTimeNum + CFG_VIP.GetAddRoomCount(Vip, roomType);
            //return ChallengeDic.ContainsKey(roomType) ? allCount - ChallengeDic[roomType] : allCount;
            return ChallengeDic.ContainsKey(roomType) ? ChallengeDic[roomType] : CFG_XMLConfig.xmlConfig.DayRoomDic[(CFG_Room.RoomType)roomType].challengeTimeNum + CFG_VIP.GetAddRoomCount(Vip, roomType);
        }

        #endregion

        /// <summary>
        /// 是否加入了行会
        /// </summary>
        /// <returns></returns>
        public bool IsJoinedGuild()
        {
            return GuildId != 0;
        }


        /// <summary>
        /// 是否解锁了伙伴槽位
        /// </summary>
        /// <returns></returns>
        public bool IsUnlockedPartnerSlot()
        {
            return Level >= CFG_XMLConfig.xmlConfig.OpenSlotDic[1];
        }

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public int GetDefaultOpenBoxEquipLevel()
        {
            //return NowOpenLevel == 0 ? Level : NowOpenLevel;
            //int slotIndexSelect = GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.DefaultOpenSoltCondition.ToString() + Uid);
            //if (slotIndexSelect == 0)
            //{
            //    return Level;
            //}
            //else
            //{
            //    var slotInfo = PartnerM.Instance.GetOnePartnerSlotInfo(slotIndexSelect);
            //    var partnetInfo = PartnerM.Instance.GetOnePartnerInfoByUid(slotInfo.PartnerUid);
            //    return partnetInfo.Level;
            //}

            return Level;
        }

        /// <summary>
        /// 获取默认开箱装备职业。
        /// 只在开箱请求时传参使用，宝箱设置界面默认值需要显示默认全部。
        /// </summary>
        /// <returns></returns>
        public int GetDefaultOpenBoxEquipJob()
        {
            //if (GuidanceM.Instance.isGuiding)
            //    return Job;

            //int unlockedJob = IsUnlockedPartnerSlot() ? NowOpenJob : Job;
            //return unlockedJob;

            int slotIndexSelect = GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.DefaultOpenSoltCondition.ToString() + Uid);
            if (slotIndexSelect == 0)
            {
                bool hasPartnetInSlots = false;
                for (int i = 0; i < 4; i++)
                {
                    if (HasPartnerInSlot(i + 1))
                    {
                        hasPartnetInSlots = true;
                        break;
                    }
                }

                if (hasPartnetInSlots)
                {
                    return 0;
                }
                else
                {
                    return Job;
                }
            }
            else if(slotIndexSelect == 1)
            {
                return Job;
            }
            else
            {
                var slotInfo = PartnerM.Instance.GetOnePartnerSlotInfo(slotIndexSelect - 1);
                var partnetInfo = PartnerM.Instance.GetOnePartnerInfoByUid(slotInfo.PartnerUid);
                return CFG_Partner.GetById(partnetInfo.CfgId).job % 10;
            }
        }

        public bool HasPartnerInSlot(int slotIndex)
        {
            if (Level < CFG_XMLConfig.xmlConfig.OpenSlotDic[slotIndex])
                return false;

            var slotInfo = PartnerM.Instance.GetOnePartnerSlotInfo(slotIndex);
            if (slotInfo == null || slotInfo.PartnerUid == 0)
                return false;
            return true;
        }

        /// <summary>
        /// 宝箱是否可以升级
        /// </summary>
        /// <returns></returns>
        public bool IsCanUpgradeBox()
        {
            int maxLevel = CFG_BoxLevel.GetMaxLevel();
            if (BoxNowlevel == 0 || BoxNowlevel >= maxLevel)
                return false;

            var cfgBoxLevelItem = CFG_BoxLevel.GetById(BoxNowlevel);

            return MainBattleId >= cfgBoxLevelItem.BattleID && Level >= cfgBoxLevelItem.Rolelevel && GameUtil.Instance.ComCoinToCheck(cfgBoxLevelItem.Consume[0].id, cfgBoxLevelItem.Consume[0].num, null, false);
        }

        public void CheckBoxLevelRed()
        {
            if (!GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.OpenBox, false))
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.BoxLevel1.ToString(), IsCanUpgradeBox() ? 1 : 0);
        }
    }

}

