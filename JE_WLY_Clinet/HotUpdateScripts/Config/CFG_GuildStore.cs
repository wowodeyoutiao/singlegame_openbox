using JEngine.Helper;
using System.Collections.Generic;

/// <summary>
/// 引导
/// </summary>
public class CFGItem_GuildStore : CFGItemBase
{
    public int UUID;

    public override long GetIdx()
    {
        return (int)UUID;
    }

    public GoodsItem[] Goods;

    public GoodsItem[] Price;

    public class GoodsItem
    {
        public int id;
        public int num;
    }


    public int GuildLimit;

    public int PersonalLimit;

    public int Unlocklevel;

}

/// <summary>
/// Guidance 配置表
/// </summary>
public class CFG_GuildStore : CFGTableBase<CFGItem_GuildStore>
{
    private static CFG_GuildStore _AllCFG;
    private static CFG_GuildStore AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_GuildStore, CFGItem_GuildStore>();
            }
            return _AllCFG;
        }
    }

    public static CFG_GuildStore Get()
    {
        return AllCFG;
    }

    public static Dictionary<long, CFGItem_GuildStore> GetDic()
    {
        return AllCFG.DicDatas;
    }

    public static CFGItem_GuildStore GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static List<CFGItem_GuildStore> GetList()
    {
        return AllCFG.DicList;
    }
}
