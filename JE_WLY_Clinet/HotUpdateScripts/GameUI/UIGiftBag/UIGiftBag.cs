using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using SuperScrollView;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

public class UIGiftBag : UIBase<UIModelGiftBag, UIViewGiftBag>
{
    private UIViewGiftBag uIViewGiftBag;

    private UIModelGiftBag.ActiveTitle curActiveTitle;

    private List<CFGItem_Activity> showActivityList;
    private Dictionary<int, string> _GiftLimetDesDic = new Dictionary<int, string>
    {
        [1] = "每日限购",
        [2] = "每周限购",
        [3] = "每月限购",
    };

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        uIViewGiftBag = GetView();

        uIViewGiftBag.btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGiftBag>();
        });

        // 按钮
        for (int i = 0; i < GetModel().ActiveTitlesDic.Count; i++)
        {
            var titleData = GetModel().ActiveTitlesDic.ElementAt(i).Value;

            uIViewGiftBag.tabBtns[i].transform.Find("Normal/Label").GetComponent<Text>().text = titleData.activeName;
            uIViewGiftBag.tabBtns[i].transform.Find("Select/Label").GetComponent<Text>().text = titleData.activeName;

            uIViewGiftBag.tabBtns[i].gameObject.name = titleData.activeType.ToString();

            uIViewGiftBag.tabBtns[i].gameObject.SetActive(titleData.ActiveOpen);
            uIViewGiftBag.tabBtns[i].onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnClickTabButton(titleData);
            });
        }

        uIViewGiftBag.infiniteScroll_Gifts.InitListView(0, OnActivityFill);
    }

    private void BindRed()
    {
        for (int i = 0; i < GetModel().ActiveTitlesDic.Count; i++)
        {
            var titleData = GetModel().ActiveTitlesDic.ElementAt(i).Value;
            if (titleData.ActiveOpen)
            {
                var index = i;
                RedPointManager.Instance.RPS.Register(titleData.RedKey, (node) =>
                {
                    uIViewGiftBag.tabBtns[index].transform.Find("Dot").gameObject.SetActive(node.pointCount > 0);
                });
            }
        }
    }

    private void ClearBindRed()
    {
        foreach (var item in GetModel().ActiveTitlesDic.Values)
        {
            if (!item.ActiveOpen) continue;
            RedPointManager.Instance.RPS.Register(item.RedKey, null);
        }
    }

    private void OnClickTabButton(UIModelGiftBag.ActiveTitle idex)
    {
        for (int i = 0; i < uIViewGiftBag.tabBtns.Length; i++)
        {
            var sel = uIViewGiftBag.tabBtns[i].gameObject.name == idex.activeType.ToString();
            uIViewGiftBag.tabBtns[i].transform.Find("Normal").gameObject.SetActive(!sel);
            uIViewGiftBag.tabBtns[i].transform.Find("Select").gameObject.SetActive(sel);
        }

        if(!string.IsNullOrEmpty(idex.iconName))
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.GiftBagBanner, idex.iconName, uIViewGiftBag.logoImg);

       
        uIViewGiftBag.normalNodeGo.SetActive(true);

        curActiveTitle = idex;

        if (idex.cfgId > 0)
            showActivityList = CFG_Activity.GetActiveListByType((CFG_Activity.Type)idex.cfgId);

        switch (idex.activeType)
        {
            case ActivityType.PartnerGift: // 伙伴礼包
            case ActivityType.PetGift: // 宠物礼包
            case ActivityType.Gemstone: // 魂石礼包
            case ActivityType.Element: // 元素礼包
            case ActivityType.BlessingGift: // 福袋礼包
            case ActivityType.RefineGift: // 洗练礼包
                OnGetShopInfo(idex.shopType);
                break;

            default:
                break;
        }
    }

    //[Subscriber(ThreadMode.Main)]
    // 整点刷新礼包
    private void OnSync_ZeroRefresh(GameEvent.Refresh_ZeroTime _e)
    {
        OnClickTabButton(curActiveTitle);
    }

    #region 每日礼包、每日特价、伙伴礼包、宠物礼包、元素礼包
    private void OnGetShopInfo(int shopType)
    {
        UITool.Open<UILoading>((ui) => {
            HttpC.OnGetShopInfo(shopType, (_response) =>
            {
                UITool.Close<UILoading>();
                if (_response.code == 0 && !isNullGo)
                {
                    uIViewGiftBag.infiniteScroll_Gifts.SetListItemCount(showActivityList.Count);
                    uIViewGiftBag.infiniteScroll_Gifts.RefreshAllShownItem();
                }
            });
        }, new UILoading.Loading { _time = 0 });
       
    }

    private LoopListViewItem2 OnActivityFill(LoopListView2 loopListView, int idx)
    {
        if (idx < 0 || idx >= showActivityList.Count)
            return null;

        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        int shopId = showActivityList[idx].PreLimit;
        var storeInfo = CFG_Store.GetById(shopId);

        JUIExtensions.txt(go.transform.Find("txt_Name")).text = showActivityList[idx].Name;

        string limitStr = string.Empty;
        if (storeInfo.Limite.limittype != 999)
        {
            if (storeInfo.Limite.times > 0)
            {
                var curTimes = StoreM.Instance.GetBuyAllNum(shopId);
                if (curTimes < storeInfo.Limite.times)
                    limitStr = $"<color=#a9865d>（再购买<color=#FF0000>{storeInfo.Limite.times - curTimes}</color>次后免费）</color>";
            }
            limitStr = $"{limitStr}{_GiftLimetDesDic[storeInfo.Limite.limittype]}:{StoreM.Instance.GetBuyNumByID(shopId)}/{storeInfo.Limite.limitnum}";
        }
        JUIExtensions.txt(go.transform.Find("txt_Limet")).text = limitStr;

        var obj = go.transform.Find("Awards");
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

        var buyBtnTrans = go.transform.Find("Btn_Buy");
        var freeBtnTrans = go.transform.Find("Btn_Free");
        var videoBtnTrans = go.transform.Find("Btn_Video");

        bool soldOut = storeInfo.Limite.limittype != 999 && StoreM.Instance.GetBuyNumByID(shopId) >= storeInfo.Limite.limitnum;
        go.transform.Find("txt_tip").gameObject.SetActive(soldOut);
        buyBtnTrans.gameObject.SetActive(false);
        freeBtnTrans.gameObject.SetActive(false);
        videoBtnTrans.gameObject.SetActive(false);

        if(!soldOut)
        {
            // limittype 为 999
            if(storeInfo.Limite.limittype == 999)
            {
                ShowBuyButton(buyBtnTrans, storeInfo.Limite.needmoney, storeInfo.Limite.Num);
                AddButtonListeners(buyBtnTrans, storeInfo);
            }
            else
            {
                // 使用广告
                var useAd = storeInfo.Limite.needmoney == 101;
                // 判断总次数，购买总次数为X次后免费
                var isFree = (storeInfo.Limite.Num == 0) || (storeInfo.Limite.times != 0 && storeInfo.Limite.times <= StoreM.Instance.GetBuyAllNum(shopId));
                if(isFree)
                {
                    freeBtnTrans.gameObject.SetActive(true);
                    AddButtonListeners(freeBtnTrans, storeInfo, isFree:true);
                }
                else if(useAd)
                {
                    videoBtnTrans.gameObject.SetActive(true);
                    AddButtonListeners(videoBtnTrans, storeInfo, isVideo:true);
                }
                else
                {
                    ShowBuyButton(buyBtnTrans, storeInfo.Limite.needmoney, storeInfo.Limite.Num);
                    AddButtonListeners(buyBtnTrans, storeInfo);
                }
            }
        }

        return item;
    }

    private void ShowBuyButton(Transform trans, int needmoney, int num)
    {
        trans.gameObject.SetActive(true);
        JUIExtensions.txt(trans.Find("txt")).text = GameUtil.Instance.FormatNumWithUnit(num);
        GameUtil.Instance.SetShopItemIcon(needmoney, JUIExtensions.img(trans.Find("icon")));
    }

    private void AddButtonListeners(Transform trans, CFGItem_Store storeInfo, bool isVideo = false, bool isFree = false)
    {
        JUIExtensions.btn(trans).onClick.RemoveAllListeners();
        JUIExtensions.btn(trans).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (!isVideo)
                OnGiftItemClick(storeInfo, isFree);
            else
                OnVideoButtonClick(storeInfo.ID);
        });
    }

    private void OnVideoButtonClick(int shopId)
    {
        GameUtil.Instance.GameToBuy(shopId, () =>
        {
            ShowAwards(shopId);
            OnGetShopInfo(curActiveTitle.shopType);
            RedPointManager.Instance.RPS.Dispatch(curActiveTitle.RedKey, GiftBagM.Instance.GetGiftRedNum(curActiveTitle.cfgId));
            SendShuShuClick("广告领取");
        });
    }

    private void OnGiftItemClick(CFGItem_Store storeInfo, bool isFree = false)
    {
        bool _free = isFree;
        if(isFree || storeInfo.Limite.needmoney == 100)
        {
            GameUtil.Instance.GameToBuy(storeInfo.ID, () =>
            {
                ShowAwards(storeInfo.ID);
                OnGetShopInfo(curActiveTitle.shopType);
                RedPointManager.Instance.RPS.Dispatch(curActiveTitle.RedKey, GiftBagM.Instance.GetGiftRedNum(curActiveTitle.cfgId));
                SendShuShuClick(_free ? "免费领取" : "购买");
            }, isFree);
            return;
        }

        if (GameUtil.Instance.ComCoinToCheck(storeInfo.Limite.needmoney, storeInfo.Limite.Num, null))
        {
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = $"是否消耗<color=#2C7E40>{GameUtil.Instance.FormatNumWithUnit(storeInfo.Limite.Num)}{CFG_Item.GetById(storeInfo.Limite.needmoney).Name}</color>购买该礼包？",
                YesFun = () =>
                {
                    GameUtil.Instance.GameToBuy(storeInfo.ID, () =>
                    {
                        ShowAwards(storeInfo.ID);
                        OnGetShopInfo(curActiveTitle.shopType);
                        RedPointManager.Instance.RPS.Dispatch(curActiveTitle.RedKey, GiftBagM.Instance.GetGiftRedNum(curActiveTitle.cfgId));
                        SendShuShuClick("元宝道具购买");
                    });
                },
            });
        }
    }

    private void ShowAwards(int shopId)
    {
        var storeInfo = CFG_Store.GetById(shopId);
        var info = new ActivityAwardInfo();
        for(int i = 0; i <storeInfo.Solditems.Length; i++)
        {
            info.AddAwardItem(storeInfo.Solditems[i].itemID, storeInfo.Solditems[i].Num);
        }
        info.des = string.Empty;
        UITool.Open<UIActivityAwards>(null, info);
    }

    /// <summary>
    /// 领取类型
    /// </summary>
    /// <param name="getType"></param>
    private void SendShuShuClick(string getType)
    {
        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "礼包界面",
            component = curActiveTitle.activeName,
            position = "",
            event_click = getType,
        });
    }
    #endregion

    //private void ScrollToButton(ActivityType activityType)
    //{
    //    //int index = 0;
    //    //foreach (var item in GetModel().ActiveTitlesDic.Values)
    //    //{
    //    //    if (item.ActiveOpen)
    //    //        index++;

    //    //    if (item.activeType == activityType)
    //    //        break;
    //    //}

    //    //if (index == 0) return;

    //    //float width = GetView().scrollTabBtns.GetComponent<RectTransform>().rect.width;
    //    //var hlg = GetView().tabBtnsTrans.GetComponent<HorizontalLayoutGroup>();
    //    //int openedActiveNum = GetActiveOpenedNums();
    //    //float totalWidth = hlg.padding.left + hlg.padding.right + (openedActiveNum * 160) + ((openedActiveNum - 1) * hlg.spacing);
    //    //float moveSpace = totalWidth - width;
    //    //float needMoveSpace = (index - 1) * (160 + hlg.spacing);
    //    //GetView().scrollTabBtns.horizontalNormalizedPosition = Mathf.Clamp01(needMoveSpace / moveSpace);
    //}

    private int GetActiveOpenedNums()
    {
        int count = 0;
        foreach (var item in GetModel().ActiveTitlesDic.Values)
        {
            if(item.ActiveOpen)
                count++;
        }
        return count;
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_ZeroRefresh);

        BindRed();

        if (args != null && args.Length > 0)
        {
            var ativityType = (ActivityType)args[0];
            //ScrollToButton(ativityType);
            OnClickTabButton(GetModel().ActiveTitlesDic[ativityType]);
        }
        else
        {
            OnClickTabButton(GetModel().ActiveTitlesDic[ActivityType.PartnerGift]);
            //GetView().scrollTabBtns.horizontalNormalizedPosition = 0;
        }
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_ZeroRefresh);

        ClearBindRed();
    }

    public override void OnDestroy()
    {
        ClearBindRed();
    }
}
