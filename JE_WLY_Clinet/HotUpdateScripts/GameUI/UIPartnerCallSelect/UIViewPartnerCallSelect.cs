using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;

public class UIViewPartnerCallSelect : MonoBehaviour, IUIView
{
    public Text titleText;

    public Button closeBtn;
    public Button okBtn;

    public Transform containerTrans;
    public Transform itemTrans;
    public Text priceText;

    public void Awake()
    {

    }

    public void Init(GameObject go)
    {

    }

    public Transform GetItemTrans()
    {
        var go = GameObject.Instantiate(itemTrans.gameObject);
        go.transform.SetParent(containerTrans);
        return go.transform;
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIPartnerCallSelect>();
    }
}