using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Event;
using JEngine.Helper;
using System;
using JEngine.MyNet;
using pb;
using SuperScrollView;
using JEngine.UI;
using HotUpdateScripts.ProtoM;

public class UIGuildApplication : UIBase<UIModelGuildApplication, UIViewGuildApplication>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_SortCombatPower.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isDescendingSortByCombatPower = !isDescendingSortByCombatPower;
            SortList();
            GetView().loopList_ApplicationList.SetListItemCount(applyList.Count);
            GetView().loopList_ApplicationList.RefreshAllShownItem();
        });

        GetView().btn_SortLevel.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isDescendingSortByCombatPower = !isDescendingSortByCombatPower;
            SortList();
            GetView().loopList_ApplicationList.SetListItemCount(applyList.Count);
            GetView().loopList_ApplicationList.RefreshAllShownItem();
        });

        GetView().loopList_ApplicationList.InitListView(0, OnItemRefresh);

        GetView().Btn_RefuseAll.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isAllAgreeOrRefuse = true;
            OnRefuseApply(GetAllMembersId());
            // 全部拒绝后清空界面列表
            ClearList();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "行会",
                component = "成员",
                position = "入会申请",
                event_click = "全部拒绝",
            });
        });

        GetView().Btn_AgreeAll.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isAllAgreeOrRefuse = true;
            OnAgreeApply(GetAllMembersId());
            // 全部同意后清空界面列表
            ClearList();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "行会",
                component = "成员",
                position = "入会申请",
                event_click = "全部同意",
            });
        });

        GetView().btn_RootPopup.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_RootPopup.gameObject.SetActive(false);
        });

        GetView().btn_Observe.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_RootPopup.gameObject.SetActive(false);
            OnObserveMeneber();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = "行会",
                position = "成员",
                event_click = "点击角色名观察角色信息",
            });
        });

        GetView().Btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGuildApplication>();
        });

    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        OnRequestApplicationList();
    }

    private LoopListViewItem2 OnItemRefresh(LoopListView2 listView, int index)
    {
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        Transform t = item.transform;
        t.Find("bg_1").gameObject.SetActive(index % 2 == 0);
        t.Find("bg_2").gameObject.SetActive(index % 2 == 1);
        int idx = index;

        JUIExtensions.btn(t.Find("btn_UserName")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("btn_UserName")).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnMemberNameButtonClicked(idx);
        });
        t.Find("btn_UserName/txt").GetComponent<Text>().text = applyList[index].ApplicantName;
        t.Find("btn_CombatPower").GetComponent<Text>().text = applyList[index].ApplicantTotalCombat.ToString();
        t.Find("txt_Level").GetComponent<Text>().text = applyList[index].ApplicantLevel.ToString();

        JUIExtensions.btn(t.Find("opt/btn_refuse")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("opt/btn_refuse")).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isAllAgreeOrRefuse = false;
            OnRefuseApply(new long[] { applyList[idx].ApplicantId });

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "行会",
                component = "成员",
                position = "入会申请",
                event_click = "单个拒绝",
            });
        });
        
        JUIExtensions.btn(t.Find("opt/btn_agree")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("opt/btn_agree")).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isAllAgreeOrRefuse = false;
            OnAgreeApply(new long[] { applyList[idx].ApplicantId });
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "行会",
                component = "成员",
                position = "入会申请",
                event_click = "单个同意",
            });
        });

        return item;
    }

    bool isAllAgreeOrRefuse;

    private void OnRequestApplicationList()
    {
        TcpC.C2S_GuildApplyListRequest();
    }

    List<GuildApplyInfo> applyList = new List<GuildApplyInfo>();
    // 按战力值排序 true 降序 false 升序
    bool isDescendingSortByCombatPower = true;
    // 按等级排序 true 降序 false 升序
    bool isDescendingSortByLevel = true;

    GuildApplyInfo currentSelectedApplyInfo;

    [Subscriber(ThreadMode.Main)]
    private void Event_GuildApplyList(GameEvent.Sync_GuildApplyList _e)
    {
        applyList = _e.applyList;
        RefreshUI();
    }

    private void ClearList()
    {
        applyList.Clear();
        GuildM.Instance.applyList.Clear();
        GetView().loopList_ApplicationList.SetListItemCount(0);
        GetView().loopList_ApplicationList.RefreshAllShownItem();
        GuildM.Instance.CheckGuildApplyRed();
    }

    private long[] GetAllMembersId()
    {
        long[] ids = new long[applyList.Count];
        for (int i = 0; i < applyList.Count; i++)
        {
            ids[i] = applyList[i].ApplicantId;
        }
        return ids;
    }

    private void RefreshUI()
    {
        GetView().loopList_ApplicationList.gameObject.SetActive(applyList.Count > 0);
        GetView().noApply1.SetActive(applyList.Count == 0);
        GetView().noApply2.SetActive(applyList.Count == 0);
        if (applyList.Count > 0)
        {
            SortList();
            GetView().loopList_ApplicationList.SetListItemCount(applyList.Count);
            GetView().loopList_ApplicationList.RefreshAllShownItem();
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_GuildHandleApply(GameEvent.Sync_GuildHandleApply _e)
    {
        if (!isAllAgreeOrRefuse)
            OnRequestApplicationList();
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_ObservePlayerInfo(GameEvent.Sync_ObservePlayerInfo _e)
    {
        UITool.Open<UIOtherRole>(null, _e.observePlayerInfo);
    }

    private void SortList()
    {
        applyList.Sort((a, b) =>
        {
            if (a.ApplicantTotalCombat != b.ApplicantTotalCombat)
            {
                if (isDescendingSortByLevel)
                {
                    return b.ApplicantTotalCombat.CompareTo(a.ApplicantTotalCombat);
                }
                else
                {
                    return a.ApplicantTotalCombat.CompareTo(b.ApplicantTotalCombat);
                }
            }
            else if (a.ApplicantLevel != b.ApplicantLevel)
            {
                if (isDescendingSortByLevel)
                {
                    return b.ApplicantLevel.CompareTo(a.ApplicantLevel);
                }
                else
                {
                    return a.ApplicantLevel.CompareTo(b.ApplicantLevel);
                }
            }
            else
            {
                return 0;
            }
        });
    }

    private void OnMemberNameButtonClicked(int index)
    {
        currentSelectedApplyInfo = applyList[index];
        RectTransformUtility.ScreenPointToLocalPointInRectangle(GetView().rect_PopupButton, Input.mousePosition, null, out Vector2 localPoint);
        GetView().rect_ButtonRoot.anchoredPosition = localPoint;
        GetView().btn_RootPopup.gameObject.SetActive(true);
    }

    private void OnRefuseApply(long[] ids)
    {
        TcpC.C2S_GuildHandleApplyRequest(ids, 2);
    }

    private void OnAgreeApply(long[] ids)
    {
        TcpC.C2S_GuildHandleApplyRequest(ids, 1);
    }

    private void OnObserveMeneber()
    {
        TcpC.C2S_ObservePlayerInfoRequest(currentSelectedApplyInfo.ApplicantId);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        currentSelectedApplyInfo = null;
        isDescendingSortByCombatPower = true;
        isDescendingSortByLevel = true;
        isAllAgreeOrRefuse = false;
        GetView().btn_RootPopup.gameObject.SetActive(false);
    }

    public override void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(this);
    }
}