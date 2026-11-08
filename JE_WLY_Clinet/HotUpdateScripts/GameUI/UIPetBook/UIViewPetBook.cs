using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewPetBook : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Dropdown dropdown_Qulite;
	public LoopGridView infiniteScroll;


    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPetBook>();
	}
}