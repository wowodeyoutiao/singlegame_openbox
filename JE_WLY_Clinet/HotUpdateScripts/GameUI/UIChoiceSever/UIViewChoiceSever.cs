using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewChoiceSever : MonoBehaviour,IUIView
{
    public LoopListView2 InfiniteScroll_Info;
    public LoopListView2 InfiniteScroll_Btns;
	public Button btn_Close;

	public Transform root_Btns;
    public void Awake()
    {
		root_Btns = InfiniteScroll_Btns.transform.Find("Viewport/Content");
    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIChoiceSever>();
	}
}