using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewCommonItemBuy : MonoBehaviour,IUIView
{
	public Text txt_Title;
	public Text txt_subTitle;

	public Button btn_Close;
	public Button btn_add;
	public Button btn_jian;
	public InputField input_Buy;
	public Button btn_Sun10;
	public Button btn_Add10;
	public Text txt_tip;
	public Button btn_buy;
	public Text txt_BtnName;

	public Image img_Consume;
	public Text txt_Consume;
	public Text txt_ConsumeDesc;
	public Text txt_ConsumeDesc2;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UICommonItemBuy>();
	}
}