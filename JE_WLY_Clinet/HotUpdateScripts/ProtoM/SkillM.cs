using JEngine.Helper;
using JEngine.Singleton;
using pb;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;

namespace HotUpdateScripts.ProtoM
{
    public class SkillClass
    {
        public int SkillID;
        public int CfgId;
        /// <summary>
        /// 进阶ID
        /// </summary>
        public int advance;

        /// <summary>
        /// CD
        /// </summary>
        public long Cooldown;

        public bool isAuto = false;
        public float EndCDTime = 0;
        public long Uid = 0;
    }

    /// <summary>
    /// 个人技能信息
    /// </summary>
    public class SkillM : SingletonTemplate<SkillM>
    {
        //拥有技能
        private List<SkillClass> SkillClassList = new List<SkillClass>();

        //装备技能
        private Dictionary<int, int> SkillEquipDic = new Dictionary<int, int>();

        //自己技能信息
        private List<SkillInfo> SkillInfoList = new List<SkillInfo>();

        //设置拥有技能
        public void SetSkillInfo(List<SkillInfo> Skills)
        {
            SkillInfoList = Skills;
            SkillClassList.Clear();
            for (int i = 0; i < Skills.Count; i++)
            {
                AddSkillClassList(Skills[i]);
            }

            SkillClassList.Sort((a, b) => { return a.SkillID.CompareTo(b.SkillID); });
        }

        public List<SkillInfo> GetSkillInfoList()
        {
            return SkillInfoList;
        }

        public SkillInfo GetSkillInfoBySkillID(int skillID)
        {
            return SkillInfoList.Find((a) => { return CFG_Skill.GetById(a.CfgId).SkillID == skillID; });
        }

        public SkillInfo GetSkillInfoByCFGID(int cfgID)
        {
            return SkillInfoList.Find((a) => { return a.CfgId == cfgID; });
        }

        //设置已装备技能
        public void SetSkillEquipDic(int[] skills)
        {
            SkillEquipDic.Clear();
            for (int i = 0; i < skills.Length; i++)
            {
                SkillEquipDic.Add(i, skills[i]);
            }
        }

        //添加拥有技能列表
        private void AddSkillClassList(SkillInfo item)
        {
            SkillClassList.Add(new SkillClass
            {
                SkillID = CFG_Skill.GetById(item.CfgId).SkillID,
                CfgId = item.CfgId,
                advance = item.Advanced,
            });
            SkillClassList.Sort((a, b) => { return a.SkillID.CompareTo(b.SkillID); });
        }

        //获得新技能
        public void AddNewSkillInfo(SkillInfo item)
        {
            AddSkillClassList(item);
            SkillInfoList.Add(item);
        }

        //获取拥有技能
        public List<SkillClass> GetSkillInfo()
        {
            return SkillClassList;
        }

        //获取一个技能
        public SkillClass GetOneSkillInfo(int skillId)
        {
            return SkillClassList.Find((a) => { return a.SkillID == skillId; });
        }

        //升级 技能
        public void SetUpOneSkill(int cfgId, int newCfgId)
        {
            for (int i = 0; i < SkillInfoList.Count; i++)
            {
                if (SkillInfoList[i].CfgId == cfgId)
                {
                    SkillInfoList[i].CfgId = newCfgId;
                    break;
                }
            }
            UpdateSkillClassList(cfgId, newCfgId);
            RefrshSoltSkill(cfgId, newCfgId);
        }

        public void UpdateSkillClassList(int cfgId, int newCfgId)
        {
            for (int i = 0; i < SkillClassList.Count; i++)
            {
                if (SkillClassList[i].CfgId == cfgId)
                {
                    SkillClassList[i].CfgId = newCfgId;
                    break;
                }
            }
        }

        /// <summary>
        /// 刷新技能槽位
        /// </summary>
        /// <param name="cfgId"></param>
        /// <param name="newCfgId"></param>
        public void RefrshSoltSkill(int cfgId, int newCfgId)
        {
            for (int i = 0; i < SkillEquipDic.Count; i++)
            {
                if (SkillEquipDic.ElementAt(i).Value == cfgId)
                {
                    SkillEquipDic[SkillEquipDic.ElementAt(i).Key] = newCfgId;
                    break;
                }
            }
        }

        //升星 技能
        public void SetUpAdvanceSkill(int cfgId, int advance)
        {
            var oneSkill = SkillClassList.Find((a) => { return a.CfgId == cfgId; });
            oneSkill.advance = advance;

            for (int i = 0; i < SkillInfoList.Count; i++)
            {
                if (SkillInfoList[i].CfgId == cfgId)
                {
                    SkillInfoList[i].Advanced = advance;
                    break;
                }
            }
        }

        private int GetNowlvBySkillId(int skillId)
        {
            int lv = 1;
            for (int i = 0; i < SkillClassList.Count; i++)
            {
                if (SkillClassList[i].SkillID == skillId)
                {
                    lv = Mathf.Max(CFG_Skill.GetById(SkillClassList[i].CfgId).SkillLevel, lv);
                }
            }
            return lv;
        }

