using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameFightOver : MonoBehaviour,IUIView
{
	public Button btn_mask;
    public GameObject obj_win;
    public Text txt_WinDesc;
    public GameObject obj_lose;
    public Text txt_LoseDesc;

    public Transform root_Awards1;
	public Text txt_Awards1Tip;
    public Transform root_Awards2;
	public Text txt_Awards2Tip;

	public GameObject obj_line;

	public Transform root_useTime;
	public Text txt_useTime;
	public GameObject obj_newRecord;

	public Text txt_point;

	public GameObject obj_space;

    public Button btn_yellow;
    public Text txt_yellow;

    public Button btn_blue;
    public Text txt_blue;

	public Text txt_num;

    public Text txt_Time;
	public Text txt_CloseTip;

	public Transform root_LuckAward;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGameFightOver>();
	}
}