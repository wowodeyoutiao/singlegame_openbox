using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using SuperScrollView;

using System;
using System.Collections.Generic;

using UnityEngine;

using static CFG_Room;

public class UIFightRecord : UIBase<UIModelFightRecord, UIViewFightRecord>
{
    private CFG_Room.RoomType roomType;
    private List<pb.Information> PVPPlayers;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFightRecord>();
        });

        //GetView().btn_add.onClick.AddListener(() =>
        //{
        //    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        //    UITool.Open<UIPVPCountBuy>(null, (int)roomType);
        //});

        GetView().infiniteScroll.InitListView(0, OnFillRecord);
    }
    private LoopListViewItem2 OnFillRecord(LoopListView2 loopListView, int idx)
    {
        if (PVPPlayers == null || idx < 0 || idx >= PVPPlayers.Count)
            return null;

        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        Transform go = item.transform;

        var player = PVPPlayers[idx];
        JUIExtensions.txt(go.Find("txt_Name")).text = player.Name;
        JUIExtensions.txt(go, "Head/img_Lv/Txt_lv").text = player.Level.ToString();
        JUIExtensions.txt(go, "power/txt").text = GameUtil.Instance.FormatPowerWithUnit(player.Power);
        JUIExtensions.txt(go, "fen/txt").text = GameUtil.Instance.FormatNumWithUnit(player.Score);

        JUIExtensions.txt(go, "txt_Record").text = GameUtil.Instance.GetTimePassedFromStep(player.MatchTime);

        GameUtil.Instance.SetJobIcon(player.Job, JUIExtensions.img(go, "Head/img_job"));
        //ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + player.Job + player.Gender, JUIExtensions.img(go, "Head/img"));
        (int headId, int frameId) = GameUtil.Instance.GetHeadAndFrameId(player.AvatarV2, player.Frame);
        GameUtil.Instance.SetRoleHead(headId, player.Job, JUIExtensions.img(go.transform, "Head/img_head"), JUIExtensions.img(go.transform, "Head/img_head_quality"));
        GameUtil.Instance.SetRoleHeadFrame(frameId, player.Job, JUIExtensions.img(go.transform, "Head/img_head_kuang"));
        // 1 胜利， 0失败

        go.Find("obj_Win").gameObject.SetActive(player.Result == 1);
        go.Find("obj_Lose").gameObject.SetActive(player.Result == 0);
        /*
        if (player.Result == 1)
        {
            var btn = JUIExtensions.btn(go, "obj_Win/Btn_Play");
            btn.onClick.RemoveAllListeners();
            btn.onClick.AddListener(() =>
            {
                OnClickPlayButton(player.OpponentUid);
            });
        }
        else
        {
            var btn = JUIExtensions.btn(go, "obj_Lose/Btn_Play");
            btn.onClick.RemoveAllListeners();
            btn.onClick.AddListener(() =>
            {
                OnClickPlayButton(player.OpponentUid);
            });
        }
        */
        return item;
    }

    /*
    private void OnClickPlayButton(long uid)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (RoleM.Instance.GetPVPMatchTimes((int)roomType) == 0)
        {
            UITool.ShowTips("挑战次数不足");
            return;
        }
        if (roomType == CFG_Room.RoomType.PVP_One)
            TcpC.C2S_PVPChallenge(uid);
        else if (roomType == CFG_Room.RoomType.PVP_More)
            UITool.Open<UIFightFormation>(null, roomType, uid);
    }
    */

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        if (args[0] != null)
            PVPPlayers = (List<pb.Information>)args[0];
        roomType = (RoomType)args[1];

        PVPPlayers.Sort((a, b) => b.MatchTime.CompareTo(a.MatchTime));

        GetView().txt_Title.text = roomType == CFG_Room.RoomType.PVP_One ? "即时挑战模式" : "回合挑战模式";
        GetView().infiniteScroll.SetListItemCount(PVPPlayers.Count);
        GetView().infiniteScroll.RefreshAllShownItem();
        OnChallengeShow();
    }

    private void OnChallengeShow()
    {
        GetView().txt_ChallangeNum.text = "今日可以挑战次数：" + RoleM.Instance.GetPVPMatchTimes((int)roomType);
        CFG_XMLConfig.XML_Config.XML_Item_PVP.PVPseason endPVPSeason = null;
        if (roomType == CFG_Room.RoomType.PVP_One)
            endPVPSeason = RoleM.Instance.GetPVPOneSeason();

        var baseInfo = roomType == CFG_Room.RoomType.PVP_One ? CFG_XMLConfig.xmlConfig.PVPInfo : CFG_XMLConfig.xmlConfig.GVGInfo;
        var isOpne = (endPVPSeason != null || roomType == CFG_Room.RoomType.PVP_More) && GameUtil.Instance.CheakCurTime(baseInfo.dailyBegin, baseInfo.dailyEnd);
        //GetView().btn_add.gameObject.SetActive(isOpne && RoleM.Instance.GetCanNumBuyPVP((int)roomType) > 0);
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
        JEvent.defaultEvent.Unregister(this);
    }

    public override void OnDestroy()
    {

    }
}