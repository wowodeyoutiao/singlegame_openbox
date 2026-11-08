using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;
using HotUpdateScripts.ProtoM;
using DG.Tweening;

public class UITips : UIBase<UIModelTips, UIViewTips>
{
    public override UILayer GetLayer()
    {
        return UILayer.Worning;
    }

    public class FlyItem
    {
        public GameObject obj_Item;
        public Text txt;
        public CanvasGroup canGroup;
    }
    private Queue<FlyItem> ObjQueue = new Queue<FlyItem>();

    public override void OnEnter(params object[] args)
    {
       // OnTipShow();
        OnClearTime();
        TimeJAction = GameHelper.Instance.AddInvoke(new TimeItem(DelayTime, DelayTime, OnTipShow,true));
    }

    private TimeItem TimeJAction = null;
    private void OnClearTime()
    { 
        if (GameHelper.Instance != null && TimeJAction != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeJAction);
            TimeJAction = null;
        }
    }

    private FlyItem GetFlyItem()
    {
        FlyItem item;
        if (ObjQueue.Count > 0)
        {
            item = ObjQueue.Dequeue();
            item.obj_Item.transform.SetParent(uiGo.transform);
        }
        else
        {
            var go = GameObject.Instantiate(GetView().obj_item, uiGo.transform);
            go.SetActive(true);
            item = new FlyItem
            {
                obj_Item = go,
                txt = go.transform.Find("txt").GetComponent<Text>(),
                canGroup = go.GetComponent<CanvasGroup>(),
            };
        }
        item.canGroup.alpha = 1;
        return item;
    }

    private void SetFlyItem(FlyItem item)
    {
        item.obj_Item.transform.SetParent(GetView().root_Hide);
        ObjQueue.Enqueue(item);
    }

    private FlyItem EndItem  = null;
    private float FlyTime = 1.5f;
    private float DelayTime = 0.5f;

    public void OnTipShow()
    {
        if (!TipM.Instance.CheckTips() || isNullGo)
        {
            return;
        }
        var FlyItem = GetFlyItem();
        EndItem = FlyItem;
        FlyItem.obj_Item.transform.DOKill();

        FlyItem.obj_Item.transform.localPosition = Vector3.up * 200;
        FlyItem.txt.text = TipM.Instance.GetTips();

        FlyItem.canGroup.DOFade(0, FlyTime).SetEase(Ease.Linear).SetDelay(DelayTime);
        FlyItem.obj_Item.transform.DOLocalMoveY(330, FlyTime).SetEase(Ease.Linear).SetDelay(DelayTime).OnComplete(() =>
        {
            SetFlyItem(FlyItem);

            if (FlyItem == EndItem)
            {
                if (!TipM.Instance.CheckTips())
                    UITool.Close<UITips>();
                //else
                //    OnTipShow();
            }
        });

        //GameHelper.Instance.AddInvoke(new TimeItem(DelayTime, OnTipShow));
    }

    public void SetTip()
    { 
    }

    public override void OnExit()
    {
        OnClearTime();
        EndItem = null;
    }

    public override void OnDestroy()
    {
        OnExit();
        TipM.Instance.ClearTips();
    }
}