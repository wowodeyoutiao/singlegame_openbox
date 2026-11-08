using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPetUp : MonoBehaviour,IUIView
{
	public Text txt_title;

	public Image img_PetBg;
	public Image img_PetIcon;
	public Text txt_PetName;
	public Text txt_PetLv;

    public Button btnClose;

    public Button btnUpLv;

	public GameObject UI_UpLV;
    public Image img_lvUsed;
    public Text txt_lvUsed;
    public Text txt_nowLv;
	public Text txt_nextLv;
	public Text txt_nowLvDes;
    public Text txt_nextLvDes;

	public GameObject UI_Break;
	public Text txt_nowBreak;
    public Text txt_nowBreakDes;
    public Text txt_nextBreak;
	public Text txt_nextBreakDes;

	public Text txt_BreakLimitDes;
	public Image img_BreakPro;
	public Text txt_BreakPro;

	public Transform root_BreakTouch;

	public Transform root_Btns;

	public GameObject obj_BreakPro;
	public GameObject obj_BreakMax;
	public GameObject obj_RedUpStar;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		btnUpLv = go.transform.Find("Center/UI_UpLV/Btn_Up").GetComponent<Button>();
		//btnClose = go.transform.Find("Center/Btn_Close").GetComponent<Button>();
	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPetUp>();
	}
}