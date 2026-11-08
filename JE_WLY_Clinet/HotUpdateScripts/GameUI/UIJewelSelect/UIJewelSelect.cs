using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using pb;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;
using SuperScrollView;

public class UIJewelSelect : UIBase<UIModelJewelSelect, UIViewJewelSelect>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIJewelSelect>();
        });

        GetView().btn_Select.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.SetGemstoneChange(new GemstoneChangeRequest { ItemUid = selectItem.CfgId, GemstoneSlot = (GemstoneSlot)SelectIdx, Slot = (EquipSlot)_SelectEquipSlot, OpType = 1,PartnerSlotId = TeamIdx });
            UITool.Close<UIJewelSelect>();
        });

        GetView().btn_Syn.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIJewelSelect>();
            UITool.Open<UISynthesis>(null, selectItem);
        });

        GetView().infiniteScroll.InitGridView(0, OnFillItem);
    }

    LoopGridViewItem OnFillItem(LoopGridView gridView, int sordIdx, int row, int column)
    {
        LoopGridViewItem _LoopGridViewItem = gridView.NewListViewItem("Item");
        GameObject go = _LoopGridViewItem.gameObject;
        var touch = JUIExtensions.exbtn(go.transform);
        touch.ClearAllListeners();
        go.transform.Find("select").gameObject.SetActive(sordIdx < siftList.Count && selectItem == siftList[sordIdx]);
        if (sordIdx < siftList.Count)
        {
            ItemInfo _item = siftList[sordIdx];
            var cfg_item = CFG_Item.GetById(_item.CfgId);
            go.transform.Find("icon").gameObject.SetActive(true);
            GameUtil.Instance.SetItemIcon(_item.CfgId, JUIExtensions.img(go.transform, "icon"), JUIExtensions.img(go.transform), JUIExtensions.img(go.transform, "icon/nature"));
            JUIExtensions.txt(go.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID(_item.CfgId));
            touch.OnClick = () => { OnItemClick(_item, go); };
            touch.OnPress = () => { UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = _item.CfgId }); };
        }
        else
        {
            GameUtil.Instance.SetItemBg(0, JUIExtensions.img(go.transform));
            go.transform.Find("icon").gameObject.SetActive(false);
            JUIExtensions.txt(go.transform.Find("txt_num")).text = "";
        }
        return _LoopGridViewItem;

    }

    private GameObject seletGo;
    private ItemInfo selectItem;
    private void OnItemClick(ItemInfo _ItemInfo, GameObject go)
    {
        if (selectItem == _ItemInfo)
        {
            selectItem = null;
            go.transform.Find("select").gameObject.SetActive(false);
        }
        else
        {
            if (go != seletGo && seletGo != null)
            {
                seletGo.transform.Find("select").gameObject.SetActive(false);
            }
            selectItem = _ItemInfo;
            seletGo = go;
            go.transform.Find("select").gameObject.SetActive(true);
        }
        OnBtnShow();
    }

    private bool IsInit = true;
    private void OnSiftRefreshClick()
    {
        if (IsInit)
        {
            IsInit = false;
            int Count = CFG_XMLConfig.xmlConfig.BaseValue.bagNum.itemnum ;
            GetView().infiniteScroll.SetListItemCount(Count);
        }
        else
            GetView().infiniteScroll.RefreshAllShownItem();

    }

    private void OnBtnShow()
    {
        GetView().Obj_tip.SetActive(selectItem == null);
        GetView().btn_Select.gameObject.SetActive(selectItem != null);
        GetView().btn_Syn.gameObject.SetActive(selectItem != null && CFG_Item.GetById(selectItem.CfgId).ComposeMark == 1);
    }

    private List<ItemInfo> siftList = new List<ItemInfo>();

    private int _SelectEquipSlot;
    private int SelectIdx;
    private int TeamIdx;
    public override void OnEnter(params object[] args)
    {
        _SelectEquipSlot = (int)args[0];
        SelectIdx = (int)args[1];
        TeamIdx = (int)args[2];
        var cfgids = CFG_Item.GetCfgIdsByEquipSolt(_SelectEquipSlot);
        siftList.Clear();
        selectItem = null;

        for (int i = 0; i < cfgids.Count; i++)
        {
            var item = BagM.Instance.GetItemByCFGID(cfgids[i]);
            if (item !=  null)
            siftList.Add(item);
        }
        OnSiftRefreshClick();
        OnBtnShow();
    }

    public override void OnExit()
    {
        siftList.Clear();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}