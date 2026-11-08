using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPet : MonoBehaviour,IUIView
{
	public Button btn_Close;
	public Button btn_PetBook;
	public Button btn_PetGift;
	public Transform root_Pets;
	public Transform root_Btns;

	public Button btn_help;
	public Button btn_help2;

    public TouchEventCom touch_Pet;
	public Text txt_PetName;
	public Image img_Pet;
	public GameObject obj_Petlv;
	public Text txt_Petlv;
	public Image img_PetMark;
	public GameObject obj_Exp;

    public GameObject obj_Time;
	public Text txt_Time;

    public Image img_fly;
    public Image img_Drag;

	public GameObject obj_fly;
	public GameObject obj_redDot;

	public Image img_res;
	public Text txt_res;
	public Button btn_res;

	public Button btn_Yc;

    public struct PetItem 
	{
		public TouchEventCom touchCom;
		public GameObject obj_lock;
		public Text txt_lock;
		public Image img_Icon;
		public GameObject obj_lv;
		public Text txt_lv;
		public Image img_Pet;
		public GameObject obj_NoMove;
		public GameObject obj_Exp;
    }

	public PetItem[] PetItems;

    public void Awake()
    {
        PetItems= new PetItem[root_Pets.childCount];
        for (int i = 0; i < root_Pets.childCount; i++)
		{
			var go = root_Pets.GetChild(i);
            PetItems[i] = new PetItem {
                touchCom = go.GetComponent<TouchEventCom>(),
				obj_lock = go.Find("lock").gameObject,
				txt_lock = go.Find("lock/txt_lock").GetComponent<Text>(),
				img_Icon = go.Find("icon").GetComponent<Image>(),
				obj_lv = go.Find("img_lv").gameObject,
				txt_lv = go.Find("img_lv/txt_lv").GetComponent<Text>(),
                img_Pet = go.Find("img_lv/img_pet").GetComponent<Image>(),
                obj_NoMove = go.Find("lock_Move").gameObject,
                obj_Exp = go.Find("img_exp").gameObject
            };
        }
    }
	public void Init(GameObject go)
	{
	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPet>();
	}
}