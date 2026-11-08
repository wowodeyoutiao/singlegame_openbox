using JEngine.Helper;

using System.Collections.Generic;

/// <summary>
/// 道具表
/// </summary>
public class CFGItem_Store : CFGItemBase
{
    /// <summary>
    /// 商品ID
    /// </summary>
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// 商城类别
    /// </summary>
    public int StoreType;

    /// <summary>
    /// 物品编号
    /// </summary>
    public int GoodsID;

    /// <summary>
    /// 前置条件
    /// </summary>
    public decimal Prelimit;

    /// <summary>
    /// 界面描述
    /// </summary>
    public string Desc;

    /// <summary>
    /// 图集
    /// </summary>
    public string Icon;

    /// <summary>
    /// 图片名称
    /// </summary>
    public string Iconname;

    //库存检测
    public int Serverdetect;

    /// <summary>
    /// 商品内容
    /// </summary>
    public CFGItem_Solditem[] Solditems;
    /// <summary>
    /// 商品内容
    /// </summary>
    public class CFGItem_Solditem
    {
        public int itemID;
        public int Num;
    }

    /// <summary>
    /// 限制条件
    /// </summary>
    public CFGItem_StoreLimite Limite;

    /// <summary>
    /// 消耗递增
    /// </summary>
    public CFGItem_Consumechange[] Consumechange;
    public class CFGItem_Consumechange
    {
        public int key;
        public int value;
    }

    /// <summary>
    /// 限制条件
    /// </summary>
    /// 人民币的needmoney定义为0，limittype填999代表没有购买限制
    /// 非人民币的其中limittype定义为：{0：代表终生限制；1代表每日限制；2代表每周限制；3代表每月限制
    public class CFGItem_StoreLimite
    {
        /// <summary>
        /// 对应Item ID 人民币为0
        /// </summary>
        public int needmoney;
        public int Num;
        public int limittype;
        public int limitnum;
        public int times;
    }

    /// <summary>
    /// 物品时效 
    /// </summary>
    public CFGItem_StoreTimeliness Timeliness;

    /// <summary>
    /// 物品时效 1：代表从道具获得按照天数进行回收的道具。 2：代表按照自然天数进行回收的道具。
    /// </summary>
    public class CFGItem_StoreTimeliness
    {
        public int type;
        public long starttime;
        public long endtime;
        public long time;
    }

}

/// <summary>
/// Store 配置表
/// </summary>
public class CFG_Store : CFGTableBase<CFGItem_Store>
{

    private static CFG_Store _AllCFG;
    private static CFG_Store AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Store, CFGItem_Store>();
            }
            return _AllCFG;
        }
    }
    public static CFG_Store Get()
    {
        return AllCFG;
    }
    public static Dictionary<long, CFGItem_Store> GetDic()
    {
        return AllCFG.DicDatas;
    }

    public static CFGItem_Store GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static Dictionary<StoreType, List<CFGItem_Store>> StoreTypeListsDic = new Dictionary<StoreType, List<CFGItem_Store>>();
    public static List<CFGItem_Store> GetListByStoreType(StoreType _stype)
    {
        if (!StoreTypeListsDic.ContainsKey(_stype))
        {
            List<CFGItem_Store> _list = new List<CFGItem_Store>();

            for (int i = 0; i < Get().DicList.Count; i++)
            {
                if (Get().DicList[i].StoreType == (int)_stype)
                {
                    _list.Add(Get().DicList[i]);
                }
            }
            _list.Sort((a, b) => { return a.GoodsID.CompareTo(b.GoodsID); });
            StoreTypeListsDic.Add(_stype, _list);
        }

        return StoreTypeListsDic[_stype];
    }

    public static Dictionary<StoreType, string> StoreNameDic = new Dictionary<StoreType, string>()
    {
        [StoreType.RMB] = "充值",
        [StoreType.Coupon] = "元宝",
        [StoreType.Coin] = "金币",
        [StoreType.Diamanta] = "战魂",
        [StoreType.SiLian] = "试炼",
        [StoreType.Ticket] = "门票",
    };
    public enum StoreType : int
    {
        No = 0,
        /// <summary>
        /// 代表充值商城
        /// </summary>
        RMB = 1,

        /// <summary>
        /// 点券商城
        /// </summary>
        Coupon = 2,

        /// <summary>
        /// 金币商城
        /// </summary>
        Coin = 3,

        /// <summary>
        /// 战魂商城
        /// </summary>
        Diamanta = 4,

        /// <summary>
        /// 试炼商城
        /// </summary>
        SiLian = 5,

        /// <summary>
        /// VIP礼包商城
        /// </summary>
        VIPGift = 11,

        /// <summary>
        /// 门票商城
        /// </summary>
        Ticket = 30,

        /// <summary>
        /// 星愿召唤
        /// </summary>
        StarWish = 36,

        /// <summary>
        /// 疯狂洗练
        /// </summary>
        CrazyRefine = 37,
    }

    public enum StroeID
    {
        /// <summary>
        /// 首领挑战
        /// </summary>
        BossChange = 14001,

        /// <summary>
        /// 遗迹挑战次数
        /// </summary>
        RelicChange = 14002,

        /// <summary>
        /// 守卫海岛月卡
        /// </summary>
        HaiDaoCard = 19001,

        /// <summary>
        /// 五行挑战次数
        /// </summary>
        WuXingChange = 28001,

        /// <summary>
        /// 属性槽位3
        /// </summary>
        RefineRandom = 34003,

        /// <summary>
        /// 属性槽位4
        /// </summary>
        RefineSelect = 34004,
    }
}
