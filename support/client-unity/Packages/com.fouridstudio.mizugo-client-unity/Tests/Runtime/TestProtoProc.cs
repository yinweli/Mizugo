using System;
using System.Collections;
using NUnit.Framework;

namespace Mizugo
{
    using static UnityEngine.GraphicsBuffer;
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    internal class TestProtoProc
    {
        [Test]
        [TestCaseSource("EncodeCases")]
        public void EncodeDecode(ProtoMsg message)
        {
            var target = new ProtoProc();
            var encode = target.Encode(message);
            var decode = target.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        [Test]
        public void EncodeFailed()
        {
            var proc = new ProtoProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                proc.Encode(null);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                proc.Encode(new object());
            });
        }

        [Test]
        public void DecodeFailed()
        {
            var target = new ProtoProc();

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
        public void Process(ProtoMsg message)
        {
            var target = new ProtoProc();
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
            var target = new ProtoProc();

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
                target.Process(new ProtoMsg { MessageID = 1 });
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
            Assert.Throws<ArgumentException>(() =>
            {
                JsonProc.Unmarshal<ProtoTest>(new object(), out var _, out var _);
            });
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(ProtoProc.Marshal(1, new ProtoTest { Data = "test1" }));
                yield return new TestCaseData(ProtoProc.Marshal(2, new ProtoTest { Data = "test2" }));
            }
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(ProtoProc.Marshal(1, new ProtoTest { Data = "test1" }));
                yield return new TestCaseData(ProtoProc.Marshal(2, new ProtoTest { Data = "test2" }));
            }
        }

        public static IEnumerable MarshalCases
        {
            get
            {
                yield return new TestCaseData(1, new ProtoTest { Data = "test1" });
                yield return new TestCaseData(2, new ProtoTest { Data = "test2" });
            }
        }
    }
}
