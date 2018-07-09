package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/libcomm"
	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

type Callback struct{}

func newPool(redisServer string) *redis.Pool {
	return &redis.Pool{MaxIdle: 10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisServer)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

var (
	infoLog *log.Logger
	db      *sql.DB
	pool    *redis.Pool
)

const (
	kTTLInterval = 2592000 // TTL life time = 30*24*3600 second
)

const (
	CMD_UPDATE_LOCATION       = 0x2007
	CMD_UPDATE_LOCATION_ACK   = 0x2008
	CMD_GET_USER_LOCATION     = 0x2009
	CMD_GET_USER_LOCATION_ACK = 0x200A
	CMD_REFRESH_LOCATION      = 0x2035
	CMD_REFRESH_LOCATION_ACK  = 0x2036
)
const (
	DB_RET_SUCCESS     = 0
	DB_RET_EXEC_FAILED = 1
	DB_RET_NOT_EXIST   = 100
)

const (
	RET_SUCCESS            = 0
	ERR_SYSERR_START       = 100
	ERR_SERVER_BUSY        = 100
	ERR_INTERNAL_ERROR     = 101
	ERR_UNFORMATTED_PACKET = 102
	ERR_NO_ACCESS          = 103
	ERR_INVALID_CLIENT     = 104
	ERR_INVALID_SESSION    = 105
	ERR_INVALID_PARAM      = 106
)

const (
	STRGOOGLE           = "google"
	LOCATIONTOPLACEHMAP = "location#placeid"
	MULTILANGNAMEHMAP   = "country#lang"
	MULTIPLACEINFOHMAP  = "place#lang"
)

var languageShortName map[string]string = map[string]string{"Chinese": "zh",
	"English":    "en",
	"Japanese":   "ja",
	"Korean":     "ko",
	"Spanish":    "es",
	"Portugues":  "pt",
	"French":     "fr",
	"German":     "de",
	"Italian":    "it",
	"Russian":    "ru",
	"Arabic":     "ar",
	"Chinese_yy": "zh-TW"}

type tagPlaceInformation struct {
	Country      []byte
	Admin1       []byte
	Admin2       []byte
	Admin3       []byte
	Locality     []byte
	SubLocality  []byte
	Neighborhood []byte
}

type tagLocationInfo struct {
	Uid        uint32
	Allowed    uint8
	TimeStamp  uint64
	PlaceId    uint32
	Latitude   []byte
	Longtitude []byte
	Source     []byte
	PlaceInfo  tagPlaceInformation
}

func GetLocationPlaceIdKey(placeInfo *tagPlaceInformation) string {
	return fmt.Sprintf("#%s#%s#%s#%s#%s#%s#%s",
		placeInfo.Country,
		placeInfo.Admin1,
		placeInfo.Admin2,
		placeInfo.Admin3,
		placeInfo.Locality,
		placeInfo.SubLocality,
		placeInfo.Neighborhood)
}

func GetMultiCountryNameKey(countryName string, languageType string) string {
	return fmt.Sprintf("#%s#%s", countryName, languageType)
}

func GetMultiPlaceInfoKey(placeId int64, shortName string) string {
	return fmt.Sprintf("%v#%s", placeId, shortName)
}

func GetUserPlaceKey(uid uint32) string {
	return fmt.Sprintf("%v#place", uid)
}

func MatchLocationPlaceIdInRedis(englishPlace *tagPlaceInformation) (placeId int64) {
	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	field := GetLocationPlaceIdKey(englishPlace)
	r, err := redisConn.Do("HGET", LOCATIONTOPLACEHMAP, field)
	placeId, err = redis.Int64(r, err)
	if err != nil {
		infoLog.Printf("Reids exec HGET field=%s err=%v", field, err)
		placeId = 0
		return
	} else {
		infoLog.Println("Redis exec HGet succ placeId =", placeId)
		return
	}
}

func MatchLocationPlaceId(locationInfo *tagLocationInfo) (placeId int64) {
	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	field := GetLocationPlaceIdKey(&locationInfo.PlaceInfo)
	r, err := redisConn.Do("HGET", LOCATIONTOPLACEHMAP, field)
	placeId, err = redis.Int64(r, err)
	if err != nil {
		infoLog.Printf("Reids exec HGET field=%s err=%v", field, err)
		// load from db
		placeId = MatchLocationPlaceIdInDb(locationInfo, &redisConn)
		return
	} else {
		infoLog.Println("Redis exec HGet succ placeId =", placeId)
		return
	}
}

func MatchLocationPlaceIdInDb(locationInfo *tagLocationInfo, conn *redis.Conn) (placeId int64) {
	if (locationInfo == nil) || (conn == nil) {
		infoLog.Println("MatchLocationPlaceIdInDb input param err")
		placeId = 0 // placeId must not be 0
		return
	}
	err := db.QueryRow("select PLACEID from HT_LOCATION_PLACE WHERE COUNTRY=? and ADMINISTRATIVE1=? and ADMINISTRATIVE2=? and "+
		"ADMINISTRATIVE3=? and LOCALITY=? and SUBLOCALITY=? and NEIGHBORHOOD=?",
		locationInfo.PlaceInfo.Country,
		locationInfo.PlaceInfo.Admin1,
		locationInfo.PlaceInfo.Admin2,
		locationInfo.PlaceInfo.Admin3,
		locationInfo.PlaceInfo.Locality,
		locationInfo.PlaceInfo.SubLocality,
		locationInfo.PlaceInfo.Neighborhood).Scan(&placeId)
	switch {
	case err == sql.ErrNoRows:
		infoLog.Println("MatchLocationPlaceInfoInDb mysql Query empty exec insert")
		// 执行插入语句 使用mysql 的函数now()获取当前时间戳
		r, err := db.Exec("insert into HT_LOCATION_PLACE (STATE,COUNTRY,ADMINISTRATIVE1,ADMINISTRATIVE2,ADMINISTRATIVE3,"+
			"LOCALITY,SUBLOCALITY,NEIGHBORHOOD,LATITUDE,LONGITUDE,UPDATETIME) value"+
			"(0, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())",
			locationInfo.PlaceInfo.Country,
			locationInfo.PlaceInfo.Admin1,
			locationInfo.PlaceInfo.Admin2,
			locationInfo.PlaceInfo.Admin3,
			locationInfo.PlaceInfo.Locality,
			locationInfo.PlaceInfo.SubLocality,
			locationInfo.PlaceInfo.Neighborhood,
			locationInfo.Latitude,
			locationInfo.Longtitude)
		if err != nil {
			infoLog.Println("MatchLocationPlaceInfoInDb insert faield err =", err)
			placeId = 0
			return
		} else {
			placeId, err = r.LastInsertId()
			if err != nil {
				infoLog.Println("MatchLocationPlaceInfoInDb Get last insert id faied err =", err)
				placeId = 0
				return
			}
			infoLog.Println("MatchLocationPlaceInfoInDb succ last id =", placeId)
			// 查询db 得到id 后重新将id 插入redis
			field := GetLocationPlaceIdKey(&locationInfo.PlaceInfo)
			_, err := (*conn).Do("HSET", LOCATIONTOPLACEHMAP, field, placeId)
			if err != nil {
				infoLog.Printf("HSET location#palceid field=%s place=%v faield err=%v", field, placeId, err)
			} else {
				infoLog.Printf("HSET locatio#palceid field=%s place=%v succ", field, placeId)
			}
		}
	case err != nil:
		infoLog.Println("MatchLocationPlaceInfoInDb rows.Scan failed err =", err)
		placeId = 0
		return
	default: // 查询成功将查结果存入redis并返回插入结果
		infoLog.Printf("placeId = %v\n", placeId)
		// 查询db 得到id 后重新将id 插入redis
		field := GetLocationPlaceIdKey(&locationInfo.PlaceInfo)
		_, err := (*conn).Do("HSET", LOCATIONTOPLACEHMAP, field, placeId)
		if err != nil {
			infoLog.Printf("HSET location#palceid field=%s place=%v faield err=%v", field, placeId, err)
		} else {
			infoLog.Printf("HSET location#palceid field=%s place=%v succ", field, placeId)
		}
		return
	}
	return
}

func UpdateLocationInDb(locationInfo *tagLocationInfo) {
	if locationInfo == nil {
		infoLog.Println("UpdateLocationInDb input nil")
		return
	}
	r, err := db.Exec("INSERT INTO HT_USER_LOCATION2 (USERID, ALLOWED, LATITUDE, LONGITUDE,"+
		"PLACEID, COUNTRY, ADMINISTRATIVE1, ADMINISTRATIVE2, ADMINISTRATIVE3, LOCALITY, "+
		"SUBLOCALITY, NEIGHBORHOOD, UPDATETIME) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())"+
		"ON DUPLICATE KEY UPDATE ALLOWED=VALUES(ALLOWED),LATITUDE=VALUES(LATITUDE),LONGITUDE=VALUES(LONGITUDE),"+
		"PLACEID=VALUES(PLACEID),COUNTRY=VALUES(COUNTRY),ADMINISTRATIVE1=VALUES(ADMINISTRATIVE1),ADMINISTRATIVE2=VALUES(ADMINISTRATIVE2),"+
		"ADMINISTRATIVE3=VALUES(ADMINISTRATIVE3),LOCALITY=VALUES(LOCALITY), SUBLOCALITY=VALUES(SUBLOCALITY),"+
		"NEIGHBORHOOD=VALUES(NEIGHBORHOOD), UPDATETIME=VALUES(UPDATETIME)",
		locationInfo.Uid,
		locationInfo.Allowed,
		locationInfo.Latitude,
		locationInfo.Longtitude,
		locationInfo.PlaceId,
		locationInfo.PlaceInfo.Country,
		locationInfo.PlaceInfo.Admin1,
		locationInfo.PlaceInfo.Admin2,
		locationInfo.PlaceInfo.Admin3,
		locationInfo.PlaceInfo.Locality,
		locationInfo.PlaceInfo.SubLocality,
		locationInfo.PlaceInfo.Neighborhood)
	if err != nil {
		infoLog.Println("UpdateLocationInDb insert faield err =", err)
		return
	} else {
		affectRow, err := r.RowsAffected()
		if err != nil {
			infoLog.Println("UpdateLocationInDb affectRow faied err =", err)
		}
		infoLog.Println("affectRow =", affectRow)
		// 更新HT_USER_LOCATION2 成功 继续更新HT_USER_PROPERTY
		curVersion := time.Now().Unix() // 得到秒
		_, err = db.Exec("update HT_USER_PROPERTY set LOCATIONVERSION = ? where USERID = ?", curVersion, locationInfo.Uid)
		if err != nil {
			infoLog.Println("update HT_USER_PROPERTY failed err =", err)
		}
		return
	}
}

func MatchLanguageShortName(languageType string) (shortName string) {
	shortName, ok := languageShortName[languageType]
	if !ok { // 不存在对应的short name 直接返回 "en"
		shortName = "en"
		return
	}
	return
}

func GetMultiCountryName(countryName string, languageType string) (name string) {
	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	lang := languageType
	_, ok := languageShortName[lang]
	if !ok {
		lang = "English"
	}
	field := GetMultiCountryNameKey(countryName, lang)
	r, err := redisConn.Do("HGET", MULTILANGNAMEHMAP, field)
	name, err = redis.String(r, err)
	if err != nil {
		infoLog.Printf("Reids exec HGET field=%s err=%v\n", field, err)
		// load from db
		name = GetMultiCountryNameInDb(countryName, languageType, &redisConn)
	} else {
		infoLog.Println("Redis exec HGet succ [filed name] =", field, name)
	}
	return
}

func GetMultiCountryNameInDb(countryName string, languageType string, conn *redis.Conn) (name string) {
	if len(countryName) == 0 || len(languageType) == 0 || conn == nil {
		infoLog.Printf("GetMultiCountryNameInDb input err countryName=%s languageType=%s", countryName, languageType)
		return
	}
	formatString := fmt.Sprintf("select %s from HT_MULTI_COUNTRY WHERE CountryName = ?", languageType)
	err := db.QueryRow(formatString, countryName).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		infoLog.Printf("Not exist countryName=%s languageType=%s name", countryName, languageType)
	case err != nil:
		infoLog.Printf("get countryName=%s languageType=%s err=%v", countryName, languageType, err)
	default:
		infoLog.Printf("countryName=%s languageType=%s name=%s", countryName, languageType, name)
		// 查询db 得到id 后重新将id 插入redis
		field := GetMultiCountryNameKey(countryName, languageType)
		_, err := (*conn).Do("HSET", MULTILANGNAMEHMAP, field, name)
		if err != nil {
			infoLog.Printf("HSET country#lang field=%s name=%s faield err=%v", field, name, err)
		} else {
			infoLog.Printf("HSET country#lang field=%s name=%s succ", field, name)
		}
	}
	return
}

