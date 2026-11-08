using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewAutoDelSetting : MonoBehaviour,IUIView
{
	public Dropdown dropdown_Qulite;
	public Dropdown dropdown_Level;

	public Button btn_Qulie;
	public GameObject obj_Quile_on;
	public Button btn_Close;
	public Button btn_Lv;
	public GameObject obj_Lv_on;

	public Button btn_Power;
	public GameObject obj_power_on;

	public Button btn_Sure;

	public Button btn_Help;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIAutoDelSetting>();
	}
}