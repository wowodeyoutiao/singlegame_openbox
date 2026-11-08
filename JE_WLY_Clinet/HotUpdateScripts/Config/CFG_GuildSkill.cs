using JEngine.Helper;
using System.Collections.Generic;
using System.Linq;

/// <summary>
/// 引导
/// </summary>
public class CFGItem_GuildSkill : CFGItemBase
{
    public int UUID;

    public override long GetIdx()
    {
        return (int)UUID;
    }

    public int SkillID;

    public string SkillName;

    public int SkillLevel;

    public int EXP;

    public ValueAddItem[] ValueAdd;

    public struct ValueAddItem
    {
        public int key;
        public int value;
    }

}

/// <summary>
/// Guidance 配置表
/// </summary>
public class CFG_GuildSkill : CFGTableBase<CFGItem_GuildSkill>
{
    private static CFG_GuildSkill _AllCFG;
    private static CFG_GuildSkill AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_GuildSkill, CFGItem_GuildSkill>();
            }
            return _AllCFG;
        }
    }

    private static Dictionary<int, List<CFGItem_GuildSkill>> _dicGuildSkill;

    public static Dictionary<int, List<CFGItem_GuildSkill>> DicGuildSkill
    {
        get
        {
            if (_dicGuildSkill == null)
            {
                _dicGuildSkill = new Dictionary<int, List<CFGItem_GuildSkill>>();
                foreach (var item in AllCFG.DicList)
                {
                    if (_dicGuildSkill.ContainsKey(item.SkillID))
                    {
                        if (_dicGuildSkill[item.SkillID] == null)
                        {
                            _dicGuildSkill[item.SkillID] = new List<CFGItem_GuildSkill>();
                        }
                        _dicGuildSkill[item.SkillID].Add(item);
                    }
                    else
                    {
                        _dicGuildSkill.Add(item.SkillID, new List<CFGItem_GuildSkill>());
                        _dicGuildSkill[item.SkillID].Add(item);
                    }
                }
            }
            return _dicGuildSkill;
        }
    }

    public static List<int> GetSkikllIdList()
    {
        return DicGuildSkill.Keys.ToList();
    }

    public static List<CFGItem_GuildSkill> GetSkillListByLevel(int skillLevel)
    {
        return AllCFG.DicList.FindAll(v => v.SkillLevel == skillLevel);
    }

    public static int GetMaxLevelById(int id)
    {
        return DicGuildSkill[id].Count;
    }

    public static bool IsMaxLevel(int id, int level)
    {
        return level >= GetMaxLevelById(id);
    }

    public static CFG_GuildSkill Get()
    {
        return AllCFG;
    }

    public static Dictionary<long, CFGItem_GuildSkill> GetDic()
    {
        return AllCFG.DicDatas;
    }

    public static CFGItem_GuildSkill GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static List<CFGItem_GuildSkill> GetList()
    {
        return AllCFG.DicList;
    }

    public static List<CFGItem_GuildSkill> GetSkillListById(int skillId)
    {
        return AllCFG.DicList.FindAll(x => x.SkillID == skillId);
    }

    public static CFGItem_GuildSkill GetSkillByIdAndLevel(int skillId, int level)
    {
        return AllCFG.DicList.Find(x => x.SkillID == skillId && x.SkillLevel == level);
    }

}
