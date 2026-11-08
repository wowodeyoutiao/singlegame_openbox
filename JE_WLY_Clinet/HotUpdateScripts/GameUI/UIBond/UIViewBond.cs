using MetaFramework.UI;
using SuperScrollView;
using UnityEngine;
using UnityEngine.UI;

public class UIViewBond : MonoBehaviour, IUIView
{
    public Button closeBtn;
    public Button helpBtn;
    public Transform tabBtnTrans;
    [HideInInspector] public Button[] tabBtns;

    // left
    public LoopListView2 leftLoopListView;

    // right
    public Text rightTitleText;
    public LoopListView2 rightLoopListView;

    // gray
    public Image grayImg;

    public void Awake()
    {

    }

    public void Init(GameObject go)
    {

    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIBond>();
    }
}