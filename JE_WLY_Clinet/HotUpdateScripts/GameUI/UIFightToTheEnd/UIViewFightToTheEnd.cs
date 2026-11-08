using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewFightToTheEnd : MonoBehaviour,IUIView
{
	public Text txt_Title;
    public Button btn_Help;
	public LoopListView2 looplist_Map;
    public Button btn_Rank;
    public Button btn_Fight;
	public Text txt_Tip;
    public Button btn_Close;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFightToTheEnd>();
	}
}