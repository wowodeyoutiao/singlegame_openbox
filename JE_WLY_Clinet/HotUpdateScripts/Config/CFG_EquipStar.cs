using JEngine.Helper;

public class CFGItem_EquipStar : CFGItemBase
{
    public int UUID;

    public override long GetIdx()
    {
        return UUID;
    }


    public int Job;

    public int StarLevel;

    public int SuccessRate;

    public Material_Item[] MaterialID;

    public struct Material_Item
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
/// EquipStar 配置表
/// </summary>
public class CFG_EquipStar : CFGTableBase<CFGItem_EquipStar>
{

    private static CFG_EquipStar _AllCFG;
    private static CFG_EquipStar AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_EquipStar, CFGItem_EquipStar>();
            }
            return _AllCFG;
        }
    }
    public static CFG_EquipStar Get()
    {
        return AllCFG;
    }

    public static CFGItem_EquipStar GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static CFGItem_EquipStar GetByJobAndLv(int job, int lv)
    {
        foreach (var item in Get().DicDatas.Values)
        {
            if (item.Job == job && item.StarLevel == lv)
                return item;
        }
        return default;
    }

    public static CFGItem_EquipStar.Attribute_Item GetAttrByJobAndLvSolt(int job, int lv, int solt)
    {
        foreach (var item in Get().DicDatas.Values)
        {
            if (item.Job == job && item.StarLevel == lv && item.Attribute != null)
            {
                foreach (var att in item.Attribute)
                {
                    if (att.slot == solt)
                        return att;
                }
            }
        }
        return default;
    }


    public static int GetMaxStar(int job)
    {
        int Count = 0;
        foreach (var item in Get().DicDatas)
        {
            if (item.Value.Job == job)
                Count++;
        }
        return Count - 1;
    }

}