func ReformLocationCityName(placeInfo *tagPlaceInformation) (cityName string) {
	cityName = string(placeInfo.Locality)
	if len(cityName) > 0 {
		return
	}
	cityName = string(placeInfo.SubLocality)
	if len(cityName) > 0 {
		return
	}
	cityName = string(placeInfo.Neighborhood)
	if len(cityName) > 0 {
		return
	}
	cityName = string(placeInfo.Admin3)
	if len(cityName) > 0 {
		return
	}
	cityName = string(placeInfo.Admin2)
	if len(cityName) > 0 {
		return
	}
	cityName = string(placeInfo.Admin1)
	if len(cityName) > 0 {
		return
	}
	return
}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	infoLog.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet, ok := p.(*common.XTHeadPacket)
	if !ok { // 不是XTHeadPacket报文
		infoLog.Println("packet can not change to xtpacket")
		return false
	}

	head, err := packet.GetHead()
	if err != nil {
		//SendResp(c, head, uint8(ERR_INVALID_PARAM))
		infoLog.Println("Get head failed", err)
		return false
	}
	attr := "golbs/recv_req_count"
	libcomm.AttrAdd(attr, 1)

	//infoLog.Printf("OnMessage:[%#v] len=%v payLoad=%v\n", head, len(packet.GetBody()), packet.GetBody())
	infoLog.Printf("OnMessage:[%#v] len=%v\n", head, len(packet.GetBody()))
	_, err = packet.CheckXTPacketValid()
	if err != nil {
		SendResp(c, head, uint8(ERR_INVALID_PARAM))
		infoLog.Println("Invalid packet", err)
		return false
	}

	switch head.Cmd {
	case CMD_UPDATE_LOCATION:
		go ProcUpdateLocation(c, p)
	case CMD_GET_USER_LOCATION:
		go ProcGetUserLocation(c, p)
	case CMD_REFRESH_LOCATION:
		go ProcRefreshLocation(c, p)
	default:
		infoLog.Println("UnHandle Cmd =", head.Cmd)
	}
	return true
}

