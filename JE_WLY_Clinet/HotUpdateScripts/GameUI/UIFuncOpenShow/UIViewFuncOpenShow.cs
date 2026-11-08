using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using Spine.Unity;
using System;

public class UIViewFuncOpenShow : MonoBehaviour,IUIView
{
	//点击区域
	public Button btn_Click;
	//特效层
	public Transform root_Effect;
    //功能Icon
    public Transform tf_Icon;
    //功能Icon
    public Image img_Icon;

    public GameObject goTitleRoot;


    public void Awake()
    {
        
    }

	public void Init(GameObject go)
	{

    }

	private Action mAutoAction = null;
	public void SetAutoAction(Action autoAction)
	{
        mAutoAction = autoAction;
        Invoke("AutoSendPlayEffect", 3);
    }
	private void AutoSendPlayEffect()
	{
        mAutoAction?.Invoke();
    }


    public void OnDestroy()
	{
		 UITool.RemoveUI<UIFuncOpenShow>();
	}
}