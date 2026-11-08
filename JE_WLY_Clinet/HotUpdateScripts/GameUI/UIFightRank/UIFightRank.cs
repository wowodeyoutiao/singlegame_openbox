using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using SuperScrollView;
using JEngine.UI;
using HotUpdateScripts;
using JEngine.Helper;
using JEngine.Event;

public class UIFightRank : UIBase<UIModelFightRank, UIViewFightRank>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFightRank>();
        });

        //GetView().btn_FightInfo.onClick.AddListener(() =>
        //{
        //    OnClickFightRecordButton();

        //});

        GetView().btn_add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIPVPCountBuy>(null, (int)_roomType);

        });

        GetView().infiniteScroll.InitListView(0, OnFillRank);
    }

    //private void OnClickFightRecordButton()
    //{
    //    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

    //    TcpC.C2S_PvpRecorderRequest(_roomType == CFG_Room.RoomType.PVP_One ? 1 : 2);
    //}

    private LoopListViewItem2 OnFillRank(LoopListView2 loopListView, int idx)
    {
        if (PVPPlayers == null || idx < 0 || idx > PVPPlayers.Count)
        {
            return null;
        }
        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        Transform go = item.transform;

        var player = PVPPlayers[idx];
        JUIExtensions.txt(go.Find("txt_Name")).text = player.Name;
        JUIExtensions.txt(go, "Head/img_Lv/Txt_lv").text = player.Level.ToString();
        JUIExtensions.txt(go, "power/txt").text = GameUtil.Instance.FormatPowerWithUnit(player.Combat);
        JUIExtensions.txt(go, "fen/txt").text = GameUtil.Instance.FormatNumWithUnit(player.Score);

        JUIExtensions.txt(go, "award/award0/txt").text = GameUtil.Instance.FormatNumWithUnit(CFG_XMLConfig.xmlConfig.PVPInfo.WinScore);

        GameUtil.Instance.SetJobIcon(player.Job, JUIExtensions.img(go, "Head/img_job"));
        //ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + player.Job + player.Gender, JUIExtensions.img(go, "Head/img_head"));

        (int headId, int frameId) = GameUtil.Instance.GetHeadAndFrameId(player.AvatarV2, player.Frame);
        GameUtil.Instance.SetRoleHead(headId, player.Job, JUIExtensions.img(go.transform, "Head/img_head"), JUIExtensions.img(go.transform, "Head/img_head_quality"));
        GameUtil.Instance.SetRoleHeadFrame(frameId, player.Job, JUIExtensions.img(go.transform, "Head/img_head_kuang"));

        for (int i = 0; i < 2; i++)
        {
            var root = go.Find("award").GetChild(i);
            var cfg_item = CFG_XMLConfig.xmlConfig.PVPInfo.winAwards[i];
            GameUtil.Instance.SetItemIcon(cfg_item.id, JUIExtensions.img(root));
            int _num = cfg_item.mark == 0 ? cfg_item.num : cfg_item.num * RoleM.Instance.Level;
            JUIExtensions.txt(root.Find("txt")).text = GameUtil.Instance.FormatNumWithUnit(_num);
        }

        var btn = JUIExtensions.btn(go, "Btn_Play");
        btn.onClick.RemoveAllListeners();
        btn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (RoleM.Instance.GetPVPMatchTimes((int)_roomType) == 0)
            {
                UITool.ShowTips("挑战次数不足");
                return;
            }
            if (_roomType == CFG_Room.RoomType.PVP_One)
            {
                TcpC.C2S_PVPChallenge(player.Uid);

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "即时挑战",
                    position = "点击进入选择对手",
                    event_click = $"开始挑战",
                });
            }
            else if (_roomType == CFG_Room.RoomType.PVP_More)
                UITool.Open<UIFightFormation>(null, _roomType, player.Uid);
        });
        return item;
    }

    private CFG_Room.RoomType _roomType;
    private List<pb.PvpMatchingResponse.PVPPlayer> PVPPlayers;


    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);

        PVPPlayers = (List<pb.PvpMatchingResponse.PVPPlayer>)args[0];
        _roomType = (CFG_Room.RoomType)args[1];
        GetView().txt_Title.text = _roomType == CFG_Room.RoomType.PVP_One ? "即时挑战模式" : "回合挑战模式";
        GetView().infiniteScroll.SetListItemCount(PVPPlayers.Count);
        GetView().infiniteScroll.RefreshAllShownItem();
        OnChallengeShow();
    }

    private void OnChallengeShow()
    {
        GetView().txt_ChallangeNum.text = "今日可以挑战次数：" + RoleM.Instance.GetPVPMatchTimes((int)_roomType);

        CFG_XMLConfig.XML_Config.XML_Item_PVP.PVPseason endPVPSeason = null;
        if (_roomType == CFG_Room.RoomType.PVP_One)
            endPVPSeason = RoleM.Instance.GetPVPOneSeason();

        var baseInfo = _roomType == CFG_Room.RoomType.PVP_One ? CFG_XMLConfig.xmlConfig.PVPInfo : CFG_XMLConfig.xmlConfig.GVGInfo;
        var isOpne = (endPVPSeason != null || _roomType == CFG_Room.RoomType.PVP_More) && GameUtil.Instance.CheakCurTime(baseInfo.dailyBegin, baseInfo.dailyEnd);
        GetView().btn_add.gameObject.SetActive(isOpne && RoleM.Instance.GetCanNumBuyPVP((int)_roomType) > 0);
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_PVP(GameEvent.Sync_PVPInfo _e)
    {
        if (isNullGo)
            return;
        OnChallengeShow();
    }

    public override void OnExit()
    {
        PVPPlayers = null;
        JEvent.defaultEvent.Unregister(this);
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}