package main

import "fmt"
import "ice_tag"
import "ice_tagalias/protobuf/TagAliasBatchQuery"
import "ice_tagalias/protobuf/tagalias_common"
import "github.com/golang/protobuf/proto"
import "time"

func check_user() {
	ice_handle := ice_tag.New("--Ice.Config=./client.conf", "TagAliasAccessProxy")
	if ice_handle.IceHandle == nil {
		fmt.Printf("init ice failed.\n");
		return
	}

	defer ice_handle.IceFree();
	
	query := new(TagAliasBatchQuery.TagAliasQuery)
	query.ReqNo = proto.Int64(10000)
	appkey := []byte("1fa5281c6c1d4cf5bb0bbbe0")
	query.Appkey = appkey
	platform := tagalias_common.PLATFORM_IPHONE
	query.Platform = &platform
	ops := &TagAliasBatchQuery.OpsUnit{}
	cmd := tagalias_common.QUERY_ACTION_CHECK_TAGS_HAS_USERS
	ops.Cmd = &cmd
	tag := []byte("aaaa")
	ops.Tags = append(ops.Tags, tag)
	query.Query = append(query.Query, ops)
	
	request, err := proto.Marshal(query)
	if err != nil {
		fmt.Printf("Marshal failed. %v \n", err)
		return 
	}
	
	response := make([]byte, 512)
	_, rsz := ice_handle.Request(request, response)	
	
	fmt.Printf("rsz: %d\n", rsz)
	fmt.Printf("resonpse: %v\n", response[:rsz])

	response_obj := new(TagAliasBatchQuery.TagAlsResp)
	err = proto.Unmarshal(response[:rsz], response_obj)
	if err != nil {
		fmt.Printf("Unmarshal failed. %v \n", err)
		return 		
	}
	
	fmt.Printf("Code: %d\n", response_obj.Code)
}

func validateTags() {
	ice_handle := ice_tag.New("--Ice.Config=./client.conf", "TagAliasAccessProxy")
	if ice_handle.IceHandle == nil {
		fmt.Printf("init ice failed.\n");
		return
	}

	defer ice_handle.IceFree();
	
	query := new(TagAliasBatchQuery.TagAliasQuery)
	query.ReqNo = proto.Int64(10000)
	appkey := []byte("1fa5281c6c1d4cf5bb0bbbe0")
	query.Appkey = appkey
	platform := tagalias_common.PLATFORM_IPHONE
	query.Platform = &platform
	ops := &TagAliasBatchQuery.OpsUnit{}
	cmd := tagalias_common.QUERY_ACTION_CHECK_TAGS_HAS_USERS
	ops.Cmd = &cmd
	tag := []byte("aaaa")
	ops.Tags = append(ops.Tags, tag)
	query.Query = append(query.Query, ops)
	
	request, err := proto.Marshal(query)
	if err != nil {
		fmt.Printf("Marshal failed. %v \n", err)
		return 
	}

	response := make([]byte, 512)
	_, rsz := ice_handle.ValidateTags(request, response)	
	
	fmt.Printf("rsz: %d\n", rsz)
	fmt.Printf("resonpse: %v\n", response[:rsz])

	response_obj := new(TagAliasBatchQuery.TagAlsResp)
	err = proto.Unmarshal(response[:rsz], response_obj)
	if err != nil {
		fmt.Printf("Unmarshal failed. %v \n", err)
		return 		
	}	
}

func main() {
	check_user()
	//validateTags()
	time.Sleep(60 * time.Second)
}


