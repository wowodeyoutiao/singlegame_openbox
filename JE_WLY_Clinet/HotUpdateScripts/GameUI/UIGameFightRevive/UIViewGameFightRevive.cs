using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameFightRevive : MonoBehaviour,IUIView
{
	public Button Btn_Revive;
	public Button Btn_Back;
	public Text Text_des;
	public Text Text_Tip;
	public Text txt_Back;
	public Text txt_Revive;


    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGameFightRevive>();
	}
}