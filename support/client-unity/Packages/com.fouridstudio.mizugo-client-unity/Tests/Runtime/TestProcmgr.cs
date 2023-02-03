using NUnit.Framework;
using System;
using System.Collections;

namespace Mizugo
{
    using MessageID = Int32;

    internal class TestProcmgr
    {
        [Test]
        [TestCase(1, 9999)]
        [TestCase(2, "9999")]
        [TestCase(3, null)]
        public void Add(MessageID messageID, object param)
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

        [Test]
        [TestCase(1)]
        [TestCase(2)]
        public void Del(MessageID messageID)
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
