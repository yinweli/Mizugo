using System;
using System.Collections;
using System.Text;
using NUnit.Framework;

namespace Mizugo
{
    using MessageID = Int32;

    internal class TestProcJson
    {
        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeDecode(JsonMsg message)
        {
            var target = new ProcJson();
            var encode = target.Encode(message);
            var decode = target.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        public void EncodeFailed()
        {
            var target = new ProcJson();

            Assert.Throws<ArgumentNullException>(() =>
            {
                target.Encode(null);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                target.Encode(new object());
            });
        }

        [Test]
        public void DecodeFailed()
        {
            var target = new ProcJson();

            Assert.Throws<ArgumentNullException>(() =>
            {
                target.Decode(null);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                target.Encode(new object());
            });
        }

        [Test]
        [TestCaseSource("ProcessCases")]
        public void Process(JsonMsg message)
        {
            var target = new ProcJson();
            var valid = false;

            target.Add(
                message.MessageID,
                (object param) =>
                {
                    valid = TestUtil.EqualsByJson(message, param);
                }
            );
            target.Process(message);
            Assert.IsTrue(valid);
        }

        [Test]
        public void ProcessFailed()
        {
            var target = new ProcJson();

            Assert.Throws<ArgumentNullException>(() =>
            {
                target.Process(null);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                target.Process(new object());
            });
            Assert.Throws<UnprocessException>(() =>
            {
                target.Process(new JsonMsg { MessageID = 1 });
            });
        }

        [Test]
        [TestCaseSource("MarshalCases")]
        public void Marshal(MessageID messageID, JsonTest message)
        {
            var marshal = ProcJson.Marshal(messageID, message);

            ProcJson.Unmarshal<JsonTest>(marshal, out var resultID, out var result);
            Assert.AreEqual(messageID, resultID);
            Assert.IsTrue(TestUtil.EqualsByJson(message, result));
        }

        [Test]
        public void MarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProcJson.Marshal(1, null);
            });
        }

        [Test]
        public void UnmarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProcJson.Unmarshal<JsonTest>(null, out var _, out var _);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                ProcJson.Unmarshal<JsonTest>(new object(), out var _, out var _);
            });
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(ProcJson.Marshal(1, Encoding.UTF8.GetBytes("test")));
                yield return new TestCaseData(ProcJson.Marshal(2, new byte[] { 0, 1, 2, }));
            }
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(ProcJson.Marshal(1, Encoding.UTF8.GetBytes("test")));
                yield return new TestCaseData(ProcJson.Marshal(2, new byte[] { 0, 1, 2, }));
            }
        }

        public static IEnumerable MarshalCases
        {
            get
            {
                yield return new TestCaseData(1, new JsonTest { Data = "test1" });
                yield return new TestCaseData(2, new JsonTest { Data = "test2" });
            }
        }
    }
}
