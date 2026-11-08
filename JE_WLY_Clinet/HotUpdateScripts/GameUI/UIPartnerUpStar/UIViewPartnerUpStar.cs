using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPartnerUpStar : MonoBehaviour,IUIView
{
	public Button btn_Close;

	public Transform root_Spine;
	public Text txt_name;
    public Text txt_Star;

    public Text txt_nowStar;
    public Text txt_nowDes;

    public Text txt_nextStar;
    public Text txt_nextDes;

    public Text txt_nextSkill;
    public Text txt_nextSkillName;
    public Button btn_Skill;

    public Button btn_SkillUI;
    public Text txt_ui_skillName;
    public Text txt_ui_skilltype;
    public Text txt_ui_skilldes;

    public Image img_used1;
    public Image img_used2;

	public Text txt_used1;
    public Text txt_used2;

	public GameObject obj_used;

	public Button btn_UpStar;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPartnerUpStar>();
	}
}