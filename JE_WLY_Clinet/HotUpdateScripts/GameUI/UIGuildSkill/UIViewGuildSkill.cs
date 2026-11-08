using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewGuildSkill : MonoBehaviour,IUIView
{
    public Text txt_Title;
	public Image img_Item;
    public Text txt_ItemCount;
	public Button btn_Help;
	public LoopListView2 loopList_Skill;
	public Text txt_TimeTip;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuildSkill>();
	}
}