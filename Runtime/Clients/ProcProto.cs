using System;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;

namespace Mizugo
{
    /// <summary>
    /// proto處理器, 封包結構使用Proto
    /// 訊息定義: support/proto-mizugo/proto.proto
    /// </summary>
    public partial class ProcProto : Procmgr, ICodec
    {
        public object Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not Proto temp)
                throw new ArgumentException("input");

            return temp.ToByteArray();
        }

        public object Decode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            return Proto.Parser.ParseFrom(temp);
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not Proto message)
                throw new ArgumentException("input");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnprocessException(message.MessageID);

            process(message);
        }
    }

    public partial class ProcProto
    {
        public static Proto Marshal(int messageID, IMessage message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            return new Proto { MessageID = messageID, Message = Any.Pack(message) };
        }

        public static void Unmarshal<T>(object input, out int messageID, out T message)
            where T : IMessage, new()
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not Proto temp)
                throw new ArgumentException("input");

            messageID = temp.MessageID;
            message = temp.Message.Unpack<T>();
        }
    }
}
