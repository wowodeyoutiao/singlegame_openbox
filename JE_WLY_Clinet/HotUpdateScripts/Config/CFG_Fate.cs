using JEngine.Helper;

public class CFGItem_Fate : CFGItemBase
{
    public int FateLevel;

    public override long GetIdx()
    {
        return FateLevel;
    }

    /// <summary>
    /// 重数显示
    /// </summary>
    public string Multiplicity;

    /// <summary>
    /// 级数显示
    /// </summary>
    public string Level;

    /// <summary>
    /// 升级经验
    /// </summary>
    public int FateExp;


    public class UpValueItem
    {
        public int key;
        public int value;
        public int job;
    }

    /// <summary>
    /// 升级属性
    /// </summary>
    public UpValueItem[] UpValue;


    public class AdvancedNeedItem
    {
        public int id;
        public int num;
    }
    /// <summary>
    /// 进阶消耗
    /// </summary>
    public AdvancedNeedItem[] AdvancedNeed;

    public UpValueItem[] AdvancedValue;
}

/// <summary>
/// Fate 配置表
/// </summary>
public class CFG_Fate : CFGTableBase<CFGItem_Fate>
{

    private static CFG_Fate _AllCFG;
    private static CFG_Fate AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Fate, CFGItem_Fate>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Fate Get()
    {
        return AllCFG;
    }

    public static CFGItem_Fate GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    //升级
    public static CFGItem_Fate.UpValueItem GetUpLvByLvAndJob(int lv, int job)
    {
        if (lv > Get().DicList.Count)
        {
            return null;
        }
        var info = GetById(lv).UpValue;
        if (info != null)
        {
            for (int i = 0; i < info.Length; i++)
            {
                if (info[i].job == job)
                    return info[i];
            }
        }
        return null;
    }

    public static int GetAllExp(int lv)
    {
        int allExp = 0;

        int baseNum = lv / 10;
        if (lv % 10 == 0)
        {
            baseNum--;
        }
        int min = baseNum * 10 + 1;
        int max = baseNum * 10 + 10;

        for (int i = min; i <= max; i++)
        {
            allExp += GetById(i).FateExp;
        }
        return allExp;
    }

    //升阶
    public static CFGItem_Fate.UpValueItem GetUplvLastValueByLv(int lv)
    {
        CFGItem_Fate.UpValueItem info = null;
        for (int i = 0; i < Get().DicList.Count; i++)
        {
            if (Get().DicList[i].FateLevel < lv && Get().DicList[i].AdvancedValue != null)
                info = Get().DicList[i].AdvancedValue[0];
        }
        return info;
    }

}
