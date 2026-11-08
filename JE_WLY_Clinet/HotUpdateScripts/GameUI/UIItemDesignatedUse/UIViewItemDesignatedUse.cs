using MetaFramework.UI;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UIViewItemDesignatedUse : MonoBehaviour, IUIView
{
    public class DesignatItem
    {
        public GameObject go;
        public Button btn;
        public GameObject selGo;
        public Image bgImg;
        public Button itemBtn;
        public Image iconImg;
        public Text nameText;
        public Outline nameOutline;

        public DesignatItem(Transform trans)
        {
            go = trans.gameObject;
            btn = trans.Find("SelButton").GetComponent<Button>();
            selGo = trans.Find("SelButton/SelImg").gameObject;
            bgImg = trans.Find("BgImg").GetComponent<Image>();
            itemBtn = trans.Find("BgImg").GetComponent<Button>();
            iconImg = trans.Find("IconImg").GetComponent<Image>();
            nameText = trans.Find("NameText").GetComponent<Text>();
            nameOutline = trans.Find("NameText").GetComponent<Outline>();
        }

        public void SetActive(bool active)
        {
            go.SetActive(active);
        }
    }

    public Text titleText;
    public Button backBtn;
    public Button okBtn;

    public Transform containterTrans;
    public List<DesignatItem> designatItems = new List<DesignatItem>();
    private GameObject itemGo;

    public void Awake()
    {
        for(int i = 1; i < containterTrans.childCount; i++)
        {
            designatItems.Add(new DesignatItem(containterTrans.GetChild(i)));
        }
    }

    public DesignatItem GetDesignatItem(int index)
    {
        if (designatItems.Count < index)
        {
            var go = GameObject.Instantiate(itemGo);
            go.transform.localPosition = Vector3.zero;
            go.transform.localScale = Vector3.one;
            designatItems.Add(new DesignatItem(go.transform));
        }

        return designatItems[index];
    }

    public void Init(GameObject go)
    {

    }

    public void OnDestroy()
    {
        designatItems.Clear();
    }
}