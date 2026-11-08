using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewRole_Skill : MonoBehaviour,IUIView
{
	public Text txt_title;
	public LoopListView2 infiniteScroll_Skill;
	public GameObject obj_NoSkill;
	public Button btn_Close;
	public Button btn_Build;
	public Text txt_BuildLock;

	public GameObject obj_BuildRed;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIRole_Skill>();
	}
}