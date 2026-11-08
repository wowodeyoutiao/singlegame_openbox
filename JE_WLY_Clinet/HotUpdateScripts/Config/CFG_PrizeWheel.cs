using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_PrizeWheel : CFGItemBase
{
    public int UID;

    public override long GetIdx()
    {
        return UID;
    }

    /// <summary>
    /// 奖励轮次
    /// </summary>
    public int Round;

    /// <summary>
    /// 充值额度
    /// </summary>
    public float rate;

}

/// <summary>
/// PrizeWheel 配置表
/// </summary>
public class CFG_PrizeWheel : CFGTableBase<CFGItem_PrizeWheel>
{

    private static CFG_PrizeWheel _AllCFG;
    private static CFG_PrizeWheel AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_PrizeWheel, CFGItem_PrizeWheel>();
            }
            return _AllCFG;
        }
    }
    public static CFG_PrizeWheel Get()
    {
        return AllCFG;
    }

    public static CFGItem_PrizeWheel GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static List<CFGItem_PrizeWheel> GetRounds(int id)
    {
        var _list = new List<CFGItem_PrizeWheel>();

        foreach (var item in Get().DicDatas)
        {
            if (item.Value.Round == id)
            {
                _list.Add(item.Value);
            }
        }
        _list.Sort((a, b) => { return a.UID.CompareTo(b.UID); });
        return _list;
    }

}
