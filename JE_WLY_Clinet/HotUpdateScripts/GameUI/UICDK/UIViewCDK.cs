using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewCDK :MonoBehaviour, IUIView
{
	public Button btn_Close;
	public Button btn_Sure;
	public InputField Input_CDK;
	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}

     void OnDestroy()
    {
		UITool.RemoveUI<UICDK>();
    }

}