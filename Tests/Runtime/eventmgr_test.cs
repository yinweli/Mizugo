using NUnit.Framework;
using System.Collections;

namespace Mizugo
{
    internal class TestEventmgr
    {
        [Test, TestCaseSource("AddCases")]
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

        public static IEnumerable AddCases
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
        public void AddNull()
        {
            var eventmgr = new Eventmgr();
            var eventID = (EventID)1;

            eventmgr.Add(eventID, null);
            eventmgr.Process(eventID, null);
        }

        [Test, TestCaseSource("DelCases")]
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

        public static IEnumerable DelCases
        {
            get
            {
                yield return new TestCaseData(1, 9999);
                yield return new TestCaseData(2, "9999");
                yield return new TestCaseData(3, new object());
                yield return new TestCaseData(4, null);
            }
        }
    }
}
