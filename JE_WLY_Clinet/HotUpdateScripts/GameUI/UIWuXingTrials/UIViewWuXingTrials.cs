using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;
using JEngine.UI;

public class UIViewWuXingTrials : MonoBehaviour, IUIView
{
	public LoopListView2 loopLayers;
	[HideInInspector] public ScrollRect scrollTower;

    public Transform cloudRoot;
    [HideInInspector] public RectTransform[] clouds;

	public Button btn_Help;
	public Transform rankRoot;
    public Button btn_Rank;
    [HideInInspector] public RankItem[] rankItems;
	public Button btn_Gift;
	public Button btn_FuNeng;

	public Transform awardRoot;
	public Text txt_ChallengeTimes;
	public Button btn_AddChallengeTimes;
	public Button btn_Challenge;
    [HideInInspector] public Image img_Challenge;
	public Button btn_Sweep;
	public Text txt_ResetTime;

	public CanvasGroup levelDescCanvas;
	public Text txt_LevelDesc;

    public Button btn_Close;

    public struct RankItem
	{
		public Text txt_Name;
		public Text txt_Layer;
	}


	public void Awake()
    {
		scrollTower = loopLayers.GetComponent<ScrollRect>();
        img_Challenge = btn_Challenge.GetComponent<Image>();

        clouds = new RectTransform[cloudRoot.childCount];
		for (int i = 0; i < cloudRoot.childCount; i++)
			clouds[i] = cloudRoot.GetChild(i).GetComponent<RectTransform>();

		rankItems = new RankItem[rankRoot.childCount - 1];
		for (int i = 1; i < rankRoot.childCount; i++)
		{
            rankItems[i - 1] = new RankItem
            {
                txt_Name = JUIExtensions.txt(rankRoot.GetChild(i).Find("txt_name")),
                txt_Layer = JUIExtensions.txt(rankRoot.GetChild(i).Find("txt_layer")),
            };
        }

    }

    public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIWuXingTrials>();
	}
}