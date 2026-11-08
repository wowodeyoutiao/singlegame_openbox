using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewCreateGuild : MonoBehaviour,IUIView
{
	public Button Btn_Close;
	public Button Btn_Yes;
	public InputField Input_Name;
	public Image img_icon;
	public Text txt_icon;

	public void Awake()
    {
		
    }

	public void Init(GameObject go)
	{
        
    }

	public void OnDestroy()
	{
		 UITool.RemoveUI<UICreateGuild>();
	}
}