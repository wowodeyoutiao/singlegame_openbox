using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewRefinePreview : MonoBehaviour,IUIView
{
	public Text txt_Title;

	public Transform root_buttons;
	[HideInInspector] public Button[] buttons;

	public Text txt_MinValue;
	public Text txt_MaxValue;

	public Button btn_Close;

	public void Awake()
    {
        buttons = root_buttons.GetComponentsInChildren<Button>(true);
    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIRefinePreview>();
	}
}