using System;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// proto處理器, 封包結構使用ProtoMsg
    /// 訊息定義: support/proto/mizugo/protomsg.proto
    /// </summary>
    public partial class ProtoProc : Procmgr, ICodec
    {
        public object Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not ProtoMsg temp)
                throw new ArgumentException("input");

            return temp.ToByteArray();
        }

        public object Decode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            return ProtoMsg.Parser.ParseFrom(temp);
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not ProtoMsg message)
                throw new ArgumentException("input");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnprocessException(message.MessageID);

            process(message);
        }
    }

    public partial class ProtoProc
    {
        /// <summary>
        /// proto訊息序列化
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        /// <returns>訊息物件</returns>
        public static ProtoMsg Marshal(MessageID messageID, IMessage message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            return new ProtoMsg { MessageID = messageID, Message = Any.Pack(message) };
        }

        /// <summary>
        /// proto訊息反序列化
        /// </summary>
        /// <typeparam name="T">訊息類型</typeparam>
        /// <param name="input">輸入物件</param>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        public static void Unmarshal<T>(object input, out MessageID messageID, out T message)
            where T : IMessage, new()
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not ProtoMsg data)
                throw new ArgumentException("input");

            messageID = data.MessageID;
            message = data.Message.Unpack<T>();
        }
    }
}
