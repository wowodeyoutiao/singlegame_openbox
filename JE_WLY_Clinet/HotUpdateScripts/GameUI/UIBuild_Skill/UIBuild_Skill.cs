using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Event;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;
using JEngine.UI;
using JEngine.MyNet;
using pb;

public class UIBuild_Skill : UIBase<UIModelBuild_Skill, UIViewBuild_Skill>
{
    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }
    public Camera mainCam;
    public override void OnInit()
    {
        mainCam = GameObject.Find("UICamera").GetComponent<Camera>();
        

        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            int idx = 10;
            if (_SelectSkillClass == UIBuild.SelectTwo.Skill_Advance)
                idx = 11;
            else if (_SelectSkillClass == UIBuild.SelectTwo.Skill_ChangeSkill)
                idx = 12;

            UITool.Open<UIComRule>(null, idx);
        });


        GetView().btn_topRes1Add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //UITool.Close<UIBuild>();
            UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
        });

        GetView().btn_topRes2Add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Shop_SL))
                return;
            //UITool.Close<UIBuild>();
            UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
        });

        GetView().btn_topRes3Add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //UITool.Close<UIBuild>();
            UITool.Open<UISpecialOffer>(null, ActivityType.GiftPay);
        });

        GetView().btn_topRes4Add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            //UITool.Close<UIBuild>();
            UITool.Open<UIShop>(null, CFG_Store.StoreType.RMB);
        });

        GetView().btn_UpLv.onClick.AddListener(OnUpLvBtnClick);

        GetView().btn_Advance.onClick.AddListener(OnAdvanceBtnClick);

        //技能配置按钮
        //for (int i = 0; i < 13; i++)
        //{
        //    GetView().root_SkillSolt.GetChild(i).GetComponent<Button>().onClick.AddListener(()=> { 
        //    });
        //}


    }

    private void OnBindRed()
    {
    }

    private void ClearRed()
    {
    }

    //按钮开启
    private void OnBtnLock()
    {
      
    }

    private UIBuild.SelectTwo _SelectSkillClass;

    //伙伴ID
    private long PartnerId;
    private int TeamIdx;

    //面板选择
    public void OnSelcetUIBtnClick(UIBuild.SelectTwo eClass,int teamIdx)
    {
        TeamIdx = teamIdx;
        PartnerId = teamIdx == 0 ? 0 : PartnerM.Instance.GetOnePartnerSlotInfo(teamIdx).PartnerUid;
        _SelectSkillClass = eClass;
        OnAllSkillShow();
        GetView().obj_Res2.SetActive(_SelectSkillClass == UIBuild.SelectTwo.Skill_UpLv);
        GetView().obj_Res3.SetActive(_SelectSkillClass == UIBuild.SelectTwo.Skill_Advance);
    }

    private int SelectId;

    //技能点击选择
    private void OnSkillSelectClick(int idx,GameObject go)
    {
        for (int i = 0; i < GetView().SkillItems.Length; i++)
        {
            GetView().SkillItems[i].obj_Select.SetActive(GetView().SkillItems[i].obj_Select == go);
        }
        SelectId = idx;

        switch (_SelectSkillClass)
        {
            case UIBuild.SelectTwo.Skill_UpLv:
                OnUpLvUI();
                break;
            case UIBuild.SelectTwo.Skill_Advance:
                OnAdvanceUI();
                break;
            case UIBuild.SelectTwo.Skill_ChangeSkill:

               
                break;
            default:
                break;
        }
    }

    //升阶
    private void OnAdvanceUI()
    {
        var skillInfo = PartnerId == 0 ? SkillM.Instance.GetSkillInfoBySkillID(SelectId) :PartnerM.Instance.GetOneSkillInfo(PartnerId, SelectId);
        var cfgInfo = CFG_Skill.GetById(skillInfo.CfgId);

        GameUtil.Instance.SetItemIcon((int)cfgInfo.SkillIcon, GetView().img_ItemAdvance);
        GameUtil.Instance.SetSkillBg((int)cfgInfo.UUID, GetView().img_ItemBgAdvance);
        GetView().txt_advSkillName.text = cfgInfo.Name;
        GetView().txt_advance_Lv.text = cfgInfo.SkillLevel.ToString();
        GetView().txt_advance_advNum.text = skillInfo.Advanced + "阶";

        var nowAdvCfgInfo = CFG_Skill.GetAdvanceByUidAndAdvan(skillInfo.CfgId, skillInfo.Advanced);
        GetView().txt_AdvanceNowLv.text = skillInfo.Advanced + "阶";
        GetView().txt_AdvanceNowDes.text = skillInfo.Advanced == 0 ? "无" :nowAdvCfgInfo.UpperlevelDes;

        var nextAdvCfgInfo = CFG_Skill.GetAdvanceByUidAndAdvan(skillInfo.CfgId, skillInfo.Advanced + 1);
        nextAdvCfgInfo = nextAdvCfgInfo == null ? nowAdvCfgInfo : nextAdvCfgInfo;

        GameUtil.Instance.SetItemIcon(nextAdvCfgInfo.UpperlevelConsume[0].id, GetView().img_AdvanCoin1Num);
        GetView().txt_AdvanCoin1Num.text = GameUtil.Instance.GetCheakNumStr(nextAdvCfgInfo.UpperlevelConsume[0].id, nextAdvCfgInfo.UpperlevelConsume[0].num);
        GameUtil.Instance.SetItemIcon(nextAdvCfgInfo.UpperlevelConsume[1].id, GetView().img_AdvanCoin2Num);

        GetView().txt_AdvanCoin2Num.text = GameUtil.Instance.GetCheakNumStr(nextAdvCfgInfo.UpperlevelConsume[1].id, nextAdvCfgInfo.UpperlevelConsume[1].num);

        GetView().txt_AdvanceNextLv.text = nextAdvCfgInfo.CanupperLevel + "阶";
        GetView().txt_AdvanceNextDes.text = nextAdvCfgInfo.UpperlevelDes;

        bool isMax = skillInfo.Advanced == nextAdvCfgInfo.CanupperLevel;
        GetView().btn_Advance.gameObject.SetActive(!isMax);
        GetView().txt_AdvanTip.text = !isMax ? $"需技能{nextAdvCfgInfo.SkillLevel}级" : "已满阶";
    }

    private void OnAdvanceBtnClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        var skillInfo = PartnerId == 0 ? SkillM.Instance.GetSkillInfoBySkillID(SelectId) : PartnerM.Instance.GetOneSkillInfo(PartnerId, SelectId);
        var cfgInfo = CFG_Skill.GetById(skillInfo.CfgId);
        var nextInfo = CFG_Skill.GetAdvanceByUidAndAdvan(skillInfo.CfgId, skillInfo.Advanced + 1);

        if (cfgInfo.SkillLevel < nextInfo.SkillLevel)
        {
            UITool.ShowTips("技能等级不足");
            return;
        }

        for (int i = 0; i < nextInfo.UpperlevelConsume.Length; i++)
        {

            if (!GameUtil.Instance.ComCoinToCheck(nextInfo.UpperlevelConsume[i].id, nextInfo.UpperlevelConsume[i].num, () => { UITool.Close<UIBuild>(); }))
            {
                return;
            }
        }
        TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_SkillAdvanced, MType.Request, TCPClientMgr.ProtoSerialize(new SkillAdvancedRequest { 
            SkillId = skillInfo.CfgId,
            PartnerUid = PartnerId
        }));

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "打造",
            component = "技能",
            position = "进阶界面",
            event_click = "进阶",
        });
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_SkillAdvanced(GameEvent.Sync_SkillAdvanced _e)
    {
        OnRefreshAllSkill();
        OnAdvanceUI();

        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_SkillAdvancedSuc);
    }


    #region 升级
    //升级界面
    public void OnUpLvUI()
    {
        var skillInfo = PartnerId == 0 ? SkillM.Instance.GetSkillInfoBySkillID(SelectId) : PartnerM.Instance.GetOneSkillInfo(PartnerId, SelectId);
        var cfgInfo = CFG_Skill.GetById(skillInfo.CfgId);

        GameUtil.Instance.SetItemIcon((int)cfgInfo.SkillIcon, GetView().img_upLv_Item);
        GameUtil.Instance.SetSkillBg(cfgInfo.UUID, GetView().img_upLv_ItemBg);
        GetView().txt_upLv_skillName.text = cfgInfo.Name;
        GetView().txt_upLv_Lv.text = cfgInfo.SkillLevel.ToString();
        GetView().txt_upLv_advanLv.text = skillInfo.Advanced + "阶";

        GetView().txt_upLv_NowLv.text = cfgInfo.SkillLevel.ToString();
        GetView().txt_upLv_NowDes.text = cfgInfo.SkillDes;

        GameUtil.Instance.SetItemIcon(cfgInfo.PromoteConsume[0].id, GetView().img_upLv_Coin1Num);
        GetView().txt_upLv_Coin1Num.text = GameUtil.Instance.GetCheakNumStr(cfgInfo.PromoteConsume[0].id, cfgInfo.PromoteConsume[0].num);
        GameUtil.Instance.SetItemIcon(cfgInfo.PromoteConsume[1].id, GetView().img_upLv_Coin2Num);
        GetView().txt_upLv_Coin2Num.text = GameUtil.Instance.GetCheakNumStr(cfgInfo.PromoteConsume[1].id,cfgInfo.PromoteConsume[1].num);

        var nextInfo = CFG_Skill.GetNextLvByUid(skillInfo.CfgId);
        GetView().txt_upLv_NextLv.text = nextInfo.SkillLevel.ToString();
        GetView().txt_upLv_NextDes.text = nextInfo.SkillDes;

        bool isMax = cfgInfo.SkillLevel == nextInfo.SkillLevel;
        GetView().btn_UpLv.gameObject.SetActive(!isMax);
        GetView().txtj_UpLvTip.text = !isMax ? $"需{(PartnerId == 0 ? "角色" : "伙伴")}{nextInfo.RoleLevellimit}级" : "已满级";

        GameHelper.Instance.SetDeayFrame(() => {
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().rect_upLv_NowDes);
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().rect_upLv_NextDes);
        });
    }

    //升级点击
    private void OnUpLvBtnClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        var skillInfo = PartnerId == 0 ? SkillM.Instance.GetSkillInfoBySkillID(SelectId) : PartnerM.Instance.GetOneSkillInfo(PartnerId, SelectId);
        var cfgInfo = CFG_Skill.GetById(skillInfo.CfgId);
        var nextInfo = CFG_Skill.GetNextLvByUid(skillInfo.CfgId);

        if (nextInfo.RoleLevellimit > (PartnerId == 0 ? RoleM.Instance.Level : PartnerM.Instance.GetOnePartnerInfoByUid(PartnerId).Level))
        {
            UITool.ShowTips("当前"+(PartnerId == 0 ? "角色" : "伙伴") +"等级不足");
            return;
        }

        for (int i = 0; i < cfgInfo.PromoteConsume.Length; i++)
        {
            if (!GameUtil.Instance.ComCoinToCheck(cfgInfo.PromoteConsume[i].id, cfgInfo.PromoteConsume[i].num, () => { UITool.Close<UIBuild>(); }))
            {
                return;
            }
        }

        TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_SkillUpgrade, MType.Request, TCPClientMgr.ProtoSerialize(new SkillUpgradeRequest { 
            SkillId = skillInfo.CfgId ,
            PartnerUid = PartnerId
        }));

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "打造",
            component = "技能",
            position = "升级界面",
            event_click = "升级",
        });
    }
    
    [Subscriber(ThreadMode.Main)]
    //技能 升级
    private void OnSync_SkillUpgrade(GameEvent.Sync_SkillUpgrade _e)
    {
        OnRefreshAllSkill();
        OnUpLvUI();

        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_SkillUpgradeSuc);
    }

    #endregion


    //技能信息
    private void OnRefreshAllSkill()
    {
        var skillList = PartnerId == 0 ? SkillM.Instance.GetSkillInfoList() : PartnerM.Instance.GetSkillInfoByPartnerId(PartnerId);
        skillList.Sort((a, b) => { return CFG_Skill.GetById(a.CfgId).SkillID.CompareTo(CFG_Skill.GetById(b.CfgId).SkillID); });

        int idx = 0;
        for (int i = 0; i < skillList.Count; i++)
        {
            var cfgInfo = CFG_Skill.GetById(skillList[i].CfgId);
            if (cfgInfo.Skilltype != (int)CFG_Skill.SkillType.Normall && cfgInfo.ActiveorPassive == 1)
            {

                GetView().SkillItems[idx].txt_Lv.text = cfgInfo.SkillLevel == 0 ? "" : cfgInfo.SkillLevel.ToString();

                if (_SelectSkillClass == UIBuild.SelectTwo.Skill_ChangeSkill)
                    GetView().SkillItems[idx].txt_Lv.text = cfgInfo.Name;

                GetView().SkillItems[idx].txt_Adv.text = skillList[i].Advanced == 0 ? "" : skillList[i].Advanced + "阶";

                GetView().SkillItems[idx].obj_equip.SetActive(_SelectSkillClass == UIBuild.SelectTwo.Skill_ChangeSkill && SkillM.Instance.CheakSkillEquip(cfgInfo.UUID));

                GetView().SkillItems[idx].dot.SetActive((_SelectSkillClass == UIBuild.SelectTwo.Skill_UpLv && RedDotM.Instance.GetOneSkillUpLvRedNum(cfgInfo.SkillID, TeamIdx, false) > 0)
                    || (_SelectSkillClass == UIBuild.SelectTwo.Skill_Advance && RedDotM.Instance.GetOneSkillAdvanceRedNum(cfgInfo.SkillID, TeamIdx, false) > 0));
                idx++;
            }
        }
       
    }
    
    //技能按钮状态
    private void OnAllSkillShow()
    {
        var skillList = PartnerId == 0 ? SkillM.Instance.GetSkillInfoList() : PartnerM.Instance.GetSkillInfoByPartnerId(PartnerId);
        skillList.Sort((a, b) => { return CFG_Skill.GetById(a.CfgId).SkillID.CompareTo(CFG_Skill.GetById(b.CfgId).SkillID); });

        int idx = 0;

        for (int i = 0; i < skillList.Count; i++)
        {
            var cfgInfo = CFG_Skill.GetById(skillList[i].CfgId);

            //主动和非普攻才能显示
            if (cfgInfo.Skilltype != (int)CFG_Skill.SkillType.Normall && cfgInfo.ActiveorPassive == 1)
            {
                GameUtil.Instance.SetItemIcon((int)cfgInfo.SkillIcon, GetView().SkillItems[idx].img_icon);
                GameUtil.Instance.SetSkillBg((int)cfgInfo.UUID, GetView().SkillItems[idx].img_bg);

                GameObject selectObj = GetView().SkillItems[idx].obj_Select;
                selectObj.SetActive(false);

                GetView().SkillItems[idx]._touch.RemoveAllListeners();

                if (_SelectSkillClass != UIBuild.SelectTwo.Skill_ChangeSkill)
                {
                    GetView().SkillItems[idx]._touch.AddClickListener((a) =>
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        OnSkillSelectClick(cfgInfo.SkillID, selectObj);
                    });
                    //默认选择第一个
                    if (idx == 0)
                    {
                        OnSkillSelectClick(cfgInfo.SkillID, selectObj);
                    }
                }
                else
                {
                    //放置移动
                    GetView().SkillItems[idx]._touch.AddDragBeginListener((p) => { OnDragBegin(p, cfgInfo.UUID, GetView().SkillItems[idx].img_bg.transform); });
                    GetView().SkillItems[idx]._touch.AddDragListener((p) => { OnDrag(p, cfgInfo.UUID); });
                    GetView().SkillItems[idx]._touch.AddDragEndListener((p) => { OnDragEnd(p, cfgInfo.UUID); });
                    GetView().SkillItems[idx]._touch.AddPointerUpListener((p) => { OnDragEnd(p, cfgInfo.UUID,false); });
                    if (idx == 0)
                    {
                        OnChangSkill();

                        if (RedPointManager.Instance.GetRedNum(RedPointConst.Skill_Setting.ToString()) != 0)
                        {
                            RedPointManager.Instance.SetRedClose(RedPointConst.Skill_Setting.ToString(), () => {
                                RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Setting);
                            });
                        }
                    }
                }

                GetView().SkillItems[idx]._touch.gameObject.SetActive(true);
                idx++;
            }
        }
      

        if (idx < GetView().SkillItems.Length)
        {
            for (int i = idx; i < GetView().SkillItems.Length; i++)
            {
                GetView().SkillItems[i]._touch.gameObject.SetActive(false);
            }
        }

        GetView().UpLvUI.SetActive(_SelectSkillClass == UIBuild.SelectTwo.Skill_UpLv && idx > 0);
        GetView().AdvanceUI.SetActive(_SelectSkillClass == UIBuild.SelectTwo.Skill_Advance && idx > 0);
        GetView().ChangeSkillUI.SetActive(_SelectSkillClass == UIBuild.SelectTwo.Skill_ChangeSkill && idx > 0);

        GetView().obj_NoSkill.SetActive(idx == 0);
        GetView().txt_tip.text = idx != 0 ? "": (PartnerId == 0 ? "未获得主动技能" : (CFG_Partner.GetById(PartnerM.Instance.GetOnePartnerInfoByUid(PartnerId).CfgId).Skills == null ? "橙色品质以上伙伴会带有技能" : "当前伙伴尚未获得主动技能")) ;

        OnRefreshAllSkill();
    }

    private void OnDragBegin(UnityEngine.EventSystems.PointerEventData pEvent,int cfgId,Transform baseObj)
    {
        if (SkillM.Instance.CheakSkillEquip(cfgId))
            return;
    
        GetView().obj_dragSolt.gameObject.SetActive(true);
        GetView().obj_dragSolt.position = baseObj.position;

        var cfgInfo = CFG_Skill.GetById(cfgId);
        GameUtil.Instance.SetItemIcon((int)cfgInfo.SkillIcon, JUIExtensions.img(GetView().obj_dragSolt.Find("icon")));
        GameUtil.Instance.SetSkillBg((int)cfgInfo.UUID, JUIExtensions.img(GetView().obj_dragSolt));
    }

    private void OnDrag(UnityEngine.EventSystems.PointerEventData pEvent, int cfgId)
    {
        if (SkillM.Instance.CheakSkillEquip(cfgId))
            return;

        var dragPos = mainCam.ScreenToWorldPoint(pEvent.position);
        GetView().obj_dragSolt.transform.position = new Vector3(dragPos.x, dragPos.y, GetView().obj_dragSolt.transform.position.z);
    }

    private void OnDragEnd(UnityEngine.EventSystems.PointerEventData pEvent, int cfgId,bool isEnd = true)
    {
        if (SkillM.Instance.CheakSkillEquip(cfgId))
            return;

        GetView().obj_dragSolt.transform.gameObject.SetActive(false);
        GameObject target = pEvent.pointerCurrentRaycast.gameObject;
        if (target != null && target.tag == "skill" && isEnd)
        {
            int solt = int.Parse(target.gameObject.name.Replace("solt_", ""));
            if (solt == 0 || CFG_XMLConfig.xmlConfig.SkillSlotDic[solt] > RoleM.Instance.Level)
            {
                UITool.ShowTips("槽位未解锁");
                return;
            }
            if (((CFG_Skill.GetById(cfgId).Skilltype == (int)CFG_Skill.SkillType.Job || CFG_Skill.GetById(cfgId).Skilltype == (int)CFG_Skill.SkillType.JobCtr) && solt !=0 ) 
                || ((CFG_Skill.GetById(cfgId).Skilltype == (int)CFG_Skill.SkillType.Super || CFG_Skill.GetById(cfgId).Skilltype == (int)CFG_Skill.SkillType.SuperCtr) && solt == 0))
            {
                TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_SkillChange, MType.Request, TCPClientMgr.ProtoSerialize(new SkillChangeRequest
                {
                    SkillId = cfgId,
                    OpType = 1,
                    Slot = (SkillsSlot)solt
                }));

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "打造",
                    component = "技能",
                    position = "设置界面",
                    event_click = "技能上阵",
                });
            }
        }
    }

    //技能配置
    private void OnChangSkill()
    {
        //目前屏蔽0号槽位 等其他功能再开放
        for (int i = 0; i < 13; i++)
        {
            bool isNull = i == 0 ? true: SkillM.Instance.GetSkillEquipBySolt(i) == 0;
            GetView().root_SkillSolt.Find("solt_" + i+ "/icon").gameObject.SetActive(!isNull);

            int openLevel = CFG_XMLConfig.xmlConfig.SkillSlotDic[i];
            GetView().root_SkillSolt.Find("solt_" + i + "/lock").gameObject.SetActive(i == 0 || RoleM.Instance.Level < openLevel);
            if (RoleM.Instance.Level < openLevel && i != 0)
            {
                JUIExtensions.txt(GetView().root_SkillSolt.Find("solt_" + i + "/lock/Text")).text = openLevel + "级\n解锁";
            } 
            else if (i != 0)
            {
                JUIExtensions.txt(GetView().root_SkillSolt.Find("solt_" + i + "/lock/Text")).text = "特技暂未开放";
            }
            JUIExtensions.btn(GetView().root_SkillSolt.Find("solt_" + i)).onClick.RemoveAllListeners();
            if (!isNull)
            {
                int cfgID = SkillM.Instance.GetSkillEquipBySolt(i);
                var cfgInfo = CFG_Skill.GetById(cfgID);
                GameUtil.Instance.SetItemIcon((int)cfgInfo.SkillIcon, JUIExtensions.img(GetView().root_SkillSolt.Find("solt_" + i + "/icon")));
                
                var skillInfo = SkillM.Instance.GetSkillInfoByCFGID(cfgID);
                if (skillInfo == default)
                    skillInfo = PartnerM.Instance.GetSkillInfoBySlotSkillID(cfgID);

                JUIExtensions.txt(GetView().root_SkillSolt.Find("solt_" + i+ "/icon/txt_Star")).text = skillInfo.Advanced == 0 ? "" : skillInfo.Advanced + "阶";
                JUIExtensions.txt(GetView().root_SkillSolt.Find("solt_" + i + "/icon/txt_Reinforce")).text = cfgInfo.SkillLevel == 0 ? "" : cfgInfo.SkillLevel.ToString();

                int idx = i;
                JUIExtensions.btn(GetView().root_SkillSolt.Find("solt_" + i)).onClick.AddListener(() =>
                {
                    if (!isNull)
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        OnSkillEuqipOut(idx);
                    }
                });

            }
            
        }
    }

    private void OnSkillEuqipOut(int idx)
    {
        TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_SkillChange, MType.Request, TCPClientMgr.ProtoSerialize(new SkillChangeRequest {
            SkillId = SkillM.Instance.GetSkillEquipBySolt(idx),
            OpType = 2,
            Slot = (SkillsSlot)idx
        }));
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_SkillChange(GameEvent.Sync_SkillChange _e)
    {
        OnChangSkill();
        OnRefreshAllSkill();
    }

    private enum ResID 
    {
        Coin = 1,
        SY = 23,
        JNJZ =24,
    }
    //顶部资源显示
    private void OnTopResouse()
    {
        GameUtil.Instance.SetItemIcon((int)ResID.Coin, GetView().img_topRes1);
        GetView().txt_topRes1.text = GameUtil.Instance.FormatNumWithUnit(RoleM.Instance.Coin);

        GameUtil.Instance.SetItemIcon((int)ResID.SY, GetView().img_topRes2);
        GetView().txt_topRes2.text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID((int)ResID.SY));

        GameUtil.Instance.SetItemIcon((int)ResID.JNJZ, GetView().img_topRes3);
        GetView().txt_topRes3.text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID((int)ResID.JNJZ));

        //GameUtil.Instance.SetItemIcon((int)ResID.DianQuan, GetView().img_topRes4);
        //GetView().txt_topRes4.text = GameUtil.Instance.FormatNumWithUnit(RoleM.Instance.Coupon);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_TopResouse(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        OnBtnLock();
        OnTopResouse();
        OnRefreshAllSkill();

        //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Advance.ToString(), RedDotM.Instance.GetAllSkillAdvanceRedNum());
        //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_UpLv.ToString(), RedDotM.Instance.GetAllSkillUpLvRedNum());
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnEvent_Source_Refresh(Event_Source_Refresh _e)
    {
        if (isNullGo)
            return;

        OnTopResouse();
        OnRefreshAllSkill();

        //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Advance.ToString(), RedDotM.Instance.GetAllSkillAdvanceRedNum());
        //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_UpLv.ToString(), RedDotM.Instance.GetAllSkillUpLvRedNum());
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnEvent_Sync_TakeSource(GameEvent.Sync_TakeSource _e)
    {
        OnEvent_Source_Refresh(default);
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);
        UniGameEvent.AddEventListener<GameEvent.Sync_TakeSource>(OnEvent_Sync_TakeSource);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_TopResouse);

        OnTopResouse();
        //OnSelcetUIBtnClick(SkillClass.UpLv);
    }

    public void OnShow(bool isShow)
    {
        if (isShow)
        {
            OnEnter();
        }
        else
        {
            OnExit();
        }
        uiGo.SetActive(isShow);
    }

    public override void OnExit()
    {
        ClearRed();
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_TakeSource>(OnEvent_Sync_TakeSource);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_TopResouse);
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}