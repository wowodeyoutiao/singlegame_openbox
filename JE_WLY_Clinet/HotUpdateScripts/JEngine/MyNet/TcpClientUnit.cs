using BestHTTP.WebSocket;
using Cysharp.Threading.Tasks;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.LitJson;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
/// <summary>
/// WebSocket 客户端 (BestHTTP 适配版)
/// </summary>
using UnityEngine;

namespace JEngine.MyNet
{

    //type - package类型，1个byte，取值如下。
    //0x01: 客户端到服务器的握手请求以及服务器到客户端的握手响应
    //0x02: 客户端到服务器的握手ack
    //0x03: 心跳包
    //0x04: 数据包
    //0x05: 服务器主动断开连接通知
    public enum PType : byte
    {
        Handshake = 0x01,
        HandshakeAck = 0x02,
        Heartbeat = 0x03,
        Data = 0x04,
        Kick = 0x05
    }
    public class Packet
    {
        public PType Type;
        public int Length;
        public byte[] Data;
        public long sendId;
    }
    public enum MType : byte
    {
        Request = 0x00,
        Notify = 0x01,
        Response = 0x02,
        Push = 0x03,
    }

    public class Message
    {
        public MType Type; //类型
        public long ID; //类型（request Id）
        public string Route;// 当有压缩为 false 时候 为字符串协议号（目前使用int协议号，暂不使用）
        public byte[] Payload;// 数据信息
        public bool compressed; //压缩 (目前默认为false)
        public Message(int _id, MType _mtype, byte[] _Payload)
        {
            Type = _mtype;
            ID = _id;
            Payload = _Payload;
        }
        public Message()
        {
        }
    }

    public class HeartClass
    {
        public int code;
        public Sys sys;
    }

    public class Sys
    {
        public float heartbeat;
    }
    // 原有枚举和数据结构保持不变...
    // 保留 Packet、Message、HeartClass 等原有数据结构

    public class TcpClientUnit
    {
        private WebSocket webSocket;
        private string connectUrl;

        // 心跳相关
        private float m_heatBeatInterval = 2.5f;
        private int m_heatBeatTimeout = 5;
        private DateTime lastHeartbeatTime;

        // 回调函数
        private Action<Message> OnMessageReceived;
        private Action<TcpClientUnit, bool> OnConnectionChanged;
        private Action<TcpClientUnit> m_cbOnSocketError;

        // 发送队列
        private Queue<Packet> sendQueue = new Queue<Packet>();
        private bool isSending;

        // 消息解析缓存
        private List<byte> receiveBuffer = new List<byte>(1024);
        private int packetLength = 0;
        private PType currentType;

        public bool IsConnected => webSocket != null && webSocket.IsOpen;

        public void SetCbOnReceiveProto(Action<Message> cb)
        {
            OnMessageReceived = cb;
        }

        public void SetSocketErrorCallBack(Action<TcpClientUnit> errorCall)
        {
            m_cbOnSocketError = errorCall;
        }

        private int m_socketType;
        public TcpClientUnit(int socketType)
        {
            m_socketType = socketType;
        }

        public void Connect(WebClientGroup.IPInfo ipInfo, Action<TcpClientUnit, bool> callback)
        {
            if (IsConnected) return;

            connectUrl = GameM.Instance.HttpUrl.Contains("https") ? $"wss://{ipInfo.ip}" : $"ws://{ipInfo.ip}";
            if (!string.IsNullOrEmpty(ipInfo.port))
            {
                connectUrl += $":{ipInfo.port}";
            }
            Debug.Log("开始链接 " + connectUrl);
            ipInfo.ipUsed = true;
            webSocket = new WebSocket(new Uri(connectUrl));
            OnConnectionChanged = callback;

            // 设置微信小游戏特殊配置
#if UNITY_WECHAT && !UNITY_EDITOR
            webSocket.InternalRequest.DisableCache = true;
#endif
            webSocket.OnOpen += OnWebSocketOpen;
            webSocket.OnMessage += OnWebSocketMessageReceived;
            webSocket.OnBinary += OnWebSocketBinaryReceived;
            webSocket.OnClosed += OnWebSocketClosed;
            webSocket.OnError += OnWebSocketError;

            webSocket.Open();
        }

        private void OnWebSocketOpen(WebSocket ws)
        {
            Debug.Log("发送握手包 ");
            // 发送握手包
            SendHandshake();

        }

        private void OnWebSocketBinaryReceived(WebSocket ws, byte[] data)
        {
            //UnityEngine.Debug.Log("处理二进制消息 ");
            // 处理二进制消息
            //JEngine.Core.CoroutineMgr.Instance.StartCoroutine(ProcessReceivedData(data));
            ProcessReceivedData(data);
        }

