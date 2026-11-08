using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Tsp;

using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using SuperScrollView;

using System.Collections;
using System.Collections.Generic;
using System.Linq;

using UnityEngine;
using UnityEngine.UI;

using static GameEvent;

public class UIRank : UIBase<UIModelRank, UIViewRank>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIRank>();
        });

        for (int i = 0; i < GetView().root_btns.childCount; i++)
        {
            var trans = GetView().root_btns.GetChild(i);
            int idx = i;
            trans.GetComponent<Button>().onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                if (selectRank == rankClasses[idx]) return;

                if (RankToSwitchDic.ContainsKey(rankClasses[idx]) && 
                (GameUtil.Instance.CheckOpenServerIsLock(RankToSwitchDic[rankClasses[idx]], true) || GameUtil.Instance.CheakIsLock(RankToSwitchDic[rankClasses[idx]], true)))
                    return;
                OnRankShow(rankClasses[idx]);
            });
        }

        GetView().InfiniteScroll_Rank.InitListView(0, OnFillRank);
        OnBindRed();
    }

    private CFG_RankingAward.RankClass[] rankClasses =
    {
        CFG_RankingAward.RankClass.Power,
        CFG_RankingAward.RankClass.ZhanHun,//即时赛 PVP
        CFG_RankingAward.RankClass.TurnBasedDay,//即时赛 PVP
        CFG_RankingAward.RankClass.WuXing,//即时赛 PVP
    };

    //锁定条件绑定
    private Dictionary<CFG_RankingAward.RankClass, CFG_XMLConfig.XML_Config.SwitchOpne> RankToSwitchDic = new ()
    {
        [CFG_RankingAward.RankClass.ZhanHun] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP,
        [CFG_RankingAward.RankClass.TurnBasedDay] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG,
    };

    private Dictionary<CFG_RankingAward.RankClass, CFG_RankingAward.RankClass> RankAwardWeekDic = new Dictionary<CFG_RankingAward.RankClass, CFG_RankingAward.RankClass>
    {
        [CFG_RankingAward.RankClass.Power] = CFG_RankingAward.RankClass.PowerWeek,
        [CFG_RankingAward.RankClass.ZhanHun] = CFG_RankingAward.RankClass.RealTimeWeek,
        [CFG_RankingAward.RankClass.TurnBasedDay] = CFG_RankingAward.RankClass.TurnBasedWeek,
    };

    // 红点绑定
    private Dictionary<CFG_RankingAward.RankClass, RedPointConst> RedDic = new Dictionary<CFG_RankingAward.RankClass, RedPointConst>
    {
        [CFG_RankingAward.RankClass.Power] = RedPointConst.Rank_Power,
        [CFG_RankingAward.RankClass.ZhanHun] = RedPointConst.Rank_ZhanHun,
        [CFG_RankingAward.RankClass.TurnBasedDay] = RedPointConst.Rank_HuiHe,
    };

    public override void OnLayerChange()
    {
        GameUtil.Instance.SetEffectLayer(GetView().root_eff, _canvasOrder + 1);
    }

    private void OnBtnLock()
    {
        for (int i = 0; i < GetView().root_btns.childCount; i++)
        {
            var trans = GetView().root_btns.GetChild(i);
            var LockKey = rankClasses[i];

            System.Text.StringBuilder tips = new("");
            trans.Find("Lock").gameObject.SetActive(
                RankToSwitchDic.ContainsKey(LockKey) &&
                ((GameUtil.Instance.CheckOpenServerIsLock(RankToSwitchDic[LockKey], false) || GameUtil.Instance.CheakIsLock(RankToSwitchDic[LockKey], false, tips)))
            );
            if (RankToSwitchDic.ContainsKey(LockKey))
            {
                JUIExtensions.txt(trans.Find("Lock/Label")).text = GameUtil.Instance.CheckOpenServerIsLock(RankToSwitchDic[LockKey], false) ?
                    $"{GameUtil.Instance.RemainingDaysUnlockOpen(RankToSwitchDic[LockKey])}天后开启"
                    : tips.ToString();// CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[RankToSwitchDic[LockKey]] + "级解锁";
            }
        }
    }

    private void OnBindRed()
    {
        for (int i = 0; i < RedDic.Count; i++)
        {
            int idx = i;
            RedPointManager.Instance.RPS.Register(RedDic.ElementAt(idx).Value.ToString(), (node) =>
            {
            if (isNullGo)
                return;

            bool islock = false;
                if (RedDic.ElementAt(idx).Key == CFG_RankingAward.RankClass.ZhanHun)
                    islock = GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false) || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false);
                GetView().root_btns.GetChild(idx).Find("dot").gameObject.SetActive(node.pointCount > 0);
            });
        }
    }

    private Dictionary<int, CFGItem_RankingAward> DicCfgRank;
    private HTTP_RankDetail RankListInfo;

    private CFG_RankingAward.RankClass selectRank = CFG_RankingAward.RankClass.Power;

    private List<UIBoxInfoShow.Buttons> awardTypesButton = new List<UIBoxInfoShow.Buttons>
    {
        new UIBoxInfoShow.Buttons { Name = "日结奖励" },
        new UIBoxInfoShow.Buttons { Name = "周结奖励" },
    };

    private void OpenAwardUI(int Rankidx)
    {
        if (Rankidx <= 0)
            return;

        // 日结奖励
        List<UIBoxInfoShow.BoxItem[]> items = new List<UIBoxInfoShow.BoxItem[]>();
        UIBoxInfoShow.BoxItem[] boxItemsDay = new UIBoxInfoShow.BoxItem[DicCfgRank[Rankidx].award.Length];
        for (int j = 0; j < DicCfgRank[Rankidx].award.Length; j++)
        {
            boxItemsDay[j] = new UIBoxInfoShow.BoxItem { Id = int.Parse(DicCfgRank[Rankidx].award[j].ElementAt(0).Key), Num = DicCfgRank[Rankidx].award[j].ElementAt(0).Value };
        }

        // 周结奖励
        var _DicCfgRank = CFG_RankingAward.GetRankDicByRankClass(RankAwardWeekDic[selectRank]);
        UIBoxInfoShow.BoxItem[] boxItemsWeek = new UIBoxInfoShow.BoxItem[_DicCfgRank[Rankidx].award.Length];
        for (int j = 0; j < _DicCfgRank[Rankidx].award.Length; j++)
        {
            boxItemsWeek[j] = new UIBoxInfoShow.BoxItem { Id = int.Parse(_DicCfgRank[Rankidx].award[j].ElementAt(0).Key), Num = _DicCfgRank[Rankidx].award[j].ElementAt(0).Value };
        }

        items.Add(boxItemsDay);
        items.Add(boxItemsWeek);

        string tip = "提示：日结奖励每晚22点发放；周结奖励每周一0点发放。";
        UITool.Open<UIBoxInfoShow>(null, awardTypesButton, items, tip);
    }

    private void OnRankShow(CFG_RankingAward.RankClass rankClass)
    {
        bool hasKey = RedDic.ContainsKey(rankClass);
        if (hasKey)
        {
            RedPointManager.Instance.SetRedClose(RedDic[rankClass].ToString(), () =>
            {
                RedPointManager.Instance.RPS.Dispatch(RedDic[rankClass]);
            });
        }

        HttpC.GetRankList(rankClass, (_r) =>
        {
            if (isNullGo)
                return;

            selectRank = rankClass;

            if(selectRank != CFG_RankingAward.RankClass.WuXing)
            {
                for (int i = 0; i < GetView().root_btns.childCount; i++)
                {
                    GetView().root_btns.GetChild(i).Find("Select").gameObject.SetActive(selectRank == rankClasses[i]);
                }

                DicCfgRank = CFG_RankingAward.GetRankDicByRankClass(rankClass);
            }

            RankListInfo = _r.data;

            OnRefreshUIShow();

            //RedDotM.Instance.SetRandRedNumByType(RedDic[rankClass].ToString(), () =>
            //{
            //    RedPointManager.Instance.RPS.Dispatch(RedDic[rankClass].ToString(), RedDotM.Instance.GetRandRedNumByType(RedDic[rankClass].ToString()));
            //});
        });

    }

    private void OnRefreshUIShow()
    {
        for (int i = 0; i < 3; i++)
        {
            HTTP_RankUserInfo userInfo = RankListInfo.rank_all.Length > i ? RankListInfo.rank_all[i] : default;
            OnRankItemShow(userInfo, (i + 1), GetView().root_top3.Find("Rank_" + (i + 1)).gameObject, RankListInfo.rank_all.Length <= i);
        }

        if (selectRank == CFG_RankingAward.RankClass.WuXing)
            OnRankItemShow(RankListInfo.rank_own, RankListInfo.rank_own.num, GetView().root_MyRank2.gameObject);
        else
            OnRankItemShow(RankListInfo.rank_own, RankListInfo.rank_own.num, GetView().root_MyRank1.gameObject);

        GetView().InfiniteScroll_Rank.SetListItemCount(97, true);
        GetView().InfiniteScroll_Rank.RefreshAllShownItem();
        GetView().InfiniteScroll_Rank.MovePanelToItemIndex(0, 0);
    }

    private void OnRankItemShow(HTTP_RankUserInfo useInfo, int Rankidx, GameObject go, bool isNull = false)
    {
        JUIExtensions.txt(go.transform.Find("txt_Rank")).text = Rankidx <= 0 ? "未上榜" : Rankidx.ToString();

        bool isWuXing = selectRank == CFG_RankingAward.RankClass.WuXing;
        Transform rootAward = go.transform.Find("root_award");
        if (rootAward)
        {
            rootAward.gameObject.SetActive(!isWuXing);
        }

        Transform rootPower = go.transform.Find("root_power");
        if (rootPower)
        {
            rootPower.gameObject.SetActive(isWuXing);
        }

        if (isWuXing)
        {
            JUIExtensions.txt(rootPower, "txt_Num").text = GameUtil.Instance.FormatPowerWithUnit(useInfo.value % 100000000000);
        }
        else
        {
            for (int i = 0; i < 1; i++)
            {
                rootAward.gameObject.SetActive(!isNull && DicCfgRank.ContainsKey(Rankidx));
                if (rootAward.Find("award_" + i) != null && DicCfgRank.ContainsKey(Rankidx))
                {
                    rootAward.Find("award_" + i).gameObject.SetActive(Rankidx > 0 && DicCfgRank[Rankidx].award.Length > i);
                    if (Rankidx > 0 && DicCfgRank[Rankidx].award.Length > i)
                    {
                        GameUtil.Instance.SetItemIcon(int.Parse(DicCfgRank[Rankidx].award[i].ElementAt(0).Key), JUIExtensions.img(rootAward.Find("award_" + i + "/icon")));
                        JUIExtensions.txt(rootAward.Find("award_" + i + "/txt_Num")).text = GameUtil.Instance.FormatNumWithUnit(DicCfgRank[Rankidx].award[i].ElementAt(0).Value);
                    }
                }
            }
        }

        go.transform.Find("HeadBG").gameObject.SetActive(!isNull);

        go.transform.Find("img_des").gameObject.SetActive(!isNull);

        if (go.transform.Find("emptySeat"))
            go.transform.Find("emptySeat").gameObject.SetActive(isNull);

        if(rootPower)
            rootPower.gameObject.SetActive(!isNull);

        if (isNull)
        {
            JUIExtensions.txt(go.transform.Find("txt_Name")).text = "";
            JUIExtensions.txt(go.transform.Find("img_des/txt_des")).text = "";
        }
        else
        {
            var HeadBtn = JUIExtensions.btn(go.transform.Find("HeadBG"));
            HeadBtn.onClick.RemoveAllListeners();
            HeadBtn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (useInfo.uid > 0 && useInfo.uid != RoleM.Instance.Uid)
                {
                    TcpC.C2S_ObservePlayerInfoRequest(useInfo.uid);
                }
            });
            JUIExtensions.txt(go.transform.Find("txt_Name")).text = useInfo.name;
            JUIExtensions.txt(go.transform.Find("img_des/txt_des")).text = GetStrValueDes(useInfo.value);

            int job = useInfo.job;
            int sex = useInfo.sex;
            if (Rankidx == 0)
            {
                job = RoleM.Instance.Job;
                sex = RoleM.Instance.Sex;
                JUIExtensions.txt(go.transform.Find("txt_Name")).text = RoleM.Instance.Nickname;
            }
            //ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + job + sex, JUIExtensions.img(go.transform, "HeadBG/img_head"));

            (int headId, int frameId) = GameUtil.Instance.GetHeadAndFrameId(useInfo.avatar_v2, useInfo.frame);
            GameUtil.Instance.SetRoleHead(headId, job, JUIExtensions.img(go.transform, "HeadBG/img_head"), JUIExtensions.img(go.transform, "HeadBG/img_head_quality"));
            GameUtil.Instance.SetRoleHeadFrame(frameId, job, JUIExtensions.img(go.transform, "HeadBG/img_head_kuang"));
        }

        if (rootAward)
        {
            JUIExtensions.btn(go.transform.Find("root_award/btn_info")).onClick.RemoveAllListeners();
            JUIExtensions.btn(go.transform.Find("root_award/btn_info")).onClick.AddListener(delegate
            {
                OpenAwardUI(Rankidx);
            });
        }
    }

    private string GetStrValueDes(long value)
    {
        string str = "";
        switch (selectRank)
        {
            case CFG_RankingAward.RankClass.Level:
                str = "等级  ";
                break;
            case CFG_RankingAward.RankClass.Power:
                str = "战力  ";
                break;
            case CFG_RankingAward.RankClass.ZhanHun:
            case CFG_RankingAward.RankClass.TurnBasedDay:
                str = "积分  ";
                break;
            case CFG_RankingAward.RankClass.SiLian:
            case CFG_RankingAward.RankClass.WuXing:
                str = "最高层数  ";
                break;
            default:
                break;
        }

        if (selectRank == CFG_RankingAward.RankClass.Power)
            str += GameUtil.Instance.FormatPowerWithUnit(value);
        else if (selectRank == CFG_RankingAward.RankClass.WuXing)
        {
            long layer = value / 100000000000;
            str += layer;
        }
        else
            str += GameUtil.Instance.FormatNumWithUnit(value);
        return str;
    }

    private LoopListViewItem2 OnFillRank(LoopListView2 loopListView, int idx)
    {
        string itemRefName = selectRank == CFG_RankingAward.RankClass.WuXing ? "Item2" : "Item1";
        LoopListViewItem2 item = loopListView.NewListViewItem(itemRefName);
        GameObject go = item.gameObject;
        int rankIdx = idx + 3;
        HTTP_RankUserInfo userInfo = RankListInfo.rank_all.Length > rankIdx ? RankListInfo.rank_all[rankIdx] : default;
        OnRankItemShow(userInfo, rankIdx + 1, go, RankListInfo.rank_all.Length <= rankIdx);
        return item;
    }

    private TimeItem RefreshUILoop;

    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(Refresh_ZeroTime);
        UniGameEvent.AddEventListener<GameEvent.Refresh_LvUp>(Refresh_LvUp);

        OnBtnLock();
        if (args.Length > 0)
            selectRank = (CFG_RankingAward.RankClass)(args[0]);
        else
            selectRank = CFG_RankingAward.RankClass.Power;

        if (RankToSwitchDic.ContainsKey(selectRank) && (GameUtil.Instance.CheckOpenServerIsLock(RankToSwitchDic[selectRank], false) ||GameUtil.Instance.CheakIsLock(RankToSwitchDic[selectRank], false)))
            selectRank = CFG_RankingAward.RankClass.Power;

        bool isWuXing = selectRank == CFG_RankingAward.RankClass.WuXing;
        GetView().root_title.gameObject.SetActive(isWuXing);
        GetView().root_btns.gameObject.SetActive(!isWuXing);
        GetView().root_MyRank1.gameObject.SetActive(!isWuXing);
        GetView().root_MyRank2.gameObject.SetActive(isWuXing);

        if(isWuXing)
            GetView().txt_Desc.text = $"榜单每5分钟更新一次。每周一0点榜单进行重置。";
        else
            GetView().txt_Desc.text = $"榜单每5分钟更新一次。日结奖励每晚22点发放；周结奖励每周一0点发放。";

        RefreshUILoop = GameHelper.Instance.AddInvoke(new TimeItem(0, 300, () =>
            {
                if (isNullGo)
                    return;

                OnRankShow(selectRank);
            }, _isLoop: true));
    }

    //[Subscriber(ThreadMode.Main)]
    private void Refresh_ZeroTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        OnBtnLock();
    }

    //[Subscriber(ThreadMode.Main)]
    private void Refresh_LvUp(GameEvent.Refresh_LvUp _e)
    {
        if (isNullGo) return;
        OnBtnLock();
    }


    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(Refresh_ZeroTime);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_LvUp>(Refresh_LvUp);

        GameHelper.Instance?.RemoveInvoke(RefreshUILoop);
        RefreshUILoop = null;
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}