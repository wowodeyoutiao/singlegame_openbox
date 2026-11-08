using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewWarFightMap : MonoBehaviour,IUIView
{
	public Transform root_Btns;
	public Button btn_Close;
    public Button btn_Help;

    public BtnItem[] Btns;
	public struct BtnItem 
	{
        public Button btn;
		public Text txt;
		public Image img_type;
    }
    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
        Btns = new BtnItem[root_Btns.childCount];
		for (int i = 0; i < root_Btns.childCount; i++)
		{
			Btns[i] = new BtnItem {
				btn = root_Btns.GetChild(i).GetComponent<Button>(),
                txt = root_Btns.GetChild(i).Find("Text").GetComponent<Text>(),
                img_type = root_Btns.GetChild(i).Find("img_type").GetComponent<Image>()
            };
        }
    }
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIWarFightMap>();
	}
}