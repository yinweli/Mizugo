using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;

namespace Mizugo
{
    /// <summary>
    /// TCP客戶端組件
    /// </summary>
    internal class TCPClient
    {
        public TCPClient(string ip, int port, Logger logger = null)
        {
            this.ip = ip;
            this.port = port;
            // TODO: process object
            this.logger = logger;
        }

        /// <summary>
        /// 連線處理
        /// </summary>
        public async void Connect()
        {
            if (client != null)
                throw new Exception("not disconnect");

            var host = await Dns.GetHostAddressesAsync(ip);

            if (host.Length <= 0)
                throw new Exception("host address failed");

            client = new TcpClient();
            client.NoDelay = true;

            await client.ConnectAsync(host[0], port);

            if (client.Connected == false)
                throw new Exception("connect failed");

            stream = client.GetStream();
            recvStream = new StreamReader(stream);
            sendStream = new StreamWriter(stream);
            recvThread = new Thread(recvLoop);
            sendThread = new Thread(sendLoop);

            info(
                "TCPClient connect success { Host: {0}:{1}, NoDelay: {2}, ReceiveTimeout: {3}, ReceiveBufferSize: {4}, SendTimeout: {5}, SendBufferSize: {6} }",
                ip,
                port,
                client.NoDelay,
                client.ReceiveTimeout,
                client.ReceiveBufferSize,
                client.SendTimeout,
                client.SendBufferSize
            );
        }

        /// <summary>
        /// 斷線處理
        /// </summary>
        public void Disconnect()
        {
            close = true;

            if (recvThread != null)
            {
                try
                {
                    recvThread.Join();
                    recvThread = null;
                } // try
                catch (Exception e)
                {
                    error("TCPClient disconnect failed: {0}", e.ToString());
                } // catch
            } // if

            if (sendThread != null)
            {
                try
                {
                    sendThread.Join();
                    sendThread = null;
                } // try
                catch (Exception e)
                {
                    error("TCPClient disconnect failed: {0}", e.ToString());
                } // catch
            } // if

            if (recvStream != null)
            {
                try
                {
                    recvStream.Close();
                    recvStream = null;
                } // try
                catch (Exception e)
                {
                    error("TCPClient disconnect failed: {0}", e.ToString());
                } // catch
            } // if

            if (sendStream != null)
            {
                try
                {
                    sendStream.Close();
                    sendStream = null;
                } // try
                catch (Exception e)
                {
                    error("TCPClient disconnect failed: {0}", e.ToString());
                } // catch
            } // if

            if (client != null)
            {
                client.Close();
                client = null;
            } // if

            close = false;

            info("TCPClient disconnect success");
        }

        public void Update()
        {
            // TODO: update
        }

        public void Send()
        {
            // TODO: send
        }

        private void recvLoop()
        {
            // TODO: recv
        }

        private void sendLoop()
        {
            // TODO: send
        }

        /// <summary>
        /// 記錄一般訊息
        /// </summary>
        /// <param name="message">訊息字串</param>
        private void info(string format, params object[] args)
        {
            if (logger != null)
                logger.Info(new StringBuilder().AppendFormat(format, args).ToString());
        }

        /// <summary>
        /// 記錄錯誤訊息
        /// </summary>
        /// <param name="message">訊息字串</param>
        private void error(string format, params object[] args)
        {
            if (logger != null)
                logger.Error(new StringBuilder().AppendFormat(format, args).ToString());
        }

        /// <summary>
        /// 連線位址
        /// </summary>
        public string ip = string.Empty;

        /// <summary>
        /// 連線埠號
        /// </summary>
        public int port = 0;

        // TODO: process object

        /// <summary>
        /// 日誌物件
        /// </summary>
        private Logger logger = null;

        /// <summary>
        /// 連線物件
        /// </summary>
        private TcpClient client = null;

        /// <summary>
        /// 網路資料流物件
        /// </summary>
        private NetworkStream stream = null;

        /// <summary>
        /// 網路讀取流物件
        /// </summary>
        private StreamReader recvStream = null;

        /// <summary>
        /// 網路寫入流物件
        /// </summary>
        private StreamWriter sendStream = null;

        /// <summary>
        /// 接收執行緒物件
        /// </summary>
        private Thread recvThread = null;

        /// <summary>
        /// 傳送執行緒物件
        /// </summary>
        private Thread sendThread = null;

        /// <summary>
        /// 關閉旗標
        /// </summary>
        private volatile bool close = false;
    }
}
