using NUnit.Framework;
using System.Collections;

namespace Mizugo
{
    internal class TestEQueue
    {
        [Test, TestCaseSource("EnqueueCases")]
        public void Enqueue(EventID eventID, object param)
        {
            var equeue = new EQueue();

            equeue.Enqueue(eventID, param);
            Assert.IsTrue(equeue.Dequeue(out var result));
            Assert.AreEqual(eventID, result.eventID);
            Assert.AreEqual(param, result.param);
            Assert.IsFalse(equeue.Dequeue(out var _));
        }

        public static IEnumerable EnqueueCases
        {
            get
            {
                yield return new TestCaseData(1, 9999);
                yield return new TestCaseData(2, "9999");
                yield return new TestCaseData(3, new object());
                yield return new TestCaseData(4, null);
            }
        }

        [Test]
        public void DequeueMulti()
        {
            var equeue = new EQueue();

            Assert.IsFalse(equeue.Dequeue(out var _));
            Assert.IsFalse(equeue.Dequeue(out var _));
        }
    }
}
