using System;
using System.IO;
using System.Security.Cryptography;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// des-cbc編碼/解碼
    /// </summary>
    public class DesCBC : ICodec
    {
        public DesCBC(PaddingMode padding, string key, string iv)
        {
            this.padding = padding;
            this.key = Encoding.UTF8.GetBytes(key);
            this.iv = Encoding.UTF8.GetBytes(iv);
        }

        public object Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            using var provider = new DESCryptoServiceProvider();

            provider.Key = key;
            provider.IV = iv;
            provider.Mode = CipherMode.CBC;
            provider.Padding = padding;

            using var memoryStream = new MemoryStream();
            using var cryptoStream = new CryptoStream(memoryStream, provider.CreateEncryptor(), CryptoStreamMode.Write);

            cryptoStream.Write(temp, 0, temp.Length);
            cryptoStream.FlushFinalBlock();
            return memoryStream.ToArray();
        }

        public object Decode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            using var provider = new DESCryptoServiceProvider();

            provider.Key = key;
            provider.IV = iv;
            provider.Mode = CipherMode.CBC;
            provider.Padding = padding;

            using var memoryStream = new MemoryStream();
            using var cryptoStream = new CryptoStream(memoryStream, provider.CreateDecryptor(), CryptoStreamMode.Write);

            cryptoStream.Write(temp, 0, temp.Length);
            cryptoStream.FlushFinalBlock();
            return memoryStream.ToArray();
        }

        private readonly PaddingMode padding = PaddingMode.PKCS7;
        private readonly byte[] key = null;
        private readonly byte[] iv = null;
    }
}
