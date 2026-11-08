using JEngine.Core;
using AppDomain = ILRuntime.Runtime.Enviorment.AppDomain;
using UnityEngine;

namespace JEngine.Helper
{
    public class MethodDelegateRegister : IRegisterHelper
    {
        public void Register(AppDomain appdomain)
        {
            /*  这里注册你所需要的MethodDelegate，参考报错的内容黏贴就好，举个例子：
             *
                appdomain.DelegateManager.RegisterMethodDelegate<object>();
             *
             */
            appdomain.DelegateManager.RegisterMethodDelegate<UnityEngine.EventSystems.PointerEventData>();
            appdomain.DelegateManager.RegisterMethodDelegate<UnityEngine.Vector2>();
            appdomain.DelegateManager.RegisterMethodDelegate<System.Int32, UnityEngine.GameObject>();
            appdomain.DelegateManager.RegisterMethodDelegate<System.Collections.IEnumerator>();
            appdomain.DelegateManager.RegisterMethodDelegate<System.Int32, UnityEngine.GameObject>();

            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.HTTPRequest, BestHTTP.HTTPResponse>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, System.String>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, BestHTTP.WebSocket.OnWebSocketOpenDelegate>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, BestHTTP.WebSocket.OnWebSocketBinaryDelegate>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, BestHTTP.WebSocket.OnWebSocketClosedDelegate>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, BestHTTP.WebSocket.OnWebSocketErrorDelegate>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, BestHTTP.WebSocket.OnWebSocketMessageDelegate>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, System.Byte[]>();
            appdomain.DelegateManager.RegisterMethodDelegate<BestHTTP.WebSocket.WebSocket, System.UInt16, System.String>();
            appdomain.DelegateManager.RegisterMethodDelegate<Spine.TrackEntry>();


        }
    }
}