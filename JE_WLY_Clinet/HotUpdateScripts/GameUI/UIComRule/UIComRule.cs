using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using SuperScrollView;

public class UIComRule : UIBase<UIModelComRule, UIViewComRule>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(()=> {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIComRule>();
        });

        //GetView().infiniteScroll.InitListView(0, OnFillItem);
    }

    private LoopListViewItem2 OnFillItem(LoopListView2 listView, int idx)
    {
        if (idx < 0 || idx >= strDess.Length)
        {
            return null;
        }
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        JUIExtensions.txt(item.transform.Find("txt")).text = strDess[idx];
        GameHelper.Instance.SetDeayFrame(() => {
            LayoutRebuilder.ForceRebuildLayoutImmediate(item.transform.Find("txt").GetComponent<RectTransform>());
            item.GetComponent<RectTransform>().sizeDelta = new Vector2(item.GetComponent<RectTransform>().sizeDelta.x, item.transform.Find("txt").GetComponent<RectTransform>().sizeDelta.y);
        });
        return item;
    }

    private int OnHeight(int idx)
    {
        return Mathf.CeilToInt(strDess[idx].Length / 15f) * 30;
    }

    private string[] strDess;
    public override void OnEnter(params object[] args)
    {
        strDess = CFG_RegularDes.GetById((int)args[0]).describe.Split('|');

        if (args.Length > 1)
            GetView().txt_name.text = (string)args[1];
        else
            GetView().txt_name.text = "规则说明";

        string str = "";
        for (int i = 0; i < strDess.Length; i++)
        {
            str += $"{strDess[i]}\n";
        }
        GetView().txt_des.text = str;
        //GetView().infiniteScroll.SetListItemCount(strDess.Length - 1);
        //GetView().infiniteScroll.RefreshAllShownItem();

        GameHelper.Instance.SetDeayFrame(() => {
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_des.GetComponent<RectTransform>());
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_des.transform.parent.GetComponent<RectTransform>());
        });
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}