using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;

public class UIViewActivityAwards : MonoBehaviour, IUIView
{
    public Button closeBtn;
    public Image titleImg;
    public Text desText;
    public Transform scrollViewTrans;
    public Transform containerTrans;
    public GameObject itemGo;
    public Text timeText;
    public Image LightImg;

    public GameObject GetItem()
    {
        var item = GameObject.Instantiate(itemGo, containerTrans);
        item.transform.localScale = Vector3.one;

        return item;
    }

    public void Clear()
    {
        for (int i = containerTrans.childCount - 1; i > 0; i--)
        {
            DestroyImmediate(containerTrans.GetChild(i).gameObject);
        }
    }

    public void Awake()
    {
        
    }

    public void Init(GameObject go)
    {
        
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIActivityAwards>();
    }
}