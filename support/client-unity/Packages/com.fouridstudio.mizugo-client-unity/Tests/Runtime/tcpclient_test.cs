using NUnit.Framework;
using System;

namespace Mizugo
{
    using MessageID = Int32;

    internal class TestTCPClient
    {
        [Test]
        [TestCase("google.com", 80)]
        [TestCase("github.com", 80)]
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
            client.Update();
            Assert.IsTrue(vaildConnect);
            client.Disconnect();
            TestUtil.Sleep();
            client.Update(); // 這次處理異常
            client.Update(); // 這次處理斷線事件
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
            TestUtil.Sleep();
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

    internal class TestTCPClientJson { }

    internal class TestTCPClientProto { }

    internal class TestTCPClientPList { }
}
