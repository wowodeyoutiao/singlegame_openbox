using System.Collections.Generic;
using UnityEngine;
using MetaFramework.UI;
using pb;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using SuperScrollView;
using HotUpdateScripts;

public class UISynthesisSelect : UIBase<UIModelSynthesisSelect, UIViewSynthesisSelect>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private int select_Job = 0;
    private int select_Sex = 0;
    private int select_Qulite = 0;
    private int select_Level = 0;
    //筛选列表
    private List<ItemInfo> siftList = new List<ItemInfo>();

    private ItemInfo _selectItem;
    private ItemInfo selectItem { set 
        {
            GetView().btn_Syn.gameObject.SetActive(value != null);
            GetView().btn_AutoSyn.gameObject.SetActive(value != null);
            GetView().Obj_tip.SetActive(value == null);
            _selectItem = value;
        }
        get { return _selectItem; } }
    private List<ItemInfo> AllItemsList = new List<ItemInfo>();


    public override void OnInit()
    {

        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UISynthesisSelect>(); 
        });
        GetView().btn_Syn.onClick.AddListener(OnBtnSynClilck);

        GetView().btn_AutoSyn.onClick.AddListener(()=> {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UISynthesisSelect>();
            
            var add = new UISynthesis.SelectSyn { id = selectItem.CfgId, num = 1, Uid = selectItem.Uid };
            UITool.Get<UISynthesis>().OnAutoAdd(add);
        });

        GetView().btn_Refresh.onClick.AddListener(OnSiftRefreshClick);

        GetView().dropdown_Job.AddOptions(GetModel().DroJobList);
        GetView().dropdown_Job.onValueChanged.AddListener(SetSortJob);

        GetView().dropdown_Sex.AddOptions(GetModel().DroSexList);
        GetView().dropdown_Sex.onValueChanged.AddListener(SetSortSex);

        GetView().dropdown_Qulite.AddOptions(GetModel().DroQualityList);
        GetView().dropdown_Qulite.onValueChanged.AddListener(SetSortQulite);

        GetView().dropdown_Level.AddOptions(GetModel().DroLevelList);
        GetView().dropdown_Level.onValueChanged.AddListener(SetSortlevel);

        GetView().infiniteScroll.InitGridView(0, OnFillItem) ;
    }

    private void OnBtnSynClilck()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (selectItem == null)
        {
            UITool.ShowTips("未选中道具");
            return;
        }
        UISynthesis.selectSyns.Add(new UISynthesis.SelectSyn { id = selectItem.CfgId, Uid = selectItem.Uid, num = 1 });
        UITool.Get<UISynthesis>().OnUISelectShow();
        UITool.Close<UISynthesisSelect>();
    }

    //private bool IsInit = true;
    private void OnSiftRefreshClick()
    {
        GetSiftInfo();

        if (GetView().infiniteScroll.ItemTotalCount != CFG_XMLConfig.xmlConfig.BaseValue.bagNum.equipnum)
            GetView().infiniteScroll.SetListItemCount(CFG_XMLConfig.xmlConfig.BaseValue.bagNum.equipnum);

        GetView().infiniteScroll.RefreshAllShownItem();
    }

    private LoopGridViewItem OnFillItem(LoopGridView gridView, int sordIdx, int row, int column)
    {
        LoopGridViewItem _LoopGridViewItem = gridView.NewListViewItem("Item");
        GameObject go = _LoopGridViewItem.gameObject;
        JUIExtensions.touch(go.transform).RemoveAllListeners();
        go.transform.Find("select").gameObject.SetActive(sordIdx < siftList.Count && selectItem == siftList[sordIdx]);

        if (sordIdx < siftList.Count && selectItem == siftList[sordIdx])
            seletGo = go;
        if (sordIdx < siftList.Count)
        {
            ItemInfo _item = siftList[sordIdx];
            var cfg_item = CFG_Item.GetById(_item.CfgId);
            go.transform.Find("icon").gameObject.SetActive(true);
            GameUtil.Instance.SetItemIcon(_item.CfgId, JUIExtensions.img(go.transform, "icon"), JUIExtensions.img(go.transform), JUIExtensions.img(go.transform, "icon/nature"));
            
            JUIExtensions.touch(go.transform).AddClickListener((p) => {
                if (!UITool.Check<UIItemShow>() && !p.dragging)
                    OnItemClick(_item, go); 
            });

            JUIExtensions.touch(go.transform).AddPointerDownListener((p) =>
            {
                _TouchTime = GameHelper.Instance.AddInvoke(new TimeItem(GameConst.TouchTime, () => {
                    if (!p.dragging)
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = _item.CfgId });
                }));
            });

            JUIExtensions.touch(go.transform).AddPointerUpListener((p) =>
            {
                ClearTouch();
            });

            JUIExtensions.txt(go.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(_item.Num);
        }
        else
        {
            go.transform.Find("icon").gameObject.SetActive(false);
            JUIExtensions.txt(go.transform.Find("txt_num")).text = "";

            GameUtil.Instance.SetItemBg(0, JUIExtensions.img(go.transform));
        }
            
        

        return _LoopGridViewItem;
          
    }

    private TimeItem _TouchTime;
    private void ClearTouch()
    {
        if (_TouchTime != null)
        {
            GameHelper.Instance.RemoveInvoke(_TouchTime);
            _TouchTime = null;
        }
    }

    private GameObject seletGo;
    private void OnItemClick(ItemInfo _ItemInfo,GameObject go)
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
    }

    //获取筛选信息
    private void GetSiftInfo()
    {
        siftList.Clear();

        for (int i = 0; i < AllItemsList.Count; i++)
        {
            var cfg = CFG_Item.GetById(AllItemsList[i].CfgId);
            if (select_Job == 0 || cfg.Occupylimit == select_Job)
            {
                if (select_Sex == 0 || cfg.Genderlimit == select_Sex)
                {
                    if (cfg.Quality <= select_Qulite || select_Qulite == 0)
                    {
                        if (select_Level == 0 || cfg.Levellimit <= select_Level)
                        {
                            siftList.Add(AllItemsList[i]);
                        }
                    }
                }
            }
        }
        selectItem = null;
    }
   

    private void SetSortJob(int idx)
    {
        select_Job = idx;
        OnSiftRefreshClick();
    }
    private void SetSortSex(int idx)
    {
        select_Sex = idx;
        OnSiftRefreshClick();
    }
    private void SetSortQulite(int idx)
    {
        select_Qulite = idx;
        OnSiftRefreshClick();
    }
    private void SetSortlevel(int idx)
    {
        select_Level = idx * 10;
        OnSiftRefreshClick();
    }
    
    public override void OnEnter(params object[] args)
    {
        int uiType = (int)args[0];
        int cid = UISynthesis.selectSyns.Count == 0 ? 0 : UISynthesis.selectSyns[0].id;

        string titleStr = "材料合成";
        int sonType = 1;
        if (uiType == (int)UISynthesis.UISelect.Jewel)
        {
            titleStr = "魂石合成";
            sonType = 3;
        }
        else if (uiType == (int)UISynthesis.UISelect.Used)
        {
            titleStr = "消耗合成";
            sonType = 4;
        }

        GetView().txt_Tiltle.text = titleStr;

        AllItemsList = BagM.Instance.GetSynListByType(sonType, UISynthesis.selectSyns,true);

        GetView().dropdown_Job.value = cid == 0 ? 0: (int)CFG_Item.GetById(cid).Occupylimit;
        GetView().dropdown_Level.value = cid == 0 ? 0 : (int)CFG_Item.GetById(cid).Levellimit / 10;
        GetView().dropdown_Qulite.value = cid == 0  ? 0 : (int)CFG_Item.GetById(cid).Quality;
        GetView().dropdown_Sex.value = cid == 0 ? 0 : (int)CFG_Item.GetById(cid).Genderlimit;

        OnSiftRefreshClick();
    }

    public override void OnExit()
    {
        ClearTouch();
    }
    public override void OnDestroy()
    {
        OnExit();
    }
}