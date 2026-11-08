using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewMessageBoxOne : MonoBehaviour,IUIView
{
	public Button Btn_Yes;
	public Text Text_des;
	public Text Text_Yes;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIMessageBoxOne>();
	}
}