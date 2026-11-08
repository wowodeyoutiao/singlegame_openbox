using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Ocsp;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.MyNet;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using Spine;
using Spine.Unity;

using System;
using System.Collections.Generic;

using UnityEditorInternal;

using UnityEngine;
using UnityEngine.UI;

using static CFG_Store;
using static GameEvent;

public class UIBuild_Equip : UIBase<UIModelBuild_Equip, UIViewBuild_Equip>
{

    //资源ID
    private enum ResID
    {
        /// <summary>
        /// 金币
        /// </summary>
        Coin = 1,

        /// <summary>
        /// 强化石
        /// </summary>
        QHS = 21,

        /// <summary>
        /// 升星石
        /// </summary>
        SXS = 22,

        /// <summary>
        /// 洗练石
        /// </summary>
        XLS = 57,
    }

    public enum EquipSelectClickFrom
    {
        TagButtonClick,
        EquipClick,
        Refresh,
    }

    public class BuyInfo
    {
        public int buyId;
        public int num;
    }

    private BuyInfo _BuyInfo = new BuyInfo();


    private int AllEquipCount = 12;


    private TimeItem timeItem;

    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }

    public override void OnInit()
    {

        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnHelpBtnClick();
        });

        //装备选择
        for (int i = 1; i <= AllEquipCount; i++)
        {
            GetView().root_allEquip.Find("item_" + i + "/txt_Name").GetComponent<Text>().text = GameM.Instance.EquipName[i];
            EquipSlot eSlot = (EquipSlot)i;
            GetView().root_allEquip.Find("item_" + i).GetComponent<Button>().onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                if (StateM.Instance.isHasRefineResult)
                {
                    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                    {
                        ShowStr = "当前有洗练结果未保留，是否确定保留？",
                        YesFun = () =>
                         {
                             OnSelectRefineResult(true);
                             OnSelectEquipClick(eSlot, EquipSelectClickFrom.EquipClick);
                         },
                        NoFun = () =>
                         {
                             OnSelectRefineResult(false);
                             OnSelectEquipClick(eSlot, EquipSelectClickFrom.EquipClick);
                         },
                    });
                    return;
                }
                OnSelectEquipClick(eSlot, EquipSelectClickFrom.EquipClick);
            });
        }

        for (int i = 0; i < GetView().usedItems.Length; i++)
        {
            int idx = i;
            GetView().usedItems[i].btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                //UITool.Close<UIBuild>();
                UITool.Open<UIShop>(null, (idx == GetView().usedItems.Length - 1 ? CFG_Store.StoreType.RMB : CFG_Store.StoreType.Coupon));
            });
        }

        GetView().btn_Reinforce.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //强化大师
            UITool.Open<UIEquipMaster>(null, CFG_EquipMaster.CategoryClass.Reinforce, TeamIdx);
        });

        GetView().btn_AddReinforceOne.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //强化一次
            OnAddRefinforceClick(false);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "打造",
                component = "装备",
                position = "强化",
                event_click = "强化1次",
            });
        });

        GetView().btn_AddReinforceTen.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //强化到顶
            EquipEnhanceTimes++;
            OnAddRefinforceClick(true);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "打造",
                component = "装备",
                position = "强化",
                event_click = "强化10次",
            });
        });

        GetView().btn_AddReinforceEnd.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //强化停止
            OnStopAutoRefinforce();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "打造",
                component = "装备",
                position = "强化",
                event_click = "停止强化",
            });
        });

        GetView().btn_Star.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //升星大师
            UITool.Open<UIEquipMaster>(null, CFG_EquipMaster.CategoryClass.UpStar, TeamIdx);
        });

        GetView().btn_AddStarOne.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //升星一次
            OnAddStarClick(false);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "打造",
                component = "装备",
                position = "升星",
                event_click = "升星1次",
            });
        });

        GetView().btn_AddStarTop.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //自动升星
            OnAddStarClick(true);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "打造",
                component = "装备",
                position = "升星",
                event_click = "自动升星",
            });
        });

        GetView().btn_EndStar.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //升星停止
            OnStopAutoUpStar();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "打造",
                component = "装备",
                position = "升星",
                event_click = "停止自动升星",
            });
        });

        GetView().dropDown_Star.onValueChanged.AddListener(SetUpStar);

        GetView().btn_RefineAttributePreview.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //洗练属性预览
            UITool.Open<UIRefinePreview>(null, _SelectEquipSlot);
        });

        //洗练
        GetView().btn_Refine.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnRefine();
        });

        //随机洗练
        GetView().tog_RefineRandom.onValueChanged.AddListener((isOn) =>
        {
            if (isOn)
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (refineType == RefineType.REFINE_TYPE_RANDOM) return;

                OnRefineRandomOrDirect(RefineType.REFINE_TYPE_RANDOM);
            }
        });

        //定向洗练
        GetView().tog_RefineFixed.onValueChanged.AddListener((isOn) =>
        {
            if (isOn)
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (refineType == RefineType.REFINE_TYPE_DIRECTIONAL) return;

                var attrData = RoleM.Instance.GetEquipSlotOneAttributeData((CharacterSlot)(TeamIdx + 1), _SelectEquipSlot, refinedSlot);
                if (attrData == null)
                {
                    UITool.ShowTips("空槽位无法进行定向洗炼！");
                    GetView().tog_RefineRandom.isOn = true;
                    return;
                }
                OnRefineRandomOrDirect(RefineType.REFINE_TYPE_DIRECTIONAL);
            }
        });

        for (int i = 0; i < GetView().tog_RefineNows.Length; i++)
        {
            int idx = i + 1;
            GetView().tog_RefineNows[i].onValueChanged.RemoveAllListeners();
            GetView().tog_RefineNows[i].onValueChanged.AddListener((isOn) =>
            {
                if (isOn)
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                    if ((RefinedSlot)idx == refinedSlot)
                        return;
                    else
                        GetView().tog_RefineRandom.isOn = true;

                    refinedSlot = (RefinedSlot)idx;
                    GetView().txt_RefineResult.transform.parent.localPosition = new Vector3(188, refineResultPosX[(int)refinedSlot - 1], 0);

                    var attrData = RoleM.Instance.GetEquipSlotOneAttributeData((CharacterSlot)(TeamIdx + 1), _SelectEquipSlot, refinedSlot);
                    if (attrData == null)
                    {
                        GetView().tog_RefineRandom.isOn = true;
                    }
                }
            });
        }

        GetView().btn_RefineResuleDontSave.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //洗练放弃
            OnSelectRefineResult(false, true);
        });

        GetView().btn_RefineResuleSave.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //洗练选择
            isSaveResult = true;
            OnSelectRefineResult(true);
        });

        //魂石 镶嵌

        GetView().btn_AutoInset.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //一键镶嵌
            TcpC.SetAllGemstoneChange(_SelectEquipSlot, TeamIdx);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "打造",
                component = "装备",
                position = "镶嵌宝石",
                event_click = "一键镶嵌",
            });
        });

        GetView().btn_InsetNatrue.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //自然大师
            UITool.Open<UIEqipInsetRule>(null, UIEqipInsetRule.UISelect.NatureManster, _SelectEquipSlot, TeamIdx);
        });

        GetView().btn_InsetInfo.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //魂石大师
            UITool.Open<UIEqipInsetRule>(null, UIEqipInsetRule.UISelect.InsetMonster, _SelectEquipSlot, TeamIdx);
        });

        for (int i = 0; i < 4; i++)
        {
            int idex = i;
            JUIExtensions.btn(GetView().root_Inset.GetChild(i)).onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnInsetOnClick(idex, true);
            });

            JUIExtensions.btn(GetView().root_Inset.GetChild(i).Find("icon/del")).onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnInsetOnClick(idex, false);
            });
        }

        // 升星符
        GetView().addTalismanBtn.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIStarTalismanUse>();
        });

        GetView().addTalismanBtn1.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIStarTalismanUse>();
        });
    }

    //装备基础信息
    private void OnEquipUIBase()
    {
        for (int i = 1; i <= AllEquipCount; i++)
        {
            EquipSlot eSlot = (EquipSlot)i;
            var go = GetView().root_allEquip.Find("item_" + i);

            var itemInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(eSlot).Item : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx, i);
            bool isNullEquip = itemInfo == null;
            go.Find("icon").gameObject.SetActive(!isNullEquip);
            go.Find("txt_Name").gameObject.SetActive(isNullEquip);

            if (!isNullEquip)
            {
                GameUtil.Instance.SetItemIcon(itemInfo.CfgId, JUIExtensions.img(go.Find("icon")), JUIExtensions.img(go), null, (int)itemInfo.Equip.Quality);
            }
            else
            {
                GameUtil.Instance.SetItemBg(0, JUIExtensions.img(go));
            }
        }

        OnEquipLvUIRefresh();
    }


    //等级 级刷新
    private void OnEquipLvUIRefresh()
    {
        for (int i = 1; i <= AllEquipCount; i++)
        {
            var go = GetView().root_allEquip.Find("item_" + i);

            var EquipInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)i) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, i);
            JUIExtensions.txt(go.Find("txt_Reinforce")).text = "+" + EquipInfo.Enhance;
            JUIExtensions.txt(go.Find("txt_Star")).text = EquipInfo.Stars + "";

            for (int j = 0; j < EquipInfo.Gemstone.Length; j++)
            {
                go.Find("inset").GetChild(j).gameObject.SetActive((int)EquipInfo.Gemstone[j] != 0);
                if ((int)EquipInfo.Gemstone[j] != 0)
                {
                    GameUtil.Instance.SetItemIcon((int)EquipInfo.Gemstone[j], JUIExtensions.img(go.Find("inset").GetChild(j)));
                }
            }
        }
    }

    // 装备选择按钮红点
    private void OnEquipSoltBtnRedCheak()
    {
        GameHelper.Instance.SetDeayFrame(() => {
            for (int i = 1; i <= AllEquipCount; i++)
            {
                GetView().root_allEquip.Find("item_" + i).Find("dot").gameObject.SetActive(RedNum(i) > 0);
            }

            int RedNum(int Idx)
            {
                if (_SelectEquipClass == UIBuild.SelectTwo.Equip_Reinforce)
                    return RedDotM.Instance.GetOneEquipSlotReinforceRedNum((EquipSlot)Idx, TeamIdx);
                else if (_SelectEquipClass == UIBuild.SelectTwo.Equip_Inset)
                    return RedDotM.Instance.GetOneEquipInsetRedNum(Idx, TeamIdx);
                else
                    return 0;
            }
        });
       
    }

    private EquipSlot _SelectEquipSlot;
    //装备选择
    private void OnSelectEquipClick(EquipSlot eSlot, EquipSelectClickFrom clickFrom)
    {
        _SelectEquipSlot = eSlot;

        for (int i = 1; i <= AllEquipCount; i++)
        {
            GetView().root_allEquip.Find("item_" + i).Find("select").gameObject.SetActive(i == (int)_SelectEquipSlot);
        }

        switch (_SelectEquipClass)
        {
            case UIBuild.SelectTwo.Equip_Reinforce:
                OnReinforceShow();
                break;
            case UIBuild.SelectTwo.Equip_Star:
                OnUpStarShow();

                if (RedPointManager.Instance.GetRedNum(RedPointConst.Equip_Star.ToString()) != 0)
                {
                    RedPointManager.Instance.SetRedClose(RedPointConst.Equip_Star.ToString(), () => {
                        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Star);
                    });
                }

                break;
            case UIBuild.SelectTwo.Equip_Refine:
                GetView().txt_RefineEquipSlotTitle.text = GameM.Instance.EquipName[(int)_SelectEquipSlot] + "槽位";
                switch (clickFrom)
                {
                    case EquipSelectClickFrom.TagButtonClick:
                        unlockAttributeSlot = 0;
                        OnRefineShow();
                        SelectAttributeToggle(0);
                        OnRefineRandomOrDirect(refineType);
                        GetView().tog_RefineRandom.isOn = true;
                        break;
                    case EquipSelectClickFrom.EquipClick:
                        OnRefineShow();
                        SelectAttributeToggle(0);
                        GetView().tog_RefineRandom.isOn = true;
                        break;
                    case EquipSelectClickFrom.Refresh:
                        OnRefineShow();
                        break;
                    default:
                        break;
                }
                break;
            case UIBuild.SelectTwo.Equip_Inset:
                OnInSetUI();
                break;
            default:
                break;
        }
    }

    #region 装备强化

    /// <summary>
    /// 强化
    /// </summary>
    private void OnReinforceShow()
    {
        bool isNullEquip = TeamIdx == 0 ? RoleM.Instance.CheckIsNullEquip(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx, (int)_SelectEquipSlot) == null;
        var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot);

        GetView().img_ItemReinforce.enabled = !isNullEquip;
        if (isNullEquip)
        {
            GameUtil.Instance.SetItemBg(0, GetView().img_ItemBgReinforce);
            GetView().txt_ReinforceItemName.text = GameM.Instance.EquipName[(int)_SelectEquipSlot] + "·槽位";
            JUIExtensions.txt(GetView().img_ItemBgReinforce.transform, "txt_soltName").text = GameM.Instance.EquipName[(int)_SelectEquipSlot];
        }
        else
        {
            GameUtil.Instance.SetItemIcon(_SoltInfo.Item.CfgId, GetView().img_ItemReinforce, GetView().img_ItemBgReinforce, null, (int)_SoltInfo.Item.Equip.Quality);
            GetView().txt_ReinforceItemName.text = GameUtil.Instance.GetItemName(_SoltInfo.Item.CfgId);
            JUIExtensions.txt(GetView().img_ItemBgReinforce.transform, "txt_soltName").text = "";
        }

        GetView().txt_ReinforceLv.text = "+" + _SoltInfo.Enhance;
        GetView().txt_StarReinforce.text = _SoltInfo.Stars + "";

        GetView().txt_ReinforceNowLv.text = _SoltInfo.Enhance.ToString();
        int Quali = isNullEquip ? 1 : (int)CFG_Item.GetById(_SoltInfo.Item.CfgId).Quality;

        int job = TeamIdx == 0 ? RoleM.Instance.Job : CFG_Partner.GetById(_PartnerInfo.CfgId).job % 10;
        //int job = RoleM.Instance.Job;
        var _cfgInfo = CFG_EquipStrengthen.GetByJobAndLv(job, _SoltInfo.Enhance);
        var Solt = CFG_EquipStrengthen.GetAttrByJobAndLvSolt(job, _SoltInfo.Enhance, (int)_SelectEquipSlot);
        GetView().txt_ReinforceNowDes.text = _SoltInfo.Enhance == 0 ? "无" : CFG_XMLConfig.GetAttributeStrByID(Solt.attribute, (int)(CFG_XMLConfig.xmlConfig.EquipEnhanceDic[Quali] * Solt.value));


        GameUtil.Instance.SetItemIcon(_cfgInfo.Consume[0].id, GetView().img_ReinforceGold);
        GetView().txt_ReinforceGoldNum.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(_cfgInfo.Consume[0].num), _cfgInfo.Consume[0].id, _cfgInfo.Consume[0].num);

        GameUtil.Instance.SetItemIcon(_cfgInfo.Consume[1].id, GetView().img_ReinforceQHS);
        GetView().txt_ReinforceQHSNum.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(_cfgInfo.Consume[1].num), _cfgInfo.Consume[1].id, _cfgInfo.Consume[1].num);

        int nextLv = CFG_EquipStrengthen.JobStLvDic[job].Count - 1 == _SoltInfo.Enhance ? _SoltInfo.Enhance : _SoltInfo.Enhance + 1;
        GetView().txt_ReinforceNextLv.text = nextLv.ToString();
        var nextSolt = CFG_EquipStrengthen.GetAttrByJobAndLvSolt(job, nextLv, (int)_SelectEquipSlot);


        GetView().txt_ReinforceNextDes.text = CFG_XMLConfig.GetAttributeStrByID(nextSolt.attribute, (int)(CFG_XMLConfig.xmlConfig.EquipEnhanceDic[Quali] * nextSolt.value));

        int canMaxlv = TeamIdx == 0 ? RoleM.Instance.Level : _PartnerInfo.Level;
        canMaxlv = canMaxlv * CFG_XMLConfig.xmlConfig.BaseValue.equipStrenghthen;

        bool isCanEnhance = _SoltInfo.Enhance < canMaxlv && nextLv != _SoltInfo.Enhance;

        bool isRedDot = true;

        for (int i = 0; i < _cfgInfo.Consume.Length; i++)
        {
            if (!GameUtil.Instance.ComCoinToCheck(_cfgInfo.Consume[i].id, _cfgInfo.Consume[i].num, null, false))
            {
                isRedDot = false;
                break;
            }
        }
        GetView().btn_AddReinforceOne.transform.Find("dot").gameObject.SetActive(isRedDot);
        GetView().btn_AddReinforceTen.transform.Find("dot").gameObject.SetActive(isRedDot);
        if (!isAutoReinforce)
        {
            GetView().btn_AddReinforceOne.gameObject.SetActive(isCanEnhance);
            GetView().btn_AddReinforceTen.gameObject.SetActive(isCanEnhance);
            GetView().Obj_AddReinforceEnd.SetActive(false);
        }
        GetView().txt_ReinforceTip.text = nextLv == _SoltInfo.Enhance ? "已满级" : (_SoltInfo.Enhance < canMaxlv ? "" :
            "需" + (TeamIdx == 0 ? "主角" : "当前伙伴") + "提升至"
        + nextLv / CFG_XMLConfig.xmlConfig.BaseValue.equipStrenghthen + "级");

    }

    //自动强化（10次）
    private bool isAutoReinforce = false;
    private void OnAddRefinforceClick(bool isAuto)
    {
        isAutoReinforce = isAuto;


        var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot);
        int job = TeamIdx == 0 ? RoleM.Instance.Job : CFG_Partner.GetById(_PartnerInfo.CfgId).job % 10;

        var _cfgInfo = CFG_EquipStrengthen.GetByJobAndLv(job, _SoltInfo.Enhance);
        //检测货币
        for (int i = 0; i < _cfgInfo.Consume.Length; i++)
        {
            if (!GameUtil.Instance.ComCoinToCheck(_cfgInfo.Consume[i].id, _cfgInfo.Consume[i].num, () => { UITool.Close<UIBuild>(); }))
            {
                if (isAutoReinforce)
                {
                    GetView().btn_AddReinforceOne.gameObject.SetActive(true);
                    GetView().btn_AddReinforceTen.gameObject.SetActive(true);
                    GetView().Obj_AddReinforceEnd.SetActive(false);
                    isAutoReinforce = false;
                }
                return;
            }
        }

        if (isAutoReinforce)
        {
            GetView().btn_AddReinforceOne.gameObject.SetActive(false);
            GetView().btn_AddReinforceTen.gameObject.SetActive(false);
            GetView().Obj_AddReinforceEnd.SetActive(true);
        }

        TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Enhance, MType.Request, TCPClientMgr.ProtoSerialize(new EnhanceRequest { Slot = _SelectEquipSlot, PartnerSlotId = TeamIdx }));

    }

    //终止强化
    private void OnStopAutoRefinforce()
    {
        EquipEnhanceTimes = 0;
        if (TimeItemRefinforece != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeItemRefinforece);
            TimeItemRefinforece = null;
        }
        isAutoReinforce = false;
        OnReinforceShow();
    }

    private int GetNextSolt(int Solt)
    {
        int nextSolt = Solt + 1;
        if (nextSolt > AllEquipCount)
            nextSolt = 1;

        int Enhance = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)nextSolt).Enhance : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, nextSolt).Enhance;
        int canMaxlv = TeamIdx == 0 ? RoleM.Instance.Level : _PartnerInfo.Level;
        if (Enhance >= CFG_XMLConfig.xmlConfig.BaseValue.equipStrenghthen * canMaxlv)
        {
            nextSolt = GetNextSolt(nextSolt);
        }
        return nextSolt;
    }

    private TimeItem TimeItemRefinforece = null;
    private int EquipEnhanceTimes;

    [Subscriber(ThreadMode.Main)]
    private void OnRefinforeceCallBack(GameEvent.Replace_EquipEnhance _e)
    {
        OnEquipLvUIRefresh();

        bool isCheakRed = true;
        if (isAutoReinforce)
        {
            //自动强化判定
            bool isAllTarget = true;
            int lv = TeamIdx == 0 ? RoleM.Instance.Level : _PartnerInfo.Level;

            //for (int i = 1; i <= AllEquipCount; i++)
            //{
            //    int Enhance = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)i).Enhance : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, i).Enhance;
            //    if (Enhance < lv * CFG_XMLConfig.xmlConfig.BaseValue.equipStrenghthen)
            //    {
            //        isAllTarget = false;
            //        break;
            //    }
            //}
            //展示当前强化面板
            // OnReinforceShow();

            int Enhance = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot).Enhance : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot).Enhance;
            if (Enhance < lv * CFG_XMLConfig.xmlConfig.BaseValue.equipStrenghthen)
                isAllTarget = false;

            if (!isAllTarget && EquipEnhanceTimes < 10)
            {
                isCheakRed = false;
                //_SelectEquipSlot = (EquipSlot)GetNextSolt((int)_SelectEquipSlot);
                //展示下一个强化界面
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_StrengthenSuc);
                OnSelectEquipClick(_SelectEquipSlot, EquipSelectClickFrom.Refresh);

                TimeItemRefinforece = GameHelper.Instance.AddInvoke(new TimeItem(0.2f, delegate
                {
                    OnClearTime();
                    //点击下一个强化按钮
                    EquipEnhanceTimes++;
                    OnAddRefinforceClick(isAutoReinforce);
                }));
            }
            else
                OnStopAutoRefinforce();
        }
        else
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_StrengthenSuc);
            OnSelectEquipClick(_SelectEquipSlot, EquipSelectClickFrom.Refresh);
        }

        if (isCheakRed)
            OnEquipSoltBtnRedCheak();

    }
    #endregion

    private bool IsAutoUpStar = false;
    //升星
    private void OnUpStarShow()
    {
        bool isNullEquip = TeamIdx == 0 ? RoleM.Instance.CheckIsNullEquip(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx, (int)_SelectEquipSlot) == null;
        var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot);

        GetView().img_Star_Item.enabled = !isNullEquip;
        if (isNullEquip)
        {
            GameUtil.Instance.SetItemBg(0, GetView().img_Star_ItemBg);
            GetView().txt_Star_ItemName.text = GameM.Instance.EquipName[(int)_SelectEquipSlot] + "·槽位";
            JUIExtensions.txt(GetView().img_Star_ItemBg.transform, "txt_soltName").text = GameM.Instance.EquipName[(int)_SelectEquipSlot];
        }
        else
        {
            GameUtil.Instance.SetItemIcon(_SoltInfo.Item.CfgId, GetView().img_Star_Item, GetView().img_Star_ItemBg, null, (int)_SoltInfo.Item.Equip.Quality);
            GetView().txt_Star_ItemName.text = GameUtil.Instance.GetItemName(_SoltInfo.Item.CfgId);
            JUIExtensions.txt(GetView().img_Star_ItemBg.transform, "txt_soltName").text = "";
        }

        GetView().txt_Star_ReinforceLv.text = "+" + _SoltInfo.Enhance;
        GetView().txt_Star_StarLv.text = _SoltInfo.Stars + "";

        GetView().txt_Star_NowLv.text = _SoltInfo.Stars + "";

        int job = TeamIdx == 0 ? RoleM.Instance.Job : CFG_Partner.GetById(_PartnerInfo.CfgId).job % 10;
        //int job = RoleM.Instance.Job;
        var _cfgInfo = CFG_EquipStar.GetByJobAndLv(job, _SoltInfo.Stars);

        string SuccessRateStr = "成功率" + (_cfgInfo.SuccessRate / 100) + "%";

        if (RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.EQUIP_STAR_PERCENTAGE) > 0) // 使用升星符
        {
            RefreshTalisman(true, _cfgInfo.SuccessRate);
            SuccessRateStr += "(升星符+" + (RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.EQUIP_STAR_PERCENTAGE) / 10000f * (_cfgInfo.SuccessRate / 100)) + "%)";
        }else
            RefreshTalisman(false, _cfgInfo.SuccessRate);

        GetView().txt_Star_Des.text = SuccessRateStr;

        var Solt = CFG_EquipStar.GetAttrByJobAndLvSolt(job, _SoltInfo.Stars, (int)_SelectEquipSlot);
        GetView().txt_Star_NowDes.text = _SoltInfo.Stars == 0 ? "无" : CFG_XMLConfig.GetAttributeStrByID(Solt.attribute, Solt.value);

        GameUtil.Instance.SetItemIcon(_cfgInfo.MaterialID[0].id, GetView().img_Star_GoldNum);
        GetView().txt_Star_GoldNum.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(_cfgInfo.MaterialID[0].num), _cfgInfo.MaterialID[0].id, _cfgInfo.MaterialID[0].num);

        GameUtil.Instance.SetItemIcon(_cfgInfo.MaterialID[1].id, GetView().img_Star_SXSNum);
        GetView().txt_Star_SXSNum.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(_cfgInfo.MaterialID[1].num), _cfgInfo.MaterialID[1].id, _cfgInfo.MaterialID[1].num);

        int nextLv = _SoltInfo.Stars == CFG_EquipStar.GetMaxStar(job) ? _SoltInfo.Stars : _SoltInfo.Stars + 1;
        GetView().txt_Star_NextLv.text = nextLv + "";
        Solt = CFG_EquipStar.GetAttrByJobAndLvSolt(job, nextLv, (int)_SelectEquipSlot);
        GetView().txt_Star_NextDes.text = CFG_XMLConfig.GetAttributeStrByID(Solt.attribute, Solt.value);

        bool isCanUpStar = _SoltInfo.Stars < CFG_EquipStar.GetMaxStar(job);
        GetView().obj_StarTip.SetActive(!isCanUpStar);

        if (!IsAutoUpStar)
        {
            GetView().btn_AddStarOne.gameObject.SetActive(isCanUpStar);
            GetView().btn_AddStarTop.gameObject.SetActive(isCanUpStar);

            if (isCanUpStar)
            {
                List<string> _dropDown = new List<string>();
                for (int i = _SoltInfo.Stars + 1; i <= CFG_EquipStar.GetMaxStar(job); i++)
                {
                    _dropDown.Add(i + "星");
                }
                GetView().dropDown_Star.options.Clear();
                GetView().dropDown_Star.AddOptions(_dropDown);
                GetView().dropDown_Star.value = 0;
                TargetStar = _SoltInfo.Stars + 1;
            }

            GetView().dropDown_Star.gameObject.SetActive(isCanUpStar);
            GetView().obj_EndStar.SetActive(false);
        }
    }

    private void RefreshTalisman(bool inUse, int successRate)
    {
        GetView().noTalismanGo.SetActive(!inUse);
        GetView().talismanTrans.gameObject.SetActive(inUse);

        var num = RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.EQUIP_STAR_PERCENTAGE);
        var itemId = CFG_ItemUse.GetStarTalismanIdByAttr(num);
        if(itemId == 0)
        {
            GetView().noTalismanGo.SetActive(true);
            GetView().talismanTrans.gameObject.SetActive(false);
            return;
        }
        var itemCfg = CFG_Item.GetById(itemId);
        GameUtil.Instance.SetItemIcon(itemId, JUIExtensions.img(GetView().talismanTrans.Find("IconImg")), JUIExtensions.img(GetView().talismanTrans.Find("BgImg")));

        JUIExtensions.txt(GetView().talismanTrans.Find("AdvanceText")).text = "(成功率+" + (RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.EQUIP_STAR_PERCENTAGE) / 10000f * (successRate / 100)) + "%)";

        ClearTimer();
        if (inUse)
            AddTimer(successRate);
    }

    private void AddTimer(int successRate)
    {
        var endTime = (int)(RoleM.Instance.GetBuffRemainingTime(14) * 0.001f) - (int)Time.time + 10;
        RefrehRemainingTimeText(endTime - 10);
        timeItem = GameHelper.Instance.AddInvoke(new TimeItem(endTime, 1, () =>
        {
            endTime--;
            RefrehRemainingTimeText(endTime - 10);
            if (endTime <= 10)
                Reset();
        }));

        void RefrehRemainingTimeText(int endTime)
        {
            if(endTime <= 0)
                Reset();
            else
            {
                JUIExtensions.txt(GetView().talismanTrans.Find("TimeText")).text = $"(剩余{GameUtil.Instance.GetTimeSToDHOrHMS(endTime)})";
            }  
        }

        void Reset()
        {
            ClearTimer();
            GetView().txt_Star_Des.text = "成功率" + (successRate / 100) + "%";
            GetView().noTalismanGo.SetActive(true);
            GetView().talismanTrans.gameObject.SetActive(false);
        }
    }

    private void ClearTimer()
    {
        if (timeItem != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(timeItem);
            timeItem = null;
        }
    }

    private int TargetStar;
    //选择目标升星
    private void SetUpStar(int idx)
    {
        var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot);
        TargetStar = idx + 1 + _SoltInfo.Stars;
    }

    //升星选择点击
    private void OnAddStarClick(bool isAuto)
    {
        IsAutoUpStar = isAuto;

        var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot);
        int job = TeamIdx == 0 ? RoleM.Instance.Job : CFG_Partner.GetById(_PartnerInfo.CfgId).job % 10;
        var _cfgInfo = CFG_EquipStar.GetByJobAndLv(job, _SoltInfo.Stars);
        //检测货币
        for (int i = 0; i < _cfgInfo.MaterialID.Length; i++)
        {
            if (!GameUtil.Instance.ComCoinToCheck(_cfgInfo.MaterialID[i].id, _cfgInfo.MaterialID[i].num, () => { UITool.Close<UIBuild>(); }))
            {
                if (IsAutoUpStar)
                {
                    GetView().obj_EndStar.SetActive(false);
                    GetView().btn_AddStarOne.gameObject.SetActive(true);
                    GetView().btn_AddStarTop.gameObject.SetActive(true);
                    IsAutoUpStar = false;
                }
                return;
            }
        }

        if (IsAutoUpStar)
        {
            GetView().obj_EndStar.SetActive(true);
            GetView().btn_AddStarOne.gameObject.SetActive(false);
            GetView().btn_AddStarTop.gameObject.SetActive(false);
        }

        TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Stars, MType.Request, TCPClientMgr.ProtoSerialize(new StarsRequest { Slot = _SelectEquipSlot, PartnerSlotId = TeamIdx }));

    }

    private void OnClearTime()
    {
        if (TimeItemRefinforece != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeItemRefinforece);
            TimeItemRefinforece = null;
        }
    }

    private void OnStopAutoUpStar()
    {
        OnClearTime();
        IsAutoUpStar = false;
        OnUpStarShow();
    }

    [Subscriber(ThreadMode.Main)]
    private void OnStarCallBack(GameEvent.Replace_EquipStar _e)
    {
        OnEquipLvUIRefresh();
        bool isCheakRed = true;
        if (IsAutoUpStar)
        {
            //自动升星判定
            var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot);
            bool isCanUpStar = _SoltInfo.Stars < TargetStar;
            //展示当前升星面板
            OnUpStarShow();

            if (isCanUpStar)
            {
                isCheakRed = false;
                TimeItemRefinforece = GameHelper.Instance.AddInvoke(new TimeItem(0.25f, delegate
                {
                    OnClearTime();
                    //点击下一个强化按钮
                    OnAddStarClick(IsAutoUpStar);
                }));
            }
            else
            {
                OnStopAutoUpStar();
            }
        }
        else
            OnUpStarShow();

        if (isCheakRed)
            OnEquipSoltBtnRedCheak();

    }

    #region 装备洗练

    private RefinedSlot refinedSlot;
    Color defaultOutlineColor => new Color(0.416f, 0.353f, 0.271f);

    /// <summary>
    /// 洗练
    /// </summary>
    /// <param name="isReselectAttribute">重置选中项</param>
    private void OnRefineShow(bool isInit = true)
    {
        if(isInit)
        {
            StateM.Instance.isHasRefineResult = false;

            GetView().txt_RefineResult.text = "暂无词条";
            JUIExtensions.outline(GetView().txt_RefineResult.transform).effectColor = defaultOutlineColor;
            GetView().obj_RefineResultMax.SetActive(false);
            Debug.Log("<color=yellow>洗练选择按钮隐藏 1111 ！！！</color>");
            GetView().obj_RefineResuleButtons.SetActive(false);
            GetView().btn_Refine.gameObject.SetActive(true);
            GetView().obj_ComsumeRoot.SetActive(true);


            SetEnabledRefineToggle(true);
            GetView().tog_RefineRandom.interactable = true;
            GetView().tog_RefineFixed.interactable = true;
        }


        SetRefineItemPreview(_SelectEquipSlot, GetView().root_RefineItemPreview);

        var slotInfo = RoleM.Instance.GetOneEquipSlotData((CharacterSlot)(TeamIdx + 1), _SelectEquipSlot);
        for (int i = 1; i <= 4; i++)
        {
            int idx = i;
            bool isUnlock = GameUtil.Instance.EquipAttributeSlotIsUnlock((CharacterSlot)(TeamIdx + 1), _SelectEquipSlot, (RefinedSlot)i);
            bool isHasData = TryGetOneSloeAttributeData(slotInfo, (RefinedSlot)i, out Stats stats);

            GetView().tog_RefineNows[i - 1].gameObject.SetActive(isUnlock);
            Transform item = GetView().root_RefineNowRoot.Find("Item_" + i);
            item.Find("lock").gameObject.SetActive(!isUnlock);

            (int type, int limit) = CFG_XMLConfig.GetOpenEquipSlotAttribute(i);

            item.Find("img_icon").gameObject.SetActive(type == 3 && !isUnlock);
            Button btnUnlock = JUIExtensions.btn(item.Find("btn_unlock"));
            btnUnlock.gameObject.SetActive(type == 3 && !isUnlock);
            btnUnlock.onClick.RemoveAllListeners();
            btnUnlock.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                if (StateM.Instance.isHasRefineResult)
                {
                    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                    {
                        ShowStr = "当前有洗练结果未保留，是否确定保留？",
                        YesFun = () =>
                        {
                            OnSelectRefineResult(true);
                        },
                        NoFun = () =>
                        {
                            OnSelectRefineResult(false);
                        }
                    });
                    return;
                }

                int storeId = idx == 3 ? 34003 : 34004;
                var refineSlot = new RefineSlot
                {
                    character_slot = (CharacterSlot)(TeamIdx + 1),
                    equip_slot = _SelectEquipSlot,
                    refined_slot = (RefinedSlot)idx
                };

                if (!GameUtil.Instance.ComCoinToCheck(storeId, null))
                {
                    return;
                }

                Debug.Log("ComCoinToCheck 洗炼：");

                unlockAttributeSlot = idx;
                GameUtil.Instance.GameToBuy(storeId, null, refineSlot: refineSlot);
            });

            bool isMax = isHasData && CFG_RandomDistribV2.IsMaxValue((int)stats.attr, stats.Num);
            item.Find("txt/max").gameObject.SetActive(isUnlock && isMax);

            string attrInfo;
            Color outline = defaultOutlineColor;

            if (isUnlock)
            {
                if (isHasData)
                {
                    int quality = CFG_RandomDistribV2.GetItem((int)stats.attr, stats.Num).QualityColor;
                    attrInfo = GameUtil.Instance.GetStrByQuality(quality, CFG_XMLConfig.GetAttributeStrByID((int)stats.attr, stats.Num));
                    outline = GameUtil.Instance.GetQualityOutlineColor(quality);
                }
                else
                {
                    attrInfo = "暂无词条";
                }
            }
            else
            {
                attrInfo = GameUtil.Instance.GetEquipAttributeSlotIsUnlockDesc((RefinedSlot)i);
                if(type == 3)
                {
                    attrInfo = GameUtil.Instance.SetAndCheckStr(attrInfo, 5, limit);
                }
            }
            JUIExtensions.txt(item, "txt").text = attrInfo;
            JUIExtensions.outline(item, "txt").effectColor = outline;
        }
    }

    private void SetRefineItemPreview(EquipSlot equipSlot, Transform item)
    {
        var itemInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(equipSlot).Item : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx, (int)equipSlot);
        bool isNullEquip = itemInfo == null;
        item.Find("icon").gameObject.SetActive(!isNullEquip);

        if (!isNullEquip)
        {
            GameUtil.Instance.SetItemIcon(itemInfo.CfgId, JUIExtensions.img(item.Find("icon")), JUIExtensions.img(item), null, (int)itemInfo.Equip.Quality);
            JUIExtensions.txt(item, "txt_Name").text = "";
        }
        else
        {
            GameUtil.Instance.SetItemBg(0, JUIExtensions.img(item));
            JUIExtensions.txt(item, "txt_Name").text = GameM.Instance.EquipName[(int)equipSlot];
        }
    }

    //选择洗练属性
    private void SelectAttributeToggle(int index)
    {
        for (int i = 0; i < 4; i++)
        {
            GetView().tog_RefineNows[i].isOn = i == index;
        }
    }

    private bool TryGetOneSloeAttributeData(EquipRefineData equipRefineData, RefinedSlot refinedSlot, out Stats stats)
    {
        if (equipRefineData == null)
        {
            stats = null;
            return false;
        }

        var data = equipRefineData.refined_stats.Find(x => x.refined_slot == refinedSlot);
        bool isNullData = data == null;
        if (isNullData)
            stats = null;
        else
            stats = data.stats;

        return !isNullData;
    }

    RefineType refineType;

    //洗练类型选择
    private void OnRefineRandomOrDirect(RefineType refineType)
    {
        this.refineType = refineType;
        if (refineType == RefineType.REFINE_TYPE_RANDOM)
        {
            GameUtil.Instance.SetItemIcon(
                CFG_XMLConfig.xmlConfig.BaseValue.RefineRandomID,
                GetView().img_RefineConsumeIcon);

            GetView().txt_RefineConsumeNum.text = GameUtil.Instance.SetAndCheckStr(
                GameUtil.Instance.FormatNumWithUnit(CFG_XMLConfig.xmlConfig.BaseValue.RefineRandomNum),
                CFG_XMLConfig.xmlConfig.BaseValue.RefineRandomID,
                CFG_XMLConfig.xmlConfig.BaseValue.RefineRandomNum);
        }
        else
        {
            GameUtil.Instance.SetItemIcon(
                CFG_XMLConfig.xmlConfig.BaseValue.RefineAssignID,
                GetView().img_RefineConsumeIcon);

            GetView().txt_RefineConsumeNum.text = GameUtil.Instance.SetAndCheckStr(
                GameUtil.Instance.FormatNumWithUnit(CFG_XMLConfig.xmlConfig.BaseValue.RefineAssignNum),
                CFG_XMLConfig.xmlConfig.BaseValue.RefineAssignID,
                CFG_XMLConfig.xmlConfig.BaseValue.RefineAssignNum);
        }
    }

    bool isRefining = false;
    /// <summary>
    /// 洗练
    /// </summary>
    private void OnRefine()
    {
        if (StateM.Instance.isHasRefineResult)
        {
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = "当前有洗练结果未保留，是否确定保留？",
                YesFun = () =>
                {
                    OnSelectRefineResult(true);
                },
                NoFun = () =>
                {
                    OnSelectRefineResult(false);
                }
            });
            return;
        }

        if (isRefining)
            return;
        isRefining = true;

        //检测货币
        int consumeID = refineType == RefineType.REFINE_TYPE_RANDOM ? CFG_XMLConfig.xmlConfig.BaseValue.RefineRandomID : CFG_XMLConfig.xmlConfig.BaseValue.RefineAssignID;
        int consumeNum = refineType == RefineType.REFINE_TYPE_RANDOM ? CFG_XMLConfig.xmlConfig.BaseValue.RefineRandomNum : CFG_XMLConfig.xmlConfig.BaseValue.RefineAssignNum;
        if (!GameUtil.Instance.ComCoinToCheck(consumeID, consumeNum, null))
        {
            isRefining = false;
            return;
        }

        TcpC.C2S_RefineRequest(refineType, (CharacterSlot)(TeamIdx + 1), _SelectEquipSlot, refinedSlot);
    }

    /// <summary>
    /// 选择洗练结果
    /// </summary>
    /// <param name="isSaveResult">true 保留结果， false 放弃结果</param>
    /// <param name="tipMaxResult">true 提示放弃洗练的最大结果， false 不提示</param>
    public void OnSelectRefineResult(bool isSaveResult, bool tipMaxResult = false)
    {
        if (!StateM.Instance.isHasRefineResult)
            return;

        if (isSaveResult)
        {
            SelectResult(RoleM.Instance.refineResuleCache.new_attributes);
        }
        else
        {
            var slotInfo = RoleM.Instance.GetOneEquipSlotData((CharacterSlot)(TeamIdx + 1), _SelectEquipSlot);
            TryGetOneSloeAttributeData(slotInfo, refinedSlot, out Stats stats);

            if (tipMaxResult && CFG_RandomDistribV2.IsMaxValue((int)RoleM.Instance.refineResuleCache.new_attributes.attr, RoleM.Instance.refineResuleCache.new_attributes.Num))
            {
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = "当前属性已达到极值，是否确定放弃？",
                    YesFun = () =>
                    {
                        SelectResult(stats);
                    },
                });
                return;
            }

            SelectResult(stats);
        }

        void SelectResult(Stats stats)
        {
            TcpC.C2S_RefineSelectRequest((CharacterSlot)(TeamIdx + 1), _SelectEquipSlot, refinedSlot, stats);
            RoleM.Instance.refineResuleCache = null;
            StateM.Instance.isHasRefineResult = false;
        }
    }

    private void OnSync_EquipSlotRefineInfo(GameEvent.EquipSlotRefineInfo _e)
    {
        OnRefineShow();
    }

    int[] refineResultPosX = { 80, 35, -10, -55 };

    //洗练消息返回
    private void OnSync_EquipSlotRefine(GameEvent.EquipSlotRefine _e)
    {
        isRefining = false;
        if (_e.success)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_RemoveResult);
            StateM.Instance.isHasRefineResult = true;
            GetView().obj_RefineResuleButtons.SetActive(true);
            GetView().btn_Refine.gameObject.SetActive(false);
            GetView().obj_ComsumeRoot.SetActive(false);

            //GetView().btn_RefineResuleSave.interactable = true;
            //GetView().btn_RefineResuleDontSave.interactable = true;
            //JUIExtensions.img(GetView().btn_RefineResuleSave.transform).SetGray(false);
            //JUIExtensions.img(GetView().btn_RefineResuleDontSave.transform).SetGray(false);
            SetEnabledRefineToggle(false);
            GetView().tog_RefineRandom.interactable = false;
            GetView().tog_RefineFixed.interactable = false;

            var stats = RoleM.Instance.refineResuleCache.new_attributes;
            bool isMax = CFG_RandomDistribV2.IsMaxValue((int)stats.attr, stats.Num);
            int quality = CFG_RandomDistribV2.GetItem((int)stats.attr, stats.Num).QualityColor;
            GetView().txt_RefineResult.text = GameUtil.Instance.GetStrByQuality(quality, CFG_XMLConfig.GetAttributeStrByID((int)stats.attr, stats.Num));
            //GetView().txt_RefineResult.text = GameUtil.Instance.GetStrByQuality(quality, attrInfo);
            JUIExtensions.outline(GetView().txt_RefineResult.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor(quality);
            GetView().obj_RefineResultMax.SetActive(isMax);

            //SpineUtil.Instance.PlayOnce(GetView().spine_PreviewResult, "t_00", null);
            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (isNullGo) return;
                GetView().ps_PreviewResult.Play();
            }, 0.5f);
            SpineUtil.Instance.PlayOnce(GetView().spine_Hammer, "bomm", null);
        }
    }

    //洗练选择返回刷新
    private void OnSync_EquipSlotRefineSelectResult(GameEvent.EquipSlotRefineSelectResult _e)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_RemoveCover);
        GetView().obj_RefineResuleButtons.SetActive(false);
        GetView().btn_Refine.gameObject.SetActive(true);
        GetView().obj_ComsumeRoot.SetActive(true);

        //GetView().btn_RefineResuleSave.interactable = false;
        //GetView().btn_RefineResuleDontSave.interactable = false;
        //JUIExtensions.img(GetView().btn_RefineResuleSave.transform).SetGray(true);
        //JUIExtensions.img(GetView().btn_RefineResuleDontSave.transform).SetGray(true);
        SetEnabledRefineToggle(true);
        GetView().tog_RefineRandom.interactable = true;
        GetView().tog_RefineFixed.interactable = true;

        GetView().obj_RefineResultMax.SetActive(false);
        Debug.Log("<color=cyan>洗练选择按钮隐藏 2222 ！！！</color>");

        GetView().txt_RefineResult.text = "暂无词条";
        JUIExtensions.outline(GetView().txt_RefineResult.transform).effectColor = defaultOutlineColor;
        OnSelectEquipClick(_SelectEquipSlot, EquipSelectClickFrom.Refresh);
        if (isSaveResult)
        {
            //SpineUtil.Instance.PlayOnce(GetView().spine_SelectResult, $"t_{(int)_e.refinedSlot}", null);
            GetView().root_RefineNowRoot.Find($"Item_{(int)_e.refinedSlot}/eff").GetComponent<ParticleSystem>().Play();
            isSaveResult = false;
        }
    }

    private void SetEnabledRefineToggle(bool enabled)
    {
        for (int i = 0; i < GetView().tog_RefineNows.Length; i++)
        {
            GetView().tog_RefineNows[i].interactable = enabled;
        }
    }

    bool isSaveResult;
    int unlockAttributeSlot;

    private void PlayUnlockSlotAnimation()
    {
        var spine = JUIExtensions.skeleton(GetView().root_RefineNowRoot.Find($"Item_{unlockAttributeSlot}/spine"));
        SpineUtil.Instance.PlayOnce(spine, "suo", null);
        unlockAttributeSlot = 0;
    }

    #endregion

    //帮助
    private void OnHelpBtnClick()
    {
        int idx = 7;
        if (_SelectEquipClass == UIBuild.SelectTwo.Equip_Refine)
            idx = 9;
        else if (_SelectEquipClass == UIBuild.SelectTwo.Equip_Star)
            idx = 8;
        else if (_SelectEquipClass == UIBuild.SelectTwo.Equip_Inset)
            idx = 17;

        UITool.Open<UIComRule>(null, idx);
    }

    private UIBuild.SelectTwo _SelectEquipClass;

    //面板选择
    private void OnSelcetUIBtnClick(UIBuild.SelectTwo eClass)
    {
        var _eClass = eClass;
        if (StateM.Instance.isHasRefineResult)
        {
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = "当前有洗练结果未保留，是否确定保留？",
                YesFun = () =>
                {
                    OnSelectRefineResult(true);
                    OnSelcetUIBtnClick(_eClass);
                },
                NoFun = () =>
                {
                    OnSelectRefineResult(false);
                    OnSelcetUIBtnClick(_eClass);
                }
            });
            return;
        }

        _SelectEquipClass = eClass;

        GetView().ReinforceUI.SetActive(eClass == UIBuild.SelectTwo.Equip_Reinforce);
        GetView().RefineUI.SetActive(eClass == UIBuild.SelectTwo.Equip_Refine);
        GetView().StarUI.SetActive(eClass == UIBuild.SelectTwo.Equip_Star);
        GetView().InsetUI.SetActive(eClass == UIBuild.SelectTwo.Equip_Inset);

        GetView().usedItems[1].obj.SetActive(eClass == UIBuild.SelectTwo.Equip_Reinforce);
        GetView().usedItems[2].obj.SetActive(eClass == UIBuild.SelectTwo.Equip_Refine);
        GetView().usedItems[3].obj.SetActive(eClass == UIBuild.SelectTwo.Equip_Star);

        _SelectEquipSlot = OnEquipHideOrShow();

        switch (eClass)
        {
            case UIBuild.SelectTwo.Equip_Reinforce:
                isAutoReinforce = false;
                break;
            case UIBuild.SelectTwo.Equip_Star:
                IsAutoUpStar = false;

                //RedPointManager.Instance.SetRedClose(RedPointConst.Equip_Star.ToString(), () => {
                //    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Star.ToString(),
                //        RedDotM.Instance.GetEquipStarRedNum());
                //});

                break;
                case UIBuild.SelectTwo.Equip_Refine:

                //RedPointManager.Instance.SetRedClose(RedPointConst.Equip_Refine.ToString(), () => {
                //    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Refine.ToString(),
                //        RedDotM.Instance.GetEquipRefineRedNum());
                //});
                break;
            default:
                break;
        }

        OnSelectEquipClick(_SelectEquipSlot, EquipSelectClickFrom.TagButtonClick);
    }

    private EquipSlot OnEquipHideOrShow()
    {
        int idx = 0;
        for (int i = 1; i <= AllEquipCount; i++)   
        {
            EquipSlot eSlot = (EquipSlot)i;
            //var item = GetView().root_allEquip.Find("item_" + i);
            bool isBest = RoleM.Instance.CheckIsBestEquip(eSlot);
            //item.gameObject.SetActive(isBest || _SelectEquipClass != UIBuild.SelectTwo.Equip_Refine);
            //if (isBest && _SelectEquipClass == UIBuild.SelectTwo.Equip_Refine && idx == 0)
            //{
            //    idx = i;
            //}

            if (isBest) idx = i;
        }

        //GetView().tip_equipTip.gameObject.SetActive(idx == 0 && _SelectEquipClass == UIBuild.SelectTwo.Equip_Refine);
        //if (_SelectEquipClass != UIBuild.SelectTwo.Equip_Refine)
        //{
        //    idx = (int)EquipSlot.Weapon;
        //}

        return idx == 0 ? EquipSlot.Weapon : (EquipSlot)idx;
    }

    #region 魂石 镶嵌
    private void OnInSetUI()
    {
        bool isNullEquip = TeamIdx == 0 ? RoleM.Instance.CheckIsNullEquip(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx, (int)_SelectEquipSlot) == null;
        var _SoltInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot);

        JUIExtensions.img(GetView().Inset_Item.Find("icon")).enabled = !isNullEquip;
        if (isNullEquip)
        {
            GameUtil.Instance.SetItemBg(0, JUIExtensions.img(GetView().Inset_Item));
            JUIExtensions.txt(GetView().Inset_Item.Find("txt_Name")).text = GameM.Instance.EquipName[(int)_SelectEquipSlot];
        }
        else
        {
            GameUtil.Instance.SetItemIcon(_SoltInfo.Item.CfgId, JUIExtensions.img(GetView().Inset_Item.Find("icon")), JUIExtensions.img(GetView().Inset_Item),null, (int)_SoltInfo.Item.Equip.Quality);
            JUIExtensions.txt(GetView().Inset_Item.Find("txt_Name")).text = GameUtil.Instance.GetItemName(_SoltInfo.Item.CfgId);
        }

        int lv = RoleM.Instance.Level;// TeamIdx == 0 ? RoleM.Instance.Level : _PartnerInfo.Level;
        for (int i = 0; i < 4; i++)
        {
            var go = GetView().root_Inset.GetChild(i);
            bool isOpen = CFG_XMLConfig.xmlConfig.Diamond_OpenlvDic[i + 1] <= lv;
            go.Find("lock").gameObject.SetActive(!isOpen);
            if (!isOpen)
            {
                JUIExtensions.txt(go.Find("lock/txt")).text = CFG_XMLConfig.xmlConfig.Diamond_OpenlvDic[i + 1] + "级\n开启";
            }

            bool isInset = _SoltInfo.Gemstone[i] != 0;
            go.Find("icon").gameObject.SetActive(isOpen && isInset);
            if (isInset)
            {
                JUIExtensions.txt(go, "icon/txt_Name").text = GameUtil.Instance.GetItemName(_SoltInfo.Gemstone[i]);
                GameUtil.Instance.SetItemIcon(_SoltInfo.Gemstone[i], JUIExtensions.img(go, "icon"), JUIExtensions.img(go), JUIExtensions.img(go, "icon/Img_natute"));
            }
            else
            {
                GameUtil.Instance.SetItemBg(0, JUIExtensions.img(go));
            }
            go.Find("add").gameObject.SetActive(isOpen && !isInset);
        }

        GetView().txt_InsetTip.text = "提示:"+ CFG_XMLConfig.xmlConfig.Diamond_PositionDes[(int)_SelectEquipSlot];
    }

    private void OnInsetOnClick(int SelectIdx,bool ItemShow)
    {
        int lv = RoleM.Instance.Level;// TeamIdx == 0 ? RoleM.Instance.Level : _PartnerInfo.Level;
        bool isOpen = CFG_XMLConfig.xmlConfig.Diamond_OpenlvDic[SelectIdx + 1] <= lv;
        if (!isOpen)
            return;

        var cid = TeamIdx == 0 ? (int)RoleM.Instance.GetSlotInfoBySolt(_SelectEquipSlot).Gemstone[SelectIdx] : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, (int)_SelectEquipSlot).Gemstone[SelectIdx];

        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (cid != 0)
        {
            if (ItemShow)
            {
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cid });
            }
            else
            {
                TcpC.SetGemstoneChange(new GemstoneChangeRequest { ItemUid = cid, GemstoneSlot = (GemstoneSlot)SelectIdx, Slot = _SelectEquipSlot, OpType = 2,PartnerSlotId = TeamIdx });
            }
        }
        else
        {
            UITool.Open<UIJewelSelect>(null, _SelectEquipSlot, SelectIdx,TeamIdx);
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_TopResouse(GameEvent.EquipInsetRefresh _e)
    {
        if (isNullGo)
            return;
        OnInSetUI();
        OnEquipLvUIRefresh();
        OnEquipSoltBtnRedCheak();
    }

    #endregion

    private int[] resIDs = { (int)ResID.Coin, (int)ResID.QHS, (int)ResID.XLS, (int)ResID.SXS };
    //顶部资源显示
    private void OnTopResouse(bool isInit = false)
    {

        for (int i = 0; i < GetView().usedItems.Length; i++)
        {
            if (isInit)
                GameUtil.Instance.SetItemIcon(resIDs[i], GetView().usedItems[i].img);

            GetView().usedItems[i].txt.text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID(resIDs[i]));
        }
    }

    private int TeamIdx = 0;
    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnSync_BagResouse);
        UniGameEvent.AddEventListener<Sync_TakeSource>(OnSync_TakeSource);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_TopResouse);
        UniGameEvent.AddEventListener<GameEvent.EquipSlotRefineInfo>(OnSync_EquipSlotRefineInfo);
        UniGameEvent.AddEventListener<GameEvent.EquipSlotRefine>(OnSync_EquipSlotRefine);
        UniGameEvent.AddEventListener<GameEvent.EquipSlotRefineSelectResult>(OnSync_EquipSlotRefineSelectResult);
        UniGameEvent.AddEventListener<GameEvent.CommonNotification>(OnRefresh_CommonNotificationPush);
        UniGameEvent.AddEventListener<GameEvent.Sync_TakeSource>(On_Refresh_Buffer);
    }

    /// <summary>
    /// 槽位伙伴信息
    /// </summary>
    private PartnerSlotInfo _PartnerSlotInfo;

    /// <summary>
    /// 伙伴信息
    /// </summary>
    private PartnerInfo _PartnerInfo;

    public void OnClickShow(int idx, UIBuild.SelectTwo selectTwo)
    {
        TeamIdx = idx;
        if (TeamIdx != 0)
        {
            _PartnerSlotInfo = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
            _PartnerInfo = PartnerM.Instance.GetOnePartnerInfoByUid(_PartnerSlotInfo.PartnerUid);
        }

        OnTopResouse(true);
        OnEquipUIBase();
        OnSelcetUIBtnClick(selectTwo);
        OnEquipSoltBtnRedCheak();
    }

    public void OnShow(bool isShow)
    {
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

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnSync_BagResouse);
        UniGameEvent.RemoveEventListener<Sync_TakeSource>(OnSync_TakeSource);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_TopResouse);
        UniGameEvent.RemoveEventListener<GameEvent.EquipSlotRefineInfo>(OnSync_EquipSlotRefineInfo);
        UniGameEvent.RemoveEventListener<GameEvent.EquipSlotRefine>(OnSync_EquipSlotRefine);
        UniGameEvent.RemoveEventListener<GameEvent.EquipSlotRefineSelectResult>(OnSync_EquipSlotRefineSelectResult);
        UniGameEvent.RemoveEventListener<GameEvent.CommonNotification>(OnRefresh_CommonNotificationPush);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_TakeSource>(On_Refresh_Buffer);

        ClearTimer();
        OnClearTime();
        _PartnerInfo = null;
        _PartnerSlotInfo = null;
        EquipEnhanceTimes = 0;
        StateM.Instance.isHasRefineResult = false;


        if (ActivityM.Instance.GetAcitvityEndTime(ActivityType.RefineFund) > 0)
            HttpC.OnGetGrowthFundInfo((int)CFG_Activity.Type.RefineFund);
    }

    private void On_Refresh_Buffer(GameEvent.Sync_TakeSource e)
    {
        OnUpStarShow();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_TopResouse(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        OnTopResouse();
        OnEquipSoltBtnRedCheak();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_BagResouse(Event_Source_Refresh _e)
    {
        if (isNullGo)
            return;

        OnTopResouse();
        OnEquipSoltBtnRedCheak();
        OnRefineShow(false);
        OnRefineRandomOrDirect(refineType);
    }

    private void OnSync_TakeSource(Sync_TakeSource _e)
    {
        OnTopResouse();
        OnRefineShow(false);
        OnRefineRandomOrDirect(refineType);
    }

    private void OnRefresh_CommonNotificationPush(GameEvent.CommonNotification _e)
    {
        if ((CFG_XMLConfig.XML_Config.SwitchOpne)_e.moduleId == CFG_XMLConfig.XML_Config.SwitchOpne.Refined)
        {
            TcpC.C2S_RefineInfoRequest();
            PlayUnlockSlotAnimation();
        }
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}