using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewBag :MonoBehaviour, IUIView
{
	public Button Btn_Close;
	public Button Btn_Del;
	public Button Btn_Refresh;

	public Transform Btns;
	public Button[] TopBtns;

	public LoopGridView infiniteScroll;
	public Button btn_syn;

	public Button btn_AutoDel;
	public GameObject obj_AutoDel_on;
	public Text txt_AutoDe;
	public Button btn_AutoDeSet;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		TopBtns = new Button[Btns.childCount];
		for (int i = 0; i < Btns.childCount; i++)
		{
			TopBtns[i] = Btns.GetChild(i).GetComponent<Button>();
		}
	}
	private void OnDestroy()
	{
		UITool.RemoveUI<UIBag>();
	}
}