using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewMailDetail :MonoBehaviour, IUIView
{

	public Button btn_Close;
	public Button btn_DelAndGet;
	public Text txt_BtnGetAndDel;

	public LoopListView2 infiniteScroll_Award;
	public Text txt_NoAward;
	public Text txt_des;
	public Text txt_Title;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}

    private void OnDestroy()
    {
		UITool.RemoveUI<UIMailDetail>();
    }

}