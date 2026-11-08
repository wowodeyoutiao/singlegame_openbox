using JEngine.Helper;

public class CFGItem_SkillEffect : CFGItemBase
{
    public int TargetID;

    public override long GetIdx()
    {
        return TargetID;
    }

    /// <summary>
    /// 起手特效
    /// </summary>
    public string StartEffect;

    /// <summary>
    /// 目标受击特效
    /// </summary>
    public string HittedEffect;

    /// <summary>
    /// 技能弹道特效
    /// </summary>
    public string BallisticEffect;

    /// <summary>
    /// 目标死亡特效
    /// </summary>
    public string DeadEffect;

    /// <summary>
    /// 目标身体特效
    /// </summary>
    public string BodyEffect;

    /// <summary>
    /// 起手时间延迟
    /// </summary>
    public decimal StartDelay;

    /// <summary>
    /// 起手持续时间
    /// </summary>
    public decimal StartRemainTime;

    /// <summary>
    /// 弹道时间延迟
    /// </summary>
    public decimal BallisticDelay;

    /// <summary>
    /// 弹道持续时间
    /// </summary>
    public decimal BallisticRemainTime;

    /// <summary>
    /// 目标时间延迟
    /// </summary>
    public decimal BodyDelay;

    /// <summary>
    /// 目标持续时间
    /// </summary>
    public decimal BodyRemainTime;

    public decimal isAmendStart;



    /// <summary>
    /// 模型常在特效
    /// </summary>
    public string ResidentEffect;

    /// <summary>
    /// 模型常在特效
    /// </summary>
    public string MoveEffect;


    public int[] Position;
    public int[] BallisticPosition;
    public int[] Scale;
    public float MovePre;
    public int PosType;
    public int MainStartDelay;
    public int MainBodyDelay;
    public int MainBallisticDelay;

}

/// <summary>
/// SkillEffect 配置表
/// </summary>
public class CFG_SkillEffect : CFGTableBase<CFGItem_SkillEffect>
{

    private static CFG_SkillEffect _AllCFG;
    private static CFG_SkillEffect AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_SkillEffect, CFGItem_SkillEffect>();
            }
            return _AllCFG;
        }
    }

    public static CFG_SkillEffect Get()
    {
        return AllCFG;
    }

    public static CFGItem_SkillEffect GetById(int id)
    {
        if (!AllCFG.DicDatas.ContainsKey(id))
            return null;
        return AllCFG.DicDatas[id];
    }

    public static string GetBallisticEffect(int id)
    {
        return GetById(id).BallisticEffect;
    }


    public enum SkillEffectType
    {
        //普通
        Normall = 1,
        //职业，
        Job = 2,
        //特别
        Super = 3,
    }


}
