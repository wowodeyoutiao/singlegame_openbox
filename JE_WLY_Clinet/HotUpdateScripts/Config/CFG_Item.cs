using JEngine.Helper;
using System.Collections.Generic;

/// <summary>
/// 道具表
/// </summary>
[System.Serializable]
public class CFGItem_Item : CFGItemBase
{
    /// <summary>
    /// 道具ID
    /// </summary>
    public int ID;
    /// <summary>
    /// 名字
    /// </summary>
    public string Name;

    /// <summary>
    /// 道具类型
    /// </summary>
    public decimal? Type;

    /// <summary>
    /// 最大堆叠数量
    /// </summary>
    public decimal? Maxnum;

    /// <summary>
    /// 描述
    /// </summary>
    public string Desc;

    /// <summary>
    /// 图集
    /// </summary>
    public string Icon;

    /// <summary>
    /// 图片名称
    /// </summary>
    public string Iconname { get;set; }

    /// <summary>
    /// 品质
    /// </summary>
    public int Quality;

    /// <summary>
    /// 排序优先级 用来表示当前道具在包裹中的排列优先级，优先级数值越高，代表着该道具在包裹中的排序越靠前
    /// </summary>
    public int Priority;

    /// <summary>
    /// 物品时效 
    /// </summary>
    public CFGItem_ItemTimeliness Timeliness;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// 物品时效 1：代表从道具获得按照天数进行回收的道具。 2：代表按照自然天数进行回收的道具。
    /// </summary>
    public class CFGItem_ItemTimeliness
    {
        public int type;
        public long starttime;
        public long endtime;
        public long time;
    }
    /// <summary>
    /// 装备部位
    /// </summary>
    public decimal? Equiptype = 0; 
    /// <summary>
    /// 等级限制
    /// </summary>
    public decimal? Levellimit = 0; 
    /// <summary>
    /// 职业限制
    /// </summary>
    public decimal? Occupylimit = 0;  
    /// <summary>
    /// 性别限制
    /// </summary>
    public decimal? Genderlimit = 0; 
    /// <summary>
    /// 物理攻击
    /// </summary>
    public decimal? Phyatk = 0;  
    /// <summary>
    /// 魔法攻击
    /// </summary>
    public decimal? Magatk = 0; 
    /// <summary>
    /// 物理防御
    /// </summary>
    public decimal? Phydef = 0;  
    /// <summary>
    /// 魔法防御
    /// </summary>
    public decimal? Magdef = 0; 
    /// <summary>
    /// 生命值
    /// </summary>
    public decimal? HP = 0;
    /// <summary>
    /// 命中值
    /// </summary>
    public decimal? Hit = 0; 
    /// <summary>
    /// 闪避值
    /// </summary>
    public decimal? Miss = 0;
    /// <summary>
    /// 暴击值
    /// </summary>
    public decimal? Crit = 0;
    /// <summary>
    /// 韧性值
    /// </summary>
    public decimal? Decrit = 0;
    /// <summary>
    /// 对人增伤
    /// </summary>  
    public decimal? PVPdamage = 0;
    /// <summary>
    /// 对怪增伤
    /// </summary>
    public decimal? PVEdamage = 0;

    /// <summary>
    /// 
    /// </summary>
    public decimal? HPPercent = 0;
    
    /// <summary>
    /// 立即使用
    /// </summary>
    public decimal? Usenow = 0;

    /// <summary>
    /// 道具子类
    /// </summary>
    public int SonType;

    /// <summary>
    /// 合成标记
    /// </summary>
    public int ComposeMark;
    /// <summary>
    /// 自然属性标记
    /// </summary>
    public int NatureMark;

    /// <summary>
    /// 物品等级
    /// </summary>
    public decimal itemlevel;

    public int[] DiamondPosition;

    public string CardIcon;
}

/// <summary>
/// Item 配置表
/// </summary>
public class CFG_Item : CFGTableBase<CFGItem_Item>
{

