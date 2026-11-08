using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using SuperScrollView;

using System.Collections.Generic;
using System.Text;
using UnityEngine;
using UnityEngine.UI;

using static UIViewGameDayDoPro;

public class UIActivityEntrance : UIBase<UIModelActivityEntrance, UIViewActivityEntrance>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    List<CFGItem_LimitedActivity> cfg_LimitedActivities;

    public override void OnInit()
    {
        cfg_LimitedActivities = CFG_LimitedActivity.Get().DicList;

        GetView().loopList.InitListView(0, OnRefreshItem);
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIActivityEntrance>();
        });
    }

    private LoopListViewItem2 OnRefreshItem(LoopListView2 list, int index)
    {
        LoopListViewItem2 item = null;
        if (index < cfg_LimitedActivities.Count)
        {
            item = list.NewListViewItem("Item");
            SetItem(item.transform,index);
        }
        else if (index == cfg_LimitedActivities.Count) 
        {
            item = list.NewListViewItem("ItemEmpty");
        }

        return item;
    }

    void SetItem(Transform item, int index)
    {
        SetItemImage(item, index);
        int id = index + 1;
        bool isOpenLock = IsOpenLock(id, false);
        bool isOpenTime = IsOpenTime(id);
        Button btnGo = JUIExtensions.btn(item.Find("btn_Go"));
        Transform mask = item.Find("btn_Go/mask");
        bool isLock = isOpenLock || !isOpenTime;
        mask.gameObject.SetActive(isLock);
        if (isLock)
            JUIExtensions.txt(mask, "txt").text = GetOpenLockDesc(id);

        JUIExtensions.img(btnGo.transform).SetGray(isLock);
        JUIExtensions.txt(item.Find("btn_Go/txt")).color = isLock ? GameConst.GrayTextColor : buttonTextColor;
        ShowAwards(item.Find("Awards"), cfg_LimitedActivities[index].GiftPre);
        btnGo.onClick.RemoveAllListeners();
        btnGo.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (IsOpenLock(id, true)) return;

            if (isOpenTime)
                OpenUI(id);
            else
                UITool.ShowTips(GetOpenTime(id));
        });

        Button btnHelp = JUIExtensions.btn(item.Find("btn_Help"));
        btnHelp.onClick.RemoveAllListeners();
        btnHelp.onClick.AddListener(() =>
        {
            switch (id)
            {
                case 1:
                    UITool.Open<UIComRule>(null, 15);
                    break;
                case 2:
                    UITool.Open<UIComRule>(null, 49);
                    break;
                case 3:
                    UITool.Open<UIComRule>(null, 48);
                    break;
                default:
                    break;
            }
        });
    }

    void SetItemImage(Transform item, int index)
    {
        switch (index)
        {
            case 0:
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ActivityEntranceAtlas, "icon_xshd_1", JUIExtensions.img(item.Find("img_picture")));
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ActivityEntranceAtlas, "word_xshd_sjbs", JUIExtensions.img(item.Find("img_title")));
                break;
            case 1:
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ActivityEntranceAtlas, "icon_xshd_2", JUIExtensions.img(item.Find("img_picture")));
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ActivityEntranceAtlas, "word_xshd_yzdd", JUIExtensions.img(item.Find("img_title")));
                break;
            case 2:
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ActivityEntranceAtlas, "icon_xshd_3", JUIExtensions.img(item.Find("img_picture")));
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ActivityEntranceAtlas, "word_xshd_qxzc", JUIExtensions.img(item.Find("img_title")));
                break;
            default:
                break;
        }
    }

    // 奖励显示数量
    const int awardCount = 3;
    Color buttonTextColor => new Color(0.498f, 0.227f, 0.012f);

    void ShowAwards(Transform root, int[] ids)
    {
        for (int i = 0; i < awardCount; i++)
        {
            Transform item = root.GetChild(i);
            item.gameObject.SetActive(i < ids.Length);
            int idx = i;
            if (i < ids.Length)
            {
                GameUtil.Instance.SetItemIcon(ids[i], JUIExtensions.img(item.Find("icon")), JUIExtensions.img(item));
                JUIExtensions.btn(item).onClick.RemoveAllListeners();
                JUIExtensions.btn(item).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = ids[idx] });
                });
            }
        }
    }

    /// <summary>
    /// 活动开启时间详情
    /// </summary>
    /// <param name="id"></param>
    /// <returns></returns>
    string GetOpenTime(int id)
    {
        return id switch
        {
            1 => $"每日{CFG_XMLConfig.xmlConfig.WorldBoss.startTime}~{CFG_XMLConfig.xmlConfig.WorldBoss.endTime}开启",
            2 => $"每日{CFG_XMLConfig.xmlConfig.FightEnd.startTime}~{CFG_XMLConfig.xmlConfig.FightEnd.endTime}开启",
            3 => GetWarFightOpenTime(),
            _ => "",
        };

        static string GetWarFightOpenTime()
        {
            var xml = CFG_XMLConfig.xmlConfig.WarFight;
            string Str = string.Empty;
            for (int i = 0; i < xml.openDays.Length; i++)
            {

                Str += GameUtil.Instance.GetWeekDayName(xml.openDays[i]);
                if (i != xml.openDays.Length - 1)
                {
                    Str += "、";
                }

            }
            return $"每{Str}{xml.startTimes[0]}~{xml.endTimes[0]} {xml.startTimes[1]}~{xml.endTimes[1]}开放";
        }
    }

    /// <summary>
    /// 活动是否在开启时间
    /// </summary>
    /// <param name="id"></param>
    /// <returns></returns>
    bool IsOpenTime(int id)
    {
        switch (id)
        {
            case 1:
                return GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.WorldBoss.startTime, CFG_XMLConfig.xmlConfig.WorldBoss.endTime);
            case 2:
                return GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.FightEnd.startTime, CFG_XMLConfig.xmlConfig.FightEnd.endTime);
            case 3:
                var xml = CFG_XMLConfig.xmlConfig.WarFight;
                for (int i = 0; i < xml.openDays.Length; i++)
                {
                    if (GameUtil.Instance.CheakCurTimeWeek(xml.openDays[i]))
                    {
                        return GameUtil.Instance.CheakCurTime(xml.startTimes[0], xml.endTimes[0]) || GameUtil.Instance.CheakCurTime(xml.startTimes[1], xml.endTimes[1]);
                    }
                }
                return false;
            default:
                return true;
        }
    }

    /// <summary>
    /// 开服时间或等级是否未解锁？
    /// </summary>
    /// <param name="id"></param>
    /// <returns>true.未解锁 false.已结束</returns>
    bool IsOpenLock(int id, bool showTip)
    {
        return id switch
        {
            1 => GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss, showTip)
                                || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss, showTip),
            2 => GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd, showTip)
                                || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd, showTip),
            3 => GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight, showTip)
                                || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight, showTip),
            _ => false,
        };
    }

    string GetOpenLockDesc(int id)
    {
        StringBuilder tips = new StringBuilder("");
        switch (id)
        {
            case 1:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss, false))
                    return $"{GameUtil.Instance.RemainingDaysUnlockOpen(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss)}天后开启";
                else if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss, false, tips))
                    return tips.ToString();//$"{CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss]}级解锁";
                else
                    return "活动未开启";
            case 2:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd, false))
                    return $"{GameUtil.Instance.RemainingDaysUnlockOpen(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd)}天后开启";
                else if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd, false, tips))
                    return tips.ToString();//$"{CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd]}级解锁";
                else
                    return "活动未开启";

            case 3:
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight, false))
                    return $"{GameUtil.Instance.RemainingDaysUnlockOpen(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight)}天后开启";
                else if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight, false, tips))
                    return tips.ToString();//$"{CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.WarFight]}级解锁";
                else
                    return "活动未开启";
            default:
                return "";
        }
    }

    /// <summary>
    /// 跳转到对应的界面
    /// </summary>
    /// <param name="id"></param>
    void OpenUI(int id)
    {
        switch (id)
        {
            case 1:
                UITool.Open<UIWorldBoss>();
                UITool.Close<UIActivityEntrance>();
                break;
            case 2:
                UITool.Open<UIFightToTheEnd>();
                UITool.Close<UIActivityEntrance>();
                break;
            case 3:
                UITool.Open<UIWarFightMap>();
                UITool.Close<UIActivityEntrance>();
                break;
            default:
                break;
        }
    }

    public override void OnEnter(params object[] args)
    {
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(ZeroRefresh);
        GetView().loopList.SetListItemCount(cfg_LimitedActivities.Count + 1);
        GetView().loopList.RefreshAllShownItem();
        OnAutoGetInfo();
    }

    private void ZeroRefresh(GameEvent.Refresh_ZeroTime _e)
    {
        GetView().loopList.RefreshAllShownItem();
    }

    private TimeItem uiUpdate;

    private void OnAutoGetInfo()
    {
        uiUpdate = GameHelper.Instance.AddInvoke(new TimeItem(-1, 2, () =>
        {
            if (isNullGo)
                return;

            GetView().loopList.RefreshAllShownItem();
        }, true, true));
    }

    private void ClearJation()
    {
        if (GameHelper.Instance != null && uiUpdate != null)
        {
            GameHelper.Instance.RemoveInvoke(uiUpdate);
            uiUpdate = null;
        }
    }

    public override void OnExit()
    {
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(ZeroRefresh);
        ClearJation();
    }

    public override void OnDestroy()
    {

    }
}