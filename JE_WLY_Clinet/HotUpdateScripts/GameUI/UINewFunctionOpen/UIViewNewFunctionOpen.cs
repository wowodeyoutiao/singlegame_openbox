using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewNewFunctionOpen : MonoBehaviour,IUIView
{
	public Button btn_Close;
    public ScrollRect scrollRect;
	public Text txt_Desc;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UINewFunctionOpen>();
	}
}