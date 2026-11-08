using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewTitleLuckChoice : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Transform root_Awards;

	public Button btn_Select;
	public Transform txt_tip;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UITitleLuckChoice>();
	}
}