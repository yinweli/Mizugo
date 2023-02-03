using NUnit.Framework;
using System.Text;

namespace Mizugo
{
    internal class TestBase64
    {
        [Test]
        public void Base64_()
        {
            var encode = Base64.Encode(src);
            var output = Base64.Decode(encode);

            Assert.AreEqual(src, output);
        }

        [SetUp]
        public void SetUp()
        {
            src = Encoding.UTF8.GetBytes("testdata");
        }

        private byte[] src = null;
    }
}