func SendResp(c *gotcp.Conn, reqHead *common.XTHead, ret uint8) bool {
	head := new(common.XTHead)
	if reqHead != nil {
		*head = *reqHead
	}

	head.Cmd = reqHead.Cmd + 1 // ack cmd = req cmd + 1
	head.Len = 1               // sizeof(uint8)
	buf := make([]byte, common.XTHeadLen+head.Len)
	err := common.SerialXTHeadToSlice(head, buf[:])
	if err != nil {
		infoLog.Println("SerialXTHeadToSlice failed")
		return false
	}
	buf[common.XTHeadLen] = ret // return code
	resp := common.NewXTHeadPacket(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func SendRespWithPayLoad(c *gotcp.Conn, reqHead *common.XTHead, payLoad []byte) bool {
	head := new(common.XTHead)
	if reqHead != nil {
		*head = *reqHead
	}

	head.Cmd = reqHead.Cmd + 1      // ack cmd = req cmd + 1
	head.Len = uint32(len(payLoad)) //
	buf := make([]byte, common.XTHeadLen+head.Len)
	err := common.SerialXTHeadToSlice(head, buf[:])
	if err != nil {
		infoLog.Println("SerialXTHeadToSlice failed")
		return false
	}
	copy(buf[common.XTHeadLen:], payLoad) // return code
	resp := common.NewXTHeadPacket(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func ProcUpdateLocation(c *gotcp.Conn, p gotcp.Packet) bool {
	packet, ok := p.(*common.XTHeadPacket)
	if !ok { // 不是XTHeadPacket报文
		infoLog.Println("packet can not change to xtpacket")
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Printf("ProcUpdateLocation get packet failed\n")
		return false
	}

	attr := "golbs/update_location_req_count"
	libcomm.AttrAdd(attr, 1)

	body := packet.GetBody()
	var locationInfo tagLocationInfo
	locationInfo.Uid = common.UnMarshalUint32(&body)
	languageType := common.UnMarshalSlice(&body)
	locationInfo.Allowed = common.UnMarshalUint8(&body)
	if locationInfo.Allowed > 1 || locationInfo.Uid != head.From {
		SendResp(c, head, uint8(ERR_INVALID_PARAM))
		infoLog.Printf("ProcUpdateLocation error allowed =%v uid=%v from=%v", locationInfo.Allowed, locationInfo.Uid, head.From)
		return false
	}

	var needUpdate int
	if locationInfo.Allowed == 1 {
		locationInfo.Latitude = common.UnMarshalSlice(&body)
		locationInfo.Longtitude = common.UnMarshalSlice(&body)
		locationInfo.PlaceInfo.Country = common.UnMarshalSlice(&body)
		locationInfo.PlaceInfo.Admin1 = common.UnMarshalSlice(&body)
		locationInfo.PlaceInfo.Admin2 = common.UnMarshalSlice(&body)
		locationInfo.PlaceInfo.Admin3 = common.UnMarshalSlice(&body)
		locationInfo.PlaceInfo.Locality = common.UnMarshalSlice(&body)
		locationInfo.PlaceInfo.SubLocality = common.UnMarshalSlice(&body)
		locationInfo.PlaceInfo.Neighborhood = common.UnMarshalSlice(&body)
		locationInfo.Source = common.UnMarshalSlice(&body)
		needUpdate = strings.Compare(STRGOOGLE, string(locationInfo.Source))
	}
	infoLog.Printf("uid=%v lang=%s allow=%d lat=%s lng=%s country=%s"+
		" admin1=%s admin2=%s admin3=%s locality=%s subloc=%s neibhood=%s source=%s needUpdate=%v",
		locationInfo.Uid,
		languageType,
		locationInfo.Allowed,
		locationInfo.Latitude,
		locationInfo.Longtitude,
		locationInfo.PlaceInfo.Country,
		locationInfo.PlaceInfo.Admin1,
		locationInfo.PlaceInfo.Admin2,
		locationInfo.PlaceInfo.Admin3,
		locationInfo.PlaceInfo.Locality,
		locationInfo.PlaceInfo.SubLocality,
		locationInfo.PlaceInfo.Neighborhood,
		locationInfo.Source,
		needUpdate)
	// 检查参数是否合法
	if !IsValidParam(head, &locationInfo) {
		SendResp(c, head, uint8(ERR_INVALID_PARAM))
		infoLog.Printf("error input allowed =%v uid=%v from=%v", locationInfo.Allowed, locationInfo.Uid, head.From)
		return false
	}

	var placeId int64
	if (locationInfo.Allowed == 1) &&
		(needUpdate == 0) && // update only if google
		(len(locationInfo.PlaceInfo.Country) != 0) { // update if Country is not empty
		placeId = MatchLocationPlaceId(&locationInfo)
		// 从redis或数据库加载placeId 仍然加载失败直接返回内部错误
		if placeId == 0 { // load failed
			SendResp(c, head, uint8(ERR_INTERNAL_ERROR))
			infoLog.Println("ProcUpdateLocation load placeId failed [from, allow] =", head.From, locationInfo.Allowed)
			return false
		}

	}

	locationInfo.PlaceId = uint32(placeId)
	// 更新HT_USER_LOCATION2 和 HT_USER_PROPERTY表
	UpdateLocationInDb(&locationInfo)
	// 删除Redis中的位置信息
	DelSingleUserLocationInRedis(locationInfo.Uid)

	// 获取多语言的位置
	multiLocation := locationInfo
	if locationInfo.Allowed == 1 {
		if strings.Compare(string(languageType), string("English")) != 0 {
			multiLocation.PlaceInfo = GetMultiPlaceInfo(&(locationInfo.PlaceInfo), string(languageType))
			// 如果查询db 多语言位置也失败 这使用English地理位置
			if len(multiLocation.PlaceInfo.Country) == 0 {
				multiLocation.PlaceInfo = locationInfo.PlaceInfo
			}

		}
		multiCountryName := GetMultiCountryName(string(locationInfo.PlaceInfo.Country), string(languageType))
		if len(multiCountryName) > 0 {
			multiLocation.PlaceInfo.Country = []byte(multiCountryName)
			infoLog.Println("multi country =", multiCountryName)
		}
	}

	var respPayLoad []byte
	common.MarshalUint8(uint8(0), &respPayLoad)
	common.MarshalUint8(locationInfo.Allowed, &respPayLoad)
	city := ReformLocationCityName(&(multiLocation.PlaceInfo))
	if locationInfo.Allowed == 1 {
		common.MarshalSlice(multiLocation.PlaceInfo.Country, &respPayLoad)
		common.MarshalSlice([]byte(city), &respPayLoad)
	}
	SendRespWithPayLoad(c, head, respPayLoad)
	infoLog.Printf("from=%v to=%v cmd=%v seq=%v ret=%d allow=%d country=%s city=%s\n",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		0,
		locationInfo.Allowed,
		multiLocation.PlaceInfo.Country,
		city)
	return true
}

func GetMultiPlaceInfo(englishPlace *tagPlaceInformation, languageType string) (multiPlace tagPlaceInformation) {
	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	shortLang := MatchLanguageShortName(languageType)
	placeId := MatchLocationPlaceIdInRedis(englishPlace) //之前更新地址位置时会立刻更新redis 所以子查询redis即可
	if placeId > 0 {
		field := GetMultiPlaceInfoKey(placeId, shortLang)
		r, err := redisConn.Do("HGET", MULTIPLACEINFOHMAP, field)
		strPlace, err := redis.String(r, err)
		if err != nil {
			infoLog.Printf("Reids exec HGET field=%s err=%v", field, err)
			// load from db
			var ret bool
			multiPlace, ret = GetMultiPlaceInfoInDb(placeId, shortLang, &redisConn)
			if !ret {
				infoLog.Println("Get MultiPlaceInfoInDb failed [placeId, shortLang]", placeId, shortLang)
				return
			}

		} else {
			infoLog.Println("Redis exec HGet succ strPalce=", strPlace)
			fieldArry := strings.Split(strPlace, "#")
			infoLog.Println("field size =", len(fieldArry))
			for i, v := range fieldArry {
				// infoLog.Printf("index=%v value=%s", i, v)
				switch i {
				case 0:
					multiPlace.Country = []byte(v)
					// infoLog.Println("Country =", multiPlace.Country)
				case 1:
					multiPlace.Admin1 = []byte(v)
					// infoLog.Println("Admin1 =", multiPlace.Admin1)
				case 2:
					multiPlace.Admin2 = []byte(v)
					// infoLog.Println("Admin2 =", multiPlace.Admin2)
				case 3:
					multiPlace.Admin3 = []byte(v)
					// infoLog.Println("Admin3 =", multiPlace.Admin3)
				case 4:
					multiPlace.Locality = []byte(v)
					// infoLog.Println("Locality =", multiPlace.Locality)
				case 5:
					multiPlace.SubLocality = []byte(v)
					// infoLog.Println("SubLocality =", multiPlace.SubLocality)
				case 6:
					multiPlace.Neighborhood = []byte(v)
					// infoLog.Println("Neighborhood =", multiPlace.Neighborhood)
				default:
					infoLog.Printf("unknow filed index=%v value=%v", i, v)
				}
			}
			return
		}

	} else {
		infoLog.Println("MatchLocationPlaceIdInRedis failed")
	}
	return
}

func GetMultiPlaceInfoInDb(placeId int64, shortLang string, conn *redis.Conn) (multiPlace tagPlaceInformation, ret bool) {
	if conn == nil {
		infoLog.Println("GetMultiPlaceInfoInDb input param err")
		ret = false
		return
	}
	var country, admin1, admin2, admin3, loca, subLoca, neigh string
	err := db.QueryRow("select COUNTRY,ADMINISTRATIVE1,ADMINISTRATIVE2,ADMINISTRATIVE3,LOCALITY,SUBLOCALITY,NEIGHBORHOOD from HT_MULTILANG_PLACE where PLACEID=? and LANGTYPE=?", placeId, shortLang).Scan(&country, &admin1, &admin2, &admin3, &loca, &subLoca, &neigh)
	switch {
	case err == sql.ErrNoRows:
		infoLog.Println("Db mysql Query empty [palce shortLang] =", placeId, shortLang)
		ret = false
		return
	case err != nil:
		infoLog.Println("GetMultiPlaceInfoInDb rows.Scan failed [place shortLang err] =", placeId, shortLang, err)
		ret = false
		return
	default: // 查询成功将查结果存入redis并返回插入结果
		infoLog.Printf("palceId=%v shortLang=%s country=%s admin1=%s admin2=%s admin3=%s loca=%s subLoac=%s neigh=%s \n", placeId, shortLang, country, admin1, admin2, admin3, loca, subLoca, neigh)
		// 设置返回值
		multiPlace.Country = []byte(country)
		multiPlace.Admin1 = []byte(admin1)
		multiPlace.Admin2 = []byte(admin2)
		multiPlace.Admin3 = []byte(admin3)
		multiPlace.Locality = []byte(loca)
		multiPlace.SubLocality = []byte(subLoca)
		multiPlace.Neighborhood = []byte(neigh)
		// 查询db 得到MultiPlaceInfo 后重新将插入redis
		field := GetMultiPlaceInfoKey(placeId, shortLang)
		value := country + "#" + admin1 + "#" + admin2 + "#" + admin3 + "#" + loca + "#" + subLoca + "#" + neigh
		_, err := (*conn).Do("HSET", MULTIPLACEINFOHMAP, field, value)
		if err != nil {
			infoLog.Printf("HSET placeid#lang field=%s value=%v faield err=%v", field, value, err)
		} else {
			infoLog.Printf("HSET palceid#lang field=%s value=%v succ", field, value)
		}
		return
	}
	return
}

func IsValidParam(head *common.XTHead, locationInfo *tagLocationInfo) bool {
	if head == nil || locationInfo == nil {
		return false
	}
	if head.From != locationInfo.Uid {
		return false
	}
	if locationInfo.Allowed != 1 {
		return true
	}
	if (len(locationInfo.Latitude) == 0) ||
		(len(locationInfo.Longtitude) == 0) ||
		(strings.Compare(string("0"), string(locationInfo.Latitude)) == 0) ||
		(strings.Compare(string("1"), string(locationInfo.Latitude)) == 0) ||
		(strings.Compare(string("0.0"), string(locationInfo.Latitude)) == 0) ||
		(strings.Compare(string("1.1"), string(locationInfo.Latitude)) == 0) ||
		(strings.Compare(string("0"), string(locationInfo.Longtitude)) == 0) ||
		(strings.Compare(string("1"), string(locationInfo.Longtitude)) == 0) ||
		(strings.Compare(string("0.0"), string(locationInfo.Longtitude)) == 0) ||
		(strings.Compare(string("1.1"), string(locationInfo.Longtitude)) == 0) ||
		(strings.Compare(string(locationInfo.Latitude), string(locationInfo.Longtitude)) == 0) {
		return false
	}
	return true
}

func ProcGetUserLocation(c *gotcp.Conn, p gotcp.Packet) bool {
	packet, ok := p.(*common.XTHeadPacket)
	if !ok { // 不是XTHeadPacket报文
		infoLog.Println("packet can not change to xtpacket")
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Printf("ProcGetUserLocation get head failed\n")
		return false
	}

	attr := "golbs/get_location_req_count"
	libcomm.AttrAdd(attr, 1)

	body := packet.GetBody()
	languageType := common.UnMarshalSlice(&body)
	userCount := common.UnMarshalUint16(&body)
	infoLog.Printf("Recv from=%v to=%v cmd=%v seq=%v count=%v lang=%s\n",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		userCount,
		languageType)

	if userCount == 0 || userCount > 200 {
		SendResp(c, head, uint8(ERR_INTERNAL_ERROR))
		infoLog.Println("input err count =", userCount)
		return false
	}

	listUserId := make([]uint32, userCount)
	for i := 0; i < int(userCount); i++ {
		uid := common.UnMarshalUint32(&body)
		if uid != 0 {
			listUserId[i] = uid
		} else {
			infoLog.Printf("input err uid = 0 from=%v\n", head.From)
		}
	}
	if len(listUserId) == 0 {
		infoLog.Printf("uid list empty from=%v\n", head.From)
		SendResp(c, head, uint8(RET_SUCCESS))
		return false
	}

	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	var listLocationInfo []tagLocationInfo
	var missUidList []uint32
	strLocation := GetUserLocationFromRedis(listUserId, &redisConn)
	if len(strLocation) == 0 {
		infoLog.Println("GetUserLocationFromRedis failed empty result load from db")
		missUidList = listUserId // 一个都没有找到直接从数据库中查询全部
	} else { // 获取到部分位置信息
		for i, v := range strLocation {
			if len(v) == 0 { // 长度为0 说明没有找到此用户地址位置信息添加到missUidList中
				missUidList = append(missUidList, listUserId[i])
			} else {
				// 解析位置 地址位置使用'#'分隔
				fieldArry := strings.Split(v, "#")
				infoLog.Println("place size =", len(fieldArry))
				if len(fieldArry) != 12 { // 一共12个字段
					infoLog.Println("GetUserLocationFromRedis format err [i, v] =", i, v)
					missUidList = append(missUidList, listUserId[i])
					continue
				}
				var locational tagLocationInfo
				if uid, err := strconv.Atoi(fieldArry[0]); err != nil {
					infoLog.Println("GetUserLocationFromRedis Conv uid err =", err)
					missUidList = append(missUidList, listUserId[i])
					continue
				} else {
					locational.Uid = uint32(uid)
				}

				if allow, err := strconv.Atoi(fieldArry[1]); err != nil {
					infoLog.Println("GetUserLocationFromRedis Conv Allowed err =", err)
					missUidList = append(missUidList, listUserId[i])
					continue
				} else {
					locational.Allowed = uint8(allow)
				}
				if locational.Allowed == 1 { // 允许定位
					for i, v := range fieldArry {
						// infoLog.Printf("index=%v value=%s \n", i, v)
						switch i { // 前面已经处理了index=0、1的元素
						case 0: // does not need the value
						case 1: // does not need the value
						case 2:
							locational.Latitude = []byte(v)
						case 3:
							locational.Longtitude = []byte(v)
						case 4:
							locational.PlaceInfo.Country = []byte(v)
						case 5:
							locational.PlaceInfo.Admin1 = []byte(v)
						case 6:
							locational.PlaceInfo.Admin2 = []byte(v)
						case 7:
							locational.PlaceInfo.Admin3 = []byte(v)
						case 8:
							locational.PlaceInfo.Locality = []byte(v)
						case 9:
							locational.PlaceInfo.SubLocality = []byte(v)
						case 10:
							locational.PlaceInfo.Neighborhood = []byte(v)
						default:
							infoLog.Printf("unknow filed index=%v value=%v", i, v)
						}
					}
				}
				if ts, err := strconv.Atoi(fieldArry[11]); err != nil {
					infoLog.Println("GetUserLocationFromRedis Conv TimeStamp err =", err)
					missUidList = append(missUidList, listUserId[i])
					continue
				} else {
					locational.TimeStamp = uint64(ts)
				}
				listLocationInfo = append(listLocationInfo, locational) // 得到一个用户的地址位置添加到结果集中
			}
		}
	}

	if len(missUidList) != 0 { // 有用户cache miss 从DB中load
		ret := GetUserLocationFromDb(missUidList, &listLocationInfo, &redisConn)
		if !ret {
			infoLog.Println("GetUserLocationFromDb failed")
		}
	}
	for i, v := range listLocationInfo {
		if v.Allowed == 1 {
			if strings.Compare(string(languageType), "English") != 0 { // 获取多语言的位置
				multiPlace := GetMultiPlaceInfo(&(v.PlaceInfo), string(languageType))
				if len(multiPlace.Country) != 0 { //存在多语言的位置才替换
					listLocationInfo[i].PlaceInfo = multiPlace
				}
			}
			multiCountryName := GetMultiCountryName(string(listLocationInfo[i].PlaceInfo.Country), string(languageType))
			if len(multiCountryName) != 0 {
				listLocationInfo[i].PlaceInfo.Country = []byte(multiCountryName)
			}
		}
	}
	// 返回响应
	var respPayLoad []byte
	common.MarshalUint8(0, &respPayLoad) // rsp code
	common.MarshalUint16(uint16(len(listLocationInfo)), &respPayLoad)
	for _, v := range listLocationInfo {
		var subRspPayLoad []byte
		common.MarshalUint32(v.Uid, &subRspPayLoad)
		common.MarshalUint64(v.TimeStamp, &subRspPayLoad)
		common.MarshalUint8(v.Allowed, &subRspPayLoad)
		if v.Allowed == 1 {
			common.MarshalSlice(v.Latitude, &subRspPayLoad)
			common.MarshalSlice(v.Longtitude, &subRspPayLoad)
			common.MarshalSlice(v.PlaceInfo.Country, &subRspPayLoad)
			cityName := ReformLocationCityName(&(v.PlaceInfo))
			common.MarshalSlice([]byte(cityName), &subRspPayLoad)
		}
		common.MarshalSlice(subRspPayLoad, &respPayLoad)
	}
	SendRespWithPayLoad(c, head, respPayLoad)
	infoLog.Printf("ret=%d count=%v\n", 0, len(listLocationInfo))
	return true
}

func GetUserLocationFromRedis(listUserId []uint32, redisConn *redis.Conn) (strLocation []string) {
	if len(listUserId) == 0 {
		infoLog.Println("listUserId empty")
		return
	}

	strUidList := make([]string, len(listUserId))
	for i, _ := range listUserId {
		strUidList[i] = GetUserPlaceKey(listUserId[i])
	}
	infoLog.Println("[len strUidList] =", len(strUidList), strUidList)
	s := make([]interface{}, len(strUidList))
	for i, v := range strUidList {
		s[i] = v
		infoLog.Printf("index=%v key=%s", i, v)
	}
	values, err := redis.Values((*redisConn).Do("MGET", s...)) // get place
	if err != nil {
		infoLog.Println("Reids failed MGet err =", err)
		return
	}

	if err := redis.ScanSlice(values, &strLocation); err != nil {
		infoLog.Println("ScanSlice failed err=", err)
		return
	}
	infoLog.Println("result size= ", len(strLocation))
	for i, v := range strLocation {
		infoLog.Printf("index=%v result=%v\n", i, v)
	}
	return
}

func GetUserLocationFromDb(listUserId []uint32, listLocationInfo *[]tagLocationInfo, conn *redis.Conn) bool {
	if len(listUserId) == 0 {
		return true
	}
	var strUids string
	for i, v := range listUserId {
		if i == 0 {
			strUids += strconv.Itoa(int(v))
		} else {
			strUids += "," + strconv.Itoa(int(v))
		}
	}
	infoLog.Println("GetUserLocationFromDb uids =", strUids)
	rows, err := db.Query("select t1.USERID, t1.ALLOWED, t1.LATITUDE, t1.LONGITUDE, t1.COUNTRY," +
		"t1.ADMINISTRATIVE1, t1.ADMINISTRATIVE2, t1.ADMINISTRATIVE3, t1.LOCALITY," +
		"t1.SUBLOCALITY, t1.NEIGHBORHOOD, t2.LOCATIONVERSION from HT_USER_LOCATION2 " +
		"as t1 left join (HT_USER_PROPERTY as t2) on (t2.USERID = t1.USERID)" +
		" where t1.USERID in (" + strUids + ");")
	if err != nil {
		infoLog.Printf("GetUserLocationFromDb uids=%s failed\n", strUids)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var strUid, strAllow, strLa, strLong, strCo string
		var strA1, strA2, strA3, strLoc, strSubLoc, strNe, strVer string
		if err := rows.Scan(&strUid, &strAllow, &strLa, &strLong, &strCo, &strA1, &strA2, &strA3, &strLoc, &strSubLoc, &strNe, &strVer); err != nil {
			infoLog.Println("GetUserLocationFromDb rows.Scan failed")
			continue
		}
		// set redis
		value := strUid + "#" + strAllow + "#" + strLa + "#" + strLong + "#" + strCo + "#"
		value += strA1 + "#" + strA2 + "#" + strA3 + "#" + strLoc + "#" + strSubLoc + "#"
		value += strNe + "#" + strVer

		uid, err := strconv.Atoi(strUid)
		if err != nil {
			infoLog.Println("GetUserLocationFromDb strconv failed uid =", strUid)
			continue
		}
		strKey := GetUserPlaceKey(uint32(uid))
		_, err = (*conn).Do("SETEX", strKey, kTTLInterval, value)
		if err != nil {
			infoLog.Printf("SETEX key=%s value=%s faield err=%v", strKey, value, err)
		} else {
			infoLog.Printf("SETEX key=%s value=%s succ err=%v", strKey, value, err)
		}

		// add to listLocationInfo
		allow, err := strconv.Atoi(strAllow)
		if err != nil {
			infoLog.Println("GetUserLocationFromDb strconv failed allow =", strAllow)
			continue
		}
		ts, err := strconv.Atoi(strVer)
		if err != nil {
			infoLog.Println("GetUserLocationFromDb strconv failed ver =", strVer)
			continue
		}
		loca := tagLocationInfo{Uid: uint32(uid),
			Allowed:    uint8(allow),
			Latitude:   []byte(strLa),
			Longtitude: []byte(strLong),
			PlaceInfo: tagPlaceInformation{Admin1: []byte(strA1),
				Admin2:       []byte(strA2),
				Admin3:       []byte(strA3),
				Country:      []byte(strCo),
				Locality:     []byte(strLoc),
				SubLocality:  []byte(strSubLoc),
				Neighborhood: []byte(strNe),
			},
			TimeStamp: uint64(ts),
		}
		*listLocationInfo = append(*listLocationInfo, loca)
	}
	return true
}

func ProcRefreshLocation(c *gotcp.Conn, p gotcp.Packet) bool {
	packet, ok := p.(*common.XTHeadPacket)
	if !ok { // 不是XTHeadPacket报文
		infoLog.Println("packet can not change to xtpacket")
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Printf("ProcRefreshLocation get head failed\n")
		return false
	}

	attr := "golbs/refersh_location_req_count"
	libcomm.AttrAdd(attr, 1)

	body := packet.GetBody()
	languageType := common.UnMarshalSlice(&body)
	uid := common.UnMarshalUint32(&body)
	ts := common.UnMarshalUint64(&body)
	infoLog.Printf("ProcRefreshLocation from=%v to=%v cmd=%v seq=%v uid=%v ts=%v\n",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		uid,
		ts)
	if head.From != uid {
		infoLog.Println("[from uid] not equal", head.From, uid)
	}

	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	var locationInfo tagLocationInfo
	strLocation := GetSingleUserLocationFromRedis(uid, &redisConn)
	var isNeedLoadFromDb bool = true
	if len(strLocation) != 0 {
		// 解析位置 地址位置使用'#'分隔
		fieldArry := strings.Split(strLocation, "#")
		infoLog.Println("place size =", len(fieldArry))
		if len(fieldArry) != 12 { // 一共12个字段
			infoLog.Println("GetSingleUserLocationFromRedis format err")
			goto LOADFROMDB
		}
		locationInfo.Uid = uid
		if allow, err := strconv.Atoi(fieldArry[1]); err != nil {
			infoLog.Println("GetUserLocationFromRedis Conv Allowed err =", err)
			goto LOADFROMDB
		} else {
			locationInfo.Allowed = uint8(allow)
		}
		for i, v := range fieldArry {
			// infoLog.Printf("index=%v value=%s", i, v)
			switch i { // 前面已经处理了index=0、1的元素
			case 2:
				locationInfo.Latitude = []byte(v)
			case 3:
				locationInfo.Longtitude = []byte(v)
			case 4:
				locationInfo.PlaceInfo.Country = []byte(v)
			case 5:
				locationInfo.PlaceInfo.Admin1 = []byte(v)
			case 6:
				locationInfo.PlaceInfo.Admin2 = []byte(v)
			case 7:
				locationInfo.PlaceInfo.Admin3 = []byte(v)
			case 8:
				locationInfo.PlaceInfo.Locality = []byte(v)
			case 9:
				locationInfo.PlaceInfo.SubLocality = []byte(v)
			case 10:
				locationInfo.PlaceInfo.Neighborhood = []byte(v)
			default:
				infoLog.Printf("unknow filed index=%v value=%v", i, v)
			}
		}
		if ts, err := strconv.Atoi(fieldArry[11]); err != nil {
			infoLog.Println("GetUserLocationFromRedis Conv TimeStamp err =", err)
			goto LOADFROMDB
		} else {
			locationInfo.TimeStamp = uint64(ts)
		}
		isNeedLoadFromDb = false // get from redis succ no need load from db
	}
LOADFROMDB:
	if isNeedLoadFromDb { // load from db
		ret := GetSingleUserLocationFromDb(uid, ts, &locationInfo, &redisConn)
		switch ret {
		case DB_RET_NOT_EXIST:
			infoLog.Printf("from=%v to=%v cmd=%v seq=%v ts=%v no need to update\n",
				head.From,
				head.To,
				head.Cmd,
				head.Seq,
				ts)
			respPayLoad := make([]byte, 1)
			common.MarshalUint8(uint8(1), &respPayLoad)
			msg := []byte("not need to update")
			common.MarshalSlice(msg, &respPayLoad)
			SendRespWithPayLoad(c, head, respPayLoad)
			return true
		case DB_RET_EXEC_FAILED:
			infoLog.Printf("from=%v to=%v cmd=%v seq=%v ts=%v db exec failed\n",
				head.From,
				head.To,
				head.Cmd,
				head.Seq,
				ts)
			SendResp(c, head, uint8(ERR_INTERNAL_ERROR))
			return true
		default: // get succ
			infoLog.Println("GetSingleUserLocationFromDb succ uid", uid)
		}
	}
	multiLocation := locationInfo
	if multiLocation.Allowed == 1 {
		if strings.Compare(string(languageType), "English") != 0 {
			multiLocation.PlaceInfo = GetMultiPlaceInfo(&(locationInfo.PlaceInfo), string(languageType))
			// 多语言地址位置为空继续使用English位置
			if len(multiLocation.PlaceInfo.Country) == 0 {
				multiLocation.PlaceInfo = locationInfo.PlaceInfo
			}
		}
		multiCountryName := GetMultiCountryName(string(locationInfo.PlaceInfo.Country), string(languageType))
		if len(multiCountryName) > 0 {
			multiLocation.PlaceInfo.Country = []byte(multiCountryName)
		}
	}
	// 返回响应
	var respPayLoad []byte
	common.MarshalUint8(uint8(0), &respPayLoad)
	var subRspPayLoad []byte
	common.MarshalUint32(multiLocation.Uid, &subRspPayLoad)
	common.MarshalUint64(multiLocation.TimeStamp, &subRspPayLoad)
	common.MarshalUint8(multiLocation.Allowed, &subRspPayLoad)
	if multiLocation.Allowed == 1 {
		common.MarshalSlice(multiLocation.Latitude, &subRspPayLoad)
		common.MarshalSlice(multiLocation.Longtitude, &subRspPayLoad)
		common.MarshalSlice(multiLocation.PlaceInfo.Country, &subRspPayLoad)
		cityName := ReformLocationCityName(&(multiLocation.PlaceInfo))
		common.MarshalSlice([]byte(cityName), &subRspPayLoad)
	}
	common.MarshalSlice(subRspPayLoad, &respPayLoad)

	SendRespWithPayLoad(c, head, respPayLoad)
	infoLog.Printf("ProcRefreshLocation ret=%d \n", 0)
	return true
}

func DelSingleUserLocationInRedis(uid uint32) (ret bool) {
	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	key := GetUserPlaceKey(uid)
	//infoLog.Println("key=", key)
	n, err := redisConn.Do("Del", key) // get uid
	count, err := redis.Int(n, err)
	if err != nil {
		infoLog.Println("Reids failed Del [key err] =", key, err)
		ret = false
	} else {
		infoLog.Println("[key count]=", key, count)
		ret = true
	}
	return
}

func GetSingleUserLocationFromRedis(uid uint32, redisConn *redis.Conn) (strLocation string) {
	key := GetUserPlaceKey(uid)
	infoLog.Println("key=", key)
	n, err := (*redisConn).Do("Get", key) // get uid
	strLocation, err = redis.String(n, err)
	if err != nil {
		infoLog.Println("Reids failed Get [key err] =", key, err)
	} else {
		infoLog.Println("[key value]=", key, strLocation)
	}

	return
}

func GetSingleUserLocationFromDb(uid uint32, ts uint64, locationInfo *tagLocationInfo, conn *redis.Conn) (ret uint32) {
	infoLog.Println("GetSingleUserLocationFromDb uid =", uid)
	var strUid, strAllow, strLa, strLong, strCo string
	var strA1, strA2, strA3, strLoc, strSubLoc, strNe, strVer string
	err := db.QueryRow("select t1.USERID, t1.ALLOWED, t1.LATITUDE, t1.LONGITUDE, t1.COUNTRY,"+
		"t1.ADMINISTRATIVE1, t1.ADMINISTRATIVE2, t1.ADMINISTRATIVE3, t1.LOCALITY,"+
		"t1.SUBLOCALITY, t1.NEIGHBORHOOD, t2.LOCATIONVERSION from HT_USER_LOCATION2 "+
		"as t1 left join (HT_USER_PROPERTY as t2) on (t2.USERID = t1.USERID)"+
		" where t1.USERID = ? AND t2.LOCATIONVERSION > ?;", uid, ts).Scan(&strUid, &strAllow, &strLa, &strLong, &strCo, &strA1, &strA2, &strA3, &strLoc, &strSubLoc, &strNe, &strVer)
	switch {
	case err == sql.ErrNoRows:
		infoLog.Printf("GetSingleUserLocationFromDb not found uid=%v ts=%v\n", uid, ts)
		ret = uint32(DB_RET_NOT_EXIST)
		return
	case err != nil:
		infoLog.Println("GetSingleUserLocationFromDb [uid  ts err] =", uid, ts, err)
		ret = uint32(DB_RET_EXEC_FAILED)
		return
	default:
		infoLog.Printf("uid=%s allow=%s la=%s long=%s co=%s A1=%s A2=%s A3=%s Loc=%s SubLod=%s Ne=%s Ver=%s\n",
			strUid,
			strAllow,
			strLa,
			strLong,
			strCo,
			strA1,
			strA2,
			strA3,
			strLoc,
			strSubLoc,
			strNe,
			strVer)
		// add to listLocationInfo
		allow, err := strconv.Atoi(strAllow)
		if err != nil {
			infoLog.Println("GetSingleUserLocationFromDb strconv failed allow =", strAllow)
			ret = uint32(DB_RET_NOT_EXIST)
			return
		}

		ver, err := strconv.Atoi(strVer)
		if err != nil {
			infoLog.Println("GetSingleUserLocationFromDb strconv failed ver =", strVer)
			ret = uint32(DB_RET_NOT_EXIST)
			return
		}

		// set redis
		value := strUid + "#" + strAllow + "#" + strLa + "#" + strLong + "#" + strCo + "#"
		value += strA1 + "#" + strA2 + "#" + strA3 + "#" + strLoc + "#" + strSubLoc + "#"
		value += strNe + "#" + strVer

		strKey := GetUserPlaceKey(uid)
		_, err = (*conn).Do("SETEX", strKey, kTTLInterval, value)
		if err != nil {
			infoLog.Printf("SETEX key=%s value=%s faield err=%v", strKey, value, err)
		} else {
			infoLog.Printf("SETEX key=%s value=%s succ err=%v", strKey, value, err)
		}
		// return result
		ret = uint32(DB_RET_SUCCESS) // only set here
		locationInfo.Uid = uid
		locationInfo.Allowed = uint8(allow)
		locationInfo.Latitude = []byte(strLa)
		locationInfo.Longtitude = []byte(strLong)
		locationInfo.PlaceInfo.Admin1 = []byte(strA1)
		locationInfo.PlaceInfo.Admin2 = []byte(strA2)
		locationInfo.PlaceInfo.Admin3 = []byte(strA3)
		locationInfo.PlaceInfo.Country = []byte(strCo)
		locationInfo.PlaceInfo.Locality = []byte(strLoc)
		locationInfo.PlaceInfo.SubLocality = []byte(strSubLoc)
		locationInfo.PlaceInfo.Neighborhood = []byte(strNe)
		locationInfo.TimeStamp = uint64(ver)
	}
	return
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	infoLog.Println("OnClose:", c.GetExtraData())
}

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ServerConf string `short:"c" long:"conf" description:"Server Config" optional:"no"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
	}

	if options.ServerConf == "" {
		log.Fatalln("Must input config file name")
	}

	// log.Println("config name =", options.ServerConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.ServerConf)
	if err != nil {
		log.Printf("load config file=%s failed", options.ServerConf)
		return
	}
	// 配置文件只读 设置此标识提升性能
	cfg.BlockMode = false
	// 定义一个文件
	fileName := cfg.Section("LOG").Key("path").MustString("/home/ht/server.log")
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
		return
	}

	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	// 配置log的Flag参数
	infoLog.SetFlags(infoLog.Flags() | log.LstdFlags)

	// init redis pool
	redisIp := cfg.Section("REDIS").Key("redis_ip").MustString("127.0.0.1")
	redisPort := cfg.Section("REDIS").Key("redis_port").MustInt(6379)
	infoLog.Printf("redis ip=%v port=%v", redisIp, redisPort)
	pool = newPool(redisIp + ":" + strconv.Itoa(redisPort))

	// init mysql
	mysqlHost := cfg.Section("MYSQL").Key("mysql_host").MustString("127.0.0.1")
	mysqlUser := cfg.Section("MYSQL").Key("mysql_user").MustString("IMServer")
	mysqlPasswd := cfg.Section("MYSQL").Key("mysql_passwd").MustString("hello")
	mysqlDbName := cfg.Section("MYSQL").Key("mysql_db").MustString("HT_IMDB")
	mysqlPort := cfg.Section("MYSQL").Key("mysql_port").MustString("3306")

	infoLog.Printf("mysql host=%v user=%v passwd=%v dbname=%v port=%v",
		mysqlHost,
		mysqlUser,
		mysqlPasswd,
		mysqlDbName,
		mysqlPort)

	db, err = sql.Open("mysql", mysqlUser+":"+mysqlPasswd+"@"+"tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDbName+"?charset=utf8&timeout=90s")
	if err != nil {
		infoLog.Println("open mysql failed")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// creates a tcp listener
	serverIp := cfg.Section("LOCAL_SERVER").Key("bind_ip").MustString("127.0.0.1")
	serverPort := cfg.Section("LOCAL_SERVER").Key("bind_port").MustInt(8990)
	infoLog.Printf("serverIp=%v serverPort=%v", serverIp, serverPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	sendChanLimit := cfg.Section("CHANLIMIT").Key("max_send_chan_count").MustUint(1000)
	recvChanLimit := cfg.Section("CHANLIMIT").Key("max_recv_chan_count").MustUint(1000)
	config := &gotcp.Config{
		PacketSendChanLimit:    uint32(sendChanLimit),
		PacketReceiveChanLimit: uint32(recvChanLimit),
	}
	srv := gotcp.NewServer(config, &Callback{}, &common.XTHeadProtocol{})

	// starts service
	go srv.Start(listener, time.Second)
	infoLog.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	infoLog.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
	db.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
