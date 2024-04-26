using System;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// base64編碼/解碼
    /// </summary>
    public class Base64 : ICodec
    {
        public object Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            var base64String = Convert.ToBase64String(temp);
            var base64Bytes = Encoding.UTF8.GetBytes(base64String);
            return base64Bytes;
        }

        public object Decode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            var base64String = Encoding.UTF8.GetString(temp);
            var base64Bytes = Convert.FromBase64String(base64String);
            return base64Bytes;
        }
    }
}
