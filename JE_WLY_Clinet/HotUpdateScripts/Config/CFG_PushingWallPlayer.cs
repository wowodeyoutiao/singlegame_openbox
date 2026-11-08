using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_PushingWallPlayer : CFGItemBase
{
    public int id;

    public override long GetIdx()
    {
        return id;
    }

    /// <summary>
    /// VIP名称
    /// </summary>
    public int skillpower;

    /// <summary>
    /// 力量值
    /// </summary>
    public int power;
}

/// <summary>
/// VIP 配置表
/// </summary>
public class CFG_PushingWallPlayer : CFGTableBase<CFGItem_PushingWallPlayer>
{
    private static CFG_PushingWallPlayer _AllCFG;
    private static CFG_PushingWallPlayer AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_PushingWallPlayer, CFGItem_PushingWallPlayer>();
            }
            return _AllCFG;
        }
    }
    public static CFG_PushingWallPlayer Get()
    {
        return AllCFG;
    }

    public static CFGItem_PushingWallPlayer GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }
}
