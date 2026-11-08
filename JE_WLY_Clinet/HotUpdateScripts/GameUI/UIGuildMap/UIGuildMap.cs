using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Event;
using HotUpdateScripts.ProtoM;
using SuperScrollView;
using System;
using JEngine.UI;
using JEngine.Helper;

public class UIGuildMap : UIBase<UIModelGuildMap, UIViewGuildMap>
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
            UITool.Open<UIComRule>(null, 36);
        });

        GetView().btn_AddItem.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGuildDonation>();
        });

        GetView().btn_GoMap.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameHelper.Instance.StarGame(battleId);
        });

        GetView().loopList_Awards.InitListView(0, OnItemRefresh);
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

    readonly int battleId = 1800001;
    CFGItem_Room cfg_room;
    CFGItem_Battle cfg_battle;

    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_RefreshPrayItem);

        cfg_room = CFG_Room.GetByRoomType(CFG_Room.RoomType.Guild, -1);
        cfg_battle = CFG_Battle.GetById(battleId);
        ShowDropItems();
        ShowPrayItemCount(RoleM.Instance.GuildContribution);
        ShowIntroduce();
        CheckButtonOpen();
    }

    private LoopListViewItem2 OnItemRefresh(LoopListView2 listView, int index)
    {
        int idx = index;
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        Transform t = item.transform;
        GameUtil.Instance.SetItemIcon(cfg_battle.Preview[index], JUIExtensions.img(t.Find("icon")), JUIExtensions.img(t));
        JUIExtensions.btn(t).onClick.RemoveAllListeners();
        JUIExtensions.btn(t).onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfg_battle.Preview[idx] });
        });
        return item;
    }

    public void ShowPrayItemCount(long count)
    {
        GetView().txt_ItemCount.text = count.ToString();
    }

    public void ShowIntroduce()
    {
        GetView().txt_MapDesc.text = cfg_room.Suggest;
        GetView().txt_PlayingRule.text = cfg_room.CopyDes;
    }

    private void CheckButtonOpen()
    {
        //bool isOpen = GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.Guild.guildcopy.starttime, CFG_XMLConfig.xmlConfig.Guild.guildcopy.endtime);
        bool isOpen = GuildM.Instance.IsGuildCopyOpen();
        GetView().btn_GoMap.gameObject.SetActive(isOpen);
        GetView().btn_Gray.gameObject.SetActive(!isOpen);
        GetView().txt_OpenTime.gameObject.SetActive(!isOpen);
        if (!isOpen)
        {
            int seconds = GameUtil.Instance.CheckNextOpenTime(CFG_XMLConfig.xmlConfig.Guild.guildcopy.starttime, CFG_XMLConfig.xmlConfig.Guild.guildcopy.endtime);
            ShowNextOpenTime(seconds);
        }
        GuildM.Instance.CheckGuildMapRed();
    }

    TimeItem timeItem;

    private void ShowNextOpenTime(int seconds)
    {
        ClearTime();
        int timeSeconds = seconds;
        GetView().txt_OpenTime.text = $"（{GameUtil.Instance.GetTimeSToHMS(timeSeconds)}后地图开启）";

        timeItem = GameHelper.Instance.AddInvoke(new TimeItem(-1, 1, () =>
        {
            timeSeconds--;
            GetView().txt_OpenTime.text = $"（{GameUtil.Instance.GetTimeSToHMS(timeSeconds)}后地图开启）";
            if (timeSeconds == 0)
            {
                ClearTime();
                CheckButtonOpen();
            }

        }, true, true));
    }

    private void ClearTime()
    {
        if (timeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(timeItem);
            timeItem = null;
        }
    }

    private void ShowDropItems()
    {
        GetView().loopList_Awards.SetListItemCount(cfg_battle.Preview.Length, false);
        GetView().loopList_Awards.RefreshAllShownItem();
    }

    // 刷新工会经验点
    //[Subscriber(ThreadMode.Main)]
    private void OnSync_RefreshPrayItem(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        ShowPrayItemCount(RoleM.Instance.GuildContribution);
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_RefreshPrayItem);
        isShow = false;
        ClearTime();
    }

    public override void OnDestroy()
    {

    }
}