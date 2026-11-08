using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewMail :MonoBehaviour, IUIView
{
	public Button btn_refresh;
	public Button btn_DelAll;
	public Button btn_AllGet;
	public Button btn_Close;
	public LoopListView2 infiniteScroll;

	public GameObject tip_NoMail;
	public Text txt_MailTip;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
    private void OnDestroy()
    {
		UITool.RemoveUI<UIMail>();
    }
}