using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_TitleActivity : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }


    /// <summary>
    /// 共鸣属性
    /// </summary>
    public Dictionary<string, int>[] Attribute;

    /// <summary>
    /// 共鸣条件
    /// </summary>
    public string Limit;

    /// <summary>
    /// 共鸣星级
    /// </summary>
    public int StarLimit;

    /// <summary>
    /// 部位
    /// </summary>
    public int Slot;

}

/// <summary>
/// CFG_TitleActivity 配置表
/// </summary>
public class CFG_TitleActivity : CFGTableBase<CFGItem_TitleActivity>
{
    private static CFG_TitleActivity _AllCFG;
    private static CFG_TitleActivity AllCFG
    {
        get
        {
            if (_AllCFG == null)
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_TitleActivity, CFGItem_TitleActivity>();

            return _AllCFG;
        }

    }

    public static CFG_TitleActivity Get()
    {
        return AllCFG;
    }

    public static CFGItem_TitleActivity GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static Dictionary<int, List<CFGItem_TitleActivity>> TitleActivityListDic = new Dictionary<int, List<CFGItem_TitleActivity>>();

    public static List<CFGItem_TitleActivity> GetActiveListByType(int soltIdx)
    {
        if (!TitleActivityListDic.ContainsKey(soltIdx))
        {
            List<CFGItem_TitleActivity> _list = new List<CFGItem_TitleActivity>();
            var list = Get().DicList;
            for (int i = 0; i < list.Count; i++)
            {
                if (list[i].Slot == soltIdx)
                {
                    _list.Add(list[i]);
                }
            }
            TitleActivityListDic.Add(soltIdx, _list);
        }
        return TitleActivityListDic[soltIdx];
    }
}
