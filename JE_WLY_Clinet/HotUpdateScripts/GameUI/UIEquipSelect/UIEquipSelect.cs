using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using SuperScrollView;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using static UIFightFormation;

public class UIEquipSelect : UIBase<UIModelEquipSelect, UIViewEquipSelect>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIEquipSelect>();
        });

        GetView().btn_Equip.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (_SelectItem == null)
            {
                UITool.ShowTips("请先选择一件装备");
                return;
            }
            TcpC.C2S_ChangeEquip(_SelectItem.Uid, 1, (int)CFG_Item.GetById(_SelectItem.CfgId).Equiptype, _TeamSolt);
            UITool.Close<UIEquipSelect>();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "角色",
                component = "点击装备图标",
                position = "点击替换",
                event_click = "确认穿戴",
            });
        });

        GetView().btn_OpenBox.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIEquipSelect>();
            UITool.Close<UIRole>();
            JEvent.defaultEvent.Post(new GameEvent.OnHall_OpenBox());
        });

        GetView().lgv_Equip.InitGridView(0, OnFillItem);
    }

    private LoopGridViewItem OnFillItem(LoopGridView gridView, int sordIdx, int row, int column)
    {
        LoopGridViewItem _LoopGridViewItem = gridView.NewListViewItem("Item");
        GameObject go = _LoopGridViewItem.gameObject;
        JUIExtensions.touch(go.transform).RemoveAllListeners();

        if (sordIdx < 0 || sordIdx >= _ShowEquipList.Count)
            return null;

        ItemInfo _item = _ShowEquipList[sordIdx];
        var cfg_item = CFG_Item.GetById(_item.CfgId);

        GameUtil.Instance.SetItemIcon(_item.CfgId, JUIExtensions.img(go.transform, "IconImg"), JUIExtensions.img(go.transform), _Quality :(int) _item.Equip.Quality);

        JUIExtensions.btn(go.transform).onClick.RemoveAllListeners();
        JUIExtensions.btn(go.transform).onClick.AddListener(() => {
            UITool.Open<UIItemUseOrDel>(null, _item, false, _TeamSolt);
        });

        JUIExtensions.btn(go.transform.Find("SelButton")).onClick.RemoveAllListeners();
        JUIExtensions.btn(go.transform.Find("SelButton")).onClick.AddListener(() => {
            if (_SelectItem == _item)
                return;
            _SelectItem = _item;
            if (_SelectObj != null)
            {
                _SelectObj.SetActive(false);
            }
            _SelectObj = go.transform.Find("SelButton/on").gameObject;
            _SelectObj.SetActive(true);

        });

        go.transform.Find("SelButton/on").gameObject.SetActive(_SelectItem == _item);
        if (_SelectItem == _item)
        {
            _SelectObj = go.transform.Find("SelButton/on").gameObject;
        }

        JUIExtensions.txt(go.transform.Find("NameText")).text = GameUtil.Instance.GetItemName(_item.CfgId, (int)_item.Equip.Quality);
        JUIExtensions.outline(go.transform.Find("NameText")).effectColor = GameUtil.Instance.GetQualityOutlineColor((int)_item.Equip.Quality);

        var bestEffectTrans = go.transform.Find("BestEffect");
        var bestEffectGo = bestEffectTrans.gameObject;
        if (_item.Equip != null && (int)_item.Equip.Quality >= 4 && _item.Equip.Attributes != null && _item.Equip.Attributes.Count > 0)
        {
            bestEffectGo.SetActive(true);
            SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, (int)_item.Equip.Quality);
        }
        else
            bestEffectGo.SetActive(false);

        return _LoopGridViewItem;

    }

    private List<ItemInfo> _ShowEquipList = new List<ItemInfo>();

    private ItemInfo _SelectItem;
    private GameObject _SelectObj;
    private int _TeamSolt;

    public override void OnEnter(params object[] args)
    {
        _SelectItem = null;
        _SelectObj = null;

        int _equipSolt = (int)args[0];
        _TeamSolt = (int)args[1];

        int _job = RoleM.Instance.Job;
        int _lv = RoleM.Instance.Level;
        if (_TeamSolt != 0)
        {
            var partner = PartnerM.Instance.GetOnePartnerInfoByUid(PartnerM.Instance.GetOnePartnerSlotInfo(_TeamSolt).PartnerUid);
            _lv = partner.Level;
            _job = CFG_Partner.GetById(partner.CfgId).job % 10;
        }

        var EquipsList = BagM.Instance.GetEquipItemList(true);
        _ShowEquipList.Clear();
        for (int i = 0; i < EquipsList.Count; i++)
        {
            var cfgEquip = CFG_Item.GetById(EquipsList[i].CfgId);
            if ((cfgEquip.Occupylimit == 0 || cfgEquip.Occupylimit == _job) && cfgEquip.Levellimit <= _lv && _equipSolt == cfgEquip.Equiptype)
            {
                _ShowEquipList.Add(EquipsList[i]);
            }
        }
        GetView().lgv_Equip.SetListItemCount(_ShowEquipList.Count);
        GetView().lgv_Equip.RefreshAllShownItem();

        GetView().btn_Equip.gameObject.SetActive(_ShowEquipList.Count != 0);
        GetView().btn_OpenBox.gameObject.SetActive(_ShowEquipList.Count == 0);
        GetView().obj_NoTip.gameObject.SetActive(_ShowEquipList.Count == 0);
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}