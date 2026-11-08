
using System.Collections;

namespace HallBattle
{
    /// <summary>
    /// 状态接口
    /// </summary>
    /// <typeparam name="T">拥有者类型</typeparam>
    public interface IState<T> where T : class
    {
        /// <summary>
        /// 状态拥有者
        /// </summary>
        T Owner { get; set; }

        /// <summary>
        /// 初始化状态
        /// </summary>
        void Initialize();

        /// <summary>
        /// 进入状态
        /// </summary>
        void Enter();

        /// <summary>
        /// 退出状态（返回协程用于处理延时）
        /// </summary>
        IEnumerator Exit();

        /// <summary>
        /// 更新状态
        /// </summary>
        void Update(float dTime);

        /// <summary>
        /// 固定时间间隔更新
        /// </summary>
        void FixedUpdate();
    }

}