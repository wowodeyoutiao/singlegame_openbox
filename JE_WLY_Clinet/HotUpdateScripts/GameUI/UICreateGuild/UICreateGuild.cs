using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Event;
using JEngine.Helper;
using System;
using JEngine.MyNet;
using pb;
using HotUpdateScripts.ProtoM;

public class UICreateGuild : UIBase<UIModelCreateGuild, UIViewCreateGuild>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UICreateGuild>();
        });

        GetView().Btn_Yes.onClick.AddListener(() => {
            OnClickCreateButton();
        });
    }

    private void OnClickCreateButton()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (string.IsNullOrEmpty(GetView().Input_Name.text))
        {
            UITool.ShowTips("名字不可为空");
            return;
        }

        if (!GameUtil.Instance.ComCoinToCheck(CFG_XMLConfig.xmlConfig.Guild.creatcost.itemid, CFG_XMLConfig.xmlConfig.Guild.creatcost.neednum, () => UITool.Close<UICreateGuild>(), true))
        {
            return;
        }

        HttpC.SendGuildCreateRequest(GetView().Input_Name.text, "");

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "社交",
            component = "行会",
            position = "行会列表",
            event_click = "创建行会",
        });
        //TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_GuildCreate, MType.Request, TCPClientMgr.ProtoSerialize(new GuildCreateRequest { Name = GetView().Input_Name.text, Description = "" }));
        //GetView().Input_Name.text = "";
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);

        GetView().Input_Name.text = "";

        GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.Guild.creatcost.itemid, GetView().img_icon);
        GetView().txt_icon.text = GameUtil.Instance.GetCheakNumStr(CFG_XMLConfig.xmlConfig.Guild.creatcost.itemid, CFG_XMLConfig.xmlConfig.Guild.creatcost.neednum);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_Refresh(GameEvent.Sync_GuildCreate _e)
    {
        UITool.Close<UICreateGuild>();
    }

    public override void OnDestroy()
    {
        
    }
}