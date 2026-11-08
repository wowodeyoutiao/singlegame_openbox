using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewBoxInfoShow : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Transform awardRoot;
	public Transform buttonRoot;
	public Text txt_Tip;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIBoxInfoShow>();
	}
}