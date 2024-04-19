using System.Collections;
using System.Text;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestBase64
    {
        [Test]
        [TestCaseSource("Base64Cases")]
        public void Base64_(object input)
        {
            var target = new Base64();
            var encode = target.Encode(input);
            var output = target.Decode(encode);

            Assert.AreEqual(input, output);
        }

        public static IEnumerable Base64Cases
        {
            get
            {
                yield return new TestCaseData(Encoding.UTF8.GetBytes("testdata"));
                yield return new TestCaseData(Encoding.UTF8.GetBytes("somedata"));
            }
        }
    }
}
