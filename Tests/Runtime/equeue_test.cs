using NUnit.Framework;

namespace Mizugo
{
    internal class TestEQueue
    {
        [Test]
        [TestCase(1, 9999)]
        [TestCase(2, "9999")]
        [TestCase(3, null)]
        public void Enqueue(EventID eventID, object param)
        {
            var equeue = new EQueue();

            equeue.Enqueue(eventID, param);
            Assert.IsTrue(equeue.Dequeue(out var result));
            Assert.AreEqual(eventID, result.eventID);
            Assert.AreEqual(param, result.param);
            Assert.IsFalse(equeue.Dequeue(out var _));
        }

        [Test]
        public void DequeueMulti()
        {
            var equeue = new EQueue();

            Assert.IsFalse(equeue.Dequeue(out var _));
            Assert.IsFalse(equeue.Dequeue(out var _));
        }

        [Test]
        public void Misc()
        {
            var equeue = new EQueue();

            Assert.IsTrue(equeue.IsEmpty);

            equeue.Enqueue(EventID.Connect, null);
            equeue.Enqueue(EventID.Disconnect, null);

            Assert.IsFalse(equeue.IsEmpty);
        }
    }
}
