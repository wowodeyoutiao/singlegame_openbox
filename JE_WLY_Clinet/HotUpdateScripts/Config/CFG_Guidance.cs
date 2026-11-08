using JEngine.Helper;
using System.Collections.Generic;


// 触发条件类型
public enum GuidanceTriggerType
{
    None = 0, // 无
    Level = 1, // 等级
    Task = 2, // 任务
    MainBattle = 3, // 主线
    Battle = 4, // 战役
}

// 引导行为类型
public enum GuidanceGuideType
{
    None = 0, // 无
    Text = 1, // 提示文本
    Click = 2, // 点击UI
    Animator = 3, // 引导动画
    TextAndClick = 4, // 点击+文本
    TextAndAnimator = 5, // 动画+文本
    All = 6, // 点击+动画+文本
    BattleETCJoystick = 7, //战役摇杆
}

// 完成条件类型
public enum GuidanceCompleteType
{
    None = 0, // 无
    TimeOut = 1, // 时间到
    Clicked = 2, // 点击完成
    Task = 3, // 任务完成
    MainBattle = 4, // 引导动画
    Battle = 5, // 战役
    ReceiveServer = 6, // 服务器
    Drag = 7, // 滑动手势
}


/// <summary>
/// 引导
/// </summary>
public class CFGItem_Guidance : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return (int)ID;
    }

    public int Stage;
    public string Describe;
    public int TriggerType;
    public int TriggerValue;
    public int GuideType;
    public int GuideParam;
    public string ClassParam;
    public string TargetParam;
    public int CompleteType;
    public int CompleteValue;
    public int DelayTime;
    public int JumpToId;
    public int NoticeServer;
    public int ReconnectJumpToStage;
    public int[] ReconnectSteps;
    public int HalfMask;
}

/// <summary>
/// Guidance 配置表
/// </summary>
public class CFG_Guidance : CFGTableBase<CFGItem_Guidance>
{
    private static CFG_Guidance _AllCFG;
    private static CFG_Guidance AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Guidance, CFGItem_Guidance>();
            }
            return _AllCFG;
        }
    }

    public static CFG_Guidance Get()
    {
        return AllCFG;
    }

    public static Dictionary<long, CFGItem_Guidance> GetDic()
    {
        return AllCFG.DicDatas;
    }

    public static CFGItem_Guidance GetById(int id)
    {
        if (AllCFG.DicDatas.ContainsKey(id))
            return AllCFG.DicDatas[id];

        return null;
    }

    public static int GetType(int id)
    {
        return AllCFG.DicDatas[id].Stage;
    }

    public static CFGItem_Guidance GetFirstByStage(int stage)
    {
        for (int i = 0; i < AllCFG.DicList.Count; i++)
        {
            if (AllCFG.DicList[i].Stage == (int)stage)
            {
                return AllCFG.DicList[i];
            }
        }
        return null;
    }

    /// <summary>
    /// 最后一个引导ID
    /// </summary>
    /// <returns></returns>
    public static int LastGuideId()
    {
        return AllCFG.DicList.Count;
    }
}
