using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewRoleTitile : MonoBehaviour,IUIView
{
    public Transform root_head;
    public Transform root_Team;
    public Button btn_Close;
    public Button btn_TitleActive;

    public UIViewRole.TeamItem[] teamItems;

    public Button btn_TitleResonate;
    public LoopListView2 loopView_Title;
    public Text txt_DressArr;
    public Text txt_BaseArr;

    public Image img_TitleName;
    public Transform root_Pice;
    public Transform root_Stars;

    public Transform root_UpStar;
    public Transform root_UpStarItem;
    public Text txt_NowStar;
    public Text txt_NextStar;
    public Text txt_UpPro;
    public Button btn_UpStar;
    public Button btn_LuckDel;

    public Button btn_TitleEquip;
    public Text txt_TitleEquip;
    public Button btn_TitleEquipDel;

    public Button btn_Help;
    public Text txt_nowStarLv;
    public Text txt_nextStarLv;
    public Text txt_NoTip;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
        teamItems = new UIViewRole.TeamItem[5];
        for (int i = 1; i <= root_Team.childCount; i++)
        {
            SetTeamItem(i, root_Team.GetChild(i - 1));
            teamItems[i].txt_Star = root_Team.GetChild(i - 1).Find("root_Star").GetComponent<Text>();
        }
        SetTeamItem(0, root_head);
        void SetTeamItem(int i, Transform root)
        {
            teamItems[i] = new UIViewRole.TeamItem
            {
                _btn = root.GetComponent<Button>(),
                _img_head = root.Find("img_head").GetComponent<Image>(),
                _img_head_quality = i == 0 ? root.Find("img_head_quality").GetComponent<Image>() : null,
                _img_job = root.Find("img_job").GetComponent<Image>(),
                _obj_Lv = root.Find("img_Lv").gameObject,
                _txt_lv = root.Find("img_Lv/Txt_lv").GetComponent<Text>(),
                _obj_lock = root.Find("obj_Lock").gameObject,
                _txt_Lock = root.Find("obj_Lock/txt_Lock").GetComponent<Text>(),
                _obj_on = root.Find("on").gameObject,
                _img_frame = root.Find("img_head_kuang").GetComponent<Image>(),
            };
        }
    }
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIRoleTitile>();
	}
}