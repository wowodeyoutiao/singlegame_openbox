using JEngine.Core;
using System;
using System.Collections;
using System.Reflection;

namespace HallBattle
{
    public class HallBattleStateMachine<T> where T : class
    {
        /// <summary>
        /// 当前状态
        /// </summary>
        public IState<T> CurrentState { get; private set; }

        /// <summary>
        /// 前一个状态
        /// </summary>
        public IState<T> PreviousState { get; private set; }

        /// <summary>
        /// 状态机拥有者
        /// </summary>
        public T Owner { get; private set; }

        /// <summary>
        /// 是否正在切换状态
        /// </summary>
        private bool isSwitchingState = false;

        public HallBattleStateMachine(T owner)
        {
            Owner = owner;
        }

        private void Awake()
        {
        }

        public void Update(float time)
        {
            if (CurrentState != null)
            {
                CurrentState.Update(time);
            }
        }

        public void FixedUpdate()
        {
            if (CurrentState != null)
            {
                CurrentState.FixedUpdate();
            }
        }

        /// <summary>
        /// 改变状态
        /// </summary>
        /// <typeparam name="TState">要切换到的状态类型</typeparam>
        /// <param name="exitDelay">退出当前状态的延时</param>
        public void ChangeState<TState>(float exitDelay = 0f) where TState : IState<T>, new()
        {
            if (isSwitchingState) return;

            CoroutineMgr.Instance.StartCoroutine(ChangeStateCoroutine<TState>(exitDelay));
        }

        /// <summary>
        /// 改变状态的协程（处理延时）
        /// </summary>
        private IEnumerator ChangeStateCoroutine<TState>(float exitDelay) where TState : IState<T>, new()
        {
            isSwitchingState = true;

            // 退出当前状态
            if (CurrentState != null)
            {
                // 如果指定了退出延时，覆盖状态中的默认延时
                if (exitDelay > 0 && CurrentState is BaseState<T> baseState)
                {
                    baseState.exitDelay = exitDelay;
                }

                yield return CoroutineMgr.Instance.StartCoroutine(CurrentState.Exit());
                PreviousState = CurrentState;
            }

            // 创建并进入新状态
            TState newState = new TState();
            newState.Owner = Owner;
            newState.Initialize();
            newState.Enter();

            CurrentState = newState;
            isSwitchingState = false;
        }

        /// <summary>
        /// 返回到上一个状态
        /// </summary>
        public void RevertToPreviousState(float exitDelay = 0f)
        {
            if (PreviousState != null)
            {
                Type previousStateType = PreviousState.GetType();
                MethodInfo changeStateMethod = GetType().GetMethod("ChangeState");
                MethodInfo genericMethod = changeStateMethod.MakeGenericMethod(previousStateType);
                genericMethod.Invoke(this, new object[] { exitDelay });
            }
        }

        /// <summary>
        /// 检查当前是否是指定状态
        /// </summary>
        public bool IsInState<TState>() where TState : IState<T>
        {
            return CurrentState is TState;
        }
    }
}