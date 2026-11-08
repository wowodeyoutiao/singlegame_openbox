using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewLoginDes : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public GameObject obj_Privacy;
	public GameObject obj_Agree;
	public Text txt_Title;

	public Transform root_Des;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UILoginDes>();
	}
}