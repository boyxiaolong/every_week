package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

const (
	msgTypeCommonReply          = 7102
	msgTypeDeliverAnnouncement  = 7103
	msgTypeSendMail             = 7104
	msgTypeGetPlayerInfo        = 7105
	msgTypeProhibitTalking      = 7106
	msgTypeBanAccount           = 7107
	msgTypeCloseServerGroup     = 7108
	msgTypeGetServerStatus      = 7109
	msgTypeGetServerStatusReply = 7110
	msgTypeRemoveTokenCache     = 7111
	msgTypeSendScrollingMsg     = 7112
	msgTypeGetPlayerInfoReply   = 7113
	msgTypeReloadConfig         = 7115
	msgTypeGetAccountCount      = 7116
	msgTypeGetAccountCountReply = 7117
	msgTypeReducePlayerResource = 7118
	msgTypeSendMailEx           = 7119
	msgTypePlayerRename         = 7120
	msgTypeGuildRename          = 7121
	msgTypeGuildChangeShortName = 7122
	msgTypeUpdateAllKingdomInfo = 7123
	msgTypeKickoffPlayers       = 7124
	msgTypeSetRoleRule          = 7125
	msgTypeSetKingdomInfo       = 7126
	msgGetAllKingdomInfo        = 7127
	msgGetAllKingdomInfoReply   = 7128
	msgGetCastleCount           = 7129
	msgGetCastleCountReply      = 7130
	msgSetKingdomNationMap      = 7131
)

func main() {
	host := flag.String("ip", "127.0.0.1", "gmserver ip")
	port := flag.String("port", "10180", "gmserver port")
	action := flag.String("action", "announce", `do action:
	announce[发布公告], mail[发送邮件], mailex[发送文字库邮件], playerinfo[获取玩家信息],
	banaccount[封号], removetoken[移除登录缓存], closeservers[关闭服务器组],
	serverstatus[获取服务器状态], scrollingmsg[走马灯], bantalking[禁言],
	reload[重载配置], accountcount[获取帐号数], reduceres[扣点], playerrename[玩家改名],
	guildrename[工会改名], guildrename2[修改工会短名字], updatekingdoms[更新所有王国信息],
	kickoffplayers[踢玩家下线], setrolerule[设置创号规则], setkingdom[设置单个王国信息],
	getkingdoms[获得所有王国信息], castlecount[获得指定王国城堡数], setnations[设置王国-国家代码映射]`)

	flag.Parse()

	// connect
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if handleWrite(conn, *action) > 0 {
		handleRead(conn)
	}
}

func handleWrite(conn net.Conn, action string) int {
	// read json file
	actionFile := "./" + action + ".json"
	data, err := ioutil.ReadFile(actionFile)
	if err != nil {
		fmt.Println("can't read ", actionFile, ", error: ", err)
		return -1
	}

	var msg []byte
	switch action {
	case "announce":
		msg = makeJSONMsg(msgTypeDeliverAnnouncement, data)
	case "mail":
		msg = makeJSONMsg(msgTypeSendMail, data)
	case "mailex":
		msg = makeJSONMsg(msgTypeSendMailEx, data)
	case "playerinfo":
		msg = makeJSONMsg(msgTypeGetPlayerInfo, data)
	case "banaccount":
		msg = makeJSONMsg(msgTypeBanAccount, data)
	case "removetoken":
		msg = makeJSONMsg(msgTypeRemoveTokenCache, data)
	case "closeservers":
		msg = makeJSONMsg(msgTypeCloseServerGroup, data)
	case "serverstatus":
		msg = makeJSONMsg(msgTypeGetServerStatus, data)
	case "scrollingmsg":
		msg = makeJSONMsg(msgTypeSendScrollingMsg, data)
	case "bantalking":
		msg = makeJSONMsg(msgTypeProhibitTalking, data)
	case "reload":
		msg = makeJSONMsg(msgTypeReloadConfig, data)
	case "accountcount":
		msg = makeJSONMsg(msgTypeGetAccountCount, data)
	case "reduceres":
		msg = makeJSONMsg(msgTypeReducePlayerResource, data)
	case "playerrename":
		msg = makeJSONMsg(msgTypePlayerRename, data)
	case "guildrename":
		msg = makeJSONMsg(msgTypeGuildRename, data)
	case "guildrename2":
		msg = makeJSONMsg(msgTypeGuildChangeShortName, data)
	case "updatekingdoms":
		msg = makeJSONMsg(msgTypeUpdateAllKingdomInfo, data)
	case "kickoffplayers":
		msg = makeJSONMsg(msgTypeKickoffPlayers, data)
	case "setrolerule":
		msg = makeJSONMsg(msgTypeSetRoleRule, data)
	case "setkingdom":
		msg = makeJSONMsg(msgTypeSetKingdomInfo, data)
	case "getkingdoms":
		msg = makeJSONMsg(msgGetAllKingdomInfo, data)
	case "castlecount":
		msg = makeJSONMsg(msgGetCastleCount, data)
	case "setnations":
		msg = makeJSONMsg(msgSetKingdomNationMap, data)
	default:
		fmt.Println("Unimplement action: ", action)
		return -1
	}

	fmt.Println(string(msg[4:]))
	sendLen, e := conn.Write(msg)
	if e != nil {
		fmt.Println("send message error: ", e.Error())
		return -1
	}

	return sendLen
}

