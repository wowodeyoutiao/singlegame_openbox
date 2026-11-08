using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewAwards : MonoBehaviour,IUIView
{
    public Transform root_Awards;
    public Button btn_back;
    public Text txt_Time;
	public Text txt_des;
    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIAwards>();
	}
}