using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewOtherRole : MonoBehaviour, IUIView
{
	public Transform root_head;
	public Text txt_allPower;
	public Transform root_Team;

    public Button btn_Close;

    public Text txt_RoleName;
    public Text txt_RolePower;
    public Text txt_lv;
    public Transform root_Equips;

    public Transform img_Pet;

    public Button btn_Attribute;
    public Text txt_attr0;
    public Text txt_attr1;
    public Text txt_attr2;
    public Text txt_attr3;
    public Transform arrow_Attribute;

    public Transform root_Role;
    public Image img_Role;

	public struct TeamItem
    {
        public Image _img_head;
        public Image _img_job;
        public GameObject _obj_Lv;
        public Text _txt_lv;
        public Button _btn;
		public GameObject _obj_lock;
        public Text _txt_Lock;
        public Text txt_Star;
        public GameObject _obj_on;
        public Image img_frame;
        public Image _img_head_quality;

    }

    [HideInInspector] public TeamItem[] teamItems;

    public void Awake()
	{

	}

	public void Init(GameObject go)
	{
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
                _img_head =  root.Find("img_head").GetComponent<Image>(),
                _img_head_quality = i == 0 ? root.Find("img_head_quality").GetComponent<Image>() : null,
                _img_job = root.Find("img_job").GetComponent<Image>(),
				_obj_Lv = root.Find("img_Lv").gameObject,
				_txt_lv = root.Find("img_Lv/Txt_lv").GetComponent<Text>(),
				_obj_lock = root.Find("obj_Lock").gameObject,
                _txt_Lock = root.Find("obj_Lock/txt_Lock").GetComponent<Text>(),
                _obj_on = root.Find("on").gameObject,
                img_frame = root.Find("img_head_kuang").GetComponent<Image>(),
            };
        }

    }

    private void OnDestroy()
    {
		UITool.RemoveUI<UIOtherRole>();
    }
}