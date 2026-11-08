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

public class UILimitetRank : UIBase<UIModelLimitetRank, UIViewLimitetRank>
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
            UITool.Close<UILimitetRank>();
        });

        /*
        for (int i = 0; i < GetView().root_btns.childCount; i++)
        {
            var trans = GetView().root_btns.GetChild(i);
            int idx = i;
            trans.GetComponent<Button>().onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (selectRank == rankClasses[idx]) return;
                OnRankShow(rankClasses[idx]);
            });
        }
        */

        GetView().InfiniteScroll_Rank.InitListView(0, OnFillRank);
    }

    private CFG_RankingAward.RankClass[] rankClasses =
    {
        CFG_RankingAward.RankClass.Fight,
        CFG_RankingAward.RankClass.MainPass,
    };


    public override void OnLayerChange()
    {
        GameUtil.Instance.SetEffectLayer(GetView().root_eff, _canvasOrder + 1);
    }

   

    private Dictionary<int, CFGItem_RankingAward> DicCfgRank;
    private HTTP_RankDetail RankListInfo;

    private CFG_RankingAward.RankClass selectRank = CFG_RankingAward.RankClass.Power;

    private List<UIBoxInfoShow.Buttons> awardTypesButton = new List<UIBoxInfoShow.Buttons>
    {
    };

    private void OpenAwardUI(int Rankidx)
    {
        if (Rankidx <= 0)
            return;

        UIBoxInfoShow.BoxItem[] boxItemsDay = new UIBoxInfoShow.BoxItem[DicCfgRank[Rankidx].award.Length];
        for (int j = 0; j < DicCfgRank[Rankidx].award.Length; j++)
        {
            boxItemsDay[j] = new UIBoxInfoShow.BoxItem { Id = int.Parse(DicCfgRank[Rankidx].award[j].ElementAt(0).Key), Num = DicCfgRank[Rankidx].award[j].ElementAt(0).Value };
        }

        List<UIBoxInfoShow.BoxItem[]> items = new (){ boxItemsDay };

        UITool.Open<UIBoxInfoShow>(null, awardTypesButton, items, string.Empty);
    }

    private void OnRankShow(CFG_RankingAward.RankClass rankClass)
    {
        HttpC.GetRankList(rankClass, (_r) =>
        {
            if (isNullGo)
                return;

            selectRank = rankClass;
           
            //for (int i = 0; i < GetView().root_btns.childCount; i++)
            //{
            //    GetView().root_btns.GetChild(i).Find("Select").gameObject.SetActive(selectRank == rankClasses[i]);
            //}

            DicCfgRank = CFG_RankingAward.GetRankDicByRankClass(rankClass);

            RankListInfo = _r.data;

            OnRefreshUIShow();

            GetView().txt_Desc.text = $"<color=#fffd44>{GameUtil.Instance.GetTimeToDHorHM(RankM.Instance.GetLimitRankTime((int)selectRank))}</color>后进行榜单奖励结算，奖励将会通过邮件附件发放。";
        });

    }

    private void OnRefreshUIShow()
    {
        for (int i = 0; i < 3; i++)
        {
            HTTP_RankUserInfo userInfo = RankListInfo.rank_all.Length > i ? RankListInfo.rank_all[i] : default;
            OnRankItemShow(userInfo, (i + 1), GetView().root_top3.Find("Rank_" + (i + 1)).gameObject, RankListInfo.rank_all.Length <= i);
        }

        OnRankItemShow(RankListInfo.rank_own, RankListInfo.rank_own.num, GetView().root_MyRank.gameObject);

        GetView().InfiniteScroll_Rank.SetListItemCount(100, true);
        GetView().InfiniteScroll_Rank.RefreshAllShownItem();
        GetView().InfiniteScroll_Rank.MovePanelToItemIndex(0, 0);
    }

    private void OnRankItemShow(HTTP_RankUserInfo useInfo, int Rankidx, GameObject go, bool isNull = false)
    {
        JUIExtensions.txt(go.transform.Find("txt_Rank")).text = Rankidx <= 0 ? "未上榜" : Rankidx.ToString();

        Transform rootAward = go.transform.Find("root_award");

        if (rootAward && !rootAward.Find("btn_info"))
        {
            rootAward.gameObject.SetActive(DicCfgRank.ContainsKey(Rankidx) && Rankidx > 0);
            if (DicCfgRank.ContainsKey(Rankidx) && Rankidx > 0)
            {
                for (int i = 0; i < 3; i++)
                {
                    var root_award = rootAward.Find("award_" + i);
                    root_award.gameObject.SetActive(DicCfgRank[Rankidx].award.Length > i);
                    if (DicCfgRank[Rankidx].award.Length > i)
                    {
                        var itemID = int.Parse(DicCfgRank[Rankidx].award[i].ElementAt(0).Key);
                        GameUtil.Instance.SetItemIcon(itemID, JUIExtensions.img(root_award.Find("icon")), JUIExtensions.img(root_award));
                        JUIExtensions.txt(root_award.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(DicCfgRank[Rankidx].award[i].ElementAt(0).Value);
                        var awardBtn = JUIExtensions.btn(root_award);
                        awardBtn.onClick.RemoveAllListeners();
                        awardBtn.onClick.AddListener(delegate
                        {
                            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = itemID });
                        });
                    }
                }
            }
           
        }
        
        go.transform.Find("HeadBG").gameObject.SetActive(!isNull);
        go.transform.Find("img_des").gameObject.SetActive(!isNull);

        if (go.transform.Find("emptySeat"))
            go.transform.Find("emptySeat").gameObject.SetActive(isNull);


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
            JUIExtensions.txt(go.transform.Find("img_des/txt_des")).text = useInfo.value == 0 ? "暂无数据" : GetStrValueDes(useInfo.value);

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

        if (rootAward && rootAward.Find("btn_info"))
        {
            JUIExtensions.btn(rootAward.Find("btn_info")).onClick.RemoveAllListeners();
            JUIExtensions.btn(rootAward.Find("btn_info")).onClick.AddListener(delegate
            {
                OpenAwardUI(Rankidx);
            });
        }
    }

    private string GetStrValueDes(long value)
    {
        int layer =(int)( value >> 32);

        if (CFG_RankingAward.RankClass.Fight == selectRank)
        {
            var CFGBattle = CFG_Battle.GetById(layer);
            return "战役" + CFGBattle.Chapter + "-" + CFGBattle.Layers;
        }
        return CFG_MainBattle.GetById(layer).name;
    }

    private LoopListViewItem2 OnFillRank(LoopListView2 loopListView, int rankIdx) 
    {
        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        HTTP_RankUserInfo userInfo = RankListInfo.rank_all.Length > rankIdx ? RankListInfo.rank_all[rankIdx] : default;
        OnRankItemShow(userInfo, rankIdx + 1, go, RankListInfo.rank_all.Length <= rankIdx);
        return item;
    }

    private TimeItem RefreshUILoop;

    public override void OnEnter(params object[] args)
    {
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(Refresh_ZeroTime);
        //OnBtnsOpne();
        selectRank = (CFG_RankingAward.RankClass)args[0];
        OnRankShow(selectRank);

        RefreshUILoop = GameHelper.Instance.AddInvoke(new TimeItem(0, 60, () =>
        {
            if (isNullGo)
                return;
            OnRankShow(selectRank);
        }, _isLoop: true, isDelay: true));
    }

    private void OnBtnsOpne()
    {
        selectRank = CFG_RankingAward.RankClass.No;
        for (int i = 0; i < GetView().root_btns.childCount; i++)
        {
            var isOpen = RankM.Instance.GetLimitRankTime((int)rankClasses[i]) > 0;
            GetView().root_btns.GetChild(i).gameObject.SetActive(isOpen);
            if (isOpen && selectRank == CFG_RankingAward.RankClass.No)
                selectRank = rankClasses[i];
        }
        OnRankShow(selectRank);
    }

    private void Refresh_ZeroTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        if(RankM.Instance.GetLimitRankTime((int)selectRank) <= 0)
        {
            UITool.Close<UILimitetRank>();
            return;
        }

        //OnBtnsOpne();
        OnRankShow(selectRank);
    }

    public override void OnExit()
    {
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(Refresh_ZeroTime);
        if (GameHelper.Instance != null)
            GameHelper.Instance.RemoveInvoke(RefreshUILoop);
        RefreshUILoop = null;
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}