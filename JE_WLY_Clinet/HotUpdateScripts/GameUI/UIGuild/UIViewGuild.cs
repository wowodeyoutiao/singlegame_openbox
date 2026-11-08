using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewGuild : MonoBehaviour,IUIView
{
    public Text txt_title;
    public Transform root_GuildChildPanel;
    public Transform root_FriendChildPanel;
    public Transform root_ExpandChildPanel;

    public Transform root_SubitemButtons;
    public Button[] btns_Subitem;

    public Transform root_BottomButtons;
    public Button[] btns_Bottom; 

    public Button btn_Close;

    public void Awake()
    {
		
    }

	public void Init(GameObject go)
	{
        btns_Subitem = new Button[root_SubitemButtons.childCount];
        for (int i = 0; i < root_SubitemButtons.childCount; i++)
        {
            btns_Subitem[i] = root_SubitemButtons.GetChild(i).GetComponent<Button>();
        }

        btns_Bottom = new Button[root_BottomButtons.childCount];
        for (int i = 0; i < root_BottomButtons.childCount; i++)
        {
            btns_Bottom[i] = root_BottomButtons.GetChild(i).GetComponent<Button>();
        }
    }

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuild>();
	}
}