using NUnit.Framework;
using System.Collections;

namespace Mizugo
{
    internal class ProcmgrSuite
    {
        [Test, TestCaseSource("ProcmgrTestCases")]
        public void ProcmgrTest(int messageID)
        {
            var procmgr = new ProcmgrTester();

            procmgr.Add(messageID, procmgr.Trigger);
            Assert.AreEqual((OnTrigger)procmgr.Trigger, procmgr.Get(messageID));

            procmgr.Del(messageID);
            Assert.IsNull(procmgr.Get(messageID));
        }

        public static IEnumerable ProcmgrTestCases
        {
            get
            {
                yield return new TestCaseData(1);
                yield return new TestCaseData(2);
            }
        }

        /// <summary>
        /// 測試用的訊息處理器以字串為核心來進行處理
        /// </summary>
        private class ProcmgrTester : Procmgr
        {
            public override byte[] Encode(object input)
            {
                throw new System.NotImplementedException();
            }

            public override object Decode(byte[] input)
            {
                throw new System.NotImplementedException();
            }

            public override bool Process(object message)
            {
                throw new System.NotImplementedException();
            }

            public void Trigger(object _) { }
        }
    }
}
