using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;

public class UIViewFightFormation : MonoBehaviour,IUIView
{
    public Transform root_Panel;
	public Transform root_GoTeams;
    [HideInInspector] public UpTeam[] goTeams;
    public Button btn_UnloadAll;
	public Button btn_UpLoadAll;

    public Transform root_AllTeams;
    [HideInInspector] public MyTeam[] allTeams;

    public Button btn_Left;
    public Button btn_Right;
    public Button btn_Play;
    public Button btn_Close;

    // 已上阵的
    public class UpTeam
    {
        public TouchEventCom touchEvent;
        public Text txt_Pos;
        public Text txt_PosDesc;
        public GameObject head;
        public Image img_Quality;
        public Image img_Head;
        public Image img_Job;
        public Text txt_Lv;
        public Text txt_Star;
        public Button btn_Remove;

        public void ShowHead(bool isShow)
        {
            head.SetActive(isShow);
        }

        public void ShowPositionText(bool isShow)
        {
            txt_Pos.gameObject.SetActive(isShow);
        }
    }
    
	// 可选择的
    public class MyTeam
	{
		public GameObject go;
        public TouchEventCom touchEvent;
        public Image img_Quality;
        public Image img_Head;
		public Image img_Job;
		public Text txt_Lv;
		public Text txt_Star;
		public GameObject obj_Mask;
        public Text txt_MaskTip;
        public Text txt_Name;
        public Outline outline_Name;
		public Image img_Power;
        public Text txt_Power;

        public void Show(bool isShow)
        {
            go.SetActive(isShow);
        }
    }


    public void Awake()
    {
        goTeams = new UpTeam[root_GoTeams.childCount];
        for (int i = 0; i < root_GoTeams.childCount; i++) 
        {
            Transform t = root_GoTeams.GetChild(i);
            goTeams[i] = new UpTeam
            {
                touchEvent = JUIExtensions.touch(t),
                txt_Pos = JUIExtensions.txt(t, "txt_pos"),
                txt_PosDesc = JUIExtensions.txt(t, "txt_posdesc"),
                head = t.Find("head").gameObject,
                img_Quality = JUIExtensions.img(t, "head/img_head_kuang"),
                img_Head = JUIExtensions.img(t, "head/img_head"),
                img_Job = JUIExtensions.img(t, "head/img_job"),
                txt_Lv = JUIExtensions.txt(t, "head/img_Lv/txt_lv"),
                txt_Star = JUIExtensions.txt(t, "head/root_star"),
                btn_Remove = JUIExtensions.btn(t, "head/btn_remove"),
            };
        }

        allTeams = new MyTeam[root_AllTeams.childCount];
        for (int i = 0; i < root_AllTeams.childCount; i++) 
        {
            Transform t = root_AllTeams.GetChild(i);
            allTeams[i] = new MyTeam
            {
                go = t.gameObject,
                touchEvent = JUIExtensions.touch(t),
                img_Quality = JUIExtensions.img(t, "img_head_kuang"),
                img_Head = JUIExtensions.img(t, "img_head"),
                img_Job = JUIExtensions.img(t, "img_job"),
                txt_Lv = JUIExtensions.txt(t, "img_Lv/txt_lv"),
                txt_Star = JUIExtensions.txt(t, "root_star"),
                obj_Mask = t.Find("img_mask").gameObject,
                txt_MaskTip = JUIExtensions.txt(t, "img_mask/txt_pos"),
                txt_Name = JUIExtensions.txt(t, "txt_name"),
                outline_Name = JUIExtensions.outline(t, "txt_name"),
                img_Power = JUIExtensions.img(t, "img_fight/img"),
                txt_Power = JUIExtensions.txt(t, "img_fight/txt_power"),
            };
        }
    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFightFormation>();
	}
}