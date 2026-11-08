using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_VIP : CFGItemBase
{
    public int VIPLevel;

    public override long GetIdx()
    {
        return VIPLevel;
    }

    /// <summary>
    /// VIP名称
    /// </summary>
    public string VIPName;

    /// <summary>
    /// 经验值
    /// </summary>
    public int NeedEXP;

    /// <summary>
    /// VIP特权内容
    /// </summary>
    public DailyTimeItem[] Pcontent;

    /// <summary>
    /// 解锁伙伴槽位
    /// </summary>
    public int Partnerslot;

    /// <summary>
    /// 活力状态增加/分
    /// </summary>
    public int energyadd;

    /// <summary>
    /// 特权展示
    /// </summary>
    public string Pdisplay;

    /// <summary>
    /// 特权展示
    /// </summary>
    public DailyTimeItem[] dailytime;

    /// <summary>
    /// 折扣率
    /// </summary>
    public int Discount;

    public class DailyTimeItem
    {
        public int key;
        public int value;
    }

}

/// <summary>
/// VIP 配置表
/// </summary>
public class CFG_VIP : CFGTableBase<CFGItem_VIP>
{
    private static CFG_VIP _AllCFG;
    private static CFG_VIP AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_VIP, CFGItem_VIP>();
            }
            return _AllCFG;
        }
    }
    public static CFG_VIP Get()
    {
        return AllCFG;
    }

    public static CFGItem_VIP GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static bool IsMaxVipLevel(int lv)
    {
        return Get().DicList.Count <= lv;
    }

    public static int GetMaxVipLevel()
    {
        return Get().DicList.Count;
    }

    public static int GetAddRoomCount(int lv, int room)
    {
        if (lv > 0)
        {
            var dailytime = GetById(lv).dailytime;
            if (dailytime != null)
            {
                for (int i = 0; i < dailytime.Length; i++)
                {
                    if (dailytime[i].key == room)
                        return dailytime[i].value;
                }
            }

        }
        return 0;
    }

    private static Dictionary<int, int> _PartnerSlotVipDic;
    private static Dictionary<int, int> PartnerSlotVipDic
    {
        get
        {
            if (_PartnerSlotVipDic == null)
            {
                _PartnerSlotVipDic = new Dictionary<int, int>();
                for (int i = 0; i < Get().DicList.Count; i++)
                {
                    if (!_PartnerSlotVipDic.ContainsKey(Get().DicList[i].Partnerslot))
                    {
                        _PartnerSlotVipDic.Add(Get().DicList[i].Partnerslot, Get().DicList[i].VIPLevel);
                    }
                }
            }
            return _PartnerSlotVipDic;
        }
    }
    //根据槽位获取VIP解锁等级
    public static int GetLimetVipBySolt(int solt)
    {
        return PartnerSlotVipDic[solt];
    }
    public static int GetPartnerSlot(int vip)
    {
        if (vip == 0)
            return 0;
        return GetById(vip).Partnerslot;
    }
}
