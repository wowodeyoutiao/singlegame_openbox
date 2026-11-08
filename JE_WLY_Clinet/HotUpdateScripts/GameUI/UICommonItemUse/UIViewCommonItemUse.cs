using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewCommonItemUse : MonoBehaviour,IUIView
{
	public InputField input_Use;
    public Button btn_jian;
	public Button btn_add;
	public Button btn_Max;
    public Button btn_use;
    public Button btn_Close;


	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UICommonItemUse>();
	}
}