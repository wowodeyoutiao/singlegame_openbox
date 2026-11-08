using JEngine.Helper;

/// <summary>
/// 
/// </summary>
[System.Serializable]
public class CFGItem_Actionset : CFGItemBase
{

    public int SkillID;

    //public int Lasttime;

    public string Music;

    public override long GetIdx()
    {
        return SkillID;
    }

}

/// <summary>
/// 配置表
/// </summary>
public class CFG_Actionset : CFGTableBase<CFGItem_Actionset>
{

    private static CFG_Actionset _AllCFG;
    private static CFG_Actionset AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Actionset, CFGItem_Actionset>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Actionset Get()
    {
        return AllCFG;
    }

    public static CFGItem_Actionset GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static bool ChealkId(int id)
    {
        return AllCFG.DicDatas.ContainsKey(id);
    }

}
