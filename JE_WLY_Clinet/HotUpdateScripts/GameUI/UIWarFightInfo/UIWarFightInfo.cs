using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using System;
using System.Collections;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UIWarFightInfo : UIBase<UIModelWarFightInfo, UIViewWarFightInfo>
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
            UITool.Close<UIWarFightInfo>();
        });
        GetView().btn_Play.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (!FightPlayM.Instance.CheakWarFightOpenTime())
            {
                UITool.ShowTips("当前时间尚未开放");
                return;
            }
            TcpC.GetGameJoin(_RoomID, 0);

            // 数数上报点击事件
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "群雄战场",
                position = $"{CFG_Room.GetById(_RoomID).Name}",
                event_click = $"进入挑战",
            });
        });

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
        GetView().txt_Tip.text = $"提示：{Str}\n{xml.startTimes[0]}~{xml.endTimes[0]} {xml.startTimes[1]}~{xml.endTimes[1]}开放";
    }

    private int _RoomID;

    private TimeItem RefreshUILoop;
    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        _RoomID = (int)args[0];
        OnUIInfoShow();
        OnPlayerNumShow();

        ClearRefreshUILoop();
        RefreshUILoop = GameHelper.Instance.AddInvoke(new TimeItem(0, 5, () =>
        {
           TcpC.Get_C2S_WarRoomList();
        }, _isLoop: true));
    }

    private void ClearRefreshUILoop()
    {
        if (GameHelper.Instance!= null && RefreshUILoop != null)
        {
            GameHelper.Instance.RemoveInvoke(RefreshUILoop);
            RefreshUILoop = null;
        }
    }

    private void OnUIInfoShow()
    {
        var CFGRoom = CFG_Room.GetById(_RoomID);
        GetView().txt_Title.text = CFGRoom.Name;
        var CFGMonster = CFG_Monster.GetById((int)CFGRoom.Show);
        ResouseHelper.Instance.LoadSpriteByAtlas(CFGMonster.Photo, CFGMonster.Photoname, GetView().img_Monster);
        ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/WarFight", CFGRoom.Name.Contains("地狱")? "bg_qxzc_2": "bg_qxzc_1", GetView().img_MonsterBg);

        Transform itemRoot = GetView().root_Awards;
        for (int i = 0; i < 8; i++)
        {
            var go = itemRoot.Find("Item_" + i);
            go.gameObject.SetActive(i < CFGRoom.Droppreview.Length);

            if (i < CFGRoom.Droppreview.Length)
            {
                int cfgId = CFGRoom.Droppreview[i];
                GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                });
            }
        }
    }

    private void OnPlayerNumShow()
    {
        var roomInfo = FightPlayM.Instance.GetOneWarFightRoomInfo(_RoomID);
        GetView().txt_PlayerNum.text = $"人数：<color={(roomInfo.CurPlayerNum == roomInfo.MaxPlayerNum ? "#ff7979" : "#aff69b")}>{roomInfo.CurPlayerNum}</color>/{roomInfo.MaxPlayerNum}" ;
    }


    [Subscriber(ThreadMode.Main)]
    private void OnRefresh_WarRoomInfo(GameEvent.Refresh_WarRoomInfo _e)
    {
        if (isNullGo)
            return;
        OnPlayerNumShow();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        ClearRefreshUILoop();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}