using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_Title : CFGItemBase
{
    public int UID;

    public override long GetIdx()
    {
        return UID;
    }

    public struct JosnItem
    {
        public int key;
        public int value;
    }
    /// <summary>
    /// 激活属性
    /// </summary>
    public JosnItem2[] Goods;

    public JosnItem[] GetGoodItemsByStar(int star)
    {
        for (int i = 0; i < Goods.Length; i++)
        {
            if (Goods[i].id == star)
                return Goods[i].values;
        }
        return null;
    }

    /// <summary>
    /// 佩戴属性
    /// </summary>
    public JosnItem[] Limit;

    public struct JosnItem2
    {
        public int id;
        public JosnItem[] values;
    }

    /// <summary>
    /// 激活消耗
    /// </summary>
    public JosnItem2[] Activity;


    /// <summary>
    /// 升星几率
    /// </summary>
    public JosnItem[] StarRate;

    /// <summary>
    /// 升星消耗
    /// </summary>
    public JosnItem2[] StarConsume;

    /// <summary>
    /// 称号名称
    /// </summary>
    public string Name;

    /// <summary>
    /// 时效
    /// </summary>
    public int Period;

    public int Slot;

    /// <summary>
    /// 图标路径
    /// </summary>
    public string Path;

    public int Icon;
}

/// <summary>
/// CFG_Title 配置表
/// </summary>
public class CFG_Title : CFGTableBase<CFGItem_Title>
{
    private static CFG_Title _AllCFG;
    private static CFG_Title AllCFG
    {
        get
        {
            if (_AllCFG == null)
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Title, CFGItem_Title>();

            return _AllCFG;
        }

    }

    public static CFG_Title Get()
    {
        return AllCFG;
    }

    public static CFGItem_Title GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static Dictionary<int, List<CFGItem_Title>> TitleListDic = new Dictionary<int, List<CFGItem_Title>>();

    public static List<CFGItem_Title> GetActiveListByType(int soltIdx)
    {
        if (!TitleListDic.ContainsKey(soltIdx))
        {
            List<CFGItem_Title> _list = new List<CFGItem_Title>();
            var list = Get().DicList;
            for (int i = 0; i < list.Count; i++)
            {
                if (list[i].Slot == soltIdx)
                {
                    _list.Add(list[i]);
                }
            }
            TitleListDic.Add(soltIdx, _list);
        }
        return TitleListDic[soltIdx];
    }
}
