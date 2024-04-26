using System;
using System.Collections;
using Google.Protobuf;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestProcProto
    {
        [Test]
        [TestCaseSource("EncodeCases")]
        public void Encode(Proto input)
        {
            var target = new ProcProto();
            var encode = target.Encode(input);
            var decode = target.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(input, decode));
        }

        [Test]
        public void EncodeFailed()
        {
            var target = new ProcProto();

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
            var target = new ProcProto();

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
        public void Process(Proto input)
        {
            var target = new ProcProto();
            var valid = false;

            target.Add(
                input.MessageID,
                (object param) =>
                {
                    valid = TestUtil.EqualsByJson(input, param);
                }
            );
            target.Process(input);
            Assert.IsTrue(valid);
        }

        [Test]
        public void ProcessFailed()
        {
            var target = new ProcProto();

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
                target.Process(new Proto { MessageID = 1 });
            });
        }

        [Test]
        [TestCaseSource("MarshalCases")]
        public void Marshal(int messageID, IMessage message)
        {
            var marshal = ProcProto.Marshal(messageID, message);

            ProcProto.Unmarshal<ProtoTest>(marshal, out var resultID, out var result);
            Assert.AreEqual(messageID, resultID);
            Assert.IsTrue(TestUtil.EqualsByJson(message, result));
        }

        [Test]
        public void MarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProcProto.Marshal(1, null);
            });
        }

        [Test]
        public void UnmarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProcJson.Unmarshal<ProtoTest>(null, out var _, out var _);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                ProcJson.Unmarshal<ProtoTest>(new object(), out var _, out var _);
            });
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(ProcProto.Marshal(1, new ProtoTest { Data = "test1" }));
                yield return new TestCaseData(ProcProto.Marshal(2, new ProtoTest { Data = "test2" }));
            }
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(ProcProto.Marshal(1, new ProtoTest { Data = "test1" }));
                yield return new TestCaseData(ProcProto.Marshal(2, new ProtoTest { Data = "test2" }));
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
