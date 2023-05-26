using NUnit.Framework;
using System;
using System.Diagnostics;
using System.Net.Sockets;

namespace Mizugo
{
    using MessageID = Int32;

    internal class TestTCPClient
    {
        [Test]
        [TestCase("google.com", 80)]
        public void Connect(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());
            var validConnect = false;
            var validDisconnect = false;

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    validConnect = true;
                }
            );
            client.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    validDisconnect = true;
                }
            );

            client.Connect(host, port);
            TestUtil.Sleep();
            client.Disconnect();
            TestUtil.Sleep();

            while (client.IsUpdate)
                client.Update();

            Assert.IsTrue(validConnect);
            Assert.IsTrue(validDisconnect);
        }

        [Test]
        [TestCase("google.com", 80)]
        public void ConnectFailed(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());
            var validAlreadyStart = false;
            var validTimeout = false;
            client.AddEvent(
                EventID.Error,
                (object e) =>
                {
                    if (e is AlreadyStartException)
                        validAlreadyStart = true;

                    if (e is SocketException exception && exception.SocketErrorCode == SocketError.TimedOut)
                        validTimeout = true;
                }
            );

            client.Connect(host, port);
            TestUtil.Sleep();
            client.Connect(host, port);
            TestUtil.Sleep();
            client.Disconnect();
            TestUtil.Sleep();
            client.Connect(host, port + 1);
            TestUtil.Sleep();

            while (client.IsUpdate)
                client.Update();

            Assert.IsTrue(validAlreadyStart);
            Assert.IsTrue(validTimeout);
        }

        [Test]
        public void DisconnectFailed()
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            client.Disconnect();
        }

        [Test]
        [TestCase("google.com", 80)]
        public void UpdateFailed(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            client.Update();
        }

        [Test]
        [TestCase(EventID.Connect)]
        [TestCase(EventID.Disconnect)]
        public void Event(EventID eventID)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            client.AddEvent(eventID, (object _) => { });
            client.DelEvent(eventID);
        }

        [Test]
        [TestCase(1)]
        [TestCase(2)]
        public void Process(MessageID messageID)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            client.AddProcess(messageID, (object _) => { });
            client.DelProcess(messageID);
        }

        [Test]
        [TestCase("google.com", 80)]
        public void Misc(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            Assert.IsFalse(client.IsConnect);
            client.Connect(host, port);
            TestUtil.Sleep();
            Assert.AreEqual(host, client.Host);
            Assert.AreEqual(port, client.Port);
            Assert.IsTrue(client.IsConnect);
            client.Disconnect();
        }
    }

    internal class TestTCPClientCycle
    {
        /// <summary>
        /// 這項測試需要啟動測試伺服器才能執行
        /// </summary>
        [Test]
        [TestCase("127.0.0.1", 9001, "key-@@@@", 1000)]
        public void Test(string host, int port, string key, int count)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc().SetBase64(true).SetDesCBC(true, key, key));
            var stopwatch = new Stopwatch();
            var actual = 0;

            void SendMPListQ()
            {
                client.Send(JsonProc.Marshal((int)MsgID.JsonQ, new MJsonQ { Time = stopwatch.ElapsedMilliseconds }));
                actual++;
            }

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    SendMPListQ();
                }
            );
            client.AddEvent(EventID.Error, TestUtil.Log);
            client.AddProcess(
                (int)MsgID.JsonA,
                (object param) =>
                {
                    JsonProc.Unmarshal<MJsonA>(param, out var messageID, out var message);
                    TestUtil.Log("duration: " + (stopwatch.ElapsedMilliseconds - message.From.Time));
                    TestUtil.Log("count: " + message.Count);

                    if (actual < count)
                        SendMPListQ();
                    else
                        client.Disconnect();
                }
            );

            stopwatch.Start();
            client.Connect(host, port);
            TestUtil.Sleep();

            while (client.IsUpdate)
                client.Update();

            Assert.AreEqual(count, actual);
        }
    }

    internal class TestTCPClientJson
    {
        /// <summary>
        /// 這項測試需要啟動測試伺服器才能執行
        /// </summary>
        [Test]
        [TestCase("127.0.0.1", 9001, "key-@@@@")]
        public void Test(string host, int port, string key)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc().SetBase64(true).SetDesCBC(true, key, key));
            var stopwatch = new Stopwatch();
            var validConnect = false;
            var validDisconnect = false;
            var validRecv = false;
            var validSend = false;
            var validMessage = false;

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    validConnect = true;
                }
            );
            client.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    validDisconnect = true;
                }
            );
            client.AddEvent(
                EventID.Recv,
                (object _) =>
                {
                    validRecv = true;
                    client.Disconnect();
                }
            );
            client.AddEvent(
                EventID.Send,
                (object _) =>
                {
                    validSend = true;
                }
            );
            client.AddEvent(EventID.Error, TestUtil.Log);
            client.AddProcess(
                (int)MsgID.JsonA,
                (object param) =>
                {
                    JsonProc.Unmarshal<MJsonA>(param, out var messageID, out var message);
                    TestUtil.Log("duration: " + (stopwatch.ElapsedMilliseconds - message.From.Time));
                    TestUtil.Log("count: " + message.Count);
                    validMessage = true;
                }
            );

            stopwatch.Start();
            client.Connect(host, port);
            TestUtil.Sleep();
            client.Send(JsonProc.Marshal((int)MsgID.JsonQ, new MJsonQ { Time = stopwatch.ElapsedMilliseconds }));
            TestUtil.Sleep();

            while (client.IsUpdate)
                client.Update();

            Assert.IsTrue(validConnect);
            Assert.IsTrue(validDisconnect);
            Assert.IsTrue(validRecv);
            Assert.IsTrue(validSend);
            Assert.IsTrue(validMessage);
        }
    }

    internal class TestTCPClientProto
    {
        /// <summary>
        /// 這項測試需要啟動測試伺服器才能執行
        /// </summary>
        [Test]
        [TestCase("127.0.0.1", 9002, "key-####")]
        public void Test(string host, int port, string key)
        {
            var client = new TCPClient(new Eventmgr(), new ProtoProc().SetBase64(true).SetDesCBC(true, key, key));
            var stopwatch = new Stopwatch();
            var validConnect = false;
            var validDisconnect = false;
            var validRecv = false;
            var validSend = false;
            var validMessage = false;

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    validConnect = true;
                }
            );
            client.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    validDisconnect = true;
                }
            );
            client.AddEvent(
                EventID.Recv,
                (object _) =>
                {
                    validRecv = true;
                    client.Disconnect();
                }
            );
            client.AddEvent(
                EventID.Send,
                (object _) =>
                {
                    validSend = true;
                }
            );
            client.AddEvent(EventID.Error, TestUtil.Log);
            client.AddProcess(
                (int)MsgID.ProtoA,
                (object param) =>
                {
                    ProtoProc.Unmarshal<MProtoA>(param, out var messageID, out var message);
                    TestUtil.Log("duration: " + (stopwatch.ElapsedMilliseconds - message.From.Time));
                    TestUtil.Log("count: " + message.Count);
                    validMessage = true;
                }
            );

            stopwatch.Start();
            client.Connect(host, port);
            TestUtil.Sleep();
            client.Send(ProtoProc.Marshal((int)MsgID.ProtoQ, new MProtoQ { Time = stopwatch.ElapsedMilliseconds }));
            TestUtil.Sleep();

            while (client.IsUpdate)
                client.Update();

            Assert.IsTrue(validConnect);
            Assert.IsTrue(validDisconnect);
            Assert.IsTrue(validRecv);
            Assert.IsTrue(validSend);
            Assert.IsTrue(validMessage);
        }
    }
}
