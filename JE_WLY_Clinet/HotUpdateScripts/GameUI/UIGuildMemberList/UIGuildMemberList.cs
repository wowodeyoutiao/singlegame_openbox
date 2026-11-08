using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using SuperScrollView;
using System;
using JEngine.UI;
using System.Reflection;
using JEngine.Event;
using pb;
using JEngine.Helper;

public class UIGuildMemberList : UIBase<UIModelGuildMemberList, UIViewGuildMemberList>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public override void OnInit()
    {
        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 37);
        });

        GetView().btn_SortCombatPower.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isDescendingSortByCombatPower = !isDescendingSortByCombatPower;
            isSortByLevel = false;
            SortList();
            GetView().loopList_MemberList.SetListItemCount(memberInfos.Count);
            GetView().loopList_MemberList.RefreshAllShownItem();
        });

        GetView().btn_SortLevel.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isDescendingSortByLevel = !isDescendingSortByLevel;
            isSortByLevel = true;
            SortList(true);
            GetView().loopList_MemberList.SetListItemCount(memberInfos.Count);
            GetView().loopList_MemberList.RefreshAllShownItem();
        });

        GetView().loopList_MemberList.InitListView(0, OnItemRefresh);

        GetView().btn_QuitGuild.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            QuitGuild();
        });

        GetView().btn_GuildApplication.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGuildApplication>();
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
        });

        GetView().btn_AppointVicePresident.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_RootPopup.gameObject.SetActive(false);
            OnAppointVicePresident();
        });

        GetView().btn_DismissVicePresident.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_RootPopup.gameObject.SetActive(false);
            OnDismissVicePresident();
        });

        GetView().btn_PleaseGoOut.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_RootPopup.gameObject.SetActive(false);
            OnPleaseGoOut();
        });

    }

    bool isShow;

    public void OnShow(bool isShow)
    {
        if (isShow == this.isShow)
            return;

        this.isShow = isShow;

        if (isShow)
        {
            OnEnter();
        }
        else
        {
            OnExit();
        }
        uiGo.SetActive(isShow);
    }

    bool isScrollListToTop;


    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        isScrollListToTop = true;
        isDescendingSortByCombatPower = true;
        isDescendingSortByLevel = true;
        isSortByLevel = false;

        OnBindRed();
        TcpC.C2S_GuildInfoRequest(RoleM.Instance.GuildId);
        OnRequestGuildList();
    }

    private void OnBindRed()
    {
        RedPointManager.Instance.RPS.Register(RedPointConst.Guild_Apply.ToString(), (node) =>
        {
            GetView().btn_GuildApplication.transform.Find("dot").gameObject.SetActive(node.pointCount > 0);
        });
    }

    private void ClearBindRed()
    {
        RedPointManager.Instance.RPS.Register(RedPointConst.Guild_Apply.ToString(), null);
    }


    private void OnRequestGuildList()
    {
        TcpC.C2S_GuildMemberListRequest(RoleM.Instance.GuildId);
    }

    List<GuildMemberInfo> memberInfos = new List<GuildMemberInfo>();

    [Subscriber(ThreadMode.Main)]
    private void Event_RefreshList(GameEvent.Sync_GuildMemberList _e)
    {
        if (isNullGo)
            return;

        memberInfos = _e.membersInfo;
        SortList(isSortByLevel);
        GetView().loopList_MemberList.SetListItemCount(memberInfos.Count, isScrollListToTop);
        GetView().loopList_MemberList.RefreshAllShownItem();
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_GuildPromote(GameEvent.Sync_GuildPromote _e)
    {
        if (isNullGo)
            return;

        isScrollListToTop = false;
        TcpC.C2S_GuildInfoRequest(RoleM.Instance.GuildId);
        OnRequestGuildList();
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_GuildKick(GameEvent.Sync_GuildKick _e)
    {
        if (isNullGo)
            return;

        isScrollListToTop = false;
        OnRequestGuildList();
    }

    GuildInfo guildInfo;
    // 我的职位：1 会长 2 副会长 3 成员
    int myPosition;

    [Subscriber(ThreadMode.Main)]
    // 工会信息
    private void Event_GuildInfo(GameEvent.Sync_GuildInfo _e)
    {
        if (isNullGo)
            return;

        guildInfo = _e.guildInfo;
        myPosition = _e.Role;

        // 会长需要显示申请按钮
        ShowGuildApplyButton(IsPresident() || IsVicePresident());
    }


    [Subscriber(ThreadMode.Main)]
    private void Event_GuildHandleApply(GameEvent.Sync_GuildHandleApply _e)
    {
        if (isNullGo)
            return;

        OnRequestGuildList();
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_ObservePlayerInfo(GameEvent.Sync_ObservePlayerInfo _e)
    {
        if (isNullGo)
            return;

        UITool.Open<UIOtherRole>(null, _e.observePlayerInfo);
    }

    private void ShowGuildApplyButton(bool isShow)
    {
        GetView().btn_GuildApplication.gameObject.SetActive(isShow);
    }

    // 按战力值排序 true 降序 false 升序
    bool isDescendingSortByCombatPower = true;
    // 按等级排序 true 降序 false 升序
    bool isDescendingSortByLevel = true;

    bool isSortByLevel;

    private void SortList(bool isByLevel = false)
    {
        memberInfos.Sort((a, b) =>
        {
            if (a.MemberLevel != b.MemberLevel)
            {
                if (isByLevel)
                {
                    if (isDescendingSortByLevel)
                        return b.MemberLevel.CompareTo(a.MemberLevel);
                    else
                        return a.MemberLevel.CompareTo(b.MemberLevel);
                }
                else
                {
                    if(isDescendingSortByCombatPower)
                        return b.MemberTotalCombat.CompareTo(a.MemberTotalCombat);
                    else
                        return a.MemberTotalCombat.CompareTo(b.MemberTotalCombat);
                }
            }
            else
            {
                return 0;
            }
        });
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

            currentSelectedMemberInfo = memberInfos[idx];
            if (currentSelectedMemberInfo.MemberId != RoleM.Instance.Uid)
                OnMemberNameButtonClicked(idx);
        });

        t.Find("img_President").gameObject.SetActive(memberInfos[index].Position == 1);
        t.Find("img_VicePresident").gameObject.SetActive(memberInfos[index].Position == 2);

        JUIExtensions.txt(t.Find("btn_UserName/txt")).text = memberInfos[index].MemberName;
        t.Find("btn_CombatPower").GetComponent<Text>().text = memberInfos[index].MemberTotalCombat.ToString();
        JUIExtensions.txt(t.Find("txt_Level")).text = memberInfos[index].MemberLevel.ToString();
        var time = memberInfos[index].MemberOnline;

        t.Find("state/img_online").gameObject.SetActive(time == 0);
        t.Find("state/img_offline").gameObject.SetActive(time != 0);

        if (time != 0)
            JUIExtensions.txt(t.Find("state/img_offline/txt")).text = GetOnlineStateString(time);

        return item;
    }

    private string GetOnlineStateString(long time)
    {
        long timeSeconds = GameUtil.Instance.GetCurrentTimeSeconds() - time;
        var timeSpan = GameUtil.Instance.GetTimeSpan(timeSeconds);
        if(timeSpan.Days >= 7)
            return $"超过7天";
        else if (timeSpan.Days >= 3)
            return $"3天前";
        else if (timeSpan.Days >= 1)
            return $"1天前";
        else
            return $"{timeSpan.Hours}小时之前";
    }

    private void QuitGuild()
    {
        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = $"退出行会之后需要{CFG_XMLConfig.xmlConfig.Guild.interlimitTime}小时之后才能再次加入行会，是否确定退出？",
            YesFun = () =>
            {
                TcpC.C2S_GuildLeaveRequest(RoleM.Instance.GuildId);

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "社交",
                    component = "行会",
                    position = "成员",
                    event_click = "退出行会",
                });
            },
        });
    }

    GuildMemberInfo currentSelectedMemberInfo;

    private void OnMemberNameButtonClicked(int index)
    {
        GetView().btn_AppointVicePresident.gameObject.SetActive(IsPresident() && currentSelectedMemberInfo.Position == 3);
        GetView().btn_DismissVicePresident.gameObject.SetActive(IsPresident() && currentSelectedMemberInfo.Position == 2);
        GetView().btn_PleaseGoOut.gameObject.SetActive(IsPresident() || (IsVicePresident() && currentSelectedMemberInfo.Position == 3)); // 副会长只能请出普通成员

        RectTransformUtility.ScreenPointToLocalPointInRectangle(GetView().rect_PopupButton, Input.mousePosition, JUIExtensions.uiCamera, out Vector2 localPoint);
        GetView().rect_ButtonRoot.anchoredPosition = localPoint;
        GetView().btn_RootPopup.gameObject.SetActive(true);
    }

    /// <summary>
    /// 是否是会长
    /// </summary>
    /// <returns></returns>
    public bool IsPresident()
    {
        return myPosition == 1;
    }

    /// <summary>
    /// 是否是副会长
    /// </summary>
    /// <returns></returns>
    public bool IsVicePresident()
    {
        return myPosition == 2;
    }

    /// <summary>
    /// 观察成员信息
    /// </summary>
    private void OnObserveMeneber()
    {
        TcpC.C2S_ObservePlayerInfoRequest(currentSelectedMemberInfo.MemberId);
    }

    private void OnAppointVicePresident()
    {
        TcpC.C2S_GuildPromoteRequest(RoleM.Instance.GuildId, currentSelectedMemberInfo.MemberId, 2);
    }

    private void OnDismissVicePresident()
    {
        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = "是否确定解除该玩家副会长职务？",
            YesFun = () => {
                TcpC.C2S_GuildPromoteRequest(RoleM.Instance.GuildId, currentSelectedMemberInfo.MemberId, 3);
            },
        });
    }

    private void OnPleaseGoOut()
    {
        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = "是否确定要将当前玩家踢出行会？",
            YesFun = () => {
                TcpC.C2S_GuildKickRequest(RoleM.Instance.GuildId, currentSelectedMemberInfo.MemberId);
            },
        });
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        ClearBindRed();
        isShow = false;
        currentSelectedMemberInfo = null;
        GetView().btn_RootPopup.gameObject.SetActive(false);
    }


    public override void OnDestroy()
    {

    }
}