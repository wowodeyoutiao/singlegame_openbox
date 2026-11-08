using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewBuild_Skill : MonoBehaviour,IUIView
{
	public Text txt_tip;
	public GameObject obj_NoSkill;

    public Image img_topRes1;
	public Text txt_topRes1;
	public Button btn_topRes1Add;

	public Image img_topRes2;
	public Text txt_topRes2;
	public Button btn_topRes2Add;

	public Image img_topRes3;
	public Text txt_topRes3;
	public Button btn_topRes3Add;

	public Image img_topRes4;
	public Text txt_topRes4;
	public Button btn_topRes4Add;

	public GameObject obj_Res1;
	public GameObject obj_Res2;
	public GameObject obj_Res3;
	public GameObject obj_Res4;

	public Button btn_Help;

	public Transform root_skills;
	public Transform tip_equipTip;

	//Advance

	public GameObject AdvanceUI;

	public Image img_ItemBgAdvance;
	public Image img_ItemAdvance;
	public Text txt_advance_advNum;
	public Text txt_advance_Lv;
	public Text txt_advSkillName;

	public Text txt_AdvanceNowLv;
	public Text txt_AdvanceNowDes;

	public Text txt_AdvanceNextLv;
	public Text txt_AdvanceNextDes;

	public Text txt_AdvanCoin1Num;
	public Text txt_AdvanCoin2Num;
	public Image img_AdvanCoin1Num;
	public Image img_AdvanCoin2Num;

	public Button btn_Advance;
	public Text txt_AdvanTip;

	//UpLvUI
	public GameObject UpLvUI;

	public Image img_upLv_ItemBg;
	public Image img_upLv_Item;
	public Text txt_upLv_advanLv;
	public Text txt_upLv_Lv;
	public Text txt_upLv_skillName;

	public Text txt_upLv_NowLv;
	public Text txt_upLv_NowDes;

	public RectTransform rect_upLv_NextDes;
	public RectTransform rect_upLv_NowDes;

	public Text txt_upLv_NextLv;
	public Text txt_upLv_NextDes;

	public Text txt_upLv_Coin1Num;
	public Text txt_upLv_Coin2Num;
	public Image img_upLv_Coin1Num;
	public Image img_upLv_Coin2Num;

	public Button btn_UpLv;
	public Text txtj_UpLvTip;

	//ChangeSkillUI
	public GameObject ChangeSkillUI;
	public Transform root_SkillSolt;
	public Transform obj_dragSolt;


	public struct SkillItem
	{
		public TouchEventCom _touch;
		public Image img_bg;
		public Image img_icon;
		public Text txt_Adv;
		public Text txt_Lv;
		public GameObject obj_equip;
		public GameObject obj_Select;
		public GameObject dot;
    }

	public SkillItem[] SkillItems;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		rect_upLv_NowDes = txt_upLv_NowDes.transform.parent.GetComponent<RectTransform>();
		rect_upLv_NextDes = txt_upLv_NextDes.transform.parent.GetComponent<RectTransform>();

        SkillItems = new SkillItem[root_skills.childCount];

		for (int i = 0; i < root_skills.childCount; i++)
		{
			var root = root_skills.GetChild(i);
			SkillItems[i] = new SkillItem { 
				_touch = root.GetComponent<TouchEventCom>(),
                img_bg = root.GetComponent<Image>(),
                img_icon = root.Find("icon").GetComponent<Image>(),
                txt_Adv = root.Find("txt_Advance").GetComponent<Text>(),
                txt_Lv = root.Find("txt_lv").GetComponent<Text>(),
                obj_equip = root.Find("txt_equip").gameObject,
                obj_Select = root.Find("select").gameObject,
                dot = root.Find("dot").gameObject
            };

        }
    }
	public void OnDestroy()
	{
		 //UITool.RemoveUI<UIBuild_Skill>();
	}
}