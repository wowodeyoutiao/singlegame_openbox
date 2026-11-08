using JEngine.Helper;

public class CFGItem_Compose : CFGItemBase
{
    public int ComposeID;

    public override long GetIdx()
    {
        return ComposeID;
    }

    public int[] Raw;

    public class NeedGoldItem
    {
        public int id;
        public int num;
    }
    public NeedGoldItem NeedGold;

    public decimal NeedNum { get; set; }
    public decimal SuccessRate { get; set; }

}

/// <summary>
/// Compose 配置表
/// </summary>
public class CFG_Compose : CFGTableBase<CFGItem_Compose>
{

    private static CFG_Compose _AllCFG;
    private static CFG_Compose AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Compose, CFGItem_Compose>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Compose Get()
    {
        return AllCFG;
    }

    public static CFGItem_Compose GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static CFGItem_Compose GetByRawId(int id)
    {
        var AllList = Get().DicList;
        for (int i = 0; i < AllList.Count; i++)
        {
            var raws = AllList[i].Raw;
            for (int j = 0; j < raws.Length; j++)
            {
                if (id == raws[j])
                    return AllList[i];
            }
        }
        return null;
    }

}
