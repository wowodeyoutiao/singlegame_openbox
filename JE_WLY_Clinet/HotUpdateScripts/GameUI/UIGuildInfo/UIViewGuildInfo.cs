using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewGuildInfo : MonoBehaviour,IUIView
{
	public Button btn_Help;
	public Image img_Head;
    public Text txt_Leader;
	public Text txt_PeopleCount;

	public Text txt_ContributionDegreeCurrentLevel;
	public Text txt_ContributionDegreeNextLevel;
	public Slider slider_LevelProgress;

	public Text txt_Progress;

	public Text txt_AttributeCurrentLevel;
	public Text txt_AttributeNextLevel;

	public LoopListView2 loopList_Attributes;

	public Image img_CoinContributionItem1;
    public Text txt_CoinContributionItemCount1;
    public Image img_CoinContributionItem2;
    public Text txt_CoinContributionItemCount2;
    public Text txt_CoinContributionTimes;
    public Button btn_CoinContribution;
    public Button btn_CoinContributionGray;

	public Image img_RollContributionItem1;
    public Text txt_RollContributionItemCount1;
    public Image img_RollContributionItem2;
    public Text txt_RollContributionItemCount2;
    public Text txt_RollContributionTimes;
    public Button btn_RollContribution;
    public Button btn_RollContributionGray;

    public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuildInfo>();
	}
}