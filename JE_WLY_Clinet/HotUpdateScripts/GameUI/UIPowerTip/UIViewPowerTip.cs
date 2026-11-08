using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPowerTip : MonoBehaviour,IUIView
{
	public Text txt_base;
	public Text txt_add;
	public CanvasGroup cGroup_add;

	public Transform eff;

	public Animator _Animator;
   

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPowerTip>();
	}
}