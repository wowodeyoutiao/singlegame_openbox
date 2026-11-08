using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewFirstPay : MonoBehaviour,IUIView
{
	public Image day_image;
    public Text txt_Des;
    public Text txt_Time;

    public Button btn_close;
    public Button btn_Pay;
	public Button btn_get;
	public Button btn_double_buy;

	public Transform root_Awards;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFirstPay>();
	}
}