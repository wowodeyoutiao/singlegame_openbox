using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using System;

public class UIViewGameFightWarMap : MonoBehaviour,IUIView
{
	public Text txt_Tiltle;
	public Button btn_Close;

	public TouchEventCom touch_Map;
	public RectTransform rect_Map;

	public Transform root_ToPos;
	public Text txt_ToPos;

	public Transform root_MyPos;

	public GameObject obj_MyTeam;
	public GameObject obj_PlayerEnemy;
    public GameObject obj_Boss;

    public Action UpateCallBack;

    public LineRenderer line;

    public Image img_Map;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
        rect_Map = touch_Map.transform.GetComponent<RectTransform>();
    }

    public void OnDestroy()
    {
        UpateCallBack = null;
        UITool.RemoveUI<UIGameFightWarMap>();
    }

    public void LateUpdate()
    {
        UpateCallBack?.Invoke();

    }
}