using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_EquipMaster : CFGItemBase
{
    public int UUID;

    public override long GetIdx()
    {
        return UUID;
    }

    public int Category;

    public int Limit;


    public AttributeItem[] Attribute;

    public struct AttributeItem
    {
        public int key;
        public int value;
    }

}

/// <summary>
/// EquipMaster 配置表
/// </summary>
public class CFG_EquipMaster : CFGTableBase<CFGItem_EquipMaster>
{
    private static CFG_EquipMaster _AllCFG;
    private static CFG_EquipMaster AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_EquipMaster, CFGItem_EquipMaster>();
            }
            return _AllCFG;
        }
    }
    public static CFG_EquipMaster Get()
    {
        return AllCFG;
    }

    public static CFGItem_EquipMaster GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static List<CFGItem_EquipMaster> GetMastersByCatedory(CategoryClass cClass)
    {
        List<CFGItem_EquipMaster> masters = new List<CFGItem_EquipMaster>();
        foreach (var item in Get().DicDatas.Values)
        {
            if (item.Category == (int)cClass)
            {
                masters.Add(item);
            }
        }
        return masters;
    }

    public enum CategoryClass
    {
        //强化大师
        Reinforce = 1,

        //升星大师
        UpStar = 2,

        /// <summary>
        /// 魂石大师
        /// </summary>
        InsetMonster = 3,

        /// <summary>
        /// 自然大师
        /// </summary>
        NatureManster = 4,
    }

}
