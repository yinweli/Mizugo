using System;
using System.Linq;
using System.Text;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;

namespace Mizugo
{
    /// <summary>
    /// raven處理器, 封包結構使用RavenS, RavenC
    /// 訊息定義: support/proto-mizugo/raven.proto
    /// </summary>
    public partial class ProcRaven : Procmgr, ICodec
    {
        public object Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not RavenS temp)
                throw new ArgumentException("input");

            return temp.ToByteArray();
        }

        public object Decode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not byte[] temp)
                throw new ArgumentException("input");

            return RavenC.Parser.ParseFrom(temp);
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not RavenC message)
                throw new ArgumentException("input");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnprocessException(message.MessageID);

            process(message);
        }
    }

    public partial class ProcRaven
    {
        public static RavenS Marshal(int messageID, IMessage header, IMessage request)
        {
            if (header == null)
                throw new ArgumentNullException("header");

            if (request == null)
                throw new ArgumentNullException("request");

            return new RavenS
            {
                MessageID = messageID,
                Header = Any.Pack(header),
                Request = Any.Pack(request),
            };
        }

        public static void Unmarshal<H, Q>(object input, out RavenData<H, Q> output)
            where H : IMessage, new()
            where Q : IMessage, new()
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not RavenC temp)
                throw new ArgumentException("input");

            output = new RavenData<H, Q>
            {
                messageID = temp.MessageID,
                errID = temp.ErrID,
                header = temp.Header.Unpack<H>(),
                request = temp.Request.Unpack<Q>(),
                respond = temp.Respond.ToArray(),
            };
        }
    }

    public class RavenData<H, Q>
        where H : IMessage, new()
        where Q : IMessage, new()
    {
        public int messageID; // 訊息編號
        public int errID; // 錯誤編號
        public H header; // 標頭資料
        public Q request; // 要求資料
        public Any[] respond; // 回應列表

        public T GetRespond<T>()
            where T : IMessage, new()
        {
            foreach (var itor in respond)
            {
                if (itor.TryUnpack<T>(out var result))
                    return result;
            } // for

            return default;
        }

        public T GetRespondAt<T>(int index)
            where T : IMessage, new()
        {
            if (index < respond.Length)
            {
                if (respond[index].TryUnpack<T>(out var result))
                    return result;
            } // if

            return default;
        }

        public int Size()
        {
            var size = sizeof(int) * 2;

            size += header.CalculateSize();
            size += request.CalculateSize();
            size += respond.Sum(itor => itor.CalculateSize());
            return size;
        }

        public string Detail()
        {
            var builder = new StringBuilder();

            builder.AppendLine($"messageID: {messageID}");
            builder.AppendLine($"errID: {errID}");
            builder.AppendLine($"header: {JsonFormatter.ToDiagnosticString(header)}");
            builder.AppendLine($"request: {JsonFormatter.ToDiagnosticString(request)}");

            for (var i = 0; i < respond.Length; i++)
                builder.AppendLine($"respond[{i}]: {JsonFormatter.ToDiagnosticString(respond[i])}");

            return builder.ToString();
        }
    }
}