        private void OnWebSocketMessageReceived(WebSocket ws, string message)
        {
            // 处理文本消息（如有需要）
        }

        private void OnWebSocketClosed(WebSocket ws, ushort code, string message)
        {
            Debug.LogError($"WebSocket Close: {code}  ");
            Cleanup();
            //OnConnectionChanged?.Invoke(this, false);

            if (m_cbOnSocketError != null)
            {
                Debug.Log("OnWebSocketClosed");
                m_cbOnSocketError(this);
            }
        }

        private void OnWebSocketError(WebSocket ws, string ex)
        {
            Debug.LogError($"WebSocket Error: {ex} ----");
            string errorMsg = string.Empty;
            //Debug.Log("An error occured: " + ws.InternalRequest.State);

            Cleanup();
            //OnConnectionChanged?.Invoke(this, false);
            if (m_cbOnSocketError != null)
            {
                m_cbOnSocketError(this);
            }
        }

        public class Handshark
        {
            public string version;// = "1.1.1",
            public string type;// = "js-websocket"}
        }

        private void SendHandshake()
        {
            var handshake = new Handshark
            {
                version = "1.1.1",
                type = "unity-websocket"
            };
            var json = JsonUtility.ToJson(handshake);
            var packet = new Packet
            {
                Type = PType.Handshake,
                Data = Encoding.UTF8.GetBytes(json),
                Length = json.Length
            };

            Debug.Log("packet data " + packet.Data);
            SendPacket(packet);
        }

        private void SendHandSharkPackAck()
        {
            Packet p = new Packet { Type = PType.HandshakeAck };
            Debug.Log("SendHandSharkPackAck-----");
            SendPacket(p);
            StartHeartbeat();
            OnConnectionChanged?.Invoke(this, true);
        }

        private CancellationTokenSource _heartbeatCts;
        private void StartHeartbeat()
        {
            lastHeartbeatTime = DateTime.Now;
            // 使用协程处理心跳
            //JEngine.Core.CoroutineMgr.Instance.StartCoroutine(HeartbeatRoutine());

            StopHeartbeat();                 // 保险起见先停一次
            _heartbeatCts = new CancellationTokenSource();
            HeartbeatLoopAsync(_heartbeatCts.Token).Forget();
        }

        public void StopHeartbeat()
        {
            _heartbeatCts?.Cancel();
            _heartbeatCts?.Dispose();
            _heartbeatCts = null;
        }

        private async UniTaskVoid HeartbeatLoopAsync(CancellationToken token)
        {
            while (IsConnected)                       // 连接中就一直跑
            {
                var now = DateTime.Now;
                if ((now - lastHeartbeatTime).TotalSeconds > m_heatBeatInterval)
                {
                    SendPacket(new Packet { Type = PType.Heartbeat });
                    lastHeartbeatTime = now;
                }

                // 每 0.5 秒检查一次，完全 0 GC
                await UniTask.Delay(
                    TimeSpan.FromSeconds(0.5f),
                    DelayType.DeltaTime,
                    PlayerLoopTiming.Update,
                    token);
            }
        }

        //private long HeartSendTime =0;
        //private System.Collections.IEnumerator HeartbeatRoutine()
        //{
        //    while (IsConnected)
        //    {
        //        if ((DateTime.Now - lastHeartbeatTime).TotalSeconds > m_heatBeatInterval)
        //        {
        //            SendPacket(new Packet { Type = PType.Heartbeat });
        //            lastHeartbeatTime = DateTime.Now;
        //        }
        //        yield return new WaitForSeconds(0.5f);
        //    }
        //}

        public void SendPacket(Packet packet)
        {
            if (!IsConnected) return;

            //byte[] data = PackageEncode(packet);
            sendQueue.Enqueue(packet);

            OnSendQueueProcessor();
        }

        private void OnSendQueueProcessor()
        {
            if (!isSending && sendedID <= 0)
            {
                SendQueueProcessor().Forget();
                //JEngine.Core.CoroutineMgr.Instance.StartCoroutine(SendQueueProcessor());
            }
        }
        //已经发送的自增ID
        private long sendedID = 0;
        private async UniTaskVoid SendQueueProcessor()
        {
            isSending = true;
            while (sendQueue.Count > 0 && IsConnected && sendedID <= 0)
            {
                Packet packet = sendQueue.Dequeue();
                webSocket.Send(PackageEncode(packet));

                packet = null;
                // 控制发送频率
                if (sendQueue.Count > 0)
                    await UniTask.Yield(PlayerLoopTiming.LastPostLateUpdate);
                //yield return new WaitForEndOfFrame() ;
            }
            isSending = false;
        }

