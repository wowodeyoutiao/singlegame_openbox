using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Helper;
using JEngine.UI;
using HotUpdateScripts.ProtoM;
using SuperScrollView;
using HotUpdateScripts;
using JEngine.Event;
using pb;

public class UIPartnerUp : UIBase<UIModelPartnerUp, UIViewPartnerUp>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private int select_Qulite = 0;
    private string guidanceName = null;
    public override void OnInit()
    {
        guidanceName = null;

        GetView().btn_Close.onClick.AddListener(() => { 
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPartnerUp>(); 
        });

        GetView().dropdown_Qulite.AddOptions(GetModel().DroStrList);
        GetView().dropdown_Qulite.onValueChanged.AddListener(SetSortQulite);

        GetView().infiniteScroll.InitGridView(0, OnFillItem);

        GetView().btn_up.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            EquipOrReplacePartner();
        });

        NewPlayerGuideWidget.Instance.AddGuideGameObject("UIPartnerUp_EnterClick", GetView().btn_up.transform);
    }

    private bool CheakIsUpStar(pb.PartnerInfo partner)
    {
        var cfgInfo = CFG_Partner.GetById(partner.CfgId);
        if (partner.Star != cfgInfo.StarConsume.Length)
        {
            var userdData = cfgInfo.GetStarConsumeByLv(partner.Star + 1);

            bool isCanUp = true;
            for (int j = 0; j < userdData.Length; j++)
            {
                if (!GameUtil.Instance.ComCoinToCheck(userdData[j].key, userdData[j].value, null, false))
                {
                    isCanUp = false;
                    break;
                }
            }
            if (isCanUp)
                return true;
        }
        return false;
    }


    LoopGridViewItem OnFillItem(LoopGridView gridView, int sordIdx, int row, int column)
    {
        LoopGridViewItem _LoopGridViewItem = gridView.NewListViewItem("Item");
        Transform go = _LoopGridViewItem.transform;
        go.name = "Item" + sordIdx;
        var touch = JUIExtensions.exbtn(go.transform);
        touch.ClearAllListeners();

        go.Find("select").gameObject.SetActive(false);
        var _item = siftList[sordIdx];
        JUIExtensions.txt(go.Find("txt_Star")).text = _item.Star + "星";
        JUIExtensions.txt(go.Find("img_lv/txt_lv")).text = _item.Level.ToString() ;
        //JUIExtensions.txt(go.Find("txt_Name")).text = CFG_Partner.GetById(_item.CfgId).Name;

        GameUtil.Instance.SetJobIcon(CFG_Partner.GetById(_item.CfgId).job % 10, JUIExtensions.img(go, "img_job"));
        GameUtil.Instance.SetItemIcon(CFG_Partner.GetById(_item.CfgId).ItemID, JUIExtensions.img(go.Find("icon")));

        go.Find("select").gameObject.SetActive(false);
        go.Find("unselect").gameObject.SetActive(PartnerM.Instance.IsInSlot(_item.Uid));

        go.Find("txt_StarTip").gameObject.SetActive(CheakIsUpStar(_item));

        touch.OnClick = () => { OnItemClick(_item.Uid, go.gameObject); };
        touch.OnPress = () => { UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFG_Partner.GetById(_item.CfgId).ItemID }); };

        GameUtil.Instance.SetItemBg(CFG_Partner.GetById(_item.CfgId).Quality, JUIExtensions.img(go));

        if (!PartnerM.Instance.IsInSlot(_item.Uid) && guidanceName == null)
        {
            //UnityEngine.Debug.LogError($"首个可上阵item {go.name}");
            guidanceName = go.name;
        }

        return _LoopGridViewItem;
    }

    private GameObject SelectObj = null;
    public long Uid = 0;
    private void OnItemClick(long _Uid, GameObject go)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (SelectObj == go || go.transform.Find("unselect").gameObject.activeSelf)
        {
            return;
        }
        else
        {
            if (SelectObj != null)
                SelectObj.transform.Find("select").gameObject.SetActive(false);

            Uid = _Uid;
            SelectObj = go;
            go.transform.Find("select").gameObject.SetActive(true);
        }
        OnSelectState();
    }

    private void OnSelectState()
    {
        GetView().btn_up.gameObject.SetActive(Uid != 0);
        GetView().Obj_tip.gameObject.SetActive(Uid == 0);
    }

    private void SetSortQulite(int idx)
    {
        select_Qulite = GetModel().DroQualityList[idx];
        SetSiftList();
    }

    private List<pb.PartnerInfo> AllPartners;
    private List<pb.PartnerInfo> siftList = new List<pb.PartnerInfo>();

    private int SoltIdx = 0;
    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(CommonResponse);

        SoltIdx = (int)args[0];

        GetView().txt_Btn.text = PartnerM.Instance.GetOnePartnerSlotInfo(SoltIdx).PartnerUid == 0? "上阵" :"替换";

        SelectObj = null;
        Uid = 0;

        AllPartners =  PartnerM.Instance.GetPartners();
        AllPartners.Sort((a, b) => {
            int qualityA = CFG_Partner.GetById(a.CfgId).Quality;
            int qualityB = CFG_Partner.GetById(b.CfgId).Quality;

            if (qualityA != qualityB)
            {
                return qualityB.CompareTo(qualityA);
            }
            else
            {
                return a.CfgId.CompareTo(b.CfgId);
            }
        });

        GetView().dropdown_Qulite.value = 0;
        SetSortQulite(0);
        OnSelectState();
    }

    private void EquipOrReplacePartner()
    {
        if (GetView().txt_Btn.text == "替换")
        {
            TcpC.C2S_PartnerReplace(SoltIdx, Uid);
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "角色",
                component = "伙伴上阵界面",
                position = "",
                event_click = "伙伴替换",
            });
        }
        else
        {
            TcpC.C2S_PartnerCombat(SoltIdx, Uid);
            UITool.Close<UIPartnerUp>();
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "角色",
                component = "伙伴上阵界面",
                position = "",
                event_click = "伙伴上阵",
            });
        }
    }

    private void SetSiftList()
    {
        siftList.Clear();
        for (int i = 0; i < AllPartners.Count; i++)
        {
            if (CFG_Partner.GetById(AllPartners[i].CfgId).Quality >= select_Qulite || select_Qulite == 0)
            {
                siftList.Add(AllPartners[i]);
            }
        }

        GetView().infiniteScroll.SetListItemCount(siftList.Count);
        GetView().infiniteScroll.RefreshAllShownItem();

        if (guidanceName != null)
        {
            var curRoot = GetView().infiniteScroll.transform.Find($"Viewport/Content/{guidanceName}");
            if (curRoot != null)
            {
                NewPlayerGuideWidget.Instance.AddGuideGameObject("UIPartnerUp_FirstRolePos", curRoot.Find("icon"));
                NewPlayerGuideWidget.Instance.AddGuideGameObject("UIPartnerUp_FirstRole", curRoot);
            }
        }
    }

    //[Subscriber(ThreadMode.Main)]
    private void CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if ( _e.cid == CmdId.C2S_PartnerReplace)
        {
            if  (_e.code == 0)
            {
                UITool.Close<UIPartnerUp>();
            }
            
        }
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(CommonResponse);
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}