using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class State<TStage>
{
    #region Token
    public TStage StageToken { get; private set; } = default(TStage);
    #endregion Token

    #region Event Action
    public Action OnEnterState { get; private set; } = null;
    public Action<float, float> OnTickState { get; private set; } = null;
    public Action OnExitState { get; private set; } = null;
    #endregion Event Action

    #region

    public State(TStage _Stage, Action _OnEnter, Action<float, float> _OnTick, Action _OnExit)
    {
        StageToken = _Stage;
        OnEnterState = _OnEnter;
        OnTickState = _OnTick;
        OnExitState = _OnExit;
    }

    #endregion Stage Event Pkg Define

    #region Dispose

    public void Dispose()
    {
        OnEnterState = null;
        OnTickState = null;
        OnExitState = null;
    }

    #endregion Dispose
}
