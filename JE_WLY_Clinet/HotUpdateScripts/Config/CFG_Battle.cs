using JEngine.Helper;
using System.Collections.Generic;

public class CFGItem_Battle : CFGItemBase
{
    public int UID;

    public override long GetIdx()
    {
        return UID;
    }

    /// <summary>
    /// 副本类型
    /// </summary>
    public int Copytype;

    /// <summary>
    /// 战力限制
    /// </summary>
    public int Powerlimit;

    /// <summary>
    /// 等级限制
    /// </summary>
    public int Levellimit;

    /// <summary>
    /// 章节
    /// </summary>
    public int Chapter;

    /// <summary>
    /// 闯关层数
    /// </summary>
    public int Layers;

    /// <summary>
    /// 刷怪配置
    /// </summary>
    public AwardItem[] Monster { get; set; }

    /// <summary>
    /// 挑战时间
    /// </summary>
    //public int Timelimit;

    public struct AwardItem
    {
        public int id;
        public int num;
        public int rate;
    }

    /// <summary>
    /// 结算奖励
    /// </summary>
    public AwardItem[] Settlement;

    /// <summary>
    /// 产出预览
    /// </summary>
    public int[] Preview;

    /// <summary>
    /// 章节奖励
    /// </summary>
    public AwardItem[] ChapterAward;

    /// <summary>
    /// 星级奖励
    /// </summary>
    public StarAwardItem[] StarAward;

    public struct StarAwardItem
    {
        public int key;
        public AwardItem[] values;
    }

    public AwardItem[] GetStarAwardByNum(int num)
    {
        for (int i = 0; i < StarAward.Length; i++)
        {
            if (StarAward[i].key == num)
            {
                return StarAward[i].values;
            }
        }
        return null;
    }

    public AwardItem[] Randomaward;

    public decimal MapID;

    public decimal MonsterShow;

}

/// <summary>
/// Battle 配置表
/// </summary>
public class CFG_Battle : CFGTableBase<CFGItem_Battle>
{
    private static CFG_Battle _AllCFG;
    private static CFG_Battle AllCFG
    {
        get
        {
            if (_AllCFG == null)
            {
                _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Battle, CFGItem_Battle>();
            }
            return _AllCFG;
        }
    }

    public static CFG_Battle Get()
    {
        return AllCFG;
    }

    public static CFGItem_Battle GetById(int id)
    {
        return AllCFG.DicDatas[id];
    }

    private static List<CFGItem_Battle> _ALlFightList;

    //战役列表
    public static List<CFGItem_Battle> AllFightList
    {
        get
        {
            if (_ALlFightList == null)
            {
                _ALlFightList = new List<CFGItem_Battle>();
                for (int i = 0; i < AllCFG.DicList.Count; i++)
                {
                    if (AllCFG.DicList[i].Copytype == (int)CFG_Room.RoomType.Fight)
                    {
                        _ALlFightList.Add(AllCFG.DicList[i]);
                    }
                }
                _ALlFightList.Sort((a, b) => { return a.UID.CompareTo(b.UID); });
            }
            return _ALlFightList;
        }
    }

    //获取战役下一个
    public static CFGItem_Battle GetNextFightById(int id)
    {
        for (int i = 0; i < AllFightList.Count; i++)
        {
            if (id < AllFightList[i].UID)
            {
                return AllFightList[i];
            }
        }
        return AllFightList[(AllFightList.Count - 1)];
    }

    //战役章节奖励
    private static List<CFGItem_Battle> _AllFightChapterAwardList;

    /// <summary>
    /// 战役章节奖励
    /// </summary>
    public static List<CFGItem_Battle> AllFightChapterAwardList
    {
        get
        {
            if (_AllFightChapterAwardList == null)
            {
                _AllFightChapterAwardList = new List<CFGItem_Battle>();
                for (int i = 0; i < AllFightList.Count; i++)
                {
                    if (AllFightList[i].ChapterAward != null)
                    {
                        _AllFightChapterAwardList.Add(AllCFG.DicList[i]);
                    }
                }
            }
            return _AllFightChapterAwardList;
        }
    }

