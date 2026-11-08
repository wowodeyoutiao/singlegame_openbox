using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewDayDoAward : MonoBehaviour,IUIView
{
	public Button Btn_No;
	public Button Btn_Yes;
	public Transform root_award;
	public Text txt_num;
	public Image img_icon;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIDayDoAward>();
	}
}