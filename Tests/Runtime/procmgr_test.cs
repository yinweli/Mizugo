using NUnit.Framework;
using System.Collections;

namespace Mizugo
{
    internal class TestProcmgr
    {
        [Test, TestCaseSource("AddCases")]
        public void Add(int messageID, object param)
        {
            var procmgr = new EmptyProc();
            var expected = param;
            var valid = false;

            procmgr.Add(
                messageID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );

            var process = procmgr.Get(messageID);

            Assert.IsNotNull(process);
            process(param);
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

        [Test, TestCaseSource("DelCases")]
        public void Del(int messageID)
        {
            var procmgr = new EmptyProc();

            procmgr.Add(messageID, (object param) => { });
            procmgr.Del(messageID);
            Assert.IsNull(procmgr.Get(messageID));
        }

        public static IEnumerable DelCases
        {
            get
            {
                yield return new TestCaseData(1);
                yield return new TestCaseData(2);
            }
        }

        private class EmptyProc : Procmgr
        {
            public override byte[] Encode(object input)
            {
                throw new System.NotImplementedException();
            }

            public override object Decode(byte[] input)
            {
                throw new System.NotImplementedException();
            }

            public override void Process(object message)
            {
                throw new System.NotImplementedException();
            }
        }
    }
}
