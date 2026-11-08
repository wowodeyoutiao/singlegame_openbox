using JEngine.Helper;

using MetaFramework.UI;

using Spine.Unity;

using SuperScrollView;

using System.Collections;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UIViewStarDetails : MonoBehaviour,IUIView
{
	public Text txt_Name;
	public Transform root_Role;
	public LoopListView2 looplist_Items;
	public Text txt_NoSkill;
	public Button btn_Close;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}

    public void OnDestroy()
	{
		 UITool.RemoveUI<UIStarDetails>();
	}
}