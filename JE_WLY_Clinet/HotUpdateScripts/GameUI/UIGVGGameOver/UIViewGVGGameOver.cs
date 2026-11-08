using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGVGGameOver : MonoBehaviour,IUIView
{
	public Transform root_Rank;
	public Button btn_back;
	public Text txt_Time;
	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGVGGameOver>();
	}
}