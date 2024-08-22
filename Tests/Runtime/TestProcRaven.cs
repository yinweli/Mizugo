using System;
using System.Collections;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestProcRaven
    {
        [Test]
        [TestCaseSource("EncodeCases")]
        public void Encode(RavenS input)
        {
            var target = new ProcRaven();
            var encode = target.Encode(input);

            Assert.AreEqual(input.ToByteArray(), encode);
        }

        [Test]
        public void EncodeFailed()
        {
            var target = new ProcRaven();

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
        [TestCaseSource("DecodeCases")]
        public void Decode(RavenC input)
        {
            var target = new ProcRaven();
            var decode = target.Decode(input.ToByteArray());

            Assert.IsTrue(TestUtil.EqualsByJson(input, decode));
        }

        [Test]
        public void DecodeFailed()
        {
            var target = new ProcRaven();

            Assert.Throws<ArgumentNullException>(() =>
            {
                target.Decode(null);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                target.Decode(new object());
            });
        }

        [Test]
        [TestCaseSource("ProcessCases")]
        public void Process(RavenC input)
        {
            var target = new ProcRaven();
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
            var target = new ProcRaven();

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
                target.Process(new RavenC { MessageID = 1 });
            });
        }

        [Test]
        [TestCaseSource("MarshalCases")]
        public void Marshal(int messageID, IMessage header, IMessage request)
        {
            var marshal = ProcRaven.Marshal(messageID, header, request);

            Assert.AreEqual(messageID, marshal.MessageID);
            Assert.AreEqual(Any.Pack(header), marshal.Header);
            Assert.AreEqual(Any.Pack(request), marshal.Request);
        }

        [Test]
        public void MarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProcRaven.Marshal(1, null, new RavenTest());
            });
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProcRaven.Marshal(1, new RavenTest(), null);
            });
        }

        [Test]
        [TestCaseSource("UnmarshalCases")]
        public void Unmarshal(RavenC input)
        {
            ProcRaven.Unmarshal<RavenTest, RavenTest>(input, out var result);
            Assert.AreEqual(input.MessageID, result.messageID);
            Assert.AreEqual(input.ErrID, result.errID);
            Assert.AreEqual(input.Header, Any.Pack(result.header));
            Assert.AreEqual(input.Request, Any.Pack(result.request));
            Assert.AreEqual(input.Respond[0], Any.Pack(result.GetRespond<RavenTest>()));
            Assert.AreEqual(input.Respond[1], Any.Pack(result.GetRespondAt<RavenTest>(1)));
            Assert.Null(result.GetRespond<ProtoTest>());
            Assert.Null(result.GetRespondAt<ProtoTest>(0));
            Assert.Null(result.GetRespondAt<RavenTest>(3));
        }

        [Test]
        public void UnmarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                ProcRaven.Unmarshal<RavenTest, RavenTest>(null, out var _);
            });
            Assert.Throws<ArgumentException>(() =>
            {
                ProcRaven.Unmarshal<RavenTest, RavenTest>(new object(), out var _);
            });
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(
                    new RavenS()
                    {
                        MessageID = 1,
                        Header = Any.Pack(new RavenTest() { Data = "header1" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest1" }),
                    }
                );
                yield return new TestCaseData(
                    new RavenS()
                    {
                        MessageID = 2,
                        Header = Any.Pack(new RavenTest() { Data = "header2" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest2" }),
                    }
                );
            }
        }

        public static IEnumerable DecodeCases
        {
            get
            {
                yield return new TestCaseData(
                    new RavenC()
                    {
                        MessageID = 1,
                        ErrID = 1,
                        Header = Any.Pack(new RavenTest() { Data = "header1" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest1" }),
                        Respond = { Any.Pack(new RavenTest() { Data = "respond1" }), Any.Pack(new RavenTest() { Data = "respond2" }) },
                    }
                );
                yield return new TestCaseData(
                    new RavenC()
                    {
                        MessageID = 2,
                        ErrID = 2,
                        Header = Any.Pack(new RavenTest() { Data = "header2" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest2" }),
                        Respond = { Any.Pack(new RavenTest() { Data = "respond1" }), Any.Pack(new RavenTest() { Data = "respond2" }) },
                    }
                );
            }
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(
                    new RavenC()
                    {
                        MessageID = 1,
                        ErrID = 1,
                        Header = Any.Pack(new RavenTest() { Data = "header1" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest1" }),
                        Respond = { Any.Pack(new RavenTest() { Data = "respond1" }), Any.Pack(new RavenTest() { Data = "respond2" }) },
                    }
                );
                yield return new TestCaseData(
                    new RavenC()
                    {
                        MessageID = 2,
                        ErrID = 2,
                        Header = Any.Pack(new RavenTest() { Data = "header2" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest2" }),
                        Respond = { Any.Pack(new RavenTest() { Data = "respond1" }), Any.Pack(new RavenTest() { Data = "respond2" }) },
                    }
                );
            }
        }

        public static IEnumerable MarshalCases
        {
            get
            {
                yield return new TestCaseData(1, new RavenTest() { Data = "header1" }, new RavenTest() { Data = "request1" });
                yield return new TestCaseData(2, new RavenTest() { Data = "header2" }, new RavenTest() { Data = "request2" });
            }
        }

        public static IEnumerable UnmarshalCases
        {
            get
            {
                yield return new TestCaseData(
                    new RavenC()
                    {
                        MessageID = 1,
                        ErrID = 1,
                        Header = Any.Pack(new RavenTest() { Data = "header1" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest1" }),
                        Respond = { Any.Pack(new RavenTest() { Data = "respond1" }), Any.Pack(new RavenTest() { Data = "respond2" }) },
                    }
                );
                yield return new TestCaseData(
                    new RavenC()
                    {
                        MessageID = 2,
                        ErrID = 2,
                        Header = Any.Pack(new RavenTest() { Data = "header2" }),
                        Request = Any.Pack(new RavenTest() { Data = "reuest2" }),
                        Respond = { Any.Pack(new RavenTest() { Data = "respond1" }), Any.Pack(new RavenTest() { Data = "respond2" }) },
                    }
                );
            }
        }
    }
}
