using System;
using System.Collections.Concurrent;
using System.Net;
using System.Net.Sockets;
using System.Threading;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// TCP客戶端組件
    /// </summary>
    public class TCPClient : IClient
    {
        /// <summary>
        /// 接收處理器
        /// </summary>
        private class RecvHandler
        {
            /// <summary>
            /// 啟動處理
            /// </summary>
            /// <param name="stream">網路流物件</param>
            /// <param name="procmgr">訊息處理器</param>
            /// <param name="equeue">事件佇列</param>
            public void Start(NetworkStream stream, IProcmgr procmgr, EQueue equeue)
            {
                try
                {
                    if (thread != null)
                        throw new AlreadyStartException("recv handler");

                    thread = new Thread(() =>
                    {
                        while (true)
                        {
                            try
                            {
                                var byteOfSize = new byte[Define.headerSize];

                                if (stream.Read(byteOfSize, 0, Define.headerSize) == 0)
                                    throw new EofException();

                                var size = BitConverter.ToUInt16(byteOfSize, 0);

                                if (size <= 0)
                                    throw new PacketZeroException("recv");

                                if (size > Define.packetSize)
                                    throw new PacketLimitException("recv");

                                var byteOfData = new byte[size];

                                if (stream.Read(byteOfData, 0, size) == 0)
                                    throw new ReceiveException();

                                var message = procmgr.Decode(byteOfData);

                                equeue.Enqueue(EventID.Message, message);
                                equeue.Enqueue(EventID.Recv, null);
                            } // try
                            catch (Exception e)
                            {
                                equeue.Enqueue(EventID.Error, e);
                                equeue.Enqueue(EventID.Disconnect, null);
                                return;
                            } // catch
                        } // while
                    });
                    thread.IsBackground = true;
                    thread.Start();
                } // try
                catch (Exception e)
                {
                    equeue.Enqueue(EventID.Error, e);
                } // catch
            }

            /// <summary>
            /// 關閉處理
            /// </summary>
            public void Close()
            {
                thread?.Join();
                thread = null;
            }

            /// <summary>
            /// 執行緒物件
            /// </summary>
            private Thread thread = null;
        }

        /// <summary>
        /// 傳送處理器
        /// </summary>
        private class SendHandler
        {
            /// <summary>
            /// 啟動處理
            /// </summary>
            /// <param name="stream">網路流物件</param>
            /// <param name="procmgr">訊息處理器</param>
            /// <param name="equeue">事件佇列</param>
            public void Start(NetworkStream stream, IProcmgr procmgr, EQueue equeue)
            {
                try
                {
                    if (thread != null)
                        throw new AlreadyStartException("send handler");

                    queue = new BlockingCollection<object>();
                    thread = new Thread(() =>
                    {
                        while (true)
                        {
                            try
                            {
                                var message = queue.Take();

                                if (message == null)
                                    continue;

                                var packet = procmgr.Encode(message);
                                var size = packet.Length;

                                if (size <= 0)
                                    throw new PacketZeroException("send");

                                if (size > Define.packetSize)
                                    throw new PacketLimitException("send");

                                var byteOfSize = BitConverter.GetBytes(size);

                                stream.Write(byteOfSize, 0, Define.headerSize);
                                stream.Write(packet, 0, size);
                                stream.Flush();
                                equeue.Enqueue(EventID.Send, null);
                            } // try
                            catch (InvalidOperationException) // 這是因為關閉處理的CompleteAdding引發的, 所以不算錯誤
                            {
                                return;
                            } // catch
                            catch (Exception e)
                            {
                                equeue.Enqueue(EventID.Error, e);
                                return;
                            } // catch
                        } // while
                    });
                    thread.IsBackground = true;
                    thread.Start();
                } // try
                catch (Exception e)
                {
                    equeue.Enqueue(EventID.Error, e);
                } // catch
            }

            /// <summary>
            /// 關閉處理
            /// </summary>
            public void Close()
            {
                queue?.CompleteAdding(); // 佇列必須先結束才可能結束執行緒
                queue = null;
                thread?.Join();
                thread = null;
            }

            /// <summary>
            /// 新增訊息
            /// </summary>
            /// <param name="message">訊息物件</param>
            public void Add(object message)
            {
                if (queue == null)
                    return;

                if (queue.IsAddingCompleted)
                    return;

                queue.Add(message);
            }

            /// <summary>
            /// 執行緒物件
            /// </summary>
            private Thread thread = null;

            /// <summary>
            /// 封包佇列
            /// </summary>
            private BlockingCollection<object> queue = null;
        }

        /// <summary>
        /// 封包資料
        /// </summary>
        private class Packet
        {
            /// <summary>
            /// 封包長度
            /// </summary>
            public ushort Size = 0;

            /// <summary>
            /// 封包資料
            /// </summary>
            public byte[] Data = null;
        }

        public TCPClient(IEventmgr eventmgr, IProcmgr procmgr)
        {
            this.eventmgr = eventmgr;
            this.procmgr = procmgr;
        }

        public void Connect(string host, int port)
        {
            if (eventmgr == null)
                throw new ArgumentNullException("eventmgr");

            if (procmgr == null)
                throw new ArgumentNullException("procmgr");

            if (client != null)
                throw new AlreadyStartException("tcp client");

            this.host = host;
            this.port = port;

            equeue = new EQueue();
            client = new TcpClient();
            client.NoDelay = true;
            client.ReceiveBufferSize = Define.packetSize;
            client.SendBufferSize = Define.packetSize;

            var addr = Dns.GetHostAddresses(host);
            var callback = new AsyncCallback(
                (IAsyncResult result) =>
                {
                    client.EndConnect(result);
                    stream = client.GetStream();
                    recvHandler = new RecvHandler();
                    recvHandler.Start(stream, procmgr, equeue);
                    sendHandler = new SendHandler();
                    sendHandler.Start(stream, procmgr, equeue);
                    equeue.Enqueue(EventID.Connect, null);
                }
            );

            client.BeginConnect(addr, port, callback, this);
        }

        public void Disconnect()
        {
            stream?.Close();
            stream = null;
            recvHandler?.Close();
            recvHandler = null;
            sendHandler?.Close();
            sendHandler = null;
            client?.Close();
            client = null;
        }

        public void Update()
        {
            if (equeue == null)
                return;

            if (equeue.Dequeue(out var data) == false)
                return;

            if (data.eventID != EventID.Message)
            {
                eventmgr.Process(data.eventID, data.param);
                return;
            } // if

            try
            {
                procmgr.Process(data.param);
            } // try
            catch (Exception e)
            {
                equeue.Enqueue(EventID.Error, e);
            } // catch
        }

        public void Send(object message)
        {
            sendHandler?.Add(message);
        }

        public void AddEvent(EventID eventID, OnTrigger onEvent)
        {
            if (eventID != EventID.Message)
                eventmgr?.Add(eventID, onEvent);
        }

        public void DelEvent(EventID eventID)
        {
            if (eventID != EventID.Message)
                eventmgr?.Del(eventID);
        }

        public void AddProcess(MessageID messageID, OnTrigger onProcess)
        {
            procmgr?.Add(messageID, onProcess);
        }

        public void DelProcess(MessageID messageID)
        {
            procmgr?.Del(messageID);
        }

        public string GetHost()
        {
            return host;
        }

        public int GetPort()
        {
            return port;
        }

        public bool IsConnect()
        {
            return client != null && client.Connected;
        }

        /// <summary>
        /// 連線位址
        /// </summary>
        private string host = string.Empty;

        /// <summary>
        /// 連線埠號
        /// </summary>
        private int port = 0;

        /// <summary>
        /// 事件處理器
        /// </summary>
        private IEventmgr eventmgr = null;

        /// <summary>
        /// 訊息處理器
        /// </summary>
        private IProcmgr procmgr = null;

        /// <summary>
        /// 事件佇列
        /// </summary>
        private EQueue equeue = null;

        /// <summary>
        /// 連線物件
        /// </summary>
        private TcpClient client = null;

        /// <summary>
        /// 網路流物件
        /// </summary>
        private NetworkStream stream = null;

        /// <summary>
        /// 接收處理物件
        /// </summary>
        private RecvHandler recvHandler = null;

        /// <summary>
        /// 傳送處理物件
        /// </summary>
        private SendHandler sendHandler = null;
    }
}
