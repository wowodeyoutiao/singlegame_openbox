using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameOverAwards : MonoBehaviour,IUIView
{
	public Button btn_Mask;
	public Transform root_Cards;
	public Text txt_Time;

	public GameObject obj_CloseTip;
    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGameOverAwards>();
	}
}