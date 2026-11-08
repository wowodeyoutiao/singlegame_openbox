using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewMessageDestoryRole : MonoBehaviour,IUIView
{
	public Button Btn_Yes;
	public Button Btn_No;
	public InputField Intput;


    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIMessageDestoryRole>();
	}
}