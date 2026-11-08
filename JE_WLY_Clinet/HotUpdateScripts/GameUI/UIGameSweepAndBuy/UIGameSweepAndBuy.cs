using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;

public class UIGameSweepAndBuy : UIBase<UIModelGameSweepAndBuy, UIViewGameSweepAndBuy>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().BtnClose.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGameSweepAndBuy>();
        });

        GetView().btnAdd.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnNumChang(true);
        });

        GetView().btnJian.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnNumChang(false);
        });

        GetView().btnMax.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnNumMaxOrMin(true);
        });

        GetView().btnMin.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnNumMaxOrMin(false);
        });

        GetView().inputField.onValueChanged.AddListener((string value) =>
        {
            int num = 1;
            if(int.TryParse(value, out int _num))
                num = _num;

            if (num > MaxNum)
                GetView().inputField.text = MaxNum.ToString();
            else if (num <= 1) 
                GetView().inputField.text = "1";
        } );

        GetView().btnyes.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if(opType == OpType.WuXingSweep)
            {
                TcpC.C2S_ElementSweepForWardRequest(int.Parse(GetView().inputField.text));
                UITool.Close<UIGameSweepAndBuy>();
                ShuShuClickSend();
                return;
            }

            if (BattleID == 0)
            {
                //int shopID = opType == OpType.RelicBuy ? (int)CFG_Store.StroeID.RelicChange : (int)CFG_Store.StroeID.BossChange;
                int shopID = opType switch
                {
                    OpType.RelicBuy => (int)CFG_Store.StroeID.RelicChange,
                    OpType.BossBuy => (int)CFG_Store.StroeID.BossChange,
                    _ => 0
                };

                int num = int.Parse(GetView().inputField.text);
                if (!GameUtil.Instance.ComCoinToCheck(shopID, delegate { UITool.Close<UIGameSweepAndBuy>(); },true, num))
                {
                    return;
                }
                GameUtil.Instance.GameToBuy(shopID, num, () =>
                {
                    GameHelper.Instance.SetDeayTimeDo(delegate {
                        TcpC.C2S_InstanceInfo();
                    },1.5f);
                    
                });
            }
            else
            {
                TcpC.C2S_CopySweep(BattleID, int.Parse(GetView().inputField.text));
                ShuShuClickSend();
            }
            UITool.Close<UIGameSweepAndBuy>();
        });
    }

    public enum OpType 
    {
        RelicSweep,
        RelicBuy,
        BossSweep,
        BossBuy,
        WuXingSweep,
    }

    private void ShuShuClickSend()
    {
        if(opType == OpType.BossSweep)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = "挑战",
                position = "首领",
                event_click = $"扫荡",
            });
        }
        else if (opType == OpType.RelicSweep)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = "挑战",
                position = "遗迹",
                event_click = $"扫荡",
            });
        }
        else if (opType == OpType.WuXingSweep)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = "五行试炼界面",
                position = "",
                event_click = $"单层扫荡",
            });
        }

    }

    private void OnNumMaxOrMin(bool isMax)
    {
        int maxNum = MaxNum;

        if(opType == OpType.RelicBuy)
        {
            var cfg = CFG_Store.GetById((int)CFG_Store.StroeID.RelicChange);
            maxNum = Mathf.Min((int)(BagM.Instance.GetNumByCfgID(cfg.Limite.needmoney) / cfg.Limite.Num), MaxNum);
        }
        else if(opType == OpType.BossBuy)
        {
            var cfg = CFG_Store.GetById((int)CFG_Store.StroeID.BossChange);
            maxNum = Mathf.Min((int)(BagM.Instance.GetNumByCfgID(cfg.Limite.needmoney) / cfg.Limite.Num), MaxNum);
        }

        int num = isMax ? maxNum : 1;
        num = num == 0 ? 1 : num;
        GetView().inputField.text = num.ToString();
        OnUsedNumShow();
    }

    private void OnNumChang(bool isAdd)
    {
        int num = int.Parse(GetView().inputField.text);
        if (isAdd && num < MaxNum)
        {
            num ++;
        }
        else if (!isAdd && num > 1)
        {
            num -- ;
        }
        GetView().inputField.text = num.ToString();
        OnUsedNumShow();
    }

    private int MaxNum = 0;

    private int BattleID = 0;
    OpType opType ;
    public override void OnEnter(params object[] args)
    {
        opType = (OpType)args[0];
        if (opType == OpType.RelicSweep)
        {
            MaxNum = FightPlayM.Instance.GetInfo().relic_map.ChallengeTime;
            GetView().txt_des.text = "可扫荡次数：" + MaxNum;
            GetView().txt_btn.text = "扫荡";
            BattleID = (int)args[1];
        }
        else if (opType == OpType.BossSweep)
        {
            MaxNum = FightPlayM.Instance.GetInfo().boss_map.ChallengeTime;
            GetView().txt_des.text = "可扫荡次数：" + MaxNum;
            GetView().txt_btn.text = "扫荡";
            BattleID = (int)args[1];
        }
        else if (opType == OpType.WuXingSweep)
        {
            MaxNum = TurnBasedM.Instance.elementTimes;
            GetView().txt_des.text = "可扫荡次数：" + MaxNum;
            GetView().txt_btn.text = "扫荡";
            BattleID = 0;
        }
        else if (opType == OpType.RelicBuy)
        {
            BattleID = 0;
            MaxNum = StoreM.Instance.GetCanBuyNum((int)CFG_Store.StroeID.RelicChange);
            GetView().txt_des.text = "今日可购次数：" + MaxNum;
            GetView().txt_btn.text = "购买";
        }
        else if (opType == OpType.BossBuy)
        {
            BattleID = 0;
            MaxNum = StoreM.Instance.GetCanBuyNum((int)CFG_Store.StroeID.BossChange);
            GetView().txt_des.text = "今日可购次数：" + MaxNum;
            GetView().txt_btn.text = "购买";
        }

        bool showUse = opType != OpType.RelicSweep && opType != OpType.BossSweep && opType != OpType.WuXingSweep;
        GetView().obj_used.SetActive(showUse);
        GameHelper.Instance.SetDeayFrame(() => {
            OnNumMaxOrMin(true);
        });

    }

    private void OnUsedNumShow()
    {
        var cfgStroe = opType switch
        {
            OpType.RelicBuy => CFG_Store.GetById((int)CFG_Store.StroeID.RelicChange),
            OpType.BossBuy => CFG_Store.GetById((int)CFG_Store.StroeID.BossChange),
            _ => null,
        };

        if (cfgStroe == null) return;
        GameUtil.Instance.SetItemIcon(cfgStroe.Limite.needmoney, GetView().img_used);
        GetView().txt_usedNum.text = GameUtil.Instance.GetCheakNumStr(cfgStroe.Limite.needmoney, cfgStroe.Limite.Num * int.Parse(GetView().inputField.text));
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}