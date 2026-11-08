using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewCom_GVGBattle : MonoBehaviour,IUIView
{
	public Transform root_Items;
	public RootItem[] rootItems;

	public class RootItem 
	{
		public GameObject go;
		public Text txt_Name;
		public Text txt_Num;
	}

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		rootItems = new RootItem[6];

		for (int i = 0; i < 6; i++)
        {
			rootItems[i] = new RootItem {
				go = root_Items.GetChild(i).gameObject,
				txt_Name = root_Items.GetChild(i).Find("txt_name").GetComponent<Text>(),
				txt_Num = root_Items.GetChild(i).Find("txt_Num").GetComponent<Text>()
			};
		}
	}
	//public void OnDestroy()
	//{
	//	 UITool.RemoveUI<UICom_GVGBattle>();
	//}
}