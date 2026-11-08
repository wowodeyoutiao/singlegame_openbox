using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPartnerBook : MonoBehaviour, IUIView
{
    public Button closeBtn;
    public GameObject contentGo;
    public Transform itemTrans;
    public Image grayImg;
    public Material grayMat;

    public void Awake()
    {
        grayMat = grayImg.material;
    }

    public void Init(GameObject go)
    {

    }

    public Transform GetItemTrans()
    {
        var go = GameObject.Instantiate(itemTrans.gameObject);
        go.transform.SetParent(itemTrans.parent);
        return go.transform;
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIPartnerBook>();
    }
}