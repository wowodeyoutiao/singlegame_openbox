using pb;
using System;
using System.Collections.Generic;
using UnityEngine;

namespace JEngine.MyNet
{
    /// <summary>
    /// Socket组 链接 会同时保持n条线
    /// </summary>
    public class WebClientGroup
    {
        public class IPInfo
        {
            public int index;
            public string ip;
            public string port;

            public bool ipUsed = false; //这个ip是否已经链接

            public IPInfo(int _index, string _ip, string _port)
            {
                index = _index;
                ip = _ip;
                port = _port;
            }
        }

        public enum SocketGroupStatus
        {
            WaitForAwake, //等待启用
            StartConnect, //开始连接(n条线同时开始,有一条线连上就切换到正常)
            NormalUse,    //正常使用(正常使用的时候会开始自动补连)
            TimeOut,      //连接超时(不会再有socket异常,只有超时连不上(补连失败))
        }

        public int socketType { get { return m_socketType; } }
        public SocketGroupStatus status { get { return m_status; } }
        public bool isConnect
        {
            get
            {
                bool returnValue = m_socketConnected.Count != 0;
                for (int i = 0; i < m_socketConnected.Count; i++)
                {
                    returnValue &= m_socketConnected[i].IsConnected;
                }
                return returnValue;
            }
        }

        int m_socketType;//网络类型

        SocketGroupStatus m_status; //socket组当前状态

        List<IPInfo> m_IPPool = new List<IPInfo>(); //ip池

        int m_IpIndexNow = 0; //当前ip下标

        List<TcpClientUnit> m_socketConnected = new List<TcpClientUnit>(); //已经成功链接的socket

        List<TcpClientUnit> m_socketConnecting = new List<TcpClientUnit>(); //正在链接的socket

        int m_aliveCount = 1; //连接数

        Action<WebClientGroup> m_cbOnConnectResult;
        Action<WebClientGroup> m_cbOnSocketError;


        //超时相关
        bool m_checkConnectTimeOut = false;
        float m_maxTimeOut = 5; //超时时间
        float m_connectTime = 0;

        //收包相关
        //long m_lastPacketTimestamp = -1; //最后一个收包时间戳
        //收包
        private System.Object m_receiveObject = new object();
        private LinkedList<Message> m_ReceivePackList = new LinkedList<Message>();

        //发包相关
        Packet m_msgSending;
        //int m_sendThreadCount; //发送线程数(所有线程报发送失败才算发送失败)
        private System.Object m_lockObject = new object();

        bool m_NeedHeatBeat = true;

        public WebClientGroup(int socketType)
        {
            m_socketType = socketType;
        }

        public void ClearIpPool()
        {
            m_IPPool.Clear();
        }

        public void AddIPToPool(string ip, string port)
        {
            int index = m_IPPool.Count;
            m_IPPool.Add(new IPInfo(index, ip, port));
        }

        public void SetAliveCount(int count)
        {
            m_aliveCount = count;
        }

        public void SetCbOnSocketError(Action<WebClientGroup> onSocketError)
        {
            m_cbOnSocketError = onSocketError;
        }

        //启用连接
        public void BeginConnect(Action<WebClientGroup> onSocketResult)
        {
            if (m_status == SocketGroupStatus.NormalUse)
            {
                onSocketResult(this);
                return;
            }

            m_cbOnConnectResult += onSocketResult;
            if (m_status == SocketGroupStatus.StartConnect)
            {
                return;
            }
            m_status = SocketGroupStatus.StartConnect;
            CheckConnectCount();
        }


        public void SendMsg(Message msg, short id, Action<byte[]> _type = null)
        {
            if (m_msgSending == null && isConnect)
            {
                _SendOneMsg(msg, id, _type);
            }
        }

