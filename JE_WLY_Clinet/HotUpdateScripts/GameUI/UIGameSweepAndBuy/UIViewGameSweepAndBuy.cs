using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameSweepAndBuy : MonoBehaviour,IUIView
{
	public Button btnAdd;
	public Button btnJian;
	public Button btnMin;
	public Button btnMax;
	public Button btnyes;

	public Button BtnClose;
	public InputField inputField;
	public Text txt_des;

	public Text txt_btn;

    public GameObject obj_used;
	public Image img_used;
	public Text txt_usedNum;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		btnAdd = go.transform.Find("BG/Center/InputField/Btn_Add").GetComponent<Button>();
		btnJian = go.transform.Find("BG/Center/InputField/Btn_Jian").GetComponent<Button>();
		btnMin = go.transform.Find("BG/Center/InputField/Btn_Min").GetComponent<Button>();
		btnMax = go.transform.Find("BG/Center/InputField/Btn_Max").GetComponent<Button>();
		btnyes = go.transform.Find("BG/Center/Btn_yes").GetComponent<Button>();

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGameSweepAndBuy>();
	}
}