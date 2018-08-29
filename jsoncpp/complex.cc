#include <stdio.h>
#include <stdlib.h>
#include <json/json.h>


#define __POUT__(node) do {								\
	Json::FastWriter writer;							\
	printf("%s\n", writer.write(node).c_str());	\
} while(0)			

#define __SOUT__(node) printf("%s\n", node.toStyledString().c_str())

int
main(int argc, char** argv)
{
	Json::Value root;
	Json::Value cluster;
	Json::Value clusters;
	Json::Value node;
	Json::Value slot;
	
	for (unsigned int idc = 0; idc < 2; idc++) {
		Json::Value nodes;
		cluster["id"] = idc;
		cluster["stat"] = 0;
		for (unsigned int idn = 0; idn < 2; idn++) {
			Json::Value slots;
			node["id"] = idn;
			node["stat"] = 0;
			for (unsigned int idx = 0; idx < 2; idx++) {
				slot["id"] = idx;
				slot["stat"] = 0;
				slot["users"] = idx * 10000;
				slots.append(slot);
			}
			node["slots"] = slots;
			nodes.append(node);
		}
		cluster["nodes"] = nodes;
		clusters.append(cluster);
	}

	root["clusters"] = clusters;
	
	//printf("JSON: %s\n", node.toStyledString().c_str());
	__POUT__(root);

	return 0;
}
