using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using MetaFramework.UI;
using pb;
using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class UIBoxUpgrade : UIBase<UIModelBoxUpgrade, UIViewBoxUpgrade>
{
    /// <summary>
    /// 要显示的等级类型
    /// </summary>
    public enum QualityRateShowType
    {
        /// <summary>
        /// 显示当前等级
        /// </summary>
        CurrentLevel,

        /// <summary>
        /// 显示下一等级
        /// </summary>
        NextLevel,
    }

    /// <summary>
    /// 宝箱升级状态
    /// </summary>
    public enum UpgradeState
    {
        /// <summary>
        /// 宝箱可升级
        /// </summary>
        CanUpgrade,

        /// <summary>
        /// 宝箱升级中
        /// </summary>
        Upgrading,

        /// <summary>
        /// 道具不足，宝箱不可升级
        /// </summary>
        CannotUpgrade,
    }

    /// <summary>
    /// 当前等级宝箱配置
    /// </summary>
    CFGItem_BoxLevel cfgBoxLevelItem;

    /// <summary>
    /// 宝箱升级状态
    /// </summary>
    UpgradeState upgradeState = UpgradeState.CannotUpgrade;

    /// <summary>
    /// 宝箱最大等级
    /// </summary>
    int maxLevel;

    /// <summary>
    /// 宝箱升级剩余所用时间
    /// </summary>
    //long boxUpgradeRemainingTime;

    //const string tip_CannotUpgrade = "角色等级不足，无法升级！";

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Upgrade.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnUpgrade();
        }); 
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIBoxUpgrade>();
        });
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnSync_SourceRefresh);
        UniGameEvent.AddEventListener<GameEvent.OnApplicationFocus>(OnSync_FocusChanged);

        maxLevel = CFG_BoxLevel.GetMaxLevel();
        RefreshUI();
    }

    private void RefreshUI()
    {
        cfgBoxLevelItem = CFG_BoxLevel.GetById(RoleM.Instance.BoxNowlevel);
        GetUpgredeState();

        ShowBoxLevelAndLimit();
        ShowQualityRate();
        ShowRoleLevelLimitTip();
        ShowButtonNameAndConsume();
        ShowUpgradeTime();
    }

    private void GetUpgredeState()
    {
        if (RoleM.Instance.BoxNextRecoveryTime > 0) // 升级中
        {
            upgradeState = UpgradeState.Upgrading;
        }
        else if (GameUtil.Instance.ComCoinToCheck(cfgBoxLevelItem.Consume[0].id, cfgBoxLevelItem.Consume[0].num, null, false)) // 检查道具是否够升级使用
        {
            if (upgradeState == UpgradeState.Upgrading)
            {
                PlayUpEffect();
            }
            upgradeState = UpgradeState.CanUpgrade;
        }
        else
        {
            if (upgradeState == UpgradeState.Upgrading)
            {
                PlayUpEffect();
            }
            upgradeState = UpgradeState.CannotUpgrade;
        }
    }

    /// <summary>
    /// 显示宝箱等级和开箱次数限制
    /// </summary>
    private void ShowBoxLevelAndLimit()
    {
        bool isMax = RoleM.Instance.BoxNowlevel >= maxLevel;
        int nextLv = isMax ? RoleM.Instance.BoxNowlevel : RoleM.Instance.BoxNowlevel + 1;

        GetView().txt_CurrentLevelValue.text = "当前等级：" + RoleM.Instance.BoxNowlevel;
        GetView().txt_NextLevelValue.text = "下一等级：" + nextLv;
        GetView().txt_CurrentOpenBoxLimit.text = cfgBoxLevelItem.Maxtime.ToString();
        GetView().txt_NextOpenBoxLimit.text = isMax ? cfgBoxLevelItem.Maxtime.ToString() : CFG_BoxLevel.GetById(nextLv).Maxtime.ToString();

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.FirstPay.ToString(), 0);
    }

    /// <summary>
    /// 显示宝箱品质概率
    /// </summary>
    private void ShowQualityRate()
    {
        ShowQualityRateItem(QualityRateShowType.CurrentLevel, cfgBoxLevelItem.qualityrandom);
        if (RoleM.Instance.BoxNowlevel < maxLevel)
        {
            CFGItem_BoxLevel.QualityRandomItem[] nextLvRateItems = CFG_BoxLevel.GetById(RoleM.Instance.BoxNowlevel + 1).qualityrandom;
            ShowQualityRateItem(QualityRateShowType.NextLevel, nextLvRateItems);
        }
        else 
        {
            ShowQualityRateItem(QualityRateShowType.NextLevel, cfgBoxLevelItem.qualityrandom);
        }
    }

    private void ShowQualityRateItem(QualityRateShowType qualityRateShowType, CFGItem_BoxLevel.QualityRandomItem[] rateItems)
    {
        switch (qualityRateShowType)
        {
            case QualityRateShowType.CurrentLevel:
                for (int i = 0; i < GetView().txtList_CurrentQualityRandom.Count; i++)
                {
                    GetView().txtList_CurrentQualityRandom[i].text = GetRateString(i, rateItems);
                }
                break;
            case QualityRateShowType.NextLevel:
                for (int i = 0; i < GetView().txtList_NextQualityRandom.Count; i++)
                {
                    GetView().txtList_NextQualityRandom[i].text = GetRateString(i, rateItems);
                }
                break;
            default:
                break;
        }
    }

    private string GetRateString(int index, CFGItem_BoxLevel.QualityRandomItem[] qualityRandoms)
    {
        int totalRate = 0;
        foreach (var item in qualityRandoms)
        {
            totalRate += item.rate;
        }
        float rate = (qualityRandoms[index].rate / (float)totalRate) * 100;
        return rate.ToString("0.00") + "%";
    }

    TimeItem EggTime;

    private void ShowUpgradeTime()
    {
        ClearEggTime();
        int needTime = cfgBoxLevelItem.Needtime;
        var mainBattleCfg = CFG_MainBattle.GetById(cfgBoxLevelItem.BattleID);

        //显示配置的时间
        if (cfgBoxLevelItem.BattleID <= RoleM.Instance.MainBattleId)
            GetView().txt_UpgradeCond.text = $"<color=#FFFFFF>{mainBattleCfg.name}</color>";
        else
            GetView().txt_UpgradeCond.text = $"<color=#FF0000>{mainBattleCfg.name}</color>";

        switch (upgradeState)
        {
            case UpgradeState.CanUpgrade:
            case UpgradeState.CannotUpgrade:
                break;

            case UpgradeState.Upgrading: //显示倒计时
                ShowTimeRemaining();
                EggTime = GameHelper.Instance.AddInvoke(new TimeItem(-1, 1, ShowTimeRemaining, true, true));
                break;
            default:
                break;
        }

        void ShowTimeRemaining()
        {
            // 显示需要消耗的加速道券数量
            ShowButtonNameAndConsume();
        }

    }

    private void ClearEggTime()
    {
        if (GameHelper.Instance != null && EggTime != null)
        {
            GameHelper.Instance.RemoveInvoke(EggTime);
            EggTime = null;
        }
    }

    private void ShowRoleLevelLimitTip()
    {
        GetView().txt_CannotUpgrade.text = IsRoleLvEnough() ? null : $"角色等级不足，需要达到{cfgBoxLevelItem.Rolelevel}级。";
    }

    private bool IsRoleLvEnough()
    {
        return RoleM.Instance.Level >= cfgBoxLevelItem.Rolelevel;
    }

    /// <summary>
    /// 需要消耗的加速卷轴数量
    /// </summary>
    /// <returns></returns>
    private int GetConsume_JSJZ()
    {
        long seconds = RoleM.Instance.BoxNextRecoveryTime - GameUtil.Instance.GetCurrentTimeSeconds();
        return (int)Mathf.Ceil((float)seconds / CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_accelerate);
    }

    /// <summary>
    /// 显示升级宝箱需要的物品数量
    /// </summary>
    private void ShowButtonNameAndConsume()
    {
        if (upgradeState == UpgradeState.Upgrading)
        {
            GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_itemId, GetView().img_Coin);
            int needItemNum = GetConsume_JSJZ();
            GetView().txt_Coin.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(needItemNum), CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_itemId, needItemNum) + " （持有：" + GameUtil.Instance.FormatNumWithUnit(RoleM.Instance.BoxQuick) + "）";
            GetView().txt_UpgradeButtonName.text = "加速";
            GetView().txt_ExchangeTime.text = $"（1张加速券={CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_accelerate / 60}分钟加速）";
        }
        else
        {
            GameUtil.Instance.SetItemIcon(cfgBoxLevelItem.Consume[0].id, GetView().img_Coin);
            bool isEnough = GameUtil.Instance.ComCoinToCheck(cfgBoxLevelItem.Consume[0].id, cfgBoxLevelItem.Consume[0].num, null, false);
            var str = GameUtil.Instance.FormatNumWithUnit(cfgBoxLevelItem.Consume[0].num);
            GetView().txt_Coin.text = !isEnough ? $"<color=#F64747>{str}</color>" : $"<color=#40B023>{str}</color>";
            GetView().txt_UpgradeButtonName.text = "升级";
        }
        GetView().txt_ExchangeTime.gameObject.SetActive(upgradeState == UpgradeState.Upgrading);

    }

    Coroutine coroutineEffPlay;
    private void PlayUpEffect()
    {
        coroutineEffPlay = CoroutineMgr.Instance.StartCoroutine(PlayUpLevelEffect());
    }

    private void StopPlayEffect()
    {
        if (coroutineEffPlay != null)
        {
            CoroutineMgr.Instance.StopCoroutine(coroutineEffPlay);
            coroutineEffPlay = null;
        }

        for (int i = 0; i < GetView().effect_UpLevel.Count; i++)
        {
            GetView().effect_UpLevel[i].SetActive(false);
        }
    }

    WaitForSeconds waitTimeEffectPlay = new WaitForSeconds(0.2f);
    private IEnumerator PlayUpLevelEffect()
    {
        for (int i = 0; i < GetView().effect_UpLevel.Count; i++)
        {
            GetView().effect_UpLevel[i].SetActive(false);
            GetView().effect_UpLevel[i].SetActive(true);
            yield return waitTimeEffectPlay;
        }
        yield return new WaitForSeconds(2);
        StopPlayEffect();
    }

    private void OnUpgrade()
    {
        FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<CFG_XMLConfig.XML_Config.SwitchOpne>() { CFG_XMLConfig.XML_Config.SwitchOpne.OpenBox });

        if (RoleM.Instance.BoxNowlevel >= maxLevel)
        {
            UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack()
            {
                showStr = "祝福等级已达到最大，无需升级。",
            });
            return;
        }

        if (!IsRoleLvEnough())
        {
            UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack()
            {
                showStr = $"角色等级不足，需要达到{cfgBoxLevelItem.Rolelevel}级。",
            });
            return;
        }

        if(cfgBoxLevelItem.BattleID > RoleM.Instance.MainBattleId)
        {
            var cfg = CFG_MainBattle.GetById(cfgBoxLevelItem.BattleID);
            UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack()
            {
                showStr = $"主线关卡进度不满足条件，请先通关{cfg.name}再试噢!",
            });
            return;
        }

        // 升级中, 点击可加速
        if (upgradeState == UpgradeState.Upgrading)
        {
            //int needItemNum = GetConsume_JSJZ();
            //if (GameUtil.Instance.ComCoinToCheck(CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_itemId, needItemNum, () => { UITool.Close<UIBoxUpgrade>(); }, true))
            //{
            //    TcpC.C2S_HomeUpgradeBoxQuick(needItemNum);
            //}
            UITool.Open<UICommonItemUse>(null, CFG_XMLConfig.xmlConfig.BaseValue.BoxLevelUp_itemId);
        }
        else // 检查是否可以升级
        {
            if (GameUtil.Instance.ComCoinToCheck(cfgBoxLevelItem.Consume[0].id, cfgBoxLevelItem.Consume[0].num, null, false))
            {
                if (StateM.Instance.isSendUpgradeOrSpeedUpBox)
                    return;

                TcpC.C2S_HomeUpgradeBox();
                StateM.Instance.isSendUpgradeOrSpeedUpBox = true;

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "宝箱界面",
                    position = "宝箱升级界面",
                    event_click = $"升级宝箱",
                });
            }
            else
            {
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = "金币不足，是否前往商店购买金币？",
                    YesFun = () => {
                        UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                    },
                });
            }
        }
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.cid == CmdId.C2S_HomeUpgradeBoxQuick)
        {
            StateM.Instance.isSendUpgradeOrSpeedUpBox = false;
            if (_e.code == 0)
            {
                ShowButtonNameAndConsume();
            }
        }

        if (_e.cid == CmdId.C2S_HomeUpgradeBox)
        {
            StateM.Instance.isSendUpgradeOrSpeedUpBox = false;
        }
    }

    [Subscriber(ThreadMode.Main)]
    // 更新时间状态
    private void OnSync_UserBase(GameEvent.Sync_UserBoxInfo _e)
    {
        if (isNullGo) return;
        RefreshUI();
    }

    [Subscriber(ThreadMode.Main)]
    // 更新时间状态
    private void OnSync_FocusChanged(GameEvent.OnApplicationFocus _e)
    {
        if (isNullGo || !_e.focus) return;
        RefreshUI();
    }

    private void OnSync_SourceRefresh(Event_Source_Refresh _e)
    {
        if (isNullGo)
            return;

        ShowButtonNameAndConsume();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnSync_SourceRefresh);
        UniGameEvent.RemoveEventListener<GameEvent.OnApplicationFocus>(OnSync_FocusChanged);

        ClearEggTime();
        StopPlayEffect();
        cfgBoxLevelItem = null;
        StateM.Instance.isSendUpgradeOrSpeedUpBox = false;
    }

    public override void OnDestroy()
    {
        ClearEggTime();
    }
}