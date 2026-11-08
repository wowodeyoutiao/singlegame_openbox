using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_EquipStrengthen : CFGItemBase
{
    public int UUID;

    public override long GetIdx()
    {
        return UUID;
    }


    public int Job;

    public int Strengthenlevel;

    public Consume_Item[] Consume;

    public struct Consume_Item
    {
        public int id;
        public int num;
    }

    public Attribute_Item[] Attribute;

    public struct Attribute_Item
    {
        public int slot;
        public int attribute;
        public int value;
    }
}

/// <summary>
/// EquipStrengthen 配置表
/// </summary>
public class CFG_EquipStrengthen : CFGTableBase<CFGItem_EquipStrengthen>
{

    private static CFG_EquipStrengthen _AllCFG;
    private static CFG_EquipStrengthen AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_EquipStrengthen, CFGItem_EquipStrengthen>();
            }
            return _AllCFG;
        }
    }
    public static CFG_EquipStrengthen Get()
    {
        return AllCFG;
    }

    public static CFGItem_EquipStrengthen GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static Dictionary<int, Dictionary<int, CFGItem_EquipStrengthen>> _JobStLvDic;
    public static Dictionary<int, Dictionary<int, CFGItem_EquipStrengthen>> JobStLvDic
    {
        get
        {
            if (_JobStLvDic == null)
            {
                _JobStLvDic = new Dictionary<int, Dictionary<int, CFGItem_EquipStrengthen>>();
                var _list = Get().DicList;
                for (int i = 0; i < _list.Count; i++)
                {
                    if (!_JobStLvDic.ContainsKey(_list[i].Job))
                    {
                        _JobStLvDic.Add(_list[i].Job, new Dictionary<int, CFGItem_EquipStrengthen>());
                    }
                    _JobStLvDic[_list[i].Job].Add(_list[i].Strengthenlevel, _list[i]);
                }
            }
            return _JobStLvDic;
        }
    }

    public static CFGItem_EquipStrengthen GetByJobAndLv(int job, int lv)
    {
        var _list = Get().DicList;

        if (JobStLvDic.ContainsKey(job) && JobStLvDic[job].ContainsKey(lv))
        {
            return JobStLvDic[job][lv];
        }

        for (int i = 0; i < _list.Count; i++)
        {
            if (_list[i].Job == job && _list[i].Strengthenlevel == lv)
                return _list[i];
        }
        return default;
    }

    public static CFGItem_EquipStrengthen.Attribute_Item GetAttrByJobAndLvSolt(int job, int lv, int solt)
    {
        var _list = JobStLvDic[job][lv];

        if (_list.Attribute != null)
        {
            foreach (var att in _list.Attribute)
            {
                if (att.slot == solt)
                {
                    return att;
                }
            }
        }
        return default;
    }

}
