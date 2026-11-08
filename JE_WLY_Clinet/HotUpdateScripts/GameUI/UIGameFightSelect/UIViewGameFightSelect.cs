using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameFightSelect : MonoBehaviour,IUIView
{
	public Toggle tog_AttRole;
	public Button btn_Close;

	public GameObject obj_Player;
	public GameObject obj_Monster;

	public Transform root_lay0;
	public Transform root_lay1;

	public AttactItem[] AttactOnes;
	public AttactItem[] AttactTwos;

	public class AttactItem 
	{
		public GameObject obj;
		public Image img_Head;
		public Slider img_Hp;
		public Text txt_Hp;
		public Text txt_Name;
		public Text txt_Lv;
	}

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		AttactOnes = new AttactItem[5];
		AttactTwos = new AttactItem[5];

		for (int i = 0; i < 5; i++)
        {
			Transform root = root_lay0.GetChild(i);
			AttactOnes[i] = new AttactItem
			{
				obj = root.gameObject,
				img_Head = root.Find("head_bg/Img").GetComponent<Image>(),
				img_Hp = root.Find("Slider").GetComponent<Slider>(),
				txt_Hp = root.Find("Slider/txt").GetComponent<Text>(),
				txt_Name = root.Find("txt_Name").GetComponent<Text>(),
				txt_Lv = root.Find("img_Lv/Txt_lv").GetComponent<Text>(),
			};

			root = root_lay1.GetChild(i);
			AttactTwos[i] = new AttactItem
			{
				obj = root.gameObject,
				img_Head = root.Find("head_bg/Img").GetComponent<Image>(),
				img_Hp = root.Find("Slider").GetComponent<Slider>(),
				txt_Hp = root.Find("Slider/txt").GetComponent<Text>(),
				txt_Name = root.Find("txt_Name").GetComponent<Text>(),
				txt_Lv = root.Find("img_Lv/Txt_lv").GetComponent<Text>(),
			};
		}
	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGameFightSelect>();
	}
}