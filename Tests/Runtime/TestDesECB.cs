using System.Collections;
using System.Security.Cryptography;
using System.Text;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestDesECB
    {
        [Test]
        [TestCaseSource("DesECBCases")]
        public void DesECBZeros(object input)
        {
            var target = new DesECB(PaddingMode.Zeros, key);
            var encode = target.Encode(input);
            var output = target.Decode(encode);

            Assert.AreEqual(input, output);
        }

        [Test]
        [TestCaseSource("DesECBCases")]
        public void DesECBPKCS7(object input)
        {
            var target = new DesECB(PaddingMode.PKCS7, key);
            var encode = target.Encode(input);
            var output = target.Decode(encode);

            Assert.AreEqual(input, output);
        }

        public static IEnumerable DesECBCases
        {
            get
            {
                yield return new TestCaseData(Encoding.UTF8.GetBytes("testdata"));
                yield return new TestCaseData(Encoding.UTF8.GetBytes("somedata"));
            }
        }

        private readonly string key = "thisakey"; // 密鑰必須為8位
    }
}
