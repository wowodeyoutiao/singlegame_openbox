using MetaFramework.UI;

using Spine.Unity;

using System.Collections;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UIViewRole :MonoBehaviour, IUIView
{
	public Transform root_head;
	public Text txt_allPower;
	public Transform root_Team;

    public Button btn_Close;
	public Transform root_SelectBtns;

    public Text txt_RoleName;
    public Text txt_RolePower;
    public Text txt_lv;
    public Transform root_Equips;

    public Transform root_RecastDisassemblingBtn;
    public Button btn_Recast;
    public Button btn_Disassembling;
    public Button btn_PreviewStarInfo;

    public Button btn_PartnerSet;
    public Text txt_PartnerSet;

    public Button btn_Pet;
    public GameObject obj_PetLock;
    public Text txt_petLock;

    public Button btn_AllEquip;

    public Button btn_Attribute;
    public Text txt_attr0;
    public Text txt_attr1;
    public Text txt_attr2;
    public Text txt_attr3;

    public Transform root_Role;
    public Image img_Role;

    public GameObject obj_PetDot;
    public GameObject obj_EqiuipDot;
    public Image img_Arrow;

    public Image img_title;

    public struct SelectBtn 
	{
		public Button _btn;
		public Text _txtOn;
        public Text _txtOff;
		public GameObject _objOn;
        public GameObject _objDot;
        public GameObject _objDotRoot;
        public GameObject _objFunctionOpenNewFlag;
        public GameObject _objMask;
        public Text _txtOpen;
    }

    [HideInInspector] public SelectBtn[] selectBtns;

	public struct TeamItem
    {
        public Image _img_head;

        /// <summary>
        /// 主角头像品质低图，不用于伙伴。
        /// </summary>
        public Image _img_head_quality;
        public Image _img_job;
        public GameObject _obj_Lv;
        public Text _txt_lv;
        public Button _btn;
		public GameObject _obj_lock;
        public Text _txt_Lock;
        public Text txt_Star;
        public GameObject _obj_on;

        /// <summary>
        /// 主角：用于头像框，伙伴：用于品质框。
        /// </summary>
        public Image _img_frame;

        public SkeletonGraphic spine_Unlock;
    }

    [HideInInspector] public TeamItem[] teamItems;

    public void Awake()
	{

	}
	public void Init(GameObject go)
	{
		selectBtns = new SelectBtn[root_SelectBtns.childCount];
		for (int i = 0; i < root_SelectBtns.childCount; i++)
		{
			selectBtns[i] = new SelectBtn { 
				_btn = root_SelectBtns.GetChild(i).GetComponent<Button>(),
                _txtOff = root_SelectBtns.GetChild(i).Find("Label").GetComponent<Text>(),
                _txtOn = root_SelectBtns.GetChild(i).Find("Select/Label").GetComponent<Text>(),
                _objOn = root_SelectBtns.GetChild(i).Find("Select").gameObject,
                _objDot = root_SelectBtns.GetChild(i).Find("redDotRoot/dot").gameObject,
                _objDotRoot = root_SelectBtns.GetChild(i).Find("redDotRoot").gameObject,
                _objFunctionOpenNewFlag = root_SelectBtns.GetChild(i).Find("functionOpenNewFlag").gameObject,
                _objMask = root_SelectBtns.GetChild(i).Find("Mask").gameObject,
                _txtOpen = root_SelectBtns.GetChild(i).Find("Mask/txt").GetComponent<Text>(),
            };
        }

        teamItems = new TeamItem[5];
		for (int i = 1; i <= root_Team.childCount; i++)
		{
			SetTeamItem(i, root_Team.GetChild(i - 1));
            teamItems[i].txt_Star = root_Team.GetChild(i - 1).Find("root_Star").GetComponent<Text>();
        }
        SetTeamItem(0, root_head);
        void SetTeamItem(int i,Transform root)
		{
            teamItems[i] = new TeamItem
            {
                _btn = root.GetComponent<Button>(),
                _img_head = root.Find("img_head").GetComponent<Image>(),
                _img_head_quality = i == 0 ? root.Find("img_head_quality").GetComponent<Image>() : null,
                _img_job = root.Find("img_job").GetComponent<Image>(),
                _obj_Lv = root.Find("img_Lv").gameObject,
                _txt_lv = root.Find("img_Lv/Txt_lv").GetComponent<Text>(),
                _obj_lock = root.Find("obj_Lock").gameObject,
                _txt_Lock = root.Find("obj_Lock/txt_Lock").GetComponent<Text>(),
                _obj_on = root.Find("on").gameObject,
                _img_frame = root.Find("img_head_kuang").GetComponent<Image>(),
                spine_Unlock = i == 0 ? null : root.Find("spine_unlock").GetComponent<SkeletonGraphic>(),
            };
        }

    }

    private void OnDestroy()
    {
		UITool.RemoveUI<UIRole>();
    }
}