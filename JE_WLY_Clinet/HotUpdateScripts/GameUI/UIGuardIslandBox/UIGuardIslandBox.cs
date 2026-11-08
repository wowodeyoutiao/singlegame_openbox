using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using pb;
using static UIFightFormation;

public class UIGuardIslandBox : UIBase<UIModelGuardIslandBox, UIViewGuardIslandBox>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGuardIslandBox>();
        });

        for (int i = 0; i < GetView().infoItems.Length; i++)
        {
            int idx = i+1;
            GetView().infoItems[i].toggle.onValueChanged.AddListener((isOn) => {
                if (isOn)
                    OnTooggleChange(idx);
            });
        }

        OnUIShow();

        GetView().Btn_Yes.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (ToggleIdx == 1)
            {
                TcpC.C2S_GuardIslandMultiplier();
            }
            else
            {
                if (GameUtil.Instance.GetNowToTimeSecoud(1, 0, 00, 00) < 1800)
                {
                    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                    {
                        ShowStr = $"是否消耗距离下一天时间不足30分钟，确认花费对应的资源进入挂机副本吗？",
                        YesFun = () => {
                            OnUsedPlay();
                        },
                    });
                }
                else
                    OnUsedPlay();

                void OnUsedPlay()
                    {
                        void OnPlay()
                        {
                            TcpC.GetGameJoin(CFG_Room.GetByRoomType(CFG_Room.RoomType.PVE_Peace, RoleM.Instance.Level).ID, 0);
                            UITool.Close<UIGuardIslandBox>();
                        }
                        var cfgInfo = CFG_XMLConfig.xmlConfig.GuardIslandDic[ToggleIdx];
                        if (cfgInfo.type == 3)
                        {
                            var itemInfo = CFG_Store.GetById(cfgInfo.id);
                            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                            {
                                ShowStr = $"是否消耗{itemInfo.Limite.Num}{GameUtil.Instance.GetItemName(itemInfo.Limite.needmoney)}",
                                YesFun = () =>
                                {
                                    if (!GameUtil.Instance.ComCoinToCheck(cfgInfo.id, () => { UITool.Close<UIGuardIslandBox>(); }))
                                        return;

                                    GameUtil.Instance.GameToBuy(cfgInfo.id, OnPlay);
                                },
                            });
                        }
                        else
                            GameUtil.Instance.GameToBuy(cfgInfo.id, OnPlay);
                    }
                
            }
        });
    }
    private int ToggleIdx = 0;
    private void OnTooggleChange(int idx)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        ToggleIdx = idx;
    }

    private void OnUIShow()
    {
        for (int i = 0; i < GetView().infoItems.Length; i++)
        {
            var cfgInfo = CFG_XMLConfig.xmlConfig.GuardIslandDic[i+1];
            GetView().infoItems[i].go.SetActive((!GameM.Instance.IsAudit && cfgInfo.type == 5 )|| cfgInfo.type != 5);
            if (cfgInfo.type != 1)
            {
                var cfg_Store = CFG_Store.GetById(cfgInfo.id);
                GameUtil.Instance.SetShopItemIcon(cfg_Store.Limite.needmoney, GetView().infoItems[i].img_icon);
                GetView().infoItems[i].txt.text = cfgInfo.type == 2 ? "广告" : cfg_Store.Limite.Num.ToString();
            }
            GetView().infoItems[i].txt_des.text = cfgInfo.multiplier + "倍效率";
        }
    }
    public override void OnEnter(params object[] args)
    {
        GetView().infoItems[0].toggle.isOn = true;
        ToggleIdx = 1; // 默认选中第一个
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}