using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewFightAward : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public LoopGridView infiniteScroll;


    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFightAward>();
	}
}