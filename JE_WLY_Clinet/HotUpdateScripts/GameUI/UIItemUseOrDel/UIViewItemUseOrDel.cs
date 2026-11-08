using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewItemUseOrDel :MonoBehaviour, IUIView
{
	public Text txt_name;
	public Text txt_LevelLimit;
	public Text txt_OccupyLimit;
	public Text txt_GenderLimit;
	public Text txt_Des;
	public Text txt_attribute0;
	public Text txt_attribute1;
	public Text txt_attributeBast0;
    public Text txt_attributeBast1;
    public GameObject line_attribute;
	public GameObject line_attributeBest;
	public RectTransform rect_attribute;
	public RectTransform rect_base;

    public GameObject obj_Des;
    public GameObject Obj_equip;
	public GameObject obj_equipUp;
    public GameObject obj_equipDown;

    public Image Img_icon;
	public Image Img_Icon_bg;
	public Image Img_Nature;

	public Button Btn_Use;
	public Text txt_btn_Use;
	public Button Btn_Del;
	public Button Btn_Close;
	public Text txt_delBtn;

	public Button Btn_Syn;
	public Button Btn_Retain;

    //使用多个
    public GameObject UIUseNumChoice;
	public InputField input_Used;
	public Button btn_UsedMore;
	public Button btn_Jian;
	public Button btn_Add;
	public Button btn_Max;
	public Button btn_closeUsedMore;

	//已穿戴
	public Text now_txt_name;
	public Text now_txt_LevelLimit;
	public Text now_txt_OccupyLimit;
	public Text now_txt_GenderLimit;
	public Text now_txt_Des;
	public Text now_txt_attribute0;
	public Text now_txt_attribute1;
    public Text now_txt_attributeBest0;
    public Text now_txt_attributeBest1;
    public GameObject now_line_attribute;
	public GameObject now_line_attributeBest;
    public GameObject now_obj_Des;

    public RectTransform now_rect_attribute;
	public RectTransform now_rect_base;
	public Image now_Img_icon;
	public Image now_Img_Icon_bg;
	public Image now_Img_Nature;

	// 极品特效
	public Transform baseBestEffect;
    public Transform compBestEffect;

	// base
	public Transform baseTrans;

    public void Awake()
    {

	}
	public void Init(GameObject go)
	{
		//rect_attribute = txt_attribute.transform.GetComponent<RectTransform>();
	}

    private void OnDestroy()
    {
		UITool.RemoveUI<UIItemUseOrDel>();
    }
}