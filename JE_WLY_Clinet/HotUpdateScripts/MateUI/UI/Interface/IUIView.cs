using UnityEngine;

namespace MetaFramework.UI
{
    public interface IUIView
    {
        /// <summary>
        /// 用于Class Bind 时候赋值
        /// </summary>
        void Awake();

        /// <summary>
        /// 可自己手写获取组件
        /// </summary>
        /// <param name="go"></param>
        void Init(GameObject go);
    }
}