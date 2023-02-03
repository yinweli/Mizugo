using NUnit.Framework;

namespace Mizugo
{
    internal class TestEventmgr
    {
        [Test]
        [TestCase(1, 9999)]
        [TestCase(2, "9999")]
        [TestCase(3, null)]
        public void Add(EventID eventID, object param)
        {
            var eventmgr = new Eventmgr();
            var expected = param;
            var valid = false;

            eventmgr.Add(
                eventID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );
            eventmgr.Process(eventID, param);
            Assert.IsTrue(valid);
        }

        [Test]
        public void AddNull()
        {
            var eventmgr = new Eventmgr();
            var eventID = (EventID)1;

            eventmgr.Add(eventID, null);
            eventmgr.Process(eventID, null);
        }

        [Test]
        [TestCase(1, 9999)]
        [TestCase(2, "9999")]
        [TestCase(3, null)]
        public void Del(EventID eventID, object param)
        {
            var eventmgr = new Eventmgr();
            var expected = param;
            var valid = false;

            eventmgr.Add(
                eventID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );
            eventmgr.Del(eventID);
            eventmgr.Process(eventID, param);
            Assert.IsFalse(valid);
        }
    }
}
