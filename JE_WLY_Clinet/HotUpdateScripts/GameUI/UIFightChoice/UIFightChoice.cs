using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Cmp;
using BestHTTP.SecureProtocol.Org.BouncyCastle.Crypto;

using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using SuperScrollView;

using System.Collections.Generic;
using System.Text;

using Unity.Jobs;

using UnityEngine;

using static CFG_XMLConfig.XML_Config;
using static GameEvent;

public class UIFightChoice : UIBase<UIModelFightChoice, UIViewFightChoice>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    //一级选择
    public enum SelectType { Challange, Limit, FuBen, WuXing,WarFight }
    private SelectType selectType = SelectType.Challange;

    //二级选择
    public enum ChanllengeType 
    {
        /// <summary>
        /// 战役
        /// </summary>
        Fight,

        /// <summary>
        /// 首领
        /// </summary>
        Boss, 

        /// <summary>
        /// 试练塔
        /// </summary>
        SLTower, 

        /// <summary>
        /// 遗迹
        /// </summary>
        Relic, 

        /// <summary>
        /// 乱斗之王
        /// </summary>
        PvP,

        /// <summary>
        /// 世界BOSS
        /// </summary>
        WorldBoss,

        /// <summary>
        /// 一战到底
        /// </summary>
        FightEnd,

        /// <summary>
        /// 每日必做
        /// </summary>
        DayDo,

        /// <summary>
        /// 守卫海岛
        /// </summary>
        HaiDao
    }
    private ChanllengeType chanllengeType = ChanllengeType.Fight;

    //页面标题
    private Dictionary<ChanllengeType, string> ChanllengeToStrDic = new Dictionary<ChanllengeType, string>()
    {
        { ChanllengeType.Fight, "战役" },
        { ChanllengeType.Boss, "首领" },
        { ChanllengeType.SLTower, "试炼\n塔" },
        { ChanllengeType.Relic, "遗迹" },
        { ChanllengeType.PvP, "乱斗\n之王" },
        { ChanllengeType.WorldBoss, "世界\nBOSS" },
        { ChanllengeType.FightEnd, "一战\n到底" },
        { ChanllengeType.DayDo, "每日\n必做" },
        { ChanllengeType.HaiDao, "挂机\n升级" },
    };

    private Dictionary<SelectType, CFG_XMLConfig.XML_Config.SwitchOpne> SelectSwitchDic = new()
    {
        [SelectType.WuXing] = CFG_XMLConfig.XML_Config.SwitchOpne.WuXing,
        [SelectType.WarFight] = CFG_XMLConfig.XML_Config.SwitchOpne.WarFight,
    };

    private Dictionary<ChanllengeType, CFG_XMLConfig.XML_Config.SwitchOpne> ChanllengeToSwitchDic = new ()
    {
        [ChanllengeType.Fight] = CFG_XMLConfig.XML_Config.SwitchOpne.Play,
        [ChanllengeType.Boss] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss,
        [ChanllengeType.SLTower] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower,
        [ChanllengeType.Relic] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_Relice,
        [ChanllengeType.PvP] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP,
        [ChanllengeType.DayDo] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_QHSAndGold,
        [ChanllengeType.WorldBoss] = CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss,
        [ChanllengeType.FightEnd] = CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd,
        [ChanllengeType.HaiDao] = CFG_XMLConfig.XML_Config.SwitchOpne.DefendIsland,
    };

    //二级按钮对应的挑战类型
    private Dictionary<SelectType, ChanllengeType[]> SelectToChallengeDic = new ()
    {
        { SelectType.Challange, new ChanllengeType[] { ChanllengeType.Fight, ChanllengeType.Boss, ChanllengeType.SLTower, ChanllengeType.Relic, ChanllengeType.HaiDao } },
        { SelectType.Limit, new ChanllengeType[] { ChanllengeType.PvP, ChanllengeType.WorldBoss, ChanllengeType.FightEnd } },
        { SelectType.FuBen, new ChanllengeType[] { ChanllengeType.DayDo } }
    };


    private Dictionary<SelectType, GameObject> ReddotRootDic = new();
    private Dictionary<SelectType, GameObject> FunNewFlagRootDic = new();

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFightChoice>();
        });

        //主选择
        for (int i = 0; i < GetView().root_SelectBtns.childCount; i++)
        {
            int idx = i;
            JUIExtensions.btn(GetView().root_SelectBtns.GetChild(idx)).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnSlectBtnClick(idx);
            });
        }

        ReddotRootDic[SelectType.Challange] = GetView().root_SelectBtns.GetChild(0).Find("redDotRoot").gameObject;
        ReddotRootDic[SelectType.Limit] = GetView().root_SelectBtns.GetChild(1).Find("redDotRoot").gameObject;
        ReddotRootDic[SelectType.FuBen] = GetView().root_SelectBtns.GetChild(2).Find("redDotRoot").gameObject;
        ReddotRootDic[SelectType.WuXing] = GetView().root_SelectBtns.GetChild(3).Find("redDotRoot").gameObject;
        ReddotRootDic[SelectType.WarFight] = GetView().root_SelectBtns.GetChild(4).Find("redDotRoot").gameObject;

        FunNewFlagRootDic[SelectType.Challange] = GetView().root_SelectBtns.GetChild(0).Find("functionOpenNewFlag").gameObject;
        FunNewFlagRootDic[SelectType.Limit] = GetView().root_SelectBtns.GetChild(1).Find("functionOpenNewFlag").gameObject;
        FunNewFlagRootDic[SelectType.FuBen] = GetView().root_SelectBtns.GetChild(2).Find("functionOpenNewFlag").gameObject;
        FunNewFlagRootDic[SelectType.WuXing] = GetView().root_SelectBtns.GetChild(3).Find("functionOpenNewFlag").gameObject;
        FunNewFlagRootDic[SelectType.WarFight] = GetView().root_SelectBtns.GetChild(4).Find("functionOpenNewFlag").gameObject;

        //挑战二级页面点击
        for (int i = 0; i < GetView().root_ChallengeBtns.childCount; i++)
        {
            int idx = i;
            JUIExtensions.btn(GetView().root_ChallengeBtns.GetChild(idx)).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnChallengeBtnClick((int)SelectToChallengeDic[selectType][idx]);
            });
        }

        GetView().btn_ZYPlay.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameHelper.Instance.StarGame(BattleId);
            //TcpC.GetGameJoin(CFG_Room.GetByRoomType(CFG_Room.RoomType.PVE_NoPeace, RoleM.Instance.Level).ID, 0);
        });

        GetView().btn_FightAward.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIFightChapterAward>();
        });

        //试炼塔
        GetView().InfiniteScroll_SLTower.InitListView(0, OnFillSLTowerItem);

        GetView().btn_SLTower_Paly.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (FightPlayM.Instance.GetInfo().trial_tower_map.BattleId == NextSLTowerId)
            {
                UITool.ShowTips("已经达到最高层 ");
                return;
            }
            GameHelper.Instance.StarGame(NextSLTowerId);
        });
        //打开失试炼基金
        GetView().btn_SLTower_SLShop.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIShop>(null, CFG_Store.StoreType.SiLian);
        });

        GetView().btn_SLTower_SLJJ.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UITrialFund>();
        });

        //首领
        GetView().InfiniteScroll_SL.InitListView(0, OnFillBossItem) ;

        GetView().btn_SL_CountAdd.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.SLTZJZ);

            //if (StoreM.Instance.GetCanBuyNum((int)CFG_Store.StroeID.BossChange) == 0)
            //{
            //    UITool.ShowTips("今日购买已达到上限");
            //    return;
            //}
            //UITool.Open<UIGameSweepAndBuy>(null, UIGameSweepAndBuy.OpType.BossBuy);
        });

        GetView().btn_PVP_ZhanHun.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIShop>(null, CFG_Store.StoreType.Diamanta);
        });

        //遗迹
        GetView().btn_YJ_Add.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnRelicShow(RelicIdx +1);
        });
        GetView().btn_YJ_Jian.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnRelicShow(RelicIdx - 1);
        });

        GetView().btn_YJ_CountAdd.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.YJTZJZ);
            //if (StoreM.Instance.GetCanBuyNum((int)CFG_Store.StroeID.RelicChange) == 0)
            //{
            //    UITool.ShowTips("今日购买已达到上限");
            //    return;
            //}
            //UITool.Open<UIGameSweepAndBuy>(null, UIGameSweepAndBuy.OpType.RelicBuy);
        });

        GetView().btn_YJ_Play.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameHelper.Instance.StarGame(BattleId);
        });

        GetView().btn_YJ_SD.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIRelicChallange>(null, BattleId, RelicSlectID);
        });


        GetView().btn_help.onClick.AddListener(OnHelpOnClick);
    }

    private void HedeAllChallengeButtonRed()
    {
        for (int i = 0; i < GetView().root_ChallengeBtns.childCount; i++)
            GetView().root_ChallengeBtns.GetChild(i).Find("dot").gameObject.SetActive(false);

    }

    private void OnBindRed()
    {
        RedPointManager.Instance.RPS.Register(RedPointConst.Challenge.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().root_SelectBtns.GetChild(0).Find("redDotRoot/dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.Limet.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().root_SelectBtns.GetChild(1).Find("redDotRoot/dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.DayDo.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().root_SelectBtns.GetChild(2).Find("redDotRoot/dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.WuXing.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().root_SelectBtns.GetChild(3).Find("redDotRoot/dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.WarFight.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().root_SelectBtns.GetChild(4).Find("redDotRoot/dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.Challenge_Fight.ToString(), (node) =>
        {
            if (isNullGo || selectType != SelectType.Challange) return;
            GetView().root_ChallengeBtns.GetChild(0).Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.Challenge_Boss.ToString(), (node) =>
        {
            if (isNullGo || selectType != SelectType.Challange) return;
            GetView().root_ChallengeBtns.GetChild(1).Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.Challenge_Tower.ToString(), (node) =>
        {
            if (isNullGo || selectType != SelectType.Challange) return;
            GetView().root_ChallengeBtns.GetChild(2).Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.Challenge_Relice.ToString(), (node) =>
        {
            if (isNullGo || selectType != SelectType.Challange) return;
            GetView().root_ChallengeBtns.GetChild(3).Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.Challenge_OnHook.ToString(), (node) =>
        {
            if (isNullGo || selectType != SelectType.Challange) return;
            GetView().root_ChallengeBtns.GetChild(4).Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.Limet.ToString(), (node) =>
        {
            if (isNullGo || selectType != SelectType.Limit) return;
            GetView().root_ChallengeBtns.GetChild(0).Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

        RedPointManager.Instance.RPS.Register(RedPointConst.DayDo.ToString(), (node) =>
        {
            if (isNullGo || selectType != SelectType.FuBen) return;
            GetView().root_ChallengeBtns.GetChild(0).Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

        //RedPointManager.Instance.RPS.Register(RedPointConst.Challenge_Tower_Award.ToString(), (node) =>
        //{
        //    if (isNullGo || selectType != SelectType.Challange) return;
        //    GetView().btn_SLTower_SLJJ.transform.Find("dot").gameObject.SetActive(node.pointCount > 0);
        //});

        RedPointManager.Instance.RPS.Register(RedPointConst.Challenge_Fight_Award.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().btn_FightAward.transform.Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

    }


    private void OnHelpOnClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        int helpIdx = 0;
        switch (chanllengeType)
        {
            case ChanllengeType.Fight:
                helpIdx = 26;
                break;
            case ChanllengeType.Boss:
                helpIdx = 27;
                break;
            case ChanllengeType.SLTower:
                helpIdx = 28;
                break;
            case ChanllengeType.Relic:
                helpIdx = 29;
                break;
            case ChanllengeType.PvP:
                break;
            case ChanllengeType.DayDo:
                helpIdx = 18;
                break;
            default:
                break;
        }
        UITool.Open<UIComRule>(null, helpIdx);
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnSync_Event_Source_Refresh);

        selectType = args.Length > 0 ? (SelectType)args[0] : SelectType.Challange;
        chanllengeType = args.Length > 1 ? (ChanllengeType)args[1] : ChanllengeType.Fight;
        OnSlectBtnClick((int)selectType,true);
        SetTagButtonLockState();

        NewPlayerGuideWidget.Instance.AddGuideGameObject("FightChoice_FightButton", GetView().btn_ZYPlay.transform);

        UniGameEvent.AddEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);
        refreshFuncNewFlag();
    }

    private void SetTagButtonLockState()
    {
        for (int i = 0; i < GetView().root_SelectBtns.childCount; i++) 
        {
            Transform mask = GetView().root_SelectBtns.GetChild(i).Find("Mask");

            var switchOpne = CFG_XMLConfig.XML_Config.SwitchOpne.No;

            //直接跳转
            if (SelectSwitchDic.ContainsKey((SelectType)i))
            {
                switchOpne = SelectSwitchDic[(SelectType)i];
            }
            //内部跳转
            else if (SelectToChallengeDic.ContainsKey((SelectType)i) && ChanllengeToSwitchDic.ContainsKey(SelectToChallengeDic[(SelectType)i][0]))
            {
                switchOpne = ChanllengeToSwitchDic[SelectToChallengeDic[(SelectType)i][0]];
            }

            if (switchOpne != CFG_XMLConfig.XML_Config.SwitchOpne.No)
            {
                StringBuilder tips = new("");
                bool isLockOpenServer = GameUtil.Instance.CheckOpenServerIsLock(switchOpne, false);
                bool isLockLv = GameUtil.Instance.CheakIsLock(switchOpne, false, tips);
                mask.gameObject.SetActive(isLockOpenServer || isLockLv);

                if (isLockOpenServer)
                {
                    JUIExtensions.txt(mask.Find("txt")).text = $"{GameUtil.Instance.RemainingDaysUnlockOpen(switchOpne)}天后开启";
                }
                else if (isLockLv)
                {
                    JUIExtensions.txt(mask.Find("txt")).text = tips.ToString();//$"{CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[switchOpne]}级解锁";
                }
            }
            else
                mask.gameObject.SetActive(false);
        }
    }

    /// <summary>
    /// 一级按钮点击事件
    /// </summary>
    /// <param name="idx"></param>
    /// <param name="Jump"></param>
    private void OnSlectBtnClick(int idx,bool Jump = false)
    {
        if (!SelectToChallengeDic.ContainsKey((SelectType)idx) || SelectToChallengeDic[(SelectType)idx].Length <= 0)
        {
            if (SelectSwitchDic.ContainsKey((SelectType)idx) &&
                   (GameUtil.Instance.CheckOpenServerIsLock(SelectSwitchDic[(SelectType)idx]) || GameUtil.Instance.CheakIsLock(SelectSwitchDic[(SelectType)idx])))
                return;

            if (idx == (int)SelectType.WuXing)
            {
                UITool.Open<UIWuXingTrials>();
                UITool.Close<UIFightChoice>();
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.WuXing });
            }
            else if (idx == (int)SelectType.WarFight)
            {
                UITool.Open<UIWarFightMap>();
                UITool.Close<UIFightChoice>();

                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.WarFight });

                if (RedPointManager.Instance.GetRedNum(RedPointConst.WarFight.ToString()) != 0)
                {
                    RedPointManager.Instance.SetRedClose(RedPointConst.WarFight.ToString(), () => {
                        RedPointManager.Instance.RPS.Dispatch(RedPointConst.WarFight);
                    });
                }
            }
            //UITool.ShowTips("暂未开放");
            return;
        }

        if (!Jump)
            chanllengeType = SelectToChallengeDic[(SelectType)idx][0];

        if (ChanllengeToSwitchDic.ContainsKey(chanllengeType) && (GameUtil.Instance.CheckOpenServerIsLock(ChanllengeToSwitchDic[chanllengeType]) || GameUtil.Instance.CheakIsLock(ChanllengeToSwitchDic[chanllengeType])))
            return;

        //二级按钮显示状态

        for (int i = 0; i < GetView().root_ChallengeBtns.childCount; i++)
        {
            GetView().root_ChallengeBtns.GetChild(i).gameObject.SetActive(i < SelectToChallengeDic[(SelectType)idx].Length);
            if (i < SelectToChallengeDic[(SelectType)idx].Length)
            {
                var secoundKey = SelectToChallengeDic[(SelectType)idx][i];
                JUIExtensions.txt(GetView().root_ChallengeBtns.GetChild(i), "Label").text = ChanllengeToStrDic[secoundKey];
                StringBuilder tips = new("");
                bool isLock = ChanllengeToSwitchDic.ContainsKey(secoundKey) && (GameUtil.Instance.CheckOpenServerIsLock(ChanllengeToSwitchDic[secoundKey], false) || GameUtil.Instance.CheakIsLock(ChanllengeToSwitchDic[secoundKey], false, tips));
                GetView().root_ChallengeBtns.GetChild(i).Find("lock").gameObject.SetActive(isLock);
                if (isLock)
                {
                    JUIExtensions.txt(GetView().root_ChallengeBtns.GetChild(i), "lock/txt").text = GameUtil.Instance.CheckOpenServerIsLock(ChanllengeToSwitchDic[secoundKey], false) ?
                       GameUtil.Instance.RemainingDaysUnlockOpen(ChanllengeToSwitchDic[secoundKey]) + "天\n后开启" :
                        "主线\n解锁";//CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[ChanllengeToSwitchDic[secoundKey]] + "级\n解锁";
                }
            }
        }

        //默认选中
        selectType = (SelectType)idx;
        for (int i = 0; i < GetView().root_SelectBtns.childCount; i++)
        {
            GetView().root_SelectBtns.GetChild(i).Find("Select").gameObject.SetActive(idx == i);
        }

        //二级按钮选择
        OnChallengeBtnClick((int)chanllengeType);
    }

    private UIDayDo _UIDayDo;
    //二级页面点击
    private void OnChallengeBtnClick(int idx, bool isRefresh = false)
    {
        if (ChanllengeToSwitchDic.ContainsKey((ChanllengeType)idx) && (GameUtil.Instance.CheckOpenServerIsLock(ChanllengeToSwitchDic[(ChanllengeType)idx], !isRefresh) || GameUtil.Instance.CheakIsLock(ChanllengeToSwitchDic[(ChanllengeType)idx], !isRefresh)))
            return;

        chanllengeType = (ChanllengeType)idx;
        for (int i = 0; i < SelectToChallengeDic[selectType].Length; i++)
        {
            GetView().root_ChallengeBtns.GetChild(i).Find("Select").gameObject.SetActive(chanllengeType == SelectToChallengeDic[selectType][i]);
        }

        GetView().txt_title.text = ChanllengeToStrDic[chanllengeType].Replace("\n", "");
        GetView().btn_help.gameObject.SetActive(chanllengeType != ChanllengeType.PvP && chanllengeType != ChanllengeType.HaiDao  );
        GetView().obj_SL.SetActive(chanllengeType == ChanllengeType.Boss);
        GetView().obj_SLTower.SetActive(chanllengeType == ChanllengeType.SLTower);
        GetView().obj_ZY.SetActive(chanllengeType == ChanllengeType.Fight);
        GetView().obj_YJ.SetActive(chanllengeType == ChanllengeType.Relic);
        GetView().obj_PvP.SetActive(chanllengeType == ChanllengeType.PvP);

        GetView().Obj_HaiDao.gameObject.SetActive(chanllengeType == ChanllengeType.HaiDao);

        if (_UIDayDo != null &&  chanllengeType != ChanllengeType.DayDo)
            _UIDayDo.OnShow(false);

        switch (chanllengeType)
        {
            case ChanllengeType.Fight:
                OnFightUI();
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Play });
                RedPointManager.Instance.SetRedClose(RedPointConst.Challenge_Fight_Click.ToString(), () => {
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Fight_Click);
                });
                break;
            case ChanllengeType.Boss:
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Play_Boss });
                OnBossUI();
                break;
            case ChanllengeType.SLTower:
                OnSLTowerUI();
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Play_SlTower });
                RedPointManager.Instance.SetRedClose(RedPointConst.Challenge_Tower_Click.ToString(), () => {
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Tower_Click);
                });
                break;
            case ChanllengeType.Relic:
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Play_Relice });
                OnRelicUI(isRefresh);
                break;

            case ChanllengeType.PvP:
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Play_PVP, SwitchOpne.Play_GVG });
                OnPVPUI();
                break;
            case ChanllengeType.WorldBoss:
                if (!GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss))
                {
                    FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.WorldBoss });
                    UITool.Open<UIWorldBoss>();
                    UITool.Close<UIFightChoice>();
                }

                break;
            case ChanllengeType.FightEnd:
                if (!GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd))
                {
                    FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.FightEnd });
                    UITool.Open<UIFightToTheEnd>();
                    UITool.Close<UIFightChoice>();
                }
                break;
            case ChanllengeType.DayDo:
                FuncOpenShowWidget.Instance.HideFuncNewFlag(FuBenGroup);
                if (_UIDayDo == null)
                {
                    UITool.LoadUII<UIDayDo>(GetView().root_UI, (ui) => {
                         _UIDayDo = (UIDayDo)ui ;
                         _UIDayDo.OnShow(true);
                     });
                }
                else
                {
                    _UIDayDo.OnShow(true);
                }
                break;
            case ChanllengeType.HaiDao:
                OnHaiDaoUI();
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.DefendIsland });
                break;
            default:
                break;
        }
        HedeAllChallengeButtonRed();
        OnBindRed();
    }

    #region 挑战


    private int BattleId = 0;
    //战役页面
    private void OnFightUI()
    {
        var fightInfo = FightPlayM.Instance.GetInfo().battle_map;
        var CFGBattle = CFG_Battle.GetNextFightById(fightInfo.BattleId);
        for (int i = 0; i < GetView().root_ZYAward.childCount; i++)
        { 
            int idx = i;
            var go = GetView().root_ZYAward.GetChild(idx);
            go.gameObject.SetActive(idx < CFGBattle.Preview.Length);
            if (idx < CFGBattle.Preview.Length)
            {
                GameUtil.Instance.SetItemIcon(CFGBattle.Preview[idx], JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() => { 
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFGBattle.Preview[idx] });
                } );
            }
        }

        var CFGMonster = CFG_Monster.GetById((int)CFGBattle.MonsterShow);
        ResouseHelper.Instance.LoadSpriteByAtlas(CFGMonster.Photo, CFGMonster.Photoname, GetView().img_ZYBoss);
        GetView().txt_ZYName.text = "关卡" + CFGBattle.Chapter + "-" + CFGBattle.Layers;
        GetView().txt_ZYlimet.text = $"角色达到{(RoleM.Instance.Level < CFGBattle.Levellimit ? $"<color=#DA4E38>{CFGBattle.Levellimit}</color>" : CFGBattle.Levellimit)}级";
        GetView().btn_ZYPlay.interactable = RoleM.Instance.Level >= CFGBattle.Levellimit && fightInfo.BattleId != CFGBattle.UID;
        BattleId = CFGBattle.UID;
    }


    private List<CFGItem_Battle> fightList;
    private int NextSLTowerId = 0;
    //试炼塔
    private void OnSLTowerUI()
    {
        var fightInfo = FightPlayM.Instance.GetInfo().trial_tower_map;
        fightList = CFG_Battle.GetSLTowerList();

        NextSLTowerId = fightList[0].UID;
        int moveIdx = 0;
        for (int i = 0; i < fightList.Count; i++)
        {
            if (fightList[i].UID > fightInfo.BattleId)
            {
                moveIdx = i;
                NextSLTowerId = fightList[i].UID;
            }
        }

        if(GetView().InfiniteScroll_SLTower.ItemTotalCount != fightList.Count)
            GetView().InfiniteScroll_SLTower.SetListItemCount(fightList.Count);
        else
            GetView().InfiniteScroll_SLTower.RefreshAllShownItem();

        GetView().InfiniteScroll_SLTower.MovePanelToItemIndex(moveIdx, 0);

        //GetView().btn_SLTower_Paly.interactable = NextSLTowerId != fightInfo.BattleId;

    }

    //试炼塔Item
    private LoopListViewItem2 OnFillSLTowerItem(LoopListView2 listView,int index)
    {
        LoopListViewItem2 viewitem = listView.NewListViewItem("Item");
        Transform item = viewitem.transform;

        var cfgInfo = fightList[index];
        JUIExtensions.txt(item, "bg_top/txt_Name").text = "第" + cfgInfo.Layers + "层";
        JUIExtensions.txt(item, "txt_limit").text = GameUtil.Instance.FormatNumWithUnit(cfgInfo.Powerlimit);

        JUIExtensions.txt(item,"txt_limit").color = JUIExtensions.ToColor(cfgInfo.Powerlimit > RoleM.Instance.TotalCombat ? "#E96C5C" : "#88DF5F");
        JUIExtensions.outline(item, "txt_limit").effectColor = JUIExtensions.ToColor(cfgInfo.Powerlimit > RoleM.Instance.TotalCombat ? "#9E2522" : "#42712C");

        var root_award = item.transform.Find("awards");
        for (int i = 0; i < root_award.childCount; i++)
        {
            int idx = i;
            var go = root_award.GetChild(idx);
            go.gameObject.SetActive(idx < cfgInfo.Preview.Length);
            if (idx < cfgInfo.Preview.Length)
            {
                GameUtil.Instance.SetItemIcon(cfgInfo.Preview[idx], JUIExtensions.img(go,"icon"), JUIExtensions.img(go));
                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() => {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgInfo.Preview[idx] });
                });
            }
        }

        JUIExtensions.txt(item.transform, "txt_State").text = cfgInfo.UID > NextSLTowerId ? "" :(cfgInfo.UID < NextSLTowerId ? "已通关" :"当前关");
        item.transform.Find("Lock").gameObject.SetActive(cfgInfo.UID > NextSLTowerId);

        var CFGMonster = CFG_Monster.GetById((int)cfgInfo.MonsterShow);
        ResouseHelper.Instance.LoadSpriteByAtlas(CFGMonster.Icon, CFGMonster.Iconname, JUIExtensions.img(item.transform, "img_left/img_boss"));

        ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/FightChoice", cfgInfo.UID > NextSLTowerId ? "bg_titlt_0" : (cfgInfo.UID < NextSLTowerId ? "bg_titlt_2" : "bg_titlt_1"), JUIExtensions.img(item.transform, "img_left"));

        return viewitem;
    }

    private int NextBossId = 0;
    //首领
    private void OnBossUI()
    {
        var fightInfo = FightPlayM.Instance.GetInfo().boss_map;
        fightList = CFG_Battle.BossList;

        NextBossId = fightList[0].UID;
        int moveIdx = 0;
        for (int i = 0; i < fightList.Count; i++)
        {
            if (fightList[i].UID > fightInfo.BattleId)
            {
                moveIdx = i;
                NextBossId = fightList[i].UID;
            }
        }
        if (GetView().InfiniteScroll_SL.ItemTotalCount != fightList.Count)
            GetView().InfiniteScroll_SL.SetListItemCount(fightList.Count);

        GetView().InfiniteScroll_SL.RefreshAllShownItem();

        GetView().InfiniteScroll_SL.MovePanelToItemIndex(moveIdx,0);

        GetView().txt_SL_Count.text = "今日剩余次数:" + fightInfo.ChallengeTime;
    }

    //首领Item
    private LoopListViewItem2 OnFillBossItem(LoopListView2 loopListView,int index)
    {
        LoopListViewItem2 viewitem = loopListView.NewListViewItem("Item");
        GameObject item = viewitem.gameObject;

        var cfgInfo = fightList[index];
        JUIExtensions.txt(item.transform, "txt_Name").text = CFG_Monster.GetById(cfgInfo.Monster[0].id).Name;
        ResouseHelper.Instance.LoadSpriteByAtlas(CFG_Monster.GetById(cfgInfo.Monster[0].id).Icon, CFG_Monster.GetById(cfgInfo.Monster[0].id).Iconname, JUIExtensions.img(item.transform, "bg_head/img"));

        bool isLimit = RoleM.Instance.TotalCombat >= cfgInfo.Powerlimit;
        
        JUIExtensions.txt(item.transform, "txt_limit").text = "战力需求:"+ (isLimit ?  GameUtil.Instance.FormatNumWithUnit(cfgInfo.Powerlimit): $"<color=#DA4E38>{GameUtil.Instance.FormatNumWithUnit(cfgInfo.Powerlimit)}</color>");
        JUIExtensions.txt(item.transform, "txt_State").text = isLimit ? "":"战力未达标";

        var root_award = item.transform.Find("awards");
        for (int i = 0; i < root_award.childCount; i++)
        {
            int idx = i;
            var go = root_award.GetChild(idx);
            go.gameObject.SetActive(idx < cfgInfo.Preview.Length);
            if (idx < cfgInfo.Preview.Length)
            {
                GameUtil.Instance.SetItemIcon(cfgInfo.Preview[idx], JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() => {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgInfo.Preview[idx] });
                });
            }
        }

        //ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/FightChoice", cfgInfo.UID > NextSLTowerId ? "bg_titlt_0" : (cfgInfo.UID < NextSLTowerId ? "bg_titlt_2" : "bg_titlt_1"), JUIExtensions.img(item.transform, "img_left"));

        item.transform.Find("btn_Play").gameObject.SetActive(cfgInfo.UID == NextBossId);
        item.transform.Find("btn_sd").gameObject.SetActive(cfgInfo.UID < NextBossId);

        JUIExtensions.btn(item.transform.Find("btn_sd")).onClick.RemoveAllListeners();
        JUIExtensions.btn(item.transform.Find("btn_Play")).onClick.RemoveAllListeners();
        JUIExtensions.btn(item.transform.Find("btn_Play")).onClick.AddListener( () =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameHelper.Instance.StarGame(cfgInfo.UID);
        });

        JUIExtensions.btn(item.transform.Find("btn_sd")).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (FightPlayM.Instance.GetInfo().boss_map.ChallengeTime <= 0)
            {
                UITool.ShowTips("挑战次数不足");
                return;
            }
            UITool.Open<UIGameSweepAndBuy>(null, UIGameSweepAndBuy.OpType.BossSweep, cfgInfo.UID);
            //TcpC.C2S_CopySweep(cfgInfo.UID,1);
        });

        return viewitem;
    }

    private int RelicSlectID = 0;
    //遗迹页面
    private void OnRelicUI(bool isRefsh)
    {
        var fightInfo = FightPlayM.Instance.GetInfo().relic_map;
        GetView().txt_YJ_Num.text = "今日剩余次数:"+ fightInfo.ChallengeTime;
        int cfgID = fightInfo.BattleId == 0 ? CFG_Battle.RelicDic[1][0].UID : fightInfo.BattleId;

        RelicSlectID = cfgID;
        if (fightInfo.BattleId != 0)
        {
            for (int i = 0; i < CFG_Battle.RelicList.Count; i++)
            {
                if (RelicSlectID < CFG_Battle.RelicList[i].UID)
                {
                    RelicSlectID = CFG_Battle.RelicList[i].UID;
                    break;
                }
            }
        }
        
        OnRelicShow(isRefsh && RelicIdx != 0 ? RelicIdx: CFG_Battle.GetById(RelicSlectID).Chapter);
    }

    private int RelicIdx = 0;
    private void OnRelicShow(int pageidx)
    {
        RelicIdx = pageidx;
        GetView().txt_YJ_Name.text = "第" + RelicIdx + "章";
        var cfgPageInfo = CFG_Battle.RelicDic[pageidx];
        var fightInfo = FightPlayM.Instance.GetInfo().relic_map;

        int defaultID = 0;
        //获取每章默认选择
        if (pageidx <= CFG_Battle.GetById(RelicSlectID).Chapter )
        {
            for (int i = 0; i < cfgPageInfo.Count; i++)
            {
                if (FightPlayM.Instance.GetRelicLevelStar(cfgPageInfo[i].UID) < 3)
                {
                    defaultID = cfgPageInfo[i].UID;
                }
                if (RelicSlectID == cfgPageInfo[i].UID)
                {
                    defaultID = cfgPageInfo[i].UID;
                    break;
                }
            }
            defaultID = defaultID == 0 ? cfgPageInfo[cfgPageInfo.Count - 1].UID : defaultID;
        }
        
        BattleId = defaultID;

        int StarAllNum = 0;
        for (int i = 0; i < cfgPageInfo.Count; i++)
        {
            int battId = cfgPageInfo[i].UID;
            var go = GetView().root_YJLevel.Find("btn_zj_" + (i + 1));
            JUIExtensions.txt(go, "img_pro/txt_pro").text = cfgPageInfo[i].Chapter + "-" + cfgPageInfo[i].Layers;

            int StarNum = FightPlayM.Instance.GetRelicLevelStar(battId);
            StarAllNum += StarNum;
            for (int j = 0; j < 3; j++)
            {
                go.Find("root_star").GetChild(j).gameObject.SetActive(StarNum > j);
            }
            JUIExtensions.btn(go).onClick.RemoveAllListeners();
            ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/FightChoice", defaultID == battId ? "bg_yiji0" : "bg_yiji1", JUIExtensions.img(go));

            JUIExtensions.btn(go).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (battId > RelicSlectID)
                { 
                    UITool.ShowTips("尚未解锁");
                    return;
                }
                UITool.Open<UIRelicChallange>(null, battId, RelicSlectID);
            });
        }

        //遗迹进度条
        GetView().slider_YJ.value = StarAllNum / 30.0f;
        GetView().txt_YJ_StarNum.text = StarAllNum + "星";

        for (int i = 1; i <= 3; i++)
        {
            int StarNum = i * 10;
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, FightPlayM.Instance.GetRelicChapterStarGet(pageidx) >= StarNum  ? "icon_baoxiang1" : "icon_baoxiang2", JUIExtensions.img(GetView().slider_YJ.transform.Find("YJBox_" + i )));
            var btn = JUIExtensions.btn(GetView().slider_YJ.transform.Find("YJBox_" + i));
            GetView().slider_YJ.transform.Find("YJBox_" + i + "/dot").gameObject.SetActive(!(FightPlayM.Instance.GetRelicChapterStarGet(pageidx) >= StarNum) && StarAllNum >= StarNum);
            btn.onClick.RemoveAllListeners();
            btn.onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (!(FightPlayM.Instance.GetRelicChapterStarGet(pageidx) >= StarNum) && StarAllNum >= StarNum)
                    TcpC.C2S_BattleReward(CFG_Room.RoomType.Relic, cfgPageInfo[0].Chapter);
                else
                {
                    var StarInfo = cfgPageInfo[cfgPageInfo.Count -1 ].GetStarAwardByNum(StarNum);
                    List<UIBoxInfoShow.BoxItem[]> items = new List<UIBoxInfoShow.BoxItem[]>();
                    UIBoxInfoShow.BoxItem[] boxItems = new UIBoxInfoShow.BoxItem[StarInfo.Length];
                    for (int j = 0; j < StarInfo.Length; j++)
                    {
                        boxItems[j] = new UIBoxInfoShow.BoxItem { Id = StarInfo[j].id, Num = StarInfo[j].num };
                    }
                    items.Add(boxItems);
                    UITool.Open<UIBoxInfoShow>(null, null, items);
                }
            });
        }

        //遗迹按钮
        GetView().btn_YJ_Play.interactable = defaultID != 0;
        GetView().btn_YJ_Play.gameObject.SetActive(FightPlayM.Instance.GetRelicLevelStar(defaultID) != 3 && defaultID != 0);

        GetView().obj_YJLock.SetActive(!(FightPlayM.Instance.GetRelicLevelStar(defaultID) != 3 && defaultID != 0));

        GetView().btn_YJ_SD.gameObject.SetActive(FightPlayM.Instance.GetRelicLevelStar(defaultID) == 3);

        GetView().btn_YJ_Add.gameObject.SetActive(RelicIdx < CFG_Battle.RelicDic.Count);
        GetView().btn_YJ_Jian.gameObject.SetActive(RelicIdx > 1);
    }
    #endregion

    #region PVP

    private void OnPVPUI()
    {
        HttpC.GetRankList(CFG_RankingAward.RankClass.ZhanHun, (_r) =>
        {
            if (isNullGo)
                return;

            bool isNoOne = _r.data.rank_all.Length == 0;
            GetView().root_PvPOne.Find("bg_head/img_head").gameObject.SetActive(!isNoOne);
            GetView().root_PvPOne.Find("bg_head/txt_empty").gameObject.SetActive(isNoOne);

            if (!isNoOne)
            {
                //ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + _r.data.rank_all[0].job + _r.data.rank_all[0].sex,JUIExtensions.img(GetView().root_PvPOne, "bg_head/img_head"));
                (int headId, int frameId) = GameUtil.Instance.GetHeadAndFrameId(_r.data.rank_all[0].avatar_v2, _r.data.rank_all[0].frame);
                GameUtil.Instance.SetRoleHead(headId, _r.data.rank_all[0].job, JUIExtensions.img(GetView().root_PvPOne, "bg_head/img_head"), JUIExtensions.img(GetView().root_PvPOne, "bg_head/img_head_quality"));
                GameUtil.Instance.SetRoleHeadFrame(frameId, _r.data.rank_all[0].job, JUIExtensions.img(GetView().root_PvPOne, "bg_head/img_head_kuang"));
            }
            else
            {
                GameUtil.Instance.SetDefauleRoleHead(1,JUIExtensions.img(GetView().root_PvPOne, "bg_head/img_head_kuang"), JUIExtensions.img(GetView().root_PvPOne, "bg_head/img_head_quality"));
            }
        });

        HttpC.GetRankList(CFG_RankingAward.RankClass.TurnBasedDay, (_r) =>
        {
            if (isNullGo)
                return;

            bool isNoOne = _r.data.rank_all.Length == 0;
            GetView().root_PvPMore.Find("bg_head/img_head").gameObject.SetActive(!isNoOne);
            GetView().root_PvPMore.Find("bg_head/txt_empty").gameObject.SetActive(isNoOne);

            if (!isNoOne)
            {
                //ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + _r.data.rank_all[0].job + _r.data.rank_all[0].sex, JUIExtensions.img(GetView().root_PvPMore, "bg_head/img_head"));
                (int headId, int frameId) = GameUtil.Instance.GetHeadAndFrameId(_r.data.rank_all[0].avatar_v2, _r.data.rank_all[0].frame);
                GameUtil.Instance.SetRoleHead(headId, _r.data.rank_all[0].job, JUIExtensions.img(GetView().root_PvPMore, "bg_head/img_head"), JUIExtensions.img(GetView().root_PvPMore, "bg_head/img_head_quality"));
                GameUtil.Instance.SetRoleHeadFrame(frameId, _r.data.rank_all[0].job, JUIExtensions.img(GetView().root_PvPMore, "bg_head/img_head_kuang"));
            }
            else
            {
                GameUtil.Instance.SetDefauleRoleHead(1, JUIExtensions.img(GetView().root_PvPMore, "bg_head/img_head_kuang"), JUIExtensions.img(GetView().root_PvPMore, "bg_head/img_head_quality"));
            }
        });

        void OnItem(Transform root, CFGItem_Room cfg)
        {
            Transform itemRoot = root.Find("img_award/awards");

            //奖励
            CFG_XMLConfig.XML_Config.XML_Item_PVP.PVPreward[] pVPrewards = null;
            if (cfg.Roomtype == CFG_Room.RoomType.PVP_One)
            {
                pVPrewards = CFG_XMLConfig.xmlConfig.PVPInfo.winAwards;
            }
            else if(cfg.Roomtype == CFG_Room.RoomType.PVP_More)
            {
                pVPrewards = CFG_XMLConfig.xmlConfig.PVPInfo.loseAwards;
            }

            for (int i = 0; i < itemRoot.childCount; i++)
            {
                var go = itemRoot.GetChild(i);
                int cfgId = pVPrewards[i].id;
                GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(go.Find("icon")), JUIExtensions.img(go));

                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                });
            }
            //else
            //{
            //    int jidu = (GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).Month - 1) / 3 + 1;
            //    pVPrewards = CFG_XMLConfig.xmlConfig.GVGInfo.GVGReward[jidu];
            //}

            //挑战次数
            JUIExtensions.txt(root, "txt_Count").text = "" + RoleM.Instance.GetPVPMatchTimes((int)cfg.Roomtype);

            //CFG_XMLConfig.XML_Config.XML_Item_PVP.PVPseason endPVPSeason = null;
            /*if (cfg.Roomtype == CFG_Room.RoomType.PVP_One)
            {
                endPVPSeason = RoleM.Instance.GetPVPOneSeason();
                JUIExtensions.txt(root, "img_time/txt_sjTime").text = endPVPSeason == null ? "赛季已结束" : $"赛季截至日期:<color=#CAFB04>{endPVPSeason.endtime}</color>";
            }
            else
            {
                JUIExtensions.txt(root, "img_time/txt_sjTime").text = "";
            }*/
            //JUIExtensions.txt(root, "img_time/txt_sjTime").text = "";
            var baseInfo = cfg.Roomtype == CFG_Room.RoomType.PVP_One ? CFG_XMLConfig.xmlConfig.PVPInfo : CFG_XMLConfig.xmlConfig.GVGInfo;
            JUIExtensions.txt(root, "txt_sjTime").text = $"(每日:{baseInfo.dailyBegin}-{baseInfo.dailyEnd})";

            JUIExtensions.btn(root.Find("Btn_Recoed")).onClick.RemoveAllListeners();
            JUIExtensions.btn(root.Find("Btn_Recoed")).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                TcpC.C2S_PvpRecorderRequest(cfg.Roomtype == CFG_Room.RoomType.PVP_One ? 1 : 2);
            });

            JUIExtensions.btn(root.Find("Btn_help")).onClick.RemoveAllListeners();
            JUIExtensions.btn(root.Find("Btn_help")).onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIComRule>(null, (cfg.Roomtype == CFG_Room.RoomType.PVP_One ? 3 : 4));
            });

            //var isOpne = (endPVPSeason != null || cfg.Roomtype == CFG_Room.RoomType.PVP_More) && GameUtil.Instance.CheakCurTime(baseInfo.dailyBegin, baseInfo.dailyEnd);
            //var isOpne = GameUtil.Instance.CheakCurTime(baseInfo.dailyBegin, baseInfo.dailyEnd);

            var isOpne = true;
            //root.Find("Btn_Play/off").gameObject.SetActive(!isOpne);
            //root.Find("Btn_AddCount/off").gameObject.SetActive(!isOpne);

            //root.Find("Btn_Play/dot").gameObject.SetActive(isOpne && RoleM.Instance.GetPVPMatchTimes((int)cfg.Roomtype) > 0);

            StringBuilder tips = new("");
            Transform playbtnMask = root.Find("Btn_Play/mask");
            bool isLockOpenDay = GameUtil.Instance.CheckOpenServerIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false);
            bool isLockLevel = GameUtil.Instance.CheakIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false, tips);
            playbtnMask.gameObject.SetActive(isLockOpenDay || isLockLevel);
            root.Find("Btn_Play/dot").gameObject.SetActive(!isLockOpenDay && !isLockLevel && RoleM.Instance.GetPVPMatchTimes((int)cfg.Roomtype) > 0);

            if (isLockOpenDay)
            {
                int remainday = GameUtil.Instance.RemainingDaysUnlockOpen(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP);
                JUIExtensions.txt(playbtnMask.Find("txt_OpenDay")).text = $"{remainday}天后开启";
            }
            else if (isLockLevel)
            {
                int openLevel = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP];
                JUIExtensions.txt(playbtnMask.Find("txt_OpenDay")).text = tips.ToString(); //$"{openLevel}级解锁";
            }

            if (isOpne)
            {
                JUIExtensions.btn(root, "Btn_Play").onClick.RemoveAllListeners();
                JUIExtensions.btn(root, "Btn_Play").onClick.AddListener(() => {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                    if (GameUtil.Instance.CheckOpenServerIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                        return;

                    if (GameUtil.Instance.CheakIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                        return;

                    if (RoleM.Instance.GetPVPMatchTimes((int)cfg.Roomtype) == 0)
                    {
                        UITool.ShowTips("挑战次数不足");
                        return;
                    }
                    OnEnterPVPGame((int)cfg.Roomtype);
                });

                //root.Find("Btn_AddCount/off").gameObject.SetActive(RoleM.Instance.GetCanNumBuyPVP((int)cfg.Roomtype) == 0);
                JUIExtensions.btn(root, "Btn_AddCount").onClick.RemoveAllListeners();
                JUIExtensions.btn(root, "Btn_AddCount").onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    //UITool.Open<UIPVPCountBuy>(null, (int)cfg.Roomtype);

                    if (cfg.Roomtype == CFG_Room.RoomType.PVP_One)
                    {
                        if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false))
                        {
                            UITool.ShowTips("功能暂未开启");
                            return;
                        }

                        if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false))
                        {
                            UITool.ShowTips("功能暂未开启");
                            return;
                        }

                        GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.PVETZJZ);
                    }
                    else if(cfg.Roomtype == CFG_Room.RoomType.PVP_More)
                    {
                        if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                        {
                            UITool.ShowTips("功能暂未开启");
                            return;
                        }

                        if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                        {
                            UITool.ShowTips("功能暂未开启");
                            return;
                        }

                        GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.PVPTZJZ);
                    }
                });
            }

            JUIExtensions.btn(root, "bg_head").onClick.RemoveAllListeners();
            JUIExtensions.btn(root, "bg_head").onClick.AddListener(() =>
            {
                if (cfg.Roomtype == CFG_Room.RoomType.PVP_One)
                    UITool.Open<UIRank>(null, CFG_RankingAward.RankClass.ZhanHun);
                else if(cfg.Roomtype == CFG_Room.RoomType.PVP_More)
                    UITool.Open<UIRank>(null, CFG_RankingAward.RankClass.TurnBasedDay);
            });
        }

        var cfg = CFG_Room.GetByRoomType(CFG_Room.RoomType.PVP_One, RoleM.Instance.Level);
        OnItem(GetView().root_PvPOne, cfg);

        cfg = CFG_Room.GetByRoomType(CFG_Room.RoomType.PVP_More, RoleM.Instance.Level);
        OnItem(GetView().root_PvPMore, cfg);
    }

    private void OnEnterPVPGame(int Roomtype)
    {
        GameRoomM.Instance._GameState = GameRoomM.GameState.No;

        if (Roomtype == (int)CFG_Room.RoomType.PVP_One)
            TcpC.C2S_PvpMatching(1);
        else
            TcpC.C2S_PvpMatching(2);
        //UITool.ShowTips("暂未开放");
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_PVP(GameEvent.Sync_PVPInfo _e)
    {
        if (isNullGo)
            return;
        if ((GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false) || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false)))
            return;

        OnPVPUI();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_CheckOpenTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        SetTagButtonLockState();

        TcpC.C2S_InstanceInfo();
        TcpC.OnPVPInfo((int)pb.MatchType.SINGLE_COMBAT);
        TcpC.OnPVPInfo((int)pb.MatchType.FREE_FOR_ALL);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_Event_Source_Refresh(Event_Source_Refresh _e)
    {
        if (isNullGo) return;

        if (UITool.IsOpen<UICommonItemBuy>())
        {
            UICommonItemBuy uiUseItem = UITool.Get<UICommonItemBuy>();
            if (uiUseItem.GetItemID() == 56)
            {
                long canusenum = CFG_ItemUse.GetDailyUseLimit(56) - BagM.Instance.GetDailyUsedItem(56);
                uiUseItem.UpdateDescription2($"（今日还可使用：{canusenum}）");
            }
        }
    }

    #endregion

    #region 守卫海岛

    //TODO
    private void OnHaiDaoUI()
    {
        int idxID1 = 0;
        int idxID2 = 0;
        int playUid = CFG_Room.GetByRoomType(CFG_Room.RoomType.PVE_Peace, RoleM.Instance.Level) == null ? 1 : CFG_Room.GetByRoomType(CFG_Room.RoomType.PVE_Peace, RoleM.Instance.Level).ID;
        if (playUid > 5)
        {
            idxID1 = 5; 
            idxID2 = playUid;
        }
        else
        {
            idxID1 = playUid;
            idxID2 = 6;
        }

        ShowHaiDaoItem(idxID1, GetView().Obj_HaiDao.Find("root_1"), 1);
        ShowHaiDaoItem(idxID2, GetView().Obj_HaiDao.Find("root_2"), 2);

        void ShowHaiDaoItem(int uid,Transform root, int type)
        {
            var cfgInfo = CFG_Room.GetById(uid);
            JUIExtensions.txt(root, "txt_Limit_Des").text = cfgInfo.Suggest;
            JUIExtensions.txt(root, "img_Des/txt_des").text = cfgInfo.CopyDes;

            //奖励显示
            Transform itemRoot = root.Find("img_award/awards");
            for (int i = 0; i < itemRoot.childCount; i++)
            {
                var go = itemRoot.GetChild(i);
                go.gameObject.SetActive(i < cfgInfo.Droppreview.Length);

                if (i < cfgInfo.Droppreview.Length)
                {
                    int cfgId = cfgInfo.Droppreview[i];
                    GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
                    JUIExtensions.btn(go).onClick.RemoveAllListeners();
                    JUIExtensions.btn(go).onClick.AddListener(() =>
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                    });
                }
            }

            JUIExtensions.btn(root, "Btn_help").onClick.RemoveAllListeners();
            JUIExtensions.btn(root, "Btn_help").onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIComRule>(null, idxID1 == uid ? 40 :41);
            });

            int gametype = type;
            bool isNotOpenTime = GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.DefendIsland, false);
            bool isOpneLv = !isNotOpenTime && RoleM.Instance.Level >= cfgInfo.MinEntryLimit && RoleM.Instance.Level <= cfgInfo.MaxEntryLimit;

            int cdTime = isOpneLv ? CDTimeRemaining(gametype) : 0;
            root.Find("Btn_Play").gameObject.SetActive(isOpneLv && cdTime <= 0);
            root.Find("Btn_SD").gameObject.SetActive(isOpneLv && cdTime <= 0);
            bool showCD = isOpneLv && cdTime > 0;
            root.Find("txt_cdtime").gameObject.SetActive(showCD);

            if(showCD)
            {
                JUIExtensions.txt(root.Find("txt_cdtime")).text = $"{cdTime}秒后可进入";
                if (gametype == 1)
                {
                    ClearCDTimer1();
                    cdTimeItem1 = GameHelper.Instance.AddInvoke(new TimeItem(cdTime + 1, 1, () =>
                    {
                        if (isNullGo) return;
                        JUIExtensions.txt(root.Find("txt_cdtime")).text = $"{cdTime}秒后可进入";
                        if (cdTime == 0)
                        {
                            ShowHaiDaoItem(idxID1, GetView().Obj_HaiDao.Find("root_1"), 1);
                        }
                        cdTime--;
                    }));
                }
                else if (gametype == 2)
                {
                    ClearCDTimer2();
                    cdTimeItem2 = GameHelper.Instance.AddInvoke(new TimeItem(cdTime + 1, 1, () =>
                    {
                        if (isNullGo) return;
                        JUIExtensions.txt(root.Find("txt_cdtime")).text = $"{cdTime}秒后可进入";
                        if (cdTime == 0)
                        {
                            ShowHaiDaoItem(idxID2, GetView().Obj_HaiDao.Find("root_2"), 2);
                        }
                        cdTime--;
                    }));
                }
            }

            JUIExtensions.txt(root, "txt_ChanllengeTime").text = !isOpneLv ? "": "今日剩余时间:" + GameUtil.Instance.GetTimeSToHMS(FightPlayM.Instance.GetInfo().afk.DailyTime) ;

            root.Find("Btn_AddCount").gameObject.SetActive(isOpneLv);
            JUIExtensions.btn(root.Find("Btn_AddCount")).onClick.RemoveAllListeners();
            JUIExtensions.btn(root.Find("Btn_AddCount")).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.hangUpTime);
            });

            var _btn = JUIExtensions.btn(root.Find("Btn_Play"));
            _btn.onClick.RemoveAllListeners();
            _btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (FightPlayM.Instance.GetInfo().afk.DailyTime == 0)
                {
                    UITool.ShowTips("剩余时间不足！");
                    return;
                }

                if (FightPlayM.Instance.GetInfo().afk.AfkMultiplierID < 0)
                {
                    //UITool.ShowTips("挑战门票不足！");
                    UITool.Open<UIGuardIslandBox>();
                    return;
                }
                TcpC.GetGameJoin(CFG_Room.GetByRoomType(CFG_Room.RoomType.PVE_Peace, RoleM.Instance.Level).ID, 0);

                string lv = playUid > 5 ? "101~150级" : "50~100级";

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = $"挂机升级（{lv}）",
                    event_click = $"挑战",
                });
            });

            _btn = JUIExtensions.btn(root.Find("Btn_SD"));
            _btn.onClick.RemoveAllListeners();
            _btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                if (!PrivilegeM.Instance.CheckQuicklyIsiand())
                {
                    UITool.ShowTips("尚未解锁终身卡，不能扫荡！");
                    return;
                }

                if (FightPlayM.Instance.GetInfo().afk.DailyTime == 0)
                {
                    UITool.ShowTips("剩余时间不足！");
                    return;
                }

                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = "扫荡将按平均效率进行结算是否确认扫荡？",
                    YesFun = () => {
                        TcpC.C2S_CopySweep(CFG_Room.GetByRoomType(CFG_Room.RoomType.PVE_Peace, RoleM.Instance.Level).ID);

                        string lv = playUid > 5 ? "101~150级" : "50~100级";
                        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                        {
                            page = "玩法",
                            component = "挑战",
                            position = $"挂机升级（{lv}）",
                            event_click = $"扫荡",
                        });
                    },
                });
            });

            //CFGItem_Store cFGItem_Store = CFG_Store.GetById((int)CFG_Store.StroeID.HaiDaoCard);

            //root.Find("Btn_AddCount").gameObject.SetActive(isOpneLv && StoreM.Instance.GetCanBuyNum((int)CFG_Store.StroeID.HaiDaoCard) > 0);

            root.Find("obj_Limit").gameObject.SetActive(!isOpneLv);
            if (isNotOpenTime)
                JUIExtensions.txt(root.Find("obj_Limit")).text = $"{GameUtil.Instance.RemainingDaysUnlockOpen(CFG_XMLConfig.XML_Config.SwitchOpne.DefendIsland)}天后开启";
            else
                JUIExtensions.txt(root.Find("obj_Limit")).text = $"角色等级不满足需求";

            string cardStr = "";
            if (isOpneLv && FightPlayM.Instance.GetInfo().afk.AfkMultiplierID >= 0)
            {
                foreach (var item in CFG_XMLConfig.xmlConfig.GuardIslandDic.Values)
                {
                    if (item.id == (int)FightPlayM.Instance.GetInfo().afk.AfkMultiplierID)
                    {
                        cardStr = item.multiplier + "倍效率"; ;
                        break;
                    }
                }
            }

            JUIExtensions.txt(root, "txt_time").text = cardStr;

        }

        int CDTimeRemaining(int type)
        {
            string cdTimeKey = GameConst.PlayerPrefsKey.HangUpCDTime.ToString() + type + "_" + GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear;
            long cdTime = GameUtil.Instance.GetPlayerPrefsLong(cdTimeKey, 0);
            return Mathf.Max(0, CFG_XMLConfig.xmlConfig.GuardIsland.CDTime - (int)(GameUtil.Instance.GetCurrentTimeSeconds() - cdTime));
        }
    }

    TimeItem cdTimeItem1, cdTimeItem2;

    public void ClearCDTimer1()
    {
        if(GameHelper.Instance!=null&& cdTimeItem1 != null)
        {
            GameHelper.Instance.RemoveInvoke(cdTimeItem1);
            cdTimeItem1 = null;
        }
    }

    public void ClearCDTimer2()
    {
        if(GameHelper.Instance!=null&& cdTimeItem2 != null)
        {
            GameHelper.Instance.RemoveInvoke(cdTimeItem2);
            cdTimeItem2 = null;
        }
    }

    #endregion

    [Subscriber(ThreadMode.Main)]
    private void Refresh_InstanceInfo(GameEvent.Refresh_InstanceInfo _e)
    {
        if (isNullGo) return;
        if (selectType == SelectType.Challange || chanllengeType == ChanllengeType.HaiDao)
            OnChallengeBtnClick((int)chanllengeType,true);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnSync_Event_Source_Refresh);
        UniGameEvent.RemoveEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);

        if (_UIDayDo != null)
            _UIDayDo.OnExit();

        ClearCDTimer1();
        ClearCDTimer2();

        if(ActivityM.Instance.GetAcitvityEndTime(ActivityType.CopyFund) > 0)
        {
            HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.CopyFund);
        }
    }

    public override void OnDestroy()
    {
        OnExit();
    }

    #region 功能开启角标

    private void OnRefresh_FunctionNewFlag(UIFunctionOpenOperate operate)
    {
        refreshFuncNewFlag();
    }

    private List<SwitchOpne> ChallangeGroup = new() { SwitchOpne.Play, SwitchOpne.Play_Boss, SwitchOpne.Play_SlTower, SwitchOpne.Play_Relice, SwitchOpne.DefendIsland };
    private List<SwitchOpne> LimitGroup = new() { SwitchOpne.Play_PVP, SwitchOpne.Play_GVG, SwitchOpne.WorldBoss, SwitchOpne.FightEnd };
    private List<SwitchOpne> FuBenGroup = new() { SwitchOpne.Play_QHSAndGold, SwitchOpne.Play_ExpAndSXS, SwitchOpne.Play_SY, SwitchOpne.Play_HS, SwitchOpne.ExpDungeon, SwitchOpne.GoldDungeon };
    private List<SwitchOpne> WuXingGroup = new() { SwitchOpne.WuXing };
    private List<SwitchOpne> WarFightGroup = new() { SwitchOpne.WarFight };
    private void refreshFuncNewFlag()
    {
        bool challangeFlag = FuncOpenShowWidget.Instance.FunctionIsLight(ChallangeGroup);
        ReddotRootDic[SelectType.Challange].SetActive(!challangeFlag);
        FunNewFlagRootDic[SelectType.Challange].SetActive(challangeFlag);

        if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Play_PVP, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Play_PVP, false)))
        {
            bool limitFlag = FuncOpenShowWidget.Instance.FunctionIsLight(LimitGroup);
            ReddotRootDic[SelectType.Limit].SetActive(!limitFlag);
            FunNewFlagRootDic[SelectType.Limit].SetActive(limitFlag);
        }

        if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Play_QHSAndGold, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Play_QHSAndGold, false)))
        {
            bool fubenFlag = FuncOpenShowWidget.Instance.FunctionIsLight(FuBenGroup);
            ReddotRootDic[SelectType.FuBen].SetActive(!fubenFlag);
            FunNewFlagRootDic[SelectType.FuBen].SetActive(fubenFlag);
        }

        if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.WuXing, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.WuXing, false)))
        {
            bool wuxingFlag = FuncOpenShowWidget.Instance.FunctionIsLight(WuXingGroup);
            ReddotRootDic[SelectType.WuXing].SetActive(!wuxingFlag);
            FunNewFlagRootDic[SelectType.WuXing].SetActive(wuxingFlag);
        }

        if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.WarFight, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.WarFight, false)))
        {
            bool warFightFlag = FuncOpenShowWidget.Instance.FunctionIsLight(WarFightGroup);
            ReddotRootDic[SelectType.WarFight].SetActive(!warFightFlag);
            FunNewFlagRootDic[SelectType.WarFight].SetActive(warFightFlag);
        }
    }
    #endregion
}