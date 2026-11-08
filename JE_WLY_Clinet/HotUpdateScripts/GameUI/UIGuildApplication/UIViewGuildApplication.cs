using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewGuildApplication : MonoBehaviour,IUIView
{
    // 战力排序
    public Button btn_SortCombatPower;
    // 等级排序
    public Button btn_SortLevel;

    public LoopListView2 loopList_ApplicationList;
    public GameObject noApply1;
    public GameObject noApply2;

    public Button Btn_RefuseAll;
    public Button Btn_AgreeAll;

    public Button btn_RootPopup;
    // 观察
    public Button btn_Observe;

    public Button Btn_Close;

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
		 UITool.RemoveUI<UIGuildApplication>();
	}
}