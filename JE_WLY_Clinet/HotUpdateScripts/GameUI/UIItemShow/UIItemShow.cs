using System.Collections.Generic;
using MetaFramework.UI;
using pb;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;
using System;
using JEngine.UI;
using UnityEngine.UI;
using System.Text;
using UnityEngine;

public class UIItemShow : UIBase<UIModelItemShow, UIViewItemShow>
{
    private int CfgId { get; set; }

    public override UILayer GetLayer()
    {
        return UILayer.Top;
    }
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIItemShow>();
        });

        GetView().Btn_PreviewStarInfo.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            var cfgPartner= CFG_Partner.GetByItemCfgID(CfgId);
            UITool.Open<UIStarDetails>(null, cfgPartner);
            UITool.Close<UIItemShow>();
        });
    }
    private void BaseItemInfo()
    {
        var cfgInfo = CFG_Item.GetById(CfgId);

        int Quality = ( cfgInfo.Type == (int)CFG_Item.ItemType.Equip && equipInfo !=null) ? (int)equipInfo.Quality : -1;

        GetView().txt_name.text = GameUtil.Instance.GetItemName(CfgId, Quality);
        JUIExtensions.outline(GetView().txt_name.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor(Quality == -1 ? (int)cfgInfo.Quality : Quality);

        GetView().txt_Des.text = cfgInfo.Desc == null ? "" : cfgInfo.Desc;
        GetView().obj_Des.SetActive(cfgInfo.Desc != null);

        GameUtil.Instance.SetItemIcon(CfgId, GetView().Img_icon,GetView().Img_Icon_bg,GetView().Img_Nature, Quality);

        GetView().txt_LevelLimit.text = (cfgInfo.Levellimit == 0 ? $"无" 
            : $"{cfgInfo.Levellimit}");

        //GetView().txt_GenderLimit.text = "性别限制：              " + (cfgInfo.Genderlimit == 0 ? $"<color=#4EC823>无</color>" 
        //    : $"<color=#{(cfgInfo.Genderlimit != RoleM.Instance.Sex ? "CE4031" : "4EC823")}>{CFG_Item.GetSexStrById(cfgInfo.Genderlimit)}</color>");
        if (IsIgnoreJob)
        {
            GetView().txt_OccupyLimit.text = cfgInfo.Occupylimit == 0 ? $"无" : CFG_Item.GetJobStrById((int)cfgInfo.Occupylimit);
        }
        else
        {
            int Job = TeamIdx == 0 ? RoleM.Instance.Job : (CFG_Partner.GetById(PartnerM.Instance.GetOnePartnerInfoByUid(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid).CfgId).job % 10);
            GetView().txt_OccupyLimit.text = cfgInfo.Occupylimit == 0 ? $"无"
                : $"<color=#{(cfgInfo.Occupylimit != Job ? "CE4031" : "4EC823")}>{CFG_Item.GetJobStrById((int)cfgInfo.Occupylimit)}</color>";
        }

        //bool IsEquip = (int)cfgInfo.Type == (int)CFG_Item.ItemType.Equip;
        //var equipRefineData = RoleM.Instance.GetOneEquipSlotData((CharacterSlot)(characterIndex), (EquipSlot)(equipSlotPos));

        bool isShowBestAttr = equipRefineData != null && equipRefineData.refined_stats != null && equipRefineData.refined_stats.Count > 0;

        GetView().obj_Bestattribute.gameObject.SetActive(isShowBestAttr);
        (GetView().txt_attribute0.text, GetView().txt_attribute1.text) = GetItemBaseAttribute(CfgId, Quality);
        GetView().obj_attribute.gameObject.SetActive(GetView().txt_attribute0.text != "");

        var bestEffectTrans = GetView().bestEffect;
        var bestEffectGo = bestEffectTrans.gameObject;
        bestEffectGo.SetActive(false);
        //if (equipInfo != null)
        //{
        //    (GetView().txt_Bestattribute0.text, GetView().txt_Bestattribute1.text) = GetItemBestAttribute(equipSlotPos);
        //    GetView().obj_Bestattribute.gameObject.SetActive(GetView().txt_Bestattribute0.text != "");

        //    if ((int)equipInfo.Quality >= 4 && equipInfo.Attributes != null && equipInfo.Attributes.Count > 0)
        //    {
        //        bestEffectGo.SetActive(true);
        //        SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, (int)equipInfo.Quality);
        //    }
        //}

        if (isShowBestAttr)
        {
            //(GetView().txt_Bestattribute0.text, GetView().txt_Bestattribute1.text) = GetItemBestAttribute(equipSlotPos);
            var sttrs = GetItemBestAttribute(equipSlotPos);

            for (int i = 0; i < 4; i++)
            {
                Transform txtTrans = GetView().obj_Bestattribute.transform.Find($"txt_attribute{i}");
                txtTrans.gameObject.SetActive(i < sttrs.Count);
                if (i < sttrs.Count)
                {
                    JUIExtensions.txt(txtTrans).text = sttrs[i];
                    txtTrans.Find("max").gameObject.SetActive(CFG_RandomDistribV2.IsMaxValue((int)equipRefineData.refined_stats[i].stats.attr, equipRefineData.refined_stats[i].stats.Num));
                }
            }
            GetView().obj_Bestattribute.gameObject.SetActive(GetView().txt_Bestattribute0.text != "");

            if (Quality >= 4)
            {
                bestEffectGo.SetActive(true);
                SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, Quality);
            }
        }

        GetView().Btn_PreviewStarInfo.gameObject.SetActive(cfgInfo.Type == 8);

    }

    private (string, string) GetItemBaseAttribute(int cfgID,int Quality)
    {
        var cfgInfo = CFG_Item.GetById(cfgID);

        float equipValue = Quality > 0? CFG_XMLConfig.xmlConfig.EquipQualityDic[Quality]._value: 1;

        StringBuilder sb1 = new StringBuilder();
        StringBuilder sb2 = new StringBuilder();
        int idx = 0;
        foreach (var item in CFG_XMLConfig.xmlConfig.EquipDisplayDic)
        {
            if (cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) != null && (int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) > 0)
            {
                int _value = (int)Math.Round((int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) * equipValue);
                //int qulite = CFG_RandomDistrib.GetColor(Quality > 0 ? Quality :(int)cfgInfo.Quality , CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value);
                string arrtStr = CFG_XMLConfig.GetAttributeStrByID((int)CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value, true);
                if (idx % 2 == 0)
                    sb1.AppendLine(arrtStr);
                else
                    sb2.AppendLine(arrtStr);
            }
        }

        return (sb1.ToString().TrimEnd('\n'), sb2.ToString().TrimEnd('\n'));
    }

    /*
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
                    //int qulite = CFG_RandomDistrib.GetColor((int)eInfo.Quality, item.SID, BestDic[item.SID]);
                    string arrtStr = CFG_XMLConfig.GetAttributeStrByID(item.SID, BestDic[item.SID],true);
                    if (idx % 2 == 0)
                        bestStr1 = $"{bestStr1}{arrtStr}\n";
                    else
                        bestStr2 = $"{bestStr2}{arrtStr}\n";
                }
            }
        }
        return (bestStr1, bestStr2);
    }
    */
    private List<string> GetItemBestAttribute(int equipSlotPos)
    {
        List<string> bestAttributes = new List<string>();
        if (equipSlotPos > 0 && equipSlotPos < 13)
        {
            StringBuilder sb = new StringBuilder();
            equipRefineData.refined_stats.Sort((a, b) =>
            {
                return a.refined_slot.CompareTo(b.refined_slot);
            });

            for (int i = 0; i < equipRefineData.refined_stats.Count; i++)
            {
                var attrData = equipRefineData.refined_stats[i];
                int quality = CFG_RandomDistribV2.GetItem((int)attrData.stats.attr, attrData.stats.Num).QualityColor;
                bestAttributes.Add(GameUtil.Instance.GetStrByQuality(quality, CFG_XMLConfig.GetAttributeStrByID((int)attrData.stats.attr, attrData.stats.Num)));
            }
        }

        return bestAttributes;
    }

    private bool isMaxBestAttribute(out int maxAttributePos)
    {
        if (equipSlotPos > 0 && equipSlotPos < 13)
        {
            for (int i = 0; i < equipRefineData.refined_stats.Count; i++)
            {
                var attrData = equipRefineData.refined_stats[i];
                bool isMax = CFG_RandomDistribV2.IsMaxValue((int)attrData.stats.attr, attrData.stats.Num);
                if (isMax)
                {
                    maxAttributePos = i;
                    return true;
                }
            }
            maxAttributePos = 0;
            return false;
        }
        else
        {
            maxAttributePos = 0;
            return false;
        }
    }

    //private bool IsShowBestAttribute()
    //{
    //    return characterIndex != 0 && equipSlotPos != 0;
    //}

    private EquipInfo equipInfo;
    private EquipRefineData equipRefineData;

    private int TeamIdx;
    //忽略职业判断，观察其他玩家装备时，不用判断职业类型。
    private bool IsIgnoreJob;

    private int characterIndex;
    private int equipSlotPos;

    public override void OnEnter(params object[] args)
    {
        //CfgId = (int)args[0];
        //if (args.Length > 1)
        //    equipInfo = args[1] as EquipInfo;
        //else
        //    equipInfo = null;

        //TeamIdx = args.Length > 2 ? ((args[2] is null) ? 0 : (int)args[2]) : 0;
        //IsIgnoreJob = args.Length > 3 && (bool)args[3];

        Params @params = (Params)args[0];
        CfgId = @params.cfgId;
        equipInfo = @params.equipInfo;
        TeamIdx = @params.teamIdx;
        IsIgnoreJob = @params.isIgnoreJob;
        characterIndex = @params.characterIndex;
        equipSlotPos = @params.equipSlotPos;
        equipRefineData = @params.equipRefineData;

        BaseItemInfo();

        var callBack = @params.callback;
        GetView().obj_BtnNode.SetActive(callBack != null);

        var btnName = @params.buttonName;
        if (callBack != null)
        {
            GetView().Btn_Yes.GetComponentInChildren<Text>().text = btnName;
            GetView().Btn_Yes.onClick.RemoveAllListeners();
            GetView().Btn_Yes.onClick.AddListener(() =>
            {
                callBack();
            });
        }
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }

    public struct Params
    {
        public int cfgId;
        public EquipInfo equipInfo;
        public int teamIdx;
        public bool isIgnoreJob;
        public int characterIndex;
        public int equipSlotPos;
        public Action callback;
        public string buttonName;
        public EquipRefineData equipRefineData;
    }
}