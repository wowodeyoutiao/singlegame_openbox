using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewEquipOff : MonoBehaviour,IUIView
{
	public Text txt_name;
	public Text txt_LevelLimit;
	public Text txt_OccupyLimit;
	public Text txt_GenderLimit;
	public Text txt_Des;
	public Text txt_attribute0;
	public Text txt_attribute1;

	public Image Img_icon;
	public Image Img_Icon_bg;

	public Button Btn_Off;
	public Button Btn_Close;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIEquipOff>();
	}
}