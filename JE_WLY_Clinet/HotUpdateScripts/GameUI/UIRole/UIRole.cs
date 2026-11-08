using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using Spine.Unity;
using System;
using System.Collections.Generic;
using System.Text;
using UnityEngine;
using UnityEngine.UI;
using static CFG_XMLConfig.XML_Config;
using static GameEvent;


public class UIRole : UIBase<UIModelRole, UIViewRole>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {

        //for (int i = 0; i < GetView().selectBtns.Length; i++)
        //{
        //    int idx = i;
        //    GetView().selectBtns[i]._btn.onClick.AddListener(delegate 
        //    { 
        //        AudioManager.Instance.PlaySound(AudioKey.Sound.Button); 
        //        OnSelectBtnClick(idx); 
        //    });
        //}

        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIRole>();
        });

        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            int idx = i;
            GetView().teamItems[i]._btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if(idx > 0 && CFG_XMLConfig.CheckPartnerSoltIsOpenServer(idx))
                {
                    UITool.ShowTips($"{CFG_XMLConfig.RemainingDaysPartnerSlotUnlockOpen(idx)}天后开启");
                    return;
                }

                if (idx > 0)
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


        GetView().btn_AllEquip.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_EquipBest(TeamIdx == 0 ? ObjectType.PLAYER : ObjectType.PARTNER, TeamIdx);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = openFrom,
                position = $"",
                event_click = $"一件穿戴装备",
            });
        });

        GetView().btn_Attribute.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            SetArrtInfo(!AttrShow);
        });

        GetView().btn_PartnerSet.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
            {
                UITool.ShowTips("战斗中不可操作");
                return;
            }

            UITool.Open<UIPartnerUp>(null, TeamIdx);
        });

        GetView().btn_Pet.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Pet))
            {
                UITool.Open<UIComRule>(null, 42);
                return;
            }
            if (GameUtil.Instance.CheakIsLock(SwitchOpne.Pet, false))
            {
                UITool.Open<UIComRule>(null, 42);
                return;
            }

            UITool.Open<UIPet>();
        });

        for (int i = 0; i < GetView().selectBtns.Length; i++)
        {
            var txt = GetView().selectBtns[i]._txtOff;
            GetView().selectBtns[i]._btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnSelectBtnClick(txt.text);
            });
        }

        for (int i = 0; i < GetView().root_Equips.childCount; i++)
        {
            int idx = i + 1;
            JUIExtensions.btn(GetView().root_Equips.Find("equip_" + idx)).onClick.AddListener(() => {
                var itemInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)idx).Item : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx, idx);
                if (itemInfo != null)
                {
                    Action callback = () =>
                    {
                        UITool.Close<UIItemShow>();
                        UITool.Open<UIEquipSelect>(null, CFG_Item.GetById(itemInfo.CfgId).Equiptype, TeamIdx);

                    };

                    var equipRefineData = RoleM.Instance.GetOneEquipSlotData((CharacterSlot)(TeamIdx + 1), (EquipSlot)idx);

                    UITool.Open<UIItemShow>(null, new UIItemShow.Params
                    {
                        cfgId = itemInfo.CfgId,
                        equipInfo = itemInfo.Equip,
                        teamIdx = TeamIdx,
                        isIgnoreJob = false,
                        callback = callback,
                        buttonName = "替换",
                        characterIndex = TeamIdx + 1,
                        equipSlotPos = idx,
                        equipRefineData = equipRefineData,
                    });
                }
            });
        }

        GetView().btn_Recast.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            long partnerUid = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid;
            var partner = PartnerM.Instance.GetOnePartnerInfoByUid(partnerUid);

            if(partner.Level < PartnerM.recastOpenLevel)
            {
                UITool.ShowTips($"伙伴等级超过{PartnerM.recastOpenLevel}级可重铸");
                return;
            }

            OnRecastPartner();
        });

        GetView().btn_Disassembling.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnDisassemblingPartner();
        });

        GetView().btn_PreviewStarInfo.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            var partnerUid = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid;
            var partner = PartnerM.Instance.GetOnePartnerInfoByUid(partnerUid);
            UITool.Open<UIStarDetails>(null, CFG_Partner.GetById(partner.CfgId));
        });
    }

    private void OnBindRed()
    {
        RedPointManager.Instance.RPS.Register(RedPointConst.Role_Equip.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().obj_EqiuipDot.gameObject.SetActive(RedDotM.Instance.GetOneTeamBeatEquipRed(TeamIdx) > 0);
        });

        // 宠物上阵
        RedPointManager.Instance.RPS.Register(RedPointConst.Partner_Into.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().btn_PartnerSet.transform.Find("dot").gameObject.SetActive(TeamIdx > 0 && RedDotM.Instance.OnePartnerIntoRedDot(TeamIdx) > 0);
        });

        // 打造 技能
        RedPointManager.Instance.RPS.Register(RedPointConst.Role_Skill.ToString(), (node) =>
        {
            if (isNullGo) return;
            //GetView().selectBtns[1]._objDot.SetActive(node.pointCount > 0);
            //if (TeamIdx == 0)
            //{
            //    GetView().selectBtns[2]._objDot.SetActive(node.pointCount > 0);
            //    //GetView().selectBtns[4]._objDot.SetActive(node.pointCount > 0);
            //}
            //else
            //{
            //    GetView().selectBtns[3]._objDot.SetActive(node.pointCount > 0);
            //    //GetView().selectBtns[5]._objDot.SetActive(node.pointCount > 0);
            //}
        });

        // 打造
        RedPointManager.Instance.RPS.Register(RedPointConst.Role_Build.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().selectBtns[1]._objDot.SetActive(node.pointCount > 0);
        });

        if(TeamIdx == 0)
        {
            // 宠物
            RedPointManager.Instance.RPS.Register(RedPointConst.Role_Pet.ToString(), (node) =>
            {
                if (isNullGo) return;
                if (TeamIdx == 0)
                    GetView().selectBtns[2]._objDot.SetActive(node.pointCount > 0);
            });
        }
        else
        {
            // 升星
            RedPointManager.Instance.RPS.Register(RedPointConst.Partner_UpStar.ToString(), (node) =>
            {
                if (isNullGo) return;
                if (TeamIdx > 0)
                    GetView().selectBtns[2]._objDot.SetActive(RedDotM.Instance.GetOnePartnerUpStarRedDot(TeamIdx) > 0);
            });
        }

        // 称号
        RedPointManager.Instance.RPS.Register(RedPointConst.Role_Title.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().selectBtns[3]._objDot.SetActive(RedDotM.Instance.GetOneTeamTitleActiveRedNum(TeamIdx,false) > 0 || RedDotM.Instance.GetOneTeamTitleActiveRedNum(TeamIdx, true) > 0);
        });
    }

    private string[] RoleBtns = {  "技能", "打造", "宠物", "称号", "外观" };
    private string[] PartnerBtns = { "技能", "打造", "升星", "称号" };
    private string[] PartnerNoBtns = { "打造" };

    //private SwitchOpne[] btrnType =
    //{
    //    SwitchOpne.No,
    //    SwitchOpne.Build_SkillLv,
    //    SwitchOpne.Build_SkillLv,
    //    SwitchOpne.Pet,
    //    SwitchOpne.Build_SkillStar,
    //    SwitchOpne.Title,
    //};
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
    private void OnSelectBtnClick(string btnStr)
    {
        switch (btnStr)
        {
            case "羁绊":
                //if (GameUtil.Instance.CheakIsLock(SwitchOpne.ExpPool))
                //    return;
                UITool.Open<UIBond>();
                break;
            case "宠物":
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Pet))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Pet))
                    return;
                
                UITool.Open<UIPet>();
                break;
            case "打造":
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Build_SkillLv))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Build_SkillLv))
                    return;

                UITool.Open<UIBuild>();
                break;
            case "技能":
                if (TeamIdx != 0 && GameUtil.Instance.CheakIsLock(SwitchOpne.Role_SkillAndCall))
                    return;

                UITool.Open<UIRole_Skill>(null,TeamIdx);
                break;
            case "升星":
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Partner_Star))
                    return;

                UITool.Open<UIPartnerUpStar>(null, TeamIdx);
                break;
            case "称号":
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Title))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Title))
                    return;
                FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Title });
                UITool.Open<UIRoleTitile>(); 
                break;
            case "外观":
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.AvatarAndFeame))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.AvatarAndFeame))
                    return;

                UITool.Open<UIRoleAvatar>();
                break;
            default:
                UITool.ShowTips("暂未开放");
                break;
        }
    }


    private void SetSelectBtnState()
    {
        var btns = TeamIdx == 0 ? RoleBtns : (PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid != 0 ? PartnerBtns : PartnerNoBtns);
        for (int i = 0; i < GetView().selectBtns.Length; i++)
        {
            GetView().selectBtns[i]._btn.gameObject.SetActive(i < btns.Length);
            if (i < btns.Length)
            {
                GetView().selectBtns[i]._txtOff.text = btns[i];
            }
        }
    }

    private void SetSelectBtnLock(int index, bool isLock, string openCondition = "")
    {
        for (int i = 0; i < GetView().selectBtns.Length; i++)
        {
            if(index == i)
            {
                GetView().selectBtns[i]._objMask.SetActive(isLock);
                if (isLock)
                {
                    GetView().selectBtns[i]._txtOpen.text = openCondition;
                }
            }
        }
    }

    private void ClearSelectButtonLock()
    {
        for (int i = 0; i < GetView().selectBtns.Length; i++)
        {
            GetView().selectBtns[i]._objMask.SetActive(false);
            GetView().selectBtns[i]._txtOpen.text = "";
        }
    }

    private void SetPartnerLock()
    {
        for (int soltIdx = 1; soltIdx < GetView().teamItems.Length; soltIdx++)
        {
            bool isLock = CFG_XMLConfig.CheckPartnerSoltIsOpenServer(soltIdx, false);
            (bool isOtherLock, string strtxttip, string strtip) = CheakIsLock(CFG_XMLConfig.GetOpenSoltLv(soltIdx), CFG_XMLConfig.GetOpenSolMainStaget(soltIdx));
            if (!isLock)
                isLock = isOtherLock;

            GetView().teamItems[soltIdx]._obj_lock.SetActive(isLock);
            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(soltIdx);
            GetView().teamItems[soltIdx].spine_Unlock.gameObject.SetActive(!isLock && solt.PartnerUid == 0);

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
                GetView().teamItems[i]._img_head.gameObject.SetActive( solt.PartnerUid != 0);
                GetView().teamItems[i]._img_job.gameObject.SetActive(solt.PartnerUid != 0);
                GetView().teamItems[i]._obj_Lv.gameObject.SetActive(solt.PartnerUid != 0);
                GetView().teamItems[i].txt_Star.gameObject.SetActive(solt.PartnerUid != 0);

                bool isLock = CFG_XMLConfig.CheckPartnerSoltIsOpenServer(i, false);
                (bool isOtherLock, string strtxttip, string strtip) = CheakIsLock(CFG_XMLConfig.GetOpenSoltLv(i), CFG_XMLConfig.GetOpenSolMainStaget(i));
                if (!isLock)
                    isLock = isOtherLock;
                GetView().teamItems[i].spine_Unlock.gameObject.SetActive(!isLock && solt.PartnerUid == 0);

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
        GetView().txt_allPower.text = GameUtil.Instance.FormatPowerWithUnit(RoleM.Instance.TotalCombat);

        if (isInit)
        {
            OnTeamItemClick(TeamIdx);
        }
    }

    private void OnPetShow()
    {
        var myPet = PetM.Instance.GetOnePetByIdx(0);
        GetView().btn_Pet.transform.Find("noPet").gameObject.SetActive(myPet.CfgId == 0);
        GetView().btn_Pet.transform.Find("Pet").gameObject.SetActive(myPet.CfgId != 0);

        if (myPet.CfgId != 0)
        {
            var root = GetView().btn_Pet.transform.Find("Pet");
            string pooName = "Pet_"+ myPet.CfgId;
            _ = GameHelper.Instance.SetAsyncObj(pooName, GameConst.GetMonsterModelByID(pooName), (_mainRole) =>
            {
                if (root.childCount > 0)
                {
                    PoolHelper.Instance.RecycObj(root.GetChild(0).gameObject);
                }
                _mainRole.transform.SetParent(root);
                _mainRole.transform.localScale = Vector3.one;
                _mainRole.transform.localPosition = Vector3.zero;
                _mainRole.SetActive(true);
            });
        }
        
    }

    public override void OnLayerChange()
    {
        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            GameUtil.Instance.SetEffectLayer(GetView().teamItems[i]._obj_on.transform, _canvasOrder + 1);
        }
    }

    private int TeamIdx = -1;
    private void OnTeamItemClick(int idx)
    {
        TeamIdx = idx;

        for (int i = 0; i < GetView().selectBtns.Length; i++)
        {
            GetView().selectBtns[i]._objDot.SetActive(false);
        }

        OnLockCheak();
        OnBindRed();
        OnBaseInfo();

        for (int i = 0; i < GetView().teamItems.Length; i++)
        {
            GetView().teamItems[i]._obj_on.SetActive(i == idx);
        }

        GetView().btn_Pet.gameObject.SetActive(TeamIdx == 0);
        GetView().btn_PartnerSet.gameObject.SetActive(TeamIdx != 0);
        GetView().btn_AllEquip.gameObject.SetActive(TeamIdx == 0 || PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid != 0);

        if (TeamIdx != 0)
        {

            GetView().txt_PartnerSet.text = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid == 0 ? "上阵" : "替换";
            GetView().btn_PreviewStarInfo.gameObject.SetActive(PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid != 0);
        }
        else
        {
            OnPetShow();
            GetView().btn_PreviewStarInfo.gameObject.SetActive(false);
        }

        //ShowRecastDisassembleButton();

        ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/Role", TeamIdx == 0 ? "bg_juese" : "bg_huoban", GetView().img_Role);
    }

    private void ShowRecastDisassembleButton()
    {
        bool isShow = TeamIdx != 0 && PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid != 0;
        GetView().root_RecastDisassemblingBtn.gameObject.SetActive(isShow);
        if (isShow)
        {
            long partnerUid = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid;
            var partner = PartnerM.Instance.GetOnePartnerInfoByUid(partnerUid);
            GetView().btn_Recast.gameObject.SetActive(partner.Level > 1 || partner.Star > 0);
            var cfg_Partner = CFG_Partner.GetById(partner.CfgId);
            GetView().btn_Disassembling.gameObject.SetActive(partner.Star >= 5 && cfg_Partner.Quality >= 7);
        }
    }

    //装备槽位信息
    private void OnBaseInfo(bool isRefreshBase = false)
    {
        string Str = string.Empty;
        if (TeamIdx == 0)
        {
            Str =  RoleM.Instance.Nickname;
            if (!string.IsNullOrEmpty(GuildM.Instance.GuildName))
                Str = $"<size=24><color=#4DAE1F>[{GuildM.Instance.GuildName}]</color></size>\n{Str}";
        }
        else
        {
            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
            if (solt.PartnerUid != 0)
            {
                var cfg_Partner = CFG_Partner.GetById(PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid).CfgId);
                string name = cfg_Partner.Name;
                Str = GameUtil.Instance.GetFightStrByQuality(cfg_Partner.Quality, name);
            }
        }

        GetView().txt_RoleName.text = Str;

        for (int i = 0; i < GetView().root_Equips.childCount; i++)
        {
            int idx = i + 1;
            var go = GetView().root_Equips.Find("equip_" + idx);

            var equipInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)idx) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamIdx, idx);
            JUIExtensions.txt(go.Find("txt_Starlv")).text = equipInfo.Stars.ToString();
            go.Find("txt_Starlv").gameObject.SetActive(equipInfo.Stars > 0);
            JUIExtensions.txt(go.Find("txt_lv")).text = equipInfo.Enhance == 0 ? "" : "+" + equipInfo.Enhance;

            var itemInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)idx).Item : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx, idx);
           
            var bastEffectTrans = go.Find("BestEffect");
            if (itemInfo != null && itemInfo.Equip != null && (int)itemInfo.Equip.Quality >= 4 && itemInfo.Equip.Attributes != null && itemInfo.Equip.Attributes.Count > 0)
            {
                bastEffectTrans.gameObject.SetActive(true);
                SpineUtil.Instance.PlayItemSpineLoop(bastEffectTrans, (int)itemInfo.Equip.Quality);
            }
            else
                bastEffectTrans.gameObject.SetActive(false);

            for (int j = 0; j < equipInfo.Gemstone.Length; j++)
            {
                go.Find("inset").GetChild(j).gameObject.SetActive(equipInfo.Gemstone[j] != 0);
                if (equipInfo.Gemstone[j] != 0)
                {
                    GameUtil.Instance.SetItemIcon(equipInfo.Gemstone[j], JUIExtensions.img(go.Find("inset").GetChild(j)));
                }
            }
        }

        OnTitleShow();
        if (!isRefreshBase)
            OnEquipShow();
    }

    /// <summary>
    /// 称号显示
    /// </summary>
    private void OnTitleShow()
    {
        GameUtil.Instance.SetTitleIcon(TitleM.Instance.GetTitleBySolt(TeamIdx),GetView().img_title);
    }

    //设置装备信息
    private void OnEquipShow(bool isRefsh = false)
    {
        string Str = "";
        string lvStr = "";
        if (TeamIdx == 0)
        {
            Str = GameUtil.Instance.FormatPowerWithUnit(RoleM.Instance.Combat);
            lvStr = RoleM.Instance.Level.ToString();

            string poolKey = "MainRole" + RoleM.Instance.Job + RoleM.Instance.Sex;
            SetRoleSpine(poolKey, GameConst.GetModelByJobAndSex((HotUpdateScripts.ProtoM.JobType)RoleM.Instance.Job, (SexType)RoleM.Instance.Sex));
        }
        else
        {
            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
            string poolKey = "";
            if (solt.PartnerUid != 0 && PartnerM.Instance.GetOnePartnerStats(solt.PartnerUid) != null)
            {
                Str = GameUtil.Instance.FormatPowerWithUnit(PartnerM.Instance.GetOnePartnerStats(solt.PartnerUid).Combat);
                lvStr = PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid).Level.ToString();
                poolKey = "Partner_" + PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid).CfgId;
                //poolKey = "Partner_1";
            }

            if (!isRefsh)
                SetRoleSpine(poolKey, poolKey == "" ? "" :GameConst.GetMonsterModelByID(poolKey));
        }

        GetView().txt_RolePower.text = Str;
        GetView().txt_RolePower.gameObject.SetActive(!string.IsNullOrEmpty(Str));

        GetView().txt_lv.text = lvStr;
        GetView().txt_lv.gameObject.SetActive(!string.IsNullOrEmpty(lvStr));


        for (int i = 0; i < GetView().root_Equips.childCount; i++)
        {
            int idx = i + 1;
            var go = GetView().root_Equips.Find("equip_" + idx);

            var itemInfo = TeamIdx == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)idx).Item : PartnerM.Instance.GetOnePartnerOneEquip(TeamIdx,idx);
            JUIExtensions.img(go.Find("icon")).enabled = itemInfo != null;
            JUIExtensions.img(go.Find("iconbg")).enabled = itemInfo != null;
            JUIExtensions.btn(go).enabled = itemInfo != null;

            var bastEffectTrans = go.Find("BestEffect");
            if (itemInfo != null && itemInfo.Equip != null && (int)itemInfo.Equip.Quality >= 4 && itemInfo.Equip.Attributes != null && itemInfo.Equip.Attributes.Count > 0)
            {
                bastEffectTrans.gameObject.SetActive(true);
                SpineUtil.Instance.PlayItemSpineLoop(bastEffectTrans, (int)itemInfo.Equip.Quality);
            }
            else
                bastEffectTrans.gameObject.SetActive(false);

            bool hasEquip = itemInfo != null;
            go.Find("defaultImg").gameObject.SetActive(!hasEquip);

            if (hasEquip)
            {
                //JUIExtensions.txt(go.Find("txt_name")).text = "";
                GameUtil.Instance.SetItemIcon(itemInfo.CfgId, JUIExtensions.img(go, "icon"), JUIExtensions.img(go),null,(int)itemInfo.Equip.Quality);
            }
            else
            {
                //JUIExtensions.txt(go.Find("txt_name")).text = GameM.Instance.EquipName[idx];
                ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/Common", "bg_zhuangbei", JUIExtensions.img(go));
            }
        }
        SetArrtInfo(false);
        SetSelectBtnState();
    }

    private bool AttrShow = false;

    //设置玩家属性信息
    private void SetArrtInfo(bool isShow)
    {
        AttrShow = isShow;
        if (TeamIdx == 0)
        {
            GetView().btn_Attribute.gameObject.SetActive(true);
            (GetView().txt_attr0.text, GetView().txt_attr1.text, GetView().txt_attr2.text, GetView().txt_attr3.text) = GetRoleAttribute((isShow ? -1 :4), RoleM.Instance.GetRoleAllStats());
        }
        else
        {
            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
            GetView().btn_Attribute.gameObject.SetActive(solt.PartnerUid != 0);
            if (solt.PartnerUid != 0 && PartnerM.Instance.GetOnePartnerStats(solt.PartnerUid) != null)
                (GetView().txt_attr0.text, GetView().txt_attr1.text, GetView().txt_attr2.text, GetView().txt_attr3.text) = GetRoleAttribute((isShow ? -1 : 4), PartnerM.Instance.GetOnePartnerStats(solt.PartnerUid).Attributes);
        }
        GetView().img_Arrow.transform.localScale = !isShow ? Vector3.one : new Vector3(1, -1, 1);
        LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().btn_Attribute.GetComponent<RectTransform>());
    }


    private void SetRoleSpine(string poolKey,string path)
    {

        if (GetView().root_Role.childCount > 0)
        {
            if (poolKey == GetView().root_Role.GetChild(0).gameObject.name)
                return;
        }

        void ClearRole()
        {
            int Count = GetView().root_Role.childCount;
            for (int i = 0; i < Count; i++)
            {
                PoolHelper.Instance.RecycObj(GetView().root_Role.GetChild(0).gameObject);
            }
        }

        if (!string.IsNullOrEmpty(poolKey))
        {
            _ = GameHelper.Instance.SetAsyncObj(poolKey, path, (_mainRole) =>
            {
                ClearRole();
                _mainRole.transform.SetParent(GetView().root_Role);
                _mainRole.transform.localScale = Vector3.one;
                _mainRole.transform.localPosition = Vector3.zero;
                _mainRole.SetActive(true);

                SpineUtil.Instance.PlayLoop(_mainRole.transform.Find("spine").GetComponent<SkeletonGraphic>(), GameConst.ModelAnim.stand.ToString());
            });
        }
        else
        {
            ClearRole();
        }

    }


    //private int SelectIdx = 0;
    //private  void OnSelectBtnClick(int idx)
    //{
    //    SelectIdx = idx;
    //}

    private void OnSelectButtonClick(int index)
    {
        var txt = GetView().selectBtns[index]._txtOff;
        OnSelectBtnClick(txt.text);
    }

    private (string,string, string, string) GetRoleAttribute(int Count, List<Stats> StatsList)
    {
        int idx = 0;
        string Str0 = "";
        string Str1 = "";
        string Str2 = "";
        string Str3 = "";

        int Amount = 4;
        for (int i = 0; i < CFG_XMLConfig.xmlConfig.RolePanelList.Count; i++)
        {
            var item = CFG_XMLConfig.xmlConfig.RolePanelList[i];
            if (item.Show == 0)
                continue;

            if (idx % Amount == 0)
            {
                Str0 = $"{Str0}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            else if(idx % Amount == 1)
            {
                Str1 = $"{Str1}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            else if (idx % Amount == 2)
            {
                Str2 = $"{Str2}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            else
            {
                Str3 = $"{Str3}{item.Name}:<color=#ffffff>{CFG_XMLConfig.GetAttrValueToStr(item.Type, GetValue(item.SID))}</color>\n";
            }
            idx++;

            if (idx == Count)
                break;
        }

        long GetValue(int SId)
        {
            for (int i = 0; i < StatsList.Count; i++)
            {
                if (StatsList[i].attr == (pb.Attribute) SId)
                    return StatsList[i].Num;
            }
            return 0;
        }
        
        return (Str0, Str1, Str2, Str3);
    }

    string openFrom;

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);
        UniGameEvent.AddEventListener<GameEvent.Refresh_LvUp>(OnRefresh_LvUp);
        UniGameEvent.AddEventListener<GameEvent.Refresh_UseAvatar>(OnRefrersh_Avatar);

        if (args.Length > 0)
        {
            int idx = (int)args[0];
            TeamIdx = idx < 0 ? 0 : idx;
            openFrom = idx < 0 ? "任务跳转" : "首页头像框";

            if (args.Length > 1)
            {
                ((Action)args[1])?.Invoke();
            }
        }
        else
        {
            TeamIdx = 0;
            openFrom = "角色";
        }

        SetPartnerLock();
        SetBaseRole(true);

        NewPlayerGuideWidget.Instance.AddGuideGameObject("UIRole_BtnClose", GetView().btn_Close.transform);
        NewPlayerGuideWidget.Instance.AddGuideGameObject("UIRole_SetPartner", GetView().btn_PartnerSet.transform);

        UniGameEvent.AddEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);

        refreshFuncNewFlag();
    }

    private void OnLockCheak()
    {
        ClearSelectButtonLock();
        if (TeamIdx == 0)
        {
            bool isLockOpenDayPet = GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Pet, false);
            StringBuilder tipsStage = new("");
            bool isLockByLvAndStage = GameUtil.Instance.CheakIsLock(SwitchOpne.Pet, false, tipsStage);

            StringBuilder tipsAvatar = new("");
            bool isLockOpenDayAvatar = GameUtil.Instance.CheakIsLock(SwitchOpne.AvatarAndFeame, false, tipsAvatar);

            GetView().obj_PetLock.SetActive(isLockOpenDayPet);
            if (isLockOpenDayPet)
            {
                GetView().txt_petLock.text = $"{GameUtil.Instance.RemainingDaysUnlockOpen(SwitchOpne.Pet)}天后开启";
                SetSelectBtnLock(2, true, $"{GameUtil.Instance.RemainingDaysUnlockOpen(SwitchOpne.Pet)}天后开启");
            }
            else
            {
                SetSelectBtnLock(2, false);
            }

            SetSelectBtnLock(2, isLockByLvAndStage, tipsStage.ToString());
            SetSelectBtnLock(4, isLockOpenDayAvatar, tipsAvatar.ToString());

            //else if (isLockByLvAndStage)
            //{
            //    //GetView().obj_PetLock.SetActive(isLockByLvAndStage);
            //    //GetView().txt_petLock.text = tipsStage.ToString();
            //    SetSelectBtnLock(2, isLockByLvAndStage, tipsStage.ToString());
            //}
            //else if (isLockOpenDayAvatar)
            //{
            //    //GetView().obj_PetLock.SetActive(true);
            //    //GetView().txt_petLock.text = tipsAvatar.ToString();
            //    SetSelectBtnLock(4, true, tipsAvatar.ToString());
            //}
            //else
            //{
            //    GetView().obj_PetLock.SetActive(false);
            //}
            //else
            //{
            //    StringBuilder tips = new("");
            //    bool isLockByLvAndStage = GameUtil.Instance.CheakIsLock(SwitchOpne.AvatarAndFeame, false, tips);

            //    GetView().obj_PetLock.SetActive(isLockByLvAndStage);
            //    GetView().txt_petLock.text = tips.ToString();
            //    SetSelectBtnLock(4, isLockByLvAndStage, tips.ToString());
            //}
        }
        else
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Partner_Star, false))
            {
                SetSelectBtnLock(2, true, $"{GameUtil.Instance.RemainingDaysUnlockOpen(SwitchOpne.Partner_Star)}天后开启");
            }
            else
            {
                StringBuilder tips = new("");
                bool isLockByLvAndStage = GameUtil.Instance.CheakIsLock(SwitchOpne.Partner_Star, false, tips);

                //int skillLockLv = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[SwitchOpne.Partner_Star];
                SetSelectBtnLock(2, isLockByLvAndStage, tips.ToString());
            }
        }

        if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Title, false))
        {
            SetSelectBtnLock(3, true, $"{GameUtil.Instance.RemainingDaysUnlockOpen(SwitchOpne.Title)}天后开启");
        }
        else
        {
            StringBuilder tips = new("");
            bool isLockByLvAndStage = GameUtil.Instance.CheakIsLock(SwitchOpne.Title, false, tips);
            //int skillLockLv = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[SwitchOpne.Title];
            SetSelectBtnLock(3, isLockByLvAndStage, tips.ToString());
        }

    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_LvUp(GameEvent.Refresh_LvUp _e)
    {
        if (isNullGo)
            return;
        OnLockCheak();
    }

    [Subscriber(ThreadMode.Main)]
    private void Refresh_TitleList(GameEvent.Refresh_TitleList _e)
    {
        if (isNullGo)
            return;
        OnTitleShow();
    }

    private void OnRecastPartner()
    {
        long partnerUid = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid;
        UITool.Open<UIPartnetRecastDisassembling>(null, partnerUid, TeamIdx, UIPartnetRecastDisassembling.PartnerOptionType.Recast);
    }

    private void OnDisassemblingPartner()
    {
        long partnerUid = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).PartnerUid;
        UITool.Open<UIPartnetRecastDisassembling>(null, partnerUid, TeamIdx, UIPartnetRecastDisassembling.PartnerOptionType.Disassembling);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_UserBase(GameEvent.Sync_UserBase _e)
    {
        SetPartnerLock();
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_UserStatsPush(GameEvent.Sync_UserStatsPush _e)
    {
        if (isNullGo)return;

        SetBaseRole(false);
        OnEquipShow(true);
    }

    [Subscriber(ThreadMode.Main)]
    private void OnRefresh_PartnerSlot(GameEvent.Refresh_PartnerSlot _e)
    {
        if (isNullGo) return;
        SetBaseRole(true);
        OnBindRed();
    }

    [Subscriber(ThreadMode.Main)]
    private void OnRefresh_PartnerList(GameEvent.Refresh_PartnerList _e)
    {
        if (isNullGo && _e.isOne && TeamIdx != 0) return;
            SetBaseRole(true);
    }

    //宠物刷新
    [Subscriber(ThreadMode.Main)]
    private void Refresh_PetsInfo(GameEvent.Refresh_PetsInfo _e)
    {
        if (isNullGo && TeamIdx != 0) return;
        OnPetShow();
        OnBindRed();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.code == 0 && _e.cid == CmdId.C2S_UpgradeLevel)
        {
            SetBaseRole(true); 
            OnBindRed();
        }

        // 伙伴重铸
        if (_e.code == 0 && _e.cid == CmdId.C2S_PartnerReforge)
        {
            SetBaseRole(true);
        }

        // 伙伴分解
        if (_e.code == 0 && _e.cid == CmdId.C2S_PartnerDisassemble)
        {
            SetBaseRole(true);
        }
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_CheckOpenTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        SetPartnerLock();
        OnLockCheak();
    }

    //槽位魂石镶嵌
    [Subscriber(ThreadMode.Main)]
    private void EquipInsetRefresh(GameEvent.EquipInsetRefresh _e)
    {
        if (isNullGo) return;
        OnBaseInfo(true);
    }

    [Subscriber(ThreadMode.Main)]
    //槽位强化
    private void Replace_EquipEnhance(GameEvent.Replace_EquipEnhance _e)
    {
        if (isNullGo) return;
        OnBaseInfo(true);
    }

    [Subscriber(ThreadMode.Main)]
    //槽位升星
    private void Replace_EquipStar(GameEvent.Replace_EquipStar _e)
    {
        if (isNullGo) return;
        OnBaseInfo(true);
    }

    [Subscriber(ThreadMode.Main)]
    //称号装备
    private void Title_Equip(GameEvent.Title_Equip _e)
    {
        if (isNullGo) return;
        OnBaseInfo();
    }

    private void OnRefrersh_Avatar(GameEvent.Refresh_UseAvatar _e)
    {
        if (isNullGo) return;
        //(string atlas, string icon) = CFG_RoleIcon.GetAvatarIcon(RoleM.Instance.CurrentAvatarId, RoleM.Instance.Job);
        //ResouseHelper.Instance.LoadSpriteByAtlas(atlas, icon, GetView().teamItems[0]._img_head);
        SetBaseRole(true);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_LvUp>(OnRefresh_LvUp);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_UseAvatar>(OnRefrersh_Avatar);
        UniGameEvent.RemoveEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);
    }

    public override void OnDestroy()
    {
        OnExit();
    }

    #region 功能开启角标

    private void OnRefresh_FunctionNewFlag(UIFunctionOpenOperate operate)
    {
        refreshFuncNewFlag();
    }

    //private List<SwitchOpne> skillGroup = new() { SwitchOpne.Build_SkillLv, SwitchOpne.Build_SkillStar};
    private List<SwitchOpne> makeGroup = new() { SwitchOpne.Build, SwitchOpne.Build_EquipStar, SwitchOpne.Refined, SwitchOpne.Build_EquipInster, SwitchOpne.Build_SkillLv, SwitchOpne.Build_SkillStar };
    private List<SwitchOpne> petGroup = new() { SwitchOpne.Pet};
    private List<SwitchOpne> titleGroup = new() { SwitchOpne.Title};
    private void refreshFuncNewFlag()
    {
        //if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Build, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Build, false)))
        //{
            bool makeFlag = FuncOpenShowWidget.Instance.FunctionIsLight(makeGroup);
            GetView().selectBtns[1]._objDotRoot.SetActive(!makeFlag);
            GetView().selectBtns[1]._objFunctionOpenNewFlag.SetActive(makeFlag);
        //}

        if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Pet, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Pet, false)))
        {
            bool petFlag = FuncOpenShowWidget.Instance.FunctionIsLight(petGroup);
            GetView().selectBtns[2]._objDotRoot.SetActive(!petFlag);
            GetView().selectBtns[2]._objFunctionOpenNewFlag.SetActive(petFlag);
        }

        if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Title, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Title, false)))
        {
            bool titleFlag = FuncOpenShowWidget.Instance.FunctionIsLight(titleGroup);
            GetView().selectBtns[3]._objDotRoot.SetActive(!titleFlag);
            GetView().selectBtns[3]._objFunctionOpenNewFlag.SetActive(titleFlag);
        }
    }
    #endregion
}