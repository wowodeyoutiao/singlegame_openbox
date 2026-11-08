using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewNoticeRun : MonoBehaviour,IUIView
{
	public RectTransform rect_Bg;
	public Text txt_Notice;
	public RectTransform rect_Notice;
    public Text txt_Jump;
    public RectTransform rect_Jump;

    public Button btn_Notice;

	public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UINoticeRun>();
	}
}