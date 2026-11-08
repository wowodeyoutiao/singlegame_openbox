using JEngine.Helper;
using System.Collections.Generic;

/// <summary>
/// 特权卡
/// </summary>
public class CFGItem_Privilege : CFGItemBase
{
    public int PID;
    public override long GetIdx()
    {
        return PID;
    }
    /// <summary>
    /// 商品ID
    /// </summary>
    public int SID; 

    public string Pname;

    public Award[] Dailyaward;

    public Pcontent_Award[] Pcontent;

    public int Energystate; // 活力状态
    public string Pdisplay; // 特权展示
    public int Duration; // 持续时间
    public int SkipAd; // 跳过广告
    public int StaminaUp; // 体力上限
    public int QuicklyIsiand; // 扫荡海岛
    public int QuicklyDaily; // 扫荡日常副本
    public int AutoBattle; // 自动战斗

    public struct Award 
    {
        public int id;
        public int num;
    }

    public struct Pcontent_Award
    {
        public int key;
        public int value;
    }

}

/// <summary>
/// CFG_Privilege 配置表
/// </summary>
public class CFG_Privilege : CFGTableBase<CFGItem_Privilege>
{
    private static CFG_Privilege _AllCFG;
    private static CFG_Privilege AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Privilege, CFGItem_Privilege>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Privilege Get()
    {
        return AllCFG;
    }

    public static Dictionary<long, CFGItem_Privilege> GetDic()
    {
        return AllCFG.DicDatas;
    }


    public static CFGItem_Privilege GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }
}
