using JEngine.Helper;

public class CFGItem_TreasureBox : CFGItemBase
{
    public int UID;

    public override long GetIdx()
    {
        return UID;
    }

    /// <summary>
    /// 等级范围
    /// </summary>
    public int[] LevelRange;

    /// <summary>
    /// 抽取消耗
    /// </summary>
    public Consume_Item[] Consume;

    public struct Consume_Item
    {
        public int type;
        public int id;
        public int num;
    }


    /// <summary>
    /// 抽取消耗
    /// </summary>
    public Common_Item[] CommonPool;

    public struct Common_Item
    {
        public int position;
        public int id;
        public int num;
        public int priority;
    }

    /// <summary>
    /// 必中奖池配置
    /// </summary>
    public Common_Item[] SpecialPool;

    /// <summary>
    /// 必中次数配置
    /// </summary>
    public int[] Node;

    /// <summary>
    /// 最大抽取次数
    /// </summary>
    public int MaxTime;
}

/// <summary>
/// TreasureBox 配置表
/// </summary>
public class CFG_TreasureBox : CFGTableBase<CFGItem_TreasureBox>
{

    private static CFG_TreasureBox _AllCFG;
    private static CFG_TreasureBox AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_TreasureBox, CFGItem_TreasureBox>();
            }
            return _AllCFG;
        }
    }
    public static CFG_TreasureBox Get()
    {
        return AllCFG;
    }

    public static CFGItem_TreasureBox GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static CFGItem_TreasureBox GetBox(int lv)
    {
        foreach (var item in Get().DicDatas.Values)
        {
            if (item.LevelRange[0] <= lv && lv <= item.LevelRange[1])
                return item;
        }
        return null;
    }

}
