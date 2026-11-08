using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewSynthesisSelect : MonoBehaviour, IUIView
{
	public Text txt_Tiltle;
	public Dropdown dropdown_Qulite;
	public Dropdown dropdown_Level;
	public Dropdown dropdown_Job;
	public Dropdown dropdown_Sex;
	public Button btn_Close;
	public Button btn_Refresh;
	public Button btn_AutoSyn;
	public Button btn_Syn;

	public LoopGridView infiniteScroll;

	public GameObject Obj_tip;
	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
    private void OnDestroy()
    {
		UITool.RemoveUI<UISynthesisSelect>();
    }

}