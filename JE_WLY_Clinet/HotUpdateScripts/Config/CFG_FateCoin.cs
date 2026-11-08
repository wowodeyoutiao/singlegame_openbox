using JEngine.Helper;

public class CFGItem_FateCoin : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// 事件描述
    /// </summary>
    public string Des;

    /// <summary>
    /// 图集
    /// </summary>
    public string Gallery;

    /// <summary>
    /// 图标
    /// </summary>
    public string Icon;

    /// <summary>
    /// 作用类型 1:正 2:反 3:无
    /// </summary>
    public int EffectType;

}

/// <summary>
/// 幸运金币 配置表
/// </summary>
public class CFG_FateCoin : CFGTableBase<CFGItem_FateCoin>
{

    private static CFG_FateCoin _AllCFG;
    private static CFG_FateCoin AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_FateCoin, CFGItem_FateCoin>();
            }
            return _AllCFG;
        }
    }
    public static CFG_FateCoin Get()
    {
        return AllCFG;
    }

    public static CFGItem_FateCoin GetById(int id)
    {
        if (!AllCFG.DicDatas.ContainsKey(id))
            return null;
        return AllCFG.DicDatas[id];
    }


}
