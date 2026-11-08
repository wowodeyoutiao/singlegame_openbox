using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewWarFightInfo : MonoBehaviour,IUIView
{
	public Transform root_Awards;
	public Button btn_Close;
	public Text txt_PlayerNum;
    public Text txt_Title;
    public Button btn_Play;
	public Image img_Monster;
	public Image img_MonsterBg;
	public Text txt_Tip;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIWarFightInfo>();
	}
}