        private void ProcessReceivedData(byte[] data)
        {
            receiveBuffer.AddRange(data);
            //Debug.Log("receiveBuffer +" + receiveBuffer.Count);
            while (receiveBuffer.Count >= 4)
            {
                if (packetLength == 0)
                {
                    // 解析包头
                    currentType = (PType)receiveBuffer[0];
                    byte[] lengthBytes = {
                        receiveBuffer[0],
                        receiveBuffer[1],
                        receiveBuffer[2],
                        receiveBuffer[3]
                    };
                    packetLength = HeadConvertBigEndianToDecimal(lengthBytes);// BitConverter.ToInt32(lengthBytes, 0);
                    //Debug.Log("packetLength +" + packetLength);
                }

                if (receiveBuffer.Count >= packetLength + 4)
                {
                    // 提取完整数据包
                    byte[] packetData = receiveBuffer.GetRange(4, packetLength).ToArray();
                    //Debug.Log("packetData +" +packetData.Length);
                    receiveBuffer.RemoveRange(0, packetLength + 4);

                    ProcessPacket(currentType, packetData);
                    packetLength = 0;
                }
                else
                {
                    break;
                }
            }
        }

        private GameEvent.Refresh_Ping _Refresh_Ping = new GameEvent.Refresh_Ping { ping = 0 };
        private void ProcessPacket(PType type, byte[] data)
        {
            //Debug.Log("type----" + type);
            switch (type)
            {
                case PType.Handshake:
                    HandleHandshake(data);
                    SendHandSharkPackAck();
                    break;
                case PType.HandshakeAck:
                    //目前不会收到
                    break;
                case PType.Heartbeat:
                    //long a = DateTime.Now.Ticks - HeartSendTime -
                    //_Refresh_Ping.ping =(int)(DateTime.Now - lastHeartbeatTime).TotalMilliseconds ;
                    lastHeartbeatTime = DateTime.Now;

                    JEvent.defaultEvent.Post(_Refresh_Ping);
                    break;
                case PType.Data:
                    DecodeMessage(data);
                    break;
                // 其他类型处理...

                case PType.Kick:
                    Debug.LogError("服务器断开链接");
                    break;
            }
        }

        /// <summary>
        /// 从packe 解析出message (当有mess ID 时，从本地的senid 中找对应协议，没有用Route 对应协议ID)
        /// </summary>
        /// <param name="bytes"></param>
        private void DecodeMessage(byte[] bytes)
        {

            if (bytes == null)
            {
                Debug.LogError("DecodeMessage  为空");
                return;
            }
            Message _mess = new Message();
            int offset = 1;
            byte flag = bytes[0];
            _mess.Type = (MType)((flag >> 1) & 0x07);

            if (_mess.Type == MType.Request || _mess.Type == MType.Response)
            {
                int messId = 0;
                for (int i = offset; i < bytes.Length; i++)
                {
                    byte b = bytes[i];
                    messId += (int)(b & 0x7F) << (int)(7 * (i - offset));
                    if (b < 128)
                    {
                        offset = i + 1;
                        break;
                    }
                }
                _mess.ID = messId;

                //if (messId == sendedID)
                //{
                //    Debug.Log("收到的 sendID" + messId);
                //    sendedID = 0; // 重置已发送ID
                //}
            }

            // Debug.Log(" _mess.compressed  "+ _mess.Type);
            if (_mess.Type == MType.Request || _mess.Type == MType.Notify || _mess.Type == MType.Push)
            {
                if ((flag & 0x01) == 1)
                {
                    _mess.compressed = true;
                    byte[] aa = new byte[2];
                    Array.Copy(bytes, offset, aa, 0, 2);
                    int code = ConvertBigEndianToDecimal(aa);
                    _mess.Route = code.ToString();
                    offset += 2;
                }
                else
                {

                    _mess.compressed = false;
                    byte rl = bytes[offset];
                    offset++;

                    _mess.Route = Encoding.UTF8.GetString(bytes, offset, (int)rl);
                    offset += (int)rl;
                }
            }

            // Array.Copy(bytes, offset, _mess.Payload, 0, bytes.Length - offset);
            _mess.Payload = bytes.Skip(offset).ToArray();

            //Debug.Log($"_mess.Type {_mess.Type} _mess.ID {_mess.ID}_mess.Route{_mess.Route}");
            if (OnMessageReceived != null)
            {
                OnMessageReceived(_mess);
            }

            //if (_mess.ID > 0 && sendedID == 0)
            //{
            //    // 发送队列处理
            //    OnSendQueueProcessor();
            //}
        }

        private static int HeadConvertBigEndianToDecimal(byte[] bytes)
        {

            if (bytes == null || bytes.Length == 0)
            {
                throw new ArgumentException("字节数组不能为空");
            }
            byte[] tBytes = new byte[bytes.Length - 1];
            Array.Copy(bytes, 1, tBytes, 0, bytes.Length - 1);

            int result = 0;
            for (int i = 0; i < tBytes.Length; i++)
            {
                result = (result << 8) | tBytes[i];
            }
            return result;
        }

