using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;

using MetaFramework.UI;

using pb;

using System.Collections;
using System.Collections.Generic;

using UnityEngine;

using static GameEvent;

public class UIExpPool : UIBase<UIModelExpPool, UIViewExpPool>
{
    private string[] JumpStrs = new string[] {"装备\n回收", "宠物\n升级", "VIP\n任务"};
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIExpPool>();
        });

        for (int i = 0; i < GetView().JumpBtns.Length; i++)
        {
            GetView().JumpBtns[i].txt.text = JumpStrs[i];
            int idx = i;
            GetView().JumpBtns[i].btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnJumpBtnClick(idx);
            });
        }

        for (int i = 0; i < GetView().Roles.Length; i++)
        {
            int idx = i;
            GetView().Roles[i].btn_up.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnRoleBtnUpClick(idx);
            });
        }

        GetView().btn_help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 25);
        });
    }

    private void OnJumpBtnClick(int idx)
    {
       if (idx == 0)
        {
            UITool.Open<UIMoreItemsDel>();
        }
        else if (idx == 1)
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet))
                return;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet))
                return;

            UITool.Open<UIPet>();
        }
        else if (idx == 2)
        {
            UITool.Open<UIVIP>();
        }
        else 
        {
            UITool.ShowTips("敬请期待");
            return;
        }
       UITool.Close<UIExpPool>();
    }

    private void OnRoleBtnUpClick(int idx)
    {
        long exp = 0;
        if (idx == 0)
        {
            exp = CFG_Role.GetByJobAndLevel(RoleM.Instance.Job, RoleM.Instance.Level).NeedEXP;
        }
        else
        {
            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(idx);
            var cfg_partner = CFG_Partner.GetById(PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid).CfgId);
            exp = CFG_Role.GetByJobAndLevel(cfg_partner.job, PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid).Level).NeedEXP;
        }
        if (!GameUtil.Instance.ComCoinToCheck((int)CFG_Item.ItemID.EXP, exp, null, false))
        {
            UITool.ShowTips("经验不足");
            return;
        }
        TcpC.C2S_UpgradeLevel(idx, (idx == 0 ? RoleM.Instance.Uid : PartnerM.Instance.GetOnePartnerSlotInfo(idx).PartnerUid));

        string roletype = idx == 0 ? "主角" : "伙伴";
        string rolename = idx == 0 ? RoleM.Instance.Nickname : CFG_Partner.GetById(PartnerM.Instance.GetOnePartnerInfoByUid(PartnerM.Instance.GetOnePartnerSlotInfo(idx).PartnerUid).CfgId).Name;
        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "经验池",
            position = "",
            event_click = $"{roletype}升级",
        });
    }

    private void OnOffLineExpShow()
    {
        var myPet = PetM.Instance.GetOnePetByIdx(0);
        string Des = "0/分钟";
        if (myPet.CfgId != 0)
        {
            var offlineAwards = CFG_Pet.GetById(myPet.CfgId).GetOfflineAwardByStar(myPet.Star);
            if (offlineAwards != null)
            {
                for (int i = 0; i < offlineAwards.Length; i++)
                {
                    if (offlineAwards[i].id == 2)
                    {
                        Des = $"{GameUtil.Instance.FormatNumWithUnit(offlineAwards[i].num)}/{offlineAwards[i].gap / 60}分钟";
                    }
                }
            }
        }
        GetView().txt_offline.text = Des;
    }

    public void OnExpShow(bool isInit = false)
    {
        float proExp = (float)RoleM.Instance.Exp / CFG_Role.GetByJobAndLevel(RoleM.Instance.Job, RoleM.Instance.Level).NeedEXP;
        GetView().img_Pro.fillAmount = proExp;
        GetView().slider_Pro.value = proExp;

        GetView().txt_expPro.text = GameUtil.Instance.GetPercentStr(proExp);
        GetView().txt_allExp.text = $"{GameUtil.Instance.FormatNumWithUnit(RoleM.Instance.Exp)}/{GameUtil.Instance.FormatNumWithUnit(CFG_Role.GetByJobAndLevel(RoleM.Instance.Job, RoleM.Instance.Level).NeedEXP)}";

        for (int i = 0; i < GetView().Roles.Length; i++)
        {
            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(i);
            GetView().Roles[i].go.SetActive(i == 0 || (solt != null && solt.PartnerUid != 0));

            long exp = 0;
            if (i == 0)
            {
                if (isInit)
                {
                    GameUtil.Instance.SetJobIcon(RoleM.Instance.Job, GetView().Roles[i].img_job);
                    GetView().Roles[i].txt_name.text = RoleM.Instance.Nickname;
                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + RoleM.Instance.Job + RoleM.Instance.Sex, GetView().Roles[i].img_head);
                }
                GetView().Roles[i].txt_level.text = "lv." + RoleM.Instance.Level;
                GetView().Roles[i].btn_up.gameObject.SetActive(CFG_Role.MaxLvByJob(RoleM.Instance.Job) != RoleM.Instance.Level);
                exp = CFG_Role.GetByJobAndLevel(RoleM.Instance.Job, RoleM.Instance.Level).NeedEXP;
            }
            else if (solt != null && solt.PartnerUid != 0)
            {
                var partner = PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid);
                var cfg_partner = CFG_Partner.GetById(partner.CfgId);
                if (isInit)
                {
                    GameUtil.Instance.SetJobIcon(cfg_partner.job % 10, GetView().Roles[i].img_job);
                    GetView().Roles[i].txt_name.text = cfg_partner.Name;
                    GameUtil.Instance.SetItemIcon(cfg_partner.ItemID, GetView().Roles[i].img_head);
                }
                GetView().Roles[i].txt_level.text = "lv." + partner.Level;
                GetView().Roles[i].btn_up.gameObject.SetActive(partner.Level < RoleM.Instance.Level);
                GetView().Roles[i].obj_tip.SetActive(partner.Level >= RoleM.Instance.Level);
                
                exp = CFG_Role.GetByJobAndLevel(cfg_partner.job, partner.Level).NeedEXP;
                GameUtil.Instance.SetPartnetQualityFrame(cfg_partner.Quality, GetView().Roles[i].img_quality);
            }
            GetView().Roles[i].txt_up.text = $"需{GameUtil.Instance.GetCheakNumStr((int)CFG_Item.ItemID.EXP,exp)}经验";
            GetView().Roles[i].obj_Red.SetActive(exp <= RoleM.Instance.Exp);
        }
    }
    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_TopResouse);

        OnOffLineExpShow();
        OnExpShow(true);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_TopResouse);

    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_TopResouse(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;
        OnExpShow();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.code == 0 && _e.cid == CmdId.C2S_UpgradeLevel)
        {
            UITool.ShowTips("升级成功");
            OnExpShow();
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_TopResouse(GameEvent.Refresh_PartnerList _e)
    {
        if (isNullGo)
            return;
        OnExpShow();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}