        //个人信息面板获取技能
        public List<CFGItem_Skill> GetRolePanelSkillInfo()
        {
            List<CFGItem_Skill> _ShowSkillInfo = new List<CFGItem_Skill>();
            var SkillDic = CFG_Skill.GetSkillsByJob(RoleM.Instance.Job);
            foreach (var item in SkillDic)
            {
                var skill = item.Value[GetNowlvBySkillId(item.Key)];
                if (skill.Skilltype == (int)CFG_Skill.SkillType.Job || skill.Skilltype == (int)CFG_Skill.SkillType.JobCtr)
                    _ShowSkillInfo.Add(item.Value[GetNowlvBySkillId(item.Key)]);
            }
            _ShowSkillInfo.Sort((a, b) => { return a.UUID.CompareTo(b.UUID); });
            return _ShowSkillInfo;
        }

        //判断该技能是否被装备
        public bool CheakSkillEquip(int cfgId)
        {
            int skillId = CFG_Skill.GetById(cfgId).SkillID;
            foreach (int id in SkillEquipDic.Values)
            {
                if (id != 0 && CFG_Skill.GetById(id).SkillID == skillId)
                    return true;
            }
            return false;
        }

        //判断该技能是否被装备
        public int GetSkillEquipIdx(int cfgId)
        {
            int idx = 0;
            var _List = SkillEquipDic.Values.ToArray();
            for (int i = 0; i < _List.Length; i++)
            {
                if (_List[i] == cfgId)
                {
                    return idx;
                }
                idx++;
            }
            return -1;
        }

        //根据技能槽位获取装备技能
        public int GetSkillEquipBySolt(int idx)
        {
            return SkillEquipDic[idx];
        }

        //设置一个槽位技能信息
        public void SetOneSkillEquip(int cfgId, int op, int solt)
        {
            //装备
            if (op == 1)
            {
                SkillEquipDic[solt] = cfgId;
            }
            else
            {
                SkillEquipDic[solt] = 0;
            }
        }

        //战斗技能控制信息
        private Dictionary<int, SkillClass> GameSkillDic = new Dictionary<int, SkillClass>();
        private List<int> Skills = new List<int>();

        /// <summary>
        /// 战斗技能CD时间
        /// </summary>
        private Dictionary<int, long> SkillCDdic = new Dictionary<int, long>();

        /// <summary>
        /// 初始化战斗技能
        /// </summary>
        /// <param name="isInit"></param>
        public void OnInitGameSkill(bool isInit = true)
        {
            Skills.Clear();
            GameSkillDic.Clear();

            if (isInit)
            {
                SkillCDdic.Clear();
            }

            for (int i = 0; i < 13; i++)
            {
                int cfgId = GetSkillEquipBySolt(i);

                if (cfgId == 0)
                    continue;

                GameSkillDic.Add(cfgId, new SkillClass
                {
                    CfgId = cfgId,
                    isAuto = GameUtil.Instance.GetPlayerPrefsInt(RoleM.Instance.Uid + "GameSkillDic" + i) == 0,
                    Uid = PartnerM.Instance.GetParnerUidBySlotSkillID(cfgId),
                });

                if (!SkillCDdic.ContainsKey(cfgId))
                    SkillCDdic.Add(cfgId, GameUtil.Instance.GetCurrentTimeSeconds());

                Skills.Add(cfgId);
            }

            GameSkillDic.Add(GameConst.GetNormallSkillByJob((JobType)RoleM.Instance.Job), new SkillClass
            {
                CfgId = GameConst.GetNormallSkillByJob((JobType)RoleM.Instance.Job),
                Uid = 0
            });
        }

        //设置战斗技能是否开启
        public void SetGameSkillOpen(int cfgId, bool isOpen, int idx)
        {
            if (idx != 0)
            {
                GameUtil.Instance.SetPlayerPrefs(RoleM.Instance.Uid + "GameSkillDic" + idx, isOpen ? 0 : 1);
            }
            GameSkillDic[cfgId].isAuto = isOpen;

        }

        public bool GetGameSkillOpen(int cfgId, int idx = -1)
        {
            return GameSkillDic[cfgId].isAuto;
        }

        //获取技能的归属ID
        public long GetUiDBySkillId(int cfgId)
        {
            if (!GameSkillDic.ContainsKey(cfgId))
                return 0;
            return GameSkillDic[cfgId].Uid;
        }

        //设置战斗技能CD时间
        public void SetGameSkillDealy(int cfgId, long deayTime)
        {
            SkillCDdic[cfgId] = deayTime;
        }

        /// <summary>
        /// 获取技能时间
        /// </summary>
        /// <param name="cfgId"></param>
        /// <param name="deayTime"></param>
        public long GetGameSkillDealy(int cfgId)
        {
            return SkillCDdic[cfgId];
        }

        //获取战斗内技能
        public List<int> GetGameFightSKill()
        {
            return Skills;
        }

        private int ItemCompare(int a, int b)
        {
            var skillA = CFG_Skill.GetById(a);
            var skillB = CFG_Skill.GetById(b);

            // 检查是否为超级技能
            bool isSuperA = skillA.Skilltype == (int)CFG_Skill.SkillType.Super || skillA.Skilltype == (int)CFG_Skill.SkillType.SuperCtr;
            bool isSuperB = skillB.Skilltype == (int)CFG_Skill.SkillType.Super || skillB.Skilltype == (int)CFG_Skill.SkillType.SuperCtr;

            // 超级技能优先
            if (isSuperA && !isSuperB)
                return 1;
            if (!isSuperA && isSuperB)
                return -1;

            // 比较攻击距离
            if (skillA.MaxDistance > skillB.MaxDistance)
                return 1;
            if (skillA.MaxDistance < skillB.MaxDistance)
                return -1;

            return 0;
        }
    }

}

