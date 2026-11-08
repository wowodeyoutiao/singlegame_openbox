using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using System;
using JEngine.Helper;
using Spine.Unity;

public class UIViewCom_TopRes : MonoBehaviour,IUIView
{
	//顶部货币显示
	public Text txt_coupon;
	public Text txt_gold;
	public Text txt_diamanta;
	public Button btn_coupon;
	public Button btn_gold;
	public Button Btn_Set;

	public Button btn_AddTip;
	public GameObject obj_Add;
	public Image img_addTip;

	public GameObject obj_SetUI;
	public Transform root_Set;

    public Button btn_MyHead;

    public Image img_MyHeadQuality;
    public Image img_MyHeadFrame;
    public Image img_MyHead;

    public Image img_MyJob;
	public Text txt_MyName;
    public Button btn_Vip;
	public Text txt_Lv;
	public Text txt_VIP;
	public Text txt_Power;
	public GameObject obj_MyBuferr;
	public Image img_MyBuferr;

    public Transform root_Team;
	[HideInInspector] public TeamItem[] TeamItems;
    [HideInInspector] public GameObject obj_MyOn;

	public Text txt_ping;

    [HideInInspector] public Action UpateCallBack;

    public struct TeamItem 
	{
		public Image img_Head;
		public Text txt_Lv;
		public Image img_Job;
		public Text txt_Lock;
		public GameObject obj_Lock;
        public GameObject obj_Lv;
		public GameObject starGo;
		public Text startText;
		public Button _btn;
		public GameObject _obj_on;
        public Image img_quality;
		public GameObject obj_Buffer;
		public Image img_Buffer;
		public ParticleSystem ps_Unlock;
		public SkeletonGraphic spine_Unlock;
    }

    public void Awake()
    {
		obj_MyOn = btn_MyHead.transform.Find("on").gameObject;

        TeamItems = new TeamItem[root_Team.childCount];

        for (int i = 0; i < root_Team.childCount; i++)
		{
			Transform root = root_Team.GetChild(i);
            TeamItems[i] = new TeamItem
			{
                _btn = root.GetComponent<Button>(),
				img_Head = root.Find("img_head").GetComponent<Image>(),
				txt_Lv = root.Find("img_Lv/Txt_lv").GetComponent<Text>(),
				img_Job = root.Find("obj_Team").GetComponent<Image>(),
				obj_Lock = root.Find("obj_Lock").gameObject,
				txt_Lock = root.Find("obj_Lock/txt_Lock").GetComponent<Text>(),
				obj_Lv = root.Find("img_Lv").gameObject,
				starGo = root.Find("root_Star").gameObject,
				startText = root.Find("root_Star/StarText").GetComponent<Text>(),
				_obj_on = root.Find("on").gameObject,
				img_quality = root.Find("img_head_kuang").GetComponent<Image>(),
				obj_Buffer = root.Find("buffer").gameObject,
				img_Buffer = root.Find("buffer/world").GetComponent<Image>(),
				ps_Unlock = root.Find("eff_other_uiCom_boom").GetComponent<ParticleSystem>(),
				spine_Unlock = root.Find("spine_unlock").GetComponent<SkeletonGraphic>(),
            };

        }
    }

	public void Init(GameObject go)
	{

	}

    public void LateUpdate()
    {
        UpateCallBack?.Invoke();
    }
}