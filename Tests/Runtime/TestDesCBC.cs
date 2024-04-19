using System.Collections;
using System.Security.Cryptography;
using System.Text;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestDesCBC
    {
        [Test]
        [TestCaseSource("DesCBCCases")]
        public void DesCBCZeros(object input)
        {
            var target = new DesCBC(PaddingMode.Zeros, key, iv);
            var encode = target.Encode(input);
            var output = target.Decode(encode);

            Assert.AreEqual(input, output);
        }

        [Test]
        [TestCaseSource("DesCBCCases")]
        public void DesCBCPKCS7(object input)
        {
            var target = new DesCBC(PaddingMode.PKCS7, key, iv);
            var encode = target.Encode(input);
            var output = target.Decode(encode);

            Assert.AreEqual(input, output);
        }

        public static IEnumerable DesCBCCases
        {
            get
            {
                yield return new TestCaseData(Encoding.UTF8.GetBytes("testdata"));
                yield return new TestCaseData(Encoding.UTF8.GetBytes("somedata"));
            }
        }

        private readonly string key = "thisakey"; // 密鑰必須為8位
        private readonly string iv = "this-iv-"; // 初始向量必須為8位
    }
}
