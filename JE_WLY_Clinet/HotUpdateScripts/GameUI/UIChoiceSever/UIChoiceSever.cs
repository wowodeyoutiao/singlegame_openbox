using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.UI;
using JEngine.Helper;
using System;
using SuperScrollView;
using HotUpdateScripts;

public class UIChoiceSever : UIBase<UIModelChoiceSever, UIViewChoiceSever>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }
    public override void OnInit()
    {

        GetView().InfiniteScroll_Btns.InitListView(0, OnFillBtnItem);

        GetView().InfiniteScroll_Info.InitListView(0, OnFillInfoItem) ;

        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIChoiceSever>();
        });
    }

    private Action callBack;
    private int severCount;


    public override void OnEnter(params object[] args)
    {
        callBack = args[0] as Action;
        severCount = (GameM.Instance.SeverCount + 19) / 20 + 1;

        if (GetView().InfiniteScroll_Btns.ItemTotalCount != severCount)
            GetView().InfiniteScroll_Btns.SetListItemCount(severCount);
        GetView().InfiniteScroll_Btns.RefreshAllShownItem();

        OnSeverBtnCkick(GameM.Instance.SeverInfoDic[0].Count == 0 ? 1 : 0);
    }

    private LoopListViewItem2 OnFillBtnItem (LoopListView2 loopListView, int idx)
    {
        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        var str = "最近登录";
        if (idx != 0)
        {
            int starIdx = (severCount - idx)* 20;
            str = (starIdx - 19) + "服 - " + (starIdx) + "服";
        }
        JUIExtensions.txt(go.transform, "bg/txt").text = str;
        JUIExtensions.txt(go.transform, "on/txt").text = str;

        go.transform.Find("on").gameObject.SetActive(_SelectDicIdx == idx);
        JUIExtensions.btn(go.transform).onClick.RemoveAllListeners();

        JUIExtensions.btn(go.transform).onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnSeverBtnCkick(idx);
        });

        return item;
    }

    private int _SelectDicIdx = 0;

    private void OnSeverBtnCkick(int idx)
    {
        _SelectDicIdx = idx;
        GetView().InfiniteScroll_Btns.RefreshAllShownItem();
        if (idx != 0)
        {
            HttpC.OnGetSeverInfo(new HTTP_SeverRequest { pageIndex = idx, pageSize = 20 }, (_r) =>
            {
                GetView().InfiniteScroll_Info.SetListItemCount(GameM.Instance.SeverInfoDic[_SelectDicIdx].Count);
                GetView().InfiniteScroll_Info.RefreshAllShownItem();
            });
        }
        else
        {
            GetView().InfiniteScroll_Info.SetListItemCount(GameM.Instance.SeverInfoDic[_SelectDicIdx].Count);
            GetView().InfiniteScroll_Info.RefreshAllShownItem();
        }
        
    }

    private LoopListViewItem2 OnFillInfoItem(LoopListView2 loopListView, int idx)
    {
        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        JUIExtensions.txt(go.transform, "txt").text = GameM.Instance.SeverInfoDic[_SelectDicIdx][idx].name;

        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, "img_" + GameM.Instance.SeverInfoDic[_SelectDicIdx][idx].status, JUIExtensions.img(go.transform, "img"));

        JUIExtensions.btn(go.transform).onClick.RemoveAllListeners();
        JUIExtensions.btn(go.transform).onClick.AddListener(() => {

            GameM.Instance.SeverDicIdx = _SelectDicIdx;
            GameM.Instance.SeverIdx = idx;
            callBack?.Invoke();
            UITool.Close<UIChoiceSever>();
        });
        return item;
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}