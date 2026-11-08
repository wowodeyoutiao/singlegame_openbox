using JEngine.Singleton;
using System;
using System.Collections.Generic;
using UnityEngine;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 商店信息
    /// </summary>
    public class StoreM : SingletonTemplate<StoreM>
    {
        public class StoreInfo
        {
            public long buyNum { get; set; }

            public int inventoryNum { get; set; }

            // 购买总数量
            public int buyAllNum { get; set; }
        }

        private Dictionary<long, StoreInfo> StoreBuyNumDic = new Dictionary<long, StoreInfo>();
        public void SetSoreInfo(HTTP_StoreList[] data)
        {
            for (int i = 0; i < data.Length; i++)
            {
                if (StoreBuyNumDic.ContainsKey(data[i].goods_uuid))
                {
                    StoreBuyNumDic[data[i].goods_uuid] = new StoreInfo { buyNum = data[i].buy, inventoryNum = data[i].date, buyAllNum = data[i].buy_all_num };
                }
                else
                {
                    StoreBuyNumDic.Add(data[i].goods_uuid, new StoreInfo { buyNum = data[i].buy, inventoryNum = data[i].date, buyAllNum = data[i].buy_all_num });
                }
            }
        }

        /// <summary>
        /// 检测库存是否充足
        /// </summary>
        public bool CheckInventory(int mailId)
        {
            var cfg = CFG_Store.GetById(mailId);
            if (cfg.Serverdetect == 1)
            {
                if (cfg.Limite.limittype == 999)
                {
                    return GetInventoryNumAddByID(mailId) == -1 || GetInventoryNumAddByID(mailId) > 0;
                }
                else
                {
                    return GetInventoryNumAddByID(mailId) != 0 && GetBuyNumByID(mailId) < cfg.Limite.limitnum;
                }
            }
            else
            {
                if (cfg.Limite.limittype == 999)
                {
                    return true;
                }
                else
                {
                    return GetBuyNumByID(mailId) < cfg.Limite.limitnum;
                }
            }
        }

        //-1 为无限
        public int GetCanBuyNum(int mailId)
        {
            var cfg = CFG_Store.GetById(mailId);
            if (cfg.Serverdetect == 1)
            {
                if (cfg.Limite.limittype == 999)
                {
                    return -1;
                }
                else
                {
                    return (int)(GetBuyNumByID(mailId) - cfg.Limite.limitnum);
                }
            }
            else
            {
                if (cfg.Limite.limittype == 999)
                {
                    return -1;
                }
                else
                {
                    return (int)(cfg.Limite.limitnum - GetBuyNumByID(mailId));
                }
            }
        }

        public int GetBuyAllNum(int id)
        {
            if (StoreBuyNumDic.ContainsKey(id))
                return StoreBuyNumDic[id].buyAllNum;

            return 0;
        }

        public long GetBuyNumByID(long id)
        {
            if (StoreBuyNumDic.ContainsKey(id))
            {
                return StoreBuyNumDic[id].buyNum;
            }
            else
            {
                return 0;
            }
        }

        public void SetBuyNumAddByID(long id, int num)
        {
            if (StoreBuyNumDic.ContainsKey(id))
            {
                StoreBuyNumDic[id].buyNum += num;
            }
            else
            {
                StoreBuyNumDic.Add(id, new StoreInfo { buyNum = 1 });
            }
        }

        //库存数量为服务器控制
        public int GetInventoryNumAddByID(long id)
        {
            if (StoreBuyNumDic.ContainsKey(id))
            {
                return StoreBuyNumDic[id].inventoryNum;
            }
            else
            {
                return 0;
            }
        }

        #region 提取账号
        //1-微信，2-支付宝
        private int AccountType { get; set; }
        private string AccountStr { get; set; }
        private string AccountNameStr { get; set; }

        public void SetExtractAccount(int _AccountType, string _Account, string _Name)
        {
            AccountType = _AccountType;
            AccountStr = _Account;
            AccountNameStr = _Name;
        }

        public string GetExtractAccountName()
        {
            return AccountNameStr;
        }

        public string GetExtractAccount()
        {
            return AccountStr;
        }

        public int GetExtractAccountType()
        {
            return AccountType;
        }
        #endregion

        private static long RMBBuyEventId = 0;
        private static Dictionary<long, Action> BuyEventDic = new Dictionary<long, Action>();

        public long SetBuyEvent(Action action = null)
        {
            RMBBuyEventId++;
            BuyEventDic.Add(RMBBuyEventId, action);
            return RMBBuyEventId;
        }

        public void OnDoneBuyEvent(long id, bool isRemove = false)
        {
            if (BuyEventDic.ContainsKey(id))
            {
                if (!isRemove)
                {
                    BuyEventDic[id]?.Invoke();
                }
                BuyEventDic.Remove(id);
            }
            else
            {
                Debug.LogError($"OnDoneBuyEvent id {id}not find");
            }
        }
    }

    /// <summary>
    /// 请求商城购买列表
    /// </summary>
    public struct HTTP_StoreListRequset
    {
        //对应商城表 StoreType 0代表全部
        public int store_type { get; set; }
    }

    /// <summary>
    /// 请求商城购买返回
    /// </summary>
    public struct HTTP_StoreListResponse
    {
        public int code { get; set; }
        public HTTP_StoreList[] data { get; set; }

    }

    public struct HTTP_StoreList
    {
        /// <summary>
        /// 已购买数量
        /// </summary>
        public long buy { get; set; }
        /// <summary>
        /// 商品id
        /// </summary>
        public long goods_uuid { get; set; }

        public int date { get; set; }

        public int other_num { get; set; }

        // 购买总数量
        public int buy_all_num { get; set; }


    }

    /// <summary>
    /// 购买物品
    /// </summary>
    public struct HTTP_StoreOperateRequset
    {
        public long goods_uuid { get; set; }
        public int num { get; set; }

        //额外信息
        public string extra { get; set; }
        public string platform { get; set; }
        public long request_id { get; set; }

        public pb.RefineSlot refine_slot { get; set; }
    }

    /// <summary>
    /// 购买物品请求微信小程序(安卓)
    /// </summary>
    public struct HTTP_StoreOperateWXItem
    {
        public string mode;
        public int buyQuantity;
        public int env;
        public string currencyType;
        public string platform;
        public string productId;
        public int goodsPrice;
    }


    /// <summary>
    /// 购买物品请求微信小程序(安卓)
    /// </summary>
    public struct HTTP_WXCustomerItem
    {
        public string source_type;
        public string out_trade_no;
        public long user_id;
        public int server_id;
    }

    /// <summary>
    /// 购买物品返回
    /// </summary>
    public struct HTTP_StoreOperateResponse
    {
        public int code { get; set; }
        public bool data { get; set; }
        public string error { get; set; }

        /// <summary>
        /// 订单号
        /// </summary>
        public string ordersn { get; set; }

        //微信支付:签名数据JSON字符串
        public string signData { get; set; }
        //微信支付:支付签名
        public string paySig { get; set; }

        public string signature { get; set; }

        public string platform { get; set; }

        //微信 IOS 支付
        public string prepay_internal_id { get; set; }
    }

    /// <summary>
    /// 提取账号
    /// </summary>
    public struct HTTP_StoreExtractAccountResponse
    {
        public int code { get; set; }
        public int type { get; set; }
        public string detail { get; set; }
        public string data { get; set; }
        public string name { get; set; }
    }

    /// <summary>
    /// 提取账号 更新
    /// </summary>
    public struct HTTP_StoreExtractAccountUpdateResponse
    {
        public int type { get; set; }
        public string detail { get; set; }
        public string name { get; set; }


    }

}

