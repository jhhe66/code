package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lice_warper -lIce -lIceUtil

#include "ice_warper.h"
#include <stdlib.h>
#include <stdio.h>
*/
import "C"
import "fmt"
import "unsafe"
import "ice_warper/protobuf/Broadcast"
//import "ice_warper/protobuf/QueryUser"
import "github.com/golang/protobuf/proto"
import "time"

func ValidApps() {
	var	req [512]byte
	var	res [512]byte
	var rsz *C.uint
	var conf_file string

	conf_file = "--Ice.Config=./client.conf"
	
	conf_path := C.CString(conf_file)
	ice_handle := C.ice_warper_init(conf_path)

	defer C.free(unsafe.Pointer(conf_path))
	defer C.ice_warper_free(ice_handle)
	
	if ice_handle == nil {
		fmt.Printf("ice init failed. \n");
		return 
	}
	
	query := new(Broadcast.Validates)
	id := proto.Int64(10000)
	query.ReqNo = id
	appkey := []byte("1fa5281c6c1d4cf5bb0bbbe0")
	platform := Broadcast.PLATFORM_IPHONE
	apps := &Broadcast.StatApps{}
	apps.Appkey = appkey
	apps.Platform = &platform

	query.Apps = append(query.Apps, apps)
	
//	query := &Broadcast.Validates{
//		ReqNo : proto.Int64(10000),
//		Apps : []* Broadcast.StatApps{
//			{ Appkey : []byte("hello"), Platform : 0,},
//		},
//	}
	
	request, err := proto.Marshal(query)
	if err != nil {
		fmt.Printf("Marshal failed. %v \n", err)
		return 
	}

//	resonpse_2 := new(Broadcast.Validates)
//	err = proto.Unmarshal(request, resonpse_2)
//	if err != nil {
//		fmt.Printf("Unmarshal resonpse_2 failed. %v \n", err)
//		return 				
//	}
//	fmt.Printf("ReqNo:%d\n", resonpse_2.GetReqNo())
//	fmt.Printf("Appkey:%s\n", resonpse_2.GetApps()[0].GetAppkey())
//	fmt.Printf("Appkey:%s\n", resonpse_2.GetApps()[0].GetPlatform())

	fmt.Printf("Len: %d\n", copy(req[:], request))
	
	fmt.Printf("req: %v\n", req)
	rsz = new(C.uint)
	*rsz = C.uint(len(res))
	//fmt.Printf("%d\n", C.ice_warper_valid_appkeys(ice_handle, (*C.char)(unsafe.Pointer(&req[0])), (*C.char)(unsafe.Pointer(&res[0])), C.uint(len(res))))
	fmt.Printf("%d\n", C.ice_warper_valid_appkeys(ice_handle, (*C.char)(unsafe.Pointer(&req[0])), (*C.char)(unsafe.Pointer(&res[0])), rsz))
	
	fmt.Printf("res: %s\n", res)
	fmt.Printf("res: %v\n", res)
	fmt.Printf("res_len: %v\n", *rsz)

	response := &Broadcast.Validates{}
	err = proto.Unmarshal(res[:*rsz], response)
	if err != nil {
		fmt.Printf("Unmarshal failed. %v \n", err)
		return 		
	}
	
	fmt.Printf("UsersLen: %d\n", len(response.GetApps()))
	fmt.Printf("Users: %d\n", response.GetApps()[0].GetValid())	
}

