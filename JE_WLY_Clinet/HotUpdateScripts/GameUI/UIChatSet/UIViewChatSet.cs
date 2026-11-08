using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Helper;
using HotUpdateScripts;
using System;

public class UIViewChatSet : MonoBehaviour,IUIView
{
	public Toggle toggle_WorldVoice;
	public Toggle toggle_GuildVoice;
	public Toggle toggle_WorldMSG;
	public Toggle toggle_GuildMSG;
    public Button Btn_Close;
	
	public void Awake()
    {

    }

	public void Init(GameObject go)
	{

    }


    public void OnDestroy()
	{
		 UITool.RemoveUI<UIChatSet>();
	}
}