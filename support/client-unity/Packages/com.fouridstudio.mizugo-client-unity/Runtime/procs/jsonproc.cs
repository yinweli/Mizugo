using Newtonsoft.Json;
using System;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// json處理器
    /// </summary>
    public class JsonProc : Procmgr
    {
        public override byte[] Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg message)
                throw new InvalidMessageException("encode");

            var json = JsonConvert.SerializeObject(message);
            var jsonBytes = Encoding.UTF8.GetBytes(json);
            var base64 = Convert.ToBase64String(jsonBytes);
            var base64Bytes = Encoding.UTF8.GetBytes(base64);

            return base64Bytes;
        }

        public override object Decode(byte[] input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            var base64 = Encoding.UTF8.GetString(input);
            var base64Bytes = Convert.FromBase64String(base64);
            var json = Encoding.UTF8.GetString(base64Bytes);
            var message = JsonConvert.DeserializeObject<JsonMsg>(json);

            return message;
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg message)
                throw new InvalidMessageException("process");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnProcessException(message.MessageID);

            process(message.Message);
        }
    }
}
