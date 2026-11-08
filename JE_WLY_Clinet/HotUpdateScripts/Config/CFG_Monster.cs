using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_Monster : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// Monster名称
    /// </summary>
    public string Name;

    /// <summary>
    /// 等级
    /// </summary>
    public int Level;

    /// <summary>
    /// boss标识
    /// </summary>
    public int Isboss;

    /// <summary>
    /// 是否显示血条
    /// </summary>
    public int Isdisplaybar;

    /// <summary>
    /// 是否显示等级
    /// </summary>
    public int Isdisplaylevel;

    public int AttackDistance;

    public string Model;

    public string AttackEffect;

    public string Photo;
    public string Photoname { get; set; }

    public string Icon;

    public string Iconname { get; set; }

    public int Delay;

    public Dictionary<string, int>[] BassProperty;

    public int Quality;

    public decimal Barnum;

    public int[] Skills;
    //public class DorpItem
    //{
    //    public int id;
    //}

    //public DorpItem[] Dorp;

}

/// <summary>
/// Monster 配置表
/// </summary>
public class CFG_Monster : CFGTableBase<CFGItem_Monster>
{
    private static CFG_Monster _AllCFG;
    private static CFG_Monster AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Monster, CFGItem_Monster>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Monster Get()
    {
        return AllCFG;
    }

    public static CFGItem_Monster GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static string GetHUDNameById(int _cfgID)
    {
        CFGItem_Monster info = GetById(_cfgID); 
        //return $"{(info.Isdisplaylevel == 1 ? "Lv" + info.Level : "")}{GameUtil.Instance.GetStrByQuality(info.Quality, info.Name)}";
        return $"{GameUtil.Instance.GetFightStrByQuality(info.Quality, info.Name)}";
    }

    public static string GetMonsterNameById(int _cfgID)
    {
        CFGItem_Monster info = GetById(_cfgID);
        return $"{GameUtil.Instance.GetStrByQuality(info.Quality, info.Name)}";
    }
    public static string GetMonsterNameByLv(int _cfgID)
    {
        CFGItem_Monster info = GetById(_cfgID);
        return $"{(info.Isdisplaylevel == 1 ? "Lv" + info.Level : "")}";
    }

    public static bool IsRomate(int _cfgID)
    {
        CFGItem_Monster info = GetById(_cfgID);
        return info.Model.ToLower().Contains("romate");
    }
}
