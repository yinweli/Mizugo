using System;
using System.Collections;
using NUnit.Framework;

namespace Mizugo
{
    using MessageID = Int32;

    internal class TestRaven
    {
        [Test]
        [TestCaseSource("RavenQCases")]
        public void RavenQ(MessageID messageID, RavenTest header, RavenTest request)
        {
            var output1 = Raven.RavenQBuilder(messageID, header, request);
            var output2 = Raven.RavenQParser<RavenTest, RavenTest>(output1);

            Assert.AreEqual(messageID, output2.messageID);
            Assert.AreEqual(header, output2.header);
            Assert.AreEqual(request, output2.request);
            UnityEngine.Debug.Log(output2.Detail());
        }

        [Test]
        public void RavenQFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                Raven.RavenQBuilder(0, null, new RavenTest());
            });
            Assert.Throws<ArgumentNullException>(() =>
            {
                Raven.RavenQBuilder(0, new RavenTest(), null);
            });
            Assert.Throws<ArgumentNullException>(() =>
            {
                Raven.RavenQParser<RavenTest, RavenTest>(null);
            });
        }

        [Test]
        [TestCaseSource("RavenACases")]
        public void RavenA(MessageID messageID, int errID, RavenTest header, RavenTest request, params RavenTest[] respond)
        {
            var output1 = Raven.RavenABuilder(messageID, errID, header, request, respond);
            var output2 = Raven.RavenAParser<RavenTest, RavenTest>(output1);

            Assert.AreEqual(messageID, output2.messageID);
            Assert.AreEqual(errID, output2.errID);
            Assert.AreEqual(header, output2.header);
            Assert.AreEqual(request, output2.request);
            Assert.AreEqual(respond[0], output2.GetRespond<RavenTest>());
            Assert.AreEqual(respond[1], output2.GetRespondAt<RavenTest>(1));
            UnityEngine.Debug.Log(output2.Detail());
        }

        [Test]
        public void RavenAFailed()
        {
            Assert.Throws<ArgumentNullException>(() =>
            {
                Raven.RavenABuilder(0, 0, null, new RavenTest());
            });
            Assert.Throws<ArgumentNullException>(() =>
            {
                Raven.RavenABuilder(0, 0, new RavenTest(), null);
            });
            Assert.Throws<ArgumentNullException>(() =>
            {
                Raven.RavenABuilder(0, 0, new RavenTest(), null, null, new RavenTest());
            });
            Assert.Throws<ArgumentNullException>(() =>
            {
                Raven.RavenQParser<RavenTest, RavenTest>(null);
            });
        }

        public static IEnumerable RavenQCases
        {
            get
            {
                yield return new object[]
                {
                    1,
                    new RavenTest() { Data = "header" },
                    new RavenTest() { Data = "request" },
                };
            }
        }

        public static IEnumerable RavenACases
        {
            get
            {
                yield return new object[]
                {
                    1,
                    2,
                    new RavenTest() { Data = "header" },
                    new RavenTest() { Data = "request" },
                    new RavenTest[]
                    {
                        new RavenTest() { Data = "respond1" },
                        new RavenTest() { Data = "respond2" },
                    },
                };
            }
        }
    }
}
