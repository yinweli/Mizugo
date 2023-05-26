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
        public TCPClient(IEventmgr eventmgr, IProcmgr procmgr)
        {
            this.eventmgr = eventmgr;
            this.procmgr = procmgr;
            this.equeue = new EQueue();
        }

        public void Connect(string host, int port)
        {
            try
            {
                if (eventmgr == null)
                    throw new ArgumentNullException("eventmgr");

                if (procmgr == null)
                    throw new ArgumentNullException("procmgr");

                if (client != null)
                    throw new AlreadyStartException("tcp client");

                this.host = host;
                this.port = port;

                client = new TcpClient();
                client.NoDelay = true;
                client.ReceiveBufferSize = Define.bufferSize;
                client.SendBufferSize = Define.bufferSize;
                connecting = true;

                var addr = Dns.GetHostAddresses(host);
                var callback = new AsyncCallback(
                    (IAsyncResult result) =>
                    {
                        try
                        {
                            client.EndConnect(result);

                            if (client.Connected)
                            {
                                stream = client.GetStream();
                                recvHandler = new RecvHandler();
                                recvHandler.Start(equeue, stream, procmgr);
                                sendHandler = new SendHandler();
                                sendHandler.Start(equeue, stream, procmgr);
                                equeue.Enqueue(EventID.Connect, null);
                            } // if
                        } // try
                        catch (Exception e)
                        {
                            equeue?.Enqueue(EventID.Error, e);
                        } // catch

                        connecting = false;
                    }
                );

                client.BeginConnect(addr, port, callback, client);
            } // try
            catch (Exception e)
            {
                equeue?.Enqueue(EventID.Error, e);
            } // catch
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
            if (equeue.Dequeue(out var data) == false)
                return;

            try
            {
                if (data.eventID != EventID.Message)
                    eventmgr?.Process(data.eventID, data.param);
                else
                    procmgr?.Process(data.param);
            } // try
            catch (Exception e)
            {
                equeue?.Enqueue(EventID.Error, e);
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

        public string Host
        {
            get { return host; }
        }

        public int Port
        {
            get { return port; }
        }

        public bool IsConnect
        {
            get { return (client != null && client.Connected) || connecting; }
        }

        public bool IsUpdate
        {
            get { return IsConnect || equeue.IsEmpty == false; }
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
        /// 客戶端物件
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

        /// <summary>
        /// 連線中旗標
        /// </summary>
        private bool connecting = false;

        /// <summary>
        /// 接收處理器
        /// </summary>
        private class RecvHandler
        {
            /// <summary>
            /// 啟動處理
            /// </summary>
            /// <param name="equeue">事件佇列</param>
            /// <param name="stream">網路流物件</param>
            /// <param name="procmgr">訊息處理器</param>
            public void Start(EQueue equeue, NetworkStream stream, IProcmgr procmgr)
            {
                if (thread != null)
                    throw new AlreadyStartException("recv handler");

                thread = new Thread(() =>
                {
                    // 如果想要改用ArrayPool, 需要改到IProcmgr.Encode, IProcmgr.Decode, Des加密/解密的函式以及相關的函式等, 影響很大, 所以現在先不動

                    var header = new byte[Define.headerSize];
                    var packet = (byte[])null;
                    var size = (ushort)0;
                    var read = (ushort)0;

                    while (true)
                    {
                        try
                        {
                            if (stream.Read(header, 0, Define.headerSize) != Define.headerSize)
                                throw new RecvHeaderException();

                            size = BitConverter.ToUInt16(header, 0);

                            if (size <= 0)
                                throw new PacketZeroException("recv");

                            if (size > Define.packetSize)
                                throw new PacketLimitException("recv");

                            packet = new byte[size];

                            while (size > read)
                            {
                                var readnow = stream.Read(packet, read, size - read);

                                if (readnow == 0)
                                    break;

                                read += (ushort)readnow;
                            } // while

                            if (size != read)
                                throw new RecvPacketException();

                            var message = procmgr.Decode(packet);

                            equeue.Enqueue(EventID.Message, message);
                            equeue.Enqueue(EventID.Recv, null);
                            read = 0;
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
            /// <param name="equeue">事件佇列</param>
            /// <param name="stream">網路流物件</param>
            /// <param name="procmgr">訊息處理器</param>
            public void Start(EQueue equeue, NetworkStream stream, IProcmgr procmgr)
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

                            if (packet.Length <= 0)
                                throw new PacketZeroException("send");

                            if (packet.Length > Define.packetSize)
                                throw new PacketLimitException("send");

                            var buffer = BitConverter.GetBytes((ushort)packet.Length);

                            stream.Write(buffer, 0, buffer.Length);
                            stream.Write(packet, 0, packet.Length);
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
            /// 封包佇列
            /// </summary>
            private BlockingCollection<object> queue = null;

            /// <summary>
            /// 執行緒物件
            /// </summary>
            private Thread thread = null;
        }
    }
}
