using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewTitleActivity : MonoBehaviour,IUIView
{
    public Button btn_Close;
    public Text txt_Title;
    public LoopListView2 loopView;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UITitleActivity>();
	}
}