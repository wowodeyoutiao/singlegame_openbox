using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGuardIslandBox : MonoBehaviour,IUIView
{
	public Button btn_close;
	public Button Btn_Yes;
	public Transform root_Item;

	public InfoItem[] infoItems;
    public struct InfoItem 
	{
		public Image img_icon;
		public Text txt;
		public Text txt_des;
		public Toggle toggle;
		public GameObject go;
    }
    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
        infoItems = new InfoItem[root_Item.childCount];

		for (int i = 0; i < root_Item.childCount; i++)
		{
			infoItems[i] = new InfoItem
			{
				go = root_Item.GetChild(i).gameObject,
                img_icon = root_Item.GetChild(i).Find("img_icom").GetComponent<Image>(),
                txt = root_Item.GetChild(i).Find("Text_des").GetComponent<Text>(),
				txt_des = root_Item.GetChild(i).Find("Toggle/Label").GetComponent<Text>(),
                toggle = root_Item.GetChild(i).Find("Toggle").GetComponent<Toggle>()
            };
		}

    }
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuardIslandBox>();
	}
}