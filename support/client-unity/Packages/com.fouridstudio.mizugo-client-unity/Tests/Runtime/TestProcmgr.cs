using System;
using System.Collections;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestProcmgr
    {
        [Test]
        [TestCase(1, 9999)]
        [TestCase(2, "9999")]
        [TestCase(3, null)]
        public void Add(int messageID, object param)
        {
            var target = new ProcEmpty();
            var expected = param;
            var valid = false;

            target.Add(
                messageID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );

            var process = target.Get(messageID);

            Assert.IsNotNull(process);
            process(param);
            Assert.IsTrue(valid);
        }

        [Test]
        [TestCase(1)]
        [TestCase(2)]
        public void Del(int messageID)
        {
            var target = new ProcEmpty();

            target.Add(messageID, (object param) => { });
            target.Del(messageID);
            Assert.IsNull(target.Get(messageID));
        }

        private class ProcEmpty : Procmgr
        {
            public override void Process(object message)
            {
                throw new System.NotImplementedException();
            }
        }
    }
}
