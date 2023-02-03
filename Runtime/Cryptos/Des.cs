using System.IO;
using System.Security.Cryptography;

namespace Mizugo
{
    /// <summary>
    /// des-ecb加解密
    /// </summary>
    public class DesECB
    {
        /// <summary>
        /// 加密
        /// </summary>
        /// <param name="padding">填充方式</param>
        /// <param name="key">密鑰</param>
        /// <param name="src">來源資料</param>
        /// <returns>結果資料</returns>
        public static byte[] Encrypt(PaddingMode padding, byte[] key, byte[] src)
        {
            using var provider = new DESCryptoServiceProvider();

            provider.Key = key;
            provider.Mode = CipherMode.ECB;
            provider.Padding = padding;

            using var memoryStream = new MemoryStream();
            using var cryptoStream = new CryptoStream(memoryStream, provider.CreateEncryptor(), CryptoStreamMode.Write);

            cryptoStream.Write(src, 0, src.Length);
            cryptoStream.FlushFinalBlock();
            return memoryStream.ToArray();
        }

        /// <summary>
        /// 解密
        /// </summary>
        /// <param name="padding">填充方式</param>
        /// <param name="key">密鑰</param>
        /// <param name="src">來源資料</param>
        /// <returns>結果資料</returns>
        public static byte[] Decrypt(PaddingMode padding, byte[] key, byte[] src)
        {
            using var provider = new DESCryptoServiceProvider();

            provider.Key = key;
            provider.Mode = CipherMode.ECB;
            provider.Padding = padding;

            using var memoryStream = new MemoryStream();
            using var cryptoStream = new CryptoStream(memoryStream, provider.CreateDecryptor(), CryptoStreamMode.Write);

            cryptoStream.Write(src, 0, src.Length);
            cryptoStream.FlushFinalBlock();
            return memoryStream.ToArray();
        }
    }

    /// <summary>
    /// des-cbc加解密
    /// </summary>
    public class DesCBC
    {
        /// <summary>
        /// 加密
        /// </summary>
        /// <param name="padding">填充方式</param>
        /// <param name="key">密鑰</param>
        /// <param name="iv">初始向量</param>
        /// <param name="src">來源資料</param>
        /// <returns>結果資料</returns>
        public static byte[] Encrypt(PaddingMode padding, byte[] key, byte[] iv, byte[] src)
        {
            using var provider = new DESCryptoServiceProvider();

            provider.Key = key;
            provider.IV = iv;
            provider.Mode = CipherMode.CBC;
            provider.Padding = padding;

            using var memoryStream = new MemoryStream();
            using var cryptoStream = new CryptoStream(memoryStream, provider.CreateEncryptor(), CryptoStreamMode.Write);

            cryptoStream.Write(src, 0, src.Length);
            cryptoStream.FlushFinalBlock();
            return memoryStream.ToArray();
        }

        /// <summary>
        /// 解密
        /// </summary>
        /// <param name="padding">填充方式</param>
        /// <param name="key">密鑰</param>
        /// <param name="iv">初始向量</param>
        /// <param name="src">來源資料</param>
        /// <returns>結果資料</returns>
        public static byte[] Decrypt(PaddingMode padding, byte[] key, byte[] iv, byte[] src)
        {
            using var provider = new DESCryptoServiceProvider();

            provider.Key = key;
            provider.IV = iv;
            provider.Mode = CipherMode.CBC;
            provider.Padding = padding;

            using var memoryStream = new MemoryStream();
            using var cryptoStream = new CryptoStream(memoryStream, provider.CreateDecryptor(), CryptoStreamMode.Write);

            cryptoStream.Write(src, 0, src.Length);
            cryptoStream.FlushFinalBlock();
            return memoryStream.ToArray();
        }
    }
}
