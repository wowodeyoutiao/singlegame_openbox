using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using Spine.Unity;

public class UIViewFuNeng : MonoBehaviour,IUIView
{
	public Button btn_Help;

	public Transform root_resFuNeng;
	public Image img_FuNeng;
	public Text txt_FuNeng;
    public Button btn_AddFuNeng;

    public Transform root_resGold;
    public Image img_Gold;
    public Text txt_Gold;
    public Button btn_AddGold;

    public Transform root_resJinJie;
    public Image img_JinJie;
    public Text txt_JinJie;
    public Button btn_AddJinJie;

    public Transform root_Elements;
    [HideInInspector] public Element[] elements;

    public SkeletonGraphic spineJinJie;

    public class Element
    {
        public Button btn;
        public GameObject selected;
        public Text txt_Lv;
        public GameObject main;
        public int posId;

        public void Select(bool select)
        {
            selected.SetActive(select);
        }
    }

    public GameObject root_ConsumeLayout;
    public Text txt_ConsumeDesc;
    public Image img_ConsumeFuNeng1;
    public Text txt_ConsumeFuNeng1;
    public Image img_ConsumeFuNeng2;
    public Text txt_ConsumeFuNeng2;

    public Text txt_MaxLevel;

    public Text txt_Probalility;

    public GameObject root_Buttons;
    public Button btn_FuNeng;
    public Button btn_AutoFuNeng;
    public Text txt_AutoFuNeng;
    public Button btn_JinJie;

    public Button btn_Close;

    public RectTransform root_AttributeWindow;
    public Text txt_AttributeInfo;
    public Text txt_NoAttribute;

    public void Awake()
    {
        elements = new Element[root_Elements.childCount];
        for (int i = 0; i < root_Elements.childCount; i++)
        {
            Transform t = root_Elements.GetChild(i);
            elements[i] = new Element
            {
                btn = JUIExtensions.btn(t),
                selected = t.Find("selected").gameObject,
                txt_Lv = JUIExtensions.txt(t.Find("img_lv/txt_lv")),
                main = t.Find("img_mianAttr").gameObject,
            };
        }
    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFuNeng>();
	}
}