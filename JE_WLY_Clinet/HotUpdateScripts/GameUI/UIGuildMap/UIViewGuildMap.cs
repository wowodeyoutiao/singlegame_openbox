using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;
using UnityEngine.PlayerLoop;
using JEngine.Helper;

public class UIViewGuildMap : MonoBehaviour,IUIView
{
    public Image img_Item;
    public Text txt_ItemCount;
    public Button btn_Help;
    public Button btn_AddItem;
    public Button btn_GoMap;
    public Button btn_Gray;
    public Text txt_OpenTime;
    public Text txt_MapDesc;
    public Text txt_PlayingRule;
    public LoopListView2 loopList_Awards;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuildMap>();
	}
}