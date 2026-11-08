using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Helper;
using HotUpdateScripts;
using System;
using SuperScrollView;

public class UIViewPartnetRecastDisassembling : MonoBehaviour,IUIView
{
	public Text txt_Title;
	public Button btn_Help;

	public Text txt_Name;
	public Image img_IconBG;
	public Image img_Icon;
	public Image img_Job;
	public Text txt_Level;

    public Text txt_Desc;
	public LoopListView2 loopList_Item;
	public Transform root_Toggle;
    public Toggle toggle_RecastFree;
	public Text txt_FreeRecastDesc;
    public Toggle toggle_RecastPay;
	public Text txt_PayRecastDesc;
	public Image img_PayRecastConsume;
	public Text txt_PayRecastConsume;
    public Button btn_Recast;
    public Button btn_Disassembling;
    public Button btn_Close;
	
	public void Awake()
    {

    }

	public void Init(GameObject go)
	{

    }


    public void OnDestroy()
	{
		 UITool.RemoveUI<UIPartnetRecastDisassembling>();
	}
}