using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;
using pb;
using JEngine.Helper;
using JEngine.MyNet;
using HotUpdateScripts.ProtoM;

public class UIEquipOff : UIBase<UIModelEquipOff, UIViewEquipOff>
{
    private ItemInfo _itemInfo { get;  set; }

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIEquipOff>();
        });

        GetView().Btn_Off.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ChangeEquip, MType.Request, 
                TCPClientMgr.ProtoSerialize(new ChangeEquipRequest() 
                {
                    ItemUid = _itemInfo.Uid, OpType = 2, Slot = (EquipSlot)CFG_Item.GetById(_itemInfo.CfgId).Equiptype 
                }));
            UITool.Close<UIEquipOff>();
        });
    }
    private void BaseItemInfo()
    {
        var cfgInfo = CFG_Item.GetById(_itemInfo.CfgId);
        GetView().txt_name.text = GameUtil.Instance.GetItemName(_itemInfo.CfgId);
        GetView().txt_Des.text = cfgInfo.Desc == null ? "无" : cfgInfo.Desc;
        GameUtil.Instance.SetItemIcon(_itemInfo.CfgId, GetView().Img_icon,GetView().Img_Icon_bg);

        GetView().txt_LevelLimit.text = "等级限制：              "+ (cfgInfo.Levellimit == 0 ? "<color=#4EC823>无</color>" 
            : $"<color=#{(cfgInfo.Levellimit > RoleM.Instance.Level ? "CE4031" : "4EC823")}>{cfgInfo.Levellimit}</color>");

        GetView().txt_GenderLimit.text = "性别限制：              "+ (cfgInfo.Genderlimit == 0 ? $"<color=#4EC823>无</color>" 
            : $"<color=#{(cfgInfo.Genderlimit != RoleM.Instance.Sex ? "CE4031" : "4EC823")}>{CFG_Item.GetSexStrById((int)cfgInfo.Genderlimit)}</color>");
        
        GetView().txt_OccupyLimit.text = "职业限制：              " +( cfgInfo.Occupylimit == 0 ? $"<color=#4EC823>无</color>" 
            : $"<color=#{(cfgInfo.Occupylimit != RoleM.Instance.Job ? "CE4031" : "4EC823")}>{CFG_Item.GetJobStrById((int)cfgInfo.Occupylimit)}</color>");


        (GetView().txt_attribute0.text, GetView().txt_attribute1.text) = GetItemAttribute(_itemInfo.CfgId, _itemInfo.Equip);

        LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_attribute0.GetComponent<RectTransform>());
        LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_attribute1.GetComponent<RectTransform>());
    }

    private (string, string) GetItemAttribute(int cfgID, EquipInfo eInfo)
    {
        var cfgInfo = CFG_Item.GetById(cfgID);

        string Str = "";
        foreach (var item in CFG_XMLConfig.xmlConfig.EquipDisplayDic)
        {
            if (cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) != null && (int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo) > 0)
            {
                int _value = (int)cfgInfo.GetType().GetField(item.Key).GetValue(cfgInfo);
                int qulite = CFG_RandomDistrib.GetColor((int)cfgInfo.Quality, CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value);
                string arrtStr = CFG_XMLConfig.GetAttributeStrByID((int)CFG_XMLConfig.xmlConfig.EquipDisplayToSeverDic[item.Key], _value);
                Str = $"{ Str }{ GameUtil.Instance.GetStrByQuality(qulite, arrtStr)}\n";
            }
        }
        string bestStr = "";
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
                    int qulite = CFG_RandomDistrib.GetColor((int)cfgInfo.Quality, item.SID, BestDic[item.SID]);
                    string arrtStr = CFG_XMLConfig.GetAttributeStrByID(item.SID, BestDic[item.SID]);
                    
                    bestStr = $"{bestStr}{GameUtil.Instance.GetStrByQuality(qulite, arrtStr)}\n";
                }
            }
        }
        return (Str, bestStr);
    }

    public override void OnEnter(params object[] args)
    {
        _itemInfo =  RoleM.Instance.GetSlotInfoBySolt((EquipSlot)args[0]).Item ;
        BaseItemInfo();
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}