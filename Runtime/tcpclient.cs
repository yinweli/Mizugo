using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;

namespace Mizugo
{
    /// <summary>
    /// TCP�Ȥ�ݲե�
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
        /// �s�u�B�z
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
        /// �_�u�B�z
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
        /// �O���@��T��
        /// </summary>
        /// <param name="message">�T���r��</param>
        private void info(string format, params object[] args)
        {
            if (logger != null)
                logger.Info(new StringBuilder().AppendFormat(format, args).ToString());
        }

        /// <summary>
        /// �O�����~�T��
        /// </summary>
        /// <param name="message">�T���r��</param>
        private void error(string format, params object[] args)
        {
            if (logger != null)
                logger.Error(new StringBuilder().AppendFormat(format, args).ToString());
        }

        /// <summary>
        /// �s�u��}
        /// </summary>
        public string ip = string.Empty;

        /// <summary>
        /// �s�u��
        /// </summary>
        public int port = 0;

        // TODO: process object

        /// <summary>
        /// ��x����
        /// </summary>
        private Logger logger = null;

        /// <summary>
        /// �s�u����
        /// </summary>
        private TcpClient client = null;

        /// <summary>
        /// ������Ƭy����
        /// </summary>
        private NetworkStream stream = null;

        /// <summary>
        /// ����Ū���y����
        /// </summary>
        private StreamReader recvStream = null;

        /// <summary>
        /// �����g�J�y����
        /// </summary>
        private StreamWriter sendStream = null;

        /// <summary>
        /// �������������
        /// </summary>
        private Thread recvThread = null;

        /// <summary>
        /// �ǰe���������
        /// </summary>
        private Thread sendThread = null;

        /// <summary>
        /// �����X��
        /// </summary>
        private volatile bool close = false;
    }
}
