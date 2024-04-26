using System;
using System.Text;
using Newtonsoft.Json;

namespace Mizugo
{
    /// <summary>
    /// json處理器, 封包結構使用Json
    /// 訊息定義: support/proto-mizugo/msg-go/msgs-json/json.go
    /// 訊息定義: support/proto-mizugo/msg-cs/msgs-json/Json.cs
    /// </summary>
    public partial class ProcJson : Procmgr, ICodec
    {
        public object Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not Json temp)
                throw new ArgumentException("input");

            return Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(temp));
        }

        public object Decode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            return JsonConvert.DeserializeObject<Json>(Encoding.UTF8.GetString(temp));
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not Json message)
                throw new ArgumentException("input");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnprocessException(message.MessageID);

            process(message);
        }
    }

    public partial class ProcJson
    {
        public static Json Marshal(int messageID, object message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            var json = JsonConvert.SerializeObject(message);
            var jsonBytes = Encoding.UTF8.GetBytes(json);
            return new Json { MessageID = messageID, Message = jsonBytes };
        }

        public static void Unmarshal<T>(object input, out int messageID, out T message)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not Json temp)
                throw new ArgumentException("input");

            var json = Encoding.UTF8.GetString(temp.Message);

            messageID = temp.MessageID;
            message = JsonConvert.DeserializeObject<T>(json);
        }
    }
}
