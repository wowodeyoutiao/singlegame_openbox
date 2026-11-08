using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPhoneLogin : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Button btn_Sure;
	public InputField Input_Phone;
	public Toggle tog_agree;
	public Button btn_agree;
	public Button btn_yinsi;
	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPhoneLogin>();
	}
}