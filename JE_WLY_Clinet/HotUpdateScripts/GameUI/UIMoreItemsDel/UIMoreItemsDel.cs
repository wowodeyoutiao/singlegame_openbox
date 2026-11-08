using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using pb;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.Event;
using JEngine.MyNet;
using JEngine.UI;
using SuperScrollView;
using HotUpdateScripts;

public class UIMoreItemsDel : UIBase<UIModelMoreItemsDel, UIViewMoreItemsDel>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    //private int select_Job = 0;
    //private int select_Sex = 0;
    private int select_Qulite = 0;
    private int select_Level = 0;
    //筛选列表
    private List<ItemInfo> siftList = new List<ItemInfo>();
    private List<ItemInfo> selectList = new List<ItemInfo>();

    private int MaxDelCount = 100;

    public override void OnInit()
    {

        GetView().btn_Close.onClick.AddListener(() => { UITool.Close<UIMoreItemsDel>(); });
        GetView().btn_Del.onClick.AddListener(OnDestroyEquipClick);
        //GetView().btn_Refresh.onClick.AddListener(OnSiftRefreshClick);

        //GetView().dropdown_Job.AddOptions(GetModel().DroJobList);
        //GetView().dropdown_Job.onValueChanged.AddListener(SetSortJob);

        //GetView().dropdown_Sex.AddOptions(GetModel().DroSexList);
        //GetView().dropdown_Sex.onValueChanged.AddListener(SetSortSex);

        GetView().dropdown_Qulite.AddOptions(GetModel().DroQualityList);
        GetView().dropdown_Qulite.onValueChanged.AddListener(SetSortQulite);

        GetView().dropdown_Level.AddOptions(GetModel().DroLevelList);
        GetView().dropdown_Level.onValueChanged.AddListener(SetSortlevel);

        GetView().tog_Select.onValueChanged.AddListener((bool isOn) =>
        {
            OnSiftRefreshClick(GetView().tog_Select.isOn);
        });

        GetView().infiniteScroll.InitGridView(0, OnFillItem);
    }

    private void OnDestroyEquipClick()
    {
        ItemDisassemblyRequest sendRequest = new ItemDisassemblyRequest();
        sendRequest.ItemUids = new long[selectList.Count];
        for (int i = 0; i < selectList.Count; i++)
        {
            sendRequest.ItemUids[i] = selectList[i].Uid;
        }

        TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_ItemDisassembly, MType.Request, TCPClientMgr.ProtoSerialize(sendRequest));

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "背包",
            position = "装备批量回收",
            event_click = "全部回收",
        });
    }

    private bool IsInit = true;
    private void OnSiftRefreshClick(bool isAdd = false, bool dontClearSelected = false)
    {
        GetSiftInfo(isAdd, dontClearSelected);

        if (IsInit)
        {
            IsInit = false;
            int Count = CFG_XMLConfig.xmlConfig.BaseValue.bagNum.equipnum ;
            GetView().infiniteScroll.SetListItemCount(Count);
        }
        else
            GetView().infiniteScroll.RefreshAllShownItem();

        SetDelBtnState();
    }

    LoopGridViewItem OnFillItem(LoopGridView gridView, int sordIdx, int row, int column)
    {
        LoopGridViewItem _LoopGridViewItem = gridView.NewListViewItem("Item");
        GameObject go = _LoopGridViewItem.gameObject;
        JUIExtensions.exbtn(go.transform).ClearAllListeners();


        var bestEffectTrans = go.transform.Find("BestEffect");
        var bestEffectGo = bestEffectTrans.gameObject;

        bool isSelect = false;
        if (sordIdx < siftList.Count)
        {
            ItemInfo _item = siftList[sordIdx];
            var cfg_item = CFG_Item.GetById(_item.CfgId);
            go.transform.Find("icon").gameObject.SetActive(true);
            GameUtil.Instance.SetItemIcon(_item.CfgId, JUIExtensions.img(go.transform, "icon"), JUIExtensions.img(go.transform), JUIExtensions.img(go.transform, "icon/nature"),(int) _item.Equip.Quality);
            GameUtil.Instance.SetItemJobIcon(_item.CfgId, JUIExtensions.img(go.transform, "icon/img_job"));

            isSelect = selectList.Contains(_item);

            JUIExtensions.exbtn(go.transform).OnClick = () =>
            {
                OnItemClick(_item, go);
            };
            JUIExtensions.exbtn(go.transform).OnPress = () =>
            {
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = _item.CfgId, equipInfo = _item.Equip });
            };

            (int powerNum, int soltIdx) = GameUtil.Instance.GetEquipIsBest(_item);
                go.transform.Find("icon/up").gameObject.SetActive(powerNum > 0 && soltIdx >= 0);
                go.transform.Find("icon/down").gameObject.SetActive(powerNum < 0 && soltIdx >= 0);

            if (_item.Equip != null && (int)_item.Equip.Quality >= 4 && _item.Equip.Attributes != null && _item.Equip.Attributes.Count > 0)
            {
                bestEffectGo.SetActive(true);
                SpineUtil.Instance.PlayItemSpineLoop(bestEffectTrans, (int)_item.Equip.Quality);
            }
            else
                bestEffectGo.SetActive(false);
        }
        else
        {
            GameUtil.Instance.SetItemBg(0, JUIExtensions.img(go.transform));

            go.transform.Find("icon").gameObject.SetActive(false);
            bestEffectGo.SetActive(false);
        }

        go.transform.Find("select").gameObject.SetActive(isSelect);

        JUIExtensions.txt(go.transform.Find("txt_num")).text = "";
        return _LoopGridViewItem;
    }

    private void OnItemClick(ItemInfo _ItemInfo,GameObject go)
    {
        bool isSelect = selectList.Contains(_ItemInfo);
        if (isSelect)
        {
            bool isRemove = selectList.Remove(_ItemInfo);
            go.transform.Find("select").gameObject.SetActive(!isRemove);
        }
        else
        {
            if (MaxDelCount <= selectList.Count)
            {
                UITool.ShowTips("单次批量最多回收100件");
                return;
            }
            selectList.Add(_ItemInfo);
            go.transform.Find("select").gameObject.SetActive(true);
        }
        SetDelBtnState();
    }

    private void SetDelBtnState()
    {
        GetView().btn_Del.gameObject.SetActive(selectList.Count > 0);
        GetView().Obj_tip.gameObject.SetActive(selectList.Count == 0);
    }

    //获取筛选信息
    private void GetSiftInfo(bool isAdd = false, bool dontClearSelected = false)
    {
        siftList.Clear();

        if (!dontClearSelected)
            selectList.Clear();

        var AllItemsList = BagM.Instance.GetEquipItemList(true);
        for (int i = 0; i < AllItemsList.Count; i++)
        {
            var cfg = CFG_Item.GetById(AllItemsList[i].CfgId);

            if ((int)AllItemsList[i].Equip.Quality <= select_Qulite || select_Qulite == 0)
            {
                if (select_Level == 0 || cfg.Levellimit <= select_Level)
                {
                    siftList.Add(AllItemsList[i]);
                    if (isAdd && selectList.Count < MaxDelCount)
                    {
                        selectList.Add(AllItemsList[i]);
                    }
                }
            }
        }

        if (!dontClearSelected)
        {
            for (int i = selectList.Count - 1; i >= 0; i--)
            {
                if (!siftList.Contains(selectList[i]))
                {
                    selectList.RemoveAt(i);
                }
            }
        }
    }


    //private void SetSortJob(int idx)
    //{
    //    select_Job = idx;
    //}
    //private void SetSortSex(int idx)
    //{
    //    select_Sex = idx;
    //}
    private void SetSortQulite(int idx)
    {
        select_Qulite = idx;
        Debug.Log("SetSortQulite   " + idx);
        OnSiftRefreshClick(GetView().tog_Select.isOn);
    }
    private void SetSortlevel(int idx)
    {
        select_Level = idx * 10;
        OnSiftRefreshClick(GetView().tog_Select.isOn);
    }

    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Sync_TakeSource>(Event_TakeSource);

        //GetView().dropdown_Job.value = 0;
        //GetView().dropdown_Level.value = 0;
        GetView().dropdown_Qulite.value = 0;
        GetView().dropdown_Sex.value = 0;
        OnSiftRefreshClick();
        GetView().tog_Select.isOn = false;
    }

    //[Subscriber(ThreadMode.Main)]
    private void Event_TakeSource(GameEvent.Sync_TakeSource _e)
    {
        //GetView().tog_Select.isOn = false;
        OnSiftRefreshClick(dontClearSelected: true);
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_TakeSource>(Event_TakeSource);
    }
    public override void OnDestroy()
    {
        OnExit();
    }
}