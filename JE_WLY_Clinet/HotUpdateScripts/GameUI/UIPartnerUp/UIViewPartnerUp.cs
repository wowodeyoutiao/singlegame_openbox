using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewPartnerUp : MonoBehaviour,IUIView
{
    public Dropdown dropdown_Qulite;
    public Button btn_up;
    public Button btn_Close;

    public LoopGridView infiniteScroll;

    public GameObject Obj_tip;

    public Text txt_Btn;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPartnerUp>();
	}
}