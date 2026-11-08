using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using MetaFramework.UI;
using UnityEngine;

public class UIGuildShopItemBuy : UIBase<UIModelGuildShopItemBuy, UIViewGuildShopItemBuy>
{
    public struct GuildShopItemInfo
    {
        public GuildShopItemInfo(int itemId, int merchandiseId, int price, int limit)
        {
            this.itemId = itemId;
            this.merchandiseId = merchandiseId;
            this.price = price;
            this.limit = limit;
        }

        public int itemId;
        public int merchandiseId;
        public int price;
        public int limit;
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
            UITool.Close<UIGuildShopItemBuy>();
        });

        GetView().btn_buy.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnGuildShopBuyRequest();
        });

        GetView().btn_jian.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyCount--;
            _BuyCount = Mathf.Clamp(_BuyCount, _MinCount, _MaxCount);
            GetView().input_Buy.text = _BuyCount.ToString();
            ShowConsume();
        });

        GetView().btn_add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyCount++;
            _BuyCount = Mathf.Clamp(_BuyCount, _MinCount, _MaxCount);
            GetView().input_Buy.text = _BuyCount.ToString();
            ShowConsume();
        });

        GetView().input_Buy.onValueChanged.AddListener((v) =>
        {
            if(int.TryParse(v, out int num))
                _BuyCount = Mathf.Clamp(int.Parse(v), _MinCount, _MaxCount);
            else
                _BuyCount = _MinCount;

            GetView().input_Buy.text = _BuyCount.ToString();
            ShowConsume();
        });
    }

    private int _BuyCount = 1;
    private int _MinCount = 1;
    private int _MaxCount = 1;

    GuildShopItemInfo guildShopItemInfo;

    /// <summary>
    /// 进入界面
    /// </summary>
    /// <param name="args">[0] itemId, [1] price</param>
    public override void OnEnter(params object[] args)
    {
        guildShopItemInfo = (GuildShopItemInfo)args[0];

        GetView().input_Buy.text = string.Empty;

        GameHelper.Instance.SetDeayFrame(() =>
        {
            _BuyCount = 1;
            _MaxCount = Mathf.Min(Mathf.FloorToInt(RoleM.Instance.GuildContribution / (float)guildShopItemInfo.price), guildShopItemInfo.limit);
            _MaxCount = Mathf.Max(_MaxCount, _MinCount);
            GetView().input_Buy.text = System.Text.Encoding.UTF8.GetString(System.Text.Encoding.UTF8.GetBytes(_BuyCount.ToString()));
        });
     
    }

    private void OnGuildShopBuyRequest()
    {
        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = $"是否确定消耗{guildShopItemInfo.price * _BuyCount}行会贡献兑换{CFG_Item.GetById(guildShopItemInfo.itemId).Name}？",
            YesFun = () => {
                TcpC.C2S_GuildShopBuyRequest(guildShopItemInfo.merchandiseId, _BuyCount);
                UITool.Close<UIGuildShopItemBuy>();

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "行会",
                    position = "行会商店",
                    event_click = $"商品购买（{CFG_Item.GetById(guildShopItemInfo.itemId).Name}）",
                });
            },
        });
    }

    private void ShowConsume()
    {
        GameUtil.Instance.SetItemIcon((int)CFG_Item.ItemID.GuildContribution, GetView().img_Consume);
        GetView().txt_Consume.text = $"{guildShopItemInfo.price * _BuyCount}行会贡献";
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}