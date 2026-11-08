using JEngine.Core;
using System;
using AppDomain = ILRuntime.Runtime.Enviorment.AppDomain;

namespace JEngine.Helper
{
    public class DelegateConvertorRegister: IRegisterHelper
    {
        public void Register(AppDomain appdomain)
        {
            /*  这里注册你所需要的DelegateConvertor，参考报错的内容黏贴就好，举个例子：
             *
                appdomain.DelegateManager.RegisterDelegateConvertor<System.Predicate<System.String>>(act =>
                {
                    return new System.Predicate<System.String>(obj =>
                    {
                        return ((System.Func<System.String, System.Boolean>)act)(obj);
                    });
                });
             *
             */

            appdomain.DelegateManager.RegisterDelegateConvertor<global::EventTriggerListener.PointerEventDelegate>((act) =>
            {
                return new global::EventTriggerListener.PointerEventDelegate((eventData) =>
                {
                    ((Action<UnityEngine.EventSystems.PointerEventData>)act)(eventData);
                });
            });
            appdomain.DelegateManager.RegisterDelegateConvertor<UnityEngine.Events.UnityAction<UnityEngine.Vector2>>((act) =>
            {
                return new UnityEngine.Events.UnityAction<UnityEngine.Vector2>((arg0) =>
                {
                    ((Action<UnityEngine.Vector2>)act)(arg0);
                });
            });

            #region BestHTTP
            appdomain.DelegateManager.RegisterDelegateConvertor<BestHTTP.OnRequestFinishedDelegate>((act) =>
            {
                return new BestHTTP.OnRequestFinishedDelegate((originalRequest, response) =>
                {
                    ((Action<BestHTTP.HTTPRequest, BestHTTP.HTTPResponse>)act)(originalRequest, response);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<BestHTTP.WebSocket.OnWebSocketOpenDelegate>((act) =>
            {
                return new BestHTTP.WebSocket.OnWebSocketOpenDelegate((webSocket) =>
                {
                    ((Action<BestHTTP.WebSocket.WebSocket>)act)(webSocket);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<BestHTTP.WebSocket.OnWebSocketMessageDelegate>((act) =>
            {    
                return new BestHTTP.WebSocket.OnWebSocketMessageDelegate((webSocket, message) =>
                {
                    ((Action<BestHTTP.WebSocket.WebSocket, System.String>)act)(webSocket, message);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<BestHTTP.WebSocket.OnWebSocketBinaryDelegate>((act) =>
            {    
                return new BestHTTP.WebSocket.OnWebSocketBinaryDelegate((webSocket, data) =>   
                {
                    ((Action<BestHTTP.WebSocket.WebSocket, System.Byte[]>)act)(webSocket, data); 
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<BestHTTP.WebSocket.OnWebSocketClosedDelegate>((act) =>
            {   
                return new BestHTTP.WebSocket.OnWebSocketClosedDelegate((webSocket, code, message) =>   
                {
                    ((Action<BestHTTP.WebSocket.WebSocket, System.UInt16, System.String>)act)(webSocket, code, message);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<BestHTTP.WebSocket.OnWebSocketErrorDelegate>((act) =>
            {   return new BestHTTP.WebSocket.OnWebSocketErrorDelegate((webSocket, reason) =>    
                {
                    ((Action<BestHTTP.WebSocket.WebSocket, System.String>)act)(webSocket, reason);
                });
            });

            #endregion

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Comparison<ILRuntime.Runtime.Intepreter.ILTypeInstance>>((act) =>
            {
                return new System.Comparison<ILRuntime.Runtime.Intepreter.ILTypeInstance>((x, y) =>
                {
                    return ((Func<ILRuntime.Runtime.Intepreter.ILTypeInstance, ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int32>)act)(x, y);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>>((act) =>
            {
                return new UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>((arg0) =>
                {
                    ((Action<UnityEngine.EventSystems.PointerEventData>)act)(arg0);
                });
            });
            appdomain.DelegateManager.RegisterDelegateConvertor<DG.Tweening.TweenCallback>((act) =>
            {
                return new DG.Tweening.TweenCallback(() =>
                {
                    ((Action)act)();
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<DG.Tweening.Core.DOGetter<System.Int32>>((act) =>
            {
                return new DG.Tweening.Core.DOGetter<System.Int32>(() =>
                {
                    return ((Func<System.Int32>)act)();
                });
            });
            appdomain.DelegateManager.RegisterDelegateConvertor<DG.Tweening.Core.DOSetter<System.Int32>>((act) =>
            {
                return new DG.Tweening.Core.DOSetter<System.Int32>((pNewValue) =>
                {
                    ((Action<System.Int32>)act)(pNewValue);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Comparison<System.Int32>>((act) =>
            {
                return new System.Comparison<System.Int32>((x, y) =>
                {
                    return ((Func<System.Int32, System.Int32, System.Int32>)act)(x, y);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<Spine.AnimationState.TrackEntryDelegate>((act) =>
            {
                return new Spine.AnimationState.TrackEntryDelegate((trackEntry) =>
                {
                    ((Action<Spine.TrackEntry>)act)(trackEntry);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.String>>((act) =>
            {
                return new System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.String>((input) =>
                {
                    return ((Func<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.String>)act)(input);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int32>>((act) =>
            {
                return new System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int32>((input) =>
                {
                    return ((Func<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int32>)act)(input);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int64>>((act) =>
            {
                return new System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int64>((input) =>
                {
                    return ((Func<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int64>)act)(input);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Single>>((act) =>
            {
                return new System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Single>((input) =>
                {
                    return ((Func<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Single>)act)(input);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Double>>((act) =>
            {
                return new System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Double>((input) =>
                {
                    return ((Func<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Double>)act)(input);
                });
            });

            appdomain.DelegateManager.RegisterDelegateConvertor<System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, ILRuntime.Runtime.Intepreter.ILTypeInstance>>((act) =>
            {
                return new System.Converter<ILRuntime.Runtime.Intepreter.ILTypeInstance, ILRuntime.Runtime.Intepreter.ILTypeInstance>((input) =>
                {
                    return ((Func<ILRuntime.Runtime.Intepreter.ILTypeInstance, ILRuntime.Runtime.Intepreter.ILTypeInstance>)act)(input);
                });
            });

        }
    }
}