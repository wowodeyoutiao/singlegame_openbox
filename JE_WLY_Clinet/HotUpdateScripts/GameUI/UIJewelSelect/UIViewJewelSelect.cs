using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewJewelSelect : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Button btn_Select;
	public Button btn_Syn;

	public GameObject Obj_tip;
	public LoopGridView infiniteScroll;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIJewelSelect>();
	}
}