using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using Spine.Unity;

public class UIViewBuild_Equip : MonoBehaviour,IUIView
{

	public Transform root_res;

	public Button btn_Help;

	public Transform root_allEquip;
	public Transform tip_equipTip;


	// 强化
	public GameObject ReinforceUI;
	public Button btn_Reinforce;

	public Image img_ItemBgReinforce;
	public Image img_ItemReinforce;
	public Text txt_StarReinforce;
	public Text txt_ReinforceLv;
	public Text txt_ReinforceItemName;

	public Text txt_ReinforceNowLv;
	public Text txt_ReinforceNowDes;

	public Text txt_ReinforceNextLv;
	public Text txt_ReinforceNextDes;

	public Image img_ReinforceGold;
	public Text txt_ReinforceGoldNum;
	public Image img_ReinforceQHS;
	public Text txt_ReinforceQHSNum;

	public Button btn_AddReinforceOne;
	public Button btn_AddReinforceTen;
	public Button btn_AddReinforceEnd;
	[HideInInspector] public GameObject Obj_AddReinforceEnd;

	public Text txt_ReinforceTip;

	//升星
	public GameObject StarUI;

	public Button btn_Star;

	public Image img_Star_ItemBg;
	public Image img_Star_Item;
	public Text txt_Star_StarLv;
	public Text txt_Star_ReinforceLv;
	public Text txt_Star_ItemName;

	public Text txt_Star_Des;

	public Text txt_Star_NowLv;
	public Text txt_Star_NowDes;

	public Text txt_Star_NextLv;
	public Text txt_Star_NextDes;

	public Image img_Star_GoldNum;
	public Text txt_Star_GoldNum;
	public Image img_Star_SXSNum;
	public Text txt_Star_SXSNum;

	public Button btn_AddStarOne;
	public Button btn_AddStarTop;
	public Dropdown dropDown_Star;

	public GameObject obj_StarTip;
	public GameObject obj_EndStar;
	public Button btn_EndStar;

	//洗练
	public GameObject RefineUI;

	public Text txt_RefineEquipSlotTitle;
	public Button btn_RefineAttributePreview;

	public Transform root_RefineItemPreview;

	public Transform root_RefineNowRoot;
	[HideInInspector] public Toggle[] tog_RefineNows;
    public Text txt_RefineResult;
	public GameObject obj_RefineResultMax;
	public GameObject obj_ComsumeRoot;
	public Text txt_RefineConsumeNum;
	public Image img_RefineConsumeIcon;

	public Toggle tog_RefineRandom;
	public Toggle tog_RefineFixed;

	public Button btn_Refine;

	public GameObject obj_RefineResuleButtons;
	public Button btn_RefineResuleDontSave;
    public Button btn_RefineResuleSave;

	public SkeletonGraphic spine_Hammer;
	public ParticleSystem ps_PreviewResult;
	public SkeletonGraphic spine_SelectResult;

    //镶嵌 魂石
    public GameObject InsetUI;
	public Transform Inset_Item;
	public Button btn_InsetInfo;
	public Button btn_InsetNatrue;
	public Button btn_AutoInset;
	public Text txt_InsetTip;
	public Transform root_Inset;

	// 升星符
	public GameObject noTalismanGo;
	public Button addTalismanBtn;

	public Transform talismanTrans;
    public Button addTalismanBtn1;
    public Text talismanTimeText;

    public struct UsedItem
	{
		public Image img;
		public Text txt;
		public Button btn;
		public GameObject obj;
	}

    [HideInInspector] public UsedItem[] usedItems;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		Obj_AddReinforceEnd = btn_AddReinforceEnd.transform.parent.gameObject;

		usedItems = new UsedItem[root_res.childCount];
		for (int i = 0; i < root_res.childCount; i++)
		{
			usedItems[i] = new UsedItem {
                obj = root_res.GetChild(i).gameObject,
                btn = root_res.GetChild(i).Find("add").GetComponent<Button>(),
                img = root_res.GetChild(i).Find("tip").GetComponent<Image>(),
				txt = root_res.GetChild(i).Find("txt").GetComponent<Text>()
            };

        }

		tog_RefineNows = new Toggle[4];
        for (int i = 0; i < 4; i++)
		{
			tog_RefineNows[i] = root_RefineNowRoot.Find($"Item_{i+1}/toggle").GetComponent<Toggle>();
        }

    }
	public void OnDestroy()
	{
		 //UITool.RemoveUI<UIBuild_Equip>();
	}
}