using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts;
using JEngine.Helper;
using System;
using HotUpdateScripts.ProtoM;

public class UIChatSet : UIBase<UIModelChatSet, UIViewChatSet>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(()=> {
            SaveSetting();
            ChatM.Instance.LoadHallChat();
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIChatSet>();
        });
    }

    public override void OnEnter(params object[] args)
    {
        LoadSetting();
    }

    void LoadSetting()
    {
        GetView().toggle_WorldVoice.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.ChatOption.WorldVoice.ToString(), 1));
        GetView().toggle_GuildVoice.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.ChatOption.GuildVoice.ToString(), 1));
        GetView().toggle_WorldMSG.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.ChatOption.WorldMSG.ToString(), 1));
        GetView().toggle_GuildMSG.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.ChatOption.GuildMSG.ToString(), 1));
    }

    void SaveSetting()
    {
        GameUtil.Instance.SetPlayerPrefs(GameConst.ChatOption.WorldVoice.ToString(), GetView().toggle_WorldVoice.isOn ? 1 : 0);
        GameUtil.Instance.SetPlayerPrefs(GameConst.ChatOption.GuildVoice.ToString(), GetView().toggle_GuildVoice.isOn ? 1 : 0);
        GameUtil.Instance.SetPlayerPrefs(GameConst.ChatOption.WorldMSG.ToString(), GetView().toggle_WorldMSG.isOn ? 1 : 0);
        GameUtil.Instance.SetPlayerPrefs(GameConst.ChatOption.GuildMSG.ToString(), GetView().toggle_GuildMSG.isOn ? 1 : 0);
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}