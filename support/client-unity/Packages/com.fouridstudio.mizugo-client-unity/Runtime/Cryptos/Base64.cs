using System;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// base64編解碼
    /// </summary>
    public class Base64
    {
        /// <summary>
        /// 編碼
        /// </summary>
        /// <param name="src">來源資料</param>
        /// <returns>結果資料</returns>
        public static byte[] Encode(byte[] src)
        {
            var base64String = Convert.ToBase64String(src);
            var base64Bytes = Encoding.UTF8.GetBytes(base64String);
            return base64Bytes;
        }

        /// <summary>
        /// 解碼
        /// </summary>
        /// <param name="src">來源資料</param>
        /// <returns>結果資料</returns>
        public static byte[] Decode(byte[] src)
        {
            var base64String = Encoding.UTF8.GetString(src);
            var base64Bytes = Convert.FromBase64String(base64String);
            return base64Bytes;
        }
    }
}
