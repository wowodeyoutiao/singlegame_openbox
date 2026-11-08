using Cysharp.Threading.Tasks;

using DG.Tweening;

using HotUpdateScripts.ProtoM;

using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using SuperScrollView;

using System;
using System.Collections.Generic;
using System.Linq;

using UnityEngine;
using UnityEngine.UI;

public class UIActive : UIBase<UIModelActive, UIViewActive>
{
    private string curDateTime;
    private bool isHappyTab;
    private int curHappyTab;
    private UIViewActive _uiView;
    private List<CFGItem_Activity> showActivityList;
    private TimeItem timeItem;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnLayerChange()
    {
        GameUtil.Instance.SetEffectLayer(GetView().flyEffectGo.transform, _canvasOrder + 1);
    }

    public override void OnInit()
    {
        _uiView = GetView();

        _uiView.closeBtn.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIActive>();
        });

        // tab
        for (int i = 0; i < GetModel().ActiveTitlesDic.Count; i++)
        {
            var titleData = GetModel().ActiveTitlesDic.ElementAt(i).Value;

            //_uiView.tabBtns[i].transform.Find("Buttom/Label").GetComponent<Text>().text = titleData.activeName;

            _uiView.tabBtns[i].gameObject.name = titleData.activeType.ToString();

            _uiView.tabBtns[i].gameObject.SetActive(titleData.ActiveOpen);
            _uiView.tabBtns[i].onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (titleData.activeType != activityType)
                    OnClickTabButton(titleData);
            });
        }

        OninitSign();
        OninitSignHappy();
        OnInitStarWish();
        OnInitFunc();
    }

    #region red Point
    private void BindRed()
    {
        for (int i = 0; i < GetModel().ActiveTitlesDic.Count; i++)
        {
            var titleData = GetModel().ActiveTitlesDic.ElementAt(i).Value;

            if (titleData.ActiveOpen)
            {
                int index = i;
                if (!string.IsNullOrEmpty(titleData.RedKey))
                {
                    RedPointManager.Instance.RPS.Register(titleData.RedKey, (node) =>
                    {
                        _uiView.tabBtns[index].transform.Find("Dot").gameObject.SetActive(node.pointCount > 0);
                    });
                }
            }
        }
    }

    private void ClearBindRed()
    {
        foreach (var item in GetModel().ActiveTitlesDic.Values)
        {
            if (!item.ActiveOpen || string.IsNullOrEmpty(item.RedKey)) continue;
            RedPointManager.Instance.RPS.Register(item.RedKey, null);
        }
    }

    #endregion

    private bool CheckActivityEnd(ActivityType activityType)
    {
        if (ActivityM.Instance.GetAcitvityEndTime(activityType) <= 0)
        {
            var info = new UIMessageBoxOne.CallBack
            {
                showStr = "活动已结束！",
                cType = 0,
                endFun = () =>
                {
                    JEvent.defaultEvent.Post(new GameEvent.Refresh_HallActivityButton() { index = 1, forceHide = true });
                    UITool.Close<UIActive>();
                },
            };
            UITool.Open<UIMessageBoxOne>(null, info);
            return true;
        }

        return false;
    }

    ActivityType activityType;

    private void OnClickTabButton(UIModelActive.ActiveTitle index)
    {
        activityType = index.activeType;

        for (int i = 0; i < _uiView.tabBtns.Length; i++)
        {
            _uiView.tabBtns[i].transform.Find("Select").gameObject.SetActive(_uiView.tabBtns[i].gameObject.name == index.activeType.ToString());
        }

        _uiView.sign7Go.SetActive(index.activeType == ActivityType.Sign7);
        _uiView.happy7Go.SetActive(index.activeType == ActivityType.SignHappy7);

        _uiView.starWishSummonTrans.gameObject.SetActive(index.activeType == ActivityType.StarWish
            || index.activeType == ActivityType.CrazyRefine);

        _uiView.rootFund.SetActive(index.activeType == ActivityType.RefineFund
            || index.activeType == ActivityType.CallFund
            || index.activeType == ActivityType.CopyFund);

        Clear();
        if (index.cfgId > 0)
        {
            showActivityList = CFG_Activity.GetActiveListByType((CFG_Activity.Type)index.cfgId);
            if (index.activeType == ActivityType.Sign7)
                _uiView.signDesText.text = $"活动说明：{CFG_Activity.GetActiveTypeDic()[index.cfgId].RuleDes}";
            //else if (index.activeType == ActivityType.StarWish || index.activeType == ActivityType.CrazyRefine)
                //_uiView.txt_StarWishDesc.text = $"活动说明：{CFG_Activity.GetActiveTypeDic()[index.cfgId].RuleDes}";
            else if (index.activeType == ActivityType.RefineFund)
                storeId = 38001;
            else if (index.activeType == ActivityType.CallFund)
                storeId = 39001;
            else if (index.activeType == ActivityType.CopyFund)
                storeId = 40001;

            AddTimer();
        }
        else
        {
            if (index.activeType == ActivityType.Sign7)
                _uiView.signDesText.text = "";
            //else if (index.activeType == ActivityType.StarWish || index.activeType == ActivityType.CrazyRefine)
            //    _uiView.txt_StarWishDesc.text = "";
            else if (index.activeType == ActivityType.RefineFund)
                storeId = 38001;
            else if (index.activeType == ActivityType.CallFund)
                storeId = 39001;
            else if (index.activeType == ActivityType.CopyFund)
                storeId = 40001;

            _uiView.activeTimeText.text = $"永久存在";
        }

        isHappyTab = index.activeType == ActivityType.SignHappy7;
        GetView().flyEffectGo.SetActive(index.activeType == ActivityType.Sign7 || index.activeType == ActivityType.SignHappy7);

        switch (index.activeType)
        {
            case ActivityType.Sign7:
                OnSign7Show();
                break;
            case ActivityType.SignHappy7:
                OnSignHappyShow();
                break;
            case ActivityType.StarWish:
                if (!CheckActivityEnd(index.activeType))
                    OnStarWishShow();
                break;
            case ActivityType.CrazyRefine:
                if (!CheckActivityEnd(index.activeType))
                    OnStarWishShow();
                break;
            case ActivityType.RefineFund:
                _uiView.txt_FundTitle.text = "洗炼次数";
                ResouseHelper.Instance.LoadSprite($"Common/Atlas/Activity/bg_Seven_days_banner_6", _uiView.img_FundLogo).Forget();
                fundActivityType = (CFG_Activity.Type)index.cfgId;
                if (activityAwardInfo != null)
                    activityAwardInfo.Clear();
                if (!CheckActivityEnd(index.activeType))
                {
                    fundPrice = CFG_Store.GetById(storeId).Limite.Num;
                    ShowFund();
                }
                break;
            case ActivityType.CallFund:
                _uiView.txt_FundTitle.text = "召唤次数";
                ResouseHelper.Instance.LoadSprite($"Common/Atlas/Activity/bg_Seven_days_banner_5", _uiView.img_FundLogo).Forget();
                fundActivityType = (CFG_Activity.Type)index.cfgId;
                if (activityAwardInfo != null)
                    activityAwardInfo.Clear();
                if (!CheckActivityEnd(index.activeType))
                {
                    fundPrice = CFG_Store.GetById(storeId).Limite.Num;
                    ShowFund();
                }
                break;
            case ActivityType.CopyFund:
                ResouseHelper.Instance.LoadSprite($"Common/Atlas/Activity/bg_Seven_days_banner_7", _uiView.img_FundLogo).Forget();
                _uiView.txt_FundTitle.text = "参与次数";
                fundActivityType = (CFG_Activity.Type)index.cfgId;
                if (activityAwardInfo != null)
                    activityAwardInfo.Clear();
                if (!CheckActivityEnd(index.activeType))
                {
                    fundPrice = CFG_Store.GetById(storeId).Limite.Num;
                    ShowFund();
                }
                break;
            default:
                break;
        }
    }

    #region 七日签到
    private void OninitSign()
    {
        if (!GetModel().ActiveTitlesDic[ActivityType.Sign7].ActiveOpen)
            return;

        _uiView.infiniteScroll_Sign7.InitListView(0, OnSign7Fill);
    }

    private LoopListViewItem2 OnSign7Fill(LoopListView2 listView, int idx)
    {
        if (idx < 0 || idx >= showActivityList.Count)
        {
            return null;
        }
        var info = showActivityList[idx];
        var item = listView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        var obj = go.transform.Find("Awards");
        JUIExtensions.txt(go.transform.Find("txt_day")).text = $"第{info.Content}天";
        for (int i = 0; i < obj.childCount; i++)
        {
            obj.GetChild(i).gameObject.SetActive(i < info.CommonItem.Length);

            if (i < info.CommonItem.Length)
            {
                int cfgID = int.Parse(info.CommonItem[i].ElementAt(0).Key);
                GameUtil.Instance.SetItemIcon(cfgID, JUIExtensions.img(obj.GetChild(i).transform.Find("icon")), JUIExtensions.img(obj.GetChild(i)));
                JUIExtensions.txt(obj.GetChild(i).transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(info.CommonItem[i].ElementAt(0).Value);

                JUIExtensions.btn(obj.GetChild(i)).onClick.RemoveAllListeners();
                JUIExtensions.btn(obj.GetChild(i)).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgID });
                });
            }
        }
        var _State = (ActivityM.SignState)ActivityM.Instance.GetSignStateByDay(info.Content).status;

        JUIExtensions.txt(go.transform.Find("txt_GetState")).text = _State == ActivityM.SignState.NoGet ? "时间未到" : "";

        go.transform.Find("Geted").gameObject.SetActive(_State == ActivityM.SignState.Geted);

        JUIExtensions.btn(go.transform.Find("Btn_Get")).gameObject.SetActive(_State == ActivityM.SignState.CanGet);
        JUIExtensions.btn(go.transform.Find("Btn_Get")).onClick.RemoveAllListeners();
        JUIExtensions.btn(go.transform.Find("Btn_Get")).onClick.AddListener(delegate { OnSignBtnClick(info, false); });

        var compensationConsumeId = int.Parse(info.CompensationConsume[0].ElementAt(0).Key);
        JUIExtensions.btn(go.transform.Find("Btn_GetPass")).gameObject.SetActive(_State == ActivityM.SignState.PayGet);
        JUIExtensions.btn(go.transform.Find("Btn_GetPass")).onClick.RemoveAllListeners();
        JUIExtensions.btn(go.transform.Find("Btn_GetPass")).onClick.AddListener(delegate { OnSignPassBtnClick(idx, info, compensationConsumeId == 51002); });

        GameUtil.Instance.SetItemIcon(compensationConsumeId, JUIExtensions.img(go.transform.Find("Btn_GetPass/icon")));
        if (compensationConsumeId == 51002)
            JUIExtensions.txt(go.transform.Find("Btn_GetPass/txt")).text = "";
        else
            JUIExtensions.txt(go.transform.Find("Btn_GetPass/txt")).text = GameUtil.Instance.FormatNumWithUnit(info.CompensationConsume[0].ElementAt(0).Value);

        return item;
    }

    public void OnSignBtnClick(CFGItem_Activity info, bool useVideo)
    {
        if (CheckActivityEnd(ActivityType.Sign7))
            return;

        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        HttpC.OnSign7Get(info.UUID, useVideo, (_response) =>
        {
            if (_response.code == 0 && !isNullGo)
            {
                OnSign7Show();

                var awardInfo = new ActivityAwardInfo();
                for (int i = 0; i < info.CommonItem.Length; i++)
                {
                    awardInfo.AddAwardItem(info.CommonItem[i].ElementAt(0).Key, info.CommonItem[i].ElementAt(0).Value);
                    int itemId = int.Parse(info.CommonItem[i].ElementAt(0).Key);
                    if (CFG_Item.GetById(itemId).Type == (int)CFG_Item.ItemType.Pet)
                        UITool.ShowTips($"获得{GameUtil.Instance.GetItemName(itemId)}");
                }
                awardInfo.des = string.Empty;
                UITool.Open<UIActivityAwards>(null, awardInfo);
            }
            else
            {
                UITool.ShowTips(_response.data);
            }
        });
    }

    public void OnSignPassBtnClick(int idx, CFGItem_Activity info, bool useVideo)
    {
        if (CheckActivityEnd(ActivityType.Sign7))
            return;

        if (!useVideo) // 使用非广告补签时
        {
            var used = showActivityList[idx].CompensationConsume[0].ElementAt(0);
            if (!GameUtil.Instance.ComCoinToCheck(int.Parse(used.Key), used.Value, delegate { UITool.Close<UIActive>(); }))
            {
                return;
            }
            OnSignBtnClick(info, useVideo);
        }
        else
        {
            if (SDKCtr.Instance.CurrentPlatform == SDKCtr.PlatformType.WebGL_WX && !GameM.Instance.SkipAD)
            {
                HttpC.ShuShuSendAdShow(new HTTP_WindowShushuAdShowRequest
                {
                    ad_type = "激励广告",
                    ad_position = "51002",
                });

                SDKCtr.Instance.ShowRewardedVideo(() =>
                {
                    OnSignBtnClick(info, useVideo);
                    AdClose(true);
                },
                () =>
                {
                    UITool.ShowTips("广告观看取消");
                    AdClose(false);
                });

                void AdClose(bool completed)
                {
                    HttpC.ShuShuSendAdClose(new HTTP_WindowShushuAdCloseRequest
                    {
                        ad_type = "激励广告",
                        ad_position = "51002",
                        view_complete = completed,
                    });
                }
            }
            else
                OnSignBtnClick(info, useVideo);
        }
    }
    private void OnSign7Show()
    {
        UITool.Open<UILoading>((ui) =>
        {
            HttpC.OnSign7Info((_response) =>
            {
                UITool.Close<UILoading>();
                if (_response.code == 0 && !isNullGo)
                {
                    if (GetView().infiniteScroll_Sign7.ItemTotalCount != showActivityList.Count)
                        GetView().infiniteScroll_Sign7.SetListItemCount(showActivityList.Count);
                    GetView().infiniteScroll_Sign7.RefreshAllShownItem();
                    GetView().infiniteScroll_Sign7.MovePanelToItemIndex(OnSign7MovePosition(), 0);
                }
            }, () => UITool.Close<UILoading>());
        }, new UILoading.Loading { _time = 0 });

    }

    private int OnSign7MovePosition()
    {
        int position = 0;

        for (int i = 0; i < showActivityList.Count; i++)
        {
            var state = (ActivityM.SignState)ActivityM.Instance.GetSignStateByDay(showActivityList[i].Content).status;
            if (state == ActivityM.SignState.CanGet)
            {
                position = i;
                break;
            }
            else if (state == ActivityM.SignState.Geted)
            {
                position = i;
            }
        }

        return position;
    }
    #endregion

    #region 七日狂欢
    private int curDay;
    private List<TaskInfo> allHappy7TaskInfos;
    private List<TaskInfo> happy7TaskInfos = new List<TaskInfo>();

    private Coroutine happyAwardCoroutine;

    private void OninitSignHappy()
    {
        for (int i = 0; i < _uiView.happy7TabBtns.Length; i++)
        {
            var index = i;
            _uiView.happy7TabBtns[i].onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnClickHappyTabButton(index);
            });
        }

        _uiView.infiniteScroll_happy7.InitListView(0, OnSignHappy7Fill);
    }

    private void ShowDayRedPoint(int index)
    {
        _uiView.happy7TabBtns[index].transform.GetChild(2).gameObject.SetActive(true);
    }

    private void RefreshHappyTabButtonState()
    {
        for (int i = 0; i < _uiView.happy7TabBtns.Length; i++)
        {
            _uiView.happy7TabBtns[i].transform.GetChild(3).gameObject.SetActive(i > curDay - 1);
        }
    }

    private void OnClickHappyTabButton(int index)
    {
        if (CheckActivityEnd(ActivityType.SignHappy7))
            return;

        if (curDay < index + 1)
        {
            UITool.ShowTips("活动开启时间未到，请明天再试！");
            return;
        }

        for (int i = 0; i < _uiView.happy7TabBtns.Length; i++)
        {
            _uiView.happy7TabBtns[i].transform.GetChild(2).gameObject.SetActive(false);
            _uiView.happy7TabBtns[i].transform.GetChild(1).gameObject.SetActive(index == i);
        }

        ShowHappyDailyGift(index);

        allHappy7TaskInfos = TaskM.Instance.GetTasksByClass(TaskClass.SEVENDAY);

        happy7TaskInfos.Clear();
        if (allHappy7TaskInfos != null && allHappy7TaskInfos.Count > 0)
        {
            for (int i = 0; i < allHappy7TaskInfos.Count; i++)
            {
                TaskInfo taskInfo = allHappy7TaskInfos[i];
                var cfgInfo = CFG_Task.GetById(taskInfo.CfgId);
                if (cfgInfo.SevenDay == index + 1)
                    happy7TaskInfos.Add(taskInfo);

                if (taskInfo.Status == TaskStatus.COMPLETED)
                    ShowDayRedPoint((int)cfgInfo.SevenDay - 1);

                if (cfgInfo.SevenDay > curDay)
                    curDay = (int)cfgInfo.SevenDay;
            }
        }
        RefreshHappyTabButtonState();
        curHappyTab = index;
        _uiView.infiniteScroll_happy7.SetListItemCount(happy7TaskInfos.Count);
        _uiView.infiniteScroll_happy7.RefreshAllShownItem();
    }

    private void OnSignHappyShow()
    {
        OnClickHappyTabButton(curHappyTab);

        // 大奖
        SignHappyBigAwardShow();
    }

    private void SignHappyBigAwardShow()
    {
        var bigAwards = TaskM.Instance.GetTasksByClass(TaskClass.SEVENDAY_ACTIVITY);
        var taskInfo = bigAwards[0];
        var taskCfg = CFG_Task.GetById(taskInfo.CfgId);
        (int amout, int targetnum) = CFG_Task.GetTaskPro(taskInfo.CfgId, taskInfo.Target.Amount, taskInfo.Target.TargetNum);
        _uiView.happy7Slider.value = amout * 1.0f / targetnum;
        _uiView.happy7SchduleText.text = $"{amout}/{targetnum}";

        var award = _uiView.happy7SliderAwardTrans;
        int CfgId = int.Parse(taskCfg.Reward[0].ElementAt(0).Key);
        GameUtil.Instance.SetItemIcon(CfgId, JUIExtensions.img(award.Find("icon")), JUIExtensions.img(award));
        JUIExtensions.txt(award.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(taskCfg.Reward[0].ElementAt(0).Value);

        var spineGo = award.Find("spine").gameObject;
        var ani = award.GetComponent<Animation>();

        award.Find("GetImg").gameObject.SetActive(taskInfo.Status == TaskStatus.REWARDED);
        JUIExtensions.btn(award).onClick.RemoveAllListeners();
        JUIExtensions.btn(award).onClick.AddListener(() =>
        {
            if (CheckActivityEnd(ActivityType.SignHappy7))
                return;

            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (taskInfo.Status == TaskStatus.COMPLETED)
            {
                OnClickHappyGetButton(taskInfo, taskCfg, null);

                ClearTween();
                award.localRotation = Quaternion.Euler(Vector3.zero);
                ani.enabled = false;
                spineGo.SetActive(false);
            }
            else
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CfgId });
        });

        ClearTween();
        award.localRotation = Quaternion.Euler(Vector3.zero);
        spineGo.SetActive(taskInfo.Status == TaskStatus.COMPLETED);
        if (taskInfo.Status == TaskStatus.COMPLETED)
        {
            happyAwardCoroutine = CoroutineMgr.Instance.StartCoroutine(SetHappyAwardCoroutine(award, ani));
        }
    }

    private System.Collections.IEnumerator SetHappyAwardCoroutine(Transform trans, Animation ani)
    {
        while (true)
        {
            ani.enabled = true;
            yield return new WaitForSeconds(1);
            trans.localRotation = Quaternion.Euler(Vector3.zero);
            ani.enabled = false;
            yield return new WaitForSeconds(1);
        }
    }

    // 特惠礼包
    private void ShowHappyDailyGift(int index)
    {
        int shopId = showActivityList[index].PreLimit;
        var storeInfo = CFG_Store.GetById(shopId);

        var obj = _uiView.happy7GiftAwardTrans;
        for (int i = 0; i < obj.childCount; i++)
        {
            var obj_Item = obj.GetChild(i);
            obj_Item.gameObject.SetActive(i < storeInfo.Solditems.Length);

            if (i < storeInfo.Solditems.Length)
            {
                int cfgID = storeInfo.Solditems[i].itemID;
                GameUtil.Instance.SetItemIcon(cfgID, JUIExtensions.img(obj_Item.transform.Find("icon")), JUIExtensions.img(obj_Item.transform));
                JUIExtensions.txt(obj_Item.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(storeInfo.Solditems[i].Num);

                JUIExtensions.btn(obj_Item.transform).onClick.RemoveAllListeners();
                JUIExtensions.btn(obj_Item.transform).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgID });
                });
            }
        }

        bool soldOut = storeInfo.Limite.limittype != 999 && StoreM.Instance.GetBuyNumByID(shopId) >= storeInfo.Limite.limitnum;
        _uiView.happy7GiftAwardBtn.gameObject.SetActive(!soldOut);
        _uiView.GetGo.SetActive(soldOut);

        if (!soldOut)
        {
            JUIExtensions.txt(_uiView.happy7GiftAwardBtn.transform.Find("txt")).text = $"￥{GameUtil.Instance.FormatNumWithUnit(storeInfo.Limite.Num)}";
            JUIExtensions.btn(_uiView.happy7GiftAwardBtn).onClick.RemoveAllListeners();
            JUIExtensions.btn(_uiView.happy7GiftAwardBtn).onClick.AddListener(() =>
            {
                if (CheckActivityEnd(ActivityType.SignHappy7))
                    return;

                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnGiftItemClick(shopId);
            });
        }
    }

    private void OnGiftItemClick(int shopId, bool isFree = false)
    {
        GameUtil.Instance.GameToBuy(shopId, () =>
        {
            _uiView.happy7GiftAwardBtn.gameObject.SetActive(false);
            _uiView.GetGo.SetActive(true);
            ShowAwards(shopId);
        });
    }

    private LoopListViewItem2 OnSignHappy7Fill(LoopListView2 listView, int idx)
    {
        var item = listView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        var taskInfo = happy7TaskInfos[idx];
        CFGItem_Task taskCfg = CFG_Task.GetById(taskInfo.CfgId);

        // 标题
        JUIExtensions.txt(go.transform.Find("TitleText")).text = CFG_Task.StrTaskDes(taskInfo.CfgId);

        // 礼物
        var index = 0;
        //var cfgId = 0;
        var awardTrans = go.transform.Find("Awards");
        for (int i = 0; i < awardTrans.childCount; i++)
        {
            var award = awardTrans.GetChild(i);
            award.gameObject.SetActive(i < taskCfg.Reward.Length);
            if (i < taskCfg.Reward.Length)
            {
                index++;
                AddAward(award, int.Parse(taskCfg.Reward[i].ElementAt(0).Key), taskCfg.Reward[i].ElementAt(0).Value);
            }
        }

        var activityRewardTrans = awardTrans.GetChild(index);
        activityRewardTrans.gameObject.SetActive(true);
        AddAward(activityRewardTrans, taskCfg.ActivityReward.key, taskCfg.ActivityReward.value);
        // 进度
        (int amout, int targetnum) = CFG_Task.GetTaskPro(taskInfo.CfgId, taskInfo.Target.Amount, taskInfo.Target.TargetNum);
        if (amout > targetnum) amout = targetnum;
        JUIExtensions.txt(go.transform.Find("sliderText")).text = $"{GameUtil.Instance.FormatNumWithUnit(amout)}/{GameUtil.Instance.FormatNumWithUnit(targetnum)}";
        go.transform.Find("Slider").GetComponent<Slider>().value = amout * 1.0f / targetnum;

        var getBtnTrans = go.transform.Find("Btn_Get"); // 领取
        var goBtnTrans = go.transform.Find("Btn_Go"); // 前往
        var finishBtnTrans = go.transform.Find("Geted"); // 已完成

        if (taskInfo.Status == TaskStatus.IN_PROGRESS)
        {
            JUIExtensions.txt(go.transform.Find("TitleText")).color = new Color(169f / 255, 134f / 255, 93f / 255);
            JUIExtensions.txt(go.transform.Find("sliderText")).color = new Color(169f / 255, 134f / 255, 93f / 255);
        }
        else
        {
            JUIExtensions.txt(go.transform.Find("TitleText")).color = new Color(64f / 255, 176f / 255, 35f / 255);
            JUIExtensions.txt(go.transform.Find("sliderText")).color = new Color(64f / 255, 176f / 255, 35f / 255);
        }

        finishBtnTrans.gameObject.SetActive(taskInfo.Status == TaskStatus.REWARDED);
        goBtnTrans.gameObject.SetActive(taskInfo.Status == TaskStatus.IN_PROGRESS);
        getBtnTrans.gameObject.SetActive(taskInfo.Status == TaskStatus.COMPLETED);
        go.transform.Find("GetedImg").gameObject.SetActive(taskInfo.Status == TaskStatus.REWARDED);

        JUIExtensions.btn(goBtnTrans).onClick.RemoveAllListeners();
        JUIExtensions.btn(getBtnTrans).onClick.RemoveAllListeners();

        if (taskInfo.Status == TaskStatus.COMPLETED)
        {
            JUIExtensions.btn(getBtnTrans).onClick.AddListener(() =>
            {
                OnClickHappyGetButton(taskInfo, taskCfg, activityRewardTrans);
            });
        }

        if (taskInfo.Status == TaskStatus.IN_PROGRESS)
        {
            JUIExtensions.btn(goBtnTrans).onClick.AddListener(() =>
            {
                OnClickHappyGoToButton(taskInfo.CfgId, taskCfg.Type);
            });
        }

        return item;
    }

    private void AddAward(Transform trans, int cfgId, int count)
    {
        GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(trans.Find("icon")), JUIExtensions.img(trans));
        JUIExtensions.txt(trans.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(count);

        JUIExtensions.btn(trans).onClick.RemoveAllListeners();
        JUIExtensions.btn(trans).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
        });
    }

    private void OnClickHappyGetButton(TaskInfo taskInfo, CFGItem_Task taskCfg, Transform itemTrans)
    {
        if (CheckActivityEnd(ActivityType.SignHappy7))
            return;

        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        TcpC.C2S_TaskComplate(taskInfo.CfgId);

        if (itemTrans != null)
            Fly(itemTrans);

        var info = new ActivityAwardInfo();
        for (int i = 0; i < taskCfg.Reward.Length; i++)
        {
            info.AddAwardItem(taskCfg.Reward[i].ElementAt(0).Key, taskCfg.Reward[i].ElementAt(0).Value);
        }
        info.des = string.Empty;
        UITool.Open<UIActivityAwards>(null, info);
    }

    private void OnClickHappyGoToButton(int configID, int taskType)
    {
        if (CheckActivityEnd(ActivityType.SignHappy7))
            return;

        if (taskType == (int)pb.TargetType.SEVENDAY_TASK_GIFT_PURCHASED)
        {
            // 购买七日任务礼包，查看那一天没有购买，跳转到那一天
            var dayIndex = -1;
            for (int i = 0; i < showActivityList.Count; i++)
            {
                int shopId = showActivityList[i].PreLimit;
                var storeInfo = CFG_Store.GetById(shopId);
                bool soldOut = storeInfo.Limite.limittype != 999 && StoreM.Instance.GetBuyNumByID(shopId) >= storeInfo.Limite.limitnum;
                if (!soldOut)
                {
                    dayIndex = i;
                    break;
                }
            }

            if (dayIndex >= 0)
                OnClickHappyTabButton(dayIndex);
        }
        else
        {
            CFG_Task.JumpUI(configID, jumpHandle: () =>
            {
                UITool.Close<UIActive>();
            });
        }
    }

    private void Fly(Transform itemTrans)
    {
        Vector3 startPos = Vector3.zero;
        startPos = GetItemPostion(itemTrans, startPos);
        startPos.z = 0;
        var endPos = _uiView.starTrans.localPosition;

        var effectGo = _uiView.GetFlyEffectGo();
        effectGo.transform.localScale = Vector3.one * 10;
        effectGo.transform.localPosition = startPos;

        Vector3[] pathPoints = new Vector3[] { startPos, startPos - new Vector3(50, 5, 0), endPos };

        effectGo.SetActive(true);
        effectGo.transform.DOLocalPath(pathPoints, duration: 0.5f, pathType: PathType.CatmullRom).OnComplete(() =>
        {
            if (effectGo != null)
            {
                effectGo.SetActive(false);
                GameObject.DestroyImmediate(effectGo);
            }
        });
    }

    private Vector3 GetItemPostion(Transform trans, Vector3 res)
    {
        if (trans == _uiView.happy7Go.transform)
            return res;

        res += trans.localPosition;
        return GetItemPostion(trans.parent, res);
    }

    private void ClearTween()
    {
        if (happyAwardCoroutine != null)
        {
            CoroutineMgr.Instance.StopCoroutine(happyAwardCoroutine);
            happyAwardCoroutine = null;
        }
    }
    #endregion

    #region 星愿召唤、疯狂洗练 公共方法

    private int lastBuyAwardIndex = -1;

    /// <summary>
    /// 最后一个购买过的奖励索引
    /// </summary>
    private void GetLastBuyAwardIndex()
    {
        lastBuyAwardIndex = -1;
        for (int i = showActivityList.Count - 1; i >= 0; i--)
        {
            int storeId = showActivityList[i].PreLimit;
            var storeInfo = CFG_Store.GetById(storeId);
            // 是否售罄或已购买
            long buyCount = StoreM.Instance.GetBuyNumByID(storeId);
            bool soldOut = storeInfo.Limite.limittype != 999 && buyCount >= storeInfo.Limite.limitnum;
            if (soldOut)
            {
                lastBuyAwardIndex = i;
                return;
            }
        }
    }

    private void ShowAwards(int shopId)
    {
        var storeInfo = CFG_Store.GetById(shopId);
        var info = new ActivityAwardInfo();
        for (int i = 0; i < storeInfo.Solditems.Length; i++)
        {
            info.AddAwardItem(storeInfo.Solditems[i].itemID, storeInfo.Solditems[i].Num);
        }
        info.des = string.Empty;
        UITool.Open<UIActivityAwards>(null, info);
    }

    #endregion


    #region 星愿召唤
    private void OnInitStarWish()
    {
        if (!GetModel().ActiveTitlesDic[ActivityType.RefineFund].ActiveOpen)
            return;

        _uiView.loop_GrowthFund.InitListView(0, OnFundFill);

        _uiView.btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            int id = 50;
            if (activityType == ActivityType.CrazyRefine)
                id = 51;
            UITool.Open<UIComRule>(null, id);
        });
    }

    private LoopGridViewItem OnStarWishFill(LoopGridView view, int index, int row, int column)
    {
        var item = view.NewListViewItem("Item");
        OnStarWishItemSetting(item.transform, index);
        return item;
    }

    private void OnStarWishItemSetting(Transform item, int index)
    {
        int newindex = index;
        if (index % 4 == 2)
            newindex = index + 1;
        else if (index % 4 == 3)
            newindex = index - 1;

        Transform leftArrow = item.Find("Arrrow/l");
        Transform rightArrow = item.Find("Arrrow/r");
        Transform downArrow = item.Find("Arrrow/d");

        bool isLast = index == showActivityList.Count - 1;
        leftArrow.gameObject.SetActive(newindex % 4 == 2 && !isLast);
        rightArrow.gameObject.SetActive(newindex % 4 == 0 && !isLast);
        downArrow.gameObject.SetActive((newindex % 4 == 1 || newindex % 4 == 3) && !isLast);

        Transform awardTrans = item.Find("Awards");
        var cfgActiveItem = showActivityList[newindex];
        var cfgStoreItem = CFG_Store.GetById(cfgActiveItem.PreLimit);

        bool isCanGet = newindex == lastBuyAwardIndex + 1;
        item.Find("BgHighLight").gameObject.SetActive(isCanGet);

        bool isGeted = newindex <= lastBuyAwardIndex;
        item.Find("Btn_Geted").gameObject.SetActive(isGeted);
        Button btn_Get = JUIExtensions.btn(item.Find("Btn_Get"));
        btn_Get.gameObject.SetActive(!isGeted);

        bool isCannotBuy = newindex > lastBuyAwardIndex + 1;
        Transform suo = btn_Get.transform.Find("price/suo");
        suo.gameObject.SetActive(isCannotBuy);
        Image iconPrice = JUIExtensions.img(btn_Get.transform.Find("price/item"));

        bool isFree = cfgStoreItem.Limite.Num == 0;
        bool isShowPriceIcon = !isGeted && cfgStoreItem.Limite.needmoney != 100 && !isFree;
        iconPrice.gameObject.SetActive(isShowPriceIcon);
        if (isShowPriceIcon)
            GameUtil.Instance.SetItemIcon(cfgStoreItem.Limite.needmoney, iconPrice);

        if (!isGeted)
        {
            Transform red = btn_Get.transform.Find("dot");
            red.gameObject.SetActive(isCanGet && isFree);

            var txt = JUIExtensions.txt(btn_Get.transform.Find("price/txt"));
            btn_Get.transform.Find("free").gameObject.SetActive(isFree);
            btn_Get.transform.Find("free/suo").gameObject.SetActive(isFree && isCannotBuy);
            if (isFree)
                txt.text = "免费";
            else if (cfgStoreItem.Limite.needmoney == 100)
                txt.text = $"￥{GameUtil.Instance.FormatNumWithUnit(cfgStoreItem.Limite.Num)}";
            else
                txt.text = GameUtil.Instance.FormatNumWithUnit(cfgStoreItem.Limite.Num);
        }

        btn_Get.onClick.RemoveAllListeners();
        btn_Get.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if(activityType == ActivityType.StarWish && CheckActivityEnd(ActivityType.StarWish))
                return;
            else if(activityType == ActivityType.CrazyRefine && CheckActivityEnd(ActivityType.CrazyRefine))
                return;

            if (!isCanGet)
            {
                UITool.ShowTips("该档位无法操作，请按照顺序领取！");
                return;
            }

            OnStarWishItemClick(cfgStoreItem.ID);
        });

        for (int i = 0; i < awardTrans.childCount; i++)
        {
            bool show = i < cfgStoreItem.Solditems.Length;
            Transform awardItem = awardTrans.GetChild(i);
            awardItem.gameObject.SetActive(show);
            if (show)
            {
                var soldItem = cfgStoreItem.Solditems[i];
                GameUtil.Instance.SetItemIcon(soldItem.itemID, JUIExtensions.img(awardItem.Find("icon")), JUIExtensions.img(awardItem));
                JUIExtensions.txt(awardItem.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(soldItem.Num);
                JUIExtensions.btn(awardItem).onClick.RemoveAllListeners();
                JUIExtensions.btn(awardItem).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = soldItem.itemID });
                });
            }
        }
    }

    private void OnStarWishItemClick(int shopId)
    {
        GameUtil.Instance.GameToBuy(shopId, () =>
        {
            OnStarWishShow();

            if (activityType == ActivityType.StarWish)
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.ActiveCenter_StarWish);
            else if (activityType == ActivityType.CrazyRefine)
                RedPointManager.Instance.RPS.Dispatch(RedPointConst.ActiveCenter_CrazyRefine);

            ShowAwards(shopId);
        });
    }

    private void OnStarWishShow()
    {
        GetView().starWishSummonTrans.Find("TitleImage1").gameObject.SetActive(activityType == ActivityType.StarWish);
        GetView().starWishSummonTrans.Find("TitleImage2").gameObject.SetActive(activityType == ActivityType.CrazyRefine);

        GetLastBuyAwardIndex();
        _uiView.loopGrid_StarWish.SetListItemCount(showActivityList.Count);
        _uiView.loopGrid_StarWish.RefreshAllShownItem();
        int moveIndex = lastBuyAwardIndex == showActivityList.Count ? showActivityList.Count : lastBuyAwardIndex + 1;
        _uiView.loopGrid_StarWish.MovePanelToItemByIndex(moveIndex, offsetY: -8);
    }

    #endregion


    #region 基金

    int storeId = 38001;
    private bool isBuy = false;
    private HTTP_GrowthFundInfo[] fundInfos;
    private Dictionary<string, bool> addedAwardItem = new Dictionary<string, bool>();
    private ActivityAwardInfo activityAwardInfo;

    private void OnInitFunc()
    {
        _uiView.btnFundBuy.onClick.AddListener(() =>
        {
            GameUtil.Instance.GameToBuy(storeId, () =>
            {
                activityAwardInfo?.Clear();
                ShowFund();

                //HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                //{
                //    page = "礼包界面",
                //    component = curActiveTitle.activeName,
                //    position = "",
                //    event_click = $"购买解锁成长基金特权",
                //});
            });
        });

        if (!GetModel().ActiveTitlesDic[ActivityType.StarWish].ActiveOpen)
            return;

        _uiView.loopGrid_StarWish.InitGridView(0, OnStarWishFill);
    }

    CFG_Activity.Type fundActivityType;

    int fundNum;
    int fundPrice;

    private void ShowFund()
    {
        UITool.Open<UILoading>((ui) => {
            HttpC.OnGetGrowthFundInfo((int)fundActivityType, (_response) =>
            {
                UITool.Close<UILoading>();
                if (_response.code == 0 && !isNullGo)
                {
                    activityAwardInfo ??= new ActivityAwardInfo(5);

                    addedAwardItem.Clear();
                    bool isbuy = _response.buy == 2;
                    _uiView.btnFundBuy.gameObject.SetActive(!isbuy);
                    _uiView.txt_FundTitle2.gameObject.SetActive(isbuy);
                    //_uiView.titleLockGo.SetActive(_response.buy != 2);

                    isBuy = _response.buy == 2;
                    fundInfos = _response.data;
                    fundNum = _response.num;
                    _uiView.txt_TaskCount.text = $"当前次数：{GameUtil.Instance.FormatNumWithUnit(fundNum)}";

                    //if(activityType == ActivityType.RefineFund)
                    //{
                    //    _uiView.txt_TaskCount.text = $"洗炼次数：{fundNum}";
                    //}
                    //else if (activityType == ActivityType.CallFund)
                    //{
                    //    _uiView.txt_TaskCount.text = $"召唤次数：{fundNum}";
                    //}
                    //else if (activityType == ActivityType.CopyFund)
                    //{
                    //    _uiView.txt_TaskCount.text = $"参与次数：{fundNum}";
                    //}

                    if (fundInfos != null && fundInfos.Length > 0)
                    {
                        Array.Sort(fundInfos, (a, b) => a.day.CompareTo(b.day));
                    }
                    _uiView.loop_GrowthFund.SetListItemCount(showActivityList.Count);
                    _uiView.txt_BtnPrice.text = "￥" + fundPrice;
                }
                else
                    _uiView.loop_GrowthFund.SetListItemCount(0);

                //uiView.btnFundBuy.gameObject.SetActive(!isBuy);
                _uiView.loop_GrowthFund.RefreshAllShownItem();
                GetView().loop_GrowthFund.MovePanelToItemIndex(OnGrowthFundMovePosition(), 0);
            });
        }, new UILoading.Loading { _time = 0 });

    }

    private LoopListViewItem2 OnFundFill(LoopListView2 loopListView, int idx)
    {
        if (idx < 0 || idx >= showActivityList.Count)
            return null;

        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        var itemTrans = go.transform;

        var activityData = showActivityList[idx];
        var isAchieve = activityData.PreLimit <= fundNum;
        itemTrans.Find("Achieve").gameObject.SetActive(isAchieve);
        itemTrans.Find("BGNoGet").gameObject.SetActive(!isAchieve);
        itemTrans.Find("BGGet").gameObject.SetActive(isAchieve);

        JUIExtensions.txt(itemTrans.Find("LevelImage/Text")).text = activityData.PreLimit.ToString();

        // 上下进度条
        itemTrans.Find("SliderUpImg").gameObject.SetActive(idx != 0);
        itemTrans.Find("SliderDownImg").gameObject.SetActive(idx != showActivityList.Count - 1);
        var sliderUpGo = itemTrans.Find("SliderUpImg/Image").gameObject;
        var sliderDownGo = itemTrans.Find("SliderDownImg/Image").gameObject;

        sliderUpGo.SetActive(idx != 0 && isAchieve);
        sliderDownGo.SetActive(isAchieve);

        var isLeftGeted = false;
        var isRightGeted = false;
        if (fundInfos != null && fundInfos.Length > idx)
        {
            isLeftGeted = fundInfos[idx].status == 2;
            isRightGeted = fundInfos[idx].status == 2 && fundInfos[idx].buy == 2;
        }

        // 奖励
        var obj = itemTrans.Find("Award1s");
        var key = $"{idx}_free";
        if (!addedAwardItem.TryGetValue(key, out bool hasAdd))
            addedAwardItem.Add(key, true);

        for (int i = 0; i < obj.childCount; i++)
        {
            var obj_Item = obj.GetChild(i);
            obj_Item.gameObject.SetActive(i < activityData.CommonItem.Length);

            if (i < activityData.CommonItem.Length)
            {
                int cfgID = int.Parse(activityData.CommonItem[i].ElementAt(0).Key);

                GameUtil.Instance.SetItemIcon(cfgID, JUIExtensions.img(obj_Item.transform.Find("icon")), JUIExtensions.img(obj_Item));
                JUIExtensions.txt(obj_Item.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(activityData.CommonItem[i].ElementAt(0).Value);

                obj_Item.transform.Find("finishIconImg").gameObject.SetActive(isLeftGeted);

                var spineTrans = obj_Item.Find("spine");
                spineTrans.gameObject.SetActive(!isLeftGeted && isAchieve);
                if (!isLeftGeted && isAchieve)
                    SpineUtil.Instance.PlayItemSpineLoop(obj_Item.Find("spine"), (int)CFG_Item.GetById(cfgID).Quality);

                if (!isLeftGeted && isAchieve && !hasAdd)
                    activityAwardInfo.AddAwardItem(activityData.CommonItem[i].ElementAt(0).Key, activityData.CommonItem[i].ElementAt(0).Value);

                JUIExtensions.btn(obj_Item).onClick.RemoveAllListeners();
                JUIExtensions.btn(obj_Item).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    if (!isLeftGeted && isAchieve)
                        GetAllGrowthFundAward();
                    else
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgID });
                });
            }
        }

        obj = itemTrans.Find("Award2s");
        key = $"{idx}_buy";
        hasAdd = false;
        if (!addedAwardItem.TryGetValue(key, out hasAdd))
            addedAwardItem.Add(key, true);

        for (int i = 0; i < obj.childCount; i++)
        {
            var obj_Item = obj.GetChild(i);
            obj_Item.gameObject.SetActive(i < activityData.PayItem.Length);

            if (i < activityData.PayItem.Length)
            {
                int cfgID = int.Parse(activityData.PayItem[i].ElementAt(0).Key);

                GameUtil.Instance.SetItemIcon(cfgID, JUIExtensions.img(obj_Item.transform.Find("icon")), JUIExtensions.img(obj_Item));
                JUIExtensions.txt(obj_Item.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(activityData.PayItem[i].ElementAt(0).Value);

                obj_Item.transform.Find("LockImg").gameObject.SetActive(!isBuy);
                obj_Item.transform.Find("finishIconImg").gameObject.SetActive(isRightGeted);

                var spineTrans = obj_Item.Find("spine");
                spineTrans.gameObject.SetActive(!isRightGeted && isAchieve && isBuy);
                if (!isRightGeted && isAchieve && isBuy)
                    SpineUtil.Instance.PlayItemSpineLoop(spineTrans, (int)CFG_Item.GetById(cfgID).Quality);

                if (!isRightGeted && isAchieve && isBuy && !hasAdd)
                    activityAwardInfo.AddAwardItem(activityData.PayItem[i].ElementAt(0).Key, activityData.PayItem[i].ElementAt(0).Value);

                JUIExtensions.btn(obj_Item).onClick.RemoveAllListeners();
                JUIExtensions.btn(obj_Item).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    if (!isRightGeted && isAchieve && isBuy)
                        GetAllGrowthFundAward();
                    else
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgID });
                });
            }
        }

        JUIExtensions.btn(itemTrans.Find("Achieve/AchieveButton")).onClick.RemoveAllListeners();
        JUIExtensions.btn(itemTrans.Find("Achieve/AchieveButton1")).onClick.RemoveAllListeners();
        JUIExtensions.btn(itemTrans.Find("BtnBuy")).onClick.RemoveAllListeners();
        // 按钮
        if (!isLeftGeted)
        {
            JUIExtensions.btn(itemTrans.Find("Achieve/AchieveButton")).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                GetAllGrowthFundAward();
            });
        }

        if (!isRightGeted)
        {
            if (!isBuy)
            {
                JUIExtensions.btn(itemTrans.Find("BtnBuy")).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    PayAndUnlock();
                });
            }

            JUIExtensions.btn(itemTrans.Find("Achieve/AchieveButton1")).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (isBuy)
                    GetAllGrowthFundAward();
                else
                    PayAndUnlock();
            });
        }

        return item;
    }

    // 一键领取奖励
    private void GetAllGrowthFundAward()
    {
        //HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        //{
        //    page = "礼包界面",
        //    component = curActiveTitle.activeName,
        //    position = "",
        //    event_click = $"免费领取礼包",
        //});

        // 展示奖励
        HttpC.OnGetGrowthFundAward((int)fundActivityType, (_r) =>
        {
            if (isNullGo)
                return;

            // 展示奖励
            activityAwardInfo.des = string.Empty;
            UITool.Open<UIActivityAwards>(null, activityAwardInfo);

            ShowFund();
        });
    }

    // 购买
    private void PayAndUnlock()
    {
        var storeInfo = CFG_Store.GetById(storeId);
        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = $"即将花费<color=#E03042>￥{GameUtil.Instance.FormatNumWithUnit(storeInfo.Limite.Num)}</color>解锁成长基金特权！",
            YesFun = () =>
            {
                GameUtil.Instance.GameToBuy(storeId, () =>
                {
                    activityAwardInfo?.Clear();

                    ShowFund();

                    //HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                    //{
                    //    page = "礼包界面",
                    //    component = curActiveTitle.activeName,
                    //    position = "",
                    //    event_click = $"购买解锁成长基金特权",
                    //});
                });
            },
        });
    }

    private int OnGrowthFundMovePosition()
    {
        int position = 0;

        var roleLevel = RoleM.Instance.Level;
        for (int i = 0; i < showActivityList.Count; i++)
        {
            var isAchieve = showActivityList[i].PreLimit <= roleLevel;
            var isLeftGeted = false;
            var isRightGeted = false;
            if (fundInfos != null && fundInfos.Length > i)
            {
                isLeftGeted = fundInfos[i].status == 2;
                isRightGeted = fundInfos[i].status == 2 && fundInfos[i].buy == 2;
            }

            if ((!isLeftGeted && isAchieve) || (!isRightGeted && isAchieve && isBuy))
            {
                position = i;
                break;
            }
            else if (isLeftGeted)
                position = i;
        }

        return position;
    }
    #endregion

    #region active time
    private void AddTimer()
    {
        _uiView.activeTimeText.text = "";
        _uiView.txt_FundTime.text = "";
        _uiView.txt_StarWishTime.text = "";

        var endTime = ActivityM.Instance.GetAcitvityEndTime(activityType);
        SetActiveTimeText(endTime);
        timeItem = GameHelper.Instance.AddInvoke(new TimeItem(endTime, 1, () =>
        {
            endTime--;
            SetActiveTimeText(endTime);
            if (endTime < 0)
                Clear();
        }));
    }

    private void SetActiveTimeText(long endTime)
    {
        if (endTime <= 0)
        {
            if(activityType == ActivityType.RefineFund || activityType == ActivityType.CallFund || activityType == ActivityType.CopyFund)
            {
                _uiView.txt_FundTime.text = "已结束";
            }
            else if(activityType == ActivityType.StarWish || activityType == ActivityType.CrazyRefine)
            {
                _uiView.txt_StarWishTime.text = "已结束";
            }
            else
            {
                _uiView.activeTimeText.text = "已结束";
            }
        }
        else
        {
            var timeStr = GameUtil.Instance.GetTimeSToDHOrHMS(endTime);
            if (activityType == ActivityType.RefineFund || activityType == ActivityType.CallFund || activityType == ActivityType.CopyFund)
            {
                _uiView.txt_FundTime.text = timeStr + "后结束";
            }
            else if (activityType == ActivityType.StarWish || activityType == ActivityType.CrazyRefine)
            {
                _uiView.txt_StarWishTime.text = timeStr + "后结束";
            }
            else
            {
                _uiView.activeTimeText.text = timeStr + "后结束";
            }
        }
    }

    private void Clear()
    {
        if (timeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(timeItem);
            timeItem = null;
        }
    }

    #endregion

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        BindRed();
        curDay = 1;
        curHappyTab = 0;
        isHappyTab = false;
        Clear();
        ClearTween();

        curDateTime = DateTime.Now.ToString("yyyyMMdd");
        if (ActivityM.Instance.GetAcitvityState(ActivityType.Sign7))
            OnClickTabButton(GetModel().ActiveTitlesDic[ActivityType.Sign7]);
        else if (ActivityM.Instance.GetAcitvityState(ActivityType.SignHappy7))
            OnClickTabButton(GetModel().ActiveTitlesDic[ActivityType.SignHappy7]);
        else if (ActivityM.Instance.GetAcitvityState(ActivityType.StarWish))
            OnClickTabButton(GetModel().ActiveTitlesDic[ActivityType.StarWish]);
        else if (ActivityM.Instance.GetAcitvityState(ActivityType.CrazyRefine))
            OnClickTabButton(GetModel().ActiveTitlesDic[ActivityType.CrazyRefine]);
        else
            OnClickTabButton(GetModel().ActiveTitlesDic[ActivityType.Sign7]);

    }

    [Subscriber(ThreadMode.Main)]
    private void Task_Refresh(GameEvent.Refresh_Task _e)
    {
        if (isNullGo) return;

        if (ActivityM.Instance.GetAcitvityEndTime(ActivityType.SignHappy7) <= 0)
            return;

        var str = DateTime.Now.ToString("yyyyMMdd");
        if (curDateTime != str)
        {
            curDay++;
            curDateTime = str;
        }

        if (uiGo.activeSelf && isHappyTab && (_e.TaskClass == TaskClass.SEVENDAY || _e.TaskClass == TaskClass.SEVENDAY_ACTIVITY))
            OnSignHappyShow();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        Clear();
        ClearTween();
        ClearBindRed();
    }

    public override void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(this);
        Clear();
        ClearTween();
        ClearBindRed();
    }
}