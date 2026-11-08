using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using HotUpdateScripts;
using JEngine.UI;
using SuperScrollView;
using JEngine.Event;
using static GameEvent;

//称号
public class UIRoleTitile : UIBase<UIModelRoleTitile, UIViewRoleTitile>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {

        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIRoleTitile>();
        });

        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            int idx = i;
            GetView().teamItems[i]._btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (idx > 0 && CFG_XMLConfig.CheckPartnerSoltIsOpenServer(idx))
                {
                    UITool.ShowTips($"{CFG_XMLConfig.RemainingDaysPartnerSlotUnlockOpen(idx)}天后开启");
                    return;
                }
                if(idx > 0)
                {
                    (bool isLock, string strtxttip, string strtip) = CheakIsLock(CFG_XMLConfig.GetOpenSoltLv(idx), CFG_XMLConfig.GetOpenSolMainStaget(idx));

                    if (isLock)
                    {
                        UITool.ShowTips(strtip);
                        return;
                    }
                }

                //if (idx > 0 && CFG_XMLConfig.GetOpenSoltLv(idx) > RoleM.Instance.Level)
                //{
                //    UITool.ShowTips("角色等级不足");
                //    return;
                //}

                OnTeamItemClick(idx);
            });
        }

        //称号激活
        GetView().btn_TitleActive.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            for (int i = 0; i < CFGInfo.Activity.Length; i++)
            {
                var Iteminfo = CFGInfo.Activity[i].values[0];
                if (BagM.Instance.GetNumByCfgID(Iteminfo.key) < Iteminfo.value)
                {
                    UITool.ShowTips("称号碎片不足");
                    return;
                }
            }
            TcpC.C2S_TitleCompose(_SelectCFGId);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "角色",
                component = "称号界面",
                position = CFGInfo.Name,
                event_click = "激活称号",
            });
        });

        //称号称号装备
        GetView().btn_TitleEquip.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_TitleEquip( _SelectCFGId, _TeamIdx);
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "角色",
                component = "称号界面",
                position = CFGInfo.Name,
                event_click = "佩戴称号",
            });
        });

        GetView().btn_TitleEquipDel.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_TitleEquip(0, _TeamIdx);
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "角色",
                component = "称号界面",
                position = CFGInfo.Name,
                event_click = "卸下称号",
            });
        });

        //称号共鸣
        GetView().btn_TitleResonate.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UITitleActivity>(null,_TeamIdx);
        });

        //称号升星
        GetView().btn_UpStar.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            for (int i = 0; i < CFGInfo.StarConsume.Length; i++)
            {
                var Iteminfo = CFGInfo.StarConsume[i].values[0];
                if (BagM.Instance.GetNumByCfgID(Iteminfo.key) < Iteminfo.value)
                {
                    UITool.ShowTips("称号碎片不足");
                    return;
                }
            }

            TcpC.C2S_TitleStart(_SelectCFGId,_TitleUpStarID);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "角色",
                component = "称号界面",
                position = CFGInfo.Name,
                event_click = "称号升星",
            });
        });

        //幸运符选择
        GetView().root_UpStarItem.GetComponent<Button>().onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UITitleLuckChoice>();
            
        });

        GetView().btn_LuckDel.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnLuckItemSelect(0);
        });

        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 47);
        });

        GetView().loopView_Title.InitListView(0, LeftViewFill);
    }

    public (bool, string, string) CheakIsLock(int level, int mainBattleID)
    {
        string strtxttip = "";
        string strtip = "";
        if (level > RoleM.Instance.Level)
        {
            strtxttip = level + "级解锁";
            strtip = level + "级解锁";
            return (true, strtxttip, strtip);
        }


        if (mainBattleID > RoleM.Instance.MainBattleId)
        {
            var battleInfo = CFG_MainBattle.GetById(mainBattleID);
            if (battleInfo == null) return (false, "", "");//没有找到配置信息，就当满足开启条件

            strtxttip = "主线解锁";
            strtip = $"通关{battleInfo.name}解锁";

            return (true, strtxttip, strtip);
        }

        return (false, "", "");
    }
    public override void OnLayerChange()
    {
        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            GameUtil.Instance.SetEffectLayer(GetView().teamItems[i]._obj_on.transform, _canvasOrder + 1);
        }
    }

    private int _TeamIdx = -1;
    private void OnTeamItemClick(int idx)
    {
        _TeamIdx = idx;
        _SelectCFGId = 0;

        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            GetView().teamItems[i]._obj_on.SetActive(i == idx);
        }
        OnSelectTitles();

    }

    //槽位解锁
    private void SetPartnerLock()
    {
        for (int soltIdx = 1; soltIdx < GetView().teamItems.Length; soltIdx++)
        {
            bool isLock = CFG_XMLConfig.CheckPartnerSoltIsOpenServer(soltIdx, false);

            (bool isOtherLock, string strtxttip, string strtip) = CheakIsLock(CFG_XMLConfig.GetOpenSoltLv(soltIdx), CFG_XMLConfig.GetOpenSolMainStaget(soltIdx));
            if (!isLock)
                isLock = isOtherLock;

            GetView().teamItems[soltIdx]._obj_lock.SetActive(isLock);
            if (isLock)
            {
                if (CFG_XMLConfig.CheckPartnerSoltIsOpenServer(soltIdx, false))
                    GetView().teamItems[soltIdx]._txt_Lock.text = $"{CFG_XMLConfig.RemainingDaysPartnerSlotUnlockOpen(soltIdx)}天后解锁";
                else
                    GetView().teamItems[soltIdx]._txt_Lock.text = strtip.Replace("通关主线关卡", "通关主线\n");
            }
        }
    }

    //设置基础信息
    private void SetBaseRole(bool isInit = false)
    {
        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            if (i == 0)
            {
                if (isInit)
                {
                    GameUtil.Instance.SetJobIcon(RoleM.Instance.Job, GetView().teamItems[i]._img_job);
                    GameUtil.Instance.SetRoleHead(RoleM.Instance.CurrentAvatarId, RoleM.Instance.Job, GetView().teamItems[i]._img_head, GetView().teamItems[i]._img_head_quality);
                    GameUtil.Instance.SetRoleHeadFrame(RoleM.Instance.CurrentFrameId, RoleM.Instance.Job, GetView().teamItems[i]._img_frame);
                }

                GetView().teamItems[i]._txt_lv.text = RoleM.Instance.Level.ToString();
            }
            else
            {
                var solt = PartnerM.Instance.GetOnePartnerSlotInfo(i);
                GetView().teamItems[i]._img_head.gameObject.SetActive(solt.PartnerUid != 0);
                GetView().teamItems[i]._img_job.gameObject.SetActive(solt.PartnerUid != 0);
                GetView().teamItems[i]._obj_Lv.gameObject.SetActive(solt.PartnerUid != 0);
                GetView().teamItems[i].txt_Star.gameObject.SetActive(solt.PartnerUid != 0);

                int quality = 0;
                if (solt.PartnerUid != 0)
                {
                    var partner = PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid);
                    GetView().teamItems[i]._txt_lv.text = partner.Level.ToString();
                    GetView().teamItems[i].txt_Star.text = partner.Star.ToString();

                    quality = CFG_Partner.GetById(partner.CfgId).Quality;

                    if (isInit)
                    {
                        GameUtil.Instance.SetJobIcon(CFG_Partner.GetById(partner.CfgId).job % 10, GetView().teamItems[i]._img_job);
                        GameUtil.Instance.SetItemIcon(CFG_Partner.GetById(partner.CfgId).ItemID, GetView().teamItems[i]._img_head);
                    }
                }

                GameUtil.Instance.SetPartnetQualityFrame(quality, GetView().teamItems[i]._img_frame);
            }

        }
        if (isInit)
        {
            OnTeamItemClick(_TeamIdx);
        }
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        OnUniEventAddOrMove(true);
        if (args.Length > 0)
            _TeamIdx = (int)args[0];
        else
            _TeamIdx = 0;

            SetPartnerLock();
        SetBaseRole(true);
        RedBind();
    }

    private void RedBind()
    {
        //RedPointManager.Instance.RPS.Register(RedPointConst.Role_Title_UpStar.ToString(), (node) =>
        //{
        //    GetView().btn_UpStar.transform.Find("red").gameObject.SetActive(node.pointCount > 0);
        //});
    }

    private void ClearRedBind()
    {
        //RedPointManager.Instance.RPS.Register(RedPointConst.Role_Title_UpStar.ToString(), null);
    }

    private List<CFGItem_Title> cFGItem_Titles;
    private void OnSelectTitles(bool isRefresh = false)
    {
        cFGItem_Titles = CFG_Title.GetActiveListByType(_TeamIdx + 1);

        cFGItem_Titles.Sort((a, b) => {
            bool aExists = TitleM.Instance.CheckTitleIsEquip(a.UID, _TeamIdx);
            bool bExists = TitleM.Instance.CheckTitleIsEquip(b.UID, _TeamIdx);

            if (aExists != bExists)
                return aExists ? -1 : 1;
            else
            {
                aExists = TitleM.Instance.CheckTitleIsGet(a.UID);
                bExists = TitleM.Instance.CheckTitleIsGet(b.UID);
                if (aExists != bExists)
                    return aExists ? -1 : 1;
            }
            return a.UID.CompareTo(b.UID);
        });

        GetView().loopView_Title.SetListItemCount(cFGItem_Titles.Count);
        GetView().loopView_Title.RefreshAllShownItem();

        if (isRefresh)
        {
            int moveIdx = 0;
            for (int i = 0; i < cFGItem_Titles.Count; i++)
            {
                if (cFGItem_Titles[i].UID == _SelectCFGId)
                {
                    moveIdx = i;
                    break;
                }
            }
            GetView().loopView_Title.MovePanelToItemIndex(moveIdx, 0);
        }
        OnTitleClick(_SelectCFGId == 0 ? cFGItem_Titles[0].UID : _SelectCFGId);
        //GameHelper.Instance.SetDeayFrame(() => { 
        //    OnTitleClick(_SelectCFGId == 0 ? cFGItem_Titles[0].UID : _SelectCFGId);
        //});
    }

    private Dictionary<Transform, int> _LoopItemObjDic = new Dictionary<Transform, int>();

    private int _SelectCFGId = 0;
    private LoopListViewItem2 LeftViewFill(LoopListView2 listView, int idx)
    {
        if (cFGItem_Titles.Count <= idx || idx < 0)
            return null;

        var itemLoop = listView.NewListViewItem("Item");
        var item = itemLoop.transform;
        item.name = "item" + idx;

        var _CFGInfo = cFGItem_Titles[idx];

        _LoopItemObjDic[item] = _CFGInfo.UID;

        item.Find("on").gameObject.SetActive(_CFGInfo.UID == _SelectCFGId);


        ResouseHelper.Instance.LoadSpriteByAtlas(_CFGInfo.Path, _CFGInfo.Icon.ToString(), JUIExtensions.img(item.Find("img_Icon")));
        OnTitleItemUpdate(item, _CFGInfo);
        var btn = JUIExtensions.btn(item);
        btn.onClick.RemoveAllListeners();
        btn.onClick.AddListener(() =>
        {
            OnTitleClick(_CFGInfo.UID);
        });

       
        return itemLoop;
    }

    //单个称号显示刷新
    private void OnTitleItemUpdate(Transform item, CFGItem_Title _CFGInfo)
    {
        bool isGet = TitleM.Instance.CheckTitleIsGet(_CFGInfo.UID);
        item.Find("img_dress").gameObject.SetActive(TitleM.Instance.CheckTitleIsEquip(_CFGInfo.UID, _TeamIdx));
        JUIExtensions.img(item.Find("img_Icon")).SetGray(!isGet);
        item.Find("obj_Star").gameObject.SetActive(isGet);
        JUIExtensions.txt(item.Find("obj_Star/txt")).text = TitleM.Instance.GetTitleStar(_CFGInfo.UID).ToString();

        bool isNewGet = false;
        for (int i = 0; i < _CFGInfo.Activity.Length; i++)
        {
            var Iteminfo = _CFGInfo.Activity[i].values[0];
            if (BagM.Instance.IsGetTitlePiceNew(Iteminfo.key))
            {
                isNewGet = true;
            }
        }
        item.Find("new").gameObject.SetActive(isNewGet);
        item.Find("red").gameObject.SetActive(RedDotM.Instance.GetOneTitleAtiveRedNum(_CFGInfo.UID, _CFGInfo.Activity) > 0);
        item.Find("img_time").gameObject.SetActive(isGet);
        if (isGet)
            JUIExtensions.txt(item.Find("img_time/txt_Time")).text = $"有效期:{(TitleM.Instance.GetTitleDay(_CFGInfo.UID) < 0 ? "永久" : TitleM.Instance.GetTitleDay(_CFGInfo.UID) + "天")}";
    }

    private CFGItem_Title CFGInfo;
    private void OnTitleClick(int cfgId)
    {
        _SelectCFGId = cfgId;
        CFGInfo = CFG_Title.GetById(_SelectCFGId);

        for (int i = 0; i < CFGInfo.Activity.Length; i++)
        {
            BagM.Instance.RemoveTitlePiceNew(CFGInfo.Activity[i].values[0].key);
        }
        foreach (var _item in _LoopItemObjDic)
        {
            _item.Key.Find("on").gameObject.SetActive(_item.Value == cfgId);
            if (_item.Value == cfgId)
            {
                OnTitleItemUpdate(_item.Key, CFGInfo);
            }
        }
        OnOnTitleInfo();
        OnTitlePiceShow();
        OnTitleStarShow();
    }


    //称号基础信息
    private void OnOnTitleInfo(bool isInit = true)
    {
        bool isGet = TitleM.Instance.CheckTitleIsGet(_SelectCFGId);
        GetView().btn_TitleActive.gameObject.SetActive(!isGet);

        bool ieEquip = TitleM.Instance.CheckTitleIsEquip(_SelectCFGId, _TeamIdx);
        GetView().btn_TitleEquip.gameObject.SetActive(!ieEquip && isGet);
        GetView().btn_TitleEquipDel.gameObject.SetActive(ieEquip && isGet);

        if (isInit)
        {
            string Str = string.Empty;
            for (int i = 0; i < CFGInfo.Limit.Length; i++)
            {
                Str += CFG_XMLConfig.GetAttributeStrByID(CFGInfo.Limit[i].key, CFGInfo.Limit[i].value) + "\n";
            }
            GetView().txt_DressArr.text = Str;

            Str = string.Empty;

            var baseGoodItems = CFGInfo.GetGoodItemsByStar(0);
            for (int i = 0; i < baseGoodItems.Length; i++)
            {
                Str += CFG_XMLConfig.GetAttributeStrByID(baseGoodItems[i].key, baseGoodItems[i].value) + "\n";
            }
            GetView().txt_BaseArr.text = Str;

            ResouseHelper.Instance.LoadSpriteByAtlas(CFGInfo.Path, CFGInfo.Icon.ToString(), GetView().img_TitleName);
        }

        GetView().txt_DressArr.SetGray(!isGet);
        GetView().txt_BaseArr.SetGray(!isGet);
    }

    /// <summary>
    /// 称号升星符ID
    /// </summary>
    private int _TitleUpStarID = 0;
    //称号星级展示
    private void OnTitleStarShow()
    {
        bool isGet = TitleM.Instance.CheckTitleIsGet(_SelectCFGId);
        GetView().root_UpStar.gameObject.SetActive(isGet);
        GetView().txt_NoTip.gameObject.SetActive(!isGet);

        int nowStar = TitleM.Instance.GetTitleStar(_SelectCFGId);
        for (int i = 0; i < GetView().root_Stars.childCount; i++)
        {
            JUIExtensions.img(GetView().root_Stars.GetChild(i)).SetGray(i >= nowStar);
        }

        if (isGet)
        {
            bool IsMax = nowStar == CFG_XMLConfig.xmlConfig.WarFight.maxStar;
            GetView().btn_UpStar.gameObject.SetActive(!IsMax);

            GetView().btn_UpStar.transform.Find("red").gameObject.SetActive(RedDotM.Instance.GetOneTitleAtiveRedNum(_SelectCFGId, CFGInfo.Activity,true) > 0);
            
            var baseGoodItems = CFGInfo.GetGoodItemsByStar(nowStar);
            GetView().txt_nowStarLv.text = nowStar.ToString();
            string Str = string.Empty;
            for (int i = 0; i < baseGoodItems.Length; i++)
            {
                Str += CFG_XMLConfig.GetAttributeStrByID(baseGoodItems[i].key, baseGoodItems[i].value) + "\n";
            }

            GetView().txt_NowStar.text = Str;

            int nextStar = IsMax ? nowStar : nowStar + 1;
            baseGoodItems = CFGInfo.GetGoodItemsByStar(nextStar);
            GetView().txt_nextStarLv.text = nextStar.ToString();
            Str = string.Empty;
            for (int i = 0; i < baseGoodItems.Length; i++)
            {
                Str += CFG_XMLConfig.GetAttributeStrByID(baseGoodItems[i].key, baseGoodItems[i].value) + "\n";
            }
            GetView().txt_NextStar.text = Str;
            OnUpStarPerShow();
        }
    }

    //幸运符选择
    public void OnLuckItemSelect( int cid)
    {
        _TitleUpStarID = cid;
        OnUpStarPerShow();
    }

    //升星概率
    private void OnUpStarPerShow()
    {
        if (_TitleUpStarID != 0 && BagM.Instance.GetNumByCfgID(_TitleUpStarID) == 0)
            _TitleUpStarID = 0;

        GetView().root_UpStarItem.Find("icon").gameObject.SetActive(_TitleUpStarID != 0);
        if (_TitleUpStarID != 0)
        {
            GameUtil.Instance.SetItemIcon(_TitleUpStarID, JUIExtensions.img(GetView().root_UpStarItem.Find("icon")), JUIExtensions.img(GetView().root_UpStarItem));
        }
        else
        {
            GameUtil.Instance.SetItemBg(0, JUIExtensions.img(GetView().root_UpStarItem));
        }
        JUIExtensions.txt(GetView().root_UpStarItem.Find("txt_lv")).text = _TitleUpStarID != 0 ? BagM.Instance.GetNumByCfgID(_TitleUpStarID).ToString() : string.Empty;

        bool IsMax = TitleM.Instance.GetTitleStar(_SelectCFGId) == CFG_XMLConfig.xmlConfig.WarFight.maxStar;
        int addRate = 0;
        if (!IsMax)
        {
            if (_TitleUpStarID != 0)
                addRate = +CFG_XMLConfig.xmlConfig.WarFight.luckTicketDic[_TitleUpStarID];
            addRate += CFGInfo.StarRate[TitleM.Instance.GetTitleStar(_SelectCFGId)].value;
        }
        GetView().txt_UpPro.text = IsMax ? "满星": "升星概率:" + GameUtil.Instance.GetPercentStr(addRate / 10000f);
    }

    private void OnTitlePiceShow(bool isInit = true)
    {
        //碎片展示
        for (int i = 0; i < CFGInfo.Activity.Length; i++)
        {
            var root = GetView().root_Pice.GetChild(CFGInfo.Activity[i].id - 1);
            int ItemID = CFGInfo.Activity[i].values[0].key;
            JUIExtensions.txt(root.Find("txt_lv")).text = BagM.Instance.GetNumByCfgID(ItemID).ToString();

            if (isInit)
            {
                GameUtil.Instance.SetItemIcon(ItemID, JUIExtensions.img(root.Find("icon")));
                var btn = JUIExtensions.btn(root);
                btn.onClick.RemoveAllListeners();
                btn.onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = ItemID });
                });
            }
        }
    }

    [Subscriber(ThreadMode.Main)]
    public void OnRefresh_TitleList(GameEvent.Refresh_TitleList e)
    {
        OnSelectTitles(true);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Title_Active);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Title_UpStar);
    }

    [Subscriber(ThreadMode.Main)]
    private void OnTitle_UpStar(GameEvent.Title_UpStar e)
    {
        //OnTitlePiceShow();
        //OnTitleStarShow();
        OnTitleClick(_SelectCFGId);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Title_UpStar);
    }

    private void OnSync_CheckOpenTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        TcpC.C2S_TitleList();
    }


    private void OnUniEventAddOrMove(bool isAdd)
    {
        if (isAdd)
        {
            UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        }
        else
        {
            UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        }
    }

    public override void OnExit()
    {
        OnUniEventAddOrMove(false);
        JEvent.defaultEvent.Unregister(this);
        ClearRedBind();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}