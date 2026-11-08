using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewEquipMaster : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Text txt_Title;
	public Text txt_Des;
	public Transform root_Items;
	public Transform Item;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIEquipMaster>();
	}
}