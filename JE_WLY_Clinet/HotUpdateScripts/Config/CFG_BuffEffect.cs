using JEngine.Helper;

public class CFGItem_BuffEffect : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// Buff类型
    /// </summary>
    public int Type;

    /// <summary>
    /// 说明
    /// </summary>
    public string Describe;

    /// <summary>
    /// Buff特效
    /// </summary>
    public string BuffEffect;

    /// <summary>
    /// 职业区分
    /// </summary>
    public int Job;

    /// <summary>
    /// Buff图集
    /// </summary>
    public string Gallery;

    /// <summary>
    /// Buff图标
    /// </summary>
    public string Icon;
}

/// <summary>
/// BuffEffect 配置表
/// </summary>
public class CFG_BuffEffect : CFGTableBase<CFGItem_BuffEffect>
{
    private static CFG_BuffEffect _AllCFG;
    private static CFG_BuffEffect AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_BuffEffect, CFGItem_BuffEffect>();
            }
            return _AllCFG;
        }
    }

    public static CFG_BuffEffect Get()
    {
        return AllCFG;
    }

    public static CFGItem_BuffEffect GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public enum BufferID
    {
        /// <summary>
        /// 被嘲讽
        /// </summary>
        CaoFeng = 25,
        /// <summary>
        /// 被眩晕
        /// </summary>
        XuanYun = 31,

        /// <summary>
        /// 冰冻
        /// </summary>
        BingDong = 15,

        /// <summary>
        /// 攻速移速
        /// </summary>
        AttackAndMoveAdd = 12,

        /// <summary>
        /// 毒蛊
        /// </summary>
        DuGu = 40,

        /// <summary>
        /// 毒砂术
        /// </summary>
        DuSha = 41,
    }

}
