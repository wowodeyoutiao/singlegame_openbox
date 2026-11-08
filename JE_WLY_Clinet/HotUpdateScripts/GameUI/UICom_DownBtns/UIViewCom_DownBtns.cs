using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;

public class UIViewCom_DownBtns : MonoBehaviour, IUIView
{
    public Image img_JyPro;
    public Text txt_JyPro;
    public BtnItem[] Btns;

    public Image grayImg;
    public Material grayMat;

    public class BtnItem
    {
        public Transform trans;
        public Button btn;
        public GameObject objDot;
        public Image iconImg;
        public GameObject selGo;

        public GameObject lockGo;
        public Image expImg;
        public Text expText;
        public Image sliderImg;
        public Transform effect;
        public GameObject mask;
        public Text txt_lock;

        public GameObject goRedDotRoot;//红点根节点
        public GameObject goFunctionOpenNewDotRoot;//功能开启New节点，该节点如果开启则需要隐藏红点
    }

    public void Awake()
    {
    }

    public void Init(GameObject go)
    {
        grayMat = grayImg.material;
        Btns = new BtnItem[7];
        for (int i = 1; i <= Btns.Length; i++)
        {
            var item = go.transform.Find("Btn_" + i);
            Btns[i - 1] = new BtnItem
            {
                trans = item,
                btn = JUIExtensions.btn(item),
                objDot = item.Find("redDotRoot/dot").gameObject,
                iconImg = item.Find("icon").GetComponent<Image>(),
                selGo = (i != Btns.Length ? item.Find("Select").gameObject : null),
                lockGo = (i == Btns.Length ? item.Find("lock").gameObject : null),
                expImg = (i == Btns.Length ? item.Find("TitleImg").GetComponent<Image>() : null),
                expText = (i == Btns.Length ? item.Find("txt_jyPro").GetComponent<Text>() : null),
                sliderImg = (i == Btns.Length ? item.Find("SliderImg").GetComponent<Image>() : null),
                effect = item.Find("spineEff"),
                mask = item.Find("mask").gameObject,
                txt_lock = JUIExtensions.txt(item, "mask/txt"),

                goRedDotRoot = item.Find("redDotRoot").gameObject,
                goFunctionOpenNewDotRoot = item.Find("functionOpenNewFlag").gameObject,
            };

            Btns[i - 1]?.goFunctionOpenNewDotRoot?.SetActive(false);//默认隐藏
            Btns[i - 1] ?.goRedDotRoot?.SetActive(true);//默认显示
        }
    }
}