using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewEquipSelect : MonoBehaviour,IUIView
{
	public Button btn_Equip;
	public Button btn_OpenBox;
    public Button btn_Close;
	public GameObject obj_NoTip;
	public LoopGridView lgv_Equip;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIEquipSelect>();
	}
}