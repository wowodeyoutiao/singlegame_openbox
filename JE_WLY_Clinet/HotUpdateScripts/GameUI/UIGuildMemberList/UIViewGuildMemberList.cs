using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewGuildMemberList : MonoBehaviour,IUIView
{

    public Button btn_Help;
	// 战力排序
    public Button btn_SortCombatPower;
    // 等级排序
    public Button btn_SortLevel;

	public LoopListView2 loopList_MemberList;

	public Button btn_QuitGuild;
	public Button btn_GuildApplication;

	public Button btn_RootPopup;
	// 观察
	public Button btn_Observe;
	// 委任副会长
	public Button btn_AppointVicePresident;
    // 解除副会长
    public Button btn_DismissVicePresident;
	// 请出行会
	public Button btn_PleaseGoOut;

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
		 UITool.RemoveUI<UIGuildMemberList>();
	}
}