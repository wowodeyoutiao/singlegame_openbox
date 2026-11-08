using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewItemShow : MonoBehaviour,IUIView
{
	public Text txt_name;
	public Text txt_LevelLimit;
	public Text txt_OccupyLimit;
	public Text txt_GenderLimit;
	public Text txt_Des;
	public Text txt_attribute0;
	public Text txt_attribute1;

    public Text txt_Bestattribute0;
    public GameObject obj_Bestattribute;

	public GameObject obj_attribute;
    public GameObject obj_Des;

    public Image Img_icon;
	public Image Img_Icon_bg;
	public Image Img_Nature;

	public Button Btn_PreviewStarInfo;
	public Button Btn_Close;

	public GameObject obj_BtnNode;
	public Button Btn_Yes;

	public Transform bestEffect;


    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIItemShow>();
	}
}