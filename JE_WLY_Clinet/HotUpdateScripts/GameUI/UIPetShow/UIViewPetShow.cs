using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPetShow : MonoBehaviour,IUIView
{
	public Text txt_name;
	public Text txt_Des;
	public Button Btn_Close;

	public Image Img_icon;
	public Image Img_Icon_bg;

	public Text txt_attribute;
    public Text txt_attribute1;

    public RectTransform rect_Center;


    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPetShow>();
	}
}