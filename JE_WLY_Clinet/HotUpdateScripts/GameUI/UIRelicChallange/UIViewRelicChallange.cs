using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewRelicChallange : MonoBehaviour,IUIView
{
	public Button btn_Sd;
	public Button btn_Play;
	public Button btn_close;

	public Text txt_title;

	public Transform root_des;
    public Transform root_award;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIRelicChallange>();
	}
}