using JEngine.Helper;

public class CFGItem_Pet : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// 名称
    /// </summary>
    public string Name;

    /// <summary>
    /// 品质
    /// </summary>
    public int Quality;

    /// <summary>
    /// 职业区分
    /// </summary>
    public int job;

    /// <summary>
    /// 图集名称
    /// </summary>
    public string Icon;

    /// <summary>
    /// 图片名称
    /// </summary>
    public string Iconname;

    public int Mark;

    ////宠物类别
    //public AttributeItem.Item[] PetKind;

    ///// <summary>
    ///// 宠物类别
    ///// </summary>
    //public int PetType {
    //    get
    //    {
    //        return PetKind[0].key;
    //    }
    //}

    public int ItemID;

    //进阶消耗
    public ConsumeItem[] StarConsume;

    //进阶限制
    public ConsumeItem[] StarLimit;


    public struct ConsumeItem
    {
        public int id;
        public int values;
    }

    public AttributeItem[] StarAttribute;
    public class AttributeItem
    {
        public class Item
        {
            public int key;
            public int value;
            public int gap;
            public int max;
        }
        public int id;
        public Item[] values;
    }

    public AttributeItem.Item[] GetStarAttributeByLv(int id)
    {
        for (int i = 0; i < StarAttribute.Length; i++)
        {
            if (StarAttribute[i].id == id)
                return StarAttribute[i].values;
        }
        return null;
    }

    public OfflineItem[] OfflineAward;

    public class OfflineItem
    {
        public class Item
        {
            public int id;
            public int num;
            public int gap;
            public int max;
        }
        public int key;
        public Item[] values;
    }

    public OfflineItem.Item[] GetOfflineAwardByStar(int idx)
    {
        if (OfflineAward == null)
            return null;
        for (int i = 0; i < OfflineAward.Length; i++)
        {
            if (OfflineAward[i].key == idx)
                return OfflineAward[i].values;
        }
        return null;
    }

    public int GetStatLimet(int id)
    {
        for (int i = 0; i < StarLimit.Length; i++)
        {
            if (StarLimit[i].id == id)
                return StarLimit[i].values;
        }
        return 0;
    }

    public int GetStarNeedAllPro(int id)
    {
        for (int i = 0; i < StarConsume.Length; i++)
        {
            if (StarConsume[i].id == id)
                return StarConsume[i].values;
        }
        return 1;
    }

    //效果描述
    public EffectDesItem[] EffectDes;

    public struct EffectDesItem
    {
        public int key;
        public string values;
    }
}

/// <summary>
/// Pet 配置表
/// </summary>
public class CFG_Pet : CFGTableBase<CFGItem_Pet>
{
    private static CFG_Pet _AllCFG;
    private static CFG_Pet AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Pet, CFGItem_Pet>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Pet Get()
    {
        return AllCFG;
    }

    public static CFGItem_Pet GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static CFGItem_Pet GetByItemId(int id)
    {
        for (int i = 0; i < AllCFG.DicList.Count; i++)
        {
            if (AllCFG.DicList[i].ItemID == id)
                return AllCFG.DicList[i];
        }
        return null;
    }

}
