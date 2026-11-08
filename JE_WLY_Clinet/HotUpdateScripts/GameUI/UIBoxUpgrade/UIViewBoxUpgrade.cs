using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewBoxUpgrade : MonoBehaviour,IUIView
{
	public Text txt_CurrentLevelValue;
	public Text txt_NextLevelValue;

	public Transform root_QualityRandom;

    [HideInInspector] public List<Text> txtList_CurrentQualityRandom = new List<Text>();
    [HideInInspector] public List<Text> txtList_NextQualityRandom = new List<Text>();
    [HideInInspector] public List<GameObject> effect_UpLevel = new List<GameObject>();

    public Text txt_CurrentOpenBoxLimit;
    public Text txt_NextOpenBoxLimit;

    public Text txt_UpgradeCond;

    public Text txt_CannotUpgrade;

    public Image img_Coin;
    public Text txt_Coin;
    // 升级或加速按钮
    public Button btn_Upgrade;
    public Text txt_UpgradeButtonName;
    public Text txt_ExchangeTime;
    public Button btn_Close;

    public void Awake()
    {
        GetQualityRandomText();
    }

	public void Init(GameObject go)
	{

	}

    void GetQualityRandomText()
    {
        for (int i = 1; i < 9; i++)
        {
            var obj = root_QualityRandom.GetChild(i);
            txtList_CurrentQualityRandom.Add(obj.Find("txt_current").GetComponent<Text>());
            txtList_NextQualityRandom.Add(obj.Find("txt_next").GetComponent<Text>());
            Transform eff = obj.Find("eff_UIBoxUpgradei_jindu");
            if (eff)
                effect_UpLevel.Add(eff.gameObject);
        }
    }

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIBoxUpgrade>();
	}
}