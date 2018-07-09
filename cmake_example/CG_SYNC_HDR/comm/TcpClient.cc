#include "TcpClient.h"

#include <sys/types.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <time.h>

CTcpClient::CTcpClient(string host, unsigned short port)
	:host_(host),
    port_(port),
	 IsConnected(false)
{
    Socket();
}

CTcpClient::~CTcpClient()
{
	Close();
}

int
CTcpClient::Send(const char* buffer, int size)
{
    int cLen = 0;
    
	if(IsConnected)
	{
		cLen = send(socketfd_, buffer, size, 0);
	}

	if(cLen == 0)
		IsConnected = false;

	return cLen;
}

int
CTcpClient::Recv(char* buffer, int size)
{
    if(!IsConnected)
	{
        return -1;
	}

	int iLen = 0;
	
	iLen = recv(socketfd_, buffer, size, 0);

	if (iLen == 0)
		IsConnected = false;

    return iLen;
}

//private

int
CTcpClient::Socket()
{
    socketfd_ = socket(AF_INET, SOCK_STREAM, 0);

    if (socketfd_ == -1)
	{
        return 0;
	}
    else
	{
        return -1;
	}
}

int
CTcpClient::Connect()
{
    int result = 0;

	if(IsConnected)
		return 0;

	sockaddr_in remote;

    remote.sin_addr.s_addr = inet_addr(host_.c_str());
    remote.sin_family = AF_INET;
    remote.sin_port = htons(port_);

	 if(this->socketfd_ > 0)
	 {
        result = connect(socketfd_, (struct sockaddr*)&remote, sizeof(struct sockaddr));

        IsConnected = result == 0 ? true : false;
	 }

    return result;
}

int
CTcpClient::Close()
{
    int result = ::close(socketfd_);

    return result;
}

int
CTcpClient::SetSendTimeout(unsigned int sec)
{
	struct timeval timeout = {sec, 0};

	return setsockopt(socketfd_, SOL_SOCKET, SO_SNDTIMEO, (char *)&timeout, sizeof(struct timeval));
}

int
CTcpClient::SetRecvTimeout(unsigned int sec)
{
	struct timeval timeout = {sec, 0};

	return setsockopt(socketfd_, SOL_SOCKET, SO_RCVTIMEO, (char *)&timeout, sizeof(struct timeval));
}

int	
CTcpClient::SetRecvBuffSize(unsigned int size)
{
	return setsockopt(socketfd_, SOL_SOCKET, SO_RCVBUF, (char*)size, sizeof(size));
}

int	
CTcpClient::SetSendBuffSize(unsigned int size)
{
	return setsockopt(socketfd_, SOL_SOCKET, SO_RCVBUF, (char*)size, sizeof(size));
}
