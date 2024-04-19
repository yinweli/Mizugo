using System;
using System.Diagnostics;
using System.Net.Sockets;
using System.Security.Cryptography;
using NUnit.Framework;

namespace Mizugo
{
    using MessageID = Int32;

    internal class TestTCPClient
    {
        [Test]
        [TestCase("google.com", 80)]
        public void Connect(string host, int port)
        {
            var target = new TCPClient();
            var eventmgr = new Eventmgr();
            var process = new JsonProc();

            target.SetEvent(eventmgr);
            target.SetProc(process);
            target.SetCodec(process);

            var validConnect = false;
            var validDisconnect = false;

            target.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    validConnect = true;
                }
            );
            target.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    validDisconnect = true;
                }
            );
            target.AddEvent(
                EventID.Error,
                (object _) =>
                {
                    validDisconnect = true;
                }
            );

            target.Connect(host, port);
            TestUtil.Sleep();
            target.Disconnect();
            TestUtil.Sleep();

            while (target.IsUpdate)
                target.Update();

            Assert.IsTrue(validConnect);
            Assert.IsTrue(validDisconnect);
        }

        [Test]
        [TestCase("google.com", 80)]
        public void ConnectFailed(string host, int port)
        {
            var target = new TCPClient();
            var eventmgr = new Eventmgr();
            var process = new JsonProc();

            target.SetEvent(eventmgr);
            target.SetProc(process);
            target.SetCodec(process);

            var validAlreadyStart = false;
            var validTimeout = false;

            target.AddEvent(
                EventID.Error,
                (object e) =>
                {
                    if (e is AlreadyStartException)
                        validAlreadyStart = true;

                    if (e is SocketException exception && exception.SocketErrorCode == SocketError.TimedOut)
                        validTimeout = true;
                }
            );

            target.Connect(host, port);
            TestUtil.Sleep();
            target.Connect(host, port);
            TestUtil.Sleep();
            target.Disconnect();
            TestUtil.Sleep();
            target.Connect(host, port + 1);
            TestUtil.Sleep();

            while (target.IsUpdate)
                target.Update();

            Assert.IsTrue(validAlreadyStart);
            Assert.IsTrue(validTimeout);
        }

        [Test]
        public void DisconnectFailed()
        {
            new TCPClient().Disconnect();
        }

        [Test]
        public void UpdateFailed()
        {
            new TCPClient().Update();
        }

        [Test]
        [TestCase(EventID.Connect)]
        [TestCase(EventID.Disconnect)]
        public void Event(EventID eventID)
        {
            var target = new TCPClient();
            var eventmgr = new Eventmgr();

            target.SetEvent(eventmgr);
            target.AddEvent(eventID, (object _) => { });
            target.DelEvent(eventID);
        }

        [Test]
        [TestCase(1)]
        [TestCase(2)]
        public void Process(MessageID messageID)
        {
            var target = new TCPClient();
            var process = new JsonProc();

            target.SetProc(process);
            target.AddProcess(messageID, (object _) => { });
            target.DelProcess(messageID);
        }

        [Test]
        [TestCase("google.com", 80)]
        public void Misc(string host, int port)
        {
            var target = new TCPClient();
            var eventmgr = new Eventmgr();
            var process = new JsonProc();

            target.SetEvent(eventmgr);
            target.SetProc(process);
            target.SetCodec(process);

            Assert.IsFalse(target.IsConnect);
            target.Connect(host, port);
            TestUtil.Sleep();
            Assert.IsTrue(target.IsConnect);
            target.Disconnect();
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
            var target = new TCPClient();
            var eventmgr = new Eventmgr();
            var process = new JsonProc();

            target.SetEvent(eventmgr);
            target.SetProc(process);
            target.SetCodec(process, new DesCBC(PaddingMode.PKCS7, key, key), new Base64());

            var stopwatch = new Stopwatch();
            var actual = 0;

            target.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    target.Send(JsonProc.Marshal((int)MsgID.JsonQ, new MJsonQ { Time = stopwatch.ElapsedMilliseconds }));
                    actual++;
                }
            );
            target.AddEvent(EventID.Error, TestUtil.Log);
            target.AddProcess(
                (int)MsgID.JsonA,
                (object param) =>
                {
                    JsonProc.Unmarshal<MJsonA>(param, out var messageID, out var message);
                    TestUtil.Log("duration: " + (stopwatch.ElapsedMilliseconds - message.From.Time));
                    TestUtil.Log("count: " + message.Count);

                    if (actual < count)
                    {
                        target.Send(JsonProc.Marshal((int)MsgID.JsonQ, new MJsonQ { Time = stopwatch.ElapsedMilliseconds }));
                        actual++;
                    }
                    else
                        target.Disconnect();
                }
            );

            stopwatch.Start();
            target.Connect(host, port);
            TestUtil.Sleep();

            while (target.IsUpdate)
                target.Update();

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
            var target = new TCPClient();
            var eventmgr = new Eventmgr();
            var process = new JsonProc();

            target.SetEvent(eventmgr);
            target.SetProc(process);
            target.SetCodec(process, new DesCBC(PaddingMode.PKCS7, key, key), new Base64());

            var stopwatch = new Stopwatch();
            var validConnect = false;
            var validDisconnect = false;
            var validRecv = false;
            var validSend = false;
            var validMessage = false;

            target.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    validConnect = true;
                }
            );
            target.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    validDisconnect = true;
                }
            );
            target.AddEvent(
                EventID.Recv,
                (object _) =>
                {
                    validRecv = true;
                    target.Disconnect();
                }
            );
            target.AddEvent(
                EventID.Send,
                (object _) =>
                {
                    validSend = true;
                }
            );
            target.AddEvent(
                EventID.Error,
                (object param) =>
                {
                    validDisconnect = true;
                    TestUtil.Log(param);
                }
            );
            target.AddProcess(
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
            target.Connect(host, port);
            TestUtil.Sleep();
            target.Send(JsonProc.Marshal((int)MsgID.JsonQ, new MJsonQ { Time = stopwatch.ElapsedMilliseconds }));
            TestUtil.Sleep();

            while (target.IsUpdate)
                target.Update();

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
            var target = new TCPClient();
            var eventmgr = new Eventmgr();
            var process = new ProtoProc();

            target.SetEvent(eventmgr);
            target.SetProc(process);
            target.SetCodec(process, new DesCBC(PaddingMode.PKCS7, key, key), new Base64());

            var stopwatch = new Stopwatch();
            var validConnect = false;
            var validDisconnect = false;
            var validRecv = false;
            var validSend = false;
            var validMessage = false;

            target.AddEvent(
                EventID.Connect,
                (object _) =>
                {
                    validConnect = true;
                }
            );
            target.AddEvent(
                EventID.Disconnect,
                (object _) =>
                {
                    validDisconnect = true;
                }
            );
            target.AddEvent(
                EventID.Recv,
                (object _) =>
                {
                    validRecv = true;
                    target.Disconnect();
                }
            );
            target.AddEvent(
                EventID.Send,
                (object _) =>
                {
                    validSend = true;
                }
            );
            target.AddEvent(
                EventID.Error,
                (object param) =>
                {
                    validDisconnect = true;
                    TestUtil.Log(param);
                }
            );
            target.AddProcess(
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
            target.Connect(host, port);
            TestUtil.Sleep();
            target.Send(ProtoProc.Marshal((int)MsgID.ProtoQ, new MProtoQ { Time = stopwatch.ElapsedMilliseconds }));
            TestUtil.Sleep();

            while (target.IsUpdate)
                target.Update();

            Assert.IsTrue(validConnect);
            Assert.IsTrue(validDisconnect);
            Assert.IsTrue(validRecv);
            Assert.IsTrue(validSend);
            Assert.IsTrue(validMessage);
        }
    }
}
