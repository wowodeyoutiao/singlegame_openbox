using JEngine.Helper;
using System.Collections.Generic;
using System.Linq;

public class CFGItem_Partner : CFGItemBase
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

   
    private int _job = 0;

    public int Job {
        get { 
            if (_job == 0) 
                _job = job % 10;
            return _job;

        }
    }

    /// <summary>
    /// 图集名称
    /// </summary>
    public string Icon;

    /// <summary>
    /// 图片名称
    /// </summary>
    public string Iconname;

    public int ItemID;

    /// <summary>
    /// 技能列表{星级:技能ID}
    /// </summary>
    public Dictionary<string, int>[] Skills { get; set; }

    //对应碎片数量
    public StarConsumeItem.Item[] ChipNum;

    public int ChipGoodID { get {
           return ChipNum[0].key;
        } 
    }

    public int ChipGoodNum
    {
        get
        {
            return ChipNum[0].value;
        }
    }

    public StarConsumeItem[] StarConsume;

    public StarConsumeItem[] StarAttribute;
    public class StarConsumeItem
    {
        public class Item
        {
            public int key;
            public int value;
        }
        public int id;
        public Item[] values;
    }

    public StarConsumeItem.Item[] GetStarAttributeByLv(int id)
    {
        for (int i = 0; i < StarAttribute.Length; i++)
        {
            if (StarAttribute[i].id == id)
                return StarAttribute[i].values;
        }
        return null;
    }

    public StarConsumeItem.Item[] GetStarConsumeByLv(int id)
    {
        for (int i = 0; i < StarConsume.Length; i++)
        {
            if (StarConsume[i].id == id)
                return StarConsume[i].values;
        }
        return null;
    }

    /// <summary>
    /// 从0星升到指定星级需要消耗的所有道具的总和
    /// </summary>
    /// <param name="star">要升级到的星级</param>
    /// <returns></returns>
    public StarConsumeItem.Item[] GetStarAllConsumeFromLv0(int star)
    {
        List<StarConsumeItem.Item> items = new List<StarConsumeItem.Item>();
        foreach (var item in StarConsume)
        {
            if (item.id <= star)
            {
                foreach (var item2 in item.values)
                {
                    var _item = items.Find(v=>v.key == item2.key);
                    if (_item != null)
                        _item.value += item2.value;
                    else
                        items.Add(new StarConsumeItem.Item() { key = item2.key, value = item2.value });
                }
            }
        }
        return items.ToArray();
    }

    public int GetNextStarSkillIdByLv(int id)
    {
        if (Skills != null)
        {
            for (int i = 0; i < Skills.Length; i++)
            {
                if (Skills[i].ElementAt(0).Key == id.ToString())
                    return Skills[i].ElementAt(0).Value;
            }
        }
        return 0;
    }
}

/// <summary>
/// Partner 配置表
/// </summary>
public class CFG_Partner : CFGTableBase<CFGItem_Partner>
{
    private static CFG_Partner _AllCFG;
    private static CFG_Partner AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Partner, CFGItem_Partner>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Partner Get()
    {
        return AllCFG;
    }

    public static List<CFGItem_Partner> GetALlList()
    {
        return AllCFG.DicList;
    }

    public static CFGItem_Partner GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static CFGItem_Partner GetByItemID(int id)
    {
        for (int i = 0; i < AllCFG.DicList.Count; i++)
        {
            if (AllCFG.DicList[i].ChipGoodID == id)
                return AllCFG.DicList[i];
        }
        return null;
    }

    public static CFGItem_Partner GetByItemCfgID(int id)
    {
        for (int i = 0; i < AllCFG.DicList.Count; i++)
        {
            if (AllCFG.DicList[i].ItemID == id)
                return AllCFG.DicList[i];
        }
        return null;
    }

    public static List<CFGItem_Partner> GetByQuality(int quality)
    {
        var list = new List<CFGItem_Partner>();

        for (int i = 0; i < AllCFG.DicList.Count; i++)
        {
            if (AllCFG.DicList[i].Quality == quality)
                list.Add(AllCFG.DicList[i]);
        }

        return list;
    }
}
