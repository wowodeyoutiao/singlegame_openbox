using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.MyNet;

using MetaFramework.UI;

using pb;

using System;
using System.ComponentModel;

using UnityEngine;

public class UICommonItemBuy : UIBase<UIModelCommonItemBuy, UIViewCommonItemBuy>
{
    public enum Option
    {
        Buy,
        Use,
    }

    public struct BuyInfo
    {
        public Option option;
        public CFGItem_Store store;
        public int itemId;
        public Action buyCallback;
        public string title;
        public string subTitle;
        public string description1;
        public string description2;
        public string bottomTip;
        /// <summary>
        /// 挑战次数兑换卷不足时，是否跳转到商店门票？
        /// </summary>
        public bool jumpShopTickets;
    }

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UICommonItemBuy>();
        });

        GetView().btn_buy.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (buyInfo.option == Option.Buy)
                OnBuyRequest();
            else if (buyInfo.option == Option.Use)
                OnUseRequest();
        });

        GetView().btn_jian.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyOrUseCount--;
            if (_BuyOrUseCount < _MinCount)
            {
                _BuyOrUseCount = _MinCount;
            }
            else if (_BuyOrUseCount > _MaxCount)
            {
                _BuyOrUseCount = _MaxCount;
            }
            GetView().input_Buy.text = _BuyOrUseCount.ToString();

            if (buyInfo.option == Option.Buy)
                ShowBuyConsume();
        });

        GetView().btn_add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyOrUseCount++;
            if (_BuyOrUseCount < _MinCount)
            {
                _BuyOrUseCount = _MinCount;
            }
            else if (_BuyOrUseCount > _MaxCount)
            {
                _BuyOrUseCount = _MaxCount;
            }
            GetView().input_Buy.text = _BuyOrUseCount.ToString();

            if (buyInfo.option == Option.Buy)
                ShowBuyConsume();
        });

        GetView().btn_Sun10.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyOrUseCount -= 10;
            if (_BuyOrUseCount < _MinCount)
            {
                _BuyOrUseCount = _MinCount;
            }
            else if (_BuyOrUseCount > _MaxCount)
            {
                _BuyOrUseCount = _MaxCount;
            }
            GetView().input_Buy.text = _BuyOrUseCount.ToString();

            if (buyInfo.option == Option.Buy)
                ShowBuyConsume();
        });

        GetView().btn_Add10.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyOrUseCount += 10;

            if (_BuyOrUseCount < _MinCount)
            {
                _BuyOrUseCount = _MinCount;
            }
            else if (_BuyOrUseCount > _MaxCount)
            {
                _BuyOrUseCount = _MaxCount;
            }

            GetView().input_Buy.text = _BuyOrUseCount.ToString();

            if (buyInfo.option == Option.Buy)
                ShowBuyConsume();
        });

        GetView().input_Buy.onValueChanged.AddListener((v) =>
        {
            long buyCount = 1;
            if (int.TryParse(v, out int num))
                buyCount = num;

            if (buyCount < _MinCount)
            {
                buyCount = _MinCount;
            }
            else if (buyCount > _MaxCount)
            {
                buyCount = _MaxCount;
            }

            _BuyOrUseCount = buyCount;
            GetView().input_Buy.text = _BuyOrUseCount.ToString();

            if (buyInfo.option == Option.Buy)
                ShowBuyConsume();
        });
    }

    private long _BuyOrUseCount = 1;
    private long _MinCount = 1;
    private long _MaxCount = 1;

    private long consume = 0;

    // 购买信息
    BuyInfo buyInfo;

    /// <summary>
    /// 进入界面
    /// </summary>
    /// <param name="args">[0] storeId, [1] title text, [2] callback</param>
    public override void OnEnter(params object[] args)
    {
        buyInfo = (BuyInfo)args[0];
        _BuyOrUseCount = 1;
        int id = buyInfo.option switch
        {
            Option.Buy => buyInfo.store.Solditems[0].itemID,
            Option.Use => buyInfo.itemId,
            _ => 0,
        };

        GetView().txt_Title.text = string.IsNullOrEmpty(buyInfo.title) ? $"{CFG_Item.GetById(id).Name}道具购买" : buyInfo.title;
        string defaultDesc = buyInfo.option == Option.Buy ? "请选择所要购买的数量：" : "请选择所要增加的数量";
        GetView().txt_subTitle.text = string.IsNullOrEmpty(buyInfo.subTitle) ? defaultDesc : buyInfo.subTitle;
        GetView().txt_tip.text = buyInfo.bottomTip;
        GetView().txt_ConsumeDesc.text = buyInfo.description1;
        GetView().txt_ConsumeDesc2.text = buyInfo.description2;
        //GetView().txt_ConsumeDesc.gameObject.SetActive(buyInfo.option == Option.Use);

        if (buyInfo.option == Option.Buy)
        {
            GetView().txt_BtnName.text = "购买";
            GetView().btn_add.gameObject.SetActive(buyInfo.store.StoreType != (int)CFG_Store.StoreType.Ticket);
            GetView().btn_Add10.gameObject.SetActive(buyInfo.store.StoreType != (int)CFG_Store.StoreType.Ticket);
            GetView().btn_jian.gameObject.SetActive(buyInfo.store.StoreType != (int)CFG_Store.StoreType.Ticket);
            GetView().btn_Sun10.gameObject.SetActive(buyInfo.store.StoreType != (int)CFG_Store.StoreType.Ticket);
            long itemNum = BagM.Instance.GetNumByCfgID(buyInfo.store.Limite.needmoney);

            long maxBuyCount = itemNum / buyInfo.store.Limite.Num;
            _MaxCount = buyInfo.store.Limite.limittype == 999 ? maxBuyCount : Math.Min(maxBuyCount, buyInfo.store.Limite.limitnum - StoreM.Instance.GetBuyNumByID(buyInfo.store.ID));

            GameUtil.Instance.SetItemIcon(buyInfo.store.Limite.needmoney, GetView().img_Consume);
            ShowBuyConsume();
        }
        else if (buyInfo.option == Option.Use)
        {
            GetView().txt_BtnName.text = "使用";
            GetView().btn_add.gameObject.SetActive(true);
            GetView().btn_Add10.gameObject.SetActive(true);
            GetView().btn_jian.gameObject.SetActive(true);
            GetView().btn_Sun10.gameObject.SetActive(true);

            long itemCount = BagM.Instance.GetNumByCfgID(buyInfo.itemId);
            long maxCount;
            if (buyInfo.itemId == 56)
            {
                long canusenum = Math.Max(0, CFG_ItemUse.GetDailyUseLimit(buyInfo.itemId) - BagM.Instance.GetDailyUsedItem(buyInfo.itemId));
                maxCount = Math.Min(itemCount, canusenum);
            }
            else
            {
                maxCount = itemCount;
            }
            _MaxCount = Math.Max(1, maxCount);
            GameUtil.Instance.SetItemIcon(buyInfo.itemId, GetView().img_Consume);
            GetView().txt_Consume.text = itemCount == 0 ? $"<color=#DA4E38>{itemCount}</color>" : GameUtil.Instance.FormatNumWithUnit(itemCount);
            GetView().btn_buy.interactable = true;
        }

        GetView().input_Buy.text = _BuyOrUseCount.ToString();
    }

    public void UpdateDescription2(string content)
    {
        GetView().txt_ConsumeDesc2.text = content;
    }

    public int GetItemID()
    {
        return buyInfo.itemId;
    }

    /// <summary>
    /// 道具只用、兑换
    /// </summary>
    private void OnUseRequest()
    {
        long itemCount = BagM.Instance.GetNumByCfgID(buyInfo.itemId);
        if (itemCount == 0)
        {
            if (buyInfo.jumpShopTickets)
            {
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = "当前包裹中没有可以使用的门票，快去礼包界面获取吧！",
                    BtnYes = "前往礼包",
                    YesFun = () =>
                    {
                        if (buyInfo.itemId == 56)
                            UITool.Open<UISpecialOffer>(null, ActivityType.TimeGift);
                        else if (buyInfo.itemId >= 25073 && buyInfo.itemId <= 25083)
                            UITool.Open<UIShop>(null, CFG_Store.StoreType.Ticket);

                        UITool.Close<UICommonItemBuy>();
                    }
                });
            }
            return;
        }

        long uid = BagM.Instance.GetItemByCFGID(buyInfo.itemId).Uid;
        TcpC.C2S_ItemUse(uid, (int)_BuyOrUseCount);
        UITool.Close<UICommonItemBuy>();
        ShuShuSend(false);
    }

    /// <summary>
    /// 商品购买
    /// </summary>
    private void OnBuyRequest()
    {
        if (!GameUtil.Instance.ComCoinToCheck(buyInfo.store.ID, null))
        {
            return;
        }

        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = $"是否确定消耗{GameUtil.Instance.FormatNumWithUnit(consume)}{CFG_Item.GetById(buyInfo.store.Limite.needmoney).Name}道具兑换{CFG_Item.GetById(buyInfo.store.Solditems[0].itemID).Name}？",
            YesFun = () =>
            {
                GameUtil.Instance.GameToBuy(buyInfo.store.ID, (int)_BuyOrUseCount, buyInfo.buyCallback);
                UITool.Close<UICommonItemBuy>();

                ShuShuSend(true);
            },
        });
    }

    private void ShowBuyConsume()
    {
        consume = 0;
        if (buyInfo.store.Consumechange != null && buyInfo.store.Consumechange.Length > 0)
        {
            var buyNum = StoreM.Instance.GetBuyNumByID(buyInfo.store.ID);
            var num = 0; // 次数
            var maxPrice = 0;
            if (buyNum == 0)
            {
                consume = buyInfo.store.Limite.Num; // 未购买过时，第一次购买价格
                num = 1;
            }

            if (_BuyOrUseCount > num)
            {
                for (int i = 0; i < buyInfo.store.Consumechange.Length; i++)
                {
                    if (buyInfo.store.Consumechange[i].key > num + buyNum + 1)
                        break;

                    maxPrice = buyInfo.store.Consumechange[i].value;
                    if (buyInfo.store.Consumechange[i].key == num + buyNum + 1)
                    {
                        consume += buyInfo.store.Consumechange[i].value;
                        num++;
                    }

                    if (num >= _BuyOrUseCount)
                        break;
                }

                if (num < _BuyOrUseCount)
                    consume += maxPrice * (_BuyOrUseCount - num);
            }
        }
        else
            consume = buyInfo.store.Limite.Num * _BuyOrUseCount;

        string price = GameUtil.Instance.ComCoinToCheck(buyInfo.store.ID, null) ? $"{GameUtil.Instance.FormatNumWithUnit(consume)}" : $"<color=#CE4031>{GameUtil.Instance.FormatNumWithUnit(consume)}</color>";
        GetView().txt_Consume.text = price;
        GetView().btn_buy.interactable = GameUtil.Instance.ComCoinToCheck(buyInfo.store.ID, null);
    }

    // 数数上报
    private void ShuShuSend(bool isBuy)
    {
        if (isBuy)
        {
            string shopTag = "";
            if (UITool.IsOpen<UIShop>())
            {
                var shopSelect = UITool.Get<UIShop>().selectShop;
                shopTag = shopSelect switch
                {
                    CFG_Store.StoreType.Coupon => "元宝",
                    CFG_Store.StoreType.Coin => "金币",
                    CFG_Store.StoreType.Diamanta => "战魂",
                    CFG_Store.StoreType.SiLian => "试炼",
                    CFG_Store.StoreType.Ticket => "门票",
                    _ => ""
                };
            }

            string solditems = string.Join(",", Array.ConvertAll(buyInfo.store.Solditems, item => CFG_Item.GetById(item.itemID).Name));
            string consume = buyInfo.store.Limite.needmoney == 100 ? "人民币" : CFG_Item.GetById(buyInfo.store.Limite.needmoney).Name;
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "商店",
                component = shopTag,
                position = $"商品内容：{solditems}",
                event_click = $"{consume}购买",
            });
        }
        else
        {
            if (buyInfo.itemId == 56)
            {
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = $"挂机升级",
                    event_click = $"挑战剩余时间兑换（15分钟）",
                });
            }
            else if (buyInfo.itemId == 25073)
            {
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = $"首领",
                    event_click = $"挑战次数兑换",
                });
            }
            else if (buyInfo.itemId == 25074)
            {
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = $"遗迹",
                    event_click = $"挑战次数兑换",
                });
            }
            else if (buyInfo.itemId >= 25075 && buyInfo.itemId <= 25080)
            {
                string type = buyInfo.itemId switch
                {
                    25075 => "强化石副本挑战",
                    25076 => "升星石副本挑战",
                    25077 => "魂石副本挑战",
                    25078 => "书页副本挑战",
                    25079 => "金币石副本挑战",
                    25080 => "经验副本挑战",
                    _ => "",
                };

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "副本",
                    position = type,
                    event_click = $"挑战次数兑换",
                });
            }
            else if (buyInfo.itemId == 25081)
            {
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "五行试炼",
                    position = "",
                    event_click = $"挑战次数兑换",
                });
            }
            else if (buyInfo.itemId == 25082)
            {
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "竞技",
                    position = "即时赛挑模式",
                    event_click = $"挑战次数兑换",
                });
            }
            else if (buyInfo.itemId == 25083)
            {
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "竞技",
                    position = "回合赛挑模式",
                    event_click = $"挑战次数兑换",
                });
            }
        }
    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {

    }
}