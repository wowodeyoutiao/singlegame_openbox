using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class MiniFsm<TStage>
{
    #region Data
    public Dictionary<TStage, State<TStage>> StateDic { get; private set; } = new Dictionary<TStage, State<TStage>>();
    public State<TStage> CurWorkState { get; private set; } = null;
    #endregion

    #region Method

    public void DoTickTask(float _DeltaTime, float _SmoothDeltaTime)
    {
        CurWorkState?.OnTickState?.Invoke(_DeltaTime, _SmoothDeltaTime);
    }

    public void ReqTransitionStage(TStage _TgtStage)
    {
        CurWorkState?.OnExitState?.Invoke();

        var _TgtPkg = GetState(_TgtStage);

        CurWorkState = _TgtPkg;

        CurWorkState?.OnEnterState?.Invoke();
    }

    public void Recovery() { ReqTransitionStage(default(TStage)); }

    #endregion

    #region Obtain

    public State<TStage> GetState(TStage _Stage)
    {
        State<TStage> _Pkg = null;
        if (StateDic.TryGetValue(_Stage, out _Pkg)) goto DoRetPkg;
        _Pkg = null;

    DoRetPkg:
        return _Pkg;
    }

    #endregion

    #region Equip Stage Pkg

    public void CreateState(TStage EState, Action _OnEnter, Action<float, float> _OnTick, Action _OnExit) 
    {
        StateDic[EState] = new State<TStage>(EState, _OnEnter, _OnTick, _OnExit);
    }

    #endregion

    #region Dispose

    public void Dispose()
    {
        foreach (var _Pkg in StateDic.Values) { _Pkg?.Dispose(); }
        StateDic?.Clear();
        StateDic = null;
        CurWorkState = null;
    }

    #endregion
}
