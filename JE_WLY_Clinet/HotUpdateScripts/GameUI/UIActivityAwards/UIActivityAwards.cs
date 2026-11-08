using DG.Tweening;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using System;
using System.Collections.Generic;
using UnityEngine;

public class UIActivityAwards : UIBase<UIModelActivityAwards, UIViewActivityAwards>
{
    private UIViewActivityAwards uIViewActivityAwards;
    private TimeItem timeItem;
    private ActivityAwardInfo awardInfo;

    private Tween tween;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        uIViewActivityAwards = GetView();
        uIViewActivityAwards.closeBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            Clear();
            OnBack();
        });
    }

    private void OnBack()
    {
        UITool.Close<UIActivityAwards>();
        awardInfo.callBack?.Invoke();

        if (uIViewActivityAwards.containerTrans.childCount > 0)
            uIViewActivityAwards.Clear();
    }

    public override void OnEnter(params object[] args)
    {
        awardInfo = args[0] as ActivityAwardInfo;
        if (awardInfo == null || awardInfo.awardItems == null || awardInfo.awardItems.Count == 0)
        {
            OnBack();
            return;
        }
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_ComReward);
        Clear();

        uIViewActivityAwards.desText.text = awardInfo.des;

        if (awardInfo.awardItems.Count < 4)
            uIViewActivityAwards.scrollViewTrans.localPosition = new Vector3(0, -130, 0);
        else if (awardInfo.awardItems.Count < 7)
            uIViewActivityAwards.scrollViewTrans.localPosition = new Vector3(0, -70, 0);
        else
            uIViewActivityAwards.scrollViewTrans.localPosition = Vector3.zero;


        for (int i = 0; i < awardInfo.awardItems.Count; i++)
        {
            var go = uIViewActivityAwards.GetItem();
            GameUtil.Instance.SetItemIcon(awardInfo.awardItems[i].id, JUIExtensions.img(go.transform.Find("IconImg")), JUIExtensions.img(go.transform.Find("BgImg")), JUIExtensions.img(go.transform.Find("NatureImg")));
            JUIExtensions.txt(go.transform, "NumText").text = GameUtil.Instance.FormatNumWithUnit(awardInfo.awardItems[i].num);
            go.SetActive(true);
        }

        var Count = awardInfo.timeCount;
        timeItem = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () =>
        {
            uIViewActivityAwards.timeText.text = Count + "秒后将自动退出";
            Count--;
            if (Count <= 0)
            {
                Clear();
                OnBack();
            }
        }));

        ShowLight();
    }

    private void ShowLight()
    {
        SetLightAlpha(0);
        tween = DOTween.To(() => 0, SetLightAlpha, 80, 0.5f).SetEase(Ease.Linear).OnComplete(() => {
            uIViewActivityAwards.LightImg.transform.DOLocalRotate(new UnityEngine.Vector3(0, 0, 15), 5f, RotateMode.FastBeyond360).SetEase(Ease.Linear);
        });
    }

    private void SetLightAlpha(int value)
    {
        var color = uIViewActivityAwards.LightImg.color;
        color.a = value / 255f;
        uIViewActivityAwards.LightImg.color = color;
    }

    private void Clear()
    {
        if (timeItem != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(timeItem);
            timeItem = null;
        }

        if(tween != null)
        {
            tween.Kill();
            tween = null;
        }
        uIViewActivityAwards.LightImg.transform.localRotation = Quaternion.Euler(Vector3.zero);
    }

    public override void OnExit()
    {
        Clear();
    }

    public override void OnDestroy()
    {
        
    }
}

/// <summary>
/// ActivityAwardItem
/// </summary>
public class ActivityAwardItem
{
    public int id;
    public int num;

    public ActivityAwardItem(int id, int num)
    {
        this.id = id;
        this.num = num;
    }
}

public class ActivityAwardInfo
{
    public string des;
    public float timeCount = 5;
    public Action callBack;
    public List<ActivityAwardItem> awardItems = new List<ActivityAwardItem>();

    public ActivityAwardInfo()
    {

    }

    public ActivityAwardInfo(float timeCount)
    {
        this.timeCount = timeCount;
    }

    public ActivityAwardInfo(string des, float timeCount, List<ActivityAwardItem> awardItems, Action callBack = null)
    {
        this.des = des;
        this.awardItems = awardItems;
        this.timeCount = timeCount;
        this.callBack = callBack;
    }

    public void AddAwardItem(string idStr, int count)
    {
        if(int.TryParse(idStr, out var id))
        {
            AddAwardItem(id, count);
        }
    }

    public void AddAwardItem(int id, int count)
    {
        awardItems.Add(new ActivityAwardItem(id, count));
    }

    public void AddAwardItems(List<ActivityAwardItem> awards)
    {
        awardItems.AddRange(awards);
    }

    public void Clear()
    {
        des = string.Empty;
        timeCount = 5;
        callBack = null;
        awardItems.Clear();
    }
}