    private static CFG_Item _AllCFG;
    private static CFG_Item AllCFG
    { 
        get 
        { 
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Item, CFGItem_Item>();
            }
            return _AllCFG;
        } 
    }
    public static CFG_Item Get()
    {
        return AllCFG;
    }

    public static CFGItem_Item GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static int GetQulite(int id)
    {
        return (int)GetById(id).Quality;
    }

    private static Dictionary<int, List<int>> EquipSoltToJewelDinc;

    public static List<int> GetCfgIdsByEquipSolt(int solt)
    {
        if (EquipSoltToJewelDinc == null)
        {
            EquipSoltToJewelDinc = new Dictionary<int, List<int>>();
            var list = Get().DicList;
            for (int i = 0; i < list.Count; i++)
            {
                var DiamondPosition = list[i].DiamondPosition;
                if (DiamondPosition != null)
                {
                    for (int j = 0; j < DiamondPosition.Length; j++)
                    {
                        if (!EquipSoltToJewelDinc.ContainsKey(DiamondPosition[j]))
                        {
                            EquipSoltToJewelDinc.Add(DiamondPosition[j], new List<int>());
                        }
                        EquipSoltToJewelDinc[DiamondPosition[j]].Add(list[i].ID);
                    }
                }
            }
        }
        return EquipSoltToJewelDinc[solt];
    }
    /// <summary>
    /// 道具是否快捷使用
    /// </summary>
    /// <param name="id"></param>
    /// <returns></returns>
    public static bool CheckIsItemTipUsed(int id)
    {
        return GetById(id).Usenow == 1;
    }

    /// <summary>
    /// 检测是否可被销毁
    /// </summary>
    /// <param name="_type"> 道具类型</param>
    /// <returns></returns>
    public static bool GeCanDel(decimal _type)
    {
        return _type >= 1 && _type <= 5;
    }

    /// <summary>
    /// 检测是否可被使用
    /// </summary>
    /// <param name="_type"> 道具类型</param>
    /// <returns></returns>
    public static bool GeCanUsed(decimal _type)
    {
        return _type == (int)ItemType.NorDelAndUse || _type == (int)ItemType.Equip || _type == (int)ItemType.TimeDelAndUse;
    }


    private static Dictionary<decimal, string> JobDic = new Dictionary<decimal, string>
    {
        [0] = "无",
        [1] = "战士",
        [2]= "法师",
        [3] = "祭祀",
    };
    public static string GetJobStrById(decimal id)
    {
        return JobDic[id];
    }

    private static Dictionary<decimal, string> SexDic = new Dictionary<decimal, string>
    {
        [0] = "无",
        [1] = "男",
        [2] = "女",
    };
    public static string GetSexStrById(decimal id)
    {
        return SexDic[id];
    }

    /// <summary>
    /// 道具类型
    /// </summary>
    public enum ItemType :int
    {
        /// <summary>
        /// 不进入包裹
        /// </summary>
        NoBag = 0,

        /// <summary>
        /// 可被销毁的道具
        /// </summary>
        NorDel = 1,

        /// <summary>
        /// 可被销毁和使用
        /// </summary>
        NorDelAndUse = 2,

        /// <summary>
        /// 可被销毁、带有时效的道具
        /// </summary>
        TimeDel = 3,

        /// <summary>
        /// 可被销毁和使用、带有时效的道具
        /// </summary>
        TimeDelAndUse = 4,

        /// <summary>
        /// 装备类
        /// </summary>
        Equip = 5,

        /// <summary>
        /// 宝石类，该雷道具可镶嵌、可合成
        /// </summary>
        Jewel = 6,

        /// <summary>
        /// 伙伴碎片
        /// </summary>
        PartnerPice = 7,

        /// <summary>
        /// 伙伴整卡
        /// </summary>
        PartnerCard = 8,

        /// <summary>
        /// 宠物
        /// </summary>
        Pet = 9,

        /// <summary>
        /// 称号碎片
        /// </summary>
        TitlePice = 10,
    }

    /// <summary>
    /// 道具时限
    /// </summary>
    public enum ItemTimeType : int
    {
        /// <summary>
        /// 无时间限制
        /// </summary>
        NoTime = 0,

        /// <summary>
        /// 从道具获得按照天数进行回收的道具
        /// </summary>
        EndTimeDay = 1,

        /// <summary>
        /// 按照自然天数进行回收的道具
        /// </summary>
        EndTime = 2,
    }


    public enum ItemID
    {
        No = 0,
        /// <summary>
        /// 金币
        /// </summary>
        Gold = 1,

        /// <summary>
        /// 经验
        /// </summary>
        EXP = 2,

        /// <summary>
        /// 行会贡献
        /// </summary>
        GuildContribution = 3,

        /// <summary>
        /// 战魂
        /// </summary>
        ZhanHun = 4,

        /// <summary>
        /// 点券
        /// </summary>
        DianQuan = 5,

        /// <summary>
        /// 试炼凭证
        /// </summary>
        SLPZ = 8,

        /// <summary>
        /// 强化石
        /// </summary>
        Strengthen = 21,

        /// <summary>
        /// 升星石
        /// </summary>
        UpStar = 22,

        /// <summary>
        /// 书页
        /// </summary>
        ShuYe = 23,

        /// <summary>
        /// 技能书卷轴
        /// </summary>
        SkillJuanZhou = 24,

        /// <summary>
        /// 15分钟挂机时间
        /// </summary>
        hangUpTime = 56,

        /// <summary>
        /// 星愿币
        /// </summary>
        StarWishCoin = 58,

        /// <summary>
        /// 元素赋能
        /// </summary>
        ElementFN = 75,

        /// <summary>
        /// 元素突破
        /// </summary>
        ElementTP = 76,

        /// <summary>
        /// 宠物经验值
        /// </summary>
        PetEXP = 83,

        /// <summary>
        /// 祈愿点
        /// </summary>
        GuildWish = 93,

        /// <summary>
        /// 无双宠物兑换卡
        /// </summary>
        PetChangeCard = 98,

        /// <summary>
        /// 加速卷轴
        /// </summary>
        JSJZ = 99,

        /// <summary>
        /// 首领挑战卷轴
        /// </summary>
        SLTZJZ = 25073,

        /// <summary>
        /// 遗迹挑战卷轴
        /// </summary>
        YJTZJZ = 25074,

        /// <summary>
        /// 强化石副本挑战卷轴
        /// </summary>
        QHSTZJZ = 25075,

        /// <summary>
        /// 升星石副本挑战卷轴
        /// </summary>
        SXSTZJZ = 25076,

        /// <summary>
        /// 魂石副本挑战卷轴
        /// </summary>
        HSTZJZ = 25077,

        /// <summary>
        /// 书页副本挑战卷轴
        /// </summary>
        SYTZJZ = 25078,

        /// <summary>
        /// 金币副本挑战卷轴
        /// </summary>
        CoinTZJZ = 25079,

        /// <summary>
        /// 经验副本挑战卷轴
        /// </summary>
        ExpTZJZ = 25080,

        /// <summary>
        /// 五行试炼挑战卷轴
        /// </summary>
        WXTZJZ = 25081,

        /// <summary>
        /// 即时赛挑战卷轴
        /// </summary>
        PVETZJZ = 25082,

        /// <summary>
        /// 回合赛挑战卷轴
        /// </summary>
        PVPTZJZ = 25083,
    }
}
