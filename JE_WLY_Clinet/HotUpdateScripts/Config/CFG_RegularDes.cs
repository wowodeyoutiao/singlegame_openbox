using JEngine.Helper;

public class CFGItem_RegularDes : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    public string describe;

}

/// <summary>
/// RegularDes 配置表
/// </summary>
public class CFG_RegularDes : CFGTableBase<CFGItem_RegularDes>
{

    private static CFG_RegularDes _AllCFG;
    private static CFG_RegularDes AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_RegularDes, CFGItem_RegularDes>();
            }
            return _AllCFG;
        }
    }
    public static CFG_RegularDes Get()
    {
        return AllCFG;
    }

    public static CFGItem_RegularDes GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

}
