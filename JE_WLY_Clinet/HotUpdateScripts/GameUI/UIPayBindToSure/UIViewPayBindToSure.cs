using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPayBindToSure : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Button btn_Sure;
	public Button btn_Change;
	public Text txt_Info;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPayBindToSure>();
	}
}