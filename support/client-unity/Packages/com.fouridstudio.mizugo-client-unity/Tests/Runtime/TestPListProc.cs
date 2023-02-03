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

    internal class TestPListProc
    {
        [Test]
        [TestCaseSource("EncodeCases")]
        public void Encode(PListMsg message)
        {
            var proc = new PListProc { KeyStr = key, IVStr = iv };
            var encode = proc.Encode(message);
            var decode = proc.Decode(encode);

            Assert.IsTrue(TestUtil.EqualsByJson(message, decode));
        }

        public static IEnumerable EncodeCases
        {
            get
            {
                yield return new TestCaseData(PListProc.Marshal(1, new PListTest { Data = "test1" }));
                yield return new TestCaseData(PListProc.Marshal(2, new PListTest { Data = "test2" }, 3, new PListTest { Data = "test3" }));
            }
        }

        [Test]
        public void EncodeFailed()
        {
            var proc = new PListProc { KeyStr = key, IVStr = iv };

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
            var proc = new PListProc { KeyStr = key, IVStr = iv };

            Assert.Throws<ArgumentNullException>(() =>
            {
                proc.Decode(null);
            });
        }

        [Test]
        [TestCaseSource("ProcessCases")]
        public void Process(PListMsg message)
        {
            var proc = new PListProc();
            var valid = true;

            foreach (var itor in message.Messages)
            {
                proc.Add(
                    itor.MessageID,
                    (object param) =>
                    {
                        valid &= TestUtil.EqualsByJson(itor, param);
                    }
                );
            } // for

            proc.Process(message);
            Assert.IsTrue(valid);
        }

        public static IEnumerable ProcessCases
        {
            get
            {
                yield return new TestCaseData(PListProc.Marshal(1, new PListTest { Data = "test1" }));
                yield return new TestCaseData(PListProc.Marshal(2, new PListTest { Data = "test2" }, 3, new PListTest { Data = "test3" }));
            }
        }

        [Test]
        public void ProcessFailed()
        {
            var proc = new PListProc();

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
                proc.Process(PListProc.Marshal(1, new PListTest()));
            });
        }

        [Test]
        public void Misc()
        {
            var proc = new PListProc { KeyStr = key, IVStr = iv };

            Assert.AreEqual(Encoding.UTF8.GetBytes(key), proc.Key);
            Assert.AreEqual(Encoding.UTF8.GetBytes(iv), proc.IV);
            Assert.AreEqual(key, proc.KeyStr);
            Assert.AreEqual(iv, proc.IVStr);
        }

        [Test]
        [TestCaseSource("MarshalSenderCases")]
        public void MarshalSender(MessageID messageID, PListTest message)
        {
            var sender = new PListSender();

            sender.Add(messageID, message);

            var marshal = PListProc.Marshal(sender);

            PListProc.Unmarshal<PListTest>(marshal.Messages[0], out var resultID, out var result);
            Assert.AreEqual(messageID, resultID);
            Assert.IsTrue(TestUtil.EqualsByJson(message, result));
        }

        public static IEnumerable MarshalSenderCases
        {
            get { yield return new TestCaseData(1, new PListTest { Data = "test1" }); }
        }

        [Test]
        public void MarshalSenderFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                PListProc.Marshal((PListSender)null);
            });
        }

        [Test]
        [TestCaseSource("MarshalListCases")]
        public void MarshalList(object[] message)
        {
            var marshal = PListProc.Marshal(message);

            for (var i = 0; i < marshal.Messages.Count; i++)
            {
                PListProc.Unmarshal<PListTest>(marshal.Messages[i], out var resultID, out var result);
                Assert.AreEqual(message[i * 2], resultID);
                Assert.IsTrue(TestUtil.EqualsByJson(message[i * 2 + 1], result));
            } // for
        }

        public static IEnumerable MarshalListCases
        {
            get
            {
                yield return new object[]
                {
                    1,
                    new PListTest { Data = "test1" }
                };
                yield return new object[]
                {
                    2,
                    new PListTest { Data = "test2" },
                    3,
                    new PListTest { Data = "test3" }
                };
            }
        }

        [Test]
        public void MarshalListFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                PListProc.Marshal((object[])null);
            });
            Assert.Throws<ArgumentOutOfRangeException>(() =>
            {
                PListProc.Marshal(1, 2, 3);
            });
            Assert.Throws<InvalidCastException>(() =>
            {
                PListProc.Marshal(new object(), null, new object(), null);
            });
            Assert.Throws<InvalidCastException>(() =>
            {
                PListProc.Marshal(1, null, 2, null);
            });
        }

        [Test]
        public void UnmarshalFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                PListProc.Unmarshal<PListTest>(null, out var _, out var _);
            });
            Assert.Throws<InvalidMessageException>(() =>
            {
                PListProc.Unmarshal<PListTest>(new object(), out var _, out var _);
            });
        }

        private string key = "thisakey";
        private string iv = "this-iv-";
    }

    internal class TestPListSender
    {
        [Test]
        public void Misc()
        {
            var sender = new PListSender();

            Assert.Throws<ArgumentNullException>(() =>
            {
                sender.Add(1, null);
            });
            Assert.AreEqual(sender, sender.Add(1, new PListTest()));
            Assert.NotNull(sender.Marshal());
        }
    }
}
