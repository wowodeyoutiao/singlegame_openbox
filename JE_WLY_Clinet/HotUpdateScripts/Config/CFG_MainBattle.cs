using JEngine.Helper;
using System.Collections.Generic;

public enum MainBattleFlayType
{
    None = 0,
    Elite = 1,
    Boss = 2
}

public class CFGItem_MainBattle : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    // 关卡名
    public string name;
    // 标识
    public int flag; // 0 普通，1 精英，2 boss
    // 消耗体力
    public int cost;
    // 首通奖励
    public Dictionary<string, int>[] firstReward;
    // 通关掉落
    public Dictionary<string, int>[] drop;
    // 怪物信息
    public int[][] monster;

}

public class CFG_MainBattle : CFGTableBase<CFGItem_MainBattle>
{

    private static CFG_MainBattle _AllCFG;
    private static CFG_MainBattle AllCFG
    {
        get
        {
            if (_AllCFG == null)
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_MainBattle, CFGItem_MainBattle>();

            return _AllCFG;
        }

    }

    public static CFG_MainBattle Get()
    {
        return AllCFG;
    }

    public static bool IsPassAll(int id)
    {
        return AllCFG.DicDatas.Count <= id;
    }

    public static CFGItem_MainBattle GetById(int id)
    {
        if (AllCFG.DicDatas.TryGetValue(id, out var item))
            return item;

        return null;
    }

    public static CFGItem_MainBattle[] GetCfgsById(int id)
    {
        if (id == 0) id = 1;
        var originId = id;
        var curId = id;
        var list = new List<CFGItem_MainBattle>();
        var totalCount = AllCFG.DicDatas.Count;
        for (int i = 0; i < 5; i++)
        {
            if (AllCFG.DicDatas.TryGetValue(curId, out var item))
            {
                list.Add(item);

                if (curId == totalCount && i == 0)
                {
                    curId--;
                    continue;
                }

                if (item.flag == (int)MainBattleFlayType.Boss)
                {
                    if (curId < originId)
                        curId = originId + 1;
                    else if (curId == originId)
                        curId++;
                    else
                        curId = originId - 1;
                }
                else
                {
                    if (curId < originId)
                        curId--;
                    else
                        curId++;
                }
            }
        }
        list.Sort((a, b) => { return a.ID.CompareTo(b.ID); });
        return list.ToArray();
    }
}