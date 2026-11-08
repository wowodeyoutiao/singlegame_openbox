using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_PushingWallSkill : CFGItemBase
{
    public int id;

    public override long GetIdx()
    {
        return id;
    }

    /// <summary>
    /// 技能名
    /// </summary>
    public string skillName;

    /// <summary>
    /// 组合
    /// </summary>
    public string combo;

    /// <summary>
    /// 力量值
    /// </summary>
    public int power;

    /// <summary>
    /// 效果（枚举）
    /// </summary>
    public int effect;

    /// <summary>
    /// 效果值
    /// </summary>
    public int effectNum;
}

/// <summary>
/// 推墙玩法技能配置表
/// </summary>
public class CFG_PushingWallSkill : CFGTableBase<CFGItem_PushingWallSkill>
{
    private static CFG_PushingWallSkill _AllCFG;
    private static CFG_PushingWallSkill AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_PushingWallSkill, CFGItem_PushingWallSkill>();
            }
            return _AllCFG;
        }
    }
    public static CFG_PushingWallSkill Get()
    {
        return AllCFG;
    }

    public static CFGItem_PushingWallSkill GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }
}
