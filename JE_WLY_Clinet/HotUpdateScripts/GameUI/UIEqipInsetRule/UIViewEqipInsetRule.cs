using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewEqipInsetRule : MonoBehaviour,IUIView
{
	public Text txt_Title;
	public Button btn_Close;
	public Transform root_Content;
	public GameObject obj_Item;
	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIEqipInsetRule>();
	}
}