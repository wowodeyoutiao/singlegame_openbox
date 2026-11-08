using HotUpdateScripts.ProtoM;

using JEngine.Helper;

using MetaFramework.UI;

using UnityEngine;


public class UIPVPCountBuy : UIBase<UIModelPVPCountBuy, UIViewPVPCountBuy>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPVPCountBuy>();
        });
        GetView().btn_buy.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            System.Action closeCall = () => {
                UITool.Close<UIPVPCountBuy>();
            };

            int storeId = roomType switch
            {
                CFG_Room.RoomType.PVP_One => _XmlInfo.storeId,
                CFG_Room.RoomType.PVP_More => _XmlInfo.storeId,
                CFG_Room.RoomType.WuXing => _Element.needitem,
                _ => 0
            };

            if (!GameUtil.Instance.ComCoinToCheck(storeId, closeCall, true, _BuyCount))
            {
                return;
            }

            System.Action backCall = () => 
            {
                UITool.Close<UIPVPCountBuy>();
                GameHelper.Instance.SetDeayTimeDo(() => 
                {
                    if (roomType == CFG_Room.RoomType.PVP_One || roomType == CFG_Room.RoomType.PVP_More)
                        TcpC.OnPVPInfo(_MatchType);
                    else if (roomType == CFG_Room.RoomType.WuXing)
                        TcpC.C2S_ElementInfoRequest();
                }, 
                1f);
            };
            GameUtil.Instance.GameToBuy(storeId, _BuyCount, backCall);
        });

        GetView().btn_jian.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyCount--;
            OnSetBuyCount(_BuyCount);
        });

        GetView().btn_add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _BuyCount++;
            OnSetBuyCount(_BuyCount);
        });

        GetView().input_Buy.onValueChanged.AddListener((num) =>
        {
            OnSetBuyCount(int.Parse(num));
        });
    }

    private void OnSetBuyCount(int num)
    {
        if (roomType == CFG_Room.RoomType.PVP_One || roomType == CFG_Room.RoomType.PVP_More)
        {
            int minNum = Mathf.Min(RoleM.Instance.GetCanNumBuyPVP(_MatchType), 1);
            _BuyCount = Mathf.Clamp(num, minNum, RoleM.Instance.GetCanNumBuyPVP(_MatchType));

            var cfg = CFG_Store.GetById(_XmlInfo.storeId);
            GameUtil.Instance.SetItemIcon(cfg.Limite.needmoney, GetView().img_Icon);
            string count = GameUtil.Instance.FormatNumWithUnit(cfg.Limite.Num * _BuyCount);
            GetView().txt_Icon.text = GameUtil.Instance.SetAndCheckStr(count, cfg.Limite.needmoney, cfg.Limite.Num * _BuyCount);
        }
        else if(roomType == CFG_Room.RoomType.WuXing)
        {
            int minNum = Mathf.Min(TurnBasedM.Instance.GetCanBuyChallengeCount_WuXing(), 1);
            _BuyCount = Mathf.Clamp(num, minNum, TurnBasedM.Instance.GetCanBuyChallengeCount_WuXing());
            var cfg = CFG_Store.GetById(CFG_XMLConfig.xmlConfig.Element.needitem);
            GameUtil.Instance.SetItemIcon(cfg.Limite.needmoney, GetView().img_Icon);
            string count = GameUtil.Instance.FormatNumWithUnit(cfg.Limite.Num * _BuyCount);
            GetView().txt_Icon.text = GameUtil.Instance.SetAndCheckStr(count, cfg.Limite.needmoney, cfg.Limite.Num * _BuyCount);
        }

        GameHelper.Instance.SetDeayFrame(() =>
        {
            GetView().input_Buy.text = "" + _BuyCount;
        });
    }

    private int _BuyCount = 1;
    private int _MatchType;

    CFG_Room.RoomType roomType;

    private CFG_XMLConfig.XML_Config.XML_Item_PVP _XmlInfo;
    private CFG_XMLConfig.XML_Config.XML_Item_Element _Element;

    public override void OnEnter(params object[] args)
    {
        _MatchType = (int)args[0];
        roomType = (CFG_Room.RoomType)_MatchType;

        switch (roomType)
        {
            case CFG_Room.RoomType.PVP_One:
                _XmlInfo = CFG_XMLConfig.xmlConfig.PVPInfo;
                break;
            case CFG_Room.RoomType.PVP_More:
                _XmlInfo = CFG_XMLConfig.xmlConfig.GVGInfo;
                break;
            case CFG_Room.RoomType.WuXing:
                _Element = CFG_XMLConfig.xmlConfig.Element;
                break;
            default:
                break;
        }

        if (roomType == CFG_Room.RoomType.PVP_One || roomType == CFG_Room.RoomType.PVP_More)
        {
            var cfg = CFG_Store.GetById(_XmlInfo.storeId);
            int enoughBuy = (int)(BagM.Instance.GetNumByCfgID(cfg.Limite.needmoney) / cfg.Limite.Num);
            OnSetBuyCount(Mathf.Min(enoughBuy, RoleM.Instance.GetCanNumBuyPVP(_MatchType)));
        }
        else if (roomType == CFG_Room.RoomType.WuXing)
        {
            var cfg = CFG_Store.GetById(CFG_XMLConfig.xmlConfig.Element.needitem);
            int enoughBuy = (int)(BagM.Instance.GetNumByCfgID(cfg.Limite.needmoney) / cfg.Limite.Num);
            OnSetBuyCount(Mathf.Min(enoughBuy, TurnBasedM.Instance.GetCanBuyChallengeCount_WuXing()));
        }

        string tip = "";
        if (roomType == CFG_Room.RoomType.PVP_One || roomType == CFG_Room.RoomType.PVP_More)
            tip = $"提示:今日可购买次数<color=#CE4031>{RoleM.Instance.GetCanNumBuyPVP(_MatchType)}</color>";
        else if(roomType == CFG_Room.RoomType.WuXing)
            tip = $"提示:今日可购买次数<color=#CE4031>{TurnBasedM.Instance.GetCanBuyChallengeCount_WuXing()}</color>";

        GetView().txt_tip.text = tip;
    }

    public override void OnExit()
    {
        _XmlInfo = null;
        _Element = null;
    }

    public override void OnDestroy()
    {

    }
}