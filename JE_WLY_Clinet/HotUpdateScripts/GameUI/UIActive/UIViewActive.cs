using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewActive : MonoBehaviour, IUIView
{
    public Button closeBtn;
    public Transform tabBtnTrans;
    [HideInInspector] public Button[] tabBtns;
    public Text activeTimeText;

    // 七日签到
    public GameObject sign7Go;
    public Text signDesText;
    public LoopListView2 infiniteScroll_Sign7;

    // 七日狂欢 
    // 大奖
    public GameObject happy7Go;
    public Slider happy7Slider;
    public Transform happy7SliderAwardTrans;
    public Text happy7SchduleText;

    public Transform happy7TabBtnTrans;
    [HideInInspector] public Button[] happy7TabBtns;
    public LoopListView2 infiniteScroll_happy7;

    public Transform starTrans;
    public GameObject flyEffectGo;

    // 特惠礼包
    public Transform happy7GiftAwardTrans;
    public Transform happy7GiftAwardBtn;
    public GameObject GetGo;

    // 星愿召唤/疯狂洗练
    public Transform starWishSummonTrans;
    public Text txt_StarWishTime;
    public Text txt_StarWishDesc;
    public Button btn_Help;
    public LoopGridView loopGrid_StarWish;

    // 成长基金
    public GameObject rootFund;
    public LoopListView2 loop_GrowthFund;
    public GameObject titleLockGo;
    public Button btnFundBuy;
    public Image img_FundLogo;
    public Text txt_BtnPrice;
    public Text txt_FundTitle;
    public Text txt_FundTitle2;
    public Text txt_FundTime;
    public Text txt_TaskCount;

    public void Awake()
    {
        tabBtns = new Button[tabBtnTrans.childCount];
        for (int i = 0; i < tabBtnTrans.childCount; i++)
            tabBtns[i] = tabBtnTrans.GetChild(i).GetComponent<Button>();

        happy7TabBtns = new Button[happy7TabBtnTrans.childCount];
        for (int i = 0; i < happy7TabBtnTrans.childCount; i++)
            happy7TabBtns[i] = happy7TabBtnTrans.GetChild(i).GetComponent<Button>();
    }
    public void Init(GameObject go)
    {

    }

    public GameObject GetFlyEffectGo()
    {
        return GameObject.Instantiate(flyEffectGo, happy7Go.transform);
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIActive>();
    }
}