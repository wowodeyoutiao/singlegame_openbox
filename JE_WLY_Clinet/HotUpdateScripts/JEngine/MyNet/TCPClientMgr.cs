using JEngine.Core;
using System;
using System.Collections.Generic;
using UnityEngine;

namespace JEngine.MyNet
{
    public enum enSocketType
    {
        SocketType_Login = 0, //登陆

    }

    public class TCPClientMgr : MonoBehaviour
    {
        public delegate void DelegateSocketError(WebClientGroup group);
        public delegate void DelegateSocketConnect(WebClientGroup group);

        public DelegateSocketError OnSocketHasError;
        public DelegateSocketConnect OnSocketConnectResult;

        private class _Clinet
        {
            public int _type;
            public WebClientGroup webClientGroup;
        }

        List<_Clinet> m_Clients = new List<_Clinet>();
        HashSet<int> m_ClientHasSet = new HashSet<int>();

        private static TCPClientMgr _instance;
        public static TCPClientMgr Instance
        {
            get
            {
                if (_instance == null)
                {
                    Log.Print(" TCPClientMgr Instance ");
                    GameObject go = new GameObject("TCPClientMgr");
                    _instance = go.AddComponent<TCPClientMgr>();
                    DontDestroyOnLoad(_instance);
                }
                return _instance;
            }
        }


        private void Awake()
        {

        }

        /// <summary>
        /// 添加TCP建立Group
        /// </summary>
        /// <param name="socketType"></param>
        WebClientGroup GetScoketGroup(int socketType)
        {
            if (!m_ClientHasSet.Contains(socketType))
            {
                var _clinet = new _Clinet { _type = socketType, webClientGroup = new WebClientGroup(socketType) };
                _clinet.webClientGroup.SetCbOnSocketError(OnSocketError);
                m_Clients.Add(_clinet);
                m_ClientHasSet.Add(socketType);
                return _clinet.webClientGroup;
            }
            for (int i = 0; i < m_Clients.Count; i++)
            {
                if (socketType == m_Clients[i]._type)
                    return m_Clients[i].webClientGroup;
            }
            return null;
        }

        /// <summary>
        /// 清除关闭所有的TCP
        /// </summary>
        public void ClearIpPool()
        {
            foreach (var kvp in m_Clients)
            {
                kvp.webClientGroup.ClearIpPool();
                kvp.webClientGroup.CloseAll();
            }
            OnSocketConnectResult = null;
            OnSocketHasError = null;
        }

        /// <summary>
        /// 添加TCP建立
        /// </summary>
        /// <param name="ip"></param>
        /// <param name="port"></param>
        /// <param name="socketType">目前只有一个可以不传</param>
        public void AddIpToPool(string ip, string port, int socketType = (int)enSocketType.SocketType_Login)
        {
            GetScoketGroup(socketType).AddIPToPool(ip, port);
        }

        public void ConnectSocket(int socketType = (int)enSocketType.SocketType_Login, int wNum = 1)
        {
            Debug.Log("开始链接Socket " + socketType);
            var _clients = GetScoketGroup(socketType);
            _clients.SetAliveCount(wNum);
            _clients.BeginConnect(OnSocketConnect);
        }

        public void CloseSocket(int SocketType = (int)enSocketType.SocketType_Login)
        {
            if (m_ClientHasSet.Contains(SocketType))
            {
                var _clients = GetScoketGroup(SocketType);
                _clients.CloseAll();
            }
        }

        void OnSocketError(WebClientGroup socket)
        {
            Debug.Log("Socket " + socket.socketType + " call error function.");
            if (OnSocketHasError != null)
            {
                OnSocketHasError(socket);
            }
        }

        void OnSocketConnect(WebClientGroup socket)
        {
            Debug.Log("Socket " + socket.socketType + " connect:" + socket.isConnect);
            if (OnSocketConnectResult != null)
            {
                OnSocketConnectResult(socket);
            }
        }

