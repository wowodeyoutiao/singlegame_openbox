using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using System;
using JEngine.Helper;
using static UIViewExpPool;
using JEngine.UI;
using pb;
using BestHTTP.SecureProtocol.Org.BouncyCastle.Utilities.Collections;
using JEngine.MyNet;


public class UICommonItemUse : UIBase<UIModelCommonItemUse, UIViewCommonItemUse>
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
            CloseUI();
        });

        GetView().btn_use.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UseItem();
        });

        GetView().btn_jian.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _UseCount--;
            _UseCount = Mathf.Clamp(_UseCount, _MinCount, _MaxCount);
            GetView().input_Use.text = _UseCount.ToString();
            //ShowConsume();
        });

        GetView().btn_add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _UseCount++;
            _UseCount = Mathf.Clamp(_UseCount, _MinCount, _MaxCount);
            GetView().input_Use.text = _UseCount.ToString();
            //ShowConsume();
        });

        GetView().btn_Max.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            _UseCount = _MaxCount;
            GetView().input_Use.text = _UseCount.ToString();
        });

        GetView().input_Use.onValueChanged.AddListener(v =>
        {
            if (int.TryParse(v, out int num))
                _UseCount = num;
            else
                _UseCount = 1;

            _UseCount = Mathf.Clamp(_UseCount, _MinCount, _MaxCount);
            GetView().input_Use.text = _UseCount.ToString();
        });
    }

    private int itemId;
    private int _UseCount = 1;
    private int _MinCount = 1;
    private int _MaxCount = 1;

    //ItemInfo itemInfo;

    /// <summary>
    /// 进入界面
    /// </summary>
    /// <param name="args">[0] cfgId</param>
    public override void OnEnter(params object[] args)
    {
        itemId = (int)args[0];

        int itemNum = Mathf.Max(1, (int)BagM.Instance.GetNumByCfgID(itemId));
        if(itemId == (int)CFG_Item.ItemID.JSJZ)
        {
            long boxUpgradeRemainingTime = RoleM.Instance.BoxNextRecoveryTime - GameUtil.Instance.GetCurrentTimeSeconds();
            _MaxCount = Mathf.Min(itemNum, (int)Mathf.Ceil((float)boxUpgradeRemainingTime / CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_accelerate));
        }
        else
        {
            _MaxCount = itemNum;

        }
        _UseCount = _MinCount;
        GetView().input_Use.text = null;
        GetView().input_Use.text = _UseCount.ToString();
    }

    private void UseItem()
    {
        // 无双宠物兑换卡
        if (itemId == (int)CFG_Item.ItemID.PetChangeCard)
        {
            var PetChangeCardItem = CFG_Summon.GetPetExchangeCard();
            if (_MaxCount < PetChangeCardItem.Consume[0].num)
                UITool.ShowTips($"{CFG_Item.GetById(itemId).Name}道具不足，无法兑换！");
            else
                HttpC.GetPetReceive(1, PetChangeCardItem.Drawtype);
        }
        else if(itemId == (int)CFG_Item.ItemID.JSJZ)
        {
            if (GameUtil.Instance.ComCoinToCheck(CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_itemId, _UseCount, () => { UITool.Close<UIBoxUpgrade>(); }, true))
            {
                if (StateM.Instance.isSendUpgradeOrSpeedUpBox)
                    return;

                TcpC.C2S_HomeUpgradeBoxQuick(_UseCount);
                StateM.Instance.isSendUpgradeOrSpeedUpBox = true;

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "宝箱界面",
                    position = "宝箱升级界面",
                    event_click = "宝箱升级加速",
                });
            }
        }
        // 其他道具
        else
        {
            long uid = BagM.Instance.GetItemByCFGID(itemId).Uid;
            TcpC.C2S_ItemUse(uid, _UseCount);
        }

        CloseUI();
    }

    private void CloseUI()
    {
        UITool.Close<UICommonItemUse>();
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}