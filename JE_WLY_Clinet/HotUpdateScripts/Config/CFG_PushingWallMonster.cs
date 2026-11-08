using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_PushingWallMonster : CFGItemBase
{
    public int id;

    public override long GetIdx()
    {
        return id;
    }

    /// <summary>
    /// 难度
    /// </summary>
    public int level;

    /// <summary>
    /// 初始力量值
    /// </summary>
    public int power;

    /// <summary>
    /// 间隔值
    /// </summary>
    public int interval;

    /// <summary>
    /// 间隔增加值
    /// </summary>
    public int increaseValue;

    /// <summary>
    /// 奖励
    /// </summary>
    public Dictionary<string, int>[] Reward;

    /// <summary>
    /// 怪物数量
    /// </summary>
    public int monsterNum;

    /// <summary>
    /// 怪物头像
    /// </summary>
    public List<int> icon;
}

/// <summary>
/// VIP 配置表
/// </summary>
public class CFG_PushingWallMonster : CFGTableBase<CFGItem_PushingWallMonster>
{
    private static CFG_PushingWallMonster _AllCFG;
    private static CFG_PushingWallMonster AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_PushingWallMonster, CFGItem_PushingWallMonster>();
            }
            return _AllCFG;
        }
    }
    public static CFG_PushingWallMonster Get()
    {
        return AllCFG;
    }

    public static CFGItem_PushingWallMonster GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static bool IsMaxVipLevel(int lv)
    {
        return Get().DicList.Count <= lv;
    }

    public static int GetMaxVipLevel()
    {
        return Get().DicList.Count;
    }
}
