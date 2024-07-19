using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Threading;

namespace Mizugo
{
    /// <summary>
    /// TCP客戶端組件
    /// </summary>
    public class TCPClient : IClient
    {
        public void Connect(string host, int port)
        {
            try
            {
                if (eventmgr == null)
                    throw new ArgumentNullException("eventmgr");

                if (procmgr == null)
                    throw new ArgumentNullException("procmgr");

                if (codec == null)
                    throw new ArgumentNullException("codec");

                if (client != null)
                    throw new AlreadyStartException("client");

                var bufferSize = (Define.headerSize + packetSize) * 10;

                client = new TcpClient
                {
                    NoDelay = true,
                    ReceiveBufferSize = bufferSize,
                    SendBufferSize = bufferSize
                };
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
                                sendHandler = new SendHandler();
                                sendHandler.Start(packetSize, queue, stream, codec);
                                recvHandler = new RecvHandler();
                                recvHandler.Start(packetSize, queue, stream, codec.Reverse());
                                queue.Enqueue(EventID.Connect, null);
                            }
                            else
                                queue.Enqueue(EventID.Disconnect, null);
                        } // try
                        catch (Exception e)
                        {
                            queue?.Enqueue(EventID.Error, e);
                        } // catch

                        connecting = false;
                    }
                );

                client.BeginConnect(addr, port, callback, client);
            } // try
            catch (Exception e)
            {
                queue?.Enqueue(EventID.Error, e);
            } // catch
        }

        public void Disconnect()
        {
            stream?.Close();
            stream = null;
            sendHandler?.Close();
            sendHandler = null;
            recvHandler?.Close();
            recvHandler = null;
            client?.Close();
            client = null;
        }

        public void Update()
        {
            if (queue == null)
                return;

            if (queue.Dequeue(out var data) == false)
                return;

            try
            {
                if (data.eventID == EventID.Message)
                    procmgr?.Process(data.param);
                else
                    eventmgr?.Process(data.eventID, data.param);
            } // try
            catch (Exception e)
            {
                queue?.Enqueue(EventID.Error, e);
            } // catch
        }

        public void Send(object message)
        {
            sendHandler?.Add(message);
        }

        public void SetEvent(IEventmgr eventmgr)
        {
            this.eventmgr = eventmgr;
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

        public void SetProc(IProcmgr procmgr)
        {
            this.procmgr = procmgr;
        }

        public void AddProcess(int messageID, OnTrigger onProcess)
        {
            procmgr?.Add(messageID, onProcess);
        }

        public void DelProcess(int messageID)
        {
            procmgr?.Del(messageID);
        }

        public void SetCodec(params ICodec[] codec)
        {
            this.codec = codec;
        }

        public void SetPacketSize(int size)
        {
            this.packetSize = size;
        }

        public bool IsConnect
        {
            get { return (client != null && client.Connected) || connecting; }
        }

        public bool IsUpdate
        {
            get { return IsConnect || queue.IsEmpty == false; }
        }

        private IEventmgr eventmgr = null;
        private IProcmgr procmgr = null;
        private ICodec[] codec = null;
        private int packetSize = Define.packetSize;
        private Queue queue = new Queue(); // 必須在建立物件時就建立好, 不然會造成各種錯誤
        private TcpClient client = null;
        private NetworkStream stream = null;
        private SendHandler sendHandler = null;
        private RecvHandler recvHandler = null;
        private bool connecting = false;

        /// <summary>
        /// 傳送處理器
        /// </summary>
        private class SendHandler
        {
            public void Start(int packetSize, Queue queue, NetworkStream stream, IEnumerable<ICodec> codec)
            {
                if (thread != null)
                    throw new AlreadyStartException("send handler");

                this.queue = new BlockingCollection<object>();
                thread = new Thread(() =>
                {
                    while (true)
                    {
                        try
                        {
                            var message = this.queue.Take();

                            if (message == null)
                                continue;

                            foreach (var itor in codec)
                                message = itor.Encode(message);

                            if (message is not byte[] packet)
                                throw new PacketNullException("send");

                            if (packet.Length <= 0)
                                throw new PacketZeroException("send");

                            if (packet.Length > packetSize)
                                throw new PacketLimitException("send");

                            var buffer = BitConverter.GetBytes(packet.Length);

                            stream.Write(buffer, 0, buffer.Length);
                            stream.Write(packet, 0, packet.Length);
                            queue.Enqueue(EventID.Send, null);
                        } // try
                        catch (InvalidOperationException) // 這是因為關閉處理的CompleteAdding引發的, 所以不算錯誤
                        {
                            return;
                        } // catch
                        catch (Exception e)
                        {
                            queue.Enqueue(EventID.Error, e);
                            return;
                        } // catch
                    } // while
                });
                thread.IsBackground = true;
                thread.Start();
            }

            public void Close()
            {
                queue?.CompleteAdding(); // 佇列必須先結束才可能結束執行緒
                queue = null;
                thread?.Join();
                thread = null;
            }

            public void Add(object message)
            {
                if (queue == null)
                    return;

                if (queue.IsAddingCompleted)
                    return;

                queue.Add(message);
            }

            private BlockingCollection<object> queue = null;
            private Thread thread = null;
        }

        /// <summary>
        /// 接收處理器
        /// </summary>
        private class RecvHandler
        {
            public void Start(int packetSize, Queue queue, NetworkStream stream, IEnumerable<ICodec> codec)
            {
                if (thread != null)
                    throw new AlreadyStartException("recv handler");

                thread = new Thread(() =>
                {
                    var header = new byte[Define.headerSize];
                    var packet = (byte[])null;
                    var packetLen = 0;
                    var read = 0;

                    while (true)
                    {
                        try
                        {
                            var headerLen = stream.Read(header, 0, Define.headerSize);

                            if (headerLen == 0)
                                throw new DisconnectException();

                            if (headerLen != Define.headerSize)
                                throw new RecvHeaderException();

                            packetLen = BitConverter.ToInt32(header, 0);

                            if (packetLen <= 0)
                                throw new PacketZeroException("recv");

                            if (packetLen > packetSize)
                                throw new PacketLimitException("recv");

                            packet = new byte[packetLen];

                            while (packetLen > read)
                            {
                                var cursor = stream.Read(packet, read, packetLen - read);

                                if (cursor == 0)
                                    break;

                                read += cursor;
                            } // while

                            if (packetLen != read)
                                throw new RecvPacketException();

                            var messsage = packet as object;

                            foreach (var itor in codec)
                                messsage = itor.Decode(messsage);

                            queue.Enqueue(EventID.Message, messsage);
                            queue.Enqueue(EventID.Recv, null);
                            read = 0;
                        } // try
                        catch (DisconnectException)
                        {
                            queue.Enqueue(EventID.Disconnect, null);
                            return;
                        } // catch
                        catch (Exception e)
                        {
                            queue.Enqueue(EventID.Error, e);
                            return;
                        } // catch
                    } // while
                });
                thread.IsBackground = true;
                thread.Start();
            }

            public void Close()
            {
                thread?.Join();
                thread = null;
            }

            private Thread thread = null;
        }

        /// <summary>
        /// 事件佇列
        /// </summary>
        private class Queue
        {
            public class Data
            {
                public EventID eventID;
                public object param;
            }

            public bool IsEmpty
            {
                get { return queue.IsEmpty; }
            }

            public void Enqueue(EventID eventID, object param)
            {
                queue.Enqueue(new Data { eventID = eventID, param = param });
            }

            public bool Dequeue(out Data data)
            {
                return queue.TryDequeue(out data);
            }

            private ConcurrentQueue<Data> queue = new ConcurrentQueue<Data>();
        }
    }
}