        void _SendOneMsg(Message message, short id, Action<byte[]> _type = null)
        {
            for (int i = 0; i < m_socketConnected.Count; i++)
            {
                if (m_socketConnected[i].IsConnected)
                {
                    (byte[] byts, long sendId) = m_socketConnected[i].MessageEncode(message, id);


                    if (sendId > 0)
                    {
                        TempProtoMap.AddProto(sendId, id);
                    }

                    Packet packet = new Packet { Type = PType.Data, Length = byts.Length, Data = byts, sendId = sendId };
                    m_socketConnected[i].SendPacket(packet);
                }
            }

            //m_msgSending = packet;
            //Debug.Log("发送消息: " + packet);
        }


        /// <summary>
        /// Socket线程收包添加消息
        /// </summary>
        /// <param name="message"></param>
        void _OnReceiveMessage(Message message)
        {
            if (message != null)
            {
                lock (m_receiveObject)
                {
                    m_ReceivePackList.AddLast(message);
                }
            }
        }

        //主线程处理接收消息
        void OnMessageDispatch(Message message)
        {
            if ((message.Type == MType.Response || message.Type == MType.Request) && (TempProtoMap.CheckID(message.ID)))
            {
                //TODO 
                TempProtoMap.GetProro(message.ID, message.Payload);
            }
            else
            {
                ProtoMap.GetProro(short.Parse(message.Route), message.Payload);
            }
            //Debug.Log($"处理消息:ID {message.ID} Route {message.Route} ");

            //Debug.Log("------------------------------主线程处理接收消息-------------------------------------------------------");

        }

        public void Update()
        {

            CheckMsgDispatch();

            if (m_status == SocketGroupStatus.WaitForAwake)
            {
                return;
            }

            CheckConnectTimeOut();
            if (m_status == SocketGroupStatus.TimeOut)
            {
                CloseAll();
                return;
            }
        }

        private Queue<Message> sendQueue = new Queue<Message>();
        /// <summary>
        /// 检测收到信息处理
        /// </summary>
        void CheckMsgDispatch()
        {
            //收包队列
            if (m_ReceivePackList != null && m_ReceivePackList.Count != 0)
            {

                lock (m_receiveObject)
                {

                    while (m_ReceivePackList.Count > 0)
                    {
                        sendQueue.Enqueue(m_ReceivePackList.First.Value);

                        m_ReceivePackList.RemoveFirst();
                    }
                }
            }

            int count = 0;
            //Debug.Log("主线程处理消息 " + sendQueue.Count);
            while (sendQueue.Count > 0 && count < 5)
            {
                OnMessageDispatch(sendQueue.Dequeue());
                count++;
            }
        }

        //超时检测
        void CheckConnectTimeOut()
        {
            if (m_checkConnectTimeOut)
            {
                m_connectTime += Time.deltaTime;
                if (m_connectTime >= m_maxTimeOut)
                {
                    m_checkConnectTimeOut = false;
                    m_connectTime = 0;
                    ConnectTimeOut();
                }
            }
        }

        //超时只有在开始连接的时候才会触发
        //服务器无法链接时的处理(10秒没有一个ip可以连上，当作服务器无法链接)
        void ConnectTimeOut()
        {
            SocketGroupStatus lastStatus = m_status;
            m_status = SocketGroupStatus.TimeOut;
            Debug.Log("ConnectTimeOut " + lastStatus);
            if (lastStatus == SocketGroupStatus.StartConnect)
            { //回调连接结果
                Debug.Log("服务器连接超时 " + m_socketType);
                Debug.Log(m_cbOnConnectResult);
                if (m_cbOnConnectResult != null)
                {
                    m_cbOnConnectResult(this);
                    m_cbOnConnectResult = null;
                }

            }
            if (lastStatus == SocketGroupStatus.NormalUse)
            { //回调连接异常
                m_cbOnSocketError(this);
            }

        }

