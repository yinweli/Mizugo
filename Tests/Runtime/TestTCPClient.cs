using NUnit.Framework;
using System;
using System.Diagnostics;

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
            var vaildConnect = false;
            var vaildDisconnect = false;

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    vaildConnect = true;
                }
            );
            client.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    vaildDisconnect = true;
                }
            );

            client.Connect(host, port);
            TestUtil.Sleep();
            client.Disconnect();
            TestUtil.Sleep();

            while (client.IsUpdate())
                client.Update();

            Assert.IsTrue(vaildConnect);
            Assert.IsTrue(vaildDisconnect);
        }

        [Test]
        [TestCase("google.com", 80)]
        public void ConnectFailed(string host, int port)
        {
            var client = new TCPClient(null, new JsonProc());
            Assert.Throws<ArgumentNullException>(() =>
            {
                client.Connect(host, port);
            });

            client = new TCPClient(new Eventmgr(), null);
            Assert.Throws<ArgumentNullException>(() =>
            {
                client.Connect(host, port);
            });

            client = new TCPClient(new Eventmgr(), new JsonProc());
            client.Connect(host, port);
            Assert.Throws<AlreadyStartException>(() =>
            {
                client.Connect(host, port);
            });
        }

        [Test]
        public void DisconnectFailed()
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            client.Disconnect();
            client.Disconnect();
        }

        [Test]
        [TestCase("google.com", 80)]
        public void UpdateFailed(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            client.Update();
            client.Connect(host, port);
            TestUtil.Sleep();
            client.Update();
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
        [TestCase("github.com", 80)]
        public void Misc(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());

            Assert.IsFalse(client.IsConnect());
            client.Connect(host, port);
            TestUtil.Sleep();
            Assert.AreEqual(host, client.GetHost());
            Assert.AreEqual(port, client.GetPort());
            Assert.IsTrue(client.IsConnect());
        }
    }

    internal class TestTCPClientCycle
    {
        /// <summary>
        /// 這項測試需要啟動測試伺服器才能執行
        /// </summary>
        [Test]
        [TestCase("127.0.0.1", 10003, "key-init", 1000)]
        public void Test(string host, int port, string key, int count)
        {
            var client = new TCPClient(new Eventmgr(), new PListProc { KeyStr = key, IVStr = key, });
            var stopwatch = new Stopwatch();
            var actual = 0;

            void SendMPListQ()
            {
                client.Send(PListProc.Marshal((int)MsgID.PlistQ, new MPListQ { Time = stopwatch.ElapsedMilliseconds }));
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
                (int)MsgID.PlistA,
                (object param) =>
                {
                    PListProc.Unmarshal<MPListA>(param, out var messageID, out var message);
                    TestUtil.Log("duration: " + (stopwatch.ElapsedMilliseconds - message.From.Time));
                    TestUtil.Log("count: " + message.Count);

                    if (actual < count)
                        SendMPListQ();
                }
            );

            stopwatch.Start();
            client.Connect(host, port);
            TestUtil.Sleep();

            while (client.IsUpdate() || actual < count)
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
        [TestCase("127.0.0.1", 10001)]
        public void Test(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new JsonProc());
            var stopwatch = new Stopwatch();
            var vaildConnect = false;
            var vaildDisconnect = false;
            var vaildRecv = false;
            var vaildSend = false;
            var validMessage = false;

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    vaildConnect = true;
                }
            );
            client.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    vaildDisconnect = true;
                }
            );
            client.AddEvent(
                EventID.Recv,
                (object _) =>
                {
                    vaildRecv = true;
                    client.Disconnect();
                }
            );
            client.AddEvent(
                EventID.Send,
                (object _) =>
                {
                    vaildSend = true;
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

            while (client.IsUpdate())
                client.Update();

            Assert.IsTrue(vaildConnect);
            Assert.IsTrue(vaildDisconnect);
            Assert.IsTrue(vaildRecv);
            Assert.IsTrue(vaildSend);
            Assert.IsTrue(validMessage);
        }
    }

    internal class TestTCPClientProto
    {
        /// <summary>
        /// 這項測試需要啟動測試伺服器才能執行
        /// </summary>
        [Test]
        [TestCase("127.0.0.1", 10002)]
        public void Test(string host, int port)
        {
            var client = new TCPClient(new Eventmgr(), new ProtoProc());
            var stopwatch = new Stopwatch();
            var vaildConnect = false;
            var vaildDisconnect = false;
            var vaildRecv = false;
            var vaildSend = false;
            var validMessage = false;

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    vaildConnect = true;
                }
            );
            client.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    vaildDisconnect = true;
                }
            );
            client.AddEvent(
                EventID.Recv,
                (object _) =>
                {
                    vaildRecv = true;
                    client.Disconnect();
                }
            );
            client.AddEvent(
                EventID.Send,
                (object _) =>
                {
                    vaildSend = true;
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

            while (client.IsUpdate())
                client.Update();

            Assert.IsTrue(vaildConnect);
            Assert.IsTrue(vaildDisconnect);
            Assert.IsTrue(vaildRecv);
            Assert.IsTrue(vaildSend);
            Assert.IsTrue(validMessage);
        }
    }

    internal class TestTCPClientPList
    {
        /// <summary>
        /// 這項測試需要啟動測試伺服器才能執行
        /// </summary>
        [Test]
        [TestCase("127.0.0.1", 10003, "key-init")]
        public void Test(string host, int port, string key)
        {
            var client = new TCPClient(new Eventmgr(), new PListProc { KeyStr = key, IVStr = key, });
            var stopwatch = new Stopwatch();
            var vaildConnect = false;
            var vaildDisconnect = false;
            var vaildRecv = false;
            var vaildSend = false;
            var validMessage = false;

            client.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    vaildConnect = true;
                }
            );
            client.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    vaildDisconnect = true;
                }
            );
            client.AddEvent(
                EventID.Recv,
                (object _) =>
                {
                    vaildRecv = true;
                    client.Disconnect();
                }
            );
            client.AddEvent(
                EventID.Send,
                (object _) =>
                {
                    vaildSend = true;
                }
            );
            client.AddEvent(EventID.Error, TestUtil.Log);
            client.AddProcess(
                (int)MsgID.PlistA,
                (object param) =>
                {
                    PListProc.Unmarshal<MPListA>(param, out var messageID, out var message);
                    TestUtil.Log("duration: " + (stopwatch.ElapsedMilliseconds - message.From.Time));
                    TestUtil.Log("count: " + message.Count);
                    validMessage = true;
                }
            );

            stopwatch.Start();
            client.Connect(host, port);
            TestUtil.Sleep();
            client.Send(PListProc.Marshal((int)MsgID.PlistQ, new MPListQ { Time = stopwatch.ElapsedMilliseconds }));
            TestUtil.Sleep();

            while (client.IsUpdate())
                client.Update();

            Assert.IsTrue(vaildConnect);
            Assert.IsTrue(vaildDisconnect);
            Assert.IsTrue(vaildRecv);
            Assert.IsTrue(vaildSend);
            Assert.IsTrue(validMessage);
        }
    }
}
