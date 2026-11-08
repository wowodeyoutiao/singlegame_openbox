using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPayBind : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Button btn_Sure;

	public InputField Input_Acc;
	public InputField Input_Name;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPayBind>();
	}
}