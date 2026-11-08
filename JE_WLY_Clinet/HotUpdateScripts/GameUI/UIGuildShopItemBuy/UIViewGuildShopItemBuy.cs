using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGuildShopItemBuy : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Button btn_add;
	public Button btn_jian;
	public InputField input_Buy;
	public Text txt_tip;
	public Button btn_buy;

	public Image img_Consume;
	public Text txt_Consume;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuildShopItemBuy>();
	}
}