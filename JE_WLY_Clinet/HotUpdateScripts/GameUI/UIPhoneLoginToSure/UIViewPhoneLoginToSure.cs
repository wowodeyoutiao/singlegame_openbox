using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPhoneLoginToSure : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Transform root_Time;
	public InputField Input_PhoneCode;
	public Button btn_newCode;
	public Text txt_time;
	public Text txt_Num;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPhoneLoginToSure>();
	}
}