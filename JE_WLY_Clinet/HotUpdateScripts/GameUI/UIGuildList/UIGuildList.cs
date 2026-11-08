using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using SuperScrollView;
using System;
using pb;
using System.Linq;
using JEngine.UI;

public class UIGuildList : UIBase<UIModelGuildList, UIViewGuildList>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public override void OnInit()
    {
        GetView().input_SearchGuild.onValueChanged.AddListener((v) =>
        {
            if (string.IsNullOrEmpty(v))
            {
                currentShowGuilds = guildInfos;
                RefreshList();
            }
        });
        GetView().btn_SearchGuild.onClick.AddListener(SearchGuild);

        GetView().btn_RefreshGuildList.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isScrollListToTop = true;
            OnRequestGuildList();
        });

        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 32);
        });

        GetView().btn_SortGuildList.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            isDescendingOrder = !isDescendingOrder;
            SortGuildList();
            RefreshList();
        });

        GetView().loopList_Guild.InitListView(0, OnItemRefresh);

        GetView().btn_CreateGuild.onClick.AddListener(CreateGuild);
        GetView().btn_ApplyAllGuild.onClick.AddListener(ApplyAllGuild);
        GetView().btn_QuitGuild.onClick.AddListener(QuitGuild);

        GetView().btn_RootPopup.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_RootPopup.gameObject.SetActive(false);
        });

        GetView().btn_Observe.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_RootPopup.gameObject.SetActive(false);
            OnObserve();
        });
    }

    bool isShow;

    public void OnShow(bool isShow)
    {
        if(isShow == this.isShow)
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

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(Event_UserBase);
        isJoinGuild = RoleM.Instance.IsJoinedGuild();
        isScrollListToTop = true;
        ShowButtons();
        OnRequestGuildList();
    }

    List<GuildInfo> guildInfos = new List<GuildInfo>();
    // 当前显示的行会（搜索到的，或者全部）
    List<GuildInfo> currentShowGuilds = new List<GuildInfo>();

    bool isDescendingOrder = true;
    bool isScrollListToTop;
    bool isJoinGuild;

    GuildInfo currentSelectedGuildInfo;

    private LoopListViewItem2 OnItemRefresh(LoopListView2 listView, int index)
    {
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        Transform t = item.transform;
        bool isMy = currentShowGuilds[index].Id == RoleM.Instance.GuildId;
        t.Find("bg_1").gameObject.SetActive(!isMy && index % 2 == 0);
        t.Find("bg_2").gameObject.SetActive(!isMy && index % 2 == 1);
        t.Find("bg_3").gameObject.SetActive(isMy);
        t.Find("txt_guildName").GetComponent<Text>().text = currentShowGuilds[index].Name;
        int idx = index;

        JUIExtensions.btn(t.Find("btn_UserName")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("btn_UserName")).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            currentSelectedGuildInfo = currentShowGuilds[index];
            if (currentSelectedGuildInfo.LeaderId != RoleM.Instance.Uid)
                SetPopupButtonPosition();
        });

        t.Find("btn_UserName/txt").GetComponent<Text>().text = currentShowGuilds[index].LeaderName;
        t.Find("txt_peopleNum").GetComponent<Text>().text = $"{currentShowGuilds[index].MemberCount}/{currentShowGuilds[index].MaxMembers}";
        t.Find("img_myguild").gameObject.SetActive(isMy && currentShowGuilds[index].Application == 0);
        t.Find("btn_Apply").gameObject.SetActive(currentShowGuilds[index].Application == 2);
        t.Find("btn_CancelApply").gameObject.SetActive(currentShowGuilds[index].Application == 1);

        JUIExtensions.btn(t.Find("btn_Apply")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("btn_Apply")).onClick.AddListener(() =>
        {
            TcpC.C2S_GuildApplyRequest(new long[1] { currentShowGuilds[idx].Id }); // 需要刷新一下列表，刷新按钮状态。

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "社交",
                component = "行会",
                position = "行会列表",
                event_click = "单个申请",
            });
        });

        JUIExtensions.btn(t.Find("btn_CancelApply")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("btn_CancelApply")).onClick.AddListener(() =>
        {
            TcpC.C2S_GuildCancelApplyRequest(currentShowGuilds[idx].Id);
        });

        return item;
    }

    public void SetPopupButtonPosition()
    {
        RectTransformUtility.ScreenPointToLocalPointInRectangle(GetView().rect_PopupButton, Input.mousePosition, JUIExtensions.uiCamera, out Vector2 localPoint);
        GetView().rect_ButtonRoot.anchoredPosition = localPoint;
        GetView().btn_RootPopup.gameObject.SetActive(true);
    }

    private void OnRequestGuildList()
    {
        TcpC.C2S_GuildListRequest(1, GuildM.PageSize);
    }

    private void OnObserve()
    {
        TcpC.C2S_ObservePlayerInfoRequest(currentSelectedGuildInfo.LeaderId);

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "社交",
            component = "行会",
            position = "行会列表",
            event_click = "点击会长名字，查看会长信息",
        });
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_RefreshList(GameEvent.Sync_GuildList _e)
    {
        if (isNullGo)
            return;

        guildInfos = _e.guildsInfo;
        currentShowGuilds = guildInfos;

        SortGuildList();
        RefreshList();
    }

    [Subscriber(ThreadMode.Main)]
    // 离开工会
    private void Event_GuildLeave(GameEvent.Sync_GuildLeave _e)
    {
        if (isNullGo)
            return;

        isScrollListToTop = false;
        ShowButtons();
        OnRequestGuildList();
    }

    [Subscriber(ThreadMode.Main)]
    // 工会创建成功
    private void Event_GuildCreate(GameEvent.Sync_GuildCreate _e)
    {
        if (isNullGo)
            return;

        isScrollListToTop = false;
        ShowButtons();
        OnRequestGuildList();
    }
    
    [Subscriber(ThreadMode.Main)]
    // 申请加入工会。申请成功回复。
    private void Event_GuildApply(GameEvent.Sync_GuildApply _e)
    {
        if (isNullGo)
            return;

        isScrollListToTop = false;
        ShowButtons();
        OnRequestGuildList();
    }
        
    [Subscriber(ThreadMode.Main)]
    // 取消申请加入工会。取消成功回复。
    private void Event_GuildCancelApply(GameEvent.Sync_GuildCancelApply _e)
    {
        if (isNullGo)
            return;

        isScrollListToTop = false;
        ShowButtons();
        OnRequestGuildList();
    }


    //[Subscriber(ThreadMode.Main)]
    // 对方是否同意加入加入工会
    private void Event_UserBase(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        if (isJoinGuild != RoleM.Instance.IsJoinedGuild())
        {
            isScrollListToTop = false;
            ShowButtons();
            OnRequestGuildList();
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_ObservePlayerInfo(GameEvent.Sync_ObservePlayerInfo _e)
    {
        if (isNullGo)
            return;

        UITool.Open<UIOtherRole>(null, _e.observePlayerInfo);
    }

    private void RefreshList()
    {
        GetView().loopList_Guild.SetListItemCount(currentShowGuilds.Count, isScrollListToTop);
        GetView().loopList_Guild.RefreshAllShownItem();
    }

    private void ShowButtons()
    {
        GetView().btn_CreateGuild.gameObject.SetActive(!RoleM.Instance.IsJoinedGuild());
        GetView().btn_ApplyAllGuild.gameObject.SetActive(!RoleM.Instance.IsJoinedGuild());
        GetView().btn_QuitGuild.gameObject.SetActive(RoleM.Instance.IsJoinedGuild());
    }

    private void SearchGuild()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        string searchName = GetView().input_SearchGuild.text.Trim().ToLower();
        if (string.IsNullOrEmpty(searchName))
        {
            UITool.ShowTips("请输入要查找的行会名称");
            return;
        }

        currentShowGuilds = guildInfos.FindAll(v => v.Name.ToLower().Contains(searchName));

        RefreshList();
    }

    private void SortGuildList()
    {
        if (isDescendingOrder)
        {
            currentShowGuilds.Sort((x, y) => y.MemberCount.CompareTo(x.MemberCount));
        }
        else
        {
            currentShowGuilds.Sort((x, y) => x.MemberCount.CompareTo(y.MemberCount));
        }
    }

    private void CreateGuild()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UICreateGuild>();
    }

    bool applyButtonIsCooling;
    private void ApplyAllGuild()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (applyButtonIsCooling)
        {
            UITool.ShowTips("操作过于频繁");
            return;
        }

        applyButtonIsCooling = true;

        if (guildInfos.Count == 0) return;

        List<long> ids = new List<long>();
        for (int i = 0; i < guildInfos.Count; i++)
        {
            if (guildInfos[i].Application == 2)
                ids.Add(guildInfos[i].Id);
        }

        GameHelper.Instance.AddInvoke(new TimeItem(5, () =>
        {
            applyButtonIsCooling = false;
        }));

        if (ids.Count == 0)
            return;

        TcpC.C2S_GuildApplyRequest(ids.ToArray());

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "社交",
            component = "行会",
            position = "行会列表",
            event_click = "一键申请",
        });
    }

    private void QuitGuild()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
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
                    position = "行会列表",
                    event_click = "退出行会",
                });
            },
        });
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(Event_UserBase);

        isShow = false;
        currentSelectedGuildInfo = null;
        isDescendingOrder = true;
        GetView().btn_RootPopup.gameObject.SetActive(false);
    }

    public override void OnDestroy()
    {

    }
}