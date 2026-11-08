using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameDayDoPro : MonoBehaviour,IUIView
{
	public Button btn_Show;
	public Button btn_Close;

	public GameObject obj_UI;
	public Text txt_Time;

	public Transform root_Tips;
	public GameObject obj_Pro;

	public Text txt_desPro;
	public Slider slider_Pro;

	public GameObject obj_pf;
	public Text txt_pf1;
	public Text txt_pf2;

	public Text txt_Pro;

	public ShowTip[] showTips;
	public class ShowTip 
	{
		public Text txt_Des;
		public Text txt_Num;
	}

	public void Awake()
    {
		
	}
	public void Init(GameObject go)
	{
		showTips = new ShowTip[root_Tips.childCount];

		for (int i = 0; i < root_Tips.childCount; i++)
		{
			showTips[i] = new ShowTip
			{
				txt_Des = root_Tips.GetChild(i).GetComponent<Text>(),
				txt_Num = root_Tips.GetChild(i).Find("txt_pro").GetComponent<Text>(),
			};

		}
	}
	//public void OnDestroy()
	//{
	//	UITool.RemoveUI<UIGameDayDoPro>();
	//}
}