#include "SimpleAmqpClient/SimpleAmqpClient.h"
#include <iostream>
#include <stdio.h>
#include <unistd.h>

using namespace AmqpClient;

int main() 
{
	AmqpClient::Channel::ptr_t _channel = Channel::CreateFromUri("amqp://push:TMQ.Push@183.232.25.113:5672/");

	if (!_channel->IsConnected()) {
		printf("rabbitmq connect failed.\n");
		return -1;
	}

	sleep(300);

	return 0;
}
