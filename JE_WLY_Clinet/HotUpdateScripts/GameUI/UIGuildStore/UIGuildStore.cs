using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Event;
using HotUpdateScripts.ProtoM;
using pb;
using JEngine.UI;
using SuperScrollView;
using JEngine.Helper;

public class UIGuildStore : UIBase<UIModelGuildStore, UIViewGuildStore>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public override void OnInit()
    {
        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 35);
        });

        GetView().btn_AddItem.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGuildDonation>();
        });

        GetView().loopGrid_Store.InitGridView(0, OnItemRefresh);
    }

    bool isShow;

    public void OnShow(bool isShow)
    {
        if (isShow == this.isShow)
            return;

        this.isShow = isShow;

        if (isShow)
        {
            OnEnter();
        }
        else
        {
            OnExit();
        }
        uiGo.SetActive(isShow);
    }

    List<CFGItem_GuildStore> storeItems = new List<CFGItem_GuildStore>();
    GuildInfo guildInfo;

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnSync_SkillUpgrede);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_RefreshPrayItem);

        storeItems = CFG_GuildStore.GetList();
        ShowPrayItemCount(RoleM.Instance.GuildContribution);
        TcpC.C2S_GuildInfoRequest(RoleM.Instance.GuildId);
    }

    private LoopGridViewItem OnItemRefresh(LoopGridView gridView, int index, int row, int column)
    {
        LoopGridViewItem item = gridView.NewListViewItem("Item");
        Transform t = item.transform;
        int idx = index;

        JUIExtensions.txt(t.Find("txt_Name")).text = GameUtil.Instance.GetItemName(storeItems[index].Goods[0].id);
        GameUtil.Instance.SetItemIcon(storeItems[index].Goods[0].id, JUIExtensions.img(t.Find("Item/img_icon")), JUIExtensions.img(t.Find("Item")));
        JUIExtensions.txt(t.Find("Item/txt_num")).text = storeItems[index].Goods[0].num.ToString();

        int personalRemaining = storeItems[index].PersonalLimit;
        int guildRemaining = storeItems[index].GuildLimit;

        if (TryGetPurchasedItem(storeItems[index].UUID, out GuildShopBuyInfo shopInfo))
        {
            personalRemaining = storeItems[index].PersonalLimit - shopInfo.BuyNum;
            guildRemaining = storeItems[index].GuildLimit - shopInfo.TotalNum;
        }

        string personalCountString = personalRemaining == 0 ? $"<color=#D2241E>{personalRemaining}</color>" : personalRemaining.ToString();
        JUIExtensions.txt(t.Find("txt_DailyRedemption")).text = $"每日兑换 {personalCountString}/{storeItems[index].PersonalLimit}";

        string guildCountString = guildRemaining == 0 ? $"<color=#D2241E>{guildRemaining}</color>" : guildRemaining.ToString();
        JUIExtensions.txt(t.Find("txt_TodayRemaining")).text = $"今日行会剩余 {guildRemaining}/{storeItems[index].GuildLimit}";

        bool islocked = guildInfo.Level < storeItems[index].Unlocklevel;
        t.Find("icon").gameObject.SetActive(!islocked);
        t.Find("txt_price").gameObject.SetActive(!islocked);
        t.Find("img_mask").gameObject.SetActive(islocked);
        t.Find("txt_Lime").gameObject.SetActive(islocked);

        JUIExtensions.btn(t.transform.Find("Item")).onClick.RemoveAllListeners();
        JUIExtensions.btn(t.transform.Find("Item")).onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (islocked)
                return;

            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = storeItems[idx].Goods[0].id });
        });

        if (islocked)
        {
            JUIExtensions.txt(t.Find("txt_Lime")).text = $"{storeItems[index].Unlocklevel}级行会解锁";
            JUIExtensions.txt(t.Find("img_mask/txt_Lime")).text = $"{storeItems[index].Unlocklevel}级行会解锁";
        }
        else
        {
            JUIExtensions.txt(t.Find("txt_price")).text = $"{GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(storeItems[index].Price[0].num), storeItems[index].Price[0].id, storeItems[index].Price[0].num)}兑换";
        }

        JUIExtensions.btn(t).onClick.RemoveAllListeners();
        JUIExtensions.btn(t).onClick.AddListener(() =>
        {
            if (islocked)
                return;

            bool isEnough = GameUtil.Instance.ComCoinToCheck(storeItems[idx].Price[0].id, storeItems[idx].Price[0].num, null, true);
            if (isEnough)
            {
                int limitNum = Mathf.Min(personalRemaining, guildRemaining);
                UITool.Open<UIGuildShopItemBuy>(null, new UIGuildShopItemBuy.GuildShopItemInfo
                {
                    itemId = storeItems[idx].Goods[0].id,
                    merchandiseId = storeItems[idx].UUID,
                    price = storeItems[idx].Price[0].num,
                    limit = limitNum,
                });
            }
        });

        return item;
    }

    private void ShowPrayItemCount(long count)
    {
        GetView().txt_ItemCount.text = count.ToString();
    }

    private void OnGuildShopBuyInfoRequest()
    {
        TcpC.C2S_GuildShopBuyInfoRequest();
    }

    // 购买后刷新
    //[Subscriber(ThreadMode.Main)]
    private void OnSync_SkillUpgrede(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.code == 0 && _e.cid == CmdId.C2S_GuildShopBuy)
        {
            // 刷新商店
            OnGuildShopBuyInfoRequest();
        }
    }

    // 刷新行会经验点
    //[Subscriber(ThreadMode.Main)]
    private void OnSync_RefreshPrayItem(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        ShowPrayItemCount(RoleM.Instance.GuildContribution);
    }

    List<GuildShopBuyInfo> buyInfos = new List<GuildShopBuyInfo>();

    // 商店信息
    [Subscriber(ThreadMode.Main)]
    private void OnSync_GuildShopBuyInfo(GameEvent.Sync_GuildShopBuyInfo _e)
    {
        if (isNullGo)
            return;

        buyInfos = _e.buyInfo;
        GetView().loopGrid_Store.SetListItemCount(storeItems.Count);
        GetView().loopGrid_Store.RefreshAllShownItem();
    }

    [Subscriber(ThreadMode.Main)]
    // 工会信息
    private void Event_GuildInfo(GameEvent.Sync_GuildInfo _e)
    {
        if (isNullGo)
            return;

        guildInfo = _e.guildInfo;
        OnGuildShopBuyInfoRequest();
    }

    private bool TryGetPurchasedItem(int storeItemId, out GuildShopBuyInfo item)
    {
        item = buyInfos.Find(x => x.ItemId == storeItemId);
        return item != null;
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnSync_SkillUpgrede);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_RefreshPrayItem);
        isShow = false;
        guildInfo = null;
    }

    public override void OnDestroy()
    {

    }
}