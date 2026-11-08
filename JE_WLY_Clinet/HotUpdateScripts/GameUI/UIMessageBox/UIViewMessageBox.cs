using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewMessageBox :MonoBehaviour, IUIView
{
	public Button Btn_No;
	public Button Btn_Yes;
	public Text Text_des;
	public Text Text_No;
	public Text Text_Yes;

	public Text txt_togTip;
	public Toggle toggle_tip;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
    private void OnDestroy()
    {
		UITool.RemoveUI<UIMessageBox>();
    }

}