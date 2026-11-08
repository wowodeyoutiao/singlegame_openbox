using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewLoading : MonoBehaviour,IUIView
{
	public Text txt_tip;
	public GameObject obj_Show0;
    public GameObject obj_Show1;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UILoading>();
	}
}