func handleRead(conn net.Conn) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}

	if len < 4 {
		fmt.Println("Recevie less bytes. ", len)
		return
	}

	msgSize := binary.LittleEndian.Uint16(buf)
	msgType := binary.LittleEndian.Uint16(buf[2:])

	if msgSize != uint16(len) {
		fmt.Println("Message size is wrong(", msgSize, msgType, ")")
		return
	}

	fmt.Println(msgSize, msgType, string(buf[4:len]))

	switch msgType {
	case msgTypeCommonReply:
		var reply commonReply
		err := json.Unmarshal(buf[4:len], &reply)
		if err != nil {
			fmt.Println("Unmarshal reply fail: ", err)
			break
		}

		fmt.Printf("reply: %+v\n", reply)

	case msgTypeGetPlayerInfoReply:
		var reply playerInfoReply
		err := json.Unmarshal(buf[4:len], &reply)
		if err != nil {
			fmt.Println("Unmarshal reply fail: ", err)
			break
		}

		fmt.Printf("reply: %+v\n", reply)

	case msgTypeGetServerStatusReply:
		var reply serverStatusArray
		err := json.Unmarshal(buf[4:len], &reply)
		if err != nil {
			fmt.Println("Unmarshal reply fail: ", err)
			break
		}

		fmt.Printf("reply: %+v\n", reply)

	case msgTypeGetAccountCountReply:
		var reply accountCountReply
		err := json.Unmarshal(buf[4:len], &reply)
		if err != nil {
			fmt.Println("Unmarshal reply fail: ", err)
			break
		}

		fmt.Printf("reply: %+v\n", reply)

	case msgGetAllKingdomInfoReply:
		var reply getAllKingdomReply
		err := json.Unmarshal(buf[4:len], &reply)
		if err != nil {
			fmt.Println("Unmarshal reply fail: ", err)
			break
		}

		fmt.Printf("reply: %+v\n", reply)

	case msgGetCastleCountReply:
		var reply getCastleCountReply
		err := json.Unmarshal(buf[4:len], &reply)
		if err != nil {
			fmt.Println("Unmarshal reply fail: ", err)
			break
		}

		fmt.Printf("reply: %+v\n", reply)

	default:
		fmt.Println("Unknown reply message. ", msgType)
	}
}

func makePacket(size uint16, typ uint16) []byte {
	p := make([]byte, 4, size)
	binary.LittleEndian.PutUint16(p, size)
	binary.LittleEndian.PutUint16(p[2:], typ)
	return p
}

type commonReply struct {
	Error uint `json:"error"`
}

type playerInfo struct {
	Iggid            uint64 `json:"iggid"`
	KingdomID        uint32 `json:"kingdom_id"`
	DBID             uint32 `json:"db_id"`
	Name             string `json:"name"`
	CastleLevel      uint32 `json:"castle_level"`
	VipLevel         uint32 `json:"vip_level"`
	Gem              uint64 `json:"gem"`
	Oil              uint64 `json:"oil"`
	Food             uint64 `json:"Food"`
	Wood             uint64 `json:"Wood"`
	Steel            uint64 `json:"steel"`
	LeagueID         uint64 `json:"league_id"`
	LeagueName       string `json:"league_name"`
	LeagueShortName  string `json:"league_short_name"`
	RegTime          int64  `json:"reg_time"`
	EnterKingdomTime int64  `json:"enter_kingdom_time"`
	OnlineTime       int64  `json:"online_time"`
	Ban              int    `json:"ban"`
}

type playerInfoReply struct {
	Error      uint       `json:"error"`
	PlayerInfo playerInfo `json:"player_info"`
}

type serverStatus struct {
	KingdomID uint   `json:"kingdom_id"`
	ServerID  uint   `json:"server_id"`
	IP        string `json:"ip"`
	Port      uint   `json:"port"`
	Status    int    `json:"status"`
	Online    int    `json:"online"`
}

type serverStatusArray struct {
	ServerStatus []serverStatus `json:"server_status"`
}

type accountCountReply struct {
	Error uint   `json:"error"`
	Count uint64 `json:"count"`
}

type itemST struct {
	ID    string `json:"id"`
	Value uint   `json:"value"`
}

type sysAnnouncement struct {
	Language uint     `json:"language"`
	Subject  string   `json:"subject"`
	Content  string   `json:"content"`
	Items    []itemST `json:"items"`
	Players  []uint64 `json:"players"`
}

type kingdomInfo struct {
	KingdomID        uint  `json:"kingdom_id"`
	KingdomBeginTime int64 `json:"kingdom_begin_time"`
	NewFlag          int   `json:"new"`
}

type getAllKingdomReply struct {
	Error    uint          `json:"error"`
	Kingdoms []kingdomInfo `json:"kingdoms"`
}

type getCastleCountReply struct {
	Error uint `json:"error"`
	Count uint `json:"count"`
}

func makeAnnouncementMsg() []byte {
	announcement := sysAnnouncement{
		Language: 0,
		Subject:  "system announcement",
		Content:  "Hello world!",
		Items: []itemST{
			{ID: "1#0", Value: 100},
			{ID: "2#0", Value: 30},
		},
		Players: []uint64{123456, 99316},
	}

	b, _ := json.Marshal(announcement)
	fmt.Println(string(b))
	return makeJSONMsg(msgTypeDeliverAnnouncement, b)
}

func makeJSONMsg(msgType int, jsonData []byte) []byte {
	p := makePacket(uint16(4+len(jsonData)), uint16(msgType))
	p = append(p, jsonData...)
	return p
}
