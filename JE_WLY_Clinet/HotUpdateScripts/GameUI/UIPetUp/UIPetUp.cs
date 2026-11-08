using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using HotUpdateScripts.ProtoM;
using pb;
using JEngine.Helper;
using JEngine.Event;
using JEngine.Core;
using System;

public class UIPetUp : UIBase<UIModelPetUp, UIViewPetUp>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btnClose.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPetUp>();
        });

        GetView().btnUpLv.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            long needExp = (long)(CFG_Role.GetByJobAndLevel(21, PetSlotInfo.Level).NeedEXP * CFG_XMLConfig.xmlConfig.PetlevelDic[CFG_Pet.GetById(PetSlotInfo.CfgId).Quality]) - PetSlotInfo.Exp;
            if (GameUtil.Instance.ComCoinToCheck((int)CFG_Item.ItemID.PetEXP, needExp, null, false))
            {
                TcpC.C2S_PetLevel();
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "角色",
                    component = "宠物",
                    position = "养成",
                    event_click = "宠物升级",
                });
            }
            else
                UITool.ShowTips($"升至下一级还需{needExp - RoleM.Instance.PetUsedExp}宠物经验");

        });

        for (int i = 0; i < GetView().root_Btns.childCount; i++)
        {
            int idx = i;
            GetView().root_Btns.GetChild(i).GetComponent<Button>().onClick.AddListener(() => {
                OnSelectBtnClick(idx);
            } );
        }

        for (int i = 0; i < GetView().root_BreakTouch.childCount; i++)
        {
            int idx = i;
            var _touch = GetView().root_BreakTouch.GetChild(i).GetComponent<TouchEventCom>();
            _touch.AddClickListener((a) =>
            {
                if (isAutoCkick)
                {
                    isAutoCkick = false;
                    return;
                }
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnBreakBtnClick(idx);
            });

            _touch.AddPointerDownListener((a) =>
            {
                OnAutoBreak(idx);
            });
           
            _touch.AddPointerUpListener((a) =>
            {
                OnClearAutoBreak();
            });
        }
    }

    public override void OnLayerChange()
    {
        GameUtil.Instance.SetEffectLayer(GetView().obj_BreakPro.transform, _canvasOrder + 1);
    }

    private void OnBindRed()
    {
        RedPointManager.Instance.RPS.Register(RedPointConst.Pet_UpStar_1.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().obj_RedUpStar.SetActive(node.pointCount > 0);
        });
    }

    private int SelectIdx = 0;
    private void OnSelectBtnClick(int idx)
    {
        SelectIdx = idx;
        for (int i = 0; i < GetView().root_Btns.childCount; i++)
        {
            GetView().root_Btns.GetChild(i).Find("on").gameObject.SetActive(SelectIdx == i);
            if (i == SelectIdx)
            { 
                GetView().txt_title.text = JUIExtensions.txt(GetView().root_Btns.GetChild(i).Find("txt")).text;
            }
        }
        if (SelectIdx == 0)
        {
            OnPetShow(false);
            OnPetLvShow();
        }
        else
        {
            SetBreakUsed(false);
            OnBreakShow();
        }

        GetView().UI_Break.SetActive(SelectIdx == 1);
        GetView().UI_UpLV.SetActive(SelectIdx == 0);
    }

    private void OnPetShow(bool isInit = true)
    {
        if (isInit)
        {
            GameUtil.Instance.SetItemIcon(CFG_Pet.GetById(PetSlotInfo.CfgId).ItemID, GetView().img_PetIcon, GetView().img_PetBg,null, CFG_Pet.GetById(PetSlotInfo.CfgId).Quality);
            GetView().txt_PetName.text = GameUtil.Instance.GetStrByQuality(CFG_Pet.GetById(PetSlotInfo.CfgId).Quality, CFG_Pet.GetById(PetSlotInfo.CfgId).Name);
        }
        GetView().txt_PetLv.text = PetSlotInfo.Level.ToString();
    }

    #region 升级

   
    private void OnPetLvShow()
    {
        int nowLv = PetSlotInfo.Level;
        var cfgLv = CFG_Role.GetByJobAndLevel(21, nowLv);
        GetView().txt_nowLv.text = nowLv.ToString();

        float addValue = CFG_Pet.GetById(PetSlotInfo.CfgId).Mark == 0 ? 1 : CFG_XMLConfig.xmlConfig.PetQualityDic[CFG_Pet.GetById(PetSlotInfo.CfgId).Quality];

        string Des = "";
        for (int i = 0; i < cfgLv.Property.Length; i++)
        {
            Des += (CFG_XMLConfig.GetAttributeStrByID(cfgLv.Property[i].Id, (int)(cfgLv.Property[i].Num * addValue)) + "\n");
        }
        GetView().txt_nowLvDes.text = Des;

        int nextLv = PetSlotInfo.Level == CFG_Role.MaxLvByJob(21) ? PetSlotInfo.Level : (PetSlotInfo.Level + 1);
        cfgLv = CFG_Role.GetByJobAndLevel(21, nextLv);
        GetView().txt_nextLv.text = nextLv.ToString();

        Des = "";
        for (int i = 0; i < cfgLv.Property.Length; i++)
        {
            Des += (CFG_XMLConfig.GetAttributeStrByID(cfgLv.Property[i].Id, (int)(cfgLv.Property[i].Num * addValue)) + "\n");
        }
        GetView().txt_nextLvDes.text = Des;
        GetView().btnUpLv.gameObject.SetActive(nowLv != nextLv);

        OnPetExpShow(true);
    }

    private void OnPetExpShow(bool isInit = false)
    {
        if (SelectIdx == 0  && PetSlotInfo.Level != CFG_Role.MaxLvByJob(21))
        {
            if (isInit)
                GameUtil.Instance.SetItemIcon((int)CFG_Item.ItemID.PetEXP, GetView().img_lvUsed);

            long needExp = (long)(CFG_Role.GetByJobAndLevel(21, PetSlotInfo.Level).NeedEXP * CFG_XMLConfig.xmlConfig.PetlevelDic[CFG_Pet.GetById(PetSlotInfo.CfgId).Quality]) - PetSlotInfo.Exp;
            GetView().txt_lvUsed.text = $"{GameUtil.Instance.GetCheakNumStr((int)CFG_Item.ItemID.PetEXP, needExp)}";
            //GetView().txt_lvUsed.text = $"{GameUtil.Instance.GetCheakNumStr((int)CFG_Item.ItemID.PetEXP, (long)needExp)}/{GameUtil.Instance.FormatNumWithUnit(RoleM.Instance.PetUsedExp)}"  ;
        }
    }

    #endregion

    #region 突破
    private TimeItem AutoCkick;
    private bool isAutoCkick = false;
    private void OnClearAutoBreak()
    {
        if (AutoCkick != null)
        {
            GameHelper.Instance.RemoveInvoke(AutoCkick);
            AutoCkick = null;
        }
    }

    private void OnAutoBreak(int idx)
    {
        OnClearAutoBreak();
        AutoCkick = GameHelper.Instance.AddInvoke(new TimeItem(0, 0.4f, () => {
            isAutoCkick = true;
            OnBreakBtnClick(idx);
        }, true, true));
    }

    private void OnBreakBtnClick(int idx)
    {
        if (BagM.Instance.GetNumByCfgID(CFG_XMLConfig.xmlConfig.PetAdvanceUseds[idx].id) < 0)
        {
            UITool.ShowTips(CFG_Item.GetById(CFG_XMLConfig.xmlConfig.PetAdvanceUseds[idx].id).Name + "不足");
            OnClearAutoBreak();
            return;
        }
        TcpC.C2S_PetStar(CFG_XMLConfig.xmlConfig.PetAdvanceUseds[idx].id);

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "角色",
            component = "宠物",
            position = "养成",
            event_click = "宠物突破",
        });
    }

    private void SetBreakUsed(bool isInit = false)
    {
        for (int i = 0; i < GetView().root_BreakTouch.childCount; i++) { 
            var go = GetView().root_BreakTouch.GetChild(i);
            if (isInit)
                GameUtil.Instance.SetItemIcon(CFG_XMLConfig.xmlConfig.PetAdvanceUseds[i].id, JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
            JUIExtensions.txt(go, "txt_lv").text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID(CFG_XMLConfig.xmlConfig.PetAdvanceUseds[i].id));
        }
    }

    private void OnBreakShow()
    {
        GetView().txt_nowBreak.text = PetSlotInfo.Star + "阶";
        var cfgInfo = CFG_Pet.GetById(PetSlotInfo.CfgId);
        var Arrt = cfgInfo.GetStarAttributeByLv(PetSlotInfo.Star);

        string Des = "";
        for (int i = 0; i < Arrt.Length; i++)
        {
            Des += (CFG_XMLConfig.GetAttributeStrByID(Arrt[i].key, Arrt[i].value) + "\n");
        }
        //var offlineAwards = cfgInfo.GetOfflineAwardByStar(PetSlotInfo.Star);
        //if (offlineAwards != null)
        //{
        //    for (int i = 0; i < offlineAwards.Length; i++)
        //    {
        //        Des += $"{CFG_Item.GetById(offlineAwards[i].id).Name}:{GameUtil.Instance.FormatNumWithUnit(offlineAwards[i].num)}/{offlineAwards[i].gap / 60}分钟\n";
        //    }
        //}
        GetView().txt_nowBreakDes.text = Des;

        int nextBreak = PetSlotInfo.Star == cfgInfo.StarConsume.Length ? PetSlotInfo.Star : (PetSlotInfo.Star + 1);
        GetView().txt_nextBreak.text = nextBreak + "阶";
        Arrt = cfgInfo.GetStarAttributeByLv(nextBreak);
        Des = "";
        for (int i = 0; i < Arrt.Length; i++)
        {
            Des += (CFG_XMLConfig.GetAttributeStrByID(Arrt[i].key, Arrt[i].value) + "\n");
        }
        //offlineAwards = cfgInfo.GetOfflineAwardByStar(nextBreak);
        //if (offlineAwards != null)
        //{
        //    for (int i = 0; i < offlineAwards.Length; i++)
        //    {
        //        Des += $"{CFG_Item.GetById(offlineAwards[i].id).Name} : {GameUtil.Instance.FormatNumWithUnit(offlineAwards[i].num)}/{offlineAwards[i].gap / 60}分钟\n";
        //    }
        //}
        GetView().txt_nextBreakDes.text = Des;

        GetView().txt_BreakLimitDes.text = "突破限制:当前宠物达到" + cfgInfo.GetStatLimet(nextBreak) + "级";

        GetView().txt_BreakLimitDes.gameObject.SetActive(PetSlotInfo.Star != nextBreak);
        GetView().obj_BreakPro.SetActive(PetSlotInfo.Star != nextBreak);
        GetView().obj_BreakMax.SetActive(PetSlotInfo.Star == nextBreak);

        int needNum = cfgInfo.GetStarNeedAllPro(nextBreak);
        GetView().txt_BreakPro.text = GameUtil.Instance.FormatNumWithUnit(PetSlotInfo.ExpStar)  + "/" + GameUtil.Instance.FormatNumWithUnit(needNum);
        GetView().img_BreakPro.fillAmount = PetSlotInfo.ExpStar / (float)needNum;

    }
    #endregion

    private PetSlotInfo PetSlotInfo { get {
            return PetM.Instance.GetOnePetByIdx(0);
        }
    }
    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(CommonResponse);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);

        OnBindRed();
        OnPetShow();
        OnSelectBtnClick(0);
        SetBreakUsed(true);
    }
    //[Subscriber(ThreadMode.Main)]
    private void OnSync_UserBase(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;
        OnPetExpShow(false);
    }

    //[Subscriber(ThreadMode.Main)]
    private void CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.cid == CmdId.C2S_PetLevel && _e.code == 0)
        {
            UITool.ShowTips("升级成功");
            OnPetShow(false);
            OnPetLvShow();
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void Refresh_PetStar(GameEvent.Refresh_PetStar _e)
    {
        if (isNullGo) return;

        if (_e.Code == 0)
        {
            UITool.ShowTips(_e.ResPer + "倍突破");
            OnBreakShow();
            SetBreakUsed(false);
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_PetGotExp);
        }
        else
            OnClearAutoBreak();
    }

    public override void OnExit()
    {
        OnClearAutoBreak();
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(CommonResponse);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);
    }

    public override void OnDestroy()
    {
        OnClearAutoBreak();
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(CommonResponse);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);
    }
}