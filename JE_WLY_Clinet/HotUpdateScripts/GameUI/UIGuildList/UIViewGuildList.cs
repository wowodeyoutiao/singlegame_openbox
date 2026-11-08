using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewGuildList : MonoBehaviour,IUIView
{
    public InputField input_SearchGuild;
    public Button btn_SearchGuild;
    public Button btn_RefreshGuildList;
    public Button btn_Help;

    public Button btn_SortGuildList;
    public LoopListView2 loopList_Guild;

    public Button btn_CreateGuild;
    public Button btn_ApplyAllGuild;
    public Button btn_QuitGuild;

    public Button btn_RootPopup;
    // 观察
    public Button btn_Observe;

    public RectTransform rect_PopupButton;
    public RectTransform rect_ButtonRoot;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuildList>();
	}
}