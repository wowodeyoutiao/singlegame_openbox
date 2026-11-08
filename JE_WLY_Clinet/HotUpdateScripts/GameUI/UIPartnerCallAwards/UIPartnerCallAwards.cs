using Cysharp.Threading.Tasks;
using DG.Tweening;
using HallBattle;
using HotUpdateScripts;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using Spine.Unity;
using System;
using System.Collections.Generic;
using UnityEngine;
using static CFG_Item;
using static UIViewExpPool;

public class UIPartnerCallAwards : UIBase<UIModelPartnerCallAwards, UIViewPartnerCallAwards>
{
    private UIViewPartnerCallAwards uIViewActivityAwards;
    private TimeItem timeItem;
    private PartnerCallAwardsInfo awardInfo;
    private GameObject goRole = null; //用于存储角色

    private Tween tween;
    private Tween rotateTween;

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
        UITool.Close<UIPartnerCallAwards>();
        awardInfo.callBack?.Invoke();

        //if (uIViewActivityAwards.containerTrans.childCount > 0)
        //    uIViewActivityAwards.Clear();
    }

    public override void OnEnter(params object[] args)
    {
        awardInfo = args[0] as PartnerCallAwardsInfo;
        if (awardInfo == null || awardInfo.awardItems == null || awardInfo.awardItems.Count == 0)
        {
            OnBack();
            return;
        }

        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_CallWhenPartnerArrive);

        Clear();

        uIViewActivityAwards.desText.text = awardInfo.des;

        var dropbagCfg = CFG_DropBag.GetById(awardInfo.awardItems[0].id);

        var itemCfg = CFG_Item.GetById(dropbagCfg.itemId);

        uIViewActivityAwards.txtCardTips.text = $"{GameUtil.Instance.GetItemName(itemCfg.ID)}x{GameUtil.Instance.FormatNumWithUnit(awardInfo.awardItems[0].num)}";
        // 角色
        string poolKey = "Partner_" + dropbagCfg.show;
        var path = GameConst.GetMonsterModelByID(poolKey);

        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_dzk_" + itemCfg.Quality, uIViewActivityAwards.imgCardQuality);

        GameHelper.Instance.SetAsyncObj(poolKey, path, (role) =>
        {
            if (isNullGo)
                return;
            goRole = role;
            role.transform.SetParent(uIViewActivityAwards.imgCardContainer.transform);
            role.name = "SpineItem";
            role.transform.localScale = Vector3.one * 0.32f;
            role.transform.localPosition = new Vector3(0, -26, 0);
            role.SetActive(true);

            var skeletonGraphic = role.GetComponentInChildren<SkeletonGraphic>();

            SpineUtil.Instance.PlayLoop(skeletonGraphic, null);
        }).Forget();

        //var Count = awardInfo.timeCount;
        //timeItem = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () =>
        //{
        //    uIViewActivityAwards.timeText.text = Count + "秒后将自动退出";
        //    Count--;
        //    if (Count <= 0)
        //    {
        //        Clear();
        //        OnBack();
        //    }
        //}));

        ShowLight();
    }

    private void ShowLight()
    {
        SetLightAlpha(0);
        tween = DOTween.To(() => 0, SetLightAlpha, 80, 0.5f).SetEase(Ease.Linear).OnComplete(() =>
        {
            rotateTween = uIViewActivityAwards.LightImg.transform.DOLocalRotate(new UnityEngine.Vector3(0, 0, 15), 5f, RotateMode.FastBeyond360).SetEase(Ease.Linear).SetLoops(-1, LoopType.Incremental);
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

        if (rotateTween != null)
        {
            rotateTween.Kill(); 
            rotateTween = null;
        }

        uIViewActivityAwards.LightImg.transform.localRotation = Quaternion.Euler(Vector3.zero);
        if (goRole != null)
        {
            GameObject.Destroy(goRole);
            goRole = null;
        }
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
public class PartnerCallAwardsItem
{
    public int id;
    public int num;

    public PartnerCallAwardsItem(int id, int num)
    {
        this.id = id;
        this.num = num;
    }
}

public class PartnerCallAwardsInfo
{
    public string des;
    public float timeCount = 5;
    public Action callBack;
    public List<PartnerCallAwardsItem> awardItems = new List<PartnerCallAwardsItem>();

    public PartnerCallAwardsInfo()
    {

    }

    public PartnerCallAwardsInfo(float timeCount)
    {
        this.timeCount = timeCount;
    }

    public PartnerCallAwardsInfo(string des, float timeCount, List<PartnerCallAwardsItem> awardItems, Action callBack = null)
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
        awardItems.Add(new PartnerCallAwardsItem(id, count));
    }

    public void AddAwardItems(List<PartnerCallAwardsItem> awards)
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