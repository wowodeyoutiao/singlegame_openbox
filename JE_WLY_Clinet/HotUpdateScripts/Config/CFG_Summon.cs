using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_Summon : CFGItemBase
{
    public int UID;

    public override long GetIdx()
    {
        return UID;
    }

    /// <summary>
    /// 抽取模块
    /// </summary>
    public int DrawModule;

    /// <summary>
    /// 抽取类别
    /// </summary>
    public int Drawtype;

    /// <summary>
    /// 等级范围
    /// </summary>
    public int[] LevelRange;

    /// <summary>
    /// 等级范围
    /// </summary>
    public int[] Node;

    /// <summary>
    /// 最大抽取次数
    /// </summary>
    public decimal MaxTime;

    /// <summary>
    /// 抽取消耗
    /// </summary>
    public ConsumeItem[] Consume;

    public class ConsumeItem
    {
        public int type;
        public int id;
        public int num;
    }

}

/// <summary>
/// Summon 配置表
/// </summary>
public class CFG_Summon : CFGTableBase<CFGItem_Summon>
{
    private static CFG_Summon _AllCFG;
    private static CFG_Summon AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Summon, CFGItem_Summon>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Summon Get()
    {
        return AllCFG;
    }

    public static CFGItem_Summon GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static List<CFGItem_Summon> PartenrGetList;

    //根据玩家等级获取可抽取的伙伴
    public static CFGItem_Summon[] GetPartenrItems()
    {
        if (PartenrGetList == null)
        {
            PartenrGetList = new List<CFGItem_Summon>();
            foreach (var item in AllCFG.DicDatas)
            {
                if (item.Value.DrawModule == 1)
                {
                    PartenrGetList.Add(item.Value);
                }
            }
        }
        CFGItem_Summon[] baseSummon = new CFGItem_Summon[3];
        for (int i = 0; i < PartenrGetList.Count; i++)
        {
            if (PartenrGetList[i].LevelRange[0] <= RoleM.Instance.Level && PartenrGetList[i].LevelRange[1] >= RoleM.Instance.Level)
            {
                baseSummon[(PartenrGetList[i].Drawtype - 1)] = PartenrGetList[i];
            }
        }
        return baseSummon;
    }

    private static List<CFGItem_Summon> PetGetList;

    //根据玩家等级获取可抽取的宠物
    public static CFGItem_Summon[] GetPetItems()
    {
        if (PetGetList == null)
        {
            PetGetList = new List<CFGItem_Summon>();
            foreach (var item in AllCFG.DicDatas)
            {
                if (item.Value.DrawModule == 2)
                {
                    PetGetList.Add(item.Value);
                }
            }
        }
        CFGItem_Summon[] baseSummon = new CFGItem_Summon[4];
        for (int i = 0; i < PetGetList.Count; i++)
        {
            if (PetGetList[i].LevelRange[0] <= RoleM.Instance.Level && PetGetList[i].LevelRange[1] >= RoleM.Instance.Level)
            {
                baseSummon[(PetGetList[i].Drawtype - 1)] = PetGetList[i];
            }
        }
        return baseSummon;
    }

    public static CFGItem_Summon GetPetExchangeCard()
    {
        var cfgItem = GetById(55);
        return cfgItem;
    }
}
