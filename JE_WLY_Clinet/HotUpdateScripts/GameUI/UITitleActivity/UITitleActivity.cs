using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using SuperScrollView;
using HotUpdateScripts.ProtoM;
using System.Linq;

//称号共鸣
public class UITitleActivity : UIBase<UIModelTitleActivity, UIViewTitleActivity>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UITitleActivity>();
        });

        GetView().loopView.InitListView(0, LoopViewFill);
    }

    private LoopListViewItem2 LoopViewFill(LoopListView2 listView, int idx)
    {
        if (_TitleActivityList.Count <= idx || idx < 0)
            return null;

        var itemLoop = listView.NewListViewItem("Item");
        var item = itemLoop.transform;
        item.name = "item" + idx;

        var _CFGInfo = _TitleActivityList[idx];

        JUIExtensions.txt(item.Find("txt_title")).text = $"<color=#{(MaxStar >= _CFGInfo.StarLimit ? "76E344" : "D6D5D4")}>·{_CFGInfo.Limit}</color>" ;

        string str = "";
        for (int i = 0; i < _CFGInfo.Attribute.Length; i++)
        {
            var arr = _CFGInfo.Attribute[i].ElementAt(0);
            str += CFG_XMLConfig.GetAttributeStrByID(int.Parse( arr.Key), arr.Value) + "  ";
        }
        JUIExtensions.txt(item.Find("txt_Des")).text = $"<color=#{(MaxStar >= _CFGInfo.StarLimit ? "76E344" : "D6D5D4")}>·{str}</color>";

        //item.transform.Find("txt_Des").GetComponent<ContentSizeFitter>().SetLayoutVertical();
        //item.GetComponent<ContentSizeFitter>().SetLayoutVertical();
        return itemLoop;
    }

    private List<CFGItem_TitleActivity> _TitleActivityList;
    private int MaxStar = 999;

    private int _SoltIdx = 0;
    public override void OnEnter(params object[] args)
    {
        _SoltIdx = (int)args[0];
        _TitleActivityList = CFG_TitleActivity.GetActiveListByType(_SoltIdx + 1);

        var cfgTitles = CFG_Title.GetActiveListByType(_SoltIdx + 1);
        for (int i = 0; i < cfgTitles.Count; i++)
        {
            if (TitleM.Instance.CheckTitleIsGet(cfgTitles[i].UID))
            {
                MaxStar = Mathf.Min(MaxStar, TitleM.Instance.GetTitleStar(cfgTitles[i].UID));
            }
            else
            {
                MaxStar = -1;
                break;
            }
        }

        GetView().loopView.SetListItemCount(_TitleActivityList.Count);
        GetView().loopView.RefreshAllShownItem();
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}