using NUnit.Framework;
using System.Collections;

namespace Mizugo
{
    internal class EventmgrSuite
    {
        [Test, TestCaseSource("EventmgrTestCases")]
        public void EventmgrTest(EventID eventID, object param)
        {
            var eventmgr = new Eventmgr();
            var tester = new TriggerTester();

            tester.Reset(param);
            eventmgr.Add(eventID, tester.Trigger);
            eventmgr.Process(eventID, param);
            Assert.IsTrue(tester.Valid());

            tester.Reset(param);
            eventmgr.Del(eventID);
            eventmgr.Process(eventID, param);
            Assert.IsFalse(tester.Valid());

            tester.Reset(param);
            eventmgr.Process(eventID, null);
            Assert.IsFalse(tester.Valid());
        }

        public static IEnumerable EventmgrTestCases
        {
            get
            {
                yield return new TestCaseData(1, 9999);
                yield return new TestCaseData(2, "9999");
                yield return new TestCaseData(3, new object());
            }
        }

        private class TriggerTester
        {
            public void Reset(object param)
            {
                expected = param;
                valid = false;
            }

            public bool Valid()
            {
                return valid;
            }

            public void Trigger(object param)
            {
                valid = expected == param;
            }

            private object expected = null;
            private bool valid = false;
        }
    }
}
