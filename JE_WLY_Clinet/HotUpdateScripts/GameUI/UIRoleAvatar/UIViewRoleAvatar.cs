using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;
using JEngine.Helper;

public class UIViewRoleAvatar : MonoBehaviour,IUIView
{
	public Button btn_Help;
	public LoopGridView avatarList;
    public Image img_Frame;
	public Image img_Quality;
	public Image img_Head;
	public Text txt_Name;
	public Outline outline_Name;
    public Text txt_Attribute;
	public Text txt_GetFrom;
	public Button btn_Use;
    public Text txt_UseOrGetState;

    public Button btn_Close;

	public Transform root_SelectBtns;
	[HideInInspector] public Button[] btn_Selects;

    public void Awake()
    {
        btn_Selects = root_SelectBtns.GetComponentsInChildren<Button>(true);
    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIRoleAvatar>();
	}
}