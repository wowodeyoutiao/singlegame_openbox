using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;

public class UIViewGameDungeons : MonoBehaviour, IUIView
{
    public ETCJoystick _ETCJoystick;

    public Button btnExti;

    public Transform root_Base;

    public Transform root_Skill;

    public TouchEventCom touch_Skill;

    public Text txt_time;
    public Text txt_MonsterTip;

    public Transform root_Pos;
    public Text txt_Pos;

    //技能面板控制
    public Button btn_SkillPanel;
    public GameObject obj_SkillPanelOn;
    public GameObject obj_SkillPaneloOff;

    public Button btn_Attack;
    public Text txt_TouchNum;

    public Button btn_FightAward;

    public Transform SkillBtn_0;

    public Transform root_Chat;

    public Image img_lz;

    public Button btn_Pos;


    public struct SkillObj 
    {
        public Transform rect;
        public Image img_icon;
        public GameObject obj_on;
        public Image img_cd;
        public GameObject obj_lock;
        public Text txt_Time;
        public Text txt_Star;
        public Text txt_Reinforce;
        public GameObject obj_Dead;
    }

    [HideInInspector] public SkillObj[] SkillBtns;

    public Button btn_AutoPlay;
    public GameObject obj_AutoPlayIcon;
    public GameObject obj_AutoPlayOn;
    //public Text txt_AutoPlay;
    public GameObject obj_AutoPlayAni;

    //操作面板
    public GameObject obj_Touch;

    public Button btn_Chat;

    public void Awake()
    {
    }

    public void Init(GameObject go)
    {
        SkillBtns = new SkillObj[13];
        SkillBtns[0] = new SkillObj {
            rect = SkillBtn_0,
            img_icon = JUIExtensions.img(SkillBtn_0.Find("Icon")),
            obj_on = SkillBtn_0.Find("on").gameObject,
            img_cd = JUIExtensions.img(SkillBtn_0.Find("cd")),
            txt_Time = JUIExtensions.txt(SkillBtn_0.Find("cd/Time")),
            obj_lock = SkillBtn_0.Find("cd/lock").gameObject,
            txt_Star = JUIExtensions.txt(SkillBtn_0.Find("Icon/txt_Star")),
            txt_Reinforce = JUIExtensions.txt(SkillBtn_0.Find("Icon/txt_Reinforce")),
            obj_Dead = SkillBtn_0.Find("deadTip").gameObject
        };

        for (int i = 1; i <= root_Skill.childCount; i++)
        {
            Transform root = root_Skill.Find("SkillBtn_" + i);
            SkillBtns[i] = new SkillObj
            {
                rect = root,
                img_icon = JUIExtensions.img(root.Find("Icon")),
                obj_on = root.Find("on").gameObject,
                img_cd = JUIExtensions.img(root.Find("cd")),
                txt_Time = JUIExtensions.txt(root.Find("cd/Time")),
                obj_lock = root.Find("cd/lock").gameObject,
                txt_Star = JUIExtensions.txt(root.Find("Icon/txt_Star")),
                txt_Reinforce = JUIExtensions.txt(root.Find("Icon/txt_Reinforce")),
                obj_Dead = root.Find("deadTip").gameObject
            };
        }
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIGameDungeons>();
    }
}

