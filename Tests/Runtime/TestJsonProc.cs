using System;
using System.Collections;
using System.Text;
using NUnit.Framework;

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
        public void Encode(JsonMsg message)
        {
            var proc = new JsonProc();
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeBase64(JsonMsg message)
        {
            var proc = new JsonProc().SetBase64(true);
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeDesCBC(JsonMsg message)
        {
            var proc = new JsonProc().SetDesCBC(true, key, key);
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeAll(JsonMsg message)
        {
            var proc = new JsonProc().SetBase64(true).SetDesCBC(true, key, key);
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
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
            var proc = new JsonProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                proc.Encode(null);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                proc.Encode(new object());
            });
        }

        [Test]
        public void DecodeFailed()
        {
            var proc = new JsonProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                proc.Decode(null);
            });
        }

        [Test]
        [TestCaseSource("ProcessCases")]
        public void Process(JsonMsg message)
        {
            var proc = new JsonProc();
            var valid = false;

            proc.Add(
                message.MessageID,
                (object param) =>
                {
                    valid = TestUtil.EqualsByJson(message, param);
                }
            );
            proc.Process(message);
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
            var proc = new JsonProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                proc.Process(null);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                proc.Process(new object());
            });
            Assert.Throws<UnprocessException>(() =>
            {
                proc.Process(new JsonMsg { MessageID = 1 });
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
                yield return new TestCaseData(1, new JsonTest { Data = "test1" });
                yield return new TestCaseData(2, new JsonTest { Data = "test2" });
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

        private string key = "thisakey";
    }
}
