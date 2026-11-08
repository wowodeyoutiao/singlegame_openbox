using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using Spine;

using SuperScrollView;

using UnityEngine;
using UnityEngine.UI;

public class UICompetitionEntrance : UIBase<UIModelCompetitionEntrance, UIViewCompetitionEntrance>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().loopList.InitListView(0, OnRefreshItem);

        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UICompetitionEntrance>();
        });
    }

    // 活动数量(实际少一个，最后一个为空白)
    int activeCount = 3;
    Color buttonTextColor => new Color(0.498f, 0.227f, 0.012f);

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(ZeroRefresh);
        GetView().loopList.SetListItemCount(activeCount);
        GetView().loopList.RefreshAllShownItem();
    }

    private LoopListViewItem2 OnRefreshItem(LoopListView2 list, int index)
    {
        LoopListViewItem2 item = null;
        if (index < 2)
        {
            item = list.NewListViewItem("Item");
            var roomType = index == 0 ? CFG_Room.RoomType.PVP_One : CFG_Room.RoomType.PVP_More;
            var cfg =  CFG_Room.GetByRoomType(roomType, RoleM.Instance.Level);
            SetPVPItem(item.transform, cfg);
        }
        // 最后一个，敬请期待
        else if(index == activeCount - 1)
        {
            item = list.NewListViewItem("ItemEmpty");
        }
        return item;
    }

    private void SetPVPItem(Transform item, CFGItem_Room cfg)
    {
        Transform itemRoot = item.Find("img_award/awards");

        CFG_XMLConfig.XML_Config.XML_Item_PVP.PVPreward[] pVPrewards = null;
        var rankClass = CFG_RankingAward.RankClass.No;

        if (cfg.Roomtype == CFG_Room.RoomType.PVP_One)
        {
            pVPrewards = CFG_XMLConfig.xmlConfig.PVPInfo.winAwards;
            rankClass = CFG_RankingAward.RankClass.ZhanHun;

            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CompetitionEntranceAtlas, "btn_event_js", JUIExtensions.img(item.Find("bg")));
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CompetitionEntranceAtlas, "word_event_jstz", JUIExtensions.img(item.Find("img_title")));
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CompetitionEntranceAtlas, "btn_event_rz_1", JUIExtensions.img(item.Find("Btn_Recoed")));
        }
        else if (cfg.Roomtype == CFG_Room.RoomType.PVP_More)
        {
            pVPrewards = CFG_XMLConfig.xmlConfig.PVPInfo.loseAwards;
            rankClass = CFG_RankingAward.RankClass.TurnBasedDay;
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CompetitionEntranceAtlas, "btn_event_hh", JUIExtensions.img(item.Find("bg")));
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CompetitionEntranceAtlas, "word_event_hhtz", JUIExtensions.img(item.Find("img_title")));
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CompetitionEntranceAtlas, "btn_event_rz_2", JUIExtensions.img(item.Find("Btn_Recoed")));
        }

        for (int i = 0; i < itemRoot.childCount; i++)
        {
            var go = itemRoot.GetChild(i);
            int cfgId = pVPrewards[i].id;
            GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(go.Find("icon")), JUIExtensions.img(go));

            JUIExtensions.btn(go).onClick.RemoveAllListeners();
            JUIExtensions.btn(go).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
            });
        }

        JUIExtensions.txt(item, "txt_Count").text = "" + RoleM.Instance.GetPVPMatchTimes((int)cfg.Roomtype);

        var baseInfo = cfg.Roomtype == CFG_Room.RoomType.PVP_One ? CFG_XMLConfig.xmlConfig.PVPInfo : CFG_XMLConfig.xmlConfig.GVGInfo;
        JUIExtensions.txt(item, "txt_sjTime").text = $"(每日:{baseInfo.dailyBegin}-{baseInfo.dailyEnd})";

        Button recoedBtn = JUIExtensions.btn(item.Find("Btn_Recoed"));
        recoedBtn.onClick.RemoveAllListeners();
        recoedBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_PvpRecorderRequest(cfg.Roomtype == CFG_Room.RoomType.PVP_One ? 1 : 2);
        });

        Button helpBtn = JUIExtensions.btn(item.Find("Btn_help"));
        helpBtn.onClick.RemoveAllListeners();
        helpBtn.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, (cfg.Roomtype == CFG_Room.RoomType.PVP_One ? 3 : 4));
        });

        Transform playbtnMask = item.Find("Btn_Play/mask");
        System.Text.StringBuilder tips = new("");
        bool isLockOpenDay = GameUtil.Instance.CheckOpenServerIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false);
        bool isLockLevel = GameUtil.Instance.CheakIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false, tips);
        playbtnMask.gameObject.SetActive(isLockOpenDay || isLockLevel);
        item.Find("Btn_Play/dot").gameObject.SetActive(!isLockOpenDay && !isLockLevel && RoleM.Instance.GetPVPMatchTimes((int)cfg.Roomtype) > 0);

        if (isLockOpenDay)
        {
            int remainday = GameUtil.Instance.RemainingDaysUnlockOpen(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP);
            JUIExtensions.txt(playbtnMask.Find("txt_OpenDay")).text = $"{remainday}天后开启";
        }
        else if (isLockLevel)
        {
            int openLevel = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP];
            JUIExtensions.txt(playbtnMask.Find("txt_OpenDay")).text = tips.ToString();//$"{openLevel}级解锁";
        }

        Button playBtn = JUIExtensions.btn(item, "Btn_Play");
        playBtn.onClick.RemoveAllListeners();
        playBtn.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (GameUtil.Instance.CheckOpenServerIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                return;

            if (GameUtil.Instance.CheakIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                return;

            if (RoleM.Instance.GetPVPMatchTimes((int)cfg.Roomtype) == 0)
            {
                UITool.ShowTips("挑战次数不足");
                return;
            }
            OnEnterPVPGame((int)cfg.Roomtype);
        });
        JUIExtensions.img(playBtn.transform).SetGray(isLockOpenDay || isLockLevel);
        JUIExtensions.txt(item.Find("Btn_Play/Text")).color = isLockOpenDay || isLockLevel ? GameConst.GrayTextColor : buttonTextColor;

        Button addCountBtn = JUIExtensions.btn(item, "Btn_AddCount");
        addCountBtn.onClick.RemoveAllListeners();
        addCountBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (cfg.Roomtype == CFG_Room.RoomType.PVP_One)
            {
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false))
                {
                    UITool.ShowTips("功能暂未开启");
                    return;
                }

                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false))
                {
                    UITool.ShowTips("功能暂未开启");
                    return;
                }

                GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.PVETZJZ);
            }
            else if (cfg.Roomtype == CFG_Room.RoomType.PVP_More)
            {
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                {
                    UITool.ShowTips("功能暂未开启");
                    return;
                }

                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, false))
                {
                    UITool.ShowTips("功能暂未开启");
                    return;
                }

                GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.PVPTZJZ);
            }
        });

        Transform _item = item;

        HttpC.GetRankList(rankClass, (_r) =>
        {
            if (isNullGo)
                return;

            bool isNoOne = _r.data.rank_all.Length == 0;
            _item.Find("bg_head/img_head").gameObject.SetActive(!isNoOne);
            _item.Find("bg_head/txt_empty").gameObject.SetActive(isNoOne);

            if (!isNoOne)
            {
                //ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + _r.data.rank_all[0].job + _r.data.rank_all[0].sex, JUIExtensions.img(_item, "bg_head/img_head"));
                (int headId, int frameId) = GameUtil.Instance.GetHeadAndFrameId(_r.data.rank_all[0].avatar_v2, _r.data.rank_all[0].frame);
                GameUtil.Instance.SetRoleHead(headId, _r.data.rank_all[0].job, JUIExtensions.img(_item, "bg_head/img_head"), JUIExtensions.img(_item, "bg_head/img_head_quality"));
                GameUtil.Instance.SetRoleHeadFrame(frameId, _r.data.rank_all[0].job, JUIExtensions.img(_item, "bg_head/img_head_kuang"));
            }
            else
            {
                GameUtil.Instance.SetDefauleRoleHead(1, JUIExtensions.img(_item, "bg_head/img_head_kuang"), JUIExtensions.img(_item, "bg_head/img_head_quality"));
            }
        });

        Button headBtn = JUIExtensions.btn(item, "bg_head");
        headBtn.onClick.RemoveAllListeners();
        headBtn.onClick.AddListener(() =>
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                return;

            if (GameUtil.Instance.CheakIsLock(cfg.Roomtype == CFG_Room.RoomType.PVP_More ? CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG : CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                return;

            if (cfg.Roomtype == CFG_Room.RoomType.PVP_One)
                UITool.Open<UIRank>(null, CFG_RankingAward.RankClass.ZhanHun);
            else if (cfg.Roomtype == CFG_Room.RoomType.PVP_More)
                UITool.Open<UIRank>(null, CFG_RankingAward.RankClass.TurnBasedDay);
        });
    }

    private void OnEnterPVPGame(int Roomtype)
    {
        GameRoomM.Instance._GameState = GameRoomM.GameState.No;

        if (Roomtype == (int)CFG_Room.RoomType.PVP_One)
            TcpC.C2S_PvpMatching(1);
        else
            TcpC.C2S_PvpMatching(2);
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_PVP(GameEvent.Sync_PVPInfo _e)
    {
        if (isNullGo)
            return;
        if ((GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false) || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, false)))
            return;
        GetView().loopList.RefreshAllShownItem();
    }

    private void ZeroRefresh(GameEvent.Refresh_ZeroTime _e)
    {
        TcpC.OnPVPInfo((int)pb.MatchType.SINGLE_COMBAT);
        TcpC.OnPVPInfo((int)pb.MatchType.FREE_FOR_ALL);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(ZeroRefresh);
    }

    public override void OnDestroy()
    {

    }
}