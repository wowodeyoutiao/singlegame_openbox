using MetaFramework.UI;
using SuperScrollView;
using UnityEngine;
using UnityEngine.UI;

public class UIViewSpecialOffer : MonoBehaviour, IUIView
{
    public Transform tabBtnsTrans;
    [HideInInspector] public Button[] tabBtns;
    [HideInInspector] public ScrollRect scrollTabBtns;

    public Button btn_Close;

    // 每日礼包、特惠礼包、宠物礼包、伙伴礼包
    public GameObject normalNodeGo;
    public Image logoImg;
    public LoopListView2 infiniteScroll_Gifts;
    public Text txt_ActiveTime;

    // 成长基金
    public GameObject growthFundNodeGo;
    public Text curRoleLevelText;
    public LoopListView2 infiniteScroll_GrowthFund;
    public GameObject titleLockGo;
    public Button growthFundBuyBtn;

    // 月卡
    public GameObject monthNodeGo;
    public Button monthLBtn;
    public Button mothRBtn;
    public GameObject monthLFinishImgObj;
    public GameObject monthRFinishImgObj;
    public Text monthTimeText;

    public void Awake()
    {

    }

    public void Init(GameObject go)
    {
        scrollTabBtns = tabBtnsTrans.GetComponentInParent<ScrollRect>();

        tabBtns = new Button[tabBtnsTrans.childCount];

        for (int i = 0; i < tabBtnsTrans.childCount; i++)
            tabBtns[i] = tabBtnsTrans.GetChild(i).GetComponent<Button>();
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UISpecialOffer>();
    }
}