using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewExpPool : MonoBehaviour,IUIView
{
	public Text txt_expPro;
	public Image img_Pro;
	public Text txt_allExp;
    public Text txt_offline;
    public Transform root_jumpBtns;
	public Button btn_Close;
    public Button btn_help;
	public Transform root_Role;
	public Slider slider_Pro;

	public struct JumpBtn
	{
		public Button btn;
        public Text txt;
		public GameObject select;
    }
	public JumpBtn[] JumpBtns;

    public struct Role
	{
		public Image img_head;
		public GameObject go;
		public Image img_job;
        public Text txt_name;
        public Text txt_level;
		public Button btn_up;
		public Text txt_up;
		public GameObject obj_tip;
		public GameObject obj_Red;
        public Image img_quality;
    }

    public Role[] Roles;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		JumpBtns = new JumpBtn[root_jumpBtns.childCount];
        for (int i = 0; i < root_jumpBtns.childCount; i++)
        {
			Transform root = root_jumpBtns.GetChild(i);
            JumpBtns[i] = new JumpBtn {
                btn = root.GetComponent<Button>(),
				txt = root.Find("Label").GetComponent<Text>(),
                select = root.Find("Select").gameObject
            };
        }

        Roles = new Role[root_Role.childCount];

		for (int i = 0; i < root_Role.childCount; i++)
		{
			Transform root = root_Role.GetChild(i);
			Roles[i] = new Role
			{
                go = root.gameObject,
                btn_up = root.Find("btn_up").GetComponent<Button>(),
                txt_up = root.Find("btn_up/txt_exp").GetComponent<Text>(),
                img_job = root.Find("bg_head/img_job").GetComponent<Image>(),
                txt_name = root.Find("img_name/txt_name").GetComponent<Text>(),
                txt_level = root.Find("img_name/Txt_lv").GetComponent<Text>(),
                obj_tip = root.Find("txt_tip").gameObject,
                img_head = root.Find("bg_head/img_head").GetComponent<Image>(),
                obj_Red = root.Find("btn_up/red").gameObject,
                img_quality = root.Find("bg_head/img_head_kuang").GetComponent<Image>(),
            };
		}
    }
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIExpPool>();
	}
}