using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using pb;
using JEngine.UI;
using SuperScrollView;
using JEngine.Helper;

public class UIGuildSkill : UIBase<UIModelGuildSkill, UIViewGuildSkill>
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
            UITool.Open<UIComRule>(null, 34);
        });

        GetView().loopList_Skill.InitListView(0, OnItemRefresh);
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

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnSync_SkillUpgrede);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_RefreshPrayItem);

        isScrollListToTop = true;
        GetView().txt_Title.text = $"祈愿点每<color=#BD341E>{CFG_XMLConfig.xmlConfig.Guild.wishvalue.refreshtime}</color>分钟恢复<color=#BD341E>1</color>点，最多累计<color=#BD341E>{CFG_XMLConfig.xmlConfig.Guild.wishvalue.maxvalue}</color>点";
        ShowPrayItemCount(RoleM.Instance.GuildWish);
        OnGuildInfoRequest();
    }

    long wishCount = -1;

    private void ShowPrayItemCount(long count)
    {
        GetView().txt_ItemCount.text = count.ToString();

        if (RoleM.Instance.GuildWish >= CFG_XMLConfig.xmlConfig.Guild.wishvalue.maxvalue)
        {
            //ClearTime();
            GetView().txt_TimeTip.text = "祈愿点已满，停止自动恢复状态";
        }
        else
        {
            if (RoleM.Instance.GuildWish != wishCount)
            {
                //ShowPrayItemTime();
                GetView().txt_TimeTip.text = $"每消耗<color=#3A44D7>1</color>点祈愿点获得<color=#3A44D7>{CFG_XMLConfig.xmlConfig.Guild.wishvalue.getitemnum}</color>点行会贡献";
                wishCount = RoleM.Instance.GuildWish;
            }
        }
    }

    private void OnGuildInfoRequest()
    {
        TcpC.C2S_GuildInfoRequest(RoleM.Instance.GuildId);
    }

    GuildInfo guildInfo;
    List<GuildSkillInfo> guildSkillList;

    [Subscriber(ThreadMode.Main)]
    // 工会信息
    private void Event_GuildInfo(GameEvent.Sync_GuildInfo _e)
    {
        if (isNullGo)
            return;

        guildInfo = _e.guildInfo;
        guildSkillList = guildInfo.Skills;
        guildSkillList.Sort((x, y) =>
        {
            return x.CfgId.CompareTo(y.CfgId);
        });
        RefreshList();
    }

    // 刷新祈愿点
    //[Subscriber(ThreadMode.Main)]
    private void OnSync_RefreshPrayItem(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        ShowPrayItemCount(RoleM.Instance.GuildWish);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_SkillUpgrede(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.code == 0 && _e.cid == CmdId.C2S_GuildUpgradeSkill)
        {
            isScrollListToTop = false;
            OnGuildInfoRequest();
        }
    }

    bool isScrollListToTop = true;

    private void RefreshList()
    {
        GetView().loopList_Skill.SetListItemCount(guildSkillList.Count, isScrollListToTop);
        GetView().loopList_Skill.RefreshAllShownItem();
    }

    //TimeItem timeItem;

    // 祈愿点恢复倒计时
    //private void ShowPrayItemTime()
    //{
    //    ClearTime();
    //    int timeSeconds = RoleM.Instance.WishRecoveryTime - (int)GameUtil.Instance.GetCurrentTimeSeconds();
    //    if (timeSeconds <= 0)
    //    {
    //        timeSeconds = 0;
    //    }
    //    else
    //    {
    //        timeItem = GameHelper.Instance.AddInvoke(new TimeItem(-1, 1, () =>
    //        {
    //            if (isNullGo) return;

    //            timeSeconds--;
    //            GetView().txt_TimeTip.text = $"{GameUtil.Instance.GetTimeSToHMS(timeSeconds)}后恢复<color=#3A44D7>1</color>点祈愿点,每消耗<color=#3A44D7>1</color>点祈愿点获得<color=#3A44D7>{CFG_XMLConfig.xmlConfig.Guild.wishvalue.getitemnum}</color>点行会贡献";
    //            if (timeSeconds <= 0)
    //            {
    //                ClearTime();
    //            }
    //        }, true, true));
    //    }

    //    GetView().txt_TimeTip.text = $"{GameUtil.Instance.GetTimeSToHMS(timeSeconds)}后恢复<color=#3A44D7>1</color>点祈愿点,每消耗<color=#3A44D7>1</color>点祈愿点获得<color=#3A44D7>{CFG_XMLConfig.xmlConfig.Guild.wishvalue.getitemnum}</color>点行会贡献";
    //}

    //private void ClearTime()
    //{
    //    if (GameHelper.Instance != null && timeItem != null)
    //    {
    //        GameHelper.Instance.RemoveInvoke(timeItem);
    //        timeItem = null;
    //    }
    //}

    private LoopListViewItem2 OnItemRefresh(LoopListView2 listView, int index)
    {
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        Transform t = item.transform;
        int idx = index;
        JUIExtensions.txt(t.Find("txt_Level")).text = $"Lv.{guildSkillList[index].Level}";

        CFGItem_GuildSkill cfg_GuildSkill = CFG_GuildSkill.GetSkillByIdAndLevel(guildSkillList[index].CfgId, guildSkillList[index].Level);
        CFG_XMLConfig.XML_Config.XML_Item_RolePanel attr = CFG_XMLConfig.GetAttributeByID(cfg_GuildSkill.ValueAdd[0].key);

        JUIExtensions.txt(t.Find("txt_CurrentLevel")).text = $"Lv.{guildSkillList[index].Level}：{attr.Name}提升{CFG_XMLConfig.GetAttrValueToStr(attr.Type, cfg_GuildSkill.ValueAdd[0].value)}";
        int nextLv = CFG_GuildSkill.IsMaxLevel(guildSkillList[index].CfgId, guildSkillList[index].Level) ? guildSkillList[index].Level : guildSkillList[index].Level + 1;
        CFGItem_GuildSkill cfg_GuildSkill_Next = CFG_GuildSkill.GetSkillByIdAndLevel(guildSkillList[index].CfgId, nextLv);
        JUIExtensions.txt(t.Find("txt_NextLevel")).text = $"Lv.{nextLv}：{attr.Name}提升{CFG_XMLConfig.GetAttrValueToStr(attr.Type, cfg_GuildSkill_Next.ValueAdd[0].value)}";

        JUIExtensions.slider(t.Find("slider")).value = GetLevelProgress(guildSkillList[index].CfgId, guildSkillList[index].Level, guildSkillList[index].Exp);

        bool isMaxLevel = guildSkillList[index].Level >= CFG_GuildSkill.GetMaxLevelById(guildSkillList[index].CfgId);
        JUIExtensions.txt(t.Find("txt_Progress")).text = isMaxLevel ? "Max" : $"{guildSkillList[index].Exp}/{cfg_GuildSkill.EXP}";
        t.Find("btn_PrayOnce").gameObject.SetActive(!isMaxLevel);
        JUIExtensions.btn(t.Find("btn_PrayOnce")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("btn_PrayOnce")).onClick.AddListener(() =>
        {
            if (RoleM.Instance.GuildWish > 0)
            {
                TcpC.C2S_GuildUpgradeSkillRequest(cfg_GuildSkill.UUID, 1); // 需要刷新一下列表，刷新按钮状态。
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "行会",
                    position = "技能",
                    event_click = "祈愿1次",
                });
            }
            else
                UITool.ShowTips("祈愿点不足");
        });

        t.Find("btn_PrayAll").gameObject.SetActive(!isMaxLevel);
        JUIExtensions.btn(t.Find("btn_PrayAll")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.Find("btn_PrayAll")).onClick.AddListener(() =>
        {
            if (RoleM.Instance.GuildWish > 0)
            {
                TcpC.C2S_GuildUpgradeSkillRequest(cfg_GuildSkill.UUID, (int)RoleM.Instance.GuildWish);
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "行会",
                    position = "技能",
                    event_click = "全部祈愿",
                });
            }
            else
                UITool.ShowTips("祈愿点不足");
        });

        t.Find("txt_MaxLevel").gameObject.SetActive(isMaxLevel);

        return item;
    }

    private float GetLevelProgress(int skillId, int level, int currentExp)
    {
        int currentEXPNeed = CFG_GuildSkill.GetSkillByIdAndLevel(skillId, level).EXP;
        return (float)currentExp / currentEXPNeed;
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnSync_SkillUpgrede);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_RefreshPrayItem);

        //ClearTime();
        isShow = false;
        wishCount = -1;
    }

    public override void OnDestroy()
    {
        //ClearTime();
    }
}