    private static List<CFGItem_Battle> _SLTowerList;
    //获取试炼塔
    public static List<CFGItem_Battle> GetSLTowerList()
    {
        if (_SLTowerList == null)
        {
            _SLTowerList = new List<CFGItem_Battle>();
            for (int i = 0; i < AllCFG.DicList.Count; i++)
            {
                if (AllCFG.DicList[i].Copytype == (int)CFG_Room.RoomType.SLTower)
                {
                    _SLTowerList.Add(AllCFG.DicList[i]);
                }
            }
            _SLTowerList.Sort((a, b) => { return b.UID.CompareTo(a.UID); });
        }
        return _SLTowerList;
    }

    private static List<CFGItem_Battle> _BossList;
    //获取首領
    public static List<CFGItem_Battle> BossList
    {
        get
        {
            if (_BossList == null)
            {
                _BossList = new List<CFGItem_Battle>();
                for (int i = 0; i < AllCFG.DicList.Count; i++)
                {
                    if (AllCFG.DicList[i].Copytype == (int)CFG_Room.RoomType.Boss)
                    {
                        _BossList.Add(AllCFG.DicList[i]);
                    }
                }
                _BossList.Sort((a, b) => { return b.UID.CompareTo(a.UID); });
            }
            return _BossList;
        }
    }

    private static List<CFGItem_Battle> _RelicList;
    //遗迹信息
    public static List<CFGItem_Battle> RelicList
    {
        get
        {
            if (_RelicList == null)
            {
                _RelicList = new List<CFGItem_Battle>();
                for (int i = 0; i < AllCFG.DicList.Count; i++)
                {
                    if (AllCFG.DicList[i].Copytype == (int)CFG_Room.RoomType.Relic)
                    {
                        _RelicList.Add(AllCFG.DicList[i]);
                    }
                }
                _RelicList.Sort((a, b) => { return a.UID.CompareTo(b.UID); });
            }
            return _RelicList;
        }
    }

    private static Dictionary<int, List<CFGItem_Battle>> _RelicDic;
    //获取遗迹信息
    public static Dictionary<int, List<CFGItem_Battle>> RelicDic
    {
        get
        {
            if (_RelicDic == null)
            {
                _RelicDic = new Dictionary<int, List<CFGItem_Battle>>();
                for (int i = 0; i < AllCFG.DicList.Count; i++)
                {
                    if (AllCFG.DicList[i].Copytype == (int)CFG_Room.RoomType.Relic)
                    {
                        if (!_RelicDic.ContainsKey(AllCFG.DicList[i].Chapter))
                            _RelicDic.Add(AllCFG.DicList[i].Chapter, new List<CFGItem_Battle>());

                        _RelicDic[AllCFG.DicList[i].Chapter].Add(AllCFG.DicList[i]);
                    }
                }
            }
            return _RelicDic;
        }
    }

    private static List<CFGItem_Battle> _YZDDList;
    //一战到底信息
    public static List<CFGItem_Battle> YZDDList
    {
        get
        {
            if (_YZDDList == null)
            {
                _YZDDList = new List<CFGItem_Battle>();
                for (int i = 0; i < AllCFG.DicList.Count; i++)
                {
                    if (AllCFG.DicList[i].Copytype == (int)CFG_Room.RoomType.YZDD)
                    {
                        _YZDDList.Add(AllCFG.DicList[i]);
                    }
                }
                _YZDDList.Sort((a, b) => { return a.UID.CompareTo(b.UID); });
            }
            return _YZDDList;
        }
    }



    public static int GetYZDDLevelCount()
    {
        int count = 0;
        foreach (var item in AllCFG.DicList)
        {
            if (item.Copytype == (int)CFG_Room.RoomType.YZDD)
                count++;
        }
        return count;
    }

}
