using System;
using System.Linq;
using System.Text;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using UnityEditor.PackageManager.Requests;
using UnityEngine.UIElements;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// Raven系列工具專為伺服器與客戶端之間的訊息傳遞協議設計
    /// 它基於 Clients/ProcProto.cs 中的 proto 處理器, 利用 Msgs/Raven.cs 中的 RavenQ 和 RavenA 對基礎協議的訊息內容進行細化處理
    /// 其中, 客戶端向伺服器發送的訊息採用 RavenQ 格式, 而伺服器向客戶端返回的訊息則採用 RavenA 格式
    /// </summary>
    public class Raven
    {
        public static object RavenQBuilder(MessageID messageID, IMessage header, IMessage request)
        {
            var message = new RavenQ() { Header = Any.Pack(header), Request = Any.Pack(request) };
            return ProcProto.Marshal(messageID, message);
        }

        public static RavenQData<H, Q> RavenQParser<H, Q>(object input)
            where H : IMessage, new()
            where Q : IMessage, new()
        {
            ProcProto.Unmarshal<RavenQ>(input, out var messageID, out var message);
            return new RavenQData<H, Q>
            {
                messageID = messageID,
                header = message.Header.Unpack<H>(),
                request = message.Request.Unpack<Q>()
            };
        }

        public static object RavenABuilder(MessageID messageID, int errID, IMessage header, IMessage request, params IMessage[] respond)
        {
            var message = new RavenA()
            {
                ErrID = errID,
                Header = Any.Pack(header),
                Request = Any.Pack(request),
                Respond = { respond.Select(Any.Pack) },
            };
            return ProcProto.Marshal(messageID, message);
        }

        public static RavenAData<H, Q> RavenAParser<H, Q>(object input)
            where H : IMessage, new()
            where Q : IMessage, new()
        {
            ProcProto.Unmarshal<RavenA>(input, out var messageID, out var message);
            return new RavenAData<H, Q>
            {
                messageID = messageID,
                errID = message.ErrID,
                header = message.Header.Unpack<H>(),
                request = message.Request.Unpack<Q>(),
                respond = message.Respond.ToArray(),
            };
        }
    }

    /// <summary>
    /// RavenQ資料
    /// </summary>
    public class RavenQData<H, Q>
        where H : IMessage, new()
        where Q : IMessage, new()
    {
        public MessageID messageID; // 訊息編號
        public H header; // 標頭資料
        public Q request; // 要求資料

        public string Detail()
        {
            var leno = sizeof(MessageID);
            var lenh = header.CalculateSize();
            var lenq = request.CalculateSize();
            var builder = new StringBuilder();

            builder.AppendLine(">> message");
            builder.AppendLine($"    messageID: {messageID}");
            builder.AppendLine($"    header: {JsonFormatter.ToDiagnosticString(header)}");
            builder.AppendLine($"    request: {JsonFormatter.ToDiagnosticString(request)}");
            builder.AppendLine(">> size");
            builder.AppendLine($"    other: {leno}");
            builder.AppendLine($"    header: {lenh}");
            builder.AppendLine($"    request: {lenq}");
            builder.AppendLine($"    total: {leno + lenh + lenq}");
            return builder.ToString();
        }
    }

    /// <summary>
    /// RavenA資料
    /// </summary>
    public class RavenAData<H, Q>
        where H : IMessage, new()
        where Q : IMessage, new()
    {
        public MessageID messageID; // 訊息編號
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
                return respond[index].Unpack<T>();

            return default;
        }

        public string Detail()
        {
            var leno = sizeof(MessageID) + sizeof(int);
            var lenh = header.CalculateSize();
            var lenq = request.CalculateSize();
            var lent = leno + lenh + lenq;
            var builder = new StringBuilder();

            builder.AppendLine(">> message");
            builder.AppendLine($"    messageID: {messageID}");
            builder.AppendLine($"    errID: {errID}");
            builder.AppendLine($"    header: {JsonFormatter.ToDiagnosticString(header)}");
            builder.AppendLine($"    request: {JsonFormatter.ToDiagnosticString(request)}");

            for (var i = 0; i < respond.Length; i++)
                builder.AppendLine($"    respond[{i}]: {JsonFormatter.ToDiagnosticString(respond[i])}");

            builder.AppendLine(">> size");
            builder.AppendLine($"    other: {leno}");
            builder.AppendLine($"    header: {lenh}");
            builder.AppendLine($"    request: {lenq}");

            for (var i = 0; i < respond.Length; i++)
            {
                var lens = respond[i].CalculateSize();
                lent += lens;
                builder.AppendLine($"    respond[{i}]: {lens}");
            } // for

            builder.AppendLine($"    total: {lent}");
            return builder.ToString();
        }
    }
}
