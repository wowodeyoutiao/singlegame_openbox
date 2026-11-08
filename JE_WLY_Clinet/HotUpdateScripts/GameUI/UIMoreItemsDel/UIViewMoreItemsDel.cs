using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewMoreItemsDel : MonoBehaviour, IUIView
{
	public Dropdown dropdown_Qulite;
	public Dropdown dropdown_Level;
	public Dropdown dropdown_Job;
	public Dropdown dropdown_Sex;
	public Button btn_Close;
	public Button btn_Refresh;
	public Button btn_Del;

	public LoopGridView infiniteScroll;

	public GameObject Obj_tip;

	public Toggle tog_Select;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
    private void OnDestroy()
    {
		UITool.RemoveUI<UIMoreItemsDel>();
    }

}