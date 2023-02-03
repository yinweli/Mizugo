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

        [SetUp]
        public void SetUp()
        {
            key = Encoding.UTF8.GetBytes("thisakey"); // 密鑰必須為8位
            iv = Encoding.UTF8.GetBytes("this-iv-"); // 初始向量必須為8位
            src = Encoding.UTF8.GetBytes("testdata");
        }

        private byte[] key = null;
        private byte[] iv = null;
        private byte[] src = null;
    }
}
