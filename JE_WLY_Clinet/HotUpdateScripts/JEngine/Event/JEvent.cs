//
// JEvent.cs
//
// Author:
//       JasonXuDeveloper（傑） <jasonxudeveloper@gmail.com>
//
// Copyright (c) 2020 JEngine
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
using System;
using System.Linq;
using System.Reflection;
using System.Collections;
using System.Collections.Generic;
using JEngine.Core;
using Cysharp.Threading.Tasks;
using System.Linq.Expressions;
using System.Collections.Concurrent;

namespace JEngine.Event
{
    public class JEvent
    {
        private struct EventHandler
        {
            public object Subscriber; // 添加订阅者引用
            public Action<object, object> Action;
            public ThreadMode ThreadMode;
        }

        // 两级字典存储：[参数类型][订阅者] = 处理列表
        private Dictionary<Type, Dictionary<object, List<EventHandler>>> _handlersByType =
            new Dictionary<Type, Dictionary<object, List<EventHandler>>>();

        // 订阅者状态跟踪
        private HashSet<object> _unsubscribed = new HashSet<object>();
        private static JEvent _instance;

        // 类型缓存：存储每个类型的所有基类+接口
        private static Dictionary<Type, List<Type>> _assignableTypesCache =
            new Dictionary<Type, List<Type>>();

        // 添加对象池
        private static readonly ObjectPool<List<InvocationItem>> s_ListPool =
            new ObjectPool<List<InvocationItem>>(
                () => new List<InvocationItem>(32),
                list => list.Clear()
            );

        // 调用项结构体（避免元组GC）
        private struct InvocationItem
        {
            public object Subscriber;
            public EventHandler Handler;
        }


        public static bool ShowLog = false;
        public static JEvent defaultEvent => _instance ??= new JEvent();

        /// <summary>
        /// 单参数广播
        /// </summary>
        public void Post(object parameter)
        {
            if (parameter == null)
            {
                if (ShowLog) Log.Print("[JEvent] <color=#ffa73>参数为null</color>");
                return;
            }

            Type paramType = parameter.GetType();
            var types = GetAssignableTypes(paramType);

            // 从对象池获取调用列表
            var invocations = s_ListPool.Get();
            try
            {
                // 第一阶段：收集需要调用的处理程序
                foreach (var type in types)
                {
                    if (!_handlersByType.ContainsKey(type))
                        continue;

                    // 优化点：直接遍历值集合减少字典查找
                    foreach (var subscriberList in _handlersByType[type].Values)
                    {
                        for (int i = 0; i < subscriberList.Count; i++)
                        {
                            var handler = subscriberList[i];
                            // 检查订阅者状态
                            if (!_unsubscribed.Contains(handler.Subscriber))
                            {
                                invocations.Add(new InvocationItem
                                {
                                    Subscriber = handler.Subscriber,
                                    Handler = handler
                                });
                            }
                        }
                    }
                }

                // 第二阶段：执行调用
                for (int i = 0; i < invocations.Count; i++)
                {
                    var item = invocations[i];
                    item.Handler.Action(item.Subscriber, parameter);
                }
            }
            finally
            {
                // 归还对象池
                s_ListPool.Release(invocations);
            }
        }

        // 获取类型的所有可分配类型（自身+基类+接口）
        private List<Type> GetAssignableTypes(Type type)
        {
            if (_assignableTypesCache.ContainsKey(type))
                return _assignableTypesCache[type];

            var types = new List<Type> { type };

            // 添加所有基类
            Type baseType = type.BaseType;
            while (baseType != null && baseType != typeof(object))
            {
                types.Add(baseType);
                baseType = baseType.BaseType;
            }

            // 添加所有接口
            types.AddRange(type.GetInterfaces());

            _assignableTypesCache[type] = types;
            return types;
        }

        public void Unregister<T>(T val) where T : new()
        {
            if (val == null) return;

            _unsubscribed.Add(val);
            if (ShowLog) Log.Print($"[JEvent] {val.GetType().Name} 取消注册");
        }

        public void Register<T>(T val) where T : new()
        {
            if (val == null) return;

            if (_unsubscribed.Contains(val))
                _unsubscribed.Remove(val);

            var type = typeof(T);
            if (IsAlreadyRegistered(val))
            {
                if (ShowLog) Log.Print($"[JEvent] {type.Name} 已注册");
                return;
            }

            RegisterType(type, val);
        }

        private bool IsAlreadyRegistered(object val)
        {
            foreach (var handlers in _handlersByType.Values)
            {
                if (handlers.ContainsKey(val)) return true;
            }
            return false;
        }

        private void RegisterType(Type type, object val)
        {
            bool classLevel = type.GetCustomAttributes(typeof(SubscriberAttribute), false).Length > 0;
            RegisterMethods(type, val, classLevel);

            //Type baseType = type.BaseType;
            //while (baseType != null && baseType != typeof(object))
            //{
            //    if (classLevel)
            //    {
            //        classLevel = baseType.GetCustomAttributes(typeof(SubscriberAttribute), false).Length > 0;
            //        RegisterMethods(baseType, val, classLevel);
            //    }
            //    baseType = baseType.BaseType;
            //}
        }

        private void RegisterMethods(Type type, object val, bool classLevel)
        {
            const BindingFlags flags = BindingFlags.Default | BindingFlags.Instance |
                                     BindingFlags.Static | BindingFlags.Public |
                                     BindingFlags.NonPublic | BindingFlags.FlattenHierarchy;

            foreach (var method in type.GetMethods(flags))
            {
                var methodAttr = method.GetCustomAttributes(typeof(SubscriberAttribute), false);
                if (!classLevel && methodAttr.Length == 0) continue;

                if (method.ContainsGenericParameters || method.IsAbstract) continue;

                var parameters = method.GetParameters();
                if (parameters.Length != 1) continue;

                Type paramType = parameters[0].ParameterType;
                var handler = new EventHandler
                {
                    Subscriber = val, // 存储订阅者
                    ThreadMode = ThreadMode.Main,
                    Action = (target, arg) => method.Invoke(target, new object[] { arg })
                };

                // 添加到存储结构
                if (!_handlersByType.TryGetValue(paramType, out var subscribers))
                {
                    subscribers = new Dictionary<object, List<EventHandler>>();
                    _handlersByType[paramType] = subscribers;
                }

                if (!subscribers.TryGetValue(val, out var handlers))
                {
                    handlers = new List<EventHandler>();
                    subscribers[val] = handlers;
                }

                handlers.Add(handler);

                if (ShowLog)
                    Log.Print($"[JEvent] 注册 {val.GetType().Name}.{method.Name} 参数: {paramType.Name}");
            }
        }
    }

    public class ObjectPool<T> where T : class
    {
        private readonly Stack<T> m_Stack = new Stack<T>();
        private readonly Func<T> m_Creator;
        private readonly Action<T> m_OnRelease;

        public ObjectPool(Func<T> creator, Action<T> onRelease = null)
        {
            m_Creator = creator;
            m_OnRelease = onRelease;
        }

        public T Get()
        {
            return m_Stack.Count == 0 ? m_Creator() : m_Stack.Pop();
        }

        public void Release(T element)
        {
            m_OnRelease?.Invoke(element);
            m_Stack.Push(element);
        }
    }
}
