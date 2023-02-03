using NUnit.Framework;
using System;
using System.Collections;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    internal class TestJsonProc
    {
        [Test]
        [TestCaseSource("EncodeCases")]
        public void Encode(JsonMsg input)
        {
            var jsonproc = new JsonProc();
            var encode = jsonproc.Encode(input);
            var decode = jsonproc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(input, decode));
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(JsonProc.Marshal(1, Encoding.UTF8.GetBytes("test")));
                yield return new TestCaseData(JsonProc.Marshal(2, new byte[] { 0, 1, 2, }));
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

        [Test]
        [TestCaseSource("ProcessCases")]
        public void Process(JsonMsg jsonMsg)
        {
            var jsonproc = new JsonProc();
            var expected = jsonMsg;
            var valid = false;

            jsonproc.Add(
                jsonMsg.MessageID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );
            jsonproc.Process(jsonMsg);
            Assert.IsTrue(valid);
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(JsonProc.Marshal(1, Encoding.UTF8.GetBytes("test")));
                yield return new TestCaseData(JsonProc.Marshal(2, new byte[] { 0, 1, 2, }));
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
            Assert.Throws<UnprocessException>(() =>
            {
                jsonproc.Process(new JsonMsg() { MessageID = 1 });
            });
        }

        [Test]
        [TestCaseSource("MarshalCases")]
        public void Marshal(MessageID messageID, JsonTest message)
        {
            var marshal = JsonProc.Marshal(messageID, message);
            JsonProc.Unmarshal<JsonTest>(marshal, out var resultID, out var result);

            Assert.AreEqual(messageID, resultID);
            Assert.IsTrue(TestUtil.EqualsByJson(message, result));
        }

        public static IEnumerable MarshalCases
        {
            get
            {
                yield return new TestCaseData(1, new JsonTest() { Data = "test1" });
                yield return new TestCaseData(2, new JsonTest() { Data = "test2" });
            }
        }

        [Test]
        public void MarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                JsonProc.Marshal(1, null);
            });
        }

        [Test]
        public void UnmarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                JsonProc.Unmarshal<JsonTest>(null, out var _, out var _);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                JsonProc.Unmarshal<JsonTest>(new object(), out var _, out var _);
            });
        }
    }
}
