using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPVPBeginTime : MonoBehaviour,IUIView
{
	public Text txt_Time;
	public Text txt_tip;
	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPVPBeginTime>();
	}
}