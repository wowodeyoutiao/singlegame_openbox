using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using pb;
using SuperScrollView;
using System;
using JEngine.UI;
using JEngine.Helper;

public class UIGuildInfo : UIBase<UIModelGuildInfo, UIViewGuildInfo>
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
            UITool.Open<UIComRule>(null, 33);
        });

        GetView().loopList_Attributes.InitListView(0, OnItemRefresh);

        GetView().btn_CoinContribution.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (GameUtil.Instance.ComCoinToCheck(CFG_XMLConfig.xmlConfig.Guild.golddonate.needitemid, CFG_XMLConfig.xmlConfig.Guild.golddonate.needitemnum, () =>{ UITool.Close<UIGuild>(); }, true))
            {
                OnDonate(1);

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "社交",
                    component = "行会",
                    position = "信息",
                    event_click = "金币捐献",
                });
            }
            else
            {
                UITool.ShowTips("金币不足");
            }
        });

        GetView().btn_RollContribution.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (GameUtil.Instance.ComCoinToCheck(CFG_XMLConfig.xmlConfig.Guild.moneydonate.needitemid, CFG_XMLConfig.xmlConfig.Guild.moneydonate.needitemnum, () => { UITool.Close<UIGuild>(); }, true))
            {
                OnDonate(2);

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "社交",
                    component = "行会",
                    position = "信息",
                    event_click = "元宝捐献",
                });
            }
            else
            {
                UITool.ShowTips("点卷不足");
            }
        });

        JUIExtensions.btn(GetView().img_CoinContributionItem1.transform.parent).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFG_XMLConfig.xmlConfig.Guild.golddonate.needitemid });
        });

        JUIExtensions.btn(GetView().img_CoinContributionItem2.transform.parent).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFG_XMLConfig.xmlConfig.Guild.golddonate.getguildexpid });
        });

        JUIExtensions.btn(GetView().img_RollContributionItem1.transform.parent).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFG_XMLConfig.xmlConfig.Guild.moneydonate.needitemid });
        });

        JUIExtensions.btn(GetView().img_RollContributionItem2.transform.parent).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFG_XMLConfig.xmlConfig.Guild.moneydonate.getguildexpid });
        });
    }

    string leaderLabel = "<color=#997C4B>会长：</color>";
    string peopleCountLabel = "<color=#997C4B>人数：</color>";
    string contributionDegreeLabel = "<color=#997C4B>贡献度:</color>";

    bool isShow;
    private Action<string> setTitleCallback;

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

    public void SetTitleListener(Action<string> callback)
    {
        setTitleCallback = callback;
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        guildSkillList = CFG_GuildSkill.GetList();
        TcpC.C2S_GuildInfoRequest(RoleM.Instance.GuildId);
        TcpC.C2S_GuildDonateInfoRequest();
    }

    GuildInfo guildInfo;

    List<CFGItem_GuildSkill> guildSkillList;

    [Subscriber(ThreadMode.Main)]
    // 工会信息
    private void Event_GuildInfo(GameEvent.Sync_GuildInfo _e)
    {
        if(isNullGo)
            return;

        guildInfo = _e.guildInfo;
        JEvent.defaultEvent.Post(new UIGuild.Title
        {
            panelType = UIGuild.GuildChildPanelType.GuildInfo,
            titleName = guildInfo.Name,
        });
        RefreshUI();
    }

    int donateNum;

    [Subscriber(ThreadMode.Main)]
    // 工会捐献数量
    private void Event_GuildDonateInfo(GameEvent.Sync_GuildDonateInfo _e)
    {
        if (isNullGo)
            return;

        donateNum = _e.num;
        ShowDonateCount();
        TcpC.C2S_GuildInfoRequest(RoleM.Instance.GuildId);
    }

    [Subscriber(ThreadMode.Main)]
    private void OnZeroRefresh()
    {
        RefreshUI();
    }

    private void RefreshUI()
    {
        GetView().txt_Leader.text = leaderLabel + guildInfo.LeaderName;
        GetView().txt_PeopleCount.text = peopleCountLabel + guildInfo.MemberCount + "/" + guildInfo.MaxMembers;
        GetView().txt_ContributionDegreeCurrentLevel.text = contributionDegreeLabel + guildInfo.Level + "级";
        int nextLevel = CFG_GuildEXP.IsMaxLevel(guildInfo.Level) ? guildInfo.Level : guildInfo.Level + 1;
        GetView().txt_ContributionDegreeNextLevel.text = nextLevel + "级";
        GetView().slider_LevelProgress.value = GetLevelProgress(guildInfo.Level);
        GetView().txt_Progress.text = $"{guildInfo.GuildExp}/{CFG_GuildEXP.GetById(guildInfo.Level).NeedEXP}";
        GetView().txt_AttributeCurrentLevel.text = guildInfo.Level + "级行会";
        GetView().txt_AttributeNextLevel.text = nextLevel + "级行会";
        GetView().loopList_Attributes.SetListItemCount(CFG_GuildEXP.GetById(guildInfo.Level).AttributeAdd.Length + 1);
        GetView().loopList_Attributes.RefreshAllShownItem();

        GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.Guild.golddonate.needitemid, GetView().img_CoinContributionItem1);
        GetView().txt_CoinContributionItemCount1.text = GameUtil.Instance.GetCheakNumStr(CFG_XMLConfig.xmlConfig.Guild.golddonate.needitemid, CFG_XMLConfig.xmlConfig.Guild.golddonate.needitemnum);
        GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.Guild.golddonate.getguildexpid, GetView().img_CoinContributionItem2);
        GetView().txt_CoinContributionItemCount2.text = CFG_XMLConfig.xmlConfig.Guild.golddonate.getitemnum.ToString();

        GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.Guild.moneydonate.needitemid, GetView().img_RollContributionItem1);
        GetView().txt_RollContributionItemCount1.text = GameUtil.Instance.GetCheakNumStr(CFG_XMLConfig.xmlConfig.Guild.moneydonate.needitemid, CFG_XMLConfig.xmlConfig.Guild.moneydonate.needitemnum);
        GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.Guild.moneydonate.getguildexpid, GetView().img_RollContributionItem2);
        GetView().txt_RollContributionItemCount2.text = CFG_XMLConfig.xmlConfig.Guild.moneydonate.getitemnum.ToString();
    }

    private void ShowDonateCount()
    {
        int donateNumRemaing = GuildM.Instance.GetRemaingDonateNum();
        string donateNumRemaingString = donateNumRemaing > 0 ? donateNumRemaing.ToString() : $"<color=#DA4E38>{donateNumRemaing}</color>";
        GetView().txt_CoinContributionTimes.text = $"今日次数 {donateNumRemaingString}/{CFG_XMLConfig.xmlConfig.Guild.donatelimitDailytime}";
        GetView().txt_RollContributionTimes.text = $"今日次数 {donateNumRemaingString}/{CFG_XMLConfig.xmlConfig.Guild.donatelimitDailytime}";
        bool isMaxContribution = donateNum >= CFG_XMLConfig.xmlConfig.Guild.donatelimitDailytime;

        GetView().btn_CoinContribution.gameObject.SetActive(!isMaxContribution);
        GetView().btn_CoinContributionGray.gameObject.SetActive(isMaxContribution);
        GetView().btn_RollContribution.gameObject.SetActive(!isMaxContribution);
        GetView().btn_RollContributionGray.gameObject.SetActive(isMaxContribution);
    }

    private LoopListViewItem2 OnItemRefresh(LoopListView2 listView, int index)
    {
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        Transform t = item.transform;
        CFGItem_GuildEXP.AttributeAddItem[] currentAttributes = CFG_GuildEXP.GetById(guildInfo.Level).AttributeAdd;
        int nextLevel = CFG_GuildEXP.IsMaxLevel(guildInfo.Level) ? guildInfo.Level : guildInfo.Level + 1;
        if(currentAttributes.Length == index)
        {
            JUIExtensions.txt(t.Find("txt_Label")).text = "行会人数上限";
            int curValue = CFG_GuildEXP.GetById(guildInfo.Level).MemberNum;
            int nextValue = CFG_GuildEXP.GetById(nextLevel).MemberNum;
            JUIExtensions.txt(t.Find("txt_CurrentValue")).text = curValue.ToString();
            JUIExtensions.txt(t.Find("txt_NextValue")).text = nextValue - curValue > 0 ? $"{nextValue}<color=#15ABEA>（+{nextValue - curValue}）</color>" : nextValue.ToString();
        }
        else
        {
            CFGItem_GuildEXP.AttributeAddItem[] nextAttributes = CFG_GuildEXP.GetById(nextLevel).AttributeAdd;
            CFG_XMLConfig.XML_Config.XML_Item_RolePanel attr = CFG_XMLConfig.GetAttributeByID(currentAttributes[index].key);
            JUIExtensions.txt(t.Find("txt_Label")).text = attr.Name;
            int curValue = currentAttributes[index].value;
            int nextValue = nextAttributes[index].value;
            JUIExtensions.txt(t.Find("txt_CurrentValue")).text = CFG_XMLConfig.GetAttrValueToStr(attr.Type, curValue);
            JUIExtensions.txt(t.Find("txt_NextValue")).text = nextValue - curValue > 0 ? $"{nextValue}<color=#15ABEA>（+{CFG_XMLConfig.GetAttrValueToStr(attr.Type, nextValue - curValue)}）</color>" : CFG_XMLConfig.GetAttrValueToStr(attr.Type, nextValue);
        }

        return item;
    }

    private float GetLevelProgress(int level)
    {
        int currentEXPNeed = CFG_GuildEXP.GetById(level).NeedEXP;
        return (float)guildInfo.GuildExp / currentEXPNeed;
    }

    private void OnDonate(int type)
    {
        TcpC.C2S_GuildDonateRequest(type);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        isShow = false;
        setTitleCallback = null;
    }

    public override void OnDestroy()
    {

    }
}