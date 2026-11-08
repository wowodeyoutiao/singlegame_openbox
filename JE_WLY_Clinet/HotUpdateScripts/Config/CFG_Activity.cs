using JEngine.Helper;

using System.Collections.Generic;

public class CFGItem_Activity : CFGItemBase
{
    public int UUID;

    public override long GetIdx()
    {
        return UUID;
    }

    /// <summary>
    /// 活动类型 1：代表7日连续登录 2：代表30日累计登录 3：代表等级奖励 4：代表推荐礼包
    /// </summary>
    public int Type;

    public int Acttag;

    /// <summary>
    /// 内容编号
    /// </summary>
    public int Content;

    /// <summary>
    /// 普通道具
    /// </summary>
    public Dictionary<string, int>[] CommonItem;

    /// <summary>
    /// 付费道具
    /// </summary>
    public Dictionary<string, int>[] PayItem;

    /// <summary>
    /// 前置条件
    /// </summary>
    public int PreLimit;

    /// <summary>
    /// 补领消耗
    /// </summary>
    public Dictionary<string, int>[] CompensationConsume;

    /// <summary>
    /// 活动时间
    /// </summary>
    public string ActivityTime;

    /// <summary>
    /// 活动规则
    /// </summary>
    public string RuleDes;

    /// <summary>
    /// 活动名称
    /// </summary>
    public string Name;

    public string Panelname;

    public decimal ResetTime;
}

/// <summary>
/// Activity 配置表
/// </summary>
public class CFG_Activity : CFGTableBase<CFGItem_Activity>
{
    private static CFG_Activity _AllCFG;
    private static CFG_Activity AllCFG
    {
        get
        {
            if (_AllCFG == null)
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Activity, CFGItem_Activity>();

            return _AllCFG;
        }

    }

    public static CFG_Activity Get()
    {
        return AllCFG;
    }

    public static CFGItem_Activity GetById(int id)
    {
        return ExcelDataManager.Instance.GetExcelItem<CFG_Activity, CFGItem_Activity>(id);
    }

    private static Dictionary<Type, List<CFGItem_Activity>> ActiveList = new Dictionary<Type, List<CFGItem_Activity>>();

    public static List<CFGItem_Activity> GetActiveListByType(Type _type)
    {
        if (!ActiveList.ContainsKey(_type))
        {
            List<CFGItem_Activity> _list = new List<CFGItem_Activity>();
            var list = Get().DicList;
            for (int i = 0; i < list.Count; i++)
            {
                if (list[i].Type == (int)_type)
                {
                    _list.Add(list[i]);
                }
            }
            ActiveList.Add(_type, _list);
        }
        return ActiveList[_type];
    }

    private static Dictionary<int, CFGItem_Activity> _ActiveDic = null;

    public static Dictionary<int, CFGItem_Activity> GetActiveTypeDic()
    {
        if (_ActiveDic == null)
        {
            _ActiveDic = new Dictionary<int, CFGItem_Activity>();

            var list = Get().DicList;
            for (int i = 0; i < list.Count; i++)
            {
                if (!_ActiveDic.ContainsKey(list[i].Type))
                {
                    _ActiveDic.Add(list[i].Type, list[i]);
                }
            }
        }
        return _ActiveDic;
    }

    /// <summary>
    /// 活动类型 
    /// </summary>
    public enum Type
    {
        /// <summary>
        /// 7日连续登录
        /// </summary>
        Sign7 = 1,

        /// <summary>
        /// 30日累计登录
        /// </summary>
        Sign30 = 2,

        /// <summary>
        /// 等级奖励
        /// </summary>
        LevelReward = 3,

        /// <summary>
        /// 推荐礼包
        /// </summary>
        GiftPack = 4,

        /// <summary>
        /// 首充礼包
        /// </summary>
        FirstPay = 5,

        /// <summary>
        /// 点券大转盘
        /// </summary>
        TurnWheel = 6,

        /// <summary>
        /// 限时悬赏
        /// </summary>
        LimitBounty = 7,

        /// <summary>
        /// 单笔充值返利
        /// </summary>
        SingleRechargeRebate = 8,

        /// <summary>
        /// 累计充值返利
        /// </summary>
        CumulativeRechargeRebate = 9,

        /// <summary>
        /// 试炼基金
        /// </summary>
        TrialFund = 10,

        /// <summary>
        /// 伙伴礼包
        /// </summary>
        PartnerGift = 11,

        /// <summary>
        /// 宠物礼包
        /// </summary>
        PetGift = 12,

        /// <summary>
        /// 成长基金
        /// </summary>
        GrowthFund = 15,

        /// <summary>
        /// 助战礼包
        /// </summary>
        BattleSupportGift = 24,

        /// <summary>
        /// 洗炼礼包
        /// </summary>
        RefineGift = 25,

        /// <summary>
        /// 疯狂洗炼
        /// </summary>
        StarWish = 26,

        /// <summary>
        /// 疯狂洗炼
        /// </summary>
        CrazyRefine = 27,

        /// <summary>
        /// 洗练基金
        /// </summary>
        RefineFund = 28,

        /// <summary>
        /// 召唤基金
        /// </summary>
        CallFund = 29,

        /// <summary>
        /// 副本基金
        /// </summary>
        CopyFund = 30,
    }
}
