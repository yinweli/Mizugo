using Newtonsoft.Json;
using System;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

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

            process(message);
        }

        /// <summary>
        /// json訊息序列化
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="input">輸入物件</param>
        /// <returns>訊息物件</returns>
        public static JsonMsg Marshal(MessageID messageID, object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            var json = JsonConvert.SerializeObject(input);
            var jsonBytes = Encoding.UTF8.GetBytes(json);

            return new JsonMsg { MessageID = messageID, Message = jsonBytes };
        }

        /// <summary>
        /// json訊息反序列化
        /// </summary>
        /// <typeparam name="T">訊息類型</typeparam>
        /// <param name="input">輸入物件</param>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        public static void Unmarshal<T>(object input, out MessageID messageID, out T message)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg jsonMsg)
                throw new InvalidMessageException("unmarshal");

            var json = Encoding.UTF8.GetString(jsonMsg.Message);

            messageID = jsonMsg.MessageID;
            message = JsonConvert.DeserializeObject<T>(json);
        }
    }
}
