#ifndef __TCPCLIENT_H_
#define __TCPCLIENT_H_

#include <string>

using namespace std;

class CTcpClient
{
public:
    CTcpClient(string host, unsigned short port);
    ~CTcpClient();

    /*
	        发送完size大小的缓冲内容才返回
	        返回-1异常
	        返回0 对方关闭连接
	 */
    int Send(const char* buffer, int size);

	 /*
	        接收完size大小的缓冲内容才返回
	        返回-1异常
	        返回0 对方关闭连接
	 */
    int Recv(char* buffer, int size);

    /*
	        连接远端地址
	        返回 -1异常 
	        返回 0 连接成功 
	 */
    int Connect();

    int Disconneced();

	int SetSendTimeout(unsigned int sec);

	int SetRecvTimeout(unsigned int sec);

	int SetRecvBuffSize(unsigned int size);
	int SetSendBuffSize(unsigned int size);

private:
    int	socketfd_;
    string	host_;
    unsigned short port_;

    bool IsConnected;

	 int Socket();

    int Close();
};

#endif
