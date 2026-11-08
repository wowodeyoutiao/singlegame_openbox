using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using DG.Tweening;
using JEngine.Helper;

public class UINoticeRun : UIBase<UIModelNoticeRun, UIViewNoticeRun>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private string JumpStr = string.Empty;

    public override void OnInit()
    {
        GetView().btn_Notice.onClick.AddListener(() =>
        {
            if (!string.IsNullOrEmpty(JumpStr))
            {
                OnJumpFun(JumpStr);
            }
        });
    }

    private void OnJumpFun(string str)
    {
        if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
        {
            UITool.ShowTips("战斗中无法操作");
            return;
        }
        switch (JumpStr)
        {
            case "UITrialFund":
                if (!GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower))
                {
                    UITool.Open<UITrialFund>();
                }

                break;
            case "UIPartnerCall":
                if (CFG_XMLConfig.GetOpenSolt(RoleM.Instance.Level) > 0)
                {
                    HttpC.GetPrivilegeDetail(1, (_r) =>
                    {
                        UITool.Close<UIShop>();
                        UITool.Open<UIPartnerCall>(null, _r);
                    });
                }

                break;
            case "UIPet":
                if (!GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet))
                    UITool.Open<UIPet>();

                break;
            case "UIFirstPay":
                if (ActivityM.Instance.GetAcitvityState(ActivityType.FirstPay))
                {
                    HttpC.GetFirstPayInfo((req) => {
                        if (req.code == 0 && !isNullGo)
                            UITool.Open<UIFirstPay>(null, req);
                    });
                }

                break;
            case "UIPrivilegeCard":
                UITool.Open<UIPrivilegeCard>();
                break;
            case "UIBuild":
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar))
                    return;
                if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar))
                    return;
                UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Star);
                break;
            case "UISynthesis":
                UITool.Open<UISynthesis>(null, null, UISynthesis.UISelect.Jewel);

                break;
            case "UIRank":
                UITool.Open<UIRank>();
                break;

            case "UIWarFightMap":

                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight) 
                    || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WarFight))
                    return;
                UITool.Open<UIWarFightMap>();
                break;

            case "UIWorldBoss":
                if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss)
                    || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss))
                    return;
                UITool.Open<UIWorldBoss>();
                break;

            case "UIGiftBag":
                UITool.Open<UISpecialOffer>(null, ActivityType.Month);
                break;
            default:
                break;
        }
    }

    public override void OnEnter(params object[] args)
    {
        SetRunNotice();
    }

    private void SetRunNotice()
    {
        if (!TipM.Instance.CheckRunNotice())
        {
            UITool.Close<UINoticeRun>();
            return;
        }

        TipM.RunNoticeClass runInfo = TipM.Instance.GetRunNotice();
        JumpStr = CFG_Lamp.GetById(runInfo.Id).Panel;

        GetView().txt_Notice.text = runInfo.Str;
        GetView().txt_Jump.gameObject.SetActive(CFG_Lamp.IsJump(runInfo.Id));
        if (CFG_Lamp.IsJump(runInfo.Id))
        {
            GetView().txt_Jump.text = "立即前往!";
        }

        GameHelper.Instance.SetDeayFrame(() =>
        {
            float startPosX = GetView().rect_Bg.rect.width / 2 + GetView().rect_Notice.rect.width / 2;
            GetView().rect_Notice.localPosition = Vector3.right * startPosX;
            float endPosX = -startPosX;
            if (CFG_Lamp.IsJump(runInfo.Id))
            {
                endPosX = -startPosX - GetView().txt_Jump.rectTransform.rect.width;
            }
            GetView().rect_Notice.DOLocalMoveX(endPosX, 81).SetSpeedBased().SetEase(Ease.Linear).OnComplete(() =>
            {
                SetRunNotice();
            });
        });
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {
        TipM.Instance.ClearRunNotice();
    }
}