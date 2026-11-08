using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewShop :MonoBehaviour, IUIView
{
	public class ToggleItem
	{
		public GameObject go;
		public Button btn;
		public Text norText;
		public Text selText;
		public GameObject selGo;
		public CFG_Store.StoreType sType;

		public ToggleItem(Transform trans)
		{
			go = trans.gameObject;
			btn = go.GetComponent<Button>();
			norText = trans.Find("Label").GetComponent<Text>();
			selText = trans.Find("Select/Label").GetComponent<Text>();
			selGo = trans.Find("Select").gameObject;

        }
    }

	public Button Btn_Close;

	public Button Btn_CouponAdd;
	public Button Btn_DiamantaAdd;
	public Button Btn_CoinAdd;
	public Text txt_Coupon;
	public Text txt_Diamanta;
	public Text txt_Coin;

	public Text txt_SLPZ;
	public Button btn_SLPZ;
	public GameObject obj_SLPZ;
	public GameObject obj_ZanHun;

	public Transform Btns;
	public ToggleItem[] toggleItems;

	public Image img_AD;
	public LoopGridView infiniteScroll;

	public Button btn_Banner;

	public GameObject ticketGo;
	public Text txt_Ticket;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
        toggleItems = new ToggleItem[Btns.childCount];
        for (int i = 0; i < Btns.childCount; i++)
        {
            toggleItems[i] = new ToggleItem(Btns.GetChild(i));
		}
	}

    private void OnDestroy()
    {
		UITool.RemoveUI<UIShop>();
    }

}