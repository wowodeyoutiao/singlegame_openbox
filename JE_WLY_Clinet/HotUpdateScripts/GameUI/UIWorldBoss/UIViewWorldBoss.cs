using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewWorldBoss : MonoBehaviour,IUIView
{
	public Button btnHelp;

	public Image img_Moster;
	public Text txt_BossName;
    public GameObject obj_Kill;
	public Slider slider_Hp;
	public Text txt_hp;
	public Text txt_Skill;
	public Transform root_Award;

    public Button btn_RegionalServerRank;
    public GameObject btn_RegionalServerRankSelect;
    public Button btn_PersonInjuryRank;
    public GameObject btn_PersonInjuryRankSelect;

	public Text txt_RankTypeDesc;
    public LoopListView2 infiniteScroll;
    public Transform root_MyItem;

    public Text txt_tip;
	public Text txt_Challenge;
	public Button btn_Play;

    public Button btn_Close;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIWorldBoss>();
	}
}