func ValidUsers() {
	var	req [512]byte
	var	res [512]byte
	var rsz *C.uint
	var conf_file string

	conf_file = "--Ice.Config=./client.conf"
	
	conf_path := C.CString(conf_file)
	ice_handle := C.ice_warper_init(conf_path)

	defer C.free(unsafe.Pointer(conf_path))
	defer C.ice_warper_free(ice_handle)
	
	if ice_handle == nil {
		fmt.Printf("ice init failed. \n");
		return 
	}
	
	query := new(Broadcast.Validates)
	id := proto.Int64(20000)
	query.ReqNo = id
	appkey := []byte("1fa5281c6c1d4cf5bb0bbbe0")
	platform := Broadcast.PLATFORM_IPHONE
	users := &Broadcast.StatUsers{}
	users.Appkey = appkey
	users.Platform = &platform
	user := &Broadcast.StatUser{}
	user.Uid = proto.Int64(999999)
	
	users.Users = append(users.Users, user)
	query.Users = append(query.Users, users)
	
//	query := &Broadcast.Validates{
//		ReqNo : proto.Int64(10000),
//		Apps : []* Broadcast.StatApps{
//			{ Appkey : []byte("hello"), Platform : 0,},
//		},
//	}
	
	request, err := proto.Marshal(query)
	if err != nil {
		fmt.Printf("Marshal failed. %v \n", err)
		return 
	}

//	resonpse_2 := new(Broadcast.Validates)
//	err = proto.Unmarshal(request, resonpse_2)
//	if err != nil {
//		fmt.Printf("Unmarshal resonpse_2 failed. %v \n", err)
//		return 				
//	}
//	fmt.Printf("ReqNo:%d\n", resonpse_2.GetReqNo())
//	fmt.Printf("Appkey:%s\n", resonpse_2.GetApps()[0].GetAppkey())
//	fmt.Printf("Appkey:%s\n", resonpse_2.GetApps()[0].GetPlatform())

	fmt.Printf("Len: %d\n", copy(req[:], request))
	
	fmt.Printf("req: %v\n", req)
	rsz = new(C.uint)
	*rsz = C.uint(len(res))
	//fmt.Printf("%d\n", C.ice_warper_valid_appkeys(ice_handle, (*C.char)(unsafe.Pointer(&req[0])), (*C.char)(unsafe.Pointer(&res[0])), C.uint(len(res))))
	fmt.Printf("%d\n", C.ice_warper_valid_users(ice_handle, (*C.char)(unsafe.Pointer(&req[0])), (*C.char)(unsafe.Pointer(&res[0])), rsz))
	
	fmt.Printf("res: %s\n", res)
	fmt.Printf("res: %v\n", res)
	fmt.Printf("res_len: %v\n", *rsz)

	response := &Broadcast.Validates{}
	err = proto.Unmarshal(res[:*rsz], response)
	if err != nil {
		fmt.Printf("Unmarshal failed. %v \n", err)
		return 		
	}
	
	fmt.Printf("UsersLen: %d\n", len(response.GetUsers()))
	fmt.Printf("Users: %d\n", response.GetUsers()[0].GetUsers()[0].GetValid())			
}

func main() {
//	var	req [512]byte
//	var	res [512]byte
//	var rsz *C.uint
//
//	conf_path := C.CString("--Ice.Config=./client.conf")
//	ice_handle := C.ice_warper_init(conf_path)
//
//	defer C.free(unsafe.Pointer(conf_path))
//	defer C.ice_warper_free(unsafe.Pointer(ice_handle))
//	
//	if ice_handle == nil {
//		fmt.Printf("ice init failed. \n");
//		return 
//	}
//	
//	query := new(Broadcast.Validates)
//	id := proto.Int64(10000)
//	query.ReqNo = id
//	appkey := []byte("1fa5281c6c1d4cf5bb0bbbe01")
//	platform := Broadcast.PLATFORM_IPHONE
//	apps := &Broadcast.StatApps{}
//	apps.Appkey = appkey
//	apps.Platform = &platform
//
//	query.Apps = append(query.Apps, apps)
//	
////	query := &Broadcast.Validates{
////		ReqNo : proto.Int64(10000),
////		Apps : []* Broadcast.StatApps{
////			{ Appkey : []byte("hello"), Platform : 0,},
////		},
////	}
//	
//	request, err := proto.Marshal(query)
//	if err != nil {
//		fmt.Printf("Marshal failed. %v \n", err)
//		return 
//	}
//
////	resonpse_2 := new(Broadcast.Validates)
////	err = proto.Unmarshal(request, resonpse_2)
////	if err != nil {
////		fmt.Printf("Unmarshal resonpse_2 failed. %v \n", err)
////		return 				
////	}
////	fmt.Printf("ReqNo:%d\n", resonpse_2.GetReqNo())
////	fmt.Printf("Appkey:%s\n", resonpse_2.GetApps()[0].GetAppkey())
////	fmt.Printf("Appkey:%s\n", resonpse_2.GetApps()[0].GetPlatform())
//
//	fmt.Printf("Len: %d\n", copy(req[:], request))
//	
//	fmt.Printf("req: %v\n", req)
//	rsz = new(C.uint)
//	*rsz = C.uint(len(res))
//	//fmt.Printf("%d\n", C.ice_warper_valid_appkeys(ice_handle, (*C.char)(unsafe.Pointer(&req[0])), (*C.char)(unsafe.Pointer(&res[0])), C.uint(len(res))))
//	fmt.Printf("%d\n", C.ice_warper_valid_appkeys(ice_handle, (*C.char)(unsafe.Pointer(&req[0])), (*C.char)(unsafe.Pointer(&res[0])), rsz))
//	
//	fmt.Printf("res: %s\n", res)
//	fmt.Printf("res: %v\n", res)
//	fmt.Printf("res_len: %v\n", *rsz)
//
//	response := &Broadcast.Validates{}
//	err = proto.Unmarshal(res[:*rsz], response)
//	if err != nil {
//		fmt.Printf("Unmarshal failed. %v \n", err)
//		return 		
//	}
//	
//	fmt.Printf("UsersLen: %d\n", len(response.GetApps()))
//	fmt.Printf("Users: %d\n", response.GetApps()[0].GetValid())
	
	
	//ValidApps()
	go ValidUsers()
	time.Sleep(10 * time.Second)
}
