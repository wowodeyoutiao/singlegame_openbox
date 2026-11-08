using System;
using UnityEngine;

public class BattleTimerWidget<TStage, TDevice> : CommonTimerWidgetBase<TDevice> where TDevice : BattleTimerWidget<TStage, TDevice>
{
    #region Fsm Actuator Widget

    #region Fsm Widget
    public static MiniFsm<TStage> CurFsmWidget { get; private set; } = new MiniFsm<TStage>();
    #endregion Fsm Widget

    #region ReqTransitionStage

    public static void ReqTransitionStage(TStage stage) => CurFsmWidget.ReqTransitionStage(stage);

    #endregion ReqTransitionStage

    #region Recovery

    public static void Recovery() => CurFsmWidget.Recovery();

    #endregion Recovery

    #region AssemblyProcedure

    protected static void AssemblyFsmAction(TStage stage, Action _OnEnter, Action<float, float> _OnTick, Action _OnExit) => CurFsmWidget.CreateState(stage, _OnEnter, _OnTick, _OnExit);

    #endregion AssemblyProcedure

    #region

    protected static void TickFsm() => CurFsmWidget.DoTickTask(Time.deltaTime, Time.smoothDeltaTime);

    #endregion

    #endregion

    #region Stage Token
    public static TStage StageToken { get => CurFsmWidget.CurWorkState.StageToken; }
    #endregion

    #region BootUp/ShutOff

    protected static void StartFsmWidget(string deviceName = null, Transform dock = null)
    {
        Start(deviceName, dock);
        Recovery();
    }

    protected static void StopFsmWidget()
    {
        Recovery();
        Stop();
    }

    #endregion
}
