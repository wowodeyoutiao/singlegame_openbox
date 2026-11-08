using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using System.Linq;
using UnityEngine;
using static CFG_Item;

public class UIFirstPay : UIBase<UIModelFirstPay, UIViewFirstPay>
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
            UITool.Close<UIFirstPay>();
        });

        GetView().btn_get.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBtnGetClick();
        });

        GetView().btn_double_buy.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBtnBuyClick();
        });

        GetView().btn_Pay.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBtnBuyClick();
        });
    }

    private void OnBtnGetClick()
    {
        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = "是否确定直接领取?",
            YesFun = () => {
                HttpC.GetFirstPayGet(1, (req) => {
                    if (req.code == 0 && !isNullGo)
                        OnHttpRefresh();
                    else
                        UITool.ShowTips(req.data);
                });

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首充界面",
                    component = "领取按钮",
                    position = "",
                    event_click = $"免费领取礼包",
                });
            },
            BtnYes = "单倍领",
            BtnNo = "再想想",
        });
    }

    private void OnHttpRefresh()
    {
        HttpC.GetFirstPayInfo((rep) =>
        {
            if (rep.code == 0 && !isNullGo)
            {
                _FirstPayInfo = rep;
                if(_FirstPayInfo.day == 3 && _FirstPayInfo.status == (int)ActivityM.SignState.Geted)
                {
                    JEvent.defaultEvent.Post(new GameEvent.Refresh_HallActivityButton() { index = 2, forceHide = true });
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Close<UIFirstPay>();
                    return;
                }

                OnUIRefresh();
            }
        });

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首充界面",
            component = "领取按钮",
            position = "",
            event_click = $"购买礼包",
        });
    }

    private void OnBtnBuyClick()
    {
        System.Action JumpFun = delegate { UITool.Close<UIFirstPay>(); };
        if (!GameUtil.Instance.ComCoinToCheck(ShowActivityInfo.PreLimit, JumpFun))
        {
            return;
        }
        GameUtil.Instance.GameToBuy(ShowActivityInfo.PreLimit, OnHttpRefresh);
    }

    private CFGItem_Activity ShowActivityInfo;
    private TimeItem TimeJaction;

    private void RemoveTime()
    {
        if (TimeJaction != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeJaction);
            TimeJaction = null;
        }
    }

    private void OnUIRefresh()
    {
        if(_FirstPayInfo.status == (int)ActivityM.SignState.Geted && _FirstPayInfo.day != 3)
        {
            ShowActivityInfo = CFG_Activity.GetActiveListByType(CFG_Activity.Type.FirstPay)[_FirstPayInfo.day];
            // 首日礼包图片资源
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FirstPayAtlas, "word_First_recharge_" + (_FirstPayInfo.day + 1), GetView().day_image);
        }
        else
        {
            ShowActivityInfo = CFG_Activity.GetActiveListByType(CFG_Activity.Type.FirstPay)[_FirstPayInfo.day - 1];
            // 首日礼包图片资源
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FirstPayAtlas, "word_First_recharge_" + _FirstPayInfo.day, GetView().day_image);
        }

        var shopInfo = CFG_Store.GetById(ShowActivityInfo.PreLimit);
        for (int i = 0; i < 3; i++)
        {
            Transform obj = GetView().root_Awards.GetChild(i);
            obj.gameObject.SetActive(i < ShowActivityInfo.CommonItem.Length);
            if (i < ShowActivityInfo.CommonItem.Length)
            {
                int itemId = int.Parse(ShowActivityInfo.CommonItem[i].ElementAt(0).Key);
                GameUtil.Instance.SetItemIcon(itemId, JUIExtensions.img(obj.Find("icon")), JUIExtensions.img(obj));

                JUIExtensions.txt(obj.transform.Find("txt_num")).text = ShowActivityInfo.CommonItem[i].ElementAt(0).Value.ToString();

                JUIExtensions.txt(obj.transform.Find("txt_quality")).text = GameUtil.Instance.GetItemName(itemId);
                JUIExtensions.outline(obj.transform.Find("txt_quality")).effectColor = GameUtil.Instance.GetQualityOutlineColorByItemID(itemId);

                obj.Find("geted").gameObject.SetActive(_FirstPayInfo.status == (int)ActivityM.SignState.Geted && _FirstPayInfo.day == 3);
                JUIExtensions.btn(obj).onClick.RemoveAllListeners();
                JUIExtensions.btn(obj).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = itemId });
                });
            }
        }

        //未购买
        if(_FirstPayInfo.status == (int)ActivityM.SignState.NoGet || _FirstPayInfo.status == (int)ActivityM.SignState.CanGet)
        {
            GetView().txt_Time.text = "";
        }
        else
        {
            RemoveTime();
            long seconds = GameUtil.Instance.GetNowToTimeSecoud(1, 0);
            GetView().txt_Time.text = $"<color=#74f374>{GameUtil.Instance.GetTimeSToHMS(seconds)}</color>后领取礼包";
            
            TimeJaction = GameHelper.Instance.AddInvoke(new TimeItem(seconds, 1, () =>
            {
                if (isNullGo) return;
                seconds--;
                GetView().txt_Time.text = $"<color=#74f374>{GameUtil.Instance.GetTimeSToHMS(seconds)}</color>后领取礼包";
                if (seconds <= 0)
                {
                    GetView().txt_Time.text = "";
                    RemoveTime();
                    OnHttpRefresh();
                    JEvent.defaultEvent.Post(new GameEvent.Refresh_HallActivityButton());
                    return;
                }
            }, true, true));
        }

        bool isFirstDay = _FirstPayInfo.day == 1;
        GetView().btn_get.gameObject.SetActive(_FirstPayInfo.status == (int)ActivityM.SignState.CanGet);
        GetView().btn_double_buy.gameObject.SetActive(_FirstPayInfo.status != (int)ActivityM.SignState.Geted && !isFirstDay);
        GetView().btn_Pay.gameObject.SetActive(_FirstPayInfo.status != (int)ActivityM.SignState.Geted && isFirstDay);

        if (_FirstPayInfo.status == (int)ActivityM.SignState.Geted)
        {
            if (_FirstPayInfo.day == 3)
            {
                GetView().txt_Time.gameObject.SetActive(false);
            }
        }

        GetView().txt_Des.text = ShowActivityInfo.RuleDes;
    }

    private HTTP_ActivityChargeResponse _FirstPayInfo;
    public override void OnEnter(params object[] args)
    {
        _FirstPayInfo = (HTTP_ActivityChargeResponse)args[0];
        OnUIRefresh();
    }

    public override void OnExit()
    {
        RemoveTime();
    }

    public override void OnDestroy()
    {
        RemoveTime();
    }
}