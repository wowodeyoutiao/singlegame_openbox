using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Ocsp;
using Cysharp.Threading.Tasks;
using DG.Tweening;

using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using System;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;
using static CFG_XMLConfig.XML_Config;
using static GameEvent;

public partial class UIHall : UIBase<UIModelHall, UIViewHall>
{
    private UIViewHall uIViewHall;

    public UIHallBg farHallBg;
    public UIHallBg midHallBg;
    public UIHallBg nearbyHallBg;

    private UICom_TopRes topUI;
    private UICom_DownBtns downUI = null;

    private int equipEffectID = 0;
    public Dictionary<int, GameObject> equipEffectsDic = new Dictionary<int, GameObject>();

    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public override void OnInit()
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(RefreshTopActiveButton);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnSync_PartnetChanged);
        UniGameEvent.AddEventListener<GameEvent.Refresh_Privilege>(RefreshPrivilege);
        UniGameEvent.AddEventListener<GameEvent.Guide_Jump>(GuideJump);
        UniGameEvent.AddEventListener<GameEvent.Refresh_HomeOpen>(OnRefreshHomeBox);
        UniGameEvent.AddEventListener<GameEvent.ChangeEquip>(Event_ChangeEquip);
        UniGameEvent.AddEventListener<GameEvent.OpenBoxDisassemblyEquip>(DisassemblyEffectFlyAsync);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_SyncUserBasePush);
        UniGameEvent.AddEventListener<GameEvent.Refresh_LvUp>(OnSync_Refresh_LvUp);
        HallBattleM.Instance.AddEventListener();

        uIViewHall = GetView();

        OnInitBg();

        OnInitHallBox();
        OnInitHallBattle();

        OnInitSummaryBtns();
        InitTopBtns();
        InitEquip();
        InitTask();
        OnBindRed();
        InitOpenBoxSetting();

        // 宠物
        ButtonAddListener(uIViewHall.petTrans.GetComponent<Button>(), () =>
        {
            if (uIViewHall.petTrans.childCount == 0) return;
            UITool.Open<UIPet>();

        });
        // 聊天
        ButtonAddListener(uIViewHall.chatTrans.GetComponent<Button>(), () => { UITool.Open<UIChat>(); });
        ButtonAddListener(uIViewHall.chatBtn, () => { UITool.Open<UIChat>(); });
    }

    private void OnInitBg()
    {
        if (farHallBg != null) return;

        farHallBg = uIViewHall.bgFarNodeTrans.gameObject.AddComponent<UIHallBg>();
        //midHallBg = uIViewHall.bgMidNodeTrans.gameObject.AddComponent<UIHallBg>();
        nearbyHallBg = uIViewHall.bgNearbyNodeTrans.gameObject.AddComponent<UIHallBg>();
        farHallBg.OnInit("bg_Scene_smyj_", "bg_Scene_bossyj_", 0.4f);
        //midHallBg.OnInit("bg_Scene_smzj_", "bg_Scene_bosszj_", 0.7f);
        nearbyHallBg.OnInit("bg_Scene_smqj_", "bg_Scene_bossqj_", 0.7f);
    }

    #region 收缩
    private void OnInitSummaryBtns()
    {
        uIViewHall.summaryTrans.GetComponent<Button>().onClick.AddListener(OnClickSummaryBtn);

        var btnTran = uIViewHall.summaryTrans.Find("Btns");
        for (int i = 1; i <= btnTran.childCount; i++)
        {
            var obj = btnTran.Find("Btn_" + i);
            if (obj != null)
            {
                int idx = i;
                ButtonAddListener(obj.GetComponent<Button>(), () => { OnClickSummaryBtnAtIndex(idx); });
            }
        }
    }

    private void OnClickSummaryBtnAtIndex(int idx)
    {
        switch (idx)
        {
            case 1:
                UITool.Open<UITask>();
                break;
            case 2:
                UITool.Open<UIRank>();
                break;
            case 3:
                UITool.Open<UISpecialOffer>();
                break;
            case 4:
                UITool.Open<UIGiftBag>();
                break;
            default:
                UITool.ShowTips("功能暂未开放");
                break;
        }
    }

    private void OnClickSummaryBtn()
    {
        var go = uIViewHall.summaryTrans.gameObject;
        go.SetActive(!go.activeSelf);

        if (go.activeSelf)
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_OpenRightPanel);
        else
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_CloseRightPanel);
    }
    #endregion

    #region 开箱指定

    private int openBoxSlotSelect;

    private void InitOpenBoxSetting()
    {
        uIViewHall.boxOpenSetBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            bool isActive = uIViewHall.rootBoxSet.gameObject.activeSelf;
            uIViewHall.rootBoxSet.gameObject.SetActive(!isActive);
            if (!isActive) Sync_OpenBoxSlotSelect();
        });

        uIViewHall.boxOpenSetSaveBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            uIViewHall.rootBoxSet.gameObject.SetActive(false);
            GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.DefaultOpenSoltCondition.ToString() + RoleM.Instance.Uid, openBoxSlotSelect);
        });

        uIViewHall.boxOpenSetMaskBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            uIViewHall.rootBoxSet.gameObject.SetActive(false);
        });

        uIViewHall.boxOpenSetCloseBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            uIViewHall.rootBoxSet.gameObject.SetActive(false);
        });

        for (int i = 0; i < uIViewHall.boxSetToggles.Length; i++)
        {
            int idx = i;
            uIViewHall.boxSetToggles[i].onValueChanged.AddListener(isOn =>
            {
                if (isOn)
                {
                    if (GameUtil.Instance.SlotIsUnLockOrHasRole(idx))
                    {
                        openBoxSlotSelect = idx;
                    }
                    else
                        SetActiveToggle(openBoxSlotSelect);
                }
            });
        }
    }

    private void SetActiveToggle(int index)
    {
        for (int i = 0; i < GetView().boxSetToggles.Length; i++)
        {
            if (i == index) GetView().boxSetToggles[i].isOn = true;
        }
    }

    private void Sync_OpenBoxSlotSelect()
    {
        openBoxSlotSelect = GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.DefaultOpenSoltCondition.ToString() + RoleM.Instance.Uid);
        SetActiveToggle(openBoxSlotSelect);
    }
    #endregion

    #region red point
    private Dictionary<RedPointConst, GameObject> RedDic;
    private void OnBindRed()
    {
        RedDic = new Dictionary<RedPointConst, GameObject>()
        {
            // 大厅右上方活动按钮
            [RedPointConst.SpecialOffer] = uIViewHall.specialOfferBtnTrans.Find("dot").gameObject, // 特惠
            [RedPointConst.GiftBag] = uIViewHall.giftBagBtnTrans.Find("dot").gameObject, // 礼包
            [RedPointConst.ActiveCenter] = uIViewHall.sevenDayBtnTrans.Find("dot").gameObject, // 开服七天乐
            [RedPointConst.FirstPay] = uIViewHall.firstPayBtnTrans.Find("dot").gameObject, // 首充

            [RedPointConst.Task] = uIViewHall.dailyBtnTrans.Find("dot").gameObject, // 日常
            [RedPointConst.Rank] = uIViewHall.rankBtnTrans.Find("dot").gameObject, // 排行榜
            [RedPointConst.Summary] = uIViewHall.summaryBtnTrans.Find("dot").gameObject, // 收缩
            // 大厅右下方功能按钮
            //[RedPointConst.WuXing] = uIViewHall.funcBtnNodeTrans.GetChild(0).Find("RedImg").gameObject
        };

        foreach (var item in RedDic)
        {
            RedPointManager.Instance.RPS.Register(item.Key.ToString(), (node) =>
            {
                item.Value.SetActive(node.pointCount > 0);
            });
        }

        RedPointManager.Instance.RPS.Register(RedPointConst.BoxLevel.ToString(), (node) =>
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.OpenBox, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.OpenBox, false))
                return;
            uIViewHall.goBlessRedDot?.SetActive(node.pointCount > 0);
            //uIViewHall.boxBtn.transform.Find("redDotRoot/RedImg").gameObject.SetActive(node.pointCount > 0);
        });
    }

    private void ClearBindRed()
    {
        if (RedDic != null)
        {
            foreach (var item in RedDic)
            {
                RedPointManager.Instance.RPS.Register(item.Key.ToString(), null);
            }
        }

        RedPointManager.Instance.RPS.Register(RedPointConst.Tower_Get.ToString(), null);
    }

    private void CheckRed()
    {
        RoleM.Instance.CheckBoxLevelRed();
    }
    #endregion

    private void ButtonAddListener(Button btn, Action act)
    {
        btn.onClick.RemoveAllListeners();
        btn.onClick.AddListener(() =>
        {
            if (!btn.gameObject.name.Equals("Btn_1"))
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            act?.Invoke();
        });
    }


    #region Top btns
    private void InitTopBtns()
    {
        // 左上角按钮
        for (int i = 1; i <= uIViewHall.topLeftBtnsTrans.childCount; i++)
        {
            var obj = uIViewHall.topLeftBtnsTrans.Find("Btn_" + i);
            if (obj != null)
            {
                int idx = i;
                ButtonAddListener(obj.GetComponent<Button>(), () => { OnTopLeftBtnClick(idx); });
            }
        }

        // 右上角按钮
        for (int i = 1; i <= uIViewHall.topRightBtnsTrans.childCount; i++)
        {
            var obj = uIViewHall.topRightBtnsTrans.Find("Btn_" + i);
            if (obj != null)
            {
                int idx = i;
                ButtonAddListener(obj.GetComponent<Button>(), () => { OnTopRightBtnClick(idx); });
            }
        }
    }

    /// <summary>
    /// 左上按钮点击
    /// </summary>
    /// <param name="idx"></param>
    private void OnTopLeftBtnClick(int idx)
    {
        switch (idx)
        {
            case 1:
                UITool.Open<UIActivityEntrance>();
                break;
            case 2:
                UITool.Open<UICompetitionEntrance>();
                break;
            default:
                break;
        }
    }

    /// <summary>
    /// 右上按钮点击
    /// </summary>
    /// <param name="idx"></param>
    private void OnTopRightBtnClick(int idx)
    {
        switch (idx)
        {
            case 1:
                OnClickSummaryBtn();
                uIViewHall.summaryTrans.SetAsLastSibling();
                break;
            case 2:
                UITool.Open<UIActive>();
                break;
            case 3:
                OpenFirstPay();
                break;
            case 4:
                UITool.Open<UILimitetRank>(null, CFG_RankingAward.RankClass.MainPass);
                break;
            case 5:
                UITool.Open<UILimitetRank>(null, CFG_RankingAward.RankClass.Fight);
                break;
            case 6:
                UITool.Open<UIBattleSupportGift>();
                break;
            default:
                UITool.ShowTips("功能暂未开放");
                break;
        }
    }

    private void OpenFirstPay()
    {
        if (!ActivityM.Instance.GetAcitvityState(ActivityType.FirstPay))
            UITool.ShowTips("首充活动已结束");

        HttpC.GetFirstPayInfo((req) =>
        {
            if (req.code == 0 && !isNullGo)
                UITool.Open<UIFirstPay>(null, req);
        });
    }

    #endregion

    #region Top UI
    private void LoadTopUI()
    {
        UITool.LoadUII<UICom_TopRes>(uIViewHall.transform, (a) =>
        {
            topUI = a as UICom_TopRes;
            ShowUI();
            topUI.OnSetLayer(_canvasOrder);
            topUI.OnHeadOnShow(curRoleIndex);

        });
    }
    #endregion

    #region Down UI
    private void LoadDownUI()
    {
        UITool.LoadUII<UICom_DownBtns>(uIViewHall.transform, (ui) =>
        {
            downUI = ui as UICom_DownBtns;

            //注册UI，用于新功能解锁时，移动到该位置
            FuncOpenShowWidget.Instance.RegisterFlyToTarget(EMainUIObject.Devlop, downUI.GetView().Btns[1].trans);   //养成
            FuncOpenShowWidget.Instance.RegisterFlyToTarget(EMainUIObject.FightChoice, downUI.GetView().Btns[0].trans);//玩法
            FuncOpenShowWidget.Instance.RegisterFlyToTarget(EMainUIObject.PartnerCall, downUI.GetView().Btns[2].trans); //召唤
            FuncOpenShowWidget.Instance.RegisterFlyToTarget(EMainUIObject.LimitedTimeActivity, uIViewHall.topLeftBtnsTrans.Find("Btn_1")); //限时活动
            FuncOpenShowWidget.Instance.RegisterFlyToTarget(EMainUIObject.Competion, uIViewHall.topLeftBtnsTrans.Find("Btn_2")); //赛事

            FuncOpenShowWidget.Instance.RegisterFlyToTarget(EMainUIObject.Zhufu, uIViewHall.boxBtn.transform); //祝福
        });
    }

    private void CheckStartGuidePartnetUp(UICom_DownBtns uiHallDown)
    {
        CheckPartnerSoltIsUnlock(1);
        IsLockToUnlock();

        bool CheckPartnerSoltIsUnlock(int soltIdx)
        {
            bool isLock = CFG_XMLConfig.CheckPartnerSoltIsOpenServer(soltIdx, false);
            if (!isLock)
                isLock = CFG_XMLConfig.GetOpenSoltLv(soltIdx) > RoleM.Instance.Level;

            return !isLock;
        }

        bool IsLockToUnlock()
        {
            return RoleM.Instance.guidePartnerRecordLevel < CFG_XMLConfig.GetOpenSoltLv(1);
        }
    }
    #endregion

    #region Equip
    private int curRoleIndex = 0;
    public void InitEquip()
    {
        // 左边按钮
        for (int i = 0; i < uIViewHall.equipTabBtnItems.Length; i++)
        {
            int idx = i;
            ButtonAddListener(uIViewHall.equipTabBtnItems[i].btn, () => { OnEquipTabBtnClick(idx); });
        }

        // 装备
        for (int i = 0; i < uIViewHall.equipItems.Length; i++)
        {
            int dressidx = i + 1;
            uIViewHall.equipItems[i].btn.onClick.RemoveAllListeners();
            uIViewHall.equipItems[i].btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (curRoleIndex == 0)
                {
                    bool isNullEquip = RoleM.Instance.CheckIsNullEquip((EquipSlot)dressidx);
                    if (!isNullEquip)
                    {
                        var ItemInfo = RoleM.Instance.GetSlotInfoBySolt((EquipSlot)dressidx).Item;
                        var equipRefineData = RoleM.Instance.GetOneEquipSlotData((CharacterSlot)(curRoleIndex + 1), (EquipSlot)dressidx);
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params
                        {
                            cfgId = ItemInfo.CfgId,
                            equipInfo = ItemInfo.Equip,
                            characterIndex = curRoleIndex + 1,
                            equipSlotPos = dressidx,
                            equipRefineData = equipRefineData,
                        });
                    }
                }
                else
                {
                    var ItemInfo = PartnerM.Instance.GetOnePartnerOneEquip(curRoleIndex, dressidx);
                    if (ItemInfo != null)
                    {
                        var equipRefineData = RoleM.Instance.GetOneEquipSlotData((CharacterSlot)(curRoleIndex + 1), (EquipSlot)dressidx);
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params
                        {
                            cfgId = ItemInfo.CfgId,
                            equipInfo = ItemInfo.Equip,
                            teamIdx = curRoleIndex,
                            characterIndex = curRoleIndex + 1,
                            equipSlotPos = dressidx,
                            equipRefineData = equipRefineData,
                        });
                    }
                }
            });
        }
    }

    public void OnClickRoleAtIndex(int index)
    {
        OnEquipTabBtnClick(index);
    }

    private void OnEquipTabBtnClick(int idx)
    {
        //if (idx > 0 && CFG_XMLConfig.GetOpenSoltLv(idx) > RoleM.Instance.Level)
        //{
        //    UITool.ShowTips("角色等级不足");
        //    return;
        //}

        if (idx > 0 && slotIsLock(idx))
            return;

        curRoleIndex = idx;
        for (int i = 0; i < uIViewHall.equipTabBtnItems.Length; i++)
        {
            uIViewHall.equipTabBtnItems[i].selGo.SetActive(idx == i);
        }

        ShowEquipUI(idx);
        OnEquipSoltsShow(idx);

        if (topUI != null)
            topUI.OnHeadOnShow(idx);
    }

    private bool slotIsLock(int soltIdx)
    {
        bool isLock = CFG_XMLConfig.CheckPartnerSoltIsOpenServer(soltIdx, false);
        (bool isOtherLock, string strtxttip, string strtip) = CheakIsLock(CFG_XMLConfig.GetOpenSoltLv(soltIdx), CFG_XMLConfig.GetOpenSolMainStaget(soltIdx));
        if (!isLock)
            isLock = isOtherLock;

        if (isLock)
        {
            if (CFG_XMLConfig.CheckPartnerSoltIsOpenServer(soltIdx, false))
                UITool.ShowTips($"{CFG_XMLConfig.RemainingDaysPartnerSlotUnlockOpen(soltIdx)}天后解锁");
            else
                UITool.ShowTips(strtip);
        }

        return isLock;
    }

    public (bool, string, string) CheakIsLock(int level, int mainBattleID)
    {
        string strtxttip = "";
        string strtip = "";
        if (level > RoleM.Instance.Level)
        {
            strtxttip = level + "级解锁";
            strtip = level + "级解锁";
            return (true, strtxttip, strtip);
        }


        if (mainBattleID > RoleM.Instance.MainBattleId)
        {
            var battleInfo = CFG_MainBattle.GetById(mainBattleID);
            if (battleInfo == null) return (false, "", "");//没有找到配置信息，就当满足开启条件

            strtxttip = "主线解锁";
            strtip = $"通关{battleInfo.name}解锁";

            return (true, strtxttip, strtip);
        }

        return (false, "", "");
    }

    private void ShowEquipUI(int idx)
    {
        for (int i = 0; i < uIViewHall.equipItems.Length; i++)
        {
            int dressIDx = i + 1;
            var equipInfo = idx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)dressIDx).Item : PartnerM.Instance.GetOnePartnerOneEquip(idx, dressIDx);
            uIViewHall.equipItems[i].iconImg.gameObject.SetActive(equipInfo != null);
            uIViewHall.equipItems[i].silhouetteGo.SetActive(equipInfo == null);
            if (equipInfo != null)
            {
                var itemCfg = CFG_Item.GetById(equipInfo.CfgId);
                GameUtil.Instance.SetItemIcon(equipInfo.CfgId, uIViewHall.equipItems[i].iconImg, uIViewHall.equipItems[i].bgImg, null, (int)equipInfo.Equip.Quality);
                uIViewHall.equipItems[i].equipLvText.text = "Lv" + itemCfg.Levellimit;

                if (equipInfo.Equip != null && (int)equipInfo.Equip.Quality >= 4 && equipInfo.Equip.Attributes != null && equipInfo.Equip.Attributes.Count > 0)
                {
                    uIViewHall.equipItems[i].baseEffectGo.SetActive(true);
                    SpineUtil.Instance.PlayItemSpineLoop(uIViewHall.equipItems[i].baseEffectTrans, (int)equipInfo.Equip.Quality);
                }
                else
                    uIViewHall.equipItems[i].baseEffectGo.SetActive(false);
            }
            else
            {
                GameUtil.Instance.SetItemBg(0, uIViewHall.equipItems[i].bgImg);
                uIViewHall.equipItems[i].baseEffectGo.SetActive(false);
            }
        }
    }

    private void SetDressEff(int solt)
    {
        GameHelper.Instance.SetAsyncObj("UIBoom", GameConst.GetEffPath("UIBoom"), (go) =>
        {
            equipEffectID++;
            if (equipEffectID >= 100000)
                equipEffectID = 0;

            if (equipEffectsDic.ContainsKey(equipEffectID))
                equipEffectID += 100;

            go.transform.SetParent(uIViewHall.equipItems[solt - 1].iconImg.transform);
            go.gameObject.SetActive(true);
            go.transform.localScale = Vector3.one * 30;
            go.transform.localPosition = Vector3.zero;
            equipEffectsDic.Add(equipEffectID, go);

            var index = equipEffectID;
            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (isNullGo) return;
                PoolHelper.Instance.RecycObj(go);
                equipEffectsDic.Remove(index);
            }, 0.5f);
        }, true).Forget();
    }

    /// <summary>
    /// 装备槽位信息
    /// </summary>
    /// <param name="TeamIdx"></param>
    private void OnEquipSoltsShow(int TeamIdx)
    {
        for (int i = 0; i < uIViewHall.equipItems.Length; i++)
        {
            int soltidx = i + 1;
            var equipInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)soltidx) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, soltidx);
            uIViewHall.equipItems[i].starText.text = equipInfo.Stars.ToString();
            uIViewHall.equipItems[i].starGo.gameObject.SetActive(equipInfo.Stars > 0);
            uIViewHall.equipItems[i].lvText.text = equipInfo.Enhance == 0 ? "" : "+" + equipInfo.Enhance;
        }
    }
    #endregion

    public override void OnEnter(params object[] args)
    {
        equipEffectID = 0;
        if (topUI == null)
        {
            LoadTopUI();
            LoadDownUI();
        }
        else
            ShowUI();

        UniGameEvent.AddEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);
        refreshFuncNewFlag();

        NewPlayerGuideWidget.Instance.AddGuideGameObject("Hall_FightButton", GetView().attackBtn.transform);

        NewPlayerGuideWidget.Instance.AddGuideGameObject("Hall_BoxShowUI", GetView().boxShowUI.transform);
        NewPlayerGuideWidget.Instance.AddGuideGameObject("Hall_Btn_Recycle", GetView().boxDelBtn.transform);
        NewPlayerGuideWidget.Instance.AddGuideGameObject("Hall_Btn_Wear", GetView().boxClothingBtn.transform);
    }

    private void ShowUI()
    {
        CheckGuide();
        GameHelper.Instance.SetDeayFrame(() =>
        {
            ChangeState(HallBattleState.Rest);
            OnSendTCPAndHTTP();
            OnHallChatUI();
            SetRoleSpine();
            SetPetSpine();
            OnEquipTabBtnClick(0);
            SetMainTaskInfo();
            RefreshMainPower();
            RefreshHallBox();
            CheckRed();
            CheckAutoSpecialSkill();

            // 礼包、特惠
            uIViewHall.specialOfferBtnTrans.gameObject.SetActive(!GameM.Instance.IsAudit);
            uIViewHall.giftBagBtnTrans.gameObject.SetActive(!GameM.Instance.IsAudit);
            GetView().rootBoxSet.gameObject.SetActive(false);
        });

        // 延迟设置公共事件
        GameHelper.Instance.SetDeayTimeDo(() =>
        {
            if (isNullGo) return;
            GameUtil.Instance.TDPublicEvent();
        }, 1f);
    }

    #region http and check active
    private void OnSendTCPAndHTTP()
    {
        HttpC.GetActivityState(() =>
        {
            CheckSignState();
            CheckFirstPayState();
        });

        HttpC.OnGetShopInfo(YesFun: (res) =>
        {
            CheckBattleSupportGiftState();
            CheckActivityRedPoint();
        });

        HttpC.OnGetMailInfo();
        TcpC.GetBagInfo();
        TcpC.C2S_RefineInfoRequest();
        TcpC.C2S_InstanceInfo();
        TcpC.C2S_TitleList();

        TcpC.OnPVPInfo((int)MatchType.SINGLE_COMBAT);
        TcpC.OnPVPInfo((int)MatchType.FREE_FOR_ALL);
        TcpC.C2S_DailyMapInfo();
        TcpC.C2S_ElementInfoRequest();

        HttpC.GetSignList();

        HttpC.GetPrivilegeDetail(1);
        HttpC.GetPrivilegeDetail(2);
        HttpC.GetTrialFundAwardInfo();
        HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.GrowthFund);
        HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.RefineFund);
        HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.CallFund);
        HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.CopyFund);
        TcpC.C2S_LotteryInfo();
    }

    private void CheckActivityRedPoint()
    {
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.GiftBag_DailyGift.ToString(), GiftBagM.Instance.GetGiftRedNum(13));
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.GiftBag_GiftPay.ToString(), GiftBagM.Instance.GetGiftRedNum(4));
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.GiftBag_AdGift.ToString(), GiftBagM.Instance.GetGiftRedNum(20));
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.GiftBag_Power.ToString(), GiftBagM.Instance.GetGiftRedNum(21));
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.GiftBag_Time.ToString(), GiftBagM.Instance.GetGiftRedNum(23));


        RedPointManager.Instance.RPS.Dispatch(RedPointConst.ActiveCenter_StarWish);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.ActiveCenter_CrazyRefine);
    }

    private void CheckSignState()
    {
        if (!ActivityM.Instance.GetAcitvityState(ActivityType.Sign7)
            && !ActivityM.Instance.GetAcitvityState(ActivityType.SignHappy7)
            && !ActivityM.Instance.GetAcitvityState(ActivityType.StarWish)
            && !ActivityM.Instance.GetAcitvityState(ActivityType.CrazyRefine)
            )
            return;

        uIViewHall.sevenDayBtnTrans.gameObject.SetActive(true && !GameM.Instance.IsAudit);
    }

    private void CheckFirstPayState()
    {
        if (!ActivityM.Instance.GetAcitvityState(ActivityType.FirstPay))
        {
            uIViewHall.firstPayBtnTrans.gameObject.SetActive(false);
            ActivityM.Instance.CloseAcitvity(ActivityType.FirstPay);
            return;
        }

        HttpC.GetFirstPayInfo((req) =>
        {
            if (req.code == 0 && req.day == 3 && req.status == (int)ActivityM.SignState.Geted)
            {
                uIViewHall.firstPayBtnTrans.gameObject.SetActive(false);
                ActivityM.Instance.CloseAcitvity(ActivityType.FirstPay);
            }
            else
                uIViewHall.firstPayBtnTrans.gameObject.SetActive(true && !GameM.Instance.IsAudit);
        });
    }

    private void CheckBattleSupportGiftState()
    {
        var cfg = CFG_Activity.GetActiveListByType(CFG_Activity.Type.BattleSupportGift)[0];
        //long seconds = GameUtil.Instance.GetRegistRoleTimeSecoud((int)cfg.ResetTime);
        uIViewHall.battleSupportGiftTrans.gameObject.SetActive(StoreM.Instance.GetBuyNumByID(33001) == 0 && !GameM.Instance.IsAudit);
        //BattleSupportGiftTime();
    }

    TimeItem battleSupportGiftTimeItem;

    /*
    private void ClearBattleSupportGiftTime()
    {
        if (GameHelper.Instance != null && battleSupportGiftTimeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(battleSupportGiftTimeItem);
            battleSupportGiftTimeItem = null;
        }
    }

    private void BattleSupportGiftTime()
    {
        ClearBattleSupportGiftTime();
        var cfg_ActivityInfo = CFG_Activity.GetActiveListByType(CFG_Activity.Type.battleSupportGift)[0];
        long seconds = GameUtil.Instance.GetRegistRoleTimeSecoud((int)cfg_ActivityInfo.ResetTime);
        JUIExtensions.txt(uIViewHall.battleSupportGiftTrans.Find("txt_time")).text = GameUtil.Instance.GetTimeToDHorHM(seconds);
        battleSupportGiftTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(seconds + 1, 1, () =>
        {
            if (isNullGo) return;

            JUIExtensions.txt(uIViewHall.battleSupportGiftTrans.Find("txt_time")).text = GameUtil.Instance.GetTimeToDHorHM(seconds);
            if (seconds <= 0)
            {
                JUIExtensions.txt(uIViewHall.battleSupportGiftTrans.Find("txt_time")).text = "";
                uIViewHall.battleSupportGiftTrans.gameObject.SetActive(false);
                ClearBattleSupportGiftTime();
                ShowOrHideActiveButton();
                return;
            }
            seconds--;
        }));
    }
    */

    #endregion

    #region chat
    public void LoadChatHistyory()
    {
        TcpC.C2S_ChatHistoryRequest(0);
    }

    private void OnHallChatUI()
    {
        // 最多两个聊天信息
        List<ChatM.ChatInfo> chatInfo = ChatM.Instance.GetLastTwoChatOfAllType();
        for (int i = 0; i < 2; i++)
        {
            uIViewHall.chatTrans.GetChild(i).gameObject.SetActive(false);
        }
        for (int i = 0; i < chatInfo.Count; i++)
        {
            GameObject go = uIViewHall.chatTrans.transform.GetChild(i).gameObject;
            go.SetActive(true);
            var oneChat = chatInfo[i];
            string content = "<color=#ffffff00>           </color>" + ChatM.Instance.GetChatByType(oneChat.ChatType, oneChat.Context, oneChat.Name);
            JUIExtensions.txt(go.transform).text = content;
            bool isEmjoy = ChatM.Instance.CheckIsEmoji(oneChat.Context);
            go.transform.Find("emjoy").gameObject.SetActive(isEmjoy);
            LayoutElement layoutElement = go.GetComponent<LayoutElement>();
            if (isEmjoy)
            {
                content += "\n";
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ChatAtlas, oneChat.Context, JUIExtensions.img(go.transform.Find("emjoy")));
                layoutElement.minHeight = 60;
            }
            else
            {
                layoutElement.minHeight = 30;
            }

            (string atlas, string icon) = ChatM.Instance.GetChatTypeImg(oneChat.ChatType);
            ResouseHelper.Instance.LoadSpriteByAtlas(atlas, icon, JUIExtensions.img(go.transform.Find("img_Tip")));
        }
    }


    [Subscriber(ThreadMode.Main)]
    private void OnRefreshChatUI(GameEvent.Sync_ChatPush _e)
    {
        if (uiGo == null || uiGo.Equals(null))
            return;
        OnHallChatUI();
    }
    #endregion

    #region 大厅任务
    private void InitTask()
    {
        uIViewHall.taskBtn.onClick.RemoveAllListeners();
        uIViewHall.taskBtn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            var mainTasks = TaskM.Instance.GetTasksByClass(TaskClass.MAINLINE);
            if (mainTasks == null || mainTasks.Count == 0)
                return;
            TaskInfo task = mainTasks[0];
            if (task.Status == TaskStatus.COMPLETED)
            {
                TcpC.C2S_TaskComplate(task.CfgId);
                SendTaskGet();
            }
            else
            {
                if ((TargetType)CFG_Task.GetById(task.CfgId).Type == TargetType.MAIN_BATTLE)
                {
                    OnClickAttackButton();
                }
                else
                {
                    CFG_Task.JumpUI(task.CfgId);
                }
            }

        });

        ButtonAddListener(uIViewHall.mainPowerBtn, () => { UITool.Open<UISpecialOffer>(null, ActivityType.PowerGift); });

        static void SendTaskGet()
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = "主线任务按钮",
                position = "",
                event_click = $"主线任务完成领取",
            });
        }
    }

    private void SetMainTaskInfo()
    {
        var alltask = TaskM.Instance.GetTasksByClass(TaskClass.MAINLINE);
        if (alltask == null || alltask.Count == 0)
        {
            //单机服：新号还没有主线任务时不要索引 [0]，否则崩在 ILRuntime 里
            uIViewHall.taskNormalTaskInfoGo?.SetActive(false);
            uIViewHall.taskEndTaskInfoGo?.SetActive(false);
            return;
        }
        alltask.Sort((a, b) => b.CfgId.CompareTo(a.CfgId));

        var taskInfo = CFG_Task.GetById(alltask[0].CfgId);
        if (taskInfo == null)
            return;
        bool isOver = alltask.Count > 0 && alltask[0].Status == TaskStatus.REWARDED && taskInfo.NextAchievement <= 0;
        uIViewHall.taskNormalTaskInfoGo?.SetActive(!isOver);
        uIViewHall.taskEndTaskInfoGo?.SetActive(isOver);

        uIViewHall.taskGuideGo.SetActive(false);
        if (isOver) return;

        //Debug.Log($"TaskID:{alltask[0].CfgId} State:{alltask[0].Status}");
        var task = alltask[0];
        uIViewHall.taskDesText.text = CFG_Task.StrTaskDes(task.CfgId);
        (int amout, int targetnum) = CFG_Task.GetTaskPro(task.CfgId, task.Target.Amount, task.Target.TargetNum);
        uIViewHall.taskProText.text = $"{GameUtil.Instance.FormatNumWithUnit(amout)}/{GameUtil.Instance.FormatNumWithUnit(targetnum)}";
        uIViewHall.taskStateGo.SetActive(task.Status == pb.TaskStatus.COMPLETED);

        uIViewHall.taskGuideGo.SetActive(task.Status != pb.TaskStatus.COMPLETED && task.CfgId <= 1015 && taskInfo.Type != 182);
        //if (task.Status == TaskStatus.COMPLETED)
        //{
        //SpineUtil.Instance.PlayLoop(uIViewHall.tfTaskFinishEffect, "animation");
        //}
        uIViewHall.taskBgImg.enabled = false;// task.Status == pb.TaskStatus.COMPLETED;

        GameUtil.Instance.SetItemIcon(int.Parse(CFG_Task.GetById(task.CfgId).Reward[0].ElementAt(0).Key), uIViewHall.taskImg);
    }

    [Subscriber(ThreadMode.Main)]
    private void Refresh_Task(GameEvent.Refresh_Task _e)
    {
        SetMainTaskInfo();
    }
    #endregion

    #region 大厅体力
    public void RefreshMainPower()
    {
        var curPower = RoleM.Instance.Energy;
        var energyInfo = CFG_XMLConfig.xmlConfig.BaseValue.energy;
        var fullValue = energyInfo.fullValue;
        if (PrivilegeM.Instance.HadPrivilege(1))
        {
            var cfg = CFG_Privilege.GetById(1);
            fullValue += cfg.StaminaUp;
        }
        uIViewHall.mainPowerSlider.value = curPower * 1.0f / fullValue;
        uIViewHall.mainPowerSliderText.text = $"{curPower}/{fullValue}";
        uIViewHall.mainPowerText.text = $"每{energyInfo.timeInterval}分钟恢复一点";
    }
    #endregion

    public void CheckNewFunctionTip()
    {
        int theDayOpenServer = GameUtil.Instance.GetTheDayOpenServer();
        if (GameUtil.Instance.HaveNewFunctionTip(theDayOpenServer) && !GameUtil.Instance.CheckIsTipNewFunction(theDayOpenServer))
        {
            GameUtil.Instance.ShowNewFunctionTip(theDayOpenServer);
            GameUtil.Instance.SetIsTipNewFunction(theDayOpenServer, true);
        }
    }

    #region game event
    [Subscriber(ThreadMode.Main)]
    private void OnSync_UserStatsPush(GameEvent.Sync_UserStatsPush _e)
    {
        ShowEquipUI(curRoleIndex);
    }

    //[Subscriber(ThreadMode.Main)]
    //打开宝箱
    //private void OnHallOpenBox(GameEvent.OnHall_OpenBox _OpenBox)
    //{
    //    OnClickFuncBtnAtIndex(1);
    //}

    [Subscriber(ThreadMode.Main)]
    private void OnRefresh_ActiveState(GameEvent.Refresh_ActiveState _e)
    {
        if (isNullGo)
            return;
        HttpC.OnLevelGiftInfo();
        HttpC.OnSign30Info();
        HttpC.OnSign7Info();
        HttpC.GetFirstPayInfo();
        HttpC.GetPrivilegeList();

        SetTopActiveButton(default);
    }

    [Subscriber(ThreadMode.Main)]
    // 设置大厅活动开关显示和隐藏
    private void SetTopActiveButton(GameEvent.Refresh_HallActivityButton _e)
    {
        ShowOrHideActiveButton();

        if (_e.forceHide)
        {
            switch (_e.index)
            {
                case 1:
                    uIViewHall.sevenDayBtnTrans.gameObject.SetActive(false);
                    break;

                case 2:
                    uIViewHall.firstPayBtnTrans.gameObject.SetActive(false);
                    ActivityM.Instance.CloseAcitvity(ActivityType.FirstPay);
                    break;

                default:
                    break;
            }
        }
    }

    //[Subscriber(ThreadMode.Main)]
    // 设置大厅活动开关显示和隐藏
    private void RefreshTopActiveButton(GameEvent.Refresh_ZeroTime _e)
    {
        HttpC.OnGetShopInfo();
        TcpC.GetBagInfo();

        ShowOrHideActiveButton();
        //SetButtomRightButtonActive();

        //检查是否有新功能开启需要表现，在每天零点
        if (hallBattleState == HallBattleState.Rest)
        {
            //NewPlayerGuideWidget.Instance.TryShowGuide();
            FuncOpenShowWidget.Instance.TryAddShowData();
            FuncOpenShowWidget.Instance.TryShow();
        }
    }

    private void ShowOrHideActiveButton()
    {
        uIViewHall.firstPayBtnTrans.gameObject.SetActive(ActivityM.Instance.GetAcitvityState(ActivityType.FirstPay) && !GameM.Instance.IsAudit);
        uIViewHall.limitetFightRankTrans.gameObject.SetActive(RankM.Instance.GetLimitRankTime((int)CFG_RankingAward.RankClass.Fight) > 0);
        uIViewHall.limitetMainRankTrans.gameObject.SetActive(RankM.Instance.GetLimitRankTime((int)CFG_RankingAward.RankClass.MainPass) > 0);
        CheckBattleSupportGiftState();
        OnLimeFightRankTime();
        OnLimeMainRankTime();
    }

    TimeItem LimeFightRankTime;

    private void ClearLimeFightRankTime()
    {
        if (GameHelper.Instance != null && LimeFightRankTime != null)
        {
            GameHelper.Instance.RemoveInvoke(LimeFightRankTime);
            LimeFightRankTime = null;
        }
    }

    TimeItem LimeMainRankTime;

    private void ClearLimeMainRankTime()
    {
        if (GameHelper.Instance != null && LimeMainRankTime != null)
        {
            GameHelper.Instance.RemoveInvoke(LimeMainRankTime);
            LimeMainRankTime = null;
        }
    }


    private void OnLimeFightRankTime()
    {
        ClearLimeFightRankTime();
        if (RankM.Instance.GetLimitRankTime((int)CFG_RankingAward.RankClass.Fight) > 0)
        {
            LimeFightRankTime = GameHelper.Instance.AddInvoke(new TimeItem(60, -1, () =>
            {
                if (isNullGo) return;
                JUIExtensions.txt(uIViewHall.limitetFightRankTrans.Find("txt_time")).text = GameUtil.Instance.GetTimeToDHorHM(RankM.Instance.GetLimitRankTime((int)CFG_RankingAward.RankClass.Fight));
            }, _isLoop: true));
        }
    }

    private void OnLimeMainRankTime()
    {
        ClearLimeMainRankTime();
        if (RankM.Instance.GetLimitRankTime((int)CFG_RankingAward.RankClass.MainPass) > 0)
        {
            LimeFightRankTime = GameHelper.Instance.AddInvoke(new TimeItem(60, -1, () =>
            {
                if (isNullGo) return;
                JUIExtensions.txt(uIViewHall.limitetMainRankTrans.Find("txt_time")).text = GameUtil.Instance.GetTimeToDHorHM(RankM.Instance.GetLimitRankTime((int)CFG_RankingAward.RankClass.MainPass));
            }, _isLoop: true));
        }
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_SyncUserBasePush(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        CheckRed();
        RefreshMainPower();

        if (GuildM.Instance.IsGuildOpen())
            GuildM.Instance.CheckGuildSkillRed();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_Refresh_LvUp(GameEvent.Refresh_LvUp _e)
    {
        if (isNullGo) return;

        if (downUI == null) return;

        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_RoleUpgrade);
        CheckStartGuidePartnetUp(downUI);
        //Debug.LogError(_e.uid);

        if (RoleM.Instance.Uid == _e.uid)
        {
            HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.GrowthFund);
            //SetButtomRightButtonActive();
            ShowLvUp();
        }

        //检查是否有新功能开启需要表现，在等级变更
        if (hallBattleState == HallBattleState.Rest)
        {
            //NewPlayerGuideWidget.Instance.TryShowGuide();
            FuncOpenShowWidget.Instance.TryAddShowData();
            FuncOpenShowWidget.Instance.TryShow();
        }

        if (hallRole != null)
            hallRole.PreLoadRoleSkills();

        CheckRed();
    }

    // 升级特效
    private void ShowLvUp()
    {
        var lv = RoleM.Instance.Level;

        var trans = uIViewHall.roleTranss[0];
        if (trans.childCount > 0)
        {
            uIViewHall.roleLvUpText.text = "";
            uIViewHall.roleLvUpTrans.SetParent(trans.GetChild(0));
            uIViewHall.roleLvUpTrans.localPosition = new Vector3(6, -94, 0);
            uIViewHall.roleLvUpTrans.gameObject.SetActive(true);

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (!isNullGo)
                {
                    uIViewHall.roleLvUpText.text = RoleM.Instance.Level.ToString();
                    DOTween.To(() => 0, (a) =>
                    {
                        SetLvUpTextColor(a);
                    }, 100, 0.5f).SetEase(Ease.Linear);
                }
            }, 0.5f);
            SpineUtil.Instance.PlayOnce(uIViewHall.roleLvUpTrans, "skill2", () =>
            {
                uIViewHall.roleLvUpTrans.gameObject.SetActive(false);
                uIViewHall.roleLvUpTrans.SetParent(uIViewHall.transform);
                uIViewHall.roleLvUpTrans.localScale = Vector3.one * 0.11f;
                uIViewHall.roleLvUpTrans.localEulerAngles = Vector3.zero;
            });
        }
    }

    private void SetLvUpTextColor(int a)
    {
        var color = uIViewHall.roleLvUpText.color;
        color.a = (100 - a) * 0.01f;
        uIViewHall.roleLvUpText.color = color;
    }

    //[Subscriber(ThreadMode.Main)]
    private void Event_ChangeEquip(GameEvent.ChangeEquip _e)
    {
        if (isNullGo) return;
        if (curRoleIndex == 0)
            ShowEquipUI(curRoleIndex);

        (_, int soltIdx) = GameUtil.Instance.GetEquipIsBest(_e.item);

        if (soltIdx >= 0 && soltIdx != curRoleIndex)
        {
            OnClickRoleAtIndex(soltIdx);

            foreach (var item in equipEffectsDic)
                PoolHelper.Instance.RecycObj(item.Value);

            equipEffectsDic.Clear();
            equipEffectID = 0;
        }

        SetDressEff(_e.soldIdx);
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_BackHall(GameEvent.BackUIHall _e)
    {
        if (_e.isInHall)
        {
            if (isNullGo || downUI == null) return;
            CheckStartGuidePartnetUp(downUI);
        }
    }
    [Subscriber(ThreadMode.Main)]
    private void OnSync_Sync_UserBoxInfo(GameEvent.Sync_UserBoxInfo _e)
    {
        // 宝箱红点
        CheckRed();
        OnRefreshHomeBox(default);
    }

    // 上阵的伙伴有变化时，刷新开箱结果的装备可穿戴状态
    //[Subscriber(ThreadMode.Main)]
    private void OnSync_PartnetChanged(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.code == 0 && (_e.cid == CmdId.C2S_PartnerReplace || _e.cid == CmdId.C2S_PartnerCombat))
            OnRefreshHomeBox(default);
    }

    [Subscriber(ThreadMode.Main)]
    //槽位魂石镶嵌
    private void EquipInsetRefresh(GameEvent.EquipInsetRefresh _e)
    {
        OnEquipSoltsShow(curRoleIndex);
    }

    [Subscriber(ThreadMode.Main)]
    //槽位强化
    private void Replace_EquipEnhance(GameEvent.Replace_EquipEnhance _e)
    {
        OnEquipSoltsShow(curRoleIndex);
    }

    [Subscriber(ThreadMode.Main)]
    //槽位升星
    private void Replace_EquipStar(GameEvent.Replace_EquipStar _e)
    {
        OnEquipSoltsShow(curRoleIndex);
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_AdvanceInfo(GameEvent.Sync_UserAdvanceInfoPush _e)
    {
        CheckGuide();
    }

    private void CheckGuide()
    {
        //if (GuidanceM.Instance.CheckShowGuide() && GuidanceM.Instance.isGuiding)
        //{
        //    GuidanceM.Instance.isGuiding = true;
        //    UITool.Open<UIGuide>();
        //}

        NewPlayerGuideWidget.Instance.TryShowGuide();
        //ShowFuncBox(GuidanceM.Instance.CanShowBox());
    }

    private void GuideJump(GameEvent.Guide_Jump _e)
    {
        //ShowFuncBox(GuidanceM.Instance.CanShowBox());
    }

    // 购买月卡啥的发消息
    private void RefreshPrivilege(GameEvent.Refresh_Privilege _e)
    {
        RefreshMainPower();
    }
    #endregion

    public override void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(RefreshTopActiveButton);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnSync_PartnetChanged);
        UniGameEvent.RemoveEventListener<MainBattleResponse>(Sever_MainBattleResponse);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_Privilege>(RefreshPrivilege);
        UniGameEvent.RemoveEventListener<GameEvent.Guide_Jump>(GuideJump);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_HomeOpen>(OnRefreshHomeBox);
        UniGameEvent.RemoveEventListener<GameEvent.ChangeEquip>(Event_ChangeEquip);
        UniGameEvent.RemoveEventListener<GameEvent.OpenBoxDisassemblyEquip>(DisassemblyEffectFlyAsync);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_SyncUserBasePush);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_LvUp>(OnSync_Refresh_LvUp);
        UniGameEvent.RemoveEventListener<GameEvent.HallBattleBossBeHit>(HallBattleBossBeHitEnt);
        UniGameEvent.RemoveEventListener<GameEvent.HallBattleDropItem>(HallBattleDropItemEnt);
        UniGameEvent.RemoveEventListener<GameEvent.HallBattleRoleAttack>(HallBattleRoleAttackEnt);
        UniGameEvent.RemoveEventListener<GameEvent.HallBallPauseOrPlay>(PauseSpecialSkill);
        UniGameEvent.RemoveEventListener<GameEvent.HallBallShowSpecialSkillName>(HallBallShowSpecialSkillNameEnt);
        HallBattleM.Instance.RemoveEventListener();

        UniGameEvent.RemoveEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);

        ClearBindRed();
        ClearPetTalkCoroutine();
        ClearLimeFightRankTime();
        ClearAutoOpenTime();
        ClearTimeOpenResult();
        ClearTimeCompareEquip();
        ClearAutoClothDelayTime();

        BattleClear();
        ClearBossRewardTimer();

        //hpChangeItemDict?.Clear();
        topUI?.OnDestroy();
        downUI?.OnDestroy();


        uIViewHall = null;
    }


    #region 功能开启角标

    private void OnRefresh_FunctionNewFlag(UIFunctionOpenOperate operate)
    {
        refreshFuncNewFlag();
    }

    private void refreshFuncNewFlag()
    {
        bool blessFlag = FuncOpenShowWidget.Instance.FunctionIsLight(new() { SwitchOpne.OpenBox });
        GetView().goBlessRedDotRoot.SetActive(!blessFlag);
        GetView().goBlessFuncNewFlag.SetActive(blessFlag);
    }
    #endregion
}