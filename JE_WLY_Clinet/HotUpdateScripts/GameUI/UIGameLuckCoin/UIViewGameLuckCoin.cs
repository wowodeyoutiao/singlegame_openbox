using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;

public class UIViewGameLuckCoin : MonoBehaviour, IUIView
{
    public Button btn_Close;
    public Button btn_Coin;

    public Image img_Bg;
    public Image img_Coin;
    public Image img_Title;
    public Image img_zs1;
    public Image img_zs2;

    public Text txt_Tip;
    public Transform root_Award;
    public Text txt_Time;

    public void Awake()
    {
    }
    public void Init(GameObject go)
    {

    }
    public void OnDestroy()
    {
        UITool.RemoveUI<UIGameLuckCoin>();
    }
}