using HallBattle;
using System.Collections;
using UnityEngine;

namespace HallBattle
{
    /// <summary>
    /// 状态基类，实现IState接口的基础功能
    /// </summary>
    /// <typeparam name="T">拥有者类型</typeparam>
    public abstract class BaseState<T> : IState<T> where T : class
    {
        /// <summary>
        /// 状态拥有者
        /// </summary>
        public T Owner { get; set; }

        /// <summary>
        /// 退出状态的延时时间（秒）
        /// </summary>
        public float exitDelay = 0f;

        public BaseState(float exitDelay = 0f)
        {
            this.exitDelay = exitDelay;
        }

        /// <summary>
        /// 初始化状态
        /// </summary>
        public virtual void Initialize() { }

        /// <summary>
        /// 进入状态
        /// </summary>
        public virtual void Enter() { }

        /// <summary>
        /// 退出状态（带延时）
        /// </summary>
        public virtual IEnumerator Exit()
        {
            // 如果有退出延时，等待指定时间
            if (exitDelay > 0)
            {
                yield return new WaitForSeconds(exitDelay);
            }
        }

        /// <summary>
        /// 更新状态
        /// </summary>
        public virtual void Update(float dTime) { }

        /// <summary>
        /// 固定时间间隔更新
        /// </summary>
        public virtual void FixedUpdate() { }
    }
}