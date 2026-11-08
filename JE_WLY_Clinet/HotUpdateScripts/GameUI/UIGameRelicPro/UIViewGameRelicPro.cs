using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameRelicPro : MonoBehaviour,IUIView
{
	public Button btn_Show;
	public Button btn_Close;

	public Transform root_Tips;
	public GameObject obj_UI;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGameRelicPro>();
	}
}