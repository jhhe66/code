package main

import "ice"
import "fmt"
import "ice_warper/protobuf/Broadcast"
import "github.com/golang/protobuf/proto"
import "time"

func ValidApps() {
	ice_handle := ice.New("--Ice.Config=./client.conf", "ValidUsersProxy")
	if ice_handle.IceHandle == nil {
		fmt.Printf("init ice failed.\n")
		return 
	}
	
	defer ice_handle.IceFree()

	query := new(Broadcast.Validates)
	id := proto.Int64(10000)
	query.ReqNo = id
	appkey := []byte("1fa5281c6c1d4cf5bb0bbbe0")
	platform := Broadcast.PLATFORM_IPHONE
	apps := &Broadcast.StatApps{}
	apps.Appkey = appkey
	apps.Platform = &platform
query.Apps = append(query.Apps, apps)

	request, err := proto.Marshal(query)
	if err != nil {
		fmt.Printf("Marshal failed. %v \n", err)
		return 
	}
	
	response := make([]byte, 512)
	_, rsz := ice_handle.VaildApps(request, response)

	response_obj := &Broadcast.Validates{}
	err = proto.Unmarshal(response[:rsz], response_obj)
	if err != nil {
		fmt.Printf("Unmarshal failed. %v \n", err)
		return 		
	}

	fmt.Printf("UsersLen: %d\n", len(response_obj.GetApps()))
	fmt.Printf("Users: %t\n", response_obj.GetApps()[0].GetValid())	
	switch response_obj.GetApps()[0].GetValid() {
		case 1:
			fmt.Printf("App had users\n")
		case 0:
			fmt.Printf("App had'nt users\n")
	}
}

func ValidUsers() {
	ice_handle := ice.New("--Ice.Config=./client.conf", "ValidUsersProxy")
	if ice_handle.IceHandle == nil {
		fmt.Printf("init ice failed.\n")
		return 
	}
	
	defer ice_handle.IceFree()

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

	request, err := proto.Marshal(query)
	if err != nil {
		fmt.Printf("Marshal failed. %v \n", err)
		return 
	}
	
	response := make([]byte, 512)
	_, rsz := ice_handle.VaildUsers(request, response)

	response_obj := &Broadcast.Validates{}
	err = proto.Unmarshal(response[:rsz], response_obj)
	if err != nil {
		fmt.Printf("Unmarshal failed. %v \n", err)
		return 		
	}

	fmt.Printf("UsersLen: %d\n", len(response_obj.GetUsers()))
	//fmt.Printf("Users: %d\n", response_obj.GetUsers()[0].GetUsers()[0].GetValid())
	switch response_obj.GetUsers()[0].GetUsers()[0].GetValid() {
		case Broadcast.VALIDSTATUS_VALID_OK:
			fmt.Printf("Is App Users.\n")
		default:
			fmt.Printf("Is not App Users.\n")
	}
}

func main() {
	go ValidApps()
	fmt.Println("\n\n")
	time.Sleep(1 * time.Second)

	go ValidUsers()
	
	time.Sleep(10 * time.Second)
}
