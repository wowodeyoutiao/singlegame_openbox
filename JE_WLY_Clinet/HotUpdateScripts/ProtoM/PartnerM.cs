using JEngine.Singleton;
using pb;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 伙伴信息
    /// </summary>
    public class PartnerM : SingletonTemplate<PartnerM>
    {
        /// <summary>
        /// 伙伴等级达到等级以后才能打开重铸界面
        /// </summary>
        public const int recastOpenLevel = 20;

        /// <summary>
        /// 付费重铸开启等级
        /// </summary>
        public const int payRecastLevel = 60;

        public int[] CanGetNums = { 0, 0, 0 };

        //伙伴信息（包含槽位上伙伴信息）
        private List<pb.PartnerInfo> partners = new List<PartnerInfo>();
        private Dictionary<long, pb.PartnerInfo> partnersDic = new Dictionary<long, PartnerInfo>();

        //伙伴槽位信息
        private List<pb.PartnerSlotInfo> partnerSlotInfos = new List<PartnerSlotInfo>();

        //伙伴属性信息 和槽位信息中对应
        private Dictionary<long, SyncUserStatsPush.PartnerStats> partnerStatsDic = new Dictionary<long, SyncUserStatsPush.PartnerStats>();

        /// <summary>
        /// 设置伙伴信息
        /// </summary>
        /// <param name="_partners"></param>
        public void SetPartners(List<PartnerInfo> _partners)
        {
            partners = _partners;

            partnersDic.Clear();
            for (int i = 0; i < partners.Count; i++)
            {
                partnersDic.Add(partners[i].Uid, partners[i]);
            }
        }

        public List<PartnerInfo> GetPartners()
        {
            return partners;
        }

        public void SetOnePartners(PartnerInfo _partner)
        {
            bool isAdd = true;
            for (int i = 0; i < partners.Count; i++)
            {
                if (partners[i].Uid == _partner.Uid)
                {
                    isAdd = false;
                    partners[i] = _partner;
                    break;
                }
            }
            //Debug.LogError($"add one: {_partner.CfgId}, isAdd:{isAdd}");
            if (isAdd)
            {
                partners.Add(_partner);
                partnersDic.Add(_partner.Uid, _partner);
            }
            else
                partnersDic[_partner.Uid] = _partner;


            BondM.Instance.CheckAddRed(_partner.CfgId);
        }


        /// <summary>
        /// 获取单个伙伴信息
        /// </summary>
        /// <param name="uid"></param>
        /// <returns></returns>
        public PartnerInfo GetOnePartnerInfoByUid(long uid)
        {
            if (partnersDic.ContainsKey(uid))
                return partnersDic[uid];
            return null;
        }

        /// <summary>
        /// 根据CID获取单个伙伴信息
        /// </summary>
        /// <param name="cid"></param>
        /// <returns></returns>
        public PartnerInfo GetOnePartnerInfoByCFGID(int cid)
        {
            for (int i = 0; i < partners.Count; i++)
            {
                if (partners[i].CfgId == cid)
                    return partners[i];
            }
            return null;
        }
        public void SetOnePartnerLv(long uid, int level)
        {
            var partner = GetOnePartnerInfoByUid(uid);
            if (partner != null)
            {
                partner.Level = level;
            }

        }

        public bool CheakPartnerIsHave(long cid)
        {
            for (int i = 0; i < partners.Count; i++)
            {
                if (partners[i].CfgId == cid)
                    return true;
            }
            return false;
        }


        /// <summary>
        /// 设置伙伴槽位信息
        /// </summary>
        /// <param name="_partnerSlotInfos"></param>
        public void SetPartnerSlotInfo(List<PartnerSlotInfo> _partnerSlotInfos)
        {
            partnerSlotInfos = _partnerSlotInfos;
        }

        /// <summary>
        /// 获取伙伴槽位信息
        /// </summary>
        /// <returns></returns>
        public List<PartnerSlotInfo> GetPartnerSlotInfo()
        {
            return partnerSlotInfos;
        }

        /// <summary>
        /// 根据伙伴ID获取槽位下标
        /// </summary>
        /// <param name="partnerUid"></param>
        /// <returns></returns>
        public int GetSoltIdxByPartnerUid(long partnerUid)
        {
            for (int i = 0; i < partnerSlotInfos.Count; i++)
            {
                if (partnerSlotInfos[i].PartnerUid == partnerUid)
                    return partnerSlotInfos[i].SlotId;
            }
            return -1;
        }

        /// <summary>
        /// 设置伙伴属性信息
        /// </summary>
        /// <param name="_partnerStats"></param>
        public void SetPartnerStatsDic(List<SyncUserStatsPush.PartnerStats> _partnerStats)
        {
            partnerStatsDic.Clear();
            for (int i = 0; i < _partnerStats.Count; i++)
            {
                partnerStatsDic.Add(_partnerStats[i].PartnerUid, _partnerStats[i]);
            }
        }

        public int GetAllSoltCount()
        {
            return partnerStatsDic.Count;
        }

        //获取单个伙伴的属性信息
        public SyncUserStatsPush.PartnerStats GetOnePartnerStats(long partnerUid)
        {
            if (partnerStatsDic.ContainsKey(partnerUid))
                return partnerStatsDic[partnerUid];
            return null;
        }

        /// <summary>
        /// 根据伙伴ID判断是否上阵
        /// </summary>
        /// <param name="partnerId"></param>
        /// <returns></returns>
        public bool IsInSlot(long partnerId)
        {
            for (int i = 0; i < partnerSlotInfos.Count; i++)
            {
                if (partnerSlotInfos[i].PartnerUid == partnerId)
                    return true;
            }
            return false;
        }

        //获取单个槽位信息
        public PartnerSlotInfo GetOnePartnerSlotInfo(int SlotId)
        {
            for (int i = 0; i < partnerSlotInfos.Count; i++)
            {
                if (partnerSlotInfos[i].SlotId == SlotId)
                    return partnerSlotInfos[i];
            }
            return null;
        }

        /// <summary>
        /// 获取槽位是否上阵伙伴
        /// </summary>
        /// <param name="SlotId"></param>
        /// <returns></returns>
        public bool CheakIsNullSolt(int SlotId)
        {
            var soltInfo = GetOnePartnerSlotInfo(SlotId);
            return soltInfo == null || soltInfo.PartnerUid == 0;
        }

        //获取单个伙伴的单个装备位的装备信息
        public ItemInfo GetOnePartnerOneEquip(int SlotId, int equipIdx)
        {
            if (SlotId == 0)
                return null;
            var soltInfo = GetOnePartnerSlotInfo(SlotId);
            if (soltInfo == null || soltInfo.PartnerUid == 0)
                return null;

            return soltInfo.Equips[equipIdx].Item;
        }

        //获取单个伙伴的单个装备位的装备信息
        public SlotInfo GetOnePartnerOneEquipSolt(int SlotId, int equipIdx)
        {
            if (SlotId == 0)
                return null;
            var soltInfo = GetOnePartnerSlotInfo(SlotId);
            return soltInfo.Equips[equipIdx];
        }

        /// <summary>
        /// 设置单个伙伴的单个装备位的强化信息
        /// </summary>
        /// <param name="SlotId"></param>
        /// <param name="equipIdx"></param>
        /// <param name="Level"></param>
        public void SetEquipeReinforeceLv(int SlotId, int equipIdx, int Level)
        {
            var Partner = GetOnePartnerSlotInfo(SlotId);
            Partner.Equips[equipIdx].Enhance = Level;
        }

        /// <summary>
        /// 设置单个伙伴的单个装备位的星级信息
        /// </summary>
        /// <param name="SlotId"></param>
        /// <param name="equipIdx"></param>
        /// <param name="Level"></param>
        public void SetEquipeStarLv(int SlotId, int equipIdx, int Level)
        {
            var Partner = GetOnePartnerSlotInfo(SlotId);
            Partner.Equips[equipIdx].Stars = Level;
        }

        public void SetRoleOneEquipInset(int SlotId, int cid, int equipSlot, int insetIdx, bool isAdd)
        {
            var Partner = GetOnePartnerSlotInfo(SlotId);
            Partner.Equips[equipSlot].Gemstone[insetIdx] = isAdd ? cid : 0;
        }

        //获取伙伴所有技能信息
        public List<SkillInfo> GetSkillInfoByPartnerId(long partnerId)
        {
            var partner = GetOnePartnerInfoByUid(partnerId);
            if (partner == null)
                return null;

            return partner.Skills;
        }

        //获取伙伴单个技能信息
        public SkillInfo GetOneSkillInfo(long partnerId, int SkiiId)
        {
            var skills = GetSkillInfoByPartnerId(partnerId);
            return skills.Find((a) => { return CFG_Skill.GetById(a.CfgId).SkillID == SkiiId; });
        }

        public void SetUpOneSkill(long partnerId, int cfgId, int newCfgId)
        {
            var skills = GetSkillInfoByPartnerId(partnerId);

            for (int i = 0; i < skills.Count; i++)
            {
                if (skills[i].CfgId == cfgId)
                {
                    skills[i].CfgId = newCfgId;
                    break;
                }
            }
        }


        //升星 技能
        public void SetUpAdvanceSkill(long partnerId, int cfgId, int advance)
        {
            var skills = GetSkillInfoByPartnerId(partnerId);

            for (int i = 0; i < skills.Count; i++)
            {
                if (skills[i].CfgId == cfgId)
                {
                    skills[i].Advanced = advance;
                    break;
                }
            }
        }

        public long GetParnerUidBySlotSkillID(int skillID)
        {
            foreach (var item in partnerStatsDic.Values)
            {
                var skill = GetSkillInfoByPartnerId(item.PartnerUid);
                for (int i = 0; i < skill.Count; i++)
                {
                    if (skill[i].CfgId == skillID)
                        return item.PartnerUid;
                }
            }
            return 0;
        }

        public SkillInfo GetSkillInfoBySlotSkillID(int skillID)
        {
            foreach (var item in partnerStatsDic.Values)
            {
                var skill = GetSkillInfoByPartnerId(item.PartnerUid);
                for (int i = 0; i < skill.Count; i++)
                {
                    if (skill[i].CfgId == skillID)
                        return skill[i];
                }
            }
            return null;
        }

        public long GetRoleAttrValueByAttr(long partnerUid, pb.Attribute _Attributes)
        {
            if (partnerStatsDic.TryGetValue(partnerUid, out var partnerStats))
            {
                if (partnerStats.Attributes != null && partnerStats.Attributes.Count > 0)
                {
                    for (int i = 0; i < partnerStats.Attributes.Count; i++)
                    {
                        if (partnerStats.Attributes[i].attr == _Attributes)
                            return partnerStats.Attributes[i].Num;
                    }
                }
            }

            return 0;
        }
    }

}

