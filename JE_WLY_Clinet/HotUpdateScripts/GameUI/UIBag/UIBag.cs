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

using static GameEvent;


public class UIBag : UIBase<UIModelBag, UIViewBag>
{
    private enum BagType :int
    { 
        No =0,
        All ,
        Equip ,
        ParentPice,
        Jewel,
        Other ,
    }
    private Dictionary<BagType, string> DicBtnName = new Dictionary<BagType, string>
    {
        [BagType.All] = "全部",
        [BagType.Equip] = "装备",
        [BagType.ParentPice] = "伙伴碎片",
        [BagType.Jewel] = "魂石",
        [BagType.Other] = "材料",
    };
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private List<ItemInfo> ShowEquipList;
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(CloseUI);
        GetView().Btn_Refresh.onClick.AddListener(OnBtnRefreshClick);
        GetView().Btn_Del.onClick.AddListener(OnDelClick);

        //GetView().btn_AutoDel.onClick.AddListener(OnAutoDelClick);
        //GetView().btn_AutoDeSet.onClick.AddListener(OnAutoDelSetClick);

        for (int i = 0; i < GetView().TopBtns.Length; i++)
        {
            int idx = i + 1;
            GetView().TopBtns[i].transform.Find("Label").GetComponent<Text>().text = DicBtnName[(BagType)idx];
            GetView().TopBtns[i].transform.Find("Select/Label").GetComponent<Text>().text = DicBtnName[(BagType)idx];
            GetView().TopBtns[i].onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnPageUIChange((BagType)idx); 
            });
        }
        GetView().infiniteScroll.InitGridView(0, OnFillItem);

        GetView().btn_syn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            var selectSyn = UISynthesis.UISelect.Material;
            switch (_SelectBagType)
            {
                //case BagType.ParentPice:
                //    selectSyn = UISynthesis.UISelect.PartnerPice;
                //    break;
                case BagType.Jewel:
                    selectSyn = UISynthesis.UISelect.Jewel;
                    break;
                default:
                    break;
            }
            UITool.Open<UISynthesis>(null,null, selectSyn);

            if (GetView().btn_syn.transform.Find("dot").gameObject.activeSelf)
            {
                RedPointManager.Instance.SetRedClose(RedPointConst.Bag_Syn.ToString(), () => {
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Bag_Syn);
                });
            }
        });
    }

    LoopGridViewItem OnFillItem(LoopGridView gridView, int index, int row, int column)
    {
        int sordIdx = (index);
        
        LoopGridViewItem _LoopGridViewItem = gridView.NewListViewItem("Item");
        GameObject go = _LoopGridViewItem.gameObject;
        go.GetComponent<Button>().onClick.RemoveAllListeners();

        var bastEffectTrans = go.transform.Find("BastEffect");
        bastEffectTrans.gameObject.SetActive(false);

        var _Icon = go.transform.Find("icon");

        go.gameObject.name = sordIdx.ToString();
        if (sordIdx < ShowEquipList.Count)
        {
            ItemInfo _item = ShowEquipList[sordIdx];
            var cfg_item = CFG_Item.GetById(_item.CfgId);
            _Icon.gameObject.SetActive(true);
            bool IsEquip = (int)CFG_Item.ItemType.Equip == cfg_item.Type;

            GameUtil.Instance.SetItemIcon(_item.CfgId, JUIExtensions.img(_Icon), JUIExtensions.img(go.transform), JUIExtensions.img(_Icon, "nature"), IsEquip ? (int)_item.Equip.Quality : -1);
            GameUtil.Instance.SetItemJobIcon(_item.CfgId, JUIExtensions.img(_Icon, "img_job"));

            JUIExtensions.txt(go.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(ShowEquipList[sordIdx].Num);

            if (IsEquip)
            {
                (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(_item);
                _Icon.transform.Find("up").gameObject.SetActive(powerNum > 0 && soltIdx >= 0);
                _Icon.transform.Find("down").gameObject.SetActive(powerNum < 0 && soltIdx >= 0);

                if (ShowEquipList[sordIdx].Equip != null && (int)ShowEquipList[sordIdx].Equip.Quality >= 4 && ShowEquipList[sordIdx].Equip.Attributes != null && ShowEquipList[sordIdx].Equip.Attributes.Count > 0)
                {
                    bastEffectTrans.gameObject.SetActive(true);
                    SpineUtil.Instance.PlayItemSpineLoop(bastEffectTrans, (int)ShowEquipList[sordIdx].Equip.Quality);
                }
            }
            else
            {
                _Icon.transform.Find("up").gameObject.SetActive(false);
                _Icon.transform.Find("down").gameObject.SetActive(false);
            }

            JUIExtensions.btn(go.transform).onClick.AddListener(() => {
                OnItemClick(_item);
            });
        }
        else
        {
            _Icon.gameObject.SetActive(false);
            JUIExtensions.txt(go.transform.Find("txt_num")).text = "";
            GameUtil.Instance.SetItemBg(0, JUIExtensions.img(go.transform));
        }
        
        return _LoopGridViewItem;
    }

    private void OnItemClick(ItemInfo item)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        UITool.Open<UIItemUseOrDel>(null, item);
    }

    private void OnDelClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIMoreItemsDel>();
    }

    #region 自动分解 目前屏蔽
    //自动分解 目前屏蔽
    private void OnAutoDelClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if ((RoleM.Instance.Level < CFG_XMLConfig.xmlConfig.BaseValue.autoDelLv && RoleM.Instance.Vip < CFG_XMLConfig.xmlConfig.BaseValue.autoDelVipLv))
        {
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = "该功能需要VIP" + CFG_XMLConfig.xmlConfig.BaseValue.autoDelVipLv + "开启，是否前往提升？",
                YesFun = () => {
                    UITool.Close<UIBag>();
                    UITool.Open<UIVIP>();
                },
            });
            return;
        }
        TcpC.C2S_AutoDisassembly(new AutoDisassemblyInfo
        {
            Open = !GetView().obj_AutoDel_on.activeSelf,
            Level = BagM.Instance.AutoDelSetInfo.Level,
            LevelSwitch = BagM.Instance.AutoDelSetInfo.LevelSwitch,
            Quality = BagM.Instance.AutoDelSetInfo.Quality,
            QualitySwitch = BagM.Instance.AutoDelSetInfo.QualitySwitch,
            HighCombat = BagM.Instance.AutoDelSetInfo.HighCombat,
        });
    }

    //自动分解设置
    private void OnAutoDelSetClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIAutoDelSetting>();
    }

    //自动分解显示
    private void OnAutoDesShow()
    {
        GetView().obj_AutoDel_on.SetActive(BagM.Instance.AutoDelSetInfo.Open);
        bool isCanAutoDel = RoleM.Instance.Level >= CFG_XMLConfig.xmlConfig.BaseValue.autoDelLv || RoleM.Instance.Vip >= CFG_XMLConfig.xmlConfig.BaseValue.autoDelVipLv;

        GetView().txt_AutoDe.text = isCanAutoDel ? "自动分解" : $"自动分解(VIP{CFG_XMLConfig.xmlConfig.BaseValue.autoDelVipLv}或{CFG_XMLConfig.xmlConfig.BaseValue.autoDelLv}级解锁)";
        GetView().btn_AutoDeSet.gameObject.SetActive(isCanAutoDel);
    }
    //[Subscriber(ThreadMode.Main)]
    //private void AutoDelSet(GameEvent.AutoDelSet _e)
    //{
    //    if (!isNullGo)
    //        OnAutoDesShow();
    //}

    //[Subscriber(ThreadMode.Main)]
    //private void AutoDelSet(GameEvent.Sync_UserBase _e)
    //{
    //    if (!isNullGo)
    //        OnAutoDesShow();
    //}
    #endregion

    private bool _RefreshCD = true;
    private void OnBtnRefreshClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (_RefreshCD)
        {
            OnPageUIChange(_SelectBagType);
            _RefreshCD = false;

            GameHelper.Instance.AddInvoke(new TimeItem(5f, () =>
            {
                _RefreshCD = true;
            }));
        }
        else
            UITool.ShowTips("操作太频繁，请稍后");
    }

    private BagType _SelectBagType = BagType.No;

    private void OnPageUIChange(BagType Idx)
    {

        for (int i = 0; i < GetView().TopBtns.Length; i++)
        {
            GetView().TopBtns[i].transform.Find("Select").gameObject.SetActive((i + 1) == (int)Idx);
        }

        bool isRefresh = _SelectBagType == Idx;
        _SelectBagType = Idx;
        int Count = CFG_XMLConfig.xmlConfig.BaseValue.bagNum.totlenum;

        if (Idx == BagType.Equip)
        {
            Count = CFG_XMLConfig.xmlConfig.BaseValue.bagNum.equipnum;
            ShowEquipList = BagM.Instance.GetEquipItemList(true);
        }
        else if (Idx == BagType.Other)
        {
            Count = CFG_XMLConfig.xmlConfig.BaseValue.bagNum.itemnum;
            ShowEquipList = BagM.Instance.GetOtherItemList(true);
        }
        else if (Idx == BagType.Jewel)
        {
            Count = CFG_XMLConfig.xmlConfig.BaseValue.bagNum.itemnum;
            ShowEquipList = BagM.Instance.JewelItemListList(true);
        }
        else if (Idx == BagType.ParentPice)
        {
            Count = CFG_XMLConfig.xmlConfig.BaseValue.bagNum.itemnum;
            ShowEquipList = BagM.Instance.GetPartnerPicesList(true);
        }
        else
            ShowEquipList = BagM.Instance.ItemList;

        GetView().infiniteScroll.SetListItemCount(Count);

        GameHelper.Instance.SetDeayTimeDo(() => {
            GetView().infiniteScroll.RefreshAllShownItem();
        },0.1f);

            //if (isRefresh)
            //    GetView().infiniteScroll.RefreshAllShownItem();

        GetView().btn_syn.gameObject.SetActive(Idx != BagType.Equip && Idx != BagType.ParentPice);
        GetView().Btn_Del.gameObject.SetActive(Idx == BagType.Equip);
        //GetView().btn_AutoDel.gameObject.SetActive(Idx == BagType.Equip);
        //if (Idx == BagType.Equip)
        //    OnAutoDesShow();
    }

    

    private void CloseUI()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Close<UIBag>();
    }

    private void OnBindRed()
    {
        RedPointManager.Instance.RPS.Register(RedPointConst.Bag_Syn.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().btn_syn.transform.Find("dot").gameObject.SetActive(node.pointCount > 0);
        });

    }
    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(Event_Refresh);
        UniGameEvent.AddEventListener<ChangeEquip>(ChangeEquip);
        UniGameEvent.AddEventListener<GameEvent.Sync_TakeSource>(Event_TakeSource);

        OnBindRed();
        OnPageUIChange(BagType.All);
    }

    //[Subscriber(ThreadMode.Main)]
    private void Event_Refresh(Event_Source_Refresh _e)
    {
        OnPageUIChange(_SelectBagType);
    }

    //[Subscriber(ThreadMode.Main)]
    private void Event_TakeSource(GameEvent.Sync_TakeSource _e)
    {
        OnPageUIChange(_SelectBagType);
    }

    //[Subscriber(ThreadMode.Main)]
    private void ChangeEquip(GameEvent.ChangeEquip _e)
    {
        OnPageUIChange(_SelectBagType);
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(Event_Refresh);
        UniGameEvent.RemoveEventListener<ChangeEquip>(ChangeEquip);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_TakeSource>(Event_TakeSource);
        _RefreshCD = true;
    }

    public override void OnDestroy()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(Event_Refresh);
        UniGameEvent.RemoveEventListener<ChangeEquip>(ChangeEquip);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_TakeSource>(Event_TakeSource);
    }
}

