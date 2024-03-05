using System;
using System.Collections;
using NUnit.Framework;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    internal class TestProtoProc
    {
        [Test]
        [TestCaseSource("EncodeCases")]
        public void Encode(ProtoMsg message)
        {
            var proc = new ProtoProc();
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeBase64(ProtoMsg message)
        {
            var proc = new ProtoProc().SetBase64(true);
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeDesCBC(ProtoMsg message)
        {
            var proc = new ProtoProc().SetDesCBC(true, key, key);
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeAll(ProtoMsg message)
        {
            var proc = new ProtoProc().SetBase64(true).SetDesCBC(true, key, key);
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(ProtoProc.Marshal(1, new ProtoTest { Data = "test1" }));
                yield return new TestCaseData(ProtoProc.Marshal(2, new ProtoTest { Data = "test2" }));
            }
        }

        [Test]
        public void EncodeFailed()
        {
            var proc = new ProtoProc();

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
            var proc = new ProtoProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                proc.Decode(null);
            });
        }

        [Test]
        [TestCaseSource("ProcessCases")]
        public void Process(ProtoMsg message)
        {
            var proc = new ProtoProc();
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
                yield return new TestCaseData(ProtoProc.Marshal(1, new ProtoTest { Data = "test1" }));
                yield return new TestCaseData(ProtoProc.Marshal(2, new ProtoTest { Data = "test2" }));
            }
        }

        [Test]
        public void ProcessFailed()
        {
            var proc = new ProtoProc();

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
                proc.Process(new ProtoMsg { MessageID = 1 });
            });
        }

        [Test]
        [TestCaseSource("MarshalCases")]
        public void Marshal(MessageID messageID, ProtoTest message)
        {
            var marshal = ProtoProc.Marshal(messageID, message);

            ProtoProc.Unmarshal<ProtoTest>(marshal, out var resultID, out var result);
            Assert.AreEqual(messageID, resultID);
            Assert.IsTrue(TestUtil.EqualsByJson(message, result));
        }

        public static IEnumerable MarshalCases
        {
            get
            {
                yield return new TestCaseData(1, new ProtoTest { Data = "test1" });
                yield return new TestCaseData(2, new ProtoTest { Data = "test2" });
            }
        }

        [Test]
        public void MarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProtoProc.Marshal(1, null);
            });
        }

        [Test]
        public void UnmarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                JsonProc.Unmarshal<ProtoTest>(null, out var _, out var _);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                JsonProc.Unmarshal<ProtoTest>(new object(), out var _, out var _);
            });
        }

        private string key = "thisakey";
    }
}
