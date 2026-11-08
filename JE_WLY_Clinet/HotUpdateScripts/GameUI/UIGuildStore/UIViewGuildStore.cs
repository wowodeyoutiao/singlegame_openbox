using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewGuildStore : MonoBehaviour,IUIView
{
    public Image img_Item;
    public Text txt_ItemCount;
    public Button btn_AddItem;
    public Button btn_Help;
    public LoopGridView loopGrid_Store;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuildStore>();
	}
}