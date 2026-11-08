using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using pb;
using JEngine.MyNet;
using System;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using HotUpdateScripts;
using JEngine.UI;

public class UIItemUseOrDel : UIBase<UIModelItemUseOrDel, UIViewItemUseOrDel>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    private ItemInfo _itemInfo;
    private int usedNum;
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            CloseUI(); 
        });

        GetView().Btn_Use.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBtnUseClick(false);
        });
        GetView().Btn_Del.onClick.AddListener(OnBtnDelClick);

        GetView().btn_Add.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnUsedNumChang(true); 
        });

        GetView().btn_Jian.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnUsedNumChang(false); });

        GetView().btn_Max.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            SetMaxUsedCount();
        });

        GetView().btn_UsedMore.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBtnUseClick(true);
        });

        GetView().btn_closeUsedMore.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().UIUseNumChoice.SetActive(false);
        });

        
        GetView().Btn_Syn.onClick.AddListener(() => {

            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIItemUseOrDel>();
            UITool.Open<UISynthesis>(null, _itemInfo);
            
        });

        GetView().Btn_Retain.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIItemUseOrDel>();
            TcpC.C2S_BoxResIntoBag(_itemInfo.Uid);
        });

        GetView().input_Used.onValueChanged.AddListener((string a) => {
            if (int.Parse(a) > _itemInfo.Num)
            {
                GetView().input_Used.text = "" + _itemInfo.Num;
            }
            usedNum = int.Parse(GetView().input_Used.text);
            });
    }

    private bool IsEquip;
    private int SoltIdx = 0;
    //基础信息
    private void BaseItemInfo()
    {
        var cfgInfo = CFG_Item.GetById(_itemInfo.CfgId);
        IsEquip = (int)CFG_Item.ItemType.Equip == cfgInfo.Type;
        int Quality = (!IsEquip ? -1 : (int)_itemInfo.Equip.Quality);

        GetView().txt_name.text = GameUtil.Instance.GetItemName(_itemInfo.CfgId, Quality);
        JUIExtensions.outline(GetView().txt_name.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor(IsEquip ? Quality : (int)cfgInfo.Quality);

        GetView().txt_Des.text = cfgInfo.Desc == null ? "" : cfgInfo.Desc;
        GetView().obj_Des.SetActive(GetView().txt_Des.text != "");

        int powerChange = 0;

        if (IsEquip)
            (powerChange, SoltIdx) = GameUtil.Instance.GetEquipIsBest(_itemInfo, SoltIdx);

        bool isGreenColorStr = (IsEquip && SoltIdx >= 0) || (cfgInfo.Levellimit == 0 || cfgInfo.Levellimit <= RoleM.Instance.Level);
        GetView().txt_LevelLimit.text = $"<color=#{(isGreenColorStr ? "47B022" : "E53E39")}>" +
           $"{(cfgInfo.Levellimit == 0 ? "无" : cfgInfo.Levellimit)}</color>";

        JUIExtensions.outline(GetView().txt_LevelLimit.transform).effectColor = isGreenColorStr ? new Color(0.13f, 0.45f, 0.06f) : new Color(0.47f, 0.14f, 0.1f);

        bool isGreenColorOutline = (IsEquip && SoltIdx >= 0) || (cfgInfo.Occupylimit == 0 || cfgInfo.Occupylimit == RoleM.Instance.Job);
        GetView().txt_OccupyLimit.text = $"<color=#{(isGreenColorOutline ? "47B022" : "E53E39")}>" +
            $"{(cfgInfo.Occupylimit == 0 ? "无" : CFG_Item.GetJobStrById((int)cfgInfo.Occupylimit) )}</color>";

        JUIExtensions.outline(GetView().txt_OccupyLimit.transform).effectColor = isGreenColorOutline ? new Color(0.13f, 0.45f, 0.06f) : new Color(0.47f, 0.14f, 0.1f);

        GameUtil.Instance.SetItemIcon(_itemInfo.CfgId, GetView().Img_icon,GetView().Img_Icon_bg, GetView().Img_Nature, Quality);

        var isPartner = cfgInfo.Type == (int)CFG_Item.ItemType.PartnerPice;
        if (isPartner)
        {
            var CFGPartner = CFG_Partner.GetByItemID(_itemInfo.CfgId);
            isPartner = CFGPartner != null && CFGPartner.ChipGoodNum <= _itemInfo.Num && !PartnerM.Instance.CheakPartnerIsHave(CFGPartner.ID);
        }
        GetView().Btn_Use.gameObject.SetActive(isPartner || (CFG_Item.GeCanUsed((int)cfgInfo.Type) && (!IsEquip || SoltIdx >= 0)));
        GetView().Btn_Del.gameObject.SetActive(CFG_Item.GeCanDel((int)cfgInfo.Type));

        GetView().Btn_Retain.gameObject.SetActive(isRetain);

        GetView().UIUseNumChoice.SetActive(false);


        GetView().txt_delBtn.text = IsEquip && CFG_Item.GeCanDel((int)cfgInfo.Type) ? "回收" : "出售";
       // GetView().txt_delBtn.text = "回收";
        GetView().Obj_equip.SetActive(IsEquip);
        GetView().txt_btn_Use.text = IsEquip && SoltIdx >= 0 ? ((SoltIdx ==  0 && RoleM.Instance.CheckIsNullEquip((EquipSlot)cfgInfo.Equiptype)
            || PartnerM.Instance.GetOnePartnerOneEquip(SoltIdx, (int)cfgInfo.Equiptype) == null)? "穿戴" : "替换") : "使用";

        GetView().line_attribute.SetActive(IsEquip || (int)cfgInfo.Type == (int)CFG_Item.ItemType.Jewel);
        GetView().line_attributeBest.SetActive(IsEquip);

        GetView().Btn_Syn.gameObject.SetActive(cfgInfo.ComposeMark == 1);

        GetView().obj_equipUp.SetActive(powerChange > 0);
        GetView().obj_equipDown.SetActive(powerChange < 0);

        if (IsEquip)
        {
            (GetView().txt_attributeBast0.text, GetView().txt_attributeBast1.text) = GetItemBestAttribute(_itemInfo.CfgId, _itemInfo.Equip);
            (GetView().txt_attribute0.text, GetView().txt_attribute1.text) = GetItemBaseAttribute(_itemInfo.CfgId, Quality);
            GetView().line_attributeBest.SetActive(GetView().txt_attributeBast0.text != "");

            var bestEffectTrans = GetView().baseBestEffect;
            var bestEffectGo = bestEffectTrans.gameObject;
            if (_itemInfo.Equip != null && (int)_itemInfo.Equip.Quality >= 4 && _itemInfo.Equip.Attributes != null && _itemInfo.Equip.Attributes.Count > 0)
            {
                bestEffectGo.SetActive(true);
                SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, (int)_itemInfo.Equip.Quality);
            }
            else
                bestEffectGo.SetActive(false);
        }
        else if ((int)cfgInfo.Type == (int)CFG_Item.ItemType.Jewel)
        {
            (GetView().txt_attribute0.text, GetView().txt_attribute1.text) = GetItemBaseAttribute(_itemInfo.CfgId, -1);
            GetView().line_attribute.SetActive(GetView().txt_attribute0.text != "");
        }

        GameHelper.Instance.SetDeayFrame(() => {
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_attribute0.GetComponent<RectTransform>());
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_attribute1.GetComponent<RectTransform>());
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_attributeBast0.GetComponent<RectTransform>());
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_attributeBast1.GetComponent<RectTransform>());
        });

        GetView().baseTrans.localPosition = new Vector3(244, -396, 0);
        GetView().now_rect_base.gameObject.SetActive(false);
        if (IsEquip && (
            (SoltIdx == 0 && !RoleM.Instance.CheckIsNullEquip((EquipSlot)cfgInfo.Equiptype) || PartnerM.Instance.GetOnePartnerOneEquip(SoltIdx, (int)cfgInfo.Equiptype) != null)
            ))
        {
            NowEquipItemInfo((EquipSlot)cfgInfo.Equiptype, SoltIdx);
        }else
        {
            GetView().baseTrans.localPosition = new Vector3(244, 0, 0); ;
        }
    }

    //已穿戴
    private void NowEquipItemInfo(EquipSlot nowSlot,int SoltIdx)
    {
        GetView().now_rect_base.gameObject.SetActive(true);
        var eqipNowInfo = SoltIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt(nowSlot).Item: PartnerM.Instance.GetOnePartnerOneEquip(SoltIdx, (int)nowSlot);
        var cfgInfo = CFG_Item.GetById(eqipNowInfo.CfgId);

        int Quality = (eqipNowInfo.Equip == null ? -1 : (int)eqipNowInfo.Equip.Quality);

        GameUtil.Instance.SetItemIcon(eqipNowInfo.CfgId, GetView().now_Img_icon,GetView().now_Img_Icon_bg,GetView().now_Img_Nature, Quality);
        GetView().now_txt_name.text = GameUtil.Instance.GetItemName(eqipNowInfo.CfgId, Quality);

        JUIExtensions.outline(GetView().now_txt_name.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor(Quality);

        GetView().now_txt_Des.text = cfgInfo.Desc == null ? "" : cfgInfo.Desc;
        GetView().now_obj_Des.SetActive(GetView().now_txt_Des.text != "");

        GetView().now_txt_LevelLimit.text = cfgInfo.Levellimit == 0 ? $"无" : $"<color=#47B022>{cfgInfo.Levellimit}</color>";
        GetView().now_txt_OccupyLimit.text = cfgInfo.Occupylimit == 0 ? $"无" : cfgInfo.Occupylimit == 0 ? $"无" : $"<color=#47B022>{CFG_Item.GetJobStrById((int)cfgInfo.Occupylimit)}</color>";

        (GetView().now_txt_attribute0.text, GetView().now_txt_attribute1.text) = GetItemBaseAttribute(eqipNowInfo.CfgId, Quality);
        (GetView().now_txt_attributeBest0.text, GetView().now_txt_attributeBest1.text) = GetItemBestAttribute(eqipNowInfo.CfgId, eqipNowInfo.Equip);

        GetView().now_line_attributeBest.SetActive(GetView().now_txt_attributeBest0.text != "");

        var bestEffectTrans = GetView().compBestEffect;
        var bestEffectGo = bestEffectTrans.gameObject;
        if (eqipNowInfo.Equip != null && (int)eqipNowInfo.Equip.Quality >= 4 && eqipNowInfo.Equip.Attributes != null && eqipNowInfo.Equip.Attributes.Count > 0)
        {
            bestEffectGo.SetActive(true);
            SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, (int)eqipNowInfo.Equip.Quality);
        }
        else
            bestEffectGo.SetActive(false);
    }

    //高级属性
    private (string, string) GetItemBaseAttribute(int cfgID, int Quality)
    {
        var cfgInfo = CFG_Item.GetById(cfgID);

        float equipValue = Quality > 0 ? CFG_XMLConfig.xmlConfig.EquipQualityDic[Quality]._value : 1;
        string Str1 = "";
        string Str2 = "";
        int idx = 0;
        foreach (var item in CFG_XMLConfig.xmlConfig.EquipDisplayDic)
        {
            if (cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) != null && (int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) > 0)
            {
                int _value = (int)Math.Round((int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) * equipValue);
                //int qulite = CFG_RandomDistrib.GetColor(Quality > 0 ? Quality : (int)cfgInfo.Quality, CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value);
                string arrtStr = CFG_XMLConfig.GetAttributeStrByID((int)CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value, true);
                if (idx % 2 == 0)
                    Str1 = $"{Str1}{arrtStr}\n";
                else
                    Str2 = $"{Str2}{arrtStr}\n";
            }
        }
        return (Str1, Str2);
    }

    //极品属性
    private (string, string) GetItemBestAttribute(int cfgID, EquipInfo eInfo)
    {
        var cfgInfo = CFG_Item.GetById(cfgID);
        string bestStr1 = "";
        string bestStr2 = "";
        int idx = 0;
        if (eInfo != null)
        {
            Dictionary<int, int> BestDic = new Dictionary<int, int>();
            foreach (var item in eInfo.Attributes)
            {
                BestDic.Add((int)item.attr, item.Num);
            }

            foreach (var item in CFG_XMLConfig.xmlConfig.RolePanelList)
            {
                if (BestDic.ContainsKey(item.SID) && BestDic[item.SID] > 0)
                {
                    int qulite = CFG_RandomDistrib.GetColor((int)eInfo.Quality, item.SID, BestDic[item.SID]);
                    string arrtStr = CFG_XMLConfig.GetAttributeStrByID(item.SID, BestDic[item.SID]);
                    if (idx % 2 == 0)
                        bestStr1 = $"{bestStr1}{arrtStr}\n";
                    //bestStr1 = $"{bestStr1}{GameUtil.Instance.GetStrByQuality(qulite, arrtStr)}\n";
                    else
                        bestStr2 = $"{bestStr2}{arrtStr}\n";
                        //bestStr2 = $"{bestStr2}{GameUtil.Instance.GetStrByQuality(qulite, arrtStr)}\n";
                }
            }
        }
        return (bestStr1, bestStr2);
    }

    //放置到背包
    private bool isRetain = false;

    int argsLength;

    public override void OnEnter(params object[] args)
    {
        _itemInfo = args[0] as ItemInfo;
        //SoltIdx = 0;

        isRetain = args.Length > 1 ? (bool)args[1] : false;

        SoltIdx = args.Length > 2 ? (int)args[2] : -1;

        argsLength = args.Length;


        BaseItemInfo();
        ForceLayout();

        //只显示没有按钮
        if (args.Length > 2)
        {
            GetView().Btn_Use.gameObject.SetActive(false);
            GetView().Btn_Del.gameObject.SetActive(false);
            GetView().Btn_Syn.gameObject.SetActive(false);
            GetView().Btn_Retain.gameObject.SetActive(false);
        }
    }

    private void ForceLayout()
    {
        LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().rect_base);
        LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().now_rect_base);
    }

    private void CloseUI()
    {
        UITool.Close<UIItemUseOrDel>();
    }

    public void OnUsedNumChang(bool isAdd)
    {
        if (isAdd)
        {
            usedNum++;
            usedNum = Mathf.Min(usedNum, (int)_itemInfo.Num);
        }
        else
        {
            usedNum--;
            usedNum = Mathf.Max(usedNum, 1);
        }
        GetView().input_Used.text = "" + usedNum;
    }

    private void SetMaxUsedCount()
    {
        usedNum = (int)_itemInfo.Num;
        GetView().input_Used.text = "" + usedNum;
    }

    public void ClothEquip()
    {
        OnBtnUseClick(false);
    }

    private void OnBtnUseClick(bool isMore =false)
    {
        //伙伴碎片使用
        if (CFG_Item.GetById(_itemInfo.CfgId).Type == (int)CFG_Item.ItemType.PartnerPice)
        {
            TcpC.C2S_PartnerCompose(CFG_Partner.GetByItemID(_itemInfo.CfgId).ID);
            CloseUI();
            return;
        }

        if (IsEquip)
        {
            TcpC.C2S_ChangeEquip(_itemInfo.Uid, 1, (int)CFG_Item.GetById(_itemInfo.CfgId).Equiptype, SoltIdx);
            CloseUI();

            string _page = argsLength == 2 ? "开箱结果" : "包裹";
            string _component = argsLength == 2 ? "点击装备道具" : "装备页";
            string _position = argsLength == 2 ? "" : "点击道具";
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = _page,
                component = _component,
                position = _position,
                event_click = "穿戴装备",
            });
        }
        else
        {
            var itemUseCfg = CFG_ItemUse.GetById(_itemInfo.CfgId);
            if (itemUseCfg != null && itemUseCfg.Type == 4)
            {
                if (itemUseCfg.Givegoods != null && itemUseCfg.Givegoods.Length > 0)
                {
                    UITool.Open<UIItemDesignatedUse>(null, _itemInfo.CfgId, _itemInfo.Uid);
                    UITool.Close<UIItemUseOrDel>();
                }
                else
                    Debug.LogError($"ItemUse表，ID：{_itemInfo.CfgId}，配置错误");
                return;
            }

            if (!isMore && _itemInfo.Num > 1 && !NeedHideCountSelectPopup())
            {
                GetView().UIUseNumChoice.SetActive(true);
                usedNum = 1;
                GetView().input_Used.text = "" + usedNum;
                return;
            }

            usedNum = isMore ? int.Parse(GetView().input_Used.text) : 1;

            // 无双宠物兑换卡
            if (_itemInfo.CfgId == (int)CFG_Item.ItemID.PetChangeCard)
            {
               var PetChangeCardItem = CFG_Summon.GetPetExchangeCard();
                if (_itemInfo.Num < PetChangeCardItem.Consume[0].num)
                    UITool.ShowTips($"{CFG_Item.GetById(_itemInfo.CfgId).Name}道具不足，无法兑换！");
                else
                {
                    HttpC.GetPetReceive(1, PetChangeCardItem.Drawtype);

                    HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                    {
                        page = "包裹",
                        component = "材料、全部页",
                        position = "点击道具",
                        event_click = "使用",
                    });
                }
                    
            }
            // 其他道具
            else 
            {
                TcpC.C2S_ItemUse(_itemInfo.Uid, usedNum);
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "包裹",
                    component = "材料、全部页",
                    position = "点击道具",
                    event_click = "使用",
                });
            }

            CloseUI();
        }
    }

    private bool NeedHideCountSelectPopup()
    {
        return _itemInfo.CfgId == (int)CFG_Item.ItemID.PetChangeCard;
    }

    //装备销毁，道具删除
    private void OnBtnDelClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        void DelFun() {
            if (IsEquip)
            {
                ItemDisassemblyRequest sendRequest = new ItemDisassemblyRequest();
                sendRequest.ItemUids = new long[] { _itemInfo.Uid };
                TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ItemDisassembly, MType.Request, TCPClientMgr.ProtoSerialize(sendRequest));

                string _page = argsLength == 2 ? "开箱结果" : "包裹";
                string _component = argsLength == 2 ? "点击装备道具" : "装备页";
                string _position = argsLength == 2 ? "" : "点击道具";
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = _page,
                    component = _component,
                    position = _position,
                    event_click = "回收",
                });
            }
            else
            {
                TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ItemDestroy, MType.Request, TCPClientMgr.ProtoSerialize(new ItemDestroyRequest()
                {
                    ItemUid = _itemInfo.Uid,
                    Num = (int)_itemInfo.Num,
                }));

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "包裹",
                    component = "装备页",
                    position = "点击道具",
                    event_click = "出售",
                });
            }
            CloseUI();

        }

        string timeStr = GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear.ToString();
        var key = isRetain ? GameConst.PlayerPrefsKey.BoxRecycle.ToString() : GameConst.PlayerPrefsKey.BagRecycle.ToString();
        if (GameUtil.Instance.GetPlayerPrefs(key) == timeStr)
        {
            DelFun();
        }
        else
        {
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = $"确认{(IsEquip ? "回收" : "出售")}全部当前道具:{GameUtil.Instance.GetItemName(_itemInfo.CfgId, _itemInfo.Equip != null ? (int)_itemInfo.Equip.Quality : -1)}？",
                isTogg = true,
                TogNoFun = () =>
                {
                    GameUtil.Instance.SetPlayerPrefs(key, "");
                },
                TogYesFun = () =>
                {
                    GameUtil.Instance.SetPlayerPrefs(key, timeStr);
                },
                YesFun = DelFun,
            });
        }
    }

    public Vector2 GetUseButtonPositionRelativeCanvas()
    {
        return JUIExtensions.GetUILocalPositionRelativeCanvas(GetView().Btn_Use.transform);
    }

    public override void OnExit()
    {
        
    }
}