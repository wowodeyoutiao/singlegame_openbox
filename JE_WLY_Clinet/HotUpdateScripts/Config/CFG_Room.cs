using JEngine.Helper;
using System.Linq;

public class CFGItem_Room : CFGItemBase
{
    public int ID;

    public override long GetIdx()
    {
        return ID;
    }

    /// <summary>
    /// Room名称
    /// </summary>
    public string Name;

    /// <summary>
    /// 游戏类型
    /// </summary>
    public int GameType;

    /// <summary>
    /// 最低进入限制
    /// </summary>
    public int MinEntryLimit;

    /// <summary>
    /// 最高进入限制
    /// </summary>
    public int MaxEntryLimit;

    /// <summary>
    /// 人数
    /// </summary>
    public int PlayerNum;

    /// <summary>
    /// 地图ID
    /// </summary>
    public int MapId;

    /// <summary>
    /// 房间类型
    /// </summary>
    public CFG_Room.RoomType Roomtype;

    /// <summary>
    /// 需要体力
    /// </summary>
    public decimal Needenerge;

    /// <summary>
    /// 增加时间
    /// </summary>
    public decimal Addtime;

    public int[] Droppreview;

    public string CopyDes;
    public string Suggest;
    public decimal Show;

    public string Music;
}

/// <summary>
/// Room 配置表
/// </summary>
public class CFG_Room : CFGTableBase<CFGItem_Room>
{
    private static CFG_Room _AllCFG;
    private static CFG_Room AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Room, CFGItem_Room>();
            }
            return _AllCFG;
        }
    }

    public static CFG_Room Get()
    {
        return AllCFG;
    }

    public static CFGItem_Room GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    public static CFGItem_Room GetByRoomType(RoomType roomType, int lv)
    {
        CFGItem_Room select = null;
        foreach (var item in Get().DicList)
        {
            if (item.Roomtype == roomType && (lv < 0 || item.MinEntryLimit <= lv && lv <= item.MaxEntryLimit))
            {

                select = item;
                break;
            }
        }
        return select;
    }

    private static CFGItem_Room[] _WarFights;
    public static CFGItem_Room[] WarFights
    {
        get
        {
            if (_WarFights == null)
            {
                _WarFights = AllCFG.DicList.Where(x => x.Roomtype == RoomType.PVE_NoPeace).ToArray();
            }
            return _WarFights;
        }
    }

    public enum RoomType
    {
        PVE_Peace = 1,
        PVE_NoPeace = 2,
        PVP_One = 3,
        PVP_More = 4,
        Tower = 5,
        WorldBoss = 6,
        /// <summary>
        /// 经验副本
        /// </summary>
        EXP = 7,
        /// <summary>
        /// 金币副本
        /// </summary>
        Gold = 8,
        /// <summary>
        /// 强化石副本
        /// </summary>
        QHS = 9,
        /// <summary>
        /// 书页副本
        /// </summary>
        SY = 10,
        /// <summary>
        /// 升星石副本
        /// </summary>
        SXS = 11,
        /// <summary>
        /// 魂石副本
        /// </summary>
        HS = 12,

        /// <summary>
        /// 新手引导
        /// </summary>
        Guide = 13,

        /// <summary>
        /// 战役
        /// </summary>
        Fight = 14,

        /// <summary>
        /// 首领副本
        /// </summary>
        Boss = 15,

        /// <summary>
        /// 试炼塔副本
        /// </summary>
        SLTower = 16,

        /// <summary>
        /// 遗迹副本
        /// </summary>
        Relic = 17,

        /// <summary>
        /// 工会战
        /// </summary>
        Guild = 18,

        /// <summary>
        /// 一战到底
        /// </summary>
        YZDD = 19,

        /// <summary>
        /// 五行试炼
        /// </summary>
        WuXing = 20,
    }

    public static int GetUIDByRoomType(RoomType roomType)
    {
        switch (roomType)
        {
            case RoomType.Fight:
                return 600;
            case RoomType.Boss:
                return 700;
            case RoomType.SLTower:
                return 800;
            case RoomType.Relic:
                return 900;
            case RoomType.Guild:
                return 1000;
            case RoomType.YZDD:
                return 1100;
            default:
                return 0;
        }
    }
}
