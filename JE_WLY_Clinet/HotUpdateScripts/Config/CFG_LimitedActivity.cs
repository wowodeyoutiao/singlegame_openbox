using JEngine.Helper;

public class CFGItem_LimitedActivity : CFGItemBase
{
    public int UID;

    public override long GetIdx()
    {
        return UID;
    }

    /// <summary>
    /// 奖励预览
    /// </summary>
    public int[] GiftPre;

    /// <summary>
    /// 功能名称
    /// </summary>
    public string Name;

    /// <summary>
    /// 活动介绍
    /// </summary>
    public string Des;
}

/// <summary>
/// CFG_LimitedActivity 配置表
/// </summary>
public class CFG_LimitedActivity : CFGTableBase<CFGItem_LimitedActivity>
{
    private static CFG_LimitedActivity _AllCFG;
    private static CFG_LimitedActivity AllCFG
    {
        get
        {
            if (_AllCFG == null)
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_LimitedActivity, CFGItem_LimitedActivity>();

            return _AllCFG;
        }

    }

    public static CFG_LimitedActivity Get()
    {
        return AllCFG;
    }

    public static CFGItem_LimitedActivity GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public enum Type
    {
        WorldBoss = 1,
        FightEnd,
        QXZC,
    }
}
