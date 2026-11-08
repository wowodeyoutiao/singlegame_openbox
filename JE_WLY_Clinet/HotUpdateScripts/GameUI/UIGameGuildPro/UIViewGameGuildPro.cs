using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGameGuildPro : MonoBehaviour,IUIView
{
	public Button btn_Show;
	public GameObject obj_UI;
	public Button btn_Close;
	public Transform root_award1;
	public Transform root_award2;
	public Text txt_NoAward;

	public AwardItem[] award1Items;
	public AwardItem[] award2Items;

    public struct AwardItem
	{
		public Button _btn;
		public Image _imgBg;
		public Image _imgIcon;
		public Text _txtNum;
    }
    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		award1Items = new AwardItem[root_award1.childCount];
		for (int i = 0; i < root_award1.childCount; i++)
		{
			award1Items[i] = new AwardItem
			{
                _btn = root_award1.GetChild(i).GetComponent<Button>(),
                _imgBg = root_award1.GetChild(i).GetComponent<Image>(),
				_imgIcon = root_award1.GetChild(i).Find("icon").GetComponent<Image>(),
                _txtNum = root_award1.GetChild(i).Find("txt_num").GetComponent<Text>(),
            };
        }

        award2Items = new AwardItem[root_award2.childCount];
		for (int i = 0; i < root_award2.childCount; i++)
		{
			award2Items[i] = new AwardItem {
                _btn = root_award2.GetChild(i).GetComponent<Button>(),
                _imgBg = root_award2.GetChild(i).GetComponent<Image>(),
                _imgIcon = root_award2.GetChild(i).Find("icon").GetComponent<Image>(),
                _txtNum = root_award2.GetChild(i).Find("txt_num").GetComponent<Text>(),
            };
		}

    }
	public void OnDestroy()
	{
		 //UITool.RemoveUI<UIGameGuildPro>();
	}
}