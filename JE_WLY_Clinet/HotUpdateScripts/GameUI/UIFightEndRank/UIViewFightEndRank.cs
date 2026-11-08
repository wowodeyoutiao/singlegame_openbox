using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewFightEndRank : MonoBehaviour,IUIView
{
	public LoopListView2 looplist_Rank;
	public Transform root_MyRank;
	public Text txt_Desc;
	public Button btn_Close;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFightEndRank>();
	}
}