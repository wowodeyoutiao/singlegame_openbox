using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using Spine.Unity;

public class UIViewSynthesis : MonoBehaviour,IUIView
{
	public Transform root_top;
	public Button btn_Close;
	public Button btn_Help;

	public Transform root_Center;

	public GameObject obj_NorTip;
	public Text txt_tip;
	public GameObject obj_lose;
	public Button btn_AutoAdd;

	public GameObject obj_SelectEnd;
	public Text txt_pro;
	public Text txt_money;
	public Image img_Money;
	public Button btn_Del;
	public Button btn_Syn;

	public SkeletonGraphic spine_anim;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UISynthesis>();
	}
}