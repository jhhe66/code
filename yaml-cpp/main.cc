#include <stdio.h>
#include <yaml-cpp/yaml.h>
#include <string>

using namespace YAML;
using std::string;

static void
__dump_jmq(const Node* node)
{
	Node jmq;

	if (node == NULL) return;
	
	jmq = *node;

	printf("id: %s\n", jmq["id"].as<string>().c_str());
	if (jmq["type"].IsDefined()) {
		printf("has type\n");
	}
	printf("srv-server: %s\n", jmq["srv"]["server"].as<string>().c_str());
	printf("vhost-user: %s\n", jmq["vhost"]["user"].as<string>().c_str());
	printf("vhost-password: %s\n", jmq["vhost"]["password"].as<string>().c_str());
	printf("exchange-content_encoding: %s\n", jmq["exchange"]["content_encoding"].as<string>().c_str());

}

static void
__read_jmq(const char* file)
{
	Node config;

	config = LoadFile(file);

	__dump_jmq(&config);
}

int
main(int argc, char** argv)
{
	//Node config = LoadFile("../python/map.yml");
	//Node config = LoadFile("../python/test.yml");

	//printf("type: %s\n", config["type"].as<string>().c_str());
	//printf("%s\n", config[0]["name"].as<string>().c_str());
	
	__read_jmq("jmq.yml");
	
	return 0;
}
