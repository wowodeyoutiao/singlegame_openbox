using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using System.Linq;

using UnityEngine;

using static CFG_Item;

public class UIBattleSupportGift : UIBase<UIModelBattleSupportGift, UIViewBattleSupportGift>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_close.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIBattleSupportGift>();
        });

        GetView().btn_Pay.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBtnBuyClick();
        });

        GetView().txt_OriginalPrice.text = "原价：1288元";
    }

    public override void OnEnter(params object[] args)
    {
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(RefreshTopActiveButton);
        OnUIRefresh();
    }

    private void OnBtnBuyClick()
    {
        if (!GameUtil.Instance.ComCoinToCheck(cfg_ActivityInfo.PreLimit, null))
        {
            return;
        }

        GameUtil.Instance.GameToBuy(cfg_ActivityInfo.PreLimit, () =>
        {
            // 购买成功后关闭礼包界面，隐藏大厅按钮
            JEvent.defaultEvent.Post(new GameEvent.Refresh_HallActivityButton());
            UITool.Close<UIBattleSupportGift>();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "助战礼包界面",
                component = "购买按钮",
                position = "",
                event_click = $"购买礼包",
            });
        });
    }


    //private TimeItem activeTimeRemaining;

    //private void RemoveTime()
    //{
    //    if (activeTimeRemaining != null && GameHelper.Instance != null)
    //    {
    //        GameHelper.Instance.RemoveInvoke(activeTimeRemaining);
    //        activeTimeRemaining = null;
    //    }
    //}

    private CFGItem_Activity cfg_ActivityInfo;

    private void OnUIRefresh()
    {
        cfg_ActivityInfo = CFG_Activity.GetActiveListByType(CFG_Activity.Type.BattleSupportGift)[0];

        for (int i = 0; i < 3; i++)
        {
            Transform obj = GetView().root_Awards.GetChild(i);
            obj.gameObject.SetActive(i < cfg_ActivityInfo.CommonItem.Length);
            if (i < cfg_ActivityInfo.CommonItem.Length)
            {
                int itemId = int.Parse(cfg_ActivityInfo.CommonItem[i].ElementAt(0).Key);
                GameUtil.Instance.SetItemIcon(itemId, JUIExtensions.img(obj.Find("icon")), JUIExtensions.img(obj));

                JUIExtensions.txt(obj.transform.Find("txt_num")).text = cfg_ActivityInfo.CommonItem[i].ElementAt(0).Value.ToString();

                JUIExtensions.txt(obj.transform.Find("txt_quality")).text = GameUtil.Instance.GetItemName(itemId);
                JUIExtensions.outline(obj.transform.Find("txt_quality")).effectColor = GameUtil.Instance.GetQualityOutlineColorByItemID(itemId);

                JUIExtensions.btn(obj).onClick.RemoveAllListeners();
                JUIExtensions.btn(obj).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = itemId });
                });
            }
        }

        //RemoveTime();
        //long seconds = GameUtil.Instance.GetRegistRoleTimeSecoud((int)cfg_ActivityInfo.ResetTime);
        GetView().txt_Des.text = cfg_ActivityInfo.RuleDes;

        //activeTimeRemaining = GameHelper.Instance.AddInvoke(new TimeItem(seconds + 1, 1, () =>
        //{
        //    if (isNullGo) return;
        //    GetView().txt_Des.text = GetTimeDesc(seconds);
        //    if (seconds <= 0)
        //    {
        //        GetView().txt_Des.text = "";
        //        JEvent.defaultEvent.Post(new GameEvent.Refresh_HallActivityButton() { index = 6, forceHide = true });
        //        UITool.Close<UIBattleSupportGift>();
        //        return;
        //    }
        //    seconds--;
        //}));
    }

    //private string GetTimeDesc(long seconds)
    //{
    //    string str = $"<color=#ff5855>{GameUtil.Instance.GetTimeSToDHMS(seconds)}</color>{cfg_ActivityInfo.RuleDes}";
    //    return str.Insert(str.IndexOf("终生"), "\n");
    //}

    private void RefreshTopActiveButton(GameEvent.Refresh_ZeroTime _e)
    {
        OnUIRefresh();
    }

    public override void OnExit()
    {
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(RefreshTopActiveButton);
        //RemoveTime();
    }

    public override void OnDestroy()
    {

    }
}