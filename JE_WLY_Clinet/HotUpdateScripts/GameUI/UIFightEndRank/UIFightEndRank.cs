using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using SuperScrollView;

using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;

using UnityEngine;
using UnityEngine.UI;


public class UIFightEndRank : UIBase<UIModelFightEndRank, UIViewFightEndRank>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFightEndRank>();
        });

        GetView().looplist_Rank.InitListView(0, OnItemFill);
    }

    public override void OnEnter(params object[] args)
    {
        OnRankShow();
    }

    private Dictionary<int, CFGItem_RankingAward> DicCfgRank;
    private HTTP_RankDetail RankListInfo;

    private LoopListViewItem2 OnItemFill(LoopListView2 list, int index)
    {
        var item = list.NewListViewItem("Item");
        HTTP_RankUserInfo userInfo = RankListInfo.rank_all.Length > index ? RankListInfo.rank_all[index] : default;
        OnRankItemShow(userInfo, index + 1, item.gameObject, index >= RankListInfo.rank_all.Length);
        return item;
    }

    private void OnRankShow()
    {
        HttpC.GetRankList(CFG_RankingAward.RankClass.FightEnd, (_r) =>
        {
            if (isNullGo)
                return;

            RankListInfo = _r.data;
            DicCfgRank = CFG_RankingAward.GetRankDicByRankClass(CFG_RankingAward.RankClass.FightEnd);

            OnRankItemShow(RankListInfo.rank_own, RankListInfo.rank_own.num, GetView().root_MyRank.gameObject);

            GetView().looplist_Rank.SetListItemCount(100);
            GetView().looplist_Rank.RefreshAllShownItem();
            GetView().looplist_Rank.MovePanelToItemIndex(0, 0);
        });
    }

    private void OnRankItemShow(HTTP_RankUserInfo useInfo, int Rankidx, GameObject go, bool isNull = false)
    {
        Transform bg = go.transform.Find("bg");
        bg.gameObject.SetActive(Rankidx % 2 == 1);
        Transform award = go.transform.Find("Awards");
        for (int i = 0; i < 3; i++)
        {
            bool isShow = !isNull && DicCfgRank.ContainsKey(Rankidx) && DicCfgRank[Rankidx].award.Length > i;
            award.Find("Item_" + i).gameObject.SetActive(isShow);
            if (isShow)
            {
                int cfgId = int.Parse(DicCfgRank[Rankidx].award[i].ElementAt(0).Key);
                GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(award.Find("Item_" + i + "/icon")), JUIExtensions.img(award.Find("Item_" + i)));
                JUIExtensions.txt(award.Find("Item_" + i + "/txt_num")).text = GameUtil.Instance.FormatNumWithUnit(DicCfgRank[Rankidx].award[i].ElementAt(0).Value);
                JUIExtensions.btn(award.Find("Item_" + i)).onClick.RemoveAllListeners();
                JUIExtensions.btn(award.Find("Item_" + i)).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                });
            }
        }

        bool isTop3 = Rankidx >= 1 && Rankidx <= 3;
        JUIExtensions.img(go.transform, "img_Rank").enabled = isTop3;
        JUIExtensions.txt(go.transform.Find("txt_Rank")).enabled = !isTop3;
        if (Rankidx <= 0)
        {
            JUIExtensions.txt(go.transform.Find("txt_Rank")).text = "未上榜";
        }
        else if (isTop3)
        {
            ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/WorldBoss", "icon_" + Rankidx, JUIExtensions.img(go.transform, "img_Rank"));
        }
        else
        {
            JUIExtensions.txt(go.transform.Find("txt_Rank")).text = Rankidx.ToString();
        }

        if (isNull)
        {
            JUIExtensions.txt(go.transform, "txt_Name").text = "虚位以待";
            JUIExtensions.txt(go.transform, "txt_Level").text = "";
            JUIExtensions.txt(go.transform, "txt_Time").text = "";
        }
        else
        {
            string serverArea = "";

            if(Rankidx > 0)
            {
                string areaId = Regex.Match(useInfo.server_name ?? "", "\\d+").ToString();
                if (string.IsNullOrEmpty(areaId))
                    serverArea = $"[{useInfo.server_name}]";
                else
                    serverArea = $"[{areaId}服]";
            }

            JUIExtensions.txt(go.transform, "txt_Name").text = $"{serverArea}{(Rankidx <= 0 ? RoleM.Instance.Nickname : useInfo.name)}";

            long level = useInfo.value / 1000000;
            long seconds = 1000000 - (useInfo.value % 1000000);
            JUIExtensions.txt(go.transform, "txt_Level").text = Rankidx <= 0 ? "" : level.ToString();
            JUIExtensions.txt(go.transform, "txt_Time").text = Rankidx <= 0 ? "" : GameUtil.Instance.GetTimeToHMS(seconds);
        }
    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {

    }
}