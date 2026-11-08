using JEngine.Helper;
using System.Collections.Generic;
public class CFGItem_EndlessTrials : CFGItemBase
{
    public int UID;

    public override long GetIdx()
    {
        return UID;
    }

    /// <summary>
    /// 等级限制
    /// </summary>
    public int Levellimit;

    /// <summary>
    /// 闯关层数
    /// </summary>
    public int Layers;

    /// <summary>
    /// 挑战时间
    /// </summary>
    public int Timelimit;

    /// <summary>
    /// 结算奖励
    /// </summary>
    public SetAward[] Settlement;

    /// <summary>
    /// 产出预览
    /// </summary>
    public int[] Preview;

    public class SetAward
    {
        public int id;
        public int num;
    }

    /// <summary>
    /// 离线收益效率
    /// </summary>
    public SetAward[] Offlineraward;

}

/// <summary>
/// EndlessTrials 配置表
/// </summary>
public class CFG_EndlessTrials : CFGTableBase<CFGItem_EndlessTrials>
{

    private static CFG_EndlessTrials _AllCFG;
    private static CFG_EndlessTrials AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_EndlessTrials, CFGItem_EndlessTrials>();
            }
            return _AllCFG;
        }
    }
    public static CFG_EndlessTrials Get()
    {
        return AllCFG;
    }

    public static CFGItem_EndlessTrials GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static Dictionary<int, CFGItem_EndlessTrials> endlessDic = new Dictionary<int, CFGItem_EndlessTrials>();

    public static CFGItem_EndlessTrials GetByLayer(int layer)
    {
        if (!endlessDic.ContainsKey(layer))
        {
            foreach (var item in Get().DicDatas.Values)
            {
                endlessDic.Add(item.Layers, item);
            }
        }
        return endlessDic[layer];
    }
}
