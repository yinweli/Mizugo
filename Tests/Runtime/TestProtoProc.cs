using NUnit.Framework;
using System;
using System.Collections;

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
        public void Encode(ProtoMsg input)
        {
            var protoproc = new ProtoProc();
            var encode = protoproc.Encode(input);
            var decode = protoproc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(input, decode));
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(ProtoProc.Marshal(1, new ProtoTest() { Data = "test1" }));
                yield return new TestCaseData(ProtoProc.Marshal(2, new ProtoTest() { Data = "test2" }));
            }
        }

        [Test]
        public void EncodeFailed()
        {
            var protoproc = new ProtoProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                protoproc.Encode(null);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                protoproc.Encode(new object());
            });
        }

        [Test]
        public void DecodeFailed()
        {
            var protoproc = new ProtoProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                protoproc.Decode(null);
            });
        }

        [Test]
        [TestCaseSource("ProcessCases")]
        public void Process(ProtoMsg protoMsg)
        {
            var protoproc = new ProtoProc();
            var expected = protoMsg;
            var valid = false;

            protoproc.Add(
                protoMsg.MessageID,
                (object param) =>
                {
                    valid = expected == param;
                }
            );
            protoproc.Process(protoMsg);
            Assert.IsTrue(valid);
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(ProtoProc.Marshal(1, new ProtoTest() { Data = "test1" }));
                yield return new TestCaseData(ProtoProc.Marshal(2, new ProtoTest() { Data = "test2" }));
            }
        }

        [Test]
        public void ProcessFailed()
        {
            var protoproc = new ProtoProc();

            Assert.Throws<ArgumentNullException>(() =>
            {
                protoproc.Process(null);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                protoproc.Process(new object());
            });
            Assert.Throws<UnprocessException>(() =>
            {
                protoproc.Process(new ProtoMsg() { MessageID = 1 });
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
                yield return new TestCaseData(1, new ProtoTest() { Data = "test1" });
                yield return new TestCaseData(2, new ProtoTest() { Data = "test2" });
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
    }
}
