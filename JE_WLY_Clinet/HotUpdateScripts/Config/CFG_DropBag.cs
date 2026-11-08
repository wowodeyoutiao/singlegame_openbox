using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_DropBag : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    public int bagId;
    public int ifBag;
    public int itemId;
    public string itemName;
    public int itemNum;
    public int show;
}

public class CFG_DropBag : CFGTableBase<CFGItem_DropBag>
{
    private static CFG_DropBag _AllCFG;
    private static List<CFGItem_DropBag> selectItems = new List<CFGItem_DropBag>();

    private static CFG_DropBag AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_DropBag, CFGItem_DropBag>();
            }
            return _AllCFG;
        }
    }
    public static CFG_DropBag Get()
    {
        return AllCFG;
    }

    public static CFGItem_DropBag GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static List<CFGItem_DropBag> GetCustomSelectList()
    {
        if (selectItems != null && selectItems.Count > 0)
            return selectItems;

        selectItems.Clear();
        var list = AllCFG.DicList;
        for (int i = 0; i < list.Count; i++)
        {
            if (list[i].bagId == 100002 && list[i].ifBag == 0)
                selectItems.Add(list[i]);
        }

        selectItems.Sort((a, b) => { return a.ID.CompareTo(b.ID); });

        return selectItems;
    }
}