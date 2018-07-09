package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "strings"
	// "time"
)

var (
	ErrNilDbObject    = errors.New("not set  object current is nil")
	ErrDbParam        = errors.New("err  param error")
	ErrNilNoSqlObject = errors.New("not set nosql object current is nil")
	ErrMemberQuit     = errors.New("member quit state is not 0")
)

var (
	RoomInfoDB        = "roominfo"
	RoomMemberTable   = "room_member"
	RoomMaxOrderTable = "room_max_order"
	pageSize          = 20
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ServerConf string `short:"c" long:"conf" description:"Server Config" optional:"no"`
}

var options Options
var (
	infoLog *log.Logger
	db      *sql.DB
	//ssdb        *gossdb.Connectors
	mongoSess *mgo.Session
)

type MemberInfoStore struct {
	Id          string `bson:"_id"` // 在Mongo中唯一标识一条记录
	RoomId      uint32
	Uid         uint32
	InviteUid   uint32
	NickName    string
	OrderId     uint32
	JoinTS      uint32
	PushSetting uint32
	ContactList uint32
	QuitStat    uint32
	UpdateTS    uint32
}

type MaxOrderInfo struct {
	RoomId   uint32 // 群聊Id
	MaxOrder uint32 // 当前最大加入序号
}

var parser = flags.NewParser(&options, flags.Default)

func main() {
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
	fileName := cfg.Section("LOG").Key("path").MustString("/home/ht/muc_data_migration.log")
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

	// 读取配置的最大RoomId
	beginRoomId := cfg.Section("MAXROOMID").Key("beginRoomId").MustInt(1)
	maxRoomId := cfg.Section("MAXROOMID").Key("roomId").MustInt(1)

	// 读取mysql 配置
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
	// 读取mongodb 配置
	// 创建mongodb对象
	mongo_url := cfg.Section("MONGO").Key("url").MustString("localhost")
	infoLog.Printf("Mongo url=%s", mongo_url)
	mongoSess, err = mgo.Dial(mongo_url)
	if err != nil {
		log.Fatalln("connect mongodb failed")
		return
	}
	defer mongoSess.Close()
	// Optional. Switch the session to a monotonic behavior.
	mongoSess.SetMode(mgo.Monotonic, true)
	// 读取mysql中的群成员信息 写入mongod中
	for startRoomId := beginRoomId; startRoomId < maxRoomId; startRoomId++ {
		MigrationRoomMemberFromMySqlToMongo(uint32(startRoomId))
	}
	return

}

func MigrationRoomMemberFromMySqlToMongo(roomId uint32) (err error) {
	if roomId == 0 {
		err = ErrDbParam
		return
	}

	// 查询整个群所有群成员
	rows, err := db.Query("select ROOMID, MEMBERID, INVITEID, NAMECARD, ORDERID, UNIX_TIMESTAMP(JOINTIME), PUSHSETTING, CONTACTLIST, QUITSTATE, UNIX_TIMESTAMP(UPDATETIME) from HT_MUC_MEMBER where ROOMID = ? order by ORDERID asc", roomId)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var roomId, memberId uint32
		var storeInviteId, storeOrderId, storeJoinTS, storePushSetting, storeContactList, storeQuitStat, storeUpdateTS sql.NullInt64
		var storeNameCard sql.NullString
		if err := rows.Scan(&roomId, &memberId, &storeInviteId, &storeNameCard, &storeOrderId, &storeJoinTS, &storePushSetting, &storeContactList, &storeQuitStat, &storeUpdateTS); err != nil {
			infoLog.Printf("MigrationRoomMemberFromMySqlToMongo rows.Scan failed")
			continue
		}

		var inviteId, orderId, joinTS, pushSetting, contactList, quitStat, updateTS uint32
		var nameCard string
		if storeInviteId.Valid {
			inviteId = uint32(storeInviteId.Int64)
		}

		if storeNameCard.Valid {
			nameCard = storeNameCard.String
		}

		if storeOrderId.Valid {
			orderId = uint32(storeOrderId.Int64)
		}

		if storeJoinTS.Valid {
			joinTS = uint32(storeJoinTS.Int64)
		}

		if storePushSetting.Valid {
			pushSetting = uint32(storePushSetting.Int64)
		}

		if storeContactList.Valid {
			contactList = uint32(storeContactList.Int64)
		}

		if storeQuitStat.Valid {
			quitStat = uint32(storeQuitStat.Int64)
		}

		if storeUpdateTS.Valid {
			updateTS = uint32(storeUpdateTS.Int64)
		}
		infoLog.Printf("roomId=%v memberId=%v inviteId=%v nameCare=%v orderId=%v joinTS=%v pushSettint=%v contactList=%v quitStat=%v updateTS=%v",
			roomId,
			memberId,
			inviteId,
			nameCard,
			orderId,
			joinTS,
			pushSetting,
			contactList,
			quitStat,
			updateTS)

		// 将群成员写入Mongo中 首先获取集合对象
		mongoConn := mongoSess.DB(RoomInfoDB).C(RoomMemberTable)
		memberInfo := &MemberInfoStore{
			Id:          bson.NewObjectId().String(),
			RoomId:      roomId,
			Uid:         memberId,
			InviteUid:   inviteId,
			NickName:    nameCard,
			OrderId:     orderId,
			JoinTS:      joinTS,
			PushSetting: pushSetting,
			ContactList: contactList,
			QuitStat:    quitStat,
			UpdateTS:    updateTS,
		}
		err = mongoConn.Insert(memberInfo)
		// 如果RoomId+Uid 作为唯一所以 如果RoomId+Uid已经存在则插入失败
		if err != nil {
			infoLog.Printf("MigrationRoomMemberFromMySqlToMongo Mongo Insert failed, roomId=%v memberId=%v inviteId=%v err=%v", roomId, memberId, inviteId, err)
			continue
		}
		// 更新Mongo中的MaxOrder
		maxOrder := &MaxOrderInfo{RoomId: roomId, MaxOrder: orderId}
		maxOrderConn := mongoSess.DB(RoomInfoDB).C(RoomMaxOrderTable)
		_, err = maxOrderConn.Upsert(bson.M{"roomid": maxOrder.RoomId}, maxOrder)
		if err != nil {
			infoLog.Printf("MigrationRoomMemberFromMySqlToMongo Mongo upsert faield roomId=%v maxOrder=%v err=%v", roomId, orderId, err)
			continue
		}
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
