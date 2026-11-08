using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewItemUesTips : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Button btn_Use;
	public Text txt_Use;
	public Button btn_Item;
	public Image img_Icon;
	public Image img_BgIcon;
	public Text txt_Num;

	public Text txt_Name;
	public Text txt_Time;

	public Text txt_btn;
	public GameObject obj_Up;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIItemUesTips>();
	}
}