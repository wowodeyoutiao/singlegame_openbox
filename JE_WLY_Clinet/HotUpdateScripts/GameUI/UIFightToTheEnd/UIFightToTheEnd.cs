using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Tsp;

using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using SuperScrollView;
using System;
using UnityEngine;


public class UIFightToTheEnd : UIBase<UIModelFightToTheEnd, UIViewFightToTheEnd>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public override void OnInit()
    {
        //GetView().btn_Help.onClick.AddListener(() =>
        //{
        //    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        //    UITool.Open<UIComRule>(null, 1);
        //});

        GetView().btn_Rank.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIFightEndRank>();
        });

        GetView().btn_Fight.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd))
                return;

            if (!GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.FightEnd.startTime, CFG_XMLConfig.xmlConfig.FightEnd.endTime))
            {
                UITool.ShowTips($"每天 {CFG_XMLConfig.xmlConfig.FightEnd.startTime}~{CFG_XMLConfig.xmlConfig.FightEnd.endTime} 开启");
                return;
            }

            if (!GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd))
            {
                GameHelper.Instance.StarGame(GetSelectedBattleId());
            }
        });

        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFightToTheEnd>();
        });

        //GetView().itemReference.SetActive(false);

        GetView().looplist_Map.InitListView(0, OnItemRefresh);
    }

    private LoopListViewItem2 OnItemRefresh(LoopListView2 list, int index)
    {
        var item = list.NewListViewItem("Item");
        int level = index + 1;

        GameObject go = item.gameObject;
        go.transform.SetAsFirstSibling();
        go.name = "Item_" + level;
        go.SetActive(true);
        float posX = index % 2 == 0 ? -102f : 118f;

        Transform t = go.transform.Find("Item");
        t.localPosition = new Vector3(posX, t.localPosition.y, 0f);
        go.transform.Find("Item/mask").gameObject.SetActive(level > lastPassLevel + 1);
        JUIExtensions.img(go.transform.Find("Item/icons/icon_passed")).enabled = level <= lastPassLevel;
        JUIExtensions.img(go.transform.Find("Item/icons/icon_nopass")).enabled = level == lastPassLevel + 1;
        JUIExtensions.img(go.transform.Find("Item/icons/icon_lock")).enabled = level > lastPassLevel + 1;

        JUIExtensions.txt(go.transform.Find("Item/img_level/txt")).text = $"第{level}关";

        string useTimeString = level > lastPassLevel ? "通关耗时：无" : $"通关耗时：{GameUtil.Instance.GetTimeToHMS(toEnd.Evaluation[index].K)}";
        JUIExtensions.txt(go.transform.Find("Item/timeinfo/txt_usetime")).text = useTimeString;

        long seconds = 0;
        if (index < allRank.Length)
            seconds = 1000000 - (allRank[index].value % 1000000);

        JUIExtensions.txt(go.transform.Find("Item/timeinfo/txt_record")).text = index < allRank.Length ? $"全区记录：{GameUtil.Instance.GetTimeToHMS(seconds)}" : $"全区记录：无";

        bool right = index % 2 == 0;
        go.transform.Find("Item/decorate_left").gameObject.SetActive(!right);
        go.transform.Find("Item/decorate_right").gameObject.SetActive(right);

        JUIExtensions.img(go.transform, "Item/img_unselect").enabled = index + 1 != selectLevel;
        JUIExtensions.img(go.transform, "Item/img_selected").enabled = index + 1 == selectLevel;

        if (right)
        {
            JUIExtensions.img(go.transform.Find("Item/decorate_right/img_lock")).enabled = level > lastPassLevel + 1;
            JUIExtensions.img(go.transform.Find("Item/decorate_right/img_unlock")).enabled = level <= lastPassLevel + 1;
        }
        else
        {
            JUIExtensions.img(go.transform.Find("Item/decorate_left/img_lock")).enabled = level > lastPassLevel + 1;
            JUIExtensions.img(go.transform.Find("Item/decorate_left/img_unlock")).enabled = level <= lastPassLevel + 1;
        }

        JUIExtensions.btn(go.transform.Find("Item")).onClick.RemoveAllListeners();
        JUIExtensions.btn(go.transform.Find("Item")).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (level <= lastPassLevel + 1)
            {
                selectLevel = level;
                SelectedLevelItem();
            }
        });
        return item;
    }

    int levelCount;
    int lastPassLevel;
    int selectLevel = 1;

    TimeItem openStateTimeItem;
    private void CheckOpenState()
    {
        ClearTime();
        int seconds = 0;
        if (GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.FightEnd.startTime, CFG_XMLConfig.xmlConfig.FightEnd.endTime))
            seconds = GameUtil.Instance.CheckNextEndTime(CFG_XMLConfig.xmlConfig.FightEnd.endTime);
        else
            seconds = GameUtil.Instance.CheckNextOpenTime(CFG_XMLConfig.xmlConfig.FightEnd.startTime, CFG_XMLConfig.xmlConfig.FightEnd.endTime);

        openStateTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(seconds, () =>
        {
            if (isNullGo) return;
            SetOpenState();
            CheckOpenState();
        }));
    }

    private void ClearTime()
    {
        if (GameHelper.Instance != null && openStateTimeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(openStateTimeItem);
            openStateTimeItem = null;
        }
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_ZeroRefresh);

        RefreshUI();
    }

    InstanceInfoResponse.ToEnd toEnd;

    private void RefreshUI()
    {
        toEnd = FightPlayM.Instance.GetInfo().to_end;
        levelCount = CFG_Battle.YZDDList.Count;
        if (toEnd.BattleId > 0)
            lastPassLevel = CFG_Battle.GetById(toEnd.BattleId).Layers;
        else
            lastPassLevel = 0;

        selectLevel = lastPassLevel == levelCount ? lastPassLevel : lastPassLevel + 1;

        SetOpenState();
        CheckOpenState();
        OnRankShow();
    }


    private HTTP_RankUserInfo[] allRank;

    private void OnRankShow()
    {
        HttpC.GetRankList(CFG_RankingAward.RankClass.ServerRecord, (_r) =>
        {
            if (isNullGo)
                return;

            allRank = _r.data.rank_all;
            Array.Sort(allRank, (a, b) => a.uid.CompareTo(b.uid));

            RefreshList();

        });
    }

    void SetOpenState()
    {
        bool isOpen = GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.FightEnd.startTime, CFG_XMLConfig.xmlConfig.FightEnd.endTime);
        //test code
        //bool isOpen = true;
        GetView().btn_Fight.gameObject.SetActive(isOpen);
        GetView().txt_Tip.gameObject.SetActive(!isOpen);
        if (!isOpen)
            GetView().txt_Tip.text = $"每天{CFG_XMLConfig.xmlConfig.FightEnd.startTime}~{CFG_XMLConfig.xmlConfig.FightEnd.endTime}可挑战";
    }

    void RefreshList()
    {
        GetView().looplist_Map.SetListItemCount(levelCount);
        GetView().looplist_Map.RefreshAllShownItem();
        GetView().looplist_Map.MovePanelToItemIndex(selectLevel - 1, 240);
    }

    int GetSelectedBattleId()
    {
        return CFG_Battle.YZDDList.Find(v => v.Layers == selectLevel).UID;
    }

    void SelectedLevelItem()
    {
        GetView().looplist_Map.RefreshAllShownItem();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_ZeroRefresh(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        TcpC.C2S_InstanceInfo();
    }


    [Subscriber(ThreadMode.Main)]
    private void Refresh_InstanceInfo(GameEvent.Refresh_InstanceInfo _e)
    {
        if (isNullGo) return;
        RefreshUI();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_ZeroRefresh);

    }

    public override void OnDestroy()
    {

    }
}