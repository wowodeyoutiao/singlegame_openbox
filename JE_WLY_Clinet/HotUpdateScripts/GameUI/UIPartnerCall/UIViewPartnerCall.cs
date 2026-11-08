using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using System.Collections.Generic;

public class UIViewPartnerCall : MonoBehaviour,IUIView
{
	public Button helpBtn;
	public Text gemstoneText;
	public Text currencyText;
	public Button gemstoneAddBtn;
	public Button currencyAddBtn;

	// 门动画
	public Transform doorSpineTrans;

    // 图鉴
    public Button handbookButton;

	// 羁绊
	public Button bondButton;

	// 召唤
	public Image sliderImg;
	public Text sliderText;
	public Button selBtn;
	public GameObject sliderSpineGo;

	// center item
	public Transform centerItemsTrans;
	public List<Transform> itemsTrans = new List<Transform>();

    // 按钮
    public Button freeBtn;
	public Button callBtn;
    public Image callPriceIconImg;
    public Text callPriceText;
    public Button refreshBtn;
    public Image refresPriceIconImg;
    public Text refresPriceText;

    public Button backBtn;
	// 倒计时
	public Text timerText;

	// 置灰的材质球
	public Image grayImg;

	// 召唤成功特效
	public GameObject callBoomEffectGo;

	// 刷新特效
	public GameObject refreshBtnEffectGo;


    public void Awake()
    {
		for(int i = 0; i < centerItemsTrans.childCount; i++)
		{
			itemsTrans.Add(centerItemsTrans.Find("Item_" + (i + 1)));
        }
    }

	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		UITool.RemoveUI<UIPartnerCall>();
	}
}