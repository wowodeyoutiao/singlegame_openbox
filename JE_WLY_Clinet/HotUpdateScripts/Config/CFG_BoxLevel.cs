using JEngine.Helper;

public class CFGItem_BoxLevel : CFGItemBase
{
    public int Level;

    public override long GetIdx()
    {
        return Level;
    }

    /// <summary>
    /// 需要的主角登记
    /// </summary>
    public int Rolelevel;

    /// <summary>
    /// 关卡需求
    /// </summary>
    public int BattleID;

    /// <summary>
    /// 升级时间
    /// </summary>
    public int Needtime;

    /// <summary>
    /// 升级消耗
    /// </summary>
    public ConsumeItem[] Consume;

    /// <summary>
    /// 最大开箱次数
    /// </summary>
    public int Maxtime;

    /// <summary>
    /// 品质随机
    /// </summary>
    public QualityRandomItem[] qualityrandom;


    public struct ConsumeItem
    {
        public int id;
        public int num;
    }

    public struct QualityRandomItem
    {
        public int id;
        public int rate;
    }
}

public class CFG_BoxLevel : CFGTableBase<CFGItem_BoxLevel>
{
    private static CFG_BoxLevel _AllCFG;
    private static CFG_BoxLevel AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_BoxLevel, CFGItem_BoxLevel>();
            }
            return _AllCFG;
        }
    }
    public static CFG_BoxLevel Get()
    {
        return AllCFG;
    }

    public static CFGItem_BoxLevel GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static int GetMaxLevel()
    {
        return AllCFG.DicList[AllCFG.DicList.Count - 1].Level;
    }

    public static CFGItem_BoxLevel GetMaxLevelItem()
    {
        return AllCFG.DicList[AllCFG.DicList.Count - 1];
    }
}