using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_Sign : CFGItemBase
{
    public int UUID;

    public override long GetIdx()
    {
        return UUID;
    }

    /// <summary>
    /// 奖励轮数
    /// </summary>
    public int Round;

    /// <summary>
    /// 天数
    /// </summary>
    public int Days;

    /// <summary>
    /// 双倍VIP等级
    /// </summary>
    public int VIPLimit;

    /// <summary>
    /// 奖励内容
    /// </summary>
    public Accumulate_Award[] Award;

    /// <summary>
    /// 累签奖励
    /// </summary>
    public Accumulate_Award[] AccumulateAward;

    public class Accumulate_Award
    {
        public int id;
        public int num;
    }
}

/// <summary>
/// Sign 配置表
/// </summary>
public class CFG_Sign : CFGTableBase<CFGItem_Sign>
{
    private static CFG_Sign _AllCFG;
    private static CFG_Sign AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Sign, CFGItem_Sign>();
            }
            return _AllCFG;
        }
    }

    public static CFG_Sign Get()
    {
        return AllCFG;
    }

    public static CFGItem_Sign GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static (List<CFGItem_Sign>, List<CFGItem_Sign>) GetSignAndAccumulateListByRound(int round)
    {
        List<CFGItem_Sign> _list = new List<CFGItem_Sign>();
        List<CFGItem_Sign> _AccumulateList = new List<CFGItem_Sign>();
        foreach (var item in Get().DicDatas.Values)
        {
            if (item.Round == round)
            {
                _list.Add(item);
                if (item.AccumulateAward != null && item.AccumulateAward.Length > 0)
                {
                    _AccumulateList.Add(item);
                }
            }
        }
        return (_list, _AccumulateList);
    }

}
