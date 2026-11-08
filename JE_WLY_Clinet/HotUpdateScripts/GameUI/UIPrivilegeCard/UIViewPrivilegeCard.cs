using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPrivilegeCard : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Transform root_Base;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPrivilegeCard>();
	}
}