        /// <summary>
        /// 序列化并返回二进制
        /// Serialize an Object and return byte[]
        /// </summary>
        /// <param name="obj"></param>
        /// <returns></returns>
        public static byte[] ProtoSerialize<T>(T obj) where T : class
        {
            try
            {
                using (var stream = new System.IO.MemoryStream())
                {
                    ProtoBuf.Serializer.Serialize(stream, obj);
                    return stream.ToArray();
                }
            }
            catch (Exception ex)
            {
                if (ex.InnerException != null) ex = ex.InnerException;
                Log.PrintError($"[StringifyHelper] 错误：{ex.Message}, {ex.Data["StackTrace"]}");
                return null;
            }
        }

        /// <summary>
        /// 通过二进制反序列化
        /// Deserialize with byte[]
        /// </summary>
        /// <param name="msg"></param>
        /// <typeparam name="T"></typeparam>
        /// <returns></returns>
        public static T ProtoDeSerialize<T>(byte[] msg) where T : class
        {
            try
            {
                return ProtoBuf.Serializer.Deserialize(typeof(T), new System.IO.MemoryStream(msg)) as T;
            }
            catch (Exception ex)
            {
                if (ex.InnerException != null) ex = ex.InnerException;
                Log.PrintError($"[StringifyHelper] 错误：{ex.Message}, {ex.Data["StackTrace"]}");
                return null;
            }
        }

        /// <summary>
        /// 协议 二进制反序列化广播
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="msg"></param>
        public static void ProtoDeSerializePost<T>(byte[] msg) where T : class
        {
            try
            {
                T post = ProtoBuf.Serializer.Deserialize(typeof(T), new System.IO.MemoryStream(msg)) as T;
                //JEvent.defaultEvent.Post(post);
                UniGameEvent.Send(post);
                post = null;
            }
            catch (Exception ex)
            {
                if (ex.InnerException != null) ex = ex.InnerException;
                Log.PrintError($"[ProtoDeSerializePost StringifyHelper] 错误：{ex.Message}, {ex.Data["StackTrace"]}");
            }
        }

        public void SendMessage(short mId, MType mType, byte[] packet, Action<byte[]> _type = null)
        {
            if (packet == null)
                packet = new byte[0];
            Message _message = new Message(0, mType, packet);
            _message.compressed = true;
            SendMessage(_message, packet, mId, _type);
        }

        void SendMessage(Message message, byte[] protos, short mid, Action<byte[]> _type, int SocketType = (int)enSocketType.SocketType_Login)
        {
            bool withLoading = (protos != null && protos.Length > 0);

            if (m_ClientHasSet.Contains(SocketType))
            {
                var client = GetScoketGroup(SocketType);
                if (client.isConnect)
                {
                    if (withLoading)
                    {
                        //UINetLoading.StartLoading(10, "");
                    }
                    client.SendMsg(message, mid, _type);
                }
                else
                {
                    if (OnSocketHasError != null)
                    {
                        OnSocketHasError(client);
                    }
                }
            }
        }

        public void Update()
        {
            for (int i = 0; i < m_Clients.Count; i++)
            {
                m_Clients[i].webClientGroup.Update();
            }
        }

        /// <summary>
        /// 判断Socket 链接状态
        /// </summary>
        /// <param name="socketType"></param>
        /// <returns></returns>
        public bool SocketConnected(int socketType = (int)enSocketType.SocketType_Login)
        {
            if (m_ClientHasSet.Contains(socketType))
            {
                return GetScoketGroup(socketType).isConnect;
            }
            return false;
        }


        public void OnDestroy()
        {
            Debug.Log("Tcp OnDestroy");
            _instance = null;
            OnSocketHasError = null;
            OnSocketConnectResult = null;
            foreach (var client in m_Clients)
            {
                client.webClientGroup.CloseAll();
            }
            m_Clients.Clear();
            m_ClientHasSet.Clear();
        }
    }
}
