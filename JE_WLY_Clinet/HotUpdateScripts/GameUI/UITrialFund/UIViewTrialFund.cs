using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using System;
using JEngine.Core;
using UnityEngine.PlayerLoop;
using SuperScrollView;

public class UIViewTrialFund : MonoBehaviour,IUIView
{
	public LoopListView2 loopListView;
    public Button btn_OneClickGet;
	public Button btn_PayUnlock;
	public Button btn_Close;
	public Text txt_AwardTip;


    public void Awake()
    {
		
    }

    public void Init(GameObject go)
	{
        
    }

    public void OnDestroy()
	{
		 UITool.RemoveUI<UITrialFund>();
	}
}