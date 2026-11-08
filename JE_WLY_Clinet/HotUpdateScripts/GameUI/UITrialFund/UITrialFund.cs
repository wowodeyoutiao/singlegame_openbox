using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;
using JEngine.Helper;
using JEngine.UI;
using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Ocsp;
using SuperScrollView;
using pb;
using System.Linq;
using HotUpdateScripts.ProtoM;
using System;

public class UITrialFund : UIBase<UIModelTrialFund, UIViewTrialFund>
{
    // 未解锁至尊奖励
    const string noPurchased = "解锁基金且领取了全部奖励后将出现下一基金";
    // 已解锁至尊奖励
    const string purchased = "领取了全部奖励后将出现下一基金";
    // 每层最大奖励数量
    const int childItemCount = 4;
    // 已初始化列表？
    bool isInitListView;
    int storeId = 15001;

    List<CFGItem_Activity> cfg_Items;
    HTTP_ActivityTrialFundAwardResponse awardInfo;

    /// <summary>
    /// 奖励道具类型
    /// </summary>
    public enum AwardType
    {
        //基础奖励
        Base,
        //至尊奖励
        Extreme,
    }

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_OneClickGet.onClick.AddListener(OnClickAward);
        GetView().btn_PayUnlock.onClick.AddListener(PayAndUnlock);
        GetView().btn_Close.onClick.AddListener(CloseUI);
    }

    public override void OnEnter(params object[] args)
    {
        cfg_Items = CFG_Activity.GetActiveListByType(CFG_Activity.Type.TrialFund);
        cfg_Items.Sort((a, b) => a.Content.CompareTo(b.Content));
        GetAwardInfo();
    }

    private void GetAwardInfo()
    {
        HttpC.GetTrialFundAwardInfo((_r) =>
        {
            if (isNullGo)
                return;

            if (_r.code == 0)
            {
                awardInfo = _r;
                ShowListView();
                ScrollTo(_r.day);
                ShowPayButton(_r.buy == 1);
                SetTipText(_r.buy == 1 ? noPurchased : purchased);
            }
        });
    }

    private void ShowListView()
    {
        if (isInitListView)
        {
            GetView().loopListView.SetListItemCount(cfg_Items.Count);
            GetView().loopListView.RefreshAllShownItem();
        }
        else
        {
            isInitListView = true;
            GetView().loopListView.InitListView(cfg_Items.Count, OnRefreshItem);
        }
    }

    void ScrollTo(int currentDay)
    {
        // awardInfo.data.Length 以经领取过的奖励数量作为索引，表示下一个奖励的位置，这个位置可能是是可领取或者不可领取的状态。
        GetView().loopListView.MovePanelToItemIndex(awardInfo.data.Length, 0);
    }

    private LoopListViewItem2 OnRefreshItem(LoopListView2 listView, int index)
    {
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        CFGItem_Activity cfg_item = cfg_Items[index];
        JUIExtensions.txt(item.transform, "floor/txt").text = $"{cfg_item.PreLimit}\n层";
        // 当前层是否可领取奖励？
        bool isCanReceivedTheDay = cfg_item.PreLimit <= awardInfo.day;
        // 是否已领取过当前层奖励？
        bool isReceived = index < awardInfo.data.Length;
        // 当前层领取状态
        int status = 3;
        HTTP_ActivityTrialFundDays[] awardsSortData = awardInfo.data;
        Array.Sort(awardsSortData, (a, b) => a.day.CompareTo(b.day));
        int currentFloorIsBuy = isReceived ? awardsSortData[index].buy : awardInfo.buy;

        if (isReceived)
        {
            status = 2;
        }
        else if(isCanReceivedTheDay)
        {
            status = 1;
        }

        for (int i = 0; i < childItemCount; i++)
        {
            Transform child = item.transform.Find($"Item{i + 1}");
            if (i == 0)
            {
                ShowAwardItems(child, cfg_item.CommonItem[0].ElementAt(0), AwardType.Base, status);
                child.gameObject.SetActive(true);
            }
            else
            {
                bool isShow = i <= cfg_item.PayItem.Length;//bug
                if (isShow)
                {
                    ShowAwardItems(child, cfg_item.PayItem[i - 1].ElementAt(0), AwardType.Extreme, status, currentFloorIsBuy);
                }
                child.gameObject.SetActive(isShow);
            }
        }

        // 显示item的领取状态
        // status 领取状态：1-待领取，2-已领取，3-不可领取，4-需要补签
        // buy 开通状态：1-没开通，2-已开通
        void ShowAwardItems(Transform item, KeyValuePair<string, int> itemData, AwardType awardType, int status, int buy = 1)
        {
            int CfgId = int.Parse(itemData.Key);
            GameUtil.Instance.SetItemIcon(CfgId, JUIExtensions.img(item.Find("icon")), JUIExtensions.img(item));
            JUIExtensions.txt(item, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(itemData.Value);

            JUIExtensions.btn(item).onClick.RemoveAllListeners();
            JUIExtensions.btn(item).onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CfgId });
            });

            item.Find("icon/lock").gameObject.SetActive(awardType == AwardType.Extreme && awardInfo.buy == 1);
            bool recerned = awardType == AwardType.Base ? status == 2 : status == 2 && buy == 2;
            item.Find("icon/received").gameObject.SetActive(recerned);
            bool canGet = awardType == AwardType.Base ? status == 1 : (status == 1 && buy == 2) || (status == 2 && buy == 1 && awardInfo.buy == 2);
            item.Find("icon/canget").gameObject.SetActive(canGet);
        }

        return item;
    }

    private void OnClickAward()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        HttpC.GetTrialFundAward((_r) =>
        {
            if (isNullGo)
                return;

            GetAwardInfo();
        });
    }

    private void PayAndUnlock()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        GameUtil.Instance.GameToBuy(storeId, () =>
        {
            GetAwardInfo();
        });
    }

    private void ShowPayButton(bool isShow)
    {
        GetView().btn_PayUnlock.gameObject.SetActive(isShow);
    }

    private void SetTipText(string value)
    {
        GetView().txt_AwardTip.text = value;
    }

    private void CloseUI()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Close<UITrialFund>();
    }

    public override void OnExit()
    {
        cfg_Items = null;
        awardInfo = default;
    }

    public override void OnDestroy()
    {
        
    }
}