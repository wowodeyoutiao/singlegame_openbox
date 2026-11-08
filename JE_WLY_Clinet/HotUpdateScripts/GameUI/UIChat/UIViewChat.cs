using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;
using System;

public class UIViewChat : MonoBehaviour,IUIView
{
    public LoopListView2 loopListView;
    public Button btn_Voice;
    public Button btn_Set;
	public InputField input_chat;
	public Button btn_SendEmoji;
	public Button btn_Send;
    public LoopGridView loopGrid_Emoji;
    public Button btn_emojiUI;
    public Image img_chatItemIcon;
    public Text txt_chatItemCount;
    public Transform btns_tagRoot;
    public Button btn_Close;
    public Button btn_SendPos;

	[HideInInspector] public Button[] ButtonTags;

    public void Awake()
    {
        ButtonTags = new Button[btns_tagRoot.childCount];
        for (int i = 0; i < btns_tagRoot.childCount; i++)
        {
            ButtonTags[i] = btns_tagRoot.GetChild(i).GetComponent<Button>();
        }
    }

	public void Init(GameObject go)
	{

    }

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIChat>();
	}
}