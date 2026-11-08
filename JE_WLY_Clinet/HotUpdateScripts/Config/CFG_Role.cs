using JEngine.Helper;
using System.Collections.Generic;
using System.Linq;

public class CFGItem_Role : CFGItemBase
{
    public int UID;
    public int Level;
    public long NeedEXP;
    public int job;

    public Dictionary<string, int>[] WarriorBassProperty;

    private PropertyItem[] _Property;
    public PropertyItem[] Property
    {
        get
        {
            if (_Property == null)
            {
                _Property = new PropertyItem[WarriorBassProperty.Length];
                for (int i = 0; WarriorBassProperty.Length > i; i++)
                {
                    var WarriorBassPropertyItem = WarriorBassProperty[i].ElementAt(0);
                    _Property[i] = new PropertyItem
                    {
                        Id = int.Parse(WarriorBassPropertyItem.Key),
                        Num = WarriorBassPropertyItem.Value,
                    };
                }
            }
            return _Property;
        }
    }

    public struct PropertyItem
    {
        public int Id;
        public int Num;
    }
    public override long GetIdx()
    {
        return UID;
    }

}

/// <summary>
/// Role 配置表
/// </summary>
public class CFG_Role : CFGTableBase<CFGItem_Role>
{

    private static CFG_Role _AllCFG;
    private static CFG_Role AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Role, CFGItem_Role>();
            }
            return _AllCFG;
        }
    }

    public static CFG_Role Get()
    {
        return AllCFG;
    }

    private static Dictionary<int, Dictionary<int, CFGItem_Role>> _JobLevelDic;

    private static Dictionary<int, Dictionary<int, CFGItem_Role>> JobLevelDic
    {
        get
        {
            if (_JobLevelDic == null)
            {
                _JobLevelDic = new Dictionary<int, Dictionary<int, CFGItem_Role>>();
                foreach (var item in AllCFG.DicDatas)
                {
                    if (!_JobLevelDic.ContainsKey(item.Value.job))
                    {
                        _JobLevelDic.Add(item.Value.job, new Dictionary<int, CFGItem_Role>());
                    }
                    _JobLevelDic[item.Value.job].Add(item.Value.Level, item.Value);
                }
            }
            return _JobLevelDic;
        }
    }

    public static CFGItem_Role GetByJobAndLevel(int Job, int Level)
    {
        return JobLevelDic[Job][Level];
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="Job">职业</param>
    /// <param name="level">等级</param>
    /// <returns>从1级升级到目标等级要消耗的所有经验</returns>
    public static long GetAllNeedExpToUpLevel(int Job, int level)
    {
        var levelDic = JobLevelDic[Job];
        long exp = 0;
        foreach (var item in levelDic)
        {
            if (item.Key < level)
                exp += item.Value.NeedEXP;
        }
        return exp;
    }

    public static int MaxLvByJob(int Job)
    {
        return JobLevelDic[Job].Count;
    }

    public static CFGItem_Role GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

}
