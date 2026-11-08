using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewBattleSupportGift : MonoBehaviour,IUIView
{
    public Text txt_Des;
    public Transform root_Awards;
    public Text txt_OriginalPrice;
    public Button btn_Pay;
    public Button btn_close;

	public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIBattleSupportGift>();
	}
}