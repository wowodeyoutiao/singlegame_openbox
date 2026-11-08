using JEngine.Helper;
using System.Collections.Generic;

/// <summary>
/// 跑马灯
/// </summary>
public class CFGItem_Lamp : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return (int)ID;
    }

    /// <summary>
    /// 走马灯公告
    /// </summary>
    public string Broad;

    /// <summary>
    /// 类型
    /// </summary>
    public int[] Type;

    /// <summary>
    /// 跳转面板
    /// </summary>
    public string Panel;

}

/// <summary>
/// Lamp 配置表
/// </summary>
public class CFG_Lamp : CFGTableBase<CFGItem_Lamp>
{
    private static CFG_Lamp _AllCFG;
    private static CFG_Lamp AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Lamp, CFGItem_Lamp>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Lamp Get()
    {
        return AllCFG;
    }
    public static Dictionary<long, CFGItem_Lamp> GetDic()
    {
        return AllCFG.DicDatas;
    }

    public static CFGItem_Lamp GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static string GetStr(int id, List<string> strs)
    {
        List<object> objList = new List<object>();

        var types = GetById(id).Type;
        for (int i = 0; i < types.Length; i++)
        {
            if (types[i] == 6)
                break;

            string str = i < strs.Count ? strs[i] : "";

            var strValue = GetStr(types[i], str);
            if (strValue != null)
            {
                objList.Add(strValue);
            }
        }

        return string.Format(GetById(id).Broad, objList.ToArray());
    }

    private static object GetStr(int _type, string str)
    {
        LampType type = (LampType)_type;
        return type switch
        {
            LampType.No => null,
            LampType.RoleName => "<color=#95F204>" + str + "</color>",
            LampType.MonsterID => "<color=#D9001B>" + CFG_Monster.GetById(int.Parse(str)).Name + "</color>",
            LampType.ItemID => GameUtil.Instance.GetItemName(int.Parse(str)),
            LampType.SeverMsg => "<color=#D900B6>" + str + "</color>",
            LampType.Str => str,
            LampType.Jump => null,
            _ => null,
        };
    }

    public static bool IsJump(int id)
    {
        var types = GetById(id).Type;
        for (int i = 0; i < types.Length; i++)
        {
            if ((int)LampType.Jump == types[i])
                return true;
        }
        return false;
    }

    public enum LampType : int
    {
        No = 0,

        RoleName = 1,

        MonsterID = 2,

        ItemID = 3,

        SeverMsg = 4,

        Str = 5,

        Jump = 6,
    }
}
