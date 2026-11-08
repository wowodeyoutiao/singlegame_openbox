using JEngine.Helper;

[System.Serializable]
public class CFGItemUse_Item : CFGItemBase
{
    /// <summary>
    /// 道具ID
    /// </summary>
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// 名字
    /// </summary>
    public string Name;

    public int Type;

    /// <summary>
    /// 每日使用上限
    /// </summary>
    public int UseLimit;

    public CFGItemUs_Givegoods[] Givegoods;
    /// <summary>
    /// 商品内容
    /// </summary>
    public class CFGItemUs_Givegoods
    {
        public int id;
        public int num;
        public int rateval;
    }

    public CFGItemUs_BuffEffect[] Buffeffect;
    /// <summary>
    /// 商品内容
    /// </summary>
    public class CFGItemUs_BuffEffect
    {
        public int key;
        public int id;
        public int num;
    }
}

/// <summary>
/// Item 配置表
/// </summary>
public class CFG_ItemUse : CFGTableBase<CFGItemUse_Item>
{

    private static CFG_ItemUse _AllCFG;
    private static CFG_ItemUse AllCFG
    { 
        get 
        { 
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_ItemUse, CFGItemUse_Item>();
            }
            return _AllCFG;
        } 
    }
    public static CFG_ItemUse Get()
    {
        return AllCFG;
    }

    public static CFGItemUse_Item GetById(int id)
    {
        if (AllCFG.DicDatas.TryGetValue(id, out var itemCfg))
            return itemCfg;

        return null;
    }

    /// <summary>
    /// 是否有这个道具的配置
    /// </summary>
    /// <param name="cfgId">道具id</param>
    /// <returns></returns>
    public static bool HasItem(int cfgId)
    {
        return AllCFG.DicDatas.ContainsKey(cfgId);
    }

    /// <summary>
    /// 道具每日使用上线
    /// </summary>
    /// <param name="cfgId">道具id</param>
    /// <returns></returns>
    public static int GetDailyUseLimit(int cfgId)
    {
        if (AllCFG.DicDatas.ContainsKey(cfgId))
            return AllCFG.DicDatas[cfgId].UseLimit;
        return 0;
    }

    public static int GetStarTalismanIdByAttr(long num)
    {
        for (int i = 25; i < 29; i++)
        {
            var cfg = GetById(i);
            if (cfg == null || cfg.Buffeffect == null || cfg.Buffeffect.Length <= 0)
                continue;

            for (int j = 0; j < cfg.Buffeffect.Length; j++)
            {
                if (cfg.Buffeffect[j].id == 14 && cfg.Buffeffect[j].num == num)
                    return i;
            }
        }
        return 0;
    }
}
