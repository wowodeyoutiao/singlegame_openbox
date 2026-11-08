using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewFightRecord : MonoBehaviour,IUIView
{
    public Button btn_Close;
    public Button btn_add;
    public Text txt_ChallangeNum;
    public LoopListView2 infiniteScroll;
    public Text txt_Title;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFightRecord>();
	}
}