        //连接数不足 补连
        void CheckConnectCount()
        {
            int countNow = m_socketConnected.Count + m_socketConnecting.Count;
            if (countNow < m_aliveCount)
            {
                for (int i = 0; i < m_aliveCount - countNow; i++)
                {
                    int nextIpindex = GetNextIp();
                    if (nextIpindex == -1)
                    {
                        Debug.Log("没有可用ip用于链接... " + m_socketType);
                        break;
                    }
                    TcpClientUnit client = new TcpClientUnit(m_socketType);
                    client.SetSocketErrorCallBack(OnSocketError);
                    client.SetCbOnReceiveProto(_OnReceiveMessage);
                    //client.SetCbOnSendProto(OnSendPacketResult);
                    m_socketConnecting.Add(client);
                    client.Connect(m_IPPool[nextIpindex], OnSocketConnectResult);
                }
            }
            if (!m_checkConnectTimeOut && m_socketConnected.Count == 0)
            { //所有线路都处于连接中...开始超时计算
                m_checkConnectTimeOut = true;
                m_connectTime = 0;
            }
        }

        //获取下一个可用ip
        int GetNextIp()
        {
            if (m_IPPool.Count == 0)
            {
                return -1;
            }
            int indexSave = m_IpIndexNow;
            while (true)
            {
                if (m_IPPool[m_IpIndexNow].ipUsed == false)
                {
                    return m_IpIndexNow;
                }
                else
                {
                    IncreaseIPIndex();
                    if (m_IpIndexNow == indexSave)
                    {
                        break;
                    }
                }
            }
            return -1;
        }

        void IncreaseIPIndex()
        {
            m_IpIndexNow += 1;
            if (m_IpIndexNow == m_IPPool.Count)
            {
                m_IpIndexNow = 0;
            }
        }

        //socket连接结果
        void OnSocketConnectResult(TcpClientUnit client, bool isConneted)
        {
            if (client.IsConnected)
            {
                Debug.Log("socket连接结果" + m_NeedHeatBeat);
                //client.sendHeatBeat = m_NeedHeatBeat;
                m_socketConnected.Add(client);
                if (m_checkConnectTimeOut)
                {
                    m_checkConnectTimeOut = false;
                    m_connectTime = 0;
                }
                if (m_status == SocketGroupStatus.StartConnect)
                {
                    Debug.Log("服务器连接成功 " + m_socketType);
                    m_status = SocketGroupStatus.NormalUse;
                    if (m_cbOnConnectResult != null)
                    {
                        m_cbOnConnectResult(this);
                        m_cbOnConnectResult = null;
                    }
                }
            }
            else
            {
                client.Cleanup();
                CheckConnectCount();
            }
            m_socketConnecting.Remove(client);
        }

        //socket异常 补连(短链接貌似应该忽略)
        //长连接改为 没有线路处于连接状态 通报与服务器断开连接
        void OnSocketError(TcpClientUnit client)
        {

            client.Cleanup();

            m_socketConnected.Remove(client);

            if (m_socketConnected.Count == 0)
            {
                Debug.Log("网络错误 连接数为0 需要关闭");
                CloseAll();
                m_cbOnSocketError(this);
            }
            else
            {
                CheckConnectCount();
            }

        }

        //发包结果[被子线程调用]
        void OnSendPacketResult(Packet msg, bool success)
        {
            if (msg == m_msgSending)
            {
                //			if (success) {
                //				Debug.Log ("消息 "+msg.packet.GetType () + " 在 "+m_socketType+" 中发送成功");
                //				m_msgSending = null;
                //			} else {
                //				m_msgSending.theadCount--;
                //				if (m_msgSending.theadCount == 0) {
                //					Debug.Log ("消息 "+msg.packet.GetType () + " 在 "+m_socketType+" 中发送失败");
                //					lock (m_lockObject) {
                //						m_msgWaitSend.Insert (0,m_msgSending);
                //					}
                //					m_msgSending = null;
                //				}
                //			}
            }
        }

        public void CloseAll()
        {
            Debug.Log("关闭连接 " + m_socketType);
            for (int i = 0; i < m_socketConnected.Count; i++)
            {
                m_socketConnected[i].Close();
            }
            for (int i = 0; i < m_socketConnecting.Count; i++)
            {
                m_socketConnecting[i].Close();
            }
            //m_lastPacketTimestamp = 0;
            m_socketConnected.Clear();
            m_socketConnecting.Clear();

            m_msgSending = null;

            m_status = SocketGroupStatus.WaitForAwake;

            m_cbOnConnectResult = null;

        }
    }
}
