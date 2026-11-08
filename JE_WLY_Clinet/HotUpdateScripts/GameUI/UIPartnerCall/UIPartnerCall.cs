using Cysharp.Threading.Tasks;
using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using MetaFramework.UI;
using pb;
using Spine.Unity;
using System;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using static CFG_Item;
using static CFG_XMLConfig.XML_Config;

public class UIPartnerCall : UIBase<UIModelPartnerCall, UIViewPartnerCall>
{
    private UIViewPartnerCall partnerCallView;
    private TimeItem timeItem;

    private int totalSliderCount = 0;

    private LotteryPoolInfo lotteryPoolInfo;
    private List<int> hadBuyList = new List<int>();

    private XML_Item_GaCha curCallPriceInfo;

    private Material originMat;

    private string doorBoom = "boom";
    private string doorNor = "men_loop";
    private string doorLow = "men_play";
    private string doorHight = "men_play_max";

    private Dictionary<int, string> spineNameDict = new()
    {
        [2] = "lv",
        [3] = "lan",
        [4] = "zi",
        [6] = "cheng",
        [7] = "hong",
        [8] = "huang",
    };

    private bool callEnable = true;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private void ButtonAddListener(Button btn, Action act)
    {
        btn.onClick.RemoveAllListeners();
        btn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            act?.Invoke();
        });
    }

    public override void OnInit()
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<LotteryInfoResponse>(Sever_LotteryInfoResponse);
        UniGameEvent.AddEventListener<LotteryRefreshResponse>(Sever_LotteryRefreshResponse);
        UniGameEvent.AddEventListener<LotteryResponse>(Sever_LotteryResponse);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);

        partnerCallView = GetView();

        ButtonAddListener(partnerCallView.helpBtn, () => { UITool.Open<UIComRule>(null, 31, "概率公示"); });
        ButtonAddListener(partnerCallView.selBtn, OnClickSelButton);
        ButtonAddListener(partnerCallView.handbookButton, OnClickHandBookButton);
        ButtonAddListener(partnerCallView.bondButton, () => { UITool.Open<UIBond>(); });
        ButtonAddListener(partnerCallView.freeBtn, OnClickCallButton);
        ButtonAddListener(partnerCallView.callBtn, OnClickCallButton);
        ButtonAddListener(partnerCallView.refreshBtn, OnClickRefreshButton);
        ButtonAddListener(partnerCallView.backBtn, () => { UITool.Close<UIPartnerCall>(); });
        ButtonAddListener(partnerCallView.gemstoneAddBtn, () => { UITool.Open<UIShop>(); });
        ButtonAddListener(partnerCallView.currencyAddBtn, () => { UITool.Open<UIGiftBag>(null, ActivityType.PartnerGift); });

        totalSliderCount = CFG_XMLConfig.xmlConfig.GaCha.maxnum;

        // 刷新
        GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.GaCha.refreshItem.useType, partnerCallView.refresPriceIconImg);
        partnerCallView.refresPriceText.text = "x" + CFG_XMLConfig.xmlConfig.GaCha.refreshItem.num;

        // 羁绊
        var redDot = partnerCallView.bondButton.transform.Find("RedPointImage").gameObject;
        RedPointManager.Instance.RPS.Register(RedPointConst.Role_Bond.ToString(), (node) =>
        {
            if (isNullGo) return;
            redDot.SetActive(BondM.Instance.GetRedNum() > 0);
        });
    }

    // 心愿自选
    private void OnClickSelButton()
    {
        UITool.Open<UIPartnerCallSelect>(null, lotteryPoolInfo.NeedItemId, (int selId) =>
        {
            lotteryPoolInfo.NeedItemId = selId;
        });
    }

    private void OnClickHandBookButton()
    {
        UITool.Open<UIPartnerBook>();
    }

    // 免费、召唤
    private void OnClickCallButton()
    {
        if (curCallPriceInfo == null)
            return;

        // 判断道具够不够
        if (curCallPriceInfo.useType != 0 && !GameUtil.Instance.ComCoinToCheck(curCallPriceInfo.useType, curCallPriceInfo.num, null, NeedJump: false))
        {
            if(GameM.Instance.IsAudit)
            {
                UITool.ShowTips("道具不足");
                return;
            }

            if (curCallPriceInfo.useType == (int)CFG_Item.ItemID.DianQuan)
                UITool.Open<UIShop>();
            else 
                UITool.Open<UIGiftBag>(null, ActivityType.PartnerGift);
            return;
        }

        if (!callEnable)
            return;

        callEnable = false;

        FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.CallRole });
        TcpC.C2S_Lottery();
    }

    // 刷新
    private void OnClickRefreshButton()
    {
        // 判断道具够不够
        if (!GameUtil.Instance.ComCoinToCheck(CFG_XMLConfig.xmlConfig.GaCha.refreshItem.useType, CFG_XMLConfig.xmlConfig.GaCha.refreshItem.num, null))
        {
            UITool.ShowTips("道具不足");
            return;
        }

        var key = RoleM.Instance.Uid + "UIPartnerCall";
        string timeStr = GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear.ToString();
        if (PlayerPrefsExtend.GetString(key) == timeStr)
        {
            RefreshLotteryInfo(false);
        }
        else
        {
            var priceCfg = CFG_Item.GetById(CFG_XMLConfig.xmlConfig.GaCha.refreshItem.useType);
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                isTogg = true,
                ShowStr = $"是否花费<color=#2C7E40>{CFG_XMLConfig.xmlConfig.GaCha.refreshItem.num}{priceCfg.Name}</color>刷新？",
                TogNoFun = () => {
                    PlayerPrefsExtend.SetString(key, "");
                },
                TogYesFun = () => {
                    PlayerPrefsExtend.SetString(key, timeStr);
                },
                YesFun = () =>
                {
                    RefreshLotteryInfo(false);
                },
            });
        }
    }

    #region timer
    private void AddTimer(long endTime)
    {
        ClearTimer();
        //Debug.LogError($"{endTime}, {GameUtil.Instance.GetCurrentTimeSeconds()}");
        endTime = endTime - GameUtil.Instance.GetCurrentTimeSeconds() + 30 * 60;
        if(endTime < 0)
            endTime = 30 * 60;
        SetTimeText(endTime);
        timeItem = GameHelper.Instance.AddInvoke(new TimeItem(endTime, 1, () =>
        {
            if(isNullGo)
            {
                ClearTimer();
                return;
            }

            endTime--;
            SetTimeText(endTime);
            if (endTime <= 1)
            {
                endTime = 30 * 60;
                ClearTimer();
                RefreshLotteryInfo(true);
            }
        }));
    }

    private void SetTimeText(long endTime)
    {
        if (endTime <= 0)
            partnerCallView.timerText.text = "刷新时间：00:00";
        else
        {
            var timeStr = GameUtil.Instance.GetTimeSToDHOrHMS(endTime);
            partnerCallView.timerText.text = $"刷新时间：<color=#fffd44>{timeStr}</color>";
        }
    }

    private void ClearTimer()
    {
        if (timeItem != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(timeItem);
            timeItem = null;
        }
    }
    #endregion

    private void RefreshCurrency()
    {
        partnerCallView.gemstoneText.text = GameUtil.Instance.FormatNumWithUnit(RoleM.Instance.Coupon);
        partnerCallView.currencyText.text = GameUtil.Instance.FormatNumWithUnit(RoleM.Instance.StarWishCoin);
    }

    public override void OnEnter(params object[] args)
    {
        ResetItem();
        // 获取列表
        TcpC.C2S_LotteryInfo();
    }

    #region 查询抽奖内容
    private void Sever_LotteryInfoResponse(LotteryInfoResponse resq)
    {
        if (resq.Code != 0)
        {
            UITool.ShowTips(resq.Error);
            return;
        }

        lotteryPoolInfo = resq.LotteryInfos;

        RefreshUI(false);
    }
    #endregion

    #region 抽奖

    private void Sever_LotteryResponse(LotteryResponse resq)
    {
        if (resq.Code != 0)
        {
            UITool.ShowTips(resq.Error);
            return;
        }
        // 修改FreeGet
        lotteryPoolInfo.FreeNum = 0;
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Call_Free.ToString(), 0);

        var count = 0;
        if(lotteryPoolInfo.HavePrizePool != null && lotteryPoolInfo.HavePrizePool.Count > 0)
        {
            for(int i = 0; i < lotteryPoolInfo.HavePrizePool.Count; i++)
            {
                if (lotteryPoolInfo.HavePrizePool[i].Id == resq.Id)
                    count++;
            }
        }

        // 修改数据，添加到已有奖励池
        var index = 0;
        for (int i = 0; i < lotteryPoolInfo.PrizePool.Count; i++)
        {
            if (lotteryPoolInfo.PrizePool[i].Id == resq.Id)
            {
                if(count > 0)
                {
                    count--;
                    continue;
                }

                index = i;
                lotteryPoolInfo.HavePrizePool.Add(lotteryPoolInfo.PrizePool[i]);
                break;
            }
        }
        lotteryPoolInfo.Progress += curCallPriceInfo.luck;


        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_CallEffect);
        // 显示选中特效
        var itemTrans = partnerCallView.itemsTrans[index];
        partnerCallView.callBoomEffectGo.transform.SetParent(itemTrans);
        partnerCallView.callBoomEffectGo.transform.localPosition = new Vector3(0, -50, 0);
        partnerCallView.callBoomEffectGo.transform.localScale = Vector3.one;
        partnerCallView.callBoomEffectGo.SetActive(true);

        SpineUtil.Instance.PlayOnce(partnerCallView.callBoomEffectGo.transform, "boom", () =>
        {
            partnerCallView.callBoomEffectGo.SetActive(false);
            partnerCallView.callBoomEffectGo.transform.SetParent(partnerCallView.transform);
        });

        var reload = lotteryPoolInfo.HavePrizePool.Count == 4;

        GameHelper.Instance.SetDeayTimeDo(() =>
        {
            if (isNullGo || resq == null)
                return;

            if(reload)
                RefreshLotteryInfo(true);

            RefreshUI(false);

            var info = new PartnerCallAwardsInfo();
            info.AddAwardItem(resq.Id, resq.ItemNum);
            info.des = string.Empty;
            UITool.Open<UIPartnerCallAwards>(null, info);
        }, 1f);
    }
    #endregion

    #region 主动刷新抽奖池
    private void RefreshLotteryInfo(bool isFree)
    {
        ResetItem();
        if(isFree)
            TcpC.C2S_LotteryInfo();
        else
            // 主动刷新抽奖池
            TcpC.C2S_LotteryRefresh();
    }

    private void Sever_LotteryRefreshResponse(LotteryRefreshResponse resq)
    {
        if (resq.Code != 0)
        {
            UITool.ShowTips(resq.Error);
            return;
        }

        lotteryPoolInfo = resq.LotteryInfos;
        ShowBirthEffect();
    }
    #endregion

    #region UI刷新
    private void ShowBirthEffect()
    {
        var flyNames = new List<int>();
        var hasRed = false;
        for (int i = 0; i < lotteryPoolInfo.PrizePool.Count; i++)
        {
            var itemCfg = CFG_Item.GetById(lotteryPoolInfo.PrizePool[i].ItemId);
            if (itemCfg.Quality >= 7)
                hasRed = true;

            flyNames.Add(itemCfg.Quality);
        }

        var spineName = doorLow;
        if (hasRed)
            spineName = doorHight;

        SpineUtil.Instance.PlayOnce(partnerCallView.doorSpineTrans, spineName, () =>
        {
            SpineUtil.Instance.PlayLoop(partnerCallView.doorSpineTrans, doorNor);
        });

        DOTween.To(() => 0, (value) => { }, 1, 1f).SetEase(Ease.Linear).OnComplete(() =>
        {
            if (isNullGo)
                return;

            // 飞球
            for (int i = 0; i < flyNames.Count; i++)
                Fly(i, flyNames[i]).Forget();
        });
    }

    private async UniTask Fly(int index, int quality)
    {
        string uiPoolKey = "Eff_Ball" + quality;

        if (!PoolHelper.Instance.CheckPool(uiPoolKey))
        {
            var baseObj = new JPrefab(GameConst.GetEffPath(uiPoolKey), true);
            await UniTask.WaitUntil(() => baseObj.Loaded);

            PoolHelper.Instance.AddPool(uiPoolKey, baseObj.Instance);
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(uiPoolKey));
        var effGo = PoolHelper.Instance.GetPoolObj(uiPoolKey);
        effGo.transform.SetParent(partnerCallView.transform);
        effGo.transform.localPosition = partnerCallView.doorSpineTrans.localPosition;

        var meshRenderers = effGo.transform.GetComponentsInChildren<Renderer>();
        if(meshRenderers != null && meshRenderers.Length > 0)
        {
            for (int i = 0; i < meshRenderers.Length; i++)
            {
                meshRenderers[i].sortingOrder = ((int)UILayer.Normal - 1) * 200 + 2;
            }
        }

        effGo.transform.localScale = Vector3.one * 0.1f;
        effGo.SetActive(true);

        var itemTrans = partnerCallView.itemsTrans[index];
        var targetPos = itemTrans.localPosition + itemTrans.parent.localPosition + itemTrans.Find("BgImg").localPosition;

        var seq = DOTween.Sequence();
        seq.Append(effGo.transform.DOScale(Vector3.one, 0.05f));
        seq.Append(effGo.transform.DOLocalMove(targetPos, 0.2f));
        seq.Append(effGo.transform.DOScale(Vector3.one * 0.1f, 0.05f));
        seq.AppendCallback(() =>
        {
            seq.Kill();
            seq = null;
            PoolHelper.Instance.RecycObj(effGo);
            if (index == 3)
                RefreshUI(true);
        });
    }

    private void RefreshUI(bool isRefresh)
    {
        callEnable = true;
        RefreshCurrency();
        AddTimer(lotteryPoolInfo.RefreshTime);
        RefreshSlider(lotteryPoolInfo.Progress);
        RefreshCallBtn(lotteryPoolInfo.FreeNum == 1, lotteryPoolInfo.HavePrizePool == null ? 0 : lotteryPoolInfo.HavePrizePool.Count);

        if (lotteryPoolInfo.HavePrizePool != null && lotteryPoolInfo.HavePrizePool.Count > 0)
        {
            for (int i = 0; i < lotteryPoolInfo.HavePrizePool.Count; i++)
            {
                hadBuyList.Add(lotteryPoolInfo.HavePrizePool[i].Id);
            }
        }

        RefreshItems(lotteryPoolInfo.PrizePool, lotteryPoolInfo.HavePrizePool);
    }

    private void RefreshSlider(int progress)
    {
        partnerCallView.sliderImg.fillAmount = progress * 1.0f / totalSliderCount;
        partnerCallView.sliderText.text = $"{progress} / {totalSliderCount}";

        partnerCallView.refreshBtnEffectGo.SetActive(progress >= totalSliderCount);
        partnerCallView.sliderSpineGo.SetActive(progress >= totalSliderCount);

        if (progress >= totalSliderCount)
            UITool.ShowTips("幸运值已满将不再增加幸运值！刷新必出红色伙伴！");
    }

    // 重置角色UI
    private void ResetItem()
    {
        for (int i = 0; i < partnerCallView.itemsTrans.Count; i++)
        {
            var itemTrans = partnerCallView.itemsTrans[i];

            var bgImg = itemTrans.Find("BgImg").GetComponent<Image>();
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_dz", bgImg);

            var spineTrans = itemTrans.Find("SpineItem");
            if (GameHelper.Instance != null && !isNullGo)
                GameHelper.Instance.ClearObjToPool(spineTrans);

            var iconImgGo = itemTrans.Find("IconImg").gameObject;
            iconImgGo.SetActive(false);

            itemTrans.Find("NameBgImg").gameObject.SetActive(false);
            itemTrans.Find("EffectSpine").gameObject.SetActive(false);
            itemTrans.Find("NewImage").gameObject.SetActive(false);
        }
    }

    // 刷新角色UI
    private void RefreshItems(List<PrizePool> prizePools, List<PrizePool> havePrizePools)
    {
        var haveIds = new List<int>();
        if (havePrizePools != null && havePrizePools.Count > 0)
        {
            for (int j = 0; j < havePrizePools.Count; j++)
                haveIds.Add(havePrizePools[j].Id);
        }

        for (int i = 0; i < prizePools.Count; i++)
        {
            var isGeted = false;
            if (haveIds != null && haveIds.Count > 0)
            {
                for (int j = 0; j < haveIds.Count; j++)
                {
                    if (prizePools[i].Id == haveIds[j])
                    {
                        isGeted = true;
                        haveIds.RemoveAt(j);
                        break;
                    }
                }
            }
            ShowItem(i, prizePools[i].Id, prizePools[i].ItemId, prizePools[i].Num, isGeted);
        }
    }

    private void ShowItem(int index, int dropbagId, int itemId, int num, bool isGeted)
    {
       // Debug.LogError($"{index}, {dropbagId}, {itemId}, {num}, {isGeted}");
        var itemCfg = CFG_Item.GetById(itemId);
        var dropbagCfg = CFG_DropBag.GetById(dropbagId);
        var itemTrans = partnerCallView.itemsTrans[index];

        if (itemCfg == null || itemTrans == null) return;

        var iconImgGo = itemTrans.Find("IconImg").gameObject;
        iconImgGo.SetActive(false);
        var bgImg = itemTrans.Find("BgImg").GetComponent<Image>();
        if (isGeted)
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_dz", bgImg);
        else
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_dz_" + itemCfg.Quality, bgImg);

        itemTrans.Find("NameBgImg").gameObject.SetActive(true);
        var nameText = itemTrans.Find("NameBgImg/NameText").GetComponent<Text>();
        nameText.text = $"{GameUtil.Instance.GetStrByQuality(itemCfg.Quality, itemCfg.Name)}x{num}";

        var newImg = itemTrans.Find("NewImage").gameObject;
        newImg.SetActive(false);
        if(!isGeted && itemCfg.Type == 8 && dropbagCfg.show > 0)
        {
            var partnerInfo = PartnerM.Instance.GetOnePartnerInfoByCFGID(dropbagCfg.show);
            itemTrans.Find("NewImage").gameObject.SetActive(partnerInfo == null);
        }

        var spineTrans = itemTrans.Find("SpineItem");
        if (GameHelper.Instance != null && !isNullGo)
            GameHelper.Instance.ClearObjToPool(spineTrans);

        if (dropbagCfg.show > 0)
        {
            string poolKey = "Partner_" + dropbagCfg.show;
            var path = GameConst.GetMonsterModelByID(poolKey);
            GameHelper.Instance.SetAsyncObj(poolKey, path, (role) =>
            {
                if (isNullGo)
                    return;

                role.transform.SetParent(spineTrans);
                role.name = "SpineItem";
                role.transform.localScale = Vector3.one * 0.6f;
                role.transform.localPosition = Vector3.zero;
                role.SetActive(true);

                var skeletonGraphic = role.GetComponentInChildren<SkeletonGraphic>();
                if (originMat == null)
                    originMat = skeletonGraphic.material;

                skeletonGraphic.material = (isGeted ? partnerCallView.grayImg.material : originMat);
                SpineUtil.Instance.PlayLoop(skeletonGraphic.transform, isGeted ? null : "stand");
            }).Forget();
        }

        // 底框动画
        var effectSpineTrans = itemTrans.Find("EffectSpine");
        effectSpineTrans.gameObject.SetActive(!isGeted);
        if (!isGeted)
            SpineUtil.Instance.PlayLoop(effectSpineTrans, spineNameDict[itemCfg.Quality]);
    }

    private void RefreshCallBtn(bool isFirstBuy, int buyCount)
    {
        partnerCallView.freeBtn.gameObject.SetActive(isFirstBuy);
        partnerCallView.callBtn.gameObject.SetActive(!isFirstBuy);
        if (!isFirstBuy)
        {
            if(CFG_XMLConfig.xmlConfig.GaCha.outputItemDict.TryGetValue(buyCount + 1, out var xml))
            {
                curCallPriceInfo = xml;
                GameUtil.Instance.SetItemIcon(curCallPriceInfo.useType, partnerCallView.callPriceIconImg);
                partnerCallView.callPriceText.text = "x" + curCallPriceInfo.num;
            }
        }
        else
            curCallPriceInfo = CFG_XMLConfig.xmlConfig.GaCha.outputItemDict[0];
    }
    #endregion

    [Subscriber(ThreadMode.Main)]
    private void OnSync_TakeSource(GameEvent.Sync_TakeSource _e)
    {
        if (isNullGo)
            return;
        RefreshCurrency();
    }


    private void OnEvent_Source_Refresh(Event_Source_Refresh _e)
    {
        if (isNullGo)
            return;
        RefreshCurrency();
    }
    public override void OnExit()
    {
        if (ActivityM.Instance.GetAcitvityEndTime(ActivityType.CallFund) > 0)
            HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.CallFund);
    }

    public override void OnDestroy()
    {
        OnExit();
        ClearTimer();
        partnerCallView = null;
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<LotteryInfoResponse>(Sever_LotteryInfoResponse);
        UniGameEvent.RemoveEventListener<LotteryResponse>(Sever_LotteryResponse);
        UniGameEvent.RemoveEventListener<LotteryRefreshResponse>(Sever_LotteryRefreshResponse);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);
    }
}