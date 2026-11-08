using DG.Tweening;
using MetaFramework.UI;
using System;
using System.Diagnostics;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 在战役中的新手引导
/// </summary>
public class InFirstPlayBattleGuidance : NewPlayerGuideStepBase
{
    public Button goClick;
    public Button goAutoGameClick;
    public Button btnChooseCard;
    public Button btnCloseCardChooseClick;
    public Button btnBackHomeClick;

    public override void Enter(UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide> parentController)
    {
        base.Enter(parentController);

        RemoveAllListeners();

        AddAllListeners();

        ParentController.GetView().maskGo?.SetActive(false);//非强制

        CancelInvoke();
        //设置首次点击区域位置
        Invoke("DelayShowAuto", 4f);
        DoNextStep();
    }

    /// <summary>
    /// 用于断线重连
    /// </summary>
    public override void ReEnter()
    {

    }

    private void DelayShowAuto()
    {
        OnClickFight();
    }

    public override void Exit()
    {
        RemoveAllListeners();

        gameObject.SetActive(false);

        UITool.Close<UINewPlayerGuide>();
    }


    private void AddAllListeners()
    {
        goClick.onClick.AddListener(OnClickFight);
        goAutoGameClick.onClick.AddListener(OnAutoGameClick);
        //btnChooseCard.onClick.AddListener(OnChooseCardClick);
        btnCloseCardChooseClick.onClick.AddListener(OnCloseChooseCardClick);
        btnBackHomeClick.onClick.AddListener(OnBackHomeClick);
    }

    private void RemoveAllListeners()
    {
        goClick.onClick.RemoveAllListeners();
        goAutoGameClick.onClick.RemoveAllListeners();
        //btnChooseCard.onClick.RemoveAllListeners();
        btnCloseCardChooseClick.onClick.RemoveAllListeners();
        btnBackHomeClick.onClick.RemoveAllListeners();

        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitOpenFightResult);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitShowCardEffectOver);
    }


    //第一步骤，模拟遥感
    private void OnClickFight()
    {
        ParentController.GetView().maskGo?.SetActive(true);//强制
        CancelInvoke();
        NewPlayerGuideWidget.Instance.SetPosition(goAutoGameClick.transform, "UIGameDungeons_AutoPlay");
        DoNextStep();
    }

    //第二步，点击自动
    private void OnAutoGameClick()
    {
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("UIGameDungeons_AutoPlay");
        ParentController.GetView().maskGo?.SetActive(false);//放开强制引导
        //等待战斗结束,弹出结算
        Steps[CurrentStep]?.SetActive(false);
        UniGameEvent.AddEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitOpenFightResult);//等待打开战斗结果UI
    }

    //第三步 把手指放到卡上
    private void WaitOpenFightResult(GameEvent.Refresh_NewPlayerGuide guide)
    {
        if (guide.args == null || guide.args.Length < 1) return;
        if (!guide.args[0].ToString().ToLower().Equals("OpenGameOverAwardUI".ToLower())) return;

        ParentController.GetView().maskGo?.SetActive(false);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitOpenFightResult);
        NewPlayerGuideWidget.Instance.SetPosition(btnChooseCard.transform, "GameOverAward_Card");
        DoNextStep();

        UniGameEvent.AddEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitHideCardFinger);//
    }

    //隐藏手指并监听动画播放结束
    private void WaitHideCardFinger(GameEvent.Refresh_NewPlayerGuide guide)
    {
        if (guide.args == null || guide.args.Length < 1) return;
        if (!guide.args[0].ToString().ToLower().Equals("BeginShowCardPlayEffect".ToLower())) return;

        ParentController.GetView().maskGo?.SetActive(true);

        Steps[CurrentStep]?.SetActive(false);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitHideCardFinger);
        UniGameEvent.AddEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitShowCardEffectOver);//等待响应抽卡效果播放完成
    }

    //抽卡效果播放完成，把手指放到关闭位置
    private void WaitShowCardEffectOver(GameEvent.Refresh_NewPlayerGuide guide)
    {
        if (guide.args == null || guide.args.Length < 1) return;
        if (!guide.args[0].ToString().ToLower().Equals("EndShowCardPlayEffect".ToLower())) return;

        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitShowCardEffectOver);
        //NewPlayerGuideWidget.Instance.SetPosition(btnCloseCardChooseClick.transform, "GameOverAward_CloseBtnPos");
        FlyButton(btnCloseCardChooseClick.transform, btnChooseCard.transform.position, "GameOverAward_CloseBtnPos", 0.5f);
        DoNextStep();
    }

    private void OnCloseChooseCardClick()
    {
        NewPlayerGuideWidget.Instance.SaveGuideStepToSever((int)eGuideStep.InPlayBattleGuidance); //发送完成步骤

        NewPlayerGuideWidget.Instance.InvokeButtonEvent("GameOverAward_CloseBtn");
        DoNextStep();

        FlyButton(btnBackHomeClick.transform, btnCloseCardChooseClick.transform.position, "GameFightOver_ExitGame", 0.5f);
    }

    //点击返回
    private void OnBackHomeClick()
    {
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("GameFightOver_ExitGame");
        DoNextStep();
    }
}