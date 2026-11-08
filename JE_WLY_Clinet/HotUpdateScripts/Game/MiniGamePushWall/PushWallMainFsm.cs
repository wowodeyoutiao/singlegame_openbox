using MetaFramework.UI;
using System;
using System.Collections;
using UnityEngine;

public enum PushWallMainStage
{
    WaitOp,
    Init
}

public class PushWallMainFsm :  BattleTimerWidget<PushWallMainStage, PushWallMainFsm>
{
    static void DoAssemblyFsmWidgetAction()
    {
        AssemblyFsmAction(PushWallMainStage.WaitOp, DoEnter_WaitOp, DoEnter_WaitOp, DoExit_WaitOp);
        AssemblyFsmAction(PushWallMainStage.Init, DoEnter_Init, DoTick_Init, () => { });

        ReqTransitionStage(PushWallMainStage.WaitOp);
    }

    int count;
    void Update()
    {
        TickFsm();

        if (count % 10 == 0)
        {
            count = 0;
            Debug.LogError(PushWallMainFsm.CurFsmWidget.CurWorkState.ToString());
        }
        
        count++;

    }

    void OnDestroy()
    {
        Debug.Log("==============");
    }

    #region WaitOp
    private static void DoEnter_WaitOp()
    {
        Debug.Log("==============");
    }
    private static void DoEnter_WaitOp(float arg1, float arg2)
    {
        //UITool.Open<UIPushWallGame>();
    }
    private static void DoExit_WaitOp()
    {

    }
    #endregion

    #region Init
    private static void DoEnter_Init()
    {

    }

    private static void DoTick_Init(float arg1, float arg2)
    {

    }
    #endregion

    #region BootUp/ShutOff
    public static void StartTimer()
    {
        DoAssemblyFsmWidgetAction();
        StartFsmWidget(@"PushWallMainCtr");
    }

    public static void ReqStartTimer()
    {
        if (PushWallMainFsm.StageToken != PushWallMainStage.WaitOp) return;

        ReqTransitionStage(PushWallMainStage.Init);
    }

    public static void EndTimer()
    {
        StopFsmWidget();
    }

    #endregion
}