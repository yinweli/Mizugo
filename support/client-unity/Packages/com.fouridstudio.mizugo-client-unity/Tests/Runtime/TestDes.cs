using NUnit.Framework;
using System.Security.Cryptography;
using System.Text;

namespace Mizugo
{
    internal class TestDes
    {
        [Test]
        [TestCase(PaddingMode.Zeros)]
        [TestCase(PaddingMode.PKCS7)]
        public void DesECB_(PaddingMode padding)
        {
            var crypto = DesECB.Encrypt(padding, key, src);
            var output = DesECB.Decrypt(padding, key, crypto);

            Assert.AreEqual(src, output);
        }

        [Test]
        [TestCase(PaddingMode.Zeros)]
        [TestCase(PaddingMode.PKCS7)]
        public void DesCBC_(PaddingMode padding)
        {
            var crypto = DesCBC.Encrypt(padding, key, iv, src);
            var output = DesCBC.Decrypt(padding, key, iv, crypto);

            Assert.AreEqual(src, output);
        }

        private byte[] key = Encoding.UTF8.GetBytes("thisakey"); // 密鑰必須為8位
        private byte[] iv = Encoding.UTF8.GetBytes("this-iv-"); // 初始向量必須為8位
        private byte[] src = Encoding.UTF8.GetBytes("testdata");
    }
}
