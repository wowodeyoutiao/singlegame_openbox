using DG.Tweening;
using HotUpdateScripts.ProtoM;
using MetaFramework.UI;
using System;
using System.Diagnostics;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 引导新角色上阵
/// </summary>
public class AddNewRoleInFormation : NewPlayerGuideStepBase
{
    public Button goFirstPartnerClick;
    public Button goOpenParternerChooseUIClick;
    public Button goChoicePartnerClick;
    public Button goEnterPartnerClick;
    public Button btnCloseUIRoleClick;
    //public Button btnGoMainBattleClick;

    bool isFirst = false;
    bool hasFreePartner = false;
    public override void Enter(UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide> parentController)
    {
        base.Enter(parentController);

        RemoveAllListeners();
        AddAllListeners();

        hasFreePartner = existFreePartner();
        if (!hasFreePartner)
        {
            NewPlayerGuideWidget.Instance.SaveGuideStepToSever((int)eGuideStep.AddNewRoleInFormation);
            UnityEngine.Debug.LogError("跳过这一步引导！");
            CurrentStep = 10;
            DoNextStep();
            return;
        }
        //if (UITool.IsOpen("UIFightChoice"))
        //{
        //    UITool.Close<UIFightChoice>();
        //}
        //设置首次点击区域位置
        NewPlayerGuideWidget.Instance.SetPosition(goFirstPartnerClick.transform, "UICom_TopRes_FirstPartner");
        //UIManager.Instance.CloseAllExceptHall
        DoNextStep();

        isFirst = true;
    }

    //此处强制关闭其他UI
    void Update()
    {
        if (isFirst)
        {
            var list = UIManager.Instance.GetAllOpenedUIByLayer(UILayer.Normal);
            if (list.Count > 0 && list[0].uiName != "UIHall")
            {
                UIManager.Instance.CloseNormal(list[0].uiName);
            }
        }
    }

    /// <summary>
    /// 用于断线重连
    /// </summary>
    public override void ReEnter()
    {

    }

    private void AddAllListeners()
    {
        goFirstPartnerClick.onClick.AddListener(OnFirstPartnerClick);
        goOpenParternerChooseUIClick.onClick.AddListener(OnOpenChooseParternerUIClick);
        goChoicePartnerClick.onClick.AddListener(OnChoicePartnerClick);
        goEnterPartnerClick.onClick.AddListener(OnEnterPartnerClick);
        btnCloseUIRoleClick.onClick.AddListener(OnCloseUIRoleClick);
        //btnGoMainBattleClick.onClick.AddListener(OnGoMainBattleClick);
    }

    private void RemoveAllListeners()
    {
        goFirstPartnerClick.onClick.RemoveAllListeners();
        goOpenParternerChooseUIClick.onClick.RemoveAllListeners();
        goChoicePartnerClick.onClick.RemoveAllListeners();
        goEnterPartnerClick.onClick.RemoveAllListeners();
        btnCloseUIRoleClick.onClick.RemoveAllListeners();
        //btnGoMainBattleClick.onClick.RemoveAllListeners();

        UniGameEvent.RemoveEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //移除UI打开监听
        //UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitClosePartnerSetUI);
    }

    public override void Exit()
    {
        CancelInvoke();
        RemoveAllListeners();
        gameObject.SetActive(false);

        UITool.Close<UINewPlayerGuide>();

        //if(hasFreePartner)
        NewPlayerGuideWidget.Instance.ShowGuide(eGuideStep.ClocseItemBoxAndGoToMainFightAgain);
    }


    //第一步骤，点击第一个伙伴头像
    private void OnFirstPartnerClick()
    {
        isFirst = false;
        Steps[CurrentStep]?.SetActive(false);

        UniGameEvent.AddEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //开始UI打开监听

        NewPlayerGuideWidget.Instance.InvokeButtonEvent("UICom_TopRes_FirstPartner");
    }

    private void WaitOpenUIFightChoice(GameEvent.UIPanel_Operate opened)
    {
        if (string.IsNullOrEmpty(opened.uiName)) return;

        if (opened.uiName.ToLower() == "UIRole".ToLower() && opened.isOpen)
        {
            UniGameEvent.RemoveEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //移除UI打开监听

            DoNextStep();

            FlyButton(goOpenParternerChooseUIClick.transform, goFirstPartnerClick.transform.position, "UIRole_SetPartner",0.5f);
        }
        else if (opened.uiName.ToLower() == "UIPartnerUp".ToLower() && opened.isOpen)
        {
            Invoke("ChoosePartner", 0.2f);
        }
    }

    private void ChoosePartner()
    {
        UniGameEvent.RemoveEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //移除UI打开监听

        DoNextStep();

        FlyButton(goChoicePartnerClick.transform, goOpenParternerChooseUIClick.transform.position, "UIPartnerUp_FirstRolePos", 0.5f);
    }

    //第二步，打开partner选择UI
    private void OnOpenChooseParternerUIClick()
    {
        Steps[CurrentStep]?.SetActive(false);//进入游戏后这个UI会被强制关闭。
        UniGameEvent.AddEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //开始UI打开监听
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("UIRole_SetPartner");
    }

    private void OnChoicePartnerClick()
    {
        if (CurrentStep != 2) return;
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("UIPartnerUp_FirstRole");
        goChoicePartnerClick.onClick.RemoveAllListeners();
        DoNextStep();

        FlyButton(goEnterPartnerClick.transform, goChoicePartnerClick.transform.position, "UIPartnerUp_EnterClick", 0.5f);
    }

    private void OnEnterPartnerClick()
    {
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("UIPartnerUp_EnterClick");
        
        //等待上阵UI关闭
        //UniGameEvent.AddEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitClosePartnerSetUI);//等待响应抽卡效果播放完成
        //Steps[CurrentStep]?.SetActive(false);
        NewPlayerGuideWidget.Instance.SaveGuideStepToSever((int)eGuideStep.AddNewRoleInFormation);
        DoNextStep();

        FlyButton(btnCloseUIRoleClick.transform, goEnterPartnerClick.transform.position, "UIRole_BtnClose", 0.5f);
    }

    //上阵UI关闭
    //private void WaitClosePartnerSetUI(GameEvent.Refresh_NewPlayerGuide guide)
    //{
    //    if (guide.args == null || guide.args.Length < 1) return;
    //    if (!guide.args[0].ToString().ToLower().Equals("ClosePartnerSetUI".ToLower())) return;

    //    UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitClosePartnerSetUI);

    //    DoNextStep();

    //    FlyButton(btnCloseUIRoleClick.transform, goEnterPartnerClick.transform.position, "UIRole_BtnClose", 0.5f);
    //}

    //关闭UIRole
    private void OnCloseUIRoleClick()
    {
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("UIRole_BtnClose");
        DoNextStep();
        //FlyButton(btnGoMainBattleClick.transform, btnCloseUIRoleClick.transform.position, "Hall_FightButton", 0.5f);
    }

    //引导再次打主线
    //private void OnGoMainBattleClick()
    //{
    //    NewPlayerGuideWidget.Instance.InvokeButtonEvent("Hall_FightButton");
    //    DoNextStep();
    //}
    //UIPartnerUp_FirstRole

    private bool existFreePartner()
    { 
        var allPartner = PartnerM.Instance.GetPartners();
        if (allPartner.Count > 4) return true;

        foreach (var item in allPartner)
        {
            if (!PartnerM.Instance.IsInSlot(item.Uid)) return true;
        }

        return false;
    }
}