using JEngine.Helper;
using UnityEngine;
public class CFGItem_PetSpeak : CFGItemBase
{
    public int id;

    public override long GetIdx()
    {
        return id;
    }

    public int type;
    public string text;
    public double weigh;
}

/// <summary>
/// EquipStar 配置表
/// </summary>
public class CFG_PetSpeak : CFGTableBase<CFGItem_PetSpeak>
{
    private static CFG_PetSpeak _AllCFG;

    private static double totalWeigh;

    private static CFG_PetSpeak AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_PetSpeak, CFGItem_PetSpeak>();
            }
            return _AllCFG;
        }
    }
    public static CFG_PetSpeak Get()
    {
        return AllCFG;
    }

    public static CFGItem_PetSpeak GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static void GetTotalWeigh()
    {
        if (totalWeigh != 0)
            return;

        for (int i = 0; i < AllCFG.DicList.Count; i++)
        {
            totalWeigh += AllCFG.DicList[i].weigh;
        }
    }

    public static CFGItem_PetSpeak GetOne()
    {
        GetTotalWeigh();

        var weigh = Random.Range(1, (float)totalWeigh + 1);
        double count = 0;
        for (int i = 0; i < AllCFG.DicList.Count; i++)
        {
            count += AllCFG.DicList[i].weigh;

            if (count > weigh)
                return AllCFG.DicList[i];
        }
        return AllCFG.DicList[0];
    }
}
