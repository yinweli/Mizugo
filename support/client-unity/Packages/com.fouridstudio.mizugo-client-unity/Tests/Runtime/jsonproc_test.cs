using NUnit.Framework;
using System;
using System.Collections;
using System.Linq;
using System.Text;

namespace Mizugo
{
    internal class TestJsonProc
    {
        [Test, TestCaseSource("EncodeCases")]
        public void Encode(JsonMsg input)
        {
            var jsonproc = new JsonProc();
            var encode = jsonproc.Encode(input);
            var decode = jsonproc.Decode(encode);

            TestUtil.AreEqualByJson(input, decode);
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(
                    new JsonMsg() { MessageID = 1, Message = Encoding.UTF8.GetBytes("test1") }
                );
                yield return new TestCaseData(
                    new JsonMsg() { MessageID = 2, Message = new byte[] { 0, 1, 2, } }
                );
            }
        }

        [Test]
        public void EncodeFailed()
        {
            var jsonproc = new JsonProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                jsonproc.Encode(null);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                jsonproc.Encode(new object());
            });
        }

        [Test]
        public void DecodeFailed()
        {
            var jsonproc = new JsonProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                jsonproc.Decode(null);
            });
        }

        [Test, TestCaseSource("ProcessCases")]
        public void Process(JsonMsg jsonMsg)
        {
            var jsonproc = new JsonProc();
            var expected = jsonMsg.Message;
            var valid = false;

            jsonproc.Add(
                jsonMsg.MessageID,
                (object param) =>
                {
                    valid = expected.SequenceEqual(param as byte[]);
                }
            );
            jsonproc.Process(jsonMsg);
            Assert.IsTrue(valid);
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(
                    new JsonMsg() { MessageID = 1, Message = Encoding.UTF8.GetBytes("test1") }
                );
                yield return new TestCaseData(
                    new JsonMsg() { MessageID = 2, Message = new byte[] { 0, 1, 2, } }
                );
            }
        }

        [Test]
        public void ProcessFailed()
        {
            var jsonproc = new JsonProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                jsonproc.Process(null);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                jsonproc.Process(new object());
            });
            Assert.Throws<UnProcessException>(() =>
            {
                jsonproc.Process(new JsonMsg() { MessageID = 1 });
            });
        }
    }
}
