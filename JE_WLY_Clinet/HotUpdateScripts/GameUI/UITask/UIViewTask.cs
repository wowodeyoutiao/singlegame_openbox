using MetaFramework.UI;
using SuperScrollView;
using UnityEngine;
using UnityEngine.UI;

public class UIViewTask : MonoBehaviour, IUIView
{
    public Transform root_topBtns;
    public Button btn_Close;

    // 任务
    public GameObject taskViewGo;
    public LoopListView2 InfiniteScroll_Task;
    public GameObject obj_DayAndWeek;
    public Text txt_HyAll;
    public Text txt_HyAllNum;
    public Transform root_ProBtn;
    public Slider slider_Pro;

    public RectTransform InfiniteScroll_Rect;

    // 每日签到
    public GameObject signDayGo;
    public LoopGridView ScrollView_SignDay;
    public LoopListView2 ScrollView_Cumulative;
    public GameObject obj_VipMax;
    public Button btn_VIPUp;

    // 标题
    public Text titleText;

    // 周重置
    public Text weekResetTime;

    public void Awake()
    {

    }
    public void Init(GameObject go)
    {
        InfiniteScroll_Rect = InfiniteScroll_Task.GetComponent<RectTransform>();
    }
    public void OnDestroy()
    {
        UITool.RemoveUI<UITask>();
    }
}