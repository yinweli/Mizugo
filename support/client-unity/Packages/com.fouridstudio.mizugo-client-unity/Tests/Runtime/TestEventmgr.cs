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
            var target = new Eventmgr();
            var expected = param;
            var valid = false;

            target.Add(
                eventID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );
            target.Process(eventID, param);
            Assert.IsTrue(valid);
        }

        [Test]
        [TestCase(1, 9999)]
        [TestCase(2, "9999")]
        [TestCase(3, null)]
        public void Del(EventID eventID, object param)
        {
            var target = new Eventmgr();
            var expected = param;
            var valid = false;

            target.Add(
                eventID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );
            target.Del(eventID);
            target.Process(eventID, param);
            Assert.IsFalse(valid);
        }
    }
}
