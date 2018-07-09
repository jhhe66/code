#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include <stdio.h>
#include <iostream>
using namespace rapidjson;
int main() {
    // 1. 把 JSON 解析至 DOM。
    const char* json = "{\"sdk_ver_android\": \"1.6.0\", \"sdk_ver_ios\": \"1.4.0\", \"sdk_ver_wp\": \"1.0.3\"}";
	Document doc;

    doc.Parse(json);
    if (doc.HasParseError())
    {
        printf("parse json failed,errorcode:%d", doc.GetParseError());
        return -1;
    }

    //doc.HasMember("sdk_ver_android")&&doc["sdk_ver_android"].IsString() ? 
	//	printf("sdk_ver_android: %s", doc["sdk_ver_android"].GetString()) : 0;
    //doc.HasMember("sdk_ver_ios")&&doc["sdk_ver_ios"].IsString() ? 
	//	printf("sdk_ver_ios: %s", doc["sdk_ver_ios"].GetString()) : 0;
    //doc.HasMember("sdk_ver_wp")&&doc["sdk_ver_wp"].IsString() ? 
	//	printf("sdk_ver_wp: %s", doc["sdk_ver_wp"].GetString()) : 0;
	Value stat_sdk;

	stat_sdk.SetString("1.2.0", doc.GetAllocator());

	doc.AddMember("sdk_stat_sdk", stat_sdk, doc.GetAllocator());
	//doc["sdk_stat_sdk"] = "1.2.0;

	StringBuffer buffer;
	Writer<StringBuffer> writer(buffer);

	doc.Accept(writer);

	printf("JSON: %s\n", buffer.GetString());

    return 0;
}
