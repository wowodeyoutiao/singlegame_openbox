using MetaFramework.UI;
using Spine.Unity;
using UnityEngine;
using UnityEngine.UI;

public class UIViewGuide : MonoBehaviour,IUIView
{
	public GameObject maskGo;
	public GameObject battleMaskGo;
	public Button closeBtn;
	public Button skipBtn;
	public GameObject skipGo;
	public Transform focusTrans;
    private RectTransform focusRect;
	public GameObject focusGo;
    public Button focusBtn;
	public SkeletonGraphic skeletonGraphic;
    public GameObject tipGo;
	public Text desText;

    public GameObject leftArrowGo;
    public GameObject rightArrowGo;

    public void Awake()
    {
		skipGo = skipBtn.gameObject;
        focusGo = focusTrans.gameObject;
		focusBtn = focusGo.GetComponent<Button>();
        focusRect = focusGo.GetComponent<RectTransform>();
		skeletonGraphic = focusGo.GetComponent<SkeletonGraphic>();
    }
	public void Init(GameObject go)
	{

	}

    public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuide>();
	}
}