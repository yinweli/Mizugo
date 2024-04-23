using System;
using System.Text;
using Newtonsoft.Json;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// json處理器, 封包結構使用JsonMs
    /// 訊息定義: support/proto/mizugo/msg-go/msgs-json/jsonmsg.go
    /// 訊息定義: support/proto/mizugo/msg-cs/msgs-json/Jsonmsg.cs
    /// </summary>
    public partial class ProcJson : Procmgr, ICodec
    {
        public object Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg temp)
                throw new ArgumentException("input");

            return Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(temp));
        }

        public object Decode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            return JsonConvert.DeserializeObject<JsonMsg>(Encoding.UTF8.GetString(temp));
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg message)
                throw new ArgumentException("input");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnprocessException(message.MessageID);

            process(message);
        }
    }

    public partial class ProcJson
    {
        public static JsonMsg Marshal(MessageID messageID, object message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            var json = JsonConvert.SerializeObject(message);
            var jsonBytes = Encoding.UTF8.GetBytes(json);
            return new JsonMsg { MessageID = messageID, Message = jsonBytes };
        }

        public static void Unmarshal<T>(object input, out MessageID messageID, out T message)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg data)
                throw new ArgumentException("input");

            var json = Encoding.UTF8.GetString(data.Message);

            messageID = data.MessageID;
            message = JsonConvert.DeserializeObject<T>(json);
        }
    }
}
