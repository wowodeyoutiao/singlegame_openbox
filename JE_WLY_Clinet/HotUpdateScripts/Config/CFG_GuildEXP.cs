using JEngine.Helper;
using System.Collections.Generic;

/// <summary>
/// 引导
/// </summary>
public class CFGItem_GuildEXP : CFGItemBase
{
    public int GuildLevel;

    public override long GetIdx()
    {
        return GuildLevel;
    }

    public int NeedEXP;

    public AttributeAddItem[] AttributeAdd;

    public struct AttributeAddItem
    {
        public int key;
        public int value;
    }

    public string PrivilegeDes;

    public int WarehouseAdd;

    public int MemberNum;
}

/// <summary>
/// Guidance 配置表
/// </summary>
public class CFG_GuildEXP : CFGTableBase<CFGItem_GuildEXP>
{
    private static CFG_GuildEXP _AllCFG;
    private static CFG_GuildEXP AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_GuildEXP, CFGItem_GuildEXP>();
            }
            return _AllCFG;
        }
    }

    public static CFG_GuildEXP Get()
    {
        return AllCFG;
    }

    public static Dictionary<long, CFGItem_GuildEXP> GetDic()
    {
        return AllCFG.DicDatas;
    }

    public static CFGItem_GuildEXP GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static int MaxLevel()
    {
        return AllCFG.DicList[AllCFG.DicList.Count - 1].GuildLevel;
    }

    public static bool IsMaxLevel(int level)
    {
        return level >= MaxLevel();
    }

    public static List<CFGItem_GuildEXP> GetList()
    {
        return AllCFG.DicList;
    }

}
