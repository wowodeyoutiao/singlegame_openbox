using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewDayDo : MonoBehaviour,IUIView
{
	public Transform root_Btns;
	public Text txt_Title;
	public Text txt_pro;
	public Button btn_Add;

	public Image img_Monster;
	public Text txt_RoomDes;

	public Transform root_Award;

	public Button btn_Paly;
	public Button btn_AutoPlay;
	public Text txt_vip;
	public GameObject obj_vip;

	public GameObject obj_lock1;
	public GameObject obj_lock2;
	public Text txt_lock2;
	public Text txt_lock1;

	public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 //UITool.RemoveUI<UIDayDo>();
	}
}