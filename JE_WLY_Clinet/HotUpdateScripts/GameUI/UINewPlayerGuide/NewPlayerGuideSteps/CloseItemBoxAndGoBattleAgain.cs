using DG.Tweening;
using MetaFramework.UI;
using System;
using System.Diagnostics;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 引导新角色上阵后执行非强制引导部分
/// </summary>
public class CloseItemBoxAndGoBattleAgain : NewPlayerGuideStepBase
{
    public Button btnWear;
    public Button btnRecycle;
    public Button btnGoFight;

    public override void Enter(UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide> parentController)
    {
        base.Enter(parentController);

        RemoveAllListeners();
        AddAllListeners();

        var tfShowBox = NewPlayerGuideWidget.Instance.GetGuideObject<Transform>("Hall_BoxShowUI");
        var tfWear = NewPlayerGuideWidget.Instance.GetGuideObject<Transform>("Hall_Btn_Wear");
        if (tfShowBox == null || !tfShowBox.gameObject.activeSelf) //没有找到ShowBox或者没有弹奖励面板
        {
            //点击继续主线战斗
            this.CurrentStep = 1;
            DoNextStep();
        }
        else if (tfWear != null && tfWear.gameObject.activeSelf)
        {
            //点击装备
            NewPlayerGuideWidget.Instance.SetPosition(btnWear.transform, "Hall_Btn_Wear");
            //UIManager.Instance.CloseAllExceptHall
            DoNextStep();
        }
        else
        {
            NewPlayerGuideWidget.Instance.SetPosition(btnWear.transform, "Hall_Btn_Recycle");
            this.CurrentStep = 0; //点击回收
            DoNextStep();
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
        btnWear.onClick.AddListener(OnWearClick);
        btnRecycle.onClick.AddListener(OnRecycleClick);
        btnGoFight.onClick.AddListener(OnGoFightClick);
    }

    private void RemoveAllListeners()
    {
        btnWear.onClick.RemoveAllListeners();
        btnRecycle.onClick.RemoveAllListeners();
        btnGoFight.onClick.RemoveAllListeners();

        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitRecycleFinish); //移除监听
    }

    public override void Exit()
    {
        RemoveAllListeners();
        gameObject.SetActive(false);

        UITool.Close<UINewPlayerGuide>();
    }


    //第一步骤，点击第一个伙伴头像
    private void OnWearClick()
    {
        Steps[CurrentStep]?.SetActive(false);

        UniGameEvent.AddEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitRecycleFinish); //开始监听

        NewPlayerGuideWidget.Instance.InvokeButtonEvent("Hall_Btn_Wear");
    }

    private void WaitRecycleFinish(GameEvent.Refresh_NewPlayerGuide guide)
    {
        if (guide.args == null || guide.args.Length < 1) return;
        if (!guide.args[0].ToString().ToLower().Equals("RecycleFinish".ToLower())) return;

        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(WaitRecycleFinish); //移除监听

        DoNextStep();

        FlyButton(btnRecycle.transform, btnWear.transform.position, "Hall_Btn_Recycle", 0.5f);
    }

    //第二步，打开partner选择UI
    private void OnRecycleClick()
    {        
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("Hall_Btn_Recycle");

        DoNextStep();

        FlyButton(btnGoFight.transform, btnRecycle.transform.position, "Hall_FightButton", 0.5f);
    }
    
    //引导再次打主线
    private void OnGoFightClick()
    {
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("Hall_FightButton");
        DoNextStep();
    }
}