        private static int ConvertBigEndianToDecimal(byte[] bytes)
        {

            if (bytes == null || bytes.Length == 0)
            {
                throw new ArgumentException("字节数组不能为空");
            }

            int result = 0;
            for (int i = 0; i < bytes.Length; i++)
            {
                result = (result << 8) | bytes[i];
            }
            return result;
        }

        public (byte[], long) MessageEncode(Message m, short id)
        {
            long sId = -1;
            byte[] bytes = MessageEncode(m, id, ref sId);
            return (bytes, sId);
        }


        private static protected long sendID = 0;
        /// <summary>
        /// message编译 
        /// </summary>
        /// <param name="m">转入message</param>
        /// <param name="id">messageId</param>
        /// <param name="_sendID"></param>
        /// <returns></returns>
        protected byte[] MessageEncode(Message m, short id, ref long _sendID)
        {
            int offset = 0;
            byte[] result = new byte[1024];
            result[offset] = (byte)((byte)m.Type << 1);
            result[offset] |= 1;
            offset++;

            //messageid 从1开始自增，  varints 编码。客户端 request/notify 需要对此id赋值。服务器response 会带上客户端的messageid。
            if (m.Type == MType.Request || m.Type == MType.Response)
            {
                sendID++;
                long n = sendID;
                while (sendID > 0)
                {
                    byte b = (byte)(n % 128);
                    n = n >> 7;
                    if (n > 0)
                    {
                        result[offset] = (byte)(b + 128);
                        offset++;
                    }
                    else
                    {
                        result[offset] = (byte)(b);
                        offset++;
                        break;
                    }
                }
                _sendID = sendID;
            }

            if (!m.compressed)
            {
                result[offset] = (byte)m.Route.Length;
                offset++;
                Buffer.BlockCopy(Encoding.UTF8.GetBytes(m.Route), 0, result, offset, m.Route.Length);
                offset += m.Route.Length;
            }
            else
            {
                //BigEndian
                result[offset] = (byte)((id >> 8) & 0xFF);
                offset++;
                result[offset] = (byte)(id & 0xFF);
                offset++;
            }

            Buffer.BlockCopy(m.Payload, 0, result, offset, m.Payload.Length);
            offset += m.Payload.Length;

            byte[] result2 = new byte[offset];
            Buffer.BlockCopy(result, 0, result2, 0, offset);
            return result2;
        }

        private void HandleHandshake(byte[] data)
        {
            var json = Encoding.UTF8.GetString(data);

            Debug.Log("----客户端到服务器的握手请求以及服务器到客户端的握手响应----" + System.Text.Encoding.UTF8.GetString(data));
            var response = JsonMapper.ToObject<HeartClass>(json);
            m_heatBeatInterval = response.sys.heartbeat / 2;
            m_heatBeatTimeout = (int)(m_heatBeatInterval * 2);
        }

        /// <summary>
        /// 第1个字节：消息类型
        /// 第2--4字节：消息长度，字节序：大端
        /// </summary>
        /// <param name="p"></param>
        /// <returns></returns>
        protected byte[] PackageEncode(Packet p)
        {
            byte[] result = new byte[p.Length + 4];
            result[0] = (byte)p.Type;
            if (p.Data != null)
            {
                byte[] b = BitConverter.GetBytes(p.Length);
                result[1] = b[2];
                result[2] = b[1];
                result[3] = b[0];

                Buffer.BlockCopy(p.Data, 0, result, 4, p.Length);
            }

            return result;
        }

        // 保留原有的编码/解码方法...
        //private byte[] PackageEncode(Packet packet)
        //{
        //    byte[] header = new byte[4];
        //    header[0] = (byte)packet.Type;
        //    byte[] lengthBytes = BitConverter.GetBytes(packet.Data.Length);
        //    Array.Copy(lengthBytes, 0, header, 1, 3);
        //    return CombineBytes(header, packet.Data);
        //}

        //private byte[] CombineBytes(params byte[][] arrays)
        //{
        //    int totalLength = 0;
        //    foreach (byte[] array in arrays)
        //        totalLength += array.Length;

        //    byte[] result = new byte[totalLength];
        //    int offset = 0;

        //    foreach (byte[] array in arrays)
        //    {
        //        Buffer.BlockCopy(array, 0, result, offset, array.Length);
        //        offset += array.Length;
        //    }
        //    return result;
        //}

        public void Cleanup()
        {
            if (webSocket != null)
            {
                webSocket.Close();
                webSocket = null;
            }
            sendQueue.Clear();
            receiveBuffer.Clear();
        }

        public void Close()
        {
            Cleanup();
            OnConnectionChanged = null;
            OnMessageReceived = null